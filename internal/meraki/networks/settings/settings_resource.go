package networksSettings

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &SettingsResource{}
var _ resource.ResourceWithImportState = &SettingsResource{}

func NewNetworksSettingsResource() resource.Resource {
	return &SettingsResource{}
}

// SettingsResource defines the resource implementation.
type SettingsResource struct {
	client *openApiClient.APIClient
}

func (r *SettingsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_settings"
}

func (r *SettingsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = settingsSchema
}

func (r *SettingsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *SettingsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *SettingsResourceModel

	tflog.Info(ctx, "[start] Create Function Call")

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// Initial create API call
	payload, payloadReqDiags := createUpdateHttpReqPayload(ctx, data)
	if payloadReqDiags != nil {
		resp.Diagnostics.Append(payloadReqDiags...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	inlineResp, httpResp, err := r.client.NetworksApi.UpdateNetworkSettings(context.Background(), data.NetworkId.ValueString()).UpdateNetworkSettingsRequest(payload).Execute()
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

	var NetworkSettings200Response SettingsResourceModel
	NetworkSettings200ResponseDiags := NetworkSettings200Response.FromGetNetworkSettings200Response(ctx, data, inlineResp)
	if NetworkSettings200ResponseDiags != nil {
		resp.Diagnostics.Append(NetworkSettings200ResponseDiags...)
		resp.Diagnostics.AddError("Create NetworkSettings200Response Error", fmt.Sprintf("\n%s", NetworkSettings200Response))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &NetworkSettings200Response)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")
	tflog.Info(ctx, "[finish] Create Function Call")
}

func (r *SettingsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *SettingsResourceModel

	tflog.Info(ctx, "[start] Read Function Call")

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	inlineResp, httpResp, err := r.client.NetworksApi.GetNetworkSettings(context.Background(), data.NetworkId.ValueString()).Execute()
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

	var NetworkSettings200Response SettingsResourceModel
	NetworkSettings200ResponseDiags := NetworkSettings200Response.FromGetNetworkSettings200Response(ctx, data, inlineResp)
	if NetworkSettings200ResponseDiags != nil {
		resp.Diagnostics.Append(NetworkSettings200ResponseDiags...)
		resp.Diagnostics.AddError("Read NetworkSettings200Response Error", fmt.Sprintf("\n%s", NetworkSettings200Response))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &NetworkSettings200Response)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
	tflog.Info(ctx, "[finish] Read Function Call")
}

func (r *SettingsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data *SettingsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// Initial create API call
	payload, payloadReqDiags := createUpdateHttpReqPayload(ctx, data)
	if payloadReqDiags != nil {
		resp.Diagnostics.Append(payloadReqDiags...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	inlineResp, httpResp, err := r.client.NetworksApi.UpdateNetworkSettings(context.Background(), data.NetworkId.ValueString()).UpdateNetworkSettingsRequest(payload).Execute()
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
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
		return
	}

	var NetworkSettings200Response SettingsResourceModel
	NetworkSettings200ResponseDiags := NetworkSettings200Response.FromGetNetworkSettings200Response(ctx, data, inlineResp)
	if NetworkSettings200ResponseDiags != nil {
		resp.Diagnostics.Append(NetworkSettings200ResponseDiags...)
		resp.Diagnostics.AddError("Update NetworkSettings200Response Error", fmt.Sprintf("\n%s", NetworkSettings200Response))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &NetworkSettings200Response)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *SettingsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data *SettingsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworkSettings := *openApiClient.NewUpdateNetworkSettingsRequest()
	updateNetworkSettings.SetLocalStatusPageEnabled(true)
	updateNetworkSettings.SetRemoteStatusPageEnabled(false)
	var v openApiClient.GetNetworkSettings200ResponseSecurePort
	v.SetEnabled(false)
	updateNetworkSettings.SetSecurePort(v)
	var l openApiClient.UpdateNetworkSettingsRequestLocalStatusPage
	var a openApiClient.UpdateNetworkSettingsRequestLocalStatusPageAuthentication
	a.SetEnabled(false)
	a.SetPassword("")
	l.SetAuthentication(a)
	updateNetworkSettings.SetLocalStatusPage(l)

	_, httpResp, err := r.client.NetworksApi.UpdateNetworkSettings(context.Background(), data.NetworkId.ValueString()).UpdateNetworkSettingsRequest(updateNetworkSettings).Execute()
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
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	resp.State.RemoveResource(ctx)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")
}

func (r *SettingsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), req.ID)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
