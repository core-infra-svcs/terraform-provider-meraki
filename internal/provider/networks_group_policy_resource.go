package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/meraki/dashboard-api-go/client"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// NetworksGroupPolicyResource defines the resource implementation.
type NetworksGroupPolicyResource struct {
	client *client.APIClient
}

func NewNetworksGroupPolicyResource() resource.Resource {
	return &NetworksGroupPolicyResource{}
}

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &NetworksGroupPolicyResource{}
var _ resource.ResourceWithImportState = &NetworksGroupPolicyResource{}

// GroupPolicyResourceModel represents a group policy.
type GroupPolicyResourceModel struct {
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
	Enabled   types.Bool        `tfsdk:"enabled" json:"enabled"`
	Monday    *ScheduleDayModel `tfsdk:"monday" json:"monday"`
	Tuesday   *ScheduleDayModel `tfsdk:"tuesday" json:"tuesday"`
	Wednesday *ScheduleDayModel `tfsdk:"wednesday" json:"wednesday"`
	Thursday  *ScheduleDayModel `tfsdk:"thursday" json:"thursday"`
	Friday    *ScheduleDayModel `tfsdk:"friday" json:"friday"`
	Saturday  *ScheduleDayModel `tfsdk:"saturday" json:"saturday"`
	Sunday    *ScheduleDayModel `tfsdk:"sunday" json:"sunday"`
}

// ScheduleDayModel represents a single day's schedule.
type ScheduleDayModel struct {
	Active types.Bool   `tfsdk:"active" json:"active"`
	From   types.String `tfsdk:"from" json:"from"`
	To     types.String `tfsdk:"to" json:"to"`
}

// BandwidthModel represents the bandwidth settings.
type BandwidthModel struct {
	Settings        types.String          `tfsdk:"settings" json:"settings"`
	BandwidthLimits *BandwidthLimitsModel `tfsdk:"bandwidth_limits" json:"bandwidthLimits"`
}

// BandwidthLimitsModel represents the bandwidth limits.
type BandwidthLimitsModel struct {
	LimitUp   types.Int64 `tfsdk:"limit_up" json:"limitUp"`
	LimitDown types.Int64 `tfsdk:"limit_down" json:"limitDown"`
}

// FirewallAndTrafficShapingModel represents the firewall and traffic shaping settings.
type FirewallAndTrafficShapingModel struct {
	Settings            types.String              `tfsdk:"settings" json:"settings"`
	L3FirewallRules     []L3FirewallRuleModel     `tfsdk:"l3_firewall_rules" json:"l3FirewallRules"`
	L7FirewallRules     []L7FirewallRuleModel     `tfsdk:"l7_firewall_rules" json:"l7FirewallRules"`
	TrafficShapingRules []TrafficShapingRuleModel `tfsdk:"traffic_shaping_rules" json:"trafficShapingRules"`
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
	DscpTagValue             types.Int64                     `tfsdk:"dscp_tag_value" json:"dscpTagValue"`
	PcpTagValue              types.Int64                     `tfsdk:"pcp_tag_value" json:"pcpTagValue"`
	PerClientBandwidthLimits *PerClientBandwidthLimitsModel  `tfsdk:"per_client_bandwidth_limits" json:"perClientBandwidthLimits"`
	Definitions              []TrafficShapingDefinitionModel `tfsdk:"definitions" json:"definitions"`
}

// PerClientBandwidthLimitsModel represents the per-client bandwidth limits.
type PerClientBandwidthLimitsModel struct {
	Settings        types.String          `tfsdk:"settings" json:"settings"`
	BandwidthLimits *BandwidthLimitsModel `tfsdk:"bandwidth_limits" json:"bandwidthLimits"`
}

// TrafficShapingDefinitionModel represents a traffic shaping definition.
type TrafficShapingDefinitionModel struct {
	Type  types.String `tfsdk:"type" json:"type"`
	Value types.String `tfsdk:"value" json:"value"`
}

// ContentFilteringModel represents the content filtering settings.
type ContentFilteringModel struct {
	AllowedUrlPatterns   UrlPatterns   `tfsdk:"allowed_url_patterns" json:"allowedUrlPatterns"`
	BlockedUrlPatterns   UrlPatterns   `tfsdk:"blocked_url_patterns" json:"blockedUrlPatterns"`
	BlockedUrlCategories UrlCategories `tfsdk:"blocked_url_categories" json:"blockedUrlCategories"`
}

type UrlPatterns struct {
	Patterns types.List   `tfsdk:"patterns" json:"patterns"`
	Settings types.String `tfsdk:"settings" json:"settings"`
}

type UrlCategories struct {
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
	Settings types.String                 `tfsdk:"settings" json:"settings"`
	Rules    []BonjourForwardingRuleModel `tfsdk:"rules" json:"rules"`
}

// BonjourForwardingRuleModel represents a Bonjour forwarding rule.
type BonjourForwardingRuleModel struct {
	Description types.String   `tfsdk:"description" json:"description"`
	VlanID      types.String   `tfsdk:"vlan_id" json:"vlanId"`
	Services    []types.String `tfsdk:"services" json:"services"`
}

