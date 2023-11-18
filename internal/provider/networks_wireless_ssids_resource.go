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

type NetworksWirelessSsidsResourceModelDot11W struct {
	Enabled  jsontypes.Bool `tfsdk:"enabled" json:"enabled"`
	Required jsontypes.Bool `tfsdk:"required" json:"required"`
}

type NetworksWirelessSsidsResourceModelDot11R struct {
	Enabled  jsontypes.Bool `tfsdk:"enabled" json:"enabled"`
	Adaptive jsontypes.Bool `tfsdk:"adaptive" json:"adaptive"`
}

type NetworksWirelessSsidsResourceModelLocalRadiusCertificateAuthenticationClientRootCaCertificate struct {
	Contents jsontypes.String `tfsdk:"contents" json:"contents"`
}

type NetworksWirelessSsidsResourceModelLocalRadiusCertificateAuthentication struct {
	Enabled                 jsontypes.Bool   `tfsdk:"enabled" json:"enabled"`
	UseLdap                 jsontypes.Bool   `tfsdk:"use_ldap" json:"useLdap"`
	UseOcsp                 jsontypes.Bool   `tfsdk:"use_ocsp" json:"useOcsp"`
	OcspResponderUrl        jsontypes.String `tfsdk:"ocsp_responder_url" json:"ocspResponderUrl"`
	ClientRootCaCertificate types.Object     `tfsdk:"client_root_ca_certificate" json:"clientRootCaCertificate"`
}

type NetworksWirelessSsidsResourceModelLocalRadiusPasswordAuthentication struct {
	Enabled jsontypes.Bool `tfsdk:"enabled" json:"enabled"`
}

type NetworksWirelessSsidsResourceModelLocalRadius struct {
	CacheTimeout              jsontypes.Int64 `tfsdk:"cache_timeout" json:"cacheTimeout"`
	PasswordAuthentication    types.Object    `tfsdk:"password_authentication" json:"passwordAuthentication"`
	CertificateAuthentication types.Object    `tfsdk:"certificate_authentication" json:"certificateAuthentication"`
}

type NetworksWirelessSsidsResourceModelServer struct {
	Host jsontypes.String `tfsdk:"host" json:"host"`
	Port jsontypes.Int64  `tfsdk:"port" json:"port"`
}

type NetworksWirelessSsidsResourceModelCredential struct {
	DistinguishedName jsontypes.String `tfsdk:"distinguished_name" json:"distinguishedName"`
	Password          jsontypes.String `tfsdk:"password" json:"password"`
}

type NetworksWirelessSsidsResourceModelServerCaCertificate struct {
	Contents jsontypes.String `tfsdk:"contents" json:"contents"`
}

type NetworksWirelessSsidsResourceModelLdap struct {
	Servers               []NetworksWirelessSsidsResourceModelServer `tfsdk:"servers" json:"servers"`
	Credentials           types.Object                               `tfsdk:"credentials" json:"credentials"`
	BaseDistinguishedName jsontypes.String                           `tfsdk:"base_distinguished_name" json:"baseDistinguishedName"`
	ServerCaCertificate   types.Object                               `tfsdk:"server_ca_certificate" json:"serverCaCertificate"`
}

type NetworksWirelessSsidsResourceModelActiveDirectoryCredential struct {
	LogonName jsontypes.String `tfsdk:"logon_name" json:"logonName"`
	Password  jsontypes.String `tfsdk:"password" json:"password"`
}

type NetworksWirelessSsidsResourceModelActiveDirectory struct {
	Servers     []NetworksWirelessSsidsResourceModelServer `tfsdk:"servers" json:"servers"`
	Credentials types.Object                               `tfsdk:"credentials" json:"credentials"`
}

type NetworksWirelessSsidsResourceModelRadiusServer struct {
	Host                     jsontypes.String `tfsdk:"host" json:"host"`
	Secret                   jsontypes.String `tfsdk:"secret" json:"secret"`
	CaCertificate            jsontypes.String `tfsdk:"ca_certificate" json:"caCertificate"`
	Port                     jsontypes.Int64  `tfsdk:"port" json:"port"`
	OpenRoamingCertificateId jsontypes.Int64  `tfsdk:"open_roaming_certificate_id" json:"openRoamingCertificateId"`
	RadsecEnabled            jsontypes.Bool   `tfsdk:"radsec_enabled" json:"radsecEnabled"`
}

type NetworksWirelessSsidsResourceModelRadiusAccountingServer struct {
	Host          jsontypes.String `tfsdk:"host" json:"host"`
	Secret        jsontypes.String `tfsdk:"secret" json:"secret"`
	CaCertificate jsontypes.String `tfsdk:"ca_certificate" json:"caCertificate"`
	Port          jsontypes.Int64  `tfsdk:"port" json:"port"`
	RadsecEnabled jsontypes.Bool   `tfsdk:"radsec_enabled" json:"radsecEnabled"`
}

type NetworksWirelessSsidsResourceModelConcentrator struct {
	Host jsontypes.String `tfsdk:"host" json:"host"`
}

type NetworksWirelessSsidsResourceModelGre struct {
	Concentrator types.Object    `tfsdk:"concentrator" json:"concentrator"`
	Key          jsontypes.Int64 `tfsdk:"key" json:"key"`
}

type NetworksWirelessSsidsResourceModelApTagsAndVlanId struct {
	Tags   types.List      `tfsdk:"tags" json:"tags"`
	VlanId jsontypes.Int64 `tfsdk:"vlan_id" json:"vlanId"`
}

type NetworksWirelessSsidsResourceModelDNSRewrite struct {
	Enabled              jsontypes.Bool `tdsdk:"enabled" json:"enabled"`
	DnsCustomNameservers types.List     `tdsdk:"dns_custom_nameservers" json:"dnsCustomNameservers"`
}

type NetworksWirelessSsidsResourceModelRadius struct {
	GuestVlan types.Object `tfsdk:"guest_vlan" json:"guestVlan"`
}

