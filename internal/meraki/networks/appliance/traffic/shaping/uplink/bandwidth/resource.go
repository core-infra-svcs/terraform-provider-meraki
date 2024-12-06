package bandwidth

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &Resource{}
var _ resource.ResourceWithImportState = &Resource{}

func NewResource() resource.Resource {
	return &Resource{}
}

// Resource defines the resource implementation.
type Resource struct {
	client *openApiClient.APIClient
}

func (r *Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_appliance_traffic_shaping_uplink_bandwidth"
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

	updateApplianceTrafficShapingUplinkBandWidth := *openApiClient.NewUpdateNetworkApplianceTrafficShapingUplinkBandwidthRequest()

	var bandwidthLimit openApiClient.UpdateNetworkApplianceTrafficShapingUplinkBandwidthRequestBandwidthLimits

	var cellular openApiClient.UpdateNetworkApplianceTrafficShapingUplinkBandwidthRequestBandwidthLimitsCellular

	if !(data.BandwidthLimitCellularLimitUp.IsUnknown() || data.BandwidthLimitCellularLimitDown.IsUnknown()) {
		if !(data.BandwidthLimitCellularLimitUp.IsNull() || data.BandwidthLimitCellularLimitDown.IsNull()) {
			cellular.SetLimitUp(int32(data.BandwidthLimitCellularLimitUp.ValueInt64()))
			cellular.SetLimitDown(int32(data.BandwidthLimitCellularLimitDown.ValueInt64()))
			bandwidthLimit.SetCellular(cellular)
		}
	}

	var wan1 openApiClient.UpdateNetworkApplianceTrafficShapingUplinkBandwidthRequestBandwidthLimitsWan1
	if !(data.BandwidthLimitWan1LimitUp.IsUnknown() || data.BandwidthLimitWan1LimitDown.IsUnknown()) {
		if !(data.BandwidthLimitWan1LimitUp.IsNull() || data.BandwidthLimitWan1LimitDown.IsNull()) {
			wan1.SetLimitUp(int32(data.BandwidthLimitWan1LimitUp.ValueInt64()))
			wan1.SetLimitDown(int32(data.BandwidthLimitWan1LimitDown.ValueInt64()))
			bandwidthLimit.SetWan1(wan1)
		}
	}

	var wan2 openApiClient.UpdateNetworkApplianceTrafficShapingUplinkBandwidthRequestBandwidthLimitsWan2
	if !(data.BandwidthLimitWan2LimitUp.IsUnknown() || data.BandwidthLimitWan2LimitDown.IsUnknown()) {
		if !(data.BandwidthLimitWan2LimitUp.IsNull() || data.BandwidthLimitWan2LimitDown.IsNull()) {
			wan2.SetLimitUp(int32(data.BandwidthLimitWan2LimitUp.ValueInt64()))
			wan2.SetLimitDown(int32(data.BandwidthLimitWan2LimitDown.ValueInt64()))
			bandwidthLimit.SetWan2(wan2)
		}
	}

	updateApplianceTrafficShapingUplinkBandWidth.SetBandwidthLimits(bandwidthLimit)

	_, httpResp, err := r.client.ApplianceApi.UpdateNetworkApplianceTrafficShapingUplinkBandwidth(context.Background(), data.NetworkId.ValueString()).UpdateNetworkApplianceTrafficShapingUplinkBandwidthRequest(updateApplianceTrafficShapingUplinkBandWidth).Execute()

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
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
		return
	}

	// Save data into Terraform state
	data, err = extractHttpResponseUplinkBandwidthResource(ctx, httpResp.Body, &resourceModelApiResponse{}, data)
	if err != nil {
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

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *resourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	_, httpResp, err := r.client.ApplianceApi.GetNetworkApplianceTrafficShapingUplinkBandwidth(context.Background(), data.NetworkId.ValueString()).Execute()

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
		return
	} else {
		resp.Diagnostics.Append()
	}

	// Decode the HTTP response body into your data model.
	// If there's an error, add it to diagnostics.
	// Save data into Terraform state
	data, err = extractHttpResponseUplinkBandwidthResource(ctx, httpResp.Body, &resourceModelApiResponse{}, data)
	if err != nil {
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

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data *resourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateApplianceTrafficShapingUplinkBandWidth := *openApiClient.NewUpdateNetworkApplianceTrafficShapingUplinkBandwidthRequest()

	var bandwidthLimit openApiClient.UpdateNetworkApplianceTrafficShapingUplinkBandwidthRequestBandwidthLimits

	var cellular openApiClient.UpdateNetworkApplianceTrafficShapingUplinkBandwidthRequestBandwidthLimitsCellular

	if !(data.BandwidthLimitCellularLimitUp.IsUnknown() || data.BandwidthLimitCellularLimitDown.IsUnknown()) {
		if !(data.BandwidthLimitCellularLimitUp.IsNull() || data.BandwidthLimitCellularLimitDown.IsNull()) {
			cellular.SetLimitUp(int32(data.BandwidthLimitCellularLimitUp.ValueInt64()))
			cellular.SetLimitDown(int32(data.BandwidthLimitCellularLimitDown.ValueInt64()))
			bandwidthLimit.SetCellular(cellular)
		}
	}

	var wan1 openApiClient.UpdateNetworkApplianceTrafficShapingUplinkBandwidthRequestBandwidthLimitsWan1
	if !(data.BandwidthLimitWan1LimitUp.IsUnknown() || data.BandwidthLimitWan1LimitDown.IsUnknown()) {
		if !(data.BandwidthLimitWan1LimitUp.IsNull() || data.BandwidthLimitWan1LimitDown.IsNull()) {
			wan1.SetLimitUp(int32(data.BandwidthLimitWan1LimitUp.ValueInt64()))
			wan1.SetLimitDown(int32(data.BandwidthLimitWan1LimitDown.ValueInt64()))
			bandwidthLimit.SetWan1(wan1)
		}
	}

	var wan2 openApiClient.UpdateNetworkApplianceTrafficShapingUplinkBandwidthRequestBandwidthLimitsWan2
	if !(data.BandwidthLimitWan2LimitUp.IsUnknown() || data.BandwidthLimitWan2LimitDown.IsUnknown()) {
		if !(data.BandwidthLimitWan2LimitUp.IsNull() || data.BandwidthLimitWan2LimitDown.IsNull()) {
			wan2.SetLimitUp(int32(data.BandwidthLimitWan2LimitUp.ValueInt64()))
			wan2.SetLimitDown(int32(data.BandwidthLimitWan2LimitDown.ValueInt64()))
			bandwidthLimit.SetWan2(wan2)
		}
	}

	updateApplianceTrafficShapingUplinkBandWidth.SetBandwidthLimits(bandwidthLimit)

	_, httpResp, err := r.client.ApplianceApi.UpdateNetworkApplianceTrafficShapingUplinkBandwidth(context.Background(), data.NetworkId.ValueString()).UpdateNetworkApplianceTrafficShapingUplinkBandwidthRequest(updateApplianceTrafficShapingUplinkBandWidth).Execute()

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
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
		return
	}

	// Decode the HTTP response body into your data model.
	// If there's an error, add it to diagnostics.
	// Save data into Terraform state
	data, err = extractHttpResponseUplinkBandwidthResource(ctx, httpResp.Body, &resourceModelApiResponse{}, data)
	if err != nil {
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

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data *resourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateApplianceTrafficShapingUplinkBandWidth := *openApiClient.NewUpdateNetworkApplianceTrafficShapingUplinkBandwidthRequest()

	var bandwidthLimit openApiClient.UpdateNetworkApplianceTrafficShapingUplinkBandwidthRequestBandwidthLimits
	var cellular openApiClient.UpdateNetworkApplianceTrafficShapingUplinkBandwidthRequestBandwidthLimitsCellular
	var wan1 openApiClient.UpdateNetworkApplianceTrafficShapingUplinkBandwidthRequestBandwidthLimitsWan1
	var wan2 openApiClient.UpdateNetworkApplianceTrafficShapingUplinkBandwidthRequestBandwidthLimitsWan2

	bandwidthLimit.SetCellular(cellular)
	bandwidthLimit.SetWan1(wan1)
	bandwidthLimit.SetWan2(wan2)
	updateApplianceTrafficShapingUplinkBandWidth.SetBandwidthLimits(bandwidthLimit)

	_, httpResp, err := r.client.ApplianceApi.UpdateNetworkApplianceTrafficShapingUplinkBandwidth(context.Background(), data.NetworkId.ValueString()).UpdateNetworkApplianceTrafficShapingUplinkBandwidthRequest(updateApplianceTrafficShapingUplinkBandWidth).Execute()

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
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

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
