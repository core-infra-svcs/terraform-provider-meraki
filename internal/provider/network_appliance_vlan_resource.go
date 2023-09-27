package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
	"strings"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &NetworksApplianceVLANsResource{}
	_ resource.ResourceWithConfigure   = &NetworksApplianceVLANsResource{}
	_ resource.ResourceWithImportState = &NetworksApplianceVLANsResource{}
)

func NewNetworksApplianceVLANsResource() resource.Resource {
	return &NetworksApplianceVLANsResource{}
}

// NetworksApplianceVLANsResource defines the resource implementation.
type NetworksApplianceVLANsResource struct {
	client *openApiClient.APIClient
}

type NetworksApplianceVLANsResourceModel struct {
	Id        jsontypes.String `tfsdk:"id" json:"-"`
	NetworkId jsontypes.String `tfsdk:"network_id" json:"networkId"`

	VlanId                 jsontypes.Int64                 `tfsdk:"vlan_id" json:"id"`
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
	FixedIpAssignments     ipNameMapping                   `tfsdk:"fixed_ip_assignments" json:"fixedIpAssignments"`
	ReservedIpRanges       []reservedIpRange               `tfsdk:"reserved_ip_ranges" json:"reservedIpRanges"`
	DnsNameservers         jsontypes.String                `tfsdk:"dns_nameservers" json:"dnsNameservers"`
	DhcpOptions            []dhcpOption                    `tfsdk:"dhcp_options" json:"dhcpOptions"`
	TemplateVlanType       jsontypes.String                `tfsdk:"template_vlan_type" json:"templateVlanType"`
	Cidr                   jsontypes.String                `tfsdk:"cidr" json:"cidr"`
	Mask                   jsontypes.Int64                 `tfsdk:"mask" json:"mask"`
	IPv6                   ipv6Configuration               `tfsdk:"ipv6" json:"ipv6"`
	MandatoryDhcp          mandatoryDhcp                   `tfsdk:"mandatory_dhcp" json:"mandatoryDhcp"`
}

type ipNameMapping struct {
	Ip   jsontypes.String `tfsdk:"ip" json:"ip"`
	Name jsontypes.String `tfsdk:"name" json:"name"`
}

type reservedIpRange struct {
	Start   jsontypes.String `tfsdk:"start" json:"start"`
	End     jsontypes.String `tfsdk:"end" json:"end"`
	Comment jsontypes.String `tfsdk:"comment" json:"comment"`
}

type dhcpOption struct {
	Code  jsontypes.String `tfsdk:"code" json:"code"`
	Type  jsontypes.String `tfsdk:"type" json:"type"`
	Value jsontypes.String `tfsdk:"value" json:"value"`
}

type ipv6Configuration struct {
	Enabled           jsontypes.Bool     `tfsdk:"enabled" json:"enabled"`
	PrefixAssignments []prefixAssignment `tfsdk:"prefix_assignments" json:"prefixAssignments"`
}

type prefixAssignment struct {
	Autonomous         jsontypes.Bool   `tfsdk:"autonomous" json:"autonomous"`
	StaticPrefix       jsontypes.String `tfsdk:"static_prefix" json:"staticPrefix"`
	StaticApplianceIp6 jsontypes.String `tfsdk:"static_appliance_ip6" json:"staticApplianceIp6"`
	Origin             origin           `tfsdk:"origin" json:"origin"`
}

type origin struct {
	Type       jsontypes.String   `tfsdk:"type" json:"type"`
	Interfaces []jsontypes.String `tfsdk:"interfaces" json:"interfaces"`
}

type mandatoryDhcp struct {
	Enabled jsontypes.Bool `tfsdk:"enabled" json:"enabled"`
}

func (r *NetworksApplianceVLANsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_appliance_vlans"
}

