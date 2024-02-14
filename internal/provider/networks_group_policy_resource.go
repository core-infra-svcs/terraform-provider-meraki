package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"io"
	"strings"

	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &NetworksGroupPolicyResource{}
var _ resource.ResourceWithImportState = &NetworksGroupPolicyResource{}

func NewNetworksGroupPolicyResource() resource.Resource {
	return &NetworksGroupPolicyResource{}
}

// NetworksGroupPolicyResource defines the resource implementation.
type NetworksGroupPolicyResource struct {
	client *openApiClient.APIClient
}

type NetworkGroupPolicyResponse struct {
	Name                      string                    `json:"name"`
	GroupPolicyId             string                    `json:"groupPolicyId"`
	Scheduling                Scheduling                `json:"scheduling"`
	Bandwidth                 Bandwidth                 `json:"bandwidth"`
	FirewallAndTrafficShaping FirewallAndTrafficShaping `json:"firewallAndTrafficShaping"`
	ContentFiltering          ContentFiltering          `json:"contentFiltering"`
	SplashAuthSettings        string                    `json:"splashAuthSettings"`
	VlanTagging               VlanTagging               `json:"vlanTagging"`
	BonjourForwarding         BonjourForwarding         `json:"bonjourForwarding"`
}

type Scheduling struct {
	Enabled   bool `json:"enabled"`
	Monday    Day  `json:"monday"`
	Tuesday   Day  `json:"tuesday"`
	Wednesday Day  `json:"wednesday"`
	Thursday  Day  `json:"thursday"`
	Friday    Day  `json:"friday"`
	Saturday  Day  `json:"saturday"`
	Sunday    Day  `json:"sunday"`
}

type Day struct {
	Active bool   `json:"active"`
	From   string `json:"from"`
	To     string `json:"to"`
}

type Bandwidth struct {
	Settings        string          `json:"settings"`
	BandwidthLimits BandwidthLimits `json:"bandwidthLimits"`
}

type FirewallAndTrafficShaping struct {
	Settings            string               `json:"settings"`
	TrafficShapingRules []TrafficShapingRule `json:"trafficShapingRules"`
	L3FirewallRules     []L3FirewallRule     `json:"l3FirewallRules"`
	L7FirewallRules     []L7FirewallRule     `json:"l7FirewallRules"`
}

type TrafficShapingRule struct {
	Definitions              []Definition             `json:"definitions"`
	PerClientBandwidthLimits PerClientBandwidthLimits `json:"perClientBandwidthLimits"`
	DscpTagValue             int                      `json:"dscpTagValue"`
	PcpTagValue              int                      `json:"pcpTagValue"`
}

type Definition struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type PerClientBandwidthLimits struct {
	Settings        string          `json:"settings"`
	BandwidthLimits BandwidthLimits `json:"bandwidthLimits"`
}

type BandwidthLimits struct {
	LimitUp   int `json:"limitUp"`
	LimitDown int `json:"limitDown"`
}

type L3FirewallRule struct {
	Comment  string `json:"comment"`
	Policy   string `json:"policy"`
	Protocol string `json:"protocol"`
	DestPort string `json:"destPort"`
	DestCidr string `json:"destCidr"`
}
type L7FirewallRule struct {
	Policy string `json:"policy"`
	Type   string `json:"type"`
	Value  string `json:"value"`
}

type ContentFiltering struct {
	AllowedUrlPatterns   AllowedUrlPatterns   `json:"allowedUrlPatterns"`
	BlockedUrlPatterns   BlockedUrlPatterns   `json:"blockedUrlPatterns"`
	BlockedUrlCategories BlockedUrlCategories `json:"blockedUrlCategories"`
}

type AllowedUrlPatterns struct {
	Settings string   `json:"settings"`
	Patterns []string `json:"patterns"`
}

type BlockedUrlPatterns struct {
	Settings string   `json:"settings"`
	Patterns []string `json:"patterns"`
}

type BlockedUrlCategories struct {
	Settings   string   `json:"settings"`
	Categories []string `json:"categories"`
}

type VlanTagging struct {
	Settings string `json:"settings"`
	VlanId   string `json:"vlanId"`
}

type BonjourForwarding struct {
	Settings string `json:"settings"`
	Rules    []Rule `json:"rules"`
}

type Rule struct {
	Description string   `json:"description"`
	VlanId      string   `json:"vlanId"`
	Services    []string `json:"services"`
}

// NetworksGroupPolicyResourceModel describes the resource data model.
type NetworksGroupPolicyResourceModel struct {
	Id                        types.String `tfsdk:"id"`
	NetworkId                 types.String `tfsdk:"network_id"`
	GroupPolicyId             types.String `tfsdk:"group_policy_id" json:"groupPolicyId"`
	Name                      types.String `tfsdk:"name" json:"name"`
	SplashAuthSettings        types.String `tfsdk:"splash_auth_settings" json:"splashAuthSettings"`
	Bandwidth                 types.Object `tfsdk:"bandwidth" json:"bandwidth"`
	BonjourForwarding         types.Object `tfsdk:"bonjour_forwarding" json:"bonjourForwarding"`
	FirewallAndTrafficShaping types.Object `tfsdk:"firewall_and_traffic_shaping" json:"firewallAndTrafficShaping"`
	Scheduling                types.Object `tfsdk:"scheduling" json:"scheduling"`
	VlanTagging               types.Object `tfsdk:"vlan_tagging" json:"vlanTagging"`
	ContentFiltering          types.Object `tfsdk:"content_filtering" json:"contentFiltering"`
}

func NetworksGroupPolicyResourceModelFirewallAndTrafficShapingAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"settings":              types.StringType,
		"l3_firewall_rules":     types.ListType{ElemType: types.ObjectType{AttrTypes: NetworksGroupPolicyResourceModelL3FirewallRuleAttrTypes()}},
		"l7_firewall_rules":     types.ListType{ElemType: types.ObjectType{AttrTypes: NetworksGroupPolicyResourceModelL7FirewallRuleAttrTypes()}},
		"traffic_shaping_rules": types.ListType{ElemType: types.ObjectType{AttrTypes: NetworksGroupPolicyResourceModelTrafficShapingRuleAttrTypes()}},
	}
}

type NetworksGroupPolicyResourceModelFirewallAndTrafficShaping struct {
	Settings            types.String `tfsdk:"settings" json:"settings"`
	L3FirewallRules     types.List   `tfsdk:"l3_firewall_rules" json:"l3FirewallRules"`
	L7FirewallRules     types.List   `tfsdk:"l7_firewall_rules" json:"l7FirewallRules"`
	TrafficShapingRules types.List   `tfsdk:"traffic_shaping_rules" json:"trafficShapingRules"`
}

func (s NetworksGroupPolicyResourceModelFirewallAndTrafficShaping) FromAPIResponse(ctx context.Context, f *FirewallAndTrafficShaping) diag.Diagnostics {
	s.Settings = types.StringValue(f.Settings)

	// l3
	l3rules := []NetworksGroupPolicyResourceModelL3FirewallRule{}
	for _, rule := range f.L3FirewallRules {
		l3rule := NetworksGroupPolicyResourceModelL3FirewallRule{}
		diags := l3rule.FromAPIResponse(ctx, &rule)
		if diags.HasError() {
			return diags
		}

		l3rules = append(l3rules, l3rule)
	}

	objectType := types.ObjectType{
		AttrTypes: NetworksGroupPolicyResourceModelL3FirewallRuleAttrTypes(),
	}

	l3RuleValue, diags := types.ListValueFrom(ctx, objectType, l3rules)
	if diags.HasError() {
		return diags
	}

	s.L3FirewallRules = l3RuleValue

	// l7
	l7rules := []NetworksGroupPolicyResourceModelL7FirewallRule{}
	for _, rule := range f.L7FirewallRules {
		l7rule := NetworksGroupPolicyResourceModelL7FirewallRule{}
		diags = l7rule.FromAPIResponse(ctx, &rule)
		if diags.HasError() {
			return diags
		}
		l7rules = append(l7rules, l7rule)
	}
	l7ObjectType := types.ObjectType{
		AttrTypes: NetworksGroupPolicyResourceModelL7FirewallRuleAttrTypes(),
	}
	l7RuleValue, diags := types.ListValueFrom(ctx, l7ObjectType, l7rules)
	if diags.HasError() {
		return diags
	}
	s.L7FirewallRules = l7RuleValue

	// traffic
	trafficShapingRules := []NetworksGroupPolicyResourceModelTrafficShapingRule{}
	for _, rule := range f.TrafficShapingRules {
		trafficShapingRule := NetworksGroupPolicyResourceModelTrafficShapingRule{}
		diags = trafficShapingRule.FromAPIResponse(ctx, &rule)
		if diags.HasError() {
			return diags
		}

		trafficShapingRules = append(trafficShapingRules, trafficShapingRule)
	}

	trafficShapingRulesObject := types.ObjectType{
		AttrTypes: NetworksGroupPolicyResourceModelTrafficShapingRuleAttrTypes(),
	}

	trafficShapingRuleValue, diags := types.ListValueFrom(ctx, trafficShapingRulesObject, trafficShapingRules)
	if diags.HasError() {
		return diags
	}

	s.TrafficShapingRules = trafficShapingRuleValue

	return diags
}

