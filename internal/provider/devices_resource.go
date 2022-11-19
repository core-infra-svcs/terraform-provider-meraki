package provider

import (
	"context"
	"encoding/json"
	"fmt"
	apiclient "github.com/core-infra-svcs/dashboard-api-go/client"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"strconv"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &DevicesResource{}
var _ resource.ResourceWithImportState = &DevicesResource{}

func NewDevicesResource() resource.Resource {
	return &DevicesResource{}
}

// DevicesResource defines the resource implementation.
type DevicesResource struct {
	client *apiclient.APIClient
}

// DevicesResourceModel describes the resource data model.
type DevicesResourceModel struct {
	Id                  types.String `tfsdk:"id"`
	Serial              types.String `tfsdk:"serial"`
	Name                types.String `tfsdk:"name"`
	Mac                 types.String `tfsdk:"mac"`
	Model               types.String `tfsdk:"model"`
	Tags                types.Set    `tfsdk:"tags"`
	LanIp               types.String `tfsdk:"lan_ip"`
	Firmware            types.String `tfsdk:"firmware"`
	Lat                 types.String `tfsdk:"lat"`
	Lng                 types.String `tfsdk:"lng"`
	Address             types.String `tfsdk:"address"`
	Notes               types.String `tfsdk:"notes"`
	Url                 types.String `tfsdk:"url"`
	Wan1Ip              types.String `tfsdk:"wan1ip"`
	Wan2Ip              types.String `tfsdk:"wan2ip"`
	MoveMapMarker       types.Bool   `tfsdk:"move_map_marker"`
	FloorPlanId         types.String `tfsdk:"floor_plan_id"`
	NetworkId           types.String `tfsdk:"network_id"`
	BeaconIdParamsUuid  types.String `tfsdk:"beacon_id_params_uuid"`
	BeaconIdParamsMajor types.Int64  `tfsdk:"beacon_id_params_major"`
	BeaconIdParamsMinor types.Int64  `tfsdk:"beacon_id_params_minor"`
	SwitchProfileId     types.String `tfsdk:"switch_profile_id"`
}

func (r *DevicesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_devices"
}

func (r *DevicesResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "Devices resource - Update the attributes of a device",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Computed:            true,
				Type:                types.StringType,
				MarkdownDescription: "Example identifier",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
			},
			"serial": {
				Description:         "The devices serial number",
				MarkdownDescription: "The devices serial number",
				Type:                types.StringType,
				Required:            true,
				Optional:            false,
				Computed:            false,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"name": {
				Description:         "The name of a device",
				MarkdownDescription: "The name of a device",
				Type:                types.StringType,
				Required:            false,
				Optional:            true,
				Computed:            true,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"tags": {
				Description:         "The list of tags of a device",
				MarkdownDescription: "The list of tags of a device",
				Type:                types.SetType{ElemType: types.StringType},
				Required:            false,
				Optional:            true,
				Computed:            true,
				Sensitive:           false,
				Attributes:          nil,
				// tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
				//					"tag": {
				//						Type:     types.StringType,
				//						Computed: true,
				//						Optional: true,
				//					},
				//				})
			},
			"lat": {
				Description:         "The latitude of a device",
				MarkdownDescription: "The latitude of a device",
				Type:                types.StringType,
				Required:            false,
				Optional:            true,
				Computed:            true,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"lng": {
				Description:         "The longitude of a device",
				MarkdownDescription: "The longitude of a device",
				Type:                types.StringType,
				Required:            false,
				Optional:            true,
				Computed:            true,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"address": {
				Description:         "The address of a device",
				MarkdownDescription: "The address of a device",
				Type:                types.StringType,
				Required:            false,
				Optional:            true,
				Computed:            true,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"notes": {
				Description:         "The notes for the device. String. Limited to 255 characters.",
				MarkdownDescription: "The notes for the device. String. Limited to 255 characters.",
				Type:                types.StringType,
				Required:            false,
				Optional:            true,
				Computed:            true,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"move_map_marker": {
				Description:         "Whether or not to set the latitude and longitude of a device based on the new address. Only applies when lat and lng are not specified.",
				MarkdownDescription: "Whether or not to set the latitude and longitude of a device based on the new address. Only applies when lat and lng are not specified.",
				Type:                types.BoolType,
				Required:            false,
				Optional:            true,
				Computed:            true,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"switch_profile_id": {
				Description:         "The ID of a switch profile to bind to the device (for available switch profiles, see the 'Switch Profiles' endpoint). Use null to unbind the switch device from the current profile. For a device to be bindable to a switch profile, it must (1) be a switch, and (2) belong to a network that is bound to a configuration template.",
				MarkdownDescription: "The ID of a switch profile to bind to the device (for available switch profiles, see the 'Switch Profiles' endpoint). Use null to unbind the switch device from the current profile. For a device to be bindable to a switch profile, it must (1) be a switch, and (2) belong to a network that is bound to a configuration template.",
				Type:                types.StringType,
				Required:            false,
				Optional:            true,
				Computed:            true,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"floor_plan_id": {
				Description:         "The floor plan to associate to this device. null disassociates the device from the floor plan.",
				MarkdownDescription: "The floor plan to associate to this device. null disassociates the device from the floor plan.",
				Type:                types.StringType,
				Required:            false,
				Optional:            true,
				Computed:            true,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"mac": {
				Description:         "The mac address of a device",
				MarkdownDescription: "The mac address of a device",
				Type:                types.StringType,
				Required:            false,
				Optional:            false,
				Computed:            true,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"model": {
				Description:         "The model of a device",
				MarkdownDescription: "The model of a device",
				Type:                types.StringType,
				Required:            false,
				Optional:            false,
				Computed:            true,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"lan_ip": {
				Description:         "The ipv4 lan ip of a device",
				MarkdownDescription: "The  ipv4 lan ip of a device",
				Type:                types.StringType,
				Required:            false,
				Optional:            false,
				Computed:            true,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"firmware": {
				Description:         "The firmware version of a device",
				MarkdownDescription: "The firmware version of a device",
				Type:                types.StringType,
				Required:            false,
				Optional:            false,
				Computed:            true,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"url": {
				Description:         "The url for the network associated with the device.",
				MarkdownDescription: "The url for the network associated with the device.",
				Type:                types.StringType,
				Required:            false,
				Optional:            false,
				Computed:            true,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"wan1ip": {
				Description:         "IP of Wan interface 1",
				MarkdownDescription: "IP of Wan interface 1",
				Type:                types.StringType,
				Required:            false,
				Optional:            false,
				Computed:            true,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"wan2ip": {
				Description:         "IP of Wan interface 2",
				MarkdownDescription: "IP of Wan interface 2",
				Type:                types.StringType,
				Required:            false,
				Optional:            false,
				Computed:            true,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"network_id": {
				Description:         "The networkId of a device",
				MarkdownDescription: "The networkId of a device",
				Type:                types.StringType,
				Required:            false,
				Optional:            false,
				Computed:            true,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"beacon_id_params_uuid": {
				Description:         "The beacon id params uuid of a device",
				MarkdownDescription: "The name of a device",
				Type:                types.StringType,
				Required:            false,
				Optional:            false,
				Computed:            true,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"beacon_id_params_major": {
				Description:         "The beacon id params major of a device",
				MarkdownDescription: "The beacon id params major of a device",
				Type:                types.Int64Type,
				Required:            false,
				Optional:            false,
				Computed:            true,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"beacon_id_params_minor": {
				Description:         "The beacon id params minor of a device",
				MarkdownDescription: "The beacon id params minor of a device",
				Type:                types.Int64Type,
				Required:            false,
				Optional:            false,
				Computed:            true,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
		},
	}, nil
}

func (r *DevicesResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*apiclient.APIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *DevicesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *DevicesResourceModel

	// Read Terraform plan data
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check for required parameters
	if len(data.Serial.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing device serial number", fmt.Sprintf("S/N: %s", data.Serial.ValueString()))
		return
	}

	// Initialize provider client and make API call
	inlineResp, httpResp, err := r.client.DevicesApi.GetDevice(context.Background(), data.Serial.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// Check for API success response code
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// collect diagnostics
	tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%s", data))
		return
	}

	// Save data into Terraform state
	extractHttpResponseDevicesResource(ctx, inlineResp, data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "created resource")
}

func (r *DevicesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *DevicesResourceModel

	// Read Terraform state data
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check for required parameters
	if len(data.Serial.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing device serial number", fmt.Sprintf("S/N: %s", data.Serial.ValueString()))
		return
	}

	// Initialize provider client and make API call
	inlineResp, httpResp, err := r.client.DevicesApi.GetDevice(context.Background(), data.Serial.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// Check for API success inlineResp code
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// collect diagnostics
	tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("State Data", fmt.Sprintf("\n%s", data))
		return
	}

	// Save data into Terraform state
	extractHttpResponseDevicesResource(ctx, inlineResp, data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *DevicesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *DevicesResourceModel

	// Read Terraform plan data
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check for required parameters
	if len(data.Serial.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing device serial number", fmt.Sprintf("S/N: %s", data.Serial.ValueString()))
		return
	}

	// Create HTTP request body
	updateDevice := apiclient.NewInlineObject()

	// Name
	updateDevice.SetName(data.Name.ValueString())

	// address
	updateDevice.SetAddress(data.Address.ValueString())

	// Tags
	if data.Tags.IsNull() != true {
		var tags []string
		for _, attribute := range data.Tags.Elements() {
			tags = append(tags, attribute.String())
		}
		updateDevice.SetTags(tags)
	}

	//	Lat
	lat, _ := strconv.ParseFloat(data.Lat.ValueString(), 32)
	updateDevice.SetLat(float32(lat))

	//	Lng
	lng, _ := strconv.ParseFloat(data.Lng.ValueString(), 32)
	updateDevice.SetLng(float32(lng))

	//	Notes
	updateDevice.SetNotes(data.Notes.ValueString())

	//	MoveMapMarker
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
	inlineResp, httpResp, err := r.client.DevicesApi.UpdateDevice(context.Background(),
		data.Serial.ValueString()).UpdateDevice(*updateDevice).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// Check for API success response code
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// collect diagnostics
	tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%s", data))
		return
	}

	// Save data into Terraform state
	extractHttpResponseDevicesResource(ctx, inlineResp, data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *DevicesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *DevicesResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check for required parameters
	if len(data.Serial.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing device serial number", fmt.Sprintf("S/N: %s", data.Serial.ValueString()))
		return
	}

	/*
		// Create HTTP request body
			deleteDevice := apiclient.NewInlineObject()

			// Name
			deleteDevice.SetAddress("")

			//	Tags
			var tags []string
			deleteDevice.SetTags(tags)

			//	Lat
			deleteDevice.SetLat(0)

			//	Lng
			deleteDevice.SetLng(0)

			//	Address
			deleteDevice.SetAddress("")

			//	Notes
			deleteDevice.SetNotes("")

			//	MoveMapMarker
			deleteDevice.SetMoveMapMarker(false)

			// SwitchProfileId
			if data.SwitchProfileId.IsNull() != true {
				deleteDevice.SetSwitchProfileId(types.StringNull().ValueString())
			}

			//	FloorPlanId
			if data.FloorPlanId.IsNull() != true {
				deleteDevice.SetFloorPlanId(types.StringNull().ValueString())
			}

			// Initialize provider client and make API call
			_, httpResp, err := r.client.DevicesApi.UpdateDevice(context.Background(),
				data.Serial.ValueString()).UpdateDevice(*deleteDevice).Execute()
			if err != nil {
				resp.Diagnostics.AddError(
					"Failed to delete resource",
					fmt.Sprintf("%v\n", err.Error()),
				)
			}

			// Check for API success response code
			if httpResp.StatusCode != 200 {
				resp.Diagnostics.AddError(
					"Unexpected HTTP Response Status Code",
					fmt.Sprintf("%v", httpResp.StatusCode),
				)
			}

			// collect diagnostics
			tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)

			// Check for errors after diagnostics collected
			if resp.Diagnostics.HasError() {
				resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%s", data))
				return
			}
	*/

	// Remove from state
	resp.State.RemoveResource(ctx)

	// Write logs using the tflog package
	tflog.Trace(ctx, "delete resource")
}

func extractHttpResponseDevicesResource(ctx context.Context, inlineResp map[string]interface{}, data *DevicesResourceModel) *DevicesResourceModel {

	// save into the Terraform state
	data.Id = types.StringValue("example-id")

	// address attribute
	if address := inlineResp["address"]; address != nil {
		data.Address = types.StringValue(address.(string))
	} else {
		data.Address = types.StringNull()
	}

	// firmware attribute
	if firmware := inlineResp["firmware"]; firmware != nil {
		data.Firmware = types.StringValue(firmware.(string))
	} else {
		data.Firmware = types.StringNull()
	}

	// mac attribute
	if mac := inlineResp["mac"]; mac != nil {
		data.Mac = types.StringValue(mac.(string))
	} else {
		data.Mac = types.StringNull()
	}

	// url attribute
	if url := inlineResp["url"]; url != nil {
		data.Url = types.StringValue(url.(string))
	} else {
		data.Url = types.StringNull()
	}

	// name attribute
	if name := inlineResp["name"]; name != nil {
		data.Name = types.StringValue(name.(string))
	} else {
		data.Name = types.StringNull()
	}

	// model attribute
	if model := inlineResp["model"]; model != nil {
		data.Model = types.StringValue(model.(string))
	} else {
		data.Model = types.StringNull()
	}

	// networkId attribute
	if networkId := inlineResp["networkId"]; networkId != nil {
		data.NetworkId = types.StringValue(networkId.(string))
	} else {
		data.NetworkId = types.StringNull()
	}

	// tags attribute
	if tags := inlineResp["tags"]; tags != nil {

		// format tags
		var sortedTagList []string
		for _, v := range tags.([]interface{}) {
			var s string
			_ = json.Unmarshal([]byte(v.(string)), &s)
			sortedTagList = append(sortedTagList, s)
		}

		// append to tag list
		var tagList []attr.Value
		for _, i := range sortedTagList {
			tag := types.StringValue(i)
			tagList = append(tagList, tag)
		}

		/*
			var tagList []attr.Value

				for _, v := range tags.([]interface{}) {
					var s string
					_ = json.Unmarshal([]byte(v.(string)), &s)
					tag := types.StringValue(s)
					tagList = append(tagList, tag)
				}
		*/

		data.Tags, _ = types.SetValue(types.StringType, tagList)
	} else {
		data.Tags = types.SetNull(types.StringType)
	}

	// floorPlanId attribute
	if floorPlanId := inlineResp["floorPlanId"]; floorPlanId != nil {
		data.FloorPlanId = types.StringValue(floorPlanId.(string))
	} else {
		data.FloorPlanId = types.StringNull()
	}

	// lat attribute
	if lat := inlineResp["lat"]; lat != nil {
		data.Lat = types.StringValue(fmt.Sprintf("%v", lat.(float64)))
	} else {
		data.Lat = types.StringNull()
	}

	// lng attribute
	if lng := inlineResp["lng"]; lng != nil {
		data.Lng = types.StringValue(fmt.Sprintf("%v", lng.(float64)))
	} else {
		data.Lng = types.StringNull()
	}

	// notes attribute
	if notes := inlineResp["notes"]; notes != nil {
		data.Notes = types.StringValue(notes.(string))
	} else {
		data.Notes = types.StringNull()
	}

	// switchProfileId attribute
	if switchProfileId := inlineResp["switchProfileId"]; switchProfileId != nil {
		data.SwitchProfileId = types.StringValue(switchProfileId.(string))
	} else {
		data.SwitchProfileId = types.StringNull()
	}

	// beaconIdParams attribute
	if beaconIdParams := inlineResp["beaconIdParams"]; beaconIdParams != nil {
		uuid := beaconIdParams.(map[string]interface{})
		data.BeaconIdParamsUuid = types.StringValue(uuid["uuid"].(string))

		major := beaconIdParams.(map[string]interface{})
		data.BeaconIdParamsMajor = types.Int64Value(major["major"].(int64))

		minor := beaconIdParams.(map[string]interface{})
		data.BeaconIdParamsMajor = types.Int64Value(minor["minor"].(int64))

	} else {
		data.BeaconIdParamsUuid = types.StringNull()
		data.BeaconIdParamsMajor = types.Int64Null()
		data.BeaconIdParamsMinor = types.Int64Null()
	}

	// moveMapMarker attribute
	if moveMapMarker := inlineResp["moveMapMarker"]; moveMapMarker != nil {
		data.MoveMapMarker = types.BoolValue(moveMapMarker.(bool))
	} else {
		data.MoveMapMarker = types.BoolNull()
	}

	// wan1Ip attribute
	if wan1Ip := inlineResp["wan1Ip"]; wan1Ip != nil {
		data.Wan1Ip = types.StringValue(wan1Ip.(string))
	} else {
		data.Wan1Ip = types.StringNull()
	}

	// wan2Ip attribute
	if wan2Ip := inlineResp["wan2Ip"]; wan2Ip != nil {
		data.Wan2Ip = types.StringValue(wan2Ip.(string))
	} else {
		data.Wan2Ip = types.StringNull()
	}

	// lanIp attribute (computed)
	if lanIp := inlineResp["lanIp"]; lanIp != nil {
		data.LanIp = types.StringValue(lanIp.(string))
	} else {
		data.LanIp = types.StringNull()
	}

	return data
}

func (r *DevicesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
