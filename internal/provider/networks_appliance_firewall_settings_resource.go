package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"

	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
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
var _ resource.Resource = &NetworksApplianceFirewallSettingsResource{}
var _ resource.ResourceWithImportState = &NetworksApplianceFirewallSettingsResource{}

func NewNetworksApplianceFirewallSettingsResource() resource.Resource {
	return &NetworksApplianceFirewallSettingsResource{}
}

// NetworksApplianceFirewallSettingsResource defines the resource implementation.
type NetworksApplianceFirewallSettingsResource struct {
	client *openApiClient.APIClient
}

// NetworksApplianceFirewallSettingsResourceModel describes the resource data model.
type NetworksApplianceFirewallSettingsResourceModel struct {
	Id                 jsontypes.String                                                 `tfsdk:"id"`
	NetworkId          jsontypes.String                                                 `tfsdk:"network_id" json:"network_id"`
	SpoofingProtection NetworksApplianceFirewallSettingsResourceModelSpoofingProtection `tfsdk:"spoofing_protection" json:"spoofingProtection"`
}

type NetworksApplianceFirewallSettingsResourceModelSpoofingProtection struct {
	IpSourceGuard NetworksApplianceFirewallSettingsResourceModelIpSourceGuard `tfsdk:"ip_source_guard" json:"ipSourceGuard"`
}

type NetworksApplianceFirewallSettingsResourceModelIpSourceGuard struct {
	Mode jsontypes.String `tfsdk:"mode" json:"mode"`
}

func (r *NetworksApplianceFirewallSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_appliance_firewall_settings"
}

func (r *NetworksApplianceFirewallSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage network appliance firewall settings.",
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
					stringvalidator.LengthBetween(1, 31),
				},
			},
			"spoofing_protection": schema.SingleNestedAttribute{
				MarkdownDescription: "Spoofing protection settings",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"ip_source_guard": schema.SingleNestedAttribute{
						MarkdownDescription: "IP source address spoofing settings",
						Required:            true,
						Attributes: map[string]schema.Attribute{
							"mode": schema.StringAttribute{
								MarkdownDescription: "Mode of protection.",
								Required:            true,
								CustomType:          jsontypes.StringType,
							},
						},
					},
				},
			},
		},
	}
}

func (r *NetworksApplianceFirewallSettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NetworksApplianceFirewallSettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *NetworksApplianceFirewallSettingsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworksApplianceFirewallSettings := *openApiClient.NewUpdateNetworkApplianceFirewallSettingsRequest()
	var spoofingProtection openApiClient.UpdateNetworkApplianceFirewallSettingsRequestSpoofingProtection
	var ipSourceGuard openApiClient.UpdateNetworkApplianceFirewallSettingsRequestSpoofingProtectionIpSourceGuard
	ipSourceGuard.SetMode(data.SpoofingProtection.IpSourceGuard.Mode.ValueString())
	spoofingProtection.SetIpSourceGuard(ipSourceGuard)
	updateNetworksApplianceFirewallSettings.SetSpoofingProtection(spoofingProtection)

	_, httpResp, err := r.client.SettingsApi.UpdateNetworkApplianceFirewallSettings(context.Background(), data.NetworkId.ValueString()).UpdateNetworkApplianceFirewallSettingsRequest(updateNetworksApplianceFirewallSettings).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
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

func (r *NetworksApplianceFirewallSettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetworksApplianceFirewallSettingsResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := r.client.SettingsApi.GetNetworkApplianceFirewallSettings(context.Background(), data.NetworkId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
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

func (r *NetworksApplianceFirewallSettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data *NetworksApplianceFirewallSettingsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworksApplianceFirewallSettings := *openApiClient.NewUpdateNetworkApplianceFirewallSettingsRequest()
	var spoofingProtection openApiClient.UpdateNetworkApplianceFirewallSettingsRequestSpoofingProtection
	var ipSourceGuard openApiClient.UpdateNetworkApplianceFirewallSettingsRequestSpoofingProtectionIpSourceGuard
	ipSourceGuard.SetMode(data.SpoofingProtection.IpSourceGuard.Mode.ValueString())
	spoofingProtection.SetIpSourceGuard(ipSourceGuard)
	updateNetworksApplianceFirewallSettings.SetSpoofingProtection(spoofingProtection)

	_, httpResp, err := r.client.SettingsApi.UpdateNetworkApplianceFirewallSettings(context.Background(), data.NetworkId.ValueString()).UpdateNetworkApplianceFirewallSettingsRequest(updateNetworksApplianceFirewallSettings).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
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
	tflog.Trace(ctx, "updated resource")
}

func (r *NetworksApplianceFirewallSettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data *NetworksApplianceFirewallSettingsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworksApplianceFirewallSettings := *openApiClient.NewUpdateNetworkApplianceFirewallSettingsRequest()
	var spoofingProtection openApiClient.UpdateNetworkApplianceFirewallSettingsRequestSpoofingProtection
	var ipSourceGuard openApiClient.UpdateNetworkApplianceFirewallSettingsRequestSpoofingProtectionIpSourceGuard
	spoofingProtection.SetIpSourceGuard(ipSourceGuard)
	updateNetworksApplianceFirewallSettings.SetSpoofingProtection(spoofingProtection)

	_, httpResp, err := r.client.SettingsApi.UpdateNetworkApplianceFirewallSettings(context.Background(), data.NetworkId.ValueString()).UpdateNetworkApplianceFirewallSettingsRequest(updateNetworksApplianceFirewallSettings).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
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

	resp.State.RemoveResource(ctx)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")

}

func (r *NetworksApplianceFirewallSettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), req.ID)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
