package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"strings"

	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

var (
	_ resource.Resource                = &DevicesResource{} // Terraform resource interface
	_ resource.ResourceWithConfigure   = &DevicesResource{} // Interface for resources with configuration methods
	_ resource.ResourceWithImportState = &DevicesResource{} // Interface for resources with import state functionality
)

func NewDevicesResource() resource.Resource {
	return &DevicesResource{}
}

// DevicesResource struct defines the structure for this resource.
// It includes an APIClient field for making requests to the Meraki API.
type DevicesResource struct {
	client *openApiClient.APIClient // APIClient instance for making API requests
}

// The DevicesResourceModel structure describes the data model.
// This struct is where you define all the attributes that are part of this resource's state.
type DevicesResourceModel struct {

	// The Id field is mandatory for all resources. It's used for resource identification and is required
	// for the acceptance tests to run.
	Id              jsontypes.String                   `tfsdk:"id"`
	Serial          jsontypes.String                   `tfsdk:"serial"`
	Name            jsontypes.String                   `tfsdk:"name"`
	Mac             jsontypes.String                   `tfsdk:"mac"`
	Model           jsontypes.String                   `tfsdk:"model"`
	Tags            jsontypes.Set[jsontypes.String]    `tfsdk:"tags"`
	LanIp           jsontypes.String                   `tfsdk:"lan_ip"`
	Firmware        jsontypes.String                   `tfsdk:"firmware"`
	Lat             jsontypes.Float64                  `tfsdk:"lat"`
	Lng             jsontypes.Float64                  `tfsdk:"lng"`
	Address         jsontypes.String                   `tfsdk:"address"`
	Notes           jsontypes.String                   `tfsdk:"notes"`
	Url             jsontypes.String                   `tfsdk:"url"`
	FloorPlanId     jsontypes.String                   `tfsdk:"floor_plan_id"`
	NetworkId       jsontypes.String                   `tfsdk:"network_id"`
	BeaconIdParams  DevicesResourceModelBeaconIdParams `tfsdk:"beacon_id_params"`
	SwitchProfileId jsontypes.String                   `tfsdk:"switch_profile_id"`
	MoveMapMarker   jsontypes.Bool                     `tfsdk:"move_map_marker"`
}

type DevicesResourceModelBeaconIdParams struct {
	Uuid  jsontypes.String `tfsdk:"uuid"`
	Major jsontypes.Int64  `tfsdk:"major"`
	Minor jsontypes.Int64  `tfsdk:"minor"`
}