func (r *NetworksApplianceVLANsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{

		MarkdownDescription: "NetworksApplianceVlans",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.StringType,
			},
			"vlan_id": schema.Int64Attribute{
				Optional:   true,
				Computed:   true,
				CustomType: jsontypes.Int64Type,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "The VLAN ID of the new VLAN (must be between 1 and 4094)",
				Required:            true,
				CustomType:          jsontypes.StringType,
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 31),
				},
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
				CustomType:          jsontypes.StringType,
			},
			"vpn_nat_subnet": schema.StringAttribute{
				MarkdownDescription: "The translated VPN subnet if VPN and VPN subnet translation are enabled on the VLAN",
				Optional:            true,
				CustomType:          jsontypes.StringType,
			},
			"dhcp_handling": schema.StringAttribute{
				MarkdownDescription: "The appliance's handling of DHCP requests on this VLAN. One of: 'Run a DHCP server', 'Relay DHCP to another server' or 'Do not respond to DHCP requests'",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"dhcp_relay_server_ips": schema.SetAttribute{
				CustomType:  jsontypes.SetType[jsontypes.String](),
				ElementType: jsontypes.StringType,
				Description: "An array of DHCP relay server IPs to which DHCP packets would get relayed for this VLAN",
				Optional:    true,
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
				CustomType:          jsontypes.StringType,
			},
			"dhcp_boot_filename": schema.StringAttribute{
				MarkdownDescription: "DHCP boot option for boot filename ",
				Optional:            true,
				CustomType:          jsontypes.StringType,
			},
			"fixed_ip_assignments": schema.SingleNestedAttribute{
				Description: "The DHCP fixed IP assignments on the VLAN. This should be an object that contains mappings from MAC addresses to objects that themselves each contain \"ip\" and \"name\" string fields. See the sample request/response for more details",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"ip": schema.StringAttribute{
						MarkdownDescription: "Enable IPv6 on VLAN.",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "Enable IPv6 on VLAN.",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
				},
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
			"dns_nameservers": schema.StringAttribute{
				MarkdownDescription: "The DNS nameservers used for DHCP responses, either \"upstream_dns\", \"google_dns\", \"opendns\", or a newline seperated string of IP addresses or domain names",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
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
			"mask": schema.Int64Attribute{
				MarkdownDescription: "Mask used for the subnet of all bound to the template networks. Applicable only for template network.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.Int64Type,
			},
			"ipv6": schema.SingleNestedAttribute{
				Description: "IPv6 configuration on the VLAN",
				Optional:    true,
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "Enable IPv6 on VLAN.",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.BoolType,
					},
					"prefix_assignments": schema.SetNestedAttribute{
						Optional:    true,
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
											Optional:    true,
										},
									},
								},
							},
						},
					},
				},
			},
			"mandatory_dhcp": schema.SingleNestedAttribute{
				Description: "Mandatory DHCP will enforce that clients connecting to this VLAN must use the IP address assigned by the DHCP server. Clients who use a static IP address won't be able to associate. Only available on firmware versions 17.0 and above",
				Optional:    true,
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

func (r *NetworksApplianceVLANsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NetworksApplianceVLANsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *NetworksApplianceVLANsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	object56 := openApiClient.NewCreateNetworkApplianceVlanRequest(data.Id.ValueString(), data.Name.ValueString())
	object56.SetCidr(data.Cidr.ValueString())
	object56.SetId(fmt.Sprintf("%v", data.VlanId.ValueInt64()))
	object56.SetApplianceIp(data.ApplianceIp.ValueString())
	object56.SetGroupPolicyId(data.GroupPolicyId.ValueString())
	object56.SetMask(int32(data.Mask.ValueInt64()))
	object56.SetName(data.Name.ValueString())
	object56.SetSubnet(data.Subnet.ValueString())
	object56.SetTemplateVlanType(data.TemplateVlanType.ValueString())
	ipv6 := openApiClient.NewUpdateNetworkApplianceSingleLanRequestIpv6()
	ipv6.SetEnabled(data.IPv6.Enabled.ValueBool())
	var prefixAssignments []openApiClient.UpdateNetworkApplianceSingleLanRequestIpv6PrefixAssignmentsInner
	for _, prefixAssignment := range data.IPv6.PrefixAssignments {
		var originInterfaces []string
		for _, originInterface := range prefixAssignment.Origin.Interfaces {
			originInterfaces = append(originInterfaces, originInterface.String())
		}
		prefixAssignments = append(prefixAssignments, openApiClient.UpdateNetworkApplianceSingleLanRequestIpv6PrefixAssignmentsInner{
			Autonomous:         prefixAssignment.Autonomous.ValueBoolPointer(),
			StaticPrefix:       prefixAssignment.StaticPrefix.ValueStringPointer(),
			StaticApplianceIp6: prefixAssignment.StaticApplianceIp6.ValueStringPointer(),
			Origin: &openApiClient.UpdateNetworkApplianceSingleLanRequestIpv6PrefixAssignmentsInnerOrigin{
				Type:       prefixAssignment.Origin.Type.ValueString(),
				Interfaces: originInterfaces,
			},
		})
	}
	ipv6.SetPrefixAssignments(prefixAssignments)
	object56.SetIpv6(*ipv6)
	dhcp := openApiClient.NewGetNetworkApplianceVlans200ResponseInnerMandatoryDhcp()
	dhcp.SetEnabled(data.MandatoryDhcp.Enabled.ValueBool())
	object56.SetMandatoryDhcp(*dhcp)

	_, httpResp, err := r.client.ApplianceApi.CreateNetworkApplianceVlan(ctx, data.NetworkId.ValueString()).CreateNetworkApplianceVlanRequest(*object56).Execute()

	// Meraki API seems to return http status code 201 as an error.
	if err != nil && httpResp.StatusCode != 201 {
		resp.Diagnostics.AddError(
			"HTTP Client Create Failure",
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
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
		return
	}

	// save inlineResp data into Terraform state.
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"Create JSON Decode issue",
			fmt.Sprintf("%v, %v", err, httpResp.Body),
		)
		return
	}

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")
}

