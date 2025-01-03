package uplink

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
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

type resourceModel struct {
	Id                             jsontypes.String                            `tfsdk:"id"`
	NetworkId                      jsontypes.String                            `tfsdk:"network_id"`
	CellularGatewayBandwidthLimits resourceModelCellularGatewayBandwidthLimits `tfsdk:"bandwidth_limits"`
}

type resourceModelCellularGatewayBandwidthLimits struct {
	LimitUp   jsontypes.Int64 `tfsdk:"limit_up"`
	LimitDown jsontypes.Int64 `tfsdk:"limit_down"`
}

// Metadata provides a way to define information about the resource.
// This method is called by the framework to retrieve metadata about the resource.
func (r *Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_networks_cellular_gateway_uplink"
}

// Schema provides a way to define the structure of the resource data.
// It is called by the framework to get the schema of the resource.
func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	// The Schema object defines the structure of the resource.
	resp.Schema = schema.Schema{

		MarkdownDescription: "Manage the uplink settings for your MG network.",

		// The Attributes map describes the fields of the resource.
		Attributes: map[string]schema.Attribute{

			// Every resource must have an ID attribute. This is computed by the framework.
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.StringType,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Network Id",
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
				MarkdownDescription: "The bandwidth settings for your MG network",
				Optional:            true,
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"limit_down": schema.Int64Attribute{
						MarkdownDescription: "The maximum download limit (integer, in Kbps). null indicates no limit.",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.Int64Type,
					},
					"limit_up": schema.Int64Attribute{
						MarkdownDescription: "The maximum upload limit (integer, in Kbps). null indicates no limit.",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.Int64Type,
					},
				},
			},
		},
	}
}

// Configure is a method of the Resource interface that Terraform calls to provide the configured provider instance to the resource.
// It passes the ResourceData that's been stored by the provider's ConfigureFunc.
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
// It takes a CreateRequest containing the planned state of the new resource and returns a CreateResponse
// with the final state of the new resource or an error.
func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *resourceModel

	// Unmarshal the plan data into the internal data model struct.
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// Check if there are any errors before proceeding.
	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworkCellularGatewayUplink := *openApiClient.NewUpdateNetworkCellularGatewayUplinkRequest()

	var bandwidthLimits openApiClient.UpdateNetworkApplianceTrafficShapingUplinkBandwidthRequestBandwidthLimitsCellular
	bandwidthLimits.SetLimitUp(int32(data.CellularGatewayBandwidthLimits.LimitUp.ValueInt64()))
	bandwidthLimits.SetLimitDown(int32(data.CellularGatewayBandwidthLimits.LimitDown.ValueInt64()))
	updateNetworkCellularGatewayUplink.SetBandwidthLimits(bandwidthLimits)

	_, httpResp, err := r.client.CellularGatewayApi.UpdateNetworkCellularGatewayUplink(context.Background(), data.NetworkId.ValueString()).UpdateNetworkCellularGatewayUplinkRequest(updateNetworkCellularGatewayUplink).Execute()

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

	// Decode the HTTP response body into your data model.
	// If there's an error, add it to diagnostics.
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	// Set ID for the new resource.
	//data.Id = jsontypes.StringValue("example-id")
	data.Id = jsontypes.StringValue(data.NetworkId.ValueString())

	// Now set the final state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Log that the resource was created.
	tflog.Trace(ctx, "created resource")
}

// Read method is responsible for reading an existing resource's state.
// It takes a ReadRequest and returns a ReadResponse with the current state of the resource or an error.
func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *resourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	// Check if there are any errors before proceeding.
	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := r.client.CellularGatewayApi.GetNetworkCellularGatewayUplink(context.Background(), data.NetworkId.ValueString()).Execute()

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

	// Decode the HTTP response body into your data model.
	// If there's an error, add it to diagnostics.
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	// Set ID for the resource.
	//data.Id = jsontypes.StringValue("example-id")
	data.Id = jsontypes.StringValue(data.NetworkId.ValueString())

	// Now set the final state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Log that the resource was read.
	tflog.Trace(ctx, "read resource")
}

// Update function is responsible for updating the state of an existing resource.
// It uses an UpdateRequest and responds with an UpdateResponse which contains the updated state of the resource or an error.
func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data *resourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// If there was an error reading the plan, return early.
	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworkCellularGatewayUplink := *openApiClient.NewUpdateNetworkCellularGatewayUplinkRequest()

	var bandwidthLimits openApiClient.UpdateNetworkApplianceTrafficShapingUplinkBandwidthRequestBandwidthLimitsCellular
	bandwidthLimits.SetLimitUp(int32(data.CellularGatewayBandwidthLimits.LimitUp.ValueInt64()))
	bandwidthLimits.SetLimitDown(int32(data.CellularGatewayBandwidthLimits.LimitDown.ValueInt64()))
	updateNetworkCellularGatewayUplink.SetBandwidthLimits(bandwidthLimits)

	_, httpResp, err := r.client.CellularGatewayApi.UpdateNetworkCellularGatewayUplink(context.Background(), data.NetworkId.ValueString()).UpdateNetworkCellularGatewayUplinkRequest(updateNetworkCellularGatewayUplink).Execute()

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

	// Decode the HTTP response body into your data model.
	// If there's an error, add it to diagnostics.
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	// Set ID for the new resource.
	//data.Id = jsontypes.StringValue("example-id")
	data.Id = jsontypes.StringValue(data.NetworkId.ValueString())

	// Now set the updated state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Log that the resource was updated.
	tflog.Trace(ctx, "updated resource")
}

// Delete function is responsible for deleting a resource.
// It uses a DeleteRequest and responds with a DeleteResponse which contains the updated state of the resource or an error.
func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data *resourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	// If there was an error reading the state, return early.
	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworkCellularGatewayUplink := *openApiClient.NewUpdateNetworkCellularGatewayUplinkRequest()

	var bandwidthLimits openApiClient.UpdateNetworkApplianceTrafficShapingUplinkBandwidthRequestBandwidthLimitsCellular
	bandwidthLimits.SetLimitUp(int32(data.CellularGatewayBandwidthLimits.LimitUp.ValueInt64()))
	bandwidthLimits.SetLimitDown(int32(data.CellularGatewayBandwidthLimits.LimitDown.ValueInt64()))
	updateNetworkCellularGatewayUplink.SetBandwidthLimits(bandwidthLimits)

	_, httpResp, err := r.client.CellularGatewayApi.UpdateNetworkCellularGatewayUplink(context.Background(), data.NetworkId.ValueString()).UpdateNetworkCellularGatewayUplinkRequest(updateNetworkCellularGatewayUplink).Execute()

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

	resp.State.RemoveResource(ctx)

	// Log that the resource was deleted.
	tflog.Trace(ctx, "removed resource")
}

// ImportState function is used to import an existing resource into Terraform.
// The function expects an ImportStateRequest and responds with an ImportStateResponse which contains
// the new state of the resource or an error.
func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	// Pass through the ID directly from the ImportStateRequest to the ImportStateResponse
	resource.ImportStatePassthroughID(ctx, path.Root("network_id"), req, resp)

}
