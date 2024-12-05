package device

import "github.com/hashicorp/terraform-plugin-framework/types"

type resourceModel struct {
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

type beaconIdParamsModel struct {
	Uuid  types.String `tfsdk:"uuid"`
	Major types.Int64  `tfsdk:"major"`
	Minor types.Int64  `tfsdk:"minor"`
}

type detailsModel struct {
	Name  types.String `tfsdk:"name"`
	Value types.String `tfsdk:"value"`
}
