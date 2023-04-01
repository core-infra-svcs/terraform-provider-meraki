package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
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
	Id                                           jsontypes.String     `tfsdk:"id"`
	NetworkId                                    jsontypes.String     `tfsdk:"network_id"`
	GroupPolicyId                                jsontypes.String     `tfsdk:"group_policy_id"`
	Name                                         jsontypes.String     `tfsdk:"name" json:"name"`
	SplashAuthSettings                           jsontypes.String     `tfsdk:"splash_auth_settings"`
	BandWidthSettings                            jsontypes.String     `tfsdk:"bandwidth_settings"`
	BandWidthLimitUp                             jsontypes.Int64      `tfsdk:"bandwidth_limit_up"`
	BandWidthLimitDown                           jsontypes.Int64      `tfsdk:"bandwidth_limit_down"`
	BonjourForwardingSettings                    jsontypes.String     `tfsdk:"bonjour_forwarding_settings"`
	BonjourForwardingRules                       []Rule               `tfsdk:"bonjour_forwarding_rules"`
	VlanTaggingSettings                          jsontypes.String     `tfsdk:"vlan_tagging_settings"`
	VlanTaggingVlanId                            jsontypes.String     `tfsdk:"vlan_tagging_vlan_id"`
	SchedulingEnabled                            jsontypes.Bool       `tfsdk:"scheduling_enabled"`
	SchedulingFridayActive                       jsontypes.Bool       `tfsdk:"scheduling_friday_active"`
	SchedulingFridayFrom                         jsontypes.String     `tfsdk:"scheduling_friday_from"`
	SchedulingFridayTo                           jsontypes.String     `tfsdk:"scheduling_friday_to"`
	SchedulingMondayActive                       jsontypes.Bool       `tfsdk:"scheduling_monday_active"`
	SchedulingMondayFrom                         jsontypes.String     `tfsdk:"scheduling_monday_from"`
	SchedulingMondayTo                           jsontypes.String     `tfsdk:"scheduling_monday_to"`
	SchedulingTuesdayActive                      jsontypes.Bool       `tfsdk:"scheduling_tuesday_active"`
	SchedulingTuesdayFrom                        jsontypes.String     `tfsdk:"scheduling_tuesday_from"`
	SchedulingTuesdayTo                          jsontypes.String     `tfsdk:"scheduling_tuesday_to"`
	SchedulingWednesdayActive                    jsontypes.Bool       `tfsdk:"scheduling_wednesday_active"`
	SchedulingWednesdayFrom                      jsontypes.String     `tfsdk:"scheduling_wednesday_from"`
	SchedulingWednesdayTo                        jsontypes.String     `tfsdk:"scheduling_wednesday_to"`
	SchedulingThursdayActive                     jsontypes.Bool       `tfsdk:"scheduling_thursday_active"`
	SchedulingThursdayFrom                       jsontypes.String     `tfsdk:"scheduling_thursday_from"`
	SchedulingThursdayTo                         jsontypes.String     `tfsdk:"scheduling_thursday_to"`
	SchedulingSaturdayActive                     jsontypes.Bool       `tfsdk:"scheduling_saturday_active"`
	SchedulingSaturdayFrom                       jsontypes.String     `tfsdk:"scheduling_saturday_from"`
	SchedulingSaturdayTo                         jsontypes.String     `tfsdk:"scheduling_saturday_to"`
	SchedulingSundayActive                       jsontypes.Bool       `tfsdk:"scheduling_sunday_active"`
	SchedulingSundayFrom                         jsontypes.String     `tfsdk:"scheduling_sunday_from"`
	SchedulingSundayTo                           jsontypes.String     `tfsdk:"scheduling_sunday_to"`
	FirewallAndTrafficShapingSettings            jsontypes.String     `tfsdk:"firewall_and_traffic_shaping_settings"`
	L3FirewallRules                              []L3FirewallRule     `tfsdk:"l3_firewall_rules"`
	L7FirewallRules                              []L7FirewallRule     `tfsdk:"l7_firewall_rules"`
	TrafficShapingRules                          []TrafficShapingRule `tfsdk:"traffic_shaping_rules"`
	ContentFilteringAllowedUrlPatternsSettings   jsontypes.String     `tfsdk:"content_filtering_allowed_url_patterns_settings"`
	UrlPatterns                                  []string             `tfsdk:"url_patterns"`
	ContentFilteringBlockedUrlCategoriesSettings jsontypes.String     `tfsdk:"content_filtering_blocked_url_categories_settings"`
	Categories                                   []string             `tfsdk:"categories"`
	ContentFilteringBlockedUrlPatternsSettings   jsontypes.String     `tfsdk:"content_filtering_blocked_url_patterns_settings"`
	BlockedUrlPatterns                           []string             `tfsdk:"blocked_url_patterns"`
}

type NetworksGroupPolicyResourceModelData struct {
	Name                      jsontypes.String          `json:"name"`
	SplashAuthSettings        jsontypes.String          `json:"splashAuthSettings"`
	BandWidth                 BandWidth                 `json:"bandwidth"`
	GroupPolicyId             jsontypes.String          `json:"groupPolicyId"`
	BonjourForwarding         BonjourForwarding         `json:"bonjourForwarding"`
	VlanTagging               VlanTagging               `json:"vlanTagging"`
	Scheduling                Scheduling                `json:"scheduling"`
	FirewallAndTrafficShaping FirewallAndTrafficShaping `json:"firewall_and_traffic_shaping"`
	ContentFiltering          ContentFiltering          `tfsdk:"content_filtering"`
}

type FirewallAndTrafficShaping struct {
	Settings            jsontypes.String      `tfsdk:"settings"`
	L3FirewallRules     []L3FirewallRule      `tfsdk:"l3_firewall_rules"`
	L7FirewallRules     []L7FirewallRule      `tfsdk:"l7_firewall_rules"`
	TrafficShapingRules []TrafficShapingRules `tfsdk:"traffic_shaping_rules"`
}

