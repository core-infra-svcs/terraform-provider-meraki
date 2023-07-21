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

var (
	_ resource.Resource                = &NetworksCellularGatewaySubnetPoolResource{} // Terraform resource interface
	_ resource.ResourceWithConfigure   = &NetworksCellularGatewaySubnetPoolResource{} // Interface for resources with configuration methods
	_ resource.ResourceWithImportState = &NetworksCellularGatewaySubnetPoolResource{} // Interface for resources with import state functionality
)

func NewNetworksCellularGatewaySubnetPoolResource() resource.Resource {
	return &NetworksCellularGatewaySubnetPoolResource{}
}

// NetworksCellularGatewaySubnetPoolResource struct defines the structure for this resource.
// It includes an APIClient field for making requests to the Meraki API.
type NetworksCellularGatewaySubnetPoolResource struct {
	client *openApiClient.APIClient // APIClient instance for making API requests
}

// The NetworksCellularGatewaySubnetPoolResourceModel structure describes the data model.
// This struct is where you define all the attributes that are part of this resource's state.
type NetworksCellularGatewaySubnetPoolResourceModel struct {
	Id             jsontypes.String `tfsdk:"id"`
	NetworkId      jsontypes.String `tfsdk:"network_id"`
	Mask           jsontypes.Int64  `tfsdk:"mask"`
	Cidr           jsontypes.String `tfsdk:"cidr"`
	DeploymentMode jsontypes.String `tfsdk:"deployment_mode"`
	Subnets        []Subnet         `tfsdk:"subnets"`
}

type Subnet struct {
	Serial      jsontypes.String `tfsdk:"serial"`
	Name        jsontypes.String `tfsdk:"name"`
	ApplianceIp jsontypes.String `tfsdk:"appliance_ip"`
	Subnet      jsontypes.String `tfsdk:"subnet"`
}

// Metadata provides a way to define information about the resource.
// This method is called by the framework to retrieve metadata about the resource.
func (r *NetworksCellularGatewaySubnetPoolResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_networks_cellular_gateway_subnet_pool"
}

// Schema provides a way to define the structure of the resource data.
// It is called by the framework to get the schema of the resource.
func (r *NetworksCellularGatewaySubnetPoolResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	// The Schema object defines the structure of the resource.
	resp.Schema = schema.Schema{

		MarkdownDescription: "NetworksCellular Gateway SubnetPool Updates the subnet pool and mask configuration for MGs in the network.",

		// The Attributes map describes the fields of the resource.
		Attributes: map[string]schema.Attribute{

			// Every resource must have an ID attribute. This is computed by the framework.
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
			"mask": schema.Int64Attribute{
				MarkdownDescription: "Mask used for the subnet of all MGs in this network.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.Int64Type,
			},
			"cidr": schema.StringAttribute{
				MarkdownDescription: "CIDR of the pool of subnets. Each MG in this network will automatically pick a subnet from this pool.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"deployment_mode": schema.StringAttribute{
				Optional:   true,
				Computed:   true,
				CustomType: jsontypes.StringType,
			},
			"subnets": schema.SetNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"serial": schema.StringAttribute{
							Optional:   true,
							Computed:   true,
							CustomType: jsontypes.StringType,
						},
						"name": schema.StringAttribute{
							Optional:   true,
							Computed:   true,
							CustomType: jsontypes.StringType,
						},
						"appliance_ip": schema.StringAttribute{
							Optional:   true,
							Computed:   true,
							CustomType: jsontypes.StringType,
						},
						"subnet": schema.StringAttribute{
							Optional:   true,
							Computed:   true,
							CustomType: jsontypes.StringType,
						},
					},
				},
			},
		},
	}
}

// Configure is a method of the Resource interface that Terraform calls to provide the configured provider instance to the resource.
// It passes the ResourceData that's been stored by the provider's ConfigureFunc.
func (r *NetworksCellularGatewaySubnetPoolResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

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
func (r *NetworksCellularGatewaySubnetPoolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *NetworksCellularGatewaySubnetPoolResourceModel

	// TODO: Make sure the plan data matches the structure of the data model.
	// Unmarshal the plan data into the internal data model struct.
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// Check if there are any errors before proceeding.
	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworkCellularGatewaySubnetPool := *openApiClient.NewInlineObject69()
	updateNetworkCellularGatewaySubnetPool.SetCidr(data.Cidr.ValueString())
	updateNetworkCellularGatewaySubnetPool.SetMask(int32(data.Mask.ValueInt64()))

	_, httpResp, err := r.client.SubnetPoolApi.UpdateNetworkCellularGatewaySubnetPool(context.Background(), data.NetworkId.ValueString()).UpdateNetworkCellularGatewaySubnetPool(updateNetworkCellularGatewaySubnetPool).Execute()

	// If there was an error during API call, add it to diagnostics.
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// Collect any HTTP diagnostics that might be useful for debugging.
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
	}

	// If it's not what you expect, add an error to diagnostics.
	// TODO: Check the HTTP response status code matches the API endpoint.
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

	if len(data.Subnets) == 0 {
		data.Subnets = nil
	}

	// Set ID for the new resource.
	data.Id = jsontypes.StringValue("example-id")

	// Now set the final state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Log that the resource was created.
	tflog.Trace(ctx, "created resource")
}

