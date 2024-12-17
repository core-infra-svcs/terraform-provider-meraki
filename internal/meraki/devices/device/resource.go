package device

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
	"net/http"
)

var (
	_ resource.Resource                = &Resource{} // Terraform resource interface
	_ resource.ResourceWithConfigure   = &Resource{} // Interface for resources with configuration methods
	_ resource.ResourceWithImportState = &Resource{} // Interface for resources with import state functionality
)

func NewResource() resource.Resource {
	return &Resource{}
}

type Resource struct {
	client *openApiClient.APIClient // APIClient instance for making API requests
}

// Metadata provides a way to define information about the resource.
func (r *Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_devices"
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
	var data *ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload, payloadErr := mapPayload(data)
	if payloadErr.HasError() {
		resp.Diagnostics.AddError("Failed to assemble mapPayload", fmt.Sprintf("%s", payloadErr.Errors()))
		return
	}

	apiCall := func() (map[string]interface{}, *http.Response, error) {
		return r.client.DevicesApi.UpdateDevice(ctx, data.Serial.ValueString()).UpdateDeviceRequest(payload).Execute()
	}

	_, httpResp, err := HandleAPICall(ctx, r.client, apiCall)
	if err != nil {
		resp.Diagnostics.AddError("Create Error", fmt.Sprintf("Could not create resource: %s", err))
		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError("Unexpected HTTP Response", fmt.Sprintf("Status Code: %v", httpResp.StatusCode))
		return
	}

	diags := mapApiResponseToState(ctx, data, httpResp)
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Trace(ctx, "created resource")
}

// Read method is responsible for reading an existing resource's state.
func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiCall := func() (map[string]interface{}, *http.Response, error) {
		return r.client.DevicesApi.GetDevice(ctx, data.Serial.ValueString()).Execute()
	}

	_, httpResp, err := HandleAPICall(ctx, r.client, apiCall)
	if err != nil {
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Could not read resource: %s", err))
		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError("Unexpected HTTP Response", fmt.Sprintf("Status Code: %v", httpResp.StatusCode))
		return
	}

	diags := mapApiResponseToState(ctx, data, httpResp)
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Trace(ctx, "read resource")
}

// Update function is responsible for updating the state of an existing resource.
func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload, payloadErr := mapPayload(data)
	if payloadErr.HasError() {
		resp.Diagnostics.AddError("Failed to assemble mapPayload", fmt.Sprintf("%s", payloadErr.Errors()))
		return
	}

	apiCall := func() (map[string]interface{}, *http.Response, error) {
		return r.client.DevicesApi.UpdateDevice(ctx, data.Serial.ValueString()).UpdateDeviceRequest(payload).Execute()
	}

	_, httpResp, err := HandleAPICall(ctx, r.client, apiCall)
	if err != nil {
		resp.Diagnostics.AddError("Update Error", fmt.Sprintf("Could not update resource: %s", err))
		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError("Unexpected HTTP Response", fmt.Sprintf("Status Code: %v", httpResp.StatusCode))
		return
	}

	diags := mapApiResponseToState(ctx, data, httpResp)
	resp.Diagnostics.Append(diags...)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	tflog.Trace(ctx, "updated resource")
}

// Delete function is responsible for deleting a resource.
func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *ResourceModel

	// Retrieve the current state
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create an empty configuration payload to reset the resource in the API
	updateDevice := openApiClient.NewUpdateDeviceRequest()

	// Reset all configurable fields to their default (nil or empty)
	updateDevice.Name = new(string)
	updateDevice.Tags = []string{}
	updateDevice.Lat = new(float32)
	updateDevice.Lng = new(float32)
	updateDevice.Address = new(string)
	updateDevice.Notes = new(string)
	updateDevice.MoveMapMarker = new(bool)
	updateDevice.FloorPlanId = nil

	// Define the API call function
	apiCall := func() (map[string]interface{}, *http.Response, error) {
		return r.client.DevicesApi.UpdateDevice(ctx, data.Serial.ValueString()).UpdateDeviceRequest(*updateDevice).Execute()
	}

	// Use the helper to handle the API call
	_, httpResp, err := HandleAPICall(ctx, r.client, apiCall)
	if err != nil {
		resp.Diagnostics.AddError(
			"Delete Error",
			fmt.Sprintf("Could not reset resource configuration in API: %s", err),
		)
		return
	}

	// Check the HTTP status code
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response",
			fmt.Sprintf("Status Code: %v", httpResp.StatusCode),
		)
		return
	}

	// Remove the resource from Terraform state
	resp.State.RemoveResource(ctx)

	// Log the successful deletion
	tflog.Trace(ctx, "resource configuration reset in API and removed from state")
}
