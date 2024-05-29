package provider

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/core-infra-svcs/terraform-provider-meraki/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/defaults"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
var _ resource.ResourceWithConfigure = &NetworksGroupPolicyResource{}
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

// GroupPolicyResourceModelScheduling represents the scheduling settings.
type GroupPolicyResourceModelScheduling struct {
	Enabled   types.Bool                           `tfsdk:"enabled" json:"enabled"`
	Monday    *GroupPolicyResourceModelScheduleDay `tfsdk:"monday" json:"monday"`
	Tuesday   *GroupPolicyResourceModelScheduleDay `tfsdk:"tuesday" json:"tuesday"`
	Wednesday *GroupPolicyResourceModelScheduleDay `tfsdk:"wednesday" json:"wednesday"`
	Thursday  *GroupPolicyResourceModelScheduleDay `tfsdk:"thursday" json:"thursday"`
	Friday    *GroupPolicyResourceModelScheduleDay `tfsdk:"friday" json:"friday"`
	Saturday  *GroupPolicyResourceModelScheduleDay `tfsdk:"saturday" json:"saturday"`
	Sunday    *GroupPolicyResourceModelScheduleDay `tfsdk:"sunday" json:"sunday"`
}

// GroupPolicyResourceModelScheduleDay represents a single day's schedule.
type GroupPolicyResourceModelScheduleDay struct {
	Active types.Bool   `tfsdk:"active" json:"active"`
	From   types.String `tfsdk:"from" json:"from"`
	To     types.String `tfsdk:"to" json:"to"`
}

// GroupPolicyResourceModelBandwidth represents the bandwidth settings.
type GroupPolicyResourceModelBandwidth struct {
	Settings        types.String                             `tfsdk:"settings" json:"settings"`
	BandwidthLimits *GroupPolicyResourceModelBandwidthLimits `tfsdk:"bandwidth_limits" json:"bandwidthLimits"`
}

// GroupPolicyResourceModelBandwidthLimits represents the bandwidth limits.
type GroupPolicyResourceModelBandwidthLimits struct {
	LimitUp   types.Int64 `tfsdk:"limit_up" json:"limitUp"`
	LimitDown types.Int64 `tfsdk:"limit_down" json:"limitDown"`
}

// GroupPolicyResourceModelFirewallAndTrafficShaping represents the firewall and traffic shaping settings.
type GroupPolicyResourceModelFirewallAndTrafficShaping struct {
	Settings            types.String                                 `tfsdk:"settings" json:"settings"`
	L3FirewallRules     []GroupPolicyResourceModelL3FirewallRule     `tfsdk:"l3_firewall_rules" json:"l3FirewallRules"`
	L7FirewallRules     []GroupPolicyResourceModelL7FirewallRule     `tfsdk:"l7_firewall_rules" json:"l7FirewallRules"`
	TrafficShapingRules []GroupPolicyResourceModelTrafficShapingRule `tfsdk:"traffic_shaping_rules" json:"trafficShapingRules"`
}

// GroupPolicyResourceModelL3FirewallRule represents a layer 3 firewall rule.
type GroupPolicyResourceModelL3FirewallRule struct {
	Comment  types.String `tfsdk:"comment" json:"comment"`
	Policy   types.String `tfsdk:"policy" json:"policy"`
	Protocol types.String `tfsdk:"protocol" json:"protocol"`
	DestPort types.String `tfsdk:"dest_port" json:"destPort"`
	DestCidr types.String `tfsdk:"dest_cidr" json:"destCidr"`
}

// GroupPolicyResourceModelL7FirewallRule represents a layer 7 firewall rule.
type GroupPolicyResourceModelL7FirewallRule struct {
	Policy types.String `tfsdk:"policy" json:"policy"`
	Type   types.String `tfsdk:"type" json:"type"`
	Value  types.String `tfsdk:"value" json:"value"`
}

// GroupPolicyResourceModelTrafficShapingRule represents a traffic shaping rule.2
type GroupPolicyResourceModelTrafficShapingRule struct {
	DscpTagValue             types.Int64                                        `tfsdk:"dscp_tag_value" json:"dscpTagValue"`
	PcpTagValue              types.Int64                                        `tfsdk:"pcp_tag_value" json:"pcpTagValue"`
	PerClientBandwidthLimits *GroupPolicyResourceModelPerClientBandwidthLimits  `tfsdk:"per_client_bandwidth_limits" json:"perClientBandwidthLimits"`
	Definitions              []GroupPolicyResourceModelTrafficShapingDefinition `tfsdk:"definitions" json:"definitions"`
}

// GroupPolicyResourceModelPerClientBandwidthLimits represents the per-client bandwidth limits.
type GroupPolicyResourceModelPerClientBandwidthLimits struct {
	Settings        types.String                             `tfsdk:"settings" json:"settings"`
	BandwidthLimits *GroupPolicyResourceModelBandwidthLimits `tfsdk:"bandwidth_limits" json:"bandwidthLimits"`
}

// GroupPolicyResourceModelTrafficShapingDefinition represents a traffic shaping definition.
type GroupPolicyResourceModelTrafficShapingDefinition struct {
	Type  types.String `tfsdk:"type" json:"type"`
	Value types.String `tfsdk:"value" json:"value"`
}

// GroupPolicyResourceModelContentFiltering represents the content filtering settings.
type GroupPolicyResourceModelContentFiltering struct {
	AllowedUrlPatterns   GroupPolicyResourceModelUrlPatterns   `tfsdk:"allowed_url_patterns" json:"allowedUrlPatterns"`
	BlockedUrlPatterns   GroupPolicyResourceModelUrlPatterns   `tfsdk:"blocked_url_patterns" json:"blockedUrlPatterns"`
	BlockedUrlCategories GroupPolicyResourceModelUrlCategories `tfsdk:"blocked_url_categories" json:"blockedUrlCategories"`
}

