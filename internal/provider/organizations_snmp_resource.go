package provider

import (
	"context"
	"encoding/json"
	"fmt"

	openApiClient "github.com/core-infra-svcs/dashboard-api-go/client"

	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &OrganizationsSnmpResource{}
var _ resource.ResourceWithImportState = &OrganizationsSnmpResource{}

func NewOrganizationsSnmpResource() resource.Resource {
	return &OrganizationsSnmpResource{}
}

// OrganizationsSnmpResource defines the resource implementation.
type OrganizationsSnmpResource struct {
	client *openApiClient.APIClient
}

// OrganizationsSnmpResourceModel describes the resource data model.
type OrganizationsSnmpResourceModel struct {
	Id              types.String                     `tfsdk:"id"`
	NetworkId       types.String                     `tfsdk:"network_id"`
	Access          types.String                     `tfsdk:"access"`
	CommunityString types.String                     `tfsdk:"community_string"`
	Users           []NetworksSnmpUsersResourceModel `tfsdk:"users"`
}

type NetworksSnmpUsersResourceModel struct {
	Username   types.String `tfsdk:"username"`
	Passphrase types.String `tfsdk:"passphrase"`
}

func (r *OrganizationsSnmpResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizations_snmp"
}
func (r *OrganizationsSnmpResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "OrganizationsSnmp resource for updating org snmp settings.",
		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{
				MarkdownDescription: "Example identifier",
				Required:            false,
				Optional:            false,
				Computed:            true,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Network Id",
				Required:            true,
			},
			"access": schema.StringAttribute{
				MarkdownDescription: "The type of SNMP access. Can be one of 'none' (disabled), 'community' (V1/V2c), or 'users' (V3).",
				Optional:            true,
				Computed:            true,
			},
			"community_string": schema.StringAttribute{
				MarkdownDescription: "The SNMP community string. Only relevant if 'access' is set to 'community'.",
				Optional:            true,
				Computed:            true,
			},
			"users": schema.SetNestedAttribute{
				Description:         "The list of SNMP users. Only relevant if 'access' is set to 'users'.",
				MarkdownDescription: "The list of SNMP users. Only relevant if 'access' is set to 'users'.",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"username": schema.StringAttribute{
							MarkdownDescription: "The username for the SNMP user",
							Computed:            true,
							Optional:            true,
						},
						"passphrase": schema.StringAttribute{
							MarkdownDescription: "The passphrase for the SNMP user.",
							Computed:            true,
							Optional:            true,
						},
					},
				},
			},
		},
	}
}

func (r *OrganizationsSnmpResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*openApiClient.APIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *OrganizationsSnmpResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *OrganizationsSnmpResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if len(data.NetworkId.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing NetworkId", fmt.Sprintf("Value: %s", data.NetworkId.ValueString()))
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworkSnmp := *openApiClient.NewInlineObject106()
	updateNetworkSnmp.SetAccess(data.Access.ValueString())
	updateNetworkSnmp.SetCommunityString(data.CommunityString.ValueString())
	if len(data.Users) > 0 {
		var usersData []openApiClient.NetworksNetworkIdSnmpUsers
		for _, user := range data.Users {
			var userData openApiClient.NetworksNetworkIdSnmpUsers
			userData.Username = user.Username.ValueString()
			userData.Passphrase = user.Passphrase.ValueString()
			usersData = append(usersData, userData)
		}
		updateNetworkSnmp.SetUsers(usersData)
	}

	inlineResp, httpResp, err := r.client.NetworksApi.UpdateNetworkSnmp(context.Background(), data.Id.ValueString()).UpdateNetworkSnmp(updateNetworkSnmp).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// Check for API success response code
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// collect diagnostics
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%s", data))
		return
	}

	// save into the Terraform state.
	extractHttpResponseNetworkSnmpSettingsResource(context.Background(), inlineResp, data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")
}

func (r *OrganizationsSnmpResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *OrganizationsSnmpResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if len(data.NetworkId.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing NetworkId", fmt.Sprintf("Value: %s", data.NetworkId.ValueString()))
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	inlineResp, httpResp, err := r.client.SnmpApi.GetNetworkSnmp(context.Background(), data.NetworkId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// Check for API success response code
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// collect diagnostics
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%s", data))
		return
	}

	// save into the Terraform state.
	extractHttpResponseNetworkSnmpSettingsResource(context.Background(), inlineResp, data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *OrganizationsSnmpResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data *OrganizationsSnmpResourceModel
	var stateData *OrganizationsSnmpResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)

	if len(data.NetworkId.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing NetworkId", fmt.Sprintf("Value: %s", data.NetworkId.ValueString()))
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworkSnmp := *openApiClient.NewInlineObject106()
	updateNetworkSnmp.SetAccess(data.Access.ValueString())
	updateNetworkSnmp.SetCommunityString(data.CommunityString.ValueString())
	if len(data.Users) > 0 {
		var usersData []openApiClient.NetworksNetworkIdSnmpUsers
		for _, user := range data.Users {
			var userData openApiClient.NetworksNetworkIdSnmpUsers
			userData.Username = user.Username.ValueString()
			userData.Passphrase = user.Passphrase.ValueString()
			usersData = append(usersData, userData)
		}
		updateNetworkSnmp.SetUsers(usersData)
	}

	inlineResp, httpResp, err := r.client.NetworksApi.UpdateNetworkSnmp(context.Background(), data.Id.ValueString()).UpdateNetworkSnmp(updateNetworkSnmp).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// collect diagnostics
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%s", data))
		return
	}

	// save into the Terraform state.
	extractHttpResponseNetworkSnmpSettingsResource(context.Background(), inlineResp, data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *OrganizationsSnmpResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	// No Delete Implementation for SNMP Resource.

}

func extractHttpResponseNetworkSnmpSettingsResource(ctx context.Context, inlineRespValue map[string]interface{}, data *OrganizationsSnmpResourceModel) *OrganizationsSnmpResourceModel {
	data.Id = types.StringValue("example-id")
	data.Access = types.StringValue(inlineRespValue["access"].(string))
	if communiteyString := inlineRespValue["communityString"]; communiteyString != nil {
		data.CommunityString = types.StringValue(inlineRespValue["communityString"].(string))
	} else {
		data.CommunityString = types.StringNull()
	}
	if users := inlineRespValue["users"]; users != nil {
		for _, tv := range users.([]interface{}) {
			var user NetworksSnmpUsersResourceModel
			_ = json.Unmarshal([]byte(tv.(string)), &user)
			data.Users = append(data.Users, user)
		}
	} else {
		data.Users = nil
	}
	return data
}

func (r *OrganizationsSnmpResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), req.ID)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
