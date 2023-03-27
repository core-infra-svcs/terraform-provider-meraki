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
	Id                        jsontypes.String          `tfsdk:"id"`
	NetworkId                 jsontypes.String          `tfsdk:"network_id"`
	GroupPolicyId             jsontypes.String          `tfsdk:"group_policy_id"`
	Name                      jsontypes.String          `tfsdk:"name" json:"name"`
	SplashAuthSettings        jsontypes.String          `tfsdk:"splash_auth_settings"`
	BandWidth                 BandWidth                 `tfsdk:"bandwidth"`
	BonjourForwarding         BonjourForwarding         `tfsdk:"bonjour_forwarding"`
	ContentFiltering          ContentFiltering          `tfsdk:"content_filtering"`
	Scheduling                Scheduling                `tfsdk:"scheduling"`
	VlanTagging               VlanTagging               `tfsdk:"vlan_tagging"`
	FirewallAndTrafficShaping FirewallAndTrafficShaping `tfsdk:"firewall_and_traffic_shaping"`
}

type FirewallAndTrafficShaping struct {
	Settings            jsontypes.String      `tfsdk:"settings"`
	L3FirewallRules     []L3FirewallRules     `tfsdk:"l3_firewall_rules"`
	L7FirewallRules     []L7FirewallRules     `tfsdk:"l7_firewall_rules"`
	TrafficShapingRules []TrafficShapingRules `tfsdk:"traffic_shaping_rules"`
}

type L3FirewallRules struct {
	Comment  jsontypes.String `tfsdk:"comment"`
	DestCidr jsontypes.String `tfsdk:"dest_cidr"`
	DestPort jsontypes.String `tfsdk:"dest_port"`
	Policy   jsontypes.String `tfsdk:"policy"`
	Protocol jsontypes.String `tfsdk:"protocol"`
}

