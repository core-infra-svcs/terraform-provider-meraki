package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
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
	Id                                           jsontypes.String          `tfsdk:"id"`
	NetworkId                                    jsontypes.String          `tfsdk:"network_id"`
	GroupPolicyId                                jsontypes.String          `tfsdk:"group_policy_id"`
	Name                                         jsontypes.String          `tfsdk:"name" json:"name"`
	SplashAuthSettings                           jsontypes.String          `tfsdk:"splash_auth_settings"`
	BandwidthSettings                            jsontypes.String          `tfsdk:"bandwidth_settings"`
	BandwidthLimitUp                             jsontypes.Int64           `tfsdk:"bandwidth_limit_up"`
	BandwidthLimitDown                           jsontypes.Int64           `tfsdk:"bandwidth_limit_down"`
	BonjourForwardingSettings                    jsontypes.String          `tfsdk:"bonjour_forwarding_settings"`
	BonjourForwardingRules                       []Rule                    `tfsdk:"bonjour_forwarding_rules"`
	FirewallAndTrafficShaping                    FirewallAndTrafficShaping `tfsdk:"firewall_and_traffic_shaping"`
	SchedulingEnabled                            jsontypes.Bool            `tfsdk:"scheduling_enabled"`
	SchedulingFridayActive                       jsontypes.Bool            `tfsdk:"scheduling_friday_active"`
	SchedulingFridayFrom                         jsontypes.String          `tfsdk:"scheduling_friday_from"`
	SchedulingFridayTo                           jsontypes.String          `tfsdk:"scheduling_friday_to"`
	SchedulingMondayActive                       jsontypes.Bool            `tfsdk:"scheduling_monday_active"`
	SchedulingMondayFrom                         jsontypes.String          `tfsdk:"scheduling_monday_from"`
	SchedulingMondayTo                           jsontypes.String          `tfsdk:"scheduling_monday_to"`
	SchedulingTuesdayActive                      jsontypes.Bool            `tfsdk:"scheduling_tuesday_active"`
	SchedulingTuesdayFrom                        jsontypes.String          `tfsdk:"scheduling_tuesday_from"`
	SchedulingTuesdayTo                          jsontypes.String          `tfsdk:"scheduling_tuesday_to"`
	SchedulingWednesdayActive                    jsontypes.Bool            `tfsdk:"scheduling_wednesday_active"`
	SchedulingWednesdayFrom                      jsontypes.String          `tfsdk:"scheduling_wednesday_from"`
	SchedulingWednesdayTo                        jsontypes.String          `tfsdk:"scheduling_wednesday_to"`
	SchedulingThursdayActive                     jsontypes.Bool            `tfsdk:"scheduling_thursday_active"`
	SchedulingThursdayFrom                       jsontypes.String          `tfsdk:"scheduling_thursday_from"`
	SchedulingThursdayTo                         jsontypes.String          `tfsdk:"scheduling_thursday_to"`
	SchedulingSaturdayActive                     jsontypes.Bool            `tfsdk:"scheduling_saturday_active"`
	SchedulingSaturdayFrom                       jsontypes.String          `tfsdk:"scheduling_saturday_from"`
	SchedulingSaturdayTo                         jsontypes.String          `tfsdk:"scheduling_saturday_to"`
	SchedulingSundayActive                       jsontypes.Bool            `tfsdk:"scheduling_sunday_active"`
	SchedulingSundayFrom                         jsontypes.String          `tfsdk:"scheduling_sunday_from"`
	SchedulingSundayTo                           jsontypes.String          `tfsdk:"scheduling_sunday_to"`
	VlanTaggingSettings                          jsontypes.String          `tfsdk:"vlan_tagging_settings"`
	VlanTaggingVlanId                            jsontypes.String          `tfsdk:"vlan_tagging_vlan_id"`
	ContentFilteringAllowUrlPatternsSettings     jsontypes.String          `tfsdk:"content_filtering_allow_url_patterns_settings"`
	ContentFilteringAllowUrlPatterns             []string                  `tfsdk:"content_filtering_allow_url_patterns"`
	ContentFilteringBlockedUrlCategoriesSettings jsontypes.String          `tfsdk:"content_filtering_blocked_url_categories_settings"`
	ContentFilteringBlockedUrlCategories         []string                  `tfsdk:"content_filtering_blocked_url_categories"`
	ContentFilteringBlockedUrlPatternsSettings   jsontypes.String          `tfsdk:"content_filtering_blocked_url_patterns_settings"`
	ContentFilteringBlockedUrlPatterns           []string                  `tfsdk:"content_filtering_blocked_url_patterns"`
}

type OutPutData struct {
	GroupPolicyId                   jsontypes.String                `json:"groupPolicyId"`
	Name                            jsontypes.String                `json:"name"`
	SplashAuthSettings              jsontypes.String                `json:"splashAuthSettings"`
	Bandwidth                       Bandwidth                       `json:"bandwidth"`
	BonjourForwarding               BonjourForwarding               `json:"bonjourForwarding"`
	VlanTagging                     VlanTagging                     `json:"vlanTagging"`
	Scheduling                      Scheduling                      `json:"scheduling"`
	OutputFirewallAndTrafficShaping OutputFirewallAndTrafficShaping `tfsdk:"firewall_and_traffic_shaping"`
}

type Bandwidth struct {
	BandwidthLimits BandwidthLimits  `tfsdk:"bandwidth_limits"`
	Settings        jsontypes.String `tfsdk:"settings"`
}

type BandwidthLimits struct {
	LimitUp   jsontypes.Int64 `tfsdk:"limit_up"`
	LimitDown jsontypes.Int64 `tfsdk:"limit_down"`
}

type Rule struct {
	Description jsontypes.String `tfsdk:"description"`
	VlanId      jsontypes.String `tfsdk:"vlan_id"`
	Services    []string         `tfsdk:"services"`
}

type BonjourForwardingRule struct {
	Description string   `json:"description"`
	VlanId      string   `json:"vlanId"`
	Services    []string `json:"services"`
}

type BonjourForwarding struct {
	BonjourForwardingSettings string `json:"settings"`
	BonjourForwardingRules    []Rule `json:"rules"`
}

type FirewallAndTrafficShaping struct {
	Settings            jsontypes.String     `tfsdk:"settings"`
	L3FirewallRules     []L3FirewallRules    `tfsdk:"l3_firewall_rules"`
	L7FirewallRules     []L7FirewallRule     `tfsdk:"l7_firewall_rules"`
	TrafficShapingRules []TrafficShapingRule `tfsdk:"traffic_shaping_rules"`
}