// Read method is responsible for reading an existing resource's state.
// It takes a ReadRequest and returns a ReadResponse with the current state of the resource or an error.
func (r *NetworksCellularGatewaySubnetPoolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetworksCellularGatewaySubnetPoolResourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	// Check if there are any errors before proceeding.
	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := r.client.SubnetPoolApi.GetNetworkCellularGatewaySubnetPool(context.Background(), data.NetworkId.ValueString()).Execute()

	// If there was an error during API call, add it to diagnostics.
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// Collect any HTTP diagnostics that might be useful for debugging.
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
	}

	// If it's not what you expect, add an error to diagnostics.
	// TODO: Check the HTTP response status code matches the API endpoint.
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

	if len(data.Subnets) == 0 {
		data.Subnets = nil
	}

	// Set ID for the resource.
	data.Id = jsontypes.StringValue("example-id")

	// Now set the final state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Log that the resource was read.
	tflog.Trace(ctx, "read resource")
}

// Update function is responsible for updating the state of an existing resource.
// It uses an UpdateRequest and responds with an UpdateResponse which contains the updated state of the resource or an error.
func (r *NetworksCellularGatewaySubnetPoolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data *NetworksCellularGatewaySubnetPoolResourceModel

	// TODO: Make sure the plan data matches the structure of the data model.
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// If there was an error reading the plan, return early.
	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworkCellularGatewaySubnetPool := *openApiClient.NewInlineObject69()
	updateNetworkCellularGatewaySubnetPool.SetCidr(data.Cidr.ValueString())
	updateNetworkCellularGatewaySubnetPool.SetMask(int32(data.Mask.ValueInt64()))

	_, httpResp, err := r.client.SubnetPoolApi.UpdateNetworkCellularGatewaySubnetPool(context.Background(), data.NetworkId.ValueString()).UpdateNetworkCellularGatewaySubnetPool(updateNetworkCellularGatewaySubnetPool).Execute()

	// If there was an error during API call, add it to diagnostics.
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// Collect any HTTP diagnostics that might be useful for debugging.
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
	}

	// If it's not what you expect, add an error to diagnostics.
	// TODO: Check the HTTP response status code matches the API endpoint.
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

	if len(data.Subnets) == 0 {
		data.Subnets = nil
	}

	// Set ID for the new resource.
	data.Id = jsontypes.StringValue("example-id")

	// Now set the updated state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Log that the resource was updated.
	tflog.Trace(ctx, "updated resource")
}

// Delete function is responsible for deleting a resource.
// It uses a DeleteRequest and responds with a DeleteResponse which contains the updated state of the resource or an error.
func (r *NetworksCellularGatewaySubnetPoolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data *NetworksCellularGatewaySubnetPoolResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	// If there was an error reading the state, return early.
	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworkCellularGatewaySubnetPool := *openApiClient.NewInlineObject69()
	updateNetworkCellularGatewaySubnetPool.SetCidr(data.Cidr.ValueString())
	updateNetworkCellularGatewaySubnetPool.SetMask(int32(data.Mask.ValueInt64()))

	_, httpResp, err := r.client.SubnetPoolApi.UpdateNetworkCellularGatewaySubnetPool(context.Background(), data.NetworkId.ValueString()).UpdateNetworkCellularGatewaySubnetPool(updateNetworkCellularGatewaySubnetPool).Execute()

	// If there was an error during API call, add it to diagnostics.
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// Collect any HTTP diagnostics that might be useful for debugging.
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
	}

	// If it's not what you expect, add an error to diagnostics.
	// TODO: Check the HTTP response status code matches the API endpoint.
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

	// TODO: The resource has been deleted, so remove it from the state.
	resp.State.RemoveResource(ctx)

	// Log that the resource was deleted.
	tflog.Trace(ctx, "removed resource")
}

// ImportState function is used to import an existing resource into Terraform.
// The function expects an ImportStateRequest and responds with an ImportStateResponse which contains
// the new state of the resource or an error.
func (r *NetworksCellularGatewaySubnetPoolResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	// Pass through the ID directly from the ImportStateRequest to the ImportStateResponse
	resource.ImportStatePassthroughID(ctx, path.Root("network_id"), req, resp)

}