type L7FirewallRules struct {
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
	Settings        jsontypes.String `tfsdk:"settings"`
	BandwidthLimits BandwidthLimits  `tfsdk:"bandwidth_limits"`
}
type BandwidthLimits struct {
	LimitDown jsontypes.Int64 `tfsdk:"limit_down"`
	LimitUp   jsontypes.Int64 `tfsdk:"limit_up"`
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
			"bandwidth": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"settings": schema.StringAttribute{
						MarkdownDescription: "Settings Bandwidth",
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
								"priority": schema.StringAttribute{
									MarkdownDescription: "A string, indicating the priority level for packets bound to your rule. Can be 'low', 'normal' or 'high'.",
									Optional:            true,
									Computed:            true,
									CustomType:          jsontypes.StringType,
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

	} else {
		data.SplashAuthSettings = jsontypes.StringNull()
	}

	if !data.BandWidth.Settings.IsUnknown() {
		var a openApiClient.NetworksNetworkIdGroupPoliciesBandwidth
		a.SetSettings(data.BandWidth.Settings.ValueString())
		if data.BandWidth.BandwidthLimits.LimitUp != jsontypes.Int64Value(0) || data.BandWidth.BandwidthLimits.LimitDown != jsontypes.Int64Value(0) {
			var v openApiClient.NetworksNetworkIdGroupPoliciesBandwidthBandwidthLimits
			v.SetLimitDown(int32(data.BandWidth.BandwidthLimits.LimitDown.ValueInt64()))
			v.SetLimitUp(int32(data.BandWidth.BandwidthLimits.LimitUp.ValueInt64()))
			a.SetBandwidthLimits(v)
		} else {
			data.BandWidth.BandwidthLimits.LimitDown = jsontypes.Int64Null()
			data.BandWidth.BandwidthLimits.LimitUp = jsontypes.Int64Null()
		}
		createNetworkGroupPolicy.SetBandwidth(a)

	} else {
		data.BandWidth.Settings = jsontypes.StringNull()
		data.BandWidth.BandwidthLimits.LimitDown = jsontypes.Int64Null()
		data.BandWidth.BandwidthLimits.LimitUp = jsontypes.Int64Null()

	}

	if !data.BonjourForwarding.Settings.IsUnknown() {
		var b openApiClient.NetworksNetworkIdGroupPoliciesBonjourForwarding
		b.SetSettings(data.BonjourForwarding.Settings.ValueString())
		var r []openApiClient.NetworksNetworkIdGroupPoliciesBonjourForwardingRules
		for _, attribute := range data.BonjourForwarding.Rules {
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
		data.BonjourForwarding.Settings = jsontypes.StringNull()
	}

	var c openApiClient.NetworksNetworkIdGroupPoliciesContentFiltering

	if !data.ContentFiltering.AllowedUrlPatterns.Settings.IsUnknown() || len(data.ContentFiltering.AllowedUrlPatterns.Patterns) > 0 {
		var aup openApiClient.NetworksNetworkIdGroupPoliciesContentFilteringAllowedUrlPatterns
		aup.SetSettings(data.ContentFiltering.AllowedUrlPatterns.Settings.ValueString())
		aup.SetPatterns(data.ContentFiltering.AllowedUrlPatterns.Patterns)
		c.SetAllowedUrlPatterns(aup)

	} else {
		data.ContentFiltering.AllowedUrlPatterns.Settings = jsontypes.StringNull()
		data.ContentFiltering.AllowedUrlPatterns.Patterns = nil
	}
	if !data.ContentFiltering.BlockedUrlCategories.Settings.IsUnknown() || len(data.ContentFiltering.BlockedUrlCategories.Categories) > 0 {
		var buc openApiClient.NetworksNetworkIdGroupPoliciesContentFilteringBlockedUrlCategories
		buc.SetSettings(data.ContentFiltering.AllowedUrlPatterns.Settings.ValueString())
		buc.SetCategories(data.ContentFiltering.BlockedUrlCategories.Categories)
		c.SetBlockedUrlCategories(buc)

	} else {
		data.ContentFiltering.BlockedUrlCategories.Settings = jsontypes.StringNull()
		data.ContentFiltering.BlockedUrlCategories.Categories = nil
	}
	if !data.ContentFiltering.BlockedUrlPatterns.Settings.IsUnknown() || len(data.ContentFiltering.BlockedUrlPatterns.Patterns) > 0 {
		var bup openApiClient.NetworksNetworkIdGroupPoliciesContentFilteringBlockedUrlPatterns
		bup.SetSettings(data.ContentFiltering.AllowedUrlPatterns.Settings.ValueString())
		bup.SetPatterns(data.ContentFiltering.BlockedUrlPatterns.Patterns)
		c.SetBlockedUrlPatterns(bup)

	} else {
		data.ContentFiltering.BlockedUrlPatterns.Settings = jsontypes.StringNull()
		data.ContentFiltering.BlockedUrlPatterns.Patterns = nil
	}
	createNetworkGroupPolicy.SetContentFiltering(c)

	if !data.Scheduling.Enabled.IsUnknown() {
		var s openApiClient.NetworksNetworkIdGroupPoliciesScheduling
		s.SetEnabled(data.Scheduling.Enabled.ValueBool())
		var f openApiClient.NetworksNetworkIdGroupPoliciesSchedulingFriday
		f.SetActive(data.Scheduling.Friday.Active.ValueBool())
		s.SetFriday(f)
		var m openApiClient.NetworksNetworkIdGroupPoliciesSchedulingMonday
		m.SetActive(data.Scheduling.Monday.Active.ValueBool())
		s.SetMonday(m)
		var tu openApiClient.NetworksNetworkIdGroupPoliciesSchedulingTuesday
		tu.SetActive(data.Scheduling.Tuesday.Active.ValueBool())
		s.SetTuesday(tu)
		var w openApiClient.NetworksNetworkIdGroupPoliciesSchedulingWednesday
		w.SetActive(data.Scheduling.Wednesday.Active.ValueBool())
		s.SetWednesday(w)
		var th openApiClient.NetworksNetworkIdGroupPoliciesSchedulingThursday
		th.SetActive(data.Scheduling.Thursday.Active.ValueBool())
		s.SetThursday(th)
		var sa openApiClient.NetworksNetworkIdGroupPoliciesSchedulingSaturday
		sa.SetActive(data.Scheduling.Saturday.Active.ValueBool())
		s.SetSaturday(sa)
		var su openApiClient.NetworksNetworkIdGroupPoliciesSchedulingSunday
		su.SetActive(data.Scheduling.Sunday.Active.ValueBool())
		s.SetSunday(su)
		createNetworkGroupPolicy.SetScheduling(s)
	} else {
		data.Scheduling.Enabled = jsontypes.BoolNull()
	}

	if !data.VlanTagging.Settings.IsUnknown() {
		var v openApiClient.NetworksNetworkIdGroupPoliciesVlanTagging
		v.SetSettings(data.VlanTagging.Settings.ValueString())
		v.SetVlanId(data.VlanTagging.VlanId.ValueString())
		createNetworkGroupPolicy.SetVlanTagging(v)
	} else {
		data.VlanTagging.Settings = jsontypes.StringNull()
	}

	if !data.FirewallAndTrafficShaping.Settings.IsUnknown() {
		var f openApiClient.NetworksNetworkIdGroupPoliciesFirewallAndTrafficShaping
		f.SetSettings(data.FirewallAndTrafficShaping.Settings.ValueString())
		if len(data.FirewallAndTrafficShaping.L3FirewallRules) > 0 {
			var l3s []openApiClient.NetworksNetworkIdGroupPoliciesFirewallAndTrafficShapingL3FirewallRules
			for _, attribute := range data.FirewallAndTrafficShaping.L3FirewallRules {
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
			data.FirewallAndTrafficShaping.L3FirewallRules = nil
		}
		if len(data.FirewallAndTrafficShaping.L7FirewallRules) > 0 {
			var l7s []openApiClient.NetworksNetworkIdGroupPoliciesFirewallAndTrafficShapingL7FirewallRules
			for _, attribute := range data.FirewallAndTrafficShaping.L7FirewallRules {
				var l7 openApiClient.NetworksNetworkIdGroupPoliciesFirewallAndTrafficShapingL7FirewallRules

				l7.SetValue(attribute.Value.ValueString())
				l7.SetPolicy(attribute.Policy.ValueString())
				l7.SetType(attribute.Type.ValueString())
				l7s = append(l7s, l7)
			}
			f.SetL7FirewallRules(l7s)
		} else {
			data.FirewallAndTrafficShaping.L7FirewallRules = nil
		}
		if len(data.FirewallAndTrafficShaping.TrafficShapingRules) > 0 {
			var tfs []openApiClient.NetworksNetworkIdGroupPoliciesFirewallAndTrafficShapingTrafficShapingRules
			for _, attribute := range data.FirewallAndTrafficShaping.TrafficShapingRules {
				var tf openApiClient.NetworksNetworkIdGroupPoliciesFirewallAndTrafficShapingTrafficShapingRules
				tf.SetDscpTagValue(int32(attribute.DscpTagValue.ValueInt64()))
				tf.SetPcpTagValue(int32(attribute.PcpTagValue.ValueInt64()))
				tf.SetPriority(attribute.Priority.ValueString())
				var perclientBamdWidthLimits openApiClient.NetworksNetworkIdApplianceTrafficShapingRulesPerClientBandwidthLimits
				var bandwidthLimits openApiClient.NetworksNetworkIdApplianceTrafficShapingRulesPerClientBandwidthLimitsBandwidthLimits
				bandwidthLimits.SetLimitDown(int32(attribute.PerClientBandwidthLimits.BandwidthLimits.LimitDown.ValueInt64()))
				bandwidthLimits.SetLimitUp(int32(attribute.PerClientBandwidthLimits.BandwidthLimits.LimitUp.ValueInt64()))
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
			data.FirewallAndTrafficShaping.TrafficShapingRules = nil
		}
		createNetworkGroupPolicy.SetFirewallAndTrafficShaping(f)
	} else {
		data.FirewallAndTrafficShaping.Settings = jsontypes.StringNull()
	}

	_, httpResp, err := r.client.NetworksApi.CreateNetworkGroupPolicy(ctx, data.NetworkId.ValueString()).CreateNetworkGroupPolicy(createNetworkGroupPolicy).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	//fmt.Println(httpResp.Body)

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

	/*
		// Save data into Terraform state
		if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
			resp.Diagnostics.AddError(
				"JSON decoding error",
				fmt.Sprintf("%v\n", err.Error()),
			)
			return
		}
	*/

	fmt.Println(data)

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

	_, httpResp, err := r.client.NetworksApi.GetNetworkGroupPolicy(ctx, data.NetworkId.ValueString(), data.GroupPolicyId.ValueString()).Execute()
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
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	data.Id = jsontypes.StringValue("example-id")
	fmt.Println(data)

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

	} else {
		data.SplashAuthSettings = jsontypes.StringNull()
	}

	if !data.BandWidth.Settings.IsUnknown() {
		var a openApiClient.NetworksNetworkIdGroupPoliciesBandwidth
		a.SetSettings(data.BandWidth.Settings.ValueString())
		if data.BandWidth.BandwidthLimits.LimitUp != jsontypes.Int64Value(0) || data.BandWidth.BandwidthLimits.LimitDown != jsontypes.Int64Value(0) {
			var v openApiClient.NetworksNetworkIdGroupPoliciesBandwidthBandwidthLimits
			v.SetLimitDown(int32(data.BandWidth.BandwidthLimits.LimitDown.ValueInt64()))
			v.SetLimitUp(int32(data.BandWidth.BandwidthLimits.LimitUp.ValueInt64()))
			a.SetBandwidthLimits(v)
		} else {
			data.BandWidth.BandwidthLimits.LimitDown = jsontypes.Int64Null()
			data.BandWidth.BandwidthLimits.LimitUp = jsontypes.Int64Null()
		}
		updateNetworkGroupPolicy.SetBandwidth(a)

	} else {
		data.BandWidth.Settings = jsontypes.StringNull()
		data.BandWidth.BandwidthLimits.LimitDown = jsontypes.Int64Null()
		data.BandWidth.BandwidthLimits.LimitUp = jsontypes.Int64Null()

	}

	if !data.BonjourForwarding.Settings.IsUnknown() {
		var b openApiClient.NetworksNetworkIdGroupPoliciesBonjourForwarding
		b.SetSettings(data.BonjourForwarding.Settings.ValueString())
		var r []openApiClient.NetworksNetworkIdGroupPoliciesBonjourForwardingRules
		for _, attribute := range data.BonjourForwarding.Rules {
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
		data.BonjourForwarding.Settings = jsontypes.StringNull()
	}

	var c openApiClient.NetworksNetworkIdGroupPoliciesContentFiltering

	if !data.ContentFiltering.AllowedUrlPatterns.Settings.IsUnknown() || len(data.ContentFiltering.AllowedUrlPatterns.Patterns) > 0 {
		var aup openApiClient.NetworksNetworkIdGroupPoliciesContentFilteringAllowedUrlPatterns
		aup.SetSettings(data.ContentFiltering.AllowedUrlPatterns.Settings.ValueString())
		aup.SetPatterns(data.ContentFiltering.AllowedUrlPatterns.Patterns)
		c.SetAllowedUrlPatterns(aup)

	} else {
		data.ContentFiltering.AllowedUrlPatterns.Settings = jsontypes.StringNull()
		data.ContentFiltering.AllowedUrlPatterns.Patterns = nil
	}
	if !data.ContentFiltering.BlockedUrlCategories.Settings.IsUnknown() || len(data.ContentFiltering.BlockedUrlCategories.Categories) > 0 {
		var buc openApiClient.NetworksNetworkIdGroupPoliciesContentFilteringBlockedUrlCategories
		buc.SetSettings(data.ContentFiltering.AllowedUrlPatterns.Settings.ValueString())
		buc.SetCategories(data.ContentFiltering.BlockedUrlCategories.Categories)
		c.SetBlockedUrlCategories(buc)

	} else {
		data.ContentFiltering.BlockedUrlCategories.Settings = jsontypes.StringNull()
		data.ContentFiltering.BlockedUrlCategories.Categories = nil
	}
	if !data.ContentFiltering.BlockedUrlPatterns.Settings.IsUnknown() || len(data.ContentFiltering.BlockedUrlPatterns.Patterns) > 0 {
		var bup openApiClient.NetworksNetworkIdGroupPoliciesContentFilteringBlockedUrlPatterns
		bup.SetSettings(data.ContentFiltering.AllowedUrlPatterns.Settings.ValueString())
		bup.SetPatterns(data.ContentFiltering.BlockedUrlPatterns.Patterns)
		c.SetBlockedUrlPatterns(bup)

	} else {
		data.ContentFiltering.BlockedUrlPatterns.Settings = jsontypes.StringNull()
		data.ContentFiltering.BlockedUrlPatterns.Patterns = nil
	}
	updateNetworkGroupPolicy.SetContentFiltering(c)

	if !data.Scheduling.Enabled.IsUnknown() {
		var s openApiClient.NetworksNetworkIdGroupPoliciesScheduling
		s.SetEnabled(data.Scheduling.Enabled.ValueBool())
		var f openApiClient.NetworksNetworkIdGroupPoliciesSchedulingFriday
		f.SetActive(data.Scheduling.Friday.Active.ValueBool())
		s.SetFriday(f)
		var m openApiClient.NetworksNetworkIdGroupPoliciesSchedulingMonday
		m.SetActive(data.Scheduling.Monday.Active.ValueBool())
		s.SetMonday(m)
		var tu openApiClient.NetworksNetworkIdGroupPoliciesSchedulingTuesday
		tu.SetActive(data.Scheduling.Tuesday.Active.ValueBool())
		s.SetTuesday(tu)
		var w openApiClient.NetworksNetworkIdGroupPoliciesSchedulingWednesday
		w.SetActive(data.Scheduling.Wednesday.Active.ValueBool())
		s.SetWednesday(w)
		var th openApiClient.NetworksNetworkIdGroupPoliciesSchedulingThursday
		th.SetActive(data.Scheduling.Thursday.Active.ValueBool())
		s.SetThursday(th)
		var sa openApiClient.NetworksNetworkIdGroupPoliciesSchedulingSaturday
		sa.SetActive(data.Scheduling.Saturday.Active.ValueBool())
		s.SetSaturday(sa)
		var su openApiClient.NetworksNetworkIdGroupPoliciesSchedulingSunday
		su.SetActive(data.Scheduling.Sunday.Active.ValueBool())
		s.SetSunday(su)
		updateNetworkGroupPolicy.SetScheduling(s)
	} else {
		data.Scheduling.Enabled = jsontypes.BoolNull()
	}

	if !data.FirewallAndTrafficShaping.Settings.IsUnknown() {
		var f openApiClient.NetworksNetworkIdGroupPoliciesFirewallAndTrafficShaping
		f.SetSettings(data.FirewallAndTrafficShaping.Settings.ValueString())
		if len(data.FirewallAndTrafficShaping.L3FirewallRules) > 0 {
			var l3s []openApiClient.NetworksNetworkIdGroupPoliciesFirewallAndTrafficShapingL3FirewallRules
			for _, attribute := range data.FirewallAndTrafficShaping.L3FirewallRules {
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
			data.FirewallAndTrafficShaping.L3FirewallRules = nil
		}
		if len(data.FirewallAndTrafficShaping.L7FirewallRules) > 0 {
			var l7s []openApiClient.NetworksNetworkIdGroupPoliciesFirewallAndTrafficShapingL7FirewallRules
			for _, attribute := range data.FirewallAndTrafficShaping.L7FirewallRules {
				var l7 openApiClient.NetworksNetworkIdGroupPoliciesFirewallAndTrafficShapingL7FirewallRules

				l7.SetValue(attribute.Value.ValueString())
				l7.SetPolicy(attribute.Policy.ValueString())
				l7.SetType(attribute.Type.ValueString())
				l7s = append(l7s, l7)
			}
			f.SetL7FirewallRules(l7s)
		} else {
			data.FirewallAndTrafficShaping.L7FirewallRules = nil
		}
		if len(data.FirewallAndTrafficShaping.TrafficShapingRules) > 0 {
			var tfs []openApiClient.NetworksNetworkIdGroupPoliciesFirewallAndTrafficShapingTrafficShapingRules
			for _, attribute := range data.FirewallAndTrafficShaping.TrafficShapingRules {
				var tf openApiClient.NetworksNetworkIdGroupPoliciesFirewallAndTrafficShapingTrafficShapingRules
				tf.SetDscpTagValue(int32(attribute.DscpTagValue.ValueInt64()))
				tf.SetPcpTagValue(int32(attribute.PcpTagValue.ValueInt64()))
				tf.SetPriority(attribute.Priority.ValueString())
				var perclientBamdWidthLimits openApiClient.NetworksNetworkIdApplianceTrafficShapingRulesPerClientBandwidthLimits
				var bandwidthLimits openApiClient.NetworksNetworkIdApplianceTrafficShapingRulesPerClientBandwidthLimitsBandwidthLimits
				bandwidthLimits.SetLimitDown(int32(attribute.PerClientBandwidthLimits.BandwidthLimits.LimitDown.ValueInt64()))
				bandwidthLimits.SetLimitUp(int32(attribute.PerClientBandwidthLimits.BandwidthLimits.LimitUp.ValueInt64()))
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
			data.FirewallAndTrafficShaping.TrafficShapingRules = nil
		}
		updateNetworkGroupPolicy.SetFirewallAndTrafficShaping(f)
	} else {
		data.FirewallAndTrafficShaping.Settings = jsontypes.StringNull()
	}

	if !data.VlanTagging.Settings.IsUnknown() {
		var v openApiClient.NetworksNetworkIdGroupPoliciesVlanTagging
		v.SetSettings(data.VlanTagging.Settings.ValueString())
		v.SetVlanId(data.VlanTagging.VlanId.ValueString())
		updateNetworkGroupPolicy.SetVlanTagging(v)
	} else {
		data.VlanTagging.Settings = jsontypes.StringNull()
	}

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

	fmt.Println(httpResp.Body)

	// Save data into Terraform state
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	data.Id = jsontypes.StringValue("example-id")
	fmt.Println(data)

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