type OutputFirewallAndTrafficShaping struct {
	Settings                  jsontypes.String
	L3FirewallRules           []L3FirewallRules
	L7FirewallRules           []L7FirewallRule
	OutputTrafficShapingRules []OutputTrafficShapingRule
}

type TrafficShapingRule struct {
	DscpTagValue                     jsontypes.Int64  `tfsdk:"dscp_tag_value"`
	PcpTagValue                      jsontypes.Int64  `tfsdk:"pcp_tag_value"`
	PerClientBandwidthLimitsSettings jsontypes.String `tfsdk:"per_client_bandwidth_limits_settings"`
	BandwidthLimitDown               jsontypes.Int64  `tfsdk:"bandwidth_limit_down"`
	BandwidthLimitUp                 jsontypes.Int64  `tfsdk:"bandwidth_limit_up"`
	Definitions                      []Definition     `tfsdk:"definitions"`
}

type OutputTrafficShapingRule struct {
	DscpTagValue             jsontypes.Int64
	PcpTagValue              jsontypes.Int64
	PerClientBandwidthLimits PerClientBandwidthLimits
	Definitions              []Definition
}

type PerClientBandwidthLimits struct {
	BandwidthLimits BandwidthLimits
	Settings        jsontypes.String
}

type Definition struct {
	Value jsontypes.String `tfsdk:"value"`
	Type  jsontypes.String `tfsdk:"type"`
}

type L3FirewallRules struct {
	Comment  jsontypes.String `tfsdk:"comment"`
	DestCidr jsontypes.String `tfsdk:"dest_cidr"`
	DestPort jsontypes.String `tfsdk:"dest_port"`
	Policy   jsontypes.String `tfsdk:"policy"`
	Protocol jsontypes.String `tfsdk:"protocol"`
}

type L7FirewallRule struct {
	Value  jsontypes.String `tfsdk:"value"`
	Type   jsontypes.String `tfsdk:"type"`
	Policy jsontypes.String `tfsdk:"policy"`
}

type Scheduling struct {
	Enabled   jsontypes.Bool `tfsdk:"enabled"`
	Friday    Schedule       `tfsdk:"friday"`
	Monday    Schedule       `tfsdk:"monday"`
	Saturday  Schedule       `tfsdk:"saturday"`
	Sunday    Schedule       `tfsdk:"sunday"`
	Thursday  Schedule       `tfsdk:"thursday"`
	Tuesday   Schedule       `tfsdk:"tuesday"`
	Wednesday Schedule       `tfsdk:"wednesday"`
}

type Schedule struct {
	From   jsontypes.String `tfsdk:"from"`
	To     jsontypes.String `tfsdk:"to"`
	Active jsontypes.Bool   `tfsdk:"active"`
}

type VlanTagging struct {
	Settings jsontypes.String `tfsdk:"settings"`
	VlanId   jsontypes.String `tfsdk:"vlan_id"`
}

type ContentFiltering struct {
	AllowedUrlPatterns   AllowedUrlPatterns   `tfsdk:"allowed_url_patterns"`
	BlockedUrlCategories BlockedUrlCategories `tfsdk:"blocked_url_categories"`
	BlockedUrlPatterns   BlockedUrlPatterns   `tfsdk:"blocked_url_patterns"`
}

type AllowedUrlPatterns struct {
	Settings jsontypes.String `tfsdk:"settings"`
	Patterns []string         `tfsdk:"patterns"`
}
type BlockedUrlCategories struct {
	Settings   jsontypes.String `tfsdk:"settings"`
	Categories []string         `tfsdk:"categories"`
}
type BlockedUrlPatterns struct {
	Settings jsontypes.String `tfsdk:"settings"`
	Patterns []string         `tfsdk:"patterns"`
}