// Schema defines the schema for the resource.
func (r *NetworksGroupPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The resource import id",
			},
			"group_policy_id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier for the group policy.",
			},
			"network_id": schema.StringAttribute{
				Required:    true,
				Description: "The network ID where the group policy is applied.",
			},
			"name": schema.StringAttribute{
				Required:    true,
				Description: "The name of the group policy.",
			},
			"scheduling": schema.SingleNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The scheduling settings of the group policy.",
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						Required: true,
					},
					"sunday": schema.SingleNestedAttribute{
						Optional: true,
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"active": schema.BoolAttribute{
								Required: true,
							},
							"from": schema.StringAttribute{
								Required: true,
							},
							"to": schema.StringAttribute{
								Required: true,
							},
						},
					},
					"monday": schema.SingleNestedAttribute{
						Optional: true,
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"active": schema.BoolAttribute{
								Required: true,
							},
							"from": schema.StringAttribute{
								Required: true,
							},
							"to": schema.StringAttribute{
								Required: true,
							},
						},
					},
					"tuesday": schema.SingleNestedAttribute{
						Optional: true,
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"active": schema.BoolAttribute{
								Required: true,
							},
							"from": schema.StringAttribute{
								Required: true,
							},
							"to": schema.StringAttribute{
								Required: true,
							},
						},
					},
					"wednesday": schema.SingleNestedAttribute{
						Optional: true,
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"active": schema.BoolAttribute{
								Required: true,
							},
							"from": schema.StringAttribute{
								Required: true,
							},
							"to": schema.StringAttribute{
								Required: true,
							},
						},
					},
					"thursday": schema.SingleNestedAttribute{
						Optional: true,
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"active": schema.BoolAttribute{
								Required: true,
							},
							"from": schema.StringAttribute{
								Required: true,
							},
							"to": schema.StringAttribute{
								Required: true,
							},
						},
					},
					"friday": schema.SingleNestedAttribute{
						Optional: true,
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"active": schema.BoolAttribute{
								Required: true,
							},
							"from": schema.StringAttribute{
								Required: true,
							},
							"to": schema.StringAttribute{
								Required: true,
							},
						},
					},
					"saturday": schema.SingleNestedAttribute{
						Optional: true,
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"active": schema.BoolAttribute{
								Required: true,
							},
							"from": schema.StringAttribute{
								Required: true,
							},
							"to": schema.StringAttribute{
								Required: true,
							},
						},
					},
				},
			},
			"bandwidth": schema.SingleNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The bandwidth settings of the group policy.",
				Attributes: map[string]schema.Attribute{
					"settings": schema.StringAttribute{
						Required: true,
					},
					"bandwidth_limits": schema.SingleNestedAttribute{
						Optional: true,
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"limit_up": schema.Int64Attribute{
								Optional:    true,
								Computed:    true,
								Description: "The upload bandwidth limit. Can be null.",
							},
							"limit_down": schema.Int64Attribute{
								Optional:    true,
								Computed:    true,
								Description: "The download bandwidth limit. Can be null.",
							},
						},
					},
				},
			},
			"firewall_and_traffic_shaping": schema.SingleNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The firewall and traffic shaping settings of the group policy.",
				Attributes: map[string]schema.Attribute{
					"settings": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"l3_firewall_rules": schema.ListNestedAttribute{
						Optional: true,
						Computed: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"comment": schema.StringAttribute{
									Optional: true,
									Computed: true,
								},
								"policy": schema.StringAttribute{
									Optional: true,
									Computed: true,
								},
								"protocol": schema.StringAttribute{
									Optional: true,
									Computed: true,
								},
								"dest_port": schema.StringAttribute{
									Optional: true,
									Computed: true,
								},
								"dest_cidr": schema.StringAttribute{
									Optional: true,
									Computed: true,
								},
							},
						},
					},
					"l7_firewall_rules": schema.ListNestedAttribute{
						Optional: true,
						Computed: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"policy": schema.StringAttribute{
									Optional: true,
									Computed: true,
								},
								"type": schema.StringAttribute{
									Optional: true,
									Computed: true,
								},
								"value": schema.StringAttribute{
									Optional: true,
									Computed: true,
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
									Optional: true,
									Computed: true,
								},
								"pcp_tag_value": schema.Int64Attribute{
									Optional: true,
									Computed: true,
								},
								"per_client_bandwidth_limits": schema.SingleNestedAttribute{
									Optional: true,
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"settings": schema.StringAttribute{
											Required: true,
										},
										"bandwidth_limits": schema.SingleNestedAttribute{
											Optional: true,
											Computed: true,
											Attributes: map[string]schema.Attribute{
												"limit_up": schema.Int64Attribute{
													Optional: true,
													Computed: true,
												},
												"limit_down": schema.Int64Attribute{
													Optional: true,
													Computed: true,
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
												Required: true,
											},
											"value": schema.StringAttribute{
												Required: true,
											},
										},
									},
								},
							},
						},
					},
				},
			},
			"content_filtering": schema.SingleNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The content filtering settings of the group policy.",
				Attributes: map[string]schema.Attribute{
					"allowed_url_patterns": schema.SingleNestedAttribute{
						Optional: true,
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"patterns": schema.ListAttribute{
								ElementType: types.StringType,
								Optional:    true,
								Computed:    true,
							},
							"settings": schema.StringAttribute{
								Optional: true,
								Computed: true,
							},
						},
					},
					"blocked_url_patterns": schema.SingleNestedAttribute{
						Optional: true,
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"patterns": schema.ListAttribute{
								ElementType: types.StringType,
								Optional:    true,
								Computed:    true,
							},
							"settings": schema.StringAttribute{
								Optional: true,
								Computed: true,
							},
						},
					},
					"blocked_url_categories": schema.SingleNestedAttribute{
						Optional: true,
						Computed: true,
						Attributes: map[string]schema.Attribute{
							"categories": schema.ListAttribute{
								ElementType: types.StringType,
								Optional:    true,
								Computed:    true,
							},
							"settings": schema.StringAttribute{
								Optional: true,
								Computed: true,
							},
						},
					},
				},
			},
			"splash_auth_settings": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The splash authentication settings of the group policy.",
			},
			"vlan_tagging": schema.SingleNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The VLAN tagging settings of the group policy.",
				Attributes: map[string]schema.Attribute{
					"settings": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
					"vlan_id": schema.StringAttribute{
						Optional: true,
						Computed: true,
					},
				},
			},
			"bonjour_forwarding": schema.SingleNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The bonjour forwarding settings of the group policy.",
				Attributes: map[string]schema.Attribute{
					"settings": schema.StringAttribute{
						Required: true,
					},
					"rules": schema.ListNestedAttribute{
						Optional: true,
						Computed: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"description": schema.StringAttribute{
									Optional: true,
									Computed: true,
								},
								"vlan_id": schema.StringAttribute{
									Optional: true,
									Computed: true,
								},
								"services": schema.ListAttribute{
									ElementType: types.StringType,
									Optional:    true,
									Computed:    true,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *NetworksGroupPolicyResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_group_policy"
}

// Configure configures the resource with the API client.
func (r *NetworksGroupPolicyResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.APIClient)

}

func updateGroupPolicyResourceState(ctx context.Context, data *GroupPolicyResourceModel, groupPolicy map[string]interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Update GroupPolicyId
	groupPolicyId, groupPolicyIdOk := groupPolicy["groupPolicyId"].(string)
	if groupPolicyIdOk {
		data.GroupPolicyId = types.StringValue(groupPolicyId)
	}

	// Safety check for nil pointer dereference on data.NetworkId
	if data.NetworkId.IsNull() && data.NetworkId.IsUnknown() {
		networkId, networkIdOk := groupPolicy["networkId"].(string)
		if networkIdOk {
			data.NetworkId = types.StringValue(networkId)
		}
	}

	// Construct ID
	if !data.NetworkId.IsNull() && !data.NetworkId.IsUnknown() && !data.GroupPolicyId.IsNull() && !data.GroupPolicyId.IsUnknown() {
		data.ID = types.StringValue(data.NetworkId.ValueString() + "," + data.GroupPolicyId.ValueString())
	} else {
		data.ID = types.StringNull()
	}

	// Update Name
	if name, nameOk := groupPolicy["name"].(string); nameOk {
		data.Name = types.StringValue(name)
	} else {
		data.Name = types.StringNull()
	}

	// Update Scheduling
	if scheduling, schedulingOk := groupPolicy["scheduling"].(map[string]interface{}); schedulingOk {
		if scheduling != nil {

			schedulingObj, schedulingDiags := updateSchedulingState(ctx, scheduling)
			if schedulingDiags.HasError() {
				diags.Append(schedulingDiags...)
			}
			data.Scheduling = schedulingObj

		}
	} else {
		data.Scheduling = types.ObjectNull(map[string]attr.Type{
			"enabled":   types.BoolType,
			"monday":    types.ObjectType{AttrTypes: map[string]attr.Type{"active": types.BoolType, "from": types.StringType, "to": types.StringType}},
			"tuesday":   types.ObjectType{AttrTypes: map[string]attr.Type{"active": types.BoolType, "from": types.StringType, "to": types.StringType}},
			"wednesday": types.ObjectType{AttrTypes: map[string]attr.Type{"active": types.BoolType, "from": types.StringType, "to": types.StringType}},
			"thursday":  types.ObjectType{AttrTypes: map[string]attr.Type{"active": types.BoolType, "from": types.StringType, "to": types.StringType}},
			"friday":    types.ObjectType{AttrTypes: map[string]attr.Type{"active": types.BoolType, "from": types.StringType, "to": types.StringType}},
			"saturday":  types.ObjectType{AttrTypes: map[string]attr.Type{"active": types.BoolType, "from": types.StringType, "to": types.StringType}},
			"sunday":    types.ObjectType{AttrTypes: map[string]attr.Type{"active": types.BoolType, "from": types.StringType, "to": types.StringType}},
		})
	}

	// Update Bandwidth
	if bandwidth, ok := groupPolicy["bandwidth"].(map[string]interface{}); ok {
		var bandwidthModel BandwidthModel

		if settings, ok := bandwidth["settings"].(string); ok {
			bandwidthModel.Settings = types.StringValue(settings)
		} else {
			bandwidthModel.Settings = types.StringNull()
		}

		if limits, ok := bandwidth["bandwidth_limits"].(map[string]interface{}); ok {
			if limitUp, ok := limits["limit_up"].(float64); ok {
				bandwidthModel.BandwidthLimits.LimitUp = types.Int64Value(int64(limitUp))
			} else {
				bandwidthModel.BandwidthLimits.LimitUp = types.Int64Null()
			}

			if limitDown, ok := limits["limit_down"].(float64); ok {
				bandwidthModel.BandwidthLimits.LimitDown = types.Int64Value(int64(limitDown))
			} else {
				bandwidthModel.BandwidthLimits.LimitDown = types.Int64Null()
			}
		} else {
			bandwidthModel.BandwidthLimits = &BandwidthLimitsModel{}
		}

		// Convert BandwidthModel to types.Object
		bandwidthObject, err := types.ObjectValueFrom(ctx, map[string]attr.Type{
			"settings": types.StringType,
			"bandwidth_limits": types.ObjectType{AttrTypes: map[string]attr.Type{
				"limit_up":   types.Float64Type,
				"limit_down": types.Float64Type,
			}},
		}, bandwidthModel)
		if err != nil {
			diags.Append(err...)
			return diags
		}
		data.Bandwidth = bandwidthObject
	} else {
		data.Bandwidth = types.ObjectNull(map[string]attr.Type{
			"settings": types.StringType,
			"bandwidth_limits": types.ObjectType{AttrTypes: map[string]attr.Type{
				"limit_up":   types.Float64Type,
				"limit_down": types.Float64Type,
			}},
		})
	}

	//SplashAuthSettings
	if splashAuthSettings, splashAuthSettingsOk := groupPolicy["splashAuthSettings"].(string); splashAuthSettingsOk {
		data.SplashAuthSettings = types.StringValue(splashAuthSettings)
	} else {
		data.SplashAuthSettings = types.StringNull()
	}

	// Update VlanTagging
	if vlanTagging, ok := groupPolicy["vlanTagging"].(map[string]interface{}); ok {
		vlanTaggingObj, vlanTaggingDiags := updateVlanTaggingState(ctx, vlanTagging)
		if vlanTaggingDiags.HasError() {
			diags.Append(vlanTaggingDiags...)
		}
		data.VlanTagging = vlanTaggingObj
	} else {
		data.VlanTagging = types.ObjectNull(map[string]attr.Type{
			"settings": types.StringType,
			"vlan_id":  types.StringType,
		})
	}

	// Update BonjourForwarding
	if bonjourForwarding, ok := groupPolicy["bonjourForwarding"].(map[string]interface{}); ok {
		settings := types.StringNull()
		if s, ok := bonjourForwarding["settings"].(string); ok {
			settings = types.StringValue(s)
		}
		rules := types.ListNull(types.ObjectType{AttrTypes: map[string]attr.Type{
			"description": types.StringType,
			"vlan_id":     types.StringType,
			"services":    types.ListType{ElemType: types.StringType},
		}})

		if r, ok := bonjourForwarding["rules"].([]interface{}); ok {
			var rulesList []attr.Value
			for _, rule := range r {
				ruleMap := rule.(map[string]interface{})
				description := types.StringNull()
				if d, ok := ruleMap["description"].(string); ok {
					description = types.StringValue(d)
				}
				vlanID := types.StringNull()
				if v, ok := ruleMap["vlanId"].(string); ok {
					vlanID = types.StringValue(v)
				}
				services := types.ListNull(types.StringType)
				if s, ok := ruleMap["services"].([]interface{}); ok {
					var serviceList []attr.Value
					for _, service := range s {
						if svc, ok := service.(string); ok {
							serviceList = append(serviceList, types.StringValue(svc))
						}
					}
					services = types.ListValueMust(types.StringType, serviceList)
				}
				ruleObj, _ := types.ObjectValue(
					map[string]attr.Type{
						"description": types.StringType,
						"vlan_id":     types.StringType,
						"services":    types.ListType{ElemType: types.StringType},
					},
					map[string]attr.Value{
						"description": description,
						"vlan_id":     vlanID,
						"services":    services,
					},
				)
				rulesList = append(rulesList, ruleObj)
			}
			rules = types.ListValueMust(types.ObjectType{AttrTypes: map[string]attr.Type{"description": types.StringType, "vlan_id": types.StringType, "services": types.ListType{ElemType: types.StringType}}}, rulesList)
		}

		bonjourForwardingObj, err := types.ObjectValue(
			map[string]attr.Type{
				"settings": types.StringType,
				"rules":    types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"description": types.StringType, "vlan_id": types.StringType, "services": types.ListType{ElemType: types.StringType}}}},
			},
			map[string]attr.Value{
				"settings": settings,
				"rules":    rules,
			},
		)
		if err.HasError() {
			diags = append(diags, err...)
		}
		data.BonjourForwarding = bonjourForwardingObj
	} else {
		data.BonjourForwarding = types.ObjectNull(map[string]attr.Type{
			"settings": types.StringType,
			"rules":    types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"description": types.StringType, "vlan_id": types.StringType, "services": types.ListType{ElemType: types.StringType}}}},
		})
	}

	// Update FirewallAndTrafficShaping
	if firewallAndTrafficShaping, ok := groupPolicy["firewallAndTrafficShaping"].(map[string]interface{}); ok {

		firewallAndTrafficShapingObject, firewallAndTrafficShapingObjectDiags := updateFirewallAndTrafficShapingRules(ctx, firewallAndTrafficShaping)
		if firewallAndTrafficShapingObjectDiags.HasError() {
			diags.Append(firewallAndTrafficShapingObjectDiags...)
		}
		data.FirewallAndTrafficShaping = firewallAndTrafficShapingObject

	} else {
		data.FirewallAndTrafficShaping = types.ObjectNull(
			map[string]attr.Type{
				"settings": types.StringType,
				"l3_firewall_rules": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
					"comment":   types.StringType,
					"policy":    types.StringType,
					"protocol":  types.StringType,
					"dest_port": types.StringType,
					"dest_cidr": types.StringType}}},
				"l7_firewall_rules": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
					"policy": types.StringType,
					"type":   types.StringType,
					"value":  types.StringType}}},
				"traffic_shaping_rules": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
					"dscp_tag_value": types.Int64Type,
					"pcp_tag_value":  types.Int64Type,
					"per_client_bandwidth_limits": types.ObjectType{AttrTypes: map[string]attr.Type{
						"settings": types.StringType,
						"bandwidth_limits": types.ObjectType{AttrTypes: map[string]attr.Type{
							"limit_up":   types.Int64Type,
							"limit_down": types.Int64Type}}}},
					"definitions": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
						"type":  types.StringType,
						"value": types.StringType}}}}}},
			},
		)
	}

	// Update ContentFiltering
	if contentFiltering, ok := groupPolicy["contentFiltering"].(map[string]interface{}); ok {
		contentFilteringObj, cfDiags := updateContentFilteringState(ctx, contentFiltering)
		if !cfDiags.HasError() {
			diags.Append(cfDiags...)
		}
		data.ContentFiltering = contentFilteringObj
	} else {
		data.ContentFiltering = types.ObjectNull(map[string]attr.Type{
			"allowed_url_patterns": types.ObjectType{AttrTypes: map[string]attr.Type{
				"patterns": types.ListType{ElemType: types.StringType},
				"settings": types.StringType,
			}},
			"blocked_url_patterns": types.ObjectType{AttrTypes: map[string]attr.Type{
				"patterns": types.ListType{ElemType: types.StringType},
				"settings": types.StringType,
			}},
			"blocked_url_categories": types.ObjectType{AttrTypes: map[string]attr.Type{
				"categories": types.ListType{ElemType: types.StringType},
				"settings":   types.StringType,
			}},
		})
	}

	return diags
}

