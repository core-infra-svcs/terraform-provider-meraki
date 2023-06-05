package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
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
	Id               jsontypes.String                                     `tfsdk:"id"`
	NetworkId        jsontypes.String                                     `tfsdk:"network_id" json:"network_id"`
	Vlan             jsontypes.Float64                                    `tfsdk:"vlan" json:"vlan"`
	UseCombinedPower jsontypes.Bool                                       `tfsdk:"use_combined_power" json:"useCombinedPower"`
	PowerExceptions  []NetworksSwitchSettingsResourceModelPowerExceptions `tfsdk:"power_exceptions" json:"powerExceptions"`
}
type NetworksSwitchSettingsResourceModelPowerExceptions struct {
	Serial    jsontypes.String `tfsdk:"serial" json:"serial"`
	PowerType jsontypes.String `tfsdk:"power_type" json:"powerType"`
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
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Network Id",
				Required:            true,
				CustomType:          jsontypes.StringType,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 31),
				},
			},
			"vlan": schema.Float64Attribute{
				MarkdownDescription: "Management VLAN",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.Float64Type,
			},
			"use_combined_power": schema.BoolAttribute{
				MarkdownDescription: "The use combined Power as the default behavior of secondary power supplies on supported devices.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.BoolType,
			},
			"power_exceptions": schema.ListNestedAttribute{
				Description: "Exceptions on a per switch basis to &quot;useCombinedPower&quot;",
				Required:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"serial": schema.StringAttribute{
							MarkdownDescription: "Serial number of the switch",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"power_type": schema.StringAttribute{
							MarkdownDescription: "Per switch exception (combined, redundant, useNetworkSetting)",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
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

	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworksSwitchSettings := *openApiClient.NewInlineObject129()
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
	} else {
		data.PowerExceptions = nil
	}
	_, httpResp, err := r.client.SettingsApi.UpdateNetworkSwitchSettings(context.Background(), data.NetworkId.ValueString()).UpdateNetworkSwitchSettings(updateNetworksSwitchSettings).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// collect diagnostics
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
	}

	// Check for API success response code
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%s", data))
		return
	}

	// Save data into Terraform state
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")
}

func (r *NetworksSwitchSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetworksSwitchSettingsResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := r.client.SettingsApi.GetNetworkSwitchSettings(context.Background(), data.NetworkId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to get resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// collect diagnostics
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
	}

	// Check for API success response code
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	} else {
		resp.Diagnostics.Append()
	}

	// Save data into Terraform state
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *NetworksSwitchSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data *NetworksSwitchSettingsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworksSwitchSettings := *openApiClient.NewInlineObject129()
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
	} else {
		data.PowerExceptions = nil
	}

	_, httpResp, err := r.client.SettingsApi.UpdateNetworkSwitchSettings(context.Background(), data.NetworkId.ValueString()).UpdateNetworkSwitchSettings(updateNetworksSwitchSettings).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// collect diagnostics
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%s", data))
		return
	}

	// Save data into Terraform state
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *NetworksSwitchSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data *NetworksSwitchSettingsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworksSwitchSettings := *openApiClient.NewInlineObject129()
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
	} else {
		data.PowerExceptions = nil
	}

	_, httpResp, err := r.client.SettingsApi.UpdateNetworkSwitchSettings(context.Background(), data.NetworkId.ValueString()).UpdateNetworkSwitchSettings(updateNetworksSwitchSettings).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// collect diagnostics
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%s", data))
		return
	}

	// Save data into Terraform state
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	resp.State.RemoveResource(ctx)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")

}

func (r *NetworksSwitchSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), req.ID)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