type GroupPolicyResourceModelUrlPatterns struct {
	Patterns types.List   `tfsdk:"patterns" json:"patterns"`
	Settings types.String `tfsdk:"settings" json:"settings"`
}

type GroupPolicyResourceModelUrlCategories struct {
	Categories types.List   `tfsdk:"categories" json:"categories"`
	Settings   types.String `tfsdk:"settings" json:"settings"`
}

// GroupPolicyResourceModelVlanTagging represents the VLAN tagging settings.
type GroupPolicyResourceModelVlanTagging struct {
	Settings types.String `tfsdk:"settings" json:"settings"`
	VlanID   types.String `tfsdk:"vlan_id" json:"vlanId"`
}

// GroupPolicyResourceModelBonjourForwarding represents the Bonjour forwarding settings.
type GroupPolicyResourceModelBonjourForwarding struct {
	Settings types.String                                    `tfsdk:"settings" json:"settings"`
	Rules    []GroupPolicyResourceModelBonjourForwardingRule `tfsdk:"rules" json:"rules"`
}

// GroupPolicyResourceModelBonjourForwardingRule represents a Bonjour forwarding rule.
type GroupPolicyResourceModelBonjourForwardingRule struct {
	Description types.String   `tfsdk:"description" json:"description"`
	VlanID      types.String   `tfsdk:"vlan_id" json:"vlanId"`
	Services    []types.String `tfsdk:"services" json:"services"`
}

// ContentFilteringSettingsDefault implements the defaults.String interface
type ContentFilteringSettingsDefault struct{}

func (d *ContentFilteringSettingsDefault) Description(ctx context.Context) string {
	return "Default value for settings"
}

func (d *ContentFilteringSettingsDefault) MarkdownDescription(ctx context.Context) string {
	return "Default value for `settings`"
}

func (d *ContentFilteringSettingsDefault) DefaultString(ctx context.Context, req defaults.StringRequest, resp *defaults.StringResponse) {
	resp.PlanValue = types.StringValue("network default")
}

// Schema defines the schema for the resource.
func (r *NetworksGroupPolicyResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier for the resource, generated by the Meraki API.",
			},
			"group_policy_id": schema.StringAttribute{
				Computed:    true,
				Description: "The unique identifier for the group policy.",
			},
			"network_id": schema.StringAttribute{
				Required:    true,
				Description: "The network ID where the group policy is applied.",
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 31),
				},
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
								Default:  utils.NewStringDefault("network default"),
								Validators: []validator.String{
									stringvalidator.OneOf("network default", "append", "override"),
								},
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
								Default:  utils.NewStringDefault("network default"),
								Validators: []validator.String{
									stringvalidator.OneOf("network default", "append", "override"),
								},
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
								Default:  utils.NewStringDefault("network default"),
								Validators: []validator.String{
									stringvalidator.OneOf("network default", "append", "override"),
								},
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

