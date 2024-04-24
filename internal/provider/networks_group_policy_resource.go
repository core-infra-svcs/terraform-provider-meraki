package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

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

// NetworksGroupPolicyResourceModel describes the resource data model.
type NetworksGroupPolicyResourceModel struct {
	Id                        jsontypes.String `tfsdk:"id"`
	NetworkId                 jsontypes.String `tfsdk:"network_id"`
	GroupPolicyId             jsontypes.String `tfsdk:"group_policy_id" json:"groupPolicyId"`
	Name                      jsontypes.String `tfsdk:"name" json:"name"`
	SplashAuthSettings        jsontypes.String `tfsdk:"splash_auth_settings" json:"splashAuthSettings"`
	Bandwidth                 types.Object     `tfsdk:"bandwidth" json:"bandwidth"`
	BonjourForwarding         types.Object     `tfsdk:"bonjour_forwarding" json:"bonjourForwarding"`
	Scheduling                types.Object     `tfsdk:"scheduling" json:"scheduling"`
	ContentFiltering          types.Object     `tfsdk:"content_filtering" json:"contentFiltering"`
	FirewallAndTrafficShaping types.Object     `tfsdk:"firewall_and_traffic_shaping" json:"firewallAndTrafficShaping"`
	VlanTagging               types.Object     `tfsdk:"vlan_tagging" json:"vlanTagging"`
}

type NetworksGroupPolicyResourceModelFirewallAndTrafficShaping struct {
	Settings            types.String `tfsdk:"settings" json:"settings"`
	L3FirewallRules     types.List   `tfsdk:"l3_firewall_rules" json:"l3FirewallRules"`
	L7FirewallRules     types.List   `tfsdk:"l7_firewall_rules" json:"l7FirewallRules"`
	TrafficShapingRules types.List   `tfsdk:"traffic_shaping_rules" json:"trafficShapingRules"`
}

type NetworksGroupPolicyResourceModelL3FirewallRule struct {
	Comment  types.String `tfsdk:"comment" json:"comment"`
	DestCidr types.String `tfsdk:"dest_cidr" json:"destCidr"`
	DestPort types.String `tfsdk:"dest_port" json:"destPort"`
	Policy   types.String `tfsdk:"policy" json:"policy"`
	Protocol types.String `tfsdk:"protocol" json:"protocol"`
}

type NetworksGroupPolicyResourceModelL7FirewallRule struct {
	Value  types.String `tfsdk:"value" json:"value"`
	Type   types.String `tfsdk:"type" json:"type"`
	Policy types.String `tfsdk:"policy" json:"policy"`
}

type NetworksGroupPolicyResourceModelTrafficShapingRule struct {
	DscpTagValue             types.Int64  `tfsdk:"dscp_tag_value" json:"dscpTagValue"`
	PcpTagValue              types.Int64  `tfsdk:"pcp_tag_value" json:"pcpTagValue"`
	PerClientBandwidthLimits types.Object `tfsdk:"per_client_bandwidth_limits" json:"perClientBandwidthLimits,,omitempty"`
	Definitions              types.List   `tfsdk:"definitions" json:"definitions"`
}

type NetworksGroupPolicyResourceModelPerClientBandwidth struct {
	BandwidthLimits types.Object `tfsdk:"bandwidth_limits" json:"bandwidthLimits,,omitempty"`
	Settings        types.String `tfsdk:"settings" json:"settings,,omitempty"`
}

type NetworksGroupPolicyResourceModelPerClientBandwidthLimits struct {
	LimitUp   types.Int64 `tfsdk:"limit_up" json:"limitUp"`
	LimitDown types.Int64 `tfsdk:"limit_down" json:"limitDown"`
}

type NetworksGroupPolicyResourceModelDefinition struct {
	Value types.String `tfsdk:"value" json:"value"`
	Type  types.String `tfsdk:"type" json:"type"`
}

type NetworksGroupPolicyResourceModelBandwidth struct {
	BandwidthLimits NetworksGroupPolicyResourceModelBandwidthLimits `tfsdk:"bandwidth_limits" json:"bandwidthLimits"`
	Settings        jsontypes.String                                `tfsdk:"settings" json:"settings"`
}

type NetworksGroupPolicyResourceModelBandwidthLimits struct {
	LimitUp   jsontypes.Int64 `tfsdk:"limit_up" json:"limitUp"`
	LimitDown jsontypes.Int64 `tfsdk:"limit_down" json:"limitDown"`
}

type NetworksGroupPolicyResourceModelBonjourForwarding struct {
	BonjourForwardingSettings types.String `tfsdk:"settings" json:"settings"`
	BonjourForwardingRules    types.List   `tfsdk:"rules" json:"rules"`
}

type NetworksGroupPolicyResourceModelRule struct {
	Description types.String `tfsdk:"description" json:"description"`
	VlanId      types.String `tfsdk:"vlan_id" json:"vlanId"`
	Services    types.Set    `tfsdk:"services" json:"services"`
}

type OutputNetworksGroupPolicyResourceModelBonjourForwarding struct {
	BonjourForwardingSettings string                                       `json:"settings"`
	BonjourForwardingRules    []OutputNetworksGroupPolicyResourceModelRule `json:"rules"`
}

type OutputNetworksGroupPolicyResourceModelRule struct {
	Description string   `json:"description"`
	VlanId      string   `json:"vlanId"`
	Services    []string `json:"services"`
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

type NetworksGroupPolicyResourceModelSchedule struct {
	From   types.String `tfsdk:"from" json:"from"`
	To     types.String `tfsdk:"to" json:"to"`
	Active types.Bool   `tfsdk:"active" json:"active"`
}

type NetworksGroupPolicyResourceModelVlanTagging struct {
	Settings types.String `tfsdk:"settings" json:"settings"`
	VlanId   types.String `tfsdk:"vlan_id" json:"vlanId"`
}

type OutputNetworksGroupPolicyResourceModelVlanTagging struct {
	Settings string `json:"settings"`
	VlanId   string `json:"vlanId"`
}

type NetworksGroupPolicyResourceModelContentFiltering struct {
	AllowedUrlPatterns   types.Object `tfsdk:"allowed_url_patterns" json:"allowedUrlPatterns"`
	BlockedUrlCategories types.Object `tfsdk:"blocked_url_categories" json:"blockedUrlCategories"`
	BlockedUrlPatterns   types.Object `tfsdk:"blocked_url_patterns" json:"blockedUrlPatterns"`
}

type NetworksGroupPolicyResourceModelAllowedUrlPatterns struct {
	Settings types.String `tfsdk:"settings" json:"settings"`
	Patterns types.List   `tfsdk:"patterns" json:"patterns"`
}
type NetworksGroupPolicyResourceModelBlockedUrlCategories struct {
	Settings   types.String `tfsdk:"settings" json:"settings"`
	Categories types.List   `tfsdk:"categories" json:"categories"`
}
type NetworksGroupPolicyResourceModelBlockedUrlPatterns struct {
	Settings types.String `tfsdk:"settings" json:"settings"`
	Patterns types.List   `tfsdk:"patterns" json:"patterns"`
}

type OutputNetworksGroupPolicyResourceModelContentFiltering struct {
	AllowedUrlPatterns   OutputNetworksGroupPolicyResourceModelAllowedUrlPatterns   `json:"allowedUrlPatterns"`
	BlockedUrlCategories OutputNetworksGroupPolicyResourceModelBlockedUrlCategories `json:"blockedUrlCategories"`
	BlockedUrlPatterns   OutputNetworksGroupPolicyResourceModelBlockedUrlPatterns   `json:"blockedUrlPatterns"`
}

type OutputNetworksGroupPolicyResourceModelAllowedUrlPatterns struct {
	Settings string   `json:"settings"`
	Patterns []string `json:"patterns"`
}
type OutputNetworksGroupPolicyResourceModelBlockedUrlCategories struct {
	Settings   string   `json:"settings"`
	Categories []string `json:"categories"`
}
type OutputNetworksGroupPolicyResourceModelBlockedUrlPatterns struct {
	Settings string   `json:"settings"`
	Patterns []string `json:"patterns"`
}

type OutputNetworksGroupPolicyResourceModelScheduling struct {
	Enabled   bool                                           `json:"enabled"`
	Friday    OutputNetworksGroupPolicyResourceModelSchedule `json:"friday"`
	Monday    OutputNetworksGroupPolicyResourceModelSchedule `json:"monday"`
	Saturday  OutputNetworksGroupPolicyResourceModelSchedule `json:"saturday"`
	Sunday    OutputNetworksGroupPolicyResourceModelSchedule `json:"sunday"`
	Thursday  OutputNetworksGroupPolicyResourceModelSchedule `json:"thursday"`
	Tuesday   OutputNetworksGroupPolicyResourceModelSchedule `json:"tuesday"`
	Wednesday OutputNetworksGroupPolicyResourceModelSchedule `json:"wednesday"`
}

type OutputNetworksGroupPolicyResourceModelSchedule struct {
	From   string `json:"from"`
	To     string `json:"to"`
	Active bool   `json:"active"`
}

type OutputNetworksGroupPolicyResourceModelPerClientBandwidthLimits struct {
	LimitUp   int64 `json:"limitUp"`
	LimitDown int64 `json:"limitDown"`
}

type OutputNetworksGroupPolicyResourceModelFirewallAndTrafficShaping struct {
	Settings            string                                                     `json:"settings"`
	L3FirewallRules     []OutputNetworksGroupPolicyResourceModelL3FirewallRule     `json:"l3FirewallRules"`
	L7FirewallRules     []OutputNetworksGroupPolicyResourceModelL7FirewallRule     `json:"l7FirewallRules"`
	TrafficShapingRules []OutputNetworksGroupPolicyResourceModelTrafficShapingRule `json:"trafficShapingRules"`
}

type OutputNetworksGroupPolicyResourceModelL3FirewallRule struct {
	Comment  string `json:"comment"`
	DestCidr string `json:"destCidr"`
	DestPort string `json:"destPort"`
	Policy   string `json:"policy"`
	Protocol string `json:"protocol"`
}

type OutputNetworksGroupPolicyResourceModelL7FirewallRule struct {
	Value  string `json:"value"`
	Type   string `json:"type"`
	Policy string `json:"policy"`
}

type OutputNetworksGroupPolicyResourceModelTrafficShapingRule struct {
	DscpTagValue             int64                                              `json:"dscpTagValue"`
	PcpTagValue              int64                                              `json:"pcpTagValue"`
	PerClientBandwidthLimits OutputNetworksGroupPolicyResourceModelBandwidth    `json:"perClientBandwidthLimits,,omitempty"`
	Definitions              []OutputNetworksGroupPolicyResourceModelDefinition `json:"definitions"`
}

type OutputNetworksGroupPolicyResourceModelDefinition struct {
	Value string `json:"value"`
	Type  string `json:"type"`
}

type OutputNetworksGroupPolicyResourceModelBandwidth struct {
	BandwidthLimits OutputNetworksGroupPolicyResourceModelPerClientBandwidthLimits `json:"bandwidthLimits"`
	Settings        string                                                         `json:"settings"`
}

func (r *NetworksGroupPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_group_policy"
}

func BandwidthLimitsAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"limit_up":   types.Int64Type,
		"limit_down": types.Int64Type,
	}
}

func PerClientBandwidthAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"bandwidth_limits": types.ObjectType{
			AttrTypes: BandwidthLimitsAttrTypes(),
		},
		"settings": types.StringType,
	}
}
func DefinitionAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"type":  types.StringType,
		"value": types.StringType,
	}
}

func VlanTagging() map[string]attr.Type {
	return map[string]attr.Type{
		"settings": types.StringType,
		"vlan_id":  types.StringType,
	}
}

func L3FirewallRulesAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"comment":   types.StringType,
		"dest_cidr": types.StringType,
		"dest_port": types.StringType,
		"policy":    types.StringType,
		"protocol":  types.StringType,
	}
}

func L7FirewallRulesAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"value":  types.StringType,
		"policy": types.StringType,
		"type":   types.StringType,
	}
}

func FirewallAndTrafficShapingAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"settings":              types.StringType,
		"l3_firewall_rules":     types.ListType{ElemType: types.ObjectType{AttrTypes: L3FirewallRulesAttrTypes()}},
		"l7_firewall_rules":     types.ListType{ElemType: types.ObjectType{AttrTypes: L7FirewallRulesAttrTypes()}},
		"traffic_shaping_rules": types.ListType{ElemType: types.ObjectType{AttrTypes: TrafficShapingRulesAttrTypes()}},
	}
}

func SchedulingAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"enabled": types.BoolType,
		"friday": types.ObjectType{
			AttrTypes: SchedulingDataAttrTypes(),
		},
		"monday": types.ObjectType{
			AttrTypes: SchedulingDataAttrTypes(),
		},
		"tuesday": types.ObjectType{
			AttrTypes: SchedulingDataAttrTypes(),
		},
		"wednesday": types.ObjectType{
			AttrTypes: SchedulingDataAttrTypes(),
		},
		"thursday": types.ObjectType{
			AttrTypes: SchedulingDataAttrTypes(),
		},
		"saturday": types.ObjectType{
			AttrTypes: SchedulingDataAttrTypes(),
		},
		"sunday": types.ObjectType{
			AttrTypes: SchedulingDataAttrTypes(),
		},
	}

}

func SchedulingDataAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"from":   types.StringType,
		"to":     types.StringType,
		"active": types.BoolType,
	}

}

func TrafficShapingRulesAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"dscp_tag_value": types.Int64Type,
		"pcp_tag_value":  types.Int64Type,
		"per_client_bandwidth_limits": types.ObjectType{
			AttrTypes: PerClientBandwidthAttrTypes(),
		},
		"definitions": types.ListType{ElemType: types.ObjectType{AttrTypes: DefinitionAttrTypes()}},
	}
}

func PatternsDataAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"settings": types.StringType,
		"patterns": types.SetType{ElemType: types.StringType},
	}

}

func CategoriesDataAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"settings":   types.StringType,
		"categories": types.SetType{ElemType: types.StringType},
	}

}

func ContentFilteringDataAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"allowed_url_patterns": types.ObjectType{
			AttrTypes: PatternsDataAttrTypes(),
		},
		"blocked_url_categories": types.ObjectType{
			AttrTypes: CategoriesDataAttrTypes(),
		},
		"blocked_url_patterns": types.ObjectType{
			AttrTypes: PatternsDataAttrTypes(),
		},
	}

}

func BonjourForwardingDataAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"settings": types.StringType,
		"rules":    types.ListType{ElemType: types.ObjectType{AttrTypes: BonjourForwardingRuleDataAttrTypes()}},
	}
}

func BonjourForwardingRuleDataAttrTypes() map[string]attr.Type {
	return map[string]attr.Type{
		"description": types.StringType,
		"vlan_id":     types.StringType,
		"services":    types.SetType{ElemType: types.StringType},
	}
}

func getBandwidth(ctx context.Context, inlineResp map[string]interface{}) (basetypes.ObjectValue, diag.Diagnostics) {

	bandwidthAttrTypes := map[string]attr.Type{
		"bandwidth_limits": types.ObjectType{
			AttrTypes: BandwidthLimitsAttrTypes(),
		},
		"settings": types.StringType,
	}

	var bandwidthData OutputNetworksGroupPolicyResourceModelBandwidth
	jsonData, _ := json.Marshal(inlineResp["bandwidth"].(map[string]interface{}))
	json.Unmarshal(jsonData, &bandwidthData)

	bandwidthLimitsMap, diags := basetypes.NewObjectValue(BandwidthLimitsAttrTypes(), map[string]attr.Value{
		"limit_up":   basetypes.NewInt64Value(bandwidthData.BandwidthLimits.LimitUp),
		"limit_down": basetypes.NewInt64Value(bandwidthData.BandwidthLimits.LimitDown),
	})
	if diags.HasError() {
		return basetypes.ObjectValue{}, diags
	}

	bandwidthMap, _ := basetypes.NewObjectValue(bandwidthAttrTypes, map[string]attr.Value{
		"bandwidth_limits": bandwidthLimitsMap,
		"settings":         basetypes.NewStringValue(bandwidthData.Settings),
	})
	objectVal, diags := types.ObjectValueFrom(ctx, bandwidthAttrTypes, bandwidthMap)
	if diags.HasError() {
		return basetypes.ObjectValue{}, diags
	}

	return objectVal, nil

}

func getVlantagging(ctx context.Context, inlineResp map[string]interface{}) (basetypes.ObjectValue, diag.Diagnostics) {

	var vlantaggingOutput OutputNetworksGroupPolicyResourceModelVlanTagging
	jsonData, _ := json.Marshal(inlineResp["vlanTagging"].(map[string]interface{}))
	json.Unmarshal(jsonData, &vlantaggingOutput)

	vlantaggingMap, _ := basetypes.NewObjectValue(VlanTagging(), map[string]attr.Value{
		"settings": basetypes.NewStringValue(vlantaggingOutput.Settings),
		"vlan_id":  basetypes.NewStringValue(vlantaggingOutput.VlanId),
	})
	vlantaggingVal, diags := types.ObjectValueFrom(ctx, VlanTagging(), vlantaggingMap)
	if diags.HasError() {
		return basetypes.ObjectValue{}, diags
	}
	return vlantaggingVal, nil
}

func getBonjourForwarding(ctx context.Context, inlineResp map[string]interface{}) (basetypes.ObjectValue, diag.Diagnostics) {
	var bonjourForwardingData OutputNetworksGroupPolicyResourceModelBonjourForwarding
	jsonData, _ := json.Marshal(inlineResp["bonjourForwarding"].(map[string]interface{}))
	json.Unmarshal(jsonData, &bonjourForwardingData)

	var bonjourForwardingDataRules []basetypes.ObjectValue
	for _, rule := range bonjourForwardingData.BonjourForwardingRules {

		servicesList, diags := types.SetValueFrom(ctx, types.StringType, rule.Services)
		if diags.HasError() {
			return basetypes.ObjectValue{}, diags
		}

		ruleDataMap, _ := basetypes.NewObjectValue(BonjourForwardingRuleDataAttrTypes(), map[string]attr.Value{
			"description": basetypes.NewStringValue(rule.Description),
			"vlan_id":     basetypes.NewStringValue(rule.VlanId),
			"services":    servicesList,
		})

		objectVal, diags := types.ObjectValueFrom(ctx, BonjourForwardingRuleDataAttrTypes(), ruleDataMap)
		if diags.HasError() {
			return basetypes.ObjectValue{}, diags
		}
		bonjourForwardingDataRules = append(bonjourForwardingDataRules, objectVal)
	}

	bonjourForwardingDataRulesList, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: BonjourForwardingRuleDataAttrTypes()}, bonjourForwardingDataRules)

	bonjourForwardingDataMap, diags := basetypes.NewObjectValue(BonjourForwardingDataAttrTypes(), map[string]attr.Value{
		"settings": basetypes.NewStringValue(bonjourForwardingData.BonjourForwardingSettings),
		"rules":    bonjourForwardingDataRulesList,
	})
	if diags.HasError() {
		return basetypes.ObjectValue{}, diags
	}

	bonjourForwardingVal, diags := types.ObjectValueFrom(ctx, BonjourForwardingDataAttrTypes(), bonjourForwardingDataMap)
	if diags.HasError() {
		return basetypes.ObjectValue{}, diags
	}
	return bonjourForwardingVal, nil
}

func getScheduling(ctx context.Context, inlineResp map[string]interface{}) (basetypes.ObjectValue, diag.Diagnostics) {

	var schedulingData OutputNetworksGroupPolicyResourceModelScheduling
	jsonData, _ := json.Marshal(inlineResp["scheduling"].(map[string]interface{}))
	json.Unmarshal(jsonData, &schedulingData)

	fridayMap, diags := basetypes.NewObjectValue(SchedulingDataAttrTypes(), map[string]attr.Value{
		"active": basetypes.NewBoolValue(schedulingData.Friday.Active),
		"to":     basetypes.NewStringValue(schedulingData.Friday.To),
		"from":   basetypes.NewStringValue(schedulingData.Friday.From),
	})
	if diags.HasError() {
		return basetypes.ObjectValue{}, diags
	}

	mondayMap, diags := basetypes.NewObjectValue(SchedulingDataAttrTypes(), map[string]attr.Value{
		"active": basetypes.NewBoolValue(schedulingData.Monday.Active),
		"to":     basetypes.NewStringValue(schedulingData.Monday.To),
		"from":   basetypes.NewStringValue(schedulingData.Monday.From),
	})
	if diags.HasError() {
		return basetypes.ObjectValue{}, diags
	}

	tuesdayMap, diags := basetypes.NewObjectValue(SchedulingDataAttrTypes(), map[string]attr.Value{
		"active": basetypes.NewBoolValue(schedulingData.Tuesday.Active),
		"to":     basetypes.NewStringValue(schedulingData.Tuesday.To),
		"from":   basetypes.NewStringValue(schedulingData.Tuesday.From),
	})
	if diags.HasError() {
		return basetypes.ObjectValue{}, diags
	}

	wednesdayMap, diags := basetypes.NewObjectValue(SchedulingDataAttrTypes(), map[string]attr.Value{
		"active": basetypes.NewBoolValue(schedulingData.Wednesday.Active),
		"to":     basetypes.NewStringValue(schedulingData.Wednesday.To),
		"from":   basetypes.NewStringValue(schedulingData.Wednesday.From),
	})
	if diags.HasError() {
		return basetypes.ObjectValue{}, diags
	}

	thursdayMap, diags := basetypes.NewObjectValue(SchedulingDataAttrTypes(), map[string]attr.Value{
		"active": basetypes.NewBoolValue(schedulingData.Thursday.Active),
		"to":     basetypes.NewStringValue(schedulingData.Thursday.To),
		"from":   basetypes.NewStringValue(schedulingData.Thursday.From),
	})
	if diags.HasError() {
		return basetypes.ObjectValue{}, diags
	}

	saturdayMap, diags := basetypes.NewObjectValue(SchedulingDataAttrTypes(), map[string]attr.Value{
		"active": basetypes.NewBoolValue(schedulingData.Saturday.Active),
		"to":     basetypes.NewStringValue(schedulingData.Saturday.To),
		"from":   basetypes.NewStringValue(schedulingData.Saturday.From),
	})
	if diags.HasError() {
		return basetypes.ObjectValue{}, diags
	}

	sundayMap, diags := basetypes.NewObjectValue(SchedulingDataAttrTypes(), map[string]attr.Value{
		"active": basetypes.NewBoolValue(schedulingData.Sunday.Active),
		"to":     basetypes.NewStringValue(schedulingData.Sunday.To),
		"from":   basetypes.NewStringValue(schedulingData.Sunday.From),
	})
	if diags.HasError() {
		return basetypes.ObjectValue{}, diags
	}

	schedulingMap, diags := basetypes.NewObjectValue(SchedulingAttrTypes(), map[string]attr.Value{
		"enabled":   basetypes.NewBoolValue(schedulingData.Enabled),
		"friday":    fridayMap,
		"monday":    mondayMap,
		"tuesday":   tuesdayMap,
		"wednesday": wednesdayMap,
		"thursday":  thursdayMap,
		"saturday":  saturdayMap,
		"sunday":    sundayMap,
	})
	if diags.HasError() {
		return basetypes.ObjectValue{}, diags
	}

	schedulingVal, diags := types.ObjectValueFrom(ctx, SchedulingAttrTypes(), schedulingMap)
	if diags.HasError() {
		return basetypes.ObjectValue{}, diags
	}

	return schedulingVal, nil
}

