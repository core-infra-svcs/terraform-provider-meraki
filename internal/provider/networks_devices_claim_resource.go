package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"

	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

var (
	_ resource.Resource                = &NetworksDevicesClaimResource{} // Terraform resource interface
	_ resource.ResourceWithImportState = &NetworksDevicesClaimResource{}
	_ resource.ResourceWithConfigure   = &NetworksDevicesClaimResource{} // Interface for resources with configuration methods
)

func NewNetworksDevicesClaimResource() resource.Resource {
	return &NetworksDevicesClaimResource{}
}

// NetworksDevicesClaimResource struct defines the structure for this resource.
// It includes an APIClient field for making requests to the Meraki API.
type NetworksDevicesClaimResource struct {
	client *openApiClient.APIClient // APIClient instance for making API requests
}

// The NetworksDevicesClaimResourceModel structure describes the data model.
// This struct is where you define all the attributes that are part of this resource's state.
type NetworksDevicesClaimResourceModel struct {
	Id        jsontypes.String   `tfsdk:"id"`
	NetworkId jsontypes.String   `tfsdk:"network_id"`
	Serials   []jsontypes.String `tfsdk:"serials"`
}

// The NetworksDevicesClaimResourceModelDevice structure describes the data model.
// This struct is where you define all the attributes that are part of this resource's state.
type NetworksDevicesClaimResourceModelDevice struct {
	// The Id field is mandatory for all resources. It's used for resource identification and is required
	// for the acceptance tests to run.
	Id              jsontypes.String                `tfsdk:"id"`
	Serial          jsontypes.String                `tfsdk:"serial"`
	Name            jsontypes.String                `tfsdk:"name"`
	Mac             jsontypes.String                `tfsdk:"mac"`
	Model           jsontypes.String                `tfsdk:"model"`
	Tags            jsontypes.Set[jsontypes.String] `tfsdk:"tags"`
	LanIp           jsontypes.String                `tfsdk:"lan_ip"`
	Firmware        jsontypes.String                `tfsdk:"firmware"`
	Lat             jsontypes.Float64               `tfsdk:"lat"`
	Lng             jsontypes.Float64               `tfsdk:"lng"`
	Address         jsontypes.String                `tfsdk:"address"`
	Notes           jsontypes.String                `tfsdk:"notes"`
	Url             jsontypes.String                `tfsdk:"url"`
	FloorPlanId     jsontypes.String                `tfsdk:"floor_plan_id"`
	NetworkId       jsontypes.String                `tfsdk:"network_id"`
	BeaconIdParams  types.Object                    `tfsdk:"beacon_id_params"`
	SwitchProfileId jsontypes.String                `tfsdk:"switch_profile_id"`
	MoveMapMarker   jsontypes.Bool                  `tfsdk:"move_map_marker"`
}
type NetworksDevicesClaimResourceModelDeviceBeaconIdParams struct {
	Uuid  jsontypes.String `tfsdk:"uuid"`
	Major jsontypes.Int64  `tfsdk:"major"`
	Minor jsontypes.Int64  `tfsdk:"minor"`
}

// Metadata provides a way to define information about the resource.
// This method is called by the framework to retrieve metadata about the resource.
func (r *NetworksDevicesClaimResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_networks_devices_claim"
}

// Schema provides a way to define the structure of the resource data.
// It is called by the framework to get the schema of the resource.
func (r *NetworksDevicesClaimResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	// The Schema object defines the structure of the resource.
	resp.Schema = schema.Schema{

		MarkdownDescription: "Claim devices into a network",

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
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 31),
				},
				CustomType: jsontypes.StringType,
			},
			"serials": schema.SetAttribute{
				MarkdownDescription: "The serials of the devices that should be claimed",
				ElementType:         jsontypes.StringType,
				CustomType:          jsontypes.SetType[jsontypes.String](),
				Computed:            true,
				Optional:            true,
			},
		},
	}
}

