package provider

import (
	"context"
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
	Tags                types.List   `tfsdk:"tags"`
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
				MarkdownDescription: "Example identifier",
				PlanModifiers: tfsdk.AttributePlanModifiers{
					resource.UseStateForUnknown(),
				},
				Type: types.StringType,
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
				Type:                types.ListType{ElemType: types.StringType},
				Required:            false,
				Optional:            true,
				Computed:            true,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
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
		return
	}

	// save into the Terraform state
	data.Id = types.StringValue("example-id")

	// name attribute
	if name := inlineResp["name"]; name != nil {
		data.Name = types.StringValue(name.(string))
	} else {
		data.Name = types.StringNull()
	}

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

	// floor plan attribute
	if floorplan := inlineResp["floorplan"]; floorplan != nil {
		data.FloorPlanId = types.StringValue(floorplan.(string))
	} else {
		data.FloorPlanId = types.StringNull()
	}

	// lat attribute
	if lat := inlineResp["lat"]; lat != nil {
		data.Lat = types.StringValue(fmt.Sprintf("%f", lat.(float64)))
	} else {
		data.Lat = types.StringNull()
	}

	// lng attribute
	if lng := inlineResp["lng"]; lng != nil {
		data.Lng = types.StringValue(fmt.Sprintf("%f", lng.(float64)))
	} else {
		data.Lng = types.StringNull()
	}

	// mac attribute
	if mac := inlineResp["mac"]; mac != nil {
		data.Mac = types.StringValue(mac.(string))
	} else {
		data.Mac = types.StringNull()
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

	// serial number is not computed or optional

	// tags attribute
	if tags := inlineResp["tags"]; tags != nil {

		// append tags to tag list
		var tagElements []attr.Value // list of tags
		for _, v := range inlineResp["tags"].([]interface{}) {
			tagElements = append(tagElements, types.StringValue(v.(string)))
		}
		data.Tags, _ = types.ListValue(types.StringType, tagElements)
	} else {
		data.Tags = types.ListNull(types.StringType)
	}

	// networkId attribute
	if networkId := inlineResp["networkId"]; networkId != nil {
		data.NetworkId = types.StringValue(networkId.(string))
	} else {
		data.NetworkId = types.StringNull()
	}

	// url attribute
	if url := inlineResp["url"]; url != nil {
		data.Url = types.StringValue(url.(string))
	} else {
		data.Url = types.StringNull()
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

	// lanIp attribute (unknown)
	data.LanIp = types.StringNull()

	// lanIp attribute (unknown)
	data.Notes = types.StringNull()

	// lanIp attribute (unknown)
	data.SwitchProfileId = types.StringNull()

	// lanIp attribute (unknown)
	data.MoveMapMarker = types.BoolNull()

	// uuid attribute (unknown)
	data.BeaconIdParamsUuid = types.StringNull()

	// major attribute  (unknown)
	data.BeaconIdParamsMajor = types.Int64Null()

	// minor attribute  (unknown)
	data.BeaconIdParamsMinor = types.Int64Null()

	// Save data into Terraform state
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
		return
	}

	// save inlineResp data into Terraform state
	data.Id = types.StringValue("example-id")

	// TODO - make reusable function with error handling,inputs: mapstring interface, key name, and tfsdk.type.
	// name attribute
	if name := inlineResp["name"]; name != nil {
		data.Name = types.StringValue(name.(string))
	} else {
		data.Name = types.StringNull()
	}

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

	// floor plan attribute
	if floorplan := inlineResp["floorplan"]; floorplan != nil {
		data.FloorPlanId = types.StringValue(floorplan.(string))
	} else {
		data.FloorPlanId = types.StringNull()
	}

	// lat attribute
	if lat := inlineResp["lat"]; lat != nil {
		data.Lat = types.StringValue(fmt.Sprintf("%f", lat.(float64)))
	} else {
		data.Lat = types.StringNull()
	}

	// lng attribute
	if lng := inlineResp["lng"]; lng != nil {
		data.Lng = types.StringValue(fmt.Sprintf("%f", lng.(float64)))
	} else {
		data.Lng = types.StringNull()
	}

	// mac attribute
	if mac := inlineResp["mac"]; mac != nil {
		data.Mac = types.StringValue(mac.(string))
	} else {
		data.Mac = types.StringNull()
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

	// serial number is not computed or optional

	// tags attribute
	if tags := inlineResp["tags"]; tags != nil {

		// append tags to tag list
		var tagElements []attr.Value // list of tags
		for _, v := range inlineResp["tags"].([]interface{}) {
			tagElements = append(tagElements, types.StringValue(v.(string)))
		}
		data.Tags, _ = types.ListValue(types.StringType, tagElements)
	} else {
		data.Tags = types.ListNull(types.StringType)
	}

	// networkId attribute
	if networkId := inlineResp["networkId"]; networkId != nil {
		data.NetworkId = types.StringValue(networkId.(string))
	} else {
		data.NetworkId = types.StringNull()
	}

	// url attribute
	if url := inlineResp["url"]; url != nil {
		data.Url = types.StringValue(url.(string))
	} else {
		data.Url = types.StringNull()
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

	// TODO - Set defaults to unknown
	// lanIp attribute (unknown)
	data.LanIp = types.StringUnknown()

	// lanIp attribute (unknown)
	data.Notes = types.StringNull()

	// lanIp attribute (unknown)
	data.SwitchProfileId = types.StringNull()

	// lanIp attribute (unknown)
	data.MoveMapMarker = types.BoolNull()

	// uuid attribute (unknown)
	data.BeaconIdParamsUuid = types.StringNull()

	// major attribute  (unknown)
	data.BeaconIdParamsMajor = types.Int64Null()

	// minor attribute  (unknown)
	data.BeaconIdParamsMinor = types.Int64Null()

	// Save updated data into Terraform state
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
	updateDevice.SetAddress(data.Name.ValueString())

	// Tags
	if data.Tags.Elements() != nil {
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

	//	Address
	updateDevice.SetAddress(data.Address.ValueString())

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
		return
	}

	// save inlineResp data into Terraform state
	data.Id = types.StringValue("example-id")

	/*
		resp.Diagnostics.AddError(
				"Response Data",
				fmt.Sprintf("\n%s", inlineResp),
			)

			// Check for errors after diagnostics collected
			if resp.Diagnostics.HasError() {
				return
			}
	*/

	/*
		// name attribute
			if name := inlineResp["name"]; name != nil {
				data.Name = types.StringValue(name.(string))
			} else if data.Name.IsNull() != true {
			} else {
				data.Name = types.StringNull()
			}


		// TODO - Treat this more like a list
		// tags attribute
		if tags := inlineResp["tags"]; tags != nil {
			// append tags to tag list
			var tagElements []attr.Value // list of tags
			for _, v := range inlineResp["tags"].([]interface{}) {
				tagElements = append(tagElements, types.StringValue(v.(string)))
			}
			data.Tags, _ = types.ListValue(types.StringType, tagElements)
		} else if data.Tags.IsNull() != true {
		} else {
			data.Tags = types.ListNull(types.StringType)
		}

	*/

	// address attribute
	if address := inlineResp["address"]; address != nil {

		data.Address = types.StringValue(address.(string))
	} else if data.Address.IsNull() != true {
	} else {
		data.Address = types.StringNull()
	}

	// lat attribute
	if latResp := inlineResp["lat"]; latResp != nil {
		data.Lat = types.StringValue(fmt.Sprintf("%f", latResp.(float64)))
	} else if data.Lat.IsNull() != true {
	} else {
		data.Lat = types.StringNull()
	}

	// lng attribute
	if lngResp := inlineResp["lng"]; lngResp != nil {
		data.Lng = types.StringValue(fmt.Sprintf("%f", lngResp.(float64)))
	} else if data.Lng.IsNull() != true {
	} else {
		data.Lng = types.StringNull()
	}

	// notes attribute
	if notes := inlineResp["notes"]; notes != nil {
		data.Notes = types.StringValue(notes.(string))
	} else if data.Notes.IsNull() != true {
		data.Notes = types.StringNull()
	} else {
		data.Notes = types.StringNull()
	}

	// moveMapMarker attribute
	if moveMapMarker := inlineResp["moveMapMarker"]; moveMapMarker != nil {
		data.MoveMapMarker = types.BoolValue(moveMapMarker.(bool))
	} else if data.MoveMapMarker.IsNull() != true {
		data.MoveMapMarker = types.BoolNull()
	} else {
		data.MoveMapMarker = types.BoolNull()
	}

	// switchProfileId attribute
	if switchProfileId := inlineResp["switchProfileId"]; switchProfileId != nil {
		data.SwitchProfileId = types.StringValue(switchProfileId.(string))
	} else if data.SwitchProfileId.IsNull() != true {
		data.SwitchProfileId = types.StringNull()
	} else {
		data.SwitchProfileId = types.StringNull()
	}

	// floor plan attribute
	if floorPlan := inlineResp["floorPlan"]; floorPlan != nil {
		data.FloorPlanId = types.StringValue(floorPlan.(string))
	} else if data.FloorPlanId.IsNull() != true {
		data.FloorPlanId = types.StringNull()
	} else {
		data.FloorPlanId = types.StringNull()
	}

	// firmware attribute
	if firmware := inlineResp["firmware"]; firmware != nil {
		data.Firmware = types.StringValue(firmware.(string))
	} else if data.Firmware.IsNull() != true {
	} else {
		data.Firmware = types.StringNull()
	}

	// mac attribute
	if mac := inlineResp["mac"]; mac != nil {
		data.Mac = types.StringValue(mac.(string))
	} else if data.Mac.IsNull() != true {
	} else {
		data.Mac = types.StringNull()
	}

	// model attribute
	if model := inlineResp["model"]; model != nil {
		data.Model = types.StringValue(model.(string))
	} else if data.Model.IsNull() != true {
	} else {
		data.Model = types.StringNull()
	}

	// networkId attribute
	if networkId := inlineResp["networkId"]; networkId != nil {
		data.NetworkId = types.StringValue(networkId.(string))
	} else if data.NetworkId.IsNull() != true {
	} else {
		data.NetworkId = types.StringNull()
	}

	// serial number is not computed or optional

	// networkId attribute
	if networkId := inlineResp["networkId"]; networkId != nil {
		data.NetworkId = types.StringValue(networkId.(string))
	} else if data.NetworkId.IsNull() != true {
		data.NetworkId = types.StringNull()
	} else {
		data.NetworkId = types.StringNull()
	}

	// url attribute
	if url := inlineResp["url"]; url != nil {
		data.Url = types.StringValue(url.(string))
	} else if data.Url.IsNull() != true {
		data.Url = types.StringNull()
	} else {
		data.Url = types.StringNull()
	}

	// wan1Ip attribute
	if wan1Ip := inlineResp["wan1Ip"]; wan1Ip != nil {
		data.Wan1Ip = types.StringValue(wan1Ip.(string))
	} else if data.Wan1Ip.ValueString() != "<nil>" {
		data.Wan1Ip = types.StringNull()
	} else {
		data.Wan1Ip = types.StringNull()
	}

	// wan2Ip attribute
	if wan2Ip := inlineResp["wan2Ip"]; wan2Ip != nil {
		data.Wan2Ip = types.StringValue(wan2Ip.(string))
	} else if data.Wan2Ip.IsNull() != true {
		data.Wan2Ip = types.StringNull()
	} else {
		data.Wan2Ip = types.StringNull()
	}

	// lanIp attribute (computed)
	if lanIp := inlineResp["lanIp"]; lanIp != nil {
		data.LanIp = types.StringValue(lanIp.(string))
	} else if data.LanIp.IsNull() != true {
		data.LanIp = types.StringNull()
	} else {
		data.LanIp = types.StringNull()
	}

	// beaconIdParamsUuid attribute (computed)
	if beaconIdParamsUuid := inlineResp["beaconIdParamsUuid"]; beaconIdParamsUuid != nil {
		data.BeaconIdParamsUuid = types.StringValue(beaconIdParamsUuid.(string))
	} else if data.BeaconIdParamsUuid.IsNull() != true {
		data.BeaconIdParamsUuid = types.StringNull()
	} else {
		data.BeaconIdParamsUuid = types.StringNull()
	}

	// beaconIdParamsMajor attribute (computed)
	if beaconIdParamsMajor := inlineResp["beaconIdParamsMajor"]; beaconIdParamsMajor != nil {
		data.BeaconIdParamsMajor = types.Int64Value(beaconIdParamsMajor.(int64))
	} else if data.BeaconIdParamsMajor.IsNull() != true {
		data.BeaconIdParamsMajor = types.Int64Null()
	} else {
		data.BeaconIdParamsMajor = types.Int64Null()
	}

	// beaconIdParamsMinor attribute (computed)
	if beaconIdParamsMinor := inlineResp["beaconIdParamsMinor"]; beaconIdParamsMinor != nil {
		data.BeaconIdParamsMinor = types.Int64Value(beaconIdParamsMinor.(int64))
	} else if data.BeaconIdParamsMinor.IsNull() != true {
		data.BeaconIdParamsMinor = types.Int64Null()
	} else {
		data.BeaconIdParamsMinor = types.Int64Null()
	}

	// Save updated data into Terraform state
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
		return
	}

	// Remove from state
	resp.State.RemoveResource(ctx)

	// Write logs using the tflog package
	tflog.Trace(ctx, "delete resource")
}

func (r *DevicesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