type NetworksWirelessSsidsResourceModelSpeedBurst struct {
	Enabled jsontypes.Bool `tfsdk:"enabled" json:"enabled"`
}

type NetworksWirelessSsidsResourceModelOauth struct {
	AllowedDomains types.List `json:"allowed_domains"`
}

type WirelessNetworkSSID struct {
	Name                             jsontypes.String                                           `tfsdk:"name"`
	AuthMode                         jsontypes.String                                           `tfsdk:"auth_mode"`
	EnterpriseAdminAccess            jsontypes.String                                           `tfsdk:"enterprise_admin_access"`
	EncryptionMode                   jsontypes.String                                           `tfsdk:"encryption_mode"`
	Psk                              jsontypes.String                                           `tfsdk:"psk"`
	WpaEncryptionMode                jsontypes.String                                           `tfsdk:"wpa_encryption_mode"`
	SplashPage                       jsontypes.String                                           `tfsdk:"splash_page"`
	RadiusCalledStationId            jsontypes.String                                           `tfsdk:"radius_called_station_id"`
	RadiusAuthenticationNasId        jsontypes.String                                           `tfsdk:"radius_authentication_nas_id"`
	RadiusFailoverPolicy             jsontypes.String                                           `tfsdk:"radius_failover_policy"`
	RadiusLoadBalancingPolicy        jsontypes.String                                           `tfsdk:"radius_load_balancing_policy"`
	RadiusAttributeForGroupPolicies  jsontypes.String                                           `tfsdk:"radius_attribute_for_group_policies"`
	IpAssignmentMode                 jsontypes.String                                           `tfsdk:"ip_assignment_mode"`
	ConcentratorNetworkId            jsontypes.String                                           `tfsdk:"concentrator_network_id"`
	SecondaryConcentratorNetworkId   jsontypes.String                                           `tfsdk:"secondary_concentrator_network_id"`
	BandSelection                    jsontypes.String                                           `tfsdk:"band_selection"`
	RadiusServerTimeout              jsontypes.Int64                                            `tfsdk:"radius_server_timeout"`
	RadiusServerAttemptsLimit        jsontypes.Int64                                            `tfsdk:"radius_server_attempts_limit"`
	RadiusAccountingInterimInterval  jsontypes.Int64                                            `tfsdk:"radius_accounting_interim_interval"`
	VlanId                           jsontypes.Int64                                            `tfsdk:"vlan_id"`
	DefaultVlanId                    jsontypes.Int64                                            `tfsdk:"default_vlan_id"`
	PerClientBandwidthLimitUp        jsontypes.Int64                                            `tfsdk:"per_client_bandwidth_limit_up"`
	PerClientBandwidthLimitDown      jsontypes.Int64                                            `tfsdk:"per_client_bandwidth_limit_down"`
	PerSsidBandwidthLimitUp          jsontypes.Int64                                            `tfsdk:"per_ssid_bandwidth_limit_up"`
	PerSsidBandwidthLimitDown        jsontypes.Int64                                            `tfsdk:"per_ssid_bandwidth_limit_down"`
	RadiusGuestVlanId                jsontypes.Int64                                            `tfsdk:"radius_guest_vlan_id"`
	MinBitrate                       jsontypes.Float64                                          `tfsdk:"min_bitrate"`
	UseVlanTagging                   jsontypes.Bool                                             `tfsdk:"use_vlan_tagging"`
	DisassociateClientsOnVpnFailover jsontypes.Bool                                             `tfsdk:"disassociate_clients_on_vpn_failover"`
	RadiusOverride                   jsontypes.Bool                                             `tfsdk:"radius_override"`
	RadiusGuestVlanEnabled           jsontypes.Bool                                             `tfsdk:"radius_guest_vlan_enabled"`
	Enabled                          jsontypes.Bool                                             `tfsdk:"enabled"`
	RadiusProxyEnabled               jsontypes.Bool                                             `tfsdk:"radius_proxy_enabled"`
	RadiusTestingEnabled             jsontypes.Bool                                             `tfsdk:"radius_testing_enabled"`
	RadiusFallbackEnabled            jsontypes.Bool                                             `tfsdk:"radius_fallback_enabled"`
	RadiusCoaEnabled                 jsontypes.Bool                                             `tfsdk:"radius_coa_enabled"`
	RadiusAccountingEnabled          jsontypes.Bool                                             `tfsdk:"radius_accounting_enabled"`
	LanIsolationEnabled              jsontypes.Bool                                             `tfsdk:"lan_isolation_enabled"`
	Visible                          jsontypes.Bool                                             `tfsdk:"visible"`
	AvailableOnAllAps                jsontypes.Bool                                             `tfsdk:"available_on_all_aps"`
	MandatoryDhcpEnabled             jsontypes.Bool                                             `tfsdk:"mandatory_dhcp_enabled"`
	AdultContentFilteringEnabled     jsontypes.Bool                                             `tfsdk:"adult_content_filtering_enabled"`
	WalledGardenEnabled              jsontypes.Bool                                             `tfsdk:"walled_garden_enabled"`
	Dot11W                           types.Object                                               `tfsdk:"dot11w"`
	Dot11R                           types.Object                                               `tfsdk:"dot11r"`
	LocalRadius                      types.Object                                               `tfsdk:"local_radius"`
	Ldap                             types.Object                                               `tfsdk:"ldap"`
	ActiveDirectory                  types.Object                                               `tfsdk:"active_directory"`
	DnsRewrite                       types.Object                                               `tfsdk:"dns_rewrite"`
	SpeedBurst                       types.Object                                               `tfsdk:"speed_burst"`
	Gre                              types.Object                                               `tfsdk:"gre"`
	Oauth                            types.Object                                               `tfsdk:"oauth"`
	SplashGuestSponsorDomains        types.List                                                 `tfsdk:"splash_guest_sponsor_domains"`
	WalledGardenRanges               types.List                                                 `tfsdk:"walled_garden_ranges"`
	AvailabilityTags                 types.List                                                 `tfsdk:"availability_tags"`
	RadiusServers                    []NetworksWirelessSsidsResourceModelRadiusServer           `tfsdk:"radius_servers"`
	RadiusAccountingServers          []NetworksWirelessSsidsResourceModelRadiusAccountingServer `tfsdk:"radius_accounting_servers"`
	ApTagsAndVlanIds                 []NetworksWirelessSsidsResourceModelApTagsAndVlanId        `tfsdk:"ap_tags_and_vlan_ids"`
}