func updateSchedulingState(ctx context.Context, scheduling map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	newSchedulingObject := SchedulingModel{}

	//enabled
	enabled := types.BoolNull()
	if e, eOk := scheduling["enabled"].(bool); eOk {
		enabled = types.BoolValue(e)
	}
	newSchedulingObject.Enabled = enabled

	// Days of the week
	updateScheduleDayModel := func(dayAttrName string) (ScheduleDayModel, diag.Diagnostics) {
		var ScheduleDayDiags diag.Diagnostics
		var scheduleDay ScheduleDayModel

		if day, dayOk := scheduling[dayAttrName].(map[string]interface{}); dayOk {

			// Active
			active := types.BoolNull()
			if a, aOk := day["active"].(bool); aOk {
				active = types.BoolValue(a)
			}
			scheduleDay.Active = active

			// To
			to := types.StringNull()
			if t, tOk := day["to"].(string); tOk {
				to = types.StringValue(t)
			}
			scheduleDay.To = to

			// From
			from := types.StringNull()
			if f, fOk := day["from"].(string); fOk {
				from = types.StringValue(f)
			}
			scheduleDay.From = from

			return scheduleDay, ScheduleDayDiags
		}

		return scheduleDay, ScheduleDayDiags
	}

	monday, monDiags := updateScheduleDayModel("monday")
	if monDiags.HasError() {
		diags.Append(monDiags...)
	}
	tuesday, tuesDiags := updateScheduleDayModel("tuesday")
	if tuesDiags.HasError() {
		diags.Append(tuesDiags...)
	}
	wednesday, wedDiags := updateScheduleDayModel("wednesday")
	if wedDiags.HasError() {
		diags.Append(wedDiags...)
	}
	thursday, thurDiags := updateScheduleDayModel("thursday")
	if thurDiags.HasError() {
		diags.Append(thurDiags...)
	}
	friday, friDiags := updateScheduleDayModel("friday")
	if tuesDiags.HasError() {
		diags.Append(friDiags...)
	}
	saturday, satDiags := updateScheduleDayModel("saturday")
	if satDiags.HasError() {
		diags.Append(satDiags...)
	}
	sunday, sunDiags := updateScheduleDayModel("sunday")
	if sunDiags.HasError() {
		diags.Append(sunDiags...)
	}

	newSchedulingObject.Monday = &monday
	newSchedulingObject.Tuesday = &tuesday
	newSchedulingObject.Wednesday = &wednesday
	newSchedulingObject.Thursday = &thursday
	newSchedulingObject.Friday = &friday
	newSchedulingObject.Saturday = &saturday
	newSchedulingObject.Sunday = &sunday

	nsod, err := types.ObjectValueFrom(ctx, map[string]attr.Type{
		"enabled":   types.BoolType,
		"monday":    types.ObjectType{AttrTypes: map[string]attr.Type{"active": types.BoolType, "from": types.StringType, "to": types.StringType}},
		"tuesday":   types.ObjectType{AttrTypes: map[string]attr.Type{"active": types.BoolType, "from": types.StringType, "to": types.StringType}},
		"wednesday": types.ObjectType{AttrTypes: map[string]attr.Type{"active": types.BoolType, "from": types.StringType, "to": types.StringType}},
		"thursday":  types.ObjectType{AttrTypes: map[string]attr.Type{"active": types.BoolType, "from": types.StringType, "to": types.StringType}},
		"friday":    types.ObjectType{AttrTypes: map[string]attr.Type{"active": types.BoolType, "from": types.StringType, "to": types.StringType}},
		"saturday":  types.ObjectType{AttrTypes: map[string]attr.Type{"active": types.BoolType, "from": types.StringType, "to": types.StringType}},
		"sunday":    types.ObjectType{AttrTypes: map[string]attr.Type{"active": types.BoolType, "from": types.StringType, "to": types.StringType}},
	}, newSchedulingObject)
	if err.HasError() {
		diags.Append(err...)
	}

	return nsod, diags

}