func (s NetworksGroupPolicyResourceModelFirewallAndTrafficShaping) ToAPIPayload(ctx context.Context) (*openApiClient.CreateNetworkGroupPolicyRequestFirewallAndTrafficShaping, diag.Diagnostics) {
	payload := openApiClient.NewCreateNetworkGroupPolicyRequestFirewallAndTrafficShaping()

	var l3Rules []NetworksGroupPolicyResourceModelL3FirewallRule

	err := s.L3FirewallRules.ElementsAs(ctx, &l3Rules, false)
	if err != nil {
		return nil, err.Errors()
	}

	for _, rule := range l3Rules {
		var irule openApiClient.CreateNetworkGroupPolicyRequestFirewallAndTrafficShapingL3FirewallRulesInner

		irule.SetPolicy(rule.Policy.ValueString())
		irule.SetDestCidr(rule.DestCidr.ValueString())
		irule.SetDestPort(rule.DestPort.ValueString())
		irule.SetComment(rule.Comment.ValueString())
		irule.SetProtocol(rule.Protocol.ValueString())

		payload.L3FirewallRules = append(payload.L3FirewallRules, irule)
	}

	var l7Rules []NetworksGroupPolicyResourceModelL7FirewallRule

	err = s.L7FirewallRules.ElementsAs(ctx, &l7Rules, false)
	if err != nil {
		return nil, err.Errors()
	}

	for _, rule := range l7Rules {
		var irule openApiClient.CreateNetworkGroupPolicyRequestFirewallAndTrafficShapingL7FirewallRulesInner

		irule.SetPolicy(rule.Policy.ValueString())
		irule.SetValue(rule.Value.ValueString())
		irule.SetType(rule.Type.ValueString())

		payload.L7FirewallRules = append(payload.L7FirewallRules, irule)
	}

	var trafficRule []NetworksGroupPolicyResourceModelTrafficShapingRule

	err = s.TrafficShapingRules.ElementsAs(ctx, &trafficRule, false)
	if err != nil {
		return nil, err.Errors()
	}

	for _, rule := range trafficRule {
		var irule openApiClient.CreateNetworkGroupPolicyRequestFirewallAndTrafficShapingTrafficShapingRulesInner

		irule.SetPcpTagValue(int32(rule.PcpTagValue.ValueInt64()))
		irule.SetDscpTagValue(int32(rule.DscpTagValue.ValueInt64()))

		var definitions []NetworksGroupPolicyResourceModelDefinition

		err = rule.Definitions.ElementsAs(ctx, &definitions, false)
		if err != nil {
			return nil, err.Errors()
		}

		for _, definition := range definitions {
			var idefinition openApiClient.UpdateNetworkApplianceTrafficShapingRulesRequestRulesInnerDefinitionsInner

			idefinition.SetType(definition.Type.ValueString())
			idefinition.SetValue(definition.Value.ValueString())

			irule.Definitions = append(irule.Definitions, idefinition)
		}

		var perClientBandwidth openApiClient.UpdateNetworkApplianceTrafficShapingRulesRequestRulesInnerPerClientBandwidthLimits
		var perclientObject NetworksGroupPolicyResourceModelPerClientBandwidthLimits
		diagnostics := rule.PerClientBandwidthLimits.As(ctx, &perclientObject, basetypes.ObjectAsOptions{})
		if diagnostics.HasError() {
			return nil, diagnostics
		}
		perClientBandwidth.SetSettings(perclientObject.Settings.ValueString())

		var perclientBandwidthLimitsObject openApiClient.UpdateNetworkApplianceTrafficShapingRulesRequestRulesInnerPerClientBandwidthLimitsBandwidthLimits
		var perclientBandwidthObject NetworksGroupPolicyResourceModelBandwidthLimits
		diagnostics = perclientObject.BandwidthLimits.As(ctx, &perclientBandwidthObject, basetypes.ObjectAsOptions{})
		if diagnostics.HasError() {
			return nil, diagnostics
		}

		perclientBandwidthLimitsObject.SetLimitUp(int32(perclientBandwidthObject.LimitUp.ValueInt64()))
		perclientBandwidthLimitsObject.SetLimitDown(int32(perclientBandwidthObject.LimitDown.ValueInt64()))

		irule.SetPerClientBandwidthLimits(perClientBandwidth)

		payload.TrafficShapingRules = append(payload.TrafficShapingRules, irule)
	}

	tflog.Info(ctx, "[finish] NetworksGroupPolicyResourceModelL3FirewallRule ToAPIPayload")
	return payload, nil
}

type NetworksGroupPolicyResourceModelL3FirewallRule struct {
	Comment  types.String `tfsdk:"comment" json:"comment"`
	DestCidr types.String `tfsdk:"dest_cidr" json:"destCidr"`
	DestPort types.String `tfsdk:"dest_port" json:"destPort"`
	Policy   types.String `tfsdk:"policy" json:"policy"`
	Protocol types.String `tfsdk:"protocol" json:"protocol"`
}

func (r NetworksGroupPolicyResourceModelL3FirewallRule) FromAPIResponse(ctx context.Context, rule *L3FirewallRule) diag.Diagnostics {
	iRule := NetworksGroupPolicyResourceModelL3FirewallRule{}
	iRule.DestPort = types.StringValue(rule.DestPort)
	iRule.Comment = types.StringValue(rule.Comment)
	iRule.Policy = types.StringValue(rule.Policy)
	iRule.DestCidr = types.StringValue(rule.DestCidr)
	iRule.Protocol = types.StringValue(rule.Protocol)

	return diag.Diagnostics{}
}

func NetworksGroupPolicyResourceModelL3FirewallRuleAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"comment":   types.StringType,
		"dest_cidr": types.StringType,
		"dest_port": types.StringType,
		"policy":    types.StringType,
		"protocol":  types.StringType,
	}
}

type NetworksGroupPolicyResourceModelL7FirewallRule struct {
	Value  types.String `tfsdk:"value" json:"value"`
	Type   types.String `tfsdk:"type" json:"type"`
	Policy types.String `tfsdk:"policy" json:"policy"`
}

func (r NetworksGroupPolicyResourceModelL7FirewallRule) FromAPIResponse(ctx context.Context, rule *L7FirewallRule) diag.Diagnostics {
	iRule := NetworksGroupPolicyResourceModelL7FirewallRule{}
	iRule.Policy = types.StringValue(rule.Policy)
	iRule.Value = types.StringValue(rule.Value)
	iRule.Type = types.StringValue(rule.Type)

	return diag.Diagnostics{}
}

func NetworksGroupPolicyResourceModelL7FirewallRuleAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"value":  types.StringType,
		"type":   types.StringType,
		"policy": types.StringType,
	}
}

type NetworksGroupPolicyResourceModelTrafficShapingRule struct {
	DscpTagValue             types.Int64  `tfsdk:"dscp_tag_value" json:"dscpTagValue"`
	PcpTagValue              types.Int64  `tfsdk:"pcp_tag_value" json:"pcpTagValue"`
	PerClientBandwidthLimits types.Object `tfsdk:"per_client_bandwidth_limits" json:"perClientBandwidthLimits,omitempty"`
	Definitions              types.List   `tfsdk:"definitions" json:"definitions"`
}

func (r NetworksGroupPolicyResourceModelTrafficShapingRule) FromAPIResponse(ctx context.Context, rule *TrafficShapingRule) diag.Diagnostics {
	trafficShappingRule := NetworksGroupPolicyResourceModelTrafficShapingRule{}
	// pcp
	trafficShappingRule.PcpTagValue = types.Int64Value(int64(rule.PcpTagValue))
	// dscp
	trafficShappingRule.DscpTagValue = types.Int64Value(int64(rule.DscpTagValue))

	// per client bandwidth
	perClientBandwidth := NetworksGroupPolicyResourceModelPerClientBandwidthLimits{}
	diags := perClientBandwidth.FromAPIResponse(ctx, rule.PerClientBandwidthLimits)
	if diags.HasError() {
		return diags
	}

	perClientBandwidthValue, diags := types.ObjectValueFrom(ctx, NetworksGroupPolicyResourceModelPerClientBandwidthLimitsAttrTypes(), perClientBandwidth)
	if diags.HasError() {
		return diags
	}

	r.PerClientBandwidthLimits = perClientBandwidthValue

	// definitions
	definitions := []NetworksGroupPolicyResourceModelDefinition{}
	for _, definition := range rule.Definitions {
		iDefinition := NetworksGroupPolicyResourceModelDefinition{}
		diags := iDefinition.FromAPIResponse(ctx, definition)
		if diags.HasError() {
			return diags
		}

		definitions = append(definitions, iDefinition)
	}

	definitionsObjectType := types.ObjectType{
		AttrTypes: NetworksGroupPolicyResourceModelDefinitionAttrTypes(),
	}
	definitionsValue, diags := types.ListValueFrom(ctx, definitionsObjectType, definitions)
	if diags.HasError() {
		return diags
	}
	r.Definitions = definitionsValue

	return diags
}

func NetworksGroupPolicyResourceModelTrafficShapingRuleAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"dscp_tag_value":              types.Int64Type,
		"pcp_tag_value":               types.Int64Type,
		"per_client_bandwidth_limits": types.ObjectType{AttrTypes: NetworksGroupPolicyResourceModelPerClientBandwidthLimitsAttrTypes()},
		"definitions":                 types.ListType{ElemType: types.ObjectType{AttrTypes: NetworksGroupPolicyResourceModelDefinitionAttrTypes()}},
	}
}

type NetworksGroupPolicyResourceModelPerClientBandwidthLimits struct {
	BandwidthLimits types.Object `tfsdk:"bandwidth_limits" json:"bandwidthLimits,,omitempty"`
	Settings        types.String `tfsdk:"settings" json:"settings,,omitempty"`
}

func (l NetworksGroupPolicyResourceModelPerClientBandwidthLimits) FromAPIResponse(ctx context.Context, limits PerClientBandwidthLimits) diag.Diagnostics {
	l.Settings = types.StringValue(limits.Settings)

	bandwidthLimits := NetworksGroupPolicyResourceModelBandwidthLimits{}
	diags := bandwidthLimits.FromAPIResponse(ctx, &limits.BandwidthLimits)

	bandwidthLimitsValue, diags := types.ObjectValueFrom(ctx, NetworksGroupPolicyResourceModelBandwidthLimitsAttrTypes(), bandwidthLimits)
	if diags.HasError() {
		return diags
	}

	l.BandwidthLimits = bandwidthLimitsValue
	return diags
}

func NetworksGroupPolicyResourceModelPerClientBandwidthLimitsAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"bandwidth_limits": types.ObjectType{AttrTypes: NetworksGroupPolicyResourceModelBandwidthLimitsAttrTypes()},
		"settings":         types.StringType,
	}
}

type NetworksGroupPolicyResourceModelDefinition struct {
	Value types.String `tfsdk:"value" json:"value"`
	Type  types.String `tfsdk:"type" json:"type"`
}

func (d NetworksGroupPolicyResourceModelDefinition) FromAPIResponse(ctx context.Context, definition Definition) diag.Diagnostics {
	d.Type = types.StringValue(definition.Type)
	d.Value = types.StringValue(definition.Value)
	return diag.Diagnostics{}
}

func NetworksGroupPolicyResourceModelDefinitionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"value": types.StringType,
		"type":  types.StringType,
	}
}

type NetworksGroupPolicyResourceModelBandwidth struct {
	BandwidthLimits types.Object `tfsdk:"bandwidth_limits" json:"bandwidthLimits"`
	Settings        types.String `tfsdk:"settings" json:"settings"`
}

func NetworksGroupPolicyResourceModelBandwidthAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"bandwidth_limits": types.ObjectType{AttrTypes: NetworksGroupPolicyResourceModelBandwidthLimitsAttrTypes()},
		"settings":         types.StringType,
	}
}

type NetworksGroupPolicyResourceModelBandwidthLimits struct {
	LimitUp   types.Int64 `tfsdk:"limit_up" json:"limitUp"`
	LimitDown types.Int64 `tfsdk:"limit_down" json:"limitDown"`
}

func NetworksGroupPolicyResourceModelBandwidthLimitsAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"limit_up":   types.Int64Type,
		"limit_down": types.Int64Type,
	}
}

type NetworksGroupPolicyResourceModelBonjourForwarding struct {
	BonjourForwardingSettings types.String `tfsdk:"settings" json:"settings"`
	BonjourForwardingRules    types.List   `tfsdk:"rules" json:"rules"`
}

func NetworksGroupPolicyResourceModelBonjourForwardingAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"settings": types.StringType,
		"rules":    types.SetType{ElemType: types.ObjectType{AttrTypes: NetworksGroupPolicyResourceModelRuleAttrTypes()}},
	}
}

func (f NetworksGroupPolicyResourceModelBonjourForwarding) FromAPIResponse(ctx context.Context, forwarding BonjourForwarding) diag.Diagnostics {
	var bonjourRules []NetworksGroupPolicyResourceModelRule
	for _, iRule := range forwarding.Rules {
		var rule NetworksGroupPolicyResourceModelRule
		diags := rule.FromAPIResponse(ctx, &iRule)
		if diags.HasError() {
			tflog.Warn(ctx, "failed to extract FromAPIResponse to PrefixAssignments")
			return diags
		}
		bonjourRules = append(bonjourRules, rule)
	}

	p, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: NetworksGroupPolicyResourceModelRuleAttrTypes()}, bonjourRules)

	f.BonjourForwardingRules = p
	f.BonjourForwardingSettings = types.StringValue(forwarding.Settings)
	return diag.Diagnostics{}
}

