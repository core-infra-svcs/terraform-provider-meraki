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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
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
	Enabled   types.Bool   `tfsdk:"enabled" json:"enabled"`
	Monday    types.Object `tfsdk:"monday" json:"monday"`
	Tuesday   types.Object `tfsdk:"tuesday" json:"tuesday"`
	Wednesday types.Object `tfsdk:"wednesday" json:"wednesday"`
	Thursday  types.Object `tfsdk:"thursday" json:"thursday"`
	Friday    types.Object `tfsdk:"friday" json:"friday"`
	Saturday  types.Object `tfsdk:"saturday" json:"saturday"`
	Sunday    types.Object `tfsdk:"sunday" json:"sunday"`
}

// GroupPolicyResourceModelScheduleDay represents a single day's schedule.
type GroupPolicyResourceModelScheduleDay struct {
	Active types.Bool   `tfsdk:"active" json:"active"`
	From   types.String `tfsdk:"from" json:"from"`
	To     types.String `tfsdk:"to" json:"to"`
}

// GroupPolicyResourceModelBandwidth represents the bandwidth settings.
type GroupPolicyResourceModelBandwidth struct {
	Settings        types.String `tfsdk:"settings" json:"settings"`
	BandwidthLimits types.Object `tfsdk:"bandwidth_limits" json:"bandwidthLimits"`
}

// GroupPolicyResourceModelBandwidthLimits represents the bandwidth limits.
type GroupPolicyResourceModelBandwidthLimits struct {
	LimitUp   types.Int64 `tfsdk:"limit_up" json:"limitUp"`
	LimitDown types.Int64 `tfsdk:"limit_down" json:"limitDown"`
}

// GroupPolicyResourceModelFirewallAndTrafficShaping represents the firewall and traffic shaping settings.
type GroupPolicyResourceModelFirewallAndTrafficShaping struct {
	Settings            types.String `tfsdk:"settings" json:"settings"`
	L3FirewallRules     types.List   `tfsdk:"l3_firewall_rules" json:"l3FirewallRules"`
	L7FirewallRules     types.List   `tfsdk:"l7_firewall_rules" json:"l7FirewallRules"`
	TrafficShapingRules types.List   `tfsdk:"traffic_shaping_rules" json:"trafficShapingRules"`
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
	DscpTagValue             types.Int64  `tfsdk:"dscp_tag_value" json:"dscpTagValue"`
	PcpTagValue              types.Int64  `tfsdk:"pcp_tag_value" json:"pcpTagValue"`
	PerClientBandwidthLimits types.Object `tfsdk:"per_client_bandwidth_limits" json:"perClientBandwidthLimits"`
	Definitions              types.List   `tfsdk:"definitions" json:"definitions"`
}

// GroupPolicyResourceModelPerClientBandwidthLimits represents the per-client bandwidth limits.
type GroupPolicyResourceModelPerClientBandwidthLimits struct {
	Settings        types.String `tfsdk:"settings" json:"settings"`
	BandwidthLimits types.Object `tfsdk:"bandwidth_limits" json:"bandwidthLimits"`
}

// GroupPolicyResourceModelTrafficShapingDefinition represents a traffic shaping definition.
type GroupPolicyResourceModelTrafficShapingDefinition struct {
	Type  types.String `tfsdk:"type" json:"type"`
	Value types.String `tfsdk:"value" json:"value"`
}

// GroupPolicyResourceModelContentFiltering represents the content filtering settings.
type GroupPolicyResourceModelContentFiltering struct {
	AllowedUrlPatterns   types.Object `tfsdk:"allowed_url_patterns" json:"allowedUrlPatterns"`
	BlockedUrlPatterns   types.Object `tfsdk:"blocked_url_patterns" json:"blockedUrlPatterns"`
	BlockedUrlCategories types.Object `tfsdk:"blocked_url_categories" json:"blockedUrlCategories"`
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
	Settings types.String `tfsdk:"settings" json:"settings"`
	Rules    types.List   `tfsdk:"rules" json:"rules"`
}