func updateFirewallAndTrafficShapingRules(ctx context.Context, firewallAndTrafficShapingRules map[string]interface{}) (types.Object, diag.Diagnostics) {
	diags := diag.Diagnostics{}

	firewallAndTrafficShapingObjectNull := types.ObjectNull(map[string]attr.Type{
		"settings":              types.StringType,
		"l3_firewall_rules":     types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"comment": types.StringType, "policy": types.StringType, "protocol": types.StringType, "dest_port": types.StringType, "dest_cidr": types.StringType}}},
		"l7_firewall_rules":     types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"policy": types.StringType, "type": types.StringType, "value": types.StringType}}},
		"traffic_shaping_rules": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"dscp_tag_value": types.Int64Type, "pcp_tag_value": types.Int64Type, "per_client_bandwidth_limits": types.ObjectType{AttrTypes: map[string]attr.Type{"settings": types.StringType, "bandwidth_limits": types.ObjectType{AttrTypes: map[string]attr.Type{"limit_up": types.Int64Type, "limit_down": types.Int64Type}}}}, "definitions": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"type": types.StringType, "value": types.StringType}}}}}},
	})

	// FirewallAndTrafficShapingRules
	if firewallAndTrafficShapingRules != nil {

		// settings
		settings := types.StringNull()
		if set, setOk := firewallAndTrafficShapingRules["settings"].(string); setOk {
			settings = types.StringValue(set)
		}

		// l3FirewallRules
		var l3FirewallRules []L3FirewallRuleModel
		if l3frs, l3frsOk := firewallAndTrafficShapingRules["l3FirewallRules"].([]interface{}); l3frsOk {

			for _, l3fwr := range l3frs {

				if i, iOk := l3fwr.(map[string]interface{}); iOk {
					rule := L3FirewallRuleModel{
						Comment:  types.StringNull(),
						Policy:   types.StringNull(),
						Protocol: types.StringNull(),
						DestPort: types.StringNull(),
						DestCidr: types.StringNull(),
					}

					if comment, commentOk := i["comment"].(string); commentOk {
						rule.Comment = types.StringValue(comment)
					}
					if policy, policyOk := i["policy"].(string); policyOk {
						rule.Policy = types.StringValue(policy)
					}
					if protocol, protocolOk := i["protocol"].(string); protocolOk {
						rule.Protocol = types.StringValue(protocol)
					}
					if destPort, destPortOk := i["destPort"].(string); destPortOk {
						rule.DestPort = types.StringValue(destPort)
					}
					if destCidr, destCidrOk := i["destCidr"].(string); destCidrOk {
						rule.DestCidr = types.StringValue(destCidr)
					}

					l3FirewallRules = append(l3FirewallRules, rule)
				}

			}
		}

		// l7FirewallRules
		var l7FirewallRules []L7FirewallRuleModel
		if l7fwrs, l7fwrsOk := firewallAndTrafficShapingRules["l7FirewallRules"].([]interface{}); l7fwrsOk {

			for _, l7fwr := range l7fwrs {

				if l7, l7Ok := l7fwr.(map[string]interface{}); l7Ok {

					// policy
					policy := types.StringNull()
					if p, policyOk := l7["policy"].(string); policyOk {
						policy = types.StringValue(p)
					}

					// type
					typ := types.StringNull()
					if t, typeOk := l7["type"].(string); typeOk {
						typ = types.StringValue(t)
					}

					//value
					val := types.StringNull()
					if v, valueOk := l7["value"].(string); valueOk {
						val = types.StringValue(v)
					}

					rule := L7FirewallRuleModel{
						Policy: policy,
						Type:   typ,
						Value:  val,
					}

					l7FirewallRules = append(l7FirewallRules, rule)

				}

			}
		}

		// trafficShapingRules
		var trafficShapingRules []TrafficShapingRuleModel
		if tsrs, tsrsOk := firewallAndTrafficShapingRules["trafficShapingRules"].([]interface{}); tsrsOk {

			for _, tsr := range tsrs {

				if sr, iOk := tsr.(map[string]interface{}); iOk {

					var trafficShapingRule TrafficShapingRuleModel

					// dscpTagValue
					if dscpTagValue, dscpTagValueOk := sr["dscpTagValue"].(float64); dscpTagValueOk {
						trafficShapingRule.DscpTagValue = types.Int64Value(int64(dscpTagValue))
					}

					// pcpTagValue
					if pcpTagValue, pcpTagValueOk := sr["pcpTagValue"].(float64); pcpTagValueOk {
						trafficShapingRule.PcpTagValue = types.Int64Value(int64(pcpTagValue))
					}

					// perClientBandwidthLimits
					if pcbl, pcblOk := sr["perClientBandwidthLimits"].(map[string]interface{}); pcblOk {
						perClientBandwidthLimits := PerClientBandwidthLimitsModel{
							BandwidthLimits: &BandwidthLimitsModel{}, // Initialize BandwidthLimits
						}

						// settings
						if set, settingsOk := pcbl["settings"].(string); settingsOk {
							perClientBandwidthLimits.Settings = types.StringValue(set)
						}

						// bandwidthLimits
						if bandwidthLimits, bandwidthLimitsOk := pcbl["bandwidthLimits"].(map[string]interface{}); bandwidthLimitsOk {
							if limitUp, limitUpOk := bandwidthLimits["limitUp"].(float64); limitUpOk {
								perClientBandwidthLimits.BandwidthLimits.LimitUp = types.Int64Value(int64(limitUp))
							}
							if limitDown, limitDownOk := bandwidthLimits["limitDown"].(float64); limitDownOk {
								perClientBandwidthLimits.BandwidthLimits.LimitDown = types.Int64Value(int64(limitDown))
							}
						}

						trafficShapingRule.PerClientBandwidthLimits = &perClientBandwidthLimits
					}

					// definitions
					if defs, defsOk := sr["definitions"].([]interface{}); defsOk {
						var definitions []TrafficShapingDefinitionModel

						for _, def := range defs {

							if ef, efOk := def.(map[string]interface{}); efOk {

								definition := TrafficShapingDefinitionModel{}

								if typ, typeOk := ef["type"].(string); typeOk {
									definition.Type = types.StringValue(typ)
								}

								if value, valueOk := ef["value"].(string); valueOk {
									definition.Value = types.StringValue(value)
								}

								definitions = append(definitions, definition)
							}

						}

						trafficShapingRule.Definitions = definitions
					}

					// update data
					trafficShapingRules = append(trafficShapingRules, trafficShapingRule)

				}

			}
		}

		// firewallAndTrafficShaping Data
		firewallAndTrafficShapingData := FirewallAndTrafficShapingModel{
			Settings:            settings,
			L3FirewallRules:     l3FirewallRules,
			L7FirewallRules:     l7FirewallRules,
			TrafficShapingRules: trafficShapingRules,
		}

		firewallAndTrafficShapingObj, err := types.ObjectValueFrom(ctx, map[string]attr.Type{
			"settings":              types.StringType,
			"l3_firewall_rules":     types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"comment": types.StringType, "policy": types.StringType, "protocol": types.StringType, "dest_port": types.StringType, "dest_cidr": types.StringType}}},
			"l7_firewall_rules":     types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"policy": types.StringType, "type": types.StringType, "value": types.StringType}}},
			"traffic_shaping_rules": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"dscp_tag_value": types.Int64Type, "pcp_tag_value": types.Int64Type, "per_client_bandwidth_limits": types.ObjectType{AttrTypes: map[string]attr.Type{"settings": types.StringType, "bandwidth_limits": types.ObjectType{AttrTypes: map[string]attr.Type{"limit_up": types.Int64Type, "limit_down": types.Int64Type}}}}, "definitions": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"type": types.StringType, "value": types.StringType}}}}}},
		}, firewallAndTrafficShapingData)
		if err.HasError() {
			diags.Append(err...)
			return firewallAndTrafficShapingObjectNull, diags
		}

		return firewallAndTrafficShapingObj, diags
	}

	return firewallAndTrafficShapingObjectNull, diags
}

