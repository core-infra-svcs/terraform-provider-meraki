package devices

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"

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
	_ resource.Resource                = &DevicesResource{} // Terraform resource interface
	_ resource.ResourceWithConfigure   = &DevicesResource{} // Interface for resources with configuration methods
	_ resource.ResourceWithImportState = &DevicesResource{} // Interface for resources with import state functionality
)

func NewDevicesResource() resource.Resource {
	return &DevicesResource{}
}

// DevicesResource struct defines the structure for this resource.
type DevicesResource struct {
	client *openApiClient.APIClient // APIClient instance for making API requests
}

// The DevicesResourceModel structure describes the data model.
type DevicesResourceModel struct {
	Id              types.String  `tfsdk:"id"`
	Serial          types.String  `tfsdk:"serial"`
	Name            types.String  `tfsdk:"name"`
	Mac             types.String  `tfsdk:"mac"`
	Model           types.String  `tfsdk:"model"`
	Tags            types.List    `tfsdk:"tags"`
	Details         types.List    `tfsdk:"details"`
	LanIp           types.String  `tfsdk:"lan_ip"`
	Firmware        types.String  `tfsdk:"firmware"`
	Lat             types.Float64 `tfsdk:"lat"`
	Lng             types.Float64 `tfsdk:"lng"`
	Address         types.String  `tfsdk:"address"`
	Notes           types.String  `tfsdk:"notes"`
	Url             types.String  `tfsdk:"url"`
	FloorPlanId     types.String  `tfsdk:"floor_plan_id"`
	NetworkId       types.String  `tfsdk:"network_id"`
	BeaconIdParams  types.Object  `tfsdk:"beacon_id_params"`
	SwitchProfileId types.String  `tfsdk:"switch_profile_id"`
	MoveMapMarker   types.Bool    `tfsdk:"move_map_marker"`
}

type DevicesResourceModelBeaconIdParams struct {
	Uuid  types.String `tfsdk:"uuid"`
	Major types.Int64  `tfsdk:"major"`
	Minor types.Int64  `tfsdk:"minor"`
}

type DevicesResourceModelDetails struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}

// Metadata provides a way to define information about the resource.
func (r *DevicesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_devices"
}

// Schema provides a way to define the structure of the resource data.
func (r *DevicesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{

		MarkdownDescription: "Manage network Devices resource. This only works for devices associated with a network.",

		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{
				Computed: true,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Network ID",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"serial": schema.StringAttribute{
				MarkdownDescription: "The devices serial number",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(14, 14),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of a device",
				Optional:            true,
				Computed:            true,
			},
			"tags": schema.ListAttribute{
				Description: "Network tags",
				Computed:    true,
				Optional:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"lat": schema.Float64Attribute{
				MarkdownDescription: "The latitude of a device",
				Optional:            true,
				Computed:            true,
			},
			"lng": schema.Float64Attribute{
				MarkdownDescription: "The longitude of a device",
				Optional:            true,
				Computed:            true,
			},
			"address": schema.StringAttribute{
				MarkdownDescription: "The address of a device",
				Optional:            true,
				Computed:            true,
			},
			"notes": schema.StringAttribute{
				MarkdownDescription: "Notes for the network",
				Optional:            true,
				Computed:            true,
			},
			"details": schema.ListNestedAttribute{
				Description: "Network tags",
				Computed:    true,
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						MarkdownDescription: "The name of a device",
						Optional:            true,
						Computed:            true,
					},
					"value": schema.StringAttribute{
						MarkdownDescription: "The value of a device",
						Optional:            true,
						Computed:            true,
					},
				}},
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"move_map_marker": schema.BoolAttribute{
				MarkdownDescription: "Whether or not to set the latitude and longitude of a device based on the new address. Only applies when lat and lng are not specified.",
				Optional:            true,
				Computed:            true,
			},
			"switch_profile_id": schema.StringAttribute{
				MarkdownDescription: "The ID of a switch profile to bind to the device (for available switch profiles, see the 'Switch Profiles' endpoint). Use null to unbind the switch device from the current profile. For a device to be bindable to a switch profile, it must (1) be a switch, and (2) belong to a network that is bound to a configuration template.",
				Optional:            true,
				Computed:            true,
			},
			"floor_plan_id": schema.StringAttribute{
				MarkdownDescription: "The floor plan to associate to this device. null disassociates the device from the floor plan.",
				Optional:            true,
				Computed:            true,
			},
			"mac": schema.StringAttribute{
				MarkdownDescription: "The mac address of a device",
				Optional:            true,
				Computed:            true,
			},
			"model": schema.StringAttribute{
				MarkdownDescription: "The model of a device",
				Optional:            true,
				Computed:            true,
			},
			"lan_ip": schema.StringAttribute{
				MarkdownDescription: "The ipv4 lan ip of a device",
				Optional:            true,
				Computed:            true,
			},
			"firmware": schema.StringAttribute{
				MarkdownDescription: "The firmware version of a device",
				Optional:            true,
				Computed:            true,
			},
			"url": schema.StringAttribute{
				MarkdownDescription: "The url for the network associated with the device.",
				Optional:            true,
				Computed:            true,
			},
			"beacon_id_params": schema.SingleNestedAttribute{
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"uuid": schema.StringAttribute{
						Computed: true,
					},
					"major": schema.Int64Attribute{
						Computed: true,
					},
					"minor": schema.Int64Attribute{
						Computed: true,
					},
				},
			},
		},
	}
}

