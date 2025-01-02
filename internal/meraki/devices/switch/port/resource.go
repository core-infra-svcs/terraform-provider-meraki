package port

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
	"io"
	"net/http"
	"strings"
	"time"
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

	resp.TypeName = req.ProviderTypeName + "_devices_switch_port"
}

// Schema provides a way to define the structure of the resource data.
func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	// The resourceSchema object defines the structure of the resource.
	resp.Schema = schema.Schema{

		MarkdownDescription: "Network Devices resource. This only works for devices associated with a network.",

		// The Attributes map describes the fields of the resource.
		Attributes: portResourceSchema,
	}
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
	var data *resourceModel
	// Unmarshal the plan data into the internal data model struct.
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// Check if there are any errors before proceeding.
	if resp.Diagnostics.HasError() {
		return
	}

	payload, diags := PortResourcePayload(context.Background(), data)
	if diags.HasError() {
		resp.Diagnostics.AddError("Resource Payload Error", fmt.Sprintf("\n%v", diags))
		return
	}

	maxRetries := r.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(r.client.GetConfig().Retry4xxErrorWaitTime)

	// API call function to be passed to retryOn4xx
	apiCall := func() (*openApiClient.GetDeviceSwitchPorts200ResponseInner, *http.Response, error) {
		inline, httpResp, err := r.client.SwitchApi.UpdateDeviceSwitchPort(context.Background(), data.Serial.ValueString(), data.PortId.ValueString()).UpdateDeviceSwitchPortRequest(payload).Execute()
		return inline, httpResp, err
	}

	apiResp, httpResp, err := utils.CustomHttpRequestRetry(ctx, maxRetries, retryDelay, apiCall)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating switch port config",
			fmt.Sprintf("Could not create group policy, unexpected error: %s", err),
		)

		if httpResp != nil {
			var responseBody string
			if httpResp.Body != nil {
				bodyBytes, readErr := io.ReadAll(httpResp.Body)
				if readErr == nil {
					responseBody = string(bodyBytes)
				}
			}
			tflog.Error(ctx, "Failed to create resource", map[string]interface{}{
				"error":          err.Error(),
				"httpStatusCode": httpResp.StatusCode,
				"responseBody":   responseBody,
			})
			resp.Diagnostics.AddError(
				"Error creating resource",
				fmt.Sprintf("HTTP Response: %v\nResponse Body: %s", httpResp, responseBody),
			)
		}
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Create HTTP Client Failure",
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

	data, diags = PortResourceState(ctx, apiResp, data)
	if diags.HasError() {
		resp.Diagnostics.AddError("Resource Response Error", fmt.Sprintf("\n%v", diags))
		return
	}

	// Now set the final state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)

	// Log that the resource was created.
	tflog.Trace(ctx, "created resource")
}

// Read method is responsible for reading an existing resource's state.
func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *resourceModel
	var diags diag.Diagnostics
	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	maxRetries := r.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(r.client.GetConfig().Retry4xxErrorWaitTime)

	// usage of CustomHttpRequestRetry with a strongly typed struct
	apiCall := func() (*openApiClient.GetDeviceSwitchPorts200ResponseInner, *http.Response, error, diag.Diagnostics) {
		inline, httpResp, err := r.client.SwitchApi.GetDeviceSwitchPort(ctx, data.Serial.ValueString(), data.PortId.ValueString()).Execute()

		return inline, httpResp, err, diags
	}

	inlineResp, httpResp, err, tfDiags := utils.CustomHttpRequestRetryStronglyTyped(ctx, maxRetries, retryDelay, apiCall)
	if err != nil {
		if tfDiags.HasError() {
			resp.Diagnostics.AddError("Diagnostics Errors", fmt.Sprintf(" %s", tfDiags.Errors()))
		}
		resp.Diagnostics.AddError("Error reading device switch port", fmt.Sprintf(" %s", err))

		if httpResp != nil {
			var responseBody string
			if httpResp.Body != nil {
				bodyBytes, readErr := io.ReadAll(httpResp.Body)
				if readErr == nil {
					responseBody = string(bodyBytes)
				} else {
					responseBody = fmt.Sprintf("Failed to read response body: %s", readErr)
				}
			} else {
				responseBody = "No response body"
			}
			resp.Diagnostics.AddError("Failed to create resource.",
				fmt.Sprintf("HTTP Status Code: %d, Response Body: %s\n", httpResp.StatusCode, responseBody))
		} else {
			resp.Diagnostics.AddError("HTTP Response is nil", "")
		}

		return
	}

	// Ensure inlineResp is not nil before dereferencing it
	if inlineResp == nil {
		fmt.Printf("Received nil response for device switch port: %s, port ID: %s\n", data.Serial.ValueString(), data.PortId.ValueString())
		return
	}

	// If it's not what you expect, add an error to diagnostics.
	if httpResp != nil && httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Use typedApiResp with the correct type for further processing
	data, diags = PortResourceState(ctx, inlineResp, data)
	if diags.HasError() {
		resp.Diagnostics.AddError("Resource Response Error", fmt.Sprintf("\n%v", diags))
		return
	}

	// Now set the final state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)

	// Log that the resource was read.
	tflog.Trace(ctx, "read resource")
}