// updateGroupPolicyResourceState updates the resource state with the provided group policy data.
func updateGroupPolicyResourceState(ctx context.Context, data *GroupPolicyResourceModel, groupPolicy map[string]interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// GroupPolicyId
	data.GroupPolicyId, diags = utils.ExtractStringAttr(groupPolicy, "groupPolicyId")
	if diags.HasError() {
		return diags
	}

	// Safety check for NetworkId
	if data.NetworkId.IsNull() && data.NetworkId.IsUnknown() {
		data.NetworkId, diags = utils.ExtractStringAttr(groupPolicy, "networkId")
		if diags.HasError() {
			return diags
		}
	}

	// Import ID
	if !data.NetworkId.IsNull() && !data.NetworkId.IsUnknown() && !data.GroupPolicyId.IsNull() && !data.GroupPolicyId.IsUnknown() {
		data.ID = types.StringValue(data.NetworkId.ValueString() + "," + data.GroupPolicyId.ValueString())
	} else {
		data.ID = types.StringNull()
	}

	// Name
	data.Name, diags = utils.ExtractStringAttr(groupPolicy, "name")
	if diags.HasError() {
		return diags
	}

	// Update Scheduling
	if scheduling, schedulingOk := groupPolicy["scheduling"].(map[string]interface{}); schedulingOk {
		if scheduling != nil {

			schedulingObj, schedulingDiags := updateGroupPolicyResourceStateScheduling(ctx, scheduling)
			// Add detailed diagnostics for scheduling extraction
			if schedulingDiags.HasError() {
				diags.AddError(
					"Scheduling Extraction Error",
					fmt.Sprintf("Failed to extract scheduling settings: %v", schedulingDiags.Errors()),
				)
				return diags
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
		var bandwidthModel GroupPolicyResourceModelBandwidth

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
			bandwidthModel.BandwidthLimits = &GroupPolicyResourceModelBandwidthLimits{}
		}

		// Convert GroupPolicyResourceModelBandwidth to types.Object
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
	data.SplashAuthSettings, diags = utils.ExtractStringAttr(groupPolicy, "splashAuthSettings")
	if diags.HasError() {
		return diags
	}

	// Update VlanTagging
	if vlanTagging, ok := groupPolicy["vlanTagging"].(map[string]interface{}); ok {
		vlanTaggingObj, vlanTaggingDiags := updateGroupPolicyResourceStateVlanTagging(ctx, vlanTagging)
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

		firewallAndTrafficShapingObject, firewallAndTrafficShapingObjectDiags := updateGroupPolicyResourceStateFirewallAndTrafficShapingRules(ctx, firewallAndTrafficShaping)
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
		contentFilteringObj, cfDiags := updateGroupPolicyResourceStateContentFiltering(ctx, contentFiltering)
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

// updateGroupPolicyResourceStateScheduling updates the resource state with the scheduling data
func updateGroupPolicyResourceStateScheduling(ctx context.Context, scheduling map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var err diag.Diagnostics

	newSchedulingModel := GroupPolicyResourceModelScheduling{}

	// Update scheduling enabled status
	newSchedulingModel.Enabled, err = utils.ExtractBoolAttr(scheduling, "enabled")
	if diags.HasError() {
		diags.Append(err...)
	}

	// Helper function to update individual days
	updateScheduleDay := func(dayAttrName string) (GroupPolicyResourceModelScheduleDay, diag.Diagnostics) {
		var scheduleDayDiags diag.Diagnostics
		var scheduleDay GroupPolicyResourceModelScheduleDay

		if day, ok := scheduling[dayAttrName].(map[string]interface{}); ok {

			scheduleDay.Active, scheduleDayDiags = utils.ExtractBoolAttr(day, "active")
			if scheduleDayDiags.HasError() {
				diags.Append(scheduleDayDiags...)
			}

			scheduleDay.To, scheduleDayDiags = utils.ExtractStringAttr(day, "to")
			if scheduleDayDiags.HasError() {
				diags.Append(scheduleDayDiags...)
			}

			scheduleDay.From, scheduleDayDiags = utils.ExtractStringAttr(day, "from")
			if scheduleDayDiags.HasError() {
				diags.Append(scheduleDayDiags...)
			}
		}

		return scheduleDay, diags
	}

	// Update each day of the week
	days := []string{"monday", "tuesday", "wednesday", "thursday", "friday", "saturday", "sunday"}
	for _, day := range days {
		daySchedule, dayDiags := updateScheduleDay(day)
		if dayDiags.HasError() {
			diags.Append(dayDiags...)
		}
		switch day {
		case "monday":
			newSchedulingModel.Monday = &daySchedule
		case "tuesday":
			newSchedulingModel.Tuesday = &daySchedule
		case "wednesday":
			newSchedulingModel.Wednesday = &daySchedule
		case "thursday":
			newSchedulingModel.Thursday = &daySchedule
		case "friday":
			newSchedulingModel.Friday = &daySchedule
		case "saturday":
			newSchedulingModel.Saturday = &daySchedule
		case "sunday":
			newSchedulingModel.Sunday = &daySchedule
		}
	}

	// Convert new scheduling model to types.Object
	newSchedulingObject, err := types.ObjectValueFrom(ctx, map[string]attr.Type{
		"enabled":   types.BoolType,
		"monday":    types.ObjectType{AttrTypes: map[string]attr.Type{"active": types.BoolType, "from": types.StringType, "to": types.StringType}},
		"tuesday":   types.ObjectType{AttrTypes: map[string]attr.Type{"active": types.BoolType, "from": types.StringType, "to": types.StringType}},
		"wednesday": types.ObjectType{AttrTypes: map[string]attr.Type{"active": types.BoolType, "from": types.StringType, "to": types.StringType}},
		"thursday":  types.ObjectType{AttrTypes: map[string]attr.Type{"active": types.BoolType, "from": types.StringType, "to": types.StringType}},
		"friday":    types.ObjectType{AttrTypes: map[string]attr.Type{"active": types.BoolType, "from": types.StringType, "to": types.StringType}},
		"saturday":  types.ObjectType{AttrTypes: map[string]attr.Type{"active": types.BoolType, "from": types.StringType, "to": types.StringType}},
		"sunday":    types.ObjectType{AttrTypes: map[string]attr.Type{"active": types.BoolType, "from": types.StringType, "to": types.StringType}},
	}, newSchedulingModel)
	if err.HasError() {
		diags.Append(err...)
	}

	return newSchedulingObject, diags
}

// updateGroupPolicyResourceStateFirewallAndTrafficShapingRules updates the resource state with the firewall and traffic shaping rules data
func updateGroupPolicyResourceStateFirewallAndTrafficShapingRules(ctx context.Context, firewallAndTrafficShapingRules map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var err diag.Diagnostics

	firewallAndTrafficShapingObjectNull := types.ObjectNull(map[string]attr.Type{
		"settings":              types.StringType,
		"l3_firewall_rules":     types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"comment": types.StringType, "policy": types.StringType, "protocol": types.StringType, "dest_port": types.StringType, "dest_cidr": types.StringType}}},
		"l7_firewall_rules":     types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"policy": types.StringType, "type": types.StringType, "value": types.StringType}}},
		"traffic_shaping_rules": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"dscp_tag_value": types.Int64Type, "pcp_tag_value": types.Int64Type, "per_client_bandwidth_limits": types.ObjectType{AttrTypes: map[string]attr.Type{"settings": types.StringType, "bandwidth_limits": types.ObjectType{AttrTypes: map[string]attr.Type{"limit_up": types.Int64Type, "limit_down": types.Int64Type}}}}, "definitions": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"type": types.StringType, "value": types.StringType}}}}}},
	})

	if firewallAndTrafficShapingRules == nil {
		return firewallAndTrafficShapingObjectNull, diags
	}

	// Extract settings
	settings, setErr := utils.ExtractStringAttr(firewallAndTrafficShapingRules, "settings")
	if setErr.HasError() {
		diags.Append(setErr...)
	}

	// Extract L3 firewall rules
	var l3FirewallRules []GroupPolicyResourceModelL3FirewallRule
	if l3frs, l3frsOk := firewallAndTrafficShapingRules["l3FirewallRules"].([]interface{}); l3frsOk {
		for _, l3fr := range l3frs {
			if l3, l3Ok := l3fr.(map[string]interface{}); l3Ok {
				rule := GroupPolicyResourceModelL3FirewallRule{
					Comment:  types.StringNull(),
					Policy:   types.StringNull(),
					Protocol: types.StringNull(),
					DestPort: types.StringNull(),
					DestCidr: types.StringNull(),
				}

				// comment
				rule.Comment, err = utils.ExtractStringAttr(l3, "comment")
				if err.HasError() {
					diags.Append(err...)
				}

				// policy
				rule.Policy, err = utils.ExtractStringAttr(l3, "policy")
				if err.HasError() {
					diags.Append(err...)
				}

				// protocol
				rule.Protocol, err = utils.ExtractStringAttr(l3, "protocol")
				if err.HasError() {
					diags.Append(err...)
				}

				// dest port
				rule.DestPort, err = utils.ExtractStringAttr(l3, "destPort")
				if err.HasError() {
					diags.Append(err...)
				}

				// dest cidr
				rule.DestCidr, err = utils.ExtractStringAttr(l3, "destCidr")
				if err.HasError() {
					diags.Append(err...)
				}

				l3FirewallRules = append(l3FirewallRules, rule)
			}
		}
	}

	// Extract L7 firewall rules
	var l7FirewallRules []GroupPolicyResourceModelL7FirewallRule
	if l7frs, l7frsOk := firewallAndTrafficShapingRules["l7FirewallRules"].([]interface{}); l7frsOk {
		for _, l7fr := range l7frs {
			if l7, l7Ok := l7fr.(map[string]interface{}); l7Ok {
				rule := GroupPolicyResourceModelL7FirewallRule{
					Policy: types.StringNull(),
					Type:   types.StringNull(),
					Value:  types.StringNull(),
				}

				// policy
				rule.Policy, err = utils.ExtractStringAttr(l7, "policy")
				if err.HasError() {
					diags.Append(err...)
				}

				// type
				rule.Type, err = utils.ExtractStringAttr(l7, "type")
				if err.HasError() {
					diags.Append(err...)
				}

				// value
				rule.Value, err = utils.ExtractStringAttr(l7, "value")
				if err.HasError() {
					diags.Append(err...)
				}

				l7FirewallRules = append(l7FirewallRules, rule)
			}
		}
	}

	// Extract traffic shaping rules
	var trafficShapingRules []GroupPolicyResourceModelTrafficShapingRule
	if tsRs, tsRsOk := firewallAndTrafficShapingRules["trafficShapingRules"].([]interface{}); tsRsOk {
		for _, tsr := range tsRs {
			if sr, srOk := tsr.(map[string]interface{}); srOk {
				var trafficShapingRule GroupPolicyResourceModelTrafficShapingRule

				// dscp tag value
				dscpTagValue, dscpTagValueErr := utils.ExtractFloat64Attr(sr, "dscpTagValue")
				if dscpTagValueErr.HasError() {
					diags.Append(dscpTagValueErr...)
				}
				trafficShapingRule.DscpTagValue = types.Int64Value(int64(dscpTagValue.ValueFloat64()))

				// pcp tag value
				pcpTagValue, pcpTagValueErr := utils.ExtractFloat64Attr(sr, "pcpTagValue")
				if pcpTagValueErr.HasError() {
					diags.Append(pcpTagValueErr...)
				}
				trafficShapingRule.PcpTagValue = types.Int64Value(int64(pcpTagValue.ValueFloat64()))

				// perClientBandwidthLimits
				if pcBl, pcBlOk := sr["perClientBandwidthLimits"].(map[string]interface{}); pcBlOk {
					perClientBandwidthLimits := GroupPolicyResourceModelPerClientBandwidthLimits{
						BandwidthLimits: &GroupPolicyResourceModelBandwidthLimits{},
					}

					// settings
					perClientBandwidthLimits.Settings, err = utils.ExtractStringAttr(pcBl, "settings")
					if err.HasError() {
						diags.Append(err...)
					}

					// bandwidth limits
					if bandwidthLimits, bandwidthLimitsOk := pcBl["bandwidthLimits"].(map[string]interface{}); bandwidthLimitsOk {

						// limit up
						limitUp, limitUpErr := utils.ExtractFloat64Attr(bandwidthLimits, "limitUp")
						if limitUpErr.HasError() {
							diags.Append(limitUpErr...)
						}
						perClientBandwidthLimits.BandwidthLimits.LimitUp = types.Int64Value(int64(limitUp.ValueFloat64()))

						// limit down
						limitDown, limitDownErr := utils.ExtractFloat64Attr(bandwidthLimits, "limitDown")
						if limitDownErr.HasError() {
							diags.Append(limitDownErr...)
						}
						perClientBandwidthLimits.BandwidthLimits.LimitDown = types.Int64Value(int64(limitDown.ValueFloat64()))

					}

					trafficShapingRule.PerClientBandwidthLimits = &perClientBandwidthLimits
				}

				// definitions
				if defs, defsOk := sr["definitions"].([]interface{}); defsOk {
					var definitionsModel []GroupPolicyResourceModelTrafficShapingDefinition

					for _, definitions := range defs {
						if def, defOk := definitions.(map[string]interface{}); defOk {
							definitionModel := GroupPolicyResourceModelTrafficShapingDefinition{}

							// type
							definitionModel.Type, err = utils.ExtractStringAttr(def, "type")
							if err.HasError() {
								diags.Append(err...)
							}

							// value
							definitionModel.Value, err = utils.ExtractStringAttr(def, "value")
							if err.HasError() {
								diags.Append(err...)
							}

							definitionsModel = append(definitionsModel, definitionModel)
						}
					}

					trafficShapingRule.Definitions = definitionsModel
				}

				trafficShapingRules = append(trafficShapingRules, trafficShapingRule)
			}
		}
	}

	firewallAndTrafficShapingData := GroupPolicyResourceModelFirewallAndTrafficShaping{
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

// updateGroupPolicyResourceStateVlanTagging updates the resource state with the vlan tagging data
func updateGroupPolicyResourceStateVlanTagging(ctx context.Context, vlanTagging map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	vlanTaggingObjectNull := types.ObjectNull(map[string]attr.Type{
		"settings": types.StringType,
		"vlan_id":  types.StringType,
	})

	if vlanTagging == nil {
		return vlanTaggingObjectNull, diags
	}

	settings := types.StringValue("network default")
	vlanID := types.StringNull()

	// settings
	if set, setOk := vlanTagging["settings"].(string); setOk {
		settings = types.StringValue(set)
	}

	// vlan id
	if vlan, vlanOk := vlanTagging["vlanId"].(string); vlanOk {
		vlanID = types.StringValue(vlan)
	}

	vlanTaggingObjectData := GroupPolicyResourceModelVlanTagging{
		VlanID:   vlanID,
		Settings: settings,
	}

	vlanTaggingObject, err := types.ObjectValueFrom(ctx, map[string]attr.Type{
		"settings": types.StringType,
		"vlan_id":  types.StringType,
	}, vlanTaggingObjectData)
	if err.HasError() {
		diags.Append(err...)
		return vlanTaggingObjectNull, diags
	}

	return vlanTaggingObject, diags
}

// updateGroupPolicyResourceStateContentFiltering updates the resource state with the content filtering data
func updateGroupPolicyResourceStateContentFiltering(ctx context.Context, contentFiltering map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var err diag.Diagnostics

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

	// url patterns
	extractUrlPatterns := func(patterns map[string]interface{}) GroupPolicyResourceModelUrlPatterns {
		var urlPatterns GroupPolicyResourceModelUrlPatterns

		// settings
		urlPatterns.Settings, err = utils.ExtractStringAttr(patterns, "settings")
		if err.HasError() {
			diags.Append(err...)
		}

		// patterns
		var patternList []string
		if patternsList, ok := patterns["patterns"].([]interface{}); ok {
			for _, pattern := range patternsList {
				if p, pOk := pattern.(string); pOk {
					patternList = append(patternList, p)
				}
			}
		}

		newPatternsObj, newPatternsObjErr := types.ListValueFrom(ctx, types.StringType, patternList)
		if newPatternsObjErr.HasError() {
			diags.Append(newPatternsObjErr...)
		}

		urlPatterns.Patterns = newPatternsObj

		return urlPatterns
	}

	// allowed url patterns
	allowedUrlPatterns := GroupPolicyResourceModelUrlPatterns{}
	if aup, aupOk := contentFiltering["allowedUrlPatterns"].(map[string]interface{}); aupOk {
		allowedUrlPatterns = extractUrlPatterns(aup)
	}

	// blocked url patterns
	blockedUrlPatterns := GroupPolicyResourceModelUrlPatterns{}
	if bup, bupOk := contentFiltering["blockedUrlPatterns"].(map[string]interface{}); bupOk {
		blockedUrlPatterns = extractUrlPatterns(bup)
	}

	// blocked url categories
	blockedUrlCategories := GroupPolicyResourceModelUrlCategories{}
	if buc, bucOk := contentFiltering["blockedUrlCategories"].(map[string]interface{}); bucOk {

		// settings
		blockedUrlCategories.Settings, err = utils.ExtractStringAttr(buc, "settings")
		if err.HasError() {
			diags.Append(err...)
		}

		// categories
		var categoriesList []string
		if categories, categoriesOk := buc["categories"].([]interface{}); categoriesOk {
			for _, category := range categories {
				if c, cOk := category.(string); cOk {
					categoriesList = append(categoriesList, c)
				}
			}
		}

		newCategoriesObj, newCategoriesObjErr := types.ListValueFrom(ctx, types.StringType, categoriesList)
		if newCategoriesObjErr.HasError() {
			diags.Append(newCategoriesObjErr...)
		}

		blockedUrlCategories.Categories = newCategoriesObj
	}

	contentFilteringObj := GroupPolicyResourceModelContentFiltering{
		AllowedUrlPatterns:   allowedUrlPatterns,
		BlockedUrlPatterns:   blockedUrlPatterns,
		BlockedUrlCategories: blockedUrlCategories,
	}

	newContentFilteringObject, contentFilteringObjErr := types.ObjectValueFrom(ctx, map[string]attr.Type{
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
	}, contentFilteringObj)
	if contentFilteringObjErr.HasError() {
		diags.Append(contentFilteringObjErr...)
		return contentFilteringObjectNull, diags
	}

	return newContentFilteringObject, diags
}

// updateGroupPolicyResourceStateHelperValidateGroupPolicyId checks if the GroupPolicyId is valid.
func (r *NetworksGroupPolicyResource) updateGroupPolicyResourceStateHelperValidateGroupPolicyId(ctx context.Context, state *GroupPolicyResourceModel, readResp *resource.ReadResponse, updateResp *resource.UpdateResponse) error {
	if state.GroupPolicyId.IsNull() || state.GroupPolicyId.IsUnknown() {
		tflog.Error(ctx, "Received empty GroupPolicyId", map[string]interface{}{
			"name":          state.Name.ValueString(),
			"groupPolicyId": state.GroupPolicyId.ValueString(),
		})

		if readResp != nil {
			readResp.Diagnostics.AddError("Received empty GroupPolicy.", fmt.Sprintf("Name: %s, ID: %s", state.Name.ValueString(), state.GroupPolicyId.ValueString()))
		}

		if updateResp != nil {
			updateResp.Diagnostics.AddError("Received empty GroupPolicy.", fmt.Sprintf("Name: %s, ID: %s", state.Name.ValueString(), state.GroupPolicyId.ValueString()))
		}

		return fmt.Errorf("invalid GroupPolicyId")
	}
	return nil
}

// updateGroupPolicyResourcePayload creates a network group policy request payload from the given GroupPolicyResourceModel data
// and returns the payload along with any diagnostics.
func updateGroupPolicyResourcePayload(data *GroupPolicyResourceModel) (client.CreateNetworkGroupPolicyRequest, diag.Diagnostics) {
	var diags diag.Diagnostics
	groupPolicy := client.CreateNetworkGroupPolicyRequest{
		Name:               data.Name.ValueString(),
		SplashAuthSettings: data.SplashAuthSettings.ValueStringPointer(),
	}

	// Extract scheduling information if present and update the group policy scheduling.
	if !data.Scheduling.IsNull() && !data.Scheduling.IsUnknown() {

		scheduling, err := updateGroupPolicyResourcePayloadScheduleDay(data.Scheduling)
		if err.HasError() {
			diags.Append(err...)
		}

		groupPolicy.Scheduling = scheduling

	}

	// Extract bandwidth limits if present and update the group policy bandwidth.
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

	// Extract firewall and traffic shaping information if present and update the group policy firewall and traffic shaping.
	if !data.FirewallAndTrafficShaping.IsNull() && !data.FirewallAndTrafficShaping.IsUnknown() {
		firewallAndTrafficShaping, err := updateGroupPolicyResourcePayloadFirewallAndTrafficShaping(data.FirewallAndTrafficShaping)
		if err.HasError() {
			diags.Append(err...)
		}

		groupPolicy.FirewallAndTrafficShaping = firewallAndTrafficShaping
	}

	// Extract content filtering information if present and update the group policy content filtering.
	if !data.ContentFiltering.IsNull() && !data.ContentFiltering.IsUnknown() {
		contentFilteringAttrs := data.ContentFiltering.Attributes()

		// Extract allowed URL patterns from content filtering attributes.
		allowedUrlPatternsObj := contentFilteringAttrs["allowed_url_patterns"].(basetypes.ObjectValue)
		allowedUrlPatternsAttrs := allowedUrlPatternsObj.Attributes()
		allowedPatterns := allowedUrlPatternsAttrs["patterns"].(types.List)
		allowedSettings := allowedUrlPatternsAttrs["settings"].(types.String)
		allowedPatternsList, allowedPatternsListErr := utils.ExtractStringsFromList(allowedPatterns)
		if allowedPatternsListErr.HasError() {
			diags.Append(allowedPatternsListErr...)
		}

		// Extract blocked URL patterns from content filtering attributes.
		blockedUrlPatternsObj := contentFilteringAttrs["blocked_url_patterns"].(basetypes.ObjectValue)
		blockedUrlPatternsAttrs := blockedUrlPatternsObj.Attributes()
		blockedPatterns := blockedUrlPatternsAttrs["patterns"].(types.List)
		blockedSettings := blockedUrlPatternsAttrs["settings"].(types.String)
		blockedPatternsList, blockedPatternsListErr := utils.ExtractStringsFromList(blockedPatterns)
		if blockedPatternsListErr.HasError() {
			diags.Append(blockedPatternsListErr...)
		}

		// Extract blocked URL categories from content filtering attributes.
		blockedUrlCategoriesObj := contentFilteringAttrs["blocked_url_categories"].(basetypes.ObjectValue)
		blockedUrlCategoriesAttrs := blockedUrlCategoriesObj.Attributes()
		blockedCategories := blockedUrlCategoriesAttrs["categories"].(types.List)
		blockedCategoriesSettings := blockedUrlCategoriesAttrs["settings"].(types.String)
		blockedCategoriesList, blockedCategoriesListErr := utils.ExtractStringsFromList(blockedCategories)
		if blockedCategoriesListErr.HasError() {
			diags.Append(blockedCategoriesListErr...)
		}

		// content filtering
		groupPolicy.ContentFiltering = &client.CreateNetworkGroupPolicyRequestContentFiltering{
			AllowedUrlPatterns: &client.CreateNetworkGroupPolicyRequestContentFilteringAllowedUrlPatterns{
				Settings: allowedSettings.ValueStringPointer(),
				Patterns: allowedPatternsList,
			},
			BlockedUrlPatterns: &client.CreateNetworkGroupPolicyRequestContentFilteringBlockedUrlPatterns{
				Settings: blockedSettings.ValueStringPointer(),
				Patterns: blockedPatternsList,
			},
			BlockedUrlCategories: &client.CreateNetworkGroupPolicyRequestContentFilteringBlockedUrlCategories{
				Settings:   blockedCategoriesSettings.ValueStringPointer(),
				Categories: blockedCategoriesList,
			},
		}
	} else {

		// Set default content filtering values if content filtering attributes are not provided.
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

	// Extract VLAN tagging information if present and update the group policy VLAN tagging.
	if !data.VlanTagging.IsNull() && !data.VlanTagging.IsUnknown() {
		vlanTagging, err := updateGroupPolicyResourcePayloadVlanTagging(data.VlanTagging)
		if err.HasError() {
			diags.Append(err...)
		}

		groupPolicy.VlanTagging = vlanTagging
	}

	// Extract Bonjour forwarding information if present and update the group policy Bonjour forwarding.
	if !data.BonjourForwarding.IsNull() && !data.BonjourForwarding.IsUnknown() {
		bonjourForwarding, err := updateGroupPolicyResourcePayloadBonjourForwarding(data.BonjourForwarding)
		if err.HasError() {
			diags.Append(err...)
		}

		groupPolicy.BonjourForwarding = bonjourForwarding
	}

	// Return the constructed group policy request payload and any diagnostics.
	return groupPolicy, diags

}

// updateGroupPolicyResourcePayloadScheduleDay extracts scheduling information from the given data
// and returns a payload for creating network group policy request scheduling along with any diagnostics.
func updateGroupPolicyResourcePayloadScheduleDay(data types.Object) (*client.CreateNetworkGroupPolicyRequestScheduling, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Extract attributes from the types.Object
	schedulingAttrs := data.Attributes()
	enabled := schedulingAttrs["enabled"].(types.Bool)

	// getScheduleDay extracts the scheduling details for a specific day from the provided attribute name
	// and returns the corresponding day payload and any diagnostics.
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

// updateGroupPolicyResourcePayloadFirewallAndTrafficShaping extracts firewall and traffic shaping information from the given data
// and returns a payload for creating network group policy request firewall and traffic shaping along with any diagnostics.
func updateGroupPolicyResourcePayloadFirewallAndTrafficShaping(data types.Object) (*client.CreateNetworkGroupPolicyRequestFirewallAndTrafficShaping, diag.Diagnostics) {
	var diags diag.Diagnostics

	firewallAndTrafficShapingAttrs := data.Attributes()
	settings := firewallAndTrafficShapingAttrs["settings"].(types.String)

	// extractL3FirewallRules extracts L3 firewall rules from the given types.List and returns a slice of L3 firewall rule objects.
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

	// extractL7FirewallRules extracts L7 firewall rules from the given types.List and returns a slice of L7 firewall rule objects.
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

	// extractTrafficShapingRules extracts traffic shaping rules from the given types.List and returns a slice of traffic shaping rule objects.
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

// updateGroupPolicyResourcePayloadBonjourForwarding extracts Bonjour forwarding information from the given data
// and returns a payload for creating network group policy request Bonjour forwarding along with any diagnostics.
func updateGroupPolicyResourcePayloadBonjourForwarding(data types.Object) (*client.CreateNetworkGroupPolicyRequestBonjourForwarding, diag.Diagnostics) {
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

// updateGroupPolicyResourcePayloadVlanTagging extracts VLAN tagging information from the given data
// and returns a payload for creating network group policy request VLAN tagging along with any diagnostics.
func updateGroupPolicyResourcePayloadVlanTagging(data types.Object) (*client.CreateNetworkGroupPolicyRequestVlanTagging, diag.Diagnostics) {
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
	var plan GroupPolicyResourceModel

	tflog.Trace(ctx, "Starting create operation for NetworksGroupPolicyResource")

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Required plan data for NetworksGroupPolicyResource", map[string]interface{}{
		"networkId": plan.NetworkId.ValueString(),
		"name":      plan.Name.ValueString(),
	})

	groupPolicy, groupPolicyDiags := updateGroupPolicyResourcePayload(&plan)
	if groupPolicyDiags.HasError() {
		tflog.Error(ctx, "Failed to create resource payload", map[string]interface{}{
			"error": groupPolicyDiags,
		})
		resp.Diagnostics.AddError(
			"Error creating group policy payload",
			fmt.Sprintf("Unexpected error: %s", groupPolicyDiags),
		)
		return
	}

	maxRetries := r.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(r.client.GetConfig().Retry4xxErrorWaitTime)

	// API call function to be passed to retryOn4xx
	apiCall := func() (map[string]interface{}, *http.Response, error) {
		return r.client.NetworksApi.CreateNetworkGroupPolicy(ctx, plan.NetworkId.ValueString()).CreateNetworkGroupPolicyRequest(groupPolicy).Execute()
	}

	// Use retryOn4xx for the API call as the meraki API backend returns HTTP 400 messages as a result of collision issues with rapid creation of postgres GroupPolicyIds.
	createdPolicy, httpResp, err := tools.RetryOn4xx(ctx, maxRetries, retryDelay, apiCall)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating group policy",
			fmt.Sprintf("Could not create group policy, unexpected error: %s", err),
		)

		if httpResp != nil {
			var responseBody string
			if httpResp != nil && httpResp.Body != nil {
				bodyBytes, readErr := io.ReadAll(httpResp.Body)
				if readErr == nil {
					responseBody = string(bodyBytes)
				}
			}
			tflog.Error(ctx, "Failed to create resource", map[string]interface{}{
				"error":          err.Error(),
				"httpStatusCode": httpResp.StatusCode,
				"responseBody":   responseBody,
			})
			resp.Diagnostics.AddError(
				"Error creating group policy",
				fmt.Sprintf("HTTP Response: %v\nResponse Body: %s", httpResp, responseBody),
			)
		}
		return
	}

	tflog.Info(ctx, "Group policy created successfully", map[string]interface{}{
		"name":          plan.Name.ValueString(),
		"groupPolicyId": createdPolicy["groupPolicyId"],
	})

	diags = updateGroupPolicyResourceState(ctx, &plan, createdPolicy)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	tflog.Trace(ctx, "Completed create operation for NetworksGroupPolicyResource")
}

// Read handles reading the group policy.
func (r *NetworksGroupPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state GroupPolicyResourceModel

	tflog.Trace(ctx, "Starting read operation for NetworksGroupPolicyResource")

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.updateGroupPolicyResourceStateHelperValidateGroupPolicyId(ctx, &state, resp, nil); err != nil {
		return
	}

	tflog.Debug(ctx, "Reading group policy", map[string]interface{}{
		"networkId":     state.NetworkId.ValueString(),
		"groupPolicyId": state.GroupPolicyId.ValueString(),
	})

	readPolicy, httpResp, err := r.client.NetworksApi.GetNetworkGroupPolicy(ctx, state.NetworkId.ValueString(), state.GroupPolicyId.ValueString()).Execute()
	if err != nil {
		var responseBody string
		if httpResp != nil && httpResp.Body != nil {
			bodyBytes, readErr := io.ReadAll(httpResp.Body)
			if readErr == nil {
				responseBody = string(bodyBytes)
			}
		}
		tflog.Error(ctx, "Failed to read resource", map[string]interface{}{
			"error":          err.Error(),
			"httpStatusCode": httpResp.StatusCode,
			"responseBody":   responseBody,
			"networkId":      state.NetworkId.ValueString(),
			"groupPolicyId":  state.GroupPolicyId.ValueString(),
		})
		resp.Diagnostics.AddError(
			"Error reading group policy",
			fmt.Sprintf("Could not read group policy, unexpected error: %s\nHTTP Response: %v\nResponse Body: %s", err, httpResp, responseBody),
		)
		return
	}

	tflog.Info(ctx, "Group policy read successfully", map[string]interface{}{
		"name":          state.Name.ValueString(),
		"groupPolicyId": state.GroupPolicyId.ValueString(),
	})

	// Update the state with the new state
	diags = updateGroupPolicyResourceState(ctx, &state, readPolicy)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "Completed read operation for NetworksGroupPolicyResource")
}

