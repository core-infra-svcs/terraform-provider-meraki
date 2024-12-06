package groupPolicy

import (
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type datasourceModel struct {
	Id        types.String `tfsdk:"id"`
	NetworkId types.String `tfsdk:"network_id"`
	List      types.List   `tfsdk:"list"`
}

type resourceModel struct {
	ID                        types.String `tfsdk:"id"`
	Name                      types.String `tfsdk:"name" json:"name"`
	GroupPolicyId             types.String `tfsdk:"group_policy_id" json:"groupPolicyId"`
	NetworkId                 types.String `tfsdk:"network_id" json:"networkId"`
	Scheduling                types.Object `tfsdk:"scheduling" json:"scheduling"`
	Bandwidth                 types.Object `tfsdk:"bandwidth" json:"bandwidth"`
	FirewallAndTrafficShaping types.Object `tfsdk:"firewall_and_traffic_shaping" json:"firewallAndTrafficShaping"`
	ContentFiltering          types.Object `tfsdk:"content_filtering" json:"contentFiltering"`
	SplashAuthSettings        types.String `tfsdk:"splash_auth_settings" json:"splashAuthSettings"`
	VlanTagging               types.Object `tfsdk:"vlan_tagging" json:"vlanTagging"`
	BonjourForwarding         types.Object `tfsdk:"bonjour_forwarding" json:"bonjourForwarding"`
}

func resourceModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"id":                           types.StringType,
		"name":                         types.StringType,
		"group_policy_id":              types.StringType,
		"network_id":                   types.StringType,
		"scheduling":                   types.ObjectType{AttrTypes: schedulingModelAttrs()},
		"bandwidth":                    types.ObjectType{AttrTypes: bandwidthModelAttrs()},
		"firewall_and_traffic_shaping": types.ObjectType{AttrTypes: firewallAndTrafficShapingModelAttrs()},
		"content_filtering":            types.ObjectType{AttrTypes: contentFilteringModelAttrs()},
		"splash_auth_settings":         types.StringType,
		"vlan_tagging":                 types.ObjectType{AttrTypes: vlanTaggingModelAttrs()},
		"bonjour_forwarding":           types.ObjectType{AttrTypes: bonjourForwardingModelAttrs()},
	}
}

// SchedulingModel represents the scheduling settings.
type SchedulingModel struct {
	Enabled   types.Bool   `tfsdk:"enabled" json:"enabled"`
	Monday    types.Object `tfsdk:"monday" json:"monday"`
	Tuesday   types.Object `tfsdk:"tuesday" json:"tuesday"`
	Wednesday types.Object `tfsdk:"wednesday" json:"wednesday"`
	Thursday  types.Object `tfsdk:"thursday" json:"thursday"`
	Friday    types.Object `tfsdk:"friday" json:"friday"`
	Saturday  types.Object `tfsdk:"saturday" json:"saturday"`
	Sunday    types.Object `tfsdk:"sunday" json:"sunday"`
}

func schedulingModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"enabled":   types.BoolType,
		"monday":    types.ObjectType{AttrTypes: scheduleDayModelAttrs()},
		"tuesday":   types.ObjectType{AttrTypes: scheduleDayModelAttrs()},
		"wednesday": types.ObjectType{AttrTypes: scheduleDayModelAttrs()},
		"thursday":  types.ObjectType{AttrTypes: scheduleDayModelAttrs()},
		"friday":    types.ObjectType{AttrTypes: scheduleDayModelAttrs()},
		"saturday":  types.ObjectType{AttrTypes: scheduleDayModelAttrs()},
		"sunday":    types.ObjectType{AttrTypes: scheduleDayModelAttrs()},
	}
}

// ScheduleDayModel represents a single day's schedule.
type ScheduleDayModel struct {
	Active types.Bool   `tfsdk:"active" json:"active"`
	From   types.String `tfsdk:"from" json:"from"`
	To     types.String `tfsdk:"to" json:"to"`
}

func scheduleDayModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"active": types.BoolType,
		"from":   types.StringType,
		"to":     types.StringType,
	}
}

// BandwidthModel represents the bandwidth settings.
type BandwidthModel struct {
	Settings        types.String `tfsdk:"settings" json:"settings"`
	BandwidthLimits types.Object `tfsdk:"bandwidth_limits" json:"bandwidthLimits"`
}

func bandwidthModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"settings":         types.StringType,
		"bandwidth_limits": types.ObjectType{AttrTypes: bandwidthLimitsModelAttrs()},
	}
}

// BandwidthLimitsModel represents the bandwidth limits.
type BandwidthLimitsModel struct {
	LimitUp   types.Int64 `tfsdk:"limit_up" json:"limitUp"`
	LimitDown types.Int64 `tfsdk:"limit_down" json:"limitDown"`
}

func bandwidthLimitsModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"limit_up":   types.Int64Type,
		"limit_down": types.Int64Type,
	}
}

// FirewallAndTrafficShapingModel represents the firewall and traffic shaping settings.
type FirewallAndTrafficShapingModel struct {
	Settings            types.String `tfsdk:"settings" json:"settings"`
	L3FirewallRules     types.List   `tfsdk:"l3_firewall_rules" json:"l3FirewallRules"`
	L7FirewallRules     types.List   `tfsdk:"l7_firewall_rules" json:"l7FirewallRules"`
	TrafficShapingRules types.List   `tfsdk:"traffic_shaping_rules" json:"trafficShapingRules"`
}

// TrafficShapingRuleModel represents a traffic shaping rule.2
type TrafficShapingRuleModel struct {
	DscpTagValue             types.Int64  `tfsdk:"dscp_tag_value" json:"dscpTagValue"`
	PcpTagValue              types.Int64  `tfsdk:"pcp_tag_value" json:"pcpTagValue"`
	PerClientBandwidthLimits types.Object `tfsdk:"per_client_bandwidth_limits" json:"perClientBandwidthLimits"`
	Definitions              types.List   `tfsdk:"definitions" json:"definitions"`
}

func firewallAndTrafficShapingModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"settings":              types.StringType,
		"l3_firewall_rules":     types.ListType{ElemType: types.ObjectType{AttrTypes: l3FirewallRuleModelAttrs()}},
		"l7_firewall_rules":     types.ListType{ElemType: types.ObjectType{AttrTypes: l7FirewallRuleModelAttrs()}},
		"traffic_shaping_rules": types.ListType{ElemType: types.ObjectType{AttrTypes: trafficShapingRuleModelAttrs()}},
	}
}

// L3FirewallRuleModel represents a layer 3 firewall rule.
type L3FirewallRuleModel struct {
	Comment  types.String `tfsdk:"comment" json:"comment"`
	Policy   types.String `tfsdk:"policy" json:"policy"`
	Protocol types.String `tfsdk:"protocol" json:"protocol"`
	DestPort types.String `tfsdk:"dest_port" json:"destPort"`
	DestCidr types.String `tfsdk:"dest_cidr" json:"destCidr"`
}

func l3FirewallRuleModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"comment":   types.StringType,
		"policy":    types.StringType,
		"protocol":  types.StringType,
		"dest_port": types.StringType,
		"dest_cidr": types.StringType,
	}
}

// L7FirewallRuleModel represents a layer 7 firewall rule.
type L7FirewallRuleModel struct {
	Policy types.String `tfsdk:"policy" json:"policy"`
	Type   types.String `tfsdk:"type" json:"type"`
	Value  types.String `tfsdk:"value" json:"value"`
}

func l7FirewallRuleModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"policy": types.StringType,
		"type":   types.StringType,
		"value":  types.StringType,
	}
}

// TrafficShapingDefinitionModel represents a traffic shaping definition.
type TrafficShapingDefinitionModel struct {
	Type  types.String `tfsdk:"type" json:"type"`
	Value types.String `tfsdk:"value" json:"value"`
}

func trafficShapingRuleModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"dscp_tag_value":              types.Int64Type,
		"pcp_tag_value":               types.Int64Type,
		"per_client_bandwidth_limits": types.ObjectType{AttrTypes: perClientBandwidthLimitsModelAttrs()},
		"definitions":                 types.ListType{ElemType: types.ObjectType{AttrTypes: trafficShapingDefinitionModelAttrs()}},
	}
}

func trafficShapingDefinitionModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"type":  types.StringType,
		"value": types.StringType,
	}
}

// PerClientBandwidthLimitsModel represents the per-client bandwidth limits.
type PerClientBandwidthLimitsModel struct {
	Settings        types.String `tfsdk:"settings" json:"settings"`
	BandwidthLimits types.Object `tfsdk:"bandwidth_limits" json:"bandwidthLimits"`
}

func perClientBandwidthLimitsModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"settings":         types.StringType,
		"bandwidth_limits": types.ObjectType{AttrTypes: bandwidthLimitsModelAttrs()},
	}
}

// ContentFilteringModel represents the content filtering settings.
type ContentFilteringModel struct {
	AllowedUrlPatterns   types.Object `tfsdk:"allowed_url_patterns" json:"allowedUrlPatterns"`
	BlockedUrlPatterns   types.Object `tfsdk:"blocked_url_patterns" json:"blockedUrlPatterns"`
	BlockedUrlCategories types.Object `tfsdk:"blocked_url_categories" json:"blockedUrlCategories"`
}

func contentFilteringModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"allowed_url_patterns":   types.ObjectType{AttrTypes: urlPatternsModelAttrs()},
		"blocked_url_patterns":   types.ObjectType{AttrTypes: urlPatternsModelAttrs()},
		"blocked_url_categories": types.ObjectType{AttrTypes: urlCategoriesModelAttrs()},
	}
}

type UrlPatternsModel struct {
	Patterns types.List   `tfsdk:"patterns" json:"patterns"`
	Settings types.String `tfsdk:"settings" json:"settings"`
}

func urlPatternsModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"patterns": types.ListType{ElemType: types.StringType},
		"settings": types.StringType,
	}
}

type UrlCategoriesModel struct {
	Categories types.List   `tfsdk:"categories" json:"categories"`
	Settings   types.String `tfsdk:"settings" json:"settings"`
}

func urlCategoriesModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"categories": types.ListType{ElemType: types.StringType},
		"settings":   types.StringType,
	}
}

// VlanTaggingModel represents the VLAN tagging settings.
type VlanTaggingModel struct {
	Settings types.String `tfsdk:"settings" json:"settings"`
	VlanID   types.String `tfsdk:"vlan_id" json:"vlanId"`
}

func vlanTaggingModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"settings": types.StringType,
		"vlan_id":  types.StringType,
	}
}

// BonjourForwardingModel represents the Bonjour forwarding settings.
type BonjourForwardingModel struct {
	Settings types.String `tfsdk:"settings" json:"settings"`
	Rules    types.List   `tfsdk:"rules" json:"rules"`
}

func bonjourForwardingModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"settings": types.StringType,
		"rules":    types.ListType{ElemType: types.ObjectType{AttrTypes: bonjourForwardingRuleModelAttrs()}},
	}
}

// BonjourForwardingRuleModel represents a Bonjour forwarding rule.
type BonjourForwardingRuleModel struct {
	Description types.String `tfsdk:"description" json:"description"`
	VlanID      types.String `tfsdk:"vlan_id" json:"vlanId"`
	Services    types.List   `tfsdk:"services" json:"services"`
}

func bonjourForwardingRuleModelAttrs() map[string]attr.Type {
	return map[string]attr.Type{
		"description": types.StringType,
		"vlan_id":     types.StringType,
		"services":    types.ListType{ElemType: types.StringType},
	}
}
