package query

import (
	"errors"
	"fmt"
	"math"
	"reflect"

	datasetsv1 "github.com/tilebox/tilebox-go/protogen/datasets/v1"
	"google.golang.org/protobuf/reflect/protoreflect"
)

// Expression is a filter expression over queryable dataset fields.
//
// Expressions use SQL- and CQL2-compatible three-valued Boolean logic. A
// comparison with a missing field evaluates to unknown, and a datapoint only
// matches when the complete expression evaluates to true.
type Expression interface {
	toProto() (*datasetsv1.FilterExpression, error)
}

// Field identifies a top-level queryable dataset field by name.
type Field string

// Equal compares the field with value for equality.
func (f Field) Equal(value any) Expression {
	return comparisonExpression{fieldName: string(f), operator: datasetsv1.FieldComparisonOperator_FIELD_COMPARISON_OPERATOR_EQUAL, value: value}
}

// NotEqual compares the field with value for inequality. Missing fields do not
// match; combine this expression with IsNull using Or to include them.
func (f Field) NotEqual(value any) Expression {
	return comparisonExpression{fieldName: string(f), operator: datasetsv1.FieldComparisonOperator_FIELD_COMPARISON_OPERATOR_NOT_EQUAL, value: value}
}

// LessThan matches fields less than value.
func (f Field) LessThan(value any) Expression {
	return comparisonExpression{fieldName: string(f), operator: datasetsv1.FieldComparisonOperator_FIELD_COMPARISON_OPERATOR_LESS_THAN, value: value}
}

// LessThanOrEqual matches fields less than or equal to value.
func (f Field) LessThanOrEqual(value any) Expression {
	return comparisonExpression{fieldName: string(f), operator: datasetsv1.FieldComparisonOperator_FIELD_COMPARISON_OPERATOR_LESS_THAN_OR_EQUAL, value: value}
}

// GreaterThan matches fields greater than value.
func (f Field) GreaterThan(value any) Expression {
	return comparisonExpression{fieldName: string(f), operator: datasetsv1.FieldComparisonOperator_FIELD_COMPARISON_OPERATOR_GREATER_THAN, value: value}
}

// GreaterThanOrEqual matches fields greater than or equal to value.
func (f Field) GreaterThanOrEqual(value any) Expression {
	return comparisonExpression{fieldName: string(f), operator: datasetsv1.FieldComparisonOperator_FIELD_COMPARISON_OPERATOR_GREATER_THAN_OR_EQUAL, value: value}
}

// IsNull matches datapoints for which the field is missing or explicitly null.
func (f Field) IsNull() Expression {
	return nullExpression{fieldName: string(f)}
}

// IsNotNull matches datapoints for which the field is present and non-null.
func (f Field) IsNotNull() Expression {
	return Not(f.IsNull())
}

// And combines two or more expressions using logical AND.
func And(first, second Expression, rest ...Expression) Expression {
	operands := make([]Expression, 0, 2+len(rest))
	operands = append(operands, first, second)
	operands = append(operands, rest...)
	return logicalExpression{operator: datasetsv1.LogicalOperator_LOGICAL_OPERATOR_AND, operands: operands}
}

// Or combines two or more expressions using logical OR.
func Or(first, second Expression, rest ...Expression) Expression {
	operands := make([]Expression, 0, 2+len(rest))
	operands = append(operands, first, second)
	operands = append(operands, rest...)
	return logicalExpression{operator: datasetsv1.LogicalOperator_LOGICAL_OPERATOR_OR, operands: operands}
}

// Not negates an expression. Unknown remains unknown under negation.
func Not(expression Expression) Expression {
	return logicalExpression{operator: datasetsv1.LogicalOperator_LOGICAL_OPERATOR_NOT, operands: []Expression{expression}}
}

// ToProtoExpression converts an expression to its protobuf representation.
func ToProtoExpression(expression Expression) (*datasetsv1.FilterExpression, error) {
	if expression == nil {
		return nil, errors.New("query expression cannot be nil")
	}
	return expression.toProto()
}

