package groupPolicy

import "github.com/hashicorp/terraform-plugin-framework/types"

// GroupPolicyModel represents a group policy.
type GroupPolicyModel struct {
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

// ScheduleDayModel represents a single day's schedule.
type ScheduleDayModel struct {
	Active types.Bool   `tfsdk:"active" json:"active"`
	From   types.String `tfsdk:"from" json:"from"`
	To     types.String `tfsdk:"to" json:"to"`
}

// BandwidthModel represents the bandwidth settings.
type BandwidthModel struct {
	Settings        types.String `tfsdk:"settings" json:"settings"`
	BandwidthLimits types.Object `tfsdk:"bandwidth_limits" json:"bandwidthLimits"`
}

// BandwidthLimitsModel represents the bandwidth limits.
type BandwidthLimitsModel struct {
	LimitUp   types.Int64 `tfsdk:"limit_up" json:"limitUp"`
	LimitDown types.Int64 `tfsdk:"limit_down" json:"limitDown"`
}

// FirewallAndTrafficShapingModel represents the firewall and traffic shaping settings.
type FirewallAndTrafficShapingModel struct {
	Settings            types.String `tfsdk:"settings" json:"settings"`
	L3FirewallRules     types.List   `tfsdk:"l3_firewall_rules" json:"l3FirewallRules"`
	L7FirewallRules     types.List   `tfsdk:"l7_firewall_rules" json:"l7FirewallRules"`
	TrafficShapingRules types.List   `tfsdk:"traffic_shaping_rules" json:"trafficShapingRules"`
}

// L3FirewallRuleModel represents a layer 3 firewall rule.
type L3FirewallRuleModel struct {
	Comment  types.String `tfsdk:"comment" json:"comment"`
	Policy   types.String `tfsdk:"policy" json:"policy"`
	Protocol types.String `tfsdk:"protocol" json:"protocol"`
	DestPort types.String `tfsdk:"dest_port" json:"destPort"`
	DestCidr types.String `tfsdk:"dest_cidr" json:"destCidr"`
}

// L7FirewallRuleModel represents a layer 7 firewall rule.
type L7FirewallRuleModel struct {
	Policy types.String `tfsdk:"policy" json:"policy"`
	Type   types.String `tfsdk:"type" json:"type"`
	Value  types.String `tfsdk:"value" json:"value"`
}

// TrafficShapingRuleModel represents a traffic shaping rule.2
type TrafficShapingRuleModel struct {
	DscpTagValue             types.Int64  `tfsdk:"dscp_tag_value" json:"dscpTagValue"`
	PcpTagValue              types.Int64  `tfsdk:"pcp_tag_value" json:"pcpTagValue"`
	PerClientBandwidthLimits types.Object `tfsdk:"per_client_bandwidth_limits" json:"perClientBandwidthLimits"`
	Definitions              types.List   `tfsdk:"definitions" json:"definitions"`
}

// PerClientBandwidthLimitsModel represents the per-client bandwidth limits.
type PerClientBandwidthLimitsModel struct {
	Settings        types.String `tfsdk:"settings" json:"settings"`
	BandwidthLimits types.Object `tfsdk:"bandwidth_limits" json:"bandwidthLimits"`
}

// TrafficShapingDefinitionModel represents a traffic shaping definition.
type TrafficShapingDefinitionModel struct {
	Type  types.String `tfsdk:"type" json:"type"`
	Value types.String `tfsdk:"value" json:"value"`
}

// ContentFilteringModel represents the content filtering settings.
type ContentFilteringModel struct {
	AllowedUrlPatterns   types.Object `tfsdk:"allowed_url_patterns" json:"allowedUrlPatterns"`
	BlockedUrlPatterns   types.Object `tfsdk:"blocked_url_patterns" json:"blockedUrlPatterns"`
	BlockedUrlCategories types.Object `tfsdk:"blocked_url_categories" json:"blockedUrlCategories"`
}

type UrlPatternsModel struct {
	Patterns types.List   `tfsdk:"patterns" json:"patterns"`
	Settings types.String `tfsdk:"settings" json:"settings"`
}

type UrlCategoriesModel struct {
	Categories types.List   `tfsdk:"categories" json:"categories"`
	Settings   types.String `tfsdk:"settings" json:"settings"`
}

// VlanTaggingModel represents the VLAN tagging settings.
type VlanTaggingModel struct {
	Settings types.String `tfsdk:"settings" json:"settings"`
	VlanID   types.String `tfsdk:"vlan_id" json:"vlanId"`
}

// BonjourForwardingModel represents the Bonjour forwarding settings.
type BonjourForwardingModel struct {
	Settings types.String `tfsdk:"settings" json:"settings"`
	Rules    types.List   `tfsdk:"rules" json:"rules"`
}

// BonjourForwardingRuleModel represents a Bonjour forwarding rule.
type BonjourForwardingRuleModel struct {
	Description types.String `tfsdk:"description" json:"description"`
	VlanID      types.String `tfsdk:"vlan_id" json:"vlanId"`
	Services    types.List   `tfsdk:"services" json:"services"`
}
