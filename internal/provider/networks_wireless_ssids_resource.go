package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
	"strings"
)

// TODO - DON'T FORGET TO DELETE ALL "TODO" COMMENTS!

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &NetworksWirelessSsidsResource{}
	_ resource.ResourceWithConfigure   = &NetworksWirelessSsidsResource{}
	_ resource.ResourceWithImportState = &NetworksWirelessSsidsResource{}
)

// TODO - This function needs to be added to the list of Data Sources in provider.go: func (p *ScaffoldingProvider) Resources
// TODO - Otherwise the provider has no idea this resource exists.
func NewNetworksWirelessSsidsResource() resource.Resource {
	return &NetworksWirelessSsidsResource{}
}

// NetworksWirelessSsidsResource defines the resource implementation.
type NetworksWirelessSsidsResource struct {
	client *openApiClient.APIClient
}

type Dot11W struct {
	Enabled  jsontypes.Bool `tfsdk:"enabled" json:"enabled"`
	Required jsontypes.Bool `tfsdk:"required" json:"required"`
}

type Dot11R struct {
	Enabled  jsontypes.Bool `tfsdk:"enabled" json:"enabled"`
	Adaptive jsontypes.Bool `tfsdk:"adaptive" json:"adaptive"`
}

type ClientRootCaCertificate struct {
	Contents jsontypes.String `tfsdk:"contents" json:"contents"`
}

type CertificateAuthentication struct {
	Enabled                 jsontypes.Bool   `tfsdk:"enabled" json:"enabled"`
	UseLdap                 jsontypes.Bool   `tfsdk:"use_ldap" json:"useLdap"`
	UseOcsp                 jsontypes.Bool   `tfsdk:"use_ocsp" json:"useOcsp"`
	OcspResponderUrl        jsontypes.String `tfsdk:"ocsp_responder_url" json:"ocspResponderUrl"`
	ClientRootCaCertificate types.Object     `tfsdk:"client_root_ca_certificate" json:"clientRootCaCertificate"`
}

type PasswordAuthentication struct {
	Enabled jsontypes.Bool `tfsdk:"enabled" json:"enabled"`
}

type LocalRadius struct {
	CacheTimeout              jsontypes.Int64 `tfsdk:"cache_timeout" json:"cacheTimeout"`
	PasswordAuthentication    types.Object    `tfsdk:"password_authentication" json:"passwordAuthentication"`
	CertificateAuthentication types.Object    `tfsdk:"certificate_authentication" json:"certificateAuthentication"`
}

type Server struct {
	Host jsontypes.String `tfsdk:"host" json:"host"`
	Port jsontypes.Int64  `tfsdk:"port" json:"port"`
}

type Credential struct {
	DistinguishedName jsontypes.String `tfsdk:"distinguished_name" json:"distinguishedName"`
	Password          jsontypes.String `tfsdk:"password" json:"password"`
}

type ServerCaCertificate struct {
	Contents jsontypes.String `tfsdk:"contents" json:"contents"`
}

type Ldap struct {
	Servers               []Server         `tfsdk:"servers" json:"servers"`
	Credentials           types.Object     `tfsdk:"credentials" json:"credentials"`
	BaseDistinguishedName jsontypes.String `tfsdk:"base_distinguished_name" json:"baseDistinguishedName"`
	ServerCaCertificate   types.Object     `tfsdk:"server_ca_certificate" json:"serverCaCertificate"`
}

type ADCredential struct {
	LogonName jsontypes.String `tfsdk:"logon_name" json:"logonName"`
	Password  jsontypes.String `tfsdk:"password" json:"password"`
}

type ActiveDirectory struct {
	Servers     []Server     `tfsdk:"servers" json:"servers"`
	Credentials types.Object `tfsdk:"credentials" json:"credentials"`
}

type RadiusServer struct {
	Host                     jsontypes.String `tfsdk:"host" json:"host"`
	Secret                   jsontypes.String `tfsdk:"secret" json:"secret"`
	CaCertificate            jsontypes.String `tfsdk:"ca_certificate" json:"caCertificate"`
	Port                     jsontypes.Int64  `tfsdk:"port" json:"port"`
	OpenRoamingCertificateId jsontypes.Int64  `tfsdk:"open_roaming_certificate_id" json:"openRoamingCertificateId"`
	RadsecEnabled            jsontypes.Bool   `tfsdk:"radsec_enabled" json:"radsecEnabled"`
}

type RadiusAccountingServer struct {
	Host          jsontypes.String `tfsdk:"host" json:"host"`
	Secret        jsontypes.String `tfsdk:"secret" json:"secret"`
	CaCertificate jsontypes.String `tfsdk:"ca_certificate" json:"caCertificate"`
	Port          jsontypes.Int64  `tfsdk:"port" json:"port"`
	RadsecEnabled jsontypes.Bool   `tfsdk:"radsec_enabled" json:"radsecEnabled"`
}

type Concentrator struct {
	Host jsontypes.String `tfsdk:"host" json:"host"`
}

type Gre struct {
	Concentrator types.Object    `tfsdk:"concentrator" json:"concentrator"`
	Key          jsontypes.Int64 `tfsdk:"key" json:"key"`
}

type ApTagsAndVlanId struct {
	Tags   types.List      `tfsdk:"tags" json:"tags"`
	VlanId jsontypes.Int64 `tfsdk:"vlan_id" json:"vlanId"`
}

type DNSRewrite struct {
	Enabled              jsontypes.Bool `tdsdk:"enabled" json:"enabled"`
	DnsCustomNameservers types.List     `tdsdk:"dns_custom_nameservers" json:"dnsCustomNameservers"`
}

type ByApTag struct {
	Tags     types.List       `tfsdk:"tags" json:"tags"`
	VlanName jsontypes.String `tfsdk:"vlan_name" json:"vlanName"`
}

type Tagging struct {
	Enabled         jsontypes.Bool   `tfsdk:"enabled" json:"enabled"`
	DefaultVlanName jsontypes.String `tfsdk:"default_vlan_name" json:"defaultVlanName"`
	ByApTags        []ByApTag        `tfsdk:"by_ap_tags" json:"byApTags"`
}

type GuestVlan struct {
	Enabled jsontypes.Bool   `tfsdk:"enabled" json:"enabled"`
	Name    jsontypes.String `tfsdk:"name" json:"name"`
}

type Radius struct {
	GuestVlan types.Object `tfsdk:"guest_vlan" json:"guestVlan"`
}

type NamedVlans struct {
	Tagging types.Object `tfsdk:"tagging" json:"tagging"`
	Radius  Radius       `tfsdk:"radius" json:"radius"`
}

type SpeedBurst struct {
	Enabled jsontypes.Bool `tfsdk:"enabled" json:"enabled"`
}

type Oauth struct {
	AllowedDomains types.List `json:"allowed_domains"`
}

type WirelessNetworkSSID struct {
	Name                             jsontypes.String         `tfsdk:"name"`
	AuthMode                         jsontypes.String         `tfsdk:"auth_mode"`
	EnterpriseAdminAccess            jsontypes.String         `tfsdk:"enterprise_admin_access"`
	EncryptionMode                   jsontypes.String         `tfsdk:"encryption_mode"`
	Psk                              jsontypes.String         `tfsdk:"psk"`
	WpaEncryptionMode                jsontypes.String         `tfsdk:"wpa_encryption_mode"`
	SplashPage                       jsontypes.String         `tfsdk:"splash_page"`
	RadiusCalledStationId            jsontypes.String         `tfsdk:"radius_called_station_id"`
	RadiusAuthenticationNasId        jsontypes.String         `tfsdk:"radius_authentication_nas_id"`
	RadiusFailoverPolicy             jsontypes.String         `tfsdk:"radius_failover_policy"`
	RadiusLoadBalancingPolicy        jsontypes.String         `tfsdk:"radius_load_balancing_policy"`
	RadiusAttributeForGroupPolicies  jsontypes.String         `tfsdk:"radius_attribute_for_group_policies"`
	IpAssignmentMode                 jsontypes.String         `tfsdk:"ip_assignment_mode"`
	ConcentratorNetworkId            jsontypes.String         `tfsdk:"concentrator_network_id"`
	SecondaryConcentratorNetworkId   jsontypes.String         `tfsdk:"secondary_concentrator_network_id"`
	BandSelection                    jsontypes.String         `tfsdk:"band_selection"`
	RadiusServerTimeout              jsontypes.Int64          `tfsdk:"radius_server_timeout"`
	RadiusServerAttemptsLimit        jsontypes.Int64          `tfsdk:"radius_server_attempts_limit"`
	RadiusAccountingInterimInterval  jsontypes.Int64          `tfsdk:"radius_accounting_interim_interval"`
	VlanId                           jsontypes.Int64          `tfsdk:"vlan_id"`
	DefaultVlanId                    jsontypes.Int64          `tfsdk:"default_vlan_id"`
	PerClientBandwidthLimitUp        jsontypes.Int64          `tfsdk:"per_client_bandwidth_limit_up"`
	PerClientBandwidthLimitDown      jsontypes.Int64          `tfsdk:"per_client_bandwidth_limit_down"`
	PerSsidBandwidthLimitUp          jsontypes.Int64          `tfsdk:"per_ssid_bandwidth_limit_up"`
	PerSsidBandwidthLimitDown        jsontypes.Int64          `tfsdk:"per_ssid_bandwidth_limit_down"`
	RadiusGuestVlanId                jsontypes.Int64          `tfsdk:"radius_guest_vlan_id"`
	MinBitrate                       jsontypes.Float64        `tfsdk:"min_bitrate"`
	UseVlanTagging                   jsontypes.Bool           `tfsdk:"use_vlan_tagging"`
	DisassociateClientsOnVpnFailover jsontypes.Bool           `tfsdk:"disassociate_clients_on_vpn_failover"`
	RadiusOverride                   jsontypes.Bool           `tfsdk:"radius_override"`
	RadiusGuestVlanEnabled           jsontypes.Bool           `tfsdk:"radius_guest_vlan_enabled"`
	Enabled                          jsontypes.Bool           `tfsdk:"enabled"`
	RadiusProxyEnabled               jsontypes.Bool           `tfsdk:"radius_proxy_enabled"`
	RadiusTestingEnabled             jsontypes.Bool           `tfsdk:"radius_testing_enabled"`
	RadiusFallbackEnabled            jsontypes.Bool           `tfsdk:"radius_fallback_enabled"`
	RadiusCoaEnabled                 jsontypes.Bool           `tfsdk:"radius_coa_enabled"`
	RadiusAccountingEnabled          jsontypes.Bool           `tfsdk:"radius_accounting_enabled"`
	LanIsolationEnabled              jsontypes.Bool           `tfsdk:"lan_isolation_enabled"`
	Visible                          jsontypes.Bool           `tfsdk:"visible"`
	AvailableOnAllAps                jsontypes.Bool           `tfsdk:"available_on_all_aps"`
	MandatoryDhcpEnabled             jsontypes.Bool           `tfsdk:"mandatory_dhcp_enabled"`
	AdultContentFilteringEnabled     jsontypes.Bool           `tfsdk:"adult_content_filtering_enabled"`
	WalledGardenEnabled              jsontypes.Bool           `tfsdk:"walled_garden_enabled"`
	Dot11W                           types.Object             `tfsdk:"dot11w"`
	Dot11R                           types.Object             `tfsdk:"dot11r"`
	LocalRadius                      types.Object             `tfsdk:"local_radius"`
	Ldap                             types.Object             `tfsdk:"ldap"`
	ActiveDirectory                  types.Object             `tfsdk:"active_directory"`
	DnsRewrite                       types.Object             `tfsdk:"dns_rewrite"`
	SpeedBurst                       types.Object             `tfsdk:"speed_burst"`
	Gre                              types.Object             `tfsdk:"gre"`
	Oauth                            types.Object             `tfsdk:"oauth"`
	SplashGuestSponsorDomains        types.List               `tfsdk:"splash_guest_sponsor_domains"`
	WalledGardenRanges               types.List               `tfsdk:"walled_garden_ranges"`
	AvailabilityTags                 types.List               `tfsdk:"availability_tags"`
	RadiusServers                    []RadiusServer           `tfsdk:"radius_servers"`
	RadiusAccountingServers          []RadiusAccountingServer `tfsdk:"radius_accounting_servers"`
	ApTagsAndVlanIds                 []ApTagsAndVlanId        `tfsdk:"ap_tags_and_vlan_ids"`
}

