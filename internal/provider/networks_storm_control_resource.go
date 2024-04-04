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
var _ resource.Resource = &NetworksStormControlResource{}
var _ resource.ResourceWithImportState = &NetworksStormControlResource{}

func NewNetworksStormControlResource() resource.Resource {
	return &NetworksStormControlResource{}
}

// NetworksStormControlResource defines the resource implementation.
type NetworksStormControlResource struct {
	client *openApiClient.APIClient
}

// NetworksStormControlResourceModel describes the resource data model.
type NetworksStormControlResourceModel struct {
	Id                      jsontypes.String `tfsdk:"id" json:"-"`
	NetworkId               jsontypes.String `tfsdk:"network_id" json:"network_id"`
	BroadcastThreshold      jsontypes.Int64  `tfsdk:"broadcast_threshold" json:"broadcastThreshold"`
	MulticastThreshold      jsontypes.Int64  `tfsdk:"multicast_threshold" json:"multicastThreshold"`
	UnknownUnicastThreshold jsontypes.Int64  `tfsdk:"unknown_unicast_threshold" json:"unknownUnicastThreshold"`
}

func (r *NetworksStormControlResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_storm_control"
}

func (r *NetworksStormControlResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "NetworksSwitchQosRule resource for updating network switch qos rule.",
		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{
				MarkdownDescription: "Network ID",
				Computed:            true,
				Optional:            true,
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
			"broadcast_threshold": schema.Int64Attribute{
				MarkdownDescription: "Broadcast Threshold",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.Int64Type,
			},
			"multicast_threshold": schema.Int64Attribute{
				MarkdownDescription: "Multicast Threshold",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.Int64Type,
			},
			"unknown_unicast_threshold": schema.Int64Attribute{
				MarkdownDescription: "Unknown Unicast Threshold",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.Int64Type,
			},
		},
	}
}

func (r *NetworksStormControlResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NetworksStormControlResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *NetworksStormControlResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworkSwitchStormControl := *openApiClient.NewUpdateNetworkSwitchStormControlRequest()

	// Set BroadcastThreshold
	if !data.BroadcastThreshold.IsUnknown() {
		updateNetworkSwitchStormControl.SetBroadcastThreshold(int32(data.BroadcastThreshold.ValueInt64()))
	}

	// Set MulticastThreshold
	if !data.MulticastThreshold.IsUnknown() {
		updateNetworkSwitchStormControl.SetMulticastThreshold(int32(data.MulticastThreshold.ValueInt64()))
	}

	// Set UnknownUnicastThreshold
	if !data.UnknownUnicastThreshold.IsUnknown() {
		updateNetworkSwitchStormControl.SetUnknownUnicastThreshold(int32(data.UnknownUnicastThreshold.ValueInt64()))
	}

	_, httpResp, err := r.client.ConfigureApi.UpdateNetworkSwitchStormControl(ctx, data.NetworkId.ValueString()).UpdateNetworkSwitchStormControlRequest(updateNetworkSwitchStormControl).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
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

	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON Decode issue",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
		return
	}

	data.Id = data.NetworkId

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")
}

func (r *NetworksStormControlResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetworksStormControlResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := r.client.ConfigureApi.GetNetworkSwitchStormControl(ctx, data.NetworkId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
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

	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON Decode issue",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
		return
	}

	data.Id = data.NetworkId

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *NetworksStormControlResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data *NetworksStormControlResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworkSwitchStormControl := *openApiClient.NewUpdateNetworkSwitchStormControlRequest()

	// Set BroadcastThreshold
	if !data.BroadcastThreshold.IsUnknown() {
		updateNetworkSwitchStormControl.SetBroadcastThreshold(int32(data.BroadcastThreshold.ValueInt64()))
	}

	// Set MulticastThreshold
	if !data.MulticastThreshold.IsUnknown() {
		updateNetworkSwitchStormControl.SetMulticastThreshold(int32(data.MulticastThreshold.ValueInt64()))
	}

	// Set UnknownUnicastThreshold
	if !data.UnknownUnicastThreshold.IsUnknown() {
		updateNetworkSwitchStormControl.SetUnknownUnicastThreshold(int32(data.UnknownUnicastThreshold.ValueInt64()))
	}

	_, httpResp, err := r.client.ConfigureApi.UpdateNetworkSwitchStormControl(ctx, data.NetworkId.ValueString()).UpdateNetworkSwitchStormControlRequest(updateNetworkSwitchStormControl).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
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

	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON Decode issue",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
		return
	}

	data.Id = data.NetworkId

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *NetworksStormControlResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *NetworksStormControlResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Initialize storm control update request
	updateNetworkSwitchStormControl := *openApiClient.NewUpdateNetworkSwitchStormControlRequest()

	// Set BroadcastThreshold
	if !data.BroadcastThreshold.IsUnknown() {
		updateNetworkSwitchStormControl.SetBroadcastThreshold(int32(100))
	}

	// Set MulticastThreshold
	if !data.MulticastThreshold.IsUnknown() {
		updateNetworkSwitchStormControl.SetMulticastThreshold(int32(100))
	}

	// Set UnknownUnicastThreshold
	if !data.UnknownUnicastThreshold.IsUnknown() {
		updateNetworkSwitchStormControl.SetUnknownUnicastThreshold(int32(100))
	}

	// Set BroadcastThreshold
	if data.BroadcastThreshold.IsUnknown() {
		updateNetworkSwitchStormControl.SetBroadcastThreshold(int32(100))
	}

	// Set MulticastThreshold
	if data.MulticastThreshold.IsUnknown() {
		updateNetworkSwitchStormControl.SetMulticastThreshold(int32(100))
	}

	// Set UnknownUnicastThreshold
	if data.UnknownUnicastThreshold.IsUnknown() {
		updateNetworkSwitchStormControl.SetUnknownUnicastThreshold(int32(100))
	}

	// Execute the storm control update request to clear storm control configuration
	_, httpResp, err := r.client.ConfigureApi.UpdateNetworkSwitchStormControl(ctx, data.NetworkId.ValueString()).UpdateNetworkSwitchStormControlRequest(updateNetworkSwitchStormControl).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
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

	resp.State.RemoveResource(ctx)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")
}

func (r *NetworksStormControlResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	// Set the imported network ID
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), req.ID)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