func (f NetworksGroupPolicyResourceModelBonjourForwarding) ToAPIPayload(ctx context.Context) (*openApiClient.CreateNetworkGroupPolicyRequestBonjourForwarding, diag.Diagnostics) {
	payload := openApiClient.NewCreateNetworkGroupPolicyRequestBonjourForwarding()
	payload.SetSettings(f.BonjourForwardingSettings.ValueString())

	bonjourForwardingRules := []NetworksGroupPolicyResourceModelRule{}
	diags := f.BonjourForwardingRules.ElementsAs(ctx, &bonjourForwardingRules, false)
	if diags.HasError() {
		return nil, diags
	}

	rulesPayload := []openApiClient.CreateNetworkGroupPolicyRequestBonjourForwardingRulesInner{}
	for _, rule := range bonjourForwardingRules {
		rulePayload := openApiClient.NewCreateNetworkGroupPolicyRequestBonjourForwardingRulesInnerWithDefaults()
		rulePayload.SetDescription(rule.Description.ValueString())
		rulePayload.SetVlanId(rule.VlanId.ValueString())

		var services []string
		diags := rule.Services.ElementsAs(ctx, &services, false)
		if diags.HasError() {
			return nil, diags
		}
		rulePayload.SetServices(services)

		rulesPayload = append(rulesPayload, *rulePayload)
	}
	payload.SetRules(rulesPayload)

	return payload, diags
}

func (n *NetworksGroupPolicyResourceModelBandwidth) FromAPIResponse(ctx context.Context, apiResponse *Bandwidth) diag.Diagnostics {
	n.Settings = types.StringValue(apiResponse.Settings)

	bandwidthLimits := NetworksGroupPolicyResourceModelBandwidthLimits{}
	diags := bandwidthLimits.FromAPIResponse(ctx, &apiResponse.BandwidthLimits)
	value, diags := types.ObjectValueFrom(ctx, NetworksGroupPolicyResourceModelBandwidthLimitsAttrTypes(), bandwidthLimits)
	n.BandwidthLimits = value
	return diags
}

func (n *NetworksGroupPolicyResourceModelBandwidth) ToAPIPayload(ctx context.Context) (*openApiClient.CreateNetworkGroupPolicyRequestBandwidth, diag.Diagnostics) {
	payload := openApiClient.NewCreateNetworkGroupPolicyRequestBandwidth()

	payload.SetSettings(n.Settings.ValueString())

	bandwidthLimitsPayload := openApiClient.NewCreateNetworkGroupPolicyRequestBandwidthBandwidthLimits()
	var bandwidthLimits NetworksGroupPolicyResourceModelBandwidthLimits
	diagnostics := n.BandwidthLimits.As(ctx, &bandwidthLimits, basetypes.ObjectAsOptions{})
	if diagnostics.HasError() {
		return nil, diagnostics
	}

	bandwidthLimitsPayload.SetLimitUp(int32(bandwidthLimits.LimitUp.ValueInt64()))
	bandwidthLimitsPayload.SetLimitDown(int32(bandwidthLimits.LimitDown.ValueInt64()))

	payload.SetBandwidthLimits(*bandwidthLimitsPayload)
	return nil, diagnostics
}

func (n *NetworksGroupPolicyResourceModelBandwidthLimits) FromAPIResponse(ctx context.Context, apiResponse *BandwidthLimits) diag.Diagnostics {
	n.LimitUp = types.Int64Value(int64(apiResponse.LimitUp))
	n.LimitDown = types.Int64Value(int64(apiResponse.LimitDown))
	return diag.Diagnostics{}
}

type NetworksGroupPolicyResourceModelRule struct {
	Description types.String `tfsdk:"description" json:"description"`
	VlanId      types.String `tfsdk:"vlan_id" json:"vlanId"`
	Services    types.Set    `tfsdk:"services" json:"services"`
}

func NetworksGroupPolicyResourceModelRuleAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"description": types.StringType,
		"vlanId":      types.StringType,
		"services":    types.SetType{ElemType: types.StringType},
	}
}

func (r NetworksGroupPolicyResourceModelRule) FromAPIResponse(ctx context.Context, s *Rule) diag.Diagnostics {
	r.Description = types.StringValue(s.Description)
	r.VlanId = types.StringValue(s.VlanId)
	var services []attr.Value
	for _, value := range s.Services {
		services = append(services, types.StringValue(value))
	}
	value, diagnostics := types.SetValue(types.StringType, services)
	if diagnostics.HasError() {
		return diagnostics
	}
	r.Services = value
	return nil
}

func (s NetworksGroupPolicyResourceModelScheduling) FromAPIResponse(ctx context.Context, scheduling *Scheduling) diag.Diagnostics {
	saturday := NetworksGroupPolicyResourceModelSchedule{}
	diags := saturday.FromAPIResponse(ctx, &scheduling.Saturday)
	if diags.HasError() {
		return diags
	}
	value, diags := types.ObjectValueFrom(ctx, NetworksGroupPolicyResourceModelScheduleAttrTypes(), saturday)
	s.Saturday = value

	sunday := NetworksGroupPolicyResourceModelSchedule{}
	diags = sunday.FromAPIResponse(ctx, &scheduling.Sunday)
	if diags.HasError() {
		return diags
	}
	value, diags = types.ObjectValueFrom(ctx, NetworksGroupPolicyResourceModelScheduleAttrTypes(), sunday)
	if diags.HasError() {
		return diags
	}
	s.Sunday = value

	monday := NetworksGroupPolicyResourceModelSchedule{}
	diags = monday.FromAPIResponse(ctx, &scheduling.Monday)
	if diags.HasError() {
		return diags
	}
	value, diags = types.ObjectValueFrom(ctx, NetworksGroupPolicyResourceModelScheduleAttrTypes(), monday)
	if diags.HasError() {
		return diags
	}
	s.Monday = value

	tuesday := NetworksGroupPolicyResourceModelSchedule{}
	diags = tuesday.FromAPIResponse(ctx, &scheduling.Tuesday)
	if diags.HasError() {
		return diags
	}
	value, diags = types.ObjectValueFrom(ctx, NetworksGroupPolicyResourceModelScheduleAttrTypes(), tuesday)
	if diags.HasError() {
		return diags
	}
	s.Tuesday = value

	wednesday := NetworksGroupPolicyResourceModelSchedule{}
	diags = wednesday.FromAPIResponse(ctx, &scheduling.Wednesday)
	if diags.HasError() {
		return diags
	}
	value, diags = types.ObjectValueFrom(ctx, NetworksGroupPolicyResourceModelScheduleAttrTypes(), wednesday)
	if diags.HasError() {
		return diags
	}
	s.Wednesday = value

	thursday := NetworksGroupPolicyResourceModelSchedule{}
	diags = thursday.FromAPIResponse(ctx, &scheduling.Thursday)
	if diags.HasError() {
		return diags
	}
	value, diags = types.ObjectValueFrom(ctx, NetworksGroupPolicyResourceModelScheduleAttrTypes(), thursday)
	if diags.HasError() {
		return diags
	}
	s.Thursday = value

	friday := NetworksGroupPolicyResourceModelSchedule{}
	diags = thursday.FromAPIResponse(ctx, &scheduling.Friday)
	if diags.HasError() {
		return diags
	}
	value, diags = types.ObjectValueFrom(ctx, NetworksGroupPolicyResourceModelScheduleAttrTypes(), friday)
	if diags.HasError() {
		return diags
	}
	s.Thursday = value

	return diags
}