type comparisonExpression struct {
	fieldName string
	operator  datasetsv1.FieldComparisonOperator
	value     any
}

func (expression comparisonExpression) toProto() (*datasetsv1.FilterExpression, error) {
	value, valueKind, err := queryValueToProto(expression.value)
	if err != nil {
		return nil, fmt.Errorf("field %q: %w", expression.fieldName, err)
	}
	if valueKind == reflect.Bool && expression.operator != datasetsv1.FieldComparisonOperator_FIELD_COMPARISON_OPERATOR_EQUAL && expression.operator != datasetsv1.FieldComparisonOperator_FIELD_COMPARISON_OPERATOR_NOT_EQUAL {
		return nil, fmt.Errorf("field %q: boolean values only support Equal and NotEqual", expression.fieldName)
	}
	return datasetsv1.FilterExpression_builder{
		Comparison: datasetsv1.FieldComparison_builder{
			FieldName: expression.fieldName,
			Operator:  expression.operator,
			Value:     value,
		}.Build(),
	}.Build(), nil
}

type nullExpression struct {
	fieldName string
}

func (expression nullExpression) toProto() (*datasetsv1.FilterExpression, error) {
	return datasetsv1.FilterExpression_builder{
		IsNull: datasetsv1.FieldNullCheck_builder{FieldName: expression.fieldName}.Build(),
	}.Build(), nil
}

type logicalExpression struct {
	operator datasetsv1.LogicalOperator
	operands []Expression
}

func (expression logicalExpression) toProto() (*datasetsv1.FilterExpression, error) {
	operands := make([]*datasetsv1.FilterExpression, len(expression.operands))
	for i, operand := range expression.operands {
		converted, err := ToProtoExpression(operand)
		if err != nil {
			return nil, fmt.Errorf("logical operand %d: %w", i, err)
		}
		operands[i] = converted
	}
	return datasetsv1.FilterExpression_builder{
		Logical: datasetsv1.LogicalExpression_builder{
			Operator: expression.operator,
			Operands: operands,
		}.Build(),
	}.Build(), nil
}

func queryValueToProto(value any) (*datasetsv1.FieldQueryValue, reflect.Kind, error) {
	if value == nil {
		return nil, reflect.Invalid, errors.New("comparison value cannot be nil")
	}
	if _, ok := value.(protoreflect.Enum); ok {
		return nil, reflect.Invalid, fmt.Errorf("unsupported comparison value type %T; protobuf enums are not queryable", value)
	}

	reflected := reflect.ValueOf(value)
	switch reflected.Kind() { //nolint:exhaustive // Unsupported kinds are rejected by the default case.
	case reflect.Bool:
		converted := reflected.Bool()
		return datasetsv1.FieldQueryValue_builder{BoolValue: &converted}.Build(), reflect.Bool, nil
	case reflect.String:
		converted := reflected.String()
		return datasetsv1.FieldQueryValue_builder{StringValue: &converted}.Build(), reflect.String, nil
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		converted := reflected.Int()
		return datasetsv1.FieldQueryValue_builder{Int64Value: &converted}.Build(), reflect.Int64, nil
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		converted := reflected.Uint()
		return datasetsv1.FieldQueryValue_builder{Uint64Value: &converted}.Build(), reflect.Uint64, nil
	case reflect.Float32, reflect.Float64:
		converted := reflected.Float()
		if math.IsNaN(converted) || math.IsInf(converted, 0) {
			return nil, reflect.Invalid, errors.New("comparison value must be finite")
		}
		return datasetsv1.FieldQueryValue_builder{DoubleValue: &converted}.Build(), reflect.Float64, nil
	default:
		return nil, reflect.Invalid, fmt.Errorf("unsupported comparison value type %T; expected a boolean, string, or numeric value", value)
	}
}
