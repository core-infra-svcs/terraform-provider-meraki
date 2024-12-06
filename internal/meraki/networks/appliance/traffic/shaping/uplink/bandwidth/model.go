package bandwidth

import "github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"

// resourceModel describes the resource data model.
type resourceModel struct {
	Id                              jsontypes.String `tfsdk:"id"`
	NetworkId                       jsontypes.String `tfsdk:"network_id" json:"network_id"`
	BandwidthLimitCellularLimitUp   jsontypes.Int64  `tfsdk:"bandwidth_limit_cellular_limit_up"`
	BandwidthLimitCellularLimitDown jsontypes.Int64  `tfsdk:"bandwidth_limit_cellular_limit_down"`
	BandwidthLimitWan2LimitUp       jsontypes.Int64  `tfsdk:"bandwidth_limit_wan2_limit_up"`
	BandwidthLimitWan2LimitDown     jsontypes.Int64  `tfsdk:"bandwidth_limit_wan2_limit_down"`
	BandwidthLimitWan1LimitUp       jsontypes.Int64  `tfsdk:"bandwidth_limit_wan1_limit_up"`
	BandwidthLimitWan1LimitDown     jsontypes.Int64  `tfsdk:"bandwidth_limit_wan1_limit_down"`
}

type resourceModelApiResponse struct {
	UplinkBandwidthLimits TrafficShapingUplinkBandWidthLimits `json:"bandwidthLimits"`
}

type TrafficShapingUplinkBandWidthLimits struct {
	Wan1     TrafficShapingUplinkBandWidthLimit `json:"wan1"`
	Wan2     TrafficShapingUplinkBandWidthLimit `json:"wan2"`
	Cellular TrafficShapingUplinkBandWidthLimit `json:"cellular"`
}

type TrafficShapingUplinkBandWidthLimit struct {
	LimitUp   jsontypes.Int64 `json:"limitUp"`
	LimitDown jsontypes.Int64 `json:"limitDown"`
}
