package jsontypes

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"reflect"
)

// JsonValue interface extends attr.Value with JSON handling capabilities.
type JsonValue interface {
	attr.Value
}

// BaseJsonValue provides a base struct for JSON values, implementing attr.Value.
type BaseJsonValue struct {
	attr.Value // Generic container for the JSON value
}

// NewBaseJsonValue creates a new BaseJsonValue instance encapsulating a generic JSON value.
func NewBaseJsonValue(value attr.Value) *BaseJsonValue {
	return &BaseJsonValue{value}
}

// Type returns the Terraform type of this BaseJsonValue.
func (b BaseJsonValue) Type(ctx context.Context) attr.Type {
	// This could be more sophisticated depending on the actual JSON structure,
	// potentially returning different Terraform types (e.g., types.ObjectType for JSON objects).
	return types.StringType
}

// Equal compares this BaseJsonValue with another attr.Value.
func (b BaseJsonValue) Equal(other attr.Value) bool {
	otherBaseJsonValue, ok := other.(JsonValue)
	if !ok {
		return false // Not the same type
	}

	// Comparing JSON values correctly requires considering the structure and contents.
	// This simplistic approach may not handle all cases accurately.
	thisJSON, err1 := json.Marshal(b.Value)
	otherJSON, err2 := json.Marshal(otherBaseJsonValue.String())
	if err1 != nil || err2 != nil {
		// Could not marshal one of the values for comparison.
		return false
	}

	return string(thisJSON) == string(otherJSON)
}

// IsNull checks if the BaseJsonValue represents a null value.
func (b BaseJsonValue) IsNull() bool {
	// Adjust according to how you represent null values within your JSON structure.
	return b.Value == nil
}

// IsUnknown checks if the BaseJsonValue represents an unknown value.
// Terraform uses "unknown" to represent values that are not yet computed during the plan phase.
func (b BaseJsonValue) IsUnknown() bool {
	// This implementation assumes we don't have a specific representation for "unknown".
	// You might need a more complex state representation to handle unknowns accurately.
	return false
}

// String returns a string representation of the BaseJsonValue.
func (b BaseJsonValue) String() string {
	// This method should return a string representation suitable for logging or debugging.
	// The JSON marshaling here is for demonstration; consider performance and error handling for production use.
	jsonBytes, err := json.Marshal(b.Value)
	if err != nil {
		// Fallback in case of marshaling error.
		return fmt.Sprintf("error marshaling JSON: %v", err)
	}
	return string(jsonBytes)
}

// MarshalJSON allows BaseJsonValue to be used with standard JSON marshaling.
func (b *BaseJsonValue) MarshalJSON() ([]byte, error) {
	return json.Marshal(b.Value)
}

