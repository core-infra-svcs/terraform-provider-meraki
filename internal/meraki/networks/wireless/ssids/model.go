package ssids

import "github.com/hashicorp/terraform-plugin-framework/types"

// NetworksWirelessSsidResourceModel represents the internal structure for the Networks Wireless SSID resource
type NetworksWirelessSsidResourceModel struct {
	Id                               types.String `tfsdk:"id" json:"id"`
	NetworkId                        types.String `tfsdk:"network_id" json:"networkId"`
	Number                           types.Int64  `tfsdk:"number" json:"number"`
	Name                             types.String `tfsdk:"name" json:"name"`
	Enabled                          types.Bool   `tfsdk:"enabled" json:"enabled"`
	AuthMode                         types.String `tfsdk:"auth_mode" json:"authMode"`
	EnterpriseAdminAccess            types.String `tfsdk:"enterprise_admin_access" json:"enterpriseAdminAccess"`
	EncryptionMode                   types.String `tfsdk:"encryption_mode" json:"encryptionMode"`
	PSK                              types.String `tfsdk:"psk" json:"psk"`
	WPAEncryptionMode                types.String `tfsdk:"wpa_encryption_mode" json:"wpaEncryptionMode"`
	Dot11w                           types.Object `tfsdk:"dot11w" json:"dot11w"`
	Dot11r                           types.Object `tfsdk:"dot11r" json:"dot11r"`
	SplashPage                       types.String `tfsdk:"splash_page" json:"splashPage"`
	SplashGuestSponsorDomains        types.List   `tfsdk:"splash_guest_sponsor_domains" json:"splashGuestSponsorDomains"`
	OAuth                            types.Object `tfsdk:"oauth" json:"oauth"`
	LocalRadius                      types.Object `tfsdk:"local_radius" json:"localRadius"`
	LDAP                             types.Object `tfsdk:"ldap" json:"ldap"`
	ActiveDirectory                  types.Object `tfsdk:"active_directory" json:"activeDirectory"`
	RadiusServers                    types.List   `tfsdk:"radius_servers" json:"radiusServers"`
	RadiusProxyEnabled               types.Bool   `tfsdk:"radius_proxy_enabled" json:"radiusProxyEnabled"`
	RadiusTestingEnabled             types.Bool   `tfsdk:"radius_testing_enabled" json:"radiusTestingEnabled"`
	RadiusCalledStationID            types.String `tfsdk:"radius_called_station_id" json:"radiusCalledStationId"`
	RadiusAuthenticationNASID        types.String `tfsdk:"radius_authentication_nas_id" json:"radiusAuthenticationNasId"`
	RadiusServerTimeout              types.Int64  `tfsdk:"radius_server_timeout" json:"radiusServerTimeout"`
	RadiusServerAttemptsLimit        types.Int64  `tfsdk:"radius_server_attempts_limit" json:"radiusServerAttemptsLimit"`
	RadiusFallbackEnabled            types.Bool   `tfsdk:"radius_fallback_enabled" json:"radiusFallbackEnabled"`
	RadiusCoaEnabled                 types.Bool   `tfsdk:"radius_coa_enabled" json:"radiusCoaEnabled"`
	RadiusFailOverPolicy             types.String `tfsdk:"radius_fail_over_policy" json:"radiusFailoverPolicy"`
	RadiusLoadBalancingPolicy        types.String `tfsdk:"radius_load_balancing_policy" json:"radiusLoadBalancingPolicy"`
	RadiusAccountingEnabled          types.Bool   `tfsdk:"radius_accounting_enabled" json:"radiusAccountingEnabled"`
	RadiusAccountingServers          types.List   `tfsdk:"radius_accounting_servers" json:"radiusAccountingServers"`
	RadiusAccountingInterimInterval  types.Int64  `tfsdk:"radius_accounting_interim_interval" json:"radiusAccountingInterimInterval"`
	RadiusAttributeForGroupPolicies  types.String `tfsdk:"radius_attribute_for_group_policies" json:"radiusAttributeForGroupPolicies"`
	IPAssignmentMode                 types.String `tfsdk:"ip_assignment_mode" json:"ipAssignmentMode"`
	UseVlanTagging                   types.Bool   `tfsdk:"use_vlan_tagging" json:"useVlanTagging"`
	ConcentratorNetworkID            types.String `tfsdk:"concentrator_network_id" json:"concentratorNetworkId"`
	SecondaryConcentratorNetworkID   types.String `tfsdk:"secondary_concentrator_network_id" json:"secondaryConcentratorNetworkId"`
	DisassociateClientsOnVpnFailOver types.Bool   `tfsdk:"disassociate_clients_on_vpn_fail_over" json:"disassociateClientsOnVpnFailover"`
	VlanID                           types.Int64  `tfsdk:"vlan_id" json:"vlanId"`
	DefaultVlanID                    types.Int64  `tfsdk:"default_vlan_id" json:"defaultVlanId"`
	ApTagsAndVlanIDs                 types.List   `tfsdk:"ap_tags_and_vlan_ids" json:"apTagsAndVlanIds"`
	WalledGardenEnabled              types.Bool   `tfsdk:"walled_garden_enabled" json:"walledGardenEnabled"`
	WalledGardenRanges               types.List   `tfsdk:"walled_garden_ranges" json:"walledGardenRanges"`
	GRE                              types.Object `tfsdk:"gre" json:"gre"`
	RadiusOverride                   types.Bool   `tfsdk:"radius_override" json:"radiusOverride"`
	RadiusGuestVlanEnabled           types.Bool   `tfsdk:"radius_guest_vlan_enabled" json:"radiusGuestVlanEnabled"`
	RadiusGuestVlanId                types.Int64  `tfsdk:"radius_guest_vlan_id" json:"radiusGuestVlanId"`
	MinBitRate                       types.Int64  `tfsdk:"min_bitrate" json:"minBitRate"`
	BandSelection                    types.String `tfsdk:"band_selection" json:"bandSelection"`
	PerClientBandwidthLimitUp        types.Int64  `tfsdk:"per_client_bandwidth_limit_up" json:"perClientBandwidthLimitUp"`
	PerClientBandwidthLimitDown      types.Int64  `tfsdk:"per_client_bandwidth_limit_down" json:"perClientBandwidthLimitDown"`
	PerSsidBandwidthLimitUp          types.Int64  `tfsdk:"per_ssid_bandwidth_limit_up" json:"perSsidBandwidthLimitUp"`
	PerSsidBandwidthLimitDown        types.Int64  `tfsdk:"per_ssid_bandwidth_limit_down" json:"perSsidBandwidthLimitDown"`
	LanIsolationEnabled              types.Bool   `tfsdk:"lan_isolation_enabled" json:"lanIsolationEnabled"`
	Visible                          types.Bool   `tfsdk:"visible" json:"visible"`
	AvailableOnAllAps                types.Bool   `tfsdk:"available_on_all_aps" json:"availableOnAllAps"`
	AvailabilityTags                 types.List   `tfsdk:"availability_tags" json:"availabilityTags"`
	MandatoryDhcpEnabled             types.Bool   `tfsdk:"mandatory_dhcp_enabled" json:"mandatoryDhcpEnabled"`
	AdultContentFilteringEnabled     types.Bool   `tfsdk:"adult_content_filtering_enabled" json:"adultContentFilteringEnabled"`
	DnsRewrite                       types.Object `tfsdk:"dns_rewrite" json:"dnsRewrite"`
	SpeedBurst                       types.Object `tfsdk:"speed_burst" json:"speedBurst"`
	SsidAdminAccessible              types.Bool   `tfsdk:"ssid_admin_accessible" json:"ssidAdminAccessible"`
	LocalAuth                        types.Bool   `tfsdk:"local_auth" json:"localAuth"`
	RadiusEnabled                    types.Bool   `tfsdk:"radius_enabled" json:"radiusEnabled"`
	AdminSplashUrl                   types.String `tfsdk:"admin_splash_url" json:"adminSplashUrl"`
	SplashTimeout                    types.String `tfsdk:"splash_timeout" json:"splashTimeout"`
	NamedVlans                       types.Object `tfsdk:"named_vlans" json:"namedVlans"`
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
	//ServerId                 types.String `tfsdk:"server_id" json:"id"`  // not in api spec and changes all the time
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
