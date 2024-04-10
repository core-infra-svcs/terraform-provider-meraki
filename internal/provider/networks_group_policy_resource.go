package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
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

// NetworksGroupPolicyResourceModel describes the resource data model.
type NetworksGroupPolicyResourceModel struct {
	Id                        jsontypes.String                                          `tfsdk:"id"`
	NetworkId                 jsontypes.String                                          `tfsdk:"network_id"`
	GroupPolicyId             jsontypes.String                                          `tfsdk:"group_policy_id" json:"groupPolicyId"`
	Name                      jsontypes.String                                          `tfsdk:"name" json:"name"`
	SplashAuthSettings        jsontypes.String                                          `tfsdk:"splash_auth_settings" json:"splashAuthSettings"`
	Bandwidth                 NetworksGroupPolicyResourceModelBandwidth                 `tfsdk:"bandwidth" json:"bandwidth"`
	BonjourForwarding         NetworksGroupPolicyResourceModelBonjourForwarding         `tfsdk:"bonjour_forwarding" json:"bonjourForwarding"`
	FirewallAndTrafficShaping NetworksGroupPolicyResourceModelFirewallAndTrafficShaping `tfsdk:"firewall_and_traffic_shaping" json:"firewallAndTrafficShaping"`
	Scheduling                NetworksGroupPolicyResourceModelScheduling                `tfsdk:"scheduling" json:"scheduling"`
	VlanTagging               NetworksGroupPolicyResourceModelVlanTagging               `tfsdk:"vlan_tagging" json:"vlanTagging"`
	ContentFiltering          NetworksGroupPolicyResourceModelContentFiltering          `tfsdk:"content_filtering" json:"contentFiltering"`
}

type NetworksGroupPolicyResourceModelFirewallAndTrafficShaping struct {
	Settings            jsontypes.String                                     `tfsdk:"settings" json:"settings"`
	L3FirewallRules     []NetworksGroupPolicyResourceModelL3FirewallRule     `tfsdk:"l3_firewall_rules" json:"l3FirewallRules"`
	L7FirewallRules     []NetworksGroupPolicyResourceModelL7FirewallRule     `tfsdk:"l7_firewall_rules" json:"l7FirewallRules"`
	TrafficShapingRules []NetworksGroupPolicyResourceModelTrafficShapingRule `tfsdk:"traffic_shaping_rules" json:"trafficShapingRules"`
}

type NetworksGroupPolicyResourceModelL3FirewallRule struct {
	Comment  jsontypes.String `tfsdk:"comment" json:"comment"`
	DestCidr jsontypes.String `tfsdk:"dest_cidr" json:"destCidr"`
	DestPort jsontypes.String `tfsdk:"dest_port" json:"destPort"`
	Policy   jsontypes.String `tfsdk:"policy" json:"policy"`
	Protocol jsontypes.String `tfsdk:"protocol" json:"protocol"`
}

type NetworksGroupPolicyResourceModelL7FirewallRule struct {
	Value  jsontypes.String `tfsdk:"value" json:"value"`
	Type   jsontypes.String `tfsdk:"type" json:"type"`
	Policy jsontypes.String `tfsdk:"policy" json:"policy"`
}

type NetworksGroupPolicyResourceModelTrafficShapingRule struct {
	DscpTagValue             jsontypes.Int64                                          `tfsdk:"dscp_tag_value" json:"dscpTagValue"`
	PcpTagValue              jsontypes.Int64                                          `tfsdk:"pcp_tag_value" json:"pcpTagValue"`
	PerClientBandwidthLimits NetworksGroupPolicyResourceModelPerClientBandwidthLimits `tfsdk:"per_client_bandwidth_limits" json:"perClientBandwidthLimits,,omitempty"`
	Definitions              []NetworksGroupPolicyResourceModelDefinition             `tfsdk:"definitions" json:"definitions"`
}

type NetworksGroupPolicyResourceModelPerClientBandwidthLimits struct {
	BandwidthLimits NetworksGroupPolicyResourceModelBandwidthLimits `tfsdk:"bandwidth_limits" json:"bandwidthLimits,,omitempty"`
	Settings        jsontypes.String                                `tfsdk:"settings" json:"settings,,omitempty"`
}

