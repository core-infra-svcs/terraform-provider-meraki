package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// NetworksApplianceVlansDataSource struct. If not, implement them.
var _ datasource.DataSource = &NetworksApplianceVlansDataSource{}

// The NewNetworksApplianceVlansDataSource function is a constructor for the data source. This function needs
// to be added to the list of Data Sources in provider.go: func (p *ScaffoldingProvider) DataSources.
// If it's not added, the provider won't be aware of this data source's existence.
func NewNetworksApplianceVlansDataSource() datasource.DataSource {
	return &NetworksApplianceVlansDataSource{}
}

// NetworksApplianceVlansDataSource struct defines the structure for this data source.
// It includes an APIClient field for making requests to the Meraki API.
// If additional fields are required (e.g., for caching or for tracking internal state), add them here.
type NetworksApplianceVlansDataSource struct {
	client *openApiClient.APIClient
}

// The NetworksApplianceVlansDataSourceModel structure describes the data model.
// This struct is where you define all the attributes that are part of this data source's state.
type NetworksApplianceVlansDataSourceModel struct {

	// The Id field is mandatory for all data sources. It's used for data source identification and is required
	// for the acceptance tests to run.
	Id        jsontypes.String                   `tfsdk:"id" json:"-"`
	NetworkId jsontypes.String                   `tfsdk:"network_id" json:"networkId"`
	List      []NetworksApplianceVLANsDataSource `tfsdk:"list" json:"-"`
}

type NetworksApplianceVLANsDataSource struct {
	VlanId                 jsontypes.String                `tfsdk:"vlan_id" json:"id"`
	Name                   jsontypes.String                `tfsdk:"name" json:"name"`
	Subnet                 jsontypes.String                `tfsdk:"subnet" json:"subnet"`
	ApplianceIp            jsontypes.String                `tfsdk:"appliance_ip" json:"applianceIp"`
	GroupPolicyId          jsontypes.String                `tfsdk:"group_policy_id" json:"groupPolicyId"`
	VpnNatSubnet           jsontypes.String                `tfsdk:"vpn_nat_subnet" json:"vpnNatSubnet"`
	DhcpHandling           jsontypes.String                `tfsdk:"dhcp_handling" json:"dhcpHandling"`
	DhcpRelayServerIps     jsontypes.Set[jsontypes.String] `tfsdk:"dhcp_relay_server_ips" json:"dhcpRelayServerIps"`
	DhcpLeaseTime          jsontypes.String                `tfsdk:"dhcp_lease_time" json:"dhcpLeaseTime"`
	DhcpBootOptionsEnabled jsontypes.Bool                  `tfsdk:"dhcp_boot_options_enabled" json:"dhcpBootOptionsEnabled"`
	DhcpBootNextServer     jsontypes.String                `tfsdk:"dhcp_boot_next_server" json:"dhcpBootNextServer"`
	DhcpBootFilename       jsontypes.String                `tfsdk:"dhcp_boot_filename" json:"dhcpBootFilename"`
	FixedIpAssignments     FixedIpAssignments              `tfsdk:"fixed_ip_assignments" json:"fixedIpAssignments"`
	ReservedIpRanges       []ReservedIPRange               `tfsdk:"reserved_ip_ranges" json:"reservedIpRanges"`
	DnsNameservers         jsontypes.String                `tfsdk:"dns_nameservers" json:"dnsNameservers"`
	DhcpOptions            []DHCPOption                    `tfsdk:"dhcp_options" json:"dhcpOptions"`
	TemplateVlanType       jsontypes.String                `tfsdk:"template_vlan_type" json:"templateVlanType"`
	Cidr                   jsontypes.String                `tfsdk:"cidr" json:"cidr"`
	Mask                   jsontypes.Int64                 `tfsdk:"mask" json:"mask"`
	Ipv6                   IPV6                            `tfsdk:"ipv6" json:"ipv6"`
	MandatoryDHCP          MandatoryDHCP                   `tfsdk:"mandatory_dhcp" json:"mandatoryDhcp"`
}

type FixedIPAssignments struct {
	IP   jsontypes.String `tfsdk:"ip" json:"IP"`
	Name jsontypes.String `tfsdk:"name" json:"Name"`
}

