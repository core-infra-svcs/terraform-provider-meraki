package devices

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"
	utils2 "github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &NetworkDevicesDataSource{}

func NewNetworkDevicesDataSource() datasource.DataSource {
	return &NetworkDevicesDataSource{}
}

// NetworkDevicesDataSource defines the data source implementation.
type NetworkDevicesDataSource struct {
	client *openApiClient.APIClient
}

// The DevicesDatasourceModel structure describes the data model.
// This struct is where you define all the attributes that are part of this resource's state.
type DevicesDatasourceModel struct {
	Id        types.String `tfsdk:"id"`
	List      types.List   `tfsdk:"list"`
	NetworkId types.String `tfsdk:"network_id"`
}

type DevicesDatasourceModelDevice struct {
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

type DevicesDatasourceModelBeaconIdParams struct {
	Uuid  jsontypes.String `tfsdk:"uuid"`
	Major jsontypes.Int64  `tfsdk:"major"`
	Minor jsontypes.Int64  `tfsdk:"minor"`
}

func (d *NetworkDevicesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_devices"
}

func (d *NetworkDevicesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Returns the list of network devices",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Network ID",
				Required:            true,
			},
			"list": schema.ListNestedAttribute{
				MarkdownDescription: "List of devices",
				Optional:            true,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"network_id": schema.StringAttribute{
							MarkdownDescription: "Network ID",
							Computed:            true,
						},
						"serial": schema.StringAttribute{
							MarkdownDescription: "The devices serial number",
							Computed:            true,
							Validators: []validator.String{
								stringvalidator.LengthBetween(8, 31),
							},
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of a device",
							Computed:            true,
						},
						"tags": schema.ListAttribute{
							Description: "Network tags",
							Computed:    true,
							ElementType: types.StringType,
						},
						"lat": schema.Float64Attribute{
							MarkdownDescription: "The latitude of a device",
							Computed:            true,
						},
						"lng": schema.Float64Attribute{
							MarkdownDescription: "The longitude of a device",
							Computed:            true,
						},
						"address": schema.StringAttribute{
							MarkdownDescription: "The address of a device",
							Computed:            true,
						},
						"notes": schema.StringAttribute{
							MarkdownDescription: "Notes for the network",
							Computed:            true,
						},
						"details": schema.ListNestedAttribute{
							Description: "Network tags",
							Computed:    true,
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
						},
						"move_map_marker": schema.BoolAttribute{
							MarkdownDescription: "Whether or not to set the latitude and longitude of a device based on the new address. Only applies when lat and lng are not specified.",
							Computed:            true,
						},
						"switch_profile_id": schema.StringAttribute{
							MarkdownDescription: "The ID of a switch profile to bind to the device (for available switch profiles, see the 'Switch Profiles' endpoint). Use null to unbind the switch device from the current profile. For a device to be bindable to a switch profile, it must (1) be a switch, and (2) belong to a network that is bound to a configuration template.",
							Computed:            true,
						},
						"floor_plan_id": schema.StringAttribute{
							MarkdownDescription: "The floor plan to associate to this device. null disassociates the device from the floor plan.",
							Computed:            true,
						},
						"mac": schema.StringAttribute{
							MarkdownDescription: "The mac address of a device",
							Computed:            true,
						},
						"model": schema.StringAttribute{
							MarkdownDescription: "The model of a device",
							Computed:            true,
						},
						"lan_ip": schema.StringAttribute{
							MarkdownDescription: "The ipv4 lan ip of a device",
							Computed:            true,
						},
						"firmware": schema.StringAttribute{
							MarkdownDescription: "The firmware version of a device",
							Computed:            true,
						},
						"url": schema.StringAttribute{
							MarkdownDescription: "The url for the network associated with the device.",
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
				},
			},
		},
	}
}

