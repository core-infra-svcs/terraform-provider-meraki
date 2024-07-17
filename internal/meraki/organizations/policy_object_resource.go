package organizations

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/meraki/dashboard-api-go/client"
	"strings"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &OrganizationPolicyObjectResource{}
	_ resource.ResourceWithConfigure   = &OrganizationPolicyObjectResource{}
	_ resource.ResourceWithImportState = &OrganizationPolicyObjectResource{}
)

type OrganizationPolicyObjectResource struct {
	client *client.APIClient
}

func NewOrganizationPolicyObjectResource() resource.Resource {
	return &OrganizationPolicyObjectResource{}
}

type OrganizationPolicyObject struct {
	Id             types.String `tfsdk:"id"`
	OrganizationID types.String `tfsdk:"organization_id"`
	ObjectId       types.String `tfsdk:"object_id"`
	Name           types.String `tfsdk:"name"`
	Category       types.String `tfsdk:"category"`
	Type           types.String `tfsdk:"type"`
	Cidr           types.String `tfsdk:"cidr"`
	Fqdn           types.String `tfsdk:"fqdn"`
	Mask           types.String `tfsdk:"mask"`
	Ip             types.String `tfsdk:"ip"`
	CreatedAt      types.String `tfsdk:"created_at"`
	UpdatedAt      types.String `tfsdk:"updated_at"`
	GroupIds       types.List   `tfsdk:"group_ids"`
	NetworkIds     types.List   `tfsdk:"network_ids"`
}

func (r *OrganizationPolicyObjectResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "meraki_organizations_policy_object"
}

func (r *OrganizationPolicyObjectResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	apiClient, ok := req.ProviderData.(*client.APIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = apiClient
}

func (r *OrganizationPolicyObjectResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "The ID of the resource",
			},
			"object_id": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "Name of a policy object, unique within the organization (alphanumeric, space, dash, or underscore characters only)",
			},
			"organization_id": schema.StringAttribute{
				Required:    true,
				Description: "The ID of the organization",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "Name of a policy object, unique within the organization (alphanumeric, space, dash, or underscore characters only)",
			},
			"category": schema.StringAttribute{
				Required:    true,
				Description: "Category of a policy object (one of: adaptivePolicy, network)",
			},
			"type": schema.StringAttribute{
				Required:    true,
				Description: "Type of a policy object (one of: adaptivePolicyIpv4Cidr, cidr, fqdn, ipAndMask)",
			},
			"cidr": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "CIDR Value of a policy object (e.g. 10.11.12.1/24\")",
			},
			"mask": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "Mask of a policy object (e.g. \"255.255.0.0\")",
			},
			"fqdn": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "Fully qualified domain name of policy object (e.g. \"example.com\")",
			},
			"ip": schema.StringAttribute{
				Computed:    true,
				Optional:    true,
				Description: "IP Address of a policy object (e.g. \"1.2.3.4\")",
			},
			"created_at": schema.StringAttribute{
				Computed:    true,
				Description: "Time Stamp of policy object creation.",
			},
			"updated_at": schema.StringAttribute{
				Computed:    true,
				Description: "Time Stamp of policy object updation.",
			},
			"group_ids": schema.ListAttribute{
				ElementType: types.Int64Type,
				Description: "The IDs of policy object groups the policy object belongs to",
				Computed:    true,
				Optional:    true,
			},
			"network_ids": schema.ListAttribute{
				ElementType: types.StringType,
				Description: "The IDs of the networks that use the policy object.",
				Computed:    true,
			},
		},
	}
}

func updateOrganizationPolicyObjectResourceState(ctx context.Context, inlineResp map[string]interface{}, state *OrganizationPolicyObject) diag.Diagnostics {
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

func OrganizationPolicyObjectResourceCreatePayload(plan OrganizationPolicyObject) (client.CreateOrganizationPolicyObjectRequest, diag.Diagnostics) {
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

func OrganizationPolicyObjectResourceUpdatePayload(plan OrganizationPolicyObject) (client.UpdateOrganizationPolicyObjectRequest, diag.Diagnostics) {
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

func (r *OrganizationPolicyObjectResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan OrganizationPolicyObject

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload, err := OrganizationPolicyObjectResourceCreatePayload(plan)
	if err.HasError() {
		resp.Diagnostics.Append(err...)
		return
	}

	inlineResp, httpResp, httpErr := r.client.OrganizationsApi.CreateOrganizationPolicyObject(ctx, plan.OrganizationID.ValueString()).CreateOrganizationPolicyObjectRequest(payload).Execute()
	if httpErr != nil {
		resp.Diagnostics.AddError(
			"Error creating policy object",
			"Could not create policy object, unexpected error: "+httpErr.Error(),
		)
		return
	}

	// Check for API success response code
	if httpResp != nil && httpResp.StatusCode != 201 {
		responseBody, _ := utils.ReadAndCloseBody(httpResp)
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			utils.NewHttpDiagnostics(httpResp, responseBody),
		)
		return
	}

	diags = updateOrganizationPolicyObjectResourceState(ctx, inlineResp, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *OrganizationPolicyObjectResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state OrganizationPolicyObject

	// Read Terraform prior state into the model
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	inlineResp, httpResp, err := r.client.OrganizationsApi.GetOrganizationPolicyObject(ctx, state.OrganizationID.ValueString(), state.ObjectId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading policy object",
			"Could not read policy object, unexpected error: "+err.Error(),
		)
		return
	}

	// Check for API success response code
	if httpResp != nil && httpResp.StatusCode != 200 {
		responseBody, _ := utils.ReadAndCloseBody(httpResp)
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			utils.NewHttpDiagnostics(httpResp, responseBody),
		)
		return
	}

	diags = updateOrganizationPolicyObjectResourceState(ctx, inlineResp, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *OrganizationPolicyObjectResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan OrganizationPolicyObject

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload, err := OrganizationPolicyObjectResourceUpdatePayload(plan)
	if err.HasError() {
		resp.Diagnostics.Append(err...)
		return
	}

	inlineResp, httpResp, httpErr := r.client.OrganizationsApi.UpdateOrganizationPolicyObject(ctx, plan.OrganizationID.ValueString(), plan.ObjectId.ValueString()).UpdateOrganizationPolicyObjectRequest(payload).Execute()
	if httpErr != nil {
		resp.Diagnostics.AddError(
			"Error updating policy object",
			"Could not update policy object, unexpected error: "+httpErr.Error(),
		)
		return
	}

	// Check for API success response code
	if httpResp != nil && httpResp.StatusCode != 200 {
		responseBody, _ := utils.ReadAndCloseBody(httpResp)
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			utils.NewHttpDiagnostics(httpResp, responseBody),
		)
		return
	}

	diags = updateOrganizationPolicyObjectResourceState(ctx, inlineResp, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set the state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *OrganizationPolicyObjectResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state OrganizationPolicyObject

	// Read Terraform plan state into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.client.OrganizationsApi.DeleteOrganizationPolicyObject(ctx, state.OrganizationID.ValueString(), state.ObjectId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting policy object",
			"Could not delete policy object, unexpected error: "+err.Error(),
		)
		return
	}

	// Check for API success response code
	if httpResp != nil && httpResp.StatusCode != 204 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			"",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)

}

func (r *OrganizationPolicyObjectResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: organization_id,policy_id, number. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("policy_id"), idParts[1])...)

	if resp.Diagnostics.HasError() {
		return
	}
}