// GroupPolicyResourceModelBonjourForwardingRule represents a Bonjour forwarding rule.
type GroupPolicyResourceModelBonjourForwardingRule struct {
	Description types.String `tfsdk:"description" json:"description"`
	VlanID      types.String `tfsdk:"vlan_id" json:"vlanId"`
	Services    types.List   `tfsdk:"services" json:"services"`
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
				Description: "The network Id where the group policy is applied.",
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
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"limit_up": schema.Int64Attribute{
								Optional:    true,
								Computed:    true,
								Description: "The upload bandwidth limit. Can be null.",
								PlanModifiers: []planmodifier.Int64{
									int64planmodifier.UseStateForUnknown(),
								},
							},
							"limit_down": schema.Int64Attribute{
								Optional:    true,
								Computed:    true,
								Description: "The download bandwidth limit. Can be null.",
								PlanModifiers: []planmodifier.Int64{
									int64planmodifier.UseStateForUnknown(),
								},
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
													PlanModifiers: []planmodifier.Int64{
														int64planmodifier.UseStateForUnknown(),
													},
												},
												"limit_down": schema.Int64Attribute{
													Optional: true,
													Computed: true,
													PlanModifiers: []planmodifier.Int64{
														int64planmodifier.UseStateForUnknown(),
													},
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
								PlanModifiers: []planmodifier.List{
									listplanmodifier.UseStateForUnknown(),
								},
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
								PlanModifiers: []planmodifier.List{
									listplanmodifier.UseStateForUnknown(),
								},
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
								PlanModifiers: []planmodifier.List{
									listplanmodifier.UseStateForUnknown(),
								},
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

		// Initialize pointers to int32 for limit up and limit down
		var limitUpInt *int32
		var limitDownInt *int32

		// Check and assign limit_up if it exists and is not null
		if limitUpAttr, ok := bandwidthLimitsAttrs["limit_up"]; ok && !limitUpAttr.IsNull() {
			limitUp := limitUpAttr.(types.Int64)
			limitUpIntVal := int32(limitUp.ValueInt64())
			limitUpInt = &limitUpIntVal
		}

		// Check and assign limit_down if it exists and is not null
		if limitDownAttr, ok := bandwidthLimitsAttrs["limit_down"]; ok && !limitDownAttr.IsNull() {
			limitDown := limitDownAttr.(types.Int64)
			limitDownIntVal := int32(limitDown.ValueInt64())
			limitDownInt = &limitDownIntVal
		}

		groupPolicy.Bandwidth = &client.CreateNetworkGroupPolicyRequestBandwidth{
			Settings: settings.ValueStringPointer(),
			BandwidthLimits: &client.CreateNetworkGroupPolicyRequestBandwidthBandwidthLimits{
				LimitUp:   limitUpInt,
				LimitDown: limitDownInt,
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
				"Error converting VLAN Id",
				fmt.Sprintf("Could not convert VLAN Id '%s' to an integer: %s", vlanIDString, err.Error()),
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

// state

// Update Scheduling Day
func updateGroupPolicyResourceStateSchedulingDay(httpResp map[string]interface{}, key string) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var day GroupPolicyResourceModelScheduleDay

	dayAttrs := map[string]attr.Type{"active": types.BoolType, "from": types.StringType, "to": types.StringType}

	d, ok := httpResp[key].(map[string]interface{})
	if ok {

		// active
		active, err := utils.ExtractBoolAttr(d, "active")
		if err.HasError() {
			diags.AddError("active Attr", fmt.Sprintf("%s", err.Errors()))
		}
		day.Active = active

		// from
		from, err := utils.ExtractStringAttr(d, "from")
		if err.HasError() {
			diags.AddError("from Attr", fmt.Sprintf("%s", err.Errors()))
		}
		day.From = from

		// to
		to, err := utils.ExtractStringAttr(d, "to")
		if err.HasError() {
			diags.AddError("to Attr", fmt.Sprintf("%s", err.Errors()))
		}
		day.To = to

	} else {
		dayObjNull := types.ObjectNull(dayAttrs)
		return dayObjNull, diags
	}

	dayObj, err := types.ObjectValueFrom(context.Background(), dayAttrs, day)
	if err.HasError() {
		diags.AddError("day object Attr", fmt.Sprintf("%s", err.Errors()))
	}

	return dayObj, diags
}

// Update Scheduling
func updateGroupPolicyResourceStateScheduling(httpResp map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var scheduling GroupPolicyResourceModelScheduling

	dayAttrs := map[string]attr.Type{"active": types.BoolType, "from": types.StringType, "to": types.StringType}

	schedulingAttrs := map[string]attr.Type{
		"enabled":   types.BoolType,
		"monday":    types.ObjectType{AttrTypes: dayAttrs},
		"tuesday":   types.ObjectType{AttrTypes: dayAttrs},
		"wednesday": types.ObjectType{AttrTypes: dayAttrs},
		"thursday":  types.ObjectType{AttrTypes: dayAttrs},
		"friday":    types.ObjectType{AttrTypes: dayAttrs},
		"saturday":  types.ObjectType{AttrTypes: dayAttrs},
		"sunday":    types.ObjectType{AttrTypes: dayAttrs},
	}

	if schedulingMap, schedulingOk := httpResp["scheduling"].(map[string]interface{}); schedulingOk {

		//  Enabled
		enabled, err := utils.ExtractBoolAttr(schedulingMap, "enabled")
		if err.HasError() {
			diags.AddError("enabled Attr", fmt.Sprintf("%s", err.Errors()))
		}
		scheduling.Enabled = enabled

		//    Friday
		friday, err := updateGroupPolicyResourceStateSchedulingDay(schedulingMap, "friday")
		if err.HasError() {
			diags.AddError("friday Attr", fmt.Sprintf("%s", err.Errors()))
		}
		scheduling.Friday = friday

		//    Monday
		monday, err := updateGroupPolicyResourceStateSchedulingDay(schedulingMap, "monday")
		if err.HasError() {
			diags.AddError("monday Attr", fmt.Sprintf("%s", err.Errors()))
		}
		scheduling.Monday = monday

		//    Saturday
		saturday, err := updateGroupPolicyResourceStateSchedulingDay(schedulingMap, "saturday")
		if err.HasError() {
			diags.AddError("saturday Attr", fmt.Sprintf("%s", err.Errors()))
		}
		scheduling.Saturday = saturday

		//    Sunday
		sunday, err := updateGroupPolicyResourceStateSchedulingDay(schedulingMap, "sunday")
		if err.HasError() {
			diags.AddError("sunday Attr", fmt.Sprintf("%s", err.Errors()))
		}
		scheduling.Sunday = sunday

		//    Thursday
		thursday, err := updateGroupPolicyResourceStateSchedulingDay(schedulingMap, "thursday")
		if err.HasError() {
			diags.AddError("thursday Attr", fmt.Sprintf("%s", err.Errors()))
		}
		scheduling.Thursday = thursday

		//    Tuesday
		tuesday, err := updateGroupPolicyResourceStateSchedulingDay(schedulingMap, "tuesday")
		if err.HasError() {
			diags.AddError("tuesday Attr", fmt.Sprintf("%s", err.Errors()))
		}
		scheduling.Tuesday = tuesday

		//    Wednesday
		wednesday, err := updateGroupPolicyResourceStateSchedulingDay(schedulingMap, "wednesday")
		if err.HasError() {
			diags.AddError("wednesday Attr", fmt.Sprintf("%s", err.Errors()))
		}
		scheduling.Wednesday = wednesday

	} else {
		schedulingObjNull := types.ObjectNull(schedulingAttrs)
		return schedulingObjNull, diags
	}

	schedulingObj, err := types.ObjectValueFrom(context.Background(), schedulingAttrs, scheduling)
	if err.HasError() {
		diags.AddError("scheduling object Attr", fmt.Sprintf("%s", err.Errors()))
	}

	return schedulingObj, diags
}

// Update Bandwidth
func updateGroupPolicyResourceStateBandwidth(httpResp map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var bandwidth GroupPolicyResourceModelBandwidth

	bandwidthLimitsAttrs := map[string]attr.Type{
		"limit_up":   types.Int64Type,
		"limit_down": types.Int64Type,
	}

	bandwidthAttrs := map[string]attr.Type{
		"settings": types.StringType,
		"bandwidth_limits": types.ObjectType{
			AttrTypes: bandwidthLimitsAttrs,
		},
	}

	if bandwidthMap, ok := httpResp["bandwidth"].(map[string]interface{}); ok {

		// settings
		settings, err := utils.ExtractStringAttr(bandwidthMap, "settings")
		if err.HasError() {
			diags.AddError("settingsAttr", fmt.Sprintf("%s", err.Errors()))
		}
		bandwidth.Settings = settings

		// bandwidth limits
		if blMap, ok := httpResp["bandwidthLimits"].(map[string]interface{}); ok {
			var bandwidthLimits GroupPolicyResourceModelBandwidthLimits

			// limit up
			limitUp, err := utils.ExtractInt32Attr(blMap, "limitUp")
			if err.HasError() {
				diags.AddError("limitUp Attr", fmt.Sprintf("%s", err.Errors()))
			}
			bandwidthLimits.LimitUp = limitUp

			// limit down
			limitDown, err := utils.ExtractInt32Attr(blMap, "limitDown")
			if err.HasError() {
				diags.AddError("limitDown Attr", fmt.Sprintf("%s", err.Errors()))
			}
			bandwidthLimits.LimitDown = limitDown

			bandwidthLimitsObj, err := types.ObjectValueFrom(context.Background(), bandwidthLimitsAttrs, bandwidthLimits)
			if err.HasError() {
				diags.AddError("bandwidthLimitsObj Attr", fmt.Sprintf("%s", err.Errors()))
			}

			bandwidth.BandwidthLimits = bandwidthLimitsObj

		} else {
			bandwidthLimitsObjNull := types.ObjectNull(bandwidthLimitsAttrs)
			bandwidth.BandwidthLimits = bandwidthLimitsObjNull
		}

	} else {
		bandwidthObjNull := types.ObjectNull(bandwidthAttrs)
		return bandwidthObjNull, diags
	}

	bandwidthObj, err := types.ObjectValueFrom(context.Background(), bandwidthAttrs, bandwidth)
	if err.HasError() {
		diags.AddError("bandwidthObj Attr", fmt.Sprintf("%s", err.Errors()))
	}

	return bandwidthObj, diags
}

// updateGroupPolicyResourceStateFirewallAndTrafficShapingRules updates the resource state with the firewall and traffic shaping rules data
func updateGroupPolicyResourceStateFirewallAndTrafficShapingRules(ctx context.Context, httpResp map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var firewallAndTrafficShapingRules GroupPolicyResourceModelFirewallAndTrafficShaping

	perClientBandwidthLimitsAttr := map[string]attr.Type{
		"settings": types.StringType,
		"bandwidth_limits": types.ObjectType{AttrTypes: map[string]attr.Type{
			"limit_up":   types.Int64Type,
			"limit_down": types.Int64Type}},
	}

	definitionsAttr := map[string]attr.Type{
		"type":  types.StringType,
		"value": types.StringType}

	trafficShapingAttrs := map[string]attr.Type{
		"dscp_tag_value":              types.Int64Type,
		"pcp_tag_value":               types.Int64Type,
		"per_client_bandwidth_limits": types.ObjectType{AttrTypes: perClientBandwidthLimitsAttr},
		"definitions":                 types.ListType{ElemType: types.ObjectType{AttrTypes: definitionsAttr}}}

	l3FirewallRulesAttr := map[string]attr.Type{
		"comment":   types.StringType,
		"policy":    types.StringType,
		"protocol":  types.StringType,
		"dest_port": types.StringType,
		"dest_cidr": types.StringType,
	}

	l7FirewallRulesAttr := map[string]attr.Type{
		"policy": types.StringType,
		"type":   types.StringType,
		"value":  types.StringType,
	}

	firewallAndTrafficShapingRulesAttrs := map[string]attr.Type{
		"settings":              types.StringType,
		"l3_firewall_rules":     types.ListType{ElemType: types.ObjectType{AttrTypes: l3FirewallRulesAttr}},
		"l7_firewall_rules":     types.ListType{ElemType: types.ObjectType{AttrTypes: l7FirewallRulesAttr}},
		"traffic_shaping_rules": types.ListType{ElemType: types.ObjectType{AttrTypes: trafficShapingAttrs}},
	}

	// trafficShapingRules
	ftsr, ok := httpResp["firewallAndTrafficShaping"].(map[string]interface{})
	if ok {

		// settings
		settings, err := utils.ExtractStringAttr(ftsr, "settings")
		if err.HasError() {
			diags.AddError("settings Attr", fmt.Sprintf("%s", err.Errors()))
		}
		firewallAndTrafficShapingRules.Settings = settings

		// l3FirewallRules
		var l3FirewallRules []types.Object
		if l3frs, l3frsOk := ftsr["l3FirewallRules"].([]interface{}); l3frsOk {
			for _, l3fr := range l3frs {
				if l3, l3Ok := l3fr.(map[string]interface{}); l3Ok {
					var rule GroupPolicyResourceModelL3FirewallRule

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

					ruleObj, err := types.ObjectValueFrom(ctx, l3FirewallRulesAttr, rule)
					if err.HasError() {
						diags.Append(err...)
					}

					l3FirewallRules = append(l3FirewallRules, ruleObj)
				}
			}

			// returns a populated or empty list instead of a null value
			if l3FirewallRules != nil {
				l3FirewallRulesList, l3FirewallRulesListErr := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: l3FirewallRulesAttr}, l3FirewallRules)
				if l3FirewallRulesListErr.HasError() {
					diags.Append(l3FirewallRulesListErr...)
				}

				firewallAndTrafficShapingRules.L3FirewallRules = l3FirewallRulesList
			} else {
				l3FirewallRulesList, err := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: l3FirewallRulesAttr}, []attr.Value{})
				if err.HasError() {
					diags.Append(err...)
				}
				firewallAndTrafficShapingRules.L3FirewallRules = l3FirewallRulesList
			}

		} else {
			l3FirewallRulesNull := types.ListNull(types.ObjectType{AttrTypes: l3FirewallRulesAttr})
			firewallAndTrafficShapingRules.L3FirewallRules = l3FirewallRulesNull
		}

		// l7FirewallRules
		var l7FirewallRules []types.Object

		if l7frs, ok := ftsr["l7FirewallRules"].([]interface{}); ok {

			for _, l7fr := range l7frs {
				if l7, ok := l7fr.(map[string]interface{}); ok {
					var rule GroupPolicyResourceModelL7FirewallRule

					// policy
					policy, err := utils.ExtractStringAttr(l7, "policy")
					if err.HasError() {
						diags.Append(err...)
					}
					rule.Policy = policy

					// type
					t, err := utils.ExtractStringAttr(l7, "type")
					if err.HasError() {
						diags.Append(err...)
					}
					rule.Type = t

					// value
					value, err := utils.ExtractStringAttr(l7, "value")
					if err.HasError() {
						diags.Append(err...)
					}
					rule.Value = value

					ruleObj, err := types.ObjectValueFrom(ctx, l7FirewallRulesAttr, rule)
					if err.HasError() {
						diags.Append(err...)
					}

					l7FirewallRules = append(l7FirewallRules, ruleObj)
				}
			}

			// returns a populated or empty list instead of a null value
			if l7FirewallRules != nil {
				l7FirewallRulesList, l7FirewallRulesListErr := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: l7FirewallRulesAttr}, l7FirewallRules)
				if l7FirewallRulesListErr.HasError() {
					diags.Append(l7FirewallRulesListErr...)
				}
				firewallAndTrafficShapingRules.L7FirewallRules = l7FirewallRulesList
			} else {
				l7FirewallRulesList, err := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: l7FirewallRulesAttr}, []attr.Value{})
				if err.HasError() {
					diags.Append(err...)
				}
				firewallAndTrafficShapingRules.L7FirewallRules = l7FirewallRulesList
			}

		} else {
			l7FirewallRulesNull := types.ListNull(types.ObjectType{AttrTypes: l7FirewallRulesAttr})
			firewallAndTrafficShapingRules.L7FirewallRules = l7FirewallRulesNull
		}

		// trafficShapingRules
		var trafficShapingRules []types.Object
		if tsRs, tsRsOk := ftsr["trafficShapingRules"].([]interface{}); tsRsOk {

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

						perClientBandwidthLimits := GroupPolicyResourceModelPerClientBandwidthLimits{}

						// settings
						if _, settingsOk := pcBl["settings"].(string); settingsOk {

							settingsVal, settingsErr := utils.ExtractStringAttr(pcBl, "settings")
							if settingsErr.HasError() {
								diags.Append(settingsErr...)
							}

							perClientBandwidthLimits.Settings = settingsVal

						}

						// bandwidth limits
						if bandwidthLimits, bandwidthLimitsOk := pcBl["bandwidthLimits"].(map[string]interface{}); bandwidthLimitsOk {

							var BandwidthLimits GroupPolicyResourceModelBandwidthLimits

							// limit up
							limitUp, limitUpErr := utils.ExtractFloat64Attr(bandwidthLimits, "limitUp")
							if limitUpErr.HasError() {
								diags.Append(limitUpErr...)
							}
							BandwidthLimits.LimitUp = types.Int64Value(int64(limitUp.ValueFloat64()))

							// limit down
							limitDown, limitDownErr := utils.ExtractFloat64Attr(bandwidthLimits, "limitDown")
							if limitDownErr.HasError() {
								diags.Append(limitDownErr...)
							}
							BandwidthLimits.LimitDown = types.Int64Value(int64(limitDown.ValueFloat64()))

							BandwidthLimitsObject, BandwidthLimitsObjectErr := types.ObjectValueFrom(ctx, map[string]attr.Type{
								"limit_up":   types.Int64Type,
								"limit_down": types.Int64Type,
							}, BandwidthLimits)

							if BandwidthLimitsObjectErr.HasError() {
								diags.Append(BandwidthLimitsObjectErr...)
							}

							perClientBandwidthLimits.BandwidthLimits = BandwidthLimitsObject

						}

						// create types.Object
						perClientBandwidthLimitsObject, perClientBandwidthLimitsObjectErr := types.ObjectValueFrom(ctx, map[string]attr.Type{
							"settings": types.StringType,
							"bandwidth_limits": types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"limit_up":   types.Int64Type,
									"limit_down": types.Int64Type,
								},
							},
						}, perClientBandwidthLimits)

						if perClientBandwidthLimitsObjectErr.HasError() {
							diags.Append(perClientBandwidthLimitsObjectErr...)
						}

						trafficShapingRule.PerClientBandwidthLimits = perClientBandwidthLimitsObject
					} else {

						perClientBandwidthLimitsObjectNull := types.ObjectNull(map[string]attr.Type{
							"settings": types.StringType,
							"bandwidth_limits": types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"limit_up":   types.Int64Type,
									"limit_down": types.Int64Type,
								},
							},
						})

						trafficShapingRule.PerClientBandwidthLimits = perClientBandwidthLimitsObjectNull
					}

					// definitions
					if defs, defsOk := sr["definitions"].([]interface{}); defsOk {
						var definitionsModel []GroupPolicyResourceModelTrafficShapingDefinition

						for _, definitions := range defs {

							definitionModel := GroupPolicyResourceModelTrafficShapingDefinition{}
							if def, defOk := definitions.(map[string]interface{}); defOk {

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

						if definitionsModel != nil {
							definitionsList, definitionsListErr := types.ListValueFrom(ctx, types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"type":  types.StringType,
									"value": types.StringType,
								},
							}, definitionsModel)
							if definitionsListErr.HasError() {
								diags.Append(definitionsListErr...)
							}
							trafficShapingRule.Definitions = definitionsList
						} else {
							definitionsListNull := types.ListNull(types.ObjectType{
								AttrTypes: map[string]attr.Type{
									"type":  types.StringType,
									"value": types.StringType,
								},
							})

							trafficShapingRule.Definitions = definitionsListNull
						}

					}

					trafficShapingRuleObj, err := types.ObjectValueFrom(ctx, trafficShapingAttrs, trafficShapingRule)
					if err.HasError() {
						diags.Append(err...)
					}

					trafficShapingRules = append(trafficShapingRules, trafficShapingRuleObj)
				}
			}

			// returns a populated or empty list instead of a null value
			if trafficShapingRules != nil {
				trafficShapingRulesList, trafficShapingRulesListErr := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: trafficShapingAttrs}, trafficShapingRules)
				if trafficShapingRulesListErr.HasError() {
					diags.Append(trafficShapingRulesListErr...)
				}
				firewallAndTrafficShapingRules.TrafficShapingRules = trafficShapingRulesList
			} else {
				trafficShapingRulesList, err := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: trafficShapingAttrs}, []attr.Value{})
				if err.HasError() {
					diags.Append(err...)
				}
				firewallAndTrafficShapingRules.TrafficShapingRules = trafficShapingRulesList
			}

		} else {
			trafficShapingRulesListNull := types.ListNull(types.ObjectType{AttrTypes: trafficShapingAttrs})
			firewallAndTrafficShapingRules.TrafficShapingRules = trafficShapingRulesListNull
		}

	} else {
		firewallAndTrafficShapingObjectNull := types.ObjectNull(firewallAndTrafficShapingRulesAttrs)
		return firewallAndTrafficShapingObjectNull, diags
	}

	firewallAndTrafficShapingObj, err := types.ObjectValueFrom(ctx, firewallAndTrafficShapingRulesAttrs, firewallAndTrafficShapingRules)
	if err.HasError() {
		diags.Append(err...)
	}

	return firewallAndTrafficShapingObj, diags
}