func updateContentFilteringState(ctx context.Context, contentFiltering map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	contentFilteringObjectNull := types.ObjectNull(map[string]attr.Type{
		"allowed_url_patterns": types.ObjectType{AttrTypes: map[string]attr.Type{
			"patterns": types.ListType{ElemType: types.StringType},
			"settings": types.StringType,
		}},
		"blocked_url_patterns": types.ObjectType{AttrTypes: map[string]attr.Type{
			"patterns": types.ListType{ElemType: types.StringType},
			"settings": types.StringType,
		}},
		"blocked_url_categories": types.ObjectType{AttrTypes: map[string]attr.Type{
			"categories": types.ListType{ElemType: types.StringType},
			"settings":   types.StringType,
		}},
	})

	// AllowedUrlPatterns
	allowedUrlPatterns := UrlPatterns{}
	if aup, aupOk := contentFiltering["allowedUrlPatterns"].(map[string]interface{}); aupOk {

		// Settings
		settings := types.StringNull()
		if s, ok := aup["settings"].(string); ok {
			settings = types.StringValue(s)
		}
		allowedUrlPatterns.Settings = settings

		// Patterns
		var patternList []string
		if patterns, patternsOk := aup["patterns"].([]interface{}); patternsOk {

			for _, pattern := range patterns {
				if p, pOk := pattern.(string); pOk {
					patternList = append(patternList, p)
				}
			}
		}

		newPatternsObj, err := types.ListValueFrom(ctx, types.StringType, patternList)
		if err.HasError() {
			diags.Append(err...)
		}

		allowedUrlPatterns.Patterns = newPatternsObj
	}

	// BlockedUrlPatterns
	blockedUrlPatterns := UrlPatterns{}
	if bup, bupOk := contentFiltering["blockedUrlPatterns"].(map[string]interface{}); bupOk {

		// Settings
		settings := types.StringNull()
		if s, ok := bup["settings"].(string); ok {
			settings = types.StringValue(s)
		}
		blockedUrlPatterns.Settings = settings

		// Patterns
		var patternList []string
		if patterns, patternsOk := bup["patterns"].([]interface{}); patternsOk {

			for _, pattern := range patterns {
				if p, pOk := pattern.(string); pOk {
					patternList = append(patternList, p)
				}
			}
		}

		newPatternsObj, err := types.ListValueFrom(ctx, types.StringType, patternList)
		if err.HasError() {
			diags.Append(err...)
		}

		blockedUrlPatterns.Patterns = newPatternsObj
	}

	// BlockedUrlCategories
	blockedUrlCategories := UrlCategories{}
	if buc, bucOk := contentFiltering["blockedUrlCategories"].(map[string]interface{}); bucOk {

		// Settings
		settings := types.StringNull()
		if s, ok := buc["settings"].(string); ok {
			settings = types.StringValue(s)
		}
		blockedUrlCategories.Settings = settings

		// Patterns
		var categoriesList []string
		if categories, categoriesOk := buc["categories"].([]interface{}); categoriesOk {

			for _, category := range categories {
				if c, cOk := category.(string); cOk {
					categoriesList = append(categoriesList, c)
				}
			}
		}

		newCategoriesObj, err := types.ListValueFrom(ctx, types.StringType, categoriesList)
		if err.HasError() {
			diags.Append(err...)
		}

		blockedUrlCategories.Categories = newCategoriesObj
	}

	// Content Filtering Object
	contentFilteringObj := ContentFilteringModel{
		AllowedUrlPatterns:   allowedUrlPatterns,
		BlockedUrlPatterns:   blockedUrlPatterns,
		BlockedUrlCategories: blockedUrlCategories,
	}

	newContentFilteringObject, contentFilteringObjErr := types.ObjectValueFrom(ctx,
		map[string]attr.Type{
			"allowed_url_patterns": types.ObjectType{AttrTypes: map[string]attr.Type{
				"patterns": types.ListType{ElemType: types.StringType},
				"settings": types.StringType,
			}},
			"blocked_url_patterns": types.ObjectType{AttrTypes: map[string]attr.Type{
				"patterns": types.ListType{ElemType: types.StringType},
				"settings": types.StringType,
			}},
			"blocked_url_categories": types.ObjectType{AttrTypes: map[string]attr.Type{
				"categories": types.ListType{ElemType: types.StringType},
				"settings":   types.StringType,
			}},
		},
		contentFilteringObj,
	)
	if contentFilteringObjErr.HasError() {
		return contentFilteringObjectNull, diags
	}

	return newContentFilteringObject, diags

}
func updateVlanTaggingState(ctx context.Context, vlanTagging map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	vlanTaggingObjectNull := types.ObjectNull(map[string]attr.Type{
		"settings": types.StringType,
		"vlan_id":  types.StringType,
	})

	if vlanTagging != nil {

		settings := types.StringValue("network default")
		vlanID := types.StringNull()

		if set, setOk := vlanTagging["settings"].(string); setOk {
			settings = types.StringValue(set)
		}

		if vlan, vlanOk := vlanTagging["vlanId"].(string); vlanOk {
			vlanID = types.StringValue(vlan)
		}

		vlanTaggingObjectData := VlanTaggingModel{VlanID: vlanID, Settings: settings}
		vlanTaggingObject, err := types.ObjectValueFrom(ctx, map[string]attr.Type{
			"settings": types.StringType,
			"vlan_id":  types.StringType,
		}, vlanTaggingObjectData)
		if err.HasError() {
			diags.Append(err...)
			return vlanTaggingObjectNull, diags
		}

		return vlanTaggingObject, nil
	}

	return vlanTaggingObjectNull, diags
}
func updateContentFilteringStateHelperExtractPatternsFromList(patterns types.List) []string {
	var result []string
	for _, pattern := range patterns.Elements() {
		result = append(result, pattern.(types.String).ValueString())
	}
	return result
}
func updateContentFilteringStateHelperValidateSettings(settings *string) *string {
	validSettings := []string{"network default", "append", "override"}
	for _, valid := range validSettings {
		if settings != nil && *settings == valid {
			return settings
		}
	}
	defaultSetting := "network default"
	return &defaultSetting
}

