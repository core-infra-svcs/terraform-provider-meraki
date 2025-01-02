package rules

import "github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"

// resourceModel describes the resource data model.
type resourceModel struct {
	Id        jsontypes.String      `tfsdk:"id"`
	NetworkId jsontypes.String      `tfsdk:"network_id" json:"network_id"`
	Rules     []l7FirewallRuleModel `tfsdk:"rules" json:"rules"`
}

type l7FirewallRuleModel struct {
	Policy jsontypes.String `tfsdk:"policy"`
	Type   jsontypes.String `tfsdk:"type"`
	Value  jsontypes.String `tfsdk:"value"`
}
