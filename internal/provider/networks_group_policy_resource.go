package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/meraki/dashboard-api-go/client"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	retryMax     = 5
	retryWaitMin = 1 * time.Second
	retryWaitMax = 30 * time.Second
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

// GroupPolicyModel represents a group policy.
type GroupPolicyModel struct {
	ID                        types.String                    `tfsdk:"id"`
	Name                      types.String                    `tfsdk:"name"`
	GroupPolicyId             types.String                    `tfsdk:"group_policy_id"`
	NetworkId                 types.String                    `tfsdk:"network_id"`
	Scheduling                *SchedulingModel                `tfsdk:"scheduling"`
	Bandwidth                 *BandwidthModel                 `tfsdk:"bandwidth"`
	FirewallAndTrafficShaping *FirewallAndTrafficShapingModel `tfsdk:"firewall_and_traffic_shaping"`
	ContentFiltering          *ContentFilteringModel          `tfsdk:"content_filtering"`
	SplashAuthSettings        types.String                    `tfsdk:"splash_auth_settings"`
	VlanTagging               *VlanTaggingModel               `tfsdk:"vlan_tagging"`
	BonjourForwarding         *BonjourForwardingModel         `tfsdk:"bonjour_forwarding"`
}

// SchedulingModel represents the scheduling settings.
type SchedulingModel struct {
	Enabled   types.Bool        `tfsdk:"enabled"`
	Monday    *ScheduleDayModel `tfsdk:"monday"`
	Tuesday   *ScheduleDayModel `tfsdk:"tuesday"`
	Wednesday *ScheduleDayModel `tfsdk:"wednesday"`
	Thursday  *ScheduleDayModel `tfsdk:"thursday"`
	Friday    *ScheduleDayModel `tfsdk:"friday"`
	Saturday  *ScheduleDayModel `tfsdk:"saturday"`
	Sunday    *ScheduleDayModel `tfsdk:"sunday"`
}

// ScheduleDayModel represents a single day's schedule.
type ScheduleDayModel struct {
	Active types.Bool   `tfsdk:"active"`
	From   types.String `tfsdk:"from"`
	To     types.String `tfsdk:"to"`
}

// BandwidthModel represents the bandwidth settings.
type BandwidthModel struct {
	Settings        types.String          `tfsdk:"settings"`
	BandwidthLimits *BandwidthLimitsModel `tfsdk:"bandwidth_limits"`
}

// BandwidthLimitsModel represents the bandwidth limits.
type BandwidthLimitsModel struct {
	LimitUp   types.Int64 `tfsdk:"limit_up"`
	LimitDown types.Int64 `tfsdk:"limit_down"`
}

// FirewallAndTrafficShapingModel represents the firewall and traffic shaping settings.
type FirewallAndTrafficShapingModel struct {
	Settings            types.String              `tfsdk:"settings"`
	L3FirewallRules     []L3FirewallRuleModel     `tfsdk:"l3_firewall_rules"`
	L7FirewallRules     []L7FirewallRuleModel     `tfsdk:"l7_firewall_rules"`
	TrafficShapingRules []TrafficShapingRuleModel `tfsdk:"traffic_shaping_rules"`
}

// L3FirewallRuleModel represents a layer 3 firewall rule.
type L3FirewallRuleModel struct {
	Comment  types.String `tfsdk:"comment"`
	Policy   types.String `tfsdk:"policy"`
	Protocol types.String `tfsdk:"protocol"`
	DestPort types.String `tfsdk:"dest_port"`
	DestCidr types.String `tfsdk:"dest_cidr"`
}

// L7FirewallRuleModel represents a layer 7 firewall rule.
type L7FirewallRuleModel struct {
	Policy types.String `tfsdk:"policy"`
	Type   types.String `tfsdk:"type"`
	Value  types.String `tfsdk:"value"`
}

// TrafficShapingRuleModel represents a traffic shaping rule.
type TrafficShapingRuleModel struct {
	DscpTagValue             types.Int64                     `tfsdk:"dscp_tag_value"`
	PcpTagValue              types.Int64                     `tfsdk:"pcp_tag_value"`
	PerClientBandwidthLimits *PerClientBandwidthLimitsModel  `tfsdk:"per_client_bandwidth_limits"`
	Definitions              []TrafficShapingDefinitionModel `tfsdk:"definitions"`
}

// PerClientBandwidthLimitsModel represents the per-client bandwidth limits.
type PerClientBandwidthLimitsModel struct {
	Settings        types.String          `tfsdk:"settings"`
	BandwidthLimits *BandwidthLimitsModel `tfsdk:"bandwidth_limits"`
}

// TrafficShapingDefinitionModel represents a traffic shaping definition.
type TrafficShapingDefinitionModel struct {
	Type  types.String `tfsdk:"type"`
	Value types.String `tfsdk:"value"`
}

// ContentFilteringModel represents the content filtering settings.
type ContentFilteringModel struct {
	AllowedUrlPatterns   UrlPatterns   `tfsdk:"allowed_url_patterns"`
	BlockedUrlPatterns   UrlPatterns   `tfsdk:"blocked_url_patterns"`
	BlockedUrlCategories UrlCategories `tfsdk:"blocked_url_categories"`
}

type UrlPatterns struct {
	Patterns types.List   `tfsdk:"patterns"`
	Settings types.String `tfsdk:"settings"`
}

type UrlCategories struct {
	Categories types.List   `tfsdk:"categories"`
	Settings   types.String `tfsdk:"settings"`
}

