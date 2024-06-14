package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/core-infra-svcs/terraform-provider-meraki/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
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
	"io"
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
	return &NetworksWirelessSsidsResource{}
}

// NetworksWirelessSsidsResource defines the resource implementation.
type NetworksWirelessSsidsResource struct {
	client *openApiClient.APIClient
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
	RadiusGuestVlanID                types.Int64   `tfsdk:"radius_guest_vlan_id" json:"radiusGuestVlanId"`
	MinBitRate                       types.Float64 `tfsdk:"min_bit_rate" json:"minBitRate"`
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
	BaseDistinguishedName types.String `tfsdk:"base_distinguished_name" json:"baseDistinguishedName"`
	ServerCaCertificate   types.Object `tfsdk:"server_ca_certificate" json:"serverCaCertificate"`
}

// LdapServer represents the structure for an LDAP server
type LdapServer struct {
	Host types.String `tfsdk:"host" json:"host"`
	Port types.Int64  `tfsdk:"port" json:"port"`
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
	DistinguishedName types.String `tfsdk:"distinguished_name" json:"distinguishedName"`
	Password          types.String `tfsdk:"password" json:"password"`
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
	VlanID types.Int64 `tfsdk:"vlan_id" json:"vlanId"`
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

// GuestVlan represents the structure for the RADIUS guest VLAN settings
type GuestVlan struct {
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
			LogonName: credentialsObject.DistinguishedName.ValueStringPointer(),
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
		vlanId, err := utils.Int32Pointer(tagAndVlan.VlanID.ValueInt64())
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

	// GuestVlan
	var guestVlan GuestVlan
	err = namedVlansObject.Radius.As(context.Background(), &guestVlan, basetypes.ObjectAsOptions{})
	if err.HasError() {
		diags = append(diags, err.Errors()...)
	}

	byApTags, err := NetworksWirelessSsidPayloadByApTags(tagging.ByApTags)
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
			GuestVlan: &openApiClient.UpdateNetworkWirelessSsidRequestNamedVlansRadiusGuestVlan{
				Enabled: guestVlan.Enabled.ValueBoolPointer(),
				Name:    guestVlan.Name.ValueStringPointer(),
			},
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

func NetworksWirelessSsidPayloadSplashGuestSponsorDomains(input types.List) ([]string, diag.Diagnostics) {
	if input.IsNull() || input.IsUnknown() || input.IsUnknown() {
		return nil, nil
	}

	var diags diag.Diagnostics
	var splashGuestSponsorDomains []string

	err := input.ElementsAs(context.Background(), splashGuestSponsorDomains, true)
	if err.HasError() {
		diags.Append(err...)
	}

	return splashGuestSponsorDomains, diags
}

func NetworksWirelessSsidPayloadByApTags(input types.List) ([]openApiClient.UpdateNetworkWirelessSsidRequestNamedVlansTaggingByApTagsInner, diag.Diagnostics) {
	if input.IsNull() || input.IsUnknown() {
		return nil, nil
	}

	var diags diag.Diagnostics
	var byApTags []openApiClient.UpdateNetworkWirelessSsidRequestNamedVlansTaggingByApTagsInner

	var byApTagsList []ApTagsAndVlanID
	err := input.ElementsAs(context.Background(), &byApTagsList, true)
	if err.HasError() {
		diags = append(diags, err.Errors()...)
		return nil, diags
	}

	for _, byApTag := range byApTagsList {
		var tags []string
		for _, tag := range byApTag.Tags.Elements() {
			tags = append(tags, tag.String())
		}

		vlan := strconv.FormatInt(byApTag.VlanID.ValueInt64(), 10)
		byApTags = append(byApTags, openApiClient.UpdateNetworkWirelessSsidRequestNamedVlansTaggingByApTagsInner{
			Tags:     tags,
			VlanName: &vlan,
		})
	}
	return byApTags, diags
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

func NetworksWirelessSsidStateDot11w(input *openApiClient.UpdateNetworkApplianceSsidRequestDot11w) *Dot11w {
	if input == nil {
		return nil
	}
	return &Dot11w{
		Enabled:  types.BoolValue(*input.Enabled),
		Required: types.BoolValue(*input.Required),
	}
}

func NetworksWirelessSsidStateDot11r(input *openApiClient.UpdateNetworkWirelessSsidRequestDot11r) *Dot11r {
	if input == nil {
		return nil
	}
	return &Dot11r{
		Enabled:  types.BoolValue(*input.Enabled),
		Adaptive: types.BoolValue(*input.Adaptive),
	}
}

func NetworksWirelessSsidStateOauth(input *openApiClient.UpdateNetworkWirelessSsidRequestOauth) *OAuth {
	if input == nil {
		return nil
	}

	var allowedDomains []types.String
	for _, domain := range input.GetAllowedDomains() {
		allowedDomains = append(allowedDomains, types.StringValue(domain))
	}

	allowedDomainsObj, err := types.ListValueFrom(context.Background(), types.StringType, allowedDomains)
	if err.HasError() {
		tflog.Error(context.Background(), fmt.Sprintf("%s", err.Errors()))
	}

	return &OAuth{
		AllowedDomains: allowedDomainsObj,
	}
}

func NetworksWirelessSsidStateLocalRadius(input *openApiClient.UpdateNetworkWirelessSsidRequestLocalRadius) *LocalRadius {
	if input == nil {
		return nil
	}

	passwordAuthentication, err := types.ObjectValueFrom(context.Background(), map[string]attr.Type{"enabled": types.BoolType}, input.PasswordAuthentication)
	if err.HasError() {
		tflog.Error(context.Background(), fmt.Sprintf("%s", err.Errors()))
	}

	// cacheTimeout
	cacheTimeout := int64(*input.CacheTimeout)

	// certificateAuthentication
	certificateAuthentication, err := types.ObjectValueFrom(context.Background(), map[string]attr.Type{"enabled": types.BoolType,
		"useLdap": types.BoolType, "useOcsp": types.BoolType, "OcspResponderUrl": types.StringType, "clientRootCaCertificate": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"contents": types.StringType,
			},
		}}, input.CertificateAuthentication)
	if err.HasError() {
		tflog.Error(context.Background(), fmt.Sprintf("%s", err.Errors()))
	}

	return &LocalRadius{
		CacheTimeout:              types.Int64Value(cacheTimeout),
		PasswordAuthentication:    passwordAuthentication,
		CertificateAuthentication: certificateAuthentication,
		/*CertificateAuthentication: &CertificateAuthentication{
			Enabled:          types.BoolValue(certificateAuthentication.Enabled),
			UseLdap:          types.BoolValue(*input.CertificateAuthentication.UseLdap),
			UseOcsp:          types.BoolValue(*input.CertificateAuthentication.UseOcsp),
			OcspResponderUrl: types.StringValue(*input.CertificateAuthentication.OcspResponderUrl),
			ClientRootCaCertificate: &CaCertificate{
				Contents: types.StringValue(*input.CertificateAuthentication.ClientRootCaCertificate.Contents),
			},
		},

		*/
	}
}

func NetworksWirelessSsidStateLdap(httpResp map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	ldapAttrs := map[string]attr.Type{
		"servers": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
			"host": types.StringType,
			"port": types.Int64Type,
		}}},
		"credentials": types.ObjectType{AttrTypes: map[string]attr.Type{
			"distinguished_name": types.StringType,
			"password":           types.StringType,
		}},
		"base_distinguished_name": types.StringType,
		"server_ca_certificate": types.ObjectType{AttrTypes: map[string]attr.Type{
			"contents": types.StringType,
		}},
	}

	ldapObject, err := utils.ExtractObjectAttr(httpResp, "ldap", ldapAttrs)
	if err.HasError() {
		diags.Append(err...)
	}

	return ldapObject, diags
}

func NetworksWirelessSsidStateActiveDirectory(httpResp map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	activeDirectoryAttrs := map[string]attr.Type{
		"servers": types.ListType{
			ElemType: types.ObjectType{
				AttrTypes: map[string]attr.Type{
					"host": types.StringType,
					"port": types.Int64Type,
				},
			},
		},
		"credentials": types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"logon_name": types.StringType,
				"password":   types.StringType,
			},
		},
	}

	activeDirectoryList, err := utils.ExtractObjectAttr(httpResp, "activeDirectory", activeDirectoryAttrs)
	if err.HasError() {
		diags.Append(err...)
	}

	return activeDirectoryList, diags
}

func NetworksWirelessSsidStateApTagsAndVlanIds(httpResp map[string]interface{}) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics

	//var apTagsAndVlanIds []attr.Value

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

func NetworksWirelessSsidStateGre(input *openApiClient.UpdateNetworkWirelessSsidRequestGre) *GRE {
	if input == nil {
		return nil
	}

	greConcentratorObject, err := types.ObjectValueFrom(context.Background(), map[string]attr.Type{"host": types.StringType}, input.Concentrator)
	if err.HasError() {
		tflog.Error(context.Background(), fmt.Sprintf("%s", err.Errors()))
	}

	return &GRE{
		Concentrator: greConcentratorObject,
		Key:          types.Int64Value(int64(input.GetKey())),
	}
}

