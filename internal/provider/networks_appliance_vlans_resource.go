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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &NetworksApplianceVlansResource{}
	_ resource.ResourceWithConfigure   = &NetworksApplianceVlansResource{}
	_ resource.ResourceWithImportState = &NetworksApplianceVlansResource{}
)

func NewNetworksApplianceVlansResource() resource.Resource {
	return &NetworksApplianceVlansResource{}
}

// NetworksApplianceVlansResource defines the resource implementation.
type NetworksApplianceVlansResource struct {
	client *openApiClient.APIClient
}

// NetworksApplianceVlansResourceModel describes the resource data model.
type NetworksApplianceVlansResourceModel struct {
	Id        jsontypes.String `tfsdk:"id" json:"-"`
	NetworkId jsontypes.String `tfsdk:"network_id" json:"networkId"`

	VlanId                 jsontypes.String                `tfsdk:"vlan_id" json:"id"`
	Name                   jsontypes.String                `tfsdk:"name" json:"name"`
	Subnet                 jsontypes.String                `tfsdk:"subnet" json:"subnet"`
	ApplianceIp            jsontypes.String                `tfsdk:"appliance_ip" json:"applianceIp"`
	GroupPolicyId          jsontypes.String                `tfsdk:"group_policy_id" json:"groupPolicyId"`
	TemplateVlanType       jsontypes.String                `tfsdk:"template_vlan_type" json:"templateVlanType"`
	Cidr                   jsontypes.String                `tfsdk:"cidr" json:"cidr"`
	DhcpHandling           jsontypes.String                `tfsdk:"dhcp_handling" json:"dhcpHandling"`
	DhcpLeaseTime          jsontypes.String                `tfsdk:"dhcp_lease_time" json:"dhcpLeaseTime"`
	DhcpBootNextServer     jsontypes.String                `tfsdk:"dhcp_boot_next_server" json:"dhcpBootNextServer"`
	DhcpBootFilename       jsontypes.String                `tfsdk:"dhcp_boot_filename" json:"dhcpBootFilename"`
	DnsNameservers         jsontypes.String                `tfsdk:"dns_nameservers" json:"dnsNameservers"`
	VpnNatSubnet           jsontypes.String                `tfsdk:"vpn_nat_subnet" json:"vpnNatSubnet"`
	Mask                   jsontypes.Int64                 `tfsdk:"mask" json:"mask"`
	DhcpBootOptionsEnabled jsontypes.Bool                  `tfsdk:"dhcp_boot_options_enabled" json:"dhcpBootOptionsEnabled"`
	DhcpRelayServerIps     jsontypes.Set[jsontypes.String] `tfsdk:"dhcp_relay_server_ips" json:"dhcpRelayServerIps"`
	ReservedIpRanges       []ReservedIPRange               `tfsdk:"reserved_ip_ranges" json:"reservedIpRanges"`
	DhcpOptions            []DHCPOption                    `tfsdk:"dhcp_options" json:"dhcpOptions"`
	Ipv6                   IPV6                            `tfsdk:"ipv6" json:"ipv6"`
}

type ReservedIPRange struct {
	Start   jsontypes.String `tfsdk:"start" json:"start"`
	End     jsontypes.String `tfsdk:"end" json:"end"`
	Comment jsontypes.String `tfsdk:"comment" json:"comment"`
}

type IPV6 struct {
	Enabled           jsontypes.Bool     `tfsdk:"enabled" json:"enabled"`
	PrefixAssignments []PrefixAssignment `tfsdk:"prefix_assignments" json:"prefixAssignments"`
}

type PrefixAssignment struct {
	Autonomous         jsontypes.Bool   `tfsdk:"autonomous" json:"autonomous"`
	StaticPrefix       jsontypes.String `tfsdk:"static_prefix" json:"staticPrefix"`
	StaticApplianceIp6 jsontypes.String `tfsdk:"static_appliance_ip6" json:"staticApplianceIp6"`
	Origin             Origin           `tfsdk:"origin" json:"origin"`
}

type Origin struct {
	Type       jsontypes.String                `tfsdk:"type" json:"type"`
	Interfaces jsontypes.Set[jsontypes.String] `tfsdk:"interfaces" json:"interfaces"`
}

