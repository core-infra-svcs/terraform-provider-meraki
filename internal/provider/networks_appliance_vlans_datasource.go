package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
	"io"
	"log"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ datasource.DataSource = &NetworksApplianceVLANsDataSource{}
)

func NewNetworksApplianceVLANDataSource() datasource.DataSource {
	return &NetworksApplianceVLANsDataSource{}
}

// NetworksApplianceVLANsDataSource defines the resource implementation.
type NetworksApplianceVLANsDataSource struct {
	client *openApiClient.APIClient
}

type NetworksApplianceVLANsDataSourceModel struct {
	Id        jsontypes.String                       `tfsdk:"id" json:"-"`
	NetworkId jsontypes.String                       `tfsdk:"network_id" json:"networkId"`
	List      []NetworksApplianceVLANDataSourceModel `tfsdk:"list"`
}

type NetworksApplianceVLANDataSourceModel struct {
	NetworkId              jsontypes.String                                                                       `tfsdk:"network_id" json:"networkId"`
	VlanId                 jsontypes.Int64                                                                        `tfsdk:"vlan_id" json:"id"`
	InterfaceId            jsontypes.String                                                                       `tfsdk:"interface_id" json:"interfaceId,omitempty"`
	Name                   jsontypes.String                                                                       `tfsdk:"name" json:"name"`
	Subnet                 jsontypes.String                                                                       `tfsdk:"subnet" json:"subnet"`
	ApplianceIp            jsontypes.String                                                                       `tfsdk:"appliance_ip" json:"applianceIp"`
	GroupPolicyId          jsontypes.String                                                                       `tfsdk:"group_policy_id" json:"groupPolicyId"`
	TemplateVlanType       jsontypes.String                                                                       `tfsdk:"template_vlan_type" json:"templateVlanType"`
	Cidr                   jsontypes.String                                                                       `tfsdk:"cidr" json:"cidr"`
	Mask                   jsontypes.Int64                                                                        `tfsdk:"mask" json:"mask"`
	DhcpRelayServerIps     []jsontypes.String                                                                     `tfsdk:"dhcp_relay_server_ips" json:"dhcpRelayServerIps"`
	DhcpHandling           jsontypes.String                                                                       `tfsdk:"dhcp_handling" json:"dhcpHandling"`
	DhcpLeaseTime          jsontypes.String                                                                       `tfsdk:"dhcp_lease_time" json:"dhcpLeaseTime"`
	DhcpBootOptionsEnabled jsontypes.Bool                                                                         `tfsdk:"dhcp_boot_options_enabled" json:"dhcpBootOptionsEnabled"`
	DhcpBootNextServer     jsontypes.String                                                                       `tfsdk:"dhcp_boot_next_server" json:"dhcpBootNextServer"`
	DhcpBootFilename       jsontypes.String                                                                       `tfsdk:"dhcp_boot_filename" json:"dhcpBootFilename"`
	FixedIpAssignments     jsontypes.Map[jsontypes.Object[NetworksApplianceVLANDataSourceModelFixedIpAssignment]] `tfsdk:"fixed_ip_assignments" json:"fixedIpAssignments"`
	ReservedIpRanges       []NetworksApplianceVLANDataSourceModelReservedIpRange                                  `tfsdk:"reserved_ip_ranges" json:"reservedIpRanges"`
	DnsNameservers         jsontypes.String                                                                       `tfsdk:"dns_nameservers" json:"dnsNameservers"`
	DhcpOptions            []NetworksApplianceVLANDataSourceModelDhcpOption                                       `tfsdk:"dhcp_options" json:"dhcpOptions"`
	VpnNatSubnet           jsontypes.String                                                                       `tfsdk:"vpn_nat_subnet" json:"vpnNatSubnet"`
	MandatoryDhcp          NetworksApplianceVLANDataSourceModelMandatoryDhcp                                      `tfsdk:"mandatory_dhcp" json:"MandatoryDhcp"`
	IPv6                   NetworksApplianceVLANDataSourceModelIpv6                                               `tfsdk:"ipv6" json:"ipv6"`
}

