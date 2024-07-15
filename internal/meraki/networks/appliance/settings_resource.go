package appliance

import (
	"context"
	"fmt"
	jsontypes2 "github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
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
	Id                   jsontypes2.String `tfsdk:"id"`
	NetworkId            jsontypes2.String `tfsdk:"network_id" json:"network_id"`
	ClientTrackingMethod jsontypes2.String `tfsdk:"client_tracking_method"`
	DeploymentMode       jsontypes2.String `tfsdk:"deployment_mode"`
	DynamicDnsPrefix     jsontypes2.String `tfsdk:"dynamic_dns_prefix"`
	DynamicDnsEnabled    jsontypes2.Bool   `tfsdk:"dynamic_dns_enabled"`
	DynamicDnsUrl        jsontypes2.String `tfsdk:"dynamic_dns_url"`
}

func (r *NetworksApplianceSettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_appliance_settings"
}

func (r *NetworksApplianceSettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage network appliance settings.",
		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{
				MarkdownDescription: "Example identifier",
				Computed:            true,
				CustomType:          jsontypes2.StringType,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Network Id",
				Required:            true,
				CustomType:          jsontypes2.StringType,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 31),
				},
			},
			"client_tracking_method": schema.StringAttribute{
				MarkdownDescription: "Client tracking method of a network",
				Required:            true,
				CustomType:          jsontypes2.StringType,
				Validators: []validator.String{
					stringvalidator.OneOf("IP address", "MAC address",
						"Unique client identifier"),
				},
			},
			"deployment_mode": schema.StringAttribute{
				MarkdownDescription: "Deployment mode of a network",
				Required:            true,
				CustomType:          jsontypes2.StringType,
				Validators: []validator.String{
					stringvalidator.OneOf("passthrough",
						"routed"),
				},
			},
			"dynamic_dns_prefix": schema.StringAttribute{
				MarkdownDescription: "Dynamic DNS url prefix for Dynamic DNS settings for a network. DDNS must be enabled to update",
				Required:            true,
				CustomType:          jsontypes2.StringType,
			},
			"dynamic_dns_enabled": schema.BoolAttribute{
				MarkdownDescription: "Dynamic DNS enabled for Dynamic DNS settings for a network",
				Required:            true,
				CustomType:          jsontypes2.BoolType,
			},
			"dynamic_dns_url": schema.StringAttribute{
				MarkdownDescription: "Dynamic DNS url. DDNS must be enabled to update",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes2.StringType,
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

	updateNetworksApplianceSettings := *openApiClient.NewUpdateNetworkApplianceSettingsRequest()
	updateNetworksApplianceSettings.SetClientTrackingMethod(data.ClientTrackingMethod.ValueString())
	updateNetworksApplianceSettings.SetDeploymentMode(data.DeploymentMode.ValueString())
	var v openApiClient.UpdateNetworkApplianceSettingsRequestDynamicDns
	v.SetEnabled(data.DynamicDnsEnabled.ValueBool())
	v.SetPrefix(data.DynamicDnsPrefix.ValueString())
	updateNetworksApplianceSettings.SetDynamicDns(v)

	inlineResp, httpResp, err := r.client.SettingsApi.UpdateNetworkApplianceSettings(context.Background(), data.NetworkId.ValueString()).UpdateNetworkApplianceSettingsRequest(updateNetworksApplianceSettings).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			utils.HttpDiagnostics(httpResp),
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

	data.ClientTrackingMethod = jsontypes2.StringValue(inlineResp.GetClientTrackingMethod())
	data.DeploymentMode = jsontypes2.StringValue(inlineResp.GetDeploymentMode())
	data.DynamicDnsPrefix = jsontypes2.StringValue(inlineResp.DynamicDns.GetPrefix())
	data.DynamicDnsEnabled = jsontypes2.BoolValue(inlineResp.DynamicDns.GetEnabled())
	data.DynamicDnsUrl = jsontypes2.StringValue(inlineResp.DynamicDns.GetUrl())
	data.Id = jsontypes2.StringValue("example-id")

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

	inlineResp, httpResp, err := r.client.SettingsApi.GetNetworkApplianceSettings(context.Background(), data.NetworkId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			utils.HttpDiagnostics(httpResp),
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

	data.ClientTrackingMethod = jsontypes2.StringValue(inlineResp.GetClientTrackingMethod())
	data.DeploymentMode = jsontypes2.StringValue(inlineResp.GetDeploymentMode())
	data.DynamicDnsPrefix = jsontypes2.StringValue(inlineResp.DynamicDns.GetPrefix())
	data.DynamicDnsEnabled = jsontypes2.BoolValue(inlineResp.DynamicDns.GetEnabled())
	data.DynamicDnsUrl = jsontypes2.StringValue(inlineResp.DynamicDns.GetUrl())
	data.Id = jsontypes2.StringValue("example-id")

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

	updateNetworksApplianceSettings := *openApiClient.NewUpdateNetworkApplianceSettingsRequest()
	updateNetworksApplianceSettings.SetClientTrackingMethod(data.ClientTrackingMethod.ValueString())
	updateNetworksApplianceSettings.SetDeploymentMode(data.DeploymentMode.ValueString())
	var v openApiClient.UpdateNetworkApplianceSettingsRequestDynamicDns
	v.SetEnabled(data.DynamicDnsEnabled.ValueBool())
	v.SetPrefix(data.DynamicDnsPrefix.ValueString())
	updateNetworksApplianceSettings.SetDynamicDns(v)

	inlineResp, httpResp, err := r.client.SettingsApi.UpdateNetworkApplianceSettings(context.Background(), data.NetworkId.ValueString()).UpdateNetworkApplianceSettingsRequest(updateNetworksApplianceSettings).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			utils.HttpDiagnostics(httpResp),
		)
		return
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

	data.ClientTrackingMethod = jsontypes2.StringValue(inlineResp.GetClientTrackingMethod())
	data.DeploymentMode = jsontypes2.StringValue(inlineResp.GetDeploymentMode())
	data.DynamicDnsPrefix = jsontypes2.StringValue(inlineResp.DynamicDns.GetPrefix())
	data.DynamicDnsEnabled = jsontypes2.BoolValue(inlineResp.DynamicDns.GetEnabled())
	data.DynamicDnsUrl = jsontypes2.StringValue(inlineResp.DynamicDns.GetUrl())
	data.Id = jsontypes2.StringValue("example-id")

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

	updateNetworksApplianceSettings := *openApiClient.NewUpdateNetworkApplianceSettingsRequest()
	var v openApiClient.UpdateNetworkApplianceSettingsRequestDynamicDns
	updateNetworksApplianceSettings.SetDynamicDns(v)

	_, httpResp, err := r.client.SettingsApi.UpdateNetworkApplianceSettings(context.Background(), data.NetworkId.ValueString()).UpdateNetworkApplianceSettingsRequest(updateNetworksApplianceSettings).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			utils.HttpDiagnostics(httpResp),
		)
		return
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
