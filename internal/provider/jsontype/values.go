package jsontype

import (
	"context"
	"encoding/json"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

type JsonValue interface {
	attr.Value
}

type Int64 struct {
	basetypes.Int64Value
}

func (i *Int64) UnmarshalJSON(bytes []byte) error {
	i.Int64Value = types.Int64Null()

	var i64 *int64
	if err := json.Unmarshal(bytes, &i64); err != nil {
		return err
	}

	if i64 != nil {
		i.Int64Value = types.Int64Value(*i64)
	}
	return nil
}

func (i Int64) Type(_ context.Context) attr.Type {
	return Int64Type
}

func (i Int64) Equal(value attr.Value) bool {
	var bv basetypes.Int64Valuable

	switch val := value.(type) {
	case basetypes.Int64Value:
		bv = val
	case Int64:
		bv = val.Int64Value
	default:
		return false
	}

	return i.Int64Value.Equal(bv)
}

type String struct {
	basetypes.StringValue
}

func StringValue(s string) String {
	return String{
		types.StringValue(s),
	}
}

func (s *String) UnmarshalJSON(bytes []byte) error {
	s.StringValue = types.StringNull()

	var str *string
	if err := json.Unmarshal(bytes, &str); err != nil {
		return err
	}

	if str != nil {
		s.StringValue = types.StringValue(*str)
	}
	return nil
}

func (s String) Type(_ context.Context) attr.Type {
	return StringType
}

func (s String) Equal(value attr.Value) bool {
	var bv basetypes.StringValue

	switch val := value.(type) {
	case basetypes.StringValue:
		bv = val
	case String:
		bv = val.StringValue
	default:
		return false
	}

	return s.StringValue.Equal(bv)
}

func StringNull() String {
	return String{
		StringValue: types.StringNull(),
	}
}

type Bool struct {
	basetypes.BoolValue
}

func (b Bool) Type(_ context.Context) attr.Type {
	return BoolType
}

func (b Bool) Equal(value attr.Value) bool {
	var bv basetypes.BoolValue

	switch val := value.(type) {
	case basetypes.BoolValue:
		bv = val
	case Bool:
		bv = val.BoolValue
	default:
		return false
	}

	return b.BoolValue.Equal(bv)
}

func (b *Bool) UnmarshalJSON(bytes []byte) error {
	b.BoolValue = types.BoolNull()

	var bl *bool
	if err := json.Unmarshal(bytes, &bl); err != nil {
		return err
	}

	if bl != nil {
		b.BoolValue = types.BoolValue(*bl)
	}
	return nil
}

func BoolValue(b bool) Bool {
	return Bool{
		BoolValue: types.BoolValue(b),
	}
}

func BoolNull() Bool {
	return Bool{
		BoolValue: types.BoolNull(),
	}
}

type Set[T JsonValue] struct {
	types.Set
}

func (s Set[T]) Type(_ context.Context) attr.Type {
	return SetType[T]()
}

func (s Set[T]) getElemType() attr.Type {
	var i T
	return i.Type(context.Background())
}

// []T != []attr.Value, even though T is an attr.Value
// We'll just re-map everything so that we can call SetValueMust
func mapToAttrs[T JsonValue](vals []T) []attr.Value {
	attrs := make([]attr.Value, len(vals))
	for i, item := range vals {
		attrs[i] = item
	}
	return attrs
}

func SetValue[T JsonValue](vals []T) Set[T] {
	s := Set[T]{}

	s.Set = types.SetValueMust(s.getElemType(), mapToAttrs(vals))
	return s
}

func (s *Set[T]) UnmarshalJSON(bytes []byte) error {
	var items []T
	if err := json.Unmarshal(bytes, &items); err != nil {
		return err
	}

	s.Set = types.SetValueMust(s.getElemType(), mapToAttrs(items))

	return nil
}
