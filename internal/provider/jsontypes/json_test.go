package jsontypes

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBool_UnmarshalJSON(t *testing.T) {
	cases := []struct {
		data     string
		expected Bool
	}{
		{data: "true", expected: BoolValue(true)},
		{data: "false", expected: BoolValue(false)},
		{data: "null", expected: BoolNull()},
	}

	for _, c := range cases {
		t.Run(c.data, func(t *testing.T) {
			var b Bool
			if err := json.Unmarshal([]byte(c.data), &b); err != nil {
				t.Fatalf("error while unmarshalling: %s", err)
			}

			if !c.expected.Equal(b) {
				t.Fatalf("expected bool: %s, received: %s", c.expected.String(), b.String())
			}
		})
	}
}

func TestString_UnmarshalJSON(t *testing.T) {
	cases := []struct {
		data     string
		expected String
	}{
		{data: `"test string!!!!"`, expected: StringValue("test string!!!!")},
		{data: "null", expected: StringNull()},
	}

	for _, c := range cases {
		t.Run(c.data, func(t *testing.T) {
			var b String
			if err := json.Unmarshal([]byte(c.data), &b); err != nil {
				t.Fatalf("error while unmarshalling: %s", err)
			}

			if !c.expected.Equal(b) {
				t.Fatalf("expected bool: %s, received: %s", c.expected.String(), b.String())
			}
		})
	}
}

func TestSet_UnmarshalJSON(t *testing.T) {
	const test = `["a", "b", "c"]`
	var s Set[String]
	if err := json.Unmarshal([]byte(test), &s); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	ctx := context.Background()
	fmt.Println(s.ElementType(ctx))

	sType, expectedType := s.Type(ctx), SetType[String]()
	if !sType.Equal(expectedType) {
		t.Fatalf("expected test set to be type %T received %T", expectedType, sType)
	}
}

func TestDynamicType_TerraformType(t *testing.T) {
	dt := DynamicType{}
	if tftypes.DynamicPseudoType.Equal(dt.TerraformType(context.Background())) != true {
		t.Errorf("Expected DynamicType to return tftypes.DynamicPseudoType")
	}
}

func TestDynamicType_ValueFromTerraform(t *testing.T) {
	dt := DynamicType{}

	// Test with a string
	strVal := tftypes.NewValue(tftypes.String, "test string")
	val, err := dt.ValueFromTerraform(context.Background(), strVal)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if val.String() != "test string" {
		t.Errorf("Expected StringValue with 'test string', got %v", val)
	}

	// TODO: Add more cases for other types like []interface{}, map[string]interface{}, etc.
}

func TestMapType_ValueFromTerraform(t *testing.T) {
	// Mocking a tftypes.Value representing a map with string keys and dynamic values
	inputMap := map[string]tftypes.Value{
		"key1": tftypes.NewValue(tftypes.String, "value1"),
	}
	tfMapVal := tftypes.NewValue(tftypes.Map{ElementType: tftypes.DynamicPseudoType}, inputMap)

	mapType := NewMapType(DynamicType{})

	// Attempt to convert the tftypes.Value to an attr.Value using our MapType
	attrVal, err := mapType.ValueFromTerraform(context.Background(), tfMapVal)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Verify attrVal is of the correct type (MapValue) and contains the expected data
	mapVal, ok := attrVal.(MapValue)
	if !ok {
		t.Fatalf("Expected result to be of type MapValue, got %T", attrVal)
	}

	// Verify the map contains the expected 'key1' with value 'value1'
	if mapVal.state != attr.ValueStateKnown {
		t.Fatalf("Expected map value state to be known, got %v", mapVal.state)
	}

	// Assuming MapValue has a method or public field to access its elements.
	// This part may need to be adjusted based on your actual implementation.
	if len(mapVal.elements) != 1 {
		t.Fatalf("Expected map to contain 1 element, got %d", len(mapVal.elements))
	}

	val, exists := mapVal.elements["key1"]
	if !exists {
		t.Fatal("Expected map to contain key 'key1'")
	}

	// Assuming DynamicValue implements attr.Value and can be directly compared or inspected.
	dynamicVal, ok := val.(DynamicValue)
	if !ok {
		t.Fatalf("Expected value to be of type DynamicValue, got %T", val)
	}

	if dynamicVal.Value != "value1" {
		t.Errorf("Expected 'key1' to be 'value1', got '%v'", dynamicVal.Value)
	}
}

