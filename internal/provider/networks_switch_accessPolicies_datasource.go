package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &NetworksSwitchAccessPoliciesDataSource{}

func NewNetworksSwitchAccessPoliciesDataSource() datasource.DataSource {
	return &NetworksSwitchAccessPoliciesDataSource{}
}

// NetworksSwitchAccessPoliciesDataSource defines the data source implementation.
type NetworksSwitchAccessPoliciesDataSource struct {
	client *openApiClient.APIClient
}

// NetworksSwitchesAccessPoliciesDataSourceModel describes the data source data model.
type NetworksSwitchesAccessPoliciesDataSourceModel struct {
	Id        types.String                                  `tfsdk:"id"`
	NetworkId jsontypes.String                              `tfsdk:"network_id"`
	List      []NetworksSwitchAccessPoliciesDataSourceModel `tfsdk:"list"`
}

type NetworkSwitchAccessPolicyRadiusServersDataSourceModelRules struct {
	Host   jsontypes.String `tfsdk:"host"`
	Port   jsontypes.Int64  `tfsdk:"port"`
	Secret jsontypes.String `tfsdk:"secret"`
}

// NetworksSwitchAccessPoliciesDataSourceModel describes the data source data model.
type NetworksSwitchAccessPoliciesDataSourceModel struct {
	Id                             jsontypes.String                                             `tfsdk:"id"`
	AccessPolicyType               jsontypes.String                                             `tfsdk:"access_policy_type"`
	Dot1xControlDirection          jsontypes.String                                             `tfsdk:"dot_1x_control_direction"`
	GuestVLANId                    jsontypes.Int64                                              `tfsdk:"guest_vlan_id"`
	HostMode                       jsontypes.String                                             `tfsdk:"host_mode"`
	IncreaseAccessSpeed            jsontypes.Bool                                               `tfsdk:"increase_access_speed"`
	Name                           jsontypes.String                                             `tfsdk:"name"`
	RadiusAccountingEnabled        jsontypes.Bool                                               `tfsdk:"radius_accounting_enabled"`
	RadiusAccountingServers        []NetworkSwitchAccessPolicyRadiusServersDataSourceModelRules `tfsdk:"radius_accounting_servers"`
	RadiusCoaSupportEnabled        jsontypes.Bool                                               `tfsdk:"radius_coa_support_enabled"`
	RadiusGroupAttribute           jsontypes.String                                             `tfsdk:"radius_group_attribute"`
	RadiusServers                  []NetworkSwitchAccessPolicyRadiusServersDataSourceModelRules `tfsdk:"radius_servers"`
	RadiusTestingEnabled           jsontypes.Bool                                               `tfsdk:"radius_testing_enabled"`
	RadiusCriticalAuth             types.Object                                                 `tfsdk:"radius_critical_auth"`
	RadiusFailedAuthVlanId         jsontypes.Int64                                              `tfsdk:"radius_failed_auth_vlan_id"`
	RadiusReauthenticationInterval jsontypes.Int64                                              `tfsdk:"radius_re_authentication_interval"`
	UrlRedirectWalledGardenEnabled jsontypes.Bool                                               `tfsdk:"url_redirect_walled_garden_enabled"`
	UrlRedirectWalledGardenRanges  types.List                                                   `tfsdk:"url_redirect_walled_garden_ranges"`
	VoiceVlanClients               jsontypes.Bool                                               `tfsdk:"voice_vlan_clients"`
}

func (d *NetworksSwitchAccessPoliciesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_switch_access_policies"
}

func (d *NetworksSwitchAccessPoliciesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List the dashboard administrators in this organization",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},

			"network_id": schema.StringAttribute{
				MarkdownDescription: "Organization ID",
				Optional:            true,
				CustomType:          jsontypes.StringType,
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 31),
				},
			},

			"list": schema.ListNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"access_policy_type": schema.StringAttribute{
							MarkdownDescription: "Access Type of the policy. Automatically 'Hybrid authentication' when hostMode is 'Multi-Domain'.",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},

						"dot_1x_control_direction": schema.StringAttribute{
							MarkdownDescription: "Supports either 'both' or 'inbound'. Set to 'inbound' to allow unauthorized egress on the switchport. Set to 'both' to control both traffic directions with authorization. Defaults to 'both'",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},

						"guest_vlan_id": schema.Int64Attribute{
							MarkdownDescription: "ID for the guest VLAN allow unauthorized devices access to limited network resources",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.Int64Type,
						},

						"host_mode": schema.StringAttribute{
							MarkdownDescription: "Choose the Host Mode for the access policy.",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},

						"increase_access_speed": schema.BoolAttribute{
							MarkdownDescription: "Enabling this option will make switches execute 802.1X and MAC-bypass authentication simultaneously so that clients authenticate faster. Only required when accessPolicyType is 'Hybrid Authentication.",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.BoolType,
						},

						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the access policy",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},

						"radius_accounting_enabled": schema.BoolAttribute{
							MarkdownDescription: "Enable to send start, interim-update and stop messages to a configured RADIUS accounting server for tracking connected clients",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.BoolType,
						},

						"radius_accounting_servers": schema.SetNestedAttribute{
							MarkdownDescription: "List of RADIUS accounting servers to require connecting devices to authenticate against before granting network access",
							Optional:            true,
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"host": schema.StringAttribute{
										MarkdownDescription: "",
										Optional:            true,
										CustomType:          jsontypes.StringType,
									},
									"port": schema.Int64Attribute{
										MarkdownDescription: "",
										Optional:            true,
										CustomType:          jsontypes.Int64Type,
									},
								},
							},
						},

						"radius_coa_support_enabled": schema.BoolAttribute{
							MarkdownDescription: "Change of authentication for RADIUS re-authentication and disconnection",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.BoolType,
						},

						"radius_group_attribute": schema.StringAttribute{
							MarkdownDescription: `Acceptable values are "" for None, or "11" for Group Policies ACL`,
							Optional:            true,
							Computed:            true,
						},

						"radius_servers": schema.SetNestedAttribute{
							MarkdownDescription: "List of RADIUS servers to require connecting devices to authenticate against before granting network access",
							Optional:            true,
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"host": schema.StringAttribute{
										MarkdownDescription: "",
										Optional:            true,
										CustomType:          jsontypes.StringType,
									},
									"port": schema.Int64Attribute{
										MarkdownDescription: "",
										Optional:            true,
										CustomType:          jsontypes.Int64Type,
									},
								},
							},
						},

						"radius_testing_enabled": schema.BoolAttribute{
							MarkdownDescription: "If enabled, Meraki devices will periodically send access-request messages to these RADIUS servers",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.BoolType,
						},

						"radius_critical_auth": schema.SetNestedAttribute{
							MarkdownDescription: "Critical auth settings for when authentication is rejected by the RADIUS server",
							Optional:            true,
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{},
							},
						},

						"radius_failed_auth_vlan_id": schema.SetNestedAttribute{
							MarkdownDescription: "VLAN that clients will be placed on when RADIUS authentication fails. Will be null if hostMode is Multi-Auth",
							Optional:            true,
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{},
							},
						},

						"radius_re_authentication_interval": schema.SetNestedAttribute{
							MarkdownDescription: "Re-authentication period in seconds. Will be null if hostMode is Multi-Auth",
							Optional:            true,
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{},
							},
						},

						"url_redirect_walled_garden_enabled": schema.BoolAttribute{
							MarkdownDescription: "Enable to restrict access for clients to a response_objectific set of IP addresses or hostnames prior to authentication",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.BoolType,
						},

						"url_redirect_walled_garden_ranges": schema.SetNestedAttribute{
							MarkdownDescription: "IP address ranges, in CIDR notation, to restrict access for clients to a specific set of IP addresses or hostnames prior to authentication",
							Optional:            true,
							Computed:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{},
							},
						},

						"voice_vlan_clients": schema.BoolAttribute{
							MarkdownDescription: "CDP/LLDP capable voice clients will be able to use this VLAN. Automatically true when hostMode is 'Multi-Domain'.",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.BoolType,
						},
					},
				},
			},
		},
	}
}

func (d *NetworksSwitchAccessPoliciesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NetworksSwitchAccessPoliciesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data NetworksSwitchesAccessPoliciesDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := d.client.SwitchApi.GetNetworkSwitchAccessPolicies(ctx, data.NetworkId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read datasource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// collect diagnostics
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("State Data", fmt.Sprintf("\n%s", data))
		return
	}

	// Save data into Terraform state
	data.Id = types.StringValue("example-id")
	if err = json.NewDecoder(httpResp.Body).Decode(&data.List); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read data source")
}