func getContentFiltering(ctx context.Context, inlineResp map[string]interface{}) (basetypes.ObjectValue, diag.Diagnostics) {

	var contentFilteringData OutputNetworksGroupPolicyResourceModelContentFiltering
	jsonData, _ := json.Marshal(inlineResp["contentFiltering"].(map[string]interface{}))
	json.Unmarshal(jsonData, &contentFilteringData)

	allowedpatternsList, diags := types.SetValueFrom(ctx, types.StringType, contentFilteringData.AllowedUrlPatterns.Patterns)
	if diags.HasError() {
		return basetypes.ObjectValue{}, diags
	}
	allowedpatternsDataMap, diags := basetypes.NewObjectValue(PatternsDataAttrTypes(), map[string]attr.Value{
		"settings": basetypes.NewStringValue(contentFilteringData.AllowedUrlPatterns.Settings),
		"patterns": allowedpatternsList,
	})
	if diags.HasError() {
		return basetypes.ObjectValue{}, diags
	}

	blockedcategoriesList, diags := types.SetValueFrom(ctx, types.StringType, contentFilteringData.BlockedUrlCategories.Categories)
	if diags.HasError() {
		return basetypes.ObjectValue{}, diags
	}
	blockedcategoriesDataMap, diags := basetypes.NewObjectValue(CategoriesDataAttrTypes(), map[string]attr.Value{
		"settings":   basetypes.NewStringValue(contentFilteringData.BlockedUrlCategories.Settings),
		"categories": blockedcategoriesList,
	})
	if diags.HasError() {
		return basetypes.ObjectValue{}, diags
	}

	blockedurlpatternsList, diags := types.SetValueFrom(ctx, types.StringType, contentFilteringData.BlockedUrlPatterns.Patterns)
	if diags.HasError() {
		return basetypes.ObjectValue{}, diags
	}
	blockedurlpatternsDataMap, diags := basetypes.NewObjectValue(PatternsDataAttrTypes(), map[string]attr.Value{
		"settings": basetypes.NewStringValue(contentFilteringData.BlockedUrlPatterns.Settings),
		"patterns": blockedurlpatternsList,
	})
	if diags.HasError() {
		return basetypes.ObjectValue{}, diags
	}

	contentFilteringMap, diags := basetypes.NewObjectValue(ContentFilteringDataAttrTypes(), map[string]attr.Value{
		"allowed_url_patterns":   allowedpatternsDataMap,
		"blocked_url_categories": blockedcategoriesDataMap,
		"blocked_url_patterns":   blockedurlpatternsDataMap,
	})
	if diags.HasError() {
		return basetypes.ObjectValue{}, diags
	}
	contentFilteringVal, diags := types.ObjectValueFrom(ctx, ContentFilteringDataAttrTypes(), contentFilteringMap)
	if diags.HasError() {
		return basetypes.ObjectValue{}, diags
	}
	return contentFilteringVal, nil
}

