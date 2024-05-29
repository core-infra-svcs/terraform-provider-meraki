package utils

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ExtractStringAttr Extracts a string attribute from a hashmap
func ExtractStringAttr(hashMap map[string]interface{}, key string) (types.String, diag.Diagnostics) {
	var diags diag.Diagnostics
	if value, ok := hashMap[key].(string); ok {
		return types.StringValue(value), diags
	}
	return types.StringNull(), diags
}

// ExtractBoolAttr Extracts a bool attribute from a hashmap
func ExtractBoolAttr(hashMap map[string]interface{}, key string) (types.Bool, diag.Diagnostics) {
	var diags diag.Diagnostics
	if value, ok := hashMap[key].(bool); ok {
		return types.BoolValue(value), diags
	}
	return types.BoolNull(), diags
}

// ExtractInt64Attr Extracts an int64 attribute from a hashmap
func ExtractInt64Attr(hashMap map[string]interface{}, key string) (types.Int64, diag.Diagnostics) {
	var diags diag.Diagnostics
	if value, ok := hashMap[key].(int64); ok {
		return types.Int64Value(value), diags
	}
	return types.Int64Null(), diags
}

// ExtractFloat64Attr Extracts a float attribute from a hashmap
func ExtractFloat64Attr(hashMap map[string]interface{}, key string) (types.Float64, diag.Diagnostics) {
	var diags diag.Diagnostics
	if value, ok := hashMap[key].(float64); ok {
		return types.Float64Value(value), diags
	}
	return types.Float64Null(), diags
}

func ExtractStringsFromList(array types.List) ([]string, diag.Diagnostics) {
	var diags diag.Diagnostics
	var result []string
	for _, element := range array.Elements() {
		result = append(result, element.(types.String).ValueString())
	}
	return result, diags
}

// NewStringDefault returns a struct that implements the defaults.String interface for a resource Schema
func NewStringDefault(defaultValue string) defaults.String {
	return &stringDefault{defaultValue: defaultValue}
}

type stringDefault struct {
	defaultValue string
}

func (d *stringDefault) Description(ctx context.Context) string {
	return fmt.Sprintf("Default value is '%s'", d.defaultValue)
}

func (d *stringDefault) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("Default value is `%s`", d.defaultValue)
}

func (d *stringDefault) DefaultString(ctx context.Context, req defaults.StringRequest, resp *defaults.StringResponse) {
	resp.PlanValue = types.StringValue(d.defaultValue)
}
