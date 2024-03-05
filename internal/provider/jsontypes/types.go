package jsontypes

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

var Int64Type = int64Type{
	Int64Type: types.Int64Type,
}

type int64Type struct {
	basetypes.Int64Type
}

func (i64 int64Type) Equal(o attr.Type) bool {
	_, ok := o.(int64Type)
	return ok
}

func (i64 int64Type) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	val, err := i64.Int64Type.ValueFromTerraform(ctx, in)

	return Int64{Int64Value: val.(types.Int64)}, err
}

var Float64Type = float64Type{
	Float64Type: types.Float64Type,
}

type float64Type struct {
	basetypes.Float64Type
}

func (f64 float64Type) Equal(o attr.Type) bool {
	_, ok := o.(float64Type)
	return ok
}

func (f64 float64Type) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	val, err := f64.Float64Type.ValueFromTerraform(ctx, in)

	return Float64{Float64Value: val.(types.Float64)}, err
}

var BoolType = boolType{
	BoolType: types.BoolType,
}

type boolType struct {
	basetypes.BoolType
}

func (mst boolType) Equal(o attr.Type) bool {
	_, ok := o.(boolType)
	return ok
}

func (mst boolType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	val, err := mst.BoolType.ValueFromTerraform(ctx, in)

	return Bool{BoolValue: val.(types.Bool)}, err
}

var StringType = stringType{
	types.StringType,
}

// stringType -
type stringType struct {
	basetypes.StringType
}

func (mst stringType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	val, err := mst.StringType.ValueFromTerraform(ctx, in)

	return String{StringValue: val.(types.String)}, err
}

func (mst stringType) Equal(o attr.Type) bool {
	_, ok := o.(stringType)
	return ok
}

func SetType[T JsonValue]() setType[T] {
	var v T
	return setType[T]{
		SetType: types.SetType{
			ElemType: v.Type(context.Background()),
		},
	}
}

type setType[T JsonValue] struct {
	types.SetType
}

func (st setType[T]) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	val, err := st.SetType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}

	set := val.(types.Set)

	var item T
	if it, et := item.Type(ctx), set.ElementType(ctx); !it.Equal(et) {
		return nil, fmt.Errorf("expected type: %T received %T", it, et)
	}

	return Set[T]{val.(types.Set)}, nil
}

func (st setType[T]) Equal(o attr.Type) bool {
	var base attr.Type

	switch typ := o.(type) {
	case setType[T]:
		base = typ.SetType
	case types.SetType:
		base = typ
	default:
		return false
	}

	return st.SetType.Equal(base)
}

// DynamicType represents a type that can dynamically adapt to the structure of the provided data.
type DynamicType struct {
	ElemType attr.Type
	Value    interface{}
}

// TerraformType returns the Terraform type of the value, which is dynamic in this case.
func (t DynamicType) TerraformType(ctx context.Context) tftypes.Type {
	return tftypes.DynamicPseudoType
}

// ElementType returns the type's element type.
func (t DynamicType) ElementType() attr.Type {
	if t.ElemType == nil {
		return nil
	}

	return t.ElemType
}

// ValueType returns an example attr.Value type that DynamicType might represent.
func (t DynamicType) ValueType(ctx context.Context) attr.Value {
	// TODO: Returning a generic placeholder types.String
	return types.String{}
}

// ValueFromTerraform returns an attr.Value that represents the given raw value.
func (t DynamicType) ValueFromTerraform(ctx context.Context, value tftypes.Value) (attr.Value, error) {
	if !value.IsKnown() {
		t.Value = nil
		return DynamicValue{Value: tftypes.UnknownValue}, nil
	}
	if value.IsNull() {
		t.Value = nil
		return DynamicValue{Value: nil}, nil
	}

	var err error
	switch {
	case value.Type().Equal(tftypes.String):
		var rawValue string
		err = value.As(&rawValue)
		t.Value = rawValue
	case value.Type().Equal(tftypes.Number):
		var rawValue float64
		err = value.As(&rawValue)
		t.Value = rawValue
	case value.Type().Equal(tftypes.Bool):
		var rawValue bool
		err = value.As(&rawValue)
		t.Value = rawValue
	default:
		return DynamicValue{}, fmt.Errorf("unsupported tftypes.Value type: %s", value.Type().String())
	}

	if err != nil {
		return DynamicValue{}, fmt.Errorf("failed to convert tftypes.Value: %w", err)
	}

	return DynamicValue{Value: t.Value}, nil
}