// updateGroupPolicyResourceStateVlanTagging updates the resource state with the vlan tagging data
func updateGroupPolicyResourceStateVlanTagging(ctx context.Context, httpResp map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var vlanTagging GroupPolicyResourceModelVlanTagging

	vlanTaggingAttr := map[string]attr.Type{
		"settings": types.StringType,
		"vlan_id":  types.StringType,
	}

	vt, ok := httpResp["vlanTagging"].(map[string]interface{})
	if ok {

		// settings
		settings, err := utils.ExtractStringAttr(vt, "settings")
		if err.HasError() {
			diags.AddError("settings Attr", fmt.Sprintf("%s", err.Errors()))
		}
		vlanTagging.Settings = settings

		// vlan id
		vlanId, err := utils.ExtractStringAttr(vt, "vlanId")
		if err.HasError() {
			diags.AddError("vlanId Attr", fmt.Sprintf("%s", err.Errors()))
		}
		vlanTagging.VlanID = vlanId

	} else {
		vlanTaggingObjNull := types.ObjectNull(vlanTaggingAttr)
		return vlanTaggingObjNull, diags
	}

	vlanTaggingObject, err := types.ObjectValueFrom(ctx, vlanTaggingAttr, vlanTagging)
	if err.HasError() {
		diags.AddError("vlanTagging obj Attr", fmt.Sprintf("%s", err.Errors()))
	}

	return vlanTaggingObject, diags
}