type NetworksGroupPolicyResourceModelDefinition struct {
	Value jsontypes.String `tfsdk:"value" json:"value"`
	Type  jsontypes.String `tfsdk:"type" json:"type"`
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
	BonjourForwardingSettings jsontypes.String                       `tfsdk:"settings" json:"settings"`
	BonjourForwardingRules    []NetworksGroupPolicyResourceModelRule `tfsdk:"rules" json:"rules"`
}

type NetworksGroupPolicyResourceModelRule struct {
	Description jsontypes.String `tfsdk:"description" json:"description"`
	VlanId      jsontypes.String `tfsdk:"vlan_id" json:"vlanId"`
	Services    []string         `tfsdk:"services" json:"services"`
}

type NetworksGroupPolicyResourceModelScheduling struct {
	Enabled   jsontypes.Bool                           `tfsdk:"enabled" json:"enabled"`
	Friday    NetworksGroupPolicyResourceModelSchedule `tfsdk:"friday" json:"friday"`
	Monday    NetworksGroupPolicyResourceModelSchedule `tfsdk:"monday" json:"monday"`
	Saturday  NetworksGroupPolicyResourceModelSchedule `tfsdk:"saturday" json:"saturday"`
	Sunday    NetworksGroupPolicyResourceModelSchedule `tfsdk:"sunday" json:"sunday"`
	Thursday  NetworksGroupPolicyResourceModelSchedule `tfsdk:"thursday" json:"thursday"`
	Tuesday   NetworksGroupPolicyResourceModelSchedule `tfsdk:"tuesday" json:"tuesday"`
	Wednesday NetworksGroupPolicyResourceModelSchedule `tfsdk:"wednesday" json:"wednesday"`
}

type NetworksGroupPolicyResourceModelSchedule struct {
	From   jsontypes.String `tfsdk:"from" json:"from"`
	To     jsontypes.String `tfsdk:"to" json:"to"`
	Active jsontypes.Bool   `tfsdk:"active" json:"active"`
}

type NetworksGroupPolicyResourceModelVlanTagging struct {
	Settings jsontypes.String `tfsdk:"settings" json:"settings"`
	VlanId   jsontypes.String `tfsdk:"vlan_id" json:"vlanId"`
}

type NetworksGroupPolicyResourceModelContentFiltering struct {
	AllowedUrlPatterns   NetworksGroupPolicyResourceModelAllowedUrlPatterns   `tfsdk:"allowed_url_patterns" json:"allowedUrlPatterns"`
	BlockedUrlCategories NetworksGroupPolicyResourceModelBlockedUrlCategories `tfsdk:"blocked_url_categories" json:"blockedUrlCategories"`
	BlockedUrlPatterns   NetworksGroupPolicyResourceModelBlockedUrlPatterns   `tfsdk:"blocked_url_patterns" json:"blockedUrlPatterns"`
}