func (s NetworksGroupPolicyResourceModelScheduling) ToAPIPayload(ctx context.Context) (*openApiClient.CreateNetworkGroupPolicyRequestScheduling, diag.Diagnostics) {
	payload := openApiClient.NewCreateNetworkGroupPolicyRequestScheduling()
	fridayPayload := openApiClient.NewCreateNetworkGroupPolicyRequestSchedulingFriday()
	friday := NetworksGroupPolicyResourceModelSchedule{}
	diags := s.Friday.As(ctx, &friday, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	fridayPayload.SetFrom(friday.From.ValueString())
	fridayPayload.SetActive(friday.Active.ValueBool())
	fridayPayload.SetTo(friday.To.ValueString())

	saturdayPayload := openApiClient.NewCreateNetworkGroupPolicyRequestSchedulingSaturday()
	saturday := NetworksGroupPolicyResourceModelSchedule{}
	diags = s.Saturday.As(ctx, &saturday, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	saturdayPayload.SetFrom(saturday.From.ValueString())
	saturdayPayload.SetActive(saturday.Active.ValueBool())
	saturdayPayload.SetTo(saturday.To.ValueString())

	sundayPayload := openApiClient.NewCreateNetworkGroupPolicyRequestSchedulingSunday()
	sunday := NetworksGroupPolicyResourceModelSchedule{}
	diags = s.Sunday.As(ctx, &sunday, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	sundayPayload.SetFrom(sunday.From.ValueString())
	sundayPayload.SetActive(sunday.Active.ValueBool())
	sundayPayload.SetTo(sunday.To.ValueString())

	mondayPayload := openApiClient.NewCreateNetworkGroupPolicyRequestSchedulingMonday()
	monday := NetworksGroupPolicyResourceModelSchedule{}
	diags = s.Monday.As(ctx, &monday, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	mondayPayload.SetFrom(monday.From.ValueString())
	mondayPayload.SetActive(monday.Active.ValueBool())
	mondayPayload.SetTo(monday.To.ValueString())

	tuesdayPayload := openApiClient.NewCreateNetworkGroupPolicyRequestSchedulingTuesday()
	tuesday := NetworksGroupPolicyResourceModelSchedule{}
	diags = s.Tuesday.As(ctx, &tuesday, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	tuesdayPayload.SetFrom(tuesday.From.ValueString())
	tuesdayPayload.SetActive(tuesday.Active.ValueBool())
	tuesdayPayload.SetTo(tuesday.To.ValueString())

	wednesdayPayload := openApiClient.NewCreateNetworkGroupPolicyRequestSchedulingWednesday()
	wednesday := NetworksGroupPolicyResourceModelSchedule{}
	diags = s.Wednesday.As(ctx, &wednesday, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	wednesdayPayload.SetFrom(wednesday.From.ValueString())
	wednesdayPayload.SetActive(wednesday.Active.ValueBool())
	wednesdayPayload.SetTo(wednesday.To.ValueString())

	thursdayPayload := openApiClient.NewCreateNetworkGroupPolicyRequestSchedulingThursday()
	thursday := NetworksGroupPolicyResourceModelSchedule{}
	diags = s.Thursday.As(ctx, &thursday, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	thursdayPayload.SetFrom(thursday.From.ValueString())
	thursdayPayload.SetActive(thursday.Active.ValueBool())
	thursdayPayload.SetTo(thursday.To.ValueString())

	payload.SetFriday(*fridayPayload)
	payload.SetSaturday(*saturdayPayload)
	payload.SetSunday(*sundayPayload)
	payload.SetMonday(*mondayPayload)
	payload.SetTuesday(*tuesdayPayload)
	payload.SetWednesday(*wednesdayPayload)
	payload.SetThursday(*thursdayPayload)
	payload.SetEnabled(s.Enabled.ValueBool())

	return payload, diags
}

type NetworksGroupPolicyResourceModelScheduling struct {
	Enabled   types.Bool   `tfsdk:"enabled" json:"enabled"`
	Friday    types.Object `tfsdk:"friday" json:"friday"`
	Monday    types.Object `tfsdk:"monday" json:"monday"`
	Saturday  types.Object `tfsdk:"saturday" json:"saturday"`
	Sunday    types.Object `tfsdk:"sunday" json:"sunday"`
	Thursday  types.Object `tfsdk:"thursday" json:"thursday"`
	Tuesday   types.Object `tfsdk:"tuesday" json:"tuesday"`
	Wednesday types.Object `tfsdk:"wednesday" json:"wednesday"`
}

func NetworksGroupPolicyResourceModelSchedulingAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"enabled":   types.BoolType,
		"saturday":  types.ObjectType{AttrTypes: NetworksGroupPolicyResourceModelScheduleAttrTypes()},
		"sunday":    types.ObjectType{AttrTypes: NetworksGroupPolicyResourceModelScheduleAttrTypes()},
		"friday":    types.ObjectType{AttrTypes: NetworksGroupPolicyResourceModelScheduleAttrTypes()},
		"monday":    types.ObjectType{AttrTypes: NetworksGroupPolicyResourceModelScheduleAttrTypes()},
		"tuesday":   types.ObjectType{AttrTypes: NetworksGroupPolicyResourceModelScheduleAttrTypes()},
		"wednesday": types.ObjectType{AttrTypes: NetworksGroupPolicyResourceModelScheduleAttrTypes()},
		"thursday":  types.ObjectType{AttrTypes: NetworksGroupPolicyResourceModelScheduleAttrTypes()},
	}
}

type NetworksGroupPolicyResourceModelSchedule struct {
	From   types.String `tfsdk:"from" json:"from"`
	To     types.String `tfsdk:"to" json:"to"`
	Active types.Bool   `tfsdk:"active" json:"active"`
}

func (s NetworksGroupPolicyResourceModelSchedule) FromAPIResponse(ctx context.Context, d *Day) diag.Diagnostics {
	s.To = types.StringValue(d.To)
	s.From = types.StringValue(d.From)
	s.Active = types.BoolValue(d.Active)

	return diag.Diagnostics{}
}

func NetworksGroupPolicyResourceModelScheduleAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"from":   types.StringType,
		"to":     types.StringType,
		"active": types.BoolType,
	}
}

type NetworksGroupPolicyResourceModelVlanTagging struct {
	Settings types.String `tfsdk:"settings" json:"settings"`
	VlanId   types.String `tfsdk:"vlan_id" json:"vlanId"`
}

func (t NetworksGroupPolicyResourceModelVlanTagging) ToAPIPayload(ctx context.Context) (*openApiClient.CreateNetworkGroupPolicyRequestVlanTagging, diag.Diagnostics) {
	tagging := openApiClient.NewCreateNetworkGroupPolicyRequestVlanTagging()

	tagging.SetVlanId(t.VlanId.ValueString())
	tagging.SetSettings(t.Settings.ValueString())

	return tagging, diag.Diagnostics{}
}

type NetworksGroupPolicyResourceModelContentFiltering struct {
	AllowedUrlPatterns   types.Object `tfsdk:"allowed_url_patterns" json:"allowedUrlPatterns"`
	BlockedUrlCategories types.Object `tfsdk:"blocked_url_categories" json:"blockedUrlCategories"`
	BlockedUrlPatterns   types.Object `tfsdk:"blocked_url_patterns" json:"blockedUrlPatterns"`
}

func NetworksGroupPolicyResourceModelContentFilteringAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"allowed_url_patterns":   types.ObjectType{AttrTypes: NetworksGroupPolicyResourceModelAllowedUrlPatternsAttrTypes()},
		"blocked_url_categories": types.ObjectType{AttrTypes: NetworksGroupPolicyResourceModelBlockedUrlCategoriesAttrTypes()},
		"blocked_url_patterns":   types.ObjectType{AttrTypes: NetworksGroupPolicyResourceModelBlockedUrlPatternsAttrTypes()},
	}
}

func (f NetworksGroupPolicyResourceModelContentFiltering) FromAPIResponse(ctx context.Context, c *ContentFiltering) diag.Diagnostics {
	allowedPatterns := NetworksGroupPolicyResourceModelAllowedUrlPatterns{}
	diags := allowedPatterns.FromAPIResponse(ctx, &c.AllowedUrlPatterns)

	if diags.HasError() {
		return diags
	}

	allowedPatternsValue, diags := types.ObjectValueFrom(ctx, NetworksGroupPolicyResourceModelAllowedUrlPatternsAttrTypes(), allowedPatterns)
	if diags.HasError() {
		return diags
	}

	blockedURLCategories := NetworksGroupPolicyResourceModelBlockedUrlCategories{}
	diags = blockedURLCategories.FromAPIResponse(ctx, &c.BlockedUrlCategories)

	if diags.HasError() {
		return diags
	}

	blockedURLCategoriesValue, diags := types.ObjectValueFrom(ctx, NetworksGroupPolicyResourceModelBlockedUrlCategoriesAttrTypes(), blockedURLCategories)
	if diags.HasError() {
		return diags
	}

	blockedPatterns := NetworksGroupPolicyResourceModelBlockedUrlPatterns{}
	diags = blockedPatterns.FromAPIResponse(ctx, &c.BlockedUrlPatterns)

	if diags.HasError() {
		return diags
	}

	blockedURLPatternsValue, diags := types.ObjectValueFrom(ctx, NetworksGroupPolicyResourceModelBlockedUrlPatternsAttrTypes(), blockedPatterns)
	if diags.HasError() {
		return diags
	}

	f.BlockedUrlCategories = blockedURLCategoriesValue
	f.AllowedUrlPatterns = allowedPatternsValue
	f.BlockedUrlPatterns = blockedURLPatternsValue
	return diags
}

func (f NetworksGroupPolicyResourceModelContentFiltering) ToAPIPayload(ctx context.Context) (*openApiClient.CreateNetworkGroupPolicyRequestContentFiltering, diag.Diagnostics) {
	contentFilteringPayload := openApiClient.NewCreateNetworkGroupPolicyRequestContentFiltering()
	blockedCategoriesPayload := openApiClient.NewCreateNetworkGroupPolicyRequestContentFilteringBlockedUrlCategories()
	var blockedCategories NetworksGroupPolicyResourceModelBlockedUrlCategories
	diags := f.BlockedUrlCategories.As(ctx, &blockedCategories, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	blockedCategoriesPayload.SetSettings(blockedCategories.Settings.ValueString())

	var categories []string
	diags = blockedCategories.Categories.ElementsAs(ctx, &categories, false)
	if diags.HasError() {
		return nil, diags
	}

	blockedCategoriesPayload.SetCategories(categories)

	// blocked patterns
	blockedPatternsPayload := openApiClient.NewCreateNetworkGroupPolicyRequestContentFilteringBlockedUrlPatterns()
	var blockedPatterns NetworksGroupPolicyResourceModelBlockedUrlPatterns
	diags = f.BlockedUrlPatterns.As(ctx, &blockedPatterns, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	blockedPatternsPayload.SetSettings(blockedPatterns.Settings.ValueString())

	var patterns []string
	diags = blockedPatterns.Patterns.ElementsAs(ctx, &patterns, false)
	if diags.HasError() {
		return nil, diags
	}

	blockedPatternsPayload.SetPatterns(patterns)

	// blocked patterns
	allowedURLPatternsPayload := openApiClient.NewCreateNetworkGroupPolicyRequestContentFilteringAllowedUrlPatterns()
	var allowedURLPatterns NetworksGroupPolicyResourceModelBlockedUrlPatterns
	diags = f.BlockedUrlPatterns.As(ctx, &allowedURLPatterns, basetypes.ObjectAsOptions{})
	if diags.HasError() {
		return nil, diags
	}

	allowedURLPatternsPayload.SetSettings(allowedURLPatterns.Settings.ValueString())

	var allowdPatterns []string
	diags = allowedURLPatterns.Patterns.ElementsAs(ctx, &allowdPatterns, false)
	if diags.HasError() {
		return nil, diags
	}

	allowedURLPatternsPayload.SetPatterns(patterns)

	contentFilteringPayload.SetBlockedUrlPatterns(*blockedPatternsPayload)
	contentFilteringPayload.SetAllowedUrlPatterns(*allowedURLPatternsPayload)
	contentFilteringPayload.SetBlockedUrlCategories(*blockedCategoriesPayload)

	return contentFilteringPayload, diags
}

type NetworksGroupPolicyResourceModelAllowedUrlPatterns struct {
	Settings types.String `tfsdk:"settings" json:"settings"`
	Patterns types.Set    `tfsdk:"patterns" json:"patterns"`
}

func (p NetworksGroupPolicyResourceModelAllowedUrlPatterns) FromAPIResponse(ctx context.Context, a *AllowedUrlPatterns) diag.Diagnostics {
	p.Settings = types.StringValue(a.Settings)

	values := []attr.Value{}

	for _, value := range a.Patterns {
		values = append(values, types.StringValue(value))
	}

	setValue, diags := types.SetValueFrom(ctx, types.StringType, values)
	if diags.HasError() {
		return diags
	}
	p.Patterns = setValue

	return diags
}

func NetworksGroupPolicyResourceModelAllowedUrlPatternsAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"settings": types.StringType,
		"patterns": types.SetType{ElemType: types.StringType},
	}
}

type NetworksGroupPolicyResourceModelBlockedUrlCategories struct {
	Settings   types.String `tfsdk:"settings" json:"settings"`
	Categories types.Set    `tfsdk:"categories" json:"categories"`
}

func (c NetworksGroupPolicyResourceModelBlockedUrlCategories) FromAPIResponse(ctx context.Context, b *BlockedUrlCategories) diag.Diagnostics {
	c.Settings = types.StringValue(b.Settings)

	values := []attr.Value{}

	for _, value := range b.Categories {
		values = append(values, types.StringValue(value))
	}

	setValue, diags := types.SetValueFrom(ctx, types.StringType, values)
	if diags.HasError() {
		return diags
	}
	c.Categories = setValue

	return diags
}

func NetworksGroupPolicyResourceModelBlockedUrlCategoriesAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"settings":   types.StringType,
		"categories": types.SetType{ElemType: types.StringType},
	}
}

type NetworksGroupPolicyResourceModelBlockedUrlPatterns struct {
	Settings types.String `tfsdk:"settings" json:"settings"`
	Patterns types.Set    `tfsdk:"patterns" json:"patterns"`
}

func (p NetworksGroupPolicyResourceModelBlockedUrlPatterns) FromAPIResponse(ctx context.Context, b *BlockedUrlPatterns) diag.Diagnostics {
	p.Settings = types.StringValue(b.Settings)

	values := []attr.Value{}

	for _, value := range b.Patterns {
		values = append(values, types.StringValue(value))
	}

	setValue, diags := types.SetValueFrom(ctx, types.StringType, values)
	if diags.HasError() {
		return diags
	}
	p.Patterns = setValue

	return diags
}

func NetworksGroupPolicyResourceModelBlockedUrlPatternsAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"settings": types.StringType,
		"patterns": types.SetType{ElemType: types.StringType},
	}
}

func (r *NetworksGroupPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_group_policy"
}

