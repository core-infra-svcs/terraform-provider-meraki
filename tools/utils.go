package tools

import (
	"fmt"
	"strconv"

	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontype"
	"github.com/hashicorp/terraform-plugin-framework/diag"
)

// MapStringValue Extracts a string from an interface and returns a Terraform type
func MapStringValue(m map[string]interface{}, key string, diags *diag.Diagnostics) jsontype.String {
	var result jsontype.String

	if v := m[key]; v != nil {
		result = jsontype.StringValue(v.(string))
	} else {
		diags.AddWarning(
			"String extraction error",
			fmt.Sprintf("Failed to extract attribute %s from API response: %s", key, v))
		result = jsontype.StringNull()
	}

	return result
}

// MapBoolValue Extracts a boolean from an interface and returns a Terraform type
func MapBoolValue(m map[string]interface{}, key string, diags *diag.Diagnostics) jsontype.Bool {
	var result jsontype.Bool
	if v := m[key]; v != nil {

		if _, ok := v.(string); ok {

			b, _ := strconv.ParseBool(v.(string))
			result = jsontype.BoolValue(b)
		} else {
			diags.AddWarning(
				"Bool extraction error",
				fmt.Sprintf("Failed to extract attribute %s from API response: %s", key, v))
			result = jsontype.BoolValue(v.(bool))
		}

	} else {
		result = jsontype.BoolNull()
	}
	return result
}
