package subnets

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func mapApiResponseToModel(apiResponse []map[string]interface{}, model *DataSourceModel) diag.Diagnostics {
	diagnostics := diag.Diagnostics{}

	// Prepare a slice to store the processed subnets.
	subnets := make([]attr.Value, len(apiResponse))

	for i, rawSubnet := range apiResponse {
		// Extract individual fields from the raw subnet map.
		subnetStr, _ := rawSubnet["subnet"].(string)
		vlanIdFloat, _ := rawSubnet["vlanId"].(float64)
		usedCountFloat, _ := rawSubnet["usedCount"].(float64)
		freeCountFloat, _ := rawSubnet["freeCount"].(float64)

		// Attempt to create an object for the subnet.
		subnetObj, diagErr := types.ObjectValue(
			ResourceAttrTypes(),
			map[string]attr.Value{
				"id":         types.StringValue(model.Id.ValueString()),
				"subnet":     types.StringValue(subnetStr),
				"vlan_id":    types.Int64Value(int64(vlanIdFloat)),
				"used_count": types.Int64Value(int64(usedCountFloat)),
				"free_count": types.Int64Value(int64(freeCountFloat)),
			},
		)
		if diagErr.HasError() {
			diagnostics.Append(diagErr...)
			continue // Skip this subnet but continue processing others.
		}

		subnets[i] = subnetObj
	}

	// Handle diagnostics if any errors occurred during subnet creation.
	if diagnostics.HasError() {
		return diagnostics
	}

	// Attempt to create the final list of subnets.
	subnetList, diagErr := types.ListValue(types.ObjectType{AttrTypes: ResourceAttrTypes()}, subnets)
	if diagErr.HasError() {
		diagnostics.Append(diagErr...)
		return diagnostics
	}

	// Assign the processed subnets to the model's resources.
	model.Resources = subnetList

	return diagnostics
}