type CreateHttpRequestData struct {
	CreateRequestData openApiClient.InlineObject87
	UpdateRequestData openApiClient.InlineObject88
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
					stringvalidator.LengthBetween(8, 31),
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
					stringvalidator.LengthBetween(8, 31),
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
			"bandwidth_settings": schema.StringAttribute{
				MarkdownDescription: "The bandwidth settings for clients bound to your group policy. How bandwidth limits are enforced. Can be 'network default', 'ignore' or 'custom'.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"network default", "ignore", "custom"}...),
				},
			},
			"bandwidth_limit_up": schema.Int64Attribute{
				MarkdownDescription: "The bandwidth limits object, specifying upload and download speed for clients bound to the group policy. These are only enforced if 'settings' is set to 'custom'. The maximum upload limit (integer, in Kbps). null indicates no limit.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.Int64Type,
			},
			"bandwidth_limit_down": schema.Int64Attribute{
				MarkdownDescription: "The bandwidth limits object, specifying upload and download speed for clients bound to the group policy. These are only enforced if 'settings' is set to 'custom'. The maximum download limit (integer, in Kbps). null indicates no limit.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.Int64Type,
			},
			"bonjour_forwarding_settings": schema.StringAttribute{
				MarkdownDescription: "The Bonjour settings for your group policy. Only valid if your network has a wireless configuration.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"bonjour_forwarding_rules": schema.SetNestedAttribute{
				Description: "A list of the Bonjour forwarding rules for your group policy. If 'settings' is set to 'custom', at least one rule must be specified.",
				Optional:    true,
				Computed:    true,
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
							Description: "A list of Bonjour services. At least one service must be specified. Available services are 'All Services', 'AirPlay', 'AFP', 'BitTorrent', 'FTP', 'iChat', 'iTunes', 'Printers', 'Samba', 'Scanners' and 'SSH'",
							ElementType: jsontypes.StringType,
							CustomType:  jsontypes.SetType[jsontypes.String](),
							Optional:    true,
							Computed:    true,
						},
					},
				},
			},
			"scheduling_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether scheduling is enabled (true) or disabled (false). Defaults to false. If true, the schedule objects for each day of the week (monday - sunday) are parsed.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.BoolType,
			},
			"scheduling_friday_active": schema.BoolAttribute{
				MarkdownDescription: "Whether the schedule is active (true) or inactive (false) during the time specified between 'from' and 'to'. Defaults to true.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.BoolType,
			},
			"scheduling_friday_from": schema.StringAttribute{
				MarkdownDescription: "The time, from '00:00' to '24:00'. Must be less than the time specified in 'to'. Defaults to '00:00'. Only 30 minute increments are allowed.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"scheduling_friday_to": schema.StringAttribute{
				MarkdownDescription: "The time, from '00:00' to '24:00'. Must be less than the time specified in 'to'. Defaults to '00:00'. Only 30 minute increments are allowed.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"scheduling_monday_active": schema.BoolAttribute{
				MarkdownDescription: "Whether the schedule is active (true) or inactive (false) during the time specified between 'from' and 'to'. Defaults to true.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.BoolType,
			},
			"scheduling_monday_from": schema.StringAttribute{
				MarkdownDescription: "The time, from '00:00' to '24:00'. Must be less than the time specified in 'to'. Defaults to '00:00'. Only 30 minute increments are allowed.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"scheduling_monday_to": schema.StringAttribute{
				MarkdownDescription: "The time, from '00:00' to '24:00'. Must be less than the time specified in 'to'. Defaults to '00:00'. Only 30 minute increments are allowed.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"scheduling_tuesday_active": schema.BoolAttribute{
				MarkdownDescription: "Whether the schedule is active (true) or inactive (false) during the time specified between 'from' and 'to'. Defaults to true.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.BoolType,
			},
			"scheduling_tuesday_from": schema.StringAttribute{
				MarkdownDescription: "The time, from '00:00' to '24:00'. Must be less than the time specified in 'to'. Defaults to '00:00'. Only 30 minute increments are allowed.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"scheduling_tuesday_to": schema.StringAttribute{
				MarkdownDescription: "The time, from '00:00' to '24:00'. Must be less than the time specified in 'to'. Defaults to '00:00'. Only 30 minute increments are allowed.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"scheduling_wednesday_active": schema.BoolAttribute{
				MarkdownDescription: "Whether the schedule is active (true) or inactive (false) during the time specified between 'from' and 'to'. Defaults to true.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.BoolType,
			},
			"scheduling_wednesday_from": schema.StringAttribute{
				MarkdownDescription: "The time, from '00:00' to '24:00'. Must be less than the time specified in 'to'. Defaults to '00:00'. Only 30 minute increments are allowed.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"scheduling_wednesday_to": schema.StringAttribute{
				MarkdownDescription: "The time, from '00:00' to '24:00'. Must be less than the time specified in 'to'. Defaults to '00:00'. Only 30 minute increments are allowed.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"scheduling_thursday_active": schema.BoolAttribute{
				MarkdownDescription: "Whether the schedule is active (true) or inactive (false) during the time specified between 'from' and 'to'. Defaults to true.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.BoolType,
			},
			"scheduling_thursday_from": schema.StringAttribute{
				MarkdownDescription: "The time, from '00:00' to '24:00'. Must be less than the time specified in 'to'. Defaults to '00:00'. Only 30 minute increments are allowed.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"scheduling_thursday_to": schema.StringAttribute{
				MarkdownDescription: "The time, from '00:00' to '24:00'. Must be less than the time specified in 'to'. Defaults to '00:00'. Only 30 minute increments are allowed.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"scheduling_saturday_active": schema.BoolAttribute{
				MarkdownDescription: "Whether the schedule is active (true) or inactive (false) during the time specified between 'from' and 'to'. Defaults to true.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.BoolType,
			},
			"scheduling_saturday_from": schema.StringAttribute{
				MarkdownDescription: "The time, from '00:00' to '24:00'. Must be less than the time specified in 'to'. Defaults to '00:00'. Only 30 minute increments are allowed.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"scheduling_saturday_to": schema.StringAttribute{
				MarkdownDescription: "The time, from '00:00' to '24:00'. Must be less than the time specified in 'to'. Defaults to '00:00'. Only 30 minute increments are allowed.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"scheduling_sunday_active": schema.BoolAttribute{
				MarkdownDescription: "Whether the schedule is active (true) or inactive (false) during the time specified between 'from' and 'to'. Defaults to true.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.BoolType,
			},
			"scheduling_sunday_from": schema.StringAttribute{
				MarkdownDescription: "The time, from '00:00' to '24:00'. Must be less than the time specified in 'to'. Defaults to '00:00'. Only 30 minute increments are allowed.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"scheduling_sunday_to": schema.StringAttribute{
				MarkdownDescription: "The time, from '00:00' to '24:00'. Must be less than the time specified in 'to'. Defaults to '00:00'. Only 30 minute increments are allowed.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"content_filtering_allow_url_patterns_settings": schema.StringAttribute{
				MarkdownDescription: "The content filtering settings for your group policy. Settings for allowed URL patterns. How URL patterns are applied. Can be 'network default', 'append' or 'override'.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"network default", "append", "override"}...),
				},
			},
			"content_filtering_allow_url_patterns": schema.SetAttribute{
				Description: "A list of URL patterns that are allowed for content filtering for your group policy",
				ElementType: jsontypes.StringType,
				CustomType:  jsontypes.SetType[jsontypes.String](),
				Required:    true,
			},
			"content_filtering_blocked_url_categories_settings": schema.StringAttribute{
				MarkdownDescription: "The content filtering settings for your group policy. Settings for blocked URL categories. How URL categories are applied. Can be 'network default', 'append' or 'override'.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"network default", "append", "override"}...),
				},
			},
			"content_filtering_blocked_url_categories": schema.SetAttribute{
				Description: "A list of URL categories to block for content filtering for your group policy",
				ElementType: jsontypes.StringType,
				CustomType:  jsontypes.SetType[jsontypes.String](),
				Required:    true,
			},
			"content_filtering_blocked_url_patterns_settings": schema.StringAttribute{
				MarkdownDescription: "The content filtering settings for your group policy. Settings for blocked URL patterns. How URL patterns are applied. Can be 'network default', 'append' or 'override'.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"network default", "append", "override"}...),
				},
			},
			"content_filtering_blocked_url_patterns": schema.SetAttribute{
				Description: "A list of URL patterns that are blocked for content filtering for your group policy",
				ElementType: jsontypes.StringType,
				CustomType:  jsontypes.SetType[jsontypes.String](),
				Required:    true,
			},
			"firewall_and_traffic_shaping": schema.SingleNestedAttribute{
				MarkdownDescription: "The firewall and traffic shaping rules and settings for your policy.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"settings": schema.StringAttribute{
						MarkdownDescription: "How firewall and traffic shaping rules are enforced. Can be 'network default', 'ignore' or 'custom'.",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.StringType,
					},
					"l3_firewall_rules": schema.SetNestedAttribute{
						Description: "An ordered array of the L3 firewall rules",
						Required:    true,
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
									Validators: []validator.String{
										stringvalidator.OneOf([]string{"allow", "deny"}...),
									},
								},
								"protocol": schema.StringAttribute{
									MarkdownDescription: "The type of protocol (must be 'tcp', 'udp', 'icmp', 'icmp6' or 'any')",
									Optional:            true,
									Computed:            true,
									CustomType:          jsontypes.StringType,
									Validators: []validator.String{
										stringvalidator.OneOf([]string{"tcp", "udp", "icmp", "icmp6", "any"}...),
									},
								},
							},
						},
					},
					"l7_firewall_rules": schema.SetNestedAttribute{
						Description: "An ordered array of the L7 firewall rules",
						Required:    true,
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
						Required: true,
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
								"per_client_bandwidth_limits_settings": schema.StringAttribute{
									MarkdownDescription: "How bandwidth limits are applied by your rule. Can be one of 'network default', 'ignore' or 'custom'.",
									Required:            true,
									CustomType:          jsontypes.StringType,
								},
								"bandwidth_limit_down": schema.Int64Attribute{
									MarkdownDescription: "The maximum download limit (integer, in Kbps).",
									Optional:            true,
									CustomType:          jsontypes.Int64Type,
								},
								"bandwidth_limit_up": schema.Int64Attribute{
									MarkdownDescription: "The maximum upload limit (integer, in Kbps).",
									Optional:            true,
									CustomType:          jsontypes.Int64Type,
								},
								"definitions": schema.SetNestedAttribute{
									Required: true,
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
			"vlan_tagging_settings": schema.StringAttribute{
				MarkdownDescription: "How VLAN tagging is applied. Can be 'network default', 'ignore' or 'custom'.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"vlan_tagging_vlan_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the vlan you want to tag. This only applies if 'settings' is set to 'custom'.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
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

	createNetworkGroupPolicy := *openApiClient.NewInlineObject87(data.Name.ValueString())

	if !data.SplashAuthSettings.IsUnknown() {
		createNetworkGroupPolicy.SetSplashAuthSettings(data.SplashAuthSettings.ValueString())

	}

	if !data.BandwidthSettings.IsUnknown() {
		var bandwidth openApiClient.NetworksNetworkIdGroupPoliciesBandwidth
		bandwidth.SetSettings(data.BandwidthSettings.ValueString())
		var bandwidthLimits openApiClient.NetworksNetworkIdGroupPoliciesBandwidthBandwidthLimits
		if !data.BandwidthLimitUp.IsUnknown() {
			bandwidthLimits.SetLimitUp(int32(data.BandwidthLimitUp.ValueInt64()))
		}
		if !data.BandwidthLimitDown.IsUnknown() {
			bandwidthLimits.SetLimitDown(int32(data.BandwidthLimitDown.ValueInt64()))
		}
		bandwidth.SetBandwidthLimits(bandwidthLimits)
		createNetworkGroupPolicy.SetBandwidth(bandwidth)
	}

	if len(data.BonjourForwardingRules) > 0 {
		var bonjourForwarding openApiClient.NetworksNetworkIdGroupPoliciesBonjourForwarding
		var bonjourForwardingRules []openApiClient.NetworksNetworkIdGroupPoliciesBonjourForwardingRules
		for _, attribute := range data.BonjourForwardingRules {
			var bonjourForwardingRule openApiClient.NetworksNetworkIdGroupPoliciesBonjourForwardingRules
			if !attribute.Description.IsUnknown() {
				bonjourForwardingRule.SetDescription(attribute.Description.ValueString())
			}
			bonjourForwardingRule.SetVlanId(attribute.VlanId.ValueString())
			bonjourForwardingRule.SetServices(attribute.Services)
			bonjourForwardingRules = append(bonjourForwardingRules, bonjourForwardingRule)
		}
		bonjourForwarding.SetRules(bonjourForwardingRules)
		if !data.BonjourForwardingSettings.IsUnknown() {
			bonjourForwarding.SetSettings(data.BonjourForwardingSettings.ValueString())
		}
		createNetworkGroupPolicy.SetBonjourForwarding(bonjourForwarding)
	}

	var firewallAndTrafficShaping openApiClient.NetworksNetworkIdGroupPoliciesFirewallAndTrafficShaping

	if !data.FirewallAndTrafficShaping.Settings.IsUnknown() {
		firewallAndTrafficShaping.SetSettings(data.FirewallAndTrafficShaping.Settings.ValueString())
	}

	if len(data.FirewallAndTrafficShaping.L3FirewallRules) > 0 {
		var l3s []openApiClient.NetworksNetworkIdGroupPoliciesFirewallAndTrafficShapingL3FirewallRules
		for _, attribute := range data.FirewallAndTrafficShaping.L3FirewallRules {
			var l3 openApiClient.NetworksNetworkIdGroupPoliciesFirewallAndTrafficShapingL3FirewallRules
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

	if len(data.FirewallAndTrafficShaping.L7FirewallRules) > 0 {
		var l7s []openApiClient.NetworksNetworkIdGroupPoliciesFirewallAndTrafficShapingL7FirewallRules
		for _, attribute := range data.FirewallAndTrafficShaping.L7FirewallRules {
			var l7 openApiClient.NetworksNetworkIdGroupPoliciesFirewallAndTrafficShapingL7FirewallRules
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

	if len(data.FirewallAndTrafficShaping.TrafficShapingRules) > 0 {
		var tfs []openApiClient.NetworksNetworkIdGroupPoliciesFirewallAndTrafficShapingTrafficShapingRules
		for _, attribute := range data.FirewallAndTrafficShaping.TrafficShapingRules {
			var tf openApiClient.NetworksNetworkIdGroupPoliciesFirewallAndTrafficShapingTrafficShapingRules
			if !attribute.DscpTagValue.IsUnknown() {
				tf.SetDscpTagValue(int32(attribute.DscpTagValue.ValueInt64()))
			}
			if !attribute.PcpTagValue.IsUnknown() {
				tf.SetPcpTagValue(int32(attribute.PcpTagValue.ValueInt64()))
			}
			var perclientBamdWidthLimits openApiClient.NetworksNetworkIdApplianceTrafficShapingRulesPerClientBandwidthLimits

			if !attribute.PerClientBandwidthLimitsSettings.IsUnknown() {
				var bandwidthLimits openApiClient.NetworksNetworkIdApplianceTrafficShapingRulesPerClientBandwidthLimitsBandwidthLimits

				if attribute.PerClientBandwidthLimitsSettings.ValueString() != "network default" {

					if !attribute.BandwidthLimitDown.Int64Value.IsUnknown() {
						bandwidthLimits.SetLimitDown(int32(attribute.BandwidthLimitDown.ValueInt64()))
					}

					if !attribute.BandwidthLimitUp.Int64Value.IsUnknown() {
						bandwidthLimits.SetLimitUp(int32(attribute.BandwidthLimitUp.ValueInt64()))
					}
					perclientBamdWidthLimits.SetBandwidthLimits(bandwidthLimits)
				} else {

					if attribute.BandwidthLimitDown.ValueInt64() != jsontypes.Int64Null().ValueInt64() {
						resp.Diagnostics.AddError("Error:", "No need to add  band width limits for network default settings")
						return
					}

					if attribute.BandwidthLimitUp.ValueInt64() != jsontypes.Int64Null().ValueInt64() {
						resp.Diagnostics.AddError("Error:", "No need to add  band width limits for network default settings")
						return
					}
				}
				perclientBamdWidthLimits.SetSettings(attribute.PerClientBandwidthLimitsSettings.ValueString())
				tf.SetPerClientBandwidthLimits(perclientBamdWidthLimits)
			}
			var defs []openApiClient.NetworksNetworkIdApplianceTrafficShapingRulesDefinitions
			if len(attribute.Definitions) > 0 {
				for _, attribute := range attribute.Definitions {
					var def openApiClient.NetworksNetworkIdApplianceTrafficShapingRulesDefinitions
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

	if !data.SchedulingEnabled.IsUnknown() {
		var schedule openApiClient.NetworksNetworkIdGroupPoliciesScheduling
		schedule.SetEnabled(data.SchedulingEnabled.ValueBool())
		if !data.SchedulingFridayActive.IsUnknown() {
			var friday openApiClient.NetworksNetworkIdGroupPoliciesSchedulingFriday
			friday.SetActive(data.SchedulingFridayActive.ValueBool())
			friday.SetFrom(data.SchedulingFridayFrom.ValueString())
			friday.SetTo(data.SchedulingFridayTo.ValueString())
			schedule.SetFriday(friday)
		}
		if !data.SchedulingMondayActive.IsUnknown() {
			var monday openApiClient.NetworksNetworkIdGroupPoliciesSchedulingMonday
			monday.SetActive(data.SchedulingMondayActive.ValueBool())
			monday.SetFrom(data.SchedulingMondayFrom.ValueString())
			monday.SetTo(data.SchedulingMondayTo.ValueString())
			schedule.SetMonday(monday)
		}
		if !data.SchedulingTuesdayActive.IsUnknown() {
			var tuesday openApiClient.NetworksNetworkIdGroupPoliciesSchedulingTuesday
			tuesday.SetActive(data.SchedulingTuesdayActive.ValueBool())
			tuesday.SetFrom(data.SchedulingTuesdayFrom.ValueString())
			tuesday.SetTo(data.SchedulingTuesdayTo.ValueString())
			schedule.SetTuesday(tuesday)
		}
		if !data.SchedulingWednesdayActive.IsUnknown() {
			var wednesday openApiClient.NetworksNetworkIdGroupPoliciesSchedulingWednesday
			wednesday.SetActive(data.SchedulingWednesdayActive.ValueBool())
			wednesday.SetFrom(data.SchedulingWednesdayFrom.ValueString())
			wednesday.SetTo(data.SchedulingWednesdayTo.ValueString())
			schedule.SetWednesday(wednesday)
		}
		if !data.SchedulingThursdayActive.IsUnknown() {
			var thursday openApiClient.NetworksNetworkIdGroupPoliciesSchedulingThursday
			thursday.SetActive(data.SchedulingThursdayActive.ValueBool())
			thursday.SetFrom(data.SchedulingThursdayFrom.ValueString())
			thursday.SetTo(data.SchedulingThursdayTo.ValueString())
			schedule.SetThursday(thursday)
		}
		if !data.SchedulingSaturdayActive.IsUnknown() {
			var saturday openApiClient.NetworksNetworkIdGroupPoliciesSchedulingSaturday
			saturday.SetActive(data.SchedulingSaturdayActive.ValueBool())
			saturday.SetFrom(data.SchedulingSaturdayFrom.ValueString())
			saturday.SetTo(data.SchedulingSaturdayTo.ValueString())
			schedule.SetSaturday(saturday)
		}
		if !data.SchedulingSundayActive.IsUnknown() {
			var sunday openApiClient.NetworksNetworkIdGroupPoliciesSchedulingSunday
			sunday.SetActive(data.SchedulingSundayActive.ValueBool())
			sunday.SetFrom(data.SchedulingSundayFrom.ValueString())
			sunday.SetTo(data.SchedulingSundayTo.ValueString())
			schedule.SetSunday(sunday)
		}
		createNetworkGroupPolicy.SetScheduling(schedule)
	}

	if !data.VlanTaggingSettings.IsUnknown() {
		if !data.VlanTaggingVlanId.IsUnknown() {
			var v openApiClient.NetworksNetworkIdGroupPoliciesVlanTagging
			v.SetSettings(data.VlanTaggingSettings.ValueString())
			v.SetVlanId(data.VlanTaggingVlanId.ValueString())
			createNetworkGroupPolicy.SetVlanTagging(v)
		}
	}
	var contentFiltering openApiClient.NetworksNetworkIdGroupPoliciesContentFiltering
	contentFilteringStatus := false

	if !data.ContentFilteringAllowUrlPatternsSettings.IsUnknown() {
		var allowedUrlPatternData openApiClient.NetworksNetworkIdGroupPoliciesContentFilteringAllowedUrlPatterns
		allowedUrlPatternData.SetSettings(data.ContentFilteringAllowUrlPatternsSettings.ValueString())
		allowedUrlPatternData.SetPatterns(data.ContentFilteringAllowUrlPatterns)
		contentFiltering.SetAllowedUrlPatterns(allowedUrlPatternData)
		contentFilteringStatus = true
	}

	if !data.ContentFilteringBlockedUrlCategoriesSettings.IsUnknown() {
		var blockedUrlCategorieData openApiClient.NetworksNetworkIdGroupPoliciesContentFilteringBlockedUrlCategories
		blockedUrlCategorieData.SetSettings(data.ContentFilteringBlockedUrlCategoriesSettings.ValueString())
		blockedUrlCategorieData.SetCategories(data.ContentFilteringBlockedUrlPatterns)
		contentFiltering.SetBlockedUrlCategories(blockedUrlCategorieData)
		contentFilteringStatus = true
	}

	if !data.ContentFilteringBlockedUrlPatternsSettings.IsUnknown() {
		var blockedUrlPatternData openApiClient.NetworksNetworkIdGroupPoliciesContentFilteringBlockedUrlPatterns
		blockedUrlPatternData.SetSettings(data.ContentFilteringBlockedUrlPatternsSettings.ValueString())
		blockedUrlPatternData.SetPatterns(data.ContentFilteringBlockedUrlPatterns)
		contentFiltering.SetBlockedUrlPatterns(blockedUrlPatternData)
		contentFilteringStatus = true
	}

	if contentFilteringStatus {
		createNetworkGroupPolicy.SetContentFiltering(contentFiltering)
	}

	inlineResp, httpResp, err := r.client.NetworksApi.CreateNetworkGroupPolicy(ctx, data.NetworkId.ValueString()).CreateNetworkGroupPolicy(createNetworkGroupPolicy).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// collect diagnostics
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
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

	data = extractHttpResponseGroupPolicyResource(ctx, inlineResp, data)

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
			"Failed to get resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// collect diagnostics
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
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

	data = extractHttpResponseGroupPolicyResource(ctx, inlineResp, data)

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

	updateNetworkGroupPolicy := *openApiClient.NewInlineObject88()
	if !data.SplashAuthSettings.IsUnknown() {
		updateNetworkGroupPolicy.SetSplashAuthSettings(data.SplashAuthSettings.ValueString())

	}

	if !data.BandwidthSettings.IsUnknown() {
		var bandwidth openApiClient.NetworksNetworkIdGroupPoliciesBandwidth
		bandwidth.SetSettings(data.BandwidthSettings.ValueString())
		var bandwidthLimits openApiClient.NetworksNetworkIdGroupPoliciesBandwidthBandwidthLimits
		if !data.BandwidthLimitUp.IsUnknown() {
			bandwidthLimits.SetLimitUp(int32(data.BandwidthLimitUp.ValueInt64()))
		}
		if !data.BandwidthLimitDown.IsUnknown() {
			bandwidthLimits.SetLimitDown(int32(data.BandwidthLimitDown.ValueInt64()))
		}
		bandwidth.SetBandwidthLimits(bandwidthLimits)
		updateNetworkGroupPolicy.SetBandwidth(bandwidth)
	}

	if len(data.BonjourForwardingRules) > 0 {
		var bonjourForwarding openApiClient.NetworksNetworkIdGroupPoliciesBonjourForwarding
		var bonjourForwardingRules []openApiClient.NetworksNetworkIdGroupPoliciesBonjourForwardingRules
		for _, attribute := range data.BonjourForwardingRules {
			var bonjourForwardingRule openApiClient.NetworksNetworkIdGroupPoliciesBonjourForwardingRules
			if !attribute.Description.IsUnknown() {
				bonjourForwardingRule.SetDescription(attribute.Description.ValueString())
			}
			bonjourForwardingRule.SetVlanId(attribute.VlanId.ValueString())
			bonjourForwardingRule.SetServices(attribute.Services)
			bonjourForwardingRules = append(bonjourForwardingRules, bonjourForwardingRule)
		}
		bonjourForwarding.SetRules(bonjourForwardingRules)
		if !data.BonjourForwardingSettings.IsUnknown() {
			bonjourForwarding.SetSettings(data.BonjourForwardingSettings.ValueString())
		}
		updateNetworkGroupPolicy.SetBonjourForwarding(bonjourForwarding)
	}

	var firewallAndTrafficShaping openApiClient.NetworksNetworkIdGroupPoliciesFirewallAndTrafficShaping

	if !data.FirewallAndTrafficShaping.Settings.IsUnknown() {
		firewallAndTrafficShaping.SetSettings(data.FirewallAndTrafficShaping.Settings.ValueString())
	}
	var l3s []openApiClient.NetworksNetworkIdGroupPoliciesFirewallAndTrafficShapingL3FirewallRules
	if len(data.FirewallAndTrafficShaping.L3FirewallRules) > 0 {

		for _, attribute := range data.FirewallAndTrafficShaping.L3FirewallRules {
			var l3 openApiClient.NetworksNetworkIdGroupPoliciesFirewallAndTrafficShapingL3FirewallRules
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

	if len(data.FirewallAndTrafficShaping.L7FirewallRules) > 0 {
		var l7s []openApiClient.NetworksNetworkIdGroupPoliciesFirewallAndTrafficShapingL7FirewallRules
		for _, attribute := range data.FirewallAndTrafficShaping.L7FirewallRules {
			var l7 openApiClient.NetworksNetworkIdGroupPoliciesFirewallAndTrafficShapingL7FirewallRules
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

	if len(data.FirewallAndTrafficShaping.TrafficShapingRules) > 0 {
		var tfs []openApiClient.NetworksNetworkIdGroupPoliciesFirewallAndTrafficShapingTrafficShapingRules
		for _, attribute := range data.FirewallAndTrafficShaping.TrafficShapingRules {
			var tf openApiClient.NetworksNetworkIdGroupPoliciesFirewallAndTrafficShapingTrafficShapingRules
			if !attribute.DscpTagValue.IsUnknown() {
				tf.SetDscpTagValue(int32(attribute.DscpTagValue.ValueInt64()))
			}
			if !attribute.PcpTagValue.IsUnknown() {
				tf.SetPcpTagValue(int32(attribute.PcpTagValue.ValueInt64()))
			}
			var perclientBamdWidthLimits openApiClient.NetworksNetworkIdApplianceTrafficShapingRulesPerClientBandwidthLimits

			if !attribute.PerClientBandwidthLimitsSettings.IsUnknown() {
				var bandwidthLimits openApiClient.NetworksNetworkIdApplianceTrafficShapingRulesPerClientBandwidthLimitsBandwidthLimits

				if attribute.PerClientBandwidthLimitsSettings.ValueString() != "network default" {

					if !attribute.BandwidthLimitDown.Int64Value.IsUnknown() {
						bandwidthLimits.SetLimitDown(int32(attribute.BandwidthLimitDown.ValueInt64()))
					}

					if !attribute.BandwidthLimitUp.Int64Value.IsUnknown() {
						bandwidthLimits.SetLimitUp(int32(attribute.BandwidthLimitUp.ValueInt64()))
					}
					perclientBamdWidthLimits.SetBandwidthLimits(bandwidthLimits)
				} else {

					if attribute.BandwidthLimitDown.ValueInt64() != jsontypes.Int64Null().ValueInt64() {
						resp.Diagnostics.AddError("Error:", "No need to add  band width limits for network default settings")
						return
					}

					if attribute.BandwidthLimitUp.ValueInt64() != jsontypes.Int64Null().ValueInt64() {
						resp.Diagnostics.AddError("Error:", "No need to add  band width limits for network default settings")
						return
					}
				}
				perclientBamdWidthLimits.SetSettings(attribute.PerClientBandwidthLimitsSettings.ValueString())
				tf.SetPerClientBandwidthLimits(perclientBamdWidthLimits)
			}
			var defs []openApiClient.NetworksNetworkIdApplianceTrafficShapingRulesDefinitions
			if len(attribute.Definitions) > 0 {
				for _, attribute := range attribute.Definitions {
					var def openApiClient.NetworksNetworkIdApplianceTrafficShapingRulesDefinitions
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

	if !data.SchedulingEnabled.IsUnknown() {
		var schedule openApiClient.NetworksNetworkIdGroupPoliciesScheduling
		schedule.SetEnabled(data.SchedulingEnabled.ValueBool())
		if !data.SchedulingFridayActive.IsUnknown() {
			var friday openApiClient.NetworksNetworkIdGroupPoliciesSchedulingFriday
			friday.SetActive(data.SchedulingFridayActive.ValueBool())
			friday.SetFrom(data.SchedulingFridayFrom.ValueString())
			friday.SetTo(data.SchedulingFridayTo.ValueString())
			schedule.SetFriday(friday)
		}
		if !data.SchedulingMondayActive.IsUnknown() {
			var monday openApiClient.NetworksNetworkIdGroupPoliciesSchedulingMonday
			monday.SetActive(data.SchedulingMondayActive.ValueBool())
			monday.SetFrom(data.SchedulingMondayFrom.ValueString())
			monday.SetTo(data.SchedulingMondayTo.ValueString())
			schedule.SetMonday(monday)
		}
		if !data.SchedulingTuesdayActive.IsUnknown() {
			var tuesday openApiClient.NetworksNetworkIdGroupPoliciesSchedulingTuesday
			tuesday.SetActive(data.SchedulingTuesdayActive.ValueBool())
			tuesday.SetFrom(data.SchedulingTuesdayFrom.ValueString())
			tuesday.SetTo(data.SchedulingTuesdayTo.ValueString())
			schedule.SetTuesday(tuesday)
		}
		if !data.SchedulingWednesdayActive.IsUnknown() {
			var wednesday openApiClient.NetworksNetworkIdGroupPoliciesSchedulingWednesday
			wednesday.SetActive(data.SchedulingWednesdayActive.ValueBool())
			wednesday.SetFrom(data.SchedulingWednesdayFrom.ValueString())
			wednesday.SetTo(data.SchedulingWednesdayTo.ValueString())
			schedule.SetWednesday(wednesday)
		}
		if !data.SchedulingThursdayActive.IsUnknown() {
			var thursday openApiClient.NetworksNetworkIdGroupPoliciesSchedulingThursday
			thursday.SetActive(data.SchedulingThursdayActive.ValueBool())
			thursday.SetFrom(data.SchedulingThursdayFrom.ValueString())
			thursday.SetTo(data.SchedulingThursdayTo.ValueString())
			schedule.SetThursday(thursday)
		}
		if !data.SchedulingSaturdayActive.IsUnknown() {
			var saturday openApiClient.NetworksNetworkIdGroupPoliciesSchedulingSaturday
			saturday.SetActive(data.SchedulingSaturdayActive.ValueBool())
			saturday.SetFrom(data.SchedulingSaturdayFrom.ValueString())
			saturday.SetTo(data.SchedulingSaturdayTo.ValueString())
			schedule.SetSaturday(saturday)
		}
		if !data.SchedulingSundayActive.IsUnknown() {
			var sunday openApiClient.NetworksNetworkIdGroupPoliciesSchedulingSunday
			sunday.SetActive(data.SchedulingSundayActive.ValueBool())
			sunday.SetFrom(data.SchedulingSundayFrom.ValueString())
			sunday.SetTo(data.SchedulingSundayTo.ValueString())
			schedule.SetSunday(sunday)
		}
		updateNetworkGroupPolicy.SetScheduling(schedule)
	}

	if !data.VlanTaggingSettings.IsUnknown() {
		if !data.VlanTaggingVlanId.IsUnknown() {
			var v openApiClient.NetworksNetworkIdGroupPoliciesVlanTagging
			v.SetSettings(data.VlanTaggingSettings.ValueString())
			v.SetVlanId(data.VlanTaggingVlanId.ValueString())
			updateNetworkGroupPolicy.SetVlanTagging(v)
		}
	}

	var contentFiltering openApiClient.NetworksNetworkIdGroupPoliciesContentFiltering
	contentFilteringStatus := false

	if !data.ContentFilteringAllowUrlPatternsSettings.IsUnknown() {
		var allowedUrlPatternData openApiClient.NetworksNetworkIdGroupPoliciesContentFilteringAllowedUrlPatterns
		allowedUrlPatternData.SetSettings(data.ContentFilteringAllowUrlPatternsSettings.ValueString())
		allowedUrlPatternData.SetPatterns(data.ContentFilteringAllowUrlPatterns)
		contentFiltering.SetAllowedUrlPatterns(allowedUrlPatternData)
		contentFilteringStatus = true
	}

	if !data.ContentFilteringBlockedUrlCategoriesSettings.IsUnknown() {
		var blockedUrlCategorieData openApiClient.NetworksNetworkIdGroupPoliciesContentFilteringBlockedUrlCategories
		blockedUrlCategorieData.SetSettings(data.ContentFilteringBlockedUrlCategoriesSettings.ValueString())
		blockedUrlCategorieData.SetCategories(data.ContentFilteringBlockedUrlPatterns)
		contentFiltering.SetBlockedUrlCategories(blockedUrlCategorieData)
		contentFilteringStatus = true
	}

	if !data.ContentFilteringBlockedUrlPatternsSettings.IsUnknown() {
		var blockedUrlPatternData openApiClient.NetworksNetworkIdGroupPoliciesContentFilteringBlockedUrlPatterns
		blockedUrlPatternData.SetSettings(data.ContentFilteringBlockedUrlPatternsSettings.ValueString())
		blockedUrlPatternData.SetPatterns(data.ContentFilteringBlockedUrlPatterns)
		contentFiltering.SetBlockedUrlPatterns(blockedUrlPatternData)
		contentFilteringStatus = true
	}

	if contentFilteringStatus {
		updateNetworkGroupPolicy.SetContentFiltering(contentFiltering)
	}

	inlineResp, httpResp, err := r.client.NetworksApi.UpdateNetworkGroupPolicy(ctx, data.NetworkId.ValueString(), data.GroupPolicyId.ValueString()).UpdateNetworkGroupPolicy(updateNetworkGroupPolicy).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// collect diagnostics
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
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

	data = extractHttpResponseGroupPolicyResource(ctx, inlineResp, data)

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
			"Failed to update resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// collect diagnostics
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
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

func extractHttpResponseGroupPolicyResource(ctx context.Context, inlineResp map[string]interface{}, data *NetworksGroupPolicyResourceModel) *NetworksGroupPolicyResourceModel {

	var outputdata OutPutData
	jsonData, _ := json.Marshal(inlineResp)
	json.Unmarshal(jsonData, &outputdata)
	data.Id = jsontypes.StringValue("example-id")
	data.GroupPolicyId = outputdata.GroupPolicyId
	data.BandwidthSettings = outputdata.Bandwidth.Settings
	data.BandwidthLimitUp = outputdata.Bandwidth.BandwidthLimits.LimitUp
	data.BandwidthLimitDown = outputdata.Bandwidth.BandwidthLimits.LimitDown
	data.SchedulingEnabled = outputdata.Scheduling.Enabled
	data.SchedulingFridayActive = outputdata.Scheduling.Friday.Active
	data.SchedulingFridayFrom = outputdata.Scheduling.Friday.From
	data.SchedulingFridayTo = outputdata.Scheduling.Friday.To
	data.SchedulingMondayActive = outputdata.Scheduling.Monday.Active
	data.SchedulingMondayFrom = outputdata.Scheduling.Monday.From
	data.SchedulingMondayTo = outputdata.Scheduling.Monday.To
	data.SchedulingTuesdayActive = outputdata.Scheduling.Tuesday.Active
	data.SchedulingTuesdayFrom = outputdata.Scheduling.Tuesday.From
	data.SchedulingTuesdayTo = outputdata.Scheduling.Tuesday.To
	data.SchedulingWednesdayActive = outputdata.Scheduling.Wednesday.Active
	data.SchedulingWednesdayFrom = outputdata.Scheduling.Wednesday.From
	data.SchedulingWednesdayTo = outputdata.Scheduling.Wednesday.To
	data.SchedulingThursdayActive = outputdata.Scheduling.Thursday.Active
	data.SchedulingThursdayFrom = outputdata.Scheduling.Thursday.From
	data.SchedulingThursdayTo = outputdata.Scheduling.Thursday.To
	data.SchedulingSaturdayActive = outputdata.Scheduling.Saturday.Active
	data.SchedulingSaturdayFrom = outputdata.Scheduling.Saturday.From
	data.SchedulingSaturdayTo = outputdata.Scheduling.Saturday.To
	data.SchedulingSundayActive = outputdata.Scheduling.Sunday.Active
	data.SchedulingSundayFrom = outputdata.Scheduling.Sunday.From
	data.SchedulingSundayTo = outputdata.Scheduling.Sunday.To
	data.VlanTaggingSettings = outputdata.VlanTagging.Settings
	data.VlanTaggingVlanId = outputdata.VlanTagging.VlanId
	data.BonjourForwardingSettings = jsontypes.StringValue(outputdata.BonjourForwarding.BonjourForwardingSettings)
	data.BonjourForwardingRules = outputdata.BonjourForwarding.BonjourForwardingRules
	var outputFirewallAndTrafficShaping OutputFirewallAndTrafficShaping
	if firewallAndTrafficShaping := inlineResp["firewallAndTrafficShaping"]; firewallAndTrafficShaping != nil {
		jsonData, _ = json.Marshal(inlineResp["firewallAndTrafficShaping"])
		json.Unmarshal(jsonData, &outputFirewallAndTrafficShaping)
		data.FirewallAndTrafficShaping.Settings = outputFirewallAndTrafficShaping.Settings
		if trafficShapingRules := inlineResp["firewallAndTrafficShaping"].(map[string]interface{})["trafficShapingRules"]; trafficShapingRules != nil {
			var outputTrafficShapingRule []OutputTrafficShapingRule
			jsonData, _ = json.Marshal(inlineResp["firewallAndTrafficShaping"].(map[string]interface{})["trafficShapingRules"])
			json.Unmarshal(jsonData, &outputTrafficShapingRule)
			for _, attribute := range outputTrafficShapingRule {
				var trafficShapingRule TrafficShapingRule
				trafficShapingRule.DscpTagValue = attribute.DscpTagValue
				trafficShapingRule.PcpTagValue = attribute.PcpTagValue
				trafficShapingRule.PerClientBandwidthLimitsSettings = attribute.PerClientBandwidthLimits.Settings
				trafficShapingRule.BandwidthLimitDown = attribute.PerClientBandwidthLimits.BandwidthLimits.LimitDown
				trafficShapingRule.BandwidthLimitUp = attribute.PerClientBandwidthLimits.BandwidthLimits.LimitUp
				if len(attribute.Definitions) > 0 {
					for _, attribute := range attribute.Definitions {
						var definition Definition
						definition.Type = attribute.Type
						definition.Value = attribute.Value
						trafficShapingRule.Definitions = append(trafficShapingRule.Definitions, definition)
					}
				} else {
					trafficShapingRule.Definitions = nil
				}
				data.FirewallAndTrafficShaping.TrafficShapingRules = nil
				data.FirewallAndTrafficShaping.TrafficShapingRules = append(data.FirewallAndTrafficShaping.TrafficShapingRules, trafficShapingRule)
			}
		} else {
			data.FirewallAndTrafficShaping.TrafficShapingRules = nil
		}
	}
	if contentFiltering := inlineResp["contentFiltering"]; contentFiltering != nil {
		var contentFilteringData ContentFiltering
		jsonData, _ = json.Marshal(inlineResp["contentFiltering"])
		json.Unmarshal(jsonData, &contentFilteringData)
		data.ContentFilteringAllowUrlPatternsSettings = contentFilteringData.AllowedUrlPatterns.Settings
		data.ContentFilteringBlockedUrlCategoriesSettings = contentFilteringData.BlockedUrlCategories.Settings
		data.ContentFilteringBlockedUrlPatternsSettings = contentFilteringData.BlockedUrlPatterns.Settings
		data.ContentFilteringAllowUrlPatterns = contentFilteringData.AllowedUrlPatterns.Patterns
		data.ContentFilteringBlockedUrlPatterns = contentFilteringData.BlockedUrlCategories.Categories
		data.ContentFilteringBlockedUrlCategories = contentFilteringData.BlockedUrlPatterns.Patterns

	} else {
		data.ContentFilteringAllowUrlPatternsSettings = jsontypes.StringNull()
		data.ContentFilteringBlockedUrlCategoriesSettings = jsontypes.StringNull()
		data.ContentFilteringBlockedUrlPatternsSettings = jsontypes.StringNull()
		data.ContentFilteringAllowUrlPatterns = make([]string, 0)
		data.ContentFilteringBlockedUrlPatterns = make([]string, 0)
		data.ContentFilteringBlockedUrlCategories = make([]string, 0)
	}

	return data

}