func updateGroupPolicyResourceStateBonjourForwarding(httpResp map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var bonjourForwarding GroupPolicyResourceModelBonjourForwarding

	rulesAttrs := map[string]attr.Type{
		"description": types.StringType,
		"vlan_id":     types.StringType,
		"services":    types.ListType{ElemType: types.StringType},
	}

	bonjourForwardingAttrs := map[string]attr.Type{
		"settings": types.StringType,
		"rules":    types.ListType{ElemType: types.ObjectType{AttrTypes: rulesAttrs}},
	}

	if bf, ok := httpResp["bonjourForwarding"].(map[string]interface{}); ok {

		// settings
		settings, err := utils.ExtractStringAttr(bf, "settings")
		if err.HasError() {
			diags.AddError("settings Attr", fmt.Sprintf("%s", err.Errors()))
		}
		bonjourForwarding.Settings = settings

		// rules
		if rs, ok := bf["rules"].([]map[string]interface{}); ok {
			var rules []types.Object
			for _, r := range rs {

				var rule GroupPolicyResourceModelBonjourForwardingRule
				// description
				description, err := utils.ExtractStringAttr(r, "description")
				if err.HasError() {
					diags.AddError("description Attr", fmt.Sprintf("%s", err.Errors()))
				}
				rule.Description = description

				// vlanId
				vlanId, err := utils.ExtractStringAttr(r, "vlanId")
				if err.HasError() {
					diags.AddError("vlanId Attr", fmt.Sprintf("%s", err.Errors()))
				}
				rule.VlanID = vlanId

				// services
				services, err := utils.ExtractListStringAttr(r, "services")
				if err.HasError() {
					diags.AddError("vlanId Attr", fmt.Sprintf("%s", err.Errors()))
				}
				rule.Services = services

				ruleObj, err := types.ObjectValueFrom(context.Background(), rulesAttrs, rule)
				if err.HasError() {
					diags.AddError("ruleObj Attr", fmt.Sprintf("%s", err.Errors()))
				}

				rules = append(rules, ruleObj)
			}

			rulesArray, err := types.ListValueFrom(context.Background(), types.ObjectType{AttrTypes: rulesAttrs}, rules)
			if err.HasError() {
				diags.AddError("rulesArray Attr", fmt.Sprintf("%s", err.Errors()))
			}
			bonjourForwarding.Rules = rulesArray

		} else {
			rulesArrayNull := types.ListNull(types.ObjectType{AttrTypes: rulesAttrs})
			bonjourForwarding.Rules = rulesArrayNull
		}

	} else {
		bonjourForwardingObjNull := types.ObjectNull(bonjourForwardingAttrs)
		return bonjourForwardingObjNull, diags
	}

	bonjourForwardingObj, err := types.ObjectValueFrom(context.Background(), bonjourForwardingAttrs, bonjourForwarding)
	if err.HasError() {
		diags.AddError("bonjourForwardingObj Attr", fmt.Sprintf("%s", err.Errors()))
	}

	return bonjourForwardingObj, diags
}