// VlanTaggingModel represents the VLAN tagging settings.
type VlanTaggingModel struct {
	Settings types.String `tfsdk:"settings"`
	VlanID   types.String `tfsdk:"vlan_id"`
}

// BonjourForwardingModel represents the Bonjour forwarding settings.
type BonjourForwardingModel struct {
	Settings types.String                 `tfsdk:"settings"`
	Rules    []BonjourForwardingRuleModel `tfsdk:"rules"`
}

// BonjourForwardingRuleModel represents a Bonjour forwarding rule.
type BonjourForwardingRuleModel struct {
	Description types.String   `tfsdk:"description"`
	VlanID      types.String   `tfsdk:"vlan_id"`
	Services    []types.String `tfsdk:"services"`
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

	/*
		// Create a retryable client
			retryableClient := retryablehttp.NewClient()
			retryableClient.RetryMax = retryMax
			retryableClient.RetryWaitMin = retryWaitMin
			retryableClient.RetryWaitMax = retryWaitMax

			// Create a new API client with the retryable HTTP client
			r.client.GetConfig().HTTPClient = retryableClient.StandardClient()
	*/

}

func updateGroupPolicyState(data *GroupPolicyModel, groupPolicy map[string]interface{}) {
	if id, ok := groupPolicy["groupPolicyId"].(string); ok {
		data.GroupPolicyId = types.StringValue(id)
	} else {
		data.GroupPolicyId = types.StringNull()
	}

	data.NetworkId = types.StringValue(data.NetworkId.ValueString())

	data.ID = types.StringValue(data.NetworkId.ValueString() + "," + data.GroupPolicyId.ValueString())

	if name, ok := groupPolicy["name"].(string); ok {
		data.Name = types.StringValue(name)
	} else {
		data.Name = types.StringNull()
	}

	if scheduling, ok := groupPolicy["scheduling"].(map[string]interface{}); ok {
		if data.Scheduling == nil {
			data.Scheduling = &SchedulingModel{}
		}
		if enabled, ok := scheduling["enabled"].(bool); ok {
			data.Scheduling.Enabled = types.BoolValue(enabled)
		} else {
			data.Scheduling.Enabled = types.BoolNull()
		}

		updateScheduleDayModel(&data.Scheduling.Monday, scheduling, "monday")
		updateScheduleDayModel(&data.Scheduling.Tuesday, scheduling, "tuesday")
		updateScheduleDayModel(&data.Scheduling.Wednesday, scheduling, "wednesday")
		updateScheduleDayModel(&data.Scheduling.Thursday, scheduling, "thursday")
		updateScheduleDayModel(&data.Scheduling.Friday, scheduling, "friday")
		updateScheduleDayModel(&data.Scheduling.Saturday, scheduling, "saturday")
		updateScheduleDayModel(&data.Scheduling.Sunday, scheduling, "sunday")
	}

	if bandwidth, ok := groupPolicy["bandwidth"].(map[string]interface{}); ok {
		if data.Bandwidth == nil {
			data.Bandwidth = &BandwidthModel{}
		}
		if settings, ok := bandwidth["settings"].(string); ok {
			data.Bandwidth.Settings = types.StringValue(settings)
		} else {
			data.Bandwidth.Settings = types.StringNull()
		}
		if bandwidthLimits, ok := bandwidth["bandwidthLimits"].(map[string]interface{}); ok {
			if data.Bandwidth.BandwidthLimits == nil {
				data.Bandwidth.BandwidthLimits = &BandwidthLimitsModel{}
			}
			if limitUp, ok := bandwidthLimits["limitUp"].(int64); ok {
				data.Bandwidth.BandwidthLimits.LimitUp = types.Int64Value(limitUp)
			} else {
				data.Bandwidth.BandwidthLimits.LimitUp = types.Int64Null()
			}
			if limitDown, ok := bandwidthLimits["limitDown"].(int64); ok {
				data.Bandwidth.BandwidthLimits.LimitDown = types.Int64Value(limitDown)
			} else {
				data.Bandwidth.BandwidthLimits.LimitDown = types.Int64Null()
			}
		}
	}

	if firewallAndTrafficShaping, ok := groupPolicy["firewallAndTrafficShaping"].(map[string]interface{}); ok {
		if data.FirewallAndTrafficShaping == nil {
			data.FirewallAndTrafficShaping = &FirewallAndTrafficShapingModel{}
		}
		if settings, ok := firewallAndTrafficShaping["settings"].(string); ok {
			data.FirewallAndTrafficShaping.Settings = types.StringValue(settings)
		} else {
			data.FirewallAndTrafficShaping.Settings = types.StringNull()
		}

		updateL3FirewallRules(&data.FirewallAndTrafficShaping.L3FirewallRules, firewallAndTrafficShaping)
		updateL7FirewallRules(&data.FirewallAndTrafficShaping.L7FirewallRules, firewallAndTrafficShaping)
		updateTrafficShapingRules(&data.FirewallAndTrafficShaping.TrafficShapingRules, firewallAndTrafficShaping)
	} else {
		data.FirewallAndTrafficShaping = nil
	}

	if splashAuthSettings, ok := groupPolicy["splashAuthSettings"].(string); ok {
		data.SplashAuthSettings = types.StringValue(splashAuthSettings)
	} else {
		data.SplashAuthSettings = types.StringNull()
	}

	if contentFiltering, ok := groupPolicy["contentFiltering"].(map[string]interface{}); ok {
		if data.ContentFiltering == nil {
			data.ContentFiltering = &ContentFilteringModel{}
		}
		updateContentFilteringModel(data, contentFiltering)
	}

	if vlanTagging, ok := groupPolicy["vlanTagging"].(map[string]interface{}); ok {
		if data.VlanTagging == nil {
			data.VlanTagging = &VlanTaggingModel{}
		}
		if settings, ok := vlanTagging["settings"].(string); ok {
			data.VlanTagging.Settings = types.StringValue(settings)
		} else {
			data.VlanTagging.Settings = types.StringNull()
		}
		if vlanId, ok := vlanTagging["vlanId"].(string); ok {
			data.VlanTagging.VlanID = types.StringValue(vlanId)
		} else {
			data.VlanTagging.VlanID = types.StringNull()
		}
	}

	if bonjourForwarding, ok := groupPolicy["bonjourForwarding"].(map[string]interface{}); ok {
		if data.BonjourForwarding == nil {
			data.BonjourForwarding = &BonjourForwardingModel{}
		}
		if settings, ok := bonjourForwarding["settings"].(string); ok {
			data.BonjourForwarding.Settings = types.StringValue(settings)
		} else {
			data.BonjourForwarding.Settings = types.StringNull()
		}
		updateBonjourForwardingRules(&data.BonjourForwarding.Rules, bonjourForwarding)
	}
}

