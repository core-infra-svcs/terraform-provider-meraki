package utils

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
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

// ExtractObjectAttr extracts an object attribute from a hashmap
func ExtractObjectAttr(hashMap map[string]interface{}, key string, attrTypes map[string]attr.Type) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	if value, ok := hashMap[key].(map[string]interface{}); ok {
		objVal, err := types.ObjectValueFrom(context.Background(), attrTypes, value)
		if err.HasError() {
			diags.Append(err...)
		}
		return objVal, diags
	}
	return types.ObjectNull(attrTypes), diags
}

// ExtractListAttr extracts a list attribute from a hashmap
func ExtractListAttr(hashMap map[string]interface{}, key string, elemType attr.Type) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	if value, ok := hashMap[key].([]interface{}); ok {
		var attrValues []attr.Value
		for _, item := range value {
			switch elemType := elemType.(type) {
			case types.ObjectType:
				if mapItem, ok := item.(map[string]interface{}); ok {
					objAttrs := make(map[string]attr.Value)
					for k, v := range mapItem {
						switch attrType := elemType.AttrTypes[k]; attrType.(type) {
						case basetypes.StringType:
							if str, ok := v.(string); ok {
								objAttrs[k] = types.StringValue(str)
							} else {
								diags.AddError("Invalid item type", fmt.Sprintf("Expected string value for %s", k))
								return types.ListNull(elemType), diags
							}
						case basetypes.Int64Type:
							if i, ok := v.(int64); ok {
								objAttrs[k] = types.Int64Value(i)
							} else if f, ok := v.(float64); ok { // Handle JSON numbers as float64
								objAttrs[k] = types.Int64Value(int64(f))
							} else {
								diags.AddError("Invalid item type", fmt.Sprintf("Expected int64 value for %s", k))
								return types.ListNull(elemType), diags
							}
						case types.ListType:
							if listItems, ok := v.([]interface{}); ok {
								if listType, ok := attrType.(types.ListType); ok {
									subListElemType := listType.ElemType
									subList, subDiags := ExtractListAttr(map[string]interface{}{k: listItems}, k, subListElemType)
									if subDiags.HasError() {
										diags.Append(subDiags...)
										return types.ListNull(elemType), diags
									}
									objAttrs[k] = subList
								} else {
									diags.AddError("Invalid element type", "Could not assert to types.ListType")
									return types.ListNull(elemType), diags
								}
							} else {
								diags.AddError("Invalid item type", fmt.Sprintf("Expected list value for %s", k))
								return types.ListNull(elemType), diags
							}
						default:
							diags.AddError("Unsupported attribute type", fmt.Sprintf("The attribute type of %s is not supported", k))
							return types.ListNull(elemType), diags
						}
					}
					objValue, err := types.ObjectValue(elemType.AttrTypes, objAttrs)
					if err.HasError() {
						diags.Append(err...)
						return types.ListNull(elemType), diags
					}
					attrValues = append(attrValues, objValue)
				} else {
					diags.AddError("Invalid item type", "Expected object value in list")
					return types.ListNull(elemType), diags
				}
			default:
				diags.AddError("Unsupported element type", "The provided element type is not supported")
				return types.ListNull(elemType), diags
			}
		}

		result, err := types.ListValueFrom(context.Background(), elemType, attrValues)
		if err.HasError() {
			diags.Append(err...)
		}
		return result, diags
	}
	return types.ListNull(elemType), diags
}

func ExtractStringsFromList(array types.List) ([]string, diag.Diagnostics) {
	var diags diag.Diagnostics
	var result []string
	for _, element := range array.Elements() {
		result = append(result, element.(types.String).ValueString())
	}
	return result, diags
}

// ExtractStringSliceAttr extracts a slice of string attributes from a hashmap
func ExtractStringSliceAttr(hashMap map[string]interface{}, key string) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	if value, ok := hashMap[key].([]interface{}); ok {
		var stringValues []attr.Value
		for _, item := range value {
			if str, ok := item.(string); ok {
				stringValues = append(stringValues, types.StringValue(str))
			} else {
				diags.AddError("Invalid item type", "Expected string value in string slice")
				return types.ListNull(types.StringType), diags
			}
		}

		result, err := types.ListValueFrom(context.Background(), types.StringType, stringValues)
		if err.HasError() {
			diags.Append(err...)
		}
		return result, diags
	}
	return types.ListNull(types.StringType), diags
}