func (r *NetworksGroupPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "NetworksGroupPolicy resource for creating updating and deleting networks group policy resource.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Example identifier",
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Network Id",
				Required:            true,
				CustomType:          jsontypes.StringType,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 31),
				},
			},
			"group_policy_id": schema.StringAttribute{
				MarkdownDescription: "Group Policy ID",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 31),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of Group Policy",
				Required:            true,
				CustomType:          jsontypes.StringType,
			},
			"splash_auth_settings": schema.StringAttribute{
				MarkdownDescription: "Whether clients bound to your policy will bypass splash authorization or behave according to the network's rules. Can be one of 'network default' or 'bypass'. Only available if your network has a wireless configuration",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"network default", "bypass"}...),
				},
			},
			"bandwidth": schema.SingleNestedAttribute{
				MarkdownDescription: "The bandwidth settings for clients bound to your group policy",
				Optional:            true,
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"settings": schema.StringAttribute{
						MarkdownDescription: "How bandwidth limits are enforced. Can be 'network default', 'ignore' or 'custom'.",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.StringType,
					},
					"bandwidth_limits": schema.SingleNestedAttribute{
						MarkdownDescription: "The bandwidth limits object, specifying upload and download speed for clients bound to the group policy. These are only enforced if 'settings' is set to 'custom'.",
						Optional:            true,
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"limit_down": schema.Int64Attribute{
								MarkdownDescription: "The maximum download limit (integer, in Kbps).",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.Int64Type,
							},
							"limit_up": schema.Int64Attribute{
								MarkdownDescription: "The maximum upload limit (integer, in Kbps).",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.Int64Type,
							},
						},
					},
				},
			},
			"bonjour_forwarding": schema.SingleNestedAttribute{
				MarkdownDescription: "The Bonjour settings for your group policy. Only valid if your network has a wireless configuration.",
				Optional:            true,
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"settings": schema.StringAttribute{
						MarkdownDescription: "How Bonjour rules are applied. Can be 'network default', 'ignore' or 'custom'.",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.StringType,
					},
					"rules": schema.SetNestedAttribute{
						Optional: true,
						Computed: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"description": schema.StringAttribute{
									MarkdownDescription: "A description for your Bonjour forwarding rule. Optional.",
									Optional:            true,
									Computed:            true,
									CustomType:          jsontypes.StringType,
								},
								"vlan_id": schema.StringAttribute{
									MarkdownDescription: "The ID of the service VLAN. Required.",
									Optional:            true,
									Computed:            true,
									CustomType:          jsontypes.StringType,
								},
								"services": schema.SetAttribute{
									CustomType: jsontypes.SetType[jsontypes.String](),
									Required:   true,
									Validators: []validator.Set{
										setvalidator.ValueStringsAre(
											stringvalidator.OneOf([]string{"All Services", "AirPlay", "AFP", "BitTorrent", "FTP", "iChat", "iTunes", "Printers", "Samba", "Scanners", "SSH"}...),
											stringvalidator.LengthAtLeast(3),
										),
									},
								},
							},
						},
					},
				},
			},
			"scheduling": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "Whether scheduling is enabled (true) or disabled (false). Defaults to false. If true, the schedule objects for each day of the week (monday - sunday) are parsed.",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.BoolType,
					},
					"friday": schema.SingleNestedAttribute{
						Optional: true,
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"from": schema.StringAttribute{
								MarkdownDescription: "The time, from '00:00' to '24:00'. Must be less than the time specified in 'to'. Defaults to '00:00'. Only 30 minute increments are allowed.",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.StringType,
							},
							"to": schema.StringAttribute{
								MarkdownDescription: "The time, from '00:00' to '24:00'. Must be greater than the time specified in 'from'. Defaults to '24:00'. Only 30 minute increments are allowed.",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.StringType,
							},
							"active": schema.BoolAttribute{
								MarkdownDescription: "Whether the schedule is active (true) or inactive (false) during the time specified between 'from' and 'to'. Defaults to true.",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.BoolType,
							},
						},
					},
					"monday": schema.SingleNestedAttribute{
						Optional: true,
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"from": schema.StringAttribute{
								MarkdownDescription: "The time, from '00:00' to '24:00'. Must be less than the time specified in 'to'. Defaults to '00:00'. Only 30 minute increments are allowed.",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.StringType,
							},
							"to": schema.StringAttribute{
								MarkdownDescription: "The time, from '00:00' to '24:00'. Must be greater than the time specified in 'from'. Defaults to '24:00'. Only 30 minute increments are allowed.",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.StringType,
							},
							"active": schema.BoolAttribute{
								MarkdownDescription: "Whether the schedule is active (true) or inactive (false) during the time specified between 'from' and 'to'. Defaults to true.",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.BoolType,
							},
						},
					},
					"saturday": schema.SingleNestedAttribute{
						Optional: true,
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"from": schema.StringAttribute{
								MarkdownDescription: "The time, from '00:00' to '24:00'. Must be less than the time specified in 'to'. Defaults to '00:00'. Only 30 minute increments are allowed.",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.StringType,
							},
							"to": schema.StringAttribute{
								MarkdownDescription: "The time, from '00:00' to '24:00'. Must be greater than the time specified in 'from'. Defaults to '24:00'. Only 30 minute increments are allowed.",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.StringType,
							},
							"active": schema.BoolAttribute{
								MarkdownDescription: "Whether the schedule is active (true) or inactive (false) during the time specified between 'from' and 'to'. Defaults to true.",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.BoolType,
							},
						},
					},
					"sunday": schema.SingleNestedAttribute{
						Optional: true,
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"from": schema.StringAttribute{
								MarkdownDescription: "The time, from '00:00' to '24:00'. Must be less than the time specified in 'to'. Defaults to '00:00'. Only 30 minute increments are allowed.",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.StringType,
							},
							"to": schema.StringAttribute{
								MarkdownDescription: "The time, from '00:00' to '24:00'. Must be greater than the time specified in 'from'. Defaults to '24:00'. Only 30 minute increments are allowed.",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.StringType,
							},
							"active": schema.BoolAttribute{
								MarkdownDescription: "Whether the schedule is active (true) or inactive (false) during the time specified between 'from' and 'to'. Defaults to true.",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.BoolType,
							},
						},
					},
					"thursday": schema.SingleNestedAttribute{
						Optional: true,
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"from": schema.StringAttribute{
								MarkdownDescription: "The time, from '00:00' to '24:00'. Must be less than the time specified in 'to'. Defaults to '00:00'. Only 30 minute increments are allowed.",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.StringType,
							},
							"to": schema.StringAttribute{
								MarkdownDescription: "The time, from '00:00' to '24:00'. Must be greater than the time specified in 'from'. Defaults to '24:00'. Only 30 minute increments are allowed.",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.StringType,
							},
							"active": schema.BoolAttribute{
								MarkdownDescription: "Whether the schedule is active (true) or inactive (false) during the time specified between 'from' and 'to'. Defaults to true.",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.BoolType,
							},
						},
					},
					"tuesday": schema.SingleNestedAttribute{
						Optional: true,
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"from": schema.StringAttribute{
								MarkdownDescription: "The time, from '00:00' to '24:00'. Must be less than the time specified in 'to'. Defaults to '00:00'. Only 30 minute increments are allowed.",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.StringType,
							},
							"to": schema.StringAttribute{
								MarkdownDescription: "The time, from '00:00' to '24:00'. Must be greater than the time specified in 'from'. Defaults to '24:00'. Only 30 minute increments are allowed.",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.StringType,
							},
							"active": schema.BoolAttribute{
								MarkdownDescription: "Whether the schedule is active (true) or inactive (false) during the time specified between 'from' and 'to'. Defaults to true.",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.BoolType,
							},
						},
					},
					"wednesday": schema.SingleNestedAttribute{
						Optional: true,
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"from": schema.StringAttribute{
								MarkdownDescription: "The time, from '00:00' to '24:00'. Must be less than the time specified in 'to'. Defaults to '00:00'. Only 30 minute increments are allowed.",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.StringType,
							},
							"to": schema.StringAttribute{
								MarkdownDescription: "The time, from '00:00' to '24:00'. Must be greater than the time specified in 'from'. Defaults to '24:00'. Only 30 minute increments are allowed.",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.StringType,
							},
							"active": schema.BoolAttribute{
								MarkdownDescription: "Whether the schedule is active (true) or inactive (false) during the time specified between 'from' and 'to'. Defaults to true.",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.BoolType,
							},
						},
					},
				},
			},
			"content_filtering": schema.SingleNestedAttribute{
				MarkdownDescription: "The content filtering settings for your group policy",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"allowed_url_patterns": schema.SingleNestedAttribute{
						MarkdownDescription: "Settings for allowed URL patterns",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"settings": schema.StringAttribute{
								MarkdownDescription: "How URL patterns are applied. Can be 'network default', 'append' or 'override'.",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.StringType,
							},
							"patterns": schema.SetAttribute{
								MarkdownDescription: "A list of URL patterns that are allowed",
								CustomType:          jsontypes.SetType[jsontypes.String](),
								Optional:            true,
							},
						},
					},
					"blocked_url_categories": schema.SingleNestedAttribute{
						MarkdownDescription: "Settings for blocked URL categories",
						Optional:            true,
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"settings": schema.StringAttribute{
								MarkdownDescription: "How URL categories are applied. Can be 'network default', 'append' or 'override'.",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.StringType,
							},
							"categories": schema.SetAttribute{
								MarkdownDescription: "A list of URL categories to block",
								CustomType:          jsontypes.SetType[jsontypes.String](),
								Optional:            true,
							},
						},
					},
					"blocked_url_patterns": schema.SingleNestedAttribute{
						MarkdownDescription: "Settings for blocked URL patterns",
						Optional:            true,
						Computed:            true,
						Attributes: map[string]schema.Attribute{
							"settings": schema.StringAttribute{
								MarkdownDescription: "How URL patterns are applied. Can be 'network default', 'append' or 'override'.",
								Optional:            true,
								Computed:            true,
								CustomType:          jsontypes.StringType,
							},
							"patterns": schema.SetAttribute{
								MarkdownDescription: "A list of URL patterns that are blocked",
								CustomType:          jsontypes.SetType[jsontypes.String](),
								Optional:            true,
							},
						},
					},
				},
			},
			"firewall_and_traffic_shaping": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"settings": schema.StringAttribute{
						MarkdownDescription: "How firewall and traffic shaping rules are enforced. Can be 'network default', 'ignore' or 'custom'.",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.StringType,
					},
					"l3_firewall_rules": schema.SetNestedAttribute{
						Optional: true,
						Computed: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"comment": schema.StringAttribute{
									MarkdownDescription: "Description of the rule (optional)",
									Optional:            true,
									Computed:            true,
									CustomType:          jsontypes.StringType,
								},
								"dest_cidr": schema.StringAttribute{
									MarkdownDescription: "Destination IP address (in IP or CIDR notation), a fully-qualified domain name (FQDN, if your network supports it) or 'any'.",
									Optional:            true,
									Computed:            true,
									CustomType:          jsontypes.StringType,
								},
								"dest_port": schema.StringAttribute{
									MarkdownDescription: "Destination port (integer in the range 1-65535), a port range (e.g. 8080-9090), or 'any'",
									Optional:            true,
									Computed:            true,
									CustomType:          jsontypes.StringType,
								},
								"policy": schema.StringAttribute{
									MarkdownDescription: "'allow' or 'deny' traffic specified by this rule",
									Optional:            true,
									Computed:            true,
									CustomType:          jsontypes.StringType,
								},
								"protocol": schema.StringAttribute{
									MarkdownDescription: "The type of protocol (must be 'tcp', 'udp', 'icmp', 'icmp6' or 'any')",
									Optional:            true,
									Computed:            true,
									CustomType:          jsontypes.StringType,
								},
							},
						},
					},
					"l7_firewall_rules": schema.SetNestedAttribute{
						Optional: true,
						Computed: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"value": schema.StringAttribute{
									MarkdownDescription: "The 'value' of what you want to block. If 'type' is 'host', 'port' or 'ipRange', 'value' must be a string matching either a hostname (e.g. somewhere.com), a port (e.g. 8080), or an IP range (e.g. 192.1.0.0/16). If 'type' is 'application' or 'applicationCategory', then 'value' must be an object with an ID for the application.",
									Optional:            true,
									Computed:            true,
									CustomType:          jsontypes.StringType,
								},
								"policy": schema.StringAttribute{
									MarkdownDescription: "The policy applied to matching traffic. Must be 'deny'.",
									Optional:            true,
									Computed:            true,
									CustomType:          jsontypes.StringType,
								},
								"type": schema.StringAttribute{
									MarkdownDescription: "Type of the L7 Rule. Must be 'application', 'applicationCategory', 'host', 'port' or 'ipRange'",
									Optional:            true,
									Computed:            true,
									CustomType:          jsontypes.StringType,
								},
							},
						},
					},
					"traffic_shaping_rules": schema.SetNestedAttribute{
						Optional: true,
						Computed: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"dscp_tag_value": schema.Int64Attribute{
									MarkdownDescription: "The DSCP tag applied by your rule. null means Do not change DSCP tag. For a list of possible tag values, use the trafficShaping/dscpTaggingOptions endpoint",
									Optional:            true,
									Computed:            true,
									CustomType:          jsontypes.Int64Type,
								},
								"pcp_tag_value": schema.Int64Attribute{
									MarkdownDescription: "The PCP tag applied by your rule. Can be 0 (lowest priority) through 7 (highest priority). null means Do not set PCP tag.",
									Optional:            true,
									Computed:            true,
									CustomType:          jsontypes.Int64Type,
								},
								"per_client_bandwidth_limits": schema.SingleNestedAttribute{
									Optional: true,
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"settings": schema.StringAttribute{
											MarkdownDescription: "How bandwidth limits are applied by your rule. Can be one of 'network default', 'ignore' or 'custom'.",
											Optional:            true,
											Computed:            true,
											CustomType:          jsontypes.StringType,
										},
										"bandwidth_limits": schema.SingleNestedAttribute{
											Optional: true,
											Computed: true,
											Attributes: map[string]schema.Attribute{
												"limit_down": schema.Int64Attribute{
													MarkdownDescription: "The maximum download limit (integer, in Kbps).",
													Optional:            true,
													Computed:            true,
													CustomType:          jsontypes.Int64Type,
												},
												"limit_up": schema.Int64Attribute{
													MarkdownDescription: "The maximum upload limit (integer, in Kbps).",
													Optional:            true,
													Computed:            true,
													CustomType:          jsontypes.Int64Type,
												},
											},
										},
									},
								},
								"definitions": schema.SetNestedAttribute{
									Optional: true,
									Computed: true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"type": schema.StringAttribute{
												MarkdownDescription: "The type of definition. Can be one of 'application', 'applicationCategory', 'host', 'port', 'ipRange' or 'localNet'.",
												Optional:            true,
												Computed:            true,
												CustomType:          jsontypes.StringType,
											},
											"value": schema.StringAttribute{
												MarkdownDescription: "If type is host, port, ipRange or localNet then value must be a string matching either a hostname (e.g. somesite.com) a port (e.g. 8080) or an IP range (192.1.0.0, 192.1.0.0/16, or 10.1.0.0/16:80). localNet also supports CIDR notation excluding custom ports If type is 'application' or 'applicationCategory', then value must be an object with the structure { id: meraki:layer7/... }, where id is the application category or application ID (for a list of IDs for your network, use the trafficShaping/applicationCategories endpoint)",
												Optional:            true,
												Computed:            true,
												CustomType:          jsontypes.StringType,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"vlan_tagging": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"settings": schema.StringAttribute{
						MarkdownDescription: "How VLAN tagging is applied. Can be 'network default', 'ignore' or 'custom'.",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.StringType,
					},
					"vlan_id": schema.StringAttribute{
						MarkdownDescription: "The ID of the vlan you want to tag. This only applies if 'settings' is set to 'custom'.",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.StringType,
					},
				},
			},
		},
	}
}

