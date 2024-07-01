package provider

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/core-infra-svcs/terraform-provider-meraki/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/float64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/objectplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &NetworksWirelessSsidsResource{}
	_ resource.ResourceWithConfigure   = &NetworksWirelessSsidsResource{}
	_ resource.ResourceWithImportState = &NetworksWirelessSsidsResource{}
)

func NewNetworksWirelessSsidsResource() resource.Resource {
	return &NetworksWirelessSsidsResource{
		typeName: "meraki_networks_wireless_ssids",
	}
}

// NetworksWirelessSsidsResource defines the resource implementation.
type NetworksWirelessSsidsResource struct {
	client   *openApiClient.APIClient
	typeName string
}

// NetworksWirelessSsidResourceModel represents the internal structure for the Networks Wireless SSID resource
type NetworksWirelessSsidResourceModel struct {
	Id                               types.String  `tfsdk:"id" json:"id"`
	NetworkId                        types.String  `tfsdk:"network_id" json:"networkId"`
	Number                           types.Int64   `tfsdk:"number" json:"number"`
	Name                             types.String  `tfsdk:"name" json:"name"`
	Enabled                          types.Bool    `tfsdk:"enabled" json:"enabled"`
	AuthMode                         types.String  `tfsdk:"auth_mode" json:"authMode"`
	EnterpriseAdminAccess            types.String  `tfsdk:"enterprise_admin_access" json:"enterpriseAdminAccess"`
	EncryptionMode                   types.String  `tfsdk:"encryption_mode" json:"encryptionMode"`
	PSK                              types.String  `tfsdk:"psk" json:"psk"`
	WPAEncryptionMode                types.String  `tfsdk:"wpa_encryption_mode" json:"wpaEncryptionMode"`
	Dot11w                           types.Object  `tfsdk:"dot11w" json:"dot11w"`
	Dot11r                           types.Object  `tfsdk:"dot11r" json:"dot11r"`
	SplashPage                       types.String  `tfsdk:"splash_page" json:"splashPage"`
	SplashGuestSponsorDomains        types.List    `tfsdk:"splash_guest_sponsor_domains" json:"splashGuestSponsorDomains"`
	OAuth                            types.Object  `tfsdk:"oauth" json:"oauth"`
	LocalRadius                      types.Object  `tfsdk:"local_radius" json:"localRadius"`
	LDAP                             types.Object  `tfsdk:"ldap" json:"ldap"`
	ActiveDirectory                  types.Object  `tfsdk:"active_directory" json:"activeDirectory"`
	RadiusServers                    types.List    `tfsdk:"radius_servers" json:"radiusServers"`
	RadiusProxyEnabled               types.Bool    `tfsdk:"radius_proxy_enabled" json:"radiusProxyEnabled"`
	RadiusTestingEnabled             types.Bool    `tfsdk:"radius_testing_enabled" json:"radiusTestingEnabled"`
	RadiusCalledStationID            types.String  `tfsdk:"radius_called_station_id" json:"radiusCalledStationId"`
	RadiusAuthenticationNASID        types.String  `tfsdk:"radius_authentication_nas_id" json:"radiusAuthenticationNasId"`
	RadiusServerTimeout              types.Int64   `tfsdk:"radius_server_timeout" json:"radiusServerTimeout"`
	RadiusServerAttemptsLimit        types.Int64   `tfsdk:"radius_server_attempts_limit" json:"radiusServerAttemptsLimit"`
	RadiusFallbackEnabled            types.Bool    `tfsdk:"radius_fallback_enabled" json:"radiusFallbackEnabled"`
	RadiusCoaEnabled                 types.Bool    `tfsdk:"radius_coa_enabled" json:"radiusCoaEnabled"`
	RadiusFailOverPolicy             types.String  `tfsdk:"radius_fail_over_policy" json:"radiusFailoverPolicy"`
	RadiusLoadBalancingPolicy        types.String  `tfsdk:"radius_load_balancing_policy" json:"radiusLoadBalancingPolicy"`
	RadiusAccountingEnabled          types.Bool    `tfsdk:"radius_accounting_enabled" json:"radiusAccountingEnabled"`
	RadiusAccountingServers          types.List    `tfsdk:"radius_accounting_servers" json:"radiusAccountingServers"`
	RadiusAccountingInterimInterval  types.Int64   `tfsdk:"radius_accounting_interim_interval" json:"radiusAccountingInterimInterval"`
	RadiusAttributeForGroupPolicies  types.String  `tfsdk:"radius_attribute_for_group_policies" json:"radiusAttributeForGroupPolicies"`
	IPAssignmentMode                 types.String  `tfsdk:"ip_assignment_mode" json:"ipAssignmentMode"`
	UseVlanTagging                   types.Bool    `tfsdk:"use_vlan_tagging" json:"useVlanTagging"`
	ConcentratorNetworkID            types.String  `tfsdk:"concentrator_network_id" json:"concentratorNetworkId"`
	SecondaryConcentratorNetworkID   types.String  `tfsdk:"secondary_concentrator_network_id" json:"secondaryConcentratorNetworkId"`
	DisassociateClientsOnVpnFailOver types.Bool    `tfsdk:"disassociate_clients_on_vpn_fail_over" json:"disassociateClientsOnVpnFailover"`
	VlanID                           types.Int64   `tfsdk:"vlan_id" json:"vlanId"`
	DefaultVlanID                    types.Int64   `tfsdk:"default_vlan_id" json:"defaultVlanId"`
	ApTagsAndVlanIDs                 types.List    `tfsdk:"ap_tags_and_vlan_ids" json:"apTagsAndVlanIds"`
	WalledGardenEnabled              types.Bool    `tfsdk:"walled_garden_enabled" json:"walledGardenEnabled"`
	WalledGardenRanges               types.List    `tfsdk:"walled_garden_ranges" json:"walledGardenRanges"`
	GRE                              types.Object  `tfsdk:"gre" json:"gre"`
	RadiusOverride                   types.Bool    `tfsdk:"radius_override" json:"radiusOverride"`
	RadiusGuestVlanEnabled           types.Bool    `tfsdk:"radius_guest_vlan_enabled" json:"radiusGuestVlanEnabled"`
	RadiusGuestVlanId                types.Int64   `tfsdk:"radius_guest_vlan_id" json:"radiusGuestVlanId"`
	MinBitRate                       types.Float64 `tfsdk:"min_bitrate" json:"minBitRate"`
	BandSelection                    types.String  `tfsdk:"band_selection" json:"bandSelection"`
	PerClientBandwidthLimitUp        types.Int64   `tfsdk:"per_client_bandwidth_limit_up" json:"perClientBandwidthLimitUp"`
	PerClientBandwidthLimitDown      types.Int64   `tfsdk:"per_client_bandwidth_limit_down" json:"perClientBandwidthLimitDown"`
	PerSsidBandwidthLimitUp          types.Int64   `tfsdk:"per_ssid_bandwidth_limit_up" json:"perSsidBandwidthLimitUp"`
	PerSsidBandwidthLimitDown        types.Int64   `tfsdk:"per_ssid_bandwidth_limit_down" json:"perSsidBandwidthLimitDown"`
	LanIsolationEnabled              types.Bool    `tfsdk:"lan_isolation_enabled" json:"lanIsolationEnabled"`
	Visible                          types.Bool    `tfsdk:"visible" json:"visible"`
	AvailableOnAllAps                types.Bool    `tfsdk:"available_on_all_aps" json:"availableOnAllAps"`
	AvailabilityTags                 types.List    `tfsdk:"availability_tags" json:"availabilityTags"`
	MandatoryDhcpEnabled             types.Bool    `tfsdk:"mandatory_dhcp_enabled" json:"mandatoryDhcpEnabled"`
	AdultContentFilteringEnabled     types.Bool    `tfsdk:"adult_content_filtering_enabled" json:"adultContentFilteringEnabled"`
	DnsRewrite                       types.Object  `tfsdk:"dns_rewrite" json:"dnsRewrite"`
	SpeedBurst                       types.Object  `tfsdk:"speed_burst" json:"speedBurst"`
	SsidAdminAccessible              types.Bool    `tfsdk:"ssid_admin_accessible" json:"ssidAdminAccessible"`
	LocalAuth                        types.Bool    `tfsdk:"local_auth" json:"localAuth"`
	RadiusEnabled                    types.Bool    `tfsdk:"radius_enabled" json:"radiusEnabled"`
	AdminSplashUrl                   types.String  `tfsdk:"admin_splash_url" json:"adminSplashUrl"`
	SplashTimeout                    types.String  `tfsdk:"splash_timeout" json:"splashTimeout"`
	NamedVlans                       types.Object  `tfsdk:"named_vlans" json:"namedVlans"`
}

// Dot11w represents the structure for 802.11w settings
type Dot11w struct {
	Enabled  types.Bool `tfsdk:"enabled" json:"enabled"`
	Required types.Bool `tfsdk:"required" json:"required"`
}

// Dot11r represents the structure for 802.11r settings
type Dot11r struct {
	Enabled  types.Bool `tfsdk:"enabled" json:"enabled"`
	Adaptive types.Bool `tfsdk:"adaptive" json:"adaptive"`
}

// OAuth represents the structure for OAuth settings
type OAuth struct {
	AllowedDomains types.List `tfsdk:"allowed_domains" json:"allowedDomains"`
}

// LocalRadius represents the structure for Local RADIUS server settings
type LocalRadius struct {
	CacheTimeout              types.Int64  `tfsdk:"cache_timeout" json:"cacheTimeout"`
	PasswordAuthentication    types.Object `tfsdk:"password_authentication" json:"passwordAuthentication"`
	CertificateAuthentication types.Object `tfsdk:"certificate_authentication" json:"certificateAuthentication"`
}

// PasswordAuthentication represents the structure for password-based authentication settings
type PasswordAuthentication struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled"`
}

// CertificateAuthentication represents the structure for certificate-based authentication settings
type CertificateAuthentication struct {
	Enabled                 types.Bool   `tfsdk:"enabled" json:"enabled"`
	UseLdap                 types.Bool   `tfsdk:"use_ldap" json:"useLdap"`
	UseOcsp                 types.Bool   `tfsdk:"use_ocsp" json:"useOcsp"`
	OcspResponderUrl        types.String `tfsdk:"ocsp_responder_url" json:"ocspResponderUrl"`
	ClientRootCaCertificate types.Object `tfsdk:"client_root_ca_certificate" json:"clientRootCaCertificate"`
}

// CaCertificate represents the structure for CA certificate settings
type CaCertificate struct {
	Contents types.String `tfsdk:"contents" json:"contents"`
}

// LDAP represents the structure for LDAP server settings
type LDAP struct {
	Servers               types.List   `tfsdk:"servers" json:"servers"`
	Credentials           types.Object `tfsdk:"credentials" json:"credentials"`
	BaseDistinguishedName types.String `tfsdk:"distinguished_name" json:"baseDistinguishedName"`
	ServerCaCertificate   types.Object `tfsdk:"server_ca_certificate" json:"serverCaCertificate"`
}

// LdapServer represents the structure for an LDAP server
type LdapServer struct {
	Host types.String `tfsdk:"host" json:"host"`
	Port types.Int64  `tfsdk:"port" json:"port"`
}

// LdapServerCaCertificate represents the structure for LDAP server certificate
type LdapServerCaCertificate struct {
	Contents types.String `tfsdk:"contents" json:"contents"`
}

// LdapCredentials represents the structure for LDAP server credentials
type LdapCredentials struct {
	DistinguishedName types.String `tfsdk:"distinguished_name" json:"distinguishedName"`
	Password          types.String `tfsdk:"password" json:"password"`
}

// ActiveDirectory represents the structure for Active Directory settings
type ActiveDirectory struct {
	Servers     types.List   `tfsdk:"servers" json:"servers"`
	Credentials types.Object `tfsdk:"credentials" json:"credentials"`
}

// ActiveDirectoryServer represents the structure for an Active Directory server
type ActiveDirectoryServer struct {
	Host types.String `tfsdk:"host" json:"host"`
	Port types.Int64  `tfsdk:"port" json:"port"`
}

// AdCredentials represents the structure for Active Directory credentials
type AdCredentials struct {
	Password  types.String `tfsdk:"password" json:"password"`
	LoginName types.String `tfsdk:"login_name" json:"loginName"`
}

// RadiusServer represents the structure for a RADIUS server
type RadiusServer struct {
	Host                     types.String `tfsdk:"host" json:"host"`
	Port                     types.Int64  `tfsdk:"port" json:"port"`
	Secret                   types.String `tfsdk:"secret" json:"secret"`
	RadSecEnabled            types.Bool   `tfsdk:"rad_sec_enabled" json:"radSecEnabled"`
	OpenRoamingCertificateID types.Int64  `tfsdk:"open_roaming_certificate_id" json:"openRoamingCertificateId"`
	CaCertificate            types.String `tfsdk:"ca_certificate" json:"caCertificate"`
}

// ApTagsAndVlanID represents the structure for AP tags and VLAN IDs
type ApTagsAndVlanID struct {
	Tags   types.List  `tfsdk:"tags" json:"tags"`
	VlanId types.Int64 `tfsdk:"vlan_id" json:"vlanId"`
}

// GRE represents the structure for GRE tunnel settings
type GRE struct {
	Concentrator types.Object `tfsdk:"concentrator" json:"concentrator"`
	Key          types.Int64  `tfsdk:"key" json:"key"`
}

// GreConcentrator represents the structure for GRE concentrator settings
type GreConcentrator struct {
	Host types.String `tfsdk:"host" json:"host"`
}

// DnsRewrite represents the structure for DNS rewrite settings
type DnsRewrite struct {
	Enabled              types.Bool `tfsdk:"enabled" json:"enabled"`
	DnsCustomNameservers types.List `tfsdk:"dns_custom_name_servers" json:"dnsCustomNameservers"`
}

// SpeedBurst represents the structure for speed burst settings
type SpeedBurst struct {
	Enabled types.Bool `tfsdk:"enabled" json:"enabled"`
}

// NamedVlans represents the structure for named VLANs configuration
type NamedVlans struct {
	Tagging types.Object `tfsdk:"tagging" json:"tagging"`
	Radius  types.Object `tfsdk:"radius" json:"radius"`
}

// Tagging represents the structure for VLAN tagging settings
type Tagging struct {
	Enabled         types.Bool   `tfsdk:"enabled" json:"enabled"`
	DefaultVlanName types.String `tfsdk:"default_vlan_name" json:"defaultVlanName"`
	ByApTags        types.List   `tfsdk:"by_ap_tags" json:"byApTags"`
}

// ByApTag represents the structure for AP tag and VLAN name association
type ByApTag struct {
	Tags     types.List   `tfsdk:"tags" json:"tags"`
	VlanName types.String `tfsdk:"vlan_name" json:"vlanName"`
}

// Radius represents the structure for RADIUS settings in named VLANs
type Radius struct {
	GuestVlan types.Object `tfsdk:"guest_vlan" json:"guestVlan"`
}

// GuestVlan represents the structure for the guest VLAN settings
type GuestVlan struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled"`
	Name    types.String `tfsdk:"name" json:"name"`
}

// RadiusGuestVlan represents the structure for the RADIUS guest VLAN settings
type RadiusGuestVlan struct {
	Enabled types.Bool   `tfsdk:"enabled" json:"enabled"`
	Name    types.String `tfsdk:"name" json:"name"`
}

func NetworksWirelessSsidPayloadDot11w(input types.Object) (*openApiClient.UpdateNetworkApplianceSsidRequestDot11w, diag.Diagnostics) {
	if input.IsNull() || input.IsUnknown() {
		return nil, nil
	}

	var diags diag.Diagnostics
	var dot11wObject Dot11w

	err := input.As(context.Background(), &dot11wObject, basetypes.ObjectAsOptions{})
	if err.HasError() {
		diags = append(diags, err.Errors()...)
		return nil, diags
	}

	return &openApiClient.UpdateNetworkApplianceSsidRequestDot11w{
		Enabled:  dot11wObject.Enabled.ValueBoolPointer(),
		Required: dot11wObject.Required.ValueBoolPointer(),
	}, diags
}

func NetworksWirelessSsidPayloadDot11r(input types.Object) (*openApiClient.UpdateNetworkWirelessSsidRequestDot11r, diag.Diagnostics) {
	if input.IsNull() || input.IsUnknown() {
		return nil, nil
	}

	var diags diag.Diagnostics
	var dot11rObject Dot11r

	err := input.As(context.Background(), &dot11rObject, basetypes.ObjectAsOptions{})
	if err.HasError() {
		diags = append(diags, err.Errors()...)
		return nil, diags
	}

	return &openApiClient.UpdateNetworkWirelessSsidRequestDot11r{
		Enabled:  dot11rObject.Enabled.ValueBoolPointer(),
		Adaptive: dot11rObject.Adaptive.ValueBoolPointer(),
	}, diags
}

func NetworksWirelessSsidPayloadOauth(input types.Object) (*openApiClient.UpdateNetworkWirelessSsidRequestOauth, diag.Diagnostics) {
	if input.IsNull() || input.IsUnknown() {
		return nil, nil
	}

	var diags diag.Diagnostics
	var oauthObject OAuth

	err := input.As(context.Background(), &oauthObject, basetypes.ObjectAsOptions{
		UnhandledNullAsEmpty: true,
	})
	if err.HasError() {
		diags = append(diags, err.Errors()...)
		return nil, diags
	}

	var allowedDomains []string
	allowedDomainsList := oauthObject.AllowedDomains.Elements()

	for _, domain := range allowedDomainsList {
		allowedDomains = append(allowedDomains, domain.String())
	}

	return &openApiClient.UpdateNetworkWirelessSsidRequestOauth{
		AllowedDomains: allowedDomains,
	}, diags
}

