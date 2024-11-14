package _switch

import (
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PortProfileModel struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled"`
	Id      types.String `tfsdk:"id" json:"id"`
	Iname   types.String `tfsdk:"iname" json:"iname"`
}

type PortResourceModel struct {
	Id                          types.String `tfsdk:"id"`
	Serial                      types.String `tfsdk:"serial" json:"serial"`
	PortId                      types.String `tfsdk:"port_id" json:"portId"`
	Name                        types.String `tfsdk:"name" json:"name"`
	Tags                        types.Set    `tfsdk:"tags" json:"tags"`
	Enabled                     types.Bool   `tfsdk:"enabled" json:"enabled"`
	PoeEnabled                  types.Bool   `tfsdk:"poe_enabled" json:"poeEnabled"`
	Type                        types.String `tfsdk:"type" json:"type"`
	Vlan                        types.Int64  `tfsdk:"vlan" json:"vlan"`
	VoiceVlan                   types.Int64  `tfsdk:"voice_vlan" json:"voiceVlan"`
	AllowedVlans                types.String `tfsdk:"allowed_vlans" json:"allowedVlans"`
	AccessPolicyNumber          types.Int64  `tfsdk:"access_policy_number" json:"accessPolicyNumber"`
	AccessPolicyType            types.String `tfsdk:"access_policy_type" json:"accessPolicyType"`
	PortScheduleId              types.String `tfsdk:"port_schedule_id" json:"portScheduleId"`
	StickyMacAllowListLimit     types.Int64  `tfsdk:"sticky_mac_allow_list_limit" json:"stickyMacWhitelistLimit"`
	MacAllowList                types.Set    `tfsdk:"mac_allow_list" json:"macWhitelist"`
	StickyMacAllowList          types.Set    `tfsdk:"sticky_mac_allow_list" json:"stickyMacWhitelist"`
	StormControlEnabled         types.Bool   `tfsdk:"storm_control_enabled" json:"stormControlEnabled"`
	AdaptivePolicyGroupId       types.String `tfsdk:"adaptive_policy_group_id" json:"adaptivePolicyGroupId"`
	PeerSgtCapable              types.Bool   `tfsdk:"peer_sgt_capable" json:"peerSgtCapable"`
	FlexibleStackingEnabled     types.Bool   `tfsdk:"flexible_stacking_enabled" json:"flexibleStackingEnabled"`
	DaiTrusted                  types.Bool   `tfsdk:"dai_trusted" json:"daiTrusted"`
	IsolationEnabled            types.Bool   `tfsdk:"isolation_enabled" json:"isolationEnabled"`
	RstpEnabled                 types.Bool   `tfsdk:"rstp_enabled" json:"rstpEnabled"`
	StpGuard                    types.String `tfsdk:"stp_guard" json:"stpGuard"`
	LinkNegotiation             types.String `tfsdk:"link_negotiation" json:"linkNegotiation"`
	LinkNegotiationCapabilities types.List   `tfsdk:"link_negotiation_capabilities" json:"linkNegotiationCapabilities"`
	Udld                        types.String `tfsdk:"udld" json:"udld"`
	Profile                     types.Object `tfsdk:"profile" json:"profile"`
}

type PortsDataSourceModel struct {
	Serial jsontypes.String    `tfsdk:"id" json:"serial"`
	Ports  []PortResourceModel `tfsdk:"list"`
}
