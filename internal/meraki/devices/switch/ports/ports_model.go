package ports

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PortResourceModel struct {
	PortId                  types.String                      `tfsdk:"port_id"`
	Name                    types.String                      `tfsdk:"name"`
	Tags                    []types.String                    `tfsdk:"tags"`
	Enabled                 types.Bool                        `tfsdk:"enabled"`
	PoeEnabled              types.Bool                        `tfsdk:"poe_enabled"`
	Type                    types.String                      `tfsdk:"type"`
	Vlan                    types.Int64                       `tfsdk:"vlan"`
	VoiceVlan               types.Int64                       `tfsdk:"voice_vlan"`
	AllowedVlans            types.String                      `tfsdk:"allowed_vlans"`
	IsolationEnabled        types.Bool                        `tfsdk:"isolation_enabled"`
	RstpEnabled             types.Bool                        `tfsdk:"rstp_enabled"`
	StpGuard                types.String                      `tfsdk:"stp_guard"`
	AccessPolicyNumber      types.Int64                       `tfsdk:"access_policy_number"`
	AccessPolicyType        types.String                      `tfsdk:"access_policy_type"`
	LinkNegotiation         types.String                      `tfsdk:"link_negotiation"`
	PortScheduleId          types.String                      `tfsdk:"port_schedule_id"`
	Udld                    types.String                      `tfsdk:"udld"`
	StickyMacWhitelistLimit types.Int64                       `tfsdk:"sticky_mac_white_list_limit"`
	StormControlEnabled     types.Bool                        `tfsdk:"storm_control_enabled"`
	MacWhitelist            []types.String                    `tfsdk:"mac_white_list"`
	StickyMacWhitelist      []types.String                    `tfsdk:"sticky_mac_white_list"`
	AdaptivePolicyGroupId   types.String                      `tfsdk:"adaptive_policy_group_id"`
	PeerSgtCapable          types.Bool                        `tfsdk:"peer_sgt_capable"`
	FlexibleStackingEnabled types.Bool                        `tfsdk:"flexible_stacking_enabled"`
	DaiTrusted              types.Bool                        `tfsdk:"dai_trusted"`
	Profile                 SwitchPortsDataSourceModelProfile `tfsdk:"profile"`
}

type SwitchPortsDataSourceModelProfile struct {
	Enabled types.Bool   `tfsdk:"enabled"`
	Id      types.String `tfsdk:"id"`
	Iname   types.String `tfsdk:"iname"`
}

// The PortsDataSourceModel structure describes the data model.
// This struct is where you define all the attributes that are part of this data source's state.
type PortsDataSourceModel struct {
	Id     types.String        `tfsdk:"id"`
	Serial types.String        `tfsdk:"serial"`
	List   []PortResourceModel `tfsdk:"list"`
}