func (r *NetworksGroupPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*openApiClient.APIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *NetworksGroupPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *NetworksGroupPolicyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Initial create API call
	payload, payloadReqDiags := CreateGroupPolicyHttpReqPayload(ctx, data)
	if payloadReqDiags != nil {
		resp.Diagnostics.Append(payloadReqDiags...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	inlineResp, httpResp, err := r.client.NetworksApi.CreateNetworkGroupPolicy(ctx, data.NetworkId.ValueString()).CreateNetworkGroupPolicyRequest(payload).Execute()

	// Meraki API seems to return http status code 201 as an error.
	if err != nil && httpResp.StatusCode != 201 {
		resp.Diagnostics.AddError(
			"HTTP Client Create Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
	}

	// Check for API success inlineResp code
	if httpResp.StatusCode != 201 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
	}

	// Check for API success response code
	if httpResp.StatusCode != 201 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%s", data))
		return
	}

	data, err = extractHttpResponseGroupPolicyResource(ctx, inlineResp, httpResp.Body, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"JSON Decode issue",
			fmt.Sprintf("%v", err.Error()),
		)
		return
	}

	data.Id = types.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")
}

func (r *NetworksGroupPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetworksGroupPolicyResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	inlineResp, httpResp, err := r.client.NetworksApi.GetNetworkGroupPolicy(ctx, data.NetworkId.ValueString(), data.GroupPolicyId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
		)
	}

	// Check for API success response code
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	} else {
		resp.Diagnostics.Append()
	}

	data, err = extractHttpResponseGroupPolicyResource(ctx, inlineResp, httpResp.Body, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"JSON Decode issue",
			fmt.Sprintf("%v", err.Error()),
		)
		return
	}

	data.Id = types.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *NetworksGroupPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data *NetworksGroupPolicyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	payload, diagnostics := UpdateNetworkGroupPolicyHttpReqPayload(ctx, data)
	if diagnostics.HasError() {
		return
	}

	inlineResp, httpResp, err := r.client.NetworksApi.UpdateNetworkGroupPolicy(ctx, data.NetworkId.ValueString(), data.GroupPolicyId.ValueString()).UpdateNetworkGroupPolicyRequest(*payload).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
	}

	// Check for API success response code
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%s", data))
		return
	}

	data, err = extractHttpResponseGroupPolicyResource(ctx, inlineResp, httpResp.Body, data)
	if err != nil {
		resp.Diagnostics.AddError(
			"JSON Decode issue",
			fmt.Sprintf("%v", err.Error()),
		)
		return
	}

	data.Id = types.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "update resource")
}

func (r *NetworksGroupPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data *NetworksGroupPolicyResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.client.NetworksApi.DeleteNetworkGroupPolicy(context.Background(), data.NetworkId.ValueString(), data.GroupPolicyId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
	}

	// Check for API success response code
	if httpResp.StatusCode != 204 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	resp.State.RemoveResource(ctx)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")

}

func (r *NetworksGroupPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: network_id, group_policy_id. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("group_policy_id"), idParts[1])...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func extractHttpResponseGroupPolicyResource(ctx context.Context, inlineResp map[string]interface{}, httpRespBody io.ReadCloser, data *NetworksGroupPolicyResourceModel) (*NetworksGroupPolicyResourceModel, error) {
	var networkGroupPolicy *NetworkGroupPolicyResponse
	if err := json.NewDecoder(httpRespBody).Decode(networkGroupPolicy); err != nil {
		return data, err
	}

	//vlan tagging
	vlanTagging := NetworksGroupPolicyResourceModelVlanTagging{
		Settings: types.StringValue(networkGroupPolicy.VlanTagging.VlanId),
		VlanId:   types.StringValue(networkGroupPolicy.VlanTagging.Settings),
	}

	vlanTaggingAttributes := map[string]attr.Type{
		"vlanId":   types.StringType,
		"settings": types.StringType,
	}

	objectVal, diags := types.ObjectValueFrom(ctx, vlanTaggingAttributes, vlanTagging)
	if diags.HasError() {
		return data, nil
	}

	data.VlanTagging = objectVal

	// bonjour forwarding

	bonjourForwarding := NetworksGroupPolicyResourceModelBonjourForwarding{}
	diags = bonjourForwarding.FromAPIResponse(ctx, networkGroupPolicy.BonjourForwarding)
	if diags.HasError() {
		return nil, nil
	}

	bonjourForwardingObject, diags := types.ObjectValueFrom(ctx, NetworksGroupPolicyResourceModelBonjourForwardingAttrTypes(), bonjourForwarding)
	if diags.HasError() {
		return nil, nil
	}

	data.BonjourForwarding = bonjourForwardingObject

	bandwidth := NetworksGroupPolicyResourceModelBandwidth{}
	diags = bandwidth.FromAPIResponse(ctx, &networkGroupPolicy.Bandwidth)
	if diags.HasError() {
		return nil, nil
	}

	bandwidthValue, diags := types.ObjectValueFrom(ctx, NetworksGroupPolicyResourceModelBandwidthAttrTypes(), bandwidth)
	if diags.HasError() {
		return nil, nil
	}

	data.Bandwidth = bandwidthValue

	// scheduling
	scheduling := NetworksGroupPolicyResourceModelScheduling{}
	diags = scheduling.FromAPIResponse(ctx, &networkGroupPolicy.Scheduling)
	if diags.HasError() {
		return nil, nil
	}

	schedulingValue, diags := types.ObjectValueFrom(ctx, NetworksGroupPolicyResourceModelSchedulingAttrTypes(), scheduling)
	if diags.HasError() {
		return nil, nil
	}

	data.Scheduling = schedulingValue

	// content filtering
	contentFiltering := NetworksGroupPolicyResourceModelContentFiltering{}
	diags = contentFiltering.FromAPIResponse(ctx, &networkGroupPolicy.ContentFiltering)

	contentFilteringValue, diags := types.ObjectValueFrom(ctx, NetworksGroupPolicyResourceModelContentFilteringAttrTypes(), contentFiltering)
	if diags.HasError() {
		return nil, nil
	}

	data.ContentFiltering = contentFilteringValue

	// FirewallAndTrafficShaping
	firewallAndTrafficShapping := NetworksGroupPolicyResourceModelFirewallAndTrafficShaping{}
	diags = firewallAndTrafficShapping.FromAPIResponse(ctx, &networkGroupPolicy.FirewallAndTrafficShaping)
	if diags.HasError() {
		return nil, nil
	}

	trafficShapingRuleValue, diags := types.ObjectValueFrom(ctx, NetworksGroupPolicyResourceModelFirewallAndTrafficShapingAttrTypes(), firewallAndTrafficShapping)
	if diags.HasError() {
		return nil, nil
	}

	data.FirewallAndTrafficShaping = trafficShapingRuleValue

	data.GroupPolicyId = types.StringValue(networkGroupPolicy.GroupPolicyId)
	data.Name = types.StringValue(networkGroupPolicy.Name)
	data.SplashAuthSettings = types.StringValue(networkGroupPolicy.SplashAuthSettings)

	return data, nil
}

