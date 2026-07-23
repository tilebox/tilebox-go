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
		{name: "protobuf enum", expression: Field("value").Equal(datasetsv1.ProcessingLevel_PROCESSING_LEVEL_L1), wantErr: "protobuf enums are not queryable"},
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
