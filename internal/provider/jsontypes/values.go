package jsontypes

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"sort"
	"strings"
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

// DynamicValue holds a dynamic JSON value.
type DynamicValue struct {
	Value interface{}
}

// Type returns the Terraform type of this dynamic value, which is always dynamic pseudo type in this context.
func (d DynamicValue) Type(ctx context.Context) attr.Type {
	return DynamicType{}
}

// ToTerraformValue converts the DynamicValue to a tftypes.Value.
func (d DynamicValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	// Assuming dynamic handling based on the actual type of d.Value
	switch v := d.Value.(type) {
	case string:
		return tftypes.NewValue(tftypes.String, v), nil
	// Handle other types as needed
	default:
		return tftypes.Value{}, fmt.Errorf("unsupported type: %T", d.Value)
	}
}

// Equal checks if two DynamicValues are equal.
func (d DynamicValue) Equal(o attr.Value) bool {
	other, ok := o.(DynamicValue)
	if !ok {
		return false
	}
	// TODO: Simplistic equality check; add deeper comparison for complex types
	return fmt.Sprintf("%v", d.Value) == fmt.Sprintf("%v", other.Value)
}

// IsNull checks if the DynamicValue is null.
func (d DynamicValue) IsNull() bool {
	return d.Value == nil
}

// IsUnknown checks if the DynamicValue is unknown.
func (d DynamicValue) IsUnknown() bool {
	// Assuming we don't have a specific representation for unknown values
	return false
}

// String returns a string representation of the DynamicValue.
func (d DynamicValue) String() string {
	return fmt.Sprintf("%v", d.Value)
}

// UnmarshalJSON custom unmarshalling to handle dynamic JSON structures.
func (d *DynamicValue) UnmarshalJSON(data []byte) error {
	var raw interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	d.Value = raw
	return nil
}

// MapValue represents a mapping of string keys to dynamic attr.Value values.
type MapValue struct {
	// elements is the mapping of known values in the Map.
	elements map[string]attr.Value
	// elementType is the type of the elements in the Map, which is dynamic in this context.
	elementType attr.Type
	// state represents whether the value is null, unknown, or known.
	state attr.ValueState
}

// NewMapValue constructs a new instance of MapValue with known elements.
func NewMapValue(elements map[string]attr.Value) MapValue {
	return MapValue{
		elements:    elements,
		elementType: DynamicType{}, // Assume all elements are dynamic
		state:       attr.ValueStateKnown,
	}
}

// ElementType returns the dynamic element type for this map.
func (m MapValue) ElementType(ctx context.Context) attr.Type {
	return m.elementType
}

// ToTerraformValue converts the MapValue to a tftypes.Value.
func (m MapValue) ToTerraformValue(ctx context.Context) (tftypes.Value, error) {
	vals := make(map[string]tftypes.Value)
	for key, val := range m.elements {
		tfVal, err := val.ToTerraformValue(ctx)
		if err != nil {
			return tftypes.Value{}, err
		}
		vals[key] = tfVal
	}
	// Always use tftypes.DynamicPseudoType for the map's element type.
	return tftypes.NewValue(tftypes.Map{ElementType: tftypes.DynamicPseudoType}, vals), nil
}

// Equal checks if the MapValue is equal to another attr.Value.
func (m MapValue) Equal(o attr.Value) bool {
	other, ok := o.(MapValue)
	if !ok {
		return false
	}

	if m.state != other.state {
		return false
	}

	if m.state == attr.ValueStateUnknown || m.state == attr.ValueStateNull {
		// If either is unknown or null, we don't need to compare elements.
		return true
	}

	// Ensure both maps have the same set of keys and each key's value is equal.
	if len(m.elements) != len(other.elements) {
		return false
	}
	for key, val := range m.elements {
		otherVal, exists := other.elements[key]
		if !exists || !val.Equal(otherVal) {
			return false
		}
	}

	return true
}

// IsNull checks if the MapValue represents a null value.
func (m MapValue) IsNull() bool {
	return m.state == attr.ValueStateNull
}

// IsUnknown checks if the MapValue represents an unknown value.
func (m MapValue) IsUnknown() bool {
	return m.state == attr.ValueStateUnknown
}

// String returns a human-readable representation of the MapValue.
func (m MapValue) String() string {
	switch m.state {
	case attr.ValueStateNull:
		return "null"
	case attr.ValueStateUnknown:
		return "unknown"
	default:
		var sb strings.Builder
		keys := make([]string, 0, len(m.elements))
		for k := range m.elements {
			keys = append(keys, k)
		}
		sort.Strings(keys) // Ensure consistent ordering
		sb.WriteString("{")
		for i, k := range keys {
			if i > 0 {
				sb.WriteString(", ")
			}
			sb.WriteString(fmt.Sprintf("%q: %s", k, m.elements[k].String()))
		}
		sb.WriteString("}")
		return sb.String()
	}
}

// Helper Functions

// NewMapNull creates a new MapValue representing a null value.
func NewMapNull() MapValue {
	return MapValue{
		state: attr.ValueStateNull,
	}
}

// NewMapUnknown creates a new MapValue representing an unknown value.
func NewMapUnknown() MapValue {
	return MapValue{
		state: attr.ValueStateUnknown,
	}
}