// Configure is a method of the Resource interface that Terraform calls to provide the configured provider instance to the resource.
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

// updateDevicesResourceState updates the resource state with the provided api data.
func updateDevicesResourceState(ctx context.Context, state *DevicesResourceModel, httpResp *http.Response) diag.Diagnostics {
	var diags diag.Diagnostics

	httpRespMap, httpRespMapErr := utils.ExtractResponseToMap(httpResp)
	if httpRespMapErr != nil {
		diags.AddError("Failed to Unmarshal HttpResp", httpRespMapErr.Error())
	}

	// "name": "My AP",
	if state.Name.IsNull() || state.Name.IsUnknown() {
		name, err := utils.ExtractStringAttr(httpRespMap, "name")
		if err != nil {
			diags.Append(err...)
		}
		state.Name = name
	}

	//  "lat": 37.4180951010362,
	if state.Lat.IsNull() || state.Lat.IsUnknown() {
		lat, err := utils.ExtractFloat64Attr(httpRespMap, "lat")
		if err != nil {
			diags.Append(err...)
		}
		state.Lat = lat
	}

	//  "lng": -122.098531723022,
	if state.Lng.IsNull() || state.Lng.IsUnknown() {
		lng, err := utils.ExtractFloat64Attr(httpRespMap, "lng")
		if err != nil {
			diags.Append(err...)
		}
		state.Lng = lng
	}

	//  "address": "1600 Pennsylvania Ave",
	if state.Address.IsNull() || state.Address.IsUnknown() {
		address, err := utils.ExtractStringAttr(httpRespMap, "address")
		if err != nil {
			diags.Append(err...)
		}
		state.Address = address
	}

	//  "notes": "My AP's note",
	if state.Notes.IsNull() || state.Notes.IsUnknown() {
		notes, err := utils.ExtractStringAttr(httpRespMap, "notes")
		if err != nil {
			diags.Append(err...)
		}
		state.Notes = notes
	}

	//  "tags": [
	//    " recently-added "
	//  ],
	if state.Tags.IsNull() || state.Tags.IsUnknown() {
		tags, err := utils.ExtractListStringAttr(httpRespMap, "tags")
		if err != nil {
			diags.Append(err...)
		}

		state.Tags = tags
	}

	//  "networkId": "N_24329156",
	if state.NetworkId.IsNull() || state.NetworkId.IsUnknown() {
		networkId, err := utils.ExtractStringAttr(httpRespMap, "networkId")
		if err != nil {
			diags.Append(err...)
		}
		state.NetworkId = networkId
	}

	//  "serial": "Q234-ABCD-5678",
	if state.Serial.IsNull() || state.Serial.IsUnknown() {
		diags.AddError("Missing Serial", "Missing Serial After State Update")
	}

	//  "model": "MR34",
	if state.Model.IsNull() || state.Model.IsUnknown() {
		model, err := utils.ExtractStringAttr(httpRespMap, "model")
		if err != nil {
			diags.Append(err...)
		}
		state.Model = model
	}

	//  "mac": "00:11:22:33:44:55",
	if state.Mac.IsNull() || state.Mac.IsUnknown() {
		mac, err := utils.ExtractStringAttr(httpRespMap, "mac")
		if err != nil {
			diags.Append(err...)
		}
		state.Mac = mac
	}

	//  "lanIp": "1.2.3.4",
	if state.LanIp.IsNull() || state.LanIp.IsUnknown() {
		lanIp, err := utils.ExtractStringAttr(httpRespMap, "lanIp")
		if err != nil {
			diags.Append(err...)
		}
		state.LanIp = lanIp
	}

	//  "firmware": "wireless-25-14",
	if state.Firmware.IsNull() || state.Firmware.IsUnknown() {
		firmware, err := utils.ExtractStringAttr(httpRespMap, "firmware")
		if err != nil {
			diags.Append(err...)
		}
		state.Firmware = firmware
	}

	//  "floorPlanId": "g_2176982374",
	if state.FloorPlanId.IsNull() || state.FloorPlanId.IsUnknown() {
		floorPlanId, err := utils.ExtractStringAttr(httpRespMap, "floorPlanId")
		if err != nil {
			diags.Append(err...)
		}
		state.FloorPlanId = floorPlanId
	}

	//  "details": [
	//    {
	//      "name": "Catalyst serial",
	//      "value": "123ABC"
	//    }
	//  ],
	if state.Details.IsNull() || state.Details.IsUnknown() {

		detailAttr := map[string]attr.Type{
			"name":  types.StringType,
			"value": types.StringType,
		}

		detailsAttrs := types.ObjectType{AttrTypes: detailAttr}

		_, ok := httpRespMap["details"].([]map[string]interface{})
		if ok {

			detailsList, err := utils.ExtractListAttr(httpRespMap, "details", detailsAttrs)
			if err.HasError() {
				tflog.Error(context.Background(), fmt.Sprintf("%s", err.Errors()))
			}

			state.Details = detailsList

		} else {
			detailsArrayObjNull := types.ListNull(detailsAttrs)
			state.Details = detailsArrayObjNull
		}

	}

	//  "beaconIdParams": {
	//    "uuid": "00000000-0000-0000-0000-000000000000",
	//    "major": 5,
	//    "minor": 3
	if state.BeaconIdParams.IsNull() || state.BeaconIdParams.IsUnknown() {
		beaconIdParamsAttrs := map[string]attr.Type{
			"uuid":  types.StringType,
			"major": types.Int64Type,
			"minor": types.Int64Type,
		}

		beaconIdParamsResp, ok := httpRespMap["beaconIdParams"].(map[string]interface{})
		if ok {
			var beaconIdParams DevicesResourceModelBeaconIdParams

			// uuid
			uuid, err := utils.ExtractStringAttr(beaconIdParamsResp, "uuid")
			if err.HasError() {
				diags.AddError("uuid Attr", fmt.Sprintf("%s", err.Errors()))
			}

			beaconIdParams.Uuid = uuid

			// major
			major, err := utils.ExtractInt32Attr(beaconIdParamsResp, "major")
			if err.HasError() {
				diags.AddError("major Attr", fmt.Sprintf("%s", err.Errors()))
			}

			beaconIdParams.Major = major

			// minor
			minor, err := utils.ExtractInt32Attr(beaconIdParamsResp, "minor")
			if err.HasError() {
				diags.AddError("minor Attr", fmt.Sprintf("%s", err.Errors()))
			}

			beaconIdParams.Minor = minor

			beaconIdParamsObj, err := types.ObjectValueFrom(ctx, beaconIdParamsAttrs, beaconIdParams)
			if err.HasError() {
				diags.AddError("beaconIdParamsObj Attr", fmt.Sprintf("%s", err.Errors()))
			}

			state.BeaconIdParams = beaconIdParamsObj
		} else {
			beaconIdParamsObjNull := types.ObjectNull(beaconIdParamsAttrs)
			state.BeaconIdParams = beaconIdParamsObjNull
		}

	}

	// url
	if state.Url.IsNull() || state.Url.IsUnknown() {
		url, err := utils.ExtractStringAttr(httpRespMap, "url")
		if err.HasError() {
			diags.AddError("url Attr", fmt.Sprintf("%s", err.Errors()))
		}
		state.Url = url
	}

	// SwitchProfileId
	if state.SwitchProfileId.IsNull() || state.SwitchProfileId.IsUnknown() {
		switchProfileId, err := utils.ExtractStringAttr(httpRespMap, "switchProfileId")
		if err.HasError() {
			diags.AddError("switchProfileId Attr", fmt.Sprintf("%s", err.Errors()))
		}
		state.SwitchProfileId = switchProfileId
	}

	// MoveMapMarker
	if state.MoveMapMarker.IsNull() || state.MoveMapMarker.IsUnknown() {
		moveMapMarker, err := utils.ExtractBoolAttr(httpRespMap, "moveMapMarker")
		if err.HasError() {
			diags.AddError("moveMapMarker Attr", fmt.Sprintf("%s", err.Errors()))
		}
		state.MoveMapMarker = moveMapMarker
	}

	// Set ID for the new resource.
	if state.Id.IsNull() || state.Id.IsUnknown() {
		state.Id = state.Serial
	}

	return diags
}