func updateGroupPolicyResourcePayload(data *GroupPolicyResourceModel) (client.CreateNetworkGroupPolicyRequest, diag.Diagnostics) {
	var diags diag.Diagnostics
	groupPolicy := client.CreateNetworkGroupPolicyRequest{
		Name:               data.Name.ValueString(),
		SplashAuthSettings: data.SplashAuthSettings.ValueStringPointer(),
	}

	if !data.Scheduling.IsNull() && !data.Scheduling.IsUnknown() {

		scheduling, err := updateScheduleDayModelPayload(data.Scheduling)
		if err.HasError() {
			diags.Append(err...)
		}

		groupPolicy.Scheduling = scheduling

	}

	if !data.Bandwidth.IsNull() && !data.Bandwidth.IsUnknown() {
		bandwidthAttrs := data.Bandwidth.Attributes()
		settings := bandwidthAttrs["settings"].(types.String)

		bandwidthLimitsObj := bandwidthAttrs["bandwidth_limits"].(basetypes.ObjectValue)
		bandwidthLimitsAttrs := bandwidthLimitsObj.Attributes()
		limitUp := bandwidthLimitsAttrs["limit_up"].(types.Int64)
		limitDown := bandwidthLimitsAttrs["limit_down"].(types.Int64)

		limitUpInt := int32(limitUp.ValueInt64())
		limitDownInt := int32(limitDown.ValueInt64())

		groupPolicy.Bandwidth = &client.CreateNetworkGroupPolicyRequestBandwidth{
			Settings: settings.ValueStringPointer(),
			BandwidthLimits: &client.CreateNetworkGroupPolicyRequestBandwidthBandwidthLimits{
				LimitUp:   &limitUpInt,
				LimitDown: &limitDownInt,
			},
		}
	}

	if !data.FirewallAndTrafficShaping.IsNull() && !data.FirewallAndTrafficShaping.IsUnknown() {
		firewallAndTrafficShaping, err := updateFirewallAndTrafficShapingModelPayload(data.FirewallAndTrafficShaping)
		if err.HasError() {
			diags.Append(err...)
		}

		groupPolicy.FirewallAndTrafficShaping = firewallAndTrafficShaping
	}

	if !data.ContentFiltering.IsNull() && !data.ContentFiltering.IsUnknown() {
		contentFilteringAttrs := data.ContentFiltering.Attributes()

		allowedUrlPatternsObj := contentFilteringAttrs["allowed_url_patterns"].(basetypes.ObjectValue)
		allowedUrlPatternsAttrs := allowedUrlPatternsObj.Attributes()
		allowedPatterns := allowedUrlPatternsAttrs["patterns"].(types.List)
		allowedSettings := allowedUrlPatternsAttrs["settings"].(types.String)

		blockedUrlPatternsObj := contentFilteringAttrs["blocked_url_patterns"].(basetypes.ObjectValue)
		blockedUrlPatternsAttrs := blockedUrlPatternsObj.Attributes()
		blockedPatterns := blockedUrlPatternsAttrs["patterns"].(types.List)
		blockedSettings := blockedUrlPatternsAttrs["settings"].(types.String)

		blockedUrlCategoriesObj := contentFilteringAttrs["blocked_url_categories"].(basetypes.ObjectValue)
		blockedUrlCategoriesAttrs := blockedUrlCategoriesObj.Attributes()
		blockedCategories := blockedUrlCategoriesAttrs["categories"].(types.List)
		blockedCategoriesSettings := blockedUrlCategoriesAttrs["settings"].(types.String)

		groupPolicy.ContentFiltering = &client.CreateNetworkGroupPolicyRequestContentFiltering{
			AllowedUrlPatterns: &client.CreateNetworkGroupPolicyRequestContentFilteringAllowedUrlPatterns{
				Settings: updateContentFilteringStateHelperValidateSettings(allowedSettings.ValueStringPointer()),
				Patterns: updateContentFilteringStateHelperExtractPatternsFromList(allowedPatterns),
			},
			BlockedUrlPatterns: &client.CreateNetworkGroupPolicyRequestContentFilteringBlockedUrlPatterns{
				Settings: updateContentFilteringStateHelperValidateSettings(blockedSettings.ValueStringPointer()),
				Patterns: updateContentFilteringStateHelperExtractPatternsFromList(blockedPatterns),
			},
			BlockedUrlCategories: &client.CreateNetworkGroupPolicyRequestContentFilteringBlockedUrlCategories{
				Settings:   updateContentFilteringStateHelperValidateSettings(blockedCategoriesSettings.ValueStringPointer()),
				Categories: updateContentFilteringStateHelperExtractPatternsFromList(blockedCategories),
			},
		}
	} else {
		networkDefault := "network default"
		groupPolicy.ContentFiltering = &client.CreateNetworkGroupPolicyRequestContentFiltering{
			AllowedUrlPatterns: &client.CreateNetworkGroupPolicyRequestContentFilteringAllowedUrlPatterns{
				Settings: &networkDefault,
				Patterns: nil,
			},
			BlockedUrlPatterns: &client.CreateNetworkGroupPolicyRequestContentFilteringBlockedUrlPatterns{
				Settings: &networkDefault,
				Patterns: nil,
			},
			BlockedUrlCategories: &client.CreateNetworkGroupPolicyRequestContentFilteringBlockedUrlCategories{
				Settings:   &networkDefault,
				Categories: nil,
			},
		}
	}

	if !data.VlanTagging.IsNull() && !data.VlanTagging.IsUnknown() {
		vlanTagging, err := updateVlanTaggingModelPayload(data.VlanTagging)
		if err.HasError() {
			diags.Append(err...)
		}

		groupPolicy.VlanTagging = vlanTagging
	}

	if !data.BonjourForwarding.IsNull() && !data.BonjourForwarding.IsUnknown() {
		bonjourForwarding, err := updateBonjourForwardingModelPayload(data.BonjourForwarding)
		if err.HasError() {
			diags.Append(err...)
		}

		groupPolicy.BonjourForwarding = bonjourForwarding
	}

	return groupPolicy, diags

}
func updateScheduleDayModelPayload(data types.Object) (*client.CreateNetworkGroupPolicyRequestScheduling, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Extract attributes from the types.Object
	schedulingAttrs := data.Attributes()
	enabled := schedulingAttrs["enabled"].(types.Bool)

	getScheduleDay := func(dayAttrName string) (interface{}, diag.Diagnostics) {
		var dayPayload interface{}
		dayObj := schedulingAttrs[dayAttrName].(types.Object)
		dayAttrs := dayObj.Attributes()

		if active, ok := dayAttrs["active"].(types.Bool); ok {
			switch dayAttrName {
			case "monday":
				dayPayload = &client.CreateNetworkGroupPolicyRequestSchedulingMonday{Active: active.ValueBoolPointer()}
			case "tuesday":
				dayPayload = &client.CreateNetworkGroupPolicyRequestSchedulingTuesday{Active: active.ValueBoolPointer()}
			case "wednesday":
				dayPayload = &client.CreateNetworkGroupPolicyRequestSchedulingWednesday{Active: active.ValueBoolPointer()}
			case "thursday":
				dayPayload = &client.CreateNetworkGroupPolicyRequestSchedulingThursday{Active: active.ValueBoolPointer()}
			case "friday":
				dayPayload = &client.CreateNetworkGroupPolicyRequestSchedulingFriday{Active: active.ValueBoolPointer()}
			case "saturday":
				dayPayload = &client.CreateNetworkGroupPolicyRequestSchedulingSaturday{Active: active.ValueBoolPointer()}
			case "sunday":
				dayPayload = &client.CreateNetworkGroupPolicyRequestSchedulingSunday{Active: active.ValueBoolPointer()}
			}
		}
		if from, ok := dayAttrs["from"].(types.String); ok {
			switch v := dayPayload.(type) {
			case *client.CreateNetworkGroupPolicyRequestSchedulingMonday:
				v.From = from.ValueStringPointer()
			case *client.CreateNetworkGroupPolicyRequestSchedulingTuesday:
				v.From = from.ValueStringPointer()
			case *client.CreateNetworkGroupPolicyRequestSchedulingWednesday:
				v.From = from.ValueStringPointer()
			case *client.CreateNetworkGroupPolicyRequestSchedulingThursday:
				v.From = from.ValueStringPointer()
			case *client.CreateNetworkGroupPolicyRequestSchedulingFriday:
				v.From = from.ValueStringPointer()
			case *client.CreateNetworkGroupPolicyRequestSchedulingSaturday:
				v.From = from.ValueStringPointer()
			case *client.CreateNetworkGroupPolicyRequestSchedulingSunday:
				v.From = from.ValueStringPointer()
			}
		}

		if to, ok := dayAttrs["to"].(types.String); ok {
			switch v := dayPayload.(type) {
			case *client.CreateNetworkGroupPolicyRequestSchedulingMonday:
				v.To = to.ValueStringPointer()
			case *client.CreateNetworkGroupPolicyRequestSchedulingTuesday:
				v.To = to.ValueStringPointer()
			case *client.CreateNetworkGroupPolicyRequestSchedulingWednesday:
				v.To = to.ValueStringPointer()
			case *client.CreateNetworkGroupPolicyRequestSchedulingThursday:
				v.To = to.ValueStringPointer()
			case *client.CreateNetworkGroupPolicyRequestSchedulingFriday:
				v.To = to.ValueStringPointer()
			case *client.CreateNetworkGroupPolicyRequestSchedulingSaturday:
				v.To = to.ValueStringPointer()
			case *client.CreateNetworkGroupPolicyRequestSchedulingSunday:
				v.To = to.ValueStringPointer()
			}
		}

		return dayPayload, diags
	}

	monday, mondayDiags := getScheduleDay("monday")
	diags = append(diags, mondayDiags...)

	tuesday, tuesdayDiags := getScheduleDay("tuesday")
	diags = append(diags, tuesdayDiags...)

	wednesday, wednesdayDiags := getScheduleDay("wednesday")
	diags = append(diags, wednesdayDiags...)

	thursday, thursdayDiags := getScheduleDay("thursday")
	diags = append(diags, thursdayDiags...)

	friday, fridayDiags := getScheduleDay("friday")
	diags = append(diags, fridayDiags...)

	saturday, saturdayDiags := getScheduleDay("saturday")
	diags = append(diags, saturdayDiags...)

	sunday, sundayDiags := getScheduleDay("sunday")
	diags = append(diags, sundayDiags...)

	payload := &client.CreateNetworkGroupPolicyRequestScheduling{
		Enabled:   enabled.ValueBoolPointer(),
		Monday:    monday.(*client.CreateNetworkGroupPolicyRequestSchedulingMonday),
		Tuesday:   tuesday.(*client.CreateNetworkGroupPolicyRequestSchedulingTuesday),
		Wednesday: wednesday.(*client.CreateNetworkGroupPolicyRequestSchedulingWednesday),
		Thursday:  thursday.(*client.CreateNetworkGroupPolicyRequestSchedulingThursday),
		Friday:    friday.(*client.CreateNetworkGroupPolicyRequestSchedulingFriday),
		Saturday:  saturday.(*client.CreateNetworkGroupPolicyRequestSchedulingSaturday),
		Sunday:    sunday.(*client.CreateNetworkGroupPolicyRequestSchedulingSunday),
	}

	return payload, diags
}
func updateFirewallAndTrafficShapingModelPayload(data types.Object) (*client.CreateNetworkGroupPolicyRequestFirewallAndTrafficShaping, diag.Diagnostics) {
	var diags diag.Diagnostics

	firewallAndTrafficShapingAttrs := data.Attributes()
	settings := firewallAndTrafficShapingAttrs["settings"].(types.String)

	// Helper function to extract L3 firewall rules
	extractL3FirewallRules := func(rulesAttr types.List) []client.CreateNetworkGroupPolicyRequestFirewallAndTrafficShapingL3FirewallRulesInner {
		var rules []client.CreateNetworkGroupPolicyRequestFirewallAndTrafficShapingL3FirewallRulesInner
		for _, rule := range rulesAttr.Elements() {
			ruleAttrs := rule.(types.Object).Attributes()
			rules = append(rules, client.CreateNetworkGroupPolicyRequestFirewallAndTrafficShapingL3FirewallRulesInner{
				Comment:  ruleAttrs["comment"].(types.String).ValueStringPointer(),
				Policy:   ruleAttrs["policy"].(types.String).ValueString(),
				Protocol: ruleAttrs["protocol"].(types.String).ValueString(),
				DestPort: ruleAttrs["dest_port"].(types.String).ValueStringPointer(),
				DestCidr: ruleAttrs["dest_cidr"].(types.String).ValueString(),
			})
		}
		return rules
	}

	// Helper function to extract L7 firewall rules
	extractL7FirewallRules := func(rulesAttr types.List) []client.CreateNetworkGroupPolicyRequestFirewallAndTrafficShapingL7FirewallRulesInner {
		var rules []client.CreateNetworkGroupPolicyRequestFirewallAndTrafficShapingL7FirewallRulesInner
		for _, rule := range rulesAttr.Elements() {
			ruleAttrs := rule.(types.Object).Attributes()
			rules = append(rules, client.CreateNetworkGroupPolicyRequestFirewallAndTrafficShapingL7FirewallRulesInner{
				Policy: ruleAttrs["policy"].(types.String).ValueStringPointer(),
				Type:   ruleAttrs["type"].(types.String).ValueStringPointer(),
				Value:  ruleAttrs["value"].(types.String).ValueStringPointer(),
			})
		}
		return rules
	}

	// Helper function to extract traffic shaping rules
	extractTrafficShapingRules := func(rulesAttr types.List) []client.CreateNetworkGroupPolicyRequestFirewallAndTrafficShapingTrafficShapingRulesInner {
		var rules []client.CreateNetworkGroupPolicyRequestFirewallAndTrafficShapingTrafficShapingRulesInner
		for _, rule := range rulesAttr.Elements() {
			ruleAttrs := rule.(types.Object).Attributes()
			dscpTagValue := int32(ruleAttrs["dscp_tag_value"].(types.Int64).ValueInt64())
			pcpTagValue := int32(ruleAttrs["pcp_tag_value"].(types.Int64).ValueInt64())

			// Extract per-client bandwidth limits
			perClientBandwidthLimitsAttrs := ruleAttrs["per_client_bandwidth_limits"].(types.Object).Attributes()
			bandwidthLimitsAttrs := perClientBandwidthLimitsAttrs["bandwidth_limits"].(types.Object).Attributes()
			limitUp := int32(bandwidthLimitsAttrs["limit_up"].(types.Int64).ValueInt64())
			limitDown := int32(bandwidthLimitsAttrs["limit_down"].(types.Int64).ValueInt64())

			pcbl := client.UpdateNetworkApplianceTrafficShapingRulesRequestRulesInnerPerClientBandwidthLimits{
				Settings:        perClientBandwidthLimitsAttrs["settings"].(types.String).ValueStringPointer(),
				BandwidthLimits: &client.UpdateNetworkApplianceTrafficShapingRulesRequestRulesInnerPerClientBandwidthLimitsBandwidthLimits{LimitUp: &limitUp, LimitDown: &limitDown},
			}

			// Extract definitions
			var definitions []client.UpdateNetworkApplianceTrafficShapingRulesRequestRulesInnerDefinitionsInner
			for _, def := range ruleAttrs["definitions"].(types.List).Elements() {
				defAttrs := def.(types.Object).Attributes()
				definitions = append(definitions, client.UpdateNetworkApplianceTrafficShapingRulesRequestRulesInnerDefinitionsInner{
					Type:  defAttrs["type"].(types.String).ValueString(),
					Value: defAttrs["value"].(types.String).ValueString(),
				})
			}

			rules = append(rules, client.CreateNetworkGroupPolicyRequestFirewallAndTrafficShapingTrafficShapingRulesInner{
				DscpTagValue:             &dscpTagValue,
				PcpTagValue:              &pcpTagValue,
				PerClientBandwidthLimits: &pcbl,
				Definitions:              definitions,
			})
		}
		return rules
	}

	payload := &client.CreateNetworkGroupPolicyRequestFirewallAndTrafficShaping{
		Settings:            settings.ValueStringPointer(),
		L3FirewallRules:     extractL3FirewallRules(firewallAndTrafficShapingAttrs["l3_firewall_rules"].(types.List)),
		L7FirewallRules:     extractL7FirewallRules(firewallAndTrafficShapingAttrs["l7_firewall_rules"].(types.List)),
		TrafficShapingRules: extractTrafficShapingRules(firewallAndTrafficShapingAttrs["traffic_shaping_rules"].(types.List)),
	}

	return payload, diags
}
func updateBonjourForwardingModelPayload(data types.Object) (*client.CreateNetworkGroupPolicyRequestBonjourForwarding, diag.Diagnostics) {
	var diags diag.Diagnostics

	bonjourForwardingAttrs := data.Attributes()
	settings := bonjourForwardingAttrs["settings"].(types.String)
	rules := bonjourForwardingAttrs["rules"].(types.List)

	// Helper function to extract Bonjour forwarding rules
	extractBonjourForwardingRules := func(rulesAttr types.List) []client.CreateNetworkGroupPolicyRequestBonjourForwardingRulesInner {
		var rules []client.CreateNetworkGroupPolicyRequestBonjourForwardingRulesInner
		for _, rule := range rulesAttr.Elements() {
			ruleAttrs := rule.(types.Object).Attributes()
			var services []string
			for _, service := range ruleAttrs["services"].(types.List).Elements() {
				services = append(services, service.(types.String).ValueString())
			}
			rules = append(rules, client.CreateNetworkGroupPolicyRequestBonjourForwardingRulesInner{
				Description: ruleAttrs["description"].(types.String).ValueStringPointer(),
				VlanId:      ruleAttrs["vlan_id"].(types.String).ValueString(),
				Services:    services,
			})
		}
		return rules
	}

	payload := &client.CreateNetworkGroupPolicyRequestBonjourForwarding{
		Settings: settings.ValueStringPointer(),
		Rules:    extractBonjourForwardingRules(rules),
	}

	return payload, diags
}
func updateVlanTaggingModelPayload(data types.Object) (*client.CreateNetworkGroupPolicyRequestVlanTagging, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Extract attributes from the types.Object
	vlanTaggingAttrs := data.Attributes()
	settings := vlanTaggingAttrs["settings"].(types.String)

	var vlanID *string
	if vlanIDAttr, ok := vlanTaggingAttrs["vlan_id"].(types.String); ok && !vlanIDAttr.IsNull() && vlanIDAttr.ValueString() != "" {
		vlanIDString := vlanIDAttr.ValueString()
		_, err := strconv.Atoi(vlanIDString)
		if err == nil {
			vlanID = &vlanIDString
		} else {
			diags.AddError(
				"Error converting VLAN ID",
				fmt.Sprintf("Could not convert VLAN ID '%s' to an integer: %s", vlanIDString, err.Error()),
			)
			return nil, diags
		}
	}

	payload := &client.CreateNetworkGroupPolicyRequestVlanTagging{
		Settings: settings.ValueStringPointer(),
		VlanId:   vlanID,
	}

	return payload, diags
}