func NetworksWirelessSsidPayloadLocalRadius(input types.Object) (openApiClient.UpdateNetworkWirelessSsidRequestLocalRadius, diag.Diagnostics) {
	var result openApiClient.UpdateNetworkWirelessSsidRequestLocalRadius
	if input.IsNull() || input.IsUnknown() {
		return result, nil
	}

	var diags diag.Diagnostics

	// Unmarshalling input to LocalRadius struct
	var localRadius LocalRadius
	err := input.As(context.Background(), &localRadius, basetypes.ObjectAsOptions{})
	if err.HasError() {
		diags = append(diags, diag.NewErrorDiagnostic("Error converting LocalRadius", fmt.Sprintf("%s", err.Errors())))
		return result, diags
	}

	// CacheTimeout
	if !localRadius.CacheTimeout.IsUnknown() && !localRadius.CacheTimeout.IsNull() {
		cacheTimeout := int32(localRadius.CacheTimeout.ValueInt64())
		result.SetCacheTimeout(cacheTimeout)
	}

	// PasswordAuthentication
	if !localRadius.PasswordAuthentication.IsUnknown() && !localRadius.PasswordAuthentication.IsNull() {
		var passwordAuth PasswordAuthentication
		err := localRadius.PasswordAuthentication.As(context.Background(), &passwordAuth, basetypes.ObjectAsOptions{})
		if err.HasError() {
			diags = append(diags, diag.NewErrorDiagnostic("Error converting PasswordAuthentication", fmt.Sprintf("%s", err.Errors())))
		} else {
			var passwordAuthentication openApiClient.UpdateNetworkWirelessSsidRequestLocalRadiusPasswordAuthentication
			passwordAuthentication.SetEnabled(passwordAuth.Enabled.ValueBool())
			result.SetPasswordAuthentication(passwordAuthentication)
		}
	}

	// CertificateAuthentication
	if !localRadius.CertificateAuthentication.IsNull() && !localRadius.CertificateAuthentication.IsUnknown() {
		var certificateAuthentication CertificateAuthentication
		err := localRadius.CertificateAuthentication.As(context.Background(), &certificateAuthentication, basetypes.ObjectAsOptions{})
		if err.HasError() {
			diags = append(diags, err.Errors()...)
		}
		var clientRootCaCertificate openApiClient.UpdateNetworkWirelessSsidRequestLocalRadiusCertificateAuthenticationClientRootCaCertificate
		if !certificateAuthentication.ClientRootCaCertificate.IsNull() {
			var clientRootCaCert CaCertificate
			err := certificateAuthentication.ClientRootCaCertificate.As(context.Background(), &clientRootCaCert, basetypes.ObjectAsOptions{})
			if err.HasError() {
				diags = append(diags, err.Errors()...)
			}
			clientRootCaCertificate.SetContents(clientRootCaCert.Contents.ValueString())
		}
		var certAuth openApiClient.UpdateNetworkWirelessSsidRequestLocalRadiusCertificateAuthentication
		certAuth.SetEnabled(certificateAuthentication.Enabled.ValueBool())
		certAuth.SetUseLdap(certificateAuthentication.UseLdap.ValueBool())
		certAuth.SetUseOcsp(certificateAuthentication.UseOcsp.ValueBool())
		certAuth.SetOcspResponderUrl(certificateAuthentication.OcspResponderUrl.ValueString())
		certAuth.SetClientRootCaCertificate(clientRootCaCertificate)
		result.SetCertificateAuthentication(certAuth)
	}

	return result, diags
}

func NetworksWirelessSsidPayloadActiveDirectory(input types.Object) (*openApiClient.UpdateNetworkWirelessSsidRequestActiveDirectory, diag.Diagnostics) {
	if input.IsNull() || input.IsUnknown() {
		return nil, nil
	}

	var diags diag.Diagnostics

	// Unmarshalling input to ActiveDirectory struct
	var activeDirectoryObject ActiveDirectory
	err := input.As(context.Background(), &activeDirectoryObject, basetypes.ObjectAsOptions{})
	if err.HasError() {
		diags = append(diags, err.Errors()...)
		return nil, diags
	}

	// Processing servers
	var activeDirectoryServers []ActiveDirectoryServer
	err = activeDirectoryObject.Servers.ElementsAs(context.Background(), &activeDirectoryServers, true)
	if err != nil {
		diags = append(diags, diag.NewErrorDiagnostic("Error converting ActiveDirectory Servers", fmt.Sprintf("%s", err.Errors())))
	}

	var activeDirectoryServersArray []openApiClient.UpdateNetworkWirelessSsidRequestActiveDirectoryServersInner
	for _, svr := range activeDirectoryServers {
		var activeDirectoryServer openApiClient.UpdateNetworkWirelessSsidRequestActiveDirectoryServersInner
		activeDirectoryServer.SetHost(svr.Host.ValueString())
		activeDirectoryServer.SetPort(int32(svr.Port.ValueInt64()))
		activeDirectoryServersArray = append(activeDirectoryServersArray, activeDirectoryServer)
	}

	// Processing credentials
	var credentialsObject AdCredentials
	err = activeDirectoryObject.Credentials.As(context.Background(), &credentialsObject, basetypes.ObjectAsOptions{})
	if err.HasError() {
		diags = append(diags, err.Errors()...)
	}

	return &openApiClient.UpdateNetworkWirelessSsidRequestActiveDirectory{
		Servers: activeDirectoryServersArray,
		Credentials: &openApiClient.UpdateNetworkWirelessSsidRequestActiveDirectoryCredentials{
			LogonName: credentialsObject.LoginName.ValueStringPointer(),
			Password:  credentialsObject.Password.ValueStringPointer(),
		},
	}, diags
}

func NetworksWirelessSsidPayloadRadiusServers(input types.List) ([]openApiClient.UpdateNetworkWirelessSsidRequestRadiusServersInner, diag.Diagnostics) {
	if input.IsNull() || input.IsUnknown() {
		return nil, nil
	}

	var diags diag.Diagnostics
	var serversList []openApiClient.UpdateNetworkWirelessSsidRequestRadiusServersInner
	var servers []RadiusServer

	err := input.ElementsAs(context.Background(), &servers, true)
	if err.HasError() {
		diags = append(diags, err.Errors()...)
		return nil, diags
	}

	for _, server := range servers {

		port, err := utils.Int32Pointer(server.Port.ValueInt64())
		if err != nil {
			diags = append(diags, diag.NewErrorDiagnostic("Error converting Port", fmt.Sprintf("%s", err.Errors())))
		}

		openRoamingCertificateId, err := utils.Int32Pointer(server.OpenRoamingCertificateID.ValueInt64())
		if err != nil {
			diags = append(diags, diag.NewErrorDiagnostic("Error converting OpenRoamingCertificateID", fmt.Sprintf("%s", err.Errors())))
		}

		serversList = append(serversList, openApiClient.UpdateNetworkWirelessSsidRequestRadiusServersInner{
			Host:                     server.Host.ValueString(),
			Port:                     port,
			Secret:                   server.Secret.ValueStringPointer(),
			CaCertificate:            server.CaCertificate.ValueStringPointer(),
			OpenRoamingCertificateId: openRoamingCertificateId,
			RadsecEnabled:            server.RadSecEnabled.ValueBoolPointer(),
		})
	}

	return serversList, diags
}

func NetworksWirelessSsidPayloadRadiusAccountingServers(input types.List) ([]openApiClient.UpdateNetworkWirelessSsidRequestRadiusAccountingServersInner, diag.Diagnostics) {
	var diags diag.Diagnostics
	var servers []openApiClient.UpdateNetworkWirelessSsidRequestRadiusAccountingServersInner
	var radiusServers []RadiusServer

	err := input.ElementsAs(context.Background(), &radiusServers, true)
	if err.HasError() {
		diags = append(diags, err.Errors()...)
		return nil, diags
	}

	for _, server := range radiusServers {
		port, err := utils.Int32Pointer(server.Port.ValueInt64())
		if err != nil {
			diags = append(diags, diag.NewErrorDiagnostic("Error converting Port", fmt.Sprintf("%s", err.Errors())))
		}
		servers = append(servers, openApiClient.UpdateNetworkWirelessSsidRequestRadiusAccountingServersInner{
			Host:          server.Host.ValueString(),
			Port:          port,
			Secret:        server.Secret.ValueStringPointer(),
			CaCertificate: server.CaCertificate.ValueStringPointer(),
			RadsecEnabled: server.RadSecEnabled.ValueBoolPointer(),
			//OpenRoamingCertificateID: server.OpenRoamingCertificateID.ValueInt64Pointer(),
		})
	}

	return servers, diags
}

func NetworksWirelessSsidPayloadApTagsAndVlanIds(input types.List) ([]openApiClient.UpdateNetworkWirelessSsidRequestApTagsAndVlanIdsInner, diag.Diagnostics) {
	if input.IsNull() || input.IsUnknown() {
		return nil, nil
	}

	var diags diag.Diagnostics
	var tagsAndVlans []openApiClient.UpdateNetworkWirelessSsidRequestApTagsAndVlanIdsInner
	var tagsAndVlansList []ApTagsAndVlanID

	err := input.ElementsAs(context.Background(), &tagsAndVlansList, true)
	if err.HasError() {
		diags = append(diags, err.Errors()...)
		return nil, diags
	}

	for _, tagAndVlan := range tagsAndVlansList {
		vlanId, err := utils.Int32Pointer(tagAndVlan.VlanId.ValueInt64())
		if err != nil {
			diags = append(diags, diag.NewErrorDiagnostic("Error converting VlanID", fmt.Sprintf("%s", err.Errors())))
		}

		var tags []string
		for _, tag := range tagAndVlan.Tags.Elements() {
			tags = append(tags, tag.String())
		}

		tagsAndVlans = append(tagsAndVlans, openApiClient.UpdateNetworkWirelessSsidRequestApTagsAndVlanIdsInner{
			Tags:   tags,
			VlanId: vlanId,
		})
	}

	return tagsAndVlans, diags
}

func NetworksWirelessSsidPayloadGre(input types.Object) (*openApiClient.UpdateNetworkWirelessSsidRequestGre, diag.Diagnostics) {
	if input.IsNull() || input.IsUnknown() {
		return nil, nil
	}

	var diags diag.Diagnostics

	var gre GRE

	err := input.As(context.Background(), &gre, basetypes.ObjectAsOptions{})
	if err.HasError() {
		diags = append(diags, err.Errors()...)
		return nil, diags
	}

	var concentrator GreConcentrator

	err = gre.Concentrator.As(context.Background(), &concentrator, basetypes.ObjectAsOptions{})
	if err.HasError() {
		diags = append(diags, err.Errors()...)
		return nil, diags
	}

	key := int32(gre.Key.ValueInt64())

	return &openApiClient.UpdateNetworkWirelessSsidRequestGre{
		Key: &key,
		Concentrator: &openApiClient.UpdateNetworkWirelessSsidRequestGreConcentrator{
			Host: concentrator.Host.ValueString(),
		},
	}, diags
}

func NetworksWirelessSsidPayloadDnsRewrite(input types.Object) (*openApiClient.UpdateNetworkWirelessSsidRequestDnsRewrite, diag.Diagnostics) {
	if input.IsNull() || input.IsUnknown() {
		return nil, nil
	}

	var diags diag.Diagnostics
	var dnsRewriteObject DnsRewrite

	err := input.As(context.Background(), &dnsRewriteObject, basetypes.ObjectAsOptions{})
	if err.HasError() {
		diags = append(diags, err.Errors()...)
		return nil, diags
	}

	var dnsCustomNameservers []string
	dnsCustomNameserversList := dnsRewriteObject.DnsCustomNameservers.Elements()

	for _, dns := range dnsCustomNameserversList {
		dnsCustomNameservers = append(dnsCustomNameservers, dns.String())
	}

	return &openApiClient.UpdateNetworkWirelessSsidRequestDnsRewrite{
		Enabled:              dnsRewriteObject.Enabled.ValueBoolPointer(),
		DnsCustomNameservers: dnsCustomNameservers,
	}, diags
}

func NetworksWirelessSsidPayloadSpeedBurst(input types.Object) (*openApiClient.UpdateNetworkWirelessSsidRequestSpeedBurst, diag.Diagnostics) {
	if input.IsNull() || input.IsUnknown() {
		return nil, nil
	}

	var diags diag.Diagnostics
	var speedBurst SpeedBurst

	err := input.As(context.Background(), &speedBurst, basetypes.ObjectAsOptions{})
	if err.HasError() {
		diags = append(diags, err.Errors()...)
		return nil, diags
	}

	return &openApiClient.UpdateNetworkWirelessSsidRequestSpeedBurst{
		Enabled: speedBurst.Enabled.ValueBoolPointer(),
	}, diags
}

func NetworksWirelessSsidPayloadNamedVlans(input types.Object) (*openApiClient.UpdateNetworkWirelessSsidRequestNamedVlans, diag.Diagnostics) {
	if input.IsNull() || input.IsUnknown() {
		return nil, nil
	}

	var diags diag.Diagnostics
	var namedVlansObject NamedVlans

	err := input.As(context.Background(), &namedVlansObject, basetypes.ObjectAsOptions{})
	if err.HasError() {
		diags = append(diags, err.Errors()...)
		return nil, diags
	}

	// Tagging
	var tagging Tagging
	err = namedVlansObject.Tagging.As(context.Background(), &tagging, basetypes.ObjectAsOptions{})
	if err.HasError() {
		diags = append(diags, err.Errors()...)
	}

	// Radius
	var radius Radius
	err = namedVlansObject.Radius.As(context.Background(), &radius, basetypes.ObjectAsOptions{})
	if err.HasError() {
		diags = append(diags, err.Errors()...)
	}

	byApTags, err := NetworksWirelessSsidPayloadByApTags(tagging.ByApTags)
	if err.HasError() {
		diags = append(diags, err.Errors()...)
	}

	// guestVlan
	guestVlan, err := NetworksWirelessSsidPayloadRadiusGuestVlan(radius.GuestVlan)
	if err.HasError() {
		diags = append(diags, err.Errors()...)
	}

	return &openApiClient.UpdateNetworkWirelessSsidRequestNamedVlans{
		Tagging: &openApiClient.UpdateNetworkWirelessSsidRequestNamedVlansTagging{
			Enabled:         tagging.Enabled.ValueBoolPointer(),
			DefaultVlanName: tagging.DefaultVlanName.ValueStringPointer(),
			ByApTags:        byApTags,
		},
		Radius: &openApiClient.UpdateNetworkWirelessSsidRequestNamedVlansRadius{
			GuestVlan: &guestVlan,
		},
	}, diags
}

func NetworksWirelessSsidPayloadLdap(input types.Object) (*openApiClient.UpdateNetworkWirelessSsidRequestLdap, diag.Diagnostics) {
	if input.IsNull() || input.IsUnknown() {
		return nil, nil
	}

	var diags diag.Diagnostics

	// Unmarshalling input to LDAP struct
	var ldapObject LDAP
	err := input.As(context.Background(), &ldapObject, basetypes.ObjectAsOptions{})
	if err.HasError() {
		diags = append(diags, err.Errors()...)
		return nil, diags
	}

	// Processing servers
	var servers []LdapServer
	err = ldapObject.Servers.ElementsAs(context.Background(), &servers, true)
	if err != nil {
		diags = append(diags, diag.NewErrorDiagnostic("Error converting Servers", fmt.Sprintf("%s", err.Errors())))
	}

	var serversArray []openApiClient.UpdateNetworkWirelessSsidRequestLdapServersInner
	for _, svr := range servers {
		var server openApiClient.UpdateNetworkWirelessSsidRequestLdapServersInner
		server.SetHost(svr.Host.ValueString())
		server.SetPort(int32(svr.Port.ValueInt64()))
		serversArray = append(serversArray, server)
	}

	// Processing credentials
	var creds LdapCredentials
	err = ldapObject.Credentials.As(context.Background(), &creds, basetypes.ObjectAsOptions{})
	if err.HasError() {
		diags = append(diags, err.Errors()...)
	}

	// Processing server CA certificate
	var serverCaCertificate CaCertificate
	err = ldapObject.ServerCaCertificate.As(context.Background(), &serverCaCertificate, basetypes.ObjectAsOptions{})
	if err.HasError() {
		diags = append(diags, err.Errors()...)
	}

	return &openApiClient.UpdateNetworkWirelessSsidRequestLdap{
		Servers: serversArray,
		Credentials: &openApiClient.UpdateNetworkWirelessSsidRequestLdapCredentials{
			DistinguishedName: creds.DistinguishedName.ValueStringPointer(),
			Password:          creds.Password.ValueStringPointer(),
		},
		BaseDistinguishedName: ldapObject.BaseDistinguishedName.ValueStringPointer(),
		ServerCaCertificate: &openApiClient.UpdateNetworkWirelessSsidRequestLdapServerCaCertificate{
			Contents: serverCaCertificate.Contents.ValueStringPointer(),
		},
	}, diags
}

