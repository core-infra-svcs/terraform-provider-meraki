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

// MapType represents a map with string keys and JsonValue values.
type mapType struct {
	basetypes.MapType
}

// MapType constructs a new MapType with the provided element type.
func MapType(elemType attr.Type) mapType {
	return mapType{
		MapType: basetypes.MapType{
			ElemType: elemType,
		},
	}
}

func (m mapType) Equal(o attr.Type) bool {
	other, ok := o.(mapType)
	if !ok {
		return false
	}
	return m.ElemType.Equal(other.ElemType)
}

func (m mapType) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	val, err := m.MapType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}
	mapVal := val.(basetypes.MapValue)
	return NewMapValue(m.ElemType, mapVal.Elements()), nil
}
