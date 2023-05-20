package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
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
var _ resource.Resource = &NetworksApplianceTrafficShappingUplinkBandWidthResource{}
var _ resource.ResourceWithImportState = &NetworksApplianceTrafficShappingUplinkBandWidthResource{}

func NewNetworksApplianceTrafficShappingUplinkBandWidthResource() resource.Resource {
	return &NetworksApplianceTrafficShappingUplinkBandWidthResource{}
}

// NetworksApplianceTrafficShappingUplinkBandWidthResource defines the resource implementation.
type NetworksApplianceTrafficShappingUplinkBandWidthResource struct {
	client *openApiClient.APIClient
}

// NetworksApplianceTrafficShappingUplinkBandWidthResourceModel describes the resource data model.
type NetworksApplianceTrafficShappingUplinkBandWidthResourceModel struct {
	Id              jsontypes.String `tfsdk:"id"`
	NetworkId       jsontypes.String `tfsdk:"network_id" json:"network_id"`
	BandwidthLimits BandwidthLimits  `tfsdk:"bandwidth_limits" json:"bandwidthLimits"`
}

type BandwidthLimits struct {
	Wan1     Limits `tfsdk:"wan1" json:"wan1"`
	Wan2     Limits `tfsdk:"wan2" json:"wan2"`
	Cellular Limits `tfsdk:"cellular" json:"cellular"`
}

type Limits struct {
	LimitUp   jsontypes.Int64 `tfsdk:"limit_up" json:"limitUp"`
	LimitDown jsontypes.Int64 `tfsdk:"limit_down" json:"limitDown"`
}

func (r *NetworksApplianceTrafficShappingUplinkBandWidthResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_appliance_traffic_shapping_uplink_bandWidth"
}

func (r *NetworksApplianceTrafficShappingUplinkBandWidthResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "NetworksApplianceTrafficShappingUplinkBandWidth resource for updating Network Appliance Traffic Shapping UplinkBandWidth.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.StringType,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Network ID",
				Required:            true,
				CustomType:          jsontypes.StringType,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 31),
				},
			},
			"bandwidth_limits": schema.SingleNestedAttribute{
				MarkdownDescription: "A mapping of uplinks to their bandwidth settings (be sure to check which uplinks are supported for your network)",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"cellular": schema.SingleNestedAttribute{
						MarkdownDescription: "The bandwidth settings for the 'cellular' uplink",
						Required:            true,
						Attributes: map[string]schema.Attribute{
							"limit_up": schema.Int64Attribute{
								MarkdownDescription: "The maximum upload limit (integer, in Kbps). null indicates no limit",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.Int64Type,
							},
							"limit_down": schema.Int64Attribute{
								MarkdownDescription: "The maximum download limit (integer, in Kbps). null indicates no limit",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.Int64Type,
							},
						},
					},
					"wan2": schema.SingleNestedAttribute{
						MarkdownDescription: "The bandwidth settings for the 'wan2' uplink",
						Required:            true,
						Attributes: map[string]schema.Attribute{
							"limit_up": schema.Int64Attribute{
								MarkdownDescription: "The maximum upload limit (integer, in Kbps). null indicates no limit",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.Int64Type,
							},
							"limit_down": schema.Int64Attribute{
								MarkdownDescription: "The maximum download limit (integer, in Kbps). null indicates no limit",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.Int64Type,
							},
						},
					},
					"wan1": schema.SingleNestedAttribute{
						MarkdownDescription: "The bandwidth settings for the 'wan1' uplink",
						Required:            true,
						Attributes: map[string]schema.Attribute{
							"limit_up": schema.Int64Attribute{
								MarkdownDescription: "The maximum upload limit (integer, in Kbps). null indicates no limit",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.Int64Type,
							},
							"limit_down": schema.Int64Attribute{
								MarkdownDescription: "The maximum download limit (integer, in Kbps). null indicates no limit",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.Int64Type,
							},
						},
					},
				},
			},
		},
	}
}

