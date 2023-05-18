package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &OrganizationsAdminsDataSource{}

func NewNetworkGroupPolicysDataSource() datasource.DataSource {
	return &NetworkGroupPolicysDataSource{}
}

// NetworkGroupPolicysDataSource defines the data source implementation.
type NetworkGroupPolicysDataSource struct {
	client *openApiClient.APIClient
}

// NetworkGroupPolicysDataSourceModel describes the data source data model.
type NetworkGroupPolicysDataSourceModel struct {
	Id        jsontypes.String                    `tfsdk:"id"`
	NetworkId jsontypes.String                    `tfsdk:"network_id"`
	List      []NetworkGroupPolicyDataSourceModel `tfsdk:"list"`
}

// NetworkGroupPolicyDataSourceModel describes the data source data model.
type NetworkGroupPolicyDataSourceModel struct {
	GroupPolicyId             jsontypes.String                    `tfsdk:"group_policy_id" json:"groupPolicyId"`
	Name                      jsontypes.String                    `tfsdk:"name" json:"name"`
	SplashAuthSettings        jsontypes.String                    `tfsdk:"splash_auth_settings" json:"splashAuthSettings"`
	BandWidth                 BandwidthDataSource                 `tfsdk:"bandwidth" json:"bandwidth"`
	BonjourForwarding         BonjourForwardingDatasource         `tfsdk:"bonjour_forwarding" json:"bonjourForwarding"`
	Scheduling                SchedulingDatasource                `tfsdk:"scheduling" json:"scheduling"`
	FirewallAndTrafficShaping FirewallAndTrafficShapingDatasource `tfsdk:"firewall_and_traffic_shaping" json:"firewallAndTrafficShaping"`
	VlanTagging               VlanTaggingDatasource               `tfsdk:"vlan_tagging" json:"vlanTagging"`
	ContentFiltering          ContentFilteringDatasource          `tfsdk:"content_filtering" json:"contentFiltering"`
}

type BandwidthDataSource struct {
	BandwidthLimitsDataSource BandwidthLimitsDataSource `tfsdk:"bandwidth_limits" json:"bandwidthLimits"`
	Settings                  jsontypes.String          `tfsdk:"settings" json:"settings"`
}

type BandwidthLimitsDataSource struct {
	LimitUp   jsontypes.Int64 `tfsdk:"limit_up" json:"limitUp"`
	LimitDown jsontypes.Int64 `tfsdk:"limit_down" json:"limitDown"`
}

type BonjourForwardingDatasource struct {
	BonjourForwardingSettings string           `tfsdk:"settings" json:"settings"`
	BonjourForwardingRules    []RuleDatasource `tfsdk:"rules" json:"rules"`
}

type RuleDatasource struct {
	Description jsontypes.String `tfsdk:"description" json:"description"`
	VlanId      jsontypes.String `tfsdk:"vlan_id" json:"vlanId"`
	Services    []string         `tfsdk:"services" json:"services"`
}

type SchedulingDatasource struct {
	Enabled   jsontypes.Bool     `tfsdk:"enabled" json:"enabled"`
	Friday    ScheduleDatasource `tfsdk:"friday" json:"friday"`
	Monday    ScheduleDatasource `tfsdk:"monday" json:"monday"`
	Saturday  ScheduleDatasource `tfsdk:"saturday" json:"saturday"`
	Sunday    ScheduleDatasource `tfsdk:"sunday" json:"sunday"`
	Thursday  ScheduleDatasource `tfsdk:"thursday" json:"thursday"`
	Tuesday   ScheduleDatasource `tfsdk:"tuesday" json:"tuesday"`
	Wednesday ScheduleDatasource `tfsdk:"wednesday" json:"wednesday"`
}

type ScheduleDatasource struct {
	From   jsontypes.String `tfsdk:"from" json:"from"`
	To     jsontypes.String `tfsdk:"to" json:"to"`
	Active jsontypes.Bool   `tfsdk:"active" json:"active"`
}

type VlanTaggingDatasource struct {
	Settings jsontypes.String `tfsdk:"settings" json:"settings"`
	VlanId   jsontypes.String `tfsdk:"vlan_id" json:"vlanId"`
}