func updateDevicesResourcePayload(plan *DevicesResourceModel) (openApiClient.UpdateDeviceRequest, diag.Diagnostics) {
	var diags diag.Diagnostics
	payload := openApiClient.NewUpdateDeviceRequest()

	//    Name
	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		payload.SetName(plan.Name.ValueString())
	}

	//    Tags
	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
		var tags []string
		for _, attribute := range plan.Tags.Elements() {
			tag := fmt.Sprint(strings.Trim(attribute.String(), "\""))
			tags = append(tags, tag)
		}
		payload.SetTags(tags)
	}

	//    Lat
	if !plan.Lat.IsNull() && !plan.Lat.IsUnknown() {
		payload.SetLat(float32(plan.Lat.ValueFloat64()))

	}

	//    Lng
	if !plan.Lng.IsNull() && !plan.Lng.IsUnknown() {
		payload.SetLng(float32(plan.Lng.ValueFloat64()))
	}

	//    Address
	if !plan.Address.IsNull() && !plan.Address.IsUnknown() {
		payload.SetAddress(plan.Address.ValueString())
	}

	//    Notes
	if !plan.Notes.IsNull() && !plan.Notes.IsUnknown() {
		payload.SetNotes(plan.Notes.ValueString())
	}

	//    MoveMapMarker
	if !plan.MoveMapMarker.IsNull() && !plan.MoveMapMarker.IsUnknown() {
		payload.SetMoveMapMarker(plan.MoveMapMarker.ValueBool())
	}

	//    SwitchProfileId
	if !plan.SwitchProfileId.IsNull() && !plan.SwitchProfileId.IsUnknown() {
		payload.SetSwitchProfileId(plan.SwitchProfileId.ValueString())
	}

	//    FloorPlanId
	if !plan.FloorPlanId.IsNull() && !plan.FloorPlanId.IsUnknown() {
		payload.SetFloorPlanId(plan.FloorPlanId.ValueString())
	}

	return *payload, diags

}

