package vlan

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &Resource{}
	_ resource.ResourceWithConfigure   = &Resource{}
	_ resource.ResourceWithImportState = &Resource{}
)

func NewNetworksApplianceVLANResource() resource.Resource {
	return &Resource{}
}

// Resource defines the resource implementation.
type Resource struct {
	client *openApiClient.APIClient
}

func (r *Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_appliance_vlan"
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
	var data *NetworksApplianceVLANModel

	// Log the received request
	tflog.Info(ctx, "[start] CREATE Function Call")
	tflog.Trace(ctx, "Create Function Called", map[string]interface{}{
		"request": req,
	})

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Initial create API call
	payload, payloadReqDiags := CreateHttpReqPayload(ctx, data)
	if payloadReqDiags != nil {
		resp.Diagnostics.Append(payloadReqDiags...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	inlineResp, httpResp, err := r.client.ApplianceApi.CreateNetworkApplianceVlan(ctx, data.NetworkId.ValueString()).CreateNetworkApplianceVlanRequest(payload).Execute()

	// Meraki API seems to return http status code 201 as an error.
	if err != nil && httpResp.StatusCode != 201 {
		resp.Diagnostics.AddError(
			"HTTP Client Create Failure",
			utils.HttpDiagnostics(httpResp),
		)
		return
	}

	// Check for API success inlineResp code
	if httpResp.StatusCode != 201 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	payloadRespDiags := CreateHttpResponse(ctx, data, inlineResp)
	if payloadRespDiags != nil {
		resp.Diagnostics.Append(payloadRespDiags...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Update to capture config items not accessible in HTTP POST
	updatePayload, updatePayloadReqDiags := UpdateHttpReqPayload(ctx, data)
	if updatePayloadReqDiags != nil {
		resp.Diagnostics.Append(updatePayloadReqDiags...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// API returns this as string, openAPI spec has set as Integer
	vlanId := fmt.Sprintf("%v", data.VlanId.ValueInt64())

	updateInlineResp, updateHttpResp, updateErr := r.client.ApplianceApi.UpdateNetworkApplianceVlan(ctx, data.NetworkId.ValueString(), vlanId).UpdateNetworkApplianceVlanRequest(*updatePayload).Execute()
	if updateErr != nil && updateHttpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"HTTP Client Update Failure",
			utils.HttpDiagnostics(updateHttpResp),
		)
		return
	}

	// Check for API success response code
	if updateHttpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", updateHttpResp.StatusCode),
		)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
		return
	}

	updatePayloadRespDiags := ReadHttpResponse(ctx, data, updateInlineResp)
	if updatePayloadRespDiags != nil {
		resp.Diagnostics.Append(updatePayloadRespDiags...)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Log the response data
	tflog.Info(ctx, "[finish] CREATE Function Call")
	tflog.Trace(ctx, "Create function completed", map[string]interface{}{
		"data": data,
	})
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetworksApplianceVLANModel

	// Log the received request
	tflog.Info(ctx, "[start] READ Function Call")
	tflog.Trace(ctx, "Read Function Called", map[string]interface{}{
		"request": req,
	})

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// API returns string, OpenAPI defines integer
	vlanId := fmt.Sprintf("%v", data.VlanId.ValueInt64())

	inlineResp, httpResp, err := r.client.ApplianceApi.GetNetworkApplianceVlan(ctx, data.NetworkId.ValueString(), vlanId).Execute()
	if err != nil && httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"HTTP Client Read Failure",
			utils.HttpDiagnostics(httpResp),
		)
		return
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

	payloadRespDiags := ReadHttpResponse(ctx, data, inlineResp)
	if payloadRespDiags != nil {
		resp.Diagnostics.Append(payloadRespDiags...)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Log the response data
	tflog.Info(ctx, "[finish] READ Function Call")
	tflog.Trace(ctx, "Read Function", map[string]interface{}{
		"data": data,
	})
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *NetworksApplianceVLANModel

	// Log the received request
	tflog.Info(ctx, "[start] UPDATE Function Call")
	tflog.Trace(ctx, "Update Function Called", map[string]interface{}{
		"request": req,
	})

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	payload, payloadReqDiags := UpdateHttpReqPayload(ctx, data)
	if payloadReqDiags != nil {
		resp.Diagnostics.Append(payloadReqDiags...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// API returns this as string, openAPI spec has set as Integer
	vlanId := fmt.Sprintf("%v", data.VlanId.ValueInt64())

	inlineResp, httpResp, err := r.client.ApplianceApi.UpdateNetworkApplianceVlan(ctx, data.NetworkId.ValueString(), vlanId).UpdateNetworkApplianceVlanRequest(*payload).Execute()
	if err != nil && httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"HTTP Client Update Failure",
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

	payloadRespDiags := ReadHttpResponse(ctx, data, inlineResp)
	if payloadRespDiags != nil {
		resp.Diagnostics.Append(payloadRespDiags...)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Log the response data
	tflog.Info(ctx, "[finish] UPDATE Function Call")
	tflog.Trace(ctx, "Update Function", map[string]interface{}{
		"data": data,
	})
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *NetworksApplianceVLANModel

	// Log the received request
	tflog.Info(ctx, "[start] DELETE Function Call")
	tflog.Trace(ctx, "Delete Function Called", map[string]interface{}{
		"request": req,
	})

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// API returns a string, OpenAPI spec defines an integer
	vlanId := fmt.Sprintf("%v", data.VlanId.ValueInt64())

	httpResp, err := r.client.ApplianceApi.DeleteNetworkApplianceVlan(ctx, data.NetworkId.ValueString(), vlanId).Execute()
	if err != nil && httpResp.StatusCode != 204 {
		resp.Diagnostics.AddError(
			"HTTP Client Delete Failure",
			utils.HttpDiagnostics(httpResp),
		)
		return
	}

	// Check for API success response code
	if httpResp.StatusCode != 204 {
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

	// Log the response data
	tflog.Info(ctx, "[finish] DELETE Function Call")
	tflog.Trace(ctx, "Delete Function", map[string]interface{}{
		"data": data,
	})

}
