package jsontypes

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"reflect"
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

func TestMap_UnmarshalJSON(t *testing.T) {
	cases := []struct {
		name     string
		data     string
		expected Map
	}{
		{
			name: "Simple string map",
			data: `{"key1": "value1", "key2": "value2"}`,
			expected: NewMapValue(StringType, map[string]attr.Value{
				"key1": StringValue("value1"),
				"key2": StringValue("value2"),
			}),
		},
		{
			name:     "Empty map",
			data:     `{}`,
			expected: NewMapValue(StringType, map[string]attr.Value{}),
		},
		{
			name:     "Null map",
			data:     `null`,
			expected: Map{MapValue: basetypes.NewMapNull(StringType)},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var m Map
			if err := json.Unmarshal([]byte(tc.data), &m); err != nil {
				t.Fatalf("unexpected error during UnmarshalJSON: %s", err)
			}

			if !reflect.DeepEqual(m, tc.expected) {
				t.Errorf("expected map to equal\n%#v\nbut got\n%#v", tc.expected, m)
			}
		})
	}
}

func TestMap_Equal(t *testing.T) {
	map1 := NewMapValue(StringType, map[string]attr.Value{
		"key1": StringValue("value1"),
	})
	map2 := NewMapValue(StringType, map[string]attr.Value{
		"key1": StringValue("value1"),
	})
	map3 := NewMapValue(StringType, map[string]attr.Value{
		"key1": StringValue("different"),
	})

	if !map1.Equal(map2) {
		t.Error("expected map1 to equal map2")
	}

	if map1.Equal(map3) {
		t.Error("expected map1 not to equal map3")
	}
}

func TestMap_Type(t *testing.T) {
	m := NewMapValue(StringType, map[string]attr.Value{
		"key": StringValue("value"),
	})
	ctx := context.Background()
	if mType := m.Type(ctx); !mType.Equal(MapType(StringType)) {
		t.Errorf("expected map type to be MapType[StringType] but got %T", mType)
	}
}