func NetworksWirelessSsidStateDnsRewrite(input *openApiClient.UpdateNetworkWirelessSsidRequestDnsRewrite) *DnsRewrite {
	if input == nil {
		return nil
	}
	return &DnsRewrite{
		Enabled: types.BoolValue(input.GetEnabled()),
		//DnsCustomNameservers: types.StringValue(input.DnsCustomNameservers),
	}
}

func NetworksWirelessSsidStateSpeedBurst(input *openApiClient.UpdateNetworkWirelessSsidRequestSpeedBurst) *SpeedBurst {
	if input == nil {
		return nil
	}
	return &SpeedBurst{
		Enabled: types.BoolValue(input.GetEnabled()),
	}
}

func NetworksWirelessSsidStateNamedVlans(input *openApiClient.UpdateNetworkWirelessSsidRequestNamedVlans) *NamedVlans {
	if input == nil {
		return nil
	}

	// tagging
	taggingObject, err := types.ObjectValueFrom(context.Background(), map[string]attr.Type{"enabled": types.BoolType, "default_vlan_name": types.StringType, "by_ap_tags": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{"tags": types.ListType{}, "vlan_name": types.StringType}}}}, input.Tagging)
	if err.HasError() {
		tflog.Error(context.Background(), fmt.Sprintf("%s", err.Errors()))
	}

	// radius
	radiusObject, err := types.ObjectValueFrom(context.Background(), map[string]attr.Type{"guest_vlan": types.ObjectType{AttrTypes: map[string]attr.Type{"enabled": types.BoolType, "name": types.StringType}}}, input.Radius)
	if err.HasError() {
		tflog.Error(context.Background(), fmt.Sprintf("%s", err.Errors()))
	}

	return &NamedVlans{
		Tagging: taggingObject,
		Radius:  radiusObject,
	}
}

func NetworksWirelessSsidStateByApTags(input *[]openApiClient.UpdateNetworkWirelessSsidRequestNamedVlansTaggingByApTagsInner) []ByApTag {
	if input == nil {
		return nil
	}
	var byApTags []ByApTag
	for _, byApTag := range *input {
		byApTags = append(byApTags, ByApTag{
			//Tags:     types.StringSliceValue(byApTag.Tags),
			VlanName: types.StringValue(byApTag.GetVlanName()),
		})
	}
	return byApTags
}

