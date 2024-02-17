package jsontypes

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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

func Int64Value(v int64) Int64 {
	return Int64{
		Int64Value: types.Int64Value(v),
	}
}

func Int64Null() Int64 {
	return Int64{
		Int64Value: types.Int64Null(),
	}
}

// Float64 wraps basetypes.Float64Value to add custom JSON unmarshaling.
type Float64 struct {
	basetypes.Float64Value
}

// UnmarshalJSON custom unmarshaler to handle JSON numbers and nulls.
func (f *Float64) UnmarshalJSON(bytes []byte) error {
	var tmp *float64 // Use pointer to detect JSON null
	if err := json.Unmarshal(bytes, &tmp); err != nil {
		return err
	}
	if tmp != nil {
		f.Float64Value = types.Float64Value(*tmp)
	} else {
		f.Float64Value = types.Float64Null()
	}
	return nil
}

// Float64SemanticEquals returns true if the given Float64Value is semantically equal to the current Float64Value.
// The underlying value *big.Float can store more precise float values then the Go built-in float64 type. This only
// compares the precision of the value that can be represented as the Go built-in float64, which is 53 bits of precision.
func (f Float64) Float64SemanticEquals(ctx context.Context, newValuable basetypes.Float64Valuable) (bool, diag.Diagnostics) {
	var diags diag.Diagnostics

	newValue, ok := newValuable.(Float64)
	if !ok {
		diags.AddError(
			"Semantic Equality Check Error",
			"An unexpected value type was received while performing semantic equality checks. "+
				"Please report this to the provider developers.\n\n"+
				"Expected Value Type: "+fmt.Sprintf("%T", f)+"\n"+
				"Got Value Type: "+fmt.Sprintf("%T", newValuable),
		)

		return false, diags
	}

	return f.ValueFloat64() == newValue.ValueFloat64(), diags
}

// Type satisfies the attr.Value interface, returns the type of Float64.
func (f Float64) Type(_ context.Context) attr.Type {
	return Float64Type
}

// Equal satisfies the attr.Value interface, needed for comparisons.
func (f Float64) Equal(value attr.Value) bool {
	var bv basetypes.Float64Valuable

	switch val := value.(type) {
	case basetypes.Float64Value:
		bv = val
	case Float64:
		bv = val.Float64Value
	default:
		return false
	}
	return f.Float64Value.Equal(bv)
}

func Float64Value(v float64) Float64 {
	return Float64{
		Float64Value: types.Float64Value(v),
	}
}

func Float64Null() Float64 {
	return Float64{
		Float64Value: types.Float64Null(),
	}

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

type Map struct {
	basetypes.MapValue
}

// UnmarshalJSON custom unmarshaler for Map values.
func (m *Map) UnmarshalJSON(bytes []byte) error {
	// Explicitly check for 'null' JSON
	if string(bytes) == "null" {
		// Assuming you have a way to set the Map state to null; adjust as per your implementation
		*m = Map{MapValue: basetypes.NewMapNull(StringType)} // Use the correct type for your context
		return nil
	}

	// Example assuming a map with string values
	var rawMap map[string]string
	if err := json.Unmarshal(bytes, &rawMap); err != nil {
		return err
	}

	elements := make(map[string]attr.Value)
	for key, val := range rawMap {
		elements[key] = StringValue(val) // Assuming StringValue correctly wraps a string as attr.Value
	}

	// Ensure elementType is correctly set before creating the map value.
	// This example hardcodes StringType for demonstration.
	m.MapValue = basetypes.NewMapValueMust(StringType, elements) // Ensure StringType is a correctly initialized attr.Type
	return nil
}

func (m Map) Type(_ context.Context) attr.Type {
	// This should return the specific mapType instance used to create this value,
	// which would typically involve storing that type information within the Map struct.
	// For simplicity, assuming string values:
	return MapType(StringType)
}

func (m Map) Equal(value attr.Value) bool {
	other, ok := value.(Map)
	if !ok {
		return false
	}
	return m.MapValue.Equal(other.MapValue)
}

// NewMapValue is a helper function to create a Map with known elements.
func NewMapValue(elemType attr.Type, elements map[string]attr.Value) Map {
	mapValue, _ := basetypes.NewMapValue(elemType, elements)
	return Map{MapValue: mapValue}
}