// Update handles updating the group policy.
func (r *NetworksGroupPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan GroupPolicyResourceModel

	tflog.Trace(ctx, "Starting update operation for NetworksGroupPolicyResource")

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check for empty GroupPolicyId
	if err := r.updateGroupPolicyResourceStateHelperValidateGroupPolicyId(ctx, &plan, nil, resp); err != nil {
		return
	}

	tflog.Debug(ctx, "Updating group policy", map[string]interface{}{
		"networkId":     plan.NetworkId.ValueString(),
		"groupPolicyId": plan.GroupPolicyId.ValueString(),
	})

	groupPolicy, groupPolicyErr := updateGroupPolicyResourcePayload(&plan)
	if groupPolicyErr.HasError() {
		tflog.Error(ctx, "Failed to create update resource payload", map[string]interface{}{
			"error": groupPolicyErr,
		})
		resp.Diagnostics.AddError(
			"Error updating group policy payload",
			fmt.Sprintf("Unexpected error: %s", groupPolicyErr),
		)
		return
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

	updatePolicy, httpResp, err := r.client.NetworksApi.UpdateNetworkGroupPolicy(ctx, plan.NetworkId.ValueString(), plan.GroupPolicyId.ValueString()).UpdateNetworkGroupPolicyRequest(groupPolicyUpdate).Execute()
	if err != nil {
		var responseBody string
		if httpResp != nil && httpResp.Body != nil {
			bodyBytes, readErr := io.ReadAll(httpResp.Body)
			if readErr == nil {
				responseBody = string(bodyBytes)
			}
		}
		tflog.Error(ctx, "Failed to update resource", map[string]interface{}{
			"error":          err.Error(),
			"httpStatusCode": httpResp.StatusCode,
			"responseBody":   responseBody,
			"networkId":      plan.NetworkId.ValueString(),
			"groupPolicyId":  plan.GroupPolicyId.ValueString(),
		})
		resp.Diagnostics.AddError(
			"Error updating group policy",
			fmt.Sprintf("Could not update group policy, unexpected error: %s\nHTTP Response: %v\nResponse Body: %s", err, httpResp, responseBody),
		)
		return
	}

	tflog.Info(ctx, "Group policy updated successfully", map[string]interface{}{
		"name":          plan.Name.ValueString(),
		"groupPolicyId": plan.GroupPolicyId.ValueString(),
	})

	// Update the state with the new plan
	diags = updateGroupPolicyResourceState(ctx, &plan, updatePolicy)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "Completed update operation for NetworksGroupPolicyResource")
}