type L3FirewallRule struct {
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

type TrafficShapingRules struct {
	DscpTagValue             jsontypes.Int64          `tfsdk:"dscp_tag_value"`
	PcpTagValue              jsontypes.Int64          `tfsdk:"pcp_tag_value"`
	Priority                 jsontypes.String         `tfsdk:"priority"`
	PerClientBandwidthLimits PerClientBandwidthLimits `tfsdk:"per_client_bandwidth_limits"`
	Definitions              []Definition             `tfsdk:"definitions"`
}

type TrafficShapingRule struct {
	DscpTagValue                     jsontypes.Int64  `tfsdk:"dscp_tag_value"`
	PcpTagValue                      jsontypes.Int64  `tfsdk:"pcp_tag_value"`
	Priority                         jsontypes.String `tfsdk:"priority"`
	PerClientBandwidthLimitsSettings jsontypes.String `tfsdk:"per_client_bandwidth_limits_settings"`
	BandwidthLimitDown               jsontypes.Int64  `tfsdk:"bandwidth_limit_down"`
	BandwidthLimitUp                 jsontypes.Int64  `tfsdk:"bandwidth_limit_up"`
	Definitions                      []Definition     `tfsdk:"definitions"`
}

type PerClientBandwidthLimits struct {
	Settings        jsontypes.String `tfsdk:"settings"`
	BandwidthLimits BandwidthLimits  `tfsdk:"bandwidth_limits"`
}

type Definition struct {
	Value jsontypes.String `tfsdk:"value"`
	Type  jsontypes.String `tfsdk:"type"`
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

type BandWidth struct {
	Settings        jsontypes.String `json:"settings"`
	BandwidthLimits BandwidthLimits  `json:"bandwidthLimits"`
}
type BandwidthLimits struct {
	LimitDown jsontypes.Int64 `json:"limitDown"`
	LimitUp   jsontypes.Int64 `json:"limitUp"`
}

type BonjourForwarding struct {
	Settings jsontypes.String `tfsdk:"settings"`
	Rules    []Rule           `tfsdk:"rules"`
}

type Rule struct {
	Description jsontypes.String `tfsdk:"description"`
	VlanId      jsontypes.String `tfsdk:"vlan_id"`
	Services    []string         `tfsdk:"services"`
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
				MarkdownDescription: "Settings Bandwidth",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"bandwidth_limit_down": schema.Int64Attribute{
				MarkdownDescription: "The maximum download limit (integer, in Kbps).",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.Int64Type,
			},
			"bandwidth_limit_up": schema.Int64Attribute{
				MarkdownDescription: "The maximum upload limit (integer, in Kbps).",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.Int64Type,
			},
			"bonjour_forwarding_settings": schema.StringAttribute{
				MarkdownDescription: "How Bonjour rules are applied. Can be 'network default', 'ignore' or 'custom'.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"bonjour_forwarding_rules": schema.SetNestedAttribute{
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
			"firewall_and_traffic_shaping_settings": schema.StringAttribute{
				MarkdownDescription: "How firewall and traffic shaping rules are enforced. Can be 'network default', 'ignore' or 'custom'.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"l3_firewall_rules": schema.SetNestedAttribute{
				Optional: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"comment": schema.StringAttribute{
							MarkdownDescription: "Description of the rule (optional)",
							Optional:            true,

							CustomType: jsontypes.StringType,
						},
						"dest_cidr": schema.StringAttribute{
							MarkdownDescription: "Destination IP address (in IP or CIDR notation), a fully-qualified domain name (FQDN, if your network supports it) or 'any'.",
							Optional:            true,

							CustomType: jsontypes.StringType,
						},
						"dest_port": schema.StringAttribute{
							MarkdownDescription: "Destination port (integer in the range 1-65535), a port range (e.g. 8080-9090), or 'any'",
							Optional:            true,

							CustomType: jsontypes.StringType,
						},
						"policy": schema.StringAttribute{
							MarkdownDescription: "'allow' or 'deny' traffic specified by this rule",
							Optional:            true,

							CustomType: jsontypes.StringType,
						},
						"protocol": schema.StringAttribute{
							MarkdownDescription: "The type of protocol (must be 'tcp', 'udp', 'icmp', 'icmp6' or 'any')",
							Optional:            true,

							CustomType: jsontypes.StringType,
						},
					},
				},
			},
			"l7_firewall_rules": schema.SetNestedAttribute{
				Optional: true,
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
						"priority": schema.StringAttribute{
							MarkdownDescription: "A string, indicating the priority level for packets bound to your rule. Can be 'low', 'normal' or 'high'.",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"per_client_bandwidth_limits_settings": schema.StringAttribute{
							MarkdownDescription: "How bandwidth limits are applied by your rule. Can be one of 'network default', 'ignore' or 'custom'.",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"bandwidth_limit_down": schema.Int64Attribute{
							MarkdownDescription: "The maximum download limit (integer, in Kbps).",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.Int64Type,
						},
						"bandwidth_limit_up": schema.Int64Attribute{
							MarkdownDescription: "The maximum upload limit (integer, in Kbps).",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.Int64Type,
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
			"content_filtering_allowed_url_patterns_settings": schema.StringAttribute{
				MarkdownDescription: "How URL patterns are applied. Can be 'network default', 'append' or 'override'.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"url_patterns": schema.SetAttribute{
				MarkdownDescription: "A list of URL patterns that are allowed",
				CustomType:          jsontypes.SetType[jsontypes.String](),
				Optional:            true,
			},
			"content_filtering_blocked_url_categories_settings": schema.StringAttribute{
				MarkdownDescription: "How URL categories are applied. Can be 'network default', 'append' or 'override'.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"categories": schema.SetAttribute{
				MarkdownDescription: "A list of URL categories to block",
				CustomType:          jsontypes.SetType[jsontypes.String](),

				Optional: true,
			},
			"content_filtering_blocked_url_patterns_settings": schema.StringAttribute{
				MarkdownDescription: "How URL categories are applied. Can be 'network default', 'append' or 'override'.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"blocked_url_patterns": schema.SetAttribute{
				MarkdownDescription: "A list of URL patterns that are blocked",
				CustomType:          jsontypes.SetType[jsontypes.String](),
				Optional:            true,
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
	var structmap *NetworksGroupPolicyResourceModelData

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	createNetworkGroupPolicy := *openApiClient.NewInlineObject87(data.Name.ValueString())

	if !data.SplashAuthSettings.IsUnknown() {
		createNetworkGroupPolicy.SetSplashAuthSettings(data.SplashAuthSettings.ValueString())

	} else {
		data.SplashAuthSettings = jsontypes.StringNull()
	}

	if !data.BandWidthSettings.IsUnknown() {
		var a openApiClient.NetworksNetworkIdGroupPoliciesBandwidth
		a.SetSettings(data.BandWidthSettings.ValueString())
		createNetworkGroupPolicy.SetBandwidth(a)

	} else {
		data.BandWidthSettings = jsontypes.StringNull()

	}

	if !data.BandWidthLimitUp.IsUnknown() && !data.BandWidthLimitDown.IsUnknown() {

		var a openApiClient.NetworksNetworkIdGroupPoliciesBandwidth
		var v openApiClient.NetworksNetworkIdGroupPoliciesBandwidthBandwidthLimits
		v.SetLimitUp(int32(data.BandWidthLimitUp.ValueInt64()))
		v.SetLimitDown(int32(data.BandWidthLimitDown.ValueInt64()))
		a.SetBandwidthLimits(v)
		createNetworkGroupPolicy.SetBandwidth(a)

	} else {
		data.BandWidthLimitUp = jsontypes.Int64Null()
		data.BandWidthLimitDown = jsontypes.Int64Null()
	}

	if !data.BonjourForwardingSettings.IsUnknown() {
		var b openApiClient.NetworksNetworkIdGroupPoliciesBonjourForwarding
		b.SetSettings(data.BonjourForwardingSettings.ValueString())
		createNetworkGroupPolicy.SetBonjourForwarding(b)

	} else {
		data.BonjourForwardingSettings = jsontypes.StringNull()
	}

	if len(data.BonjourForwardingRules) > 0 {
		var b openApiClient.NetworksNetworkIdGroupPoliciesBonjourForwarding
		var r []openApiClient.NetworksNetworkIdGroupPoliciesBonjourForwardingRules
		for _, attribute := range data.BonjourForwardingRules {
			var a openApiClient.NetworksNetworkIdGroupPoliciesBonjourForwardingRules
			if !attribute.Description.IsUnknown() {
				a.SetDescription(attribute.Description.ValueString())
			}
			if !attribute.VlanId.IsUnknown() {
				a.SetVlanId(attribute.VlanId.ValueString())
			}
			if len(attribute.Services) > 0 {
				a.SetServices(attribute.Services)
			}
			r = append(r, a)
		}
		b.SetRules(r)
		createNetworkGroupPolicy.SetBonjourForwarding(b)
	} else {
		data.BonjourForwardingRules = nil
	}

	if !data.VlanTaggingSettings.IsUnknown() && !data.VlanTaggingVlanId.IsUnknown() {
		var v openApiClient.NetworksNetworkIdGroupPoliciesVlanTagging
		v.SetSettings(data.VlanTaggingSettings.ValueString())
		v.SetVlanId(data.VlanTaggingVlanId.ValueString())
		createNetworkGroupPolicy.SetVlanTagging(v)
	} else {
		data.VlanTaggingSettings = jsontypes.StringNull()
		data.VlanTaggingVlanId = jsontypes.StringNull()
	}

	if !data.SchedulingEnabled.IsUnknown() {
		var s openApiClient.NetworksNetworkIdGroupPoliciesScheduling
		s.SetEnabled(data.SchedulingEnabled.ValueBool())
		if !data.SchedulingFridayActive.IsUnknown() {
			var f openApiClient.NetworksNetworkIdGroupPoliciesSchedulingFriday
			f.SetActive(data.SchedulingFridayActive.ValueBool())
			f.SetFrom(data.SchedulingFridayFrom.ValueString())
			f.SetTo(data.SchedulingFridayTo.ValueString())
			s.SetFriday(f)
			createNetworkGroupPolicy.SetScheduling(s)
		} else {
			data.SchedulingFridayActive = jsontypes.BoolNull()
		}
		if !data.SchedulingMondayActive.IsUnknown() {
			var f openApiClient.NetworksNetworkIdGroupPoliciesSchedulingMonday
			f.SetActive(data.SchedulingMondayActive.ValueBool())
			f.SetFrom(data.SchedulingMondayFrom.ValueString())
			f.SetTo(data.SchedulingMondayTo.ValueString())
			s.SetMonday(f)
			createNetworkGroupPolicy.SetScheduling(s)
		} else {
			data.SchedulingMondayActive = jsontypes.BoolNull()
		}
		if !data.SchedulingTuesdayActive.IsUnknown() {
			var f openApiClient.NetworksNetworkIdGroupPoliciesSchedulingTuesday
			f.SetActive(data.SchedulingTuesdayActive.ValueBool())
			f.SetFrom(data.SchedulingTuesdayFrom.ValueString())
			f.SetTo(data.SchedulingTuesdayTo.ValueString())
			s.SetTuesday(f)
			createNetworkGroupPolicy.SetScheduling(s)
		} else {
			data.SchedulingTuesdayActive = jsontypes.BoolNull()
		}
		if !data.SchedulingWednesdayActive.IsUnknown() {
			var f openApiClient.NetworksNetworkIdGroupPoliciesSchedulingWednesday
			f.SetActive(data.SchedulingWednesdayActive.ValueBool())
			f.SetFrom(data.SchedulingWednesdayFrom.ValueString())
			f.SetTo(data.SchedulingWednesdayTo.ValueString())
			s.SetWednesday(f)
			createNetworkGroupPolicy.SetScheduling(s)
		} else {
			data.SchedulingWednesdayActive = jsontypes.BoolNull()
		}
		if !data.SchedulingThursdayActive.IsUnknown() {
			var f openApiClient.NetworksNetworkIdGroupPoliciesSchedulingThursday
			f.SetActive(data.SchedulingThursdayActive.ValueBool())
			f.SetFrom(data.SchedulingThursdayFrom.ValueString())
			f.SetTo(data.SchedulingThursdayTo.ValueString())
			s.SetThursday(f)
			createNetworkGroupPolicy.SetScheduling(s)
		} else {
			data.SchedulingThursdayActive = jsontypes.BoolNull()
		}
		if !data.SchedulingSaturdayActive.IsUnknown() {
			var f openApiClient.NetworksNetworkIdGroupPoliciesSchedulingSaturday
			f.SetActive(data.SchedulingSaturdayActive.ValueBool())
			f.SetFrom(data.SchedulingSaturdayFrom.ValueString())
			f.SetTo(data.SchedulingSaturdayTo.ValueString())
			s.SetSaturday(f)
			createNetworkGroupPolicy.SetScheduling(s)
		} else {
			data.SchedulingSaturdayActive = jsontypes.BoolNull()
		}
		if !data.SchedulingSundayActive.IsUnknown() {
			var f openApiClient.NetworksNetworkIdGroupPoliciesSchedulingSunday
			f.SetActive(data.SchedulingSundayActive.ValueBool())
			f.SetFrom(data.SchedulingSundayFrom.ValueString())
			f.SetTo(data.SchedulingSundayTo.ValueString())
			s.SetSunday(f)
			createNetworkGroupPolicy.SetScheduling(s)
		} else {
			data.SchedulingSundayActive = jsontypes.BoolNull()
		}

	} else {
		data.SchedulingEnabled = jsontypes.BoolNull()
	}

	var f openApiClient.NetworksNetworkIdGroupPoliciesFirewallAndTrafficShaping
	if !data.FirewallAndTrafficShapingSettings.IsUnknown() {

		f.SetSettings(data.FirewallAndTrafficShapingSettings.ValueString())

	} else {
		data.FirewallAndTrafficShapingSettings = jsontypes.StringNull()
	}

	if len(data.L3FirewallRules) > 0 {
		var l3s []openApiClient.NetworksNetworkIdGroupPoliciesFirewallAndTrafficShapingL3FirewallRules
		for _, attribute := range data.L3FirewallRules {
			var l3 openApiClient.NetworksNetworkIdGroupPoliciesFirewallAndTrafficShapingL3FirewallRules
			l3.SetComment(attribute.Comment.ValueString())
			l3.SetDestCidr(attribute.DestCidr.ValueString())
			l3.SetDestPort(attribute.DestPort.ValueString())
			l3.SetPolicy(attribute.Policy.ValueString())
			l3.SetProtocol(attribute.Protocol.ValueString())
			l3s = append(l3s, l3)
		}
		f.SetL3FirewallRules(l3s)
	} else {
		data.L3FirewallRules = nil
	}
	if len(data.L7FirewallRules) > 0 {
		var l7s []openApiClient.NetworksNetworkIdGroupPoliciesFirewallAndTrafficShapingL7FirewallRules
		for _, attribute := range data.L7FirewallRules {
			var l7 openApiClient.NetworksNetworkIdGroupPoliciesFirewallAndTrafficShapingL7FirewallRules

			l7.SetValue(attribute.Value.ValueString())
			l7.SetPolicy(attribute.Policy.ValueString())
			l7.SetType(attribute.Type.ValueString())
			l7s = append(l7s, l7)
		}
		f.SetL7FirewallRules(l7s)
	} else {
		data.L7FirewallRules = nil
	}

	if len(data.TrafficShapingRules) > 0 {
		var tfs []openApiClient.NetworksNetworkIdGroupPoliciesFirewallAndTrafficShapingTrafficShapingRules
		for _, attribute := range data.TrafficShapingRules {
			var tf openApiClient.NetworksNetworkIdGroupPoliciesFirewallAndTrafficShapingTrafficShapingRules
			tf.SetDscpTagValue(int32(attribute.DscpTagValue.ValueInt64()))
			tf.SetPcpTagValue(int32(attribute.PcpTagValue.ValueInt64()))
			tf.SetPriority(attribute.Priority.ValueString())
			var perclientBamdWidthLimits openApiClient.NetworksNetworkIdApplianceTrafficShapingRulesPerClientBandwidthLimits
			var bandwidthLimits openApiClient.NetworksNetworkIdApplianceTrafficShapingRulesPerClientBandwidthLimitsBandwidthLimits
			bandwidthLimits.SetLimitDown(int32(attribute.BandwidthLimitDown.ValueInt64()))
			bandwidthLimits.SetLimitUp(int32(attribute.BandwidthLimitUp.ValueInt64()))
			perclientBamdWidthLimits.SetBandwidthLimits(bandwidthLimits)
			tf.SetPerClientBandwidthLimits(perclientBamdWidthLimits)
			var defs []openApiClient.NetworksNetworkIdApplianceTrafficShapingRulesDefinitions
			for _, attribute := range attribute.Definitions {
				var def openApiClient.NetworksNetworkIdApplianceTrafficShapingRulesDefinitions
				def.SetType(attribute.Type.ValueString())
				def.SetValue(attribute.Value.ValueString())
				defs = append(defs, def)
			}
			tf.SetDefinitions(defs)
			tfs = append(tfs, tf)
		}
		f.SetTrafficShapingRules(tfs)
	} else {
		data.TrafficShapingRules = nil
	}

	createNetworkGroupPolicy.SetFirewallAndTrafficShaping(f)

	var c openApiClient.NetworksNetworkIdGroupPoliciesContentFiltering
	if !data.ContentFilteringAllowedUrlPatternsSettings.IsUnknown() || len(data.UrlPatterns) > 0 {
		var aup openApiClient.NetworksNetworkIdGroupPoliciesContentFilteringAllowedUrlPatterns
		aup.SetSettings(data.ContentFilteringAllowedUrlPatternsSettings.ValueString())
		aup.SetPatterns(data.UrlPatterns)
		c.SetAllowedUrlPatterns(aup)

	} else {
		data.ContentFilteringAllowedUrlPatternsSettings = jsontypes.StringNull()
		data.UrlPatterns = nil
	}
	if !data.ContentFilteringBlockedUrlCategoriesSettings.IsUnknown() || len(data.Categories) > 0 {
		var buc openApiClient.NetworksNetworkIdGroupPoliciesContentFilteringBlockedUrlCategories
		buc.SetSettings(data.ContentFilteringBlockedUrlCategoriesSettings.ValueString())
		buc.SetCategories(data.Categories)
		c.SetBlockedUrlCategories(buc)

	} else {
		data.ContentFilteringBlockedUrlCategoriesSettings = jsontypes.StringNull()
		data.Categories = nil
	}
	if !data.ContentFilteringBlockedUrlPatternsSettings.IsUnknown() || len(data.BlockedUrlPatterns) > 0 {
		var bup openApiClient.NetworksNetworkIdGroupPoliciesContentFilteringBlockedUrlPatterns
		bup.SetSettings(data.ContentFilteringBlockedUrlPatternsSettings.ValueString())
		bup.SetPatterns(data.BlockedUrlPatterns)
		c.SetBlockedUrlPatterns(bup)

	} else {
		data.ContentFilteringBlockedUrlPatternsSettings = jsontypes.StringNull()
		data.BlockedUrlPatterns = nil
	}
	createNetworkGroupPolicy.SetContentFiltering(c)

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

	// Save data into Terraform state
	jsonData, _ := json.Marshal(inlineResp)
	json.Unmarshal(jsonData, &structmap)

	data.Id = jsontypes.StringValue("example-id")
	data.Name = structmap.Name
	data.GroupPolicyId = structmap.GroupPolicyId
	data.BandWidthLimitDown = structmap.BandWidth.BandwidthLimits.LimitDown
	data.BandWidthLimitUp = structmap.BandWidth.BandwidthLimits.LimitUp
	data.BandWidthSettings = structmap.BandWidth.Settings
	data.SplashAuthSettings = structmap.SplashAuthSettings
	data.BandWidthLimitDown = structmap.BandWidth.BandwidthLimits.LimitDown
	data.BandWidthLimitUp = structmap.BandWidth.BandwidthLimits.LimitUp
	data.BonjourForwardingRules = structmap.BonjourForwarding.Rules
	data.BonjourForwardingSettings = structmap.BonjourForwarding.Settings
	data.VlanTaggingSettings = structmap.VlanTagging.Settings
	data.VlanTaggingVlanId = structmap.VlanTagging.VlanId
	data.SchedulingEnabled = structmap.Scheduling.Enabled
	data.SchedulingFridayActive = structmap.Scheduling.Friday.Active
	data.SchedulingFridayFrom = structmap.Scheduling.Friday.From
	data.SchedulingFridayTo = structmap.Scheduling.Friday.To
	data.SchedulingMondayActive = structmap.Scheduling.Monday.Active
	data.SchedulingMondayFrom = structmap.Scheduling.Monday.From
	data.SchedulingMondayTo = structmap.Scheduling.Monday.To
	data.SchedulingTuesdayActive = structmap.Scheduling.Tuesday.Active
	data.SchedulingTuesdayFrom = structmap.Scheduling.Tuesday.From
	data.SchedulingTuesdayTo = structmap.Scheduling.Tuesday.To
	data.SchedulingWednesdayActive = structmap.Scheduling.Wednesday.Active
	data.SchedulingWednesdayFrom = structmap.Scheduling.Wednesday.From
	data.SchedulingWednesdayTo = structmap.Scheduling.Wednesday.To
	data.SchedulingThursdayActive = structmap.Scheduling.Thursday.Active
	data.SchedulingThursdayFrom = structmap.Scheduling.Thursday.From
	data.SchedulingThursdayTo = structmap.Scheduling.Thursday.To
	data.SchedulingSaturdayActive = structmap.Scheduling.Saturday.Active
	data.SchedulingSaturdayFrom = structmap.Scheduling.Saturday.From
	data.SchedulingSaturdayTo = structmap.Scheduling.Saturday.To
	data.SchedulingSundayActive = structmap.Scheduling.Sunday.Active
	data.SchedulingSundayFrom = structmap.Scheduling.Sunday.From
	data.SchedulingSundayTo = structmap.Scheduling.Sunday.To
	if len(structmap.FirewallAndTrafficShaping.L3FirewallRules) > 0 {

		data.L3FirewallRules = structmap.FirewallAndTrafficShaping.L3FirewallRules
	}
	if len(structmap.FirewallAndTrafficShaping.L7FirewallRules) > 0 {
		data.L7FirewallRules = structmap.FirewallAndTrafficShaping.L7FirewallRules
	}
	data.FirewallAndTrafficShapingSettings = structmap.FirewallAndTrafficShaping.Settings
	for _, attribute := range structmap.FirewallAndTrafficShaping.TrafficShapingRules {
		var a TrafficShapingRule
		a.Definitions = attribute.Definitions
		a.BandwidthLimitDown = attribute.PerClientBandwidthLimits.BandwidthLimits.LimitDown
		a.BandwidthLimitUp = attribute.PerClientBandwidthLimits.BandwidthLimits.LimitUp
		a.Priority = attribute.Priority
		a.DscpTagValue = attribute.DscpTagValue
		a.PcpTagValue = attribute.PcpTagValue
		data.TrafficShapingRules = append(data.TrafficShapingRules, a)
	}
	data.ContentFilteringAllowedUrlPatternsSettings = structmap.ContentFiltering.AllowedUrlPatterns.Settings
	if len(structmap.ContentFiltering.AllowedUrlPatterns.Patterns) > 0 {
		data.UrlPatterns = structmap.ContentFiltering.AllowedUrlPatterns.Patterns
	} else {
		data.UrlPatterns = nil
	}
	data.ContentFilteringBlockedUrlCategoriesSettings = structmap.ContentFiltering.BlockedUrlCategories.Settings
	if len(structmap.ContentFiltering.BlockedUrlCategories.Categories) > 0 {
		data.Categories = structmap.ContentFiltering.BlockedUrlCategories.Categories
	} else {
		data.Categories = nil
	}
	data.ContentFilteringBlockedUrlPatternsSettings = structmap.ContentFiltering.BlockedUrlPatterns.Settings
	if len(structmap.ContentFiltering.BlockedUrlPatterns.Patterns) > 0 {

		data.BlockedUrlPatterns = structmap.ContentFiltering.BlockedUrlPatterns.Patterns
	} else {
		data.BlockedUrlPatterns = nil
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")
}

func (r *NetworksGroupPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetworksGroupPolicyResourceModel
	var structmap *NetworksGroupPolicyResourceModelData

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

	// Save data into Terraform state
	// Save data into Terraform state
	jsonData, _ := json.Marshal(inlineResp)
	json.Unmarshal(jsonData, &structmap)

	data.Id = jsontypes.StringValue("example-id")
	data.Name = structmap.Name
	data.GroupPolicyId = structmap.GroupPolicyId
	data.BandWidthLimitDown = structmap.BandWidth.BandwidthLimits.LimitDown
	data.BandWidthLimitUp = structmap.BandWidth.BandwidthLimits.LimitUp
	data.BandWidthSettings = structmap.BandWidth.Settings
	data.SplashAuthSettings = structmap.SplashAuthSettings
	data.BandWidthLimitDown = structmap.BandWidth.BandwidthLimits.LimitDown
	data.BandWidthLimitUp = structmap.BandWidth.BandwidthLimits.LimitUp
	data.BonjourForwardingRules = structmap.BonjourForwarding.Rules
	data.BonjourForwardingSettings = structmap.BonjourForwarding.Settings
	data.VlanTaggingSettings = structmap.VlanTagging.Settings
	data.VlanTaggingVlanId = structmap.VlanTagging.VlanId
	data.SchedulingEnabled = structmap.Scheduling.Enabled
	data.SchedulingFridayActive = structmap.Scheduling.Friday.Active
	data.SchedulingFridayFrom = structmap.Scheduling.Friday.From
	data.SchedulingFridayTo = structmap.Scheduling.Friday.To
	data.SchedulingMondayActive = structmap.Scheduling.Monday.Active
	data.SchedulingMondayFrom = structmap.Scheduling.Monday.From
	data.SchedulingMondayTo = structmap.Scheduling.Monday.To
	data.SchedulingTuesdayActive = structmap.Scheduling.Tuesday.Active
	data.SchedulingTuesdayFrom = structmap.Scheduling.Tuesday.From
	data.SchedulingTuesdayTo = structmap.Scheduling.Tuesday.To
	data.SchedulingWednesdayActive = structmap.Scheduling.Wednesday.Active
	data.SchedulingWednesdayFrom = structmap.Scheduling.Wednesday.From
	data.SchedulingWednesdayTo = structmap.Scheduling.Wednesday.To
	data.SchedulingThursdayActive = structmap.Scheduling.Thursday.Active
	data.SchedulingThursdayFrom = structmap.Scheduling.Thursday.From
	data.SchedulingThursdayTo = structmap.Scheduling.Thursday.To
	data.SchedulingSaturdayActive = structmap.Scheduling.Saturday.Active
	data.SchedulingSaturdayFrom = structmap.Scheduling.Saturday.From
	data.SchedulingSaturdayTo = structmap.Scheduling.Saturday.To
	data.SchedulingSundayActive = structmap.Scheduling.Sunday.Active
	data.SchedulingSundayFrom = structmap.Scheduling.Sunday.From
	data.SchedulingSundayTo = structmap.Scheduling.Sunday.To
	if len(structmap.FirewallAndTrafficShaping.L3FirewallRules) > 0 {

		data.L3FirewallRules = structmap.FirewallAndTrafficShaping.L3FirewallRules
	}
	if len(structmap.FirewallAndTrafficShaping.L7FirewallRules) > 0 {
		data.L7FirewallRules = structmap.FirewallAndTrafficShaping.L7FirewallRules
	}
	data.FirewallAndTrafficShapingSettings = structmap.FirewallAndTrafficShaping.Settings
	for _, attribute := range structmap.FirewallAndTrafficShaping.TrafficShapingRules {
		var a TrafficShapingRule
		a.Definitions = attribute.Definitions
		a.BandwidthLimitDown = attribute.PerClientBandwidthLimits.BandwidthLimits.LimitDown
		a.BandwidthLimitUp = attribute.PerClientBandwidthLimits.BandwidthLimits.LimitUp
		a.Priority = attribute.Priority
		a.DscpTagValue = attribute.DscpTagValue
		a.PcpTagValue = attribute.PcpTagValue
		data.TrafficShapingRules = append(data.TrafficShapingRules, a)
	}
	data.ContentFilteringAllowedUrlPatternsSettings = structmap.ContentFiltering.AllowedUrlPatterns.Settings
	if len(structmap.ContentFiltering.AllowedUrlPatterns.Patterns) > 0 {
		data.UrlPatterns = structmap.ContentFiltering.AllowedUrlPatterns.Patterns
	} else {
		data.UrlPatterns = nil
	}
	data.ContentFilteringBlockedUrlCategoriesSettings = structmap.ContentFiltering.BlockedUrlCategories.Settings
	if len(structmap.ContentFiltering.BlockedUrlCategories.Categories) > 0 {
		data.Categories = structmap.ContentFiltering.BlockedUrlCategories.Categories
	} else {
		data.Categories = nil
	}
	data.ContentFilteringBlockedUrlPatternsSettings = structmap.ContentFiltering.BlockedUrlPatterns.Settings
	if len(structmap.ContentFiltering.BlockedUrlPatterns.Patterns) > 0 {

		data.BlockedUrlPatterns = structmap.ContentFiltering.BlockedUrlPatterns.Patterns
	} else {
		data.BlockedUrlPatterns = nil
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *NetworksGroupPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data *NetworksGroupPolicyResourceModel
	var structmap *NetworksGroupPolicyResourceModelData

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworkGroupPolicy := *openApiClient.NewInlineObject88()

	if !data.SplashAuthSettings.IsUnknown() {
		updateNetworkGroupPolicy.SetSplashAuthSettings(data.SplashAuthSettings.ValueString())

	} else {
		data.SplashAuthSettings = jsontypes.StringNull()
	}

	if !data.BandWidthSettings.IsUnknown() {
		var a openApiClient.NetworksNetworkIdGroupPoliciesBandwidth
		a.SetSettings(data.BandWidthSettings.ValueString())
		updateNetworkGroupPolicy.SetBandwidth(a)

	} else {
		data.BandWidthSettings = jsontypes.StringNull()

	}

	if !data.BandWidthLimitUp.IsUnknown() && !data.BandWidthLimitDown.IsUnknown() {

		var a openApiClient.NetworksNetworkIdGroupPoliciesBandwidth
		var v openApiClient.NetworksNetworkIdGroupPoliciesBandwidthBandwidthLimits
		v.SetLimitUp(int32(data.BandWidthLimitUp.ValueInt64()))
		v.SetLimitDown(int32(data.BandWidthLimitDown.ValueInt64()))
		a.SetBandwidthLimits(v)
		updateNetworkGroupPolicy.SetBandwidth(a)

	} else {
		data.BandWidthLimitUp = jsontypes.Int64Null()
		data.BandWidthLimitDown = jsontypes.Int64Null()
	}

	if !data.BonjourForwardingSettings.IsUnknown() {
		var b openApiClient.NetworksNetworkIdGroupPoliciesBonjourForwarding
		b.SetSettings(data.BonjourForwardingSettings.ValueString())
		updateNetworkGroupPolicy.SetBonjourForwarding(b)

	} else {
		data.BonjourForwardingSettings = jsontypes.StringNull()
	}

	if len(data.BonjourForwardingRules) > 0 {
		var b openApiClient.NetworksNetworkIdGroupPoliciesBonjourForwarding
		var r []openApiClient.NetworksNetworkIdGroupPoliciesBonjourForwardingRules
		for _, attribute := range data.BonjourForwardingRules {
			var a openApiClient.NetworksNetworkIdGroupPoliciesBonjourForwardingRules
			if !attribute.Description.IsUnknown() {
				a.SetDescription(attribute.Description.ValueString())
			}
			if !attribute.VlanId.IsUnknown() {
				a.SetVlanId(attribute.VlanId.ValueString())
			}
			if len(attribute.Services) > 0 {
				a.SetServices(attribute.Services)
			}
			r = append(r, a)
		}
		b.SetRules(r)
		updateNetworkGroupPolicy.SetBonjourForwarding(b)
	} else {
		data.BonjourForwardingRules = nil
	}

	if !data.VlanTaggingSettings.IsUnknown() && !data.VlanTaggingVlanId.IsUnknown() {
		var v openApiClient.NetworksNetworkIdGroupPoliciesVlanTagging
		v.SetSettings(data.VlanTaggingSettings.ValueString())
		v.SetVlanId(data.VlanTaggingVlanId.ValueString())
		updateNetworkGroupPolicy.SetVlanTagging(v)
	} else {
		data.VlanTaggingSettings = jsontypes.StringNull()
		data.VlanTaggingVlanId = jsontypes.StringNull()
	}

	if !data.SchedulingEnabled.IsUnknown() {
		var s openApiClient.NetworksNetworkIdGroupPoliciesScheduling
		s.SetEnabled(data.SchedulingEnabled.ValueBool())
		if !data.SchedulingFridayActive.IsUnknown() {
			var f openApiClient.NetworksNetworkIdGroupPoliciesSchedulingFriday
			f.SetActive(data.SchedulingFridayActive.ValueBool())
			f.SetFrom(data.SchedulingFridayFrom.ValueString())
			f.SetTo(data.SchedulingFridayTo.ValueString())
			s.SetFriday(f)
			updateNetworkGroupPolicy.SetScheduling(s)
		} else {
			data.SchedulingFridayActive = jsontypes.BoolNull()
		}
		if !data.SchedulingMondayActive.IsUnknown() {
			var f openApiClient.NetworksNetworkIdGroupPoliciesSchedulingMonday
			f.SetActive(data.SchedulingMondayActive.ValueBool())
			f.SetFrom(data.SchedulingMondayFrom.ValueString())
			f.SetTo(data.SchedulingMondayTo.ValueString())
			s.SetMonday(f)
			updateNetworkGroupPolicy.SetScheduling(s)
		} else {
			data.SchedulingMondayActive = jsontypes.BoolNull()
		}
		if !data.SchedulingTuesdayActive.IsUnknown() {
			var f openApiClient.NetworksNetworkIdGroupPoliciesSchedulingTuesday
			f.SetActive(data.SchedulingTuesdayActive.ValueBool())
			f.SetFrom(data.SchedulingTuesdayFrom.ValueString())
			f.SetTo(data.SchedulingTuesdayTo.ValueString())
			s.SetTuesday(f)
			updateNetworkGroupPolicy.SetScheduling(s)
		} else {
			data.SchedulingTuesdayActive = jsontypes.BoolNull()
		}
		if !data.SchedulingWednesdayActive.IsUnknown() {
			var f openApiClient.NetworksNetworkIdGroupPoliciesSchedulingWednesday
			f.SetActive(data.SchedulingWednesdayActive.ValueBool())
			f.SetFrom(data.SchedulingWednesdayFrom.ValueString())
			f.SetTo(data.SchedulingWednesdayTo.ValueString())
			s.SetWednesday(f)
			updateNetworkGroupPolicy.SetScheduling(s)
		} else {
			data.SchedulingWednesdayActive = jsontypes.BoolNull()
		}
		if !data.SchedulingThursdayActive.IsUnknown() {
			var f openApiClient.NetworksNetworkIdGroupPoliciesSchedulingThursday
			f.SetActive(data.SchedulingThursdayActive.ValueBool())
			f.SetFrom(data.SchedulingThursdayFrom.ValueString())
			f.SetTo(data.SchedulingThursdayTo.ValueString())
			s.SetThursday(f)
			updateNetworkGroupPolicy.SetScheduling(s)
		} else {
			data.SchedulingThursdayActive = jsontypes.BoolNull()
		}
		if !data.SchedulingSaturdayActive.IsUnknown() {
			var f openApiClient.NetworksNetworkIdGroupPoliciesSchedulingSaturday
			f.SetActive(data.SchedulingSaturdayActive.ValueBool())
			f.SetFrom(data.SchedulingSaturdayFrom.ValueString())
			f.SetTo(data.SchedulingSaturdayTo.ValueString())
			s.SetSaturday(f)
			updateNetworkGroupPolicy.SetScheduling(s)
		} else {
			data.SchedulingSaturdayActive = jsontypes.BoolNull()
		}
		if !data.SchedulingSundayActive.IsUnknown() {
			var f openApiClient.NetworksNetworkIdGroupPoliciesSchedulingSunday
			f.SetActive(data.SchedulingSundayActive.ValueBool())
			f.SetFrom(data.SchedulingSundayFrom.ValueString())
			f.SetTo(data.SchedulingSundayTo.ValueString())
			s.SetSunday(f)
			updateNetworkGroupPolicy.SetScheduling(s)
		} else {
			data.SchedulingSundayActive = jsontypes.BoolNull()
		}

	} else {
		data.SchedulingEnabled = jsontypes.BoolNull()
	}

	var f openApiClient.NetworksNetworkIdGroupPoliciesFirewallAndTrafficShaping
	if !data.FirewallAndTrafficShapingSettings.IsUnknown() {

		f.SetSettings(data.FirewallAndTrafficShapingSettings.ValueString())

	} else {
		data.FirewallAndTrafficShapingSettings = jsontypes.StringNull()
	}

	if len(data.L3FirewallRules) > 0 {
		var l3s []openApiClient.NetworksNetworkIdGroupPoliciesFirewallAndTrafficShapingL3FirewallRules
		for _, attribute := range data.L3FirewallRules {
			var l3 openApiClient.NetworksNetworkIdGroupPoliciesFirewallAndTrafficShapingL3FirewallRules
			l3.SetComment(attribute.Comment.ValueString())
			l3.SetDestCidr(attribute.DestCidr.ValueString())
			l3.SetDestPort(attribute.DestPort.ValueString())
			l3.SetPolicy(attribute.Policy.ValueString())
			l3.SetProtocol(attribute.Protocol.ValueString())
			l3s = append(l3s, l3)
		}
		f.SetL3FirewallRules(l3s)
	} else {
		data.L3FirewallRules = nil
	}
	if len(data.L7FirewallRules) > 0 {
		var l7s []openApiClient.NetworksNetworkIdGroupPoliciesFirewallAndTrafficShapingL7FirewallRules
		for _, attribute := range data.L7FirewallRules {
			var l7 openApiClient.NetworksNetworkIdGroupPoliciesFirewallAndTrafficShapingL7FirewallRules

			l7.SetValue(attribute.Value.ValueString())
			l7.SetPolicy(attribute.Policy.ValueString())
			l7.SetType(attribute.Type.ValueString())
			l7s = append(l7s, l7)
		}
		f.SetL7FirewallRules(l7s)
	} else {
		data.L7FirewallRules = nil
	}

	if len(data.TrafficShapingRules) > 0 {
		var tfs []openApiClient.NetworksNetworkIdGroupPoliciesFirewallAndTrafficShapingTrafficShapingRules
		for _, attribute := range data.TrafficShapingRules {
			var tf openApiClient.NetworksNetworkIdGroupPoliciesFirewallAndTrafficShapingTrafficShapingRules
			tf.SetDscpTagValue(int32(attribute.DscpTagValue.ValueInt64()))
			tf.SetPcpTagValue(int32(attribute.PcpTagValue.ValueInt64()))
			tf.SetPriority(attribute.Priority.ValueString())
			var perclientBamdWidthLimits openApiClient.NetworksNetworkIdApplianceTrafficShapingRulesPerClientBandwidthLimits
			var bandwidthLimits openApiClient.NetworksNetworkIdApplianceTrafficShapingRulesPerClientBandwidthLimitsBandwidthLimits
			bandwidthLimits.SetLimitDown(int32(attribute.BandwidthLimitDown.ValueInt64()))
			bandwidthLimits.SetLimitUp(int32(attribute.BandwidthLimitUp.ValueInt64()))
			perclientBamdWidthLimits.SetBandwidthLimits(bandwidthLimits)
			tf.SetPerClientBandwidthLimits(perclientBamdWidthLimits)
			var defs []openApiClient.NetworksNetworkIdApplianceTrafficShapingRulesDefinitions
			for _, attribute := range attribute.Definitions {
				var def openApiClient.NetworksNetworkIdApplianceTrafficShapingRulesDefinitions
				def.SetType(attribute.Type.ValueString())
				def.SetValue(attribute.Value.ValueString())
				defs = append(defs, def)
			}
			tf.SetDefinitions(defs)
			tfs = append(tfs, tf)
		}
		f.SetTrafficShapingRules(tfs)
	} else {
		data.TrafficShapingRules = nil
	}

	updateNetworkGroupPolicy.SetFirewallAndTrafficShaping(f)

	var c openApiClient.NetworksNetworkIdGroupPoliciesContentFiltering
	if !data.ContentFilteringAllowedUrlPatternsSettings.IsUnknown() || len(data.UrlPatterns) > 0 {
		var aup openApiClient.NetworksNetworkIdGroupPoliciesContentFilteringAllowedUrlPatterns
		aup.SetSettings(data.ContentFilteringAllowedUrlPatternsSettings.ValueString())
		aup.SetPatterns(data.UrlPatterns)
		c.SetAllowedUrlPatterns(aup)

	} else {
		data.ContentFilteringAllowedUrlPatternsSettings = jsontypes.StringNull()
		data.UrlPatterns = nil
	}
	if !data.ContentFilteringBlockedUrlCategoriesSettings.IsUnknown() || len(data.Categories) > 0 {
		var buc openApiClient.NetworksNetworkIdGroupPoliciesContentFilteringBlockedUrlCategories
		buc.SetSettings(data.ContentFilteringBlockedUrlCategoriesSettings.ValueString())
		buc.SetCategories(data.Categories)
		c.SetBlockedUrlCategories(buc)

	} else {
		data.ContentFilteringBlockedUrlCategoriesSettings = jsontypes.StringNull()
		data.Categories = nil
	}
	if !data.ContentFilteringBlockedUrlPatternsSettings.IsUnknown() || len(data.BlockedUrlPatterns) > 0 {
		var bup openApiClient.NetworksNetworkIdGroupPoliciesContentFilteringBlockedUrlPatterns
		bup.SetSettings(data.ContentFilteringBlockedUrlPatternsSettings.ValueString())
		bup.SetPatterns(data.BlockedUrlPatterns)
		c.SetBlockedUrlPatterns(bup)

	} else {
		data.ContentFilteringBlockedUrlPatternsSettings = jsontypes.StringNull()
		data.BlockedUrlPatterns = nil
	}
	updateNetworkGroupPolicy.SetContentFiltering(c)

	_, httpResp, err := r.client.NetworksApi.UpdateNetworkGroupPolicy(ctx, data.NetworkId.ValueString(), data.GroupPolicyId.ValueString()).UpdateNetworkGroupPolicy(updateNetworkGroupPolicy).Execute()
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

	// Save data into Terraform state
	data.Id = jsontypes.StringValue("example-id")
	data.Name = structmap.Name
	data.GroupPolicyId = structmap.GroupPolicyId
	data.BandWidthLimitDown = structmap.BandWidth.BandwidthLimits.LimitDown
	data.BandWidthLimitUp = structmap.BandWidth.BandwidthLimits.LimitUp
	data.BandWidthSettings = structmap.BandWidth.Settings
	data.SplashAuthSettings = structmap.SplashAuthSettings
	data.BandWidthLimitDown = structmap.BandWidth.BandwidthLimits.LimitDown
	data.BandWidthLimitUp = structmap.BandWidth.BandwidthLimits.LimitUp
	data.BonjourForwardingRules = structmap.BonjourForwarding.Rules
	data.BonjourForwardingSettings = structmap.BonjourForwarding.Settings
	data.VlanTaggingSettings = structmap.VlanTagging.Settings
	data.VlanTaggingVlanId = structmap.VlanTagging.VlanId
	data.SchedulingEnabled = structmap.Scheduling.Enabled
	data.SchedulingFridayActive = structmap.Scheduling.Friday.Active
	data.SchedulingFridayFrom = structmap.Scheduling.Friday.From
	data.SchedulingFridayTo = structmap.Scheduling.Friday.To
	data.SchedulingMondayActive = structmap.Scheduling.Monday.Active
	data.SchedulingMondayFrom = structmap.Scheduling.Monday.From
	data.SchedulingMondayTo = structmap.Scheduling.Monday.To
	data.SchedulingTuesdayActive = structmap.Scheduling.Tuesday.Active
	data.SchedulingTuesdayFrom = structmap.Scheduling.Tuesday.From
	data.SchedulingTuesdayTo = structmap.Scheduling.Tuesday.To
	data.SchedulingWednesdayActive = structmap.Scheduling.Wednesday.Active
	data.SchedulingWednesdayFrom = structmap.Scheduling.Wednesday.From
	data.SchedulingWednesdayTo = structmap.Scheduling.Wednesday.To
	data.SchedulingThursdayActive = structmap.Scheduling.Thursday.Active
	data.SchedulingThursdayFrom = structmap.Scheduling.Thursday.From
	data.SchedulingThursdayTo = structmap.Scheduling.Thursday.To
	data.SchedulingSaturdayActive = structmap.Scheduling.Saturday.Active
	data.SchedulingSaturdayFrom = structmap.Scheduling.Saturday.From
	data.SchedulingSaturdayTo = structmap.Scheduling.Saturday.To
	data.SchedulingSundayActive = structmap.Scheduling.Sunday.Active
	data.SchedulingSundayFrom = structmap.Scheduling.Sunday.From
	data.SchedulingSundayTo = structmap.Scheduling.Sunday.To
	if len(structmap.FirewallAndTrafficShaping.L3FirewallRules) > 0 {

		data.L3FirewallRules = structmap.FirewallAndTrafficShaping.L3FirewallRules
	}
	if len(structmap.FirewallAndTrafficShaping.L7FirewallRules) > 0 {
		data.L7FirewallRules = structmap.FirewallAndTrafficShaping.L7FirewallRules
	}
	data.FirewallAndTrafficShapingSettings = structmap.FirewallAndTrafficShaping.Settings
	for _, attribute := range structmap.FirewallAndTrafficShaping.TrafficShapingRules {
		var a TrafficShapingRule
		a.Definitions = attribute.Definitions
		a.BandwidthLimitDown = attribute.PerClientBandwidthLimits.BandwidthLimits.LimitDown
		a.BandwidthLimitUp = attribute.PerClientBandwidthLimits.BandwidthLimits.LimitUp
		a.Priority = attribute.Priority
		a.DscpTagValue = attribute.DscpTagValue
		a.PcpTagValue = attribute.PcpTagValue
		data.TrafficShapingRules = append(data.TrafficShapingRules, a)
	}
	data.ContentFilteringAllowedUrlPatternsSettings = structmap.ContentFiltering.AllowedUrlPatterns.Settings
	if len(structmap.ContentFiltering.AllowedUrlPatterns.Patterns) > 0 {
		data.UrlPatterns = structmap.ContentFiltering.AllowedUrlPatterns.Patterns
	} else {
		data.UrlPatterns = nil
	}
	data.ContentFilteringBlockedUrlCategoriesSettings = structmap.ContentFiltering.BlockedUrlCategories.Settings
	if len(structmap.ContentFiltering.BlockedUrlCategories.Categories) > 0 {
		data.Categories = structmap.ContentFiltering.BlockedUrlCategories.Categories
	} else {
		data.Categories = nil
	}
	data.ContentFilteringBlockedUrlPatternsSettings = structmap.ContentFiltering.BlockedUrlPatterns.Settings
	if len(structmap.ContentFiltering.BlockedUrlPatterns.Patterns) > 0 {

		data.BlockedUrlPatterns = structmap.ContentFiltering.BlockedUrlPatterns.Patterns
	} else {
		data.BlockedUrlPatterns = nil
	}

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
