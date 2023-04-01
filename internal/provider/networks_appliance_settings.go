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
var _ resource.Resource = &NetworksApplianceSettingsResource{}
var _ resource.ResourceWithImportState = &NetworksApplianceSettingsResource{}

func NewNetworksApplianceSettingsResource() resource.Resource {
	return &NetworksApplianceSettingsResource{}
}

// NetworksApplianceSettingsResource defines the resource implementation.
type NetworksApplianceSettingsResource struct {
	client *openApiClient.APIClient
}

// NetworksApplianceSettingsResourceModel describes the resource data model.
type NetworksApplianceSettingsResourceModel struct {
	Id                   jsontypes.String `tfsdk:"id"`
	NetworkId            jsontypes.String `tfsdk:"network_id" json:"network_id"`
	ClientTrackingMethod jsontypes.String `tfsdk:"client_tracking_method"`
	DeploymentMode       jsontypes.String `tfsdk:"deployment_mode"`
	DynamicDnsPrefix     jsontypes.String `tfsdk:"dynamic_dns_prefix"`
	DynamicDnsUrl        jsontypes.String `tfsdk:"dynamic_dns_url"`
	DynamicDnsEnabled    jsontypes.Bool   `tfsdk:"dynamic_dns_enabled"`
}

func (r *NetworksApplianceSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_appliance_settings"
}

func (r *NetworksApplianceSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "NetworksApplianceSettings resource for updating network appliance settings.",
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
			"client_tracking_method": schema.StringAttribute{
				MarkdownDescription: "Client tracking method of a network",
				Required:            true,
				CustomType:          jsontypes.StringType,
			},
			"deployment_mode": schema.StringAttribute{
				MarkdownDescription: "Deployment mode of a network",
				Required:            true,
				CustomType:          jsontypes.StringType,
			},
			"dynamic_dns_prefix": schema.StringAttribute{
				MarkdownDescription: "Dynamic DNS url prefix. DDNS must be enabled to update",
				Required:            true,
				CustomType:          jsontypes.StringType,
			},
			"dynamic_dns_url": schema.StringAttribute{
				MarkdownDescription: "Dynamic DNS url. DDNS must be enabled to update",
				Optional:            true,
				CustomType:          jsontypes.StringType,
			},
			"dynamic_dns_enabled": schema.BoolAttribute{
				MarkdownDescription: "Dynamic DNS url. DDNS must be enabled to update",
				Required:            true,
				CustomType:          jsontypes.BoolType,
			},
		},
	}
}

func (r *NetworksApplianceSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NetworksApplianceSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *NetworksApplianceSettingsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworksApplianceSettings := *openApiClient.NewInlineObject45()
	updateNetworksApplianceSettings.SetClientTrackingMethod(data.ClientTrackingMethod.ValueString())
	updateNetworksApplianceSettings.SetDeploymentMode(data.DeploymentMode.ValueString())
	var v openApiClient.NetworksNetworkIdApplianceSettingsDynamicDns
	v.SetEnabled(data.DynamicDnsEnabled.ValueBool())
	v.SetPrefix(data.DynamicDnsPrefix.ValueString())

	updateNetworksApplianceSettings.SetDynamicDns(v)
	_, httpResp, err := r.client.SettingsApi.UpdateNetworkApplianceSettings(context.Background(), data.NetworkId.ValueString()).UpdateNetworkApplianceSettings(updateNetworksApplianceSettings).Execute()
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

func (r *NetworksApplianceSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetworksApplianceSettingsResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := r.client.SettingsApi.GetNetworkApplianceSettings(context.Background(), data.NetworkId.ValueString()).Execute()
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

func (r *NetworksApplianceSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data *NetworksApplianceSettingsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworksApplianceSettings := *openApiClient.NewInlineObject45()
	updateNetworksApplianceSettings.SetClientTrackingMethod(data.ClientTrackingMethod.ValueString())
	updateNetworksApplianceSettings.SetDeploymentMode(data.DeploymentMode.ValueString())
	var v openApiClient.NetworksNetworkIdApplianceSettingsDynamicDns
	v.SetEnabled(data.DynamicDnsEnabled.ValueBool())
	v.SetPrefix(data.DynamicDnsPrefix.ValueString())

	updateNetworksApplianceSettings.SetDynamicDns(v)
	_, httpResp, err := r.client.SettingsApi.UpdateNetworkApplianceSettings(context.Background(), data.NetworkId.ValueString()).UpdateNetworkApplianceSettings(updateNetworksApplianceSettings).Execute()
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

func (r *NetworksApplianceSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data *NetworksApplianceSettingsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworksApplianceSettings := *openApiClient.NewInlineObject45()
	updateNetworksApplianceSettings.SetClientTrackingMethod(data.ClientTrackingMethod.ValueString())
	updateNetworksApplianceSettings.SetDeploymentMode(data.DeploymentMode.ValueString())
	var v openApiClient.NetworksNetworkIdApplianceSettingsDynamicDns
	v.SetEnabled(data.DynamicDnsEnabled.ValueBool())
	v.SetPrefix(data.DynamicDnsPrefix.ValueString())

	updateNetworksApplianceSettings.SetDynamicDns(v)
	_, httpResp, err := r.client.SettingsApi.UpdateNetworkApplianceSettings(context.Background(), data.NetworkId.ValueString()).UpdateNetworkApplianceSettings(updateNetworksApplianceSettings).Execute()
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

func (r *NetworksApplianceSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), req.ID)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
