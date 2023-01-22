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
var _ resource.Resource = &NetworksSwitchSettingsResource{}
var _ resource.ResourceWithImportState = &NetworksSwitchSettingsResource{}

func NewNetworksSwitchSettingsResource() resource.Resource {
	return &NetworksSwitchSettingsResource{}
}

// NetworksSwitchSettingsResource defines the resource implementation.
type NetworksSwitchSettingsResource struct {
	client *openApiClient.APIClient
}

// NetworksSwitchSettingsResourceModel describes the resource data model.
type NetworksSwitchSettingsResourceModel struct {
	Id               types.String                                         `tfsdk:"id"`
	NetworkId        types.String                                         `tfsdk:"network_id"`
	Vlan             types.Float64                                        `tfsdk:"vlan"`
	UseCombinedPower types.Bool                                           `tfsdk:"use_combined_power"`
	PowerExceptions  []NetworksSwitchSettingsResourceModelPowerExceptions `tfsdk:"power_exceptions"`
}
type NetworksSwitchSettingsResourceModelPowerExceptions struct {
	Serial    types.String `tfsdk:"serial"`
	PowerType types.String `tfsdk:"power_type"`
}

func (r *NetworksSwitchSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_switch_settings"
}

func (r *NetworksSwitchSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "NetworksSwitchSettings resource for updating network switch settings.",
		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{
				MarkdownDescription: "Example identifier",
				Optional:            true,
				Computed:            true,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Network Id",
				Required:            true,
			},
			"vlan": schema.Float64Attribute{
				MarkdownDescription: "Management VLAN",
				Optional:            true,
				Computed:            true,
			},
			"use_combined_power": schema.BoolAttribute{
				MarkdownDescription: "The use Combined Power as the default behavior of secondary power supplies on supported devices.",
				Optional:            true,
				Computed:            true,
			},
			"power_exceptions": schema.SetNestedAttribute{
				Description:         "Exceptions on a per switch basis to &quot;useCombinedPower&quot;",
				MarkdownDescription: "Exceptions on a per switch basis to &quot;useCombinedPower&quot;",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"serial": schema.StringAttribute{

							MarkdownDescription: "Serial number of the switch",
							Computed:            true,
							Optional:            true,
						},
						"power_type": schema.StringAttribute{
							MarkdownDescription: "Per switch exception (combined, redundant, useNetworkSetting)",
							Computed:            true,
							Optional:            true,
						},
					},
				},
			},
		},
	}
}

func (r *NetworksSwitchSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NetworksSwitchSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *NetworksSwitchSettingsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if len(data.NetworkId.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing NetworkId", fmt.Sprintf("Value: %s", data.NetworkId.ValueString()))
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworksSwitchSettings := *openApiClient.NewInlineObject127()
	updateNetworksSwitchSettings.SetUseCombinedPower(data.UseCombinedPower.ValueBool())
	updateNetworksSwitchSettings.SetVlan(int32(data.Vlan.ValueFloat64()))
	if len(data.PowerExceptions) > 0 {
		var powerExceptions []openApiClient.NetworksNetworkIdSwitchSettingsPowerExceptions
		for _, attribute := range data.PowerExceptions {
			var powerException openApiClient.NetworksNetworkIdSwitchSettingsPowerExceptions
			powerException.Serial = attribute.Serial.ValueString()
			powerException.PowerType = attribute.PowerType.ValueString()
			powerExceptions = append(powerExceptions, powerException)
		}
		updateNetworksSwitchSettings.SetPowerExceptions(powerExceptions)
	}

	inlineResp, httpResp, err := r.client.SettingsApi.UpdateNetworkSwitchSettings(context.Background(), data.NetworkId.ValueString()).UpdateNetworkSwitchSettings(updateNetworksSwitchSettings).Execute()
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
	tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%s", data))
		return
	}

	// save into the Terraform state.
	extractHttpResponseNetworkSwitchSettingsResource(ctx, inlineResp, data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")
}

func (r *NetworksSwitchSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetworksSwitchSettingsResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if len(data.NetworkId.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing NetworkId", fmt.Sprintf("Value: %s", data.NetworkId.ValueString()))
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	inlineResp, httpResp, err := r.client.SettingsApi.GetNetworkSwitchSettings(context.Background(), data.NetworkId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get resource",
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
	tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	}

	// save into the Terraform state.
	extractHttpResponseNetworkSwitchSettingsResource(ctx, inlineResp, data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *NetworksSwitchSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data *NetworksSwitchSettingsResourceModel
	var stateData *NetworksSwitchSettingsResourceModel

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

	updateNetworksSwitchSettings := *openApiClient.NewInlineObject127()
	updateNetworksSwitchSettings.SetUseCombinedPower(data.UseCombinedPower.ValueBool())
	updateNetworksSwitchSettings.SetVlan(int32(data.Vlan.ValueFloat64()))
	if len(data.PowerExceptions) > 0 {
		var powerExceptions []openApiClient.NetworksNetworkIdSwitchSettingsPowerExceptions
		for _, attribute := range data.PowerExceptions {
			var powerException openApiClient.NetworksNetworkIdSwitchSettingsPowerExceptions
			powerException.Serial = attribute.Serial.ValueString()
			powerException.PowerType = attribute.PowerType.ValueString()
			powerExceptions = append(powerExceptions, powerException)
		}
		updateNetworksSwitchSettings.SetPowerExceptions(powerExceptions)
	}

	inlineResp, httpResp, err := r.client.SettingsApi.UpdateNetworkSwitchSettings(context.Background(), data.NetworkId.ValueString()).UpdateNetworkSwitchSettings(updateNetworksSwitchSettings).Execute()
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
	tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%s", data))
		return
	}

	extractHttpResponseNetworkSwitchSettingsResource(ctx, inlineResp, data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *NetworksSwitchSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	// No Delete for Network Switch Settings

}

func extractHttpResponseNetworkSwitchSettingsResource(ctx context.Context, inlineRespValue map[string]interface{}, data *NetworksSwitchSettingsResourceModel) *NetworksSwitchSettingsResourceModel {

	data.Id = types.StringValue("example-id")

	if vlan := inlineRespValue["vlan"]; vlan != nil {
		data.Vlan = types.Float64Value(vlan.(float64))
	} else {
		data.Vlan = types.Float64Null()
	}

	if useCombinedPower := inlineRespValue["useCombinedPower"]; useCombinedPower != nil {
		data.UseCombinedPower = types.BoolValue(useCombinedPower.(bool))
	} else {
		data.UseCombinedPower = types.BoolNull()
	}

	if powerExceptions := inlineRespValue["powerExceptions"]; powerExceptions != nil {
		for _, tv := range powerExceptions.([]interface{}) {
			var powerException NetworksSwitchSettingsResourceModelPowerExceptions
			_ = json.Unmarshal([]byte(tv.(string)), &powerException)
			data.PowerExceptions = append(data.PowerExceptions, powerException)
		}
	} else {
		data.PowerExceptions = nil
	}

	return data
}

func (r *NetworksSwitchSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), req.ID)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
