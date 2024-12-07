package ssid

import (
	"bytes"
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
	"io"
	"net/http"
	"strings"
	"time"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &Resource{}
	_ resource.ResourceWithConfigure   = &Resource{}
	_ resource.ResourceWithImportState = &Resource{}
)

func NewResource() resource.Resource {
	return &Resource{
		typeName: "meraki_networks_wireless_ssids",
	}
}

// Resource defines the resource implementation.
type Resource struct {
	client        *openApiClient.APIClient
	typeName      string
	encryptionKey string
}

func (r *Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.typeName

}

func (r *Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	// Retrieve the encryption key and client from the provider configuration
	client, ok := req.ProviderData.(*openApiClient.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Client Type",
			fmt.Sprintf("Expected *openApiClient.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client

	// Since we are passing only the client directly for resources, we need to handle the encryption key separately
	// Retrieve the encryption key from the context
	encryptionKey, ok := ctx.Value("encryption_key").(string)
	if ok {
		r.encryptionKey = encryptionKey
	} else {
		r.encryptionKey = ""
	}

}

// Create creates the resource and sets the initial Terraform state.
func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var state, plan resourceModel

	// Read the Terraform configuration into the model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare the request payload
	payload, payloadDiags := updateNetworksWirelessSsidsResourcePayload(ctx, &plan)
	if payloadDiags.HasError() {
		tflog.Error(ctx, "Failed to create resource payload", map[string]interface{}{
			"error": payloadDiags,
		})
		resp.Diagnostics.AddError(
			"Error creating ssid payload",
			fmt.Sprintf("Unexpected error: %s", payloadDiags),
		)
		return
	}

	maxRetries := r.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(r.client.GetConfig().Retry4xxErrorWaitTime)

	inlineResp, httpResp, err := utils.CustomHttpRequestRetry[*openApiClient.GetNetworkWirelessSsids200ResponseInner](ctx, maxRetries, retryDelay, func() (*openApiClient.GetNetworkWirelessSsids200ResponseInner, *http.Response, error) {
		inline, respHttp, err := r.client.WirelessApi.UpdateNetworkWirelessSsid(context.Background(), plan.NetworkId.ValueString(), fmt.Sprint(plan.Number.ValueInt64())).UpdateNetworkWirelessSsidRequest(payload).Execute()
		return inline, respHttp, err
	})

	// Capture the response body for logging
	var responseBody string
	if httpResp != nil && httpResp.Body != nil {
		bodyBytes, readErr := io.ReadAll(httpResp.Body)
		if readErr == nil {
			responseBody = string(bodyBytes)
		}
		// Reset the response body so it can be read again later if necessary
		httpResp.Body = io.NopCloser(io.NopCloser(bytes.NewBuffer(bodyBytes)))
	}

	// Check if the error matches a specific condition
	if err != nil {
		// Terminate early if specific error condition is met
		if strings.Contains(responseBody, "Open Roaming certificate 0 not found") {
			tflog.Error(ctx, "Terminating early due to specific error condition", map[string]interface{}{
				"error":        err.Error(),
				"responseBody": responseBody,
			})
			resp.Diagnostics.AddError(
				"HTTP Call Failed",
				fmt.Sprintf("Details: %s", responseBody),
			)
			return
		}

		// Check for the specific unmarshalling error
		if strings.Contains(err.Error(), "json: cannot unmarshal number") && strings.Contains(err.Error(), "GetNetworkWirelessSsids200ResponseInner.minBitrate") {
			tflog.Warn(ctx, "Suppressing unmarshaling error: json: cannot unmarshal number into GetNetworkWirelessSsids200ResponseInner.minBitrate of type int32")
		} else {
			tflog.Error(ctx, "HTTP Call Failed", map[string]interface{}{
				"error":        err.Error(),
				"responseBody": responseBody,
			})
			resp.Diagnostics.AddError(
				"HTTP Call Failed",
				fmt.Sprintf("Details: %s", err.Error()),
			)
		}
		return
	}

	// Check for API success response code
	if httpResp != nil && httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%d", httpResp.StatusCode),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	diags = updateNetworksWirelessSsidsResourceState(ctx, &plan, &state, inlineResp, httpResp)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state, plan resourceModel

	// Read Terraform prior state into the model
	diags := req.State.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	maxRetries := r.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(r.client.GetConfig().Retry4xxErrorWaitTime)

	inlineResp, httpResp, err := utils.CustomHttpRequestRetry[*openApiClient.GetNetworkWirelessSsids200ResponseInner](ctx, maxRetries, retryDelay, func() (*openApiClient.GetNetworkWirelessSsids200ResponseInner, *http.Response, error) {
		inline, respHttp, err := r.client.WirelessApi.GetNetworkWirelessSsid(context.Background(), plan.NetworkId.ValueString(), fmt.Sprint(plan.Number.ValueInt64())).Execute()
		if err != nil {
			// Check for specific error
			if strings.Contains(err.Error(), "json: cannot unmarshal number") && strings.Contains(err.Error(), "GetNetworkWirelessSsids200ResponseInner.minBitrate") {
				tflog.Warn(ctx, "Suppressing unmarshaling error: json: cannot unmarshal number into GetNetworkWirelessSsids200ResponseInner.minBitrate of type int32")
				err = nil
			}

			// Check for specific error
			if strings.Contains(err.Error(), "Open Roaming certificate 0 not found") {
				tflog.Warn(ctx, "Suppressing unmarshaling error: Open Roaming certificate 0 not found")
				err = nil
			}
		}
		return inline, respHttp, err
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Call Failed",
			fmt.Sprintf("Details: %s", err.Error()),
		)
		return
	}

	// Ensure inlineResp and httpResp are not nil before using them
	if inlineResp == nil {
		resp.Diagnostics.AddError(
			"Received nil response",
			"Expected a valid response but received nil",
		)
		return
	}

	if httpResp == nil {
		resp.Diagnostics.AddError(
			"Received nil HTTP response",
			"Expected a valid HTTP response but received nil",
		)
		return
	}

	// Check for API success response code
	if httpResp.StatusCode != 200 {
		responseBody, _ := utils.ReadAndCloseBody(httpResp)
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			utils.NewHttpDiagnostics(httpResp, responseBody),
		)
	}

	diags = updateNetworksWirelessSsidsResourceState(ctx, &plan, &state, inlineResp, httpResp)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var state, plan resourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform plan into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare the request payload
	payload, payloadDiags := updateNetworksWirelessSsidsResourcePayload(ctx, &plan)
	if payloadDiags.HasError() {
		tflog.Error(ctx, "Failed to create resource payload", map[string]interface{}{
			"error": payloadDiags,
		})
		resp.Diagnostics.AddError(
			"Error creating ssid payload",
			fmt.Sprintf("Unexpected error: %s", payloadDiags),
		)
		return
	}

	maxRetries := r.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(r.client.GetConfig().Retry4xxErrorWaitTime)

	inlineResp, httpResp, err := utils.CustomHttpRequestRetry[*openApiClient.GetNetworkWirelessSsids200ResponseInner](ctx, maxRetries, retryDelay, func() (*openApiClient.GetNetworkWirelessSsids200ResponseInner, *http.Response, error) {
		inline, respHttp, err := r.client.WirelessApi.UpdateNetworkWirelessSsid(context.Background(), plan.NetworkId.ValueString(), fmt.Sprint(plan.Number.ValueInt64())).UpdateNetworkWirelessSsidRequest(payload).Execute()
		if err != nil {
			// Check for specific error
			if strings.Contains(err.Error(), "json: cannot unmarshal number") && strings.Contains(err.Error(), "GetNetworkWirelessSsids200ResponseInner.minBitrate") {
				tflog.Warn(ctx, "Suppressing unmarshaling error: json: cannot unmarshal number into GetNetworkWirelessSsids200ResponseInner.minBitrate of type int32")
				err = nil
			}

			// Check for specific error
			if strings.Contains(err.Error(), "Open Roaming certificate 0 not found") {
				tflog.Warn(ctx, "Suppressing unmarshaling error: Open Roaming certificate 0 not found")
				err = nil
			}

		}
		return inline, respHttp, err
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Update Call Failed",
			fmt.Sprintf("Details: %s", err.Error()),
		)
	}

	// Check for API success response code
	if httpResp != nil && httpResp.StatusCode != 200 {
		responseBody, _ := utils.ReadAndCloseBody(httpResp)
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			utils.NewHttpDiagnostics(httpResp, responseBody),
		)
	}

	diags = updateNetworksWirelessSsidsResourceState(ctx, &plan, &state, inlineResp, httpResp)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *resourceModel
	// Read Terraform plan state into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	payload := *openApiClient.NewUpdateNetworkWirelessSsidRequest()
	payload.SetEnabled(false)
	payload.SetName("")
	payload.SetAuthMode("open")
	payload.SetVlanId(1)

	maxRetries := r.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(r.client.GetConfig().Retry4xxErrorWaitTime)

	_, httpResp, err := utils.CustomHttpRequestRetry[*openApiClient.GetNetworkWirelessSsids200ResponseInner](ctx, maxRetries, retryDelay, func() (*openApiClient.GetNetworkWirelessSsids200ResponseInner, *http.Response, error) {
		inline, respHttp, err := r.client.WirelessApi.UpdateNetworkWirelessSsid(context.Background(), state.NetworkId.ValueString(), fmt.Sprint(state.Number.ValueInt64())).UpdateNetworkWirelessSsidRequest(payload).Execute()
		if err != nil {
			// Check for specific error
			if strings.Contains(err.Error(), "json: cannot unmarshal number") && strings.Contains(err.Error(), "GetNetworkWirelessSsids200ResponseInner.minBitrate") {
				tflog.Warn(ctx, "Suppressing unmarshaling error: json: cannot unmarshal number into GetNetworkWirelessSsids200ResponseInner.minBitrate of type int32")
				err = nil
			}

			// Check for specific error
			if strings.Contains(err.Error(), "Open Roaming certificate 0 not found") {
				tflog.Warn(ctx, "Suppressing unmarshaling error: Open Roaming certificate 0 not found")
				err = nil
			}

		}
		return inline, respHttp, err
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Call Failed",
			fmt.Sprintf("Details: %s", err.Error()),
		)
	}

	// Check for API success response code
	if httpResp != nil && httpResp.StatusCode != 200 {
		responseBody, _ := utils.ReadAndCloseBody(httpResp)
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			utils.NewHttpDiagnostics(httpResp, responseBody),
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")
}