// Update function is responsible for updating the state of an existing resource.
func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data *resourceModel
	// Unmarshal the plan data into the internal data model struct.
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// Check if there are any errors before proceeding.
	if resp.Diagnostics.HasError() {
		return
	}

	payload, diags := PortResourcePayload(context.Background(), data)
	if diags.HasError() {
		resp.Diagnostics.AddError("Resource Payload Error", fmt.Sprintf("\n%v", diags))
		return
	}

	maxRetries := r.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(r.client.GetConfig().Retry4xxErrorWaitTime)

	// API call function to be passed to retryOn4xx
	apiCall := func() (*openApiClient.GetDeviceSwitchPorts200ResponseInner, *http.Response, error) {
		inline, httpResp, err := r.client.SwitchApi.UpdateDeviceSwitchPort(context.Background(), data.Serial.ValueString(), data.PortId.ValueString()).UpdateDeviceSwitchPortRequest(payload).Execute()
		return inline, httpResp, err
	}

	inlineResp, httpResp, err := utils.CustomHttpRequestRetry(ctx, maxRetries, retryDelay, apiCall)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating switch port config",
			fmt.Sprintf("Could not create group policy, unexpected error: %s", err),
		)

		if httpResp != nil {
			var responseBody string
			if httpResp != nil && httpResp.Body != nil {
				bodyBytes, readErr := io.ReadAll(httpResp.Body)
				if readErr == nil {
					responseBody = string(bodyBytes)
				}
			}
			tflog.Error(ctx, "Failed to create resource", map[string]interface{}{
				"error":          err.Error(),
				"httpStatusCode": httpResp.StatusCode,
				"responseBody":   responseBody,
			})
			resp.Diagnostics.AddError(
				"Error creating resource",
				fmt.Sprintf("HTTP Response: %v\nResponse Body: %s", httpResp, responseBody),
			)
		}
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

	data, diags = PortResourceState(ctx, inlineResp, data)
	if diags.HasError() {
		resp.Diagnostics.AddError("Resource Response Error", fmt.Sprintf("\n%v", diags))
		return
	}

	// Now set the final state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, data)...)

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

	// Create HTTP request body
	payload := *openApiClient.NewUpdateDeviceSwitchPortRequest()
	payload.SetName("")
	payload.SetTags([]string{})
	payload.SetEnabled(false)
	payload.SetPoeEnabled(false)
	payload.SetType("trunk")
	payload.SetVlan(1)
	payload.SetVoiceVlan(1)
	payload.SetAllowedVlans("1")
	payload.SetAccessPolicyType("Open")

	maxRetries := r.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(r.client.GetConfig().Retry4xxErrorWaitTime)

	// API call function to be passed to retryOn4xx
	apiCall := func() (*openApiClient.GetDeviceSwitchPorts200ResponseInner, *http.Response, error) {
		inline, httpResp, err := r.client.SwitchApi.UpdateDeviceSwitchPort(context.Background(), data.Serial.ValueString(), data.PortId.ValueString()).UpdateDeviceSwitchPortRequest(payload).Execute()
		return inline, httpResp, err
	}

	_, httpResp, err := utils.CustomHttpRequestRetry(ctx, maxRetries, retryDelay, apiCall)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating switch port config",
			fmt.Sprintf("Could not create group policy, unexpected error: %s", err),
		)

		if httpResp != nil {
			var responseBody string
			if httpResp.Body != nil {
				bodyBytes, readErr := io.ReadAll(httpResp.Body)
				if readErr == nil {
					responseBody = string(bodyBytes)
				}
			}
			tflog.Error(ctx, "Failed to create resource", map[string]interface{}{
				"error":          err.Error(),
				"httpStatusCode": httpResp.StatusCode,
				"responseBody":   responseBody,
			})
			resp.Diagnostics.AddError(
				"Error creating resource",
				fmt.Sprintf("HTTP Response: %v\nResponse Body: %s", httpResp, responseBody),
			)
		}
		return
	}

	// If it's not what you expect, add an error to diagnostics.
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
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

	resp.State.RemoveResource(ctx)

	// Log that the resource was deleted.
	tflog.Trace(ctx, "removed resource")
}

// ImportState function is used to import an existing resource into Terraform.
func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: serial, port_id. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("serial"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("port_id"), idParts[1])...)

	if resp.Diagnostics.HasError() {
		return
	}
}