// Metadata provides a way to define information about the data source.
// This method is called by the framework to retrieve metadata about the data source.
func (d *NetworksApplianceVlansDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_appliance_vlans"
}

// Schema provides a way to define the structure of the data source data.
// It is called by the framework to get the schema of the data source.
func (d *NetworksApplianceVlansDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	resp.Schema = schema.Schema{

		MarkdownDescription: "NetworksApplianceVlans",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.StringType,
			},
			"vlan_id": schema.StringAttribute{
				Computed:   false,
				Optional:   true,
				CustomType: jsontypes.StringType,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "The VLAN ID of the new VLAN (must be between 1 and 4094)",
				Required:            true,
				CustomType:          jsontypes.StringType,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the new VLAN",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"subnet": schema.StringAttribute{
				MarkdownDescription: "The subnet of the VLAN",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"appliance_ip": schema.StringAttribute{
				MarkdownDescription: "The local IP of the appliance on the VLAN",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"group_policy_id": schema.StringAttribute{
				MarkdownDescription: " desired group policy to apply to the VLAN",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"template_vlan_type": schema.StringAttribute{
				MarkdownDescription: "Type of subnetting of the VLAN. Applicable only for template network.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"cidr": schema.StringAttribute{
				MarkdownDescription: "CIDR of the pool of subnets. Applicable only for template network. Each network bound to the template will automatically pick a subnet from this pool to build its own VLAN.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"dhcp_handling": schema.StringAttribute{
				MarkdownDescription: "The appliance's handling of DHCP requests on this VLAN. One of: 'Run a DHCP server', 'Relay DHCP to another server' or 'Do not respond to DHCP requests'",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"dhcp_lease_time": schema.StringAttribute{
				MarkdownDescription: "The term of DHCP leases if the appliance is running a DHCP server on this VLAN. One of: '30 minutes', '1 hour', '4 hours', '12 hours', '1 day' or '1 week'",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"dhcp_boot_next_server": schema.StringAttribute{
				MarkdownDescription: "DHCP boot option to direct boot clients to the server to load the boot file from",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"dhcp_boot_filename": schema.StringAttribute{
				MarkdownDescription: "DHCP boot option for boot filename ",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"dns_nameservers": schema.StringAttribute{
				MarkdownDescription: "The DNS nameservers used for DHCP responses, either \"upstream_dns\", \"google_dns\", \"opendns\", or a newline seperated string of IP addresses or domain names",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"vpn_nat_subnet": schema.StringAttribute{
				MarkdownDescription: "The translated VPN subnet if VPN and VPN subnet translation are enabled on the VLAN",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"mask": schema.Int64Attribute{
				MarkdownDescription: "Mask used for the subnet of all bound to the template networks. Applicable only for template network.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.Int64Type,
			},
			"dhcp_boot_options_enabled": schema.BoolAttribute{
				MarkdownDescription: "Use DHCP boot options specified in other properties",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.BoolType,
			},
			"dhcp_relay_server_ips": schema.SetAttribute{
				CustomType:  jsontypes.SetType[jsontypes.String](),
				ElementType: jsontypes.StringType,
				Description: "The IPs of the DHCP servers that DHCP requests should be relayed to",
				Computed:    true,
				Optional:    true,
			},
			"reserved_ip_ranges": schema.SetNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The DHCP reserved IP ranges on the VLAN",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"start": schema.StringAttribute{
							MarkdownDescription: "The first IP in the reserved range",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"end": schema.StringAttribute{
							MarkdownDescription: "The last IP in the reserved range",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"comment": schema.StringAttribute{
							MarkdownDescription: "A text comment for the reserved range",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
					},
				},
			},
			"dhcp_options": schema.SetNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "The list of DHCP options that will be included in DHCP responses. Each object in the list should have \"code\", \"type\", and \"value\" properties.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"code": schema.StringAttribute{
							MarkdownDescription: "The code for the DHCP option. This should be an integer between 2 and 254.",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "The type for the DHCP option. One of: 'text', 'ip', 'hex' or 'integer'",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"value": schema.StringAttribute{
							MarkdownDescription: "The value for the DHCP option",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
					},
				},
			},
			"ipv6": schema.SingleNestedAttribute{
				Description: "IPv6 configuration on the VLAN",
				Optional:    true,
				Computed:    true,
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "Enable IPv6 on VLAN.",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.BoolType,
					},
					"prefix_assignments": schema.SetNestedAttribute{
						Optional:    true,
						Computed:    true,
						Description: "Prefix assignments on the VLAN",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"autonomous": schema.BoolAttribute{
									MarkdownDescription: "Auto assign a /64 prefix from the origin to the VLAN",
									Optional:            true,
									Computed:            true,
									CustomType:          jsontypes.BoolType,
								},
								"static_prefix": schema.StringAttribute{
									MarkdownDescription: "Manual configuration of a /64 prefix on the VLAN",
									Optional:            true,
									Computed:            true,
									CustomType:          jsontypes.StringType,
								},
								"static_appliance_ip6": schema.StringAttribute{
									MarkdownDescription: "Manual configuration of the IPv6 Appliance IP",
									Optional:            true,
									Computed:            true,
									CustomType:          jsontypes.StringType,
								},
								"origin": schema.SingleNestedAttribute{
									MarkdownDescription: "The origin of the prefix",
									Optional:            true,
									Computed:            true,
									Attributes: map[string]schema.Attribute{
										"type": schema.StringAttribute{
											MarkdownDescription: "Type of the origin",
											Optional:            true,
											Computed:            true,
											CustomType:          jsontypes.StringType,
										},
										"interfaces": schema.SetAttribute{
											CustomType:  jsontypes.SetType[jsontypes.String](),
											ElementType: jsontypes.StringType,
											Description: "Interfaces associated with the prefix",
											Computed:    true,
											Optional:    true,
										},
									},
								},
							},
						},
					},
				},
			},
			"fixed_ip_assignments": schema.SingleNestedAttribute{
				Description: "The DHCP fixed IP assignments on the VLAN. This should be an object that contains mappings from MAC addresses to objects that themselves each contain \"ip\" and \"name\" string fields. See the sample request/response for more details",
				Optional:    true,
				Computed:    false,
				Attributes: map[string]schema.Attribute{
					"ip": schema.StringAttribute{
						MarkdownDescription: "Enable IPv6 on VLAN.",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.StringType,
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "Enable IPv6 on VLAN.",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.StringType,
					},
				},
			},
			"mandatory_dhcp": schema.SingleNestedAttribute{
				Description: "Mandatory DHCP will enforce that clients connecting to this VLAN must use the IP address assigned by the DHCP server. Clients who use a static IP address won't be able to associate. Only available on firmware versions 17.0 and above",
				Optional:    true,
				Computed:    false,
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "Enable Mandatory DHCP on VLAN.",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.BoolType,
					},
				},
			},
		},
	}
}