func CreateGroupPolicyHttpReqPayload(ctx context.Context, data *NetworksGroupPolicyResourceModel) (openApiClient.CreateNetworkGroupPolicyRequest, diag.Diagnostics) {
	resp := diag.Diagnostics{}

	// Log the received request
	tflog.Info(ctx, "[start] Create HTTP Request Payload Call")
	tflog.Trace(ctx, "Create Request Payload", map[string]interface{}{
		"data": data,
	})

	// Initialize the payload
	payload := openApiClient.NewCreateNetworkGroupPolicyRequest(data.Name.ValueString())

	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		// API returns this as string, openAPI spec has set as Integer
		payload.SetName(fmt.Sprintf("%v", data.Name.ValueString()))
	}

	if !data.Bandwidth.IsNull() && !data.Bandwidth.IsUnknown() {
		bandwidth := openApiClient.NewCreateNetworkGroupPolicyRequestBandwidth()
		var ibandwidth NetworksGroupPolicyResourceModelBandwidth

		diags := data.Bandwidth.As(ctx, &ibandwidth, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			resp.AddError(
				"Create Payload Failure", fmt.Sprintf("%v", diags),
			)
		}

		limits := openApiClient.NewCreateNetworkGroupPolicyRequestBandwidthBandwidthLimits()
		var ilimits NetworksGroupPolicyResourceModelBandwidthLimits

		diags = ibandwidth.BandwidthLimits.As(ctx, &ilimits, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			resp.AddError(
				"Create Payload Failure", fmt.Sprintf("%v", diags),
			)
		}
		limits.SetLimitUp(int32(ilimits.LimitUp.ValueInt64()))
		limits.SetLimitDown(int32(ilimits.LimitDown.ValueInt64()))

		bandwidth.SetBandwidthLimits(*limits)
		bandwidth.SetSettings(ibandwidth.Settings.ValueString())

		payload.SetBandwidth(*bandwidth)
	}

	if !data.BonjourForwarding.IsNull() && !data.BonjourForwarding.IsUnknown() {
		bonjourForwarding := openApiClient.NewCreateNetworkGroupPolicyRequestBonjourForwarding()
		var ibonjourForwarding NetworksGroupPolicyResourceModelBonjourForwarding

		diags := data.BonjourForwarding.As(ctx, &ibonjourForwarding, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			resp.AddError(
				"Create Payload Failure", fmt.Sprintf("%v", diags),
			)
		}

		bonjourForwarding.SetSettings(ibonjourForwarding.BonjourForwardingSettings.ValueString())

		var bonjourForwardingNetworkRules []NetworksGroupPolicyResourceModelRule
		diags = data.BonjourForwarding.As(ctx, &bonjourForwardingNetworkRules, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			resp.AddError(
				"Create Payload Failure", fmt.Sprintf("%v", diags),
			)
		}
		var rules []openApiClient.CreateNetworkGroupPolicyRequestBonjourForwardingRulesInner
		for _, value := range bonjourForwardingNetworkRules {
			rule := openApiClient.NewCreateNetworkGroupPolicyRequestBonjourForwardingRulesInnerWithDefaults()
			rule.SetVlanId(value.VlanId.ValueString())
			rule.SetDescription(value.Description.ValueString())
			var services []string
			for _, item := range value.Services.Elements() {
				services = append(services, item.String())
			}
			rule.SetServices(services)
			rules = append(rules, *rule)
		}
		bonjourForwarding.SetRules(rules)
		payload.SetBonjourForwarding(*bonjourForwarding)
	}

	if !data.ContentFiltering.IsNull() && !data.ContentFiltering.IsUnknown() {
		filtering := openApiClient.NewCreateNetworkGroupPolicyRequestContentFiltering()
		var ifiltering NetworksGroupPolicyResourceModelContentFiltering
		diags := data.ContentFiltering.As(ctx, &ifiltering, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			resp.AddError(
				"Create Payload Failure", fmt.Sprintf("%v", diags),
			)
		}

		categories := openApiClient.NewCreateNetworkGroupPolicyRequestContentFilteringBlockedUrlCategories()
		var iblockedURLCategories NetworksGroupPolicyResourceModelBlockedUrlCategories
		diags = ifiltering.BlockedUrlCategories.As(ctx, &iblockedURLCategories, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			resp.AddError(
				"Create Payload Failure", fmt.Sprintf("%v", diags),
			)
		}
		categories.SetSettings(iblockedURLCategories.Settings.ValueString())
		var urls []string
		for _, value := range iblockedURLCategories.Categories.Elements() {
			urls = append(urls, value.String())
		}
		categories.SetCategories(urls)
		filtering.SetBlockedUrlCategories(*categories)

		blockedURLPatterns := openApiClient.NewCreateNetworkGroupPolicyRequestContentFilteringBlockedUrlPatterns()
		var iblockedURLPatterns NetworksGroupPolicyResourceModelBlockedUrlPatterns
		diags = ifiltering.BlockedUrlPatterns.As(ctx, &iblockedURLPatterns, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			resp.AddError(
				"Create Payload Failure", fmt.Sprintf("%v", diags),
			)
		}
		var patterns []string
		for _, value := range iblockedURLPatterns.Patterns.Elements() {
			patterns = append(patterns, value.String())
		}
		blockedURLPatterns.SetPatterns(patterns)
		blockedURLPatterns.SetSettings(iblockedURLPatterns.Settings.ValueString())
		filtering.SetBlockedUrlPatterns(*blockedURLPatterns)

		allowedPatterns := openApiClient.NewCreateNetworkGroupPolicyRequestContentFilteringAllowedUrlPatterns()
		var iallowedPatterns NetworksGroupPolicyResourceModelAllowedUrlPatterns
		diags = ifiltering.AllowedUrlPatterns.As(ctx, &iallowedPatterns, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			resp.AddError(
				"Create Payload Failure", fmt.Sprintf("%v", diags),
			)
		}
		var allowedURLPatterns []string
		for _, value := range iallowedPatterns.Patterns.Elements() {
			allowedURLPatterns = append(allowedURLPatterns, value.String())
		}
		allowedPatterns.SetPatterns(allowedURLPatterns)
		allowedPatterns.SetSettings(iallowedPatterns.Settings.ValueString())
		filtering.SetAllowedUrlPatterns(*allowedPatterns)
		payload.SetContentFiltering(*filtering)
	}

	if !data.FirewallAndTrafficShaping.IsNull() && !data.FirewallAndTrafficShaping.IsUnknown() {
		firewallAndTrafficShapping := openApiClient.NewCreateNetworkGroupPolicyRequestFirewallAndTrafficShaping()
		var ifirewallAndTrafficShapping NetworksGroupPolicyResourceModelFirewallAndTrafficShaping
		diags := data.FirewallAndTrafficShaping.As(ctx, &ifirewallAndTrafficShapping, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			resp.AddError(
				"Create Payload Failure", fmt.Sprintf("%v", diags),
			)
		}

		firewallAndTrafficShapping.SetSettings(ifirewallAndTrafficShapping.Settings.ValueString())

		l3FirewallRulesInners := []openApiClient.CreateNetworkGroupPolicyRequestFirewallAndTrafficShapingL3FirewallRulesInner{}
		for _, l3Value := range ifirewallAndTrafficShapping.L3FirewallRules.Elements() {
			l3FirewallRule := NetworksGroupPolicyResourceModelL3FirewallRule{}
			il3FirewallRule := openApiClient.CreateNetworkGroupPolicyRequestFirewallAndTrafficShapingL3FirewallRulesInner{}

			diags := tfsdk.ValueAs(ctx, l3Value, &l3FirewallRule)
			if diags.HasError() {
				return *payload, diags
			}

			il3FirewallRule.SetPolicy(l3FirewallRule.Policy.ValueString())
			il3FirewallRule.SetProtocol(l3FirewallRule.Protocol.ValueString())
			il3FirewallRule.SetComment(l3FirewallRule.Comment.ValueString())
			il3FirewallRule.SetDestPort(l3FirewallRule.DestPort.ValueString())
			il3FirewallRule.SetDestCidr(l3FirewallRule.DestCidr.ValueString())

			l3FirewallRulesInners = append(l3FirewallRulesInners, il3FirewallRule)
		}
		firewallAndTrafficShapping.SetL3FirewallRules(l3FirewallRulesInners)

		l7FirewallRulesInners := []openApiClient.CreateNetworkGroupPolicyRequestFirewallAndTrafficShapingL7FirewallRulesInner{}
		for _, l7Value := range ifirewallAndTrafficShapping.L7FirewallRules.Elements() {
			l7FirewallRule := NetworksGroupPolicyResourceModelL7FirewallRule{}
			il7FirewallRule := openApiClient.CreateNetworkGroupPolicyRequestFirewallAndTrafficShapingL7FirewallRulesInner{}

			diags := tfsdk.ValueAs(ctx, l7Value, &l7FirewallRule)
			if diags.HasError() {
				return *payload, diags
			}

			il7FirewallRule.SetPolicy(l7FirewallRule.Policy.ValueString())
			il7FirewallRule.SetValue(l7FirewallRule.Value.ValueString())
			il7FirewallRule.SetType(l7FirewallRule.Type.ValueString())

			l7FirewallRulesInners = append(l7FirewallRulesInners, il7FirewallRule)
		}
		firewallAndTrafficShapping.SetL7FirewallRules(l7FirewallRulesInners)

		trafficRulesInner := []openApiClient.CreateNetworkGroupPolicyRequestFirewallAndTrafficShapingTrafficShapingRulesInner{}
		for _, trafficRuleValue := range ifirewallAndTrafficShapping.TrafficShapingRules.Elements() {
			trafficRule := NetworksGroupPolicyResourceModelTrafficShapingRule{}
			iTrafficRule := openApiClient.CreateNetworkGroupPolicyRequestFirewallAndTrafficShapingTrafficShapingRulesInner{}

			diags := tfsdk.ValueAs(ctx, trafficRuleValue, &trafficRule)
			if diags.HasError() {
				return *payload, diags
			}

			iTrafficRule.SetDscpTagValue(int32(trafficRule.DscpTagValue.ValueInt64()))
			iTrafficRule.SetPcpTagValue(int32(trafficRule.PcpTagValue.ValueInt64()))

			definitions := []openApiClient.UpdateNetworkApplianceTrafficShapingRulesRequestRulesInnerDefinitionsInner{}
			for _, definitionValue := range trafficRule.Definitions.Elements() {
				definition := NetworksGroupPolicyResourceModelDefinition{}
				iDefinition := openApiClient.UpdateNetworkApplianceTrafficShapingRulesRequestRulesInnerDefinitionsInner{}

				diags := tfsdk.ValueAs(ctx, definitionValue, &definition)
				if diags.HasError() {
					return *payload, diags
				}

				iDefinition.SetValue(definition.Value.ValueString())
				iDefinition.SetType(definition.Type.ValueString())

				definitions = append(definitions, iDefinition)
			}
			iTrafficRule.SetDefinitions(definitions)

			perClientBandwidthLimits := openApiClient.NewUpdateNetworkApplianceTrafficShapingRulesRequestRulesInnerPerClientBandwidthLimits()
			var iPerclientBandwidths NetworksGroupPolicyResourceModelPerClientBandwidthLimits
			diags = trafficRule.PerClientBandwidthLimits.As(ctx, &iPerclientBandwidths, basetypes.ObjectAsOptions{})
			if diags.HasError() {
				resp.AddError(
					"Create Payload Failure", fmt.Sprintf("%v", diags),
				)
			}

			iBandwidthLimits := NetworksGroupPolicyResourceModelBandwidthLimits{}
			bandwidthLimits := openApiClient.UpdateNetworkApplianceTrafficShapingRulesRequestRulesInnerPerClientBandwidthLimitsBandwidthLimits{}

			diags = iPerclientBandwidths.BandwidthLimits.As(ctx, &iBandwidthLimits, basetypes.ObjectAsOptions{})
			if diags.HasError() {
				resp.AddError(
					"Create Payload Failure", fmt.Sprintf("%v", diags),
				)
			}

			bandwidthLimits.SetLimitUp(int32(iBandwidthLimits.LimitUp.ValueInt64()))
			bandwidthLimits.SetLimitDown(int32(iBandwidthLimits.LimitDown.ValueInt64()))
			perClientBandwidthLimits.SetBandwidthLimits(bandwidthLimits)
			perClientBandwidthLimits.SetSettings(iPerclientBandwidths.Settings.ValueString())

			iTrafficRule.SetPerClientBandwidthLimits(*perClientBandwidthLimits)
			trafficRulesInner = append(trafficRulesInner, iTrafficRule)
		}

		payload.SetFirewallAndTrafficShaping(*firewallAndTrafficShapping)
	}

	if !data.SplashAuthSettings.IsNull() && !data.SplashAuthSettings.IsUnknown() {
		payload.SetSplashAuthSettings(data.SplashAuthSettings.ValueString())
	}

	if !data.VlanTagging.IsNull() && !data.VlanTagging.IsUnknown() {
		vlanTagging := openApiClient.NewCreateNetworkGroupPolicyRequestVlanTagging()
		iVlanTagging := NetworksGroupPolicyResourceModelVlanTagging{}
		diags := data.FirewallAndTrafficShaping.As(ctx, &iVlanTagging, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			resp.AddError(
				"Create Payload Failure", fmt.Sprintf("%v", diags),
			)
		}

		vlanTagging.SetVlanId(iVlanTagging.VlanId.ValueString())
		vlanTagging.SetSettings(iVlanTagging.Settings.ValueString())

		payload.SetVlanTagging(*vlanTagging)
	}

	if !data.Scheduling.IsNull() && !data.Scheduling.IsUnknown() {
		scheduling := openApiClient.NewCreateNetworkGroupPolicyRequestScheduling()
		iScheduling := NetworksGroupPolicyResourceModelScheduling{}
		diags := data.FirewallAndTrafficShaping.As(ctx, &iScheduling, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			resp.AddError(
				"Create Payload Failure", fmt.Sprintf("%v", diags),
			)
		}

		monday := openApiClient.NewCreateNetworkGroupPolicyRequestSchedulingMonday()
		imonday := NetworksGroupPolicyResourceModelSchedule{}
		diags = iScheduling.Monday.As(ctx, &imonday, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			resp.AddError(
				"Create Payload Failure", fmt.Sprintf("%v", diags),
			)
		}
		monday.SetTo(imonday.To.ValueString())
		monday.SetFrom(imonday.From.ValueString())
		monday.SetActive(imonday.Active.ValueBool())

		tuesday := openApiClient.NewCreateNetworkGroupPolicyRequestSchedulingTuesday()
		ituesday := NetworksGroupPolicyResourceModelSchedule{}
		diags = iScheduling.Monday.As(ctx, &ituesday, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			resp.AddError(
				"Create Payload Failure", fmt.Sprintf("%v", diags),
			)
		}
		tuesday.SetTo(ituesday.To.ValueString())
		tuesday.SetFrom(ituesday.From.ValueString())
		tuesday.SetActive(ituesday.Active.ValueBool())

		wednesday := openApiClient.NewCreateNetworkGroupPolicyRequestSchedulingWednesday()
		iwednesday := NetworksGroupPolicyResourceModelSchedule{}
		diags = iScheduling.Monday.As(ctx, &iwednesday, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			resp.AddError(
				"Create Payload Failure", fmt.Sprintf("%v", diags),
			)
		}
		wednesday.SetTo(iwednesday.To.ValueString())
		wednesday.SetFrom(iwednesday.From.ValueString())
		wednesday.SetActive(iwednesday.Active.ValueBool())

		thursday := openApiClient.NewCreateNetworkGroupPolicyRequestSchedulingThursday()
		ithursday := NetworksGroupPolicyResourceModelSchedule{}
		diags = iScheduling.Monday.As(ctx, &ithursday, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			resp.AddError(
				"Create Payload Failure", fmt.Sprintf("%v", diags),
			)
		}
		thursday.SetTo(ithursday.To.ValueString())
		thursday.SetFrom(ithursday.From.ValueString())
		thursday.SetActive(ithursday.Active.ValueBool())

		friday := openApiClient.NewCreateNetworkGroupPolicyRequestSchedulingFriday()
		ifriday := NetworksGroupPolicyResourceModelSchedule{}
		diags = iScheduling.Monday.As(ctx, &ifriday, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			resp.AddError(
				"Create Payload Failure", fmt.Sprintf("%v", diags),
			)
		}
		friday.SetTo(ifriday.To.ValueString())
		friday.SetFrom(ifriday.From.ValueString())
		friday.SetActive(ifriday.Active.ValueBool())

		saturday := openApiClient.NewCreateNetworkGroupPolicyRequestSchedulingSaturday()
		isaturday := NetworksGroupPolicyResourceModelSchedule{}
		diags = iScheduling.Monday.As(ctx, &isaturday, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			resp.AddError(
				"Create Payload Failure", fmt.Sprintf("%v", diags),
			)
		}
		saturday.SetTo(isaturday.To.ValueString())
		saturday.SetFrom(isaturday.From.ValueString())
		saturday.SetActive(isaturday.Active.ValueBool())

		sunday := openApiClient.NewCreateNetworkGroupPolicyRequestSchedulingSunday()
		isunday := NetworksGroupPolicyResourceModelSchedule{}
		diags = iScheduling.Monday.As(ctx, &isunday, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			resp.AddError(
				"Create Payload Failure", fmt.Sprintf("%v", diags),
			)
		}
		sunday.SetTo(isunday.To.ValueString())
		sunday.SetFrom(isunday.From.ValueString())
		sunday.SetActive(isunday.Active.ValueBool())

		scheduling.SetMonday(*monday)
		scheduling.SetTuesday(*tuesday)
		scheduling.SetWednesday(*wednesday)
		scheduling.SetThursday(*thursday)
		scheduling.SetFriday(*friday)
		scheduling.SetSaturday(*saturday)
		scheduling.SetSunday(*sunday)

		payload.SetScheduling(*scheduling)
	}

	return *payload, nil
}