func NetworksWirelessSsidPayloadByApTags(input types.List) ([]openApiClient.UpdateNetworkWirelessSsidRequestNamedVlansTaggingByApTagsInner, diag.Diagnostics) {
	if input.IsNull() || input.IsUnknown() {
		return nil, nil
	}

	var diags diag.Diagnostics
	var byApTags []openApiClient.UpdateNetworkWirelessSsidRequestNamedVlansTaggingByApTagsInner

	var byApTagsList []ByApTag
	err := input.ElementsAs(context.Background(), &byApTagsList, false)
	if err.HasError() {
		diags = append(diags, err.Errors()...)
		return nil, diags
	}

	for _, byApTag := range byApTagsList {
		var tags []string
		for _, tag := range byApTag.Tags.Elements() {
			tags = append(tags, tag.String())
		}

		byApTags = append(byApTags, openApiClient.UpdateNetworkWirelessSsidRequestNamedVlansTaggingByApTagsInner{
			Tags:     tags,
			VlanName: byApTag.VlanName.ValueStringPointer(),
		})
	}
	return byApTags, diags
}

func NetworksWirelessSsidPayloadRadiusGuestVlan(input types.Object) (openApiClient.UpdateNetworkWirelessSsidRequestNamedVlansRadiusGuestVlan, diag.Diagnostics) {
	var diags diag.Diagnostics
	var guestVlans openApiClient.UpdateNetworkWirelessSsidRequestNamedVlansRadiusGuestVlan

	var data RadiusGuestVlan

	err := input.As(context.Background(), data, basetypes.ObjectAsOptions{})
	if err.HasError() {
		return guestVlans, err
	}

	// enabled
	guestVlans.SetEnabled(data.Enabled.ValueBool())

	// name
	guestVlans.SetName(data.Name.ValueString())

	return guestVlans, diags
}

// update terraform state funcs

func networksWirelessSsidAdminSplashUrl(data *openApiClient.GetNetworkWirelessSsids200ResponseInner) (basetypes.StringValue, diag.Diagnostics) {
	var diags diag.Diagnostics

	result := types.StringNull()
	adminSplashUrl, ok := data.GetAdminSplashUrlOk()
	if ok {
		result = types.StringValue(*adminSplashUrl)
		return result, diags
	}

	return result, diags
}

func NetworksWirelessSsidStateRadiusServers(input []openApiClient.GetNetworkWirelessSsids200ResponseInnerRadiusServersInner) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	var radiusServers []attr.Value

	radiusServerAttr := map[string]attr.Type{
		"host":                        types.StringType,
		"port":                        types.Int64Type,
		"secret":                      types.StringType,
		"rad_sec_enabled":             types.BoolType,
		"ca_certificate":              types.StringType,
		"open_roaming_certificate_id": types.Int64Type,
	}

	for _, i := range input {
		radiusServer := RadiusServer{
			Host:                     types.StringValue(i.GetHost()),
			Port:                     types.Int64Value(int64(i.GetPort())),
			Secret:                   types.StringNull(), // secret is not returned for security reasons
			RadSecEnabled:            types.BoolNull(),   // not sure why this isn't returned.
			OpenRoamingCertificateID: types.Int64Value(int64(i.GetOpenRoamingCertificateId())),
			CaCertificate:            types.StringValue(i.GetCaCertificate()),
		}

		radiusServerObject, err := types.ObjectValueFrom(context.Background(), radiusServerAttr, radiusServer)
		if err.HasError() {
			diags.Append(err...)
			continue
		}

		radiusServers = append(radiusServers, radiusServerObject)
	}

	radiusServersList, err := types.ListValueFrom(context.Background(), types.ObjectType{AttrTypes: radiusServerAttr}, radiusServers)
	if err.HasError() {
		diags.Append(err...)
	}

	return radiusServersList, diags
}

func NetworksWirelessSsidStateRadiusAccountingServers(input []openApiClient.GetNetworkWirelessSsids200ResponseInnerRadiusAccountingServersInner) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	var radiusServers []attr.Value

	radiusServerAttr := map[string]attr.Type{
		"host":                        types.StringType,
		"port":                        types.Int64Type,
		"secret":                      types.StringType,
		"rad_sec_enabled":             types.BoolType,
		"ca_certificate":              types.StringType,
		"open_roaming_certificate_id": types.Int64Type,
	}

	for _, i := range input {
		radiusServer := RadiusServer{
			Host:                     types.StringValue(i.GetHost()),
			Port:                     types.Int64Value(int64(i.GetPort())),
			Secret:                   types.StringNull(), // secret is not returned for security reasons
			RadSecEnabled:            types.BoolNull(),   // not sure why this isn't returned.
			OpenRoamingCertificateID: types.Int64Value(int64(i.GetOpenRoamingCertificateId())),
			CaCertificate:            types.StringValue(i.GetCaCertificate()),
		}

		radiusServerObject, err := types.ObjectValueFrom(context.Background(), radiusServerAttr, radiusServer)
		if err.HasError() {
			diags.Append(err...)
			continue
		}

		radiusServers = append(radiusServers, radiusServerObject)
	}

	radiusServersList, err := types.ListValueFrom(context.Background(), types.ObjectType{AttrTypes: radiusServerAttr}, radiusServers)
	if err.HasError() {
		diags.Append(err...)
	}

	return radiusServersList, diags
}

func NetworksWirelessSsidStateDot11w(rawResp map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var dot11w Dot11w

	dot11wAttrs := map[string]attr.Type{
		"enabled":  types.BoolType,
		"required": types.BoolType,
	}

	if d, ok := rawResp["dot11w"].(map[string]interface{}); ok {

		// enabled
		enabled, err := utils.ExtractBoolAttr(d, "enabled")
		if diags.HasError() {
			diags.Append(err...)
		}
		dot11w.Enabled = enabled

		// required
		required, err := utils.ExtractBoolAttr(d, "required")
		if diags.HasError() {
			diags.Append(err...)
		}
		dot11w.Required = required

	} else {
		Dot11wNull := types.ObjectNull(dot11wAttrs)
		return Dot11wNull, diags
	}

	dot11wObj, err := types.ObjectValueFrom(context.Background(), dot11wAttrs, dot11w)
	if err.HasError() {
		diags.Append(err...)
	}

	return dot11wObj, diags
}

func NetworksWirelessSsidStateDot11r(rawResp map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var dot11r Dot11r

	dot11rAttrs := map[string]attr.Type{
		"enabled":  types.BoolType,
		"adaptive": types.BoolType,
	}

	if d, ok := rawResp["dot11r"].(map[string]interface{}); ok {

		// enabled
		enabled, err := utils.ExtractBoolAttr(d, "enabled")
		if diags.HasError() {
			diags.Append(err...)
		}
		dot11r.Enabled = enabled

		// adaptive
		adaptive, err := utils.ExtractBoolAttr(d, "adaptive")
		if diags.HasError() {
			diags.Append(err...)
		}
		dot11r.Adaptive = adaptive

	} else {
		dot11rNull := types.ObjectNull(dot11rAttrs)
		return dot11rNull, diags
	}

	outputObj, err := types.ObjectValueFrom(context.Background(), dot11rAttrs, dot11r)
	if err.HasError() {
		diags.Append(err...)
	}

	return outputObj, diags
}

func NetworksWirelessSsidStateOauth(rawResp map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var oauth OAuth

	oauthAttrs := map[string]attr.Type{
		"allowed_domains": types.ListType{ElemType: types.StringType},
	}

	//oauth
	if oa, ok := rawResp["oauth"].(map[string]interface{}); ok {

		// allowed domains
		if ad, ok := oa["allowed_domains"].([]string); ok {
			var allowedDomains []types.String
			for _, domain := range ad {
				allowedDomains = append(allowedDomains, types.StringValue(domain))
			}

			allowedDomainsObj, err := types.ListValueFrom(context.Background(), types.StringType, allowedDomains)
			if err.HasError() {
				diags.Append(err...)
			}
			oauth.AllowedDomains = allowedDomainsObj

		} else {
			allowedDomainsObjNull := types.ListNull(types.StringType)
			oauth.AllowedDomains = allowedDomainsObjNull
		}

	} else {
		oauthObjNull := types.ObjectNull(oauthAttrs)
		return oauthObjNull, diags
	}

	oauthObj, err := types.ObjectValueFrom(context.Background(), oauthAttrs, oauth)
	if err.HasError() {
		diags.Append(err...)
	}

	return oauthObj, diags
}

func NetworksWirelessSsidStateLocalRadius(rawResp map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	passwordAuthenticationAttrs := map[string]attr.Type{
		"enabled": types.BoolType,
	}

	contentsAttr := map[string]attr.Type{
		"contents": types.StringType,
	}

	certificateAuthenticationAttr := map[string]attr.Type{
		"enabled":                    types.BoolType,
		"use_ldap":                   types.BoolType,
		"use_ocsp":                   types.BoolType,
		"ocsp_responder_url":         types.StringType,
		"client_root_ca_certificate": types.ObjectType{AttrTypes: contentsAttr},
	}

	localRadiusAttrs := map[string]attr.Type{
		"cache_timeout":              types.Int64Type,
		"password_authentication":    types.ObjectType{AttrTypes: passwordAuthenticationAttrs},
		"certificate_authentication": types.ObjectType{AttrTypes: certificateAuthenticationAttr},
	}

	var localRadius LocalRadius

	// cacheTimeout
	cacheTimeout, err := utils.ExtractInt64Attr(rawResp, "cacheTimeOut")
	if diags.HasError() {
		diags.Append(err...)
	}
	localRadius.CacheTimeout = cacheTimeout

	// Password Authentication
	if pa, ok := rawResp["passwordAuthentication"].(map[string]interface{}); ok {
		var passwordAuth PasswordAuthentication

		// enabled
		enabled, err := utils.ExtractBoolAttr(pa, "enabled")
		if diags.HasError() {
			diags.Append(err...)
		}
		passwordAuth.Enabled = enabled

	} else {
		passwordAuthenticationObjNull := types.ObjectNull(passwordAuthenticationAttrs)
		localRadius.PasswordAuthentication = passwordAuthenticationObjNull
	}

	// certificateAuthentication
	if ca, ok := rawResp["certificateAuthentication"].(map[string]interface{}); ok {
		var certificateAuthentication CertificateAuthentication

		//   Enabled
		if _, ok := ca["enabled"].(types.Bool); ok {

			caEnabled, err := utils.ExtractBoolAttr(ca, "enabled")
			if diags.HasError() {
				diags.Append(err...)
			}

			certificateAuthentication.Enabled = caEnabled

		} else {
			caEnabledNull := types.BoolNull()
			certificateAuthentication.Enabled = caEnabledNull
		}

		//    UseLdap
		if _, ok := ca["useLdap"].(types.Bool); ok {

			useLdap, err := utils.ExtractBoolAttr(ca, "useLdap")
			if diags.HasError() {
				diags.Append(err...)
			}

			certificateAuthentication.UseLdap = useLdap

		} else {
			useLdapNull := types.BoolNull()
			certificateAuthentication.UseLdap = useLdapNull
		}

		//    UseOcsp
		if _, ok := ca["useOcsp"].(types.Bool); ok {

			useOcsp, err := utils.ExtractBoolAttr(ca, "useOcsp")
			if diags.HasError() {
				diags.Append(err...)
			}

			certificateAuthentication.UseOcsp = useOcsp

		} else {
			useOcspNull := types.BoolNull()
			certificateAuthentication.UseOcsp = useOcspNull
		}

		//    OcspResponderUrl
		if _, ok := ca["ocspResponderUrl"].(types.String); ok {

			ocspResponderUrl, err := utils.ExtractStringAttr(ca, "ocspResponderUrl")
			if diags.HasError() {
				diags.Append(err...)
			}

			certificateAuthentication.OcspResponderUrl = ocspResponderUrl

		} else {
			ocspResponderUrlNull := types.StringNull()
			certificateAuthentication.OcspResponderUrl = ocspResponderUrlNull
		}

		//    ClientRootCaCertificate
		if crca, ok := rawResp["clientRootCaCertificate"].(map[string]interface{}); ok {
			var clientRootCaCertificate CaCertificate

			// Contents
			if _, ok := crca["contents"].(types.String); ok {

				contents, err := utils.ExtractStringAttr(ca, "contents")
				if diags.HasError() {
					diags.Append(err...)
				}

				clientRootCaCertificate.Contents = contents

			} else {
				contentsNull := types.StringNull()
				clientRootCaCertificate.Contents = contentsNull
			}
		}

		certificateAuthenticationObj, err := types.ObjectValueFrom(context.Background(), certificateAuthenticationAttr, certificateAuthentication)
		if err.HasError() {
			tflog.Error(context.Background(), fmt.Sprintf("%s", err.Errors()))
		}

		localRadius.CertificateAuthentication = certificateAuthenticationObj

	} else {
		certificateAuthenticationObjNull := types.ObjectNull(certificateAuthenticationAttr)
		localRadius.CertificateAuthentication = certificateAuthenticationObjNull
	}

	outputObj, err := types.ObjectValueFrom(context.Background(), localRadiusAttrs, localRadius)
	if err.HasError() {
		diags.Append(err...)
	}

	return outputObj, diags
}

func NetworksWirelessSsidStateLdap(httpResp map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var ldap LDAP

	serverAttr := map[string]attr.Type{
		"host": types.StringType,
		"port": types.Int64Type,
	}

	credentialsAttr := map[string]attr.Type{
		"distinguished_name": types.StringType,
		"password":           types.StringType,
	}

	contentsAttr := map[string]attr.Type{
		"contents": types.StringType,
	}

	ldapAttrs := map[string]attr.Type{
		"base_distinguished_name": types.StringType,
		"servers":                 types.ListType{ElemType: types.ObjectType{AttrTypes: serverAttr}},
		"credentials":             types.ObjectType{AttrTypes: credentialsAttr},
		"server_ca_certificate":   types.ObjectType{AttrTypes: contentsAttr},
	}

	if l, ok := httpResp["ldap"].(map[string]interface{}); ok {

		// baseDistinguishedName
		baseDistinguishedName, err := utils.ExtractStringAttr(l, "baseDistinguishedName")
		if err.HasError() {
			diags.AddError("baseDistinguishedName Attribute", fmt.Sprintf("%s", err.Errors()))
		}
		ldap.BaseDistinguishedName = baseDistinguishedName

		// credentials
		if credsMap, ok := l["credentials"].(map[string]interface{}); ok {
			var creds LdapCredentials

			// loginName
			DistinguishedNameObj, err := utils.ExtractStringAttr(credsMap, "DistinguishedName")
			if err.HasError() {
				diags.AddError("DistinguishedName Attribute", fmt.Sprintf("%s", err.Errors()))
			}
			creds.DistinguishedName = DistinguishedNameObj

			// Password
			passwordObj, err := utils.ExtractStringAttr(credsMap, "password")
			if err.HasError() {
				diags.AddError("password Attribute", fmt.Sprintf("%s", err.Errors()))
			}
			creds.Password = passwordObj

			credsObj, err := types.ObjectValueFrom(context.Background(), credentialsAttr, creds)
			if err.HasError() {
				diags.AddError("credentials object Attribute", fmt.Sprintf("%s", err.Errors()))
			}

			ldap.Credentials = credsObj

		} else {
			credsNull := types.ObjectNull(credentialsAttr)
			ldap.Credentials = credsNull
		}

		// serverCaCertificate
		if serverCaCertificateMap, ok := l["serverCaCertificate"].(map[string]interface{}); ok {
			var serverCaCertificate LdapServerCaCertificate

			// contents
			contents, err := utils.ExtractStringAttr(serverCaCertificateMap, "contents")
			if err.HasError() {
				diags.AddError("contents Attribute", fmt.Sprintf("%s", err.Errors()))
			}

			serverCaCertificate.Contents = contents

			ServerCaCertObj, err := types.ObjectValueFrom(context.Background(), contentsAttr, serverCaCertificate)
			if err.HasError() {
				diags.AddError("serverCaCertificate object Attribute", fmt.Sprintf("%s", err.Errors()))
			}

			ldap.ServerCaCertificate = ServerCaCertObj

		} else {
			ServerCaCertObjNull := types.ObjectNull(contentsAttr)
			ldap.ServerCaCertificate = ServerCaCertObjNull
		}

		// servers
		if listMapArray, ok := l["servers"].([]map[string]interface{}); ok {

			var serversArray []types.Object

			for _, listMap := range listMapArray {
				var server ActiveDirectoryServer

				// host
				host, err := utils.ExtractStringAttr(listMap, "host")
				if err.HasError() {
					diags.AddError("host Attribute", fmt.Sprintf("%s", err.Errors()))
				}
				server.Host = host

				// port
				port, err := utils.ExtractInt64Attr(listMap, "port")
				if err.HasError() {
					diags.AddError("port Attribute", fmt.Sprintf("%s", err.Errors()))
				}
				server.Port = port

				serverObj, err := types.ObjectValueFrom(context.Background(), serverAttr, server)
				if err.HasError() {
					diags.AddError("server Attribute", fmt.Sprintf("%s", err.Errors()))
				}

				serversArray = append(serversArray, serverObj)
			}

			// servers Array
			serversArrayObj, err := types.ListValueFrom(context.Background(), types.ObjectType{AttrTypes: serverAttr}, serversArray)
			if err.HasError() {
				diags.AddError("servers array Attribute", fmt.Sprintf("%s", err.Errors()))
			}

			ldap.Servers = serversArrayObj

		} else {
			lisObjNull := types.ListNull(types.ObjectType{AttrTypes: serverAttr})
			ldap.Servers = lisObjNull

		}

	} else {
		ldapObjNull := types.ObjectNull(ldapAttrs)
		return ldapObjNull, diags
	}

	ldapObject, err := types.ObjectValueFrom(context.Background(), ldapAttrs, ldap)
	if err.HasError() {
		diags.AddError("ldap object Attribute", fmt.Sprintf("%s", err.Errors()))
	}

	return ldapObject, diags
}