type NetworksGroupPolicyResourceModelAllowedUrlPatterns struct {
	Settings jsontypes.String `tfsdk:"settings" json:"settings"`
	Patterns []string         `tfsdk:"patterns" json:"patterns"`
}
type NetworksGroupPolicyResourceModelBlockedUrlCategories struct {
	Settings   jsontypes.String `tfsdk:"settings" json:"settings"`
	Categories []string         `tfsdk:"categories" json:"categories"`
}
type NetworksGroupPolicyResourceModelBlockedUrlPatterns struct {
	Settings jsontypes.String `tfsdk:"settings" json:"settings"`
	Patterns []string         `tfsdk:"patterns" json:"patterns"`
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

	if !data.Bandwidth.Settings.IsUnknown() {
		var bandwidth openApiClient.CreateNetworkGroupPolicyRequestBandwidth
		bandwidth.SetSettings(data.Bandwidth.Settings.ValueString())
		var bandwidthLimits openApiClient.CreateNetworkGroupPolicyRequestBandwidthBandwidthLimits
		if !data.Bandwidth.BandwidthLimits.LimitUp.IsUnknown() {
			bandwidthLimits.SetLimitUp(int32(data.Bandwidth.BandwidthLimits.LimitUp.ValueInt64()))
		}
		if !data.Bandwidth.BandwidthLimits.LimitDown.IsUnknown() {
			bandwidthLimits.SetLimitDown(int32(data.Bandwidth.BandwidthLimits.LimitDown.ValueInt64()))
		}
		bandwidth.SetBandwidthLimits(bandwidthLimits)
		createNetworkGroupPolicy.SetBandwidth(bandwidth)
	}

	if len(data.BonjourForwarding.BonjourForwardingRules) > 0 {
		var bonjourForwarding openApiClient.CreateNetworkGroupPolicyRequestBonjourForwarding
		var bonjourForwardingRules []openApiClient.CreateNetworkGroupPolicyRequestBonjourForwardingRulesInner
		for _, attribute := range data.BonjourForwarding.BonjourForwardingRules {
			var bonjourForwardingRule openApiClient.CreateNetworkGroupPolicyRequestBonjourForwardingRulesInner
			if !attribute.Description.IsUnknown() {
				bonjourForwardingRule.SetDescription(attribute.Description.ValueString())
			}
			bonjourForwardingRule.SetVlanId(attribute.VlanId.ValueString())
			bonjourForwardingRule.SetServices(attribute.Services)
			bonjourForwardingRules = append(bonjourForwardingRules, bonjourForwardingRule)
		}
		bonjourForwarding.SetRules(bonjourForwardingRules)
		if !data.BonjourForwarding.BonjourForwardingSettings.IsUnknown() {
			bonjourForwarding.SetSettings(data.BonjourForwarding.BonjourForwardingSettings.ValueString())
		}
		createNetworkGroupPolicy.SetBonjourForwarding(bonjourForwarding)
	}

	var firewallAndTrafficShaping openApiClient.CreateNetworkGroupPolicyRequestFirewallAndTrafficShaping

	if !data.FirewallAndTrafficShaping.Settings.IsUnknown() {
		firewallAndTrafficShaping.SetSettings(data.FirewallAndTrafficShaping.Settings.ValueString())
	}

	if len(data.FirewallAndTrafficShaping.L3FirewallRules) > 0 {
		var l3s []openApiClient.CreateNetworkGroupPolicyRequestFirewallAndTrafficShapingL3FirewallRulesInner
		for _, attribute := range data.FirewallAndTrafficShaping.L3FirewallRules {
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

	if len(data.FirewallAndTrafficShaping.L7FirewallRules) > 0 {
		var l7s []openApiClient.CreateNetworkGroupPolicyRequestFirewallAndTrafficShapingL7FirewallRulesInner
		for _, attribute := range data.FirewallAndTrafficShaping.L7FirewallRules {
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

	if len(data.FirewallAndTrafficShaping.TrafficShapingRules) > 0 {
		var tfs []openApiClient.CreateNetworkGroupPolicyRequestFirewallAndTrafficShapingTrafficShapingRulesInner
		for _, attribute := range data.FirewallAndTrafficShaping.TrafficShapingRules {
			var tf openApiClient.CreateNetworkGroupPolicyRequestFirewallAndTrafficShapingTrafficShapingRulesInner
			if !attribute.DscpTagValue.IsUnknown() {
				tf.SetDscpTagValue(int32(attribute.DscpTagValue.ValueInt64()))
			}
			if !attribute.PcpTagValue.IsUnknown() {
				tf.SetPcpTagValue(int32(attribute.PcpTagValue.ValueInt64()))
			}
			var perClientBandWidthLimits openApiClient.UpdateNetworkApplianceTrafficShapingRulesRequestRulesInnerPerClientBandwidthLimits

			if !attribute.PerClientBandwidthLimits.Settings.IsUnknown() {
				var bandwidthLimits openApiClient.UpdateNetworkApplianceTrafficShapingRulesRequestRulesInnerPerClientBandwidthLimitsBandwidthLimits

				if attribute.PerClientBandwidthLimits.Settings.ValueString() != "network default" {

					if !attribute.PerClientBandwidthLimits.BandwidthLimits.LimitDown.Int64Value.IsUnknown() {
						bandwidthLimits.SetLimitDown(int32(attribute.PerClientBandwidthLimits.BandwidthLimits.LimitDown.ValueInt64()))
					}

					if !attribute.PerClientBandwidthLimits.BandwidthLimits.LimitUp.Int64Value.IsUnknown() {
						bandwidthLimits.SetLimitUp(int32(attribute.PerClientBandwidthLimits.BandwidthLimits.LimitUp.ValueInt64()))
					}
					perClientBandWidthLimits.SetBandwidthLimits(bandwidthLimits)
				}
				perClientBandWidthLimits.SetSettings(attribute.PerClientBandwidthLimits.Settings.ValueString())
				tf.SetPerClientBandwidthLimits(perClientBandWidthLimits)
			}
			var defs []openApiClient.UpdateNetworkApplianceTrafficShapingRulesRequestRulesInnerDefinitionsInner
			if len(attribute.Definitions) > 0 {
				for _, attribute := range attribute.Definitions {
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

	if !data.Scheduling.Enabled.IsUnknown() {
		var schedule openApiClient.CreateNetworkGroupPolicyRequestScheduling
		schedule.SetEnabled(data.Scheduling.Enabled.ValueBool())
		if !data.Scheduling.Friday.Active.IsUnknown() {
			var friday openApiClient.CreateNetworkGroupPolicyRequestSchedulingFriday
			friday.SetActive(data.Scheduling.Friday.Active.ValueBool())
			friday.SetFrom(data.Scheduling.Friday.From.ValueString())
			friday.SetTo(data.Scheduling.Friday.To.ValueString())
			schedule.SetFriday(friday)
		}
		if !data.Scheduling.Monday.Active.IsUnknown() {
			var monday openApiClient.CreateNetworkGroupPolicyRequestSchedulingMonday
			monday.SetActive(data.Scheduling.Monday.Active.ValueBool())
			monday.SetFrom(data.Scheduling.Monday.From.ValueString())
			monday.SetTo(data.Scheduling.Monday.To.ValueString())
			schedule.SetMonday(monday)
		}
		if !data.Scheduling.Tuesday.Active.IsUnknown() {
			var tuesday openApiClient.CreateNetworkGroupPolicyRequestSchedulingTuesday
			tuesday.SetActive(data.Scheduling.Tuesday.Active.ValueBool())
			tuesday.SetFrom(data.Scheduling.Tuesday.From.ValueString())
			tuesday.SetTo(data.Scheduling.Tuesday.To.ValueString())
			schedule.SetTuesday(tuesday)
		}
		if !data.Scheduling.Wednesday.Active.IsUnknown() {
			var wednesday openApiClient.CreateNetworkGroupPolicyRequestSchedulingWednesday
			wednesday.SetActive(data.Scheduling.Wednesday.Active.ValueBool())
			wednesday.SetFrom(data.Scheduling.Wednesday.From.ValueString())
			wednesday.SetTo(data.Scheduling.Wednesday.To.ValueString())
			schedule.SetWednesday(wednesday)
		}
		if !data.Scheduling.Thursday.Active.IsUnknown() {
			var thursday openApiClient.CreateNetworkGroupPolicyRequestSchedulingThursday
			thursday.SetActive(data.Scheduling.Thursday.Active.ValueBool())
			thursday.SetFrom(data.Scheduling.Thursday.From.ValueString())
			thursday.SetTo(data.Scheduling.Thursday.To.ValueString())
			schedule.SetThursday(thursday)
		}
		if !data.Scheduling.Saturday.Active.IsUnknown() {
			var saturday openApiClient.CreateNetworkGroupPolicyRequestSchedulingSaturday
			saturday.SetActive(data.Scheduling.Saturday.Active.ValueBool())
			saturday.SetFrom(data.Scheduling.Saturday.From.ValueString())
			saturday.SetTo(data.Scheduling.Saturday.To.ValueString())
			schedule.SetSaturday(saturday)
		}
		if !data.Scheduling.Sunday.Active.IsUnknown() {
			var sunday openApiClient.CreateNetworkGroupPolicyRequestSchedulingSunday
			sunday.SetActive(data.Scheduling.Sunday.Active.ValueBool())
			sunday.SetFrom(data.Scheduling.Sunday.From.ValueString())
			sunday.SetTo(data.Scheduling.Sunday.To.ValueString())
			schedule.SetSunday(sunday)
		}
		createNetworkGroupPolicy.SetScheduling(schedule)
	}

	if !data.VlanTagging.Settings.IsUnknown() {
		if !data.VlanTagging.VlanId.IsUnknown() {
			var v openApiClient.CreateNetworkGroupPolicyRequestVlanTagging
			v.SetSettings(data.VlanTagging.Settings.ValueString())
			v.SetVlanId(data.VlanTagging.VlanId.ValueString())
			createNetworkGroupPolicy.SetVlanTagging(v)
		}
	}
	var contentFiltering openApiClient.CreateNetworkGroupPolicyRequestContentFiltering
	contentFilteringStatus := false

	if !data.ContentFiltering.AllowedUrlPatterns.Settings.IsUnknown() {
		var allowedUrlPatternData openApiClient.CreateNetworkGroupPolicyRequestContentFilteringAllowedUrlPatterns
		allowedUrlPatternData.SetSettings(data.ContentFiltering.AllowedUrlPatterns.Settings.ValueString())
		allowedUrlPatternData.SetPatterns(data.ContentFiltering.AllowedUrlPatterns.Patterns)
		contentFiltering.SetAllowedUrlPatterns(allowedUrlPatternData)
		contentFilteringStatus = true
	}

	if !data.ContentFiltering.BlockedUrlCategories.Settings.IsUnknown() {
		var blockedUrlCategoriesData openApiClient.CreateNetworkGroupPolicyRequestContentFilteringBlockedUrlCategories
		blockedUrlCategoriesData.SetSettings(data.ContentFiltering.BlockedUrlCategories.Settings.ValueString())
		blockedUrlCategoriesData.SetCategories(data.ContentFiltering.BlockedUrlCategories.Categories)
		contentFiltering.SetBlockedUrlCategories(blockedUrlCategoriesData)
		contentFilteringStatus = true
	}

	if !data.ContentFiltering.BlockedUrlPatterns.Settings.IsUnknown() {
		var blockedUrlPatternData openApiClient.CreateNetworkGroupPolicyRequestContentFilteringBlockedUrlPatterns
		blockedUrlPatternData.SetSettings(data.ContentFiltering.BlockedUrlPatterns.Settings.ValueString())
		blockedUrlPatternData.SetPatterns(data.ContentFiltering.BlockedUrlPatterns.Patterns)
		contentFiltering.SetBlockedUrlPatterns(blockedUrlPatternData)
		contentFilteringStatus = true
	}

	if contentFilteringStatus {
		createNetworkGroupPolicy.SetContentFiltering(contentFiltering)
	}

	inlineResp, httpResp, err := r.client.NetworksApi.CreateNetworkGroupPolicy(ctx, data.NetworkId.ValueString()).CreateNetworkGroupPolicyRequest(createNetworkGroupPolicy).Execute()
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

	if !data.Bandwidth.Settings.IsUnknown() {
		var bandwidth openApiClient.CreateNetworkGroupPolicyRequestBandwidth
		bandwidth.SetSettings(data.Bandwidth.Settings.ValueString())
		var bandwidthLimits openApiClient.CreateNetworkGroupPolicyRequestBandwidthBandwidthLimits
		if !data.Bandwidth.BandwidthLimits.LimitUp.IsUnknown() {
			bandwidthLimits.SetLimitUp(int32(data.Bandwidth.BandwidthLimits.LimitUp.ValueInt64()))
		}
		if !data.Bandwidth.BandwidthLimits.LimitDown.IsUnknown() {
			bandwidthLimits.SetLimitDown(int32(data.Bandwidth.BandwidthLimits.LimitDown.ValueInt64()))
		}
		bandwidth.SetBandwidthLimits(bandwidthLimits)
		updateNetworkGroupPolicy.SetBandwidth(bandwidth)
	}

	if len(data.BonjourForwarding.BonjourForwardingRules) > 0 {
		var bonjourForwarding openApiClient.CreateNetworkGroupPolicyRequestBonjourForwarding
		var bonjourForwardingRules []openApiClient.CreateNetworkGroupPolicyRequestBonjourForwardingRulesInner
		for _, attribute := range data.BonjourForwarding.BonjourForwardingRules {
			var bonjourForwardingRule openApiClient.CreateNetworkGroupPolicyRequestBonjourForwardingRulesInner
			if !attribute.Description.IsUnknown() {
				bonjourForwardingRule.SetDescription(attribute.Description.ValueString())
			}
			bonjourForwardingRule.SetVlanId(attribute.VlanId.ValueString())
			bonjourForwardingRule.SetServices(attribute.Services)
			bonjourForwardingRules = append(bonjourForwardingRules, bonjourForwardingRule)
		}
		bonjourForwarding.SetRules(bonjourForwardingRules)
		if !data.BonjourForwarding.BonjourForwardingSettings.IsUnknown() {
			bonjourForwarding.SetSettings(data.BonjourForwarding.BonjourForwardingSettings.ValueString())
		}
		updateNetworkGroupPolicy.SetBonjourForwarding(bonjourForwarding)
	}

	var firewallAndTrafficShaping openApiClient.CreateNetworkGroupPolicyRequestFirewallAndTrafficShaping

	if !data.FirewallAndTrafficShaping.Settings.IsUnknown() {
		firewallAndTrafficShaping.SetSettings(data.FirewallAndTrafficShaping.Settings.ValueString())
	}

	if len(data.FirewallAndTrafficShaping.L3FirewallRules) > 0 {
		var l3s []openApiClient.CreateNetworkGroupPolicyRequestFirewallAndTrafficShapingL3FirewallRulesInner
		for _, attribute := range data.FirewallAndTrafficShaping.L3FirewallRules {
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

	if len(data.FirewallAndTrafficShaping.L7FirewallRules) > 0 {
		var l7s []openApiClient.CreateNetworkGroupPolicyRequestFirewallAndTrafficShapingL7FirewallRulesInner
		for _, attribute := range data.FirewallAndTrafficShaping.L7FirewallRules {
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

	if len(data.FirewallAndTrafficShaping.TrafficShapingRules) > 0 {
		var tfs []openApiClient.CreateNetworkGroupPolicyRequestFirewallAndTrafficShapingTrafficShapingRulesInner
		for _, attribute := range data.FirewallAndTrafficShaping.TrafficShapingRules {
			var tf openApiClient.CreateNetworkGroupPolicyRequestFirewallAndTrafficShapingTrafficShapingRulesInner
			if !attribute.DscpTagValue.IsUnknown() {
				tf.SetDscpTagValue(int32(attribute.DscpTagValue.ValueInt64()))
			}
			if !attribute.PcpTagValue.IsUnknown() {
				tf.SetPcpTagValue(int32(attribute.PcpTagValue.ValueInt64()))
			}
			var perClientBandWidthLimits openApiClient.UpdateNetworkApplianceTrafficShapingRulesRequestRulesInnerPerClientBandwidthLimits

			if !attribute.PerClientBandwidthLimits.Settings.IsUnknown() {
				var bandwidthLimits openApiClient.UpdateNetworkApplianceTrafficShapingRulesRequestRulesInnerPerClientBandwidthLimitsBandwidthLimits

				if attribute.PerClientBandwidthLimits.Settings.ValueString() != "network default" {

					if !attribute.PerClientBandwidthLimits.BandwidthLimits.LimitDown.Int64Value.IsUnknown() {
						bandwidthLimits.SetLimitDown(int32(attribute.PerClientBandwidthLimits.BandwidthLimits.LimitDown.ValueInt64()))
					}

					if !attribute.PerClientBandwidthLimits.BandwidthLimits.LimitUp.Int64Value.IsUnknown() {
						bandwidthLimits.SetLimitUp(int32(attribute.PerClientBandwidthLimits.BandwidthLimits.LimitUp.ValueInt64()))
					}
					perClientBandWidthLimits.SetBandwidthLimits(bandwidthLimits)
				}
				perClientBandWidthLimits.SetSettings(attribute.PerClientBandwidthLimits.Settings.ValueString())
				tf.SetPerClientBandwidthLimits(perClientBandWidthLimits)
			}
			var defs []openApiClient.UpdateNetworkApplianceTrafficShapingRulesRequestRulesInnerDefinitionsInner
			if len(attribute.Definitions) > 0 {
				for _, attribute := range attribute.Definitions {
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

	if !data.Scheduling.Enabled.IsUnknown() {
		var schedule openApiClient.CreateNetworkGroupPolicyRequestScheduling
		schedule.SetEnabled(data.Scheduling.Enabled.ValueBool())
		if !data.Scheduling.Friday.Active.IsUnknown() {
			var friday openApiClient.CreateNetworkGroupPolicyRequestSchedulingFriday
			friday.SetActive(data.Scheduling.Friday.Active.ValueBool())
			friday.SetFrom(data.Scheduling.Friday.From.ValueString())
			friday.SetTo(data.Scheduling.Friday.To.ValueString())
			schedule.SetFriday(friday)
		}
		if !data.Scheduling.Monday.Active.IsUnknown() {
			var monday openApiClient.CreateNetworkGroupPolicyRequestSchedulingMonday
			monday.SetActive(data.Scheduling.Monday.Active.ValueBool())
			monday.SetFrom(data.Scheduling.Monday.From.ValueString())
			monday.SetTo(data.Scheduling.Monday.To.ValueString())
			schedule.SetMonday(monday)
		}
		if !data.Scheduling.Tuesday.Active.IsUnknown() {
			var tuesday openApiClient.CreateNetworkGroupPolicyRequestSchedulingTuesday
			tuesday.SetActive(data.Scheduling.Tuesday.Active.ValueBool())
			tuesday.SetFrom(data.Scheduling.Tuesday.From.ValueString())
			tuesday.SetTo(data.Scheduling.Tuesday.To.ValueString())
			schedule.SetTuesday(tuesday)
		}
		if !data.Scheduling.Wednesday.Active.IsUnknown() {
			var wednesday openApiClient.CreateNetworkGroupPolicyRequestSchedulingWednesday
			wednesday.SetActive(data.Scheduling.Wednesday.Active.ValueBool())
			wednesday.SetFrom(data.Scheduling.Wednesday.From.ValueString())
			wednesday.SetTo(data.Scheduling.Wednesday.To.ValueString())
			schedule.SetWednesday(wednesday)
		}
		if !data.Scheduling.Thursday.Active.IsUnknown() {
			var thursday openApiClient.CreateNetworkGroupPolicyRequestSchedulingThursday
			thursday.SetActive(data.Scheduling.Thursday.Active.ValueBool())
			thursday.SetFrom(data.Scheduling.Thursday.From.ValueString())
			thursday.SetTo(data.Scheduling.Thursday.To.ValueString())
			schedule.SetThursday(thursday)
		}
		if !data.Scheduling.Saturday.Active.IsUnknown() {
			var saturday openApiClient.CreateNetworkGroupPolicyRequestSchedulingSaturday
			saturday.SetActive(data.Scheduling.Saturday.Active.ValueBool())
			saturday.SetFrom(data.Scheduling.Saturday.From.ValueString())
			saturday.SetTo(data.Scheduling.Saturday.To.ValueString())
			schedule.SetSaturday(saturday)
		}
		if !data.Scheduling.Sunday.Active.IsUnknown() {
			var sunday openApiClient.CreateNetworkGroupPolicyRequestSchedulingSunday
			sunday.SetActive(data.Scheduling.Sunday.Active.ValueBool())
			sunday.SetFrom(data.Scheduling.Sunday.From.ValueString())
			sunday.SetTo(data.Scheduling.Sunday.To.ValueString())
			schedule.SetSunday(sunday)
		}
		updateNetworkGroupPolicy.SetScheduling(schedule)
	}

	if !data.VlanTagging.Settings.IsUnknown() {
		if !data.VlanTagging.VlanId.IsUnknown() {
			var v openApiClient.CreateNetworkGroupPolicyRequestVlanTagging
			v.SetSettings(data.VlanTagging.Settings.ValueString())
			v.SetVlanId(data.VlanTagging.VlanId.ValueString())
			updateNetworkGroupPolicy.SetVlanTagging(v)
		}
	}
	var contentFiltering openApiClient.CreateNetworkGroupPolicyRequestContentFiltering
	contentFilteringStatus := false

	if !data.ContentFiltering.AllowedUrlPatterns.Settings.IsUnknown() {
		var allowedUrlPatternData openApiClient.CreateNetworkGroupPolicyRequestContentFilteringAllowedUrlPatterns
		allowedUrlPatternData.SetSettings(data.ContentFiltering.AllowedUrlPatterns.Settings.ValueString())
		allowedUrlPatternData.SetPatterns(data.ContentFiltering.AllowedUrlPatterns.Patterns)
		contentFiltering.SetAllowedUrlPatterns(allowedUrlPatternData)
		contentFilteringStatus = true
	}

	if !data.ContentFiltering.BlockedUrlCategories.Settings.IsUnknown() {
		var blockedUrlCategoriesData openApiClient.CreateNetworkGroupPolicyRequestContentFilteringBlockedUrlCategories
		blockedUrlCategoriesData.SetSettings(data.ContentFiltering.BlockedUrlCategories.Settings.ValueString())
		blockedUrlCategoriesData.SetCategories(data.ContentFiltering.BlockedUrlCategories.Categories)
		contentFiltering.SetBlockedUrlCategories(blockedUrlCategoriesData)
		contentFilteringStatus = true
	}

	if !data.ContentFiltering.BlockedUrlPatterns.Settings.IsUnknown() {
		var blockedUrlPatternData openApiClient.CreateNetworkGroupPolicyRequestContentFilteringBlockedUrlPatterns
		blockedUrlPatternData.SetSettings(data.ContentFiltering.BlockedUrlPatterns.Settings.ValueString())
		blockedUrlPatternData.SetPatterns(data.ContentFiltering.BlockedUrlPatterns.Patterns)
		contentFiltering.SetBlockedUrlPatterns(blockedUrlPatternData)
		contentFilteringStatus = true
	}

	if contentFilteringStatus {
		updateNetworkGroupPolicy.SetContentFiltering(contentFiltering)
	}

	inlineResp, httpResp, err := r.client.NetworksApi.UpdateNetworkGroupPolicy(ctx, data.NetworkId.ValueString(), data.GroupPolicyId.ValueString()).UpdateNetworkGroupPolicyRequest(updateNetworkGroupPolicy).Execute()
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

	if err := json.NewDecoder(httpRespBody).Decode(data); err != nil {
		return data, err
	}
	var trafficShapingRule []NetworksGroupPolicyResourceModelTrafficShapingRule
	jsonData, _ := json.Marshal(inlineResp["firewallAndTrafficShaping"].(map[string]interface{})["trafficShapingRules"])
	json.Unmarshal(jsonData, &trafficShapingRule)
	if len(trafficShapingRule) > 0 {
		for _, attribute := range trafficShapingRule {
			var trafficShapingRule NetworksGroupPolicyResourceModelTrafficShapingRule
			trafficShapingRule.DscpTagValue = attribute.DscpTagValue
			trafficShapingRule.PcpTagValue = attribute.PcpTagValue
			trafficShapingRule.PerClientBandwidthLimits.Settings = attribute.PerClientBandwidthLimits.Settings
			trafficShapingRule.PerClientBandwidthLimits.BandwidthLimits.LimitDown = attribute.PerClientBandwidthLimits.BandwidthLimits.LimitDown
			trafficShapingRule.PerClientBandwidthLimits.BandwidthLimits.LimitUp = attribute.PerClientBandwidthLimits.BandwidthLimits.LimitUp
			if len(attribute.Definitions) > 0 {
				for _, attribute := range attribute.Definitions {
					var definition NetworksGroupPolicyResourceModelDefinition
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
	}

	if data.VlanTagging.VlanId.IsUnknown() {
		data.VlanTagging.VlanId = jsontypes.StringNull()
	}
	if data.ContentFiltering.AllowedUrlPatterns.Settings.IsUnknown() {
		data.ContentFiltering.AllowedUrlPatterns.Settings = jsontypes.StringNull()
	}
	if data.ContentFiltering.BlockedUrlCategories.Settings.IsUnknown() {
		data.ContentFiltering.BlockedUrlCategories.Settings = jsontypes.StringNull()
	}
	if data.ContentFiltering.BlockedUrlPatterns.Settings.IsUnknown() {
		data.ContentFiltering.BlockedUrlPatterns.Settings = jsontypes.StringNull()
	}

	return data, nil
}
