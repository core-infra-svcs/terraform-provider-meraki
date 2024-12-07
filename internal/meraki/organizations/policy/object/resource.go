package object

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/meraki/dashboard-api-go/client"
	"strings"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &Resource{}
	_ resource.ResourceWithConfigure   = &Resource{}
	_ resource.ResourceWithImportState = &Resource{}
)

type Resource struct {
	client *client.APIClient
}

func NewResource() resource.Resource {
	return &Resource{}
}

type resourceModel struct {
	Id             types.String `tfsdk:"id" json:"id"`
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

func (r *Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = "meraki_organizations_policy_object"
}

func (r *Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
				ElementType: types.StringType,
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

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan resourceModel

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

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resourceModel

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

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan resourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state resourceModel

	// Read Terraform prior state into the model
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload, err := OrganizationPolicyObjectResourceUpdatePayload(plan)
	if err.HasError() {
		resp.Diagnostics.Append(err...)
		return
	}

	inlineResp, httpResp, httpErr := r.client.OrganizationsApi.UpdateOrganizationPolicyObject(ctx, plan.OrganizationID.ValueString(), state.ObjectId.ValueString()).UpdateOrganizationPolicyObjectRequest(payload).Execute()
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

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state resourceModel

	// Read Terraform plan state into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.client.OrganizationsApi.DeleteOrganizationPolicyObject(
		ctx,
		state.OrganizationID.ValueString(),
		state.ObjectId.ValueString(),
	).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting policy object",
			fmt.Sprintf("Could not delete policy object: %s", err.Error()),
		)
		return
	}

	// Confirm deletion with the expected status code
	if httpResp != nil && httpResp.StatusCode != 204 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("Received status code %d, expected 204", httpResp.StatusCode),
		)
		return
	}

	resp.State.RemoveResource(ctx)

}

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: organization_id,object_id, number. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("object_id"), idParts[1])...)

	if resp.Diagnostics.HasError() {
		return
	}
}