func NetworksWirelessSsidStateActiveDirectory(httpResp map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var activeDirectory ActiveDirectory

	serverAttr := map[string]attr.Type{
		"host": types.StringType,
		"port": types.Int64Type,
	}

	credentialsAttr := map[string]attr.Type{
		"login_name": types.StringType,
		"password":   types.StringType,
	}

	activeDirectoryAttrs := map[string]attr.Type{
		"servers": types.ListType{
			ElemType: types.ObjectType{AttrTypes: serverAttr},
		},
		"credentials": types.ObjectType{
			AttrTypes: credentialsAttr,
		},
	}

	if ad, ok := httpResp["activeDirectory"].(map[string]interface{}); ok {

		// servers
		if listMapArray, ok := ad["servers"].([]map[string]interface{}); ok {

			var serversArray []types.Object

			for _, listMap := range listMapArray {
				var server ActiveDirectoryServer

				// host
				host, err := utils.ExtractStringAttr(listMap, "host")
				if err.HasError() {
					diags.AddError("host Attribute", fmt.Sprintf("%s", err.Errors()))
				}
				server.Host = host

				// port
				port, err := utils.ExtractInt64Attr(listMap, "port")
				if err.HasError() {
					diags.AddError("port Attribute", fmt.Sprintf("%s", err.Errors()))
				}
				server.Port = port

				serverObj, err := types.ObjectValueFrom(context.Background(), serverAttr, server)
				if err.HasError() {
					diags.AddError("server Attribute", fmt.Sprintf("%s", err.Errors()))
				}

				serversArray = append(serversArray, serverObj)
			}

			// servers Array
			serversArrayObj, err := types.ListValueFrom(context.Background(), types.ObjectType{AttrTypes: serverAttr}, serversArray)
			if err.HasError() {
				diags.AddError("servers array Attribute", fmt.Sprintf("%s", err.Errors()))
			}

			activeDirectory.Servers = serversArrayObj

		} else {
			lisObjNull := types.ListNull(types.ObjectType{AttrTypes: serverAttr})
			activeDirectory.Servers = lisObjNull

		}

		// credentials
		if credsMap, ok := ad["credentials"].(map[string]interface{}); ok {
			var creds AdCredentials

			// loginName
			loginNameObj, err := utils.ExtractStringAttr(credsMap, "loginName")
			if err.HasError() {
				diags.AddError("loginName Attribute", fmt.Sprintf("%s", err.Errors()))
			}
			creds.LoginName = loginNameObj

			// Password
			passwordObj, err := utils.ExtractStringAttr(credsMap, "password")
			if err.HasError() {
				diags.AddError("password Attribute", fmt.Sprintf("%s", err.Errors()))
			}
			creds.Password = passwordObj

			credsObj, err := types.ObjectValueFrom(context.Background(), credentialsAttr, creds)
			if err.HasError() {
				diags.AddError("credentials object Attribute", fmt.Sprintf("%s", err.Errors()))
			}

			activeDirectory.Credentials = credsObj

		} else {
			credsNull := types.ObjectNull(credentialsAttr)
			activeDirectory.Credentials = credsNull
		}

	} else {
		activeDirectoryObjNull := types.ObjectNull(activeDirectoryAttrs)
		return activeDirectoryObjNull, diags
	}

	activeDirectoryObj, err := types.ObjectValueFrom(context.Background(), activeDirectoryAttrs, activeDirectory)
	if err.HasError() {
		diags.AddError("Active Directory Attr", fmt.Sprintf("%s", err.Errors()))
	}

	return activeDirectoryObj, diags
}

func NetworksWirelessSsidStateApTagsAndVlanIds(httpResp map[string]interface{}) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	apTagsAndVlanIdsAttr := map[string]attr.Type{
		"tags":    types.ListType{ElemType: types.StringType},
		"vlan_id": types.Int64Type,
	}

	apTagsAndVlanIdsAttrs := types.ObjectType{AttrTypes: apTagsAndVlanIdsAttr}

	apTagsAndVlanIdsList, err := utils.ExtractListAttr(httpResp, "apTagsAndVlanIds", apTagsAndVlanIdsAttrs)
	if err.HasError() {
		tflog.Error(context.Background(), fmt.Sprintf("%s", err.Errors()))
	}

	return apTagsAndVlanIdsList, diags
}

func NetworksWirelessSsidStateGre(httpResp map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var gre GRE

	concentratorAttrs := map[string]attr.Type{
		"host": types.StringType,
	}
	greAttrs := map[string]attr.Type{
		"concentrator": types.ObjectType{AttrTypes: concentratorAttrs},
		"key":          types.Int64Type,
	}

	if g, ok := httpResp["gre"].(map[string]interface{}); ok {

		// key
		gre.Key, diags = utils.ExtractInt64Attr(httpResp, "key")

		// concentrator
		if c, ok := g["concentrator"].(map[string]interface{}); ok {
			var concentrator GreConcentrator

			concentrator.Host, diags = utils.ExtractStringAttr(c, "host")

			concentratorObj, err := types.ObjectValueFrom(context.Background(), concentratorAttrs, concentrator)
			if err.HasError() {
				diags.Append(err...)
			}

			gre.Concentrator = concentratorObj
		} else {
			concentratorObjNull := types.ObjectNull(concentratorAttrs)
			gre.Concentrator = concentratorObjNull
		}

	} else {
		greObjNull := types.ObjectNull(greAttrs)
		return greObjNull, diags
	}

	greObj, err := types.ObjectValueFrom(context.Background(), greAttrs, gre)
	if err.HasError() {
		diags.Append(err...)
	}

	return greObj, diags
}

func NetworksWirelessSsidStateDnsRewrite(httpResp map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var dnsRewrite DnsRewrite

	dnsCustomNameServersAttrs := types.ListType{ElemType: types.StringType}

	dnsRewriteAttrs := map[string]attr.Type{
		"enabled":                 types.BoolType,
		"dns_custom_name_servers": dnsCustomNameServersAttrs,
	}

	dns, ok := httpResp["dnsRewrite"].(map[string]interface{})
	if ok {

		// enabled
		enabled, err := utils.ExtractBoolAttr(dns, "enabled")
		if err.HasError() {
			diags.AddError("enabled Attr", fmt.Sprintf("%s", err.Errors()))
		}

		dnsRewrite.Enabled = enabled

		// dns custom Name Servers
		dnsCustomNameservers, err := utils.ExtractListStringAttr(dns, "dnsCustomNameServers")
		if err.HasError() {
			diags.AddError("dnsCustomNameservers Attr", fmt.Sprintf("%s", err.Errors()))
		}

		dnsRewrite.DnsCustomNameservers = dnsCustomNameservers

	} else {
		dnsRewriteObjNull := types.ObjectNull(dnsRewriteAttrs)
		return dnsRewriteObjNull, diags
	}

	// dnsRewrite Terraform types Object
	dnsRewriteObj, dnsRewriteDiags := types.ObjectValueFrom(context.Background(), dnsRewriteAttrs, dnsRewrite)
	if dnsRewriteDiags.HasError() {
		diags.AddError("dnsRewriteObject Attr", fmt.Sprintf("%s", dnsRewriteDiags.Errors()))
	}

	return dnsRewriteObj, diags
}

func NetworksWirelessSsidStateSpeedBurst(httpResp map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var speedBurst SpeedBurst

	speedBurstAttrs := map[string]attr.Type{
		"enabled": types.BoolType,
	}

	sb, ok := httpResp["speedBurst"].(map[string]interface{})
	if ok {
		enabled, err := utils.ExtractBoolAttr(sb, "enabled")
		if err.HasError() {
			diags.AddError("enabled Attr", fmt.Sprintf("%s", err.Errors()))
		}

		speedBurst.Enabled = enabled
	} else {
		speedBurstObjNull := types.ObjectNull(speedBurstAttrs)
		return speedBurstObjNull, diags
	}

	speedBurstObj, speedBurstDiags := types.ObjectValueFrom(context.Background(), speedBurstAttrs, speedBurst)
	if speedBurstDiags.HasError() {
		diags.AddError("enabled Attr", fmt.Sprintf("%s", speedBurstDiags.Errors()))
	}

	return speedBurstObj, diags
}

func NetworksWirelessSsidStateNamedVlans(httpResp map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var namedVlans NamedVlans

	byApTagsAttrs := map[string]attr.Type{
		"tags":      types.ListType{ElemType: types.StringType},
		"vlan_name": types.StringType,
	}

	taggingAttrs := map[string]attr.Type{
		"enabled":           types.BoolType,
		"default_vlan_name": types.StringType,
		"by_ap_tags":        types.ListType{ElemType: types.ObjectType{AttrTypes: byApTagsAttrs}},
	}

	guestVlanAttrs := map[string]attr.Type{
		"enabled": types.BoolType,
		"name":    types.StringType,
	}

	radiusAttrs := map[string]attr.Type{
		"guest_vlan": types.ObjectType{AttrTypes: guestVlanAttrs},
	}

	namedVlansAttrs := map[string]attr.Type{
		"tagging": types.ObjectType{AttrTypes: taggingAttrs},
		"radius":  types.ObjectType{AttrTypes: radiusAttrs},
	}

	nv, ok := httpResp["namedVlans"].(map[string]interface{})
	if ok {

		// tagging
		t, ok := nv["tagging"].(map[string]interface{})
		if ok {
			var tagging Tagging

			// Enabled
			enabled, err := utils.ExtractBoolAttr(t, "enabled")
			if err.HasError() {
				diags.AddError("enabled Attr", fmt.Sprintf("%s", err.Errors()))
			}
			tagging.Enabled = enabled

			// DefaultVlanName
			defaultVlanName, err := utils.ExtractStringAttr(t, "defaultVlanName")
			if err.HasError() {
				diags.AddError("defaultVlanName Attr", fmt.Sprintf("%s", err.Errors()))
			}
			tagging.DefaultVlanName = defaultVlanName

			// ByApTags
			b, ok := nv["byApTags"].(map[string]interface{})
			if ok {
				var byApTags ByApTag

				// tags
				tags, err := utils.ExtractListStringAttr(b, "tags")
				if err.HasError() {
					diags.AddError("tags Attr", fmt.Sprintf("%s", err.Errors()))
				}
				byApTags.Tags = tags

				// vlanName
				vlanName, err := utils.ExtractStringAttr(b, "vlanName")
				if err.HasError() {
					diags.AddError("vlanName Attr", fmt.Sprintf("%s", err.Errors()))
				}
				byApTags.VlanName = vlanName

				byApTagsObj, err := types.ObjectValueFrom(context.Background(), byApTagsAttrs, byApTags)
				if err.HasError() {
					diags.AddError("byApTags Object Attr", fmt.Sprintf("%s", err.Errors()))
				}

				byApTagsArray, err := types.ListValueFrom(context.Background(), types.ObjectType{AttrTypes: byApTagsAttrs}, byApTagsObj)
				if err.HasError() {
					diags.AddError("byApTags Array Attr", fmt.Sprintf("%s", err.Errors()))
				}

				tagging.ByApTags = byApTagsArray
			} else {
				byApTagsArrayNull := types.ListNull(types.ObjectType{AttrTypes: byApTagsAttrs})
				tagging.ByApTags = byApTagsArrayNull
			}

			taggingObj, err := types.ObjectValueFrom(context.Background(), taggingAttrs, tagging)
			if err.HasError() {
				diags.AddError("tagging Object Attr", fmt.Sprintf("%s", err.Errors()))
			}
			namedVlans.Tagging = taggingObj

		} else {
			taggingObjNull := types.ObjectNull(taggingAttrs)
			namedVlans.Tagging = taggingObjNull
		}

		// radius
		r, ok := nv["radius"].(map[string]interface{})
		if ok {
			var radius Radius

			g, ok := r["guestVlan"].(map[string]interface{})
			if ok {
				var guestVlans RadiusGuestVlan

				// enabled
				enabled, err := utils.ExtractBoolAttr(g, "enabled")
				if err.HasError() {
					diags.AddError("enabled Attr", fmt.Sprintf("%s", err.Errors()))
				}
				guestVlans.Enabled = enabled

				// name
				name, err := utils.ExtractStringAttr(g, "name")
				if err.HasError() {
					diags.AddError("name Attr", fmt.Sprintf("%s", err.Errors()))
				}
				guestVlans.Name = name

				guestVlansObj, err := types.ObjectValueFrom(context.Background(), guestVlanAttrs, guestVlans)
				if err.HasError() {
					diags.AddError("guestVlans Object Attr", fmt.Sprintf("%s", err.Errors()))
				}
				radius.GuestVlan = guestVlansObj

			} else {
				guestVlansObjNull := types.ObjectNull(guestVlanAttrs)
				radius.GuestVlan = guestVlansObjNull
			}

			radiusObj, err := types.ObjectValueFrom(context.Background(), radiusAttrs, radius)
			if err.HasError() {
				diags.AddError("radius object Attr", fmt.Sprintf("%s", err.Errors()))
			}
			namedVlans.Radius = radiusObj

		} else {
			radiusObjNull := types.ObjectNull(radiusAttrs)
			namedVlans.Radius = radiusObjNull
		}

	} else {
		namedVlansObjNull := types.ObjectNull(namedVlansAttrs)
		return namedVlansObjNull, diags
	}

	namedVlansObj, err := types.ObjectValueFrom(context.Background(), radiusAttrs, namedVlans)
	if err.HasError() {
		diags.AddError("namedVlans object Attr", fmt.Sprintf("%s", err.Errors()))
	}

	return namedVlansObj, diags

}