func (d *NetworkDevicesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*openApiClient.APIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

/*

func updateDevicesDatasourceStateList(ctx context.Context, state *DevicesDatasourceModel, httpResp *http.Response) diag.Diagnostics {
	var diags diag.Diagnostics
	var device DevicesDatasourceModelDevice



	if device.BeaconIdParams.IsUnknown() {
		device.BeaconIdParams = types.ObjectNull(DevicesResourceModelBeaconIdParamsAttrTypes())
	}

	if state.Firmware.IsUnknown() {
		state.Firmware = jsontypes.StringNull()
	}
	if state.FloorPlanId.IsUnknown() {
		state.FloorPlanId = jsontypes.StringNull()
	}
	if state.LanIp.IsUnknown() {
		state.LanIp = jsontypes.StringNull()
	}
	if state.Mac.IsUnknown() {
		state.Mac = jsontypes.StringNull()
	}
	if state.Url.IsUnknown() {
		state.Url = jsontypes.StringNull()
	}
	if state.Model.IsUnknown() {
		state.Model = jsontypes.StringNull()
	}
	if state.SwitchProfileId.IsUnknown() {
		state.SwitchProfileId = jsontypes.StringNull()
	}
	if state.MoveMapMarker.IsUnknown() {
		state.MoveMapMarker = jsontypes.BoolNull()
	}

	if state.Lat.IsUnknown() {
		state.Lat = jsontypes.Float64Null()
	}

	if state.Lng.IsUnknown() {
		state.Lng = jsontypes.Float64Null()
	}

	if state.Notes.IsUnknown() {
		state.Notes = jsontypes.StringNull()
	}

	if state.Name.IsUnknown() {
		state.Name = jsontypes.StringNull()
	}

	if state.Address.IsUnknown() {
		state.Address = jsontypes.StringNull()
	}

	// Set ID for the resource.
	state.Id = jsontypes.StringValue(state.Serial.ValueString())

	state.List = append(state.List, device)

	return diags
}

*/

// updateDevicesResourceState updates the resource state with the provided api data.
func updateDevicesDatasourceState(ctx context.Context, state *DevicesDatasourceModel, inlineResp []map[string]interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	var deviceList []attr.Value

	devicesDatasourceModelBeaconIdParamsAttrs := map[string]attr.Type{
		"uuid":  types.StringType,
		"major": types.Int64Type,
		"minor": types.Int64Type,
	}

	devicesDatasourceModelDetailAttr := map[string]attr.Type{
		"name":  types.StringType,
		"value": types.StringType,
	}

	devicesDatasourceModelDeviceAttrs := map[string]attr.Type{
		"id":                types.StringType,
		"serial":            types.StringType,
		"name":              types.StringType,
		"mac":               types.StringType,
		"model":             types.StringType,
		"tags":              types.ListType{ElemType: types.StringType}, // Assuming tags are strings
		"details":           types.ListType{ElemType: types.ObjectType{AttrTypes: devicesDatasourceModelDetailAttr}},
		"lan_ip":            types.StringType,
		"firmware":          types.StringType,
		"lat":               types.Float64Type,
		"lng":               types.Float64Type,
		"address":           types.StringType,
		"notes":             types.StringType,
		"url":               types.StringType,
		"floor_plan_id":     types.StringType,
		"network_id":        types.StringType,
		"beacon_id_params":  types.ObjectType{AttrTypes: devicesDatasourceModelBeaconIdParamsAttrs},
		"switch_profile_id": types.StringType,
		"move_map_marker":   types.BoolType,
	}

	// Set ID for the new resource.
	if state.Id.IsNull() || state.Id.IsUnknown() {
		state.Id = state.NetworkId
	}

	for _, d := range inlineResp {
		var device DevicesDatasourceModelDevice

		// "name": "My AP",
		name, err := utils2.ExtractStringAttr(d, "name")
		if err != nil {
			diags.Append(err...)
		}
		device.Name = name

		//  "lat": 37.4180951010362,
		if device.Lat.IsNull() || device.Lat.IsUnknown() {
			lat, err := utils2.ExtractFloat64Attr(d, "lat")
			if err != nil {
				diags.Append(err...)
			}
			device.Lat = lat
		}

		//  "lng": -122.098531723022,
		if device.Lng.IsNull() || device.Lng.IsUnknown() {
			lng, err := utils2.ExtractFloat64Attr(d, "lng")
			if err != nil {
				diags.Append(err...)
			}
			device.Lng = lng
		}

		//  "address": "1600 Pennsylvania Ave",
		if device.Address.IsNull() || device.Address.IsUnknown() {
			address, err := utils2.ExtractStringAttr(d, "address")
			if err != nil {
				diags.Append(err...)
			}
			device.Address = address
		}

		//  "notes": "My AP's note",
		if device.Notes.IsNull() || device.Notes.IsUnknown() {
			notes, err := utils2.ExtractStringAttr(d, "notes")
			if err != nil {
				diags.Append(err...)
			}
			device.Notes = notes
		}

		//  "tags": [
		//    " recently-added "
		//  ],
		if device.Tags.IsNull() || device.Tags.IsUnknown() {
			tags, err := utils2.ExtractListStringAttr(d, "tags")
			if err != nil {
				diags.Append(err...)
			}

			device.Tags = tags
		}

		//  "networkId": "N_24329156",
		if device.NetworkId.IsNull() || device.NetworkId.IsUnknown() {
			networkId, err := utils2.ExtractStringAttr(d, "networkId")
			if err != nil {
				diags.Append(err...)
			}
			device.NetworkId = networkId
		}

		//  "serial": "Q234-ABCD-5678",
		if device.Serial.IsNull() || device.Serial.IsUnknown() {
			serial, err := utils2.ExtractStringAttr(d, "serial")
			if err != nil {
				diags.Append(err...)
			}
			device.Serial = serial
		}

		//  "model": "MR34",
		if device.Model.IsNull() || device.Model.IsUnknown() {
			model, err := utils2.ExtractStringAttr(d, "model")
			if err != nil {
				diags.Append(err...)
			}
			device.Model = model
		}

		//  "mac": "00:11:22:33:44:55",
		if device.Mac.IsNull() || device.Mac.IsUnknown() {
			mac, err := utils2.ExtractStringAttr(d, "mac")
			if err != nil {
				diags.Append(err...)
			}
			device.Mac = mac
		}

		//  "lanIp": "1.2.3.4",
		if device.LanIp.IsNull() || device.LanIp.IsUnknown() {
			lanIp, err := utils2.ExtractStringAttr(d, "lanIp")
			if err != nil {
				diags.Append(err...)
			}
			device.LanIp = lanIp
		}

		//  "firmware": "wireless-25-14",
		if device.Firmware.IsNull() || device.Firmware.IsUnknown() {
			firmware, err := utils2.ExtractStringAttr(d, "firmware")
			if err != nil {
				diags.Append(err...)
			}
			device.Firmware = firmware
		}

		//  "floorPlanId": "g_2176982374",
		if device.FloorPlanId.IsNull() || device.FloorPlanId.IsUnknown() {
			floorPlanId, err := utils2.ExtractStringAttr(d, "floorPlanId")
			if err != nil {
				diags.Append(err...)
			}
			device.FloorPlanId = floorPlanId
		}

		//  "details": [
		//    {
		//      "name": "Catalyst serial",
		//      "value": "123ABC"
		//    }
		//  ],
		if device.Details.IsNull() || device.Details.IsUnknown() {

			detailAttr := map[string]attr.Type{
				"name":  types.StringType,
				"value": types.StringType,
			}

			detailsAttrs := types.ObjectType{AttrTypes: detailAttr}

			_, ok := d["details"].([]map[string]interface{})
			if ok {

				detailsList, err := utils2.ExtractListAttr(d, "details", detailsAttrs)
				if err.HasError() {
					tflog.Error(context.Background(), fmt.Sprintf("%s", err.Errors()))
				}

				device.Details = detailsList

			} else {
				detailsArrayObjNull := types.ListNull(detailsAttrs)
				device.Details = detailsArrayObjNull
			}

		}

		//  "beaconIdParams": {
		//    "uuid": "00000000-0000-0000-0000-000000000000",
		//    "major": 5,
		//    "minor": 3
		if device.BeaconIdParams.IsNull() || device.BeaconIdParams.IsUnknown() {
			beaconIdParamsAttrs := map[string]attr.Type{
				"uuid":  types.StringType,
				"major": types.Int64Type,
				"minor": types.Int64Type,
			}

			beaconIdParamsResp, ok := d["beaconIdParams"].(map[string]interface{})
			if ok {
				var beaconIdParams DevicesResourceModelBeaconIdParams

				// uuid
				uuid, err := utils2.ExtractStringAttr(beaconIdParamsResp, "uuid")
				if err.HasError() {
					diags.AddError("uuid Attr", fmt.Sprintf("%s", err.Errors()))
				}

				beaconIdParams.Uuid = uuid

				// major
				major, err := utils2.ExtractInt32Attr(beaconIdParamsResp, "major")
				if err.HasError() {
					diags.AddError("major Attr", fmt.Sprintf("%s", err.Errors()))
				}

				beaconIdParams.Major = major

				// minor
				minor, err := utils2.ExtractInt32Attr(beaconIdParamsResp, "minor")
				if err.HasError() {
					diags.AddError("minor Attr", fmt.Sprintf("%s", err.Errors()))
				}

				beaconIdParams.Minor = minor

				beaconIdParamsObj, err := types.ObjectValueFrom(ctx, beaconIdParamsAttrs, beaconIdParams)
				if err.HasError() {
					diags.AddError("beaconIdParamsObj Attr", fmt.Sprintf("%s", err.Errors()))
				}

				device.BeaconIdParams = beaconIdParamsObj
			} else {
				beaconIdParamsObjNull := types.ObjectNull(beaconIdParamsAttrs)
				device.BeaconIdParams = beaconIdParamsObjNull
			}

		}

		// url
		if device.Url.IsNull() || device.Url.IsUnknown() {
			url, err := utils2.ExtractStringAttr(d, "url")
			if err.HasError() {
				diags.AddError("url Attr", fmt.Sprintf("%s", err.Errors()))
			}
			device.Url = url
		}

		// SwitchProfileId
		if device.SwitchProfileId.IsNull() || device.SwitchProfileId.IsUnknown() {
			switchProfileId, err := utils2.ExtractStringAttr(d, "switchProfileId")
			if err.HasError() {
				diags.AddError("switchProfileId Attr", fmt.Sprintf("%s", err.Errors()))
			}
			device.SwitchProfileId = switchProfileId
		}

		// MoveMapMarker
		if device.MoveMapMarker.IsNull() || device.MoveMapMarker.IsUnknown() {
			moveMapMarker, err := utils2.ExtractBoolAttr(d, "moveMapMarker")
			if err.HasError() {
				diags.AddError("moveMapMarker Attr", fmt.Sprintf("%s", err.Errors()))
			}
			device.MoveMapMarker = moveMapMarker
		}

		deviceObj, err := types.ObjectValueFrom(ctx, devicesDatasourceModelDeviceAttrs, device)
		if err.HasError() {
			diags.Append(err...)
		}

		deviceList = append(deviceList, deviceObj)
	}

	deviceListArray, err := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: devicesDatasourceModelDeviceAttrs}, deviceList)
	if err.HasError() {
		diags.Append(err...)
	}

	state.List = deviceListArray

	return diags
}

func (d *NetworkDevicesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DevicesDatasourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	inlineResp, httpResp, err := d.client.DevicesApi.GetNetworkDevices(ctx, data.NetworkId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			utils2.HttpDiagnostics(httpResp),
		)
	}

	// Check for API success inlineResp code
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	}

	diags := updateDevicesDatasourceState(ctx, &data, inlineResp)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read a data source")

}