type DHCPOption struct {
	Code  jsontypes.String `tfsdk:"code" json:"code"`
	Type  jsontypes.String `tfsdk:"type" json:"type"`
	Value jsontypes.String `tfsdk:"value" json:"value"`
}

func (r *NetworksApplianceVlansResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_appliance_vlans"
}

func (r *NetworksApplianceVlansResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
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
			"name": schema.StringAttribute{
				MarkdownDescription: "The destination port range of the incoming packet. Applicable only if protocol is set to TCP or UDP. Example: 70-80",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"subnet": schema.StringAttribute{
				MarkdownDescription: "The destination port range of the incoming packet. Applicable only if protocol is set to TCP or UDP. Example: 70-80",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"appliance_ip": schema.StringAttribute{
				MarkdownDescription: "The destination port range of the incoming packet. Applicable only if protocol is set to TCP or UDP. Example: 70-80",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"group_policy_id": schema.StringAttribute{
				MarkdownDescription: "The destination port range of the incoming packet. Applicable only if protocol is set to TCP or UDP. Example: 70-80",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"template_vlan_type": schema.StringAttribute{
				MarkdownDescription: "The destination port range of the incoming packet. Applicable only if protocol is set to TCP or UDP. Example: 70-80",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"cidr": schema.StringAttribute{
				MarkdownDescription: "The destination port range of the incoming packet. Applicable only if protocol is set to TCP or UDP. Example: 70-80",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"dhcp_handling": schema.StringAttribute{
				MarkdownDescription: "The destination port range of the incoming packet. Applicable only if protocol is set to TCP or UDP. Example: 70-80",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"dhcp_lease_time": schema.StringAttribute{
				MarkdownDescription: "The destination port range of the incoming packet. Applicable only if protocol is set to TCP or UDP. Example: 70-80",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"dhcp_boot_next_server": schema.StringAttribute{
				MarkdownDescription: "The destination port range of the incoming packet. Applicable only if protocol is set to TCP or UDP. Example: 70-80",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"dhcp_boot_filename": schema.StringAttribute{
				MarkdownDescription: "The destination port range of the incoming packet. Applicable only if protocol is set to TCP or UDP. Example: 70-80",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"dns_nameservers": schema.StringAttribute{
				MarkdownDescription: "The destination port range of the incoming packet. Applicable only if protocol is set to TCP or UDP. Example: 70-80",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"vpn_nat_subnet": schema.StringAttribute{
				MarkdownDescription: "The destination port range of the incoming packet. Applicable only if protocol is set to TCP or UDP. Example: 70-80",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"mask": schema.Int64Attribute{
				MarkdownDescription: "The destination port range of the incoming packet. Applicable only if protocol is set to TCP or UDP. Example: 70-80",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.Int64Type,
			},
			"dhcp_boot_options_enabled": schema.BoolAttribute{
				MarkdownDescription: "The destination port range of the incoming packet. Applicable only if protocol is set to TCP or UDP. Example: 70-80",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.BoolType,
			},
			"dhcp_relay_server_ips": schema.SetAttribute{
				CustomType:  jsontypes.SetType[jsontypes.String](),
				ElementType: jsontypes.StringType,
				Description: "Network tags",
				Computed:    true,
				Optional:    true,
			},
			"reserved_ip_ranges": schema.SetNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"start": schema.StringAttribute{
							MarkdownDescription: "Network ID",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"end": schema.StringAttribute{
							MarkdownDescription: "Network ID",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"comment": schema.StringAttribute{
							MarkdownDescription: "Network ID",
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
				Description: "",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"code": schema.StringAttribute{
							MarkdownDescription: "Network ID",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "Network ID",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"value": schema.StringAttribute{
							MarkdownDescription: "Network ID",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
					},
				},
			},
			"ipv6": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "Enables / disables the secure port.",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.BoolType,
					},
					"prefix_assignments": schema.SetNestedAttribute{
						Optional:    true,
						Computed:    true,
						Description: "",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"autonomous": schema.BoolAttribute{
									MarkdownDescription: "Network ID",
									Optional:            true,
									Computed:            true,
									CustomType:          jsontypes.BoolType,
								},
								"static_prefix": schema.StringAttribute{
									MarkdownDescription: "Network ID",
									Optional:            true,
									Computed:            true,
									CustomType:          jsontypes.StringType,
								},
								"static_appliance_ip6": schema.StringAttribute{
									MarkdownDescription: "Network ID",
									Optional:            true,
									Computed:            true,
									CustomType:          jsontypes.StringType,
								},
								"origin": schema.SingleNestedAttribute{
									Optional: true,
									Computed: true,
									Attributes: map[string]schema.Attribute{
										"type": schema.StringAttribute{
											MarkdownDescription: "Network ID",
											Optional:            true,
											Computed:            true,
											CustomType:          jsontypes.StringType,
										},
										"interfaces": schema.SetAttribute{
											CustomType:  jsontypes.SetType[jsontypes.String](),
											ElementType: jsontypes.StringType,
											Description: "Network tags",
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
		},
	}
}

func (r *NetworksApplianceVlansResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NetworksApplianceVlansResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *NetworksApplianceVlansResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	object56 := openApiClient.NewInlineObject56(data.Id.ValueString(), data.Name.ValueString())
	object56.SetCidr(data.Cidr.ValueString())
	object56.SetId(data.VlanId.ValueString())
	object56.SetApplianceIp(data.ApplianceIp.ValueString())
	object56.SetGroupPolicyId(data.GroupPolicyId.ValueString())
	object56.SetMask(int32(data.Mask.ValueInt64()))
	object56.SetName(data.Name.ValueString())
	object56.SetSubnet(data.Subnet.ValueString())
	object56.SetTemplateVlanType(data.TemplateVlanType.ValueString())
	ipv6 := openApiClient.NewNetworksNetworkIdApplianceSingleLanIpv6()
	ipv6.SetEnabled(data.Ipv6.Enabled.ValueBool())
	var prefixAssignments []openApiClient.NetworksNetworkIdApplianceSingleLanIpv6PrefixAssignments
	for _, prefixAssignment := range data.Ipv6.PrefixAssignments {
		originInterfaces := []string{}
		for _, originInterface := range prefixAssignment.Origin.Interfaces.Elements() {
			originInterfaces = append(originInterfaces, originInterface.String())
		}
		prefixAssignments = append(prefixAssignments, openApiClient.NetworksNetworkIdApplianceSingleLanIpv6PrefixAssignments{
			Autonomous:         prefixAssignment.Autonomous.ValueBoolPointer(),
			StaticPrefix:       prefixAssignment.StaticPrefix.ValueStringPointer(),
			StaticApplianceIp6: prefixAssignment.StaticApplianceIp6.ValueStringPointer(),
			Origin: &openApiClient.NetworksNetworkIdApplianceSingleLanIpv6Origin{
				Type:       prefixAssignment.Origin.Type.ValueString(),
				Interfaces: originInterfaces,
			},
		})
	}
	ipv6.SetPrefixAssignments(prefixAssignments)
	object56.SetIpv6(*ipv6)

	_, httpResp, err := r.client.ApplianceApi.CreateNetworkApplianceVlan(ctx, data.NetworkId.ValueString()).CreateNetworkApplianceVlan(*object56).Execute()
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

	// save inlineResp data into Terraform state.
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON Decode issue",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")
}

func (r *NetworksApplianceVlansResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetworksApplianceVlansResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := r.client.ApplianceApi.GetNetworkApplianceVlans(ctx, data.NetworkId.ValueString()).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}
	// collect diagnostics
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
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

func (r *NetworksApplianceVlansResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *NetworksApplianceVlansResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	object56 := openApiClient.NewInlineObject56(data.Id.ValueString(), data.Name.ValueString())
	object56.SetCidr(data.Cidr.ValueString())
	object56.SetId(data.VlanId.ValueString())
	object56.SetApplianceIp(data.ApplianceIp.ValueString())
	object56.SetGroupPolicyId(data.GroupPolicyId.ValueString())
	object56.SetMask(int32(data.Mask.ValueInt64()))
	object56.SetName(data.Name.ValueString())
	object56.SetSubnet(data.Subnet.ValueString())
	object56.SetTemplateVlanType(data.TemplateVlanType.ValueString())
	ipv6 := openApiClient.NewNetworksNetworkIdApplianceSingleLanIpv6()
	ipv6.SetEnabled(data.Ipv6.Enabled.ValueBool())
	var prefixAssignments []openApiClient.NetworksNetworkIdApplianceSingleLanIpv6PrefixAssignments
	for _, prefixAssignment := range data.Ipv6.PrefixAssignments {
		originInterfaces := []string{}
		for _, originInterface := range prefixAssignment.Origin.Interfaces.Elements() {
			originInterfaces = append(originInterfaces, originInterface.String())
		}
		prefixAssignments = append(prefixAssignments, openApiClient.NetworksNetworkIdApplianceSingleLanIpv6PrefixAssignments{
			Autonomous:         prefixAssignment.Autonomous.ValueBoolPointer(),
			StaticPrefix:       prefixAssignment.StaticPrefix.ValueStringPointer(),
			StaticApplianceIp6: prefixAssignment.StaticApplianceIp6.ValueStringPointer(),
			Origin: &openApiClient.NetworksNetworkIdApplianceSingleLanIpv6Origin{
				Type:       prefixAssignment.Origin.Type.ValueString(),
				Interfaces: originInterfaces,
			},
		})
	}
	ipv6.SetPrefixAssignments(prefixAssignments)
	object56.SetIpv6(*ipv6)

	_, httpResp, err := r.client.ApplianceApi.CreateNetworkApplianceVlan(ctx, data.NetworkId.ValueString()).CreateNetworkApplianceVlan(*object56).Execute()
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

	// save inlineResp data into Terraform state.
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON Decode issue",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *NetworksApplianceVlansResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *NetworksApplianceVlansResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	object56 := openApiClient.NewInlineObject56(data.Id.ValueString(), data.Name.ValueString())
	object56.SetCidr(data.Cidr.ValueString())
	object56.SetId(data.VlanId.ValueString())
	object56.SetApplianceIp(data.ApplianceIp.ValueString())
	object56.SetGroupPolicyId(data.GroupPolicyId.ValueString())
	object56.SetMask(int32(data.Mask.ValueInt64()))
	object56.SetName(data.Name.ValueString())
	object56.SetSubnet(data.Subnet.ValueString())
	object56.SetTemplateVlanType(data.TemplateVlanType.ValueString())
	ipv6 := openApiClient.NewNetworksNetworkIdApplianceSingleLanIpv6()
	ipv6.SetEnabled(data.Ipv6.Enabled.ValueBool())
	var prefixAssignments []openApiClient.NetworksNetworkIdApplianceSingleLanIpv6PrefixAssignments
	for _, prefixAssignment := range data.Ipv6.PrefixAssignments {
		originInterfaces := []string{}
		for _, originInterface := range prefixAssignment.Origin.Interfaces.Elements() {
			originInterfaces = append(originInterfaces, originInterface.String())
		}
		prefixAssignments = append(prefixAssignments, openApiClient.NetworksNetworkIdApplianceSingleLanIpv6PrefixAssignments{
			Autonomous:         prefixAssignment.Autonomous.ValueBoolPointer(),
			StaticPrefix:       prefixAssignment.StaticPrefix.ValueStringPointer(),
			StaticApplianceIp6: prefixAssignment.StaticApplianceIp6.ValueStringPointer(),
			Origin: &openApiClient.NetworksNetworkIdApplianceSingleLanIpv6Origin{
				Type:       prefixAssignment.Origin.Type.ValueString(),
				Interfaces: originInterfaces,
			},
		})
	}
	ipv6.SetPrefixAssignments(prefixAssignments)
	object56.SetIpv6(*ipv6)

	_, httpResp, err := r.client.ApplianceApi.CreateNetworkApplianceVlan(ctx, data.NetworkId.ValueString()).CreateNetworkApplianceVlan(*object56).Execute()
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

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%s", data))
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

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")

}

func (r *NetworksApplianceVlansResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), req.ID)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
