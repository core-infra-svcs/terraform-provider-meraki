package settings

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

var (
	_ resource.Resource                = &Resource{} // Terraform resource interface
	_ resource.ResourceWithConfigure   = &Resource{} // Interface for resources with configuration methods
	_ resource.ResourceWithImportState = &Resource{} // Interface for resources with import state functionality
)

type Resource struct {
	client *openApiClient.APIClient // APIClient instance for making API requests
}

func NewResource() resource.Resource {
	return &Resource{}
}

// Metadata provides a way to define information about the resource.
func (r *Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	// The TypeName attribute is important as it provides the user-friendly name for the resource/data source.
	// This is the name users will use to reference the resource/data source, and it's also used in the acceptance tests.
	resp.TypeName = req.ProviderTypeName + "_networks_wireless_ssids_splash_settings"
}

// Configure is a method of the Resource interface that Terraform calls to provide the configured provider instance to the resource.
func (r *Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

	// The provider must be properly configured before it can be used.
	if req.ProviderData == nil {
		return
	}

	// Here we expect the provider data to be of type *openApiClient.APIClient.
	client, ok := req.ProviderData.(*openApiClient.APIClient)

	// This is a fatal error and the provider cannot proceed without it.
	// If you see this error, it means there is an issue with the provider setup.
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	// This allows the resource to use the configured provider for any API calls it needs to make.
	r.client = client
}

// Create method is responsible for creating a new resource.
func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state *resourceModel

	// Unmarshal the plan state into the internal state model struct.
	resp.Diagnostics.Append(req.Plan.Get(ctx, &state)...)

	// Check if there are any errors before proceeding.
	if resp.Diagnostics.HasError() {
		return
	}

	payload, payloadErr := NetworksWirelessSsidsSplashSettingsResourcePayload(ctx, state)
	if payloadErr.HasError() {
		resp.Diagnostics.Append(payloadErr...)
	}

	inlineResp, httpResp, err := r.client.SettingsApi.UpdateNetworkWirelessSsidSplashSettings(context.Background(), state.NetworkId.ValueString(), state.Number.ValueString()).UpdateNetworkWirelessSsidSplashSettingsRequest(payload).Execute()
	// If there was an error during API call, add it to diagnostics.
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

		// If there were any errors up to this point, log the plan state and return.
		if resp.Diagnostics.HasError() {
			resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", state))
			return
		}
	}

	stateErr := updateNetworksWirelessSsidsSplashSettingsResourceState(ctx, *inlineResp, state)
	if stateErr.HasError() {
		resp.Diagnostics.Append(stateErr...)
	}

	// Now set the final state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	// Log that the resource was created.
	tflog.Trace(ctx, "created resource")
}

// Read method is responsible for reading an existing resource's state.
func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *resourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	// Check if there are any errors before proceeding.
	if resp.Diagnostics.HasError() {
		return
	}

	inlineResp, httpResp, err := r.client.SettingsApi.GetNetworkWirelessSsidSplashSettings(context.Background(), data.NetworkId.ValueString(), data.Number.ValueString()).Execute()

	// If there was an error during API call, add it to diagnostics.
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			utils.HttpDiagnostics(httpResp),
		)
		return
	}

	// If it's not what you expect, add an error to diagnostics.
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// If there were any errors up to this point, log the state data and return.
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("State Data", fmt.Sprintf("\n%v", data))
		return
	}

	stateErr := updateNetworksWirelessSsidsSplashSettingsResourceState(ctx, *inlineResp, data)
	if stateErr.HasError() {
		resp.Diagnostics.Append(stateErr...)
	}

	// Now set the final state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Log that the resource was read.
	tflog.Trace(ctx, "read resource")
}

// Update function is responsible for updating the state of an existing resource.
func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data *resourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// If there was an error reading the plan, return early.
	if resp.Diagnostics.HasError() {
		return
	}

	payload, payloadErr := NetworksWirelessSsidsSplashSettingsResourcePayload(ctx, data)
	if payloadErr.HasError() {
		resp.Diagnostics.Append(payloadErr...)
	}

	inlineResp, httpResp, err := r.client.SettingsApi.UpdateNetworkWirelessSsidSplashSettings(context.Background(), data.NetworkId.ValueString(), data.Number.ValueString()).UpdateNetworkWirelessSsidSplashSettingsRequest(payload).Execute()

	// If there was an error during API call, add it to diagnostics.
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			utils.HttpDiagnostics(httpResp),
		)
		return
	}

	// If it's not what you expect, add an error to diagnostics.
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// If there were any errors up to this point, log the plan data and return.
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
		return
	}

	stateErr := updateNetworksWirelessSsidsSplashSettingsResourceState(ctx, *inlineResp, data)
	if stateErr.HasError() {
		resp.Diagnostics.Append(stateErr...)
	}

	// Now set the updated state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Log that the resource was updated.
	tflog.Trace(ctx, "updated resource")
}

// Delete function is responsible for deleting a resource.
func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data *resourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	// If there was an error reading the state, return early.
	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworkWirelessSsidSplashSettings := *openApiClient.NewUpdateNetworkWirelessSsidSplashSettingsRequest()

	if !data.SplashUrl.IsUnknown() {
		if !data.SplashUrl.IsNull() {
			updateNetworkWirelessSsidSplashSettings.SetSplashUrl(data.SplashUrl.ValueString())
		}
	}

	_, httpResp, err := r.client.SettingsApi.UpdateNetworkWirelessSsidSplashSettings(context.Background(), data.NetworkId.ValueString(), data.Number.ValueString()).UpdateNetworkWirelessSsidSplashSettingsRequest(updateNetworkWirelessSsidSplashSettings).Execute()

	// If there was an error during API call, add it to diagnostics.
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			utils.HttpDiagnostics(httpResp),
		)
		return
	}

	// If it's not what you expect, add an error to diagnostics.
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	resp.State.RemoveResource(ctx)

	// Log that the resource was deleted.
	tflog.Trace(ctx, "removed resource")
}
