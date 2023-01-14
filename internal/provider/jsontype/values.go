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
	if v, ok := value.(basetypes.StringValue); ok {
		return s.StringValue.Equal(v)
	}

	v, ok := value.(String)
	if !ok {
		return false
	}

	return s.StringValue.Equal(v.StringValue)
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
	if v, ok := value.(basetypes.BoolValue); ok {
		return b.BoolValue.Equal(v)
	}

	v, ok := value.(Bool)
	if !ok {
		return false
	}

	return b.BoolValue.Equal(v.BoolValue)
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
	var item T

	ctx := context.Background()

	return Set[T]{
		types.SetValueMust(item.Type(ctx), mapToAttrs(vals)),
	}

}

func (s *Set[T]) UnmarshalJSON(bytes []byte) error {
	var items []T
	if err := json.Unmarshal(bytes, &items); err != nil {
		return err
	}

	var item T
	ctx := context.Background()

	if items == nil {
		s.Set = types.SetNull(item.Type(ctx))
	} else {
		s.Set = types.SetValueMust(item.Type(ctx), mapToAttrs(items))
	}

	return nil
}
