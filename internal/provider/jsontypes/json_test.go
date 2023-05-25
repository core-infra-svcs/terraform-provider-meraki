package jsontypes

import (
	"context"
	"encoding/json"
	"fmt"
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