// Delete handles deleting the group policy.
func (r *NetworksGroupPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state GroupPolicyResourceModel

	tflog.Trace(ctx, "Starting delete operation for NetworksGroupPolicyResource")

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check for empty GroupPolicyId
	if state.GroupPolicyId.IsNull() || state.GroupPolicyId.IsUnknown() {
		tflog.Error(ctx, "Received empty GroupPolicyId", map[string]interface{}{
			"name":          state.Name.ValueString(),
			"groupPolicyId": state.GroupPolicyId.ValueString(),
		})
		resp.Diagnostics.AddError("DELETE, Received empty GroupPolicy.", fmt.Sprintf("Name: %s, ID: %s", state.Name.ValueString(), state.GroupPolicyId.ValueString()))
		return
	}

	tflog.Debug(ctx, "Deleting group policy", map[string]interface{}{
		"networkId":     state.NetworkId.ValueString(),
		"groupPolicyId": state.GroupPolicyId.ValueString(),
	})

	httpResp, err := r.client.NetworksApi.DeleteNetworkGroupPolicy(ctx, state.NetworkId.ValueString(), state.GroupPolicyId.ValueString()).Execute()
	if err != nil {
		var responseBody string
		if httpResp != nil && httpResp.Body != nil {
			bodyBytes, readErr := io.ReadAll(httpResp.Body)
			if readErr == nil {
				responseBody = string(bodyBytes)
			}
		}
		tflog.Error(ctx, "Failed to delete resource", map[string]interface{}{
			"error":          err.Error(),
			"httpStatusCode": httpResp.StatusCode,
			"responseBody":   responseBody,
			"networkId":      state.NetworkId.ValueString(),
			"groupPolicyId":  state.GroupPolicyId.ValueString(),
		})
		resp.Diagnostics.AddError(
			"Error deleting group policy",
			fmt.Sprintf("Could not delete group policy, unexpected error: %s\nHTTP Response: %v\nResponse Body: %s", err, httpResp, responseBody),
		)
		return
	}

	tflog.Info(ctx, "Group policy deleted successfully", map[string]interface{}{
		"name":          state.Name.ValueString(),
		"groupPolicyId": state.GroupPolicyId.ValueString(),
	})

	resp.State.RemoveResource(ctx)

	tflog.Trace(ctx, "Completed delete operation for NetworksGroupPolicyResource")
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