// NetworksWirelessSsidsResourceModel describes the resource data model.
type NetworksWirelessSsidsResourceModel struct {
	Id jsontypes.String `tfsdk:"id"`

	NetworkID jsontypes.String `tfsdk:"network_id"`
	Number    jsontypes.String `tfsdk:"number"`
	Serial    jsontypes.String `tfsdk:"serial"`
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
			"serial": schema.StringAttribute{
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
				MarkdownDescription: "ssid object",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"name": schema.StringAttribute{
						MarkdownDescription: "The name of the SSID",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"auth_mode": schema.StringAttribute{
						MarkdownDescription: "The association control method for the SSID ('open', 'open-enhanced', 'psk', 'open-with-radius', 'open-with-nac', '8021x-meraki', '8021x-nac', '8021x-radius', '8021x-google', '8021x-localradius', 'ipsk-with-radius', 'ipsk-without-radius' or 'ipsk-with-nac')",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"enterprise_admin_access": schema.StringAttribute{
						MarkdownDescription: "Whether or not an SSID is accessible by 'enterprise' administrators ('access disabled' or 'access enabled').",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"encryption_mode": schema.StringAttribute{
						MarkdownDescription: "The psk encryption mode for the SSID ('wep' or 'wpa'). This param is only valid if the authMode is 'psk'",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"psk": schema.StringAttribute{
						MarkdownDescription: "The passkey for the SSID. This param is only valid if the authMode is 'psk'",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"wpa_encryption_mode": schema.StringAttribute{
						MarkdownDescription: "The types of WPA encryption. ('WPA1 only', 'WPA1 and WPA2', 'WPA2 only', 'WPA3 Transition Mode', 'WPA3 only' or 'WPA3 192-bit Security').",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"splash_page": schema.StringAttribute{
						MarkdownDescription: "The type of splash page for the SSID ('None', 'Click-through splash page', 'Billing', 'Password-protected with Meraki RADIUS', 'Password-protected with custom RADIUS', 'Password-protected with Active Directory', 'Password-protected with LDAP', 'SMS authentication', 'Systems Manager Sentry', 'Facebook Wi-Fi', 'Google OAuth', 'Sponsored guest', 'Cisco ISE' or 'Google Apps domain'). This attribute is not supported for template children.",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"radius_called_station_id": schema.StringAttribute{
						MarkdownDescription: "The template of the called station identifier to be used for RADIUS (ex. $NODE_MAC$:$VAP_NUM$).",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"radius_authentication_nas_id": schema.StringAttribute{
						MarkdownDescription: "The template of the NAS identifier to be used for RADIUS authentication (ex. $NODE_MAC$:$VAP_NUM$).",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"radius_failover_policy": schema.StringAttribute{
						MarkdownDescription: "This policy determines how authentication requests should be handled in the event that all of the configured RADIUS servers are unreachable ('Deny access' or 'Allow access').",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"radius_load_balancing_policy": schema.StringAttribute{
						MarkdownDescription: "This policy determines which RADIUS server will be contacted first in an authentication attempt and the ordering of any necessary retry attempts ('Strict priority order' or 'Round robin').",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"radius_attribute_for_group_policies": schema.StringAttribute{
						MarkdownDescription: "Specify the RADIUS attribute used to look up group policies ('Filter-Id', 'Reply-Message', 'Airespace-ACL-Name' or 'Aruba-User-Role'). Access points must receive this attribute in the RADIUS Access-Accept message",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"ip_assignment_mode": schema.StringAttribute{
						MarkdownDescription: "The client IP assignment mode ('NAT mode', 'Bridge mode', 'Layer 3 roaming', 'Ethernet over GRE', 'Layer 3 roaming with a concentrator' or 'VPN').",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"concentrator_network_id": schema.StringAttribute{
						MarkdownDescription: "The concentrator to use when the ipAssignmentMode is 'Layer 3 roaming with a concentrator' or 'VPN'.",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"secondary_concentrator_network_id": schema.StringAttribute{
						MarkdownDescription: "The secondary concentrator to use when the ipAssignmentMode is 'VPN'. If configured, the APs will switch to using this concentrator if the primary concentrator is unreachable. This param is optional. ('disabled' represents no secondary concentrator.).",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"band_selection": schema.StringAttribute{
						MarkdownDescription: "The client-serving radio frequencies of this SSID in the default indoor RF profile. ('Dual band operation', '5 GHz band only' or 'Dual band operation with Band Steering')",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"radius_server_timeout": schema.Int64Attribute{
						MarkdownDescription: "The amount of time for which a RADIUS client waits for a reply from the RADIUS server (must be between 1-10 seconds).",
						Optional:            true,
						CustomType:          jsontypes.Int64Type,
					},
					"radius_server_attempts_limit": schema.Int64Attribute{
						MarkdownDescription: "The maximum number of transmit attempts after which a RADIUS server is failed over (must be between 1-5).",
						Optional:            true,
						CustomType:          jsontypes.Int64Type,
					},
					"radius_accounting_interim_interval": schema.Int64Attribute{
						MarkdownDescription: "The interval (in seconds) in which accounting information is updated and sent to the RADIUS accounting server.",
						Optional:            true,
						CustomType:          jsontypes.Int64Type,
					},
					"vlan_id": schema.Int64Attribute{
						MarkdownDescription: "The VLAN ID used for VLAN tagging. This param is only valid when the ipAssignmentMode is 'Layer 3 roaming with a concentrator' or 'VPN'",
						Optional:            true,
						CustomType:          jsontypes.Int64Type,
					},
					"default_vlan_id": schema.Int64Attribute{
						MarkdownDescription: "The default VLAN ID used for 'all other APs'. This param is only valid when the ipAssignmentMode is 'Bridge mode' or 'Layer 3 roaming'",
						Optional:            true,
						CustomType:          jsontypes.Int64Type,
					},
					"per_client_bandwidth_limit_up": schema.Int64Attribute{
						MarkdownDescription: "The upload bandwidth limit in Kbps. (0 represents no limit.)",
						Optional:            true,
						CustomType:          jsontypes.Int64Type,
					},
					"per_client_bandwidth_limit_down": schema.Int64Attribute{
						MarkdownDescription: "The download bandwidth limit in Kbps. (0 represents no limit.)",
						Optional:            true,
						CustomType:          jsontypes.Int64Type,
					},
					"per_ssid_bandwidth_limit_up": schema.Int64Attribute{
						MarkdownDescription: "The total upload bandwidth limit in Kbps. (0 represents no limit.)",
						Optional:            true,
						CustomType:          jsontypes.Int64Type,
					},
					"per_ssid_bandwidth_limit_down": schema.Int64Attribute{
						MarkdownDescription: "The total download bandwidth limit in Kbps. (0 represents no limit.)",
						Optional:            true,
						CustomType:          jsontypes.Int64Type,
					},
					"radius_guest_vlan_id": schema.Int64Attribute{
						MarkdownDescription: "VLAN ID of the RADIUS Guest VLAN. This param is only valid if the authMode is 'open-with-radius' and addressing mode is not set to 'isolated' or 'nat' mode",
						Optional:            true,
						CustomType:          jsontypes.Int64Type,
					},
					"min_bitrate": schema.Float64Attribute{
						MarkdownDescription: "The minimum bitrate in Mbps of this SSID in the default indoor RF profile. ('1', '2', '5.5', '6', '9', '11', '12', '18', '24', '36', '48' or '54')",
						Optional:            true,
						CustomType:          jsontypes.Float64Type,
					},
					"use_vlan_tagging": schema.BoolAttribute{
						MarkdownDescription: "Whether or not traffic should be directed to use specific VLANs. This param is only valid if the ipAssignmentMode is 'Bridge mode' or 'Layer 3 roaming'",
						Optional:            true,
						CustomType:          jsontypes.BoolType,
					},
					"disassociate_clients_on_vpn_failover": schema.BoolAttribute{
						MarkdownDescription: "Disassociate clients when 'VPN' concentrator failover occurs in order to trigger clients to re-associate and generate new DHCP requests. This param is only valid if ipAssignmentMode is 'VPN'.",
						Optional:            true,
						CustomType:          jsontypes.BoolType,
					},
					"radius_override": schema.BoolAttribute{
						MarkdownDescription: "If true, the RADIUS response can override VLAN tag. This is not valid when ipAssignmentMode is 'NAT mode'.",
						Optional:            true,
						CustomType:          jsontypes.BoolType,
					},
					"radius_guest_vlan_enabled": schema.BoolAttribute{
						MarkdownDescription: "Whether or not RADIUS Guest VLAN is enabled. This param is only valid if the authMode is 'open-with-radius' and addressing mode is not set to 'isolated' or 'nat' mode",
						Optional:            true,
						CustomType:          jsontypes.BoolType,
					},
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "Whether or not the SSID is enabled",
						Optional:            true,
						CustomType:          jsontypes.BoolType,
					},
					"radius_proxy_enabled": schema.BoolAttribute{
						MarkdownDescription: "If true, Meraki devices will proxy RADIUS messages through the Meraki cloud to the configured RADIUS auth and accounting servers.",
						Optional:            true,
						CustomType:          jsontypes.BoolType,
					},
					"radius_testing_enabled": schema.BoolAttribute{
						MarkdownDescription: "If true, Meraki devices will periodically send Access-Request messages to configured RADIUS servers using identity 'meraki_8021x_test' to ensure that the RADIUS servers are reachable.",
						Optional:            true,
						CustomType:          jsontypes.BoolType,
					},
					"radius_fallback_enabled": schema.BoolAttribute{
						MarkdownDescription: "Whether or not higher priority RADIUS servers should be retried after 60 seconds.",
						Optional:            true,
						CustomType:          jsontypes.BoolType,
					},
					"radius_coa_enabled": schema.BoolAttribute{
						MarkdownDescription: "If true, Meraki devices will act as a RADIUS Dynamic Authorization Server and will respond to RADIUS Change-of-Authorization and Disconnect messages sent by the RADIUS server.",
						Optional:            true,
						CustomType:          jsontypes.BoolType,
					},
					"radius_accounting_enabled": schema.BoolAttribute{
						MarkdownDescription: "Whether or not RADIUS accounting is enabled. This param is only valid if the authMode is 'open-with-radius', '8021x-radius' or 'ipsk-with-radius'",
						Optional:            true,
						CustomType:          jsontypes.BoolType,
					},
					"lan_isolation_enabled": schema.BoolAttribute{
						MarkdownDescription: "Boolean indicating whether Layer 2 LAN isolation should be enabled or disabled. Only configurable when ipAssignmentMode is 'Bridge mode'.",
						Optional:            true,
						CustomType:          jsontypes.BoolType,
					},
					"visible": schema.BoolAttribute{
						MarkdownDescription: "Boolean indicating whether APs should advertise or hide this SSID. APs will only broadcast this SSID if set to true.",
						Optional:            true,
						CustomType:          jsontypes.BoolType,
					},
					"available_on_all_aps": schema.BoolAttribute{
						MarkdownDescription: "Boolean indicating whether all APs should broadcast the SSID or if it should be restricted to APs matching any availability tags. Can only be false if the SSID has availability tags.",
						Optional:            true,
						CustomType:          jsontypes.BoolType,
					},
					"mandatory_dhcp_enabled": schema.BoolAttribute{
						MarkdownDescription: "If true, Mandatory DHCP will enforce that clients connecting to this SSID must use the IP address assigned by the DHCP server. Clients who use a static IP address won't be able to associate.",
						Optional:            true,
						CustomType:          jsontypes.BoolType,
					},
					"adult_content_filtering_enabled": schema.BoolAttribute{
						MarkdownDescription: "Boolean indicating whether or not adult content will be blocked.",
						Optional:            true,
						CustomType:          jsontypes.BoolType,
					},
					"walled_garden_enabled": schema.BoolAttribute{
						MarkdownDescription: "Allow access to a configurable list of IP ranges, which users may access prior to sign-on.",
						Optional:            true,
						CustomType:          jsontypes.BoolType,
					},
					"dot11w": schema.SingleNestedAttribute{
						MarkdownDescription: "The current setting for Protected Management Frames (802.11w).",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								MarkdownDescription: "Whether 802.11w is enabled or not.",
								Optional:            true,
								CustomType:          jsontypes.BoolType,
							},
							"required": schema.BoolAttribute{
								MarkdownDescription: "(Optional) Whether 802.11w is required or not.",
								Optional:            true,
								CustomType:          jsontypes.BoolType,
							},
						},
					},
					"dot11r": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "The current setting for 802.11r",
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								MarkdownDescription: "Whether 802.11r is enabled or not.",
								Optional:            true,
								CustomType:          jsontypes.BoolType,
							},
							"adaptive": schema.BoolAttribute{
								MarkdownDescription: "(Optional) Whether 802.11r is adaptive or not.",
								Optional:            true,
								CustomType:          jsontypes.BoolType,
							},
						},
					},
					"local_radius": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "The current setting for Local Authentication, a built-in RADIUS server on the access point. Only valid if authMode is '8021x-localradius'.",
						Attributes: map[string]schema.Attribute{
							"cache_timeout": schema.Int64Attribute{
								MarkdownDescription: "The duration (in seconds) for which LDAP and OCSP lookups are cached.",
								Optional:            true,
								CustomType:          jsontypes.Int64Type,
							},
							"password_authentication": schema.SingleNestedAttribute{
								Optional:            true,
								MarkdownDescription: "The current setting for password-based authentication.",
								Attributes: map[string]schema.Attribute{
									"enabled": schema.BoolAttribute{
										MarkdownDescription: "Whether or not to use EAP-TTLS/PAP or PEAP-GTC password-based authentication via LDAP lookup.",
										Optional:            true,
										CustomType:          jsontypes.BoolType,
									},
								},
							},
							"certificate_authentication": schema.SingleNestedAttribute{
								Optional:            true,
								MarkdownDescription: "The current setting for certificate verification.",
								Attributes: map[string]schema.Attribute{
									"enabled": schema.BoolAttribute{
										MarkdownDescription: "Whether or not to use EAP-TLS certificate-based authentication to validate wireless clients.",
										Optional:            true,
										CustomType:          jsontypes.BoolType,
									},
									"use_ldap": schema.BoolAttribute{
										MarkdownDescription: "Whether or not to verify the certificate with LDAP.",
										Optional:            true,
										CustomType:          jsontypes.BoolType,
									},
									"use_ocsp": schema.BoolAttribute{
										MarkdownDescription: "Whether or not to verify the certificate with OCSP.",
										Optional:            true,
										CustomType:          jsontypes.BoolType,
									},
									"ocsp_responder_url": schema.StringAttribute{
										MarkdownDescription: "(Optional) The URL of the OCSP responder to verify client certificate status.",
										Optional:            true,
										CustomType:          jsontypes.StringType,
									},
									"client_root_ca_certificate": schema.SingleNestedAttribute{
										Optional:            true,
										MarkdownDescription: "The Client CA Certificate used to sign the client certificate.",
										Attributes: map[string]schema.Attribute{
											"contents": schema.StringAttribute{
												MarkdownDescription: "The contents of the Client CA Certificate. Must be in PEM or DER format.",
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
						Optional:            true,
						MarkdownDescription: "The current setting for LDAP. Only valid if splashPage is 'Password-protected with LDAP'.",
						Attributes: map[string]schema.Attribute{
							"servers": schema.SetNestedAttribute{
								Optional:            true,
								Computed:            true,
								MarkdownDescription: "The LDAP servers to be used for authentication.",
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
								Optional:            true,
								MarkdownDescription: "(Optional) The credentials of the user account to be used by the AP to bind to your LDAP server. The LDAP account should have permissions on all your LDAP servers.",
								Attributes: map[string]schema.Attribute{
									"distinguished_name": schema.StringAttribute{
										MarkdownDescription: "(Optional) The credentials of the user account to be used by the AP to bind to your Active Directory server. The Active Directory account should have permissions on all your Active Directory servers. Only valid if the splashPage is 'Password-protected with Active Directory'.",
										Optional:            true,
										CustomType:          jsontypes.StringType,
									},
									"password": schema.StringAttribute{
										MarkdownDescription: "The password to the Active Directory user account..",
										Optional:            true,
										CustomType:          jsontypes.StringType,
									},
								},
							},
							"base_distinguished_name": schema.StringAttribute{
								MarkdownDescription: "The base distinguished name of users on the LDAP server.",
								Optional:            true,
								CustomType:          jsontypes.StringType,
							},
							"server_ca_certificate": schema.SingleNestedAttribute{
								Optional:            true,
								MarkdownDescription: "The CA certificate used to sign the LDAP server's key.",
								Attributes: map[string]schema.Attribute{
									"contents": schema.StringAttribute{
										MarkdownDescription: "The contents of the CA certificate. Must be in PEM or DER format.",
										Optional:            true,
										CustomType:          jsontypes.StringType,
									},
								},
							},
						},
					},
					"active_directory": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "The current setting for Active Directory. Only valid if splashPage is 'Password-protected with Active Directory'",
						Attributes: map[string]schema.Attribute{
							"servers": schema.SetNestedAttribute{
								Optional:            true,
								Computed:            true,
								MarkdownDescription: "The Active Directory servers to be used for authentication.",
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"host": schema.StringAttribute{
											MarkdownDescription: "IP address (or FQDN) of your Active Directory server.",
											Optional:            true,
											CustomType:          jsontypes.StringType,
										},
										"port": schema.Int64Attribute{
											MarkdownDescription: "(Optional) UDP port the Active Directory server listens on. By default, uses port 3268.",
											Optional:            true,
											CustomType:          jsontypes.Int64Type,
										},
									},
								},
							},
							"credentials": schema.SingleNestedAttribute{
								Optional:            true,
								MarkdownDescription: "(Optional) The credentials of the user account to be used by the AP to bind to your Active Directory server. The Active Directory account should have permissions on all your Active Directory servers. Only valid if the splashPage is 'Password-protected with Active Directory'.",
								Attributes: map[string]schema.Attribute{
									"logon_name": schema.StringAttribute{
										MarkdownDescription: "The logon name of the Active Directory account.",
										Optional:            true,
										CustomType:          jsontypes.StringType,
									},
									"password": schema.StringAttribute{
										MarkdownDescription: "The password to the Active Directory user account.",
										Optional:            true,
										CustomType:          jsontypes.StringType,
									},
								},
							},
						},
					},
					"dns_rewrite": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "DNS servers rewrite settings",
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								MarkdownDescription: "Boolean indicating whether or not DNS server rewrite is enabled. If disabled, upstream DNS will be used",
								Optional:            true,
								CustomType:          jsontypes.BoolType,
							},
							"dns_custom_nameservers": schema.ListAttribute{
								MarkdownDescription: "User specified DNS servers (up to two servers)",
								Optional:            true,
								ElementType:         jsontypes.StringType,
							},
						},
					},
					"speed_burst": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "The SpeedBurst setting for this SSID'",
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								MarkdownDescription: "Boolean indicating whether or not to allow users to temporarily exceed the bandwidth limit for short periods while still keeping them under the bandwidth limit over time.",
								Optional:            true,
								CustomType:          jsontypes.BoolType,
							},
						},
					},
					"gre": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Ethernet over GRE settings",
						Attributes: map[string]schema.Attribute{
							"concentrator": schema.SingleNestedAttribute{
								Optional:            true,
								MarkdownDescription: "The EoGRE concentrator's settings",
								Attributes: map[string]schema.Attribute{
									"host": schema.StringAttribute{
										MarkdownDescription: "The EoGRE concentrator's IP or FQDN. This param is required when ipAssignmentMode is 'Ethernet over GRE'.",
										Optional:            true,
										CustomType:          jsontypes.StringType,
									},
								},
							},
							"key": schema.Int64Attribute{
								MarkdownDescription: "Optional numerical identifier that will add the GRE key field to the GRE header. Used to identify an individual traffic flow within a tunnel.",
								Optional:            true,
								CustomType:          jsontypes.Int64Type,
							},
						},
					},
					"oauth": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "The OAuth settings of this SSID. Only valid if splashPage is 'Google OAuth'.",
						Attributes: map[string]schema.Attribute{
							"allowed_domains": schema.ListAttribute{
								MarkdownDescription: "(Optional) The list of domains allowed access to the network.",
								Optional:            true,
								ElementType:         jsontypes.StringType,
							},
						},
					},
					"splash_guest_sponsor_domains": schema.ListAttribute{
						MarkdownDescription: "Array of valid sponsor email domains for sponsored guest splash type.",
						Optional:            true,
						ElementType:         jsontypes.StringType,
					},
					"walled_garden_ranges": schema.ListAttribute{
						MarkdownDescription: "Specify your walled garden by entering an array of addresses, ranges using CIDR notation, domain names, and domain wildcards (e.g. '192.168.1.1/24', '192.168.37.10/32', 'www.yahoo.com', '*.google.com']). Meraki's splash page is automatically included in your walled garden.",
						Optional:            true,
						ElementType:         jsontypes.StringType,
					},
					"availability_tags": schema.ListAttribute{
						MarkdownDescription: "Accepts a list of tags for this SSID. If availableOnAllAps is false, then the SSID will only be broadcast by APs with tags matching any of the tags in this list.",
						Optional:            true,
						ElementType:         jsontypes.StringType,
					},
					"radius_servers": schema.SetNestedAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The RADIUS 802.1X servers to be used for authentication. This param is only valid if the authMode is 'open-with-radius', '8021x-radius' or 'ipsk-with-radius'",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"host": schema.StringAttribute{
									MarkdownDescription: "IP address (or FQDN) of your RADIUS server",
									Optional:            true,
									CustomType:          jsontypes.StringType,
								},
								"secret": schema.StringAttribute{
									MarkdownDescription: "RADIUS client shared secret",
									Optional:            true,
									CustomType:          jsontypes.StringType,
								},
								"ca_certificate": schema.StringAttribute{
									MarkdownDescription: "Certificate used for authorization for the RADSEC Server",
									Optional:            true,
									CustomType:          jsontypes.StringType,
								},
								"port": schema.Int64Attribute{
									MarkdownDescription: "UDP port the RADIUS server listens on for Access-requests",
									Optional:            true,
									CustomType:          jsontypes.Int64Type,
								},
								"open_roaming_certificate_id": schema.Int64Attribute{
									MarkdownDescription: "The ID of the Openroaming Certificate attached to radius server.",
									Optional:            true,
									CustomType:          jsontypes.Int64Type,
								},
								"radsec_enabled": schema.BoolAttribute{
									MarkdownDescription: "Use RADSEC (TLS over TCP) to connect to this RADIUS server. Requires radiusProxyEnabled.",
									Optional:            true,
									CustomType:          jsontypes.BoolType,
								},
							},
						},
					},
					"radius_accounting_servers": schema.SetNestedAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The RADIUS accounting 802.1X servers to be used for authentication. This param is only valid if the authMode is 'open-with-radius', '8021x-radius' or 'ipsk-with-radius' and radiusAccountingEnabled is 'true'",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"host": schema.StringAttribute{
									MarkdownDescription: "IP address (or FQDN) to which the APs will send RADIUS accounting messages",
									Optional:            true,
									CustomType:          jsontypes.StringType,
								},
								"secret": schema.StringAttribute{
									MarkdownDescription: "Shared key used to authenticate messages between the APs and RADIUS server",
									Optional:            true,
									CustomType:          jsontypes.StringType,
								},
								"ca_certificate": schema.StringAttribute{
									MarkdownDescription: "Certificate used for authorization for the RADSEC Server",
									Optional:            true,
									CustomType:          jsontypes.StringType,
								},
								"port": schema.Int64Attribute{
									MarkdownDescription: "Port on the RADIUS server that is listening for accounting messages",
									Optional:            true,
									CustomType:          jsontypes.Int64Type,
								},
								"radsec_enabled": schema.BoolAttribute{
									MarkdownDescription: "Use RADSEC (TLS over TCP) to connect to this RADIUS accounting server. Requires radiusProxyEnabled.",
									Optional:            true,
									CustomType:          jsontypes.BoolType,
								},
							},
						},
					},
					"ap_tags_and_vlan_ids": schema.SetNestedAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The list of tags and VLAN IDs used for VLAN tagging. This param is only valid when the ipAssignmentMode is 'Bridge mode' or 'Layer 3 roaming'",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"vlan_id": schema.Int64Attribute{
									MarkdownDescription: "Numerical identifier that is assigned to the VLAN",
									Optional:            true,
									CustomType:          jsontypes.Int64Type,
								},
								"tags": schema.ListAttribute{
									MarkdownDescription: "Array of AP tags",
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

	request := NetworksWirelessSsidsPayload(ctx, data)

	_, httpResp, err := r.client.WirelessApi.UpdateNetworkWirelessSsid(context.Background(), data.NetworkID.ValueString(), data.Number.ValueString()).UpdateNetworkWirelessSsidRequest(request).Execute()
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

	request := NetworksWirelessSsidsPayload(ctx, data)
	_, httpResp, err := r.client.WirelessApi.UpdateNetworkWirelessSsid(context.Background(), data.NetworkID.ValueString(), data.Number.ValueString()).UpdateNetworkWirelessSsidRequest(request).Execute()
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

	// Create HTTP request body
	updateDevice := openApiClient.NewUpdateDeviceRequest()

	var name string
	var tags []string
	var lat float32
	var lng float32
	var address string
	var notes string
	var moveMapMarker bool
	//var switchProfileId string
	//var floorPlanId string

	updateDevice.Name = &name
	updateDevice.Tags = tags
	updateDevice.Lat = &lat
	updateDevice.Lng = &lng
	updateDevice.Address = &address
	updateDevice.Notes = &notes
	updateDevice.MoveMapMarker = &moveMapMarker

	// Initialize provider client and make API call
	_, httpResp, err := r.client.DevicesApi.UpdateDevice(context.Background(), data.Serial.ValueString()).UpdateDeviceRequest(*updateDevice).Execute()

	// Check for API success response code
	if httpResp.StatusCode == 404 {
		resp.Diagnostics.AddWarning(
			"no ssid number/serial or network not found in API",
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

	}

	data.Id = jsontypes.StringValue("example-id")

	resp.State.RemoveResource(ctx)

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

func NetworksWirelessSsidsPayload(ctx context.Context, data *NetworksWirelessSsidsResourceModel) openApiClient.UpdateNetworkWirelessSsidRequest {
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
	var dot11wObject NetworksWirelessSsidsResourceModelDot11W
	ssid.Dot11W.As(ctx, &dot11wObject, basetypes.ObjectAsOptions{})
	dot11w.SetEnabled(dot11wObject.Enabled.ValueBool())
	dot11w.SetRequired(dot11wObject.Required.ValueBool())
	request.SetDot11w(*dot11w)

	dot11r := openApiClient.NewUpdateNetworkWirelessSsidRequestDot11r()
	var dot11rObject NetworksWirelessSsidsResourceModelDot11R
	ssid.Dot11R.As(ctx, &dot11rObject, basetypes.ObjectAsOptions{})
	dot11r.SetEnabled(dot11rObject.Enabled.ValueBool())
	dot11r.SetAdaptive(dot11rObject.Adaptive.ValueBool())
	request.SetDot11r(*dot11r)

	localRadius := openApiClient.NewUpdateNetworkWirelessSsidRequestLocalRadius()
	var localRadiusObject NetworksWirelessSsidsResourceModelLocalRadius
	ssid.LocalRadius.As(ctx, &localRadiusObject, basetypes.ObjectAsOptions{})
	localRadius.SetCacheTimeout(int32(localRadiusObject.CacheTimeout.ValueInt64()))
	certificateAuthentication := openApiClient.NewUpdateNetworkWirelessSsidRequestLocalRadiusCertificateAuthentication()
	var certificateAuthenticationObject NetworksWirelessSsidsResourceModelLocalRadiusCertificateAuthentication
	localRadiusObject.CertificateAuthentication.As(ctx, &certificateAuthenticationObject, basetypes.ObjectAsOptions{})
	certificateAuthentication.SetEnabled(certificateAuthenticationObject.Enabled.ValueBool())
	certificateAuthentication.SetOcspResponderUrl(certificateAuthenticationObject.OcspResponderUrl.ValueString())
	certificateAuthentication.SetUseLdap(certificateAuthenticationObject.UseLdap.ValueBool())
	certificateAuthentication.SetUseOcsp(certificateAuthenticationObject.UseOcsp.ValueBool())
	var clientRootCaCertificateObject NetworksWirelessSsidsResourceModelLocalRadiusCertificateAuthenticationClientRootCaCertificate
	certificateAuthenticationObject.ClientRootCaCertificate.As(ctx, &clientRootCaCertificateObject, basetypes.ObjectAsOptions{})
	rootCACertificate := openApiClient.NewUpdateNetworkWirelessSsidRequestLocalRadiusCertificateAuthenticationClientRootCaCertificate()
	rootCACertificate.SetContents(clientRootCaCertificateObject.Contents.ValueString())
	certificateAuthentication.SetClientRootCaCertificate(*rootCACertificate)
	localRadius.SetCertificateAuthentication(*certificateAuthentication)
	request.SetLocalRadius(*localRadius)

	ldap := openApiClient.NewUpdateNetworkWirelessSsidRequestLdap()
	var ldapObject NetworksWirelessSsidsResourceModelLdap
	ssid.Ldap.As(ctx, &ldapObject, basetypes.ObjectAsOptions{})
	ldap.SetBaseDistinguishedName(ldapObject.BaseDistinguishedName.ValueString())
	ldapCredentials := openApiClient.NewUpdateNetworkWirelessSsidRequestLdapCredentials()
	var ldapCredentialsObject NetworksWirelessSsidsResourceModelCredential
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
	var ldapServerCaCertificateObject NetworksWirelessSsidsResourceModelServerCaCertificate
	ldapObject.ServerCaCertificate.As(ctx, &ldapServerCaCertificateObject, basetypes.ObjectAsOptions{})
	ldapServerCaCertificate := openApiClient.NewUpdateNetworkWirelessSsidRequestLdapServerCaCertificate()
	ldapServerCaCertificate.SetContents(ldapServerCaCertificateObject.Contents.ValueString())
	ldap.SetServerCaCertificate(*ldapServerCaCertificate)
	request.SetLdap(*ldap)

	ad := openApiClient.NewUpdateNetworkWirelessSsidRequestActiveDirectory()
	var activeDirectoryObject NetworksWirelessSsidsResourceModelActiveDirectory
	ssid.ActiveDirectory.As(ctx, &activeDirectoryObject, basetypes.ObjectAsOptions{})
	var adServers []openApiClient.UpdateNetworkWirelessSsidRequestActiveDirectoryServersInner
	for _, server := range activeDirectoryObject.Servers {
		inner := openApiClient.NewUpdateNetworkWirelessSsidRequestActiveDirectoryServersInner(server.Host.ValueString())
		inner.SetPort(int32(server.Port.ValueInt64()))
		adServers = append(adServers, *inner)
	}
	ad.SetServers(adServers)
	var adCredentialsObject NetworksWirelessSsidsResourceModelActiveDirectoryCredential
	activeDirectoryObject.Credentials.As(ctx, &adCredentialsObject, basetypes.ObjectAsOptions{})
	adCredentials := openApiClient.NewUpdateNetworkWirelessSsidRequestActiveDirectoryCredentials()
	adCredentials.SetPassword(adCredentialsObject.Password.ValueString())
	adCredentials.SetLogonName(adCredentialsObject.LogonName.ValueString())
	ad.SetCredentials(*adCredentials)
	request.SetActiveDirectory(*ad)

	var dnsRewriteObject NetworksWirelessSsidsResourceModelDNSRewrite
	ssid.DnsRewrite.As(ctx, &dnsRewriteObject, basetypes.ObjectAsOptions{})
	dnsRewrite := openApiClient.NewUpdateNetworkWirelessSsidRequestDnsRewrite()
	dnsRewrite.SetEnabled(dnsRewriteObject.Enabled.ValueBool())
	var dnsServers []string
	for _, server := range dnsRewriteObject.DnsCustomNameservers.Elements() {
		dnsServers = append(dnsServers, server.String())
	}
	dnsRewrite.SetDnsCustomNameservers(dnsServers)
	request.SetDnsRewrite(*dnsRewrite)

	var speedBurstObject NetworksWirelessSsidsResourceModelSpeedBurst
	ssid.SpeedBurst.As(ctx, &speedBurstObject, basetypes.ObjectAsOptions{})
	speedBurst := openApiClient.NewUpdateNetworkWirelessSsidRequestSpeedBurst()
	speedBurst.SetEnabled(speedBurstObject.Enabled.ValueBool())
	request.SetSpeedBurst(*speedBurst)

	var greObject NetworksWirelessSsidsResourceModelGre
	ssid.Gre.As(ctx, &greObject, basetypes.ObjectAsOptions{})
	gre := openApiClient.NewUpdateNetworkWirelessSsidRequestGre()
	gre.SetKey(int32(greObject.Key.ValueInt64()))
	var greConentratorObject NetworksWirelessSsidsResourceModelConcentrator
	greObject.Concentrator.As(ctx, &greConentratorObject, basetypes.ObjectAsOptions{})
	concentrator := openApiClient.NewUpdateNetworkWirelessSsidRequestGreConcentrator(greConentratorObject.Host.ValueString())
	gre.SetConcentrator(*concentrator)
	request.SetGre(*gre)

	var oauthObject NetworksWirelessSsidsResourceModelOauth
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

	var radiusServers []openApiClient.UpdateNetworkWirelessSsidRequestRadiusServersInner
	for _, radius := range ssid.RadiusServers {
		radiusServer := openApiClient.NewUpdateNetworkWirelessSsidRequestRadiusServersInner(radius.Host.ValueString())
		radiusServer.SetPort(int32(radius.Port.ValueInt64()))
		radiusServer.SetSecret(radius.Secret.ValueString())
		radiusServers = append(radiusServers, *radiusServer)
	}
	request.SetRadiusServers(radiusServers)

	var radiusAccountingServers []openApiClient.UpdateNetworkWirelessSsidRequestRadiusAccountingServersInner
	for _, radius := range ssid.RadiusAccountingServers {
		radiusServer := openApiClient.NewUpdateNetworkWirelessSsidRequestRadiusAccountingServersInner(radius.Host.ValueString())
		radiusServer.SetPort(int32(radius.Port.ValueInt64()))
		radiusServer.SetSecret(radius.Secret.ValueString())
		radiusServer.SetCaCertificate(radius.CaCertificate.ValueString())
		radiusServer.SetRadsecEnabled(radius.RadsecEnabled.ValueBool())
		radiusAccountingServers = append(radiusAccountingServers, *radiusServer)
	}
	request.SetRadiusAccountingServers(radiusAccountingServers)

	var apTagsAndVlanIds []openApiClient.UpdateNetworkWirelessSsidRequestApTagsAndVlanIdsInner
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

	return *request
}
