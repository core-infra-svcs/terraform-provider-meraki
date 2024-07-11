package networks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/meraki/dashboard-api-go/client"
)

type NetworksSnmpResource struct {
	client *client.APIClient
}

func NewNetworksSnmpResource() resource.Resource {
	return &NetworksSnmpResource{}
}

// NetworksSnmpResourceModel describes the resource data model.
type NetworksSnmpResourceModel struct {
	Id              jsontypes.String                 `tfsdk:"id"`
	NetworkId       jsontypes.String                 `tfsdk:"organization_id" json:"organizationId"`
	Access          jsontypes.String                 `tfsdk:"access" json:"access"`
	CommunityString jsontypes.String                 `tfsdk:"community_string" json:"communityString"`
	Users           []NetworksSnmpResourceModelUsers `tfsdk:"users" json:"users"`
}

type NetworksSnmpResourceModelUsers struct {
	Username   jsontypes.String `tfsdk:"username"`
	Passphrase jsontypes.String `tfsdk:"passphrase"`
}

func (r *NetworksSnmpResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_snmp"
}

func (r *NetworksSnmpResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "network snmp settings.",
		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{
				MarkdownDescription: "Example identifier",
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Network Id",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 31),
				},
				CustomType: jsontypes.StringType,
			},
			"community_string": schema.StringAttribute{
				MarkdownDescription: "The SNMP community string. Only relevant if 'access' is set to 'community'.",
				Computed:            true,
				Optional:            true,
				CustomType:          jsontypes.StringType,
			},
			"access": schema.StringAttribute{
				MarkdownDescription: "The type of SNMP access. Can be one of 'none' (disabled), 'community' (V1/V2c), or 'users' (V3).",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"none", "community", "users"}...),
					stringvalidator.LengthAtLeast(4),
				},
				CustomType: jsontypes.StringType,
			},
			"users": schema.SetNestedAttribute{
				Description: "The list of SNMP users. Only relevant if 'access' is set to 'users'.",
				Optional:    true,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"username": schema.StringAttribute{
							MarkdownDescription: "The username for the SNMP user",
							Computed:            true,
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
						"passphrase": schema.StringAttribute{
							MarkdownDescription: "The passphrase for the SNMP user.",
							Computed:            true,
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
					},
				},
			},
		},
	}
}

func (r *NetworksSnmpResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.APIClient)
}

func (r *NetworksSnmpResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan NetworksSnmpResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var users []client.UpdateNetworkSnmpRequestUsersInner

	for _, user := range plan.Users {

		users = append(users, client.UpdateNetworkSnmpRequestUsersInner{
			Username:   user.Username.ValueString(),
			Passphrase: user.Passphrase.ValueString(),
		})

	}

	snmpCreateSettings := client.UpdateNetworkSnmpRequest{
		Access:          plan.Access.ValueStringPointer(),
		CommunityString: plan.CommunityString.ValueStringPointer(),
		Users:           users,
	}

	_, httpResp, err := r.client.NetworksApi.UpdateNetworkSnmp(ctx, plan.NetworkId.ValueString()).UpdateNetworkSnmpRequest(snmpCreateSettings).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating SNMP settings",
			fmt.Sprintf("Could not create SNMP settings, unexpected error: %s", err),
		)
		return
	}

	// Save data into Terraform state
	if err = json.NewDecoder(httpResp.Body).Decode(&plan); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	plan.Id = plan.NetworkId

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *NetworksSnmpResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state NetworksSnmpResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Reading SNMP settings", map[string]interface{}{"network_id": state.NetworkId.ValueString()})

	_, httpResp, err := r.client.SnmpApi.GetNetworkSnmp(ctx, state.NetworkId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading SNMP settings",
			fmt.Sprintf("Could not read SNMP settings, unexpected error: %s", err.Error()),
		)
		return
	}

	// Save data into Terraform state
	if err = json.NewDecoder(httpResp.Body).Decode(&state); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	state.Id = state.NetworkId

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *NetworksSnmpResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan NetworksSnmpResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var users []client.UpdateNetworkSnmpRequestUsersInner

	for _, user := range plan.Users {

		users = append(users, client.UpdateNetworkSnmpRequestUsersInner{
			Username:   user.Username.ValueString(),
			Passphrase: user.Passphrase.ValueString(),
		})

	}

	snmpCreateSettings := client.UpdateNetworkSnmpRequest{
		Access:          plan.Access.ValueStringPointer(),
		CommunityString: plan.CommunityString.ValueStringPointer(),
		Users:           users,
	}

	_, httpResp, err := r.client.NetworksApi.UpdateNetworkSnmp(ctx, plan.NetworkId.ValueString()).UpdateNetworkSnmpRequest(snmpCreateSettings).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error updating SNMP settings",
			fmt.Sprintf("Could not update SNMP settings, unexpected error: %s", err),
		)
		return
	}

	// Save data into Terraform state
	if err = json.NewDecoder(httpResp.Body).Decode(&plan); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	plan.Id = plan.NetworkId

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *NetworksSnmpResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state NetworksSnmpResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, _, err := r.client.NetworksApi.UpdateNetworkSnmp(ctx, state.NetworkId.ValueString()).UpdateNetworkSnmpRequest(client.UpdateNetworkSnmpRequest{}).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting SNMP settings",
			fmt.Sprintf("Could not delete SNMP settings, unexpected error: %v", err),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

func (r *NetworksSnmpResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), req.ID)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