func updateScheduleDayModel(dayModel **ScheduleDayModel, scheduling map[string]interface{}, day string) {
	if dayData, ok := scheduling[day].(map[string]interface{}); ok {
		if *dayModel == nil {
			*dayModel = &ScheduleDayModel{}
		}
		if active, ok := dayData["active"].(bool); ok {
			(*dayModel).Active = types.BoolValue(active)
		} else {
			(*dayModel).Active = types.BoolNull()
		}
		if from, ok := dayData["from"].(string); ok {
			(*dayModel).From = types.StringValue(from)
		} else {
			(*dayModel).From = types.StringNull()
		}
		if to, ok := dayData["to"].(string); ok {
			(*dayModel).To = types.StringValue(to)
		} else {
			(*dayModel).To = types.StringNull()
		}
	}
}

func updateL3FirewallRules(l3FirewallRules *[]L3FirewallRuleModel, firewallAndTrafficShaping map[string]interface{}) {
	if rules, ok := firewallAndTrafficShaping["l3FirewallRules"].([]interface{}); ok {
		var newL3FirewallRules []L3FirewallRuleModel
		for _, rule := range rules {
			if r, ok := rule.(map[string]interface{}); ok {
				newL3FirewallRules = append(newL3FirewallRules, L3FirewallRuleModel{
					Comment:  types.StringValue(r["comment"].(string)),
					Policy:   types.StringValue(r["policy"].(string)),
					Protocol: types.StringValue(r["protocol"].(string)),
					DestPort: types.StringValue(r["destPort"].(string)),
					DestCidr: types.StringValue(r["destCidr"].(string)),
				})
			}
		}
		*l3FirewallRules = newL3FirewallRules
	}
}

func updateL7FirewallRules(l7FirewallRules *[]L7FirewallRuleModel, firewallAndTrafficShaping map[string]interface{}) {
	if rules, ok := firewallAndTrafficShaping["l7FirewallRules"].([]interface{}); ok {
		var newL7FirewallRules []L7FirewallRuleModel
		for _, rule := range rules {
			if r, ok := rule.(map[string]interface{}); ok {
				newL7FirewallRules = append(newL7FirewallRules, L7FirewallRuleModel{
					Policy: types.StringValue(r["policy"].(string)),
					Type:   types.StringValue(r["type"].(string)),
					Value:  types.StringValue(r["value"].(string)),
				})
			}
		}
		*l7FirewallRules = newL7FirewallRules
	}
}

func updateTrafficShapingRules(trafficShapingRules *[]TrafficShapingRuleModel, firewallAndTrafficShaping map[string]interface{}) {
	// Initialize trafficShapingRules slice
	*trafficShapingRules = []TrafficShapingRuleModel{}

	if rules, ok := firewallAndTrafficShaping["trafficShapingRules"].([]interface{}); ok {
		for _, rule := range rules {
			if r, ok := rule.(map[string]interface{}); ok {
				tsr := TrafficShapingRuleModel{}
				if dscpTagValue, ok := r["dscpTagValue"].(float64); ok {
					tsr.DscpTagValue = types.Int64Value(int64(dscpTagValue))
				} else {
					tsr.DscpTagValue = types.Int64Value(0) // Default to 0 if not set
				}
				if pcpTagValue, ok := r["pcpTagValue"].(float64); ok {
					tsr.PcpTagValue = types.Int64Value(int64(pcpTagValue))
				} else {
					tsr.PcpTagValue = types.Int64Value(0) // Default to 0 if not set
				}
				if definitions, ok := r["definitions"].([]interface{}); ok {
					var defs []TrafficShapingDefinitionModel
					for _, definition := range definitions {
						if d, ok := definition.(map[string]interface{}); ok {
							defs = append(defs, TrafficShapingDefinitionModel{
								Type:  types.StringValue(d["type"].(string)),
								Value: types.StringValue(d["value"].(string)),
							})
						}
					}
					tsr.Definitions = defs
				} else {
					tsr.Definitions = nil
				}

				if perClientBandwidthLimits, ok := r["perClientBandwidthLimits"].(map[string]interface{}); ok {
					pcbl := PerClientBandwidthLimitsModel{}
					if settings, ok := perClientBandwidthLimits["settings"].(string); ok {
						pcbl.Settings = types.StringValue(settings)
					} else {
						pcbl.Settings = types.StringNull()
					}
					if bandwidthLimits, ok := perClientBandwidthLimits["bandwidthLimits"].(map[string]interface{}); ok {
						bl := BandwidthLimitsModel{}
						if limitUp, ok := bandwidthLimits["limitUp"].(float64); ok {
							bl.LimitUp = types.Int64Value(int64(limitUp))
						} else {
							bl.LimitUp = types.Int64Value(0) // Default to 0 if not set
						}
						if limitDown, ok := bandwidthLimits["limitDown"].(float64); ok {
							bl.LimitDown = types.Int64Value(int64(limitDown))
						} else {
							bl.LimitDown = types.Int64Value(0) // Default to 0 if not set
						}
						pcbl.BandwidthLimits = &bl
					}
					tsr.PerClientBandwidthLimits = &pcbl
				} else {
					tsr.PerClientBandwidthLimits = nil
				}
				*trafficShapingRules = append(*trafficShapingRules, tsr)
			}
		}
	}
}

