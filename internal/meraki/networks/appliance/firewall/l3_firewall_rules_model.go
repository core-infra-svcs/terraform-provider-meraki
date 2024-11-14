package firewall

import "github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"

// L3FirewallRulesModel describes the resource data model.
type L3FirewallRulesModel struct {
	Id                jsontypes.String           `tfsdk:"id" json:"network_id"`
	SyslogDefaultRule jsontypes.Bool             `tfsdk:"syslog_default_rule"`
	Rules             []L3FirewallRulesRuleModel `tfsdk:"rules" json:"rules"`
}

type L3FirewallRulesRuleModel struct {
	Comment       jsontypes.String `tfsdk:"comment"`
	DestCidr      jsontypes.String `tfsdk:"dest_cidr"`
	DestPort      jsontypes.String `tfsdk:"dest_port"`
	Policy        jsontypes.String `tfsdk:"policy"`
	Protocol      jsontypes.String `tfsdk:"protocol"`
	SrcPort       jsontypes.String `tfsdk:"src_port"`
	SrcCidr       jsontypes.String `tfsdk:"src_cidr"`
	SysLogEnabled jsontypes.Bool   `tfsdk:"syslog_enabled"`
}