// updateNetworksWirelessSsidsResourceState updates the resource state with the provided api data.
func updateNetworksWirelessSsidsResourceState(ctx context.Context, state *NetworksWirelessSsidResourceModel, data *openApiClient.GetNetworkWirelessSsids200ResponseInner, httpResp map[string]interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	// Number
	if state.Number.IsNull() || state.Number.IsUnknown() {
		number := int64(*data.Number)
		state.Number = types.Int64Value(number)
	}

	// Import Id
	if !state.NetworkId.IsNull() && !state.NetworkId.IsUnknown() && !state.Number.IsNull() && !state.Number.IsUnknown() {
		id := state.NetworkId.ValueString() + "," + strconv.FormatInt(state.Number.ValueInt64(), 10)
		state.Id = types.StringValue(id)
	} else {
		state.Id = types.StringNull()
	}

	// Name
	if state.Name.IsNull() || state.Name.IsUnknown() {

		if _, ok := data.GetNameOk(); ok {
			name := data.GetName()
			state.Name = types.StringValue(name)
		} else {
			state.Name = types.StringNull()
		}

	}

	// Enabled
	if state.Enabled.IsNull() || state.Enabled.IsUnknown() {
		enabled := data.GetEnabled()
		state.Enabled = types.BoolValue(enabled)
	}

	// SplashPage
	if state.SplashPage.IsNull() || state.SplashPage.IsUnknown() {

		if _, ok := data.GetSplashPageOk(); ok {
			splashPage := data.GetSplashPage()
			state.SplashPage = types.StringValue(splashPage)
		} else {
			state.SplashPage = types.StringNull()
		}

	}

	// SsidAdminAccessible
	if state.SsidAdminAccessible.IsNull() || state.SsidAdminAccessible.IsUnknown() {

		if _, ok := data.GetSsidAdminAccessibleOk(); ok {
			ssidAdminAccessible := data.GetSsidAdminAccessible()
			state.SsidAdminAccessible = types.BoolValue(ssidAdminAccessible)
		} else {
			state.SsidAdminAccessible = types.BoolNull()
		}
	}

	// LocalAuth
	if state.LocalAuth.IsNull() || state.LocalAuth.IsUnknown() {

		if _, ok := data.GetLocalAuthOk(); ok {
			localAuth := data.GetLocalAuth()
			state.LocalAuth = types.BoolValue(localAuth)
		} else {
			state.LocalAuth = types.BoolNull()
		}

	}

	// AuthMode
	if state.AuthMode.IsNull() || state.AuthMode.IsUnknown() {

		if _, ok := data.GetAuthModeOk(); ok {
			authMode := data.GetAuthMode()
			state.AuthMode = types.StringValue(authMode)
		} else {
			state.AuthMode = types.StringNull()
		}

	}

	// EncryptionMode
	if state.EncryptionMode.IsNull() || state.EncryptionMode.IsUnknown() {

		if _, ok := data.GetEncryptionModeOk(); ok {
			encryptionMode := data.GetEncryptionMode()
			state.EncryptionMode = types.StringValue(encryptionMode)
		} else {
			state.EncryptionMode = types.StringNull()
		}

	}

	// WPAEncryptionMode
	if state.WPAEncryptionMode.IsNull() || state.WPAEncryptionMode.IsUnknown() {

		if _, ok := data.GetWpaEncryptionModeOk(); ok {
			wpaEncryptionMode := data.GetWpaEncryptionMode()
			state.WPAEncryptionMode = types.StringValue(wpaEncryptionMode)
		} else {
			state.WPAEncryptionMode = types.StringNull()
		}
	}

	// RadiusServers
	if state.RadiusServers.IsNull() || state.RadiusServers.IsUnknown() {

		if _, ok := data.GetRadiusServersOk(); ok {
			radiusServers, err := NetworksWirelessSsidStateRadiusServers(data.RadiusServers)
			if err.HasError() {
				return err
			}
			state.RadiusServers = radiusServers
		} else {
			radiusServerAttr := map[string]attr.Type{
				"host":                        types.StringType,
				"port":                        types.Int64Type,
				"secret":                      types.StringType,
				"rad_sec_enabled":             types.BoolType,
				"ca_certificate":              types.StringType,
				"open_roaming_certificate_id": types.Int64Type,
			}
			state.RadiusServers = types.ListNull(types.ObjectType{AttrTypes: radiusServerAttr})
		}

	}

	// RadiusAccountingServers
	if state.RadiusAccountingServers.IsNull() || state.RadiusAccountingServers.IsUnknown() {

		if _, ok := data.GetRadiusAccountingServersOk(); ok {
			radiusAccountingServers, err := NetworksWirelessSsidStateRadiusAccountingServers(data.RadiusAccountingServers)
			if err.HasError() {
				return err
			}
			state.RadiusAccountingServers = radiusAccountingServers
		} else {
			radiusServerAttr := map[string]attr.Type{
				"host":                        types.StringType,
				"port":                        types.Int64Type,
				"secret":                      types.StringType,
				"rad_sec_enabled":             types.BoolType,
				"ca_certificate":              types.StringType,
				"open_roaming_certificate_id": types.Int64Type,
			}
			state.RadiusAccountingServers = types.ListNull(types.ObjectType{AttrTypes: radiusServerAttr})
		}

	}

	// RadiusAccountingEnabled
	if state.RadiusAccountingEnabled.IsNull() || state.RadiusAccountingEnabled.IsUnknown() {

		if _, ok := data.GetRadiusAccountingEnabledOk(); ok {
			radiusAccountingEnabled := data.GetRadiusAccountingEnabled()
			state.RadiusAccountingEnabled = types.BoolValue(radiusAccountingEnabled)
		} else {
			state.RadiusAccountingEnabled = types.BoolNull()
		}

	}

	// RadiusEnabled
	if state.RadiusEnabled.IsNull() || state.RadiusEnabled.IsUnknown() {

		if _, ok := data.GetRadiusEnabledOk(); ok {
			radiusEnabled := data.GetRadiusEnabled()
			state.RadiusEnabled = types.BoolValue(radiusEnabled)
		} else {
			state.RadiusEnabled = types.BoolNull()
		}
	}

	// RadiusAttributeForGroupPolicies
	if state.RadiusAttributeForGroupPolicies.IsNull() || state.RadiusAttributeForGroupPolicies.IsUnknown() {

		if _, ok := data.GetRadiusAttributeForGroupPoliciesOk(); ok {
			radiusAttributeForGroupPolicies := data.GetRadiusAttributeForGroupPolicies()
			state.RadiusAttributeForGroupPolicies = types.StringValue(radiusAttributeForGroupPolicies)
		} else {
			state.RadiusAttributeForGroupPolicies = types.StringNull()
		}

	}

	// RadiusFailOverPolicy
	if state.RadiusFailOverPolicy.IsNull() || state.RadiusFailOverPolicy.IsUnknown() {

		if _, ok := data.GetRadiusFailoverPolicyOk(); ok {
			radiusFailOverPolicy := data.GetRadiusFailoverPolicy()
			state.RadiusFailOverPolicy = types.StringValue(radiusFailOverPolicy)
		} else {
			state.RadiusFailOverPolicy = types.StringNull()
		}

	}

	// RadiusLoadBalancingPolicy
	if state.RadiusLoadBalancingPolicy.IsNull() || state.RadiusLoadBalancingPolicy.IsUnknown() {

		if _, ok := data.GetRadiusLoadBalancingPolicyOk(); ok {
			radiusLoadBalancingPolicy := data.GetRadiusLoadBalancingPolicy()
			state.RadiusLoadBalancingPolicy = types.StringValue(radiusLoadBalancingPolicy)
		} else {
			state.RadiusLoadBalancingPolicy = types.StringNull()
		}

	}

	// IPAssignmentMode
	if state.IPAssignmentMode.IsNull() || state.IPAssignmentMode.IsUnknown() {

		if _, ok := data.GetIpAssignmentModeOk(); ok {
			ipAssignmentMode := data.GetIpAssignmentMode()
			state.IPAssignmentMode = types.StringValue(ipAssignmentMode)
		} else {
			state.IPAssignmentMode = types.StringNull()
		}

	}

	// AdminSplashUrl
	if state.AdminSplashUrl.IsNull() || state.AdminSplashUrl.IsUnknown() {
		if _, ok := data.GetAdminSplashUrlOk(); ok {
			state.AdminSplashUrl, diags = networksWirelessSsidAdminSplashUrl(data)
		} else {
			state.AdminSplashUrl = types.StringNull()
		}

	}

	// SplashTimeout
	if state.SplashTimeout.IsNull() || state.SplashTimeout.IsUnknown() {

		if _, ok := data.GetSplashTimeoutOk(); ok {
			splashTimeout := data.GetSplashTimeout()
			state.SplashTimeout = types.StringValue(splashTimeout)
		} else {
			state.SplashTimeout = types.StringNull()
		}

	}

	// WalledGardenEnabled
	if state.WalledGardenEnabled.IsNull() || state.WalledGardenEnabled.IsUnknown() {
		if _, ok := data.GetWalledGardenEnabledOk(); ok {
			walledGardenEnabled := data.GetWalledGardenEnabled()
			state.WalledGardenEnabled = types.BoolValue(walledGardenEnabled)
		} else {
			state.WalledGardenEnabled = types.BoolNull()
		}

	}

	// WalledGardenRanges
	if state.WalledGardenRanges.IsNull() || state.WalledGardenRanges.IsUnknown() {

		if _, ok := data.GetWalledGardenRangesOk(); ok {
			walledGardenRanges, walledGardenRangesErr := types.ListValueFrom(ctx, types.StringType, data.WalledGardenRanges)
			if walledGardenRangesErr.HasError() {
				diags.Append(walledGardenRangesErr...)
			}
			state.WalledGardenRanges = walledGardenRanges
		} else {
			state.WalledGardenRanges = types.ListNull(types.StringType)
		}

	}

	// MinBitRate
	if state.MinBitRate.IsNull() || state.MinBitRate.IsUnknown() {

		if _, ok := data.GetMinBitrateOk(); ok {
			minBitRate := float64(data.GetMinBitrate())
			state.MinBitRate = types.Float64Value(minBitRate)
		} else {
			state.MinBitRate = types.Float64Null()
		}

	}

	// BandSelection
	if state.BandSelection.IsNull() || state.BandSelection.IsUnknown() {

		if _, ok := data.GetBandSelectionOk(); ok {
			bandSelection := data.GetBandSelection()
			state.BandSelection = types.StringValue(bandSelection)
		} else {
			state.BandSelection = types.StringNull()
		}

	}

	// PerClientBandwidthLimitUp
	if state.PerClientBandwidthLimitUp.IsNull() || state.PerClientBandwidthLimitUp.IsUnknown() {

		if _, ok := data.GetPerSsidBandwidthLimitUpOk(); ok {
			perClientBandwidthLimitUp := int64(data.GetPerClientBandwidthLimitUp())
			state.PerClientBandwidthLimitUp = types.Int64Value(perClientBandwidthLimitUp)
		} else {
			state.PerClientBandwidthLimitUp = types.Int64Null()
		}

	}

	// PerClientBandwidthLimitDown
	if state.PerClientBandwidthLimitDown.IsNull() || state.PerClientBandwidthLimitDown.IsUnknown() {

		if _, ok := data.GetPerClientBandwidthLimitDownOk(); ok {
			perClientBandwidthLimitDown := int64(data.GetPerSsidBandwidthLimitDown())
			state.PerClientBandwidthLimitDown = types.Int64Value(perClientBandwidthLimitDown)
		} else {
			state.PerClientBandwidthLimitDown = types.Int64Null()
		}

	}

	// Visible
	if state.Visible.IsNull() || state.Visible.IsUnknown() {
		if _, ok := data.GetVisibleOk(); ok {
			visible := data.GetVisible()
			state.Visible = types.BoolValue(visible)
		} else {
			state.Visible = types.BoolNull()
		}

	}

	// AvailableOnAllAps
	if state.AvailableOnAllAps.IsNull() || state.AvailableOnAllAps.IsUnknown() {

		if _, ok := data.GetAvailableOnAllApsOk(); ok {
			availableOnAllAps := data.GetAvailableOnAllAps()
			state.AvailableOnAllAps = types.BoolValue(availableOnAllAps)
		} else {
			state.AvailableOnAllAps = types.BoolNull()
		}

	}

	// AvailabilityTags
	if state.AvailabilityTags.IsNull() || state.AvailabilityTags.IsUnknown() {

		if _, ok := data.GetAvailabilityTagsOk(); ok {
			availabilityTags, availabilityTagsErr := types.ListValueFrom(ctx, types.StringType, data.AvailabilityTags)
			if availabilityTagsErr.HasError() {
				diags.Append(availabilityTagsErr...)
			}
			state.AvailabilityTags = availabilityTags
		} else {
			state.AvailabilityTags = types.ListNull(types.StringType)
		}

	}

	// PerSsidBandwidthLimitUp
	if state.PerSsidBandwidthLimitUp.IsNull() || state.PerSsidBandwidthLimitUp.IsUnknown() {

		if _, ok := data.GetPerSsidBandwidthLimitUpOk(); ok {
			perSsidBandwidthLimitUp := int64(data.GetPerSsidBandwidthLimitUp())
			state.PerSsidBandwidthLimitUp = types.Int64Value(perSsidBandwidthLimitUp)
		} else {
			state.PerSsidBandwidthLimitUp = types.Int64Null()
		}

	}

	// PerSsidBandwidthLimitDown
	if state.PerSsidBandwidthLimitDown.IsNull() || state.PerSsidBandwidthLimitDown.IsUnknown() {

		if _, ok := data.GetPerSsidBandwidthLimitDownOk(); ok {
			perSsidBandwidthLimitDown := int64(data.GetPerSsidBandwidthLimitDown())
			state.PerSsidBandwidthLimitDown = types.Int64Value(perSsidBandwidthLimitDown)
		} else {
			state.PerSsidBandwidthLimitDown = types.Int64Null()
		}

	}

	// MandatoryDhcpEnabled
	if state.MandatoryDhcpEnabled.IsNull() || state.MandatoryDhcpEnabled.IsUnknown() {

		if _, ok := data.GetMandatoryDhcpEnabledOk(); ok {
			mandatoryDhcpEnabled := data.GetMandatoryDhcpEnabled()
			state.MandatoryDhcpEnabled = types.BoolValue(mandatoryDhcpEnabled)
		} else {
			state.MandatoryDhcpEnabled = types.BoolNull()
		}

	}

	// Active Directory
	if state.ActiveDirectory.IsNull() || state.ActiveDirectory.IsUnknown() {
		state.ActiveDirectory, diags = NetworksWirelessSsidStateActiveDirectory(httpResp)
		if diags.HasError() {
			return diags
		}
	}

	// PSK
	if state.PSK.IsNull() || state.PSK.IsUnknown() {
		state.PSK, diags = utils.ExtractStringAttr(httpResp, "psk")
		if diags.HasError() {
			return diags
		}
	}

	// EnterpriseAdminAccess
	if state.EnterpriseAdminAccess.IsNull() || state.EnterpriseAdminAccess.IsUnknown() {
		state.EnterpriseAdminAccess, diags = utils.ExtractStringAttr(httpResp, "enterpriseAdminAccess")
		if diags.HasError() {
			return diags
		}
	}

	// Dot11w
	if state.Dot11w.IsNull() || state.Dot11w.IsUnknown() {
		dot11wAttrs := map[string]attr.Type{
			"enabled":  types.BoolType,
			"required": types.BoolType,
		}

		state.Dot11w, diags = utils.ExtractObjectAttr(httpResp, "dot11w", dot11wAttrs)
		if diags.HasError() {
			return diags
		}
	}

	// Dot11r
	if state.Dot11r.IsNull() || state.Dot11r.IsUnknown() {
		dot11rAttrs := map[string]attr.Type{
			"enabled":  types.BoolType,
			"adaptive": types.BoolType,
		}

		state.Dot11r, diags = utils.ExtractObjectAttr(httpResp, "dot11r", dot11rAttrs)
		if diags.HasError() {
			return diags
		}
	}

	// SplashGuestSponsorDomains
	if state.SplashGuestSponsorDomains.IsNull() || state.SplashGuestSponsorDomains.IsUnknown() {
		state.SplashGuestSponsorDomains, diags = utils.ExtractStringSliceAttr(httpResp, "splashGuestSponsorDomains")
		if diags.HasError() {
			return diags
		}
	}

	// OAuth
	if state.OAuth.IsNull() || state.OAuth.IsUnknown() {
		oauthAttrs := map[string]attr.Type{
			"allowed_domains": types.ListType{ElemType: types.StringType},
		}

		state.OAuth, diags = utils.ExtractObjectAttr(httpResp, "oauth", oauthAttrs)
		if diags.HasError() {
			return diags
		}
	}

	// LocalRadius
	if state.LocalRadius.IsNull() || state.LocalRadius.IsUnknown() {
		localRadiusAttrs := map[string]attr.Type{
			"cache_timeout": types.Int64Type,
			"password_authentication": types.ObjectType{AttrTypes: map[string]attr.Type{
				"enabled": types.BoolType,
			}},
			"certificate_authentication": types.ObjectType{AttrTypes: map[string]attr.Type{
				"enabled":            types.BoolType,
				"use_ldap":           types.BoolType,
				"use_ocsp":           types.BoolType,
				"ocsp_responder_url": types.StringType,
				"client_root_ca_certificate": types.ObjectType{AttrTypes: map[string]attr.Type{
					"contents": types.StringType,
				},
				},
			}},
		}

		state.LocalRadius, diags = utils.ExtractObjectAttr(httpResp, "localRadius", localRadiusAttrs)
		if diags.HasError() {
			return diags
		}
	}

	// LDAP
	if state.LDAP.IsNull() || state.LDAP.IsUnknown() {
		state.LDAP, diags = NetworksWirelessSsidStateLdap(httpResp)
		if diags.HasError() {
			return diags
		}
	}

	// RadiusProxyEnabled
	if state.RadiusProxyEnabled.IsNull() || state.RadiusProxyEnabled.IsUnknown() {
		state.RadiusProxyEnabled, diags = utils.ExtractBoolAttr(httpResp, "radiusProxyEnabled")
		if diags.HasError() {
			return diags
		}
	}

	// RadiusTestingEnabled
	if state.RadiusProxyEnabled.IsNull() || state.RadiusProxyEnabled.IsUnknown() {
		state.RadiusTestingEnabled, diags = utils.ExtractBoolAttr(httpResp, "radiusTestingEnabled")
		if diags.HasError() {
			return diags
		}
	}

	// RadiusCalledStationID
	if state.RadiusCalledStationID.IsNull() || state.RadiusCalledStationID.IsUnknown() {
		state.RadiusCalledStationID, diags = utils.ExtractStringAttr(httpResp, "radiusCalledStationId")
		if diags.HasError() {
			return diags
		}
	}

	// RadiusAuthenticationNASID
	if state.RadiusAuthenticationNASID.IsNull() || state.RadiusAuthenticationNASID.IsUnknown() {
		state.RadiusAuthenticationNASID, diags = utils.ExtractStringAttr(httpResp, "radiusAuthenticationNasId")
		if diags.HasError() {
			return diags
		}
	}

	// RadiusServerTimeout
	if state.RadiusServerTimeout.IsNull() || state.RadiusServerTimeout.IsUnknown() {
		state.RadiusServerTimeout, diags = utils.ExtractInt64Attr(httpResp, "radiusServerTimeout")
		if diags.HasError() {
			return diags
		}
	}

	// RadiusServerAttemptsLimit
	if state.RadiusServerAttemptsLimit.IsNull() || state.RadiusServerAttemptsLimit.IsUnknown() {
		state.RadiusServerAttemptsLimit, diags = utils.ExtractInt64Attr(httpResp, "radiusServerAttemptsLimit")
		if diags.HasError() {
			return diags
		}
	}

	// RadiusFallbackEnabled
	if state.RadiusFallbackEnabled.IsNull() || state.RadiusFallbackEnabled.IsUnknown() {
		state.RadiusFallbackEnabled, diags = utils.ExtractBoolAttr(httpResp, "radiusFallbackEnabled")
		if diags.HasError() {
			return diags
		}
	}

	// RadiusCoaEnabled
	if state.RadiusCoaEnabled.IsNull() || state.RadiusCoaEnabled.IsUnknown() {
		state.RadiusCoaEnabled, diags = utils.ExtractBoolAttr(httpResp, "radiusCoaEnabled")
		if diags.HasError() {
			return diags
		}
	}

	// RadiusAccountingInterimInterval
	if state.RadiusAccountingInterimInterval.IsNull() || state.RadiusAccountingInterimInterval.IsUnknown() {
		state.RadiusAccountingInterimInterval, diags = utils.ExtractInt64Attr(httpResp, "radiusAccountingInterimInterval")
		if diags.HasError() {
			return diags
		}
	}

	// UseVlanTagging
	if state.UseVlanTagging.IsNull() || state.UseVlanTagging.IsUnknown() {
		state.UseVlanTagging, diags = utils.ExtractBoolAttr(httpResp, "useVlanTagging")
		if diags.HasError() {
			return diags
		}
	}

	// ConcentratorNetworkID
	if state.ConcentratorNetworkID.IsNull() || state.ConcentratorNetworkID.IsUnknown() {
		state.ConcentratorNetworkID, diags = utils.ExtractStringAttr(httpResp, "concentratorNetworkId")
		if diags.HasError() {
			return diags
		}
	}

	// SecondaryConcentratorNetworkID
	if state.SecondaryConcentratorNetworkID.IsNull() || state.SecondaryConcentratorNetworkID.IsUnknown() {
		state.SecondaryConcentratorNetworkID, diags = utils.ExtractStringAttr(httpResp, "secondaryConcentratorNetworkId")
		if diags.HasError() {
			return diags
		}
	}

	// DisassociateClientsOnVpnFailOver
	if state.DisassociateClientsOnVpnFailOver.IsNull() || state.DisassociateClientsOnVpnFailOver.IsUnknown() {
		state.DisassociateClientsOnVpnFailOver, diags = utils.ExtractBoolAttr(httpResp, "disassociateClientsOnVpnFailOver")
		if diags.HasError() {
			return diags
		}
	}

	// VlanID
	if state.VlanID.IsNull() || state.VlanID.IsUnknown() {
		state.VlanID, diags = utils.ExtractInt64Attr(httpResp, "vlanId")
		if diags.HasError() {
			return diags
		}
	}

	// DefaultVlanID
	if state.DefaultVlanID.IsNull() || state.DefaultVlanID.IsUnknown() {
		state.DefaultVlanID, diags = utils.ExtractInt64Attr(httpResp, "defaultVlanId")
		if diags.HasError() {
			return diags
		}
	}

	// ApTagsAndVlanIDs
	if state.ApTagsAndVlanIDs.IsNull() || state.ApTagsAndVlanIDs.IsUnknown() {
		state.ApTagsAndVlanIDs, diags = NetworksWirelessSsidStateApTagsAndVlanIds(httpResp)
		if diags.HasError() {
			return diags
		}
	}

	// GRE
	if state.GRE.IsNull() || state.GRE.IsUnknown() {
		greAttrs := map[string]attr.Type{
			"concentrator": types.ObjectType{AttrTypes: map[string]attr.Type{
				"host": types.StringType,
			}},
			"key": types.Int64Type,
		}

		state.GRE, diags = utils.ExtractObjectAttr(httpResp, "gre", greAttrs)
		if diags.HasError() {
			return diags
		}
	}

	// RadiusOverride
	if state.RadiusOverride.IsNull() || state.RadiusOverride.IsUnknown() {
		state.RadiusOverride, diags = utils.ExtractBoolAttr(httpResp, "radiusOverride")
		if diags.HasError() {
			return diags
		}
	}

	// RadiusGuestVlanEnabled
	if state.RadiusGuestVlanEnabled.IsNull() || state.RadiusGuestVlanEnabled.IsUnknown() {
		state.RadiusGuestVlanEnabled, diags = utils.ExtractBoolAttr(httpResp, "radiusGuestVlanEnabled")
		if diags.HasError() {
			return diags
		}
	}

	// RadiusGuestVlanID
	if state.RadiusGuestVlanID.IsNull() || state.RadiusGuestVlanID.IsUnknown() {
		state.RadiusGuestVlanID, diags = utils.ExtractInt64Attr(httpResp, "radiusGuestVlanId")
		if diags.HasError() {
			return diags
		}
	}

	// BandSelection
	if state.BandSelection.IsNull() || state.BandSelection.IsUnknown() {
		state.BandSelection, diags = utils.ExtractStringAttr(httpResp, "bandSelection")
		if diags.HasError() {
			return diags
		}
	}

	// LanIsolationEnabled
	if state.LanIsolationEnabled.IsNull() || state.LanIsolationEnabled.IsUnknown() {
		state.LanIsolationEnabled, diags = utils.ExtractBoolAttr(httpResp, "lanIsolationEnabled")
		if diags.HasError() {
			return diags
		}
	}

	// AdultContentFilteringEnabled
	if state.AdultContentFilteringEnabled.IsNull() || state.AdultContentFilteringEnabled.IsUnknown() {
		state.AdultContentFilteringEnabled, diags = utils.ExtractBoolAttr(httpResp, "adultContentFilteringEnabled")
		if diags.HasError() {
			return diags
		}
	}

	// DnsRewrite
	if state.DnsRewrite.IsNull() || state.DnsRewrite.IsUnknown() {
		dnsRewriteAttrs := map[string]attr.Type{
			"enabled":                 types.BoolType,
			"dns_custom_name_servers": types.ListType{ElemType: types.StringType},
		}

		if val, ok := httpResp["dnsRewrite"].(map[string]interface{}); ok && val != nil {
			enabled, _ := val["enabled"].(bool)
			dnsCustomNameServersInterface, _ := val["dns_custom_name_servers"].([]interface{})

			dnsCustomNameServers := make([]attr.Value, len(dnsCustomNameServersInterface))
			for i, dns := range dnsCustomNameServersInterface {
				dnsCustomNameServers[i] = types.StringValue(dns.(string))
			}

			dnsCustomNameServersList, diags := types.ListValue(types.StringType, dnsCustomNameServers)
			if diags.HasError() {
				return diags
			}

			dnsRewrite, diags := types.ObjectValue(
				dnsRewriteAttrs,
				map[string]attr.Value{
					"enabled":                 types.BoolValue(enabled),
					"dns_custom_name_servers": dnsCustomNameServersList,
				},
			)
			if diags.HasError() {
				return diags
			}

			state.DnsRewrite = dnsRewrite
		} else {
			state.DnsRewrite = types.ObjectNull(dnsRewriteAttrs)
		}
		if diags.HasError() {
			return diags
		}
	}

	// SpeedBurst
	if state.SpeedBurst.IsNull() || state.SpeedBurst.IsUnknown() {

		speedBurstAttrs := map[string]attr.Type{
			"enabled": types.BoolType,
		}

		if value, ok := httpResp["speedBurst"].(map[string]interface{}); ok {

			if enable, enableOk := value["enabled"].(bool); enableOk {
				s := SpeedBurst{Enabled: types.BoolValue(enable)}
				objVal, err := types.ObjectValueFrom(context.Background(), speedBurstAttrs, s)
				if err.HasError() {
					diags.Append(err...)
				}

				state.SpeedBurst = objVal
				if diags.HasError() {
					return diags
				}
			} else {
				state.SpeedBurst = types.ObjectNull(speedBurstAttrs)
			}

		} else {
			state.SpeedBurst = types.ObjectNull(speedBurstAttrs)
		}
	}

	// SsidAdminAccessible
	if state.SsidAdminAccessible.IsNull() || state.SsidAdminAccessible.IsUnknown() {
		state.SsidAdminAccessible, diags = utils.ExtractBoolAttr(httpResp, "ssidAdminAccessible")
		if diags.HasError() {
			return diags
		}
	}

	// RadiusEnabled
	if state.RadiusEnabled.IsNull() || state.RadiusEnabled.IsUnknown() {
		state.RadiusEnabled, diags = utils.ExtractBoolAttr(httpResp, "radiusEnabled")
		if diags.HasError() {
			return diags
		}
	}

	// NamedVlans
	if state.NamedVlans.IsNull() || state.NamedVlans.IsUnknown() {
		namedVlansAttrs := map[string]attr.Type{
			"tagging": types.ObjectType{AttrTypes: map[string]attr.Type{
				"enabled":           types.BoolType,
				"default_vlan_name": types.StringType,
				"by_ap_tags": types.ListType{ElemType: types.ObjectType{AttrTypes: map[string]attr.Type{
					"tags":      types.ListType{ElemType: types.StringType},
					"vlan_name": types.StringType,
				}}},
			}},
			"radius": types.ObjectType{AttrTypes: map[string]attr.Type{
				"guest_vlan": types.ObjectType{AttrTypes: map[string]attr.Type{
					"enabled": types.BoolType,
					"name":    types.StringType,
				}},
			}},
		}

		state.NamedVlans, diags = utils.ExtractObjectAttr(httpResp, "namedVlans", namedVlansAttrs)
		if diags.HasError() {
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

	if !plan.RadiusGuestVlanID.IsNull() && !plan.RadiusGuestVlanID.IsUnknown() {
		radiusGuestVlanId, err := utils.Int32Pointer(plan.RadiusGuestVlanID.ValueInt64())
		if err.HasError() {
			tflog.Error(context.Background(), fmt.Sprintf("%s", err.Errors()))
			diags.Append(err...)
		}
		payload.RadiusGuestVlanId = radiusGuestVlanId
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
	resp.TypeName = req.ProviderTypeName + "_networks_wireless_ssids"
}

func (r *NetworksWirelessSsidsResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "NetworksWirelessSsids",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The unique identifier for the resource, generated by the Meraki API.",
				Computed:    true,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier of the Meraki network to which the SSID belongs.",
				Required:            true,
			},
			"number": schema.Int64Attribute{
				MarkdownDescription: "Represents the SSID's order or position within the network's list of SSIDs.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the SSID.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Determines if the SSID is active and broadcasting.",
				Optional:            true,
				Computed:            true,
				Default:             utils.NewBoolDefault(true),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"auth_mode": schema.StringAttribute{
				MarkdownDescription: "Specifies the authentication method for the SSID.",
				Optional:            true,
				Computed:            true,
				Default:             utils.NewStringDefault("open"),
				Validators: []validator.String{
					stringvalidator.OneOf("8021x-google", "8021x-localradius", "8021x-meraki", "8021x-nac", "8021x-radius",
						"ipsk-with-nac", "ipsk-with-radius", "ipsk-without-radius", "open",
						"open-enhanced", "open-with-nac", "open-with-radius", "psk"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"enterprise_admin_access": schema.StringAttribute{
				MarkdownDescription: "Controls whether the SSID is accessible by enterprise administrators.",
				Optional:            true,
				Computed:            true,
				//Default:             utils.NewStringDefault("access disabled"),
				Validators: []validator.String{
					stringvalidator.OneOf("access disabled", "access enabled"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"encryption_mode": schema.StringAttribute{
				MarkdownDescription: "Defines the type of PSK encryption (e.g., WEP, WPA) used for securing wireless network data. This param is only valid if the authMode is 'psk'",
				Optional:            true,
				Computed:            true,
				Default:             utils.NewStringDefault("open"),
				Validators: []validator.String{
					stringvalidator.OneOf("open", "wep", "wpa", "wpa-eap"),
				},
			},
			"psk": schema.StringAttribute{
				MarkdownDescription: "The Pre-shared Key for the SSID. This param is only valid if the authMode is 'psk'",
				Optional:            true,
				Sensitive:           true,
				PlanModifiers: []planmodifier.String{
					utils.RequiresReplaceIfSensitive{},
				},
			},
			"wpa_encryption_mode": schema.StringAttribute{
				MarkdownDescription: "Specifies the WPA encryption mode.",
				Optional:            true,
				Computed:            true,
				//Default:             utils.NewStringDefault("WPA1 and WPA2"),
				Validators: []validator.String{
					stringvalidator.OneOf("WPA1 only", "WPA1 and WPA2", "WPA2 only", "WPA3 Transition Mode", "WPA3 only", "WPA3 192-bit Security"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"dot11w": schema.SingleNestedAttribute{
				MarkdownDescription: "Configuration for Protected Management Frames (802.11w).",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "Whether 802.11w is enabled or not.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"required": schema.BoolAttribute{
						MarkdownDescription: "Whether 802.11w is required or not.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"dot11r": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Settings for 802.11r, used for fast roaming.",
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "Whether 802.11r is enabled or not.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"adaptive": schema.BoolAttribute{
						MarkdownDescription: "Whether 802.11r is adaptive or not.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"splash_page": schema.StringAttribute{
				MarkdownDescription: "Defines the splash page type used for guest access management and authentication.",
				Optional:            true,
				Computed:            true,
				Default:             utils.NewStringDefault("None"),
				Validators: []validator.String{
					stringvalidator.OneOf("Billing", "Cisco ISE", "Click-through splash page", "Facebook Wi-Fi", "Google Apps domain",
						"Google OAuth", "None", "Password-protected with Active Directory", "Password-protected with LDAP", "Password-protected with Meraki RADIUS",
						"Password-protected with custom RADIUS", "SMS authentication", "Sponsored guest", "Systems Manager Sentry"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"splash_guest_sponsor_domains": schema.ListAttribute{
				MarkdownDescription: "A list of email domains allowed to sponsor guest access.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"oauth": schema.SingleNestedAttribute{
				MarkdownDescription: "Configures OAuth settings for integrating third-party authentication providers.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				Attributes: map[string]schema.Attribute{
					"allowed_domains": schema.ListAttribute{
						MarkdownDescription: "List of allowed domains for OAuth.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
						ElementType: types.StringType,
					},
				},
			},
			"local_radius": schema.SingleNestedAttribute{
				MarkdownDescription: "Local RADIUS server configuration for authentication. Only valid if authMode is '8021x-localradius'",
				Optional:            true,
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"cache_timeout": schema.Int64Attribute{
						MarkdownDescription: "The duration (in seconds) for which LDAP and OCSP lookups are cached.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"password_authentication": schema.SingleNestedAttribute{
						MarkdownDescription: "Password-based authentication settings.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								MarkdownDescription: "Whether to use password-based authentication.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
						},
					},
					"certificate_authentication": schema.SingleNestedAttribute{
						MarkdownDescription: "Certificate verification settings.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								MarkdownDescription: "Whether to use certificate-based authentication.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"use_ldap": schema.BoolAttribute{
								MarkdownDescription: "Whether to verify the certificate with LDAP.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"use_ocsp": schema.BoolAttribute{
								MarkdownDescription: "Whether to verify the certificate with OCSP.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"ocsp_responder_url": schema.StringAttribute{
								MarkdownDescription: "The URL of the OCSP responder to verify client certificate status.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"client_root_ca_certificate": schema.SingleNestedAttribute{
								MarkdownDescription: "The Client CA Certificate used to sign the client certificate.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.Object{
									objectplanmodifier.UseStateForUnknown(),
								},
								Attributes: map[string]schema.Attribute{
									"contents": schema.StringAttribute{
										MarkdownDescription: "The contents of the Client CA Certificate.",
										Optional:            true,
										Computed:            true,
										PlanModifiers: []planmodifier.String{
											stringplanmodifier.UseStateForUnknown(),
										},
									},
								},
							},
						},
					},
				},
			},
			"ldap": schema.SingleNestedAttribute{
				MarkdownDescription: "LDAP server configuration for authentication. Only valid if splashPage is 'Password-protected with LDAP'.",
				Optional:            true,
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"servers": schema.ListNestedAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The LDAP servers to be used for authentication.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"host": schema.StringAttribute{
									MarkdownDescription: "The LDAP server host.",
									Optional:            true,
									Computed:            true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
									},
								},
								"port": schema.Int64Attribute{
									MarkdownDescription: "The LDAP server port.",
									Optional:            true,
									Computed:            true,
									PlanModifiers: []planmodifier.Int64{
										int64planmodifier.UseStateForUnknown(),
									},
								},
							},
						},
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
					},
					"credentials": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "The credentials for LDAP server authentication.",
						Attributes: map[string]schema.Attribute{
							"distinguished_name": schema.StringAttribute{
								MarkdownDescription: "The distinguished name for LDAP.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"password": schema.StringAttribute{
								MarkdownDescription: "The password for LDAP.",
								Optional:            true,
								Sensitive:           true,
								PlanModifiers: []planmodifier.String{
									utils.RequiresReplaceIfSensitive{},
								},
							},
						},
					},
					"base_distinguished_name": schema.StringAttribute{
						MarkdownDescription: "The base distinguished name on the LDAP server.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"server_ca_certificate": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "The CA certificate for the LDAP server.",
						Attributes: map[string]schema.Attribute{
							"contents": schema.StringAttribute{
								MarkdownDescription: "The contents of the CA certificate.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
						},
					},
				},
			},
			"active_directory": schema.SingleNestedAttribute{
				MarkdownDescription: "Sets up Active Directory for authentication. Only valid if splashPage is 'Password-protected with Active Directory'",
				Optional:            true,
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"servers": schema.ListNestedAttribute{
						Optional: true,
						Computed: true,
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
						MarkdownDescription: "The Active Directory servers to be used for authentication.",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"host": schema.StringAttribute{
									MarkdownDescription: "The Active Directory server host.",
									Optional:            true,
									Computed:            true,
									PlanModifiers: []planmodifier.String{
										stringplanmodifier.UseStateForUnknown(),
									},
								},
								"port": schema.Int64Attribute{
									MarkdownDescription: "The Active Directory server port.",
									Optional:            true,
									Computed:            true,
									PlanModifiers: []planmodifier.Int64{
										int64planmodifier.UseStateForUnknown(),
									},
								},
							},
						},
					},
					"credentials": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "The credentials for Active Directory server authentication.",
						Attributes: map[string]schema.Attribute{
							"logon_name": schema.StringAttribute{
								MarkdownDescription: "The logon name for Active Directory.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"password": schema.StringAttribute{
								MarkdownDescription: "The password for Active Directory.",
								Optional:            true,
								Sensitive:           true,
								PlanModifiers: []planmodifier.String{
									utils.RequiresReplaceIfSensitive{},
								},
							},
						},
					},
				},
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
			},
			"radius_servers": schema.ListNestedAttribute{
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "A list of RADIUS servers for authentication. This param is only valid if the authMode is 'open-with-radius', '8021x-radius' or 'ipsk-with-radius'",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"host": schema.StringAttribute{
							MarkdownDescription: "IP address or hostname of the RADIUS server.",
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"port": schema.Int64Attribute{
							MarkdownDescription: "Port number on which the RADIUS server is listening.",
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.Int64{
								int64planmodifier.UseStateForUnknown(),
							},
						},
						"secret": schema.StringAttribute{
							MarkdownDescription: "Shared secret for the RADIUS server.",
							Optional:            true,
							Sensitive:           true,
							PlanModifiers: []planmodifier.String{
								utils.RequiresReplaceIfSensitive{},
							},
						},
						"rad_sec_enabled": schema.BoolAttribute{
							MarkdownDescription: "Indicates whether RADSEC is enabled.",
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
						"open_roaming_certificate_id": schema.Int64Attribute{
							MarkdownDescription: "OpenRoaming certificate Id associated with the RADIUS server.",
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.Int64{
								int64planmodifier.UseStateForUnknown(),
							},
						},
						"ca_certificate": schema.StringAttribute{
							MarkdownDescription: "CA certificate for the RADIUS server.",
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
					},
				},
			},
			"radius_proxy_enabled": schema.BoolAttribute{
				MarkdownDescription: "Indicates if a RADIUS proxy is enabled.",
				Optional:            true,
				Computed:            true,
				Default:             utils.NewBoolDefault(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"radius_testing_enabled": schema.BoolAttribute{
				MarkdownDescription: "Enables or disables testing for RADIUS server configurations.",
				Optional:            true,
				Computed:            true,
				Default:             utils.NewBoolDefault(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"radius_called_station_id": schema.StringAttribute{
				MarkdownDescription: "Specifies the template for the called station identifier used in RADIUS interactions.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"radius_authentication_nas_id": schema.StringAttribute{
				MarkdownDescription: "Defines the NAS identifier template for RADIUS authentication purposes.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"radius_server_timeout": schema.Int64Attribute{
				MarkdownDescription: "The duration a RADIUS client will wait for a reply from the RADIUS server.",
				Optional:            true,
				Computed:            true,
				Default:             utils.NewInt64Default(3),
				Validators: []validator.Int64{
					int64validator.Between(1, 10),
				},
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"radius_server_attempts_limit": schema.Int64Attribute{
				MarkdownDescription: "Maximum number of retry attempts for RADIUS server authentication before failing over.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
				Default: utils.NewInt64Default(3),
				Validators: []validator.Int64{
					int64validator.Between(1, 5),
				},
			},
			"radius_fallback_enabled": schema.BoolAttribute{
				MarkdownDescription: "Determines if RADIUS fallback is activated.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"radius_coa_enabled": schema.BoolAttribute{
				MarkdownDescription: "Controls the usage of RADIUS Change of Authorization (CoA) for dynamic modifications to a client's session.",
				Optional:            true,
				Computed:            true,
				Default:             utils.NewBoolDefault(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"radius_fail_over_policy": schema.StringAttribute{
				MarkdownDescription: "This policy determines how authentication requests should be handled in the event that all of the configured RADIUS servers are unreachable.",
				Optional:            true,
				Computed:            true,
				Default:             utils.NewStringDefault("Deny access"),
				Validators: []validator.String{
					stringvalidator.OneOf("Allow access", "Deny access"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"radius_load_balancing_policy": schema.StringAttribute{
				MarkdownDescription: "This policy determines which RADIUS server will be contacted first in an authentication attempt and the ordering of any necessary retry attempts.",
				Optional:            true,
				Computed:            true,
				Default:             utils.NewStringDefault("Round robin"),
				Validators: []validator.String{
					stringvalidator.OneOf("Round robin", "Strict priority order"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"radius_accounting_enabled": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether RADIUS accounting is enabled.",
				Optional:            true,
				Computed:            true,
				Default:             utils.NewBoolDefault(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"radius_accounting_servers": schema.ListNestedAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "The RADIUS accounting servers to be used for accounting services.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"host": schema.StringAttribute{
							MarkdownDescription: "IP address or hostname of the RADIUS accounting server.",
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"port": schema.Int64Attribute{
							MarkdownDescription: "Port number on which the RADIUS accounting server is listening.",
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.Int64{
								int64planmodifier.UseStateForUnknown(),
							},
						},
						"secret": schema.StringAttribute{
							MarkdownDescription: "Shared secret for the RADIUS accounting server.",
							Optional:            true,
							Sensitive:           true,
							PlanModifiers: []planmodifier.String{
								utils.RequiresReplaceIfSensitive{},
							},
						},
						"ca_certificate": schema.StringAttribute{
							MarkdownDescription: "CA certificate for the RADIUS accounting server.",
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.String{
								stringplanmodifier.UseStateForUnknown(),
							},
						},
						"rad_sec_enabled": schema.BoolAttribute{
							MarkdownDescription: "Indicates whether RADSEC (RADIUS over TLS) is enabled for secure communication with the RADIUS accounting server.",
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.Bool{
								boolplanmodifier.UseStateForUnknown(),
							},
						},
						"open_roaming_certificate_id": schema.Int64Attribute{
							MarkdownDescription: "The Open Roaming Certificate Id.",
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.Int64{
								int64planmodifier.UseStateForUnknown(),
							},
						},
					},
				},
			},
			"radius_accounting_interim_interval": schema.Int64Attribute{
				MarkdownDescription: "The interval (in seconds) in which accounting information is updated and sent to the RADIUS accounting server.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"radius_attribute_for_group_policies": schema.StringAttribute{
				MarkdownDescription: "Specify the RADIUS attribute used to look up group policies.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Default: utils.NewStringDefault("Reply-Message"),
				Validators: []validator.String{
					stringvalidator.OneOf("Airespace-ACL-Name", "Aruba-User-Role", "Filter-Id", "Reply-Message"),
				},
			},
			"ip_assignment_mode": schema.StringAttribute{
				MarkdownDescription: "The client IP assignment mode ('NAT mode', 'Bridge mode', 'Layer 3 roaming', 'Ethernet over GRE', 'Layer 3 roaming with a concentrator' or 'VPN')",
				Optional:            true,
				Computed:            true,
				Default:             utils.NewStringDefault("Bridge mode"),
				Validators: []validator.String{
					stringvalidator.OneOf("NAT mode", "Bridge mode", "Layer 3 roaming", "Ethernet over GRE", "Layer 3 roaming with a concentrator", "VPN"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"use_vlan_tagging": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether VLAN tagging is used.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"concentrator_network_id": schema.StringAttribute{
				MarkdownDescription: "The concentrator to use when the ipAssignmentMode is 'Layer 3 roaming with a concentrator' or 'VPN'.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"secondary_concentrator_network_id": schema.StringAttribute{
				MarkdownDescription: "The secondary concentrator to use when the ipAssignmentMode is 'VPN'.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"disassociate_clients_on_vpn_fail_over": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether clients should be disassociated during VPN failover.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"vlan_id": schema.Int64Attribute{
				MarkdownDescription: "The VLAN Id used for VLAN tagging.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"default_vlan_id": schema.Int64Attribute{
				MarkdownDescription: "The default VLAN Id used for 'all other APs'.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"ap_tags_and_vlan_ids": schema.ListNestedAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "A set of AP tags and corresponding VLAN IDs.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"tags": schema.ListAttribute{
							MarkdownDescription: "Array of AP tags.",
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.List{
								listplanmodifier.UseStateForUnknown(),
							},
							ElementType: types.StringType,
						},
						"vlan_id": schema.Int64Attribute{
							MarkdownDescription: "VLAN Id associated with the AP tags.",
							Optional:            true,
							Computed:            true,
							PlanModifiers: []planmodifier.Int64{
								int64planmodifier.UseStateForUnknown(),
							},
						},
					},
				},
			},
			"walled_garden_enabled": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether a walled garden is enabled for the SSID.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"walled_garden_ranges": schema.ListAttribute{
				MarkdownDescription: "List of Walled Garden ranges.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"gre": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "GRE (Generic Routing Encapsulation) tunnel configuration.",
				Attributes: map[string]schema.Attribute{
					"concentrator": schema.SingleNestedAttribute{
						Optional: true,
						Computed: true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						MarkdownDescription: "GRE tunnel concentrator configuration.",
						Attributes: map[string]schema.Attribute{
							"host": schema.StringAttribute{
								MarkdownDescription: "The GRE concentrator host.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
						},
					},
					"key": schema.Int64Attribute{
						MarkdownDescription: "The GRE key.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"radius_override": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether RADIUS attributes can override other settings.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"radius_guest_vlan_enabled": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether the RADIUS guest VLAN is enabled.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"radius_guest_vlan_id": schema.Int64Attribute{
				MarkdownDescription: "VLAN Id of the RADIUS Guest VLAN.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"min_bit_rate": schema.Float64Attribute{
				MarkdownDescription: "The minimum bitrate in Mbps of this SSID in the default indoor RF profile.",
				Optional:            true,
				Computed:            true,
				Default:             utils.NewFloat64Default(1),
				PlanModifiers: []planmodifier.Float64{
					float64planmodifier.UseStateForUnknown(),
				},
			},
			"band_selection": schema.StringAttribute{
				MarkdownDescription: "The client-serving radio frequencies of this SSID in the default indoor RF profile.",
				Optional:            true,
				Computed:            true,
				Default:             utils.NewStringDefault("Dual band operation"),
				Validators: []validator.String{
					stringvalidator.OneOf("Dual band operation", "5 GHz band only", "Dual band operation with Band Steering", "2.4 GHz band only"),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"per_client_bandwidth_limit_up": schema.Int64Attribute{
				MarkdownDescription: "The upload bandwidth limit in Kbps.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"per_client_bandwidth_limit_down": schema.Int64Attribute{
				MarkdownDescription: "The download bandwidth limit in Kbps.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"per_ssid_bandwidth_limit_up": schema.Int64Attribute{
				MarkdownDescription: "The total upload bandwidth limit in Kbps.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"per_ssid_bandwidth_limit_down": schema.Int64Attribute{
				MarkdownDescription: "The total download bandwidth limit in Kbps.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.UseStateForUnknown(),
				},
			},
			"lan_isolation_enabled": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether LAN isolation is enabled.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"visible": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether the SSID is visible.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"available_on_all_aps": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether the SSID is available on all access points.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"availability_tags": schema.ListAttribute{
				MarkdownDescription: "List of availability tags for the SSID.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
				ElementType: types.StringType,
			},
			"mandatory_dhcp_enabled": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether mandatory DHCP is enabled.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"adult_content_filtering_enabled": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether adult content filtering is enabled.",
				Optional:            true,
				Computed:            true,
				Default:             utils.NewBoolDefault(false),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"dns_rewrite": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "DNS rewrite configuration.",
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "Indicates whether DNS rewrite is enabled.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
					"dns_custom_name_servers": schema.ListAttribute{
						MarkdownDescription: "List of custom DNS nameservers.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.List{
							listplanmodifier.UseStateForUnknown(),
						},
						ElementType: types.StringType,
					},
				},
			},
			"speed_burst": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "Speed burst configuration.",
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "Indicates whether speed burst is enabled.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Bool{
							boolplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
			"ssid_admin_accessible": schema.BoolAttribute{
				MarkdownDescription: "SSID Administrator access status.",
				Optional:            true,
				Computed:            true,
				Default:             utils.NewBoolDefault(true),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"local_auth": schema.BoolAttribute{
				MarkdownDescription: "Extended local auth flag for Enterprise NAC.",
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"radius_enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether RADIUS authentication is enabled.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
			"admin_splash_url": schema.StringAttribute{
				MarkdownDescription: "URL for the admin splash page.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"splash_timeout": schema.StringAttribute{
				MarkdownDescription: "Splash page timeout.",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"named_vlans": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Object{
					objectplanmodifier.UseStateForUnknown(),
				},
				MarkdownDescription: "Configuration for named VLANs.",
				Attributes: map[string]schema.Attribute{
					"tagging": schema.SingleNestedAttribute{
						Optional: true,
						Computed: true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						MarkdownDescription: "Tagging configuration for named VLANs.",
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								MarkdownDescription: "Indicates whether VLAN tagging is enabled.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.Bool{
									boolplanmodifier.UseStateForUnknown(),
								},
							},
							"default_vlan_name": schema.StringAttribute{
								MarkdownDescription: "The default VLAN name.",
								Optional:            true,
								Computed:            true,
								PlanModifiers: []planmodifier.String{
									stringplanmodifier.UseStateForUnknown(),
								},
							},
							"by_ap_tags": schema.ListNestedAttribute{
								Optional: true,
								Computed: true,
								PlanModifiers: []planmodifier.List{
									listplanmodifier.UseStateForUnknown(),
								},
								MarkdownDescription: "Sets of AP tags and corresponding VLAN names.",
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"tags": schema.ListAttribute{
											MarkdownDescription: "Array of AP tags.",
											Optional:            true,
											Computed:            true,
											PlanModifiers: []planmodifier.List{
												listplanmodifier.UseStateForUnknown(),
											},
											ElementType: types.StringType,
										},
										"vlan_name": schema.StringAttribute{
											MarkdownDescription: "VLAN name associated with the AP tags.",
											Optional:            true,
											Computed:            true,
											PlanModifiers: []planmodifier.String{
												stringplanmodifier.UseStateForUnknown(),
											},
										},
									},
								},
							},
						},
					},
					"radius": schema.SingleNestedAttribute{
						Optional: true,
						Computed: true,
						PlanModifiers: []planmodifier.Object{
							objectplanmodifier.UseStateForUnknown(),
						},
						MarkdownDescription: "RADIUS configuration for named VLANs.",
						Attributes: map[string]schema.Attribute{
							"guest_vlan": schema.SingleNestedAttribute{
								Optional: true,
								Computed: true,
								PlanModifiers: []planmodifier.Object{
									objectplanmodifier.UseStateForUnknown(),
								},
								MarkdownDescription: "Guest VLAN configuration for RADIUS.",
								Attributes: map[string]schema.Attribute{
									"enabled": schema.BoolAttribute{
										MarkdownDescription: "Indicates whether the RADIUS guest VLAN is enabled.",
										Optional:            true,
										Computed:            true,
										PlanModifiers: []planmodifier.Bool{
											boolplanmodifier.UseStateForUnknown(),
										},
									},
									"name": schema.StringAttribute{
										MarkdownDescription: "Name of the RADIUS guest VLAN.",
										Optional:            true,
										Computed:            true,
										PlanModifiers: []planmodifier.String{
											stringplanmodifier.UseStateForUnknown(),
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

	// API call function to be passed to retryOn4xx
	apiCall := func() (*openApiClient.GetNetworkWirelessSsids200ResponseInner, *http.Response, error) {
		return r.client.WirelessApi.UpdateNetworkWirelessSsid(context.Background(), plan.NetworkId.ValueString(), fmt.Sprint(plan.Number.ValueInt64())).UpdateNetworkWirelessSsidRequest(payload).Execute()
	}

	// Use retryOn4xx for the API call as the meraki API backend returns HTTP 400 messages as a result of collision issues with rapid creation of postgres GroupPolicyIds.
	inlineResp, httpResp, err := tools.CustomHttpRequestRetryStronglyTyped(ctx, maxRetries, retryDelay, apiCall)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating wireless ssids",
			fmt.Sprintf("Could not create group policy, unexpected error: %s", err),
		)

		if httpResp != nil {
			var responseBody string
			if httpResp != nil && httpResp.Body != nil {
				bodyBytes, readErr := io.ReadAll(httpResp.Body)
				if readErr == nil {
					responseBody = string(bodyBytes)
				}
			}
			tflog.Error(ctx, "Failed to create resource", map[string]interface{}{
				"error":          err.Error(),
				"httpStatusCode": httpResp.StatusCode,
				"responseBody":   responseBody,
			})
			resp.Diagnostics.AddError(
				"Error creating group policy",
				fmt.Sprintf("HTTP Response: %v\nResponse Body: %s", httpResp, responseBody),
			)
		}
		return
	}

	tflog.Info(ctx, "ssid created successfully", map[string]interface{}{
		"name":   plan.Name.ValueString(),
		"number": plan.Number.ValueInt64(),
		"vlanId": plan.VlanID.ValueInt64(),
	})

	// Read the response body
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		diags.AddError("Error reading response body", err.Error())
	}

	// Parse the response body into a map
	var responseData map[string]interface{}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		diags.AddError("Error unmarshalling response body", err.Error())
	}

	diags = updateNetworksWirelessSsidsResourceState(ctx, &plan, inlineResp, responseData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	tflog.Trace(ctx, "Completed create operation for NetworksWirelessSsidsResource")
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

	// API call function to be passed to retryOn4xx
	apiCall := func() (*openApiClient.GetNetworkWirelessSsids200ResponseInner, *http.Response, error) {
		return r.client.WirelessApi.GetNetworkWirelessSsid(context.Background(), state.NetworkId.ValueString(), fmt.Sprint(state.Number.ValueInt64())).Execute()
	}

	// Use retryOn4xx for the API call as the meraki API backend returns HTTP 400 messages as a result of collision issues with rapid creation of postgres GroupPolicyIds.
	inlineResp, httpResp, err := tools.CustomHttpRequestRetryStronglyTyped(ctx, maxRetries, retryDelay, apiCall)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating wireless ssids",
			fmt.Sprintf("Could not create group policy, unexpected error: %s", err),
		)

		if httpResp != nil {
			var responseBody string
			if httpResp != nil && httpResp.Body != nil {
				bodyBytes, readErr := io.ReadAll(httpResp.Body)
				if readErr == nil {
					responseBody = string(bodyBytes)
				}
			}
			tflog.Error(ctx, "Failed to create resource", map[string]interface{}{
				"error":          err.Error(),
				"httpStatusCode": httpResp.StatusCode,
				"responseBody":   responseBody,
			})
			resp.Diagnostics.AddError(
				"Error creating group policy",
				fmt.Sprintf("HTTP Response: %v\nResponse Body: %s", httpResp, responseBody),
			)
		}
		return
	}

	// If there was an error during API call, add it to diagnostics.
	if err != nil {
		resp.Diagnostics.AddError(
			"Read HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Read Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// If there were any errors up to this point, log the state and return.
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("State Data", fmt.Sprintf("\n%v", state))
		return
	}

	// Read the response body
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		resp.Diagnostics.AddError("Error reading response body", err.Error())
	}

	// Parse the response body into a map
	var responseData map[string]interface{}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		resp.Diagnostics.AddError("Error unmarshalling response body", err.Error())
	}

	diags = updateNetworksWirelessSsidsResourceState(ctx, &state, inlineResp, responseData)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *NetworksWirelessSsidsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan NetworksWirelessSsidResourceModel

	tflog.Trace(ctx, "Starting update operation for NetworksWirelessSsidsResource")

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

	// API call function to be passed to retryOn4xx
	apiCall := func() (*openApiClient.GetNetworkWirelessSsids200ResponseInner, *http.Response, error) {
		return r.client.WirelessApi.UpdateNetworkWirelessSsid(context.Background(), plan.NetworkId.ValueString(), fmt.Sprint(plan.Number.ValueInt64())).UpdateNetworkWirelessSsidRequest(payload).Execute()
	}

	// Use retryOn4xx for the API call as the meraki API backend returns HTTP 400 messages as a result of collision issues with rapid creation of postgres GroupPolicyIds.
	inlineResp, httpResp, err := tools.CustomHttpRequestRetryStronglyTyped(ctx, maxRetries, retryDelay, apiCall)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating wireless ssids",
			fmt.Sprintf("Could not create group policy, unexpected error: %s", err),
		)

		if httpResp != nil {
			var responseBody string
			if httpResp != nil && httpResp.Body != nil {
				bodyBytes, readErr := io.ReadAll(httpResp.Body)
				if readErr == nil {
					responseBody = string(bodyBytes)
				}
			}
			tflog.Error(ctx, "Failed to create resource", map[string]interface{}{
				"error":          err.Error(),
				"httpStatusCode": httpResp.StatusCode,
				"responseBody":   responseBody,
			})
			resp.Diagnostics.AddError(
				"Error creating group policy",
				fmt.Sprintf("HTTP Response: %v\nResponse Body: %s", httpResp, responseBody),
			)
		}
		return
	}

	// If there was an error during API call, add it to diagnostics.
	if err != nil {
		resp.Diagnostics.AddError(
			"Update HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
	}

	// Check for API success response code
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Update Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// If there were any errors up to this point, log the state plan and return.
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("State Data", fmt.Sprintf("\n%v", plan))
		return
	}

	tflog.Info(ctx, "ssid updated successfully", map[string]interface{}{
		"name":   plan.Name.ValueString(),
		"number": plan.Number.ValueInt64(),
		"vlanId": plan.VlanID.ValueInt64(),
	})

	// Read the response body
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		diags.AddError("Error reading response body", err.Error())
	}

	// Parse the response body into a map
	var responseData map[string]interface{}
	err = json.Unmarshal(body, &responseData)
	if err != nil {
		diags.AddError("Error unmarshalling response body", err.Error())
	}

	diags = updateNetworksWirelessSsidsResourceState(ctx, &plan, inlineResp, responseData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	tflog.Trace(ctx, "Completed create operation for NetworksWirelessSsidsResource")
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

	// API call function to be passed to retryOn4xx
	apiCall := func() (*openApiClient.GetNetworkWirelessSsids200ResponseInner, *http.Response, error) {
		return r.client.WirelessApi.UpdateNetworkWirelessSsid(context.Background(), state.NetworkId.ValueString(), fmt.Sprint(state.Number.ValueInt64())).UpdateNetworkWirelessSsidRequest(payload).Execute()
	}

	// Use retryOn4xx for the API call as the meraki API backend returns HTTP 400 messages as a result of collision issues with rapid creation of postgres GroupPolicyIds.
	_, httpResp, err := tools.CustomHttpRequestRetryStronglyTyped(ctx, maxRetries, retryDelay, apiCall)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating wireless ssids",
			fmt.Sprintf("Could not create group policy, unexpected error: %s", err),
		)

		if httpResp != nil {
			var responseBody string
			if httpResp != nil && httpResp.Body != nil {
				bodyBytes, readErr := io.ReadAll(httpResp.Body)
				if readErr == nil {
					responseBody = string(bodyBytes)
				}
			}
			tflog.Error(ctx, "Failed to create resource", map[string]interface{}{
				"error":          err.Error(),
				"httpStatusCode": httpResp.StatusCode,
				"responseBody":   responseBody,
			})
			resp.Diagnostics.AddError(
				"Error creating group policy",
				fmt.Sprintf("HTTP Response: %v\nResponse Body: %s", httpResp, responseBody),
			)
		}
		return
	}

	// Check for API success response code
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Delete Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// If there were any errors up to this point, log the state and return.
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("State Data", fmt.Sprintf("\n%v", state))
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