// NetworksWirelessSsidsResourceModel describes the resource data model.
type NetworksWirelessSsidsResourceModel struct {
	Id jsontypes.String `tfsdk:"id"`

	NetworkID jsontypes.String `tfsdk:"network_id"`
	Number    jsontypes.String `tfsdk:"number"`
	SSID      types.Object     `tfsdk:"ssid"`
}

func (r *NetworksWirelessSsidsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_wireless_ssids"
}

func (r *NetworksWirelessSsidsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{

		MarkdownDescription: "NetworksWirelessSsids",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.StringType,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Network ID",
				Required:            true,
				CustomType:          jsontypes.StringType,
			},
			"number": schema.StringAttribute{
				MarkdownDescription: "Number",
				Required:            true,
				CustomType:          jsontypes.StringType,
			},
			"ssid": schema.SingleNestedAttribute{
				Required: true,
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						MarkdownDescription: "Enable or disable the interface (only for MX devices). Valid values are 'enabled', 'disabled', and 'not configured'.",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"auth_mode": schema.StringAttribute{
						MarkdownDescription: "Enable or disable the interface (only for MX devices). Valid values are 'enabled', 'disabled', and 'not configured'.",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"enterprise_admin_access": schema.StringAttribute{
						MarkdownDescription: "Enable or disable the interface (only for MX devices). Valid values are 'enabled', 'disabled', and 'not configured'.",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"encryption_mode": schema.StringAttribute{
						MarkdownDescription: "Enable or disable the interface (only for MX devices). Valid values are 'enabled', 'disabled', and 'not configured'.",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"psk": schema.StringAttribute{
						MarkdownDescription: "Enable or disable the interface (only for MX devices). Valid values are 'enabled', 'disabled', and 'not configured'.",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"wpa_encryption_mode": schema.StringAttribute{
						MarkdownDescription: "Enable or disable the interface (only for MX devices). Valid values are 'enabled', 'disabled', and 'not configured'.",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"splash_page": schema.StringAttribute{
						MarkdownDescription: "Enable or disable the interface (only for MX devices). Valid values are 'enabled', 'disabled', and 'not configured'.",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"radius_called_station_id": schema.StringAttribute{
						MarkdownDescription: "Enable or disable the interface (only for MX devices). Valid values are 'enabled', 'disabled', and 'not configured'.",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"radius_authentication_nas_id": schema.StringAttribute{
						MarkdownDescription: "Enable or disable the interface (only for MX devices). Valid values are 'enabled', 'disabled', and 'not configured'.",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"radius_failover_policy": schema.StringAttribute{
						MarkdownDescription: "Enable or disable the interface (only for MX devices). Valid values are 'enabled', 'disabled', and 'not configured'.",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"radius_load_balancing_policy": schema.StringAttribute{
						MarkdownDescription: "Enable or disable the interface (only for MX devices). Valid values are 'enabled', 'disabled', and 'not configured'.",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"radius_attribute_for_group_policies": schema.StringAttribute{
						MarkdownDescription: "Enable or disable the interface (only for MX devices). Valid values are 'enabled', 'disabled', and 'not configured'.",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"ip_assignment_mode": schema.StringAttribute{
						MarkdownDescription: "Enable or disable the interface (only for MX devices). Valid values are 'enabled', 'disabled', and 'not configured'.",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"concentrator_network_id": schema.StringAttribute{
						MarkdownDescription: "Enable or disable the interface (only for MX devices). Valid values are 'enabled', 'disabled', and 'not configured'.",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"secondary_concentrator_network_id": schema.StringAttribute{
						MarkdownDescription: "Enable or disable the interface (only for MX devices). Valid values are 'enabled', 'disabled', and 'not configured'.",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"band_selection": schema.StringAttribute{
						MarkdownDescription: "Enable or disable the interface (only for MX devices). Valid values are 'enabled', 'disabled', and 'not configured'.",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"radius_server_timeout": schema.Int64Attribute{
						MarkdownDescription: "The VLAN that management traffic should be tagged with. Applies whether usingStaticIp is true or false.",
						Optional:            true,
						CustomType:          jsontypes.Int64Type,
					},
					"radius_server_attempts_limit": schema.Int64Attribute{
						MarkdownDescription: "The VLAN that management traffic should be tagged with. Applies whether usingStaticIp is true or false.",
						Optional:            true,
						CustomType:          jsontypes.Int64Type,
					},
					"radius_accounting_interim_interval": schema.Int64Attribute{
						MarkdownDescription: "The VLAN that management traffic should be tagged with. Applies whether usingStaticIp is true or false.",
						Optional:            true,
						CustomType:          jsontypes.Int64Type,
					},
					"vlan_id": schema.Int64Attribute{
						MarkdownDescription: "The VLAN that management traffic should be tagged with. Applies whether usingStaticIp is true or false.",
						Optional:            true,
						CustomType:          jsontypes.Int64Type,
					},
					"default_vlan_id": schema.Int64Attribute{
						MarkdownDescription: "The VLAN that management traffic should be tagged with. Applies whether usingStaticIp is true or false.",
						Optional:            true,
						CustomType:          jsontypes.Int64Type,
					},
					"per_client_bandwidth_limit_up": schema.Int64Attribute{
						MarkdownDescription: "The VLAN that management traffic should be tagged with. Applies whether usingStaticIp is true or false.",
						Optional:            true,
						CustomType:          jsontypes.Int64Type,
					},
					"per_client_bandwidth_limit_down": schema.Int64Attribute{
						MarkdownDescription: "The VLAN that management traffic should be tagged with. Applies whether usingStaticIp is true or false.",
						Optional:            true,
						CustomType:          jsontypes.Int64Type,
					},
					"per_ssid_bandwidth_limit_up": schema.Int64Attribute{
						MarkdownDescription: "The VLAN that management traffic should be tagged with. Applies whether usingStaticIp is true or false.",
						Optional:            true,
						CustomType:          jsontypes.Int64Type,
					},
					"per_ssid_bandwidth_limit_down": schema.Int64Attribute{
						MarkdownDescription: "The VLAN that management traffic should be tagged with. Applies whether usingStaticIp is true or false.",
						Optional:            true,
						CustomType:          jsontypes.Int64Type,
					},
					"radius_guest_vlan_id": schema.Int64Attribute{
						MarkdownDescription: "The VLAN that management traffic should be tagged with. Applies whether usingStaticIp is true or false.",
						Optional:            true,
						CustomType:          jsontypes.Int64Type,
					},
					"min_bitrate": schema.Float64Attribute{
						MarkdownDescription: "The VLAN that management traffic should be tagged with. Applies whether usingStaticIp is true or false.",
						Optional:            true,
						CustomType:          jsontypes.Float64Type,
					},
					"use_vlan_tagging": schema.BoolAttribute{
						MarkdownDescription: "Configure the interface to have static IP settings or use DHCP.",
						Optional:            true,
						CustomType:          jsontypes.BoolType,
					},
					"disassociate_clients_on_vpn_failover": schema.BoolAttribute{
						MarkdownDescription: "Configure the interface to have static IP settings or use DHCP.",
						Optional:            true,
						CustomType:          jsontypes.BoolType,
					},
					"radius_override": schema.BoolAttribute{
						MarkdownDescription: "Configure the interface to have static IP settings or use DHCP.",
						Optional:            true,
						CustomType:          jsontypes.BoolType,
					},
					"radius_guest_vlan_enabled": schema.BoolAttribute{
						MarkdownDescription: "Configure the interface to have static IP settings or use DHCP.",
						Optional:            true,
						CustomType:          jsontypes.BoolType,
					},
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "Configure the interface to have static IP settings or use DHCP.",
						Optional:            true,
						CustomType:          jsontypes.BoolType,
					},
					"radius_proxy_enabled": schema.BoolAttribute{
						MarkdownDescription: "Configure the interface to have static IP settings or use DHCP.",
						Optional:            true,
						CustomType:          jsontypes.BoolType,
					},
					"radius_testing_enabled": schema.BoolAttribute{
						MarkdownDescription: "Configure the interface to have static IP settings or use DHCP.",
						Optional:            true,
						CustomType:          jsontypes.BoolType,
					},
					"radius_fallback_enabled": schema.BoolAttribute{
						MarkdownDescription: "Configure the interface to have static IP settings or use DHCP.",
						Optional:            true,
						CustomType:          jsontypes.BoolType,
					},
					"radius_coa_enabled": schema.BoolAttribute{
						MarkdownDescription: "Configure the interface to have static IP settings or use DHCP.",
						Optional:            true,
						CustomType:          jsontypes.BoolType,
					},
					"radius_accounting_enabled": schema.BoolAttribute{
						MarkdownDescription: "Configure the interface to have static IP settings or use DHCP.",
						Optional:            true,
						CustomType:          jsontypes.BoolType,
					},
					"lan_isolation_enabled": schema.BoolAttribute{
						MarkdownDescription: "Configure the interface to have static IP settings or use DHCP.",
						Optional:            true,
						CustomType:          jsontypes.BoolType,
					},
					"visible": schema.BoolAttribute{
						MarkdownDescription: "Configure the interface to have static IP settings or use DHCP.",
						Optional:            true,
						CustomType:          jsontypes.BoolType,
					},
					"available_on_all_aps": schema.BoolAttribute{
						MarkdownDescription: "Configure the interface to have static IP settings or use DHCP.",
						Optional:            true,
						CustomType:          jsontypes.BoolType,
					},
					"mandatory_dhcp_enabled": schema.BoolAttribute{
						MarkdownDescription: "Configure the interface to have static IP settings or use DHCP.",
						Optional:            true,
						CustomType:          jsontypes.BoolType,
					},
					"adult_content_filtering_enabled": schema.BoolAttribute{
						MarkdownDescription: "Configure the interface to have static IP settings or use DHCP.",
						Optional:            true,
						CustomType:          jsontypes.BoolType,
					},
					"walled_garden_enabled": schema.BoolAttribute{
						MarkdownDescription: "Configure the interface to have static IP settings or use DHCP.",
						Optional:            true,
						CustomType:          jsontypes.BoolType,
					},
					"dot11w": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								MarkdownDescription: "Configure the interface to have static IP settings or use DHCP.",
								Optional:            true,
								CustomType:          jsontypes.BoolType,
							},
							"required": schema.BoolAttribute{
								MarkdownDescription: "Configure the interface to have static IP settings or use DHCP.",
								Optional:            true,
								CustomType:          jsontypes.BoolType,
							},
						},
					},
					"dot11r": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								MarkdownDescription: "Configure the interface to have static IP settings or use DHCP.",
								Optional:            true,
								CustomType:          jsontypes.BoolType,
							},
							"adaptive": schema.BoolAttribute{
								MarkdownDescription: "Configure the interface to have static IP settings or use DHCP.",
								Optional:            true,
								CustomType:          jsontypes.BoolType,
							},
						},
					},
					"local_radius": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"cache_timeout": schema.Int64Attribute{
								MarkdownDescription: "The VLAN that management traffic should be tagged with. Applies whether usingStaticIp is true or false.",
								Optional:            true,
								CustomType:          jsontypes.Int64Type,
							},
							"password_authentication": schema.SingleNestedAttribute{
								Optional: true,
								Attributes: map[string]schema.Attribute{
									"enabled": schema.BoolAttribute{
										MarkdownDescription: "Configure the interface to have static IP settings or use DHCP.",
										Optional:            true,
										CustomType:          jsontypes.BoolType,
									},
								},
							},
							"certificate_authentication": schema.SingleNestedAttribute{
								Optional: true,
								Attributes: map[string]schema.Attribute{
									"enabled": schema.BoolAttribute{
										MarkdownDescription: "Configure the interface to have static IP settings or use DHCP.",
										Optional:            true,
										CustomType:          jsontypes.BoolType,
									},
									"use_ldap": schema.BoolAttribute{
										MarkdownDescription: "Configure the interface to have static IP settings or use DHCP.",
										Optional:            true,
										CustomType:          jsontypes.BoolType,
									},
									"use_ocsp": schema.BoolAttribute{
										MarkdownDescription: "Configure the interface to have static IP settings or use DHCP.",
										Optional:            true,
										CustomType:          jsontypes.BoolType,
									},
									"ocsp_responder_url": schema.StringAttribute{
										MarkdownDescription: "Enable or disable the interface (only for MX devices). Valid values are 'enabled', 'disabled', and 'not configured'.",
										Optional:            true,
										CustomType:          jsontypes.StringType,
									},
									"client_root_ca_certificate": schema.SingleNestedAttribute{
										Optional: true,
										Attributes: map[string]schema.Attribute{
											"contents": schema.StringAttribute{
												MarkdownDescription: "Enable or disable the interface (only for MX devices). Valid values are 'enabled', 'disabled', and 'not configured'.",
												Optional:            true,
												CustomType:          jsontypes.StringType,
											},
										},
									},
								},
							},
						},
					},
					"ldap": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"servers": schema.SetNestedAttribute{
								Optional: true,
								Computed: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"host": schema.StringAttribute{
											MarkdownDescription: "Enable or disable the interface (only for MX devices). Valid values are 'enabled', 'disabled', and 'not configured'.",
											Optional:            true,
											CustomType:          jsontypes.StringType,
										},
										"port": schema.Int64Attribute{
											MarkdownDescription: "The VLAN that management traffic should be tagged with. Applies whether usingStaticIp is true or false.",
											Optional:            true,
											CustomType:          jsontypes.Int64Type,
										},
									},
								},
							},
							"credentials": schema.SingleNestedAttribute{
								Optional: true,
								Attributes: map[string]schema.Attribute{
									"distinguished_name": schema.StringAttribute{
										MarkdownDescription: "Enable or disable the interface (only for MX devices). Valid values are 'enabled', 'disabled', and 'not configured'.",
										Optional:            true,
										CustomType:          jsontypes.StringType,
									},
									"password": schema.StringAttribute{
										MarkdownDescription: "Enable or disable the interface (only for MX devices). Valid values are 'enabled', 'disabled', and 'not configured'.",
										Optional:            true,
										CustomType:          jsontypes.StringType,
									},
								},
							},
							"base_distinguished_name": schema.StringAttribute{
								MarkdownDescription: "Enable or disable the interface (only for MX devices). Valid values are 'enabled', 'disabled', and 'not configured'.",
								Optional:            true,
								CustomType:          jsontypes.StringType,
							},
							"server_ca_certificate": schema.SingleNestedAttribute{
								Optional: true,
								Attributes: map[string]schema.Attribute{
									"contents": schema.StringAttribute{
										MarkdownDescription: "Enable or disable the interface (only for MX devices). Valid values are 'enabled', 'disabled', and 'not configured'.",
										Optional:            true,
										CustomType:          jsontypes.StringType,
									},
								},
							},
						},
					},
					"active_directory": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"servers": schema.SetNestedAttribute{
								Optional: true,
								Computed: true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"host": schema.StringAttribute{
											MarkdownDescription: "Enable or disable the interface (only for MX devices). Valid values are 'enabled', 'disabled', and 'not configured'.",
											Optional:            true,
											CustomType:          jsontypes.StringType,
										},
										"port": schema.Int64Attribute{
											MarkdownDescription: "The VLAN that management traffic should be tagged with. Applies whether usingStaticIp is true or false.",
											Optional:            true,
											CustomType:          jsontypes.Int64Type,
										},
									},
								},
							},
							"credentials": schema.SingleNestedAttribute{
								Optional: true,
								Attributes: map[string]schema.Attribute{
									"logon_name": schema.StringAttribute{
										MarkdownDescription: "Enable or disable the interface (only for MX devices). Valid values are 'enabled', 'disabled', and 'not configured'.",
										Optional:            true,
										CustomType:          jsontypes.StringType,
									},
									"password": schema.StringAttribute{
										MarkdownDescription: "Enable or disable the interface (only for MX devices). Valid values are 'enabled', 'disabled', and 'not configured'.",
										Optional:            true,
										CustomType:          jsontypes.StringType,
									},
								},
							},
						},
					},
					"dns_rewrite": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								MarkdownDescription: "Configure the interface to have static IP settings or use DHCP.",
								Optional:            true,
								CustomType:          jsontypes.BoolType,
							},
							"dns_custom_nameservers": schema.ListAttribute{
								MarkdownDescription: "Up to two DNS IPs.",
								Optional:            true,
								ElementType:         jsontypes.StringType,
							},
						},
					},
					"speed_burst": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								MarkdownDescription: "Configure the interface to have static IP settings or use DHCP.",
								Optional:            true,
								CustomType:          jsontypes.BoolType,
							},
						},
					},
					"gre": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"concentrator": schema.SingleNestedAttribute{
								Optional: true,
								Attributes: map[string]schema.Attribute{
									"host": schema.StringAttribute{
										MarkdownDescription: "Enable or disable the interface (only for MX devices). Valid values are 'enabled', 'disabled', and 'not configured'.",
										Optional:            true,
										CustomType:          jsontypes.StringType,
									},
								},
							},
							"key": schema.Int64Attribute{
								MarkdownDescription: "The VLAN that management traffic should be tagged with. Applies whether usingStaticIp is true or false.",
								Optional:            true,
								CustomType:          jsontypes.Int64Type,
							},
						},
					},
					"oauth": schema.SingleNestedAttribute{
						Optional: true,
						Attributes: map[string]schema.Attribute{
							"allowed_domains": schema.ListAttribute{
								MarkdownDescription: "Up to two DNS IPs.",
								Optional:            true,
								ElementType:         jsontypes.StringType,
							},
						},
					},
					"splash_guest_sponsor_domains": schema.ListAttribute{
						MarkdownDescription: "Up to two DNS IPs.",
						Optional:            true,
						ElementType:         jsontypes.StringType,
					},
					"walled_garden_ranges": schema.ListAttribute{
						MarkdownDescription: "Up to two DNS IPs.",
						Optional:            true,
						ElementType:         jsontypes.StringType,
					},
					"availability_tags": schema.ListAttribute{
						MarkdownDescription: "Up to two DNS IPs.",
						Optional:            true,
						ElementType:         jsontypes.StringType,
					},
					"radius_servers": schema.SetNestedAttribute{
						Optional: true,
						Computed: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"host": schema.StringAttribute{
									MarkdownDescription: "Enable or disable the interface (only for MX devices). Valid values are 'enabled', 'disabled', and 'not configured'.",
									Optional:            true,
									CustomType:          jsontypes.StringType,
								},
								"secret": schema.StringAttribute{
									MarkdownDescription: "Enable or disable the interface (only for MX devices). Valid values are 'enabled', 'disabled', and 'not configured'.",
									Optional:            true,
									CustomType:          jsontypes.StringType,
								},
								"ca_certificate": schema.StringAttribute{
									MarkdownDescription: "Enable or disable the interface (only for MX devices). Valid values are 'enabled', 'disabled', and 'not configured'.",
									Optional:            true,
									CustomType:          jsontypes.StringType,
								},
								"port": schema.Int64Attribute{
									MarkdownDescription: "The VLAN that management traffic should be tagged with. Applies whether usingStaticIp is true or false.",
									Optional:            true,
									CustomType:          jsontypes.Int64Type,
								},
								"open_roaming_certificate_id": schema.Int64Attribute{
									MarkdownDescription: "The VLAN that management traffic should be tagged with. Applies whether usingStaticIp is true or false.",
									Optional:            true,
									CustomType:          jsontypes.Int64Type,
								},
								"radsec_enabled": schema.BoolAttribute{
									MarkdownDescription: "Configure the interface to have static IP settings or use DHCP.",
									Optional:            true,
									CustomType:          jsontypes.BoolType,
								},
							},
						},
					},
					"radius_accounting_servers": schema.SetNestedAttribute{
						Optional: true,
						Computed: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"host": schema.StringAttribute{
									MarkdownDescription: "Enable or disable the interface (only for MX devices). Valid values are 'enabled', 'disabled', and 'not configured'.",
									Optional:            true,
									CustomType:          jsontypes.StringType,
								},
								"secret": schema.StringAttribute{
									MarkdownDescription: "Enable or disable the interface (only for MX devices). Valid values are 'enabled', 'disabled', and 'not configured'.",
									Optional:            true,
									CustomType:          jsontypes.StringType,
								},
								"ca_certificate": schema.StringAttribute{
									MarkdownDescription: "Enable or disable the interface (only for MX devices). Valid values are 'enabled', 'disabled', and 'not configured'.",
									Optional:            true,
									CustomType:          jsontypes.StringType,
								},
								"port": schema.Int64Attribute{
									MarkdownDescription: "The VLAN that management traffic should be tagged with. Applies whether usingStaticIp is true or false.",
									Optional:            true,
									CustomType:          jsontypes.Int64Type,
								},
								"radsec_enabled": schema.BoolAttribute{
									MarkdownDescription: "Configure the interface to have static IP settings or use DHCP.",
									Optional:            true,
									CustomType:          jsontypes.BoolType,
								},
							},
						},
					},
					"ap_tags_and_vlan_ids": schema.SetNestedAttribute{
						Optional: true,
						Computed: true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"vlan_id": schema.Int64Attribute{
									MarkdownDescription: "The VLAN that management traffic should be tagged with. Applies whether usingStaticIp is true or false.",
									Optional:            true,
									CustomType:          jsontypes.Int64Type,
								},
								"tags": schema.ListAttribute{
									MarkdownDescription: "Up to two DNS IPs.",
									Optional:            true,
									ElementType:         jsontypes.StringType,
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *NetworksWirelessSsidsResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NetworksWirelessSsidsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *NetworksWirelessSsidsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	request := openApiClient.NewUpdateNetworkWirelessSsidRequest()
	var ssid WirelessNetworkSSID
	data.SSID.As(ctx, &ssid, basetypes.ObjectAsOptions{})

	request.SetName(ssid.Name.ValueString())
	request.SetAuthMode(ssid.AuthMode.ValueString())
	request.SetEnterpriseAdminAccess(ssid.EnterpriseAdminAccess.ValueString())
	request.SetEncryptionMode(ssid.EncryptionMode.ValueString())
	request.SetPsk(ssid.Psk.ValueString())
	request.SetWpaEncryptionMode(ssid.WpaEncryptionMode.ValueString())
	request.SetSplashPage(ssid.SplashPage.ValueString())
	request.SetRadiusCalledStationId(ssid.RadiusCalledStationId.ValueString())
	request.SetRadiusAuthenticationNasId(ssid.RadiusAuthenticationNasId.ValueString())
	request.SetRadiusFailoverPolicy(ssid.RadiusFailoverPolicy.ValueString())
	request.SetRadiusLoadBalancingPolicy(ssid.RadiusLoadBalancingPolicy.ValueString())
	request.SetRadiusAttributeForGroupPolicies(ssid.RadiusAttributeForGroupPolicies.ValueString())
	request.SetIpAssignmentMode(ssid.IpAssignmentMode.ValueString())
	request.SetConcentratorNetworkId(ssid.ConcentratorNetworkId.ValueString())
	request.SetSecondaryConcentratorNetworkId(ssid.SecondaryConcentratorNetworkId.ValueString())
	request.SetBandSelection(ssid.BandSelection.ValueString())
	request.SetRadiusServerTimeout(int32(ssid.RadiusServerTimeout.ValueInt64()))
	request.SetRadiusServerAttemptsLimit(int32(ssid.RadiusServerAttemptsLimit.ValueInt64()))
	request.SetRadiusAccountingInterimInterval(int32(ssid.RadiusAccountingInterimInterval.ValueInt64()))
	request.SetVlanId(int32(ssid.VlanId.ValueInt64()))
	request.SetDefaultVlanId(int32(ssid.DefaultVlanId.ValueInt64()))
	request.SetPerClientBandwidthLimitUp(int32(ssid.PerClientBandwidthLimitUp.ValueInt64()))
	request.SetPerClientBandwidthLimitDown(int32(ssid.PerClientBandwidthLimitDown.ValueInt64()))
	request.SetPerSsidBandwidthLimitUp(int32(ssid.PerSsidBandwidthLimitUp.ValueInt64()))
	request.SetPerSsidBandwidthLimitDown(int32(ssid.PerSsidBandwidthLimitDown.ValueInt64()))
	request.SetRadiusGuestVlanId(int32(ssid.RadiusGuestVlanId.ValueInt64()))
	request.SetMinBitrate(float32(ssid.MinBitrate.ValueFloat64()))
	request.SetUseVlanTagging(ssid.UseVlanTagging.ValueBool())
	request.SetDisassociateClientsOnVpnFailover(ssid.DisassociateClientsOnVpnFailover.ValueBool())
	request.SetRadiusOverride(ssid.RadiusOverride.ValueBool())
	request.SetRadiusGuestVlanEnabled(ssid.RadiusGuestVlanEnabled.ValueBool())
	request.SetEnabled(ssid.Enabled.ValueBool())
	request.SetRadiusProxyEnabled(ssid.RadiusProxyEnabled.ValueBool())
	request.SetRadiusTestingEnabled(ssid.RadiusTestingEnabled.ValueBool())
	request.SetRadiusFallbackEnabled(ssid.RadiusFallbackEnabled.ValueBool())
	request.SetRadiusCoaEnabled(ssid.RadiusCoaEnabled.ValueBool())
	request.SetRadiusAccountingEnabled(ssid.RadiusAccountingEnabled.ValueBool())
	request.SetLanIsolationEnabled(ssid.LanIsolationEnabled.ValueBool())
	request.SetVisible(ssid.Visible.ValueBool())
	request.SetAvailableOnAllAps(ssid.AvailableOnAllAps.ValueBool())
	request.SetMandatoryDhcpEnabled(ssid.MandatoryDhcpEnabled.ValueBool())
	request.SetAdultContentFilteringEnabled(ssid.AdultContentFilteringEnabled.ValueBool())
	request.SetWalledGardenEnabled(ssid.WalledGardenEnabled.ValueBool())

	dot11w := openApiClient.NewUpdateNetworkApplianceSsidRequestDot11w()
	var dot11wObject Dot11W
	ssid.Dot11W.As(ctx, &dot11wObject, basetypes.ObjectAsOptions{})
	dot11w.SetEnabled(dot11wObject.Enabled.ValueBool())
	dot11w.SetRequired(dot11wObject.Required.ValueBool())
	request.SetDot11w(*dot11w)

	dot11r := openApiClient.NewUpdateNetworkWirelessSsidRequestDot11r()
	var dot11rObject Dot11R
	ssid.Dot11R.As(ctx, &dot11rObject, basetypes.ObjectAsOptions{})
	dot11r.SetEnabled(dot11rObject.Enabled.ValueBool())
	dot11r.SetAdaptive(dot11rObject.Adaptive.ValueBool())
	request.SetDot11r(*dot11r)

	localRadius := openApiClient.NewUpdateNetworkWirelessSsidRequestLocalRadius()
	var localRadiusObject LocalRadius
	ssid.LocalRadius.As(ctx, &localRadiusObject, basetypes.ObjectAsOptions{})
	localRadius.SetCacheTimeout(int32(localRadiusObject.CacheTimeout.ValueInt64()))
	certificateAuthentication := openApiClient.NewUpdateNetworkWirelessSsidRequestLocalRadiusCertificateAuthentication()
	var certificateAuthenticationObject CertificateAuthentication
	localRadiusObject.CertificateAuthentication.As(ctx, &certificateAuthenticationObject, basetypes.ObjectAsOptions{})
	certificateAuthentication.SetEnabled(certificateAuthenticationObject.Enabled.ValueBool())
	certificateAuthentication.SetOcspResponderUrl(certificateAuthenticationObject.OcspResponderUrl.ValueString())
	certificateAuthentication.SetUseLdap(certificateAuthenticationObject.UseLdap.ValueBool())
	certificateAuthentication.SetUseOcsp(certificateAuthenticationObject.UseOcsp.ValueBool())
	var clientRootCaCertificateObject ClientRootCaCertificate
	certificateAuthenticationObject.ClientRootCaCertificate.As(ctx, &clientRootCaCertificateObject, basetypes.ObjectAsOptions{})
	rootCACertificate := openApiClient.NewUpdateNetworkWirelessSsidRequestLocalRadiusCertificateAuthenticationClientRootCaCertificate()
	rootCACertificate.SetContents(clientRootCaCertificateObject.Contents.ValueString())
	certificateAuthentication.SetClientRootCaCertificate(*rootCACertificate)
	localRadius.SetCertificateAuthentication(*certificateAuthentication)
	request.SetLocalRadius(*localRadius)

	ldap := openApiClient.NewUpdateNetworkWirelessSsidRequestLdap()
	var ldapObject Ldap
	ssid.Ldap.As(ctx, &ldapObject, basetypes.ObjectAsOptions{})
	ldap.SetBaseDistinguishedName(ldapObject.BaseDistinguishedName.ValueString())
	ldapCredentials := openApiClient.NewUpdateNetworkWirelessSsidRequestLdapCredentials()
	var ldapCredentialsObject Credential
	ldapObject.Credentials.As(ctx, &ldapCredentialsObject, basetypes.ObjectAsOptions{})
	ldapCredentials.SetDistinguishedName(ldapCredentialsObject.DistinguishedName.ValueString())
	ldapCredentials.SetPassword(ldapCredentialsObject.Password.ValueString())
	ldap.SetCredentials(*ldapCredentials)
	var ldapServers []openApiClient.UpdateNetworkWirelessSsidRequestLdapServersInner
	for _, server := range ldapObject.Servers {
		inner := openApiClient.NewUpdateNetworkWirelessSsidRequestLdapServersInner(server.Host.ValueString(), int32(server.Port.ValueInt64()))
		ldapServers = append(ldapServers, *inner)
	}
	ldap.SetServers(ldapServers)
	var ldapServerCaCertificateObject ServerCaCertificate
	ldapObject.ServerCaCertificate.As(ctx, &ldapServerCaCertificateObject, basetypes.ObjectAsOptions{})
	ldapServerCaCertificate := openApiClient.NewUpdateNetworkWirelessSsidRequestLdapServerCaCertificate()
	ldapServerCaCertificate.SetContents(ldapServerCaCertificateObject.Contents.ValueString())
	ldap.SetServerCaCertificate(*ldapServerCaCertificate)
	request.SetLdap(*ldap)

	ad := openApiClient.NewUpdateNetworkWirelessSsidRequestActiveDirectory()
	var activeDirectoryObject ActiveDirectory
	ssid.ActiveDirectory.As(ctx, &activeDirectoryObject, basetypes.ObjectAsOptions{})
	var adServers []openApiClient.UpdateNetworkWirelessSsidRequestActiveDirectoryServersInner
	for _, server := range activeDirectoryObject.Servers {
		inner := openApiClient.NewUpdateNetworkWirelessSsidRequestActiveDirectoryServersInner(server.Host.ValueString())
		inner.SetPort(int32(server.Port.ValueInt64()))
		adServers = append(adServers, *inner)
	}
	ad.SetServers(adServers)
	var adCredentialsObject ADCredential
	activeDirectoryObject.Credentials.As(ctx, &adCredentialsObject, basetypes.ObjectAsOptions{})
	adCredentials := openApiClient.NewUpdateNetworkWirelessSsidRequestActiveDirectoryCredentials()
	adCredentials.SetPassword(adCredentialsObject.Password.ValueString())
	adCredentials.SetLogonName(adCredentialsObject.LogonName.ValueString())
	ad.SetCredentials(*adCredentials)
	request.SetActiveDirectory(*ad)

	var dnsRewriteObject DNSRewrite
	ssid.DnsRewrite.As(ctx, &dnsRewriteObject, basetypes.ObjectAsOptions{})
	dnsRewrite := openApiClient.NewUpdateNetworkWirelessSsidRequestDnsRewrite()
	dnsRewrite.SetEnabled(dnsRewriteObject.Enabled.ValueBool())
	var dnsServers []string
	for _, server := range dnsRewriteObject.DnsCustomNameservers.Elements() {
		dnsServers = append(dnsServers, server.String())
	}
	dnsRewrite.SetDnsCustomNameservers(dnsServers)
	request.SetDnsRewrite(*dnsRewrite)

	var speedBurstObject SpeedBurst
	ssid.SpeedBurst.As(ctx, &speedBurstObject, basetypes.ObjectAsOptions{})
	speedBurst := openApiClient.NewUpdateNetworkWirelessSsidRequestSpeedBurst()
	speedBurst.SetEnabled(speedBurstObject.Enabled.ValueBool())
	request.SetSpeedBurst(*speedBurst)

	var greObject Gre
	ssid.Gre.As(ctx, &greObject, basetypes.ObjectAsOptions{})
	gre := openApiClient.NewUpdateNetworkWirelessSsidRequestGre()
	gre.SetKey(int32(greObject.Key.ValueInt64()))
	var greConentratorObject Concentrator
	greObject.Concentrator.As(ctx, &greConentratorObject, basetypes.ObjectAsOptions{})
	concentrator := openApiClient.NewUpdateNetworkWirelessSsidRequestGreConcentrator(greConentratorObject.Host.ValueString())
	gre.SetConcentrator(*concentrator)
	request.SetGre(*gre)

	var oauthObject Oauth
	ssid.Oauth.As(ctx, &oauthObject, basetypes.ObjectAsOptions{})
	oauth := openApiClient.NewUpdateNetworkWirelessSsidRequestOauth()
	oauthDomains := []string{}
	for _, domain := range oauthObject.AllowedDomains.Elements() {
		oauthDomains = append(oauthDomains, domain.String())
	}
	oauth.SetAllowedDomains(oauthDomains)
	request.SetOauth(*oauth)

	var sponsorDomains []string
	for _, domain := range ssid.SplashGuestSponsorDomains.Elements() {
		sponsorDomains = append(sponsorDomains, domain.String())
	}
	request.SetSplashGuestSponsorDomains(sponsorDomains)
	var gardenRanges []string
	for _, gardenRange := range ssid.WalledGardenRanges.Elements() {
		gardenRanges = append(gardenRanges, gardenRange.String())
	}
	request.SetWalledGardenRanges(gardenRanges)
	var availabilityTags []string
	for _, availabilityTag := range ssid.AvailabilityTags.Elements() {
		availabilityTags = append(availabilityTags, availabilityTag.String())
	}
	request.SetAvailabilityTags(availabilityTags)

	radiusServers := []openApiClient.UpdateNetworkWirelessSsidRequestRadiusServersInner{}
	for _, radius := range ssid.RadiusServers {
		radiusServer := openApiClient.NewUpdateNetworkWirelessSsidRequestRadiusServersInner(radius.Host.ValueString())
		radiusServer.SetPort(int32(radius.Port.ValueInt64()))
		radiusServer.SetSecret(radius.Secret.ValueString())
		radiusServers = append(radiusServers, *radiusServer)
	}
	request.SetRadiusServers(radiusServers)

	radiusAccountingServers := []openApiClient.UpdateNetworkWirelessSsidRequestRadiusAccountingServersInner{}
	for _, radius := range ssid.RadiusAccountingServers {
		radiusServer := openApiClient.NewUpdateNetworkWirelessSsidRequestRadiusAccountingServersInner(radius.Host.ValueString())
		radiusServer.SetPort(int32(radius.Port.ValueInt64()))
		radiusServer.SetSecret(radius.Secret.ValueString())
		radiusServer.SetCaCertificate(radius.CaCertificate.ValueString())
		radiusServer.SetRadsecEnabled(radius.RadsecEnabled.ValueBool())
		radiusAccountingServers = append(radiusAccountingServers, *radiusServer)
	}
	request.SetRadiusAccountingServers(radiusAccountingServers)

	apTagsAndVlanIds := []openApiClient.UpdateNetworkWirelessSsidRequestApTagsAndVlanIdsInner{}
	for _, vlanID := range ssid.ApTagsAndVlanIds {
		apTagsAndVlanId := openApiClient.NewUpdateNetworkWirelessSsidRequestApTagsAndVlanIdsInner()
		apTagsAndVlanId.SetVlanId(int32(vlanID.VlanId.ValueInt64()))
		tags := []string{}
		for _, tag := range vlanID.Tags.Elements() {
			tags = append(tags, tag.String())
		}
		apTagsAndVlanId.SetTags(tags)
		apTagsAndVlanIds = append(apTagsAndVlanIds, *apTagsAndVlanId)
	}
	request.SetApTagsAndVlanIds(apTagsAndVlanIds)

	_, httpResp, err := r.client.WirelessApi.UpdateNetworkWirelessSsid(context.Background(), data.NetworkID.ValueString(), data.Number.ValueString()).UpdateNetworkWirelessSsidRequest(*request).Execute()
	// Check for API success response code
	if httpResp.StatusCode == 404 {
		resp.Diagnostics.AddWarning(
			"no ssid number or network found in API",
			tools.HttpDiagnostics(httpResp),
		)

	} else if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)

		// HTTP 400 counts as an error so moving this here
		// If there was an error during API call, add it to diagnostics.
		if err != nil {
			resp.Diagnostics.AddError(
				"HTTP Client Failure",
				tools.HttpDiagnostics(httpResp),
			)
			return
		}

		// If there were any errors up to this point, log the state data and return.
		if resp.Diagnostics.HasError() {
			resp.Diagnostics.AddError("State Data", fmt.Sprintf("\n%v", data))
			return
		}

		// Decode the HTTP response body into your data model.
		// If there's an error, add it to diagnostics.
		if err = json.NewDecoder(httpResp.Body).Decode(&data.SSID); err != nil {
			resp.Diagnostics.AddError(
				"JSON decoding error",
				fmt.Sprintf("%v\n", err.Error()),
			)
			return
		}
	}

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data.SSID)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")
}

