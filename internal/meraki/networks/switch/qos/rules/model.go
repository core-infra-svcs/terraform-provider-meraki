package rules

import "github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"

// dataSourceModel describes the resource data model.
type dataSourceModel struct {
	Id        jsontypes.String       `tfsdk:"id" json:"-"`
	NetworkId jsontypes.String       `tfsdk:"network_id" json:"network_id"`
	List      []dataSourceModelRules `tfsdk:"list"`
}

// dataSourceModelRules describes the resource data model.
type dataSourceModelRules struct {
	QosRulesId   jsontypes.String  `tfsdk:"qos_rule_id" json:"id"`
	Vlan         jsontypes.Int64   `tfsdk:"vlan" json:"vlan"`
	Dscp         jsontypes.Int64   `tfsdk:"dscp" json:"dscp"`
	DstPort      jsontypes.Float64 `tfsdk:"dst_port" json:"dstPort"`
	SrcPort      jsontypes.Float64 `tfsdk:"src_port" json:"srcPort"`
	DstPortRange jsontypes.String  `tfsdk:"dst_port_range" json:"dstPortRange"`
	Protocol     jsontypes.String  `tfsdk:"protocol" json:"protocol"`
	SrcPortRange jsontypes.String  `tfsdk:"src_port_range" json:"srcPortRange"`
}

// resourceModel describes the resource data model.
type resourceModel struct {
	Id           jsontypes.String  `tfsdk:"id" json:"-"`
	NetworkId    jsontypes.String  `tfsdk:"network_id" json:"network_id"`
	QosRulesId   jsontypes.String  `tfsdk:"qos_rule_id" json:"id"`
	Vlan         jsontypes.Int64   `tfsdk:"vlan" json:"vlan"`
	Dscp         jsontypes.Int64   `tfsdk:"dscp" json:"dscp"`
	DstPort      jsontypes.Float64 `tfsdk:"dst_port" json:"dstPort"`
	SrcPort      jsontypes.Float64 `tfsdk:"src_port" json:"srcPort"`
	DstPortRange jsontypes.String  `tfsdk:"dst_port_range" json:"dstPortRange"`
	Protocol     jsontypes.String  `tfsdk:"protocol" json:"protocol"`
	SrcPortRange jsontypes.String  `tfsdk:"src_port_range" json:"srcPortRange"`
}