type NetworksApplianceVLANDataSourceModelIpNameMapping struct {
	Ip   jsontypes.String `tfsdk:"ip" json:"ip"`
	Name jsontypes.String `tfsdk:"name" json:"name"`
}

type NetworksApplianceVLANDataSourceModelReservedIpRange struct {
	Start   jsontypes.String `tfsdk:"start" json:"start"`
	End     jsontypes.String `tfsdk:"end" json:"end"`
	Comment jsontypes.String `tfsdk:"comment" json:"comment"`
}

type NetworksApplianceVLANDataSourceModelDhcpOption struct {
	Code  jsontypes.String `tfsdk:"code" json:"code"`
	Type  jsontypes.String `tfsdk:"type" json:"type"`
	Value jsontypes.String `tfsdk:"value" json:"value"`
}

type NetworksApplianceVLANDataSourceModelFixedIpAssignment struct {
	jsontypes.BaseJsonValue
	Ip   jsontypes.String `tfsdk:"ip"  json:"ip"`
	Name jsontypes.String `tfsdk:"name"  json:"name"`
}

// NetworksApplianceVLANDataSourceModelIpv6 represents the IPv6 configuration for a VLAN resource model.
type NetworksApplianceVLANDataSourceModelIpv6 struct {
	Enabled           jsontypes.Bool                                             `tfsdk:"enabled" json:"enabled"`
	PrefixAssignments []NetworksApplianceVLANDataSourceModelIpv6PrefixAssignment `tfsdk:"prefix_assignments" json:"prefixAssignments"`
}

// NetworksApplianceVLANDataSourceModelIpv6PrefixAssignment represents a prefix assignment for an IPv6 configuration in the VLAN resource model.
type NetworksApplianceVLANDataSourceModelIpv6PrefixAssignment struct {
	Autonomous         jsontypes.Bool                                                 `tfsdk:"autonomous" json:"autonomous"`
	StaticPrefix       jsontypes.String                                               `tfsdk:"static_prefix" json:"staticPrefix"`
	StaticApplianceIp6 jsontypes.String                                               `tfsdk:"static_appliance_ip6" json:"staticApplianceIp6"`
	Origin             NetworksApplianceVLANDataSourceModelIpv6PrefixAssignmentOrigin `tfsdk:"origin" json:"origin"`
}

// NetworksApplianceVLANDataSourceModelIpv6PrefixAssignmentOrigin represents the origin data structure for a VLAN resource model.
type NetworksApplianceVLANDataSourceModelIpv6PrefixAssignmentOrigin struct {
	Type       jsontypes.String   `tfsdk:"type" json:"type"`
	Interfaces []jsontypes.String `tfsdk:"interfaces" json:"interfaces"`
}

type NetworksApplianceVLANDataSourceModelMandatoryDhcp struct {
	Enabled jsontypes.Bool `tfsdk:"enabled" json:"enabled"`
}

func (r *NetworksApplianceVLANsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_appliance_vlans"
}