// updateNetworksWirelessSsidsResourceState updates the resource state with the provided api data.
func updateNetworksWirelessSsidsResourceState(ctx context.Context, state *NetworksWirelessSsidResourceModel, data *openApiClient.GetNetworkWirelessSsids200ResponseInner, httpResp *http.Response) diag.Diagnostics {
	var diags diag.Diagnostics

	rawResp, err := tools.ExtractResponseToMap(httpResp)
	if err != nil {
		diags.AddError("Failed to Unmarshal HttpResp", err.Error())
	}

	// Number
	if state.Number.IsNull() || state.Number.IsUnknown() {
		number := int64(*data.Number)
		state.Number = types.Int64Value(number)
	}

	// Import ID
	if !state.NetworkId.IsNull() && !state.NetworkId.IsUnknown() && !state.Number.IsNull() && !state.Number.IsUnknown() {
		id := state.NetworkId.ValueString() + "," + strconv.FormatInt(state.Number.ValueInt64(), 10)
		state.Id = types.StringValue(id)
	} else {
		state.Id = types.StringNull()
	}

	// Name
	if state.Name.IsNull() || state.Name.IsUnknown() {
		state.Name, diags = utils.ExtractStringAttr(rawResp, "name")
		if diags.HasError() {
			diags.AddError("Name Attribute", state.Name.ValueString())
			return diags
		}

	}

	// Enabled
	if state.Enabled.IsNull() || state.Enabled.IsUnknown() {
		state.Enabled, diags = utils.ExtractBoolAttr(rawResp, "enabled")
		if diags.HasError() {
			diags.AddError("Enabled Attribute", "")
			return diags
		}
	}

	// SplashPage
	if state.SplashPage.IsNull() || state.SplashPage.IsUnknown() {
		state.SplashPage, diags = utils.ExtractStringAttr(rawResp, "splashPage")

		if diags.HasError() {
			tflog.Error(ctx, "SplashPage Attribute")
			return diags
		}

	}

	// SsidAdminAccessible
	if state.SsidAdminAccessible.IsNull() || state.SsidAdminAccessible.IsUnknown() {

		state.SsidAdminAccessible, diags = utils.ExtractBoolAttr(rawResp, "ssidAdminAccessible")
		if diags.HasError() {
			diags.AddError("SSIDAdminAccessible Attribute", "")
			return diags
		}
	}

	// LocalAuth
	if state.LocalAuth.IsNull() || state.LocalAuth.IsUnknown() {

		state.LocalAuth, diags = utils.ExtractBoolAttr(rawResp, "localAuth")
		if diags.HasError() {
			diags.AddError("LocalAuth Attribute", "")
			return diags
		}
	}

	// AuthMode
	if state.AuthMode.IsNull() || state.AuthMode.IsUnknown() {

		state.AuthMode, diags = utils.ExtractStringAttr(rawResp, "authMode")
		if diags.HasError() {
			diags.AddError("AuthMode Attribute", "")
			return diags
		}

	}

	// EncryptionMode
	if state.EncryptionMode.IsNull() || state.EncryptionMode.IsUnknown() {

		state.EncryptionMode, diags = utils.ExtractStringAttr(rawResp, "encryptionMode")
		if diags.HasError() {
			diags.AddError("Encryption Attribute", "")
			return diags
		}

	}

	// WPAEncryptionMode
	if state.WPAEncryptionMode.IsNull() || state.WPAEncryptionMode.IsUnknown() {

		state.WPAEncryptionMode, diags = utils.ExtractStringAttr(rawResp, "wpaEncryptionMode")
		if diags.HasError() {
			diags.AddError("WPAEncryptionMode Attribute", "")
			return diags
		}
	}

	// RadiusServers
	if state.RadiusServers.IsNull() || state.RadiusServers.IsUnknown() {

		state.RadiusServers, diags = NetworksWirelessSsidStateRadiusServers(data.RadiusServers)
		if diags.HasError() {
			diags.AddError("Radius Servers Attribute", "")
			return diags
		}

	}

	// RadiusAccountingServers
	if state.RadiusAccountingServers.IsNull() || state.RadiusAccountingServers.IsUnknown() {

		radiusAccountingServers, diags := NetworksWirelessSsidStateRadiusAccountingServers(data.RadiusAccountingServers)
		if diags.HasError() {
			diags.AddError("Radius Accounting Servers Attribute", "")
			return diags
		}
		state.RadiusAccountingServers = radiusAccountingServers
	}

	// RadiusAccountingEnabled
	if state.RadiusAccountingEnabled.IsNull() || state.RadiusAccountingEnabled.IsUnknown() {
		state.RadiusAccountingEnabled, diags = utils.ExtractBoolAttr(rawResp, "radiusAccountingEnabled")
		if diags.HasError() {
			diags.AddError("Radius Accounting Enabled Attribute", "")
			return diags
		}
	}

	// RadiusEnabled
	if state.RadiusEnabled.IsNull() || state.RadiusEnabled.IsUnknown() {

		state.RadiusEnabled, diags = utils.ExtractBoolAttr(rawResp, "radiusEnabled")
		if diags.HasError() {
			diags.AddError("Radius Enabled Attribute", "")
			return diags
		}
	}

	// RadiusAttributeForGroupPolicies
	if state.RadiusAttributeForGroupPolicies.IsNull() || state.RadiusAttributeForGroupPolicies.IsUnknown() {

		state.RadiusAttributeForGroupPolicies, diags = utils.ExtractStringAttr(rawResp, "radiusAttributeForGroupPolicies")
		if diags.HasError() {
			diags.AddError("Radius Attribute for Group Policies Attribute", "")
			return diags
		}

	}

	// RadiusFailOverPolicy
	if state.RadiusFailOverPolicy.IsNull() || state.RadiusFailOverPolicy.IsUnknown() {

		state.RadiusFailOverPolicy, diags = utils.ExtractStringAttr(rawResp, "radiusFailOverPolicy")
		if diags.HasError() {
			diags.AddError("Radius Failover Attribute", "")
			return diags
		}

	}

	// RadiusLoadBalancingPolicy
	if state.RadiusLoadBalancingPolicy.IsNull() || state.RadiusLoadBalancingPolicy.IsUnknown() {

		state.RadiusLoadBalancingPolicy, diags = utils.ExtractStringAttr(rawResp, "radiusLoadBalancingPolicy")
		if diags.HasError() {
			diags.AddError("Radius load balancing policyAttribute", "")
			return diags
		}

	}

	// IPAssignmentMode
	if state.IPAssignmentMode.IsNull() || state.IPAssignmentMode.IsUnknown() {

		state.IPAssignmentMode, diags = utils.ExtractStringAttr(rawResp, "ipAssignmentMode")
		if diags.HasError() {
			diags.AddError("IP Assignment mode Attribute", "")
			return diags
		}

	}

	// AdminSplashUrl
	if state.AdminSplashUrl.IsNull() || state.AdminSplashUrl.IsUnknown() {

		state.AdminSplashUrl, diags = networksWirelessSsidAdminSplashUrl(data)
		if diags.HasError() {
			diags.AddError("AdminSplashUrl Attribute", "")
			return diags
		}

	}

	// SplashTimeout
	if state.SplashTimeout.IsNull() || state.SplashTimeout.IsUnknown() {

		state.SplashTimeout, diags = utils.ExtractStringAttr(rawResp, "splashTimeout")
		if diags.HasError() {
			diags.AddError("SplashTimeout Attribute", "")
			return diags
		}

	}

	// WalledGardenEnabled
	if state.WalledGardenEnabled.IsNull() || state.WalledGardenEnabled.IsUnknown() {

		state.WalledGardenEnabled, diags = utils.ExtractBoolAttr(rawResp, "walledGardenEnabled")
		if diags.HasError() {
			diags.AddError("WalledGardanEnabled Attribute", "")
			return diags
		}

	}

	// WalledGardenRanges
	if state.WalledGardenRanges.IsNull() || state.WalledGardenRanges.IsUnknown() {

		state.WalledGardenRanges, diags = utils.ExtractListStringAttr(rawResp, "walledGardenRanges")
		if diags.HasError() {
			diags.AddError("Walled Garden Ranges Attribute", "")
			return diags
		}
	}

	// MinBitRate
	if state.MinBitRate.IsNull() || state.MinBitRate.IsUnknown() {

		minBitrateInt, diags := utils.ExtractInt64Attr(rawResp, "minBitrate")
		if diags.HasError() {
			diags.AddError("Min Bite Rate Attribute", "")
			return diags
		}

		// convert int64 into float64 type
		if !minBitrateInt.IsNull() && !minBitrateInt.IsUnknown() {
			minBitRate := types.Float64Value(float64(minBitrateInt.ValueInt64()))
			state.MinBitRate = minBitRate
		} else {
			minBitRateNull := types.Float64Null()
			state.MinBitRate = minBitRateNull
		}

	}

	// BandSelection
	if state.BandSelection.IsNull() || state.BandSelection.IsUnknown() {

		state.BandSelection, diags = utils.ExtractStringAttr(rawResp, "bandSelection")
		if diags.HasError() {
			diags.AddError("Band Selection Attribute", "")
			return diags
		}

	}

	// PerClientBandwidthLimitUp
	if state.PerClientBandwidthLimitUp.IsNull() || state.PerClientBandwidthLimitUp.IsUnknown() {

		state.PerClientBandwidthLimitUp, diags = utils.ExtractInt32Attr(rawResp, "ipAssignmentMode")
		if diags.HasError() {
			diags.AddError("Per client Bandwidth limit up Attribute", "")
			return diags
		}

	}

	// PerClientBandwidthLimitDown
	if state.PerClientBandwidthLimitDown.IsNull() || state.PerClientBandwidthLimitDown.IsUnknown() {

		state.PerClientBandwidthLimitDown, diags = utils.ExtractInt64Attr(rawResp, "perClientBandwidthLimitDown")
		if diags.HasError() {
			diags.AddError("Per client Bandwidth limit down Attribute", "")
			return diags
		}

	}

	// Visible
	if state.Visible.IsNull() || state.Visible.IsUnknown() {

		state.Visible, diags = utils.ExtractBoolAttr(rawResp, "perClientBandwidthLimitDown")
		if diags.HasError() {
			diags.AddError("Visible Attribute", "")
			return diags
		}

	}

	// AvailableOnAllAps
	if state.AvailableOnAllAps.IsNull() || state.AvailableOnAllAps.IsUnknown() {

		state.AvailableOnAllAps, diags = utils.ExtractBoolAttr(rawResp, "availableOnAllAps")
		if diags.HasError() {
			diags.AddError("AvailableOnAllAPs Attribute", "")
			return diags
		}

	}

	// AvailabilityTags
	if state.AvailabilityTags.IsNull() || state.AvailabilityTags.IsUnknown() {

		state.AvailabilityTags, diags = utils.ExtractListStringAttr(rawResp, "availabilityTags")
		if diags.HasError() {
			diags.AddError("AvailabilityTags Attribute", "")
			return diags
		}

	}

	// PerSsidBandwidthLimitUp
	if state.PerSsidBandwidthLimitUp.IsNull() || state.PerSsidBandwidthLimitUp.IsUnknown() {

		state.PerSsidBandwidthLimitUp, diags = utils.ExtractInt32Attr(rawResp, "perSsidBandwidthLimitUp")
		if diags.HasError() {
			diags.AddError("perSsidBandwidthLimitUp Attribute", "")
			return diags
		}

	}

	// PerSsidBandwidthLimitDown
	if state.PerSsidBandwidthLimitDown.IsNull() || state.PerSsidBandwidthLimitDown.IsUnknown() {

		state.PerSsidBandwidthLimitDown, diags = utils.ExtractInt32Attr(rawResp, "perSsidBandwidthLimitDown")
		if diags.HasError() {
			diags.AddError("perSsidBandwidthLimitDown Attribute", "")
			return diags
		}

	}

	// MandatoryDhcpEnabled
	if state.MandatoryDhcpEnabled.IsNull() || state.MandatoryDhcpEnabled.IsUnknown() {

		state.MandatoryDhcpEnabled, diags = utils.ExtractBoolAttr(rawResp, "mandatoryDhcpEnabled")
		if diags.HasError() {
			diags.AddError("mandatoryDhcpEnabled Attribute", "")
			return diags
		}

	}

	// Active Directory
	if state.ActiveDirectory.IsNull() || state.ActiveDirectory.IsUnknown() {
		state.ActiveDirectory, diags = NetworksWirelessSsidStateActiveDirectory(rawResp)
		if diags.HasError() {
			diags.AddError("Active Directory Attribute", "")
			return diags
		}
	}

	// PSK
	if state.PSK.IsNull() || state.PSK.IsUnknown() {
		state.PSK, diags = utils.ExtractStringAttr(rawResp, "psk")
		if diags.HasError() {
			diags.AddError("PSK Attribute", "")
			return diags
		}
	}

	// EnterpriseAdminAccess
	if state.EnterpriseAdminAccess.IsNull() || state.EnterpriseAdminAccess.IsUnknown() {
		state.EnterpriseAdminAccess, diags = utils.ExtractStringAttr(rawResp, "enterpriseAdminAccess")
		if diags.HasError() {
			diags.AddError("enterpriseAdminAccess Attribute", "")
			return diags
		}
	}

	// Dot11w
	if state.Dot11w.IsNull() || state.Dot11w.IsUnknown() {

		state.Dot11w, diags = NetworksWirelessSsidStateDot11w(rawResp)
		if diags.HasError() {
			diags.AddError("Dot11w Attribute", "")
			return diags
		}
	}

	// Dot11r
	if state.Dot11r.IsNull() || state.Dot11r.IsUnknown() {

		state.Dot11r, diags = NetworksWirelessSsidStateDot11r(rawResp)
		if diags.HasError() {
			diags.AddError("Dot11r Attribute", "")
			return diags
		}
	}

	// SplashGuestSponsorDomains
	if state.SplashGuestSponsorDomains.IsNull() || state.SplashGuestSponsorDomains.IsUnknown() {
		state.SplashGuestSponsorDomains, diags = utils.ExtractListStringAttr(rawResp, "splashGuestSponsorDomains")
		if diags.HasError() {
			diags.AddError("splashGuestSponsorDomains Attribute", "")
			return diags
		}
	}

	// OAuth
	if state.OAuth.IsNull() || state.OAuth.IsUnknown() {

		state.OAuth, diags = NetworksWirelessSsidStateOauth(rawResp)
		if diags.HasError() {
			diags.AddError("oAuth Attribute", "")
			return diags
		}
	}

	// LocalRadius
	if state.LocalRadius.IsNull() || state.LocalRadius.IsUnknown() {

		state.LocalRadius, diags = NetworksWirelessSsidStateLocalRadius(rawResp)
		if diags.HasError() {
			diags.AddError("LocalRadius Attribute", "")
			return diags
		}
	}

	// LDAP
	if state.LDAP.IsNull() || state.LDAP.IsUnknown() {
		state.LDAP, diags = NetworksWirelessSsidStateLdap(rawResp)
		if diags.HasError() {
			diags.AddError("LDAP Attribute", "")
			return diags
		}
	}

	// RadiusProxyEnabled
	if state.RadiusProxyEnabled.IsNull() || state.RadiusProxyEnabled.IsUnknown() {
		state.RadiusProxyEnabled, diags = utils.ExtractBoolAttr(rawResp, "radiusProxyEnabled")
		if diags.HasError() {
			diags.AddError("radiusProxyEnabled Attribute", "")
			return diags
		}
	}

	// RadiusTestingEnabled
	if state.RadiusTestingEnabled.IsNull() || state.RadiusTestingEnabled.IsUnknown() {
		state.RadiusTestingEnabled, diags = utils.ExtractBoolAttr(rawResp, "radiusTestingEnabled")
		if diags.HasError() {
			diags.AddError("radiusTestingEnabled Attribute", "")
			return diags
		}
	}

	// RadiusCalledStationID
	if state.RadiusCalledStationID.IsNull() || state.RadiusCalledStationID.IsUnknown() {
		state.RadiusCalledStationID, diags = utils.ExtractStringAttr(rawResp, "radiusCalledStationId")
		if diags.HasError() {
			diags.AddError("radiusCalledStationId Attribute", "")
			return diags
		}
	}

	// RadiusAuthenticationNASID
	if state.RadiusAuthenticationNASID.IsNull() || state.RadiusAuthenticationNASID.IsUnknown() {
		state.RadiusAuthenticationNASID, diags = utils.ExtractStringAttr(rawResp, "radiusAuthenticationNasId")
		if diags.HasError() {
			diags.AddError("radiusAuthenticationNasId Attribute", "")
			return diags
		}
	}

	// RadiusServerTimeout
	if state.RadiusServerTimeout.IsNull() || state.RadiusServerTimeout.IsUnknown() {
		state.RadiusServerTimeout, diags = utils.ExtractInt64Attr(rawResp, "radiusServerTimeout")
		if diags.HasError() {
			diags.AddError("radiusServerTimeout Attribute", "")
			return diags
		}
	}

	// RadiusServerAttemptsLimit
	if state.RadiusServerAttemptsLimit.IsNull() || state.RadiusServerAttemptsLimit.IsUnknown() {
		state.RadiusServerAttemptsLimit, diags = utils.ExtractInt64Attr(rawResp, "radiusServerAttemptsLimit")
		if diags.HasError() {
			diags.AddError("radiusServerAttemptsLimit Attribute", "")
			return diags
		}
	}

	// RadiusFallbackEnabled
	if state.RadiusFallbackEnabled.IsNull() || state.RadiusFallbackEnabled.IsUnknown() {
		state.RadiusFallbackEnabled, diags = utils.ExtractBoolAttr(rawResp, "radiusFallbackEnabled")
		if diags.HasError() {
			diags.AddError("radiusFallbackEnabled Attribute", "")
			return diags
		}
	}

	// RadiusCoaEnabled
	if state.RadiusCoaEnabled.IsNull() || state.RadiusCoaEnabled.IsUnknown() {
		state.RadiusCoaEnabled, diags = utils.ExtractBoolAttr(rawResp, "radiusCoaEnabled")
		if diags.HasError() {
			diags.AddError("radiusCoaEnabled Attribute", "")
			return diags
		}
	}

	// RadiusAccountingInterimInterval
	if state.RadiusAccountingInterimInterval.IsNull() || state.RadiusAccountingInterimInterval.IsUnknown() {
		state.RadiusAccountingInterimInterval, diags = utils.ExtractInt64Attr(rawResp, "radiusAccountingInterimInterval")
		if diags.HasError() {
			diags.AddError("radiusAccountingInterimInterval Attribute", "")
			return diags
		}
	}

	// UseVlanTagging
	if state.UseVlanTagging.IsNull() || state.UseVlanTagging.IsUnknown() {
		state.UseVlanTagging, diags = utils.ExtractBoolAttr(rawResp, "useVlanTagging")
		if diags.HasError() {
			diags.AddError("useVlanTagging Attribute", "")
			return diags
		}
	}

	// ConcentratorNetworkID
	if state.ConcentratorNetworkID.IsNull() || state.ConcentratorNetworkID.IsUnknown() {
		state.ConcentratorNetworkID, diags = utils.ExtractStringAttr(rawResp, "concentratorNetworkId")
		if diags.HasError() {
			diags.AddError("concentratorNetworkId Attribute", "")
			return diags
		}
	}

	// SecondaryConcentratorNetworkID
	if state.SecondaryConcentratorNetworkID.IsNull() || state.SecondaryConcentratorNetworkID.IsUnknown() {
		state.SecondaryConcentratorNetworkID, diags = utils.ExtractStringAttr(rawResp, "secondaryConcentratorNetworkId")
		if diags.HasError() {
			diags.AddError("secondaryConcentratorNetworkId Attribute", "")
			return diags
		}
	}

	// DisassociateClientsOnVpnFailOver
	if state.DisassociateClientsOnVpnFailOver.IsNull() || state.DisassociateClientsOnVpnFailOver.IsUnknown() {
		state.DisassociateClientsOnVpnFailOver, diags = utils.ExtractBoolAttr(rawResp, "disassociateClientsOnVpnFailOver")
		if diags.HasError() {
			diags.AddError("disassociateClientsOnVpnFailOver Attribute", "")
			return diags
		}
	}

	// VlanID
	if state.VlanID.IsNull() || state.VlanID.IsUnknown() {
		state.VlanID, diags = utils.ExtractInt64Attr(rawResp, "vlanId")
		if diags.HasError() {
			diags.AddError("vlanId Attribute", "")
			return diags
		}
	}

	// DefaultVlanID
	if state.DefaultVlanID.IsNull() || state.DefaultVlanID.IsUnknown() {
		state.DefaultVlanID, diags = utils.ExtractInt64Attr(rawResp, "defaultVlanId")
		if diags.HasError() {
			diags.AddError("defaultVlanId Attribute", "")
			return diags
		}
	}

	// ApTagsAndVlanIDs
	if state.ApTagsAndVlanIDs.IsNull() || state.ApTagsAndVlanIDs.IsUnknown() {
		state.ApTagsAndVlanIDs, diags = NetworksWirelessSsidStateApTagsAndVlanIds(rawResp)
		if diags.HasError() {
			diags.AddError("ApTagsAndVlanIDs Attribute", "")
			return diags
		}
	}

	// GRE
	if state.GRE.IsNull() || state.GRE.IsUnknown() {

		state.GRE, diags = NetworksWirelessSsidStateGre(rawResp)
		if diags.HasError() {
			diags.AddError("Gre Attribute", "")
			return diags
		}
	}

	// RadiusOverride
	if state.RadiusOverride.IsNull() || state.RadiusOverride.IsUnknown() {
		state.RadiusOverride, diags = utils.ExtractBoolAttr(rawResp, "radiusOverride")
		if diags.HasError() {
			diags.AddError("radiusOverride Attribute", "")
			return diags
		}
	}

	// RadiusGuestVlanEnabled
	if state.RadiusGuestVlanEnabled.IsNull() || state.RadiusGuestVlanEnabled.IsUnknown() {
		state.RadiusGuestVlanEnabled, diags = utils.ExtractBoolAttr(rawResp, "radiusGuestVlanEnabled")
		if diags.HasError() {
			diags.AddError("radiusGuestVlanEnabled Attribute", "")
			return diags
		}
	}

	// RadiusGuestVlanId
	if state.RadiusGuestVlanId.IsNull() || state.RadiusGuestVlanId.IsUnknown() {
		state.RadiusGuestVlanId, diags = utils.ExtractInt64Attr(rawResp, "radiusGuestVlanId")
		if diags.HasError() {
			diags.AddError("radiusGuestVlanId Attribute", "")
			return diags
		}
	}

	// LanIsolationEnabled
	if state.LanIsolationEnabled.IsNull() || state.LanIsolationEnabled.IsUnknown() {
		state.LanIsolationEnabled, diags = utils.ExtractBoolAttr(rawResp, "lanIsolationEnabled")
		if diags.HasError() {
			diags.AddError("lanIsolationEnabled Attribute", "")
			return diags
		}
	}

	// AdultContentFilteringEnabled
	if state.AdultContentFilteringEnabled.IsNull() || state.AdultContentFilteringEnabled.IsUnknown() {
		state.AdultContentFilteringEnabled, diags = utils.ExtractBoolAttr(rawResp, "adultContentFilteringEnabled")
		if diags.HasError() {
			diags.AddError("AdultContentFilteringEnabled Attribute", "")
			return diags
		}
	}

	// DnsRewrite
	if state.DnsRewrite.IsNull() || state.DnsRewrite.IsUnknown() {

		state.DnsRewrite, diags = NetworksWirelessSsidStateDnsRewrite(rawResp)
		if diags.HasError() {
			diags.AddError("DnsRewrite Attribute", "")
			return diags
		}

	}

	// SpeedBurst
	if state.SpeedBurst.IsNull() || state.SpeedBurst.IsUnknown() {
		state.SpeedBurst, diags = NetworksWirelessSsidStateSpeedBurst(rawResp)
		if diags.HasError() {
			diags.AddError("SpeedBurst Attribute", "")
			return diags
		}

	}

	// NamedVlans
	if state.NamedVlans.IsNull() || state.NamedVlans.IsUnknown() {
		state.NamedVlans, diags = NetworksWirelessSsidStateNamedVlans(rawResp)
		if diags.HasError() {
			diags.AddError("NamedVlans Attribute", "")
			return diags
		}
	}

	return diags
}