type ContentFilteringDatasource struct {
	AllowedUrlPatterns   AllowedUrlPatternsDatasource   `tfsdk:"allowed_url_patterns" json:"allowedUrlPatterns"`
	BlockedUrlCategories BlockedUrlCategoriesDatasource `tfsdk:"blocked_url_categories" json:"blockedUrlCategories"`
	BlockedUrlPatterns   BlockedUrlPatternsDatasource   `tfsdk:"blocked_url_patterns" json:"blockedUrlPatterns"`
}

type AllowedUrlPatternsDatasource struct {
	Settings jsontypes.String `tfsdk:"settings" json:"settings"`
	Patterns []string         `tfsdk:"patterns" json:"patterns"`
}
type BlockedUrlCategoriesDatasource struct {
	Settings   jsontypes.String `tfsdk:"settings" json:"settings"`
	Categories []string         `tfsdk:"categories" json:"categories"`
}
type BlockedUrlPatternsDatasource struct {
	Settings jsontypes.String `tfsdk:"settings" json:"settings"`
	Patterns []string         `tfsdk:"patterns" json:"patterns"`
}

type FirewallAndTrafficShapingDatasource struct {
	Settings            jsontypes.String               `tfsdk:"settings" json:"settings"`
	L3FirewallRules     []L3FirewallRuleDatasource     `tfsdk:"l3_firewall_rules" json:"l3FirewallRules"`
	L7FirewallRules     []L7FirewallRuleDatasource     `tfsdk:"l7_firewall_rules" json:"l7FirewallRules"`
	TrafficShapingRules []TrafficShapingRuleDatasource `tfsdk:"traffic_shaping_rules" json:"trafficShapingRules"`
}

type L3FirewallRuleDatasource struct {
	Comment  jsontypes.String `tfsdk:"comment" json:"comment"`
	DestCidr jsontypes.String `tfsdk:"dest_cidr" json:"destCidr"`
	DestPort jsontypes.String `tfsdk:"dest_port" json:"destPort"`
	Policy   jsontypes.String `tfsdk:"policy" json:"policy"`
	Protocol jsontypes.String `tfsdk:"protocol" json:"protocol"`
}

type L7FirewallRuleDatasource struct {
	Value  jsontypes.String `tfsdk:"value" json:"value"`
	Type   jsontypes.String `tfsdk:"type" json:"type"`
	Policy jsontypes.String `tfsdk:"policy" json:"policy"`
}

type TrafficShapingRuleDatasource struct {
	DscpTagValue             jsontypes.Int64                    `tfsdk:"dscp_tag_value" json:"dscpTagValue"`
	PcpTagValue              jsontypes.Int64                    `tfsdk:"pcp_tag_value" json:"pcpTagValue"`
	PerClientBandwidthLimits PerClientBandwidthLimitsDataSource `tfsdk:"per_client_bandwidth_limits" json:"perClientBandwidthLimits"`
	Definitions              []DefinitionDatasource             `tfsdk:"definitions" json:"definitions"`
}

type PerClientBandwidthLimitsDataSource struct {
	BandwidthLimitsDataSource BandwidthLimitsDataSource `tfsdk:"bandwidth_limits" json:"bandwidthLimits"`
	Settings                  jsontypes.String          `tfsdk:"settings" json:"settings"`
}

type DefinitionDatasource struct {
	Value jsontypes.String `tfsdk:"value" json:"value"`
	Type  jsontypes.String `tfsdk:"type" json:"type"`
}

func (d *NetworkGroupPolicysDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_network_group_policys"
}

func (d *NetworkGroupPolicysDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List the group policy's in this network",

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
			},
			"list": schema.ListNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"group_policy_id": schema.StringAttribute{
							MarkdownDescription: "Group Policy ID",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
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
					},
				},
			},
		},
	}
}

func (d *NetworkGroupPolicysDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*openApiClient.APIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *NetworkGroupPolicysDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var data *NetworkGroupPolicysDataSourceModel
	var model *NetworkGroupPolicyDataSourceModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	inlineResp, httpResp, err := d.client.NetworksApi.GetNetworkGroupPolicies(ctx, data.NetworkId.ValueString()).Execute()
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
	for i := 0; i < len(inlineResp); i++ {
		jsonData, _ := json.Marshal(inlineResp[i])
		json.Unmarshal(jsonData, &model)
		data.List = append(data.List, *model)
	}

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}