// Metadata provides a way to define information about the resource.
// This method is called by the framework to retrieve metadata about the resource.
func (r *DevicesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {

	// The TypeName attribute is important as it provides the user-friendly name for the resource/data source.
	// This is the name users will use to reference the resource/data source and it's also used in the acceptance tests.
	resp.TypeName = req.ProviderTypeName + "_devices"
}

// Schema provides a way to define the structure of the resource data.
// It is called by the framework to get the schema of the resource.
func (r *DevicesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {

	// The Schema object defines the structure of the resource.
	resp.Schema = schema.Schema{

		MarkdownDescription: "Manage network Devices resource. This only works for devices associated with a network.",

		// The Attributes map describes the fields of the resource.
		Attributes: map[string]schema.Attribute{

			// Every resource must have an ID attribute. This is computed by the framework.
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.StringType,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Network ID",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				CustomType: jsontypes.StringType,
			},
			"serial": schema.StringAttribute{
				MarkdownDescription: "The devices serial number",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 31),
				},
				CustomType: jsontypes.StringType,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of a device",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"tags": schema.SetAttribute{
				Description: "Network tags",
				CustomType:  jsontypes.SetType[jsontypes.String](),
				Computed:    true,
				Optional:    true,
				PlanModifiers: []planmodifier.Set{
					setplanmodifier.UseStateForUnknown(),
				},
			},
			"lat": schema.Float64Attribute{
				MarkdownDescription: "The latitude of a device",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.Float64Type,
			},
			"lng": schema.Float64Attribute{
				MarkdownDescription: "The longitude of a device",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.Float64Type,
			},
			"address": schema.StringAttribute{
				MarkdownDescription: "The address of a device",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"notes": schema.StringAttribute{
				MarkdownDescription: "Notes for the network",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"move_map_marker": schema.BoolAttribute{
				MarkdownDescription: "Whether or not to set the latitude and longitude of a device based on the new address. Only applies when lat and lng are not specified.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.BoolType,
			},
			"switch_profile_id": schema.StringAttribute{
				MarkdownDescription: "The ID of a switch profile to bind to the device (for available switch profiles, see the 'Switch Profiles' endpoint). Use null to unbind the switch device from the current profile. For a device to be bindable to a switch profile, it must (1) be a switch, and (2) belong to a network that is bound to a configuration template.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"floor_plan_id": schema.StringAttribute{
				MarkdownDescription: "The floor plan to associate to this device. null disassociates the device from the floor plan.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"mac": schema.StringAttribute{
				MarkdownDescription: "The mac address of a device",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"model": schema.StringAttribute{
				MarkdownDescription: "The model of a device",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"lan_ip": schema.StringAttribute{
				MarkdownDescription: "The ipv4 lan ip of a device",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"firmware": schema.StringAttribute{
				MarkdownDescription: "The firmware version of a device",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "The url for the network associated with the device.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"beacon_id_params": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"uuid": schema.StringAttribute{
						Optional:   true,
						Computed:   true,
						CustomType: jsontypes.StringType,
					},
					"major": schema.Int64Attribute{
						Optional:   true,
						Computed:   true,
						CustomType: jsontypes.Int64Type,
					},
					"minor": schema.Int64Attribute{
						Optional:   true,
						Computed:   true,
						CustomType: jsontypes.Int64Type,
					},
				},
			},
		},
	}
}

// Configure is a method of the Resource interface that Terraform calls to provide the configured provider instance to the resource.
// It passes the ResourceData that's been stored by the provider's ConfigureFunc.
func (r *DevicesResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {

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
func (r *DevicesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *DevicesResourceModel

	// Unmarshal the plan data into the internal data model struct.
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// Check if there are any errors before proceeding.
	if resp.Diagnostics.HasError() {
		return
	}

	// Check for required parameters
	if len(data.Serial.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing device serial number", fmt.Sprintf("S/N: %s", data.Serial.ValueString()))
		return
	}

	// Create HTTP request body
	updateDevice := openApiClient.NewUpdateDeviceRequest()

	updateDevice.SetName(data.Name.ValueString())
	updateDevice.SetAddress(data.Address.ValueString())

	// Tags
	if !data.Tags.IsNull() {
		var tags []string
		for _, attribute := range data.Tags.Elements() {
			tag := fmt.Sprint(strings.Trim(attribute.String(), "\""))
			tags = append(tags, tag)
		}
		updateDevice.SetTags(tags)
	}

	updateDevice.SetLat(float32(data.Lat.ValueFloat64()))
	updateDevice.SetLng(float32(data.Lng.ValueFloat64()))
	updateDevice.SetNotes(data.Notes.ValueString())
	updateDevice.SetMoveMapMarker(data.MoveMapMarker.ValueBool())

	// SwitchProfileId
	if len(data.SwitchProfileId.ValueString()) > 1 {
		updateDevice.SetSwitchProfileId(data.SwitchProfileId.ValueString())
	}

	//	FloorPlanId
	if len(data.FloorPlanId.ValueString()) > 1 {
		updateDevice.SetFloorPlanId(data.FloorPlanId.ValueString())
	}

	// Initialize provider client and make API call
	_, httpResp, err := r.client.DevicesApi.UpdateDevice(context.Background(), data.Serial.ValueString()).UpdateDeviceRequest(*updateDevice).Execute()
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

	// Decode the HTTP response body into your data model.
	// If there's an error, add it to diagnostics.
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}
	if data.BeaconIdParams.Major.IsUnknown() {
		data.BeaconIdParams.Major = jsontypes.Int64Null()
	}
	if data.BeaconIdParams.Minor.IsUnknown() {
		data.BeaconIdParams.Minor = jsontypes.Int64Null()
	}
	if data.BeaconIdParams.Uuid.IsUnknown() {
		data.BeaconIdParams.Uuid = jsontypes.StringNull()
	}
	if data.Firmware.IsUnknown() {
		data.Firmware = jsontypes.StringNull()
	}
	if data.FloorPlanId.IsUnknown() {
		data.FloorPlanId = jsontypes.StringNull()
	}
	if data.LanIp.IsUnknown() {
		data.LanIp = jsontypes.StringNull()
	}
	if data.Mac.IsUnknown() {
		data.Mac = jsontypes.StringNull()
	}
	if data.Url.IsUnknown() {
		data.Url = jsontypes.StringNull()
	}
	if data.Model.IsUnknown() {
		data.Model = jsontypes.StringNull()
	}
	if data.SwitchProfileId.IsUnknown() {
		data.SwitchProfileId = jsontypes.StringNull()
	}
	if data.MoveMapMarker.IsUnknown() {
		data.MoveMapMarker = jsontypes.BoolNull()
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
func (r *DevicesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *DevicesResourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	// Check if there are any errors before proceeding.
	if resp.Diagnostics.HasError() {
		return
	}

	// Check for required parameters
	if len(data.Serial.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing device serial number", fmt.Sprintf("S/N: %s", data.Serial.ValueString()))
		return
	}

	// Initialize provider client and make API call
	_, httpResp, err := r.client.DevicesApi.GetDevice(context.Background(), data.Serial.ValueString()).Execute()
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

	if data.BeaconIdParams.Major.IsUnknown() {
		data.BeaconIdParams.Major = jsontypes.Int64Null()
	}
	if data.BeaconIdParams.Minor.IsUnknown() {
		data.BeaconIdParams.Minor = jsontypes.Int64Null()
	}
	if data.BeaconIdParams.Uuid.IsUnknown() {
		data.BeaconIdParams.Uuid = jsontypes.StringNull()
	}
	if data.Firmware.IsUnknown() {
		data.Firmware = jsontypes.StringNull()
	}
	if data.FloorPlanId.IsUnknown() {
		data.FloorPlanId = jsontypes.StringNull()
	}
	if data.LanIp.IsUnknown() {
		data.LanIp = jsontypes.StringNull()
	}
	if data.Mac.IsUnknown() {
		data.Mac = jsontypes.StringNull()
	}
	if data.Url.IsUnknown() {
		data.Url = jsontypes.StringNull()
	}
	if data.Model.IsUnknown() {
		data.Model = jsontypes.StringNull()
	}
	if data.SwitchProfileId.IsUnknown() {
		data.SwitchProfileId = jsontypes.StringNull()
	}
	if data.MoveMapMarker.IsUnknown() {
		data.MoveMapMarker = jsontypes.BoolNull()
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
func (r *DevicesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data *DevicesResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// If there was an error reading the plan, return early.
	if resp.Diagnostics.HasError() {
		return
	}

	// Check for required parameters
	if len(data.Serial.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing device serial number", fmt.Sprintf("S/N: %s", data.Serial.ValueString()))
		return
	}

	// Create HTTP request body
	updateDevice := openApiClient.NewUpdateDeviceRequest()

	updateDevice.SetName(data.Name.ValueString())
	updateDevice.SetAddress(data.Address.ValueString())

	// Tags
	if !data.Tags.IsNull() {
		var tags []string
		for _, attribute := range data.Tags.Elements() {
			tag := fmt.Sprint(strings.Trim(attribute.String(), "\""))
			tags = append(tags, tag)
		}
		updateDevice.SetTags(tags)
	}

	updateDevice.SetLat(float32(data.Lat.ValueFloat64()))
	updateDevice.SetLng(float32(data.Lng.ValueFloat64()))
	updateDevice.SetNotes(data.Notes.ValueString())
	updateDevice.SetMoveMapMarker(data.MoveMapMarker.ValueBool())

	// SwitchProfileId
	if len(data.SwitchProfileId.ValueString()) > 1 {
		updateDevice.SetSwitchProfileId(data.SwitchProfileId.ValueString())
	}

	//	FloorPlanId
	if len(data.FloorPlanId.ValueString()) > 1 {
		updateDevice.SetFloorPlanId(data.FloorPlanId.ValueString())
	}

	// Initialize provider client and make API call
	_, httpResp, err := r.client.DevicesApi.UpdateDevice(context.Background(),
		data.Serial.ValueString()).UpdateDeviceRequest(*updateDevice).Execute()
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

	// Decode the HTTP response body into your data model.
	// If there's an error, add it to diagnostics.
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	if data.BeaconIdParams.Major.IsUnknown() {
		data.BeaconIdParams.Major = jsontypes.Int64Null()
	}
	if data.BeaconIdParams.Minor.IsUnknown() {
		data.BeaconIdParams.Minor = jsontypes.Int64Null()
	}
	if data.BeaconIdParams.Uuid.IsUnknown() {
		data.BeaconIdParams.Uuid = jsontypes.StringNull()
	}
	if data.Firmware.IsUnknown() {
		data.Firmware = jsontypes.StringNull()
	}
	if data.FloorPlanId.IsUnknown() {
		data.FloorPlanId = jsontypes.StringNull()
	}
	if data.LanIp.IsUnknown() {
		data.LanIp = jsontypes.StringNull()
	}
	if data.Mac.IsUnknown() {
		data.Mac = jsontypes.StringNull()
	}
	if data.Url.IsUnknown() {
		data.Url = jsontypes.StringNull()
	}
	if data.Model.IsUnknown() {
		data.Model = jsontypes.StringNull()
	}
	if data.SwitchProfileId.IsUnknown() {
		data.SwitchProfileId = jsontypes.StringNull()
	}
	if data.MoveMapMarker.IsUnknown() {
		data.MoveMapMarker = jsontypes.BoolNull()
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
func (r *DevicesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data *DevicesResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	// If there was an error reading the state, return early.
	if resp.Diagnostics.HasError() {
		return
	}

	// Check for required parameters
	if len(data.Serial.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing device serial number", fmt.Sprintf("S/N: %s", data.Serial.ValueString()))
		return
	}

	// Create HTTP request body
	updateDevice := openApiClient.NewUpdateDeviceRequest()

	var name string
	var tags []string
	var lat float32
	var lng float32
	var address string
	var notes string
	var moveMapMarker bool
	//var switchProfileId string
	//var floorPlanId string

	updateDevice.Name = &name
	updateDevice.Tags = tags
	updateDevice.Lat = &lat
	updateDevice.Lng = &lng
	updateDevice.Address = &address
	updateDevice.Notes = &notes
	updateDevice.MoveMapMarker = &moveMapMarker
	//updateDevice.SwitchProfileId = &switchProfileId
	//updateDevice.FloorPlanId = &floorPlanId

	// Initialize provider client and make API call
	_, httpResp, err := r.client.DevicesApi.UpdateDevice(context.Background(),
		data.Serial.ValueString()).UpdateDeviceRequest(*updateDevice).Execute()
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
	data.Id = jsontypes.StringValue("example-id")

	resp.State.RemoveResource(ctx)

	// Log that the resource was deleted.
	tflog.Trace(ctx, "removed resource")
}

// ImportState function is used to import an existing resource into Terraform.
// The function expects an ImportStateRequest and responds with an ImportStateResponse which contains
// the new state of the resource or an error.
func (r *DevicesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	// Pass through the ID directly from the ImportStateRequest to the ImportStateResponse
	resource.ImportStatePassthroughID(ctx, path.Root("serial"), req, resp)

}