// Configure is a method of the data source interface that Terraform calls to provide the configured provider instance to the data source.
// It passes the DataSourceData that's been stored by the provider's ConfigureFunc.
func (d *NetworksApplianceVlansDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

	// The provider must be properly configured before it can be used.
	if req.ProviderData == nil {
		return
	}

	// Here we expect the provider data to be of type *openApiClient.APIClient.
	client, ok := req.ProviderData.(*openApiClient.APIClient)

	// This is a fatal error and the provider cannot proceed without it.
	// If you see this error, it means there is an issue with the provider setup.
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	// This allows the data source to use the configured provider for any API calls it needs to make.
	d.client = client
}

// Read method is responsible for reading an existing data source's state.
func (d *NetworksApplianceVlansDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data NetworksApplianceVlansDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	// Check if there are any errors before proceeding.
	if resp.Diagnostics.HasError() {
		return
	}

	// Remember to handle any potential errors.
	_, httpResp, err := d.client.ApplianceApi.GetNetworkApplianceVlans(ctx, data.NetworkId.ValueString()).Execute()

	// If there was an error during API call, add it to diagnostics.
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read data source",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// Collect any HTTP diagnostics that might be useful for debugging.
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
	}

	// If it's not what you expect, add an error to diagnostics.
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// If there were any errors up to this point, log the state data and return.
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("State Data", fmt.Sprintf("\n%v", data))
		return
	}

	// Decode the HTTP response body into your data model.
	// If there's an error, add it to diagnostics.
	if err = json.NewDecoder(httpResp.Body).Decode(&data.List); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	// Set ID for the data source.
	data.Id = jsontypes.StringValue("example-id")

	// Now set the final state of the data source.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read data source")
}
