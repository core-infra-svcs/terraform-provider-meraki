package tools

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
)

// MapStringValue Extracts a string from an interface and returns a Terraform type
func MapStringValue(m map[string]interface{}, key string, diags *diag.Diagnostics) types.String {
	var result types.String

	if v := m[key]; v != nil {
		result = types.StringValue(v.(string))
	} else {
		diags.AddWarning(
			"String extraction error",
			fmt.Sprintf("Failed to extract attribute %s from API response: %s", key, v))
		result = types.StringNull()
	}

	return result
}

// MapBoolValue Extracts a boolean from an interface and returns a Terraform type
func MapBoolValue(m map[string]interface{}, key string, diags *diag.Diagnostics) types.Bool {
	var result types.Bool
	if v := m[key]; v != nil {

		if _, ok := v.(string); ok {

			b, _ := strconv.ParseBool(v.(string))
			result = types.BoolValue(b)
		} else {
			diags.AddWarning(
				"Bool extraction error",
				fmt.Sprintf("Failed to extract attribute %s from API response: %s", key, v))
			result = types.BoolValue(v.(bool))
		}

	} else {
		result = types.BoolNull()
	}
	return result
}

// TODO - MapCustomStructValue - Extracts data from an interface using generics to return a custom type
func MapCustomStructValue[T any](m map[string]interface{}, key string, diags *diag.Diagnostics) T {

	var results T // fmt.Println(reflect.TypeOf(results))

	// string value
	if d := m[key]; d != nil {

		if _, ok := d.(string); ok {
			_ = json.Unmarshal([]byte(d.(string)), &results)
		} else {
			diags.AddWarning(
				"String extraction error",
				fmt.Sprintf("Failed to extract attribute %s from API response: %s", key, d))
		}
	}

	/*
		   // Current Method
			if networks := inlineRespValue["networks"]; networks != nil {
				for _, tv := range networks.([]interface{}) {
					var network OrganizationsAdminResourceModelNetwork
					_ = json.Unmarshal([]byte(tv.(string)), &network)
					data.Networks = append(data.Networks, network)
				}
			} else {
				data.Networks = nil
			}
	*/

	return results
}