func UpdateNetworkGroupPolicyHttpReqPayload(ctx context.Context, data *NetworksGroupPolicyResourceModel) (*openApiClient.UpdateNetworkGroupPolicyRequest, diag.Diagnostics) {
	resp := diag.Diagnostics{}

	tflog.Info(ctx, "[start] UpdateNetworkGroupPolicyHttpReqPayload Call")

	payload := openApiClient.NewUpdateNetworkGroupPolicyRequest()

	// Name
	if !data.Name.IsUnknown() && !data.Name.IsNull() {
		payload.SetName(data.Name.ValueString())
	}

	if !data.FirewallAndTrafficShaping.IsUnknown() && !data.FirewallAndTrafficShaping.IsNull() {
		var firewallAndTrafficData NetworksGroupPolicyResourceModelFirewallAndTrafficShaping

		diags := data.FirewallAndTrafficShaping.As(ctx, &firewallAndTrafficData, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			resp.AddError(
				"Update Payload Failure", fmt.Sprintf("%v", diags),
			)
		}

		firewallAndTrafficPayload, firewallAndTrafficDiags := firewallAndTrafficData.ToAPIPayload(ctx)
		if firewallAndTrafficDiags.HasError() {
			return payload, firewallAndTrafficDiags
		}

		payload.SetFirewallAndTrafficShaping(*firewallAndTrafficPayload)
	}

	if !data.Bandwidth.IsUnknown() && !data.Bandwidth.IsNull() {
		var bandwitdh NetworksGroupPolicyResourceModelBandwidth

		diags := data.Bandwidth.As(ctx, &bandwitdh, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			resp.AddError(
				"Update Payload Failure", fmt.Sprintf("%v", diags),
			)
		}

		bandwidth, bandwidthDiags := bandwitdh.ToAPIPayload(ctx)
		if bandwidthDiags.HasError() {
			return payload, bandwidthDiags
		}

		payload.SetBandwidth(*bandwidth)
	}

	if !data.ContentFiltering.IsUnknown() && !data.ContentFiltering.IsNull() {
		var contentFiltering NetworksGroupPolicyResourceModelContentFiltering

		diags := data.ContentFiltering.As(ctx, &contentFiltering, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			resp.AddError(
				"Update Payload Failure", fmt.Sprintf("%v", diags),
			)
		}

		contentFilteringPayload, contentFilteringDiags := contentFiltering.ToAPIPayload(ctx)
		if contentFilteringDiags.HasError() {
			return payload, contentFilteringDiags
		}

		payload.SetContentFiltering(*contentFilteringPayload)
	}

	// bonjour forwarding
	if !data.BonjourForwarding.IsUnknown() && !data.BonjourForwarding.IsNull() {
		var bonjourForwarding NetworksGroupPolicyResourceModelBonjourForwarding

		diags := data.BonjourForwarding.As(ctx, &bonjourForwarding, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			resp.AddError(
				"Update Payload Failure", fmt.Sprintf("%v", diags),
			)
		}

		bonjourForwardingPayload, bonjourForwardingDiags := bonjourForwarding.ToAPIPayload(ctx)
		if bonjourForwardingDiags.HasError() {
			return payload, bonjourForwardingDiags
		}

		payload.SetBonjourForwarding(*bonjourForwardingPayload)
	}

	// vlan tagging
	if !data.VlanTagging.IsUnknown() && !data.VlanTagging.IsNull() {
		var vlanTagging NetworksGroupPolicyResourceModelVlanTagging

		diags := data.VlanTagging.As(ctx, &vlanTagging, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			resp.AddError(
				"Update Payload Failure", fmt.Sprintf("%v", diags),
			)
		}

		vlanTaggingPayload, vlanTaggingDiags := vlanTagging.ToAPIPayload(ctx)
		if vlanTaggingDiags.HasError() {
			return payload, vlanTaggingDiags
		}

		payload.SetVlanTagging(*vlanTaggingPayload)
	}

	// scheduling
	if !data.Scheduling.IsUnknown() && !data.Scheduling.IsNull() {
		var scheduling NetworksGroupPolicyResourceModelScheduling

		diags := data.Scheduling.As(ctx, &scheduling, basetypes.ObjectAsOptions{})
		if diags.HasError() {
			resp.AddError(
				"Update Payload Failure", fmt.Sprintf("%v", diags),
			)
		}

		schedulingPayload, schedulingDiags := scheduling.ToAPIPayload(ctx)
		if schedulingDiags.HasError() {
			return payload, schedulingDiags
		}

		payload.SetScheduling(*schedulingPayload)
	}

	// splash settings
	if !data.SplashAuthSettings.IsUnknown() && !data.SplashAuthSettings.IsNull() {
		payload.SetSplashAuthSettings(data.SplashAuthSettings.ValueString())
	}

	return payload, diag.Diagnostics{}
}