// UnmarshalJSON allows JSON data to be unmarshaled directly into a BaseJsonValue.
func (b *BaseJsonValue) UnmarshalJSON(data []byte) error {
	// Directly unmarshaling into the value interface{} allows flexibility,
	// but be cautious of the potential for runtime type issues.
	return json.Unmarshal(data, &b.Value)
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

type Map[T JsonValue] struct {
	types.Map
}

func (m Map[T]) Type(_ context.Context) attr.Type {
	return MapType[T]()
}

func (m Map[T]) getElemType() attr.Type {
	var i T
	return i.Type(context.Background())
}

// This function converts a map of string keys to JsonValue (T) into a map of string keys to attr.Value.
func mapToStringAttrValue[T attr.Value](vals map[string]T) map[string]attr.Value {
	attrs := make(map[string]attr.Value, len(vals))
	for key, value := range vals {
		attrs[key] = value
	}
	return attrs
}

// MapValue attempts to convert the Map[T] struct to a map[string]attr.Value suitable for Terraform.
func (m Map[T]) MapValue(ctx context.Context) (map[string]attr.Value, error) {
	var goMap map[string]T
	internalValue, _ := m.Map.ToMapValue(ctx) // This should be a tftypes.Value
	if internalValue.IsUnknown() {
		return nil, fmt.Errorf("value is unknown")
	}
	if internalValue.IsNull() {
		return nil, fmt.Errorf("value is null")
	}

	err := internalValue.ElementsAs(ctx, &goMap, false)
	if err != nil {
		return nil, fmt.Errorf("error converting tftypes.Value to go map: %w", err)
	}

	result := make(map[string]attr.Value, len(goMap))

	m.Map = types.MapValueMust(m.getElemType(), result)

	return result, nil
}

// UnmarshalJSON defines how the Map[T] should be unmarshalled from JSON.
func (m *Map[T]) UnmarshalJSON(data []byte) error {
	var tempMap map[string]T
	if err := json.Unmarshal(data, &tempMap); err != nil {
		return err
	}

	convertedMap := mapToStringAttrValue(tempMap)
	m.Map = types.MapValueMust(m.getElemType(), convertedMap)
	return nil
}

// Object is a generic wrapper around types.Object to handle JSON object values.
type Object[T JsonValue] struct {
	types.Object
}

func (o Object[T]) Type(_ context.Context) attr.Type {
	ot, err := ObjectType[T]()
	if err != nil {
		// Panic if unable to obtain the ObjectType.
		// Ensure your application can handle this or that it's an acceptable outcome.
		panic(fmt.Sprintf("Critical error obtaining ObjectType: %v", err))
	}
	return ot
}

func (o Object[T]) getElemType() attr.Type {
	var i T
	return i.Type(context.Background())
}

// getAttrTypes generates a map of string to attr.Type, representing
// the attribute types of the JSON object.
func (o Object[T]) getAttrTypes() map[string]attr.Type {
	result := make(map[string]attr.Type)

	// Obtain the type of T
	tType := reflect.TypeOf((*T)(nil)).Elem()

	// Iterate through the struct fields
	for i := 0; i < tType.NumField(); i++ {
		field := tType.Field(i)
		// Based on the field type, decide which Terraform attr.Type to use

		switch field.Type.Kind() {
		case reflect.String:
			result[field.Name] = types.StringType
		case reflect.Bool:
			result[field.Name] = types.BoolType
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			result[field.Name] = types.Int64Type
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			// Terraform does not have an unsigned integer type; map to Int64 and handle range restrictions separately if necessary
			result[field.Name] = types.Int64Type
		case reflect.Float32, reflect.Float64:
			result[field.Name] = types.Float64Type
		case reflect.Slice:
			elemType := field.Type.Elem()
			// Example for slices of strings; adjust for other element types or add more cases as needed
			if elemType.Kind() == reflect.String {
				result[field.Name] = types.ListType{ElemType: types.StringType}
			}
			// Add additional cases for slices of other types
		case reflect.Map:
			keyType := field.Type.Key()
			elemType := field.Type.Elem()
			// Example for maps with string keys and string values; adjust for other types as needed
			if keyType.Kind() == reflect.String && elemType.Kind() == reflect.String {
				result[field.Name] = types.MapType{ElemType: types.StringType}
			}
			// Add additional cases for maps with other key/value types
		case reflect.Struct:
			// For embedded structs or complex types, you may need to define a custom attr.Type or handle them as ObjectType with further inspection
			// This is a placeholder for how you might begin to handle nested structs
			// result[field.Name] = YourCustomObjectTypeHandlingMethod(field.Type)
		default:
			// Default or error for unsupported types
			// Log or handle unsupported field types as needed
		}
	}

	return result
}

/*
func (s Set[T]) getElemType() attr.Type {
	var i T
	return i.Type(context.Background())
}
*/

// ObjectValue creates a new Object instance from a map of attributes.
func ObjectValue[T JsonValue](attrs map[string]T) Object[T] {
	attrValues := make(map[string]attr.Value, len(attrs))
	for key, value := range attrs {
		attrValues[key] = value
	}
	obj := Object[T]{}
	obj.Object = types.ObjectValueMust(obj.getAttrTypes(), attrValues)
	return obj
}

func (o *Object[T]) UnmarshalJSON(bytes []byte) error {
	var attrs map[string]T
	if err := json.Unmarshal(bytes, &attrs); err != nil {
		return err
	}

	convertedAttrs := make(map[string]attr.Value, len(attrs))
	for key, value := range attrs {
		convertedAttrs[key] = value
	}
	o.Object = types.ObjectValueMust(o.getAttrTypes(), convertedAttrs)
	return nil
}