// Configure is a method of the Resource interface that Terraform calls to provide the configured provider instance to the resource.
// It passes the ResourceData that's been stored by the provider's ConfigureFunc.
func (r *NetworksDevicesClaimResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

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
func (r *NetworksDevicesClaimResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *NetworksDevicesClaimResourceModel

	// Unmarshal the plan data into the internal data model struct.
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// Check if there are any errors before proceeding.
	if resp.Diagnostics.HasError() {
		return
	}

	// serials
	var serials []string
	for _, serial := range data.Serials {
		serials = append(serials, serial.ValueString())
	}

	claimNetworkDevices := *openApiClient.NewClaimNetworkDevicesRequest(serials)

	httpResp, err := r.client.NetworksApi.ClaimNetworkDevices(ctx, data.NetworkId.ValueString()).ClaimNetworkDevicesRequest(claimNetworkDevices).Execute()

	// If there was an error during API call, add it to diagnostics.
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
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

	//// Extracting serial numbers from the response
	//var respSerials []string
	//for _, device := range inlineResp {
	//	if serial, ok := device["serial"].(string); ok {
	//		respSerials = append(respSerials, serial)
	//	}
	//}
	//
	//if len(serials) == 0 || (len(serials) == 1 && serials[0] == "") {
	//	fmt.Println("No serial numbers provided")
	//	os.Exit(0)
	//}
	//
	//var missingSerials, extraSerials []string
	//// Loop through API serials to find missing ones in provided serials
	//for _, apiSerial := range respSerials {
	//	found := false
	//	for _, providedSerial := range serials {
	//		if apiSerial == providedSerial {
	//			found = true
	//			break
	//		}
	//	}
	//	// If serial is not found in provided serials, add it to missingSerials
	//	if !found {
	//		missingSerials = append(missingSerials, apiSerial)
	//	}
	//}
	//
	//// Loop through provided serials to find extra ones not present in API serials
	//for _, providedSerial := range serials {
	//	found := false
	//	for _, apiSerial := range respSerials {
	//		if providedSerial == apiSerial {
	//			found = true
	//			break
	//		}
	//	}
	//	// If serial is not found in API serials, add it to extraSerials
	//	if !found {
	//		extraSerials = append(extraSerials, providedSerial)
	//	}
	//}
	//
	//// Print diagnostics based on the comparison results
	//if len(missingSerials) > 0 {
	//	resp.Diagnostics.AddError(
	//		"Missing serial numbers:", fmt.Sprintf("%s", missingSerials),
	//	)
	//	return
	//} else if len(extraSerials) > 0 {
	//	resp.Diagnostics.AddError(
	//		"Extra serial numbers:", fmt.Sprintf("%s", extraSerials),
	//	)
	//	return
	//} else if !reflect.DeepEqual(serials, respSerials) {
	//	resp.Diagnostics.AddError(
	//		"Serial numbers are out of order", fmt.Sprintf("%s", respSerials),
	//	)
	//	return
	//}

	// Set ID for the new resource.
	data.Id = jsontypes.StringValue(data.NetworkId.ValueString())

	// Now set the final state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Log that the resource was created.
	tflog.Trace(ctx, "created resource")
}

// Read method is responsible for reading an existing resource's state.
// It takes a ReadRequest and returns a ReadResponse with the current state of the resource or an error.
func (r *NetworksDevicesClaimResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetworksDevicesClaimResourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	// Check if there are any errors before proceeding.
	if resp.Diagnostics.HasError() {
		return
	}

	// serials
	var serials []string
	for _, serial := range data.Serials {
		serials = append(serials, serial.ValueString())
	}

	inlineResp, httpResp, err := r.client.NetworksApi.GetNetworkDevices(ctx, data.NetworkId.ValueString()).Execute()

	// If there was an error during API call, add it to diagnostics.
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
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

	var respSerials []jsontypes.String

	for _, inlineData := range inlineResp {
		var device NetworksDevicesClaimResourceModelDevice
		inlineDataBytes, err := json.Marshal(inlineData)
		err = json.Unmarshal(inlineDataBytes, &device)
		if err != nil {
			resp.Diagnostics.AddError(
				"Serial unmarshal error",
				fmt.Sprintf("%v", err.Error()),
			)
		}
		respSerials = append(respSerials, jsontypes.StringValue(device.Serial.ValueString()))
	}

	data.Serials = respSerials
	// Set ID for the resource.
	data.Id = jsontypes.StringValue(data.NetworkId.ValueString())

	// Now set the final state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Log that the resource was read.
	tflog.Trace(ctx, "read resource")
}

// Update function is responsible for updating the state of an existing resource.
// It uses an UpdateRequest and responds with an UpdateResponse which contains the updated state of the resource or an error.
func (r *NetworksDevicesClaimResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data, state *NetworksDevicesClaimResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	// If there was an error reading the plan, return early.
	if resp.Diagnostics.HasError() {
		return
	}

	// serials
	var planSerials, stateSerials, serialsToAdd, serialsToRemove []string
	for _, serial := range data.Serials {
		planSerials = append(planSerials, serial.ValueString())
	}

	for _, serial := range state.Serials {
		stateSerials = append(stateSerials, serial.ValueString())
	}

	serialsToAdd = difference(planSerials, stateSerials)
	serialsToRemove = difference(stateSerials, planSerials)

	claimNetworkDevices := *openApiClient.NewClaimNetworkDevicesRequest(serialsToAdd)

	httpResp, err := r.client.NetworksApi.ClaimNetworkDevices(ctx, data.NetworkId.ValueString()).ClaimNetworkDevicesRequest(claimNetworkDevices).Execute()

	// If there was an error during API call, add it to diagnostics.
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
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

	// deleting serials
	for _, serial := range serialsToRemove {

		se := fmt.Sprint(strings.Trim(serial, "\""))

		removeNetworkDevices := *openApiClient.NewRemoveNetworkDevicesRequest(se)

		httpResp, err := r.client.NetworksApi.RemoveNetworkDevices(ctx, data.NetworkId.ValueString()).RemoveNetworkDevicesRequest(removeNetworkDevices).Execute()

		// If there was an error during API call, add it to diagnostics.
		if err != nil {
			resp.Diagnostics.AddError(
				"HTTP Client Failure",
				tools.HttpDiagnostics(httpResp),
			)
			return
		}

		// If it's not what you expect, add an error to diagnostics.
		if httpResp.StatusCode != 204 {
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

	}

	// If there were any errors up to this point, log the plan data and return.
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
		return
	}

	// Set ID for the resource.
	data.Id = jsontypes.StringValue("example-id")

	// Now set the updated state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Log that the resource was updated.
	tflog.Trace(ctx, "updated resource")
}

// Delete function is responsible for deleting a resource.
// It uses a DeleteRequest and responds with a DeleteResponse which contains the updated state of the resource or an error.
func (r *NetworksDevicesClaimResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data *NetworksDevicesClaimResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	// deleting serials
	for _, serial := range data.Serials {

		se := fmt.Sprint(strings.Trim(serial.String(), "\""))

		removeNetworkDevices := *openApiClient.NewRemoveNetworkDevicesRequest(se)

		httpResp, err := r.client.NetworksApi.RemoveNetworkDevices(ctx, data.NetworkId.ValueString()).RemoveNetworkDevicesRequest(removeNetworkDevices).Execute()

		// If there was an error during API call, add it to diagnostics.
		if err != nil {
			resp.Diagnostics.AddError(
				"HTTP Client Failure",
				tools.HttpDiagnostics(httpResp),
			)
			return
		}

		// If it's not what you expect, add an error to diagnostics.
		if httpResp.StatusCode != 204 {
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

	}
	resp.State.RemoveResource(ctx)

	// Log that the resource was deleted.
	tflog.Trace(ctx, "removed resource")
}

func (r *NetworksDevicesClaimResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	resource.ImportStatePassthroughID(ctx, path.Root("network_id"), req, resp)

	if resp.Diagnostics.HasError() {
		return
	}
}

// difference returns the elements in `a` that aren't in `b`.
func difference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}