// FromTerraform5Value dynamically interprets and assigns the tftypes.Value to the DynamicType's Value field.
// This method is a custom implementation to handle dynamic or unknown Terraform types.
func (t *DynamicType) FromTerraform5Value(ctx context.Context, value tftypes.Value) error {
	if !value.IsKnown() {
		t.Value = tftypes.UnknownValue
		return nil
	}
	if value.IsNull() {
		t.Value = nil
		return nil
	}

	// Handling specific types by comparing tftypes.Type
	if value.Type().Equal(tftypes.String) {
		var str string
		if err := value.As(&str); err != nil {
			return fmt.Errorf("failed to convert tftypes.Value to string: %w", err)
		}
		t.Value = str
	} else {
		// For non-string or more complex types, handle them here.
		// This might involve more complex logic, especially for composite types.
		var rawValue interface{}
		if err := value.As(&rawValue); err != nil {
			return fmt.Errorf("failed to dynamically convert tftypes.Value: %w", err)
		}
		t.Value = rawValue
	}

	return nil
}

// ApplyTerraform5AttributePathStep handles attribute path stepping for the type.
func (t DynamicType) ApplyTerraform5AttributePathStep(step tftypes.AttributePathStep) (interface{}, error) {
	// This implementation does not support attribute path stepping.
	if _, ok := step.(tftypes.ElementKeyString); !ok {
		return nil, fmt.Errorf("cannot apply step %T to MapType", step)
	}

	return t.ElementType(), nil
}

// Equal checks if the type is equal to another type.
func (t DynamicType) Equal(other attr.Type) bool {
	_, ok := other.(DynamicType)
	return ok
}

// String provides a string representation of the type.
func (t DynamicType) String() string {
	return "DynamicType"
}

type TypeWithAttributeTypes interface {
	AttributeTypes() map[string]attr.Type
}

type MapType struct {
	basetypes.MapType
	ElemType attr.Type
}

func (m MapType) AttributeTypes() map[string]attr.Type {
	return map[string]attr.Type{}
}

// NewMapType constructs a new mapType with the provided element type.
func NewMapType(elemType attr.Type) MapType {
	return MapType{
		ElemType: elemType,
	}
}

func (m MapType) Equal(o attr.Type) bool {
	other, ok := o.(MapType)
	if !ok {
		return false
	}
	return m.ElemType.Equal(other.ElemType)
}

// NewMapValueMust creates a Map with a known value, converting any diagnostics
// into a panic at runtime. Access the value via the Map
// type Elements or ElementsAs methods.
func NewMapValueMust(elementType attr.Type, elements map[string]attr.Value) MapValue {
	m := NewMapValue(elements)

	return m
}

// Type returns a MapType with the same element type as `m`.
func (m MapValue) Type(ctx context.Context) attr.Type {
	return MapType{ElemType: m.ElementType(ctx)}
}

func (m MapType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	if in.Type() == nil {
		return NewMapNull(), nil
	}
	if !in.Type().Is(tftypes.Map{}) {
		return nil, fmt.Errorf("can't use %s as value of MapValue, can only use tftypes.Map values", in.String())
	}

	if !in.Type().Equal(tftypes.Map{ElementType: m.ElementType().TerraformType(ctx)}) {
		return nil, fmt.Errorf("can't use %s as value of Map with ElementType %T, can only use %s values", in.String(), m.ElementType(), m.ElementType().TerraformType(ctx).String())
	}

	// Use the recursive function to handle nested maps.
	return recursiveMapValueFromTerraformValue(ctx, in, m.ElemType)
}

// recursiveMapValueFromTerraformValue converts a tftypes.Value representing a map into an attr.Value,
// handling nested maps recursively.
func recursiveMapValueFromTerraformValue(ctx context.Context, in tftypes.Value, elemType attr.Type) (attr.Value, error) {

	if !in.IsKnown() {
		return NewMapUnknown(), nil
	}
	if in.IsNull() {
		return NewMapNull(), nil
	}

	var rawMap map[string]tftypes.Value
	err := in.As(&rawMap)
	if err != nil {
		return nil, fmt.Errorf("failed to convert tftypes.Value to map[string]tftypes.Value: %w", err)
	}

	elems := make(map[string]attr.Value)
	for key, val := range rawMap {
		// Check if the value is a map, indicating a need for recursion.
		if val.Type().Is(tftypes.Map{}) {
			nestedMapVal, nestedErr := recursiveMapValueFromTerraformValue(ctx, val, DynamicType{})
			if nestedErr != nil {
				return nil, fmt.Errorf("error processing nested map for key %s: %w", key, nestedErr)
			}
			elems[key] = nestedMapVal
		} else {
			// For non-map values, convert directly.
			elemVal, err := elemType.ValueFromTerraform(ctx, val)
			if err != nil {
				return nil, fmt.Errorf("error converting map element for key %s: %w", key, err)
			}
			elems[key] = elemVal
		}
	}

	return NewMapValue(elems), nil
}