func updateContentFilteringModel(data *GroupPolicyModel, contentFiltering map[string]interface{}) {
	if allowedUrlPatterns, ok := contentFiltering["allowedUrlPatterns"].(map[string]interface{}); ok {
		if settings, ok := allowedUrlPatterns["settings"].(string); ok {
			data.ContentFiltering.AllowedUrlPatterns.Settings = types.StringValue(settings)
		} else {
			data.ContentFiltering.AllowedUrlPatterns.Settings = types.StringNull()
		}

		// Allowed URL Patterns
		if patterns, ok := allowedUrlPatterns["patterns"].([]interface{}); ok {
			var patternList []attr.Value
			for _, pattern := range patterns {
				if p, ok := pattern.(string); ok {
					patternList = append(patternList, types.StringValue(p))
				}
			}
			data.ContentFiltering.AllowedUrlPatterns.Patterns = types.ListValueMust(types.StringType, patternList)
		}
	}

	if blockedUrlPatterns, ok := contentFiltering["blockedUrlPatterns"].(map[string]interface{}); ok {
		if settings, ok := blockedUrlPatterns["settings"].(string); ok {
			data.ContentFiltering.BlockedUrlPatterns.Settings = types.StringValue(settings)
		} else {
			data.ContentFiltering.BlockedUrlPatterns.Settings = types.StringNull()
		}

		// Blocked URL Patterns
		if patterns, ok := blockedUrlPatterns["patterns"].([]interface{}); ok {
			var patternList []attr.Value
			for _, pattern := range patterns {
				if p, ok := pattern.(string); ok {
					patternList = append(patternList, types.StringValue(p))
				}
			}
			data.ContentFiltering.BlockedUrlPatterns.Patterns = types.ListValueMust(types.StringType, patternList)
		}
	}

	if blockedUrlCategories, ok := contentFiltering["blockedUrlCategories"].(map[string]interface{}); ok {
		if settings, ok := blockedUrlCategories["settings"].(string); ok {
			data.ContentFiltering.BlockedUrlCategories.Settings = types.StringValue(settings)
		} else {
			data.ContentFiltering.BlockedUrlCategories.Settings = types.StringNull()
		}

		// Blocked URL Categories
		if categories, ok := blockedUrlCategories["categories"].([]interface{}); ok {
			var categoryList []attr.Value
			for _, category := range categories {
				if c, ok := category.(string); ok {
					categoryList = append(categoryList, types.StringValue(c))
				}
			}
			data.ContentFiltering.BlockedUrlCategories.Categories = types.ListValueMust(types.StringType, categoryList)
		}
	}
}

func updateBonjourForwardingRules(rules *[]BonjourForwardingRuleModel, bonjourForwarding map[string]interface{}) {
	if ruleList, ok := bonjourForwarding["rules"].([]interface{}); ok {
		for _, rule := range ruleList {
			if r, ok := rule.(map[string]interface{}); ok {
				bfRule := BonjourForwardingRuleModel{}
				if description, ok := r["description"].(string); ok {
					bfRule.Description = types.StringValue(description)
				} else {
					bfRule.Description = types.StringNull()
				}
				if vlanId, ok := r["vlanId"].(string); ok {
					bfRule.VlanID = types.StringValue(vlanId)
				} else {
					bfRule.VlanID = types.StringNull()
				}
				if services, ok := r["services"].([]interface{}); ok {
					for _, service := range services {
						if s, ok := service.(string); ok {
							bfRule.Services = append(bfRule.Services, types.StringValue(s))
						}
					}
				}
				*rules = append(*rules, bfRule)
			}
		}
	}
}