func updateNetworksWirelessSsidsResourcePayload(plan *NetworksWirelessSsidResourceModel) (openApiClient.UpdateNetworkWirelessSsidRequest, diag.Diagnostics) {
	var diags diag.Diagnostics
	var payload openApiClient.UpdateNetworkWirelessSsidRequest

	// Set simple attributes with checks
	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		payload.SetName(plan.Name.ValueString())
	}

	if !plan.Enabled.IsNull() && !plan.Enabled.IsUnknown() {
		payload.SetEnabled(plan.Enabled.ValueBool())
	}

	if !plan.AuthMode.IsNull() && !plan.AuthMode.IsUnknown() {
		payload.SetAuthMode(plan.AuthMode.ValueString())
	}

	if !plan.EncryptionMode.IsNull() && !plan.EncryptionMode.IsUnknown() {
		payload.SetEncryptionMode(plan.EncryptionMode.ValueString())
	}

	if !plan.PSK.IsNull() && !plan.PSK.IsUnknown() {
		payload.SetPsk(plan.PSK.ValueString())
	}

	if !plan.WPAEncryptionMode.IsNull() && !plan.WPAEncryptionMode.IsUnknown() {
		payload.SetWpaEncryptionMode(plan.WPAEncryptionMode.ValueString())
	}

	if !plan.SplashPage.IsNull() && !plan.SplashPage.IsUnknown() {
		payload.SetSplashPage(plan.SplashPage.ValueString())
	}

	if !plan.RadiusProxyEnabled.IsNull() && !plan.RadiusProxyEnabled.IsUnknown() {
		payload.SetRadiusProxyEnabled(plan.RadiusProxyEnabled.ValueBool())
	}

	if !plan.RadiusTestingEnabled.IsNull() && !plan.RadiusTestingEnabled.IsUnknown() {
		payload.SetRadiusTestingEnabled(plan.RadiusTestingEnabled.ValueBool())
	}

	if !plan.RadiusCalledStationID.IsNull() && !plan.RadiusCalledStationID.IsUnknown() {
		payload.SetRadiusCalledStationId(plan.RadiusCalledStationID.ValueString())
		payload.SetRadiusAuthenticationNasId(plan.RadiusCalledStationID.ValueString())
	}

	if !plan.RadiusFallbackEnabled.IsNull() && !plan.RadiusFallbackEnabled.IsUnknown() {
		payload.SetRadiusFallbackEnabled(plan.RadiusFallbackEnabled.ValueBool())
	}

	if !plan.RadiusCoaEnabled.IsNull() && !plan.RadiusCoaEnabled.IsUnknown() {
		payload.SetRadiusCoaEnabled(plan.RadiusCoaEnabled.ValueBool())
	}

	if !plan.RadiusFailOverPolicy.IsNull() && !plan.RadiusFailOverPolicy.IsUnknown() {
		payload.SetRadiusFailoverPolicy(plan.RadiusFailOverPolicy.ValueString())
	}

	if !plan.RadiusLoadBalancingPolicy.IsNull() && !plan.RadiusLoadBalancingPolicy.IsUnknown() {
		payload.SetRadiusLoadBalancingPolicy(plan.RadiusLoadBalancingPolicy.ValueString())
	}

	if !plan.RadiusAccountingEnabled.IsNull() && !plan.RadiusAccountingEnabled.IsUnknown() {
		payload.SetRadiusAccountingEnabled(plan.RadiusAccountingEnabled.ValueBool())
	}

	if !plan.IPAssignmentMode.IsNull() && !plan.IPAssignmentMode.IsUnknown() {
		payload.SetIpAssignmentMode(plan.IPAssignmentMode.ValueString())
	}

	if !plan.UseVlanTagging.IsNull() && !plan.UseVlanTagging.IsUnknown() {
		payload.SetUseVlanTagging(plan.UseVlanTagging.ValueBool())
	}

	if !plan.ConcentratorNetworkID.IsNull() && !plan.ConcentratorNetworkID.IsUnknown() {
		payload.SetConcentratorNetworkId(plan.ConcentratorNetworkID.ValueString())
	}

	if !plan.SecondaryConcentratorNetworkID.IsNull() && !plan.SecondaryConcentratorNetworkID.IsUnknown() {
		payload.SetSecondaryConcentratorNetworkId(plan.SecondaryConcentratorNetworkID.ValueString())
	}

	if !plan.DisassociateClientsOnVpnFailOver.IsNull() && !plan.DisassociateClientsOnVpnFailOver.IsUnknown() {
		payload.SetDisassociateClientsOnVpnFailover(plan.DisassociateClientsOnVpnFailOver.ValueBool())
	}

	if !plan.WalledGardenEnabled.IsNull() && !plan.WalledGardenEnabled.IsUnknown() {
		payload.SetWalledGardenEnabled(plan.WalledGardenEnabled.ValueBool())
	}

	if !plan.RadiusOverride.IsNull() && !plan.RadiusOverride.IsUnknown() {
		payload.SetRadiusOverride(plan.RadiusOverride.ValueBool())
	}

	if !plan.RadiusGuestVlanEnabled.IsNull() && !plan.RadiusGuestVlanEnabled.IsUnknown() {
		payload.SetRadiusGuestVlanEnabled(plan.RadiusGuestVlanEnabled.ValueBool())
	}

	if !plan.RadiusGuestVlanId.IsNull() && !plan.RadiusGuestVlanId.IsUnknown() {
		payload.SetRadiusGuestVlanId(int32(plan.RadiusGuestVlanId.ValueInt64()))
	}

	if !plan.BandSelection.IsNull() && !plan.BandSelection.IsUnknown() {
		payload.SetBandSelection(plan.BandSelection.ValueString())
	}

	if !plan.LanIsolationEnabled.IsNull() && !plan.LanIsolationEnabled.IsUnknown() {
		payload.SetLanIsolationEnabled(plan.LanIsolationEnabled.ValueBool())
	}

	if !plan.Visible.IsNull() && !plan.Visible.IsUnknown() {
		payload.SetVisible(plan.Visible.ValueBool())
	}

	if !plan.AvailableOnAllAps.IsNull() && !plan.AvailableOnAllAps.IsUnknown() {
		payload.SetAvailableOnAllAps(plan.AvailableOnAllAps.ValueBool())
	}

	if !plan.MandatoryDhcpEnabled.IsNull() && !plan.MandatoryDhcpEnabled.IsUnknown() {
		payload.SetMandatoryDhcpEnabled(plan.MandatoryDhcpEnabled.ValueBool())
	}

	if !plan.AdultContentFilteringEnabled.IsNull() && !plan.AdultContentFilteringEnabled.IsUnknown() {
		payload.SetAdultContentFilteringEnabled(plan.AdultContentFilteringEnabled.ValueBool())
	}

	if !plan.RadiusServerTimeout.IsNull() && !plan.RadiusServerTimeout.IsUnknown() {
		radiusServerTimeout, err := utils.Int32Pointer(plan.RadiusServerTimeout.ValueInt64())
		if err.HasError() {
			tflog.Error(context.Background(), fmt.Sprintf("%s", err.Errors()))
			diags.Append(err...)
		}
		payload.SetRadiusServerTimeout(*radiusServerTimeout)
	}

	if !plan.RadiusServerAttemptsLimit.IsNull() && !plan.RadiusServerAttemptsLimit.IsUnknown() {
		radiusServerAttemptsLimit, err := utils.Int32Pointer(plan.RadiusServerAttemptsLimit.ValueInt64())
		if err.HasError() {
			tflog.Error(context.Background(), fmt.Sprintf("%s", err.Errors()))
			diags.Append(err...)
		}
		payload.RadiusServerAttemptsLimit = radiusServerAttemptsLimit
	}

	if !plan.RadiusAccountingInterimInterval.IsNull() && !plan.RadiusAccountingInterimInterval.IsUnknown() {
		radiusAccountingInterimInterval, err := utils.Int32Pointer(plan.RadiusAccountingInterimInterval.ValueInt64())
		if err.HasError() {
			tflog.Error(context.Background(), fmt.Sprintf("%s", err.Errors()))
			diags.Append(err...)
		}
		payload.RadiusAccountingInterimInterval = radiusAccountingInterimInterval
	}

	if !plan.VlanID.IsNull() && !plan.VlanID.IsUnknown() {
		vlanId, err := utils.Int32Pointer(plan.VlanID.ValueInt64())
		if err.HasError() {
			tflog.Error(context.Background(), fmt.Sprintf("%s", err.Errors()))
			diags.Append(err...)
		}
		payload.VlanId = vlanId
	}

	if !plan.DefaultVlanID.IsNull() && !plan.DefaultVlanID.IsUnknown() {
		defaultVlanId, err := utils.Int32Pointer(plan.DefaultVlanID.ValueInt64())
		if err.HasError() {
			tflog.Error(context.Background(), fmt.Sprintf("%s", err.Errors()))
			diags.Append(err...)
		}
		payload.DefaultVlanId = defaultVlanId
	}

	if !plan.MinBitRate.IsNull() && !plan.MinBitRate.IsUnknown() {
		minBitRate, err := utils.Float32Pointer(plan.MinBitRate.ValueFloat64())
		if err.HasError() {
			tflog.Error(context.Background(), fmt.Sprintf("%s", err.Errors()))
			diags.Append(err...)
		}
		payload.MinBitrate = minBitRate
	}

	if !plan.PerClientBandwidthLimitUp.IsNull() && !plan.PerClientBandwidthLimitUp.IsUnknown() {
		perClientBandwidthLimitUp, err := utils.Int32Pointer(plan.PerClientBandwidthLimitUp.ValueInt64())
		if err.HasError() {
			tflog.Error(context.Background(), fmt.Sprintf("%s", err.Errors()))
			diags.Append(err...)
		}
		payload.PerClientBandwidthLimitUp = perClientBandwidthLimitUp
	}

	if !plan.PerClientBandwidthLimitDown.IsNull() && !plan.PerClientBandwidthLimitDown.IsUnknown() {
		perClientBandwidthLimitDown, err := utils.Int32Pointer(plan.PerClientBandwidthLimitDown.ValueInt64())
		if err.HasError() {
			tflog.Error(context.Background(), fmt.Sprintf("%s", err.Errors()))
			diags.Append(err...)
		}
		payload.PerClientBandwidthLimitDown = perClientBandwidthLimitDown
	}

	if !plan.PerSsidBandwidthLimitUp.IsNull() && !plan.PerSsidBandwidthLimitUp.IsUnknown() {
		perSsidBandwidthLimitUp, err := utils.Int32Pointer(plan.PerSsidBandwidthLimitUp.ValueInt64())
		if err.HasError() {
			tflog.Error(context.Background(), fmt.Sprintf("%s", err.Errors()))
			diags.Append(err...)
		}
		payload.PerSsidBandwidthLimitUp = perSsidBandwidthLimitUp
	}

	if !plan.PerSsidBandwidthLimitDown.IsNull() && !plan.PerSsidBandwidthLimitDown.IsUnknown() {
		perSsidBandwidthLimitDown, err := utils.Int32Pointer(plan.PerSsidBandwidthLimitDown.ValueInt64())
		if err.HasError() {
			tflog.Error(context.Background(), fmt.Sprintf("%s", err.Errors()))
			diags.Append(err...)
		}
		payload.PerSsidBandwidthLimitDown = perSsidBandwidthLimitDown
	}

	if !plan.WalledGardenRanges.IsNull() && !plan.WalledGardenRanges.IsUnknown() {
		walledGardenRanges, err := utils.ExtractStringsFromList(plan.WalledGardenRanges)
		if err.HasError() {
			diags.Append(err...)
		}
		payload.WalledGardenRanges = walledGardenRanges
	}

	availabilityTags, err := utils.ExtractStringsFromList(plan.AvailabilityTags)
	if err.HasError() {
		diags.Append(err...)
	}
	payload.AvailabilityTags = availabilityTags

	splashGuestSponsorDomains, err := utils.ExtractStringsFromList(plan.SplashGuestSponsorDomains)
	if err.HasError() {
		diags.Append(err...)
	}
	payload.SplashGuestSponsorDomains = splashGuestSponsorDomains

	// Check enterprise admin access
	if !plan.EnterpriseAdminAccess.IsNull() && !plan.EnterpriseAdminAccess.IsUnknown() {
		payload.SetEnterpriseAdminAccess(plan.EnterpriseAdminAccess.ValueString())
	}

	// Expand complex attributes
	dot11w, err := NetworksWirelessSsidPayloadDot11w(plan.Dot11w)
	if err.HasError() {
		diags.Append(err...)
	}
	payload.Dot11w = dot11w

	dot11r, err := NetworksWirelessSsidPayloadDot11r(plan.Dot11r)
	if err.HasError() {
		diags.Append(err...)
	}
	payload.Dot11r = dot11r

	oauth, err := NetworksWirelessSsidPayloadOauth(plan.OAuth)
	if err.HasError() {
		diags.Append(err...)
	}
	payload.Oauth = oauth

	localRadius, err := NetworksWirelessSsidPayloadLocalRadius(plan.LocalRadius)
	if err.HasError() {
		diags.Append(err...)
	}
	payload.LocalRadius = &localRadius

	ldap, err := NetworksWirelessSsidPayloadLdap(plan.LDAP)
	if err.HasError() {
		diags.Append(err...)
	}
	payload.Ldap = ldap

	activeDirectory, err := NetworksWirelessSsidPayloadActiveDirectory(plan.ActiveDirectory)
	if err.HasError() {
		diags.Append(err...)
	}
	payload.ActiveDirectory = activeDirectory

	radiusServers, err := NetworksWirelessSsidPayloadRadiusServers(plan.RadiusServers)
	if err.HasError() {
		diags.Append(err...)
	}
	payload.RadiusServers = radiusServers

	radiusAccountingServers, err := NetworksWirelessSsidPayloadRadiusAccountingServers(plan.RadiusAccountingServers)
	if err.HasError() {
		diags.Append(err...)
	}
	payload.RadiusAccountingServers = radiusAccountingServers

	apTagsAndVlanIds, err := NetworksWirelessSsidPayloadApTagsAndVlanIds(plan.ApTagsAndVlanIDs)
	if err.HasError() {
		diags.Append(err...)
	}
	payload.ApTagsAndVlanIds = apTagsAndVlanIds

	gre, err := NetworksWirelessSsidPayloadGre(plan.GRE)
	if err.HasError() {
		diags.Append(err...)
	}
	payload.Gre = gre

	dnsRewrite, err := NetworksWirelessSsidPayloadDnsRewrite(plan.DnsRewrite)
	payload.DnsRewrite = dnsRewrite
	if err.HasError() {
		diags.Append(err...)
	}

	speedBurst, err := NetworksWirelessSsidPayloadSpeedBurst(plan.SpeedBurst)
	payload.SpeedBurst = speedBurst
	if err.HasError() {
		diags.Append(err...)
	}

	namedVlans, err := NetworksWirelessSsidPayloadNamedVlans(plan.NamedVlans)
	payload.NamedVlans = namedVlans
	if err.HasError() {
		diags.Append(err...)
	}

	return payload, diags

}

func (r *NetworksWirelessSsidsResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = r.typeName

}

func (r *NetworksWirelessSsidsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier for the resource, generated by the Meraki API.",
				Computed:    true,
			},
			"active_directory": schema.SingleNestedAttribute{
				MarkdownDescription: `The current setting for Active Directory. Only valid if splashPage is 'Password-protected with Active Directory'`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"credentials": schema.SingleNestedAttribute{
						MarkdownDescription: `(Optional) The credentials of the user account to be used by the AP to bind to your Active Directory server. The Active Directory account should have permissions on all your Active Directory servers. Only valid if the splashPage is 'Password-protected with Active Directory'.`,
						Computed:            true,
						Optional:            true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"login_name": schema.StringAttribute{
								MarkdownDescription: `The login name of the Active Directory account.`,
								Computed:            true,
								Optional:            true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"password": schema.StringAttribute{
								MarkdownDescription: `The password to the Active Directory user account.`,
								Sensitive:           true,
								Computed:            true,
								Optional:            true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
						},
					},
					"servers": schema.ListNestedAttribute{
						MarkdownDescription: `The Active Directory servers to be used for authentication.`,
						Computed:            true,
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{

								"host": schema.StringAttribute{
									MarkdownDescription: `IP address (or FQDN) of your Active Directory server.`,
									Computed:            true,
									Optional:            true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
									},
								},
								"port": schema.Int64Attribute{
									MarkdownDescription: `(Optional) UDP port the Active Directory server listens on. By default, uses port 3268.`,
									Computed:            true,
									Optional:            true,
									PlanModifiers: []planmodifier.Int64{
										int64planmodifier.UseStateForUnknown(),
									},
								},
							},
						},
					},
				},
			},
			"admin_splash_url": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"adult_content_filtering_enabled": schema.BoolAttribute{
				MarkdownDescription: `Boolean indicating whether or not adult content will be blocked`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"ap_tags_and_vlan_ids": schema.ListNestedAttribute{
				MarkdownDescription: `The list of tags and VLAN IDs used for VLAN tagging. This param is only valid when the ipAssignmentMode is 'Bridge mode' or 'Layer 3 roaming'`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{

						"tags": schema.ListAttribute{
							MarkdownDescription: `Array of AP tags`,
							Computed:            true,
							Optional:            true,

							ElementType: types.StringType,
						},
						"vlan_id": schema.Int64Attribute{
							MarkdownDescription: `Numerical identifier that is assigned to the VLAN`,
							Computed:            true,
							Optional:            true,
							PlanModifiers: []planmodifier.Int64{
								int64planmodifier.UseStateForUnknown(),
							},
						},
					},
				},
			},
			"auth_mode": schema.StringAttribute{
				MarkdownDescription: `The association control method for the SSID`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(
						"8021x-google",
						"8021x-localradius",
						"8021x-meraki",
						"8021x-nac",
						"8021x-radius",
						"ipsk-with-nac",
						"ipsk-with-radius",
						"ipsk-without-radius",
						"open",
						"open-enhanced",
						"open-with-nac",
						"open-with-radius",
						"psk",
					),
				},
			},
			"availability_tags": schema.ListAttribute{
				MarkdownDescription: `List of tags for this SSID. If availableOnAllAps is false, then the SSID is only broadcast by APs with tags matching any of the tags in this list`,
				Optional:            true,
				Computed:            true,

				ElementType: types.StringType,
			},
			"available_on_all_aps": schema.BoolAttribute{
				MarkdownDescription: `Whether all APs broadcast the SSID or if it's restricted to APs matching any availability tags`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"band_selection": schema.StringAttribute{
				MarkdownDescription: `The client-serving radio frequencies of this SSID in the default indoor RF profile`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(
						"5 GHz band only",
						"Dual band operation",
						"Dual band operation with Band Steering",
					),
				},
			},
			"concentrator_network_id": schema.StringAttribute{
				MarkdownDescription: `The concentrator to use when the ipAssignmentMode is 'Layer 3 roaming with a concentrator' or 'VPN'.`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"default_vlan_id": schema.Int64Attribute{
				MarkdownDescription: `The default VLAN ID used for 'all other APs'. This param is only valid when the ipAssignmentMode is 'Bridge mode' or 'Layer 3 roaming'`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"disassociate_clients_on_vpn_fail_over": schema.BoolAttribute{
				MarkdownDescription: `Disassociate clients when 'VPN' concentrator failover occurs in order to trigger clients to re-associate and generate new DHCP requests. This param is only valid if ipAssignmentMode is 'VPN'.`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"dns_rewrite": schema.SingleNestedAttribute{
				MarkdownDescription: `DNS servers rewrite settings`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{

					"dns_custom_name_servers": schema.ListAttribute{
						MarkdownDescription: `User specified DNS servers (up to two servers)`,
						Computed:            true,
						Optional:            true,
						ElementType:         types.StringType,
					},
					"enabled": schema.BoolAttribute{
						MarkdownDescription: `Boolean indicating whether or not DNS server rewrite is enabled. If disabled, upstream DNS will be used`,
						Computed:            true,
						Optional:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"dot11r": schema.SingleNestedAttribute{
				MarkdownDescription: `The current setting for 802.11r`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{

					"adaptive": schema.BoolAttribute{
						MarkdownDescription: `(Optional) Whether 802.11r is adaptive or not.`,
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"enabled": schema.BoolAttribute{
						MarkdownDescription: `Whether 802.11r is enabled or not.`,
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"dot11w": schema.SingleNestedAttribute{
				MarkdownDescription: `The current setting for Protected Management Frames (802.11w).`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{

					"enabled": schema.BoolAttribute{
						MarkdownDescription: `Whether 802.11w is enabled or not.`,
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"required": schema.BoolAttribute{
						MarkdownDescription: `(Optional) Whether 802.11w is required or not.`,
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: `Whether or not the SSID is enabled`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"encryption_mode": schema.StringAttribute{
				MarkdownDescription: `The psk encryption mode for the SSID`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enterprise_admin_access": schema.StringAttribute{
				MarkdownDescription: `Whether or not an SSID is accessible by 'enterprise' administrators ('access disabled' or 'access enabled')`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(
						"access disabled",
						"access enabled",
					),
				},
			},
			"gre": schema.SingleNestedAttribute{
				MarkdownDescription: `Ethernet over GRE settings`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{

					"concentrator": schema.SingleNestedAttribute{
						MarkdownDescription: `The EoGRE concentrator's settings`,
						Computed:            true,
						Optional:            true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{

							"host": schema.StringAttribute{
								MarkdownDescription: `The EoGRE concentrator's IP or FQDN. This param is required when ipAssignmentMode is 'Ethernet over GRE'.`,
								Computed:            true,
								Optional:            true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
						},
					},
					"key": schema.Int64Attribute{
						MarkdownDescription: `Optional numerical identifier that will add the GRE key field to the GRE header. Used to identify an individual traffic flow within a tunnel.`,
						Computed:            true,
						Optional:            true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"ip_assignment_mode": schema.StringAttribute{
				MarkdownDescription: `The client IP assignment mode`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(
						"Bridge mode",
						"Ethernet over GRE",
						"Layer 3 roaming",
						"Layer 3 roaming with a concentrator",
						"NAT mode",
						"VPN",
					),
				},
			},
			"lan_isolation_enabled": schema.BoolAttribute{
				MarkdownDescription: `Boolean indicating whether Layer 2 LAN isolation should be enabled or disabled. Only configurable when ipAssignmentMode is 'Bridge mode'.`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"ldap": schema.SingleNestedAttribute{
				MarkdownDescription: `The current setting for LDAP. Only valid if splashPage is 'Password-protected with LDAP'.`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{

					"base_distinguished_name": schema.StringAttribute{
						MarkdownDescription: `The base distinguished name of users on the LDAP server.`,
						Computed:            true,
						Optional:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"credentials": schema.SingleNestedAttribute{
						MarkdownDescription: `(Optional) The credentials of the user account to be used by the AP to bind to your LDAP server. The LDAP account should have permissions on all your LDAP servers.`,
						Computed:            true,
						Optional:            true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{

							"distinguished_name": schema.StringAttribute{
								MarkdownDescription: `The distinguished name of the LDAP user account (example: cn=user,dc=meraki,dc=com).`,
								Computed:            true,
								Optional:            true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"password": schema.StringAttribute{
								MarkdownDescription: `The password of the LDAP user account.`,
								Sensitive:           true,
								Computed:            true,
								Optional:            true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
						},
					},
					"server_ca_certificate": schema.SingleNestedAttribute{
						MarkdownDescription: `The CA certificate used to sign the LDAP server's key.`,
						Computed:            true,
						Optional:            true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{

							"contents": schema.StringAttribute{
								MarkdownDescription: `The contents of the CA certificate. Must be in PEM or DER format.`,
								Computed:            true,
								Optional:            true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
						},
					},
					"servers": schema.ListNestedAttribute{
						MarkdownDescription: `The LDAP servers to be used for authentication.`,
						Computed:            true,
						Optional:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{

								"host": schema.StringAttribute{
									MarkdownDescription: `IP address (or FQDN) of your LDAP server.`,
									Computed:            true,
									Optional:            true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
									},
								},
								"port": schema.Int64Attribute{
									MarkdownDescription: `UDP port the LDAP server listens on.`,
									Computed:            true,
									Optional:            true,
									PlanModifiers: []planmodifier.Int64{
										int64planmodifier.UseStateForUnknown(),
									},
								},
							},
						},
					},
				},
			},
			"local_auth": schema.BoolAttribute{
				MarkdownDescription: `Extended local auth flag for Enterprise NAC`,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"local_radius": schema.SingleNestedAttribute{
				MarkdownDescription: `The current setting for Local Authentication, a built-in RADIUS server on the access point. Only valid if authMode is '8021x-localradius'.`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{

					"cache_timeout": schema.Int64Attribute{
						MarkdownDescription: `The duration (in seconds) for which LDAP and OCSP lookups are cached.`,
						Computed:            true,
						Optional:            true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"certificate_authentication": schema.SingleNestedAttribute{
						MarkdownDescription: `The current setting for certificate verification.`,
						Computed:            true,
						Optional:            true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{

							"client_root_ca_certificate": schema.SingleNestedAttribute{
								MarkdownDescription: `The Client CA Certificate used to sign the client certificate.`,
								Computed:            true,
								Optional:            true,
								PlanModifiers: []planmodifier.Object{
									objectplanmodifier.UseStateForUnknown(),
								},
								Attributes: map[string]schema.Attribute{

									"contents": schema.StringAttribute{
										MarkdownDescription: `The contents of the Client CA Certificate. Must be in PEM or DER format.`,
										Computed:            true,
										Optional:            true,
										PlanModifiers: []planmodifier.String{
											stringplanmodifier.UseStateForUnknown(),
										},
									},
								},
							},
							"enabled": schema.BoolAttribute{
								MarkdownDescription: `Whether or not to use EAP-TLS certificate-based authentication to validate wireless clients.`,
								Computed:            true,
								Optional:            true,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"ocsp_responder_url": schema.StringAttribute{
								MarkdownDescription: `(Optional) The URL of the OCSP responder to verify client certificate status.`,
								Computed:            true,
								Optional:            true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"use_ldap": schema.BoolAttribute{
								MarkdownDescription: `Whether or not to verify the certificate with LDAP.`,
								Computed:            true,
								Optional:            true,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"use_ocsp": schema.BoolAttribute{
								MarkdownDescription: `Whether or not to verify the certificate with OCSP.`,
								Computed:            true,
								Optional:            true,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
						},
					},
					"password_authentication": schema.SingleNestedAttribute{
						MarkdownDescription: `The current setting for password-based authentication.`,
						Computed:            true,
						Optional:            true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{

							"enabled": schema.BoolAttribute{
								MarkdownDescription: `Whether or not to use EAP-TTLS/PAP or PEAP-GTC password-based authentication via LDAP lookup.`,
								Computed:            true,
								Optional:            true,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
						},
					},
				},
			},
			"mandatory_dhcp_enabled": schema.BoolAttribute{
				MarkdownDescription: `Whether clients connecting to this SSID must use the IP address assigned by the DHCP server`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"min_bitrate": schema.Float64Attribute{
				MarkdownDescription: `The minimum bitrate in Mbps of this SSID in the default indoor RF profile`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.Float64{
					float64planmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: `The name of the SSID`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"named_vlans": schema.SingleNestedAttribute{
				MarkdownDescription: `Named VLAN settings.`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"radius": schema.SingleNestedAttribute{
						MarkdownDescription: `RADIUS settings. This param is only valid when authMode is 'open-with-radius' and ipAssignmentMode is not 'NAT mode'.`,
						Computed:            true,
						Optional:            true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"guest_vlan": schema.SingleNestedAttribute{
								MarkdownDescription: `guest vlan settings`,
								Computed:            true,
								Optional:            true,
								PlanModifiers: []planmodifier.Object{
									objectplanmodifier.UseStateForUnknown(),
								},
								Attributes: map[string]schema.Attribute{

									"enabled": schema.BoolAttribute{
										MarkdownDescription: `Whether or not RADIUS guest named VLAN is enabled.`,
										Computed:            true,
										Optional:            true,
										PlanModifiers: []planmodifier.Bool{
											boolplanmodifier.UseStateForUnknown(),
										},
									},
									"name": schema.StringAttribute{
										MarkdownDescription: `RADIUS guest VLAN name.`,
										Computed:            true,
										Optional:            true,
										PlanModifiers: []planmodifier.String{
											stringplanmodifier.UseStateForUnknown(),
										},
									},
								},
							},
						},
					},
					"tagging": schema.SingleNestedAttribute{
						MarkdownDescription: `VLAN tagging settings. This param is only valid when ipAssignmentMode is 'Bridge mode' or 'Layer 3 roaming'.`,
						Computed:            true,
						Optional:            true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"by_ap_tags": schema.ListNestedAttribute{
								MarkdownDescription: `The list of AP tags and VLAN names used for named VLAN tagging. If an AP has a tag matching one in the list, then traffic on this SSID will be directed to use the VLAN name associated to the tag.`,
								Computed:            true,
								Optional:            true,
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"tags": schema.ListAttribute{
											MarkdownDescription: `List of AP tags.`,
											Computed:            true,
											Optional:            true,
											ElementType:         types.StringType,
										},
										"vlan_name": schema.StringAttribute{
											MarkdownDescription: `VLAN name that will be used to tag traffic.`,
											Computed:            true,
											Optional:            true,
											PlanModifiers: []planmodifier.String{
												stringplanmodifier.UseStateForUnknown(),
											},
										},
									},
								},
							},
							"default_vlan_name": schema.StringAttribute{
								MarkdownDescription: `The default VLAN name used to tag traffic in the absence of a matching AP tag.`,
								Computed:            true,
								Optional:            true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"enabled": schema.BoolAttribute{
								MarkdownDescription: `Whether or not traffic should be directed to use specific VLAN names.`,
								Computed:            true,
								Optional:            true,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
						},
					},
				},
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: `networkId path parameter. Network ID`,
				Required:            true,
			},
			"number": schema.Int64Attribute{
				MarkdownDescription: `Unique identifier of the SSID`,
				Required:            true,
				//            Differents_types: `   parameter: schema.TypeString, item: schema.TypeInt`,
			},
			"oauth": schema.SingleNestedAttribute{
				MarkdownDescription: `The OAuth settings of this SSID. Only valid if splashPage is 'Google OAuth'.`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{

					"allowed_domains": schema.ListAttribute{
						MarkdownDescription: `(Optional) The list of domains allowed access to the network.`,
						Computed:            true,
						Optional:            true,
						ElementType:         types.StringType,
					},
				},
			},
			"per_client_bandwidth_limit_down": schema.Int64Attribute{
				MarkdownDescription: `The download bandwidth limit in Kbps. (0 represents no limit.)`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"per_client_bandwidth_limit_up": schema.Int64Attribute{
				MarkdownDescription: `The upload bandwidth limit in Kbps. (0 represents no limit.)`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"per_ssid_bandwidth_limit_down": schema.Int64Attribute{
				MarkdownDescription: `The total download bandwidth limit in Kbps (0 represents no limit)`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"per_ssid_bandwidth_limit_up": schema.Int64Attribute{
				MarkdownDescription: `The total upload bandwidth limit in Kbps (0 represents no limit)`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"psk": schema.StringAttribute{
				MarkdownDescription: `The passkey for the SSID. This param is only valid if the authMode is 'psk'`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"radius_accounting_enabled": schema.BoolAttribute{
				MarkdownDescription: `Whether or not RADIUS accounting is enabled`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"radius_accounting_interim_interval": schema.Int64Attribute{
				MarkdownDescription: `The interval (in seconds) in which accounting information is updated and sent to the RADIUS accounting server.`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"radius_accounting_servers": schema.ListNestedAttribute{
				MarkdownDescription: `List of RADIUS accounting 802.1X servers to be used for authentication`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{

						"ca_certificate": schema.StringAttribute{
							MarkdownDescription: `Certificate used for authorization for the RADSEC Server`,
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"host": schema.StringAttribute{
							MarkdownDescription: `IP address (or FQDN) to which the APs will send RADIUS accounting messages`,
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"port": schema.Int64Attribute{
							MarkdownDescription: `Port on the RADIUS server that is listening for accounting messages`,
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.Int64{
								int64planmodifier.UseStateForUnknown(),
							},
						},
						"open_roaming_certificate_id": schema.Int64Attribute{
							MarkdownDescription: `The ID of the Open roaming Certificate attached to radius server`,
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.Int64{
								int64planmodifier.UseStateForUnknown(),
							},
						},
						"rad_sec_enabled": schema.BoolAttribute{
							MarkdownDescription: `Use RADSEC (TLS over TCP) to connect to this RADIUS accounting server. Requires radiusProxyEnabled.`,
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
						"secret": schema.StringAttribute{
							MarkdownDescription: `Shared key used to authenticate messages between the APs and RADIUS server`,
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
					},
				},
			},
			"radius_attribute_for_group_policies": schema.StringAttribute{
				MarkdownDescription: `RADIUS attribute used to look up group policies`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(
						"Airespace-ACL-Name",
						"Aruba-User-Role",
						"Filter-Id",
						"Reply-Message",
					),
				},
			},
			"radius_authentication_nas_id": schema.StringAttribute{
				MarkdownDescription: `The template of the NAS identifier to be used for RADIUS authentication (ex. $NODE_MAC$:$VAP_NUM$).`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"radius_called_station_id": schema.StringAttribute{
				MarkdownDescription: `The template of the called station identifier to be used for RADIUS (ex. $NODE_MAC$:$VAP_NUM$).`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"radius_coa_enabled": schema.BoolAttribute{
				MarkdownDescription: `If true, Meraki devices will act as a RADIUS Dynamic Authorization Server and will respond to RADIUS Change-of-Authorization and Disconnect messages sent by the RADIUS server.`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"radius_enabled": schema.BoolAttribute{
				MarkdownDescription: `Whether RADIUS authentication is enabled`,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"radius_fail_over_policy": schema.StringAttribute{
				MarkdownDescription: `Policy which determines how authentication requests should be handled in the event that all of the configured RADIUS servers are unreachable`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(
						"Allow access",
						"Deny access",
					),
				},
			},
			"radius_fallback_enabled": schema.BoolAttribute{
				MarkdownDescription: `Whether or not higher priority RADIUS servers should be retried after 60 seconds.`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"radius_guest_vlan_enabled": schema.BoolAttribute{
				MarkdownDescription: `Whether or not RADIUS Guest VLAN is enabled. This param is only valid if the authMode is 'open-with-radius' and addressing mode is not set to 'isolated' or 'nat' mode`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"radius_guest_vlan_id": schema.Int64Attribute{
				MarkdownDescription: `VLAN ID of the RADIUS Guest VLAN. This param is only valid if the authMode is 'open-with-radius' and addressing mode is not set to 'isolated' or 'nat' mode`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"radius_load_balancing_policy": schema.StringAttribute{
				MarkdownDescription: `Policy which determines which RADIUS server will be contacted first in an authentication attempt, and the ordering of any necessary retry attempts`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(
						"Round robin",
						"Strict priority order",
					),
				},
			},
			"radius_override": schema.BoolAttribute{
				MarkdownDescription: `If true, the RADIUS response can override VLAN tag. This is not valid when ipAssignmentMode is 'NAT mode'.`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"radius_proxy_enabled": schema.BoolAttribute{
				MarkdownDescription: `If true, Meraki devices will proxy RADIUS messages through the Meraki cloud to the configured RADIUS auth and accounting servers.`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"radius_server_attempts_limit": schema.Int64Attribute{
				MarkdownDescription: `The maximum number of transmit attempts after which a RADIUS server is failed over (must be between 1-5).`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"radius_server_timeout": schema.Int64Attribute{
				MarkdownDescription: `The amount of time for which a RADIUS client waits for a reply from the RADIUS server (must be between 1-10 seconds).`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"radius_servers": schema.ListNestedAttribute{
				MarkdownDescription: `The RADIUS 802.1X servers to be used for authentication. This param is only valid if the authMode is 'open-with-radius', '8021x-radius' or 'ipsk-with-radius'`,
				Optional:            true,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{

						"ca_certificate": schema.StringAttribute{
							MarkdownDescription: `Certificate used for authorization for the RADSEC Server`,
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"host": schema.StringAttribute{
							MarkdownDescription: `IP address of your RADIUS server`,
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"open_roaming_certificate_id": schema.Int64Attribute{
							MarkdownDescription: `The ID of the Openroaming Certificate attached to radius server.`,
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.Int64{
								int64planmodifier.UseStateForUnknown(),
							},
						},
						"port": schema.Int64Attribute{
							MarkdownDescription: `UDP port the RADIUS server listens on for Access-requests`,
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.Int64{
								int64planmodifier.UseStateForUnknown(),
							},
						},
						"rad_sec_enabled": schema.BoolAttribute{
							MarkdownDescription: `Use RADSEC (TLS over TCP) to connect to this RADIUS server. Requires radiusProxyEnabled.`,
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
						"secret": schema.StringAttribute{
							MarkdownDescription: `RADIUS client shared secret`,
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
					},
				},
			},
			"radius_testing_enabled": schema.BoolAttribute{
				MarkdownDescription: `If true, Meraki devices will periodically send Access-Request messages to configured RADIUS servers using identity 'meraki_8021x_test' to ensure that the RADIUS servers are reachable.`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"secondary_concentrator_network_id": schema.StringAttribute{
				MarkdownDescription: `The secondary concentrator to use when the ipAssignmentMode is 'VPN'. If configured, the APs will switch to using this concentrator if the primary concentrator is unreachable. This param is optional. ('disabled' represents no secondary concentrator.)`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"speed_burst": schema.SingleNestedAttribute{
				MarkdownDescription: `The SpeedBurst setting for this SSID'`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: `Boolean indicating whether or not to allow users to temporarily exceed the bandwidth limit for short periods while still keeping them under the bandwidth limit over time.`,
						Computed:            true,
						Optional:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"splash_guest_sponsor_domains": schema.ListAttribute{
				MarkdownDescription: `Array of valid sponsor email domains for sponsored guest splash type.`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},

				ElementType: types.StringType,
			},
			"splash_page": schema.StringAttribute{
				MarkdownDescription: `The type of splash page for the SSID`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(
						"Billing",
						"Cisco ISE",
						"Click-through splash page",
						"Facebook Wi-Fi",
						"Google Apps domain",
						"Google OAuth",
						"None",
						"Password-protected with Active Directory",
						"Password-protected with LDAP",
						"Password-protected with Meraki RADIUS",
						"Password-protected with custom RADIUS",
						"SMS authentication",
						"Sponsored guest",
						"Systems Manager Sentry",
					),
				},
			},
			"splash_timeout": schema.StringAttribute{
				MarkdownDescription: `Splash page timeout`,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Computed: true,
			},
			"ssid_admin_accessible": schema.BoolAttribute{
				MarkdownDescription: `SSID Administrator access status`,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
				Computed: true,
			},
			"use_vlan_tagging": schema.BoolAttribute{
				MarkdownDescription: `Whether or not traffic should be directed to use specific VLANs. This param is only valid if the ipAssignmentMode is 'Bridge mode' or 'Layer 3 roaming'`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"visible": schema.BoolAttribute{
				MarkdownDescription: `Whether the SSID is advertised or hidden by the AP`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"vlan_id": schema.Int64Attribute{
				MarkdownDescription: `The VLAN ID used for VLAN tagging. This param is only valid when the ipAssignmentMode is 'Layer 3 roaming with a concentrator' or 'VPN'`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"walled_garden_enabled": schema.BoolAttribute{
				MarkdownDescription: `Allow users to access a configurable list of IP ranges prior to sign-on`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"walled_garden_ranges": schema.ListAttribute{
				MarkdownDescription: `Domain names and IP address ranges available in Walled Garden mode`,
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},

				ElementType: types.StringType,
			},
			"wpa_encryption_mode": schema.StringAttribute{
				MarkdownDescription: `The types of WPA encryption`,
				Computed:            true,
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.OneOf(
						"WPA1 and WPA2",
						"WPA1 only",
						"WPA2 only",
						"WPA3 192-bit Security",
						"WPA3 Transition Mode",
						"WPA3 only",
					),
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

// Create creates the resource and sets the initial Terraform state.
func (r *NetworksWirelessSsidsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan NetworksWirelessSsidResourceModel

	// Read the Terraform configuration into the model
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare the request payload
	payload, payloadDiags := updateNetworksWirelessSsidsResourcePayload(&plan)
	if payloadDiags.HasError() {
		tflog.Error(ctx, "Failed to create resource payload", map[string]interface{}{
			"error": payloadDiags,
		})
		resp.Diagnostics.AddError(
			"Error creating ssid payload",
			fmt.Sprintf("Unexpected error: %s", payloadDiags),
		)
		return
	}

	maxRetries := r.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(r.client.GetConfig().Retry4xxErrorWaitTime)

	inlineResp, httpResp, err := tools.CustomHttpRequestRetry[*openApiClient.GetNetworkWirelessSsids200ResponseInner](ctx, maxRetries, retryDelay, func() (*openApiClient.GetNetworkWirelessSsids200ResponseInner, *http.Response, error) {
		inline, respHttp, err := r.client.WirelessApi.UpdateNetworkWirelessSsid(context.Background(), plan.NetworkId.ValueString(), fmt.Sprint(plan.Number.ValueInt64())).UpdateNetworkWirelessSsidRequest(payload).Execute()
		return inline, respHttp, err
	})
	if err != nil {
		// Check for the specific unmarshaling error
		if strings.Contains(err.Error(), "json: cannot unmarshal number") && strings.Contains(err.Error(), "GetNetworkWirelessSsids200ResponseInner.minBitrate") {
			tflog.Warn(ctx, "Suppressing unmarshaling error: json: cannot unmarshal number into GetNetworkWirelessSsids200ResponseInner.minBitrate of type int32")
		} else {
			tflog.Error(ctx, "HTTP Call Failed", map[string]interface{}{
				"error": err.Error(),
			})
			resp.Diagnostics.AddError(
				"HTTP Call Failed",
				fmt.Sprintf("Details: %s", err.Error()),
			)
		}
	}

	// Check for API success response code
	if httpResp != nil && httpResp.StatusCode != 200 {
		responseBody, _ := tools.ReadAndCloseBody(httpResp)
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			tools.NewHttpDiagnostics(httpResp, responseBody),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	diags = updateNetworksWirelessSsidsResourceState(ctx, &plan, inlineResp, httpResp)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *NetworksWirelessSsidsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state NetworksWirelessSsidResourceModel

	// Read Terraform prior state into the model
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	maxRetries := r.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(r.client.GetConfig().Retry4xxErrorWaitTime)

	inlineResp, httpResp, err := tools.CustomHttpRequestRetry[*openApiClient.GetNetworkWirelessSsids200ResponseInner](ctx, maxRetries, retryDelay, func() (*openApiClient.GetNetworkWirelessSsids200ResponseInner, *http.Response, error) {
		inline, respHttp, err := r.client.WirelessApi.GetNetworkWirelessSsid(context.Background(), state.NetworkId.ValueString(), fmt.Sprint(state.Number.ValueInt64())).Execute()
		if err != nil {
			// Check for specific error
			if strings.Contains(err.Error(), "json: cannot unmarshal number") && strings.Contains(err.Error(), "GetNetworkWirelessSsids200ResponseInner.minBitrate") {
				tflog.Warn(ctx, "Suppressing unmarshaling error: json: cannot unmarshal number into GetNetworkWirelessSsids200ResponseInner.minBitrate of type int32")
				err = nil
			}
		}
		return inline, respHttp, err
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Call Failed",
			fmt.Sprintf("Details: %s", err.Error()),
		)
	}

	// Check for API success response code
	if httpResp != nil && httpResp.StatusCode != 200 {
		responseBody, _ := tools.ReadAndCloseBody(httpResp)
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			tools.NewHttpDiagnostics(httpResp, responseBody),
		)
	}

	diags = updateNetworksWirelessSsidsResourceState(ctx, &state, inlineResp, httpResp)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *NetworksWirelessSsidsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan NetworksWirelessSsidResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read Terraform plan into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Prepare the request payload
	payload, payloadDiags := updateNetworksWirelessSsidsResourcePayload(&plan)
	if payloadDiags.HasError() {
		tflog.Error(ctx, "Failed to create resource payload", map[string]interface{}{
			"error": payloadDiags,
		})
		resp.Diagnostics.AddError(
			"Error creating ssid payload",
			fmt.Sprintf("Unexpected error: %s", payloadDiags),
		)
		return
	}

	maxRetries := r.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(r.client.GetConfig().Retry4xxErrorWaitTime)

	inlineResp, httpResp, err := tools.CustomHttpRequestRetry[*openApiClient.GetNetworkWirelessSsids200ResponseInner](ctx, maxRetries, retryDelay, func() (*openApiClient.GetNetworkWirelessSsids200ResponseInner, *http.Response, error) {
		inline, respHttp, err := r.client.WirelessApi.UpdateNetworkWirelessSsid(context.Background(), plan.NetworkId.ValueString(), fmt.Sprint(plan.Number.ValueInt64())).UpdateNetworkWirelessSsidRequest(payload).Execute()
		if err != nil {
			// Check for specific error
			if strings.Contains(err.Error(), "json: cannot unmarshal number") && strings.Contains(err.Error(), "GetNetworkWirelessSsids200ResponseInner.minBitrate") {
				tflog.Warn(ctx, "Suppressing unmarshaling error: json: cannot unmarshal number into GetNetworkWirelessSsids200ResponseInner.minBitrate of type int32")
				err = nil
			}
		}
		return inline, respHttp, err
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Call Failed",
			fmt.Sprintf("Details: %s", err.Error()),
		)
	}

	// Check for API success response code
	if httpResp != nil && httpResp.StatusCode != 200 {
		responseBody, _ := tools.ReadAndCloseBody(httpResp)
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			tools.NewHttpDiagnostics(httpResp, responseBody),
		)
	}

	diags = updateNetworksWirelessSsidsResourceState(ctx, &plan, inlineResp, httpResp)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *NetworksWirelessSsidsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *NetworksWirelessSsidResourceModel
	// Read Terraform plan state into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	payload := *openApiClient.NewUpdateNetworkWirelessSsidRequest()
	payload.SetEnabled(false)
	payload.SetName("")
	payload.SetAuthMode("open")
	payload.SetVlanId(1)

	maxRetries := r.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(r.client.GetConfig().Retry4xxErrorWaitTime)

	_, httpResp, err := tools.CustomHttpRequestRetry[*openApiClient.GetNetworkWirelessSsids200ResponseInner](ctx, maxRetries, retryDelay, func() (*openApiClient.GetNetworkWirelessSsids200ResponseInner, *http.Response, error) {
		inline, respHttp, err := r.client.WirelessApi.UpdateNetworkWirelessSsid(context.Background(), state.NetworkId.ValueString(), fmt.Sprint(state.Number.ValueInt64())).UpdateNetworkWirelessSsidRequest(payload).Execute()
		if err != nil {
			// Check for specific error
			if strings.Contains(err.Error(), "json: cannot unmarshal number") && strings.Contains(err.Error(), "GetNetworkWirelessSsids200ResponseInner.minBitrate") {
				tflog.Warn(ctx, "Suppressing unmarshaling error: json: cannot unmarshal number into GetNetworkWirelessSsids200ResponseInner.minBitrate of type int32")
				err = nil
			}
		}
		return inline, respHttp, err
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Call Failed",
			fmt.Sprintf("Details: %s", err.Error()),
		)
	}

	// Check for API success response code
	if httpResp != nil && httpResp.StatusCode != 200 {
		responseBody, _ := tools.ReadAndCloseBody(httpResp)
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			tools.NewHttpDiagnostics(httpResp, responseBody),
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

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
	i, err := strconv.ParseInt(idParts[1], 10, 64)
	if err != nil {
		panic(err)
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("number"), i)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