func (r *NetworksApplianceVLANsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetworksApplianceVLANsResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := r.client.ApplianceApi.GetNetworkApplianceVlan(ctx, data.NetworkId.ValueString(), fmt.Sprintf("%v", data.VlanId.ValueInt64())).Execute()
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

	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}
	data.Id = jsontypes.StringValue("example-id")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *NetworksApplianceVLANsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *NetworksApplianceVLANsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworkApplianceVlan := openApiClient.NewUpdateNetworkApplianceVlanRequest()
	updateNetworkApplianceVlan.SetCidr(data.Cidr.ValueString())
	updateNetworkApplianceVlan.SetApplianceIp(data.ApplianceIp.ValueString())
	updateNetworkApplianceVlan.SetGroupPolicyId(data.GroupPolicyId.ValueString())
	updateNetworkApplianceVlan.SetMask(int32(data.Mask.ValueInt64()))
	updateNetworkApplianceVlan.SetName(data.Name.ValueString())
	updateNetworkApplianceVlan.SetSubnet(data.Subnet.ValueString())
	updateNetworkApplianceVlan.SetTemplateVlanType(data.TemplateVlanType.ValueString())
	ipv6 := openApiClient.NewUpdateNetworkApplianceSingleLanRequestIpv6()
	ipv6.SetEnabled(data.IPv6.Enabled.ValueBool())
	var prefixAssignments []openApiClient.UpdateNetworkApplianceSingleLanRequestIpv6PrefixAssignmentsInner
	for _, prefixAssignment := range data.IPv6.PrefixAssignments {
		var originInterfaces []string
		for _, originInterface := range prefixAssignment.Origin.Interfaces {
			originInterfaces = append(originInterfaces, originInterface.String())
		}
		prefixAssignments = append(prefixAssignments, openApiClient.UpdateNetworkApplianceSingleLanRequestIpv6PrefixAssignmentsInner{
			Autonomous:         prefixAssignment.Autonomous.ValueBoolPointer(),
			StaticPrefix:       prefixAssignment.StaticPrefix.ValueStringPointer(),
			StaticApplianceIp6: prefixAssignment.StaticApplianceIp6.ValueStringPointer(),
			Origin: &openApiClient.UpdateNetworkApplianceSingleLanRequestIpv6PrefixAssignmentsInnerOrigin{
				Type:       prefixAssignment.Origin.Type.ValueString(),
				Interfaces: originInterfaces,
			},
		})
	}
	ipv6.SetPrefixAssignments(prefixAssignments)
	updateNetworkApplianceVlan.SetIpv6(*ipv6)
	dhcp := openApiClient.NewGetNetworkApplianceVlans200ResponseInnerMandatoryDhcp()
	dhcp.SetEnabled(data.MandatoryDhcp.Enabled.ValueBool())
	updateNetworkApplianceVlan.SetMandatoryDhcp(*dhcp)

	_, httpResp, err := r.client.ApplianceApi.UpdateNetworkApplianceVlan(ctx, data.NetworkId.ValueString(), fmt.Sprintf("%v", data.VlanId.ValueInt64())).UpdateNetworkApplianceVlanRequest(*updateNetworkApplianceVlan).Execute()
	if err != nil && httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"HTTP Client Update Failure",
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
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
		return
	}

	// save inlineResp data into Terraform state.
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON Decode issue",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
		return
	}
	data.Id = jsontypes.StringValue("example-id")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *NetworksApplianceVLANsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *NetworksApplianceVLANsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.client.ApplianceApi.DeleteNetworkApplianceVlan(ctx, data.NetworkId.ValueString(), fmt.Sprintf("%v", data.VlanId.ValueInt64())).Execute()
	if err != nil && httpResp.StatusCode != 204 {
		resp.Diagnostics.AddError(
			"HTTP Client Delete Failure",
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

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")

}

func (r *NetworksApplianceVLANsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: organization_id, admin_id. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vlan_id"), idParts[1])...)
	if resp.Diagnostics.HasError() {
		return
	}
}