// Create handles the creation of the group policy.
func (r *NetworksGroupPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data GroupPolicyModel
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupPolicy := client.CreateNetworkGroupPolicyRequest{
		Name:               data.Name.ValueString(),
		SplashAuthSettings: data.SplashAuthSettings.ValueStringPointer(),
	}

	if data.Scheduling != nil {
		groupPolicy.Scheduling = &client.CreateNetworkGroupPolicyRequestScheduling{
			Enabled: data.Scheduling.Enabled.ValueBoolPointer(),
			Monday: &client.CreateNetworkGroupPolicyRequestSchedulingMonday{
				Active: data.Scheduling.Monday.Active.ValueBoolPointer(),
				From:   data.Scheduling.Monday.From.ValueStringPointer(),
				To:     data.Scheduling.Monday.To.ValueStringPointer(),
			},
			Tuesday: &client.CreateNetworkGroupPolicyRequestSchedulingTuesday{
				Active: data.Scheduling.Tuesday.Active.ValueBoolPointer(),
				From:   data.Scheduling.Tuesday.From.ValueStringPointer(),
				To:     data.Scheduling.Tuesday.To.ValueStringPointer(),
			},
			Wednesday: &client.CreateNetworkGroupPolicyRequestSchedulingWednesday{
				Active: data.Scheduling.Wednesday.Active.ValueBoolPointer(),
				From:   data.Scheduling.Wednesday.From.ValueStringPointer(),
				To:     data.Scheduling.Wednesday.To.ValueStringPointer(),
			},
			Thursday: &client.CreateNetworkGroupPolicyRequestSchedulingThursday{
				Active: data.Scheduling.Thursday.Active.ValueBoolPointer(),
				From:   data.Scheduling.Thursday.From.ValueStringPointer(),
				To:     data.Scheduling.Thursday.To.ValueStringPointer(),
			},
			Friday: &client.CreateNetworkGroupPolicyRequestSchedulingFriday{
				Active: data.Scheduling.Friday.Active.ValueBoolPointer(),
				From:   data.Scheduling.Friday.From.ValueStringPointer(),
				To:     data.Scheduling.Friday.To.ValueStringPointer(),
			},
			Saturday: &client.CreateNetworkGroupPolicyRequestSchedulingSaturday{
				Active: data.Scheduling.Saturday.Active.ValueBoolPointer(),
				From:   data.Scheduling.Saturday.From.ValueStringPointer(),
				To:     data.Scheduling.Saturday.To.ValueStringPointer(),
			},
			Sunday: &client.CreateNetworkGroupPolicyRequestSchedulingSunday{
				Active: data.Scheduling.Saturday.Active.ValueBoolPointer(),
				From:   data.Scheduling.Saturday.From.ValueStringPointer(),
				To:     data.Scheduling.Saturday.To.ValueStringPointer(),
			},
		}
	}

	if data.Bandwidth != nil {
		limitUp := int32(data.Bandwidth.BandwidthLimits.LimitUp.ValueInt64())
		limitDown := int32(data.Bandwidth.BandwidthLimits.LimitDown.ValueInt64())
		groupPolicy.Bandwidth = &client.CreateNetworkGroupPolicyRequestBandwidth{
			Settings: data.Bandwidth.Settings.ValueStringPointer(),
			BandwidthLimits: &client.CreateNetworkGroupPolicyRequestBandwidthBandwidthLimits{
				LimitUp:   &limitUp,
				LimitDown: &limitDown,
			},
		}
	}

	if data.FirewallAndTrafficShaping != nil {
		settings := data.FirewallAndTrafficShaping.Settings.ValueString()
		groupPolicy.FirewallAndTrafficShaping = &client.CreateNetworkGroupPolicyRequestFirewallAndTrafficShaping{
			Settings: &settings,
		}

		for _, rule := range data.FirewallAndTrafficShaping.L3FirewallRules {
			groupPolicy.FirewallAndTrafficShaping.L3FirewallRules = append(
				groupPolicy.FirewallAndTrafficShaping.L3FirewallRules,
				client.CreateNetworkGroupPolicyRequestFirewallAndTrafficShapingL3FirewallRulesInner{
					Comment:  rule.Comment.ValueStringPointer(),
					Policy:   rule.Policy.ValueString(),
					Protocol: rule.Protocol.ValueString(),
					DestPort: rule.DestPort.ValueStringPointer(),
					DestCidr: rule.DestCidr.ValueString(),
				},
			)
		}

		for _, rule := range data.FirewallAndTrafficShaping.L7FirewallRules {
			groupPolicy.FirewallAndTrafficShaping.L7FirewallRules = append(
				groupPolicy.FirewallAndTrafficShaping.L7FirewallRules,
				client.CreateNetworkGroupPolicyRequestFirewallAndTrafficShapingL7FirewallRulesInner{
					Policy: rule.Policy.ValueStringPointer(),
					Type:   rule.Type.ValueStringPointer(),
					Value:  rule.Value.ValueStringPointer(),
				},
			)
		}

		for _, rule := range data.FirewallAndTrafficShaping.TrafficShapingRules {

			limitUp := int32(rule.PerClientBandwidthLimits.BandwidthLimits.LimitUp.ValueInt64())
			limitDown := int32(rule.PerClientBandwidthLimits.BandwidthLimits.LimitDown.ValueInt64())

			bl := client.UpdateNetworkApplianceTrafficShapingRulesRequestRulesInnerPerClientBandwidthLimitsBandwidthLimits{
				LimitUp:   &limitUp,
				LimitDown: &limitDown,
			}

			pcbl := client.UpdateNetworkApplianceTrafficShapingRulesRequestRulesInnerPerClientBandwidthLimits{
				Settings:        rule.PerClientBandwidthLimits.Settings.ValueStringPointer(),
				BandwidthLimits: &bl,
			}

			dscpTagValue := int32(rule.DscpTagValue.ValueInt64())
			pcpTagValue := int32(rule.PcpTagValue.ValueInt64())

			tsr := client.CreateNetworkGroupPolicyRequestFirewallAndTrafficShapingTrafficShapingRulesInner{
				DscpTagValue:             &dscpTagValue,
				PcpTagValue:              &pcpTagValue,
				PerClientBandwidthLimits: &pcbl,
			}

			for _, def := range rule.Definitions {
				tsr.Definitions = append(tsr.Definitions, client.UpdateNetworkApplianceTrafficShapingRulesRequestRulesInnerDefinitionsInner{
					Type:  def.Type.ValueString(),
					Value: def.Value.ValueString(),
				})
			}

			groupPolicy.FirewallAndTrafficShaping.TrafficShapingRules = append(groupPolicy.FirewallAndTrafficShaping.TrafficShapingRules, tsr)
		}
	}

	if data.ContentFiltering != nil {
		groupPolicy.ContentFiltering = &client.CreateNetworkGroupPolicyRequestContentFiltering{}

		// Initialize and set AllowedUrlPatterns
		if data.ContentFiltering.AllowedUrlPatterns.Settings.ValueString() != "" {
			groupPolicy.ContentFiltering.AllowedUrlPatterns = &client.CreateNetworkGroupPolicyRequestContentFilteringAllowedUrlPatterns{
				Settings: data.ContentFiltering.AllowedUrlPatterns.Settings.ValueStringPointer(),
			}

			// Ranging over Allowed URL Patterns
			for _, ap := range data.ContentFiltering.AllowedUrlPatterns.Patterns.Elements() {
				groupPolicy.ContentFiltering.AllowedUrlPatterns.Patterns = append(groupPolicy.ContentFiltering.AllowedUrlPatterns.Patterns, ap.String())
			}
		}

		// Initialize and set BlockedUrlPatterns
		if data.ContentFiltering.BlockedUrlPatterns.Settings.ValueString() != "" {
			groupPolicy.ContentFiltering.BlockedUrlPatterns = &client.CreateNetworkGroupPolicyRequestContentFilteringBlockedUrlPatterns{
				Settings: data.ContentFiltering.BlockedUrlPatterns.Settings.ValueStringPointer(),
			}

			// Ranging over Blocked URL Patterns
			for _, bp := range data.ContentFiltering.BlockedUrlPatterns.Patterns.Elements() {
				groupPolicy.ContentFiltering.BlockedUrlPatterns.Patterns = append(groupPolicy.ContentFiltering.BlockedUrlPatterns.Patterns, bp.String())
			}
		}

		// Initialize and set BlockedUrlCategories
		if data.ContentFiltering.BlockedUrlCategories.Settings.ValueString() != "" {
			groupPolicy.ContentFiltering.BlockedUrlCategories = &client.CreateNetworkGroupPolicyRequestContentFilteringBlockedUrlCategories{
				Settings: data.ContentFiltering.BlockedUrlCategories.Settings.ValueStringPointer(),
			}
			// Ranging over Blocked URL Categories
			for _, bc := range data.ContentFiltering.BlockedUrlCategories.Categories.Elements() {
				groupPolicy.ContentFiltering.BlockedUrlCategories.Categories = append(groupPolicy.ContentFiltering.BlockedUrlCategories.Categories, bc.String())
			}
		}
	}

	if data.VlanTagging != nil {
		var vlanID *string
		if !data.VlanTagging.VlanID.IsNull() && data.VlanTagging.VlanID.ValueString() != "" {
			vlanIDString := data.VlanTagging.VlanID.ValueString()
			_, err := strconv.Atoi(vlanIDString)
			if err == nil {
				vlanID = &vlanIDString
			} else {
				resp.Diagnostics.AddError(
					"Error converting VLAN ID",
					fmt.Sprintf("Could not convert VLAN ID '%s' to an integer: %s", vlanIDString, err.Error()),
				)
				return
			}
		}
		groupPolicy.VlanTagging = &client.CreateNetworkGroupPolicyRequestVlanTagging{
			Settings: data.VlanTagging.Settings.ValueStringPointer(),
			VlanId:   vlanID,
		}
	}

	if data.BonjourForwarding != nil {
		groupPolicy.BonjourForwarding = &client.CreateNetworkGroupPolicyRequestBonjourForwarding{
			Settings: data.BonjourForwarding.Settings.ValueStringPointer(),
		}

		for _, rule := range data.BonjourForwarding.Rules {

			var services []string
			for _, service := range rule.Services {
				services = append(services, service.ValueString())
			}

			groupPolicy.BonjourForwarding.Rules = append(
				groupPolicy.BonjourForwarding.Rules,
				client.CreateNetworkGroupPolicyRequestBonjourForwardingRulesInner{
					Description: rule.Description.ValueStringPointer(),
					VlanId:      rule.VlanID.ValueString(),
					Services:    services,
				},
			)
		}
	}

	createdPolicy, httpResp, err := r.client.NetworksApi.CreateNetworkGroupPolicy(ctx, data.NetworkId.ValueString()).CreateNetworkGroupPolicyRequest(groupPolicy).Execute()
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
			fmt.Sprintf("Could not create group policy, unexpected error: %s\nHTTP Response: %v\nResponse Body: %s", err.Error(), httpResp, responseBody),
		)
		return
	}

	// Update the state with the new data
	updateGroupPolicyState(&data, createdPolicy)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read handles reading the group policy.
