package policy

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
	"io"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &DataSource{}

func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

type DataSource struct {
	client *openApiClient.APIClient
}

func (d *DataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_group_policies"
}

func (d *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage the group policy's in this network",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Example identifier",
				Computed:            true,
				CustomType:          types.StringType,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Network Id",
				Required:            true,
				CustomType:          types.StringType,
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
							CustomType:          types.StringType,
							Validators: []validator.String{
								stringvalidator.LengthBetween(1, 31),
							},
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of Group Policy",
							Required:            true,
							CustomType:          types.StringType,
						},
						"splash_auth_settings": schema.StringAttribute{
							MarkdownDescription: "Whether clients bound to your policy will bypass splash authorization or behave according to the network's rules. Can be one of 'network default' or 'bypass'. Only available if your network has a wireless configuration",
							Optional:            true,
							Computed:            true,
							CustomType:          types.StringType,
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
									CustomType:          types.StringType,
								},
								"bandwidth_limits": schema.SingleNestedAttribute{
									Optional: true,
									Computed: true,

									Attributes: map[string]schema.Attribute{
										"limit_down": schema.Int64Attribute{
											MarkdownDescription: "The maximum download limit (integer, in Kbps).",
											Optional:            true,
											Computed:            true,
											CustomType:          types.Int64Type,
										},
										"limit_up": schema.Int64Attribute{
											MarkdownDescription: "The maximum upload limit (integer, in Kbps).",
											Optional:            true,
											Computed:            true,
											CustomType:          types.Int64Type,
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
											CustomType:          types.StringType,
										},
										"patterns": schema.ListAttribute{
											MarkdownDescription: "A list of URL patterns that are allowed",
											ElementType:         types.StringType,
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
											CustomType:          types.StringType,
										},
										"categories": schema.ListAttribute{
											MarkdownDescription: "A list of URL categories to block",
											ElementType:         types.StringType,
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
											CustomType:          types.StringType,
										},
										"patterns": schema.ListAttribute{
											MarkdownDescription: "A list of URL patterns that are blocked",
											ElementType:         types.StringType,
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
									CustomType:          types.StringType,
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
												CustomType:          types.StringType,
											},
											"vlan_id": schema.StringAttribute{
												MarkdownDescription: "The ID of the service VLAN. Required.",
												Optional:            true,
												Computed:            true,
												CustomType:          types.StringType,
											},
											"services": schema.ListAttribute{
												ElementType: types.StringType,
												Optional:    true,
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
									CustomType:          types.StringType,
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
												CustomType:          types.StringType,
											},
											"dest_cidr": schema.StringAttribute{
												MarkdownDescription: "Destination IP address (in IP or CIDR notation), a fully-qualified domain name (FQDN, if your network supports it) or 'any'.",
												Optional:            true,
												Computed:            true,
												CustomType:          types.StringType,
											},
											"dest_port": schema.StringAttribute{
												MarkdownDescription: "Destination port (integer in the range 1-65535), a port range (e.g. 8080-9090), or 'any'",
												Optional:            true,
												Computed:            true,
												CustomType:          types.StringType,
											},
											"policy": schema.StringAttribute{
												MarkdownDescription: "'allow' or 'deny' traffic specified by this rule",
												Optional:            true,
												Computed:            true,
												CustomType:          types.StringType,
											},
											"protocol": schema.StringAttribute{
												MarkdownDescription: "The type of protocol (must be 'tcp', 'udp', 'icmp', 'icmp6' or 'any')",
												Optional:            true,
												Computed:            true,
												CustomType:          types.StringType,
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
												CustomType:          types.StringType,
											},
											"policy": schema.StringAttribute{
												MarkdownDescription: "The policy applied to matching traffic. Must be 'deny'.",
												Optional:            true,
												Computed:            true,
												CustomType:          types.StringType,
											},
											"type": schema.StringAttribute{
												MarkdownDescription: "Type of the L7 Rule. Must be 'application', 'applicationCategory', 'host', 'port' or 'ipRange'",
												Optional:            true,
												Computed:            true,
												CustomType:          types.StringType,
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
												CustomType:          types.Int64Type,
											},
											"pcp_tag_value": schema.Int64Attribute{
												MarkdownDescription: "The PCP tag applied by your rule. Can be 0 (lowest priority) through 7 (highest priority). null means Do not set PCP tag.",
												Optional:            true,
												Computed:            true,
												CustomType:          types.Int64Type,
											},
											"per_client_bandwidth_limits": schema.SingleNestedAttribute{
												Optional: true,
												Computed: true,
												Attributes: map[string]schema.Attribute{
													"settings": schema.StringAttribute{
														MarkdownDescription: "How bandwidth limits are applied by your rule. Can be one of 'network default', 'ignore' or 'custom'.",
														Optional:            true,
														Computed:            true,
														CustomType:          types.StringType,
													},
													"bandwidth_limits": schema.SingleNestedAttribute{
														Optional: true,
														Computed: true,
														Attributes: map[string]schema.Attribute{
															"limit_down": schema.Int64Attribute{
																MarkdownDescription: "The maximum download limit (integer, in Kbps).",
																Optional:            true,
																Computed:            true,
																CustomType:          types.Int64Type,
															},
															"limit_up": schema.Int64Attribute{
																MarkdownDescription: "The maximum upload limit (integer, in Kbps).",
																Optional:            true,
																Computed:            true,
																CustomType:          types.Int64Type,
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
															CustomType:          types.StringType,
														},
														"value": schema.StringAttribute{
															MarkdownDescription: "If type is host, port, ipRange or localNet then value must be a string matching either a hostname (e.g. somesite.com) a port (e.g. 8080) or an IP range (192.1.0.0, 192.1.0.0/16, or 10.1.0.0/16:80). localNet also supports CIDR notation excluding custom ports If type is 'application' or 'applicationCategory', then value must be an object with the structure { id: meraki:layer7/... }, where id is the application category or application ID (for a list of IDs for your network, use the trafficShaping/applicationCategories endpoint)",
															Optional:            true,
															Computed:            true,
															CustomType:          types.StringType,
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
									CustomType:          types.StringType,
								},
								"vlan_id": schema.StringAttribute{
									MarkdownDescription: "The ID of the vlan you want to tag. This only applies if 'settings' is set to 'custom'.",
									Optional:            true,
									Computed:            true,
									CustomType:          types.StringType,
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
									CustomType:          types.BoolType,
								},
								"friday": schema.SingleNestedAttribute{
									Optional: true,
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"from": schema.StringAttribute{
											MarkdownDescription: "The time, from '00:00' to '24:00'. Must be less than the time specified in 'to'. Defaults to '00:00'. Only 30 minute increments are allowed.",
											Optional:            true,
											Computed:            true,
											CustomType:          types.StringType,
										},
										"to": schema.StringAttribute{
											MarkdownDescription: "The time, from '00:00' to '24:00'. Must be greater than the time specified in 'from'. Defaults to '24:00'. Only 30 minute increments are allowed.",
											Optional:            true,
											Computed:            true,
											CustomType:          types.StringType,
										},
										"active": schema.BoolAttribute{
											MarkdownDescription: "Whether the schedule is active (true) or inactive (false) during the time specified between 'from' and 'to'. Defaults to true.",
											Optional:            true,
											Computed:            true,
											CustomType:          types.BoolType,
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
											CustomType:          types.StringType,
										},
										"to": schema.StringAttribute{
											MarkdownDescription: "The time, from '00:00' to '24:00'. Must be greater than the time specified in 'from'. Defaults to '24:00'. Only 30 minute increments are allowed.",
											Optional:            true,
											Computed:            true,
											CustomType:          types.StringType,
										},
										"active": schema.BoolAttribute{
											MarkdownDescription: "Whether the schedule is active (true) or inactive (false) during the time specified between 'from' and 'to'. Defaults to true.",
											Optional:            true,
											Computed:            true,
											CustomType:          types.BoolType,
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
											CustomType:          types.StringType,
										},
										"to": schema.StringAttribute{
											MarkdownDescription: "The time, from '00:00' to '24:00'. Must be greater than the time specified in 'from'. Defaults to '24:00'. Only 30 minute increments are allowed.",
											Optional:            true,
											Computed:            true,
											CustomType:          types.StringType,
										},
										"active": schema.BoolAttribute{
											MarkdownDescription: "Whether the schedule is active (true) or inactive (false) during the time specified between 'from' and 'to'. Defaults to true.",
											Optional:            true,
											Computed:            true,
											CustomType:          types.BoolType,
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
											CustomType:          types.StringType,
										},
										"to": schema.StringAttribute{
											MarkdownDescription: "The time, from '00:00' to '24:00'. Must be greater than the time specified in 'from'. Defaults to '24:00'. Only 30 minute increments are allowed.",
											Optional:            true,
											Computed:            true,
											CustomType:          types.StringType,
										},
										"active": schema.BoolAttribute{
											MarkdownDescription: "Whether the schedule is active (true) or inactive (false) during the time specified between 'from' and 'to'. Defaults to true.",
											Optional:            true,
											Computed:            true,
											CustomType:          types.BoolType,
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
											CustomType:          types.StringType,
										},
										"to": schema.StringAttribute{
											MarkdownDescription: "The time, from '00:00' to '24:00'. Must be greater than the time specified in 'from'. Defaults to '24:00'. Only 30 minute increments are allowed.",
											Optional:            true,
											Computed:            true,
											CustomType:          types.StringType,
										},
										"active": schema.BoolAttribute{
											MarkdownDescription: "Whether the schedule is active (true) or inactive (false) during the time specified between 'from' and 'to'. Defaults to true.",
											Optional:            true,
											Computed:            true,
											CustomType:          types.BoolType,
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
											CustomType:          types.StringType,
										},
										"to": schema.StringAttribute{
											MarkdownDescription: "The time, from '00:00' to '24:00'. Must be greater than the time specified in 'from'. Defaults to '24:00'. Only 30 minute increments are allowed.",
											Optional:            true,
											Computed:            true,
											CustomType:          types.StringType,
										},
										"active": schema.BoolAttribute{
											MarkdownDescription: "Whether the schedule is active (true) or inactive (false) during the time specified between 'from' and 'to'. Defaults to true.",
											Optional:            true,
											Computed:            true,
											CustomType:          types.BoolType,
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
											CustomType:          types.StringType,
										},
										"to": schema.StringAttribute{
											MarkdownDescription: "The time, from '00:00' to '24:00'. Must be greater than the time specified in 'from'. Defaults to '24:00'. Only 30 minute increments are allowed.",
											Optional:            true,
											Computed:            true,
											CustomType:          types.StringType,
										},
										"active": schema.BoolAttribute{
											MarkdownDescription: "Whether the schedule is active (true) or inactive (false) during the time specified between 'from' and 'to'. Defaults to true.",
											Optional:            true,
											Computed:            true,
											CustomType:          types.BoolType,
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

func (d *DataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state datasourceModel

	// Read configuration into state
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call API to fetch list of group policies
	inlineResp, httpResp, err := d.client.NetworksApi.GetNetworkGroupPolicies(ctx, state.NetworkId.ValueString()).Execute()
	if err != nil {
		var responseBody string
		if httpResp != nil && httpResp.Body != nil {
			bodyBytes, readErr := io.ReadAll(httpResp.Body)
			if readErr == nil {
				responseBody = string(bodyBytes)
			}
		}

		resp.Diagnostics.AddError(
			"Error reading network group policies",
			fmt.Sprintf("Could not read group policies, unexpected error: %s\nHTTP Response: %v\nResponse Body: %s", err, httpResp, responseBody),
		)
		return
	}

	// Map API response to resourceModel list
	var policies []attr.Value
	for _, inline := range inlineResp {
		var policy dataSourceListModel
		diags = ToTerraformStateData(ctx, &policy, inline)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Convert each policy to a Terraform Object
		policyObj, policyObjDiags := types.ObjectValueFrom(ctx, resourceModelAttrs(), policy)
		resp.Diagnostics.Append(policyObjDiags...)
		if resp.Diagnostics.HasError() {
			return
		}

		policies = append(policies, policyObj)
	}

	// Create a Terraform ListValue from the policy objects
	policyList, diags := types.ListValue(types.ObjectType{AttrTypes: resourceModelAttrs()}, policies)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the state
	state.List = policyList

	state.Id = state.NetworkId
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Log completion
	tflog.Trace(ctx, "Finished reading network group policies data source")
}