func (r *NetworksApplianceTrafficShappingUplinkBandWidthResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NetworksApplianceTrafficShappingUplinkBandWidthResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	var data *NetworksApplianceTrafficShappingUplinkBandWidthResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateApplianceTrafficShappingUplinkBandWidth := *openApiClient.NewInlineObject54()

	var bandwidthLimit openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkBandwidthBandwidthLimits

	var cellular openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkBandwidthBandwidthLimitsCellular
	if !(data.BandwidthLimits.Cellular.LimitUp.IsUnknown() || data.BandwidthLimits.Cellular.LimitDown.IsUnknown()) {
		if !(data.BandwidthLimits.Cellular.LimitUp.IsNull() || data.BandwidthLimits.Cellular.LimitUp.IsNull()) {
			cellular.SetLimitUp(int32(data.BandwidthLimits.Cellular.LimitUp.ValueInt64()))
			cellular.SetLimitDown(int32(data.BandwidthLimits.Cellular.LimitDown.ValueInt64()))
			bandwidthLimit.SetCellular(cellular)
		}
	}

	var wan1 openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkBandwidthBandwidthLimitsWan1
	if !(data.BandwidthLimits.Wan1.LimitUp.IsUnknown() || data.BandwidthLimits.Wan1.LimitUp.IsUnknown()) {
		if !(data.BandwidthLimits.Wan1.LimitUp.IsNull() || data.BandwidthLimits.Wan1.LimitUp.IsNull()) {
			wan1.SetLimitUp(int32(data.BandwidthLimits.Wan1.LimitUp.ValueInt64()))
			wan1.SetLimitDown(int32(data.BandwidthLimits.Wan1.LimitDown.ValueInt64()))
			bandwidthLimit.SetWan1(wan1)
		}
	}

	var wan2 openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkBandwidthBandwidthLimitsWan2
	if !(data.BandwidthLimits.Wan2.LimitUp.IsUnknown() || data.BandwidthLimits.Wan2.LimitUp.IsUnknown()) {
		if !(data.BandwidthLimits.Wan2.LimitUp.IsNull() || data.BandwidthLimits.Wan2.LimitUp.IsNull()) {
			wan2.SetLimitUp(int32(data.BandwidthLimits.Wan2.LimitUp.ValueInt64()))
			wan2.SetLimitDown(int32(data.BandwidthLimits.Wan2.LimitDown.ValueInt64()))
			bandwidthLimit.SetWan2(wan2)
		}
	}

	updateApplianceTrafficShappingUplinkBandWidth.SetBandwidthLimits(bandwidthLimit)

	inlineResp, httpResp, err := r.client.ApplianceApi.UpdateNetworkApplianceTrafficShapingUplinkBandwidth(context.Background(), data.NetworkId.ValueString()).UpdateNetworkApplianceTrafficShapingUplinkBandwidth(updateApplianceTrafficShappingUplinkBandWidth).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// collect diagnostics
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
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

	data = extractHttpResponseGroupPolicyResource(ctx, inlineResp, data)

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *NetworksApplianceTrafficShappingUplinkBandWidthResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetworksApplianceTrafficShappingUplinkBandWidthResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	inlineResp, httpResp, err := r.client.ApplianceApi.GetNetworkApplianceTrafficShapingUplinkBandwidth(context.Background(), data.NetworkId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// collect diagnostics
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
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

	data = extractHttpResponseGroupPolicyResource(ctx, inlineResp, data)

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *NetworksApplianceTrafficShappingUplinkBandWidthResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data *NetworksApplianceTrafficShappingUplinkBandWidthResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateApplianceTrafficShappingUplinkBandWidth := *openApiClient.NewInlineObject54()

	var bandwidthLimit openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkBandwidthBandwidthLimits

	var cellular openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkBandwidthBandwidthLimitsCellular
	if !(data.BandwidthLimits.Cellular.LimitUp.IsUnknown() || data.BandwidthLimits.Cellular.LimitDown.IsUnknown()) {
		if !(data.BandwidthLimits.Cellular.LimitUp.IsNull() || data.BandwidthLimits.Cellular.LimitUp.IsNull()) {
			cellular.SetLimitUp(int32(data.BandwidthLimits.Cellular.LimitUp.ValueInt64()))
			cellular.SetLimitDown(int32(data.BandwidthLimits.Cellular.LimitDown.ValueInt64()))
			bandwidthLimit.SetCellular(cellular)
		}
	}

	var wan1 openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkBandwidthBandwidthLimitsWan1
	if !(data.BandwidthLimits.Wan1.LimitUp.IsUnknown() || data.BandwidthLimits.Wan1.LimitUp.IsUnknown()) {
		if !(data.BandwidthLimits.Wan1.LimitUp.IsNull() || data.BandwidthLimits.Wan1.LimitUp.IsNull()) {
			wan1.SetLimitUp(int32(data.BandwidthLimits.Wan1.LimitUp.ValueInt64()))
			wan1.SetLimitDown(int32(data.BandwidthLimits.Wan1.LimitDown.ValueInt64()))
			bandwidthLimit.SetWan1(wan1)
		}
	}

	var wan2 openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkBandwidthBandwidthLimitsWan2
	if !(data.BandwidthLimits.Wan2.LimitUp.IsUnknown() || data.BandwidthLimits.Wan2.LimitUp.IsUnknown()) {
		if !(data.BandwidthLimits.Wan2.LimitUp.IsNull() || data.BandwidthLimits.Wan2.LimitUp.IsNull()) {
			wan2.SetLimitUp(int32(data.BandwidthLimits.Wan2.LimitUp.ValueInt64()))
			wan2.SetLimitDown(int32(data.BandwidthLimits.Wan2.LimitDown.ValueInt64()))
			bandwidthLimit.SetWan2(wan2)
		}
	}

	updateApplianceTrafficShappingUplinkBandWidth.SetBandwidthLimits(bandwidthLimit)

	inlineResp, httpResp, err := r.client.ApplianceApi.UpdateNetworkApplianceTrafficShapingUplinkBandwidth(context.Background(), data.NetworkId.ValueString()).UpdateNetworkApplianceTrafficShapingUplinkBandwidth(updateApplianceTrafficShappingUplinkBandWidth).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// collect diagnostics
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
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

	data = extractHttpResponseGroupPolicyResource(ctx, inlineResp, data)

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *NetworksApplianceTrafficShappingUplinkBandWidthResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data *NetworksApplianceTrafficShappingUplinkBandWidthResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateApplianceTrafficShappingUplinkBandWidth := *openApiClient.NewInlineObject54()

	var bandwidthLimit openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkBandwidthBandwidthLimits

	var cellular openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkBandwidthBandwidthLimitsCellular
	if !(data.BandwidthLimits.Cellular.LimitUp.IsUnknown() || data.BandwidthLimits.Cellular.LimitDown.IsUnknown()) {
		if !(data.BandwidthLimits.Cellular.LimitUp.IsNull() || data.BandwidthLimits.Cellular.LimitUp.IsNull()) {
			cellular.SetLimitUp(int32(data.BandwidthLimits.Cellular.LimitUp.ValueInt64()))
			cellular.SetLimitDown(int32(data.BandwidthLimits.Cellular.LimitDown.ValueInt64()))
			bandwidthLimit.SetCellular(cellular)
		}
	}

	var wan1 openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkBandwidthBandwidthLimitsWan1
	if !(data.BandwidthLimits.Wan1.LimitUp.IsUnknown() || data.BandwidthLimits.Wan1.LimitUp.IsUnknown()) {
		if !(data.BandwidthLimits.Wan1.LimitUp.IsNull() || data.BandwidthLimits.Wan1.LimitUp.IsNull()) {
			wan1.SetLimitUp(int32(data.BandwidthLimits.Wan1.LimitUp.ValueInt64()))
			wan1.SetLimitDown(int32(data.BandwidthLimits.Wan1.LimitDown.ValueInt64()))
			bandwidthLimit.SetWan1(wan1)
		}
	}

	var wan2 openApiClient.NetworksNetworkIdApplianceTrafficShapingUplinkBandwidthBandwidthLimitsWan2
	if !(data.BandwidthLimits.Wan2.LimitUp.IsUnknown() || data.BandwidthLimits.Wan2.LimitUp.IsUnknown()) {
		if !(data.BandwidthLimits.Wan2.LimitUp.IsNull() || data.BandwidthLimits.Wan2.LimitUp.IsNull()) {
			wan2.SetLimitUp(int32(data.BandwidthLimits.Wan2.LimitUp.ValueInt64()))
			wan2.SetLimitDown(int32(data.BandwidthLimits.Wan2.LimitDown.ValueInt64()))
			bandwidthLimit.SetWan2(wan2)
		}
	}

	updateApplianceTrafficShappingUplinkBandWidth.SetBandwidthLimits(bandwidthLimit)

	inlineResp, httpResp, err := r.client.ApplianceApi.UpdateNetworkApplianceTrafficShapingUplinkBandwidth(context.Background(), data.NetworkId.ValueString()).UpdateNetworkApplianceTrafficShapingUplinkBandwidth(updateApplianceTrafficShappingUplinkBandWidth).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// collect diagnostics
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
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

	data = extractHttpResponseGroupPolicyResource(ctx, inlineResp, data)

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	resp.State.RemoveResource(ctx)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")
}

func (r *NetworksApplianceTrafficShappingUplinkBandWidthResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), req.ID)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func extractHttpResponseGroupPolicyResource(ctx context.Context, inlineResp map[string]interface{}, data *NetworksApplianceTrafficShappingUplinkBandWidthResourceModel) *NetworksApplianceTrafficShappingUplinkBandWidthResourceModel {

	if bandwidthLimits := inlineResp["bandwidthLimits"]; bandwidthLimits != nil {
		var bandwidthLimitsData BandwidthLimits
		jsonData, _ := json.Marshal(bandwidthLimits)
		json.Unmarshal(jsonData, &bandwidthLimitsData)
		if !bandwidthLimitsData.Wan1.LimitDown.IsUnknown() {
			data.BandwidthLimits.Wan1.LimitDown = bandwidthLimitsData.Wan1.LimitDown
		} else {
			data.BandwidthLimits.Wan1.LimitDown = jsontypes.Int64Null()
		}
		if !bandwidthLimitsData.Wan1.LimitUp.IsUnknown() {
			data.BandwidthLimits.Wan1.LimitUp = bandwidthLimitsData.Wan1.LimitUp
		} else {
			data.BandwidthLimits.Wan1.LimitUp = jsontypes.Int64Null()
		}
		if !bandwidthLimitsData.Wan2.LimitDown.IsUnknown() {
			data.BandwidthLimits.Wan2.LimitDown = bandwidthLimitsData.Wan2.LimitDown
		} else {
			data.BandwidthLimits.Wan2.LimitDown = jsontypes.Int64Null()
		}
		if !bandwidthLimitsData.Wan2.LimitUp.IsUnknown() {
			data.BandwidthLimits.Wan2.LimitUp = bandwidthLimitsData.Wan2.LimitUp
		} else {
			data.BandwidthLimits.Wan2.LimitUp = jsontypes.Int64Null()
		}
		if !bandwidthLimitsData.Cellular.LimitDown.IsUnknown() {
			data.BandwidthLimits.Cellular.LimitDown = bandwidthLimitsData.Cellular.LimitDown
		} else {
			data.BandwidthLimits.Cellular.LimitDown = jsontypes.Int64Null()
		}
		if !bandwidthLimitsData.Cellular.LimitUp.IsUnknown() {
			data.BandwidthLimits.Cellular.LimitUp = bandwidthLimitsData.Cellular.LimitUp
		} else {
			data.BandwidthLimits.Cellular.LimitUp = jsontypes.Int64Null()
		}

	} else {
		data.BandwidthLimits.Wan1.LimitUp = jsontypes.Int64Null()
		data.BandwidthLimits.Wan1.LimitDown = jsontypes.Int64Null()
		data.BandwidthLimits.Wan1.LimitUp = jsontypes.Int64Null()
		data.BandwidthLimits.Wan1.LimitDown = jsontypes.Int64Null()
		data.BandwidthLimits.Wan1.LimitUp = jsontypes.Int64Null()
		data.BandwidthLimits.Wan1.LimitDown = jsontypes.Int64Null()
	}
	return data
}
