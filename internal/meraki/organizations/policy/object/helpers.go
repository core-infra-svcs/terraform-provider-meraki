package object

import (
	"context"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/meraki/dashboard-api-go/client"
)

func updateOrganizationPolicyObjectResourceState(ctx context.Context, inlineResp map[string]interface{}, state *resourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	//     "id": "1234",
	if state.ObjectId.IsNull() || state.ObjectId.IsUnknown() {
		state.ObjectId, diags = utils.ExtractStringAttr(inlineResp, "id")
		if diags.HasError() {
			diags.AddError("ObjectId Attribute", state.ObjectId.ValueString())
			return diags
		}
	}

	//    "name": "Web Servers - Datacenter 10",
	if state.Name.IsNull() || state.Name.IsUnknown() {
		state.Name, diags = utils.ExtractStringAttr(inlineResp, "name")
		if diags.HasError() {
			diags.AddError("Name Attribute", state.Name.ValueString())
			return diags
		}

	}

	//    "category": "network",
	if state.Category.IsNull() || state.Category.IsUnknown() {
		state.Category, diags = utils.ExtractStringAttr(inlineResp, "category")
		if diags.HasError() {
			diags.AddError("Category Attribute", state.Category.ValueString())
			return diags
		}
	}

	//    "type": "cidr",
	if state.Type.IsNull() || state.Type.IsUnknown() {
		state.Type, diags = utils.ExtractStringAttr(inlineResp, "type")
		if diags.HasError() {
			diags.AddError("Type Attribute", state.Type.ValueString())
			return diags
		}
	}

	//    "cidr": "10.0.0.0/24",
	if state.Cidr.IsNull() || state.Cidr.IsUnknown() {
		state.Cidr, diags = utils.ExtractStringAttr(inlineResp, "cidr")
		if diags.HasError() {
			diags.AddError("Cidr Attribute", state.Cidr.ValueString())
			return diags
		}
	}

	// mask
	if state.Mask.IsNull() || state.Mask.IsUnknown() {
		state.Mask, diags = utils.ExtractStringAttr(inlineResp, "mask")
		if diags.HasError() {
			diags.AddError("Mask Attribute", state.Mask.ValueString())
			return diags
		}
	}

	// fqdn
	if state.Fqdn.IsNull() || state.Fqdn.IsUnknown() {
		state.Fqdn, diags = utils.ExtractStringAttr(inlineResp, "fqdn")
		if diags.HasError() {
			diags.AddError("Fqdn Attribute", state.Fqdn.ValueString())
			return diags
		}
	}

	// ip
	if state.Ip.IsNull() || state.Ip.IsUnknown() {
		state.Ip, diags = utils.ExtractStringAttr(inlineResp, "ip")
		if diags.HasError() {
			diags.AddError("Ip Attribute", state.Ip.ValueString())
			return diags
		}
	}

	// "createdAt": "2018-05-12T00:00:00Z",
	if state.CreatedAt.IsNull() || state.CreatedAt.IsUnknown() {
		state.CreatedAt, diags = utils.ExtractStringAttr(inlineResp, "createdAt")
		if diags.HasError() {
			diags.AddError("CreatedAt Attribute", state.CreatedAt.ValueString())
			return diags
		}
	}

	//    "updatedAt": "2018-05-12T00:00:00Z",
	if state.UpdatedAt.IsNull() || state.UpdatedAt.IsUnknown() {
		state.UpdatedAt, diags = utils.ExtractStringAttr(inlineResp, "updatedAt")
		if diags.HasError() {
			diags.AddError("UpdatedAt Attribute", state.UpdatedAt.ValueString())
			return diags
		}
	}

	//    "groupIds": [ "8" ],
	if state.GroupIds.IsNull() || state.GroupIds.IsUnknown() {
		state.GroupIds, diags = utils.ExtractListStringAttr(inlineResp, "groupIds")
		if diags.HasError() {
			diags.AddError("GroupIds Attribute", state.GroupIds.String())
			return diags
		}
	}

	//    "networkIds": [ "L_12345", "N_123456" ]
	if state.NetworkIds.IsNull() || state.NetworkIds.IsUnknown() {
		state.NetworkIds, diags = utils.ExtractListStringAttr(inlineResp, "networkIds")
		if diags.HasError() {
			diags.AddError("NetworkIds Attribute", state.NetworkIds.String())
			return diags
		}
	}

	// Import ID
	if !state.OrganizationID.IsNull() && !state.OrganizationID.IsUnknown() && !state.ObjectId.IsNull() && !state.ObjectId.IsUnknown() {
		id := state.OrganizationID.ValueString() + "," + state.ObjectId.ValueString()
		state.Id = types.StringValue(id)
	} else {
		state.Id = types.StringNull()
	}

	return diags
}