// Create handles the creation of the group policy.
func (r *NetworksGroupPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data GroupPolicyResourceModel
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupPolicy, groupPolicyErr := updateGroupPolicyResourcePayload(&data)
	if groupPolicyErr.HasError() {
		resp.Diagnostics.AddError(
			"Error creating group policy payload",
			fmt.Sprintf("unexpected error: %s", groupPolicyErr),
		)
	}

	// The error “Group number has already been taken” suggests that the Meraki API is not handling rapid sequential requests very well.
	// This could be due to rate limiting or some delay in the backend processing.
	maxRetries := r.client.GetConfig().MaximumRetries
	retryDelay := r.client.GetConfig().Retry4xxErrorWaitTime

	createdPolicy, httpResp, err := r.client.NetworksApi.CreateNetworkGroupPolicy(ctx, data.NetworkId.ValueString()).CreateNetworkGroupPolicyRequest(groupPolicy).Execute()
	retries := 0
	for retries < maxRetries && httpResp != nil && httpResp.StatusCode == http.StatusBadRequest {
		fmt.Println(fmt.Sprintf("CREATE Retrying Max: %v, Delay: %v, Attempt:%v", maxRetries, retryDelay, retries))
		fmt.Println(fmt.Sprintf("CREATE Name: %s", data.Name.ValueString()))
		time.Sleep(time.Duration(retryDelay))
		createdPolicy, httpResp, err = r.client.NetworksApi.CreateNetworkGroupPolicy(ctx, data.NetworkId.ValueString()).CreateNetworkGroupPolicyRequest(groupPolicy).Execute()
		retries++
	}

	if err != nil {
		// Extract additional information from the HTTP response
		var responseBody string
		if httpResp != nil && httpResp.Body != nil {
			bodyBytes, readErr := io.ReadAll(httpResp.Body)
			if readErr == nil {
				responseBody = string(bodyBytes)
			}
		}

		resp.Diagnostics.AddError(
			"Error creating group policy",
			fmt.Sprintf("Could not create group policy, unexpected error: %s\nHTTP Response: %v\nResponse Body: %s", err, httpResp, responseBody),
		)
		return
	}

	// Update the state with the new data
	updateGroupPolicyResourceState(ctx, &data, createdPolicy)

	//TODO: TEMP, Check GroupPolicyId
	groupPolicyId, groupPolicyIdOk := createdPolicy["groupPolicyId"].(string)
	if groupPolicyIdOk {
		if groupPolicyId == "" {
			diags.AddError("CREATE, Missing GroupPolicy.", fmt.Sprintf("Name: %s, ID: %s", data.Name.ValueString(), groupPolicyId))
		}
	}

	fmt.Println(fmt.Sprintf("CREATE Name: %s, ID: %s", data.Name.ValueString(), groupPolicyId))

	// TODO: TEMP Check GPO ID
	if data.GroupPolicyId.IsNull() || data.GroupPolicyId.IsUnknown() {
		diags.AddError("CREATE, Received empty GroupPolicy.", fmt.Sprintf("Name: %s, ID: %s", data.Name.ValueString(), data.GroupPolicyId.ValueString()))
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read handles reading the group policy.
func (r *NetworksGroupPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data GroupPolicyResourceModel
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: TEMP Check GPO ID
	if data.GroupPolicyId.IsNull() || data.GroupPolicyId.IsUnknown() {
		diags.AddError("READ, Received empty GroupPolicy.", fmt.Sprintf("Name: %s, ID: %s", data.Name.ValueString(), data.GroupPolicyId.ValueString()))
	}

	fmt.Println(fmt.Sprintf("READ RECIEVED Name: %s, ID: %s", data.Name.ValueString(), data.GroupPolicyId.ValueString()))

	readPolicy, httpResp, err := r.client.NetworksApi.GetNetworkGroupPolicy(ctx, data.NetworkId.ValueString(), data.GroupPolicyId.ValueString()).Execute()
	if err != nil {

		// Extract additional information from the HTTP response
		var responseBody string
		if httpResp != nil && httpResp.Body != nil {
			bodyBytes, readErr := io.ReadAll(httpResp.Body)
			if readErr == nil {
				responseBody = string(bodyBytes)
			}
		}

		resp.Diagnostics.AddError(
			"Error reading group policy",
			fmt.Sprintf("Could not read group policy, unexpected error: %s\nHTTP Response: %v\nResponse Body: %s", err.Error(), httpResp, responseBody),
		)
		return
	}

	// Update the state with the new data
	updateGroupPolicyResourceState(ctx, &data, readPolicy)

	//TODO: TEMP, Check GroupPolicyId
	groupPolicyId, groupPolicyIdOk := readPolicy["groupPolicyId"].(string)
	if groupPolicyIdOk {
		if groupPolicyId == "" {
			diags.AddError("READ, Missing GroupPolicy.", fmt.Sprintf("Name: %s, ID: %s", data.Name.ValueString(), groupPolicyId))
		}
	}
	fmt.Println(fmt.Sprintf("READ Name: %s, ID: %s", data.Name.ValueString(), groupPolicyId))

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update handles updating the group policy.
func (r *NetworksGroupPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data GroupPolicyResourceModel
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: TEMP Check GPO ID
	if data.GroupPolicyId.IsNull() || data.GroupPolicyId.IsUnknown() {
		diags.AddError("UPDATE, Received empty GroupPolicy.", fmt.Sprintf("Name: %s, ID: %s", data.Name.ValueString(), data.GroupPolicyId.ValueString()))
	}

	fmt.Println(fmt.Sprintf("UPDATE Name: %s, ID: %s", data.Name.ValueString(), data.GroupPolicyId.ValueString()))

	groupPolicy, groupPolicyErr := updateGroupPolicyResourcePayload(&data)
	if groupPolicyErr.HasError() {
		resp.Diagnostics.AddError(
			"Error updating group policy payload",
			fmt.Sprintf("unexpected error: %s", groupPolicyErr),
		)
	}

	groupPolicyUpdate := client.UpdateNetworkGroupPolicyRequest{
		Name:                      &groupPolicy.Name,
		Scheduling:                groupPolicy.Scheduling,
		Bandwidth:                 groupPolicy.Bandwidth,
		FirewallAndTrafficShaping: groupPolicy.FirewallAndTrafficShaping,
		ContentFiltering:          groupPolicy.ContentFiltering,
		SplashAuthSettings:        groupPolicy.SplashAuthSettings,
		VlanTagging:               groupPolicy.VlanTagging,
		BonjourForwarding:         groupPolicy.BonjourForwarding,
	}

	updatePolicy, httpResp, err := r.client.NetworksApi.UpdateNetworkGroupPolicy(ctx, data.NetworkId.ValueString(), data.GroupPolicyId.ValueString()).UpdateNetworkGroupPolicyRequest(groupPolicyUpdate).Execute()
	if err != nil {

		// Extract additional information from the HTTP response
		var responseBody string
		if httpResp != nil && httpResp.Body != nil {
			bodyBytes, readErr := io.ReadAll(httpResp.Body)
			if readErr == nil {
				responseBody = string(bodyBytes)
			}
		}

		resp.Diagnostics.AddError(
			"Error updating group policy",
			fmt.Sprintf("Could not update group policy, unexpected error: %s\nHTTP Response: %v\nResponse Body: %s", err.Error(), httpResp, responseBody),
		)

		resp.Diagnostics.AddError(
			"group policy info",
			fmt.Sprintf("NetworkId: %s\nGroupPolicyId: %v\n", data.NetworkId.ValueString(), data.GroupPolicyId.ValueString()),
		)
		return
	}

	// Update the state with the new data
	updateGroupPolicyResourceState(ctx, &data, updatePolicy)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete handles deleting the group policy.
func (r *NetworksGroupPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data GroupPolicyResourceModel
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.client.NetworksApi.DeleteNetworkGroupPolicy(ctx, data.NetworkId.ValueString(), data.GroupPolicyId.ValueString()).Execute()
	if err != nil {

		// Extract additional information from the HTTP response
		var responseBody string
		if httpResp != nil && httpResp.Body != nil {
			bodyBytes, readErr := io.ReadAll(httpResp.Body)
			if readErr == nil {
				responseBody = string(bodyBytes)
			}
		}

		resp.Diagnostics.AddError(
			"Error deleting group policy",
			fmt.Sprintf("Could not delete group policy, unexpected error: %s\nHTTP Response: %v\nResponse Body: %s", err.Error(), httpResp, responseBody),
		)
		return
	}

	resp.State.RemoveResource(ctx)
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
