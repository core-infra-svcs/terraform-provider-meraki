package statuses

import "github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"

type dataSourceModel struct {
	Id             jsontypes.String      `tfsdk:"id"`
	OrganizationId jsontypes.String      `tfsdk:"organization_id"`
	PerPage        jsontypes.Int64       `tfsdk:"per_page"`
	StartingAfter  jsontypes.String      `tfsdk:"starting_after"`
	EndingBefore   jsontypes.String      `tfsdk:"ending_before"`
	NetworkIds     []jsontypes.String    `tfsdk:"network_ids"`
	Serials        []jsontypes.String    `tfsdk:"serials"`
	Iccids         []jsontypes.String    `tfsdk:"iccids"`
	List           []dataSourceModelList `tfsdk:"list"`
}

type dataSourceModelList struct {
	NetworkId      jsontypes.String        `tfsdk:"network_id" json:"networkId,omitempty"`
	Serial         jsontypes.String        `tfsdk:"serial"`
	Model          jsontypes.String        `tfsdk:"model"`
	LastReportedAt jsontypes.String        `tfsdk:"last_reported_at"`
	Uplinks        []dataSourceModelUplink `tfsdk:"uplinks"`
}

type dataSourceModelUplink struct {
	Interface      jsontypes.String          `tfsdk:"interface"`
	Status         jsontypes.String          `tfsdk:"status"`
	Ip             jsontypes.String          `tfsdk:"ip"`
	Provider       jsontypes.String          `tfsdk:"provider"`
	PublicIp       jsontypes.String          `tfsdk:"public_ip"`
	Model          jsontypes.String          `tfsdk:"model"`
	SignalStat     dataSourceModelSignalStat `tfsdk:"signal_stat"`
	ConnectionType jsontypes.String          `tfsdk:"connection_type"`
	Apn            jsontypes.String          `tfsdk:"apn"`
	Gateway        jsontypes.String          `tfsdk:"gateway"`
	Dns1           jsontypes.String          `tfsdk:"dns1"`
	Dns2           jsontypes.String          `tfsdk:"dns2"`
	SignalType     jsontypes.String          `tfsdk:"signal_type"`
	Iccid          jsontypes.String          `tfsdk:"iccid"`
}

type dataSourceModelSignalStat struct {
	Rsrp jsontypes.String `tfsdk:"rsrp"`
	Rsrq jsontypes.String `tfsdk:"rsrq"`
}
