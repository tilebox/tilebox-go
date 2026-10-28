package query

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	datasetsv1 "github.com/tilebox/tilebox-go/protogen/datasets/v1"
	"google.golang.org/protobuf/proto"
)

func TestComparisonExpressionValues(t *testing.T) {
	type signedCount int32
	type granuleName string

	tests := []struct {
		name  string
		value any
		want  *datasetsv1.FieldQueryValue
	}{
		{name: "bool", value: false, want: datasetsv1.FieldQueryValue_builder{BoolValue: new(false)}.Build()},
		{name: "string", value: "S2A_GRANULE", want: datasetsv1.FieldQueryValue_builder{StringValue: new("S2A_GRANULE")}.Build()},
		{name: "named string", value: granuleName("S2B_GRANULE"), want: datasetsv1.FieldQueryValue_builder{StringValue: new("S2B_GRANULE")}.Build()},
		{name: "int", value: 0, want: datasetsv1.FieldQueryValue_builder{Int64Value: new(int64(0))}.Build()},
		{name: "named int32", value: signedCount(12), want: datasetsv1.FieldQueryValue_builder{Int64Value: new(int64(12))}.Build()},
		{name: "uint64", value: uint64(42), want: datasetsv1.FieldQueryValue_builder{Uint64Value: new(uint64(42))}.Build()},
		{name: "float32", value: float32(1.5), want: datasetsv1.FieldQueryValue_builder{DoubleValue: new(float64(1.5))}.Build()},
		{name: "float64", value: 20.0, want: datasetsv1.FieldQueryValue_builder{DoubleValue: new(float64(20))}.Build()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ToProtoExpression(Field("value").Equal(tt.value))
			require.NoError(t, err)
			assert.True(t, proto.Equal(tt.want, got.GetComparison().GetValue()))
		})
	}
}

func TestExpressionToProto(t *testing.T) {
	expression := And(
		Field("cloud_cover").LessThan(20.0),
		Or(
			Field("quality").Equal(int32(2)),
			Field("quality").Equal(int64(3)),
		),
		Field("published").IsNotNull(),
	)

	got, err := ToProtoExpression(expression)
	require.NoError(t, err)

	assert.Equal(t, datasetsv1.LogicalOperator_LOGICAL_OPERATOR_AND, got.GetLogical().GetOperator())
	require.Len(t, got.GetLogical().GetOperands(), 3)
	assert.Equal(t, "cloud_cover", got.GetLogical().GetOperands()[0].GetComparison().GetFieldName())
	assert.Equal(t, datasetsv1.FieldComparisonOperator_FIELD_COMPARISON_OPERATOR_LESS_THAN, got.GetLogical().GetOperands()[0].GetComparison().GetOperator())
	assert.Equal(t, datasetsv1.LogicalOperator_LOGICAL_OPERATOR_OR, got.GetLogical().GetOperands()[1].GetLogical().GetOperator())
	assert.Equal(t, datasetsv1.LogicalOperator_LOGICAL_OPERATOR_NOT, got.GetLogical().GetOperands()[2].GetLogical().GetOperator())
	assert.Equal(t, "published", got.GetLogical().GetOperands()[2].GetLogical().GetOperands()[0].GetIsNull().GetFieldName())
}

func TestInvalidExpression(t *testing.T) {
	tests := []struct {
		name       string
		expression Expression
		wantErr    string
	}{
		{name: "nil value", expression: Field("value").Equal(nil), wantErr: "comparison value cannot be nil"},
		{name: "unsupported value", expression: Field("value").Equal([]byte("text")), wantErr: "unsupported comparison value type []uint8"},
		{name: "NaN", expression: Field("value").Equal(math.NaN()), wantErr: "comparison value must be finite"},
		{name: "infinity", expression: Field("value").Equal(math.Inf(1)), wantErr: "comparison value must be finite"},
		{name: "protobuf enum", expression: Field("value").Equal(datasetsv1.ProcessingLevel_PROCESSING_LEVEL_L1), wantErr: "use a public enum string value"},
		{name: "boolean ordering", expression: Field("value").LessThan(true), wantErr: "boolean values only support Equal and NotEqual"},
		{name: "nil expression", expression: Not(nil), wantErr: "query expression cannot be nil"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ToProtoExpression(tt.expression)
			require.ErrorContains(t, err, tt.wantErr)
		})
	}
}

func TestArrayExpressions(t *testing.T) {
	for _, tc := range []struct {
		name     string
		build    func(...string) Expression
		operator datasetsv1.FieldArrayOperator
	}{
		{"contains all", Field("use_cases").ContainsAll, datasetsv1.FieldArrayOperator_FIELD_ARRAY_OPERATOR_CONTAINS_ALL},
		{"contains any", Field("use_cases").ContainsAny, datasetsv1.FieldArrayOperator_FIELD_ARRAY_OPERATOR_OVERLAPS},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for _, count := range []int{0, 1, 100, 101} {
				values := make([]string, count)
				for i := range values {
					values[i] = "wind-power"
				}
				expression := tc.build(values...)
				if count > 0 {
					values[0] = "changed"
				}
				got, err := ToProtoExpression(expression)
				if count > 100 {
					require.ErrorContains(t, err, "at most 100 values")
					continue
				}
				require.NoError(t, err)
				assert.False(t, got.HasComparison() || got.HasLogical() || got.HasIsNull())
				predicate := got.GetArrayPredicate()
				require.NotNil(t, predicate)
				assert.Equal(t, "use_cases", predicate.GetFieldName())
				assert.Equal(t, tc.operator, predicate.GetOperator())
				require.Len(t, predicate.GetValues(), count)
				for _, value := range predicate.GetValues() {
					assert.Equal(t, "wind-power", value)
				}
				wire, err := proto.Marshal(got)
				require.NoError(t, err)
				decoded := &datasetsv1.FilterExpression{}
				require.NoError(t, proto.Unmarshal(wire, decoded))
				assert.True(t, proto.Equal(got, decoded))
			}
			_, err := ToProtoExpression(tc.build("defense", ""))
			require.ErrorContains(t, err, "must not be empty strings")
		})
	}

	got, err := ToProtoExpression(And(Field("level").Equal("L1"), Not(Or(
		Field("use_cases").ContainsAll("defense", "agriculture"),
		Field("use_cases").ContainsAny("wind-power", "雪"), //nolint:gosmopolitan // Public enum values may contain Unicode.
	))))
	require.NoError(t, err)
	operands := got.GetLogical().GetOperands()[1].GetLogical().GetOperands()[0].GetLogical().GetOperands()
	assert.Equal(t, []string{"defense", "agriculture"}, operands[0].GetArrayPredicate().GetValues())
	assert.Equal(t, []string{"wind-power", "雪"}, operands[1].GetArrayPredicate().GetValues()) //nolint:gosmopolitan // Verify the public value is unchanged.
}
