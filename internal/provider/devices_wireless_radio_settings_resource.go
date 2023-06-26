package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &DevicesWirelessRadioSettingsResource{}
	_ resource.ResourceWithConfigure   = &DevicesWirelessRadioSettingsResource{}
	_ resource.ResourceWithImportState = &DevicesWirelessRadioSettingsResource{}
)

func NewDevicesWirelessRadioSettingsResource() resource.Resource {
	return &DevicesWirelessRadioSettingsResource{}
}

// DevicesWirelessRadioSettingsResource defines the resource implementation.
type DevicesWirelessRadioSettingsResource struct {
	client *openApiClient.APIClient
}

// DevicesWirelessRadioSettingsResourceModel describes the resource data model.
type DevicesWirelessRadioSettingsResourceModel struct {
	Id                 jsontypes.String `tfsdk:"id"`
	Serial             jsontypes.String `tfsdk:"serial" json:"serial"`
	RfProfileId        jsontypes.String `tfsdk:"rf_profile_id" json:"rfProfileId"`
	TwoFourGhzSettings struct {
		Channel     jsontypes.Int64 `tfsdk:"channel" json:"channel"`
		TargetPower jsontypes.Int64 `tfsdk:"target_power" json:"targetPower"`
	} `tfsdk:"two_four_ghz_settings" json:"twoFourGhzSettings"`
	FiveGhzSettings struct {
		Channel      jsontypes.Int64 `tfsdk:"channel" json:"channel"`
		ChannelWidth jsontypes.Int64 `tfsdk:"channel_width" json:"channelWidth"`
		TargetPower  jsontypes.Int64 `tfsdk:"target_power" json:"targetPower"`
	} `tfsdk:"five_ghz_settings" json:"fiveGhzSettings"`
}

func (r *DevicesWirelessRadioSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_devices_wireless_radio_settings"
}

func (r *DevicesWirelessRadioSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "DevicesWirelessRadioSettings",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.StringType,
			},
			"serial": schema.StringAttribute{
				MarkdownDescription: "Device serial number",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"rf_profile_id": schema.StringAttribute{
				MarkdownDescription: "The ID of an RF profile to assign to the device. If the value of this parameter is null, the appropriate basic RF profile (indoor or outdoor) will be assigned to the device. Assigning an RF profile will clear ALL manually configured overrides on the device (channel width, channel, power)..",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"two_four_ghz_settings": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"channel": schema.Int64Attribute{
						MarkdownDescription: "Sets a manual channel for 2.4 GHz. Can be '1', '2', '3', '4', '5', '6', '7', '8', '9', '10', '11', '12', '13' or '14' or null for using auto channel..",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.Int64Type,
					},
					"target_power": schema.Int64Attribute{
						MarkdownDescription: "Set a manual target power for 2.4 GHz. Can be between '5' or '30' or null for using auto power range..",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.Int64Type,
					},
				},
			},
			"five_ghz_settings": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"channel": schema.Int64Attribute{
						MarkdownDescription: "Sets a manual channel for 5 GHz. Can be '36', '40', '44', '48', '52', '56', '60', '64', '100', '104', '108', '112', '116', '120', '124', '128', '132', '136', '140', '144', '149', '153', '157', '161', '165', '169', '173' or '177' or null for using auto channel.",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.Int64Type,
					},
					"channel_width": schema.Int64Attribute{
						MarkdownDescription: "Sets a manual channel for 5 GHz. Can be '0', '20', '40', '80' or '160' or null for using auto channel width.",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.Int64Type,
					},
					"target_power": schema.Int64Attribute{
						MarkdownDescription: "Set a manual target power for 5 GHz. Can be between '8' or '30' or null for using auto power range.",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.Int64Type,
					},
				},
			},
		},
	}
}

func (r *DevicesWirelessRadioSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DevicesWirelessRadioSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *DevicesWirelessRadioSettingsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	object25 := openApiClient.NewInlineObject25()
	object25.SetRfProfileId(data.RfProfileId.ValueString())
	settings := openApiClient.NewDevicesSerialWirelessRadioSettingsFiveGhzSettings()
	settings.SetChannelWidth(int32(data.FiveGhzSettings.Channel.ValueInt64()))
	settings.SetTargetPower(int32(data.FiveGhzSettings.TargetPower.ValueInt64()))
	settings.SetChannelWidth(int32(data.FiveGhzSettings.ChannelWidth.ValueInt64()))
	object25.SetFiveGhzSettings(*settings)
	ghzSettings := openApiClient.NewDevicesSerialWirelessRadioSettingsTwoFourGhzSettings()
	ghzSettings.SetTargetPower(int32(data.TwoFourGhzSettings.TargetPower.ValueInt64()))
	ghzSettings.SetChannel(int32(data.TwoFourGhzSettings.Channel.ValueInt64()))
	object25.SetTwoFourGhzSettings(*ghzSettings)

	_, httpResp, err := r.client.WirelessApi.UpdateDeviceWirelessRadioSettings(ctx, data.Serial.ValueString()).UpdateDeviceWirelessRadioSettings(*object25).Execute()

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

	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON Decode issue",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
		return
	}

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")
}

func (r *DevicesWirelessRadioSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *DevicesWirelessRadioSettingsResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := r.client.WirelessApi.GetDeviceWirelessRadioSettings(context.Background(), data.Serial.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}
	// collect diagnostics
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
	}

	// Check for API success inlineResp code
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	}

	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON Decode issue",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
		return
	}
	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *DevicesWirelessRadioSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *DevicesWirelessRadioSettingsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	object25 := openApiClient.NewInlineObject25()
	object25.SetRfProfileId(data.RfProfileId.ValueString())
	settings := openApiClient.NewDevicesSerialWirelessRadioSettingsFiveGhzSettings()
	settings.SetChannelWidth(int32(data.FiveGhzSettings.Channel.ValueInt64()))
	settings.SetTargetPower(int32(data.FiveGhzSettings.TargetPower.ValueInt64()))
	settings.SetChannelWidth(int32(data.FiveGhzSettings.ChannelWidth.ValueInt64()))
	object25.SetFiveGhzSettings(*settings)
	ghzSettings := openApiClient.NewDevicesSerialWirelessRadioSettingsTwoFourGhzSettings()
	ghzSettings.SetTargetPower(int32(data.TwoFourGhzSettings.TargetPower.ValueInt64()))
	ghzSettings.SetChannel(int32(data.TwoFourGhzSettings.Channel.ValueInt64()))
	object25.SetTwoFourGhzSettings(*ghzSettings)

	_, httpResp, err := r.client.WirelessApi.UpdateDeviceWirelessRadioSettings(ctx, data.Serial.ValueString()).UpdateDeviceWirelessRadioSettings(*object25).Execute()

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

	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON Decode issue",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
		return
	}

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *DevicesWirelessRadioSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *DevicesWirelessRadioSettingsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	object25 := openApiClient.NewInlineObject25()

	_, httpResp, err := r.client.WirelessApi.UpdateDeviceWirelessRadioSettings(ctx, data.Serial.ValueString()).UpdateDeviceWirelessRadioSettings(*object25).Execute()

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

	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON Decode issue",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
		return
	}

	data.Id = jsontypes.StringValue("example-id")

	resp.State.RemoveResource(ctx)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")
}

func (r *DevicesWirelessRadioSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("serial"), req.ID)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