func OrganizationPolicyObjectResourceCreatePayload(plan resourceModel) (client.CreateOrganizationPolicyObjectRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	payload := client.NewCreateOrganizationPolicyObjectRequestWithDefaults()

	// Name
	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		payload.SetName(plan.Name.ValueString())

	}

	// Category
	if !plan.Category.IsNull() && !plan.Category.IsUnknown() {
		payload.SetCategory(plan.Category.ValueString())

	}

	// Type
	if !plan.Type.IsNull() && !plan.Type.IsUnknown() {
		payload.SetType(plan.Type.ValueString())

	}

	// 	Cidr
	if !plan.Cidr.IsNull() && !plan.Cidr.IsUnknown() {
		payload.SetCidr(plan.Cidr.ValueString())

	}

	// Fqdn
	if !plan.Fqdn.IsNull() && !plan.Fqdn.IsUnknown() {
		payload.SetFqdn(plan.Fqdn.ValueString())

	}

	// Mask
	if !plan.Mask.IsNull() && !plan.Mask.IsUnknown() {
		payload.SetMask(plan.Mask.ValueString())

	}

	// Ip
	if !plan.Ip.IsNull() && !plan.Ip.IsUnknown() {
		payload.SetIp(plan.Ip.ValueString())

	}

	// GroupIds
	if !plan.GroupIds.IsNull() && !plan.GroupIds.IsUnknown() {

		groupIds, err := utils.ListInt64TypeToInt32Array(plan.GroupIds)
		if err.HasError() {
			diags.Append(err...)
		}

		payload.SetGroupIds(groupIds)

	}

	return *payload, diags
}

func OrganizationPolicyObjectResourceUpdatePayload(plan resourceModel) (client.UpdateOrganizationPolicyObjectRequest, diag.Diagnostics) {
	var diags diag.Diagnostics

	payload := client.NewUpdateOrganizationPolicyObjectRequestWithDefaults()

	// Name
	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		payload.SetName(plan.Name.ValueString())

	}

	// 	Cidr
	if !plan.Cidr.IsNull() && !plan.Cidr.IsUnknown() {
		payload.SetCidr(plan.Cidr.ValueString())

	}

	// Fqdn
	if !plan.Fqdn.IsNull() && !plan.Fqdn.IsUnknown() {
		payload.SetFqdn(plan.Fqdn.ValueString())

	}

	// Mask
	if !plan.Mask.IsNull() && !plan.Mask.IsUnknown() {
		payload.SetMask(plan.Mask.ValueString())

	}

	// Ip
	if !plan.Ip.IsNull() && !plan.Ip.IsUnknown() {
		payload.SetIp(plan.Ip.ValueString())

	}

	// GroupIds
	if !plan.GroupIds.IsNull() && !plan.GroupIds.IsUnknown() {

		groupIds, err := utils.ListInt64TypeToInt32Array(plan.GroupIds)
		if err.HasError() {
			diags.Append(err...)
		}

		payload.SetGroupIds(groupIds)

	}

	return *payload, diags
}