// Create method is responsible for creating a new resource.
func (r *DevicesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *DevicesResourceModel

	// Unmarshal the plan data into the internal data model struct.
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// Check if there are any errors before proceeding.
	if resp.Diagnostics.HasError() {
		return
	}

	// Create HTTP request body
	payload, payloadErr := updateDevicesResourcePayload(data)
	if payloadErr.HasError() {
		resp.Diagnostics.AddError("Failed to assemble payload", fmt.Sprintf("%s", payloadErr.Errors()))
	}

	// Initialize provider client and make API call
	maxRetries := r.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(r.client.GetConfig().Retry4xxErrorWaitTime)

	// API call function to be passed to retryOn4xx
	apiCall := func() (map[string]interface{}, *http.Response, error) {
		inline, httpResp, err := r.client.DevicesApi.UpdateDevice(ctx, data.Serial.ValueString()).UpdateDeviceRequest(payload).Execute()
		return inline, httpResp, err
	}

	inlineResp, httpResp, err := utils.CustomHttpRequestRetry(ctx, maxRetries, retryDelay, apiCall)
	if err != nil {
		resp.Diagnostics.AddError(
			"Create Error",
			fmt.Sprintf("Could not create resource, unexpected error: %s", err),
		)
		resp.Diagnostics.AddError(
			"inlineResp",
			fmt.Sprintf("%s", inlineResp),
		)

		if httpResp != nil {
			var responseBody string
			if httpResp != nil && httpResp.Body != nil {
				bodyBytes, readErr := io.ReadAll(httpResp.Body)
				if readErr == nil {
					responseBody = string(bodyBytes)
				}
			}
			tflog.Error(ctx, "Failed to create resource", map[string]interface{}{
				"error":          err.Error(),
				"httpStatusCode": httpResp.StatusCode,
				"responseBody":   responseBody,
			})
			resp.Diagnostics.AddError(
				"Error creating resource",
				fmt.Sprintf("HTTP Response: %v\nResponse Body: %s", httpResp, responseBody),
			)
			resp.Diagnostics.AddError(
				"Error",
				err.Error(),
			)
		}
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Create HTTP Client Failure",
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

	diags := updateDevicesResourceState(ctx, data, httpResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
	}

	// Now set the final state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Log that the resource was created.
	tflog.Trace(ctx, "created resource")
}

// Read method is responsible for reading an existing resource's state.
func (r *DevicesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *DevicesResourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	// Check if there are any errors before proceeding.
	if resp.Diagnostics.HasError() {
		return
	}

	// Initialize provider client and make API call
	maxRetries := r.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(r.client.GetConfig().Retry4xxErrorWaitTime)

	// API call function to be passed to retryOn4xx
	apiCall := func() (map[string]interface{}, *http.Response, error) {
		inline, httpResp, err := r.client.DevicesApi.GetDevice(ctx, data.Serial.ValueString()).Execute()
		return inline, httpResp, err
	}

	_, httpResp, err := utils.CustomHttpRequestRetry(ctx, maxRetries, retryDelay, apiCall)
	if err != nil {
		resp.Diagnostics.AddError(
			"Read Error",
			fmt.Sprintf("Could not read resource, unexpected error: %s", err),
		)

		if httpResp != nil {
			var responseBody string
			if httpResp != nil && httpResp.Body != nil {
				bodyBytes, readErr := io.ReadAll(httpResp.Body)
				if readErr == nil {
					responseBody = string(bodyBytes)
				}
			}
			tflog.Error(ctx, "Failed to read resource", map[string]interface{}{
				"error":          err.Error(),
				"httpStatusCode": httpResp.StatusCode,
				"responseBody":   responseBody,
			})
			resp.Diagnostics.AddError(
				"Error reading resource",
				fmt.Sprintf("HTTP Response: %v\nResponse Body: %s", httpResp, responseBody),
			)
			resp.Diagnostics.AddError(
				"Error",
				err.Error(),
			)
		}
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Read HTTP Client Failure",
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

	diags := updateDevicesResourceState(ctx, data, httpResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
	}

	// Now set the final state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Log that the resource was read.
	tflog.Trace(ctx, "read resource")
}

// Update function is responsible for updating the state of an existing resource.
func (r *DevicesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data *DevicesResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// If there was an error reading the plan, return early.
	if resp.Diagnostics.HasError() {
		return
	}

	// Create HTTP request body
	payload, payloadErr := updateDevicesResourcePayload(data)
	if payloadErr.HasError() {
		resp.Diagnostics.AddError("Failed to assemble payload", fmt.Sprintf("%s", payloadErr.Errors()))
	}

	// Initialize provider client and make API call
	maxRetries := r.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(r.client.GetConfig().Retry4xxErrorWaitTime)

	// API call function to be passed to retryOn4xx
	apiCall := func() (map[string]interface{}, *http.Response, error) {
		inline, httpResp, err := r.client.DevicesApi.UpdateDevice(context.Background(), data.Serial.ValueString()).UpdateDeviceRequest(payload).Execute()
		return inline, httpResp, err
	}

	_, httpResp, err := utils.CustomHttpRequestRetry(ctx, maxRetries, retryDelay, apiCall)
	if err != nil {
		resp.Diagnostics.AddError(
			"Update Error",
			fmt.Sprintf("Could not update resource, unexpected error: %s", err),
		)

		if httpResp != nil {
			var responseBody string
			if httpResp != nil && httpResp.Body != nil {
				bodyBytes, readErr := io.ReadAll(httpResp.Body)
				if readErr == nil {
					responseBody = string(bodyBytes)
				}
			}
			tflog.Error(ctx, "Failed to create resource", map[string]interface{}{
				"error":          err.Error(),
				"httpStatusCode": httpResp.StatusCode,
				"responseBody":   responseBody,
			})
			resp.Diagnostics.AddError(
				"Error updating resource",
				fmt.Sprintf("HTTP Response: %v\nResponse Body: %s", httpResp, responseBody),
			)
			resp.Diagnostics.AddError(
				"Error",
				err.Error(),
			)
		}
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Create HTTP Client Failure",
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

	diags := updateDevicesResourceState(ctx, data, httpResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
	}

	// Now set the updated state of the resource.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Log that the resource was updated.
	tflog.Trace(ctx, "updated resource")
}

// Delete function is responsible for deleting a resource.
func (r *DevicesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *DevicesResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	// If there was an error reading the state, return early.
	if resp.Diagnostics.HasError() {
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
	//var floor string

	updateDevice.Name = &name
	updateDevice.Tags = tags
	updateDevice.Lat = &lat
	updateDevice.Lng = &lng
	updateDevice.Address = &address
	updateDevice.Notes = &notes
	updateDevice.MoveMapMarker = &moveMapMarker
	//updateDevice.SwitchProfileId = &switchProfileId
	updateDevice.FloorPlanId = nil

	maxRetries := r.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(r.client.GetConfig().Retry4xxErrorWaitTime)

	// Initialize provider client and make API call

	// API call function to be passed to retryOn4xx
	apiCall := func() (map[string]interface{}, *http.Response, error) {
		inline, httpResp, err := r.client.DevicesApi.UpdateDevice(context.Background(), data.Serial.ValueString()).UpdateDeviceRequest(*updateDevice).Execute()
		return inline, httpResp, err
	}

	_, httpResp, err := utils.CustomHttpRequestRetry(ctx, maxRetries, retryDelay, apiCall)
	if err != nil {
		resp.Diagnostics.AddError(
			"Delete Error",
			fmt.Sprintf("Could not create group policy, unexpected error: %s", err),
		)

		if httpResp != nil {
			var responseBody string
			if httpResp != nil && httpResp.Body != nil {
				bodyBytes, readErr := io.ReadAll(httpResp.Body)
				if readErr == nil {
					responseBody = string(bodyBytes)
				}
			}
			tflog.Error(ctx, "Failed to create resource", map[string]interface{}{
				"error":          err.Error(),
				"httpStatusCode": httpResp.StatusCode,
				"responseBody":   responseBody,
			})
			resp.Diagnostics.AddError(
				"Error deleting resource",
				fmt.Sprintf("HTTP Response: %v\nResponse Body: %s", httpResp, responseBody),
			)
			resp.Diagnostics.AddError(
				"Error",
				err.Error(),
			)
		}
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
func (r *DevicesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	// Pass through the ID directly from the ImportStateRequest to the ImportStateResponse
	resource.ImportStatePassthroughID(ctx, path.Root("serial"), req, resp)

}