func (r *NetworksWirelessSsidsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetworksWirelessSsidsResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := r.client.WirelessApi.GetNetworkWirelessSsid(context.Background(), data.NetworkID.ValueString(), data.Number.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// Check for API success response code
	if httpResp.StatusCode == 404 {
		resp.Diagnostics.AddWarning(
			"n found in API",
			tools.HttpDiagnostics(httpResp),
		)

	} else if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)

		// HTTP 400 counts as an error so moving this here
		// If there was an error during API call, add it to diagnostics.
		if err != nil {
			resp.Diagnostics.AddError(
				"HTTP Client Failure",
				tools.HttpDiagnostics(httpResp),
			)
			return
		}

		// If there were any errors up to this point, log the state data and return.
		if resp.Diagnostics.HasError() {
			resp.Diagnostics.AddError("State Data", fmt.Sprintf("\n%v", data))
			return
		}

		// Decode the HTTP response body into your data model.
		// If there's an error, add it to diagnostics.
		if err = json.NewDecoder(httpResp.Body).Decode(&data); err != nil {
			resp.Diagnostics.AddError(
				"JSON decoding error",
				fmt.Sprintf("%v\n", err.Error()),
			)
			return
		}
	}

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data.SSID)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *NetworksWirelessSsidsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *NetworksWirelessSsidsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	request := openApiClient.NewUpdateNetworkWirelessSsidRequest()
	var ssid WirelessNetworkSSID
	data.SSID.As(ctx, &ssid, basetypes.ObjectAsOptions{})

	request.SetName(ssid.Name.ValueString())
	request.SetAuthMode(ssid.AuthMode.ValueString())
	request.SetEnterpriseAdminAccess(ssid.EnterpriseAdminAccess.ValueString())
	request.SetEncryptionMode(ssid.EncryptionMode.ValueString())
	request.SetPsk(ssid.Psk.ValueString())
	request.SetWpaEncryptionMode(ssid.WpaEncryptionMode.ValueString())
	request.SetSplashPage(ssid.SplashPage.ValueString())
	request.SetRadiusCalledStationId(ssid.RadiusCalledStationId.ValueString())
	request.SetRadiusAuthenticationNasId(ssid.RadiusAuthenticationNasId.ValueString())
	request.SetRadiusFailoverPolicy(ssid.RadiusFailoverPolicy.ValueString())
	request.SetRadiusLoadBalancingPolicy(ssid.RadiusLoadBalancingPolicy.ValueString())
	request.SetRadiusAttributeForGroupPolicies(ssid.RadiusAttributeForGroupPolicies.ValueString())
	request.SetIpAssignmentMode(ssid.IpAssignmentMode.ValueString())
	request.SetConcentratorNetworkId(ssid.ConcentratorNetworkId.ValueString())
	request.SetSecondaryConcentratorNetworkId(ssid.SecondaryConcentratorNetworkId.ValueString())
	request.SetBandSelection(ssid.BandSelection.ValueString())
	request.SetRadiusServerTimeout(int32(ssid.RadiusServerTimeout.ValueInt64()))
	request.SetRadiusServerAttemptsLimit(int32(ssid.RadiusServerAttemptsLimit.ValueInt64()))
	request.SetRadiusAccountingInterimInterval(int32(ssid.RadiusAccountingInterimInterval.ValueInt64()))
	request.SetVlanId(int32(ssid.VlanId.ValueInt64()))
	request.SetDefaultVlanId(int32(ssid.DefaultVlanId.ValueInt64()))
	request.SetPerClientBandwidthLimitUp(int32(ssid.PerClientBandwidthLimitUp.ValueInt64()))
	request.SetPerClientBandwidthLimitDown(int32(ssid.PerClientBandwidthLimitDown.ValueInt64()))
	request.SetPerSsidBandwidthLimitUp(int32(ssid.PerSsidBandwidthLimitUp.ValueInt64()))
	request.SetPerSsidBandwidthLimitDown(int32(ssid.PerSsidBandwidthLimitDown.ValueInt64()))
	request.SetRadiusGuestVlanId(int32(ssid.RadiusGuestVlanId.ValueInt64()))
	request.SetMinBitrate(float32(ssid.MinBitrate.ValueFloat64()))
	request.SetUseVlanTagging(ssid.UseVlanTagging.ValueBool())
	request.SetDisassociateClientsOnVpnFailover(ssid.DisassociateClientsOnVpnFailover.ValueBool())
	request.SetRadiusOverride(ssid.RadiusOverride.ValueBool())
	request.SetRadiusGuestVlanEnabled(ssid.RadiusGuestVlanEnabled.ValueBool())
	request.SetEnabled(ssid.Enabled.ValueBool())
	request.SetRadiusProxyEnabled(ssid.RadiusProxyEnabled.ValueBool())
	request.SetRadiusTestingEnabled(ssid.RadiusTestingEnabled.ValueBool())
	request.SetRadiusFallbackEnabled(ssid.RadiusFallbackEnabled.ValueBool())
	request.SetRadiusCoaEnabled(ssid.RadiusCoaEnabled.ValueBool())
	request.SetRadiusAccountingEnabled(ssid.RadiusAccountingEnabled.ValueBool())
	request.SetLanIsolationEnabled(ssid.LanIsolationEnabled.ValueBool())
	request.SetVisible(ssid.Visible.ValueBool())
	request.SetAvailableOnAllAps(ssid.AvailableOnAllAps.ValueBool())
	request.SetMandatoryDhcpEnabled(ssid.MandatoryDhcpEnabled.ValueBool())
	request.SetAdultContentFilteringEnabled(ssid.AdultContentFilteringEnabled.ValueBool())
	request.SetWalledGardenEnabled(ssid.WalledGardenEnabled.ValueBool())

	dot11w := openApiClient.NewUpdateNetworkApplianceSsidRequestDot11w()
	var dot11wObject Dot11W
	ssid.Dot11W.As(ctx, &dot11wObject, basetypes.ObjectAsOptions{})
	dot11w.SetEnabled(dot11wObject.Enabled.ValueBool())
	dot11w.SetRequired(dot11wObject.Required.ValueBool())
	request.SetDot11w(*dot11w)

	dot11r := openApiClient.NewUpdateNetworkWirelessSsidRequestDot11r()
	var dot11rObject Dot11R
	ssid.Dot11R.As(ctx, &dot11rObject, basetypes.ObjectAsOptions{})
	dot11r.SetEnabled(dot11rObject.Enabled.ValueBool())
	dot11r.SetAdaptive(dot11rObject.Adaptive.ValueBool())
	request.SetDot11r(*dot11r)

	localRadius := openApiClient.NewUpdateNetworkWirelessSsidRequestLocalRadius()
	var localRadiusObject LocalRadius
	ssid.LocalRadius.As(ctx, &localRadiusObject, basetypes.ObjectAsOptions{})
	localRadius.SetCacheTimeout(int32(localRadiusObject.CacheTimeout.ValueInt64()))
	certificateAuthentication := openApiClient.NewUpdateNetworkWirelessSsidRequestLocalRadiusCertificateAuthentication()
	var certificateAuthenticationObject CertificateAuthentication
	localRadiusObject.CertificateAuthentication.As(ctx, &certificateAuthenticationObject, basetypes.ObjectAsOptions{})
	certificateAuthentication.SetEnabled(certificateAuthenticationObject.Enabled.ValueBool())
	certificateAuthentication.SetOcspResponderUrl(certificateAuthenticationObject.OcspResponderUrl.ValueString())
	certificateAuthentication.SetUseLdap(certificateAuthenticationObject.UseLdap.ValueBool())
	certificateAuthentication.SetUseOcsp(certificateAuthenticationObject.UseOcsp.ValueBool())
	var clientRootCaCertificateObject ClientRootCaCertificate
	certificateAuthenticationObject.ClientRootCaCertificate.As(ctx, &clientRootCaCertificateObject, basetypes.ObjectAsOptions{})
	rootCACertificate := openApiClient.NewUpdateNetworkWirelessSsidRequestLocalRadiusCertificateAuthenticationClientRootCaCertificate()
	rootCACertificate.SetContents(clientRootCaCertificateObject.Contents.ValueString())
	certificateAuthentication.SetClientRootCaCertificate(*rootCACertificate)
	localRadius.SetCertificateAuthentication(*certificateAuthentication)
	request.SetLocalRadius(*localRadius)

	ldap := openApiClient.NewUpdateNetworkWirelessSsidRequestLdap()
	var ldapObject Ldap
	ssid.Ldap.As(ctx, &ldapObject, basetypes.ObjectAsOptions{})
	ldap.SetBaseDistinguishedName(ldapObject.BaseDistinguishedName.ValueString())
	ldapCredentials := openApiClient.NewUpdateNetworkWirelessSsidRequestLdapCredentials()
	var ldapCredentialsObject Credential
	ldapObject.Credentials.As(ctx, &ldapCredentialsObject, basetypes.ObjectAsOptions{})
	ldapCredentials.SetDistinguishedName(ldapCredentialsObject.DistinguishedName.ValueString())
	ldapCredentials.SetPassword(ldapCredentialsObject.Password.ValueString())
	ldap.SetCredentials(*ldapCredentials)
	var ldapServers []openApiClient.UpdateNetworkWirelessSsidRequestLdapServersInner
	for _, server := range ldapObject.Servers {
		inner := openApiClient.NewUpdateNetworkWirelessSsidRequestLdapServersInner(server.Host.ValueString(), int32(server.Port.ValueInt64()))
		ldapServers = append(ldapServers, *inner)
	}
	ldap.SetServers(ldapServers)
	var ldapServerCaCertificateObject ServerCaCertificate
	ldapObject.ServerCaCertificate.As(ctx, &ldapServerCaCertificateObject, basetypes.ObjectAsOptions{})
	ldapServerCaCertificate := openApiClient.NewUpdateNetworkWirelessSsidRequestLdapServerCaCertificate()
	ldapServerCaCertificate.SetContents(ldapServerCaCertificateObject.Contents.ValueString())
	ldap.SetServerCaCertificate(*ldapServerCaCertificate)
	request.SetLdap(*ldap)

	ad := openApiClient.NewUpdateNetworkWirelessSsidRequestActiveDirectory()
	var activeDirectoryObject ActiveDirectory
	ssid.ActiveDirectory.As(ctx, &activeDirectoryObject, basetypes.ObjectAsOptions{})
	var adServers []openApiClient.UpdateNetworkWirelessSsidRequestActiveDirectoryServersInner
	for _, server := range activeDirectoryObject.Servers {
		inner := openApiClient.NewUpdateNetworkWirelessSsidRequestActiveDirectoryServersInner(server.Host.ValueString())
		inner.SetPort(int32(server.Port.ValueInt64()))
		adServers = append(adServers, *inner)
	}
	ad.SetServers(adServers)
	var adCredentialsObject ADCredential
	activeDirectoryObject.Credentials.As(ctx, &adCredentialsObject, basetypes.ObjectAsOptions{})
	adCredentials := openApiClient.NewUpdateNetworkWirelessSsidRequestActiveDirectoryCredentials()
	adCredentials.SetPassword(adCredentialsObject.Password.ValueString())
	adCredentials.SetLogonName(adCredentialsObject.LogonName.ValueString())
	ad.SetCredentials(*adCredentials)
	request.SetActiveDirectory(*ad)

	var dnsRewriteObject DNSRewrite
	ssid.DnsRewrite.As(ctx, &dnsRewriteObject, basetypes.ObjectAsOptions{})
	dnsRewrite := openApiClient.NewUpdateNetworkWirelessSsidRequestDnsRewrite()
	dnsRewrite.SetEnabled(dnsRewriteObject.Enabled.ValueBool())
	var dnsServers []string
	for _, server := range dnsRewriteObject.DnsCustomNameservers.Elements() {
		dnsServers = append(dnsServers, server.String())
	}
	dnsRewrite.SetDnsCustomNameservers(dnsServers)
	request.SetDnsRewrite(*dnsRewrite)

	var speedBurstObject SpeedBurst
	ssid.SpeedBurst.As(ctx, &speedBurstObject, basetypes.ObjectAsOptions{})
	speedBurst := openApiClient.NewUpdateNetworkWirelessSsidRequestSpeedBurst()
	speedBurst.SetEnabled(speedBurstObject.Enabled.ValueBool())
	request.SetSpeedBurst(*speedBurst)

	var greObject Gre
	ssid.Gre.As(ctx, &greObject, basetypes.ObjectAsOptions{})
	gre := openApiClient.NewUpdateNetworkWirelessSsidRequestGre()
	gre.SetKey(int32(greObject.Key.ValueInt64()))
	var greConentratorObject Concentrator
	greObject.Concentrator.As(ctx, &greConentratorObject, basetypes.ObjectAsOptions{})
	concentrator := openApiClient.NewUpdateNetworkWirelessSsidRequestGreConcentrator(greConentratorObject.Host.ValueString())
	gre.SetConcentrator(*concentrator)
	request.SetGre(*gre)

	var oauthObject Oauth
	ssid.Oauth.As(ctx, &oauthObject, basetypes.ObjectAsOptions{})
	oauth := openApiClient.NewUpdateNetworkWirelessSsidRequestOauth()
	oauthDomains := []string{}
	for _, domain := range oauthObject.AllowedDomains.Elements() {
		oauthDomains = append(oauthDomains, domain.String())
	}
	oauth.SetAllowedDomains(oauthDomains)
	request.SetOauth(*oauth)

	var sponsorDomains []string
	for _, domain := range ssid.SplashGuestSponsorDomains.Elements() {
		sponsorDomains = append(sponsorDomains, domain.String())
	}
	request.SetSplashGuestSponsorDomains(sponsorDomains)
	var gardenRanges []string
	for _, gardenRange := range ssid.WalledGardenRanges.Elements() {
		gardenRanges = append(gardenRanges, gardenRange.String())
	}
	request.SetWalledGardenRanges(gardenRanges)
	var availabilityTags []string
	for _, availabilityTag := range ssid.AvailabilityTags.Elements() {
		availabilityTags = append(availabilityTags, availabilityTag.String())
	}
	request.SetAvailabilityTags(availabilityTags)

	radiusServers := []openApiClient.UpdateNetworkWirelessSsidRequestRadiusServersInner{}
	for _, radius := range ssid.RadiusServers {
		radiusServer := openApiClient.NewUpdateNetworkWirelessSsidRequestRadiusServersInner(radius.Host.ValueString())
		radiusServer.SetPort(int32(radius.Port.ValueInt64()))
		radiusServer.SetSecret(radius.Secret.ValueString())
		radiusServers = append(radiusServers, *radiusServer)
	}
	request.SetRadiusServers(radiusServers)

	radiusAccountingServers := []openApiClient.UpdateNetworkWirelessSsidRequestRadiusAccountingServersInner{}
	for _, radius := range ssid.RadiusAccountingServers {
		radiusServer := openApiClient.NewUpdateNetworkWirelessSsidRequestRadiusAccountingServersInner(radius.Host.ValueString())
		radiusServer.SetPort(int32(radius.Port.ValueInt64()))
		radiusServer.SetSecret(radius.Secret.ValueString())
		radiusServer.SetCaCertificate(radius.CaCertificate.ValueString())
		radiusServer.SetRadsecEnabled(radius.RadsecEnabled.ValueBool())
		radiusAccountingServers = append(radiusAccountingServers, *radiusServer)
	}
	request.SetRadiusAccountingServers(radiusAccountingServers)

	apTagsAndVlanIds := []openApiClient.UpdateNetworkWirelessSsidRequestApTagsAndVlanIdsInner{}
	for _, vlanID := range ssid.ApTagsAndVlanIds {
		apTagsAndVlanId := openApiClient.NewUpdateNetworkWirelessSsidRequestApTagsAndVlanIdsInner()
		apTagsAndVlanId.SetVlanId(int32(vlanID.VlanId.ValueInt64()))
		tags := []string{}
		for _, tag := range vlanID.Tags.Elements() {
			tags = append(tags, tag.String())
		}
		apTagsAndVlanId.SetTags(tags)
		apTagsAndVlanIds = append(apTagsAndVlanIds, *apTagsAndVlanId)
	}
	request.SetApTagsAndVlanIds(apTagsAndVlanIds)

	_, httpResp, err := r.client.WirelessApi.UpdateNetworkWirelessSsid(context.Background(), data.NetworkID.ValueString(), data.Number.ValueString()).UpdateNetworkWirelessSsidRequest(*request).Execute()
	// Check for API success response code
	if httpResp.StatusCode == 404 {
		resp.Diagnostics.AddWarning(
			"no ssid number or network found in API",
			tools.HttpDiagnostics(httpResp),
		)

	} else if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)

		// HTTP 400 counts as an error so moving this here
		// If there was an error during API call, add it to diagnostics.
		if err != nil {
			resp.Diagnostics.AddError(
				"HTTP Client Failure",
				tools.HttpDiagnostics(httpResp),
			)
			return
		}

		// If there were any errors up to this point, log the state data and return.
		if resp.Diagnostics.HasError() {
			resp.Diagnostics.AddError("State Data", fmt.Sprintf("\n%v", data))
			return
		}

		// Decode the HTTP response body into your data model.
		// If there's an error, add it to diagnostics.
		if err = json.NewDecoder(httpResp.Body).Decode(&data.SSID); err != nil {
			resp.Diagnostics.AddError(
				"JSON decoding error",
				fmt.Sprintf("%v\n", err.Error()),
			)
			return
		}
	}

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data.SSID)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *NetworksWirelessSsidsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *NetworksWirelessSsidsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	request := openApiClient.NewUpdateNetworkWirelessSsidRequest()
	var ssid WirelessNetworkSSID
	data.SSID.As(ctx, &ssid, basetypes.ObjectAsOptions{})

	request.SetName(ssid.Name.ValueString())
	request.SetAuthMode(ssid.AuthMode.ValueString())
	request.SetEnterpriseAdminAccess(ssid.EnterpriseAdminAccess.ValueString())
	request.SetEncryptionMode(ssid.EncryptionMode.ValueString())
	request.SetPsk(ssid.Psk.ValueString())
	request.SetWpaEncryptionMode(ssid.WpaEncryptionMode.ValueString())
	request.SetSplashPage(ssid.SplashPage.ValueString())
	request.SetRadiusCalledStationId(ssid.RadiusCalledStationId.ValueString())
	request.SetRadiusAuthenticationNasId(ssid.RadiusAuthenticationNasId.ValueString())
	request.SetRadiusFailoverPolicy(ssid.RadiusFailoverPolicy.ValueString())
	request.SetRadiusLoadBalancingPolicy(ssid.RadiusLoadBalancingPolicy.ValueString())
	request.SetRadiusAttributeForGroupPolicies(ssid.RadiusAttributeForGroupPolicies.ValueString())
	request.SetIpAssignmentMode(ssid.IpAssignmentMode.ValueString())
	request.SetConcentratorNetworkId(ssid.ConcentratorNetworkId.ValueString())
	request.SetSecondaryConcentratorNetworkId(ssid.SecondaryConcentratorNetworkId.ValueString())
	request.SetBandSelection(ssid.BandSelection.ValueString())
	request.SetRadiusServerTimeout(int32(ssid.RadiusServerTimeout.ValueInt64()))
	request.SetRadiusServerAttemptsLimit(int32(ssid.RadiusServerAttemptsLimit.ValueInt64()))
	request.SetRadiusAccountingInterimInterval(int32(ssid.RadiusAccountingInterimInterval.ValueInt64()))
	request.SetVlanId(int32(ssid.VlanId.ValueInt64()))
	request.SetDefaultVlanId(int32(ssid.DefaultVlanId.ValueInt64()))
	request.SetPerClientBandwidthLimitUp(int32(ssid.PerClientBandwidthLimitUp.ValueInt64()))
	request.SetPerClientBandwidthLimitDown(int32(ssid.PerClientBandwidthLimitDown.ValueInt64()))
	request.SetPerSsidBandwidthLimitUp(int32(ssid.PerSsidBandwidthLimitUp.ValueInt64()))
	request.SetPerSsidBandwidthLimitDown(int32(ssid.PerSsidBandwidthLimitDown.ValueInt64()))
	request.SetRadiusGuestVlanId(int32(ssid.RadiusGuestVlanId.ValueInt64()))
	request.SetMinBitrate(float32(ssid.MinBitrate.ValueFloat64()))
	request.SetUseVlanTagging(ssid.UseVlanTagging.ValueBool())
	request.SetDisassociateClientsOnVpnFailover(ssid.DisassociateClientsOnVpnFailover.ValueBool())
	request.SetRadiusOverride(ssid.RadiusOverride.ValueBool())
	request.SetRadiusGuestVlanEnabled(ssid.RadiusGuestVlanEnabled.ValueBool())
	request.SetEnabled(ssid.Enabled.ValueBool())
	request.SetRadiusProxyEnabled(ssid.RadiusProxyEnabled.ValueBool())
	request.SetRadiusTestingEnabled(ssid.RadiusTestingEnabled.ValueBool())
	request.SetRadiusFallbackEnabled(ssid.RadiusFallbackEnabled.ValueBool())
	request.SetRadiusCoaEnabled(ssid.RadiusCoaEnabled.ValueBool())
	request.SetRadiusAccountingEnabled(ssid.RadiusAccountingEnabled.ValueBool())
	request.SetLanIsolationEnabled(ssid.LanIsolationEnabled.ValueBool())
	request.SetVisible(ssid.Visible.ValueBool())
	request.SetAvailableOnAllAps(ssid.AvailableOnAllAps.ValueBool())
	request.SetMandatoryDhcpEnabled(ssid.MandatoryDhcpEnabled.ValueBool())
	request.SetAdultContentFilteringEnabled(ssid.AdultContentFilteringEnabled.ValueBool())
	request.SetWalledGardenEnabled(ssid.WalledGardenEnabled.ValueBool())

	dot11w := openApiClient.NewUpdateNetworkApplianceSsidRequestDot11w()
	var dot11wObject Dot11W
	ssid.Dot11W.As(ctx, &dot11wObject, basetypes.ObjectAsOptions{})
	dot11w.SetEnabled(dot11wObject.Enabled.ValueBool())
	dot11w.SetRequired(dot11wObject.Required.ValueBool())
	request.SetDot11w(*dot11w)

	dot11r := openApiClient.NewUpdateNetworkWirelessSsidRequestDot11r()
	var dot11rObject Dot11R
	ssid.Dot11R.As(ctx, &dot11rObject, basetypes.ObjectAsOptions{})
	dot11r.SetEnabled(dot11rObject.Enabled.ValueBool())
	dot11r.SetAdaptive(dot11rObject.Adaptive.ValueBool())
	request.SetDot11r(*dot11r)

	localRadius := openApiClient.NewUpdateNetworkWirelessSsidRequestLocalRadius()
	var localRadiusObject LocalRadius
	ssid.LocalRadius.As(ctx, &localRadiusObject, basetypes.ObjectAsOptions{})
	localRadius.SetCacheTimeout(int32(localRadiusObject.CacheTimeout.ValueInt64()))
	certificateAuthentication := openApiClient.NewUpdateNetworkWirelessSsidRequestLocalRadiusCertificateAuthentication()
	var certificateAuthenticationObject CertificateAuthentication
	localRadiusObject.CertificateAuthentication.As(ctx, &certificateAuthenticationObject, basetypes.ObjectAsOptions{})
	certificateAuthentication.SetEnabled(certificateAuthenticationObject.Enabled.ValueBool())
	certificateAuthentication.SetOcspResponderUrl(certificateAuthenticationObject.OcspResponderUrl.ValueString())
	certificateAuthentication.SetUseLdap(certificateAuthenticationObject.UseLdap.ValueBool())
	certificateAuthentication.SetUseOcsp(certificateAuthenticationObject.UseOcsp.ValueBool())
	var clientRootCaCertificateObject ClientRootCaCertificate
	certificateAuthenticationObject.ClientRootCaCertificate.As(ctx, &clientRootCaCertificateObject, basetypes.ObjectAsOptions{})
	rootCACertificate := openApiClient.NewUpdateNetworkWirelessSsidRequestLocalRadiusCertificateAuthenticationClientRootCaCertificate()
	rootCACertificate.SetContents(clientRootCaCertificateObject.Contents.ValueString())
	certificateAuthentication.SetClientRootCaCertificate(*rootCACertificate)
	localRadius.SetCertificateAuthentication(*certificateAuthentication)
	request.SetLocalRadius(*localRadius)

	ldap := openApiClient.NewUpdateNetworkWirelessSsidRequestLdap()
	var ldapObject Ldap
	ssid.Ldap.As(ctx, &ldapObject, basetypes.ObjectAsOptions{})
	ldap.SetBaseDistinguishedName(ldapObject.BaseDistinguishedName.ValueString())
	ldapCredentials := openApiClient.NewUpdateNetworkWirelessSsidRequestLdapCredentials()
	var ldapCredentialsObject Credential
	ldapObject.Credentials.As(ctx, &ldapCredentialsObject, basetypes.ObjectAsOptions{})
	ldapCredentials.SetDistinguishedName(ldapCredentialsObject.DistinguishedName.ValueString())
	ldapCredentials.SetPassword(ldapCredentialsObject.Password.ValueString())
	ldap.SetCredentials(*ldapCredentials)
	var ldapServers []openApiClient.UpdateNetworkWirelessSsidRequestLdapServersInner
	for _, server := range ldapObject.Servers {
		inner := openApiClient.NewUpdateNetworkWirelessSsidRequestLdapServersInner(server.Host.ValueString(), int32(server.Port.ValueInt64()))
		ldapServers = append(ldapServers, *inner)
	}
	ldap.SetServers(ldapServers)
	var ldapServerCaCertificateObject ServerCaCertificate
	ldapObject.ServerCaCertificate.As(ctx, &ldapServerCaCertificateObject, basetypes.ObjectAsOptions{})
	ldapServerCaCertificate := openApiClient.NewUpdateNetworkWirelessSsidRequestLdapServerCaCertificate()
	ldapServerCaCertificate.SetContents(ldapServerCaCertificateObject.Contents.ValueString())
	ldap.SetServerCaCertificate(*ldapServerCaCertificate)
	request.SetLdap(*ldap)

	ad := openApiClient.NewUpdateNetworkWirelessSsidRequestActiveDirectory()
	var activeDirectoryObject ActiveDirectory
	ssid.ActiveDirectory.As(ctx, &activeDirectoryObject, basetypes.ObjectAsOptions{})
	var adServers []openApiClient.UpdateNetworkWirelessSsidRequestActiveDirectoryServersInner
	for _, server := range activeDirectoryObject.Servers {
		inner := openApiClient.NewUpdateNetworkWirelessSsidRequestActiveDirectoryServersInner(server.Host.ValueString())
		inner.SetPort(int32(server.Port.ValueInt64()))
		adServers = append(adServers, *inner)
	}
	ad.SetServers(adServers)
	var adCredentialsObject ADCredential
	activeDirectoryObject.Credentials.As(ctx, &adCredentialsObject, basetypes.ObjectAsOptions{})
	adCredentials := openApiClient.NewUpdateNetworkWirelessSsidRequestActiveDirectoryCredentials()
	adCredentials.SetPassword(adCredentialsObject.Password.ValueString())
	adCredentials.SetLogonName(adCredentialsObject.LogonName.ValueString())
	ad.SetCredentials(*adCredentials)
	request.SetActiveDirectory(*ad)

	var dnsRewriteObject DNSRewrite
	ssid.DnsRewrite.As(ctx, &dnsRewriteObject, basetypes.ObjectAsOptions{})
	dnsRewrite := openApiClient.NewUpdateNetworkWirelessSsidRequestDnsRewrite()
	dnsRewrite.SetEnabled(dnsRewriteObject.Enabled.ValueBool())
	var dnsServers []string
	for _, server := range dnsRewriteObject.DnsCustomNameservers.Elements() {
		dnsServers = append(dnsServers, server.String())
	}
	dnsRewrite.SetDnsCustomNameservers(dnsServers)
	request.SetDnsRewrite(*dnsRewrite)

	var speedBurstObject SpeedBurst
	ssid.SpeedBurst.As(ctx, &speedBurstObject, basetypes.ObjectAsOptions{})
	speedBurst := openApiClient.NewUpdateNetworkWirelessSsidRequestSpeedBurst()
	speedBurst.SetEnabled(speedBurstObject.Enabled.ValueBool())
	request.SetSpeedBurst(*speedBurst)

	var greObject Gre
	ssid.Gre.As(ctx, &greObject, basetypes.ObjectAsOptions{})
	gre := openApiClient.NewUpdateNetworkWirelessSsidRequestGre()
	gre.SetKey(int32(greObject.Key.ValueInt64()))
	var greConentratorObject Concentrator
	greObject.Concentrator.As(ctx, &greConentratorObject, basetypes.ObjectAsOptions{})
	concentrator := openApiClient.NewUpdateNetworkWirelessSsidRequestGreConcentrator(greConentratorObject.Host.ValueString())
	gre.SetConcentrator(*concentrator)
	request.SetGre(*gre)

	var oauthObject Oauth
	ssid.Oauth.As(ctx, &oauthObject, basetypes.ObjectAsOptions{})
	oauth := openApiClient.NewUpdateNetworkWirelessSsidRequestOauth()
	oauthDomains := []string{}
	for _, domain := range oauthObject.AllowedDomains.Elements() {
		oauthDomains = append(oauthDomains, domain.String())
	}
	oauth.SetAllowedDomains(oauthDomains)
	request.SetOauth(*oauth)

	var sponsorDomains []string
	for _, domain := range ssid.SplashGuestSponsorDomains.Elements() {
		sponsorDomains = append(sponsorDomains, domain.String())
	}
	request.SetSplashGuestSponsorDomains(sponsorDomains)
	var gardenRanges []string
	for _, gardenRange := range ssid.WalledGardenRanges.Elements() {
		gardenRanges = append(gardenRanges, gardenRange.String())
	}
	request.SetWalledGardenRanges(gardenRanges)
	var availabilityTags []string
	for _, availabilityTag := range ssid.AvailabilityTags.Elements() {
		availabilityTags = append(availabilityTags, availabilityTag.String())
	}
	request.SetAvailabilityTags(availabilityTags)

	radiusServers := []openApiClient.UpdateNetworkWirelessSsidRequestRadiusServersInner{}
	for _, radius := range ssid.RadiusServers {
		radiusServer := openApiClient.NewUpdateNetworkWirelessSsidRequestRadiusServersInner(radius.Host.ValueString())
		radiusServer.SetPort(int32(radius.Port.ValueInt64()))
		radiusServer.SetSecret(radius.Secret.ValueString())
		radiusServers = append(radiusServers, *radiusServer)
	}
	request.SetRadiusServers(radiusServers)

	radiusAccountingServers := []openApiClient.UpdateNetworkWirelessSsidRequestRadiusAccountingServersInner{}
	for _, radius := range ssid.RadiusAccountingServers {
		radiusServer := openApiClient.NewUpdateNetworkWirelessSsidRequestRadiusAccountingServersInner(radius.Host.ValueString())
		radiusServer.SetPort(int32(radius.Port.ValueInt64()))
		radiusServer.SetSecret(radius.Secret.ValueString())
		radiusServer.SetCaCertificate(radius.CaCertificate.ValueString())
		radiusServer.SetRadsecEnabled(radius.RadsecEnabled.ValueBool())
		radiusAccountingServers = append(radiusAccountingServers, *radiusServer)
	}
	request.SetRadiusAccountingServers(radiusAccountingServers)

	apTagsAndVlanIds := []openApiClient.UpdateNetworkWirelessSsidRequestApTagsAndVlanIdsInner{}
	for _, vlanID := range ssid.ApTagsAndVlanIds {
		apTagsAndVlanId := openApiClient.NewUpdateNetworkWirelessSsidRequestApTagsAndVlanIdsInner()
		apTagsAndVlanId.SetVlanId(int32(vlanID.VlanId.ValueInt64()))
		tags := []string{}
		for _, tag := range vlanID.Tags.Elements() {
			tags = append(tags, tag.String())
		}
		apTagsAndVlanId.SetTags(tags)
		apTagsAndVlanIds = append(apTagsAndVlanIds, *apTagsAndVlanId)
	}
	request.SetApTagsAndVlanIds(apTagsAndVlanIds)

	_, httpResp, err := r.client.WirelessApi.UpdateNetworkWirelessSsid(context.Background(), data.NetworkID.ValueString(), data.Number.ValueString()).UpdateNetworkWirelessSsidRequest(*request).Execute()
	// Check for API success response code
	if httpResp.StatusCode == 404 {
		resp.Diagnostics.AddWarning(
			"no ssid number or network found in API",
			tools.HttpDiagnostics(httpResp),
		)

	} else if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)

		// HTTP 400 counts as an error so moving this here
		// If there was an error during API call, add it to diagnostics.
		if err != nil {
			resp.Diagnostics.AddError(
				"HTTP Client Failure",
				tools.HttpDiagnostics(httpResp),
			)
			return
		}

		// If there were any errors up to this point, log the state data and return.
		if resp.Diagnostics.HasError() {
			resp.Diagnostics.AddError("State Data", fmt.Sprintf("\n%v", data))
			return
		}

		// Decode the HTTP response body into your data model.
		// If there's an error, add it to diagnostics.
		if err = json.NewDecoder(httpResp.Body).Decode(&data.SSID); err != nil {
			resp.Diagnostics.AddError(
				"JSON decoding error",
				fmt.Sprintf("%v\n", err.Error()),
			)
			return
		}
	}

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data.SSID)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")
}

func (r *NetworksWirelessSsidsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: network_id, number. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("number"), idParts[1])...)

	if resp.Diagnostics.HasError() {
		return
	}
}