func getFirewallAndTrafficShaping(ctx context.Context, inlineResp map[string]interface{}) (basetypes.ObjectValue, diag.Diagnostics) {

	var firewallAndTrafficShapingOutput OutputNetworksGroupPolicyResourceModelFirewallAndTrafficShaping
	jsonData, _ := json.Marshal(inlineResp["firewallAndTrafficShaping"].(map[string]interface{}))
	json.Unmarshal(jsonData, &firewallAndTrafficShapingOutput)

	var trafficShapingRule []OutputNetworksGroupPolicyResourceModelTrafficShapingRule
	jsonData, _ = json.Marshal(inlineResp["firewallAndTrafficShaping"].(map[string]interface{})["trafficShapingRules"])
	json.Unmarshal(jsonData, &trafficShapingRule)

	var l3FirewallRules []basetypes.ObjectValue
	for _, rule := range firewallAndTrafficShapingOutput.L3FirewallRules {

		ruleDataMap, _ := basetypes.NewObjectValue(L3FirewallRulesAttrTypes(), map[string]attr.Value{
			"comment":   basetypes.NewStringValue(rule.Comment),
			"dest_cidr": basetypes.NewStringValue(rule.DestCidr),
			"dest_port": basetypes.NewStringValue(rule.DestPort),
			"policy":    basetypes.NewStringValue(rule.Policy),
			"protocol":  basetypes.NewStringValue(rule.Protocol),
		})

		objectVal, diags := types.ObjectValueFrom(ctx, L3FirewallRulesAttrTypes(), ruleDataMap)
		if diags.HasError() {
			return basetypes.ObjectValue{}, diags
		}
		l3FirewallRules = append(l3FirewallRules, objectVal)
	}

	l3FirewallRulesList, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: L3FirewallRulesAttrTypes()}, l3FirewallRules)

	var l7FirewallRules []basetypes.ObjectValue
	for _, rule := range firewallAndTrafficShapingOutput.L7FirewallRules {

		ruleDataMap, _ := basetypes.NewObjectValue(L7FirewallRulesAttrTypes(), map[string]attr.Value{
			"value":  basetypes.NewStringValue(rule.Value),
			"policy": basetypes.NewStringValue(rule.Policy),
			"type":   basetypes.NewStringValue(rule.Type),
		})
		objectVal, diags := types.ObjectValueFrom(ctx, L7FirewallRulesAttrTypes(), ruleDataMap)
		if diags.HasError() {
			return basetypes.ObjectValue{}, diags
		}
		l7FirewallRules = append(l7FirewallRules, objectVal)
	}

	l7FirewallRulesList, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: L7FirewallRulesAttrTypes()}, l7FirewallRules)

	var trafficShapingRules []basetypes.ObjectValue

	for _, rule := range trafficShapingRule {

		var definitions []basetypes.ObjectValue
		for _, definition := range rule.Definitions {
			definitionMap, diags := basetypes.NewObjectValue(DefinitionAttrTypes(), map[string]attr.Value{
				"type":  basetypes.NewStringValue(definition.Type),
				"value": basetypes.NewStringValue(definition.Value),
			})
			if diags.HasError() {
				return basetypes.ObjectValue{}, diags
			}
			definitionMapObjectVal, diags := types.ObjectValueFrom(ctx, DefinitionAttrTypes(), definitionMap)
			if diags.HasError() {
				return basetypes.ObjectValue{}, diags
			}
			definitions = append(definitions, definitionMapObjectVal)
		}
		definitionsList, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: DefinitionAttrTypes()}, definitions)

		bandwidthLimitsMap, diags := basetypes.NewObjectValue(BandwidthLimitsAttrTypes(), map[string]attr.Value{
			"limit_up":   basetypes.NewInt64Value(rule.PerClientBandwidthLimits.BandwidthLimits.LimitUp),
			"limit_down": basetypes.NewInt64Value(rule.PerClientBandwidthLimits.BandwidthLimits.LimitDown),
		})
		if diags.HasError() {
			return basetypes.ObjectValue{}, diags
		}
		bandwidthMap, _ := basetypes.NewObjectValue(PerClientBandwidthAttrTypes(), map[string]attr.Value{
			"bandwidth_limits": bandwidthLimitsMap,
			"settings":         basetypes.NewStringValue(rule.PerClientBandwidthLimits.Settings),
		})
		objectVal, diags := types.ObjectValueFrom(ctx, PerClientBandwidthAttrTypes(), bandwidthMap)
		if diags.HasError() {
			return basetypes.ObjectValue{}, diags
		}

		ruleDataMap, _ := basetypes.NewObjectValue(L7FirewallRulesAttrTypes(), map[string]attr.Value{
			"dscp_tag_value":              basetypes.NewInt64Value(rule.DscpTagValue),
			"pcp_tag_value":               basetypes.NewInt64Value(rule.PcpTagValue),
			"per_client_bandwidth_limits": objectVal,
			"definitions":                 definitionsList,
		})
		trafficShapingRulesObjectVal, diags := types.ObjectValueFrom(ctx, TrafficShapingRulesAttrTypes(), ruleDataMap)
		if diags.HasError() {
			return basetypes.ObjectValue{}, diags
		}

		trafficShapingRules = append(trafficShapingRules, trafficShapingRulesObjectVal)

	}

	trafficShapingRulesList, _ := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: TrafficShapingRulesAttrTypes()}, trafficShapingRules)

	firewallAndTrafficShapingMap, diags := basetypes.NewObjectValue(FirewallAndTrafficShapingAttrTypes(), map[string]attr.Value{
		"settings":              basetypes.NewStringValue(firewallAndTrafficShapingOutput.Settings),
		"l3_firewall_rules":     l3FirewallRulesList,
		"l7_firewall_rules":     l7FirewallRulesList,
		"traffic_shaping_rules": trafficShapingRulesList,
	})
	if diags.HasError() {
		return basetypes.ObjectValue{}, diags
	}

	firewallAndTrafficShapingObjectVal, diags := types.ObjectValueFrom(ctx, FirewallAndTrafficShapingAttrTypes(), firewallAndTrafficShapingMap)
	if diags.HasError() {
		return basetypes.ObjectValue{}, diags
	}

	return firewallAndTrafficShapingObjectVal, nil

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
						Default:             stringdefault.StaticString("network default"),
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
						Default:             stringdefault.StaticString("network default"),
					},
					"rules": schema.ListNestedAttribute{
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
								Default:             stringdefault.StaticString("network default"),
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
								Default:             stringdefault.StaticString("network default"),
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
								Default:             stringdefault.StaticString("network default"),
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
			"vlan_tagging": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"settings": schema.StringAttribute{
						MarkdownDescription: "How VLAN tagging is applied. Can be 'network default', 'ignore' or 'custom'.",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.StringType,
						Default:             stringdefault.StaticString("network default"),
					},
					"vlan_id": schema.StringAttribute{
						MarkdownDescription: "The ID of the vlan you want to tag. This only applies if 'settings' is set to 'custom'.",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.StringType,
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
						Default:             stringdefault.StaticString("network default"),
					},
					"l3_firewall_rules": schema.ListNestedAttribute{
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
					"l7_firewall_rules": schema.ListNestedAttribute{
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
					"traffic_shaping_rules": schema.ListNestedAttribute{
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
											Default:             stringdefault.StaticString("network default"),
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
								"definitions": schema.ListNestedAttribute{
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

	createNetworkGroupPolicy := *openApiClient.NewCreateNetworkGroupPolicyRequest(data.Name.ValueString())
	if !data.SplashAuthSettings.IsUnknown() {
		createNetworkGroupPolicy.SetSplashAuthSettings(data.SplashAuthSettings.ValueString())

	}

	if !data.VlanTagging.IsUnknown() && !data.VlanTagging.IsNull() {
		var vlanTagging NetworksGroupPolicyResourceModelVlanTagging
		vlanTaggingErr := data.Bandwidth.As(ctx, &vlanTagging, basetypes.ObjectAsOptions{})
		if vlanTaggingErr != nil {
			resp.Diagnostics.AddError(
				"Failed to unmarshal vlanTagging",
				fmt.Sprintf("%v", vlanTaggingErr),
			)
			return
		}
		if !vlanTagging.Settings.IsUnknown() {
			if !vlanTagging.VlanId.IsUnknown() {
				var v openApiClient.CreateNetworkGroupPolicyRequestVlanTagging
				v.SetSettings(vlanTagging.Settings.ValueString())
				v.SetVlanId(vlanTagging.VlanId.ValueString())
				createNetworkGroupPolicy.SetVlanTagging(v)
			}
		}
	}

	if !data.Bandwidth.IsUnknown() && !data.Bandwidth.IsNull() {
		var bandwidth NetworksGroupPolicyResourceModelBandwidth
		bandwidthDataErr := data.Bandwidth.As(ctx, &bandwidth, basetypes.ObjectAsOptions{})
		if bandwidthDataErr != nil {
			resp.Diagnostics.AddError(
				"Failed to unmarshal Bandwidth",
				fmt.Sprintf("%v", bandwidthDataErr),
			)
			return
		}
		if !bandwidth.Settings.IsUnknown() {
			var bandwidthAPI openApiClient.CreateNetworkGroupPolicyRequestBandwidth
			bandwidthAPI.SetSettings(bandwidth.Settings.ValueString())
			var bandwidthLimitsAPI openApiClient.CreateNetworkGroupPolicyRequestBandwidthBandwidthLimits
			if !bandwidth.BandwidthLimits.LimitUp.IsUnknown() {
				bandwidthLimitsAPI.SetLimitUp(int32(bandwidth.BandwidthLimits.LimitUp.ValueInt64()))
			}
			if !bandwidth.BandwidthLimits.LimitDown.IsUnknown() {
				bandwidthLimitsAPI.SetLimitDown(int32(bandwidth.BandwidthLimits.LimitDown.ValueInt64()))
			}
			bandwidthAPI.SetBandwidthLimits(bandwidthLimitsAPI)
			createNetworkGroupPolicy.SetBandwidth(bandwidthAPI)
		}
	}

	if !data.BonjourForwarding.IsUnknown() && !data.BonjourForwarding.IsNull() {

		var bonjourForwarding openApiClient.CreateNetworkGroupPolicyRequestBonjourForwarding
		var bonjourForwardingData NetworksGroupPolicyResourceModelBonjourForwarding
		bonjourForwardingDataErr := data.Bandwidth.As(ctx, &bonjourForwardingData, basetypes.ObjectAsOptions{})
		if bonjourForwardingDataErr != nil {
			resp.Diagnostics.AddError(
				"Failed to unmarshal Bonjour Forwarding",
				fmt.Sprintf("%v", bonjourForwardingDataErr),
			)
			return
		}
		var bonjourForwardingRules []openApiClient.CreateNetworkGroupPolicyRequestBonjourForwardingRulesInner
		var bonjourForwardingRulesData []NetworksGroupPolicyResourceModelRule
		err := bonjourForwardingData.BonjourForwardingRules.ElementsAs(ctx, &bonjourForwardingRulesData, false)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to unmarshal FirewallAndTrafficShaping",
				fmt.Sprintf("%v", err),
			)
			return

		}
		for _, attribute := range bonjourForwardingRulesData {
			var bonjourForwardingRule openApiClient.CreateNetworkGroupPolicyRequestBonjourForwardingRulesInner
			if !attribute.Description.IsUnknown() {
				bonjourForwardingRule.SetDescription(attribute.Description.ValueString())
			}
			bonjourForwardingRule.SetVlanId(attribute.VlanId.ValueString())
			var services []string
			serviceserr := attribute.Services.ElementsAs(ctx, &services, false)
			if serviceserr != nil {
				resp.Diagnostics.AddError(
					"Failed to unmarshal bonjourForwardingRule services",
					fmt.Sprintf("%v", serviceserr),
				)
				return

			}
			bonjourForwardingRule.SetServices(services)
			bonjourForwardingRules = append(bonjourForwardingRules, bonjourForwardingRule)
		}
		bonjourForwarding.SetRules(bonjourForwardingRules)
		if !bonjourForwardingData.BonjourForwardingSettings.IsUnknown() {
			bonjourForwarding.SetSettings(bonjourForwardingData.BonjourForwardingSettings.ValueString())
		}
		createNetworkGroupPolicy.SetBonjourForwarding(bonjourForwarding)

	}

	if !data.FirewallAndTrafficShaping.IsUnknown() && !data.FirewallAndTrafficShaping.IsNull() {

		var firewallAndTrafficShaping openApiClient.CreateNetworkGroupPolicyRequestFirewallAndTrafficShaping
		var firewallAndTrafficShapingData NetworksGroupPolicyResourceModelFirewallAndTrafficShaping
		firewallAndTrafficShapingErr := data.FirewallAndTrafficShaping.As(ctx, &firewallAndTrafficShapingData, basetypes.ObjectAsOptions{})
		if firewallAndTrafficShapingErr != nil {
			resp.Diagnostics.AddError(
				"Failed to unmarshal FirewallAndTrafficShaping",
				fmt.Sprintf("%v", firewallAndTrafficShapingErr),
			)
			return
		}
		if !firewallAndTrafficShapingData.Settings.IsUnknown() {
			firewallAndTrafficShaping.SetSettings(firewallAndTrafficShapingData.Settings.ValueString())
		}
		var l3FirewallRules []NetworksGroupPolicyResourceModelL3FirewallRule
		err := firewallAndTrafficShapingData.L3FirewallRules.ElementsAs(ctx, &l3FirewallRules, false)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to unmarshal FirewallAndTrafficShaping",
				fmt.Sprintf("%v", err),
			)
			return

		}
		if len(l3FirewallRules) > 0 {
			var l3s []openApiClient.CreateNetworkGroupPolicyRequestFirewallAndTrafficShapingL3FirewallRulesInner
			for _, attribute := range l3FirewallRules {
				var l3 openApiClient.CreateNetworkGroupPolicyRequestFirewallAndTrafficShapingL3FirewallRulesInner
				if !attribute.Comment.IsUnknown() {
					l3.SetComment(attribute.Comment.ValueString())
				}
				if !attribute.DestCidr.IsUnknown() {
					l3.SetDestCidr(attribute.DestCidr.ValueString())
				}
				if !attribute.DestPort.IsUnknown() {
					l3.SetDestPort(attribute.DestPort.ValueString())
				}
				if !attribute.Policy.IsUnknown() {
					l3.SetPolicy(attribute.Policy.ValueString())
				}
				if !attribute.Protocol.IsUnknown() {
					l3.SetProtocol(attribute.Protocol.ValueString())
				}
				l3s = append(l3s, l3)
			}
			firewallAndTrafficShaping.SetL3FirewallRules(l3s)
		}
		var l7FirewallRules []NetworksGroupPolicyResourceModelL7FirewallRule
		err = firewallAndTrafficShapingData.L7FirewallRules.ElementsAs(ctx, &l7FirewallRules, false)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to unmarshal L7FirewallRules",
				fmt.Sprintf("%v", err),
			)
			return

		}
		if len(l7FirewallRules) > 0 {
			var l7s []openApiClient.CreateNetworkGroupPolicyRequestFirewallAndTrafficShapingL7FirewallRulesInner
			for _, attribute := range l7FirewallRules {
				var l7 openApiClient.CreateNetworkGroupPolicyRequestFirewallAndTrafficShapingL7FirewallRulesInner
				if !attribute.Value.IsUnknown() {
					l7.SetValue(attribute.Value.ValueString())
				}
				if !attribute.Type.IsUnknown() {
					l7.SetType(attribute.Type.ValueString())
				}

				if !attribute.Policy.IsUnknown() {
					l7.SetPolicy(attribute.Policy.ValueString())
				}

				l7s = append(l7s, l7)
			}
			firewallAndTrafficShaping.SetL7FirewallRules(l7s)

		}
		var trafficShapingRule []NetworksGroupPolicyResourceModelTrafficShapingRule
		err = firewallAndTrafficShapingData.TrafficShapingRules.ElementsAs(ctx, &trafficShapingRule, false)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to unmarshal TrafficShapingRules",
				fmt.Sprintf("%v", err),
			)
			return

		}
		if len(trafficShapingRule) > 0 {
			var tfs []openApiClient.CreateNetworkGroupPolicyRequestFirewallAndTrafficShapingTrafficShapingRulesInner
			for _, attribute := range trafficShapingRule {
				var tf openApiClient.CreateNetworkGroupPolicyRequestFirewallAndTrafficShapingTrafficShapingRulesInner
				if !attribute.DscpTagValue.IsUnknown() {
					tf.SetDscpTagValue(int32(attribute.DscpTagValue.ValueInt64()))
				}
				if !attribute.PcpTagValue.IsUnknown() {
					tf.SetPcpTagValue(int32(attribute.PcpTagValue.ValueInt64()))
				}
				var perClientBandWidthLimits openApiClient.UpdateNetworkApplianceTrafficShapingRulesRequestRulesInnerPerClientBandwidthLimits
				var perClientBandWidthData NetworksGroupPolicyResourceModelPerClientBandwidth
				perClientBandWidthDataErr := attribute.PerClientBandwidthLimits.As(ctx, &perClientBandWidthData, basetypes.ObjectAsOptions{})
				if perClientBandWidthDataErr != nil {
					resp.Diagnostics.AddError(
						"Failed to unmarshal perClientBandWidth",
						fmt.Sprintf("%v", perClientBandWidthDataErr),
					)
					return

				}

				if !perClientBandWidthData.Settings.IsUnknown() {
					var bandwidthLimits openApiClient.UpdateNetworkApplianceTrafficShapingRulesRequestRulesInnerPerClientBandwidthLimitsBandwidthLimits

					if perClientBandWidthData.Settings.ValueString() != "network default" {
						var perClientBandWidthLimitsData NetworksGroupPolicyResourceModelPerClientBandwidthLimits
						perClientBandWidthLimitsDataErr := attribute.PerClientBandwidthLimits.As(ctx, &perClientBandWidthLimitsData, basetypes.ObjectAsOptions{})
						if perClientBandWidthLimitsDataErr != nil {
							resp.Diagnostics.AddError(
								"Failed to unmarshal perClientBandWidthLimits",
								fmt.Sprintf("%v", perClientBandWidthLimitsDataErr),
							)
							return

						}
						bandwidthLimits.SetLimitDown(int32(perClientBandWidthLimitsData.LimitDown.ValueInt64()))
						bandwidthLimits.SetLimitUp(int32(perClientBandWidthLimitsData.LimitUp.ValueInt64()))

						perClientBandWidthLimits.SetBandwidthLimits(bandwidthLimits)
					}
					perClientBandWidthLimits.SetSettings(perClientBandWidthData.Settings.ValueString())
					tf.SetPerClientBandwidthLimits(perClientBandWidthLimits)
				}
				var defs []openApiClient.UpdateNetworkApplianceTrafficShapingRulesRequestRulesInnerDefinitionsInner
				var definitions []NetworksGroupPolicyResourceModelDefinition
				err = attribute.Definitions.ElementsAs(ctx, &definitions, false)
				if err != nil {
					resp.Diagnostics.AddError(
						"Failed to unmarshal Definitions",
						fmt.Sprintf("%v", err),
					)
					return

				}
				if len(definitions) > 0 {
					for _, attribute := range definitions {
						var def openApiClient.UpdateNetworkApplianceTrafficShapingRulesRequestRulesInnerDefinitionsInner
						def.SetType(attribute.Type.ValueString())
						def.SetValue(attribute.Value.ValueString())
						defs = append(defs, def)
					}
					tf.SetDefinitions(defs)
				}
				tfs = append(tfs, tf)
			}
			firewallAndTrafficShaping.SetTrafficShapingRules(tfs)
		}

		createNetworkGroupPolicy.SetFirewallAndTrafficShaping(firewallAndTrafficShaping)
	}

	if !data.Scheduling.IsUnknown() && !data.Scheduling.IsNull() {
		var schedulingData NetworksGroupPolicyResourceModelScheduling
		schedulingDataErr := data.Scheduling.As(ctx, &schedulingData, basetypes.ObjectAsOptions{})
		if schedulingDataErr != nil {
			resp.Diagnostics.AddError(
				"Failed to unmarshal Scheduling",
				fmt.Sprintf("%v", schedulingDataErr),
			)
			return
		}

		if !schedulingData.Enabled.IsUnknown() {
			var schedule openApiClient.CreateNetworkGroupPolicyRequestScheduling
			schedule.SetEnabled(schedulingData.Enabled.ValueBool())
			var schedulingFridayDayData NetworksGroupPolicyResourceModelSchedule
			schedulingFridayDayDataErr := schedulingData.Friday.As(ctx, &schedulingFridayDayData, basetypes.ObjectAsOptions{})
			if schedulingFridayDayDataErr != nil {

				resp.Diagnostics.AddError(
					"Failed to unmarshal Friday Scheduling",
					fmt.Sprintf("%v", schedulingFridayDayDataErr),
				)
				return
			}
			if !schedulingFridayDayData.Active.IsUnknown() {
				var friday openApiClient.CreateNetworkGroupPolicyRequestSchedulingFriday
				friday.SetActive(schedulingFridayDayData.Active.ValueBool())
				friday.SetFrom(schedulingFridayDayData.From.ValueString())
				friday.SetTo(schedulingFridayDayData.To.ValueString())
				schedule.SetFriday(friday)
			}
			var schedulingMondayDayData NetworksGroupPolicyResourceModelSchedule
			schedulingMondayDayDataErr := schedulingData.Monday.As(ctx, &schedulingMondayDayData, basetypes.ObjectAsOptions{})
			if schedulingMondayDayDataErr != nil {

				resp.Diagnostics.AddError(
					"Failed to unmarshal Monday Scheduling",
					fmt.Sprintf("%v", schedulingMondayDayDataErr),
				)
				return
			}
			if !schedulingMondayDayData.Active.IsUnknown() {
				var monday openApiClient.CreateNetworkGroupPolicyRequestSchedulingMonday
				monday.SetActive(schedulingMondayDayData.Active.ValueBool())
				monday.SetFrom(schedulingMondayDayData.From.ValueString())
				monday.SetTo(schedulingMondayDayData.To.ValueString())
				schedule.SetMonday(monday)
			}
			var schedulingTuesdayDayData NetworksGroupPolicyResourceModelSchedule
			schedulingTuesdayDayDataErr := schedulingData.Tuesday.As(ctx, &schedulingTuesdayDayData, basetypes.ObjectAsOptions{})
			if schedulingTuesdayDayDataErr != nil {
				resp.Diagnostics.AddError(
					"Failed to unmarshal Tuesday Scheduling",
					fmt.Sprintf("%v", schedulingTuesdayDayDataErr),
				)
				return
			}
			if !schedulingTuesdayDayData.Active.IsUnknown() {
				var tuesday openApiClient.CreateNetworkGroupPolicyRequestSchedulingTuesday
				tuesday.SetActive(schedulingTuesdayDayData.Active.ValueBool())
				tuesday.SetFrom(schedulingTuesdayDayData.From.ValueString())
				tuesday.SetTo(schedulingTuesdayDayData.To.ValueString())
				schedule.SetTuesday(tuesday)
			}
			var schedulingWednesdayDayData NetworksGroupPolicyResourceModelSchedule
			schedulingWednesdayDayDataErr := schedulingData.Wednesday.As(ctx, &schedulingWednesdayDayData, basetypes.ObjectAsOptions{})
			if schedulingWednesdayDayDataErr != nil {
				resp.Diagnostics.AddError(
					"Failed to unmarshal Wednesday Scheduling",
					fmt.Sprintf("%v", schedulingWednesdayDayDataErr),
				)
				return
			}
			if !schedulingWednesdayDayData.Active.IsUnknown() {
				var wednesday openApiClient.CreateNetworkGroupPolicyRequestSchedulingWednesday
				wednesday.SetActive(schedulingWednesdayDayData.Active.ValueBool())
				wednesday.SetFrom(schedulingWednesdayDayData.From.ValueString())
				wednesday.SetTo(schedulingWednesdayDayData.To.ValueString())
				schedule.SetWednesday(wednesday)
			}
			var schedulingThursdayDayData NetworksGroupPolicyResourceModelSchedule
			schedulingThursdayDayDataErr := schedulingData.Thursday.As(ctx, &schedulingThursdayDayData, basetypes.ObjectAsOptions{})
			if schedulingThursdayDayDataErr != nil {
				resp.Diagnostics.AddError(
					"Failed to unmarshal Thursday Scheduling",
					fmt.Sprintf("%v", schedulingThursdayDayDataErr),
				)
				return
			}
			if !schedulingThursdayDayData.Active.IsUnknown() {
				var thursday openApiClient.CreateNetworkGroupPolicyRequestSchedulingThursday
				thursday.SetActive(schedulingThursdayDayData.Active.ValueBool())
				thursday.SetFrom(schedulingThursdayDayData.From.ValueString())
				thursday.SetTo(schedulingThursdayDayData.To.ValueString())
				schedule.SetThursday(thursday)
			}
			var schedulingSaturdayDayData NetworksGroupPolicyResourceModelSchedule
			schedulingSaturdayDayDataErr := schedulingData.Saturday.As(ctx, &schedulingSaturdayDayData, basetypes.ObjectAsOptions{})
			if schedulingSaturdayDayDataErr != nil {
				resp.Diagnostics.AddError(
					"Failed to unmarshal Saturday Scheduling",
					fmt.Sprintf("%v", schedulingSaturdayDayDataErr),
				)
				return
			}
			if !schedulingSaturdayDayData.Active.IsUnknown() {
				var saturday openApiClient.CreateNetworkGroupPolicyRequestSchedulingSaturday
				saturday.SetActive(schedulingSaturdayDayData.Active.ValueBool())
				saturday.SetFrom(schedulingSaturdayDayData.From.ValueString())
				saturday.SetTo(schedulingSaturdayDayData.To.ValueString())
				schedule.SetSaturday(saturday)
			}
			var schedulingSundayDayData NetworksGroupPolicyResourceModelSchedule
			schedulingSundayDayDataErr := schedulingData.Sunday.As(ctx, &schedulingSundayDayData, basetypes.ObjectAsOptions{})
			if schedulingSundayDayDataErr != nil {
				resp.Diagnostics.AddError(
					"Failed to unmarshal Sunday Scheduling",
					fmt.Sprintf("%v", schedulingSundayDayDataErr),
				)
				return
			}
			if !schedulingSundayDayData.Active.IsUnknown() {
				var sunday openApiClient.CreateNetworkGroupPolicyRequestSchedulingSunday
				sunday.SetActive(schedulingSundayDayData.Active.ValueBool())
				sunday.SetFrom(schedulingSundayDayData.From.ValueString())
				sunday.SetTo(schedulingSundayDayData.To.ValueString())
				schedule.SetSunday(sunday)
			}
			createNetworkGroupPolicy.SetScheduling(schedule)
		}

	}

	if !data.ContentFiltering.IsUnknown() && !data.ContentFiltering.IsNull() {
		var contentFiltering openApiClient.CreateNetworkGroupPolicyRequestContentFiltering
		var contentFilteringData NetworksGroupPolicyResourceModelContentFiltering
		contentDataErr := data.ContentFiltering.As(ctx, &contentFilteringData, basetypes.ObjectAsOptions{})
		if contentDataErr != nil {
			resp.Diagnostics.AddError(
				"Failed to unmarshal ContentFiltering",
				fmt.Sprintf("%v", contentDataErr),
			)
			return
		}
		contentFilteringStatus := false

		var contentFilteringPatternsData NetworksGroupPolicyResourceModelAllowedUrlPatterns
		contentDataPatternsErr := contentFilteringData.AllowedUrlPatterns.As(ctx, &contentFilteringPatternsData, basetypes.ObjectAsOptions{})
		if contentDataPatternsErr != nil {
			resp.Diagnostics.AddError(
				"Failed to unmarshal ContentFilteringPatterns",
				fmt.Sprintf("%v", contentDataPatternsErr),
			)
			return
		}

		if !contentFilteringPatternsData.Settings.IsUnknown() {
			var allowedUrlPatternData openApiClient.CreateNetworkGroupPolicyRequestContentFilteringAllowedUrlPatterns
			allowedUrlPatternData.SetSettings(contentFilteringPatternsData.Settings.ValueString())
			var patternsData []string
			patternsDataerr := contentFilteringPatternsData.Patterns.ElementsAs(ctx, &patternsData, false)
			if patternsDataerr != nil {
				resp.Diagnostics.AddError(
					"Failed to unmarshal AllowedURLPatterns",
					fmt.Sprintf("%v", patternsDataerr),
				)
				return

			}
			allowedUrlPatternData.SetPatterns(patternsData)
			contentFiltering.SetAllowedUrlPatterns(allowedUrlPatternData)
			contentFilteringStatus = true
		}

		var contentFilteringBlockedUrlCategoriesData NetworksGroupPolicyResourceModelBlockedUrlCategories
		contentFilteringBlockedUrlCategoriesDataErr := contentFilteringData.BlockedUrlCategories.As(ctx, &contentFilteringBlockedUrlCategoriesData, basetypes.ObjectAsOptions{})
		if contentFilteringBlockedUrlCategoriesDataErr != nil {
			resp.Diagnostics.AddError(
				"Failed to unmarshal ContentBlockedUrlCategories",
				fmt.Sprintf("%v", contentFilteringBlockedUrlCategoriesDataErr),
			)
			return
		}

		if !contentFilteringBlockedUrlCategoriesData.Settings.IsUnknown() {
			var blockedUrlCategoriesData openApiClient.CreateNetworkGroupPolicyRequestContentFilteringBlockedUrlCategories
			blockedUrlCategoriesData.SetSettings(contentFilteringBlockedUrlCategoriesData.Settings.ValueString())
			var blockedUrlCategories []string
			blockedUrlCategoriesDataerr := contentFilteringBlockedUrlCategoriesData.Categories.ElementsAs(ctx, &blockedUrlCategories, false)
			if blockedUrlCategoriesDataerr != nil {
				resp.Diagnostics.AddError(
					"Failed to unmarshal blockedUrlCategories",
					fmt.Sprintf("%v", blockedUrlCategoriesDataerr),
				)
				return

			}
			blockedUrlCategoriesData.SetCategories(blockedUrlCategories)
			contentFiltering.SetBlockedUrlCategories(blockedUrlCategoriesData)
			contentFilteringStatus = true
		}

		var contentFilteringBlockedUrlPatternsData NetworksGroupPolicyResourceModelBlockedUrlPatterns
		contentFilteringBlockedUrlPatternsDataErr := contentFilteringData.BlockedUrlPatterns.As(ctx, &contentFilteringBlockedUrlPatternsData, basetypes.ObjectAsOptions{})
		if contentFilteringBlockedUrlPatternsDataErr != nil {
			resp.Diagnostics.AddError(
				"Failed to unmarshal ContentBlockedUrlPatterns",
				fmt.Sprintf("%v", contentFilteringBlockedUrlPatternsDataErr),
			)
			return
		}

		if !contentFilteringBlockedUrlPatternsData.Settings.IsUnknown() {
			var blockedUrlPatternsData openApiClient.CreateNetworkGroupPolicyRequestContentFilteringBlockedUrlPatterns
			blockedUrlPatternsData.SetSettings(contentFilteringBlockedUrlPatternsData.Settings.ValueString())
			var blockedUrlPatterns []string
			blockedUrlPatternsDataerr := contentFilteringBlockedUrlPatternsData.Patterns.ElementsAs(ctx, &blockedUrlPatterns, false)
			if blockedUrlPatternsDataerr != nil {
				resp.Diagnostics.AddError(
					"Failed to unmarshal blockedUrlPatterns",
					fmt.Sprintf("%v", blockedUrlPatternsDataerr),
				)
				return

			}
			blockedUrlPatternsData.SetPatterns(blockedUrlPatterns)
			contentFiltering.SetBlockedUrlPatterns(blockedUrlPatternsData)
			contentFilteringStatus = true
		}

		if contentFilteringStatus {
			createNetworkGroupPolicy.SetContentFiltering(contentFiltering)
		}
	}

	// Wrap the API call in the retryAPICall function
	inlineRespMap, httpResp, err := retryAPICall(ctx, func() (interface{}, *http.Response, error) {
		resp, httpResp, err := r.client.NetworksApi.CreateNetworkGroupPolicy(ctx, data.NetworkId.ValueString()).CreateNetworkGroupPolicyRequest(createNetworkGroupPolicy).Execute()
		if err != nil {
			return nil, httpResp, err
		}
		return resp, httpResp, nil
	})

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

	inlineResp, ok := inlineRespMap.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Type Assertion Failure",
			"Failed to assert API response type to map[string]interface{}",
		)
		return
	}

	// Save data into Terraform state
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	bandwidthVal, diags := getBandwidth(ctx, inlineResp)
	if diags.HasError() {
		return
	}
	firewallAndTrafficShapingObjectVal, diags := getFirewallAndTrafficShaping(ctx, inlineResp)
	if diags.HasError() {
		return
	}
	schedulingVal, diags := getScheduling(ctx, inlineResp)
	if diags.HasError() {
		return
	}
	contentFilteringVal, diags := getContentFiltering(ctx, inlineResp)
	if diags.HasError() {
		return
	}
	bonjourForwardingVal, diags := getBonjourForwarding(ctx, inlineResp)
	if diags.HasError() {
		return
	}

	vlantaggingVal, diags := getVlantagging(ctx, inlineResp)
	if diags.HasError() {
		return
	}

	data.Bandwidth = bandwidthVal
	data.FirewallAndTrafficShaping = firewallAndTrafficShapingObjectVal
	data.Scheduling = schedulingVal
	data.ContentFiltering = contentFilteringVal
	data.BonjourForwarding = bonjourForwardingVal
	data.VlanTagging = vlantaggingVal
	data.Id = jsontypes.StringValue("example-id")

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

	inlineRespMap, httpResp, err := retryAPICall(ctx, func() (interface{}, *http.Response, error) {
		resp, httpResp, err := r.client.NetworksApi.GetNetworkGroupPolicy(ctx, data.NetworkId.ValueString(), data.GroupPolicyId.ValueString()).Execute()
		if err != nil {
			return nil, httpResp, err
		}
		return resp, httpResp, nil
	})

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

	inlineResp, ok := inlineRespMap.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Type Assertion Failure",
			"Failed to assert API response type to map[string]interface{}",
		)
		return
	}

	// Save data into Terraform state
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	bandwidthVal, diags := getBandwidth(ctx, inlineResp)
	if diags.HasError() {
		return
	}
	firewallAndTrafficShapingObjectVal, diags := getFirewallAndTrafficShaping(ctx, inlineResp)
	if diags.HasError() {
		return
	}
	schedulingVal, diags := getScheduling(ctx, inlineResp)
	if diags.HasError() {
		return
	}
	contentFilteringVal, diags := getContentFiltering(ctx, inlineResp)
	if diags.HasError() {
		return
	}
	bonjourForwardingVal, diags := getBonjourForwarding(ctx, inlineResp)
	if diags.HasError() {
		return
	}

	vlantaggingVal, diags := getVlantagging(ctx, inlineResp)
	if diags.HasError() {
		return
	}

	data.Bandwidth = bandwidthVal
	data.FirewallAndTrafficShaping = firewallAndTrafficShapingObjectVal
	data.Scheduling = schedulingVal
	data.ContentFiltering = contentFilteringVal
	data.BonjourForwarding = bonjourForwardingVal
	data.VlanTagging = vlantaggingVal
	data.Id = jsontypes.StringValue("example-id")

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

	updateNetworkGroupPolicy := *openApiClient.NewUpdateNetworkGroupPolicyRequest()

	if !data.SplashAuthSettings.IsUnknown() {
		updateNetworkGroupPolicy.SetSplashAuthSettings(data.SplashAuthSettings.ValueString())

	}

	if !data.VlanTagging.IsUnknown() && !data.VlanTagging.IsNull() {
		var vlanTagging NetworksGroupPolicyResourceModelVlanTagging
		vlanTaggingErr := data.Bandwidth.As(ctx, &vlanTagging, basetypes.ObjectAsOptions{})
		if vlanTaggingErr != nil {
			resp.Diagnostics.AddError(
				"Failed to unmarshal vlanTagging",
				fmt.Sprintf("%v", vlanTaggingErr),
			)
			return
		}
		if !vlanTagging.Settings.IsUnknown() {
			if !vlanTagging.VlanId.IsUnknown() {
				var v openApiClient.CreateNetworkGroupPolicyRequestVlanTagging
				v.SetSettings(vlanTagging.Settings.ValueString())
				v.SetVlanId(vlanTagging.VlanId.ValueString())
				updateNetworkGroupPolicy.SetVlanTagging(v)
			}
		}
	}

	if !data.Bandwidth.IsUnknown() && !data.Bandwidth.IsNull() {
		var bandwidth NetworksGroupPolicyResourceModelBandwidth
		bandwidthDataErr := data.Bandwidth.As(ctx, &bandwidth, basetypes.ObjectAsOptions{})
		if bandwidthDataErr != nil {
			resp.Diagnostics.AddError(
				"Failed to unmarshal Bandwidth",
				fmt.Sprintf("%v", bandwidthDataErr),
			)
			return
		}
		if !bandwidth.Settings.IsUnknown() {
			var bandwidthAPI openApiClient.CreateNetworkGroupPolicyRequestBandwidth
			bandwidthAPI.SetSettings(bandwidth.Settings.ValueString())
			var bandwidthLimitsAPI openApiClient.CreateNetworkGroupPolicyRequestBandwidthBandwidthLimits
			if !bandwidth.BandwidthLimits.LimitUp.IsUnknown() {
				bandwidthLimitsAPI.SetLimitUp(int32(bandwidth.BandwidthLimits.LimitUp.ValueInt64()))
			}
			if !bandwidth.BandwidthLimits.LimitDown.IsUnknown() {
				bandwidthLimitsAPI.SetLimitDown(int32(bandwidth.BandwidthLimits.LimitDown.ValueInt64()))
			}
			bandwidthAPI.SetBandwidthLimits(bandwidthLimitsAPI)
			updateNetworkGroupPolicy.SetBandwidth(bandwidthAPI)
		}
	}

	if !data.BonjourForwarding.IsUnknown() && !data.BonjourForwarding.IsNull() {

		var bonjourForwarding openApiClient.CreateNetworkGroupPolicyRequestBonjourForwarding
		var bonjourForwardingData NetworksGroupPolicyResourceModelBonjourForwarding
		bonjourForwardingDataErr := data.Bandwidth.As(ctx, &bonjourForwardingData, basetypes.ObjectAsOptions{})
		if bonjourForwardingDataErr != nil {
			resp.Diagnostics.AddError(
				"Failed to unmarshal Bonjour Forwarding",
				fmt.Sprintf("%v", bonjourForwardingDataErr),
			)
			return
		}
		var bonjourForwardingRules []openApiClient.CreateNetworkGroupPolicyRequestBonjourForwardingRulesInner
		var bonjourForwardingRulesData []NetworksGroupPolicyResourceModelRule
		err := bonjourForwardingData.BonjourForwardingRules.ElementsAs(ctx, &bonjourForwardingRulesData, false)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to unmarshal FirewallAndTrafficShaping",
				fmt.Sprintf("%v", err),
			)
			return

		}
		for _, attribute := range bonjourForwardingRulesData {
			var bonjourForwardingRule openApiClient.CreateNetworkGroupPolicyRequestBonjourForwardingRulesInner
			if !attribute.Description.IsUnknown() {
				bonjourForwardingRule.SetDescription(attribute.Description.ValueString())
			}
			bonjourForwardingRule.SetVlanId(attribute.VlanId.ValueString())
			var services []string
			serviceserr := attribute.Services.ElementsAs(ctx, &services, false)
			if serviceserr != nil {
				resp.Diagnostics.AddError(
					"Failed to unmarshal bonjourForwardingRule services",
					fmt.Sprintf("%v", serviceserr),
				)
				return

			}
			bonjourForwardingRule.SetServices(services)
			bonjourForwardingRules = append(bonjourForwardingRules, bonjourForwardingRule)
		}
		bonjourForwarding.SetRules(bonjourForwardingRules)
		if !bonjourForwardingData.BonjourForwardingSettings.IsUnknown() {
			bonjourForwarding.SetSettings(bonjourForwardingData.BonjourForwardingSettings.ValueString())
		}
		updateNetworkGroupPolicy.SetBonjourForwarding(bonjourForwarding)

	}

	if !data.FirewallAndTrafficShaping.IsUnknown() && !data.FirewallAndTrafficShaping.IsNull() {

		var firewallAndTrafficShaping openApiClient.CreateNetworkGroupPolicyRequestFirewallAndTrafficShaping
		var firewallAndTrafficShapingData NetworksGroupPolicyResourceModelFirewallAndTrafficShaping
		firewallAndTrafficShapingErr := data.FirewallAndTrafficShaping.As(ctx, &firewallAndTrafficShapingData, basetypes.ObjectAsOptions{})
		if firewallAndTrafficShapingErr != nil {
			resp.Diagnostics.AddError(
				"Failed to unmarshal FirewallAndTrafficShaping",
				fmt.Sprintf("%v", firewallAndTrafficShapingErr),
			)
			return
		}
		if !firewallAndTrafficShapingData.Settings.IsUnknown() {
			firewallAndTrafficShaping.SetSettings(firewallAndTrafficShapingData.Settings.ValueString())
		}
		var l3FirewallRules []NetworksGroupPolicyResourceModelL3FirewallRule
		err := firewallAndTrafficShapingData.L3FirewallRules.ElementsAs(ctx, &l3FirewallRules, false)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to unmarshal FirewallAndTrafficShaping",
				fmt.Sprintf("%v", err),
			)
			return

		}
		if len(l3FirewallRules) > 0 {
			var l3s []openApiClient.CreateNetworkGroupPolicyRequestFirewallAndTrafficShapingL3FirewallRulesInner
			for _, attribute := range l3FirewallRules {
				var l3 openApiClient.CreateNetworkGroupPolicyRequestFirewallAndTrafficShapingL3FirewallRulesInner
				if !attribute.Comment.IsUnknown() {
					l3.SetComment(attribute.Comment.ValueString())
				}
				if !attribute.DestCidr.IsUnknown() {
					l3.SetDestCidr(attribute.DestCidr.ValueString())
				}
				if !attribute.DestPort.IsUnknown() {
					l3.SetDestPort(attribute.DestPort.ValueString())
				}
				if !attribute.Policy.IsUnknown() {
					l3.SetPolicy(attribute.Policy.ValueString())
				}
				if !attribute.Protocol.IsUnknown() {
					l3.SetProtocol(attribute.Protocol.ValueString())
				}
				l3s = append(l3s, l3)
			}
			firewallAndTrafficShaping.SetL3FirewallRules(l3s)
		}
		var l7FirewallRules []NetworksGroupPolicyResourceModelL7FirewallRule
		err = firewallAndTrafficShapingData.L7FirewallRules.ElementsAs(ctx, &l7FirewallRules, false)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to unmarshal L7FirewallRules",
				fmt.Sprintf("%v", err),
			)
			return

		}
		if len(l7FirewallRules) > 0 {
			var l7s []openApiClient.CreateNetworkGroupPolicyRequestFirewallAndTrafficShapingL7FirewallRulesInner
			for _, attribute := range l7FirewallRules {
				var l7 openApiClient.CreateNetworkGroupPolicyRequestFirewallAndTrafficShapingL7FirewallRulesInner
				if !attribute.Value.IsUnknown() {
					l7.SetValue(attribute.Value.ValueString())
				}
				if !attribute.Type.IsUnknown() {
					l7.SetType(attribute.Type.ValueString())
				}

				if !attribute.Policy.IsUnknown() {
					l7.SetPolicy(attribute.Policy.ValueString())
				}

				l7s = append(l7s, l7)
			}
			firewallAndTrafficShaping.SetL7FirewallRules(l7s)

		}
		var trafficShapingRule []NetworksGroupPolicyResourceModelTrafficShapingRule
		err = firewallAndTrafficShapingData.TrafficShapingRules.ElementsAs(ctx, &trafficShapingRule, false)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to unmarshal TrafficShapingRules",
				fmt.Sprintf("%v", err),
			)
			return

		}
		if len(trafficShapingRule) > 0 {
			var tfs []openApiClient.CreateNetworkGroupPolicyRequestFirewallAndTrafficShapingTrafficShapingRulesInner
			for _, attribute := range trafficShapingRule {
				var tf openApiClient.CreateNetworkGroupPolicyRequestFirewallAndTrafficShapingTrafficShapingRulesInner
				if !attribute.DscpTagValue.IsUnknown() {
					tf.SetDscpTagValue(int32(attribute.DscpTagValue.ValueInt64()))
				}
				if !attribute.PcpTagValue.IsUnknown() {
					tf.SetPcpTagValue(int32(attribute.PcpTagValue.ValueInt64()))
				}
				var perClientBandWidthLimits openApiClient.UpdateNetworkApplianceTrafficShapingRulesRequestRulesInnerPerClientBandwidthLimits
				var perClientBandWidthData NetworksGroupPolicyResourceModelPerClientBandwidth
				perClientBandWidthDataErr := attribute.PerClientBandwidthLimits.As(ctx, &perClientBandWidthData, basetypes.ObjectAsOptions{})
				if perClientBandWidthDataErr != nil {
					resp.Diagnostics.AddError(
						"Failed to unmarshal perClientBandWidth",
						fmt.Sprintf("%v", perClientBandWidthDataErr),
					)
					return

				}

				if !perClientBandWidthData.Settings.IsUnknown() {
					var bandwidthLimits openApiClient.UpdateNetworkApplianceTrafficShapingRulesRequestRulesInnerPerClientBandwidthLimitsBandwidthLimits

					if perClientBandWidthData.Settings.ValueString() != "network default" {
						var perClientBandWidthLimitsData NetworksGroupPolicyResourceModelPerClientBandwidthLimits
						perClientBandWidthLimitsDataErr := attribute.PerClientBandwidthLimits.As(ctx, &perClientBandWidthLimitsData, basetypes.ObjectAsOptions{})
						if perClientBandWidthLimitsDataErr != nil {
							resp.Diagnostics.AddError(
								"Failed to unmarshal perClientBandWidthLimits",
								fmt.Sprintf("%v", perClientBandWidthLimitsDataErr),
							)
							return

						}
						bandwidthLimits.SetLimitDown(int32(perClientBandWidthLimitsData.LimitDown.ValueInt64()))
						bandwidthLimits.SetLimitUp(int32(perClientBandWidthLimitsData.LimitUp.ValueInt64()))

						perClientBandWidthLimits.SetBandwidthLimits(bandwidthLimits)
					}
					perClientBandWidthLimits.SetSettings(perClientBandWidthData.Settings.ValueString())
					tf.SetPerClientBandwidthLimits(perClientBandWidthLimits)
				}
				var defs []openApiClient.UpdateNetworkApplianceTrafficShapingRulesRequestRulesInnerDefinitionsInner
				var definitions []NetworksGroupPolicyResourceModelDefinition
				err = attribute.Definitions.ElementsAs(ctx, &definitions, false)
				if err != nil {
					resp.Diagnostics.AddError(
						"Failed to unmarshal Definitions",
						fmt.Sprintf("%v", err),
					)
					return

				}
				if len(definitions) > 0 {
					for _, attribute := range definitions {
						var def openApiClient.UpdateNetworkApplianceTrafficShapingRulesRequestRulesInnerDefinitionsInner
						def.SetType(attribute.Type.ValueString())
						def.SetValue(attribute.Value.ValueString())
						defs = append(defs, def)
					}
					tf.SetDefinitions(defs)
				}
				tfs = append(tfs, tf)
			}
			firewallAndTrafficShaping.SetTrafficShapingRules(tfs)
		}

		updateNetworkGroupPolicy.SetFirewallAndTrafficShaping(firewallAndTrafficShaping)
	}

	if !data.Scheduling.IsUnknown() && !data.Scheduling.IsNull() {
		var schedulingData NetworksGroupPolicyResourceModelScheduling
		schedulingDataErr := data.Scheduling.As(ctx, &schedulingData, basetypes.ObjectAsOptions{})
		if schedulingDataErr != nil {
			resp.Diagnostics.AddError(
				"Failed to unmarshal Scheduling",
				fmt.Sprintf("%v", schedulingDataErr),
			)
			return
		}

		if !schedulingData.Enabled.IsUnknown() {
			var schedule openApiClient.CreateNetworkGroupPolicyRequestScheduling
			schedule.SetEnabled(schedulingData.Enabled.ValueBool())
			var schedulingFridayDayData NetworksGroupPolicyResourceModelSchedule
			schedulingFridayDayDataErr := schedulingData.Friday.As(ctx, &schedulingFridayDayData, basetypes.ObjectAsOptions{})
			if schedulingFridayDayDataErr != nil {

				resp.Diagnostics.AddError(
					"Failed to unmarshal Friday Scheduling",
					fmt.Sprintf("%v", schedulingFridayDayDataErr),
				)
				return
			}
			if !schedulingFridayDayData.Active.IsUnknown() {
				var friday openApiClient.CreateNetworkGroupPolicyRequestSchedulingFriday
				friday.SetActive(schedulingFridayDayData.Active.ValueBool())
				friday.SetFrom(schedulingFridayDayData.From.ValueString())
				friday.SetTo(schedulingFridayDayData.To.ValueString())
				schedule.SetFriday(friday)
			}
			var schedulingMondayDayData NetworksGroupPolicyResourceModelSchedule
			schedulingMondayDayDataErr := schedulingData.Monday.As(ctx, &schedulingMondayDayData, basetypes.ObjectAsOptions{})
			if schedulingMondayDayDataErr != nil {

				resp.Diagnostics.AddError(
					"Failed to unmarshal Monday Scheduling",
					fmt.Sprintf("%v", schedulingMondayDayDataErr),
				)
				return
			}
			if !schedulingMondayDayData.Active.IsUnknown() {
				var monday openApiClient.CreateNetworkGroupPolicyRequestSchedulingMonday
				monday.SetActive(schedulingMondayDayData.Active.ValueBool())
				monday.SetFrom(schedulingMondayDayData.From.ValueString())
				monday.SetTo(schedulingMondayDayData.To.ValueString())
				schedule.SetMonday(monday)
			}
			var schedulingTuesdayDayData NetworksGroupPolicyResourceModelSchedule
			schedulingTuesdayDayDataErr := schedulingData.Tuesday.As(ctx, &schedulingTuesdayDayData, basetypes.ObjectAsOptions{})
			if schedulingTuesdayDayDataErr != nil {
				resp.Diagnostics.AddError(
					"Failed to unmarshal Tuesday Scheduling",
					fmt.Sprintf("%v", schedulingTuesdayDayDataErr),
				)
				return
			}
			if !schedulingTuesdayDayData.Active.IsUnknown() {
				var tuesday openApiClient.CreateNetworkGroupPolicyRequestSchedulingTuesday
				tuesday.SetActive(schedulingTuesdayDayData.Active.ValueBool())
				tuesday.SetFrom(schedulingTuesdayDayData.From.ValueString())
				tuesday.SetTo(schedulingTuesdayDayData.To.ValueString())
				schedule.SetTuesday(tuesday)
			}
			var schedulingWednesdayDayData NetworksGroupPolicyResourceModelSchedule
			schedulingWednesdayDayDataErr := schedulingData.Wednesday.As(ctx, &schedulingWednesdayDayData, basetypes.ObjectAsOptions{})
			if schedulingWednesdayDayDataErr != nil {
				resp.Diagnostics.AddError(
					"Failed to unmarshal Wednesday Scheduling",
					fmt.Sprintf("%v", schedulingWednesdayDayDataErr),
				)
				return
			}
			if !schedulingWednesdayDayData.Active.IsUnknown() {
				var wednesday openApiClient.CreateNetworkGroupPolicyRequestSchedulingWednesday
				wednesday.SetActive(schedulingWednesdayDayData.Active.ValueBool())
				wednesday.SetFrom(schedulingWednesdayDayData.From.ValueString())
				wednesday.SetTo(schedulingWednesdayDayData.To.ValueString())
				schedule.SetWednesday(wednesday)
			}
			var schedulingThursdayDayData NetworksGroupPolicyResourceModelSchedule
			schedulingThursdayDayDataErr := schedulingData.Thursday.As(ctx, &schedulingThursdayDayData, basetypes.ObjectAsOptions{})
			if schedulingThursdayDayDataErr != nil {
				resp.Diagnostics.AddError(
					"Failed to unmarshal Thursday Scheduling",
					fmt.Sprintf("%v", schedulingThursdayDayDataErr),
				)
				return
			}
			if !schedulingThursdayDayData.Active.IsUnknown() {
				var thursday openApiClient.CreateNetworkGroupPolicyRequestSchedulingThursday
				thursday.SetActive(schedulingThursdayDayData.Active.ValueBool())
				thursday.SetFrom(schedulingThursdayDayData.From.ValueString())
				thursday.SetTo(schedulingThursdayDayData.To.ValueString())
				schedule.SetThursday(thursday)
			}
			var schedulingSaturdayDayData NetworksGroupPolicyResourceModelSchedule
			schedulingSaturdayDayDataErr := schedulingData.Saturday.As(ctx, &schedulingSaturdayDayData, basetypes.ObjectAsOptions{})
			if schedulingSaturdayDayDataErr != nil {
				resp.Diagnostics.AddError(
					"Failed to unmarshal Saturday Scheduling",
					fmt.Sprintf("%v", schedulingSaturdayDayDataErr),
				)
				return
			}
			if !schedulingSaturdayDayData.Active.IsUnknown() {
				var saturday openApiClient.CreateNetworkGroupPolicyRequestSchedulingSaturday
				saturday.SetActive(schedulingSaturdayDayData.Active.ValueBool())
				saturday.SetFrom(schedulingSaturdayDayData.From.ValueString())
				saturday.SetTo(schedulingSaturdayDayData.To.ValueString())
				schedule.SetSaturday(saturday)
			}
			var schedulingSundayDayData NetworksGroupPolicyResourceModelSchedule
			schedulingSundayDayDataErr := schedulingData.Sunday.As(ctx, &schedulingSundayDayData, basetypes.ObjectAsOptions{})
			if schedulingSundayDayDataErr != nil {
				resp.Diagnostics.AddError(
					"Failed to unmarshal Sunday Scheduling",
					fmt.Sprintf("%v", schedulingSundayDayDataErr),
				)
				return
			}
			if !schedulingSundayDayData.Active.IsUnknown() {
				var sunday openApiClient.CreateNetworkGroupPolicyRequestSchedulingSunday
				sunday.SetActive(schedulingSundayDayData.Active.ValueBool())
				sunday.SetFrom(schedulingSundayDayData.From.ValueString())
				sunday.SetTo(schedulingSundayDayData.To.ValueString())
				schedule.SetSunday(sunday)
			}
			updateNetworkGroupPolicy.SetScheduling(schedule)
		}

	}

	if !data.ContentFiltering.IsUnknown() && !data.ContentFiltering.IsNull() {
		var contentFiltering openApiClient.CreateNetworkGroupPolicyRequestContentFiltering
		var contentFilteringData NetworksGroupPolicyResourceModelContentFiltering
		contentDataErr := data.ContentFiltering.As(ctx, &contentFilteringData, basetypes.ObjectAsOptions{})
		if contentDataErr != nil {
			resp.Diagnostics.AddError(
				"Failed to unmarshal ContentFiltering",
				fmt.Sprintf("%v", contentDataErr),
			)
			return
		}
		contentFilteringStatus := false

		var contentFilteringPatternsData NetworksGroupPolicyResourceModelAllowedUrlPatterns
		contentDataPatternsErr := contentFilteringData.AllowedUrlPatterns.As(ctx, &contentFilteringPatternsData, basetypes.ObjectAsOptions{})
		if contentDataPatternsErr != nil {
			resp.Diagnostics.AddError(
				"Failed to unmarshal ContentFilteringPatterns",
				fmt.Sprintf("%v", contentDataPatternsErr),
			)
			return
		}

		if !contentFilteringPatternsData.Settings.IsUnknown() {
			var allowedUrlPatternData openApiClient.CreateNetworkGroupPolicyRequestContentFilteringAllowedUrlPatterns
			allowedUrlPatternData.SetSettings(contentFilteringPatternsData.Settings.ValueString())
			var patternsData []string
			patternsDataerr := contentFilteringPatternsData.Patterns.ElementsAs(ctx, &patternsData, false)
			if patternsDataerr != nil {
				resp.Diagnostics.AddError(
					"Failed to unmarshal AllowedURLPatterns",
					fmt.Sprintf("%v", patternsDataerr),
				)
				return

			}
			allowedUrlPatternData.SetPatterns(patternsData)
			contentFiltering.SetAllowedUrlPatterns(allowedUrlPatternData)
			contentFilteringStatus = true
		}

		var contentFilteringBlockedUrlCategoriesData NetworksGroupPolicyResourceModelBlockedUrlCategories
		contentFilteringBlockedUrlCategoriesDataErr := contentFilteringData.BlockedUrlCategories.As(ctx, &contentFilteringBlockedUrlCategoriesData, basetypes.ObjectAsOptions{})
		if contentFilteringBlockedUrlCategoriesDataErr != nil {
			resp.Diagnostics.AddError(
				"Failed to unmarshal ContentBlockedUrlCategories",
				fmt.Sprintf("%v", contentFilteringBlockedUrlCategoriesDataErr),
			)
			return
		}

		if !contentFilteringBlockedUrlCategoriesData.Settings.IsUnknown() {
			var blockedUrlCategoriesData openApiClient.CreateNetworkGroupPolicyRequestContentFilteringBlockedUrlCategories
			blockedUrlCategoriesData.SetSettings(contentFilteringBlockedUrlCategoriesData.Settings.ValueString())
			var blockedUrlCategories []string
			blockedUrlCategoriesDataerr := contentFilteringBlockedUrlCategoriesData.Categories.ElementsAs(ctx, &blockedUrlCategories, false)
			if blockedUrlCategoriesDataerr != nil {
				resp.Diagnostics.AddError(
					"Failed to unmarshal blockedUrlCategories",
					fmt.Sprintf("%v", blockedUrlCategoriesDataerr),
				)
				return

			}
			blockedUrlCategoriesData.SetCategories(blockedUrlCategories)
			contentFiltering.SetBlockedUrlCategories(blockedUrlCategoriesData)
			contentFilteringStatus = true
		}

		var contentFilteringBlockedUrlPatternsData NetworksGroupPolicyResourceModelBlockedUrlPatterns
		contentFilteringBlockedUrlPatternsDataErr := contentFilteringData.BlockedUrlPatterns.As(ctx, &contentFilteringBlockedUrlPatternsData, basetypes.ObjectAsOptions{})
		if contentFilteringBlockedUrlPatternsDataErr != nil {
			resp.Diagnostics.AddError(
				"Failed to unmarshal ContentBlockedUrlPatterns",
				fmt.Sprintf("%v", contentFilteringBlockedUrlPatternsDataErr),
			)
			return
		}

		if !contentFilteringBlockedUrlPatternsData.Settings.IsUnknown() {
			var blockedUrlPatternsData openApiClient.CreateNetworkGroupPolicyRequestContentFilteringBlockedUrlPatterns
			blockedUrlPatternsData.SetSettings(contentFilteringBlockedUrlPatternsData.Settings.ValueString())
			var blockedUrlPatterns []string
			blockedUrlPatternsDataerr := contentFilteringBlockedUrlPatternsData.Patterns.ElementsAs(ctx, &blockedUrlPatterns, false)
			if blockedUrlPatternsDataerr != nil {
				resp.Diagnostics.AddError(
					"Failed to unmarshal blockedUrlPatterns",
					fmt.Sprintf("%v", blockedUrlPatternsDataerr),
				)
				return

			}
			blockedUrlPatternsData.SetPatterns(blockedUrlPatterns)
			contentFiltering.SetBlockedUrlPatterns(blockedUrlPatternsData)
			contentFilteringStatus = true
		}

		if contentFilteringStatus {
			updateNetworkGroupPolicy.SetContentFiltering(contentFiltering)
		}
	}

	inlineRespMap, httpResp, err := retryAPICall(ctx, func() (interface{}, *http.Response, error) {
		resp, httpResp, err := r.client.NetworksApi.UpdateNetworkGroupPolicy(ctx, data.NetworkId.ValueString(), data.GroupPolicyId.ValueString()).UpdateNetworkGroupPolicyRequest(updateNetworkGroupPolicy).Execute()
		if err != nil {
			return nil, httpResp, err
		}
		return resp, httpResp, nil
	})

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

	inlineResp, ok := inlineRespMap.(map[string]interface{})
	if !ok {
		resp.Diagnostics.AddError(
			"Type Assertion Failure",
			"Failed to assert API response type to map[string]interface{}",
		)
		return
	}

	// Save data into Terraform state
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	bandwidthVal, diags := getBandwidth(ctx, inlineResp)
	if diags.HasError() {
		return
	}
	firewallAndTrafficShapingObjectVal, diags := getFirewallAndTrafficShaping(ctx, inlineResp)
	if diags.HasError() {
		return
	}
	schedulingVal, diags := getScheduling(ctx, inlineResp)
	if diags.HasError() {
		return
	}
	contentFilteringVal, diags := getContentFiltering(ctx, inlineResp)
	if diags.HasError() {
		return
	}
	bonjourForwardingVal, diags := getBonjourForwarding(ctx, inlineResp)
	if diags.HasError() {
		return
	}

	vlantaggingVal, diags := getVlantagging(ctx, inlineResp)
	if diags.HasError() {
		return
	}

	data.Bandwidth = bandwidthVal
	data.FirewallAndTrafficShaping = firewallAndTrafficShapingObjectVal
	data.Scheduling = schedulingVal
	data.ContentFiltering = contentFilteringVal
	data.BonjourForwarding = bonjourForwardingVal
	data.VlanTagging = vlantaggingVal
	data.Id = jsontypes.StringValue("example-id")

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

	_, httpResp, err := retryAPICall(ctx, func() (interface{}, *http.Response, error) {
		httpResp, err := r.client.NetworksApi.DeleteNetworkGroupPolicy(context.Background(), data.NetworkId.ValueString(), data.GroupPolicyId.ValueString()).Execute()
		if err != nil {
			return nil, httpResp, err
		}
		return resp, httpResp, nil
	})

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
