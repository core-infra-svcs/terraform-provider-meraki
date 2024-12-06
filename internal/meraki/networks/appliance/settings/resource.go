package settings

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"
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
var _ resource.Resource = &Resource{}
var _ resource.ResourceWithImportState = &Resource{}

func NewNetworksApplianceSettingsResource() resource.Resource {
	return &Resource{}
}

// Resource defines the resource implementation.
type Resource struct {
	client *openApiClient.APIClient
}

// resourceModel describes the resource data model.
type resourceModel struct {
	Id                   jsontypes.String `tfsdk:"id"`
	NetworkId            jsontypes.String `tfsdk:"network_id" json:"network_id"`
	ClientTrackingMethod jsontypes.String `tfsdk:"client_tracking_method"`
	DeploymentMode       jsontypes.String `tfsdk:"deployment_mode"`
	DynamicDnsPrefix     jsontypes.String `tfsdk:"dynamic_dns_prefix"`
	DynamicDnsEnabled    jsontypes.Bool   `tfsdk:"dynamic_dns_enabled"`
	DynamicDnsUrl        jsontypes.String `tfsdk:"dynamic_dns_url"`
}

func (r *Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_appliance_settings"
}

func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage network appliance settings.",
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
			"client_tracking_method": schema.StringAttribute{
				MarkdownDescription: "Client tracking method of a network",
				Required:            true,
				CustomType:          jsontypes.StringType,
				Validators: []validator.String{
					stringvalidator.OneOf("IP address", "MAC address",
						"Unique client identifier"),
				},
			},
			"deployment_mode": schema.StringAttribute{
				MarkdownDescription: "Deployment mode of a network",
				Required:            true,
				CustomType:          jsontypes.StringType,
				Validators: []validator.String{
					stringvalidator.OneOf("passthrough",
						"routed"),
				},
			},
			"dynamic_dns_prefix": schema.StringAttribute{
				MarkdownDescription: "Dynamic DNS url prefix for Dynamic DNS settings for a network. DDNS must be enabled to update",
				Required:            true,
				CustomType:          jsontypes.StringType,
			},
			"dynamic_dns_enabled": schema.BoolAttribute{
				MarkdownDescription: "Dynamic DNS enabled for Dynamic DNS settings for a network",
				Required:            true,
				CustomType:          jsontypes.BoolType,
			},
			"dynamic_dns_url": schema.StringAttribute{
				MarkdownDescription: "Dynamic DNS url. DDNS must be enabled to update",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
		},
	}
}

func (r *Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *resourceModel

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

	data.ClientTrackingMethod = jsontypes.StringValue(inlineResp.GetClientTrackingMethod())
	data.DeploymentMode = jsontypes.StringValue(inlineResp.GetDeploymentMode())
	data.DynamicDnsPrefix = jsontypes.StringValue(inlineResp.DynamicDns.GetPrefix())
	data.DynamicDnsEnabled = jsontypes.BoolValue(inlineResp.DynamicDns.GetEnabled())
	data.DynamicDnsUrl = jsontypes.StringValue(inlineResp.DynamicDns.GetUrl())
	data.Id = data.NetworkId

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *resourceModel

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

	data.ClientTrackingMethod = jsontypes.StringValue(inlineResp.GetClientTrackingMethod())
	data.DeploymentMode = jsontypes.StringValue(inlineResp.GetDeploymentMode())
	data.DynamicDnsPrefix = jsontypes.StringValue(inlineResp.DynamicDns.GetPrefix())
	data.DynamicDnsEnabled = jsontypes.BoolValue(inlineResp.DynamicDns.GetEnabled())
	data.DynamicDnsUrl = jsontypes.StringValue(inlineResp.DynamicDns.GetUrl())
	data.Id = data.NetworkId

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data *resourceModel

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

	data.ClientTrackingMethod = jsontypes.StringValue(inlineResp.GetClientTrackingMethod())
	data.DeploymentMode = jsontypes.StringValue(inlineResp.GetDeploymentMode())
	data.DynamicDnsPrefix = jsontypes.StringValue(inlineResp.DynamicDns.GetPrefix())
	data.DynamicDnsEnabled = jsontypes.BoolValue(inlineResp.DynamicDns.GetEnabled())
	data.DynamicDnsUrl = jsontypes.StringValue(inlineResp.DynamicDns.GetUrl())
	data.Id = data.NetworkId

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data *resourceModel

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

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), req.ID)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