func (r *NetworksGroupPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data GroupPolicyModel
	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

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
	updateGroupPolicyState(&data, readPolicy)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update handles updating the group policy.
func (r *NetworksGroupPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data GroupPolicyModel
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupPolicy := client.UpdateNetworkGroupPolicyRequest{
		Name:               data.Name.ValueStringPointer(),
		SplashAuthSettings: data.SplashAuthSettings.ValueStringPointer(),
	}

	if data.Scheduling != nil {
		groupPolicy.Scheduling = &client.CreateNetworkGroupPolicyRequestScheduling{
			Enabled: data.Scheduling.Enabled.ValueBoolPointer(),
			Monday: &client.CreateNetworkGroupPolicyRequestSchedulingMonday{
				Active: data.Scheduling.Monday.Active.ValueBoolPointer(),
				From:   data.Scheduling.Monday.From.ValueStringPointer(),
				To:     data.Scheduling.Monday.To.ValueStringPointer(),
			},
			Tuesday: &client.CreateNetworkGroupPolicyRequestSchedulingTuesday{
				Active: data.Scheduling.Tuesday.Active.ValueBoolPointer(),
				From:   data.Scheduling.Tuesday.From.ValueStringPointer(),
				To:     data.Scheduling.Tuesday.To.ValueStringPointer(),
			},
			Wednesday: &client.CreateNetworkGroupPolicyRequestSchedulingWednesday{
				Active: data.Scheduling.Wednesday.Active.ValueBoolPointer(),
				From:   data.Scheduling.Wednesday.From.ValueStringPointer(),
				To:     data.Scheduling.Wednesday.To.ValueStringPointer(),
			},
			Thursday: &client.CreateNetworkGroupPolicyRequestSchedulingThursday{
				Active: data.Scheduling.Thursday.Active.ValueBoolPointer(),
				From:   data.Scheduling.Thursday.From.ValueStringPointer(),
				To:     data.Scheduling.Thursday.To.ValueStringPointer(),
			},
			Friday: &client.CreateNetworkGroupPolicyRequestSchedulingFriday{
				Active: data.Scheduling.Friday.Active.ValueBoolPointer(),
				From:   data.Scheduling.Friday.From.ValueStringPointer(),
				To:     data.Scheduling.Friday.To.ValueStringPointer(),
			},
			Saturday: &client.CreateNetworkGroupPolicyRequestSchedulingSaturday{
				Active: data.Scheduling.Saturday.Active.ValueBoolPointer(),
				From:   data.Scheduling.Saturday.From.ValueStringPointer(),
				To:     data.Scheduling.Saturday.To.ValueStringPointer(),
			},
			Sunday: &client.CreateNetworkGroupPolicyRequestSchedulingSunday{
				Active: data.Scheduling.Saturday.Active.ValueBoolPointer(),
				From:   data.Scheduling.Saturday.From.ValueStringPointer(),
				To:     data.Scheduling.Saturday.To.ValueStringPointer(),
			},
		}
	}

	if data.Bandwidth != nil {
		limitUp := int32(data.Bandwidth.BandwidthLimits.LimitUp.ValueInt64())
		limitDown := int32(data.Bandwidth.BandwidthLimits.LimitDown.ValueInt64())
		groupPolicy.Bandwidth = &client.CreateNetworkGroupPolicyRequestBandwidth{
			Settings: data.Bandwidth.Settings.ValueStringPointer(),
			BandwidthLimits: &client.CreateNetworkGroupPolicyRequestBandwidthBandwidthLimits{
				LimitUp:   &limitUp,
				LimitDown: &limitDown,
			},
		}
	}

	if data.FirewallAndTrafficShaping != nil {
		settings := data.FirewallAndTrafficShaping.Settings.ValueString()
		groupPolicy.FirewallAndTrafficShaping = &client.CreateNetworkGroupPolicyRequestFirewallAndTrafficShaping{
			Settings: &settings,
		}

		for _, rule := range data.FirewallAndTrafficShaping.L3FirewallRules {
			groupPolicy.FirewallAndTrafficShaping.L3FirewallRules = append(
				groupPolicy.FirewallAndTrafficShaping.L3FirewallRules,
				client.CreateNetworkGroupPolicyRequestFirewallAndTrafficShapingL3FirewallRulesInner{
					Comment:  rule.Comment.ValueStringPointer(),
					Policy:   rule.Policy.ValueString(),
					Protocol: rule.Protocol.ValueString(),
					DestPort: rule.DestPort.ValueStringPointer(),
					DestCidr: rule.DestCidr.ValueString(),
				},
			)
		}

		for _, rule := range data.FirewallAndTrafficShaping.L7FirewallRules {
			groupPolicy.FirewallAndTrafficShaping.L7FirewallRules = append(
				groupPolicy.FirewallAndTrafficShaping.L7FirewallRules,
				client.CreateNetworkGroupPolicyRequestFirewallAndTrafficShapingL7FirewallRulesInner{
					Policy: rule.Policy.ValueStringPointer(),
					Type:   rule.Type.ValueStringPointer(),
					Value:  rule.Value.ValueStringPointer(),
				},
			)
		}

		for _, rule := range data.FirewallAndTrafficShaping.TrafficShapingRules {

			limitUp := int32(rule.PerClientBandwidthLimits.BandwidthLimits.LimitUp.ValueInt64())
			limitDown := int32(rule.PerClientBandwidthLimits.BandwidthLimits.LimitDown.ValueInt64())

			bl := client.UpdateNetworkApplianceTrafficShapingRulesRequestRulesInnerPerClientBandwidthLimitsBandwidthLimits{
				LimitUp:   &limitUp,
				LimitDown: &limitDown,
			}

			pcbl := client.UpdateNetworkApplianceTrafficShapingRulesRequestRulesInnerPerClientBandwidthLimits{
				Settings:        rule.PerClientBandwidthLimits.Settings.ValueStringPointer(),
				BandwidthLimits: &bl,
			}

			dscpTagValue := int32(rule.DscpTagValue.ValueInt64())
			pcpTagValue := int32(rule.PcpTagValue.ValueInt64())

			tsr := client.CreateNetworkGroupPolicyRequestFirewallAndTrafficShapingTrafficShapingRulesInner{
				DscpTagValue:             &dscpTagValue,
				PcpTagValue:              &pcpTagValue,
				PerClientBandwidthLimits: &pcbl,
			}

			for _, def := range rule.Definitions {
				tsr.Definitions = append(tsr.Definitions, client.UpdateNetworkApplianceTrafficShapingRulesRequestRulesInnerDefinitionsInner{
					Type:  def.Type.ValueString(),
					Value: def.Value.ValueString(),
				})
			}

			groupPolicy.FirewallAndTrafficShaping.TrafficShapingRules = append(groupPolicy.FirewallAndTrafficShaping.TrafficShapingRules, tsr)
		}
	}

	if data.ContentFiltering != nil {
		groupPolicy.ContentFiltering = &client.CreateNetworkGroupPolicyRequestContentFiltering{}

		// Initialize and set AllowedUrlPatterns
		if data.ContentFiltering.AllowedUrlPatterns.Settings.ValueString() != "" {
			groupPolicy.ContentFiltering.AllowedUrlPatterns = &client.CreateNetworkGroupPolicyRequestContentFilteringAllowedUrlPatterns{
				Settings: data.ContentFiltering.AllowedUrlPatterns.Settings.ValueStringPointer(),
			}

			// Ranging over Allowed URL Patterns
			for _, ap := range data.ContentFiltering.AllowedUrlPatterns.Patterns.Elements() {
				groupPolicy.ContentFiltering.AllowedUrlPatterns.Patterns = append(groupPolicy.ContentFiltering.AllowedUrlPatterns.Patterns, ap.String())
			}
		}

		// Initialize and set BlockedUrlPatterns
		if data.ContentFiltering.BlockedUrlPatterns.Settings.ValueString() != "" {
			groupPolicy.ContentFiltering.BlockedUrlPatterns = &client.CreateNetworkGroupPolicyRequestContentFilteringBlockedUrlPatterns{
				Settings: data.ContentFiltering.BlockedUrlPatterns.Settings.ValueStringPointer(),
			}

			// Ranging over Blocked URL Patterns
			for _, bp := range data.ContentFiltering.BlockedUrlPatterns.Patterns.Elements() {
				groupPolicy.ContentFiltering.BlockedUrlPatterns.Patterns = append(groupPolicy.ContentFiltering.BlockedUrlPatterns.Patterns, bp.String())
			}
		}

		// Initialize and set BlockedUrlCategories
		if data.ContentFiltering.BlockedUrlCategories.Settings.ValueString() != "" {
			groupPolicy.ContentFiltering.BlockedUrlCategories = &client.CreateNetworkGroupPolicyRequestContentFilteringBlockedUrlCategories{
				Settings: data.ContentFiltering.BlockedUrlCategories.Settings.ValueStringPointer(),
			}

			// Ranging over Blocked URL Categories
			for _, bc := range data.ContentFiltering.BlockedUrlCategories.Categories.Elements() {
				groupPolicy.ContentFiltering.BlockedUrlCategories.Categories = append(groupPolicy.ContentFiltering.BlockedUrlCategories.Categories, bc.String())
			}
		}
	}

	if data.VlanTagging != nil {
		var vlanID *string
		if !data.VlanTagging.VlanID.IsNull() && data.VlanTagging.VlanID.ValueString() != "" {
			vlanIDString := data.VlanTagging.VlanID.ValueString()
			_, err := strconv.Atoi(vlanIDString)
			if err == nil {
				vlanID = &vlanIDString
			} else {
				resp.Diagnostics.AddError(
					"Error converting VLAN ID",
					fmt.Sprintf("Could not convert VLAN ID '%s' to an integer: %s", vlanIDString, err.Error()),
				)
				return
			}
		}
		groupPolicy.VlanTagging = &client.CreateNetworkGroupPolicyRequestVlanTagging{
			Settings: data.VlanTagging.Settings.ValueStringPointer(),
			VlanId:   vlanID,
		}
	}

	if data.BonjourForwarding != nil {
		groupPolicy.BonjourForwarding = &client.CreateNetworkGroupPolicyRequestBonjourForwarding{
			Settings: data.BonjourForwarding.Settings.ValueStringPointer(),
		}

		for _, rule := range data.BonjourForwarding.Rules {

			var services []string
			for _, service := range rule.Services {
				services = append(services, service.ValueString())
			}

			groupPolicy.BonjourForwarding.Rules = append(
				groupPolicy.BonjourForwarding.Rules,
				client.CreateNetworkGroupPolicyRequestBonjourForwardingRulesInner{
					Description: rule.Description.ValueStringPointer(),
					VlanId:      rule.VlanID.ValueString(),
					Services:    services,
				},
			)
		}
	}

	updatePolicy, httpResp, err := r.client.NetworksApi.UpdateNetworkGroupPolicy(ctx, data.NetworkId.ValueString(), data.GroupPolicyId.ValueString()).UpdateNetworkGroupPolicyRequest(groupPolicy).Execute()
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
		return
	}

	// Update the state with the new data
	updateGroupPolicyState(&data, updatePolicy)

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete handles deleting the group policy.
func (r *NetworksGroupPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data GroupPolicyModel
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