// updateGroupPolicyResourceStateContentFiltering updates the resource state with the content filtering data
func updateGroupPolicyResourceStateContentFiltering(ctx context.Context, httpResp map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var contentFiltering GroupPolicyResourceModelContentFiltering

	URLPatternsAttrs := map[string]attr.Type{
		"patterns": types.ListType{ElemType: types.StringType},
		"settings": types.StringType,
	}

	URLCategoriesAttrs := map[string]attr.Type{
		"categories": types.ListType{ElemType: types.StringType},
		"settings":   types.StringType,
	}

	contentFilteringAttrs := map[string]attr.Type{
		"allowed_url_patterns":   types.ObjectType{AttrTypes: URLPatternsAttrs},
		"blocked_url_patterns":   types.ObjectType{AttrTypes: URLPatternsAttrs},
		"blocked_url_categories": types.ObjectType{AttrTypes: URLCategoriesAttrs},
	}

	if cf, ok := httpResp["contentFiltering"].(map[string]interface{}); ok {

		// allowedURLPatterns
		if ap, ok := cf["allowedUrlPatterns"].(map[string]interface{}); ok {
			var allowedURLPatterns GroupPolicyResourceModelUrlPatterns

			// patterns array
			patterns, err := utils.ExtractListStringAttr(ap, "patterns")
			if diags.HasError() {
				diags.AddError("patterns array Attr", fmt.Sprintf("%s", err.Errors()))
			}
			allowedURLPatterns.Patterns = patterns

			// settings
			settings, err := utils.ExtractStringAttr(ap, "settings")
			if diags.HasError() {
				diags.AddError("settings Attr", fmt.Sprintf("%s", err.Errors()))
			}
			allowedURLPatterns.Settings = settings

			allowedURLPatternsObj, err := types.ObjectValueFrom(ctx, URLPatternsAttrs, allowedURLPatterns)
			if err.HasError() {
				diags.Append(err...)
			}
			contentFiltering.AllowedUrlPatterns = allowedURLPatternsObj

		} else {
			allowedURLPatternsObjNull := types.ObjectNull(URLPatternsAttrs)
			contentFiltering.AllowedUrlPatterns = allowedURLPatternsObjNull
		}

		// blockedURLPatterns
		if bp, ok := cf["blockedUrlPatterns"].(map[string]interface{}); ok {
			var blockedURLPatterns GroupPolicyResourceModelUrlPatterns

			// patterns array
			patterns, err := utils.ExtractListStringAttr(bp, "patterns")
			if diags.HasError() {
				diags.AddError("patterns array Attr", fmt.Sprintf("%s", err.Errors()))
			}
			blockedURLPatterns.Patterns = patterns

			// settings
			settings, err := utils.ExtractStringAttr(bp, "settings")
			if diags.HasError() {
				diags.AddError("settings Attr", fmt.Sprintf("%s", err.Errors()))
			}
			blockedURLPatterns.Settings = settings

			blockedURLPatternsObj, err := types.ObjectValueFrom(ctx, URLPatternsAttrs, blockedURLPatterns)
			if err.HasError() {
				diags.Append(err...)
			}
			contentFiltering.BlockedUrlPatterns = blockedURLPatternsObj

		} else {
			blockedURLPatternsObjNull := types.ObjectNull(URLPatternsAttrs)
			contentFiltering.BlockedUrlCategories = blockedURLPatternsObjNull
		}

		// blockedURLCategories
		if bc, ok := cf["blockedUrlCategories"].(map[string]interface{}); ok {
			var blockedURLCategories GroupPolicyResourceModelUrlCategories

			// patterns array
			patterns, err := utils.ExtractListStringAttr(bc, "categories")
			if diags.HasError() {
				diags.AddError("categories array Attr", fmt.Sprintf("%s", err.Errors()))
			}
			blockedURLCategories.Categories = patterns

			// settings
			settings, err := utils.ExtractStringAttr(bc, "settings")
			if diags.HasError() {
				diags.AddError("settings Attr", fmt.Sprintf("%s", err.Errors()))
			}
			blockedURLCategories.Settings = settings

			blockedURLCategoriesObj, err := types.ObjectValueFrom(ctx, URLCategoriesAttrs, blockedURLCategories)
			if err.HasError() {
				diags.AddError("blockedURLCategoriesObj Attr", fmt.Sprintf("%s", err.Errors()))
			}
			contentFiltering.BlockedUrlCategories = blockedURLCategoriesObj

		} else {
			blockedURLCategoriesObjNull := types.ObjectNull(URLCategoriesAttrs)
			contentFiltering.BlockedUrlCategories = blockedURLCategoriesObjNull
		}

	} else {
		contentFilteringObjNull := types.ObjectNull(contentFilteringAttrs)
		return contentFilteringObjNull, diags
	}

	contentFilteringObj, err := types.ObjectValueFrom(context.Background(), contentFilteringAttrs, contentFiltering)
	if err.HasError() {
		diags.AddError("contentFilteringObj Attr", fmt.Sprintf("%s", err.Errors()))
	}

	return contentFilteringObj, diags
}

