package tools

import (
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"testing"
)

// test extraction of untyped api response into a terraform string object
func TestMapStringValue(t *testing.T) {
	var resp diag.Diagnostics

	ApiResponse := make(map[string]interface{})
	ApiResponse["name"] = `testAdmin`

	got := MapStringValue(ApiResponse, "name", &resp)
	want := types.StringValue("testAdmin")

	if got.ValueString() != want.ValueString() {
		t.Errorf("got: %v, want: %v", got.ValueString(), want.ValueString())
	}

}

// test extraction of untyped api response into a terraform bool object
func TestMapBoolValue(t *testing.T) {
	var resp diag.Diagnostics

	ApiResponse := make(map[string]interface{})
	ApiResponse["hasApiKey"] = `true`

	got := MapBoolValue(ApiResponse, "hasApiKey", &resp)
	want := types.BoolValue(true)

	if got.ValueBool() != want.ValueBool() {
		t.Errorf("got: %v, want: %v", got.ValueBool(), want.ValueBool())
	}

}

// TODO - test extraction of untyped api response into a custom struct
/*
func TestMapCustomStructValue(t *testing.T) {
	var resp diag.Diagnostics

	type exampleStructModelTags struct {
		Tag    types.String `tfsdk:"tag"`
		Access types.String `tfsdk:"access"`
	}

	ApiResponse := make(map[string]interface{})
	ApiResponse["tags"] = `[{
	   "tag": "west",
	   "access": "read-only"
	  }]`

	got := MapCustomStructValue[exampleStructModelTags](ApiResponse, "tags", &resp)

	want := exampleStructModelTags{
		Tag:    types.StringValue("west"),
		Access: types.StringValue("read-only"),
	}

	if got.Tag.ValueString() != want.Tag.ValueString() {
		t.Errorf("got: %v want: %v", got.Tag.ValueString(), want.Tag.ValueString())
		t.Errorf("got: %v want: %v", got.Access.ValueString(), want.Access.ValueString())
	}

}

*/