func TestDynamicValue_ToTerraformValue(t *testing.T) {
	dv := DynamicValue{Value: "test"}
	tfVal, err := dv.ToTerraformValue(context.Background())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	expected := tftypes.NewValue(tftypes.String, "test")
	if !tfVal.Equal(expected) {
		t.Errorf("Expected %v, got %v", expected, tfVal)
	}
}

func TestMapValue_ToTerraformValue(t *testing.T) {
	mv := NewMapValue(map[string]attr.Value{
		"key": DynamicValue{Value: "value"},
	})
	tfVal, err := mv.ToTerraformValue(context.Background())
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	// Assert tfVal is of the correct tftypes.Type
	if !tfVal.Type().Is(tftypes.Map{ElementType: tftypes.DynamicPseudoType}) {
		t.Fatalf("Expected tfVal to be of tftypes.Map type with ElementType tftypes.DynamicPseudoType, got: %s", tfVal.Type().String())
	}

	// Extract the map from tfVal
	var gotMap map[string]tftypes.Value
	if err := tfVal.As(&gotMap); err != nil {
		t.Fatalf("Error extracting map from tftypes.Value: %v", err)
	}

	// Check for the expected key-value pair
	if val, exists := gotMap["key"]; exists {
		var gotValue string
		if err := val.As(&gotValue); err != nil {
			t.Fatalf("Error extracting string from tftypes.Value: %v", err)
		}
		if gotValue != "value" {
			t.Errorf("Expected value 'value', got '%s'", gotValue)
		}
	} else {
		t.Errorf("Expected key 'key' to exist")
	}
}

func TestMapValue_Equal(t *testing.T) {
	mv1 := NewMapValue(map[string]attr.Value{
		"key": DynamicValue{Value: "value1"},
	})
	mv2 := NewMapValue(map[string]attr.Value{
		"key": DynamicValue{Value: "value2"},
	})
	mv3 := NewMapValue(map[string]attr.Value{
		"key": DynamicValue{Value: "value1"},
	})

	if mv1.Equal(mv2) {
		t.Errorf("Expected mv1 and mv2 to be unequal")
	}
	if !mv1.Equal(mv3) {
		t.Errorf("Expected mv1 and mv3 to be equal")
	}
}

func TestFixedIpAssignmentsMapValue_ToFromTerraformValue(t *testing.T) {
	// Simulate the fixedIpAssignments data structure as a MapValue
	ipAssignment := map[string]attr.Value{
		"ip":   DynamicValue{Value: "1.2.3.4"},
		"name": DynamicValue{Value: "My favorite IP"},
	}
	fixedIpAssignments := NewMapValue(map[string]attr.Value{
		"00:11:22:33:44:55": NewMapValue(ipAssignment),
	})

	// Convert to Terraform value
	tfVal, err := fixedIpAssignments.ToTerraformValue(context.Background())
	assert.NoError(t, err, "ToTerraformValue should not produce an error")

	// Verify the Terraform value is of the correct type and structure
	expectedTfType := tftypes.Map{ElementType: tftypes.DynamicPseudoType}
	tfVal, err = fixedIpAssignments.ToTerraformValue(context.Background())
	assert.NoError(t, err, "ToTerraformValue should not produce an error")

	// Adjust the assertion to accurately reflect the dynamic nature of map elements
	assert.True(t, tfVal.Type().Equal(expectedTfType), "Expected tfVal to be of correct tftypes.Map type")

	// Adjust the assertion to accurately reflect the dynamic nature of map elements
	assert.True(t, tfVal.Type().Equal(expectedTfType), "Expected tfVal to be of correct tftypes.Map type")

	// Convert back from Terraform value to attr.Value
	roundTrippedValue, err := NewMapType(DynamicType{}).ValueFromTerraform(context.Background(), tfVal)
	assert.NoError(t, err, "ValueFromTerraform should not produce an error")

	// Verify the original and round-tripped MapValues are equal
	assert.True(t, fixedIpAssignments.Equal(roundTrippedValue), "Round-tripped value should be equal to the original")
}
