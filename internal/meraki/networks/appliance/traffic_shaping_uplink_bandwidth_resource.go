package appliance

import (
	"context"
	"encoding/json"
	"fmt"
	jsontypes2 "github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"io"

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
var _ resource.Resource = &NetworksApplianceTrafficShapingUplinkBandWidthResource{}
var _ resource.ResourceWithImportState = &NetworksApplianceTrafficShapingUplinkBandWidthResource{}

func NewNetworksApplianceTrafficShapingUplinkBandWidthResource() resource.Resource {
	return &NetworksApplianceTrafficShapingUplinkBandWidthResource{}
}

// NetworksApplianceTrafficShapingUplinkBandWidthResource defines the resource implementation.
type NetworksApplianceTrafficShapingUplinkBandWidthResource struct {
	client *openApiClient.APIClient
}

// NetworksApplianceTrafficShapingUplinkBandWidthResourceModel describes the resource data model.
type NetworksApplianceTrafficShapingUplinkBandWidthResourceModel struct {
	Id                              jsontypes2.String `tfsdk:"id"`
	NetworkId                       jsontypes2.String `tfsdk:"network_id" json:"network_id"`
	BandwidthLimitCellularLimitUp   jsontypes2.Int64  `tfsdk:"bandwidth_limit_cellular_limit_up"`
	BandwidthLimitCellularLimitDown jsontypes2.Int64  `tfsdk:"bandwidth_limit_cellular_limit_down"`
	BandwidthLimitWan2LimitUp       jsontypes2.Int64  `tfsdk:"bandwidth_limit_wan2_limit_up"`
	BandwidthLimitWan2LimitDown     jsontypes2.Int64  `tfsdk:"bandwidth_limit_wan2_limit_down"`
	BandwidthLimitWan1LimitUp       jsontypes2.Int64  `tfsdk:"bandwidth_limit_wan1_limit_up"`
	BandwidthLimitWan1LimitDown     jsontypes2.Int64  `tfsdk:"bandwidth_limit_wan1_limit_down"`
}

type NetworksApplianceTrafficShapingUplinkBandWidthResourceModelApiResponse struct {
	UplinkBandwidthLimits NetworksApplianceTrafficShapingUplinkBandWidthResourceModelUplinkBandwidthLimits `json:"bandwidthLimits"`
}

type NetworksApplianceTrafficShapingUplinkBandWidthResourceModelUplinkBandwidthLimits struct {
	Wan1     NetworksApplianceTrafficShapingUplinkBandWidthResourceModelLimits `json:"wan1"`
	Wan2     NetworksApplianceTrafficShapingUplinkBandWidthResourceModelLimits `json:"wan2"`
	Cellular NetworksApplianceTrafficShapingUplinkBandWidthResourceModelLimits `json:"cellular"`
}

type NetworksApplianceTrafficShapingUplinkBandWidthResourceModelLimits struct {
	LimitUp   jsontypes2.Int64 `json:"limitUp"`
	LimitDown jsontypes2.Int64 `json:"limitDown"`
}

func (r *NetworksApplianceTrafficShapingUplinkBandWidthResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_appliance_traffic_shaping_uplink_bandwidth"
}

func (r *NetworksApplianceTrafficShapingUplinkBandWidthResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage Network Appliance Traffic Shaping UplinkBandWidth",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes2.StringType,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Network ID",
				Required:            true,
				CustomType:          jsontypes2.StringType,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 31),
				},
			},
			"bandwidth_limit_cellular_limit_up": schema.Int64Attribute{
				MarkdownDescription: "The bandwidth settings for the 'cellular' uplink. The maximum upload limit (integer, in Kbps). null indicates no limit",
				Optional:            true,
				CustomType:          jsontypes2.Int64Type,
			},
			"bandwidth_limit_cellular_limit_down": schema.Int64Attribute{
				MarkdownDescription: "The bandwidth settings for the 'cellular' uplink. The maximum download limit (integer, in Kbps). null indicates no limit",
				Optional:            true,
				CustomType:          jsontypes2.Int64Type,
			},
			"bandwidth_limit_wan2_limit_up": schema.Int64Attribute{
				MarkdownDescription: "The bandwidth settings for the 'wan2' uplink. The maximum upload limit (integer, in Kbps). null indicates no limit",
				Optional:            true,
				CustomType:          jsontypes2.Int64Type,
			},
			"bandwidth_limit_wan2_limit_down": schema.Int64Attribute{
				MarkdownDescription: "The bandwidth settings for the 'wan2' uplink. The maximum download limit (integer, in Kbps). null indicates no limit",
				Optional:            true,
				CustomType:          jsontypes2.Int64Type,
			},
			"bandwidth_limit_wan1_limit_up": schema.Int64Attribute{
				MarkdownDescription: "The bandwidth settings for the 'wan1' uplink. The maximum upload limit (integer, in Kbps). null indicates no limit",
				Optional:            true,
				CustomType:          jsontypes2.Int64Type,
			},
			"bandwidth_limit_wan1_limit_down": schema.Int64Attribute{
				MarkdownDescription: "The bandwidth settings for the 'wan1' uplink. The maximum download limit (integer, in Kbps). null indicates no limit",
				Optional:            true,
				CustomType:          jsontypes2.Int64Type,
			},
		},
	}
}