func (r *NetworksApplianceVLANsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{

		MarkdownDescription: "Manage the VLANs for an MX network",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.StringType,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "The VLAN ID of the new VLAN (must be between 1 and 4094)",
				Required:            true,
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
						"vlan_id": schema.Int64Attribute{
							Optional:   true,
							Computed:   true,
							CustomType: jsontypes.Int64Type,
						},
						"network_id": schema.StringAttribute{
							MarkdownDescription: "The VLAN ID of the new VLAN (must be between 1 and 4094)",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
							Validators: []validator.String{
								stringvalidator.LengthBetween(8, 31),
							},
						},
						"interface_id": schema.StringAttribute{
							MarkdownDescription: "The Interface ID",
							Optional:            true,
							Computed:            true,
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
						"vpn_nat_subnet": schema.StringAttribute{
							MarkdownDescription: "The translated VPN subnet if VPN and VPN subnet translation are enabled on the VLAN",
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
						"dhcp_relay_server_ips": schema.ListAttribute{
							ElementType: jsontypes.StringType,
							Description: "An array of DHCP relay server IPs to which DHCP packets would get relayed for this VLAN",
							Optional:    true,
							Computed:    true,
						},
						"dhcp_lease_time": schema.StringAttribute{
							MarkdownDescription: "The term of DHCP leases if the appliance is running a DHCP server on this VLAN. One of: '30 minutes', '1 hour', '4 hours', '12 hours', '1 day' or '1 week'",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"dhcp_boot_options_enabled": schema.BoolAttribute{
							MarkdownDescription: "Use DHCP boot options specified in other properties",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.BoolType,
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
						"fixed_ip_assignments": schema.MapNestedAttribute{
							Optional:    true,
							Computed:    true,
							Description: "The DHCP fixed IP assignments on the VLAN, mapped by MAC address to an object containing 'ip' and 'name'.",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"ip": schema.StringAttribute{
										Description: "The IP address associated with the fixed IP assignment.",
										Optional:    true,
										Computed:    true,
										CustomType:  jsontypes.StringType,
									},
									"name": schema.StringAttribute{
										Description: "A descriptive name for the IP assignment.",
										Optional:    true,
										Computed:    true,
										CustomType:  jsontypes.StringType,
									},
								},
							},
						},
						"reserved_ip_ranges": schema.ListNestedAttribute{
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
						"dns_nameservers": schema.StringAttribute{
							MarkdownDescription: "The DNS nameservers used for DHCP responses, either \"upstream_dns\", \"google_dns\", \"opendns\", or a newline seperated string of IP addresses or domain names",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"dhcp_options": schema.ListNestedAttribute{
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
										Validators: []validator.String{
											stringvalidator.OneOf("text", "ip", "hex", "integer"),
										},
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
						"template_vlan_type": schema.StringAttribute{
							MarkdownDescription: "Type of subnetting of the VLAN. Applicable only for template network.",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
							Validators: []validator.String{
								stringvalidator.OneOf("same", "unique"),
							},
						},
						"cidr": schema.StringAttribute{
							MarkdownDescription: "CIDR of the pool of subnets. Applicable only for template network. Each network bound to the template will automatically pick a subnet from this pool to build its own VLAN.",
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
								"prefix_assignments": schema.ListNestedAttribute{
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
														Validators: []validator.String{
															stringvalidator.OneOf("independent", "internet"),
														},
													},
													"interfaces": schema.SetAttribute{
														ElementType: jsontypes.StringType,
														Description: "Interfaces associated with the prefix",
														Optional:    true,
														Computed:    true,
													},
												},
											},
										}},
								},
							},
						},
						"mandatory_dhcp": schema.SingleNestedAttribute{
							Description: "Mandatory DHCP will enforce that clients connecting to this VLAN must use the IP address assigned by the DHCP server. Clients who use a static IP address won't be able to associate. Only available on firmware versions 17.0 and above",
							Optional:    true,
							Computed:    true,
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
				},
			},
		},
	}
}

func (r *NetworksApplianceVLANsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *NetworksApplianceVLANsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *NetworksApplianceVLANsDataSourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := r.client.ApplianceApi.GetNetworkApplianceVlans(ctx, data.NetworkId.ValueString()).Execute()
	if err != nil && httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"HTTP Client Read Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
	}

	// Check for API success inlineResp code
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	}

	// Read the body into a byte slice
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return
	}
	if err != nil {
		log.Fatalf("Error reading response body: %v", err)
	}

	// Unmarshal the JSON data into the struct
	var list []NetworksApplianceVLANDataSourceModel
	if err := json.Unmarshal(body, &list); err != nil {
		resp.Diagnostics.Append(diag.NewErrorDiagnostic("unmarshal error", fmt.Sprintf("%s", err)))
	}

	data.List = list
	data.Id = jsontypes.StringValue("example-id")

	tflog.Info(ctx, "Payload", map[string]interface{}{
		"DATA": httpResp.Body,
	})

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}