// updateGroupPolicyResourceState updates the resource state with the provided api data.
func updateGroupPolicyResourceState(ctx context.Context, state *GroupPolicyResourceModel, inlineResp map[string]interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// GroupPolicyId
	if state.GroupPolicyId.IsNull() || state.GroupPolicyId.IsUnknown() {
		groupPolicyId, err := utils.ExtractStringAttr(inlineResp, "groupPolicyId")
		if err.HasError() {
			diags.AddError("groupPolicyId Attribute", fmt.Sprintf("%s", err.Errors()))
			return diags
		}

		if !groupPolicyId.IsNull() {
			state.GroupPolicyId = groupPolicyId
		}
	}

	// check for NetworkId
	if state.NetworkId.IsNull() || state.NetworkId.IsUnknown() {
		state.NetworkId, diags = utils.ExtractStringAttr(inlineResp, "networkId")
		if diags.HasError() {
			diags.AddError("networkId Attribute", "")
			return diags
		}
	}

	// Import ID
	if !state.NetworkId.IsNull() || !state.NetworkId.IsUnknown() && !state.GroupPolicyId.IsNull() || !state.GroupPolicyId.IsUnknown() {
		state.ID = types.StringValue(state.NetworkId.ValueString() + "," + state.GroupPolicyId.ValueString())
	} else {
		state.ID = types.StringNull()
	}

	// Name
	if state.Name.IsNull() || state.Name.IsUnknown() {
		state.Name, diags = utils.ExtractStringAttr(inlineResp, "name")
		if diags.HasError() {
			diags.AddError("name Attribute", "")
			return diags
		}
	}

	// Update Scheduling
	if state.Scheduling.IsNull() || state.Scheduling.IsUnknown() {
		scheduling, err := updateGroupPolicyResourceStateScheduling(inlineResp)
		if err.HasError() {
			diags.AddError("scheduling Attr", fmt.Sprintf("%s", err.Errors()))
		}
		state.Scheduling = scheduling
	}

	// Update Bandwidth
	bandwidth, bandwidthErr := updateGroupPolicyResourceStateBandwidth(inlineResp)
	if bandwidthErr.HasError() {
		diags.AddError("bandwidth Attr", fmt.Sprintf("%s", bandwidthErr.Errors()))
	}
	state.Bandwidth = bandwidth

	//SplashAuthSettings
	if state.SplashAuthSettings.IsNull() || state.SplashAuthSettings.IsUnknown() {
		state.SplashAuthSettings, diags = utils.ExtractStringAttr(inlineResp, "splashAuthSettings")
		if diags.HasError() {
			return diags
		}
	}

	// Update VlanTagging
	vlanTaggingObj, vlanTaggingObjErr := updateGroupPolicyResourceStateVlanTagging(ctx, inlineResp)
	if vlanTaggingObjErr.HasError() {
		diags.AddError("vlanTagging Attr", fmt.Sprintf("%s", vlanTaggingObjErr.Errors()))
	}
	state.VlanTagging = vlanTaggingObj

	// Update BonjourForwarding
	if state.BonjourForwarding.IsNull() || state.BonjourForwarding.IsUnknown() {
		bonjourForwarding, err := updateGroupPolicyResourceStateBonjourForwarding(inlineResp)
		if err.HasError() {
			diags.AddError("vlanTagging Attr", fmt.Sprintf("%s", err.Errors()))
		}
		state.BonjourForwarding = bonjourForwarding

	}

	// Update FirewallAndTrafficShaping
	firewallAndTrafficShapingObject, firewallAndTrafficShapingObjectDiags := updateGroupPolicyResourceStateFirewallAndTrafficShapingRules(ctx, inlineResp)
	if firewallAndTrafficShapingObjectDiags.HasError() {
		diags.Append(firewallAndTrafficShapingObjectDiags...)
	}
	state.FirewallAndTrafficShaping = firewallAndTrafficShapingObject

	// Update ContentFiltering
	contentFilteringObj, cfDiags := updateGroupPolicyResourceStateContentFiltering(ctx, inlineResp)
	if cfDiags.HasError() {
		diags.Append(cfDiags...)
	}
	state.ContentFiltering = contentFilteringObj

	return diags
}