func (r *NetworksApplianceTrafficShapingUplinkBandWidthResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NetworksApplianceTrafficShapingUplinkBandWidthResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data *NetworksApplianceTrafficShapingUplinkBandWidthResourceModel

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
	data, err = extractHttpResponseUplinkBandwidthResource(ctx, httpResp.Body, &NetworksApplianceTrafficShapingUplinkBandWidthResourceModelApiResponse{}, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	data.Id = jsontypes2.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")
}

func (r *NetworksApplianceTrafficShapingUplinkBandWidthResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetworksApplianceTrafficShapingUplinkBandWidthResourceModel

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
	data, err = extractHttpResponseUplinkBandwidthResource(ctx, httpResp.Body, &NetworksApplianceTrafficShapingUplinkBandWidthResourceModelApiResponse{}, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	data.Id = jsontypes2.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *NetworksApplianceTrafficShapingUplinkBandWidthResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data *NetworksApplianceTrafficShapingUplinkBandWidthResourceModel

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
	data, err = extractHttpResponseUplinkBandwidthResource(ctx, httpResp.Body, &NetworksApplianceTrafficShapingUplinkBandWidthResourceModelApiResponse{}, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	data.Id = jsontypes2.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *NetworksApplianceTrafficShapingUplinkBandWidthResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data *NetworksApplianceTrafficShapingUplinkBandWidthResourceModel

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

func (r *NetworksApplianceTrafficShapingUplinkBandWidthResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), req.ID)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func extractHttpResponseUplinkBandwidthResource(ctx context.Context, httpRespBody io.ReadCloser, apiResponse *NetworksApplianceTrafficShapingUplinkBandWidthResourceModelApiResponse, data *NetworksApplianceTrafficShapingUplinkBandWidthResourceModel) (*NetworksApplianceTrafficShapingUplinkBandWidthResourceModel, error) {

	if err := json.NewDecoder(httpRespBody).Decode(apiResponse); err != nil {
		return data, err
	}

	data.BandwidthLimitCellularLimitDown = apiResponse.UplinkBandwidthLimits.Cellular.LimitDown
	data.BandwidthLimitCellularLimitUp = apiResponse.UplinkBandwidthLimits.Cellular.LimitUp
	data.BandwidthLimitWan2LimitDown = apiResponse.UplinkBandwidthLimits.Wan2.LimitDown
	data.BandwidthLimitWan2LimitUp = apiResponse.UplinkBandwidthLimits.Wan2.LimitUp
	data.BandwidthLimitWan1LimitDown = apiResponse.UplinkBandwidthLimits.Wan1.LimitDown
	data.BandwidthLimitWan1LimitUp = apiResponse.UplinkBandwidthLimits.Wan1.LimitUp

	if data.BandwidthLimitWan1LimitDown.IsUnknown() {
		data.BandwidthLimitWan1LimitDown = jsontypes2.Int64Null()
	}
	if data.BandwidthLimitWan1LimitUp.IsUnknown() {
		data.BandwidthLimitWan1LimitUp = jsontypes2.Int64Null()
	}
	if data.BandwidthLimitWan2LimitDown.Int64Value.IsUnknown() {
		data.BandwidthLimitWan2LimitDown = jsontypes2.Int64Null()
	}
	if data.BandwidthLimitWan2LimitUp.IsUnknown() {
		data.BandwidthLimitWan2LimitUp = jsontypes2.Int64Null()
	}
	if data.BandwidthLimitCellularLimitDown.IsUnknown() {
		data.BandwidthLimitCellularLimitDown = jsontypes2.Int64Null()
	}
	if data.BandwidthLimitCellularLimitUp.IsUnknown() {
		data.BandwidthLimitCellularLimitUp = jsontypes2.Int64Null()
	}

	return data, nil
}
