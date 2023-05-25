package jsontypes

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
	"testing"
)

func TestString_Equal(t *testing.T) {
	a, b := StringValue("test"), StringValue("test")
	if !a.Equal(b) {
		t.Fatalf("expected %s == %s, a != b", a.String(), b.String())
	}

	base := types.StringValue("test")
	if !a.Equal(base) {
		t.Fatalf("expected %s == %s, a != base", a.String(), b.String())
	}
}