// Create handles the creation of the group policy.
func (r *NetworksGroupPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan GroupPolicyResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupPolicy, groupPolicyDiags := updateGroupPolicyResourcePayload(&plan)
	if groupPolicyDiags.HasError() {

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
		time.Sleep(retryDelay * time.Millisecond)

		inline, httpResp, err := r.client.NetworksApi.CreateNetworkGroupPolicy(ctx, plan.NetworkId.ValueString()).CreateNetworkGroupPolicyRequest(groupPolicy).Execute()
		return inline, httpResp, err
	}
	// Use retryOn4xx for the API call as the meraki API backend returns HTTP 400 messages as a result of collision issues with rapid creation of postgres GroupPolicyIds.
	inlineResp, httpResp, err := tools.CustomHttpRequestRetry(ctx, maxRetries, retryDelay, apiCall)
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

	diags = updateGroupPolicyResourceState(ctx, &plan, inlineResp)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Read handles reading the group policy.
func (r *NetworksGroupPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state GroupPolicyResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	inlineResp, httpResp, err := r.client.NetworksApi.GetNetworkGroupPolicy(ctx, state.NetworkId.ValueString(), state.GroupPolicyId.ValueString()).Execute()
	if err != nil {
		var responseBody string
		if httpResp != nil && httpResp.Body != nil {
			bodyBytes, readErr := io.ReadAll(httpResp.Body)
			if readErr == nil {
				responseBody = string(bodyBytes)
			}
		}

		resp.Diagnostics.AddError(
			"Error reading group policy",
			fmt.Sprintf("Could not read group policy, unexpected error: %s\nHTTP Response: %v\nResponse Body: %s", err, httpResp, responseBody),
		)
		return
	}

	// Update the state with the new state
	diags = updateGroupPolicyResourceState(ctx, &state, inlineResp)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update handles updating the group policy.
func (r *NetworksGroupPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan GroupPolicyResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state GroupPolicyResourceModel

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupPolicy, groupPolicyDiags := updateGroupPolicyResourcePayload(&plan)
	if groupPolicyDiags.HasError() {

		resp.Diagnostics.AddError(
			"Error creating group policy payload",
			fmt.Sprintf("Unexpected error: %s", groupPolicyDiags),
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

	inlineResp, httpResp, err := r.client.NetworksApi.UpdateNetworkGroupPolicy(ctx, plan.NetworkId.ValueString(), state.GroupPolicyId.ValueString()).UpdateNetworkGroupPolicyRequest(groupPolicyUpdate).Execute()
	if err != nil {
		var responseBody string
		if httpResp != nil && httpResp.Body != nil {
			bodyBytes, readErr := io.ReadAll(httpResp.Body)
			if readErr == nil {
				responseBody = string(bodyBytes)
			}
		}

		resp.Diagnostics.AddError(
			"Error updating group policy",
			fmt.Sprintf("Could not update group policy, unexpected error: %s\nHTTP Response: %v\nResponse Body: %s", err, httpResp, responseBody),
		)
		return
	}

	// Update the state with the new plan
	diags = updateGroupPolicyResourceState(ctx, &plan, inlineResp)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Delete handles deleting the group policy.
func (r *NetworksGroupPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state GroupPolicyResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check for empty GroupPolicyId
	if state.GroupPolicyId.IsNull() || state.GroupPolicyId.IsUnknown() {
		resp.Diagnostics.AddError("DELETE, Received empty GroupPolicy.", fmt.Sprintf("Name: %s, Id: %s", state.Name.ValueString(), state.GroupPolicyId.ValueString()))
		return
	}

	httpResp, err := r.client.NetworksApi.DeleteNetworkGroupPolicy(ctx, state.NetworkId.ValueString(), state.GroupPolicyId.ValueString()).Execute()
	if err != nil {
		var responseBody string
		if httpResp != nil && httpResp.Body != nil {
			bodyBytes, readErr := io.ReadAll(httpResp.Body)
			if readErr == nil {
				responseBody = string(bodyBytes)
			}
		}
		resp.Diagnostics.AddError(
			"Error deleting group policy",
			fmt.Sprintf("Could not delete group policy, unexpected error: %s\nHTTP Response: %v\nResponse Body: %s", err, httpResp, responseBody),
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
