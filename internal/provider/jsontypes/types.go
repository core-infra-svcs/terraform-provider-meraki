package jsontypes

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"reflect"
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

// MapType is a generic map structure for JSON objects.
// It uses Go generics to allow for any value type `V`.
func MapType[T JsonValue]() mapType[T] {
	var v T
	return mapType[T]{
		types.MapType{
			ElemType: v.Type(context.Background()),
		},
	}
}

type mapType[T JsonValue] struct {
	types.MapType
}

func (mt mapType[T]) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	val, err := mt.MapType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}

	mAp := val.(types.Map)

	var item T
	if it, et := item.Type(ctx), mAp.ElementType(ctx); !it.Equal(et) {
		return nil, fmt.Errorf("expected type: %T received %T", it, et)
	}

	return Map[T]{val.(types.Map)}, nil
}

func (mt mapType[T]) Equal(o attr.Type) bool {
	var base attr.Type

	switch typ := o.(type) {
	case mapType[T]:
		base = typ.MapType
	case types.MapType:
		base = typ
	default:
		return false
	}

	return mt.MapType.Equal(base)
}

// ObjectType is a generic structure for JSON objects.
func ObjectType[T any]() (types.ObjectType, error) {

	tType := reflect.TypeOf((*T)(nil)).Elem()

	// Check if T is a struct
	if tType.Kind() != reflect.Struct {
		return types.ObjectType{}, fmt.Errorf("ObjectType[T] requires T to be a struct, got %s", tType.Kind())
	}

	attrTypes := make(map[string]attr.Type)

	for i := 0; i < tType.NumField(); i++ {
		field := tType.Field(i)
		if field.Name == "BaseJsonValue" {
			continue
		}
		// Use the field tag or name to determine the key in attrTypes map
		fieldName := field.Tag.Get("tfsdk")
		if fieldName == "" {
			fieldName = field.Name
		}
		attrTypes[fieldName] = reflectTypeToAttrType(field.Type)
	}

	return types.ObjectType{
		AttrTypes: attrTypes,
	}, nil
}

// reflectTypeToAttrType converts a reflect.Type to an attr.Type.
// This is a simple and direct mapping; adjust according to your needs.
func reflectTypeToAttrType(t reflect.Type) attr.Type {
	switch t.Kind() {
	case reflect.String:
		return types.StringType
	case reflect.Bool:
		return types.BoolType
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return types.Int64Type
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		// Terraform does not directly support unsigned integers.
		return types.Int64Type
	case reflect.Float32, reflect.Float64:
		return types.Float64Type
	// Add other cases as necessary
	default:
		// Default fallback, might use StringType or another appropriate Terraform type.
		return types.StringType
	}
}

type objectType[T JsonValue] struct {
	types.ObjectType
}

func (ot objectType[T]) ValueFromTerraform(ctx context.Context, in tftypes.Value) (attr.Value, error) {
	val, err := ot.ObjectType.ValueFromTerraform(ctx, in)
	if err != nil {
		return nil, err
	}

	object := val.(types.Object)

	var item T

	it := item.Type(ctx)
	for attrName, attrType := range object.AttributeTypes(ctx) {
		if !it.Equal(attrType) {
			return nil, fmt.Errorf("expected all attributes to be of type %s but attribute '%s' is of type %T", it, attrName, attrType)
		}
	}

	return Object[T]{val.(types.Object)}, nil
}

func (ot objectType[T]) Equal(o attr.Type) bool {
	var base attr.Type

	switch typ := o.(type) {
	case objectType[T]:
		base = typ.ObjectType
	case types.ObjectType:
		base = typ
	default:
		return false
	}

	return ot.ObjectType.Equal(base)
}
