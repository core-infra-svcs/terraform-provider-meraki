package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
	"strconv"
	"strings"
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

// NetworksWirelessSsidsResourceModel describes the resource data model.
type NetworksWirelessSsidsResourceModel struct {
	Id                               jsontypes.String  `json:"id" tfsdk:"id"`
	NetworkID                        jsontypes.String  `json:"network_id" tfsdk:"network_id"`
	Number                           jsontypes.Int64   `json:"number" tfsdk:"number"`
	Name                             jsontypes.String  `json:"name" tfsdk:"name"`
	Enabled                          jsontypes.Bool    `json:"enabled" tfsdk:"enabled"`
	AuthMode                         jsontypes.String  `json:"authMode" tfsdk:"auth_mode"`
	EnterpriseAdminAccess            jsontypes.String  `json:"enterpriseAdminAccess" tfsdk:"enterprise_admin_access"`
	EncryptionMode                   jsontypes.String  `json:"encryptionMode" tfsdk:"encryption_mode"`
	Psk                              jsontypes.String  `json:"psk" tfsdk:"psk"`
	WpaEncryptionMode                jsontypes.String  `json:"wpaEncryptionMode" tfsdk:"wpa_encryption_mode"`
	Dot11W                           types.Object      `json:"dot11w" tfsdk:"dot11w"`
	Dot11R                           types.Object      `json:"dot11r" tfsdk:"dot11r"`
	SplashPage                       jsontypes.String  `json:"splashPage" tfsdk:"splash_page"`
	SplashGuestSponsorDomains        types.List        `json:"splashGuestSponsorDomains" tfsdk:"splash_guest_sponsor_domains"`
	Oauth                            types.Object      `json:"oauth" tfsdk:"oauth"`
	LocalRadius                      types.Object      `json:"localRadius" tfsdk:"local_radius"`
	Ldap                             types.Object      `json:"ldap" tfsdk:"ldap"`
	ActiveDirectory                  types.Object      `json:"activeDirectory" tfsdk:"active_directory"`
	RadiusServers                    types.List        `json:"radiusServers" tfsdk:"radius_servers"`
	RadiusProxyEnabled               jsontypes.Bool    `json:"radiusProxyEnabled" tfsdk:"radius_proxy_enabled"`
	RadiusTestingEnabled             jsontypes.Bool    `json:"radiusTestingEnabled" tfsdk:"radius_testing_enabled"`
	RadiusCalledStationId            jsontypes.String  `json:"radiusCalledStationId" tfsdk:"radius_called_station_id"`
	RadiusAuthenticationNasId        jsontypes.String  `json:"radiusAuthenticationNasId" tfsdk:"radius_authentication_nas_id"`
	RadiusServerTimeout              jsontypes.Int64   `json:"radiusServerTimeout" tfsdk:"radius_server_timeout"`
	RadiusServerAttemptsLimit        jsontypes.Int64   `json:"radiusServerAttemptsLimit" tfsdk:"radius_server_attempts_limit"`
	RadiusFallbackEnabled            jsontypes.Bool    `json:"radiusFallbackEnabled" tfsdk:"radius_fallback_enabled"`
	RadiusCoaEnabled                 jsontypes.Bool    `json:"radiusCoaEnabled" tfsdk:"radius_coa_enabled"`
	RadiusFailOverPolicy             jsontypes.String  `json:"radiusFailOverPolicy" tfsdk:"radius_fail_over_policy"`
	RadiusLoadBalancingPolicy        jsontypes.String  `json:"radiusLoadBalancingPolicy" tfsdk:"radius_load_balancing_policy"`
	RadiusAccountingEnabled          jsontypes.Bool    `json:"radiusAccountingEnabled" tfsdk:"radius_accounting_enabled"`
	RadiusAccountingServers          types.List        `json:"radiusAccountingServers" tfsdk:"radius_accounting_servers"`
	RadiusAccountingInterimInterval  jsontypes.Int64   `json:"radiusAccountingInterimInterval" tfsdk:"radius_accounting_interim_interval"`
	RadiusAttributeForGroupPolicies  jsontypes.String  `json:"radiusAttributeForGroupPolicies" tfsdk:"radius_attribute_for_group_policies"`
	IpAssignmentMode                 jsontypes.String  `json:"ipAssignmentMode" tfsdk:"ip_assignment_mode"`
	UseVlanTagging                   jsontypes.Bool    `json:"useVlanTagging" tfsdk:"use_vlan_tagging"`
	ConcentratorNetworkId            jsontypes.String  `json:"concentratorNetworkId" tfsdk:"concentrator_network_id"`
	SecondaryConcentratorNetworkId   jsontypes.String  `json:"secondaryConcentratorNetworkId" tfsdk:"secondary_concentrator_network_id"`
	DisassociateClientsOnVpnFailOver jsontypes.Bool    `json:"disassociateClientsOnVpnFailOver" tfsdk:"disassociate_clients_on_vpn_fail_over"`
	VlanId                           jsontypes.Int64   `json:"vlanId" tfsdk:"vlan_id"`
	DefaultVlanId                    jsontypes.Int64   `json:"defaultVlanId" tfsdk:"default_vlan_id"`
	ApTagsAndVlanIds                 types.List        `json:"apTagsAndVlanIds" tfsdk:"ap_tags_and_vlan_ids"`
	WalledGardenEnabled              jsontypes.Bool    `json:"walledGardenEnabled" tfsdk:"walled_garden_enabled"`
	WalledGardenRanges               types.List        `json:"walledGardenRanges" tfsdk:"walled_garden_ranges"`
	Gre                              types.Object      `json:"gre" tfsdk:"gre"`
	RadiusOverride                   jsontypes.Bool    `json:"radiusOverride" tfsdk:"radius_override"`
	RadiusGuestVlanEnabled           jsontypes.Bool    `json:"radiusGuestVlanEnabled" tfsdk:"radius_guest_vlan_enabled"`
	RadiusGuestVlanId                jsontypes.Int64   `json:"radiusGuestVlanId" tfsdk:"radius_guest_vlan_id"`
	MinBitrate                       jsontypes.Float64 `json:"minBitrate" tfsdk:"min_bit_rate"`
	BandSelection                    jsontypes.String  `json:"bandSelection" tfsdk:"band_selection"`
	PerClientBandwidthLimitUp        jsontypes.Int64   `json:"perClientBandwidthLimitUp" tfsdk:"per_client_bandwidth_limit_up"`
	PerClientBandwidthLimitDown      jsontypes.Int64   `json:"perClientBandwidthLimitDown" tfsdk:"per_client_bandwidth_limit_down"`
	PerSsidBandwidthLimitUp          jsontypes.Int64   `json:"perSsidBandwidthLimitUp" tfsdk:"per_ssid_bandwidth_limit_up"`
	PerSsidBandwidthLimitDown        jsontypes.Int64   `json:"perSsidBandwidthLimitDown" tfsdk:"per_ssid_bandwidth_limit_down"`
	LanIsolationEnabled              jsontypes.Bool    `json:"lanIsolationEnabled" tfsdk:"lan_isolation_enabled"`
	Visible                          jsontypes.Bool    `json:"visible" tfsdk:"visible"`
	AvailableOnAllAps                jsontypes.Bool    `json:"availableOnAllAps" tfsdk:"available_on_all_aps"`
	AvailabilityTags                 types.List        `json:"availabilityTags" tfsdk:"availability_tags"`
	MandatoryDhcpEnabled             jsontypes.Bool    `json:"mandatoryDhcpEnabled" tfsdk:"mandatory_dhcp_enabled"`
	AdultContentFilteringEnabled     jsontypes.Bool    `json:"adultContentFilteringEnabled" tfsdk:"adult_content_filtering_enabled"`
	DnsRewrite                       types.Object      `json:"dnsRewrite" tfsdk:"dns_rewrite"`
	SpeedBurst                       types.Object      `json:"speedBurst" tfsdk:"speed_burst"`
	NamedVlans                       types.Object      `json:"namedVlans" tfsdk:"named_vlans"`
}

func (m *NetworksWirelessSsidsResourceModel) UnmarshalJSON(b []byte) error {
	type Alias NetworksWirelessSsidsResourceModel

	temp := struct {
		*Alias
		RawAvailabilityTags json.RawMessage `json:"availabilityTags"`
	}{
		Alias: (*Alias)(m),
	}

	if err := json.Unmarshal(b, &temp); err != nil {
		return err
	}

	// Unmarshal the RawAvailabilityTags into a slice of strings
	var tags []string
	if err := json.Unmarshal(temp.RawAvailabilityTags, &tags); err != nil {
		return err
	}

	var items []attr.Value
	for _, tag := range tags {

		items = append(items, jsontypes.StringValue(tag))
	}

	list, _ := basetypes.NewListValue(jsontypes.StringType, items)

	// Convert the slice of strings to types.List
	m.AvailabilityTags = list

	return nil
}

type Dot11WConfig struct {
	Enabled  jsontypes.Bool `json:"enabled" tfsdk:"enabled"`
	Required jsontypes.Bool `json:"required" tfsdk:"required"`
}

func (cfg *Dot11WConfig) ToAPIModel(ctx context.Context) (*openApiClient.UpdateNetworkApplianceSsidRequestDot11w, diag.Diagnostics) {
	if cfg == nil {
		return nil, nil
	}

	dot11w := *openApiClient.NewUpdateNetworkApplianceSsidRequestDot11w()
	dot11w.SetEnabled(cfg.Enabled.ValueBool())
	dot11w.SetRequired(cfg.Required.ValueBool())

	return &dot11w, nil
}

type Dot11RConfig struct {
	Enabled  jsontypes.Bool `json:"enabled" tfsdk:"enabled"`
	Adaptive jsontypes.Bool `json:"adaptive" tfsdk:"adaptive"`
}

func (cfg *Dot11RConfig) ToAPIModel(ctx context.Context) (*openApiClient.UpdateNetworkWirelessSsidRequestDot11r, diag.Diagnostics) {
	if cfg == nil {
		return nil, nil
	}

	dot11r := *openApiClient.NewUpdateNetworkWirelessSsidRequestDot11r()
	dot11r.SetEnabled(cfg.Enabled.ValueBool())
	dot11r.SetAdaptive(cfg.Adaptive.ValueBool())

	return &dot11r, nil
}

type OauthConfig struct {
	AllowedDomains []jsontypes.String `json:"allowedDomains" tfsdk:"allowed_domains"`
}

func (cfg *OauthConfig) ToAPIModel(ctx context.Context) (*openApiClient.UpdateNetworkWirelessSsidRequestOauth, diag.Diagnostics) {
	if cfg == nil {
		return nil, nil
	}

	oauth := *openApiClient.NewUpdateNetworkWirelessSsidRequestOauth()
	var domains []string
	for _, domain := range cfg.AllowedDomains {
		domains = append(domains, domain.ValueString())
	}
	oauth.SetAllowedDomains(domains)

	return &oauth, nil
}

type LocalRadiusConfig struct {
	CacheTimeout              jsontypes.Int64                 `json:"cacheTimeout" tfsdk:"cache_timeout"`
	PasswordAuthentication    PasswordAuthenticationConfig    `json:"passwordAuthentication" tfsdk:"password_authentication"`
	CertificateAuthentication CertificateAuthenticationConfig `json:"certificateAuthentication" tfsdk:"certificate_authentication"`
}

func (cfg *LocalRadiusConfig) ToAPIModel(ctx context.Context) (*openApiClient.UpdateNetworkWirelessSsidRequestLocalRadius, diag.Diagnostics) {
	if cfg == nil {
		return nil, nil
	}

	localRadius := *openApiClient.NewUpdateNetworkWirelessSsidRequestLocalRadius()
	localRadius.SetCacheTimeout(int32(cfg.CacheTimeout.ValueInt64()))

	passwordAuth, diags := cfg.PasswordAuthentication.ToAPIModel(ctx)
	if diags.HasError() {
		return nil, diags
	}
	localRadius.SetPasswordAuthentication(*passwordAuth)

	certificateAuth, diags := cfg.CertificateAuthentication.ToAPIModel(ctx)
	if diags.HasError() {
		return nil, diags
	}
	localRadius.SetCertificateAuthentication(*certificateAuth)

	return &localRadius, nil
}

type PasswordAuthenticationConfig struct {
	Enabled jsontypes.Bool `json:"enabled" tfsdk:"enabled"`
}

func (cfg *PasswordAuthenticationConfig) ToAPIModel(ctx context.Context) (*openApiClient.UpdateNetworkWirelessSsidRequestLocalRadiusPasswordAuthentication, diag.Diagnostics) {
	if cfg == nil {
		return nil, nil
	}

	passwordAuthentication := *openApiClient.NewUpdateNetworkWirelessSsidRequestLocalRadiusPasswordAuthentication()
	passwordAuthentication.SetEnabled(cfg.Enabled.ValueBool())

	return &passwordAuthentication, nil
}

type CertificateAuthenticationConfig struct {
	Enabled                 jsontypes.Bool     `json:"enabled" tfsdk:"enabled"`
	UseLdap                 jsontypes.Bool     `json:"useLdap" tfsdk:"use_ldap"`
	UseOcsp                 jsontypes.Bool     `json:"useOcsp" tfsdk:"use_ocsp"`
	OcspResponderUrl        jsontypes.String   `json:"ocspResponderUrl" tfsdk:"ocsp_responder_url"`
	ClientRootCaCertificate CertificateContent `json:"clientRootCaCertificate" tfsdk:"client_root_ca_certificate"`
}

type CertificateContent struct {
	Contents jsontypes.String `json:"contents" tfsdk:"contents"`
}

func (cfg *CertificateAuthenticationConfig) ToAPIModel(ctx context.Context) (*openApiClient.UpdateNetworkWirelessSsidRequestLocalRadiusCertificateAuthentication, diag.Diagnostics) {
	if cfg == nil {
		return nil, nil
	}

	certificateAuthentication := *openApiClient.NewUpdateNetworkWirelessSsidRequestLocalRadiusCertificateAuthentication()

	// Set properties of certificateAuthentication based on cfg's fields
	certificateAuthentication.SetEnabled(cfg.Enabled.ValueBool())
	certificateAuthentication.SetUseLdap(cfg.UseLdap.ValueBool())
	certificateAuthentication.SetUseOcsp(cfg.UseOcsp.ValueBool())
	certificateAuthentication.SetOcspResponderUrl(cfg.OcspResponderUrl.ValueString())

	contents := *openApiClient.NewUpdateNetworkWirelessSsidRequestLocalRadiusCertificateAuthenticationClientRootCaCertificate()
	contents.SetContents(cfg.ClientRootCaCertificate.Contents.ValueString())
	certificateAuthentication.SetClientRootCaCertificate(contents)

	return &certificateAuthentication, nil
}

type LdapConfig struct {
	Servers               []ServerConfig     `json:"servers" tfsdk:"servers"`
	Credentials           Credential         `json:"credentials" tfsdk:"credentials"`
	BaseDistinguishedName jsontypes.String   `json:"baseDistinguishedName" tfsdk:"base_distinguished_name"`
	ServerCaCertificate   CertificateContent `json:"serverCaCertificate" tfsdk:"server_ca_certificate"`
}

func (cfg *LdapConfig) ToAPIModel(ctx context.Context) (*openApiClient.UpdateNetworkWirelessSsidRequestLdap, diag.Diagnostics) {
	if cfg == nil {
		return nil, nil
	}

	ldapPayload := *openApiClient.NewUpdateNetworkWirelessSsidRequestLdap()

	// Process Servers
	var serversPayload []openApiClient.UpdateNetworkWirelessSsidRequestLdapServersInner
	for _, serverConfig := range cfg.Servers {
		serverPayload, diags := serverConfig.ToAPIModel(ctx)
		if diags.HasError() {
			return nil, diags
		}
		serversPayload = append(serversPayload, *serverPayload)
	}
	ldapPayload.SetServers(serversPayload)

	// Process Credentials
	credentialsPayload, diags := cfg.Credentials.ToAPIModel(ctx)
	if diags.HasError() {
		return nil, diags
	}
	ldapPayload.SetCredentials(*credentialsPayload)

	// Process BaseDistinguishedName and ServerCaCertificate
	ldapPayload.SetBaseDistinguishedName(cfg.BaseDistinguishedName.ValueString())
	caCertificate := openApiClient.UpdateNetworkWirelessSsidRequestLdapServerCaCertificate{}
	caCertificate.SetContents(cfg.ServerCaCertificate.Contents.ValueString())

	ldapPayload.SetServerCaCertificate(caCertificate)

	return &ldapPayload, nil
}

type ServerConfig struct {
	Host jsontypes.String `json:"host" tfsdk:"host"`
	Port jsontypes.Int64  `json:"port" tfsdk:"port"`
}

func (cfg *ServerConfig) ToAPIModel(ctx context.Context) (*openApiClient.UpdateNetworkWirelessSsidRequestLdapServersInner, diag.Diagnostics) {
	if cfg == nil {
		return nil, nil
	}

	serverPayload := *openApiClient.NewUpdateNetworkWirelessSsidRequestLdapServersInner(cfg.Host.ValueString(), int32(cfg.Port.ValueInt64()))

	return &serverPayload, nil
}

type Credential struct {
	DistinguishedName jsontypes.String `json:"distinguishedName" tfsdk:"distinguished_name"`
	Password          jsontypes.String `json:"password" tfsdk:"password"`
}

func (cfg *Credential) ToAPIModel(ctx context.Context) (*openApiClient.UpdateNetworkWirelessSsidRequestLdapCredentials, diag.Diagnostics) {
	if cfg == nil {
		return nil, nil
	}

	credentialsPayload := *openApiClient.NewUpdateNetworkWirelessSsidRequestLdapCredentials()
	credentialsPayload.SetDistinguishedName(cfg.DistinguishedName.ValueString())
	credentialsPayload.SetPassword(cfg.Password.ValueString())

	return &credentialsPayload, nil
}

type ActiveDirectoryConfig struct {
	Servers     []AdServerConfig `json:"servers" tfsdk:"servers"`
	Credentials AdCredential     `json:"credentials" tfsdk:"credentials"`
}

func (cfg *ActiveDirectoryConfig) ToAPIModel(ctx context.Context) (*openApiClient.UpdateNetworkWirelessSsidRequestActiveDirectory, diag.Diagnostics) {
	if cfg == nil {
		return nil, nil
	}

	adPayload := *openApiClient.NewUpdateNetworkWirelessSsidRequestActiveDirectory()

	// Process Servers
	var serversPayload []openApiClient.UpdateNetworkWirelessSsidRequestActiveDirectoryServersInner
	for _, serverConfig := range cfg.Servers {
		serverPayload, diags := serverConfig.ToAPIModel(ctx)
		if diags.HasError() {
			return nil, diags
		}
		serversPayload = append(serversPayload, *serverPayload)
	}
	adPayload.SetServers(serversPayload)

	// Process Credentials
	credentialsPayload, diags := cfg.Credentials.ToAPIModel(ctx)
	if diags.HasError() {
		return nil, diags
	}
	adPayload.SetCredentials(*credentialsPayload)

	return &adPayload, nil
}

type AdServerConfig struct {
	Host jsontypes.String `json:"host" tfsdk:"host"`
	Port jsontypes.Int64  `json:"port" tfsdk:"port"`
}

func (cfg *AdServerConfig) ToAPIModel(ctx context.Context) (*openApiClient.UpdateNetworkWirelessSsidRequestActiveDirectoryServersInner, diag.Diagnostics) {
	if cfg == nil {
		return nil, nil
	}

	serverPayload := *openApiClient.NewUpdateNetworkWirelessSsidRequestActiveDirectoryServersInner(cfg.Host.ValueString())
	serverPayload.SetPort(int32(cfg.Port.ValueInt64()))

	return &serverPayload, nil
}

type AdCredential struct {
	LogonName jsontypes.String `json:"logonName,omitempty" tfsdk:"logon_name"`
	Password  jsontypes.String `json:"password" tfsdk:"password"`
}

func (cfg *AdCredential) ToAPIModel(ctx context.Context) (*openApiClient.UpdateNetworkWirelessSsidRequestActiveDirectoryCredentials, diag.Diagnostics) {
	if cfg == nil {
		return nil, nil
	}

	credentialsPayload := *openApiClient.NewUpdateNetworkWirelessSsidRequestActiveDirectoryCredentials()
	credentialsPayload.SetLogonName(cfg.LogonName.ValueString())
	credentialsPayload.SetPassword(cfg.Password.ValueString())

	return &credentialsPayload, nil
}

type RadiusServerConfig struct {
	Host                     jsontypes.String `json:"host" tfsdk:"host"`
	Secret                   jsontypes.String `json:"secret" tfsdk:"secret"`
	CaCertificate            jsontypes.String `json:"caCertificate" tfsdk:"ca_certificate"`
	Port                     jsontypes.Int64  `json:"port" tfsdk:"port"`
	RadSecEnabled            jsontypes.Bool   `json:"radsecEnabled" tfsdk:"rad_sec_enabled"`
	OpenRoamingCertificateId jsontypes.Int64  `json:"openRoamingCertificateId" tfsdk:"open_roaming_certificate_id"`
}

func (cfg *RadiusServerConfig) ToAPIModel(ctx context.Context) (*openApiClient.UpdateNetworkWirelessSsidRequestRadiusServersInner, diag.Diagnostics) {
	if cfg == nil {
		return nil, nil
	}

	radiusServer := *openApiClient.NewUpdateNetworkWirelessSsidRequestRadiusServersInner(cfg.Host.ValueString())
	radiusServer.SetSecret(cfg.Secret.ValueString())
	radiusServer.SetCaCertificate(cfg.CaCertificate.ValueString())
	radiusServer.SetPort(int32(cfg.Port.ValueInt64()))
	radiusServer.SetRadsecEnabled(cfg.RadSecEnabled.ValueBool())
	radiusServer.SetOpenRoamingCertificateId(int32(cfg.OpenRoamingCertificateId.ValueInt64()))

	return &radiusServer, nil
}

type TagAndVlanId struct {
	Tags   []jsontypes.String `json:"tags" tfsdk:"tags"`
	VlanId jsontypes.Int64    `json:"vlanId" tfsdk:"vlan_id"`
}

func (cfg *TagAndVlanId) ToAPIModel(ctx context.Context) (*openApiClient.UpdateNetworkWirelessSsidRequestApTagsAndVlanIdsInner, diag.Diagnostics) {
	if cfg == nil {
		return nil, nil
	}

	tagAndVlanIdPayload := *openApiClient.NewUpdateNetworkWirelessSsidRequestApTagsAndVlanIdsInner()
	var tagsPayload []string
	for _, tag := range cfg.Tags {
		tagsPayload = append(tagsPayload, tag.ValueString())
	}
	tagAndVlanIdPayload.SetTags(tagsPayload)
	tagAndVlanIdPayload.SetVlanId(int32(cfg.VlanId.ValueInt64()))

	return &tagAndVlanIdPayload, nil
}

type GreServerConfig struct {
	Host jsontypes.String `json:"host" tfsdk:"host"`
}

func (cfg *GreServerConfig) ToAPIModel(ctx context.Context) (*openApiClient.UpdateNetworkWirelessSsidRequestGreConcentrator, diag.Diagnostics) {
	if cfg == nil {
		return nil, nil
	}

	concentratorPayload := *openApiClient.NewUpdateNetworkWirelessSsidRequestGreConcentrator(cfg.Host.ValueString())

	return &concentratorPayload, nil
}

type GreConfig struct {
	Concentrator GreServerConfig `json:"concentrator" tfsdk:"concentrator"`
	Key          jsontypes.Int64 `json:"key" tfsdk:"key"`
}

func (cfg *GreConfig) ToAPIModel(ctx context.Context) (*openApiClient.UpdateNetworkWirelessSsidRequestGre, diag.Diagnostics) {
	if cfg == nil {
		return nil, nil
	}

	grePayload := *openApiClient.NewUpdateNetworkWirelessSsidRequestGre()
	grePayload.SetKey(int32(cfg.Key.ValueInt64()))

	concentratorPayload, diags := cfg.Concentrator.ToAPIModel(ctx)
	if diags.HasError() {
		return nil, diags
	}

	grePayload.SetConcentrator(*concentratorPayload)

	return &grePayload, nil
}

type DnsRewriteConfig struct {
	Enabled              jsontypes.Bool `json:"enabled" tfsdk:"enabled"`
	DnsCustomNameservers types.List     `json:"dnsCustomNameservers" tfsdk:"dns_custom_nameservers"`
}

func (cfg *DnsRewriteConfig) ToAPIModel(ctx context.Context) (*openApiClient.UpdateNetworkWirelessSsidRequestDnsRewrite, diag.Diagnostics) {
	if cfg == nil {
		return nil, nil
	}

	dnsRewritePayload := *openApiClient.NewUpdateNetworkWirelessSsidRequestDnsRewrite()
	dnsRewritePayload.SetEnabled(cfg.Enabled.ValueBool())

	var nameservers []string
	err := cfg.DnsCustomNameservers.ElementsAs(ctx, &nameservers, false)
	if err.HasError() {
		return nil, err
	}

	dnsRewritePayload.SetDnsCustomNameservers(nameservers)

	return &dnsRewritePayload, nil
}

type SpeedBurstConfig struct {
	Enabled jsontypes.Bool `json:"enabled" tfsdk:"enabled"`
}

func (cfg *SpeedBurstConfig) ToAPIModel(ctx context.Context) (*openApiClient.UpdateNetworkWirelessSsidRequestSpeedBurst, diag.Diagnostics) {
	if cfg == nil {
		return nil, nil
	}

	speedBurstPayload := *openApiClient.NewUpdateNetworkWirelessSsidRequestSpeedBurst()
	speedBurstPayload.SetEnabled(cfg.Enabled.ValueBool())

	return &speedBurstPayload, nil
}

type NamedVlansConfig struct {
	Tagging NamedVlansTaggingConfig `json:"tagging" tfsdk:"tagging"`
	Radius  NamedVlansRadiusConfig  `json:"radius" tfsdk:"radius"`
}

type NamedVlansTaggingConfig struct {
	Enabled         jsontypes.Bool   `json:"enabled" tfsdk:"enabled"`
	DefaultVlanName jsontypes.String `json:"defaultVlanName" tfsdk:"default_vlan_name"`
	ByApTags        types.List       `json:"byApTags" tfsdk:"by_ap_tags"`
}

type TagAndVlanName struct {
	Tags     types.List       `json:"tags" tfsdk:"tags"`
	VlanName jsontypes.String `json:"vlanName" tfsdk:"vlan_name"`
}

type TagValue struct {
	Value jsontypes.String `json:"value" tfsdk:"value"`
}

type NamedVlansRadiusConfig struct {
	GuestVlan NamedVlansGuestVlanConfig `json:"guestVlan" tfsdk:"guest_vlan"`
}

type NamedVlansGuestVlanConfig struct {
	Enabled jsontypes.Bool   `json:"enabled" tfsdk:"enabled"`
	Name    jsontypes.String `json:"name" tfsdk:"name"`
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
			"number": schema.Int64Attribute{
				MarkdownDescription: "Number",
				Required:            true,
				CustomType:          jsontypes.Int64Type,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the SSID",
				Optional:            true,
				CustomType:          jsontypes.StringType,
			},
			"enabled": schema.BoolAttribute{
				MarkdownDescription: "Whether or not the SSID is enabled",
				Optional:            true,
				CustomType:          jsontypes.BoolType,
			},
			"auth_mode": schema.StringAttribute{
				MarkdownDescription: "The association control method for the SSID",
				Optional:            true,
				CustomType:          jsontypes.StringType,
				Validators: []validator.String{
					stringvalidator.OneOf("8021x-google",
						"8021x-localradius", "8021x-meraki", "8021x-nac", "8021x-radius",
						"ipsk-with-nac", "ipsk-with-radius", "ipsk-without-radius", "open",
						"open-enhanced", "open-with-nac", "open-with-radius", "psk"),
				},
			},
			"enterprise_admin_access": schema.StringAttribute{
				MarkdownDescription: "Whether or not an SSID is accessible by 'enterprise' administrators",
				Optional:            true,
				CustomType:          jsontypes.StringType,
			},
			"encryption_mode": schema.StringAttribute{
				MarkdownDescription: "The psk encryption mode for the SSID",
				Optional:            true,
				CustomType:          jsontypes.StringType,
				Validators: []validator.String{
					stringvalidator.OneOf("wep", "wpa"),
				},
			},
			"psk": schema.StringAttribute{
				MarkdownDescription: "The passkey for the SSID",
				Optional:            true,
				CustomType:          jsontypes.StringType,
			},
			"wpa_encryption_mode": schema.StringAttribute{
				MarkdownDescription: "The types of WPA encryption",
				Optional:            true,
				CustomType:          jsontypes.StringType,
				Validators: []validator.String{
					stringvalidator.OneOf("WPA1 only", "WPA1 and WPA2", "WPA2 only", "WPA3 Transition Mode", "WPA3 only", "WPA3 192-bit Security"),
				},
			},
			"dot11w": schema.SingleNestedAttribute{
				MarkdownDescription: "The current setting for Protected Management Frames (802.11w)",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "Whether 802.11w is enabled or not",
						Optional:            true,
						CustomType:          jsontypes.BoolType,
					},
					"required": schema.BoolAttribute{
						MarkdownDescription: "Whether 802.11w is required or not",
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
						MarkdownDescription: "Whether 802.11r is enabled or not",
						Optional:            true,
						CustomType:          jsontypes.BoolType,
					},
					"adaptive": schema.BoolAttribute{
						MarkdownDescription: "Whether 802.11r is adaptive or not",
						Optional:            true,
						CustomType:          jsontypes.BoolType,
					},
				},
			},
			"splash_page": schema.StringAttribute{
				MarkdownDescription: "The type of splash page for the SSID",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
				Validators: []validator.String{
					stringvalidator.OneOf("Billing", "Cisco ISE", "Click-through splash page",
						"Facebook Wi-Fi", "Google Apps domain",
						"Google OAuth", "None", "Password-protected with Active Directory",
						"Password-protected with LDAP", "Password-protected with Meraki RADIUS",
						"Password-protected with custom RADIUS", "SMS authentication", "Sponsored guest",
						"Systems Manager Sentry"),
				},
			},
			"splash_guest_sponsor_domains": schema.ListAttribute{
				MarkdownDescription: "Array of valid sponsor email domains for sponsored guest splash type",
				Optional:            true,
				ElementType:         jsontypes.StringType,
			},
			"oauth": schema.SingleNestedAttribute{
				MarkdownDescription: "OAuth settings for the SSID",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"allowed_domains": schema.ListAttribute{
						MarkdownDescription: "List of allowed domains for OAuth",
						Optional:            true,
						ElementType:         jsontypes.StringType,
					},
				},
			},
			"local_radius": schema.SingleNestedAttribute{
				MarkdownDescription: "Local RADIUS server settings",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"cache_timeout": schema.Int64Attribute{
						MarkdownDescription: "The duration (in seconds) for which LDAP and OCSP lookups are cached",
						Optional:            true,
						CustomType:          jsontypes.Int64Type,
					},
					"password_authentication": schema.SingleNestedAttribute{
						MarkdownDescription: "Password-based authentication settings",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								MarkdownDescription: "Whether to use password-based authentication",
								Optional:            true,
								CustomType:          jsontypes.BoolType,
							},
						},
					},
					"certificate_authentication": schema.SingleNestedAttribute{
						MarkdownDescription: "Certificate verification settings",
						Optional:            true,
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								MarkdownDescription: "Whether to use certificate-based authentication",
								Optional:            true,
								CustomType:          jsontypes.BoolType,
							},
							"use_ldap": schema.BoolAttribute{
								MarkdownDescription: "Whether to verify the certificate with LDAP",
								Optional:            true,
								CustomType:          jsontypes.BoolType,
							},
							"use_ocsp": schema.BoolAttribute{
								MarkdownDescription: "Whether to verify the certificate with OCSP",
								Optional:            true,
								CustomType:          jsontypes.BoolType,
							},
							"ocsp_responder_url": schema.StringAttribute{
								MarkdownDescription: "The URL of the OCSP responder to verify client certificate status",
								Optional:            true,
								CustomType:          jsontypes.StringType,
							},
							"client_root_ca_certificate": schema.SingleNestedAttribute{
								MarkdownDescription: "The Client CA Certificate used to sign the client certificate",
								Optional:            true,
								Attributes: map[string]schema.Attribute{
									"contents": schema.StringAttribute{
										MarkdownDescription: "The contents of the Client CA Certificate",
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
				MarkdownDescription: "LDAP server settings",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"servers": schema.ListNestedAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The LDAP servers to be used for authentication",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"host": schema.StringAttribute{
									MarkdownDescription: "The LDAP server host",
									Optional:            true,
									CustomType:          jsontypes.StringType,
								},
								"port": schema.Int64Attribute{
									MarkdownDescription: "The LDAP server port",
									Optional:            true,
									CustomType:          jsontypes.Int64Type,
								},
							},
						},
					},
					"credentials": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "The credentials for LDAP server authentication",
						Attributes: map[string]schema.Attribute{
							"distinguished_name": schema.StringAttribute{
								MarkdownDescription: "The distinguished name for LDAP",
								Optional:            true,
								CustomType:          jsontypes.StringType,
							},
							"password": schema.StringAttribute{
								MarkdownDescription: "The password for LDAP",
								Optional:            true,
								CustomType:          jsontypes.StringType,
							},
						},
					},
					"base_distinguished_name": schema.StringAttribute{
						MarkdownDescription: "The base distinguished name on the LDAP server",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"server_ca_certificate": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "The CA certificate for the LDAP server",
						Attributes: map[string]schema.Attribute{
							"contents": schema.StringAttribute{
								MarkdownDescription: "The contents of the CA certificate",
								Optional:            true,
								CustomType:          jsontypes.StringType,
							},
						},
					},
				},
			},
			"active_directory": schema.SingleNestedAttribute{
				MarkdownDescription: "Active Directory server settings",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"servers": schema.ListNestedAttribute{
						Optional:            true,
						Computed:            true,
						MarkdownDescription: "The Active Directory servers to be used for authentication",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"host": schema.StringAttribute{
									MarkdownDescription: "The Active Directory server host",
									Optional:            true,
									CustomType:          jsontypes.StringType,
								},
								"port": schema.Int64Attribute{
									MarkdownDescription: "The Active Directory server port",
									Optional:            true,
									CustomType:          jsontypes.Int64Type,
								},
							},
						},
					},
					"credentials": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "The credentials for Active Directory server authentication",
						Attributes: map[string]schema.Attribute{
							"distinguished_name": schema.StringAttribute{
								MarkdownDescription: "The logon name for Active Directory",
								Optional:            true,
								CustomType:          jsontypes.StringType,
							},
							"password": schema.StringAttribute{
								MarkdownDescription: "The password for Active Directory",
								Optional:            true,
								CustomType:          jsontypes.StringType,
							},
						},
					},
				},
			},
			"radius_servers": schema.ListNestedAttribute{
				Optional:            true,
				MarkdownDescription: "The RADIUS servers to be used for authentication.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"host": schema.StringAttribute{
							MarkdownDescription: "IP address or hostname of the RADIUS server.",
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
						"port": schema.Int64Attribute{
							MarkdownDescription: "Port number on which the RADIUS server is listening.",
							Optional:            true,
							CustomType:          jsontypes.Int64Type,
						},
						"secret": schema.StringAttribute{
							MarkdownDescription: "Shared secret for the RADIUS server.",
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
						"ca_certificate": schema.StringAttribute{
							MarkdownDescription: "CA certificate for the RADIUS server.",
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
						"open_roaming_certificate_id": schema.Int64Attribute{
							MarkdownDescription: "OpenRoaming certificate ID associated with the RADIUS server.",
							Optional:            true,
							CustomType:          jsontypes.Int64Type,
						},
						"rad_sec_enabled": schema.BoolAttribute{
							MarkdownDescription: "Indicates whether RADSEC is enabled.",
							Optional:            true,
							CustomType:          jsontypes.BoolType,
						},
					},
				},
			},
			"radius_proxy_enabled": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether the RADIUS proxy is enabled.",
				Optional:            true,
				CustomType:          jsontypes.BoolType,
			},
			"radius_testing_enabled": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether RADIUS testing is enabled.",
				Optional:            true,
				CustomType:          jsontypes.BoolType,
			},
			"radius_called_station_id": schema.StringAttribute{
				MarkdownDescription: "The template of the called station identifier to be used for RADIUS",
				Optional:            true,
				CustomType:          jsontypes.StringType,
			},
			"radius_authentication_nas_id": schema.StringAttribute{
				MarkdownDescription: "The template of the NAS identifier to be used for RADIUS authentication",
				Optional:            true,
				CustomType:          jsontypes.StringType,
			},
			"radius_server_timeout": schema.Int64Attribute{
				MarkdownDescription: "The amount of time for which a RADIUS client waits for a reply from the RADIUS server",
				Optional:            true,
				CustomType:          jsontypes.Int64Type,
			},
			"radius_server_attempts_limit": schema.Int64Attribute{
				MarkdownDescription: "The maximum number of transmit attempts after which a RADIUS server is failed over",
				Optional:            true,
				CustomType:          jsontypes.Int64Type,
			},
			"radius_fallback_enabled": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether RADIUS fallback is enabled.",
				Optional:            true,
				CustomType:          jsontypes.BoolType,
			},
			"radius_coa_enabled": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether RADIUS Change of Authorization (CoA) is enabled.",
				Optional:            true,
				CustomType:          jsontypes.BoolType,
			},
			"radius_fail_over_policy": schema.StringAttribute{
				MarkdownDescription: "This policy determines how authentication requests should be handled in the event that all of the configured RADIUS servers are unreachable",
				Optional:            true,
				CustomType:          jsontypes.StringType,
			},
			"radius_load_balancing_policy": schema.StringAttribute{
				MarkdownDescription: "This policy determines which RADIUS server will be contacted first in an authentication attempt and the ordering of any necessary retry attempts",
				Optional:            true,
				CustomType:          jsontypes.StringType,
			},
			"radius_accounting_enabled": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether RADIUS accounting is enabled.",
				Optional:            true,
				CustomType:          jsontypes.BoolType,
			},
			"radius_accounting_servers": schema.ListNestedAttribute{
				Optional:            true,
				MarkdownDescription: "The RADIUS accounting servers to be used for accounting services.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"host": schema.StringAttribute{
							MarkdownDescription: "IP address or hostname of the RADIUS accounting server.",
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
						"port": schema.Int64Attribute{
							MarkdownDescription: "Port number on which the RADIUS accounting server is listening.",
							Optional:            true,
							CustomType:          jsontypes.Int64Type,
						},
						"secret": schema.StringAttribute{
							MarkdownDescription: "Shared secret for the RADIUS accounting server.",
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
						"ca_certificate": schema.StringAttribute{
							MarkdownDescription: "CA certificate for the RADIUS accounting server.",
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
						"rad_sec_enabled": schema.BoolAttribute{
							MarkdownDescription: "Indicates whether RADSEC (RADIUS over TLS) is enabled for secure communication with the RADIUS accounting server.",
							Optional:            true,
							CustomType:          jsontypes.BoolType,
						},
						"open_roaming_certificate_id": schema.Int64Attribute{
							MarkdownDescription: "The Open Roaming Certificate Id",
							Optional:            true,
							CustomType:          jsontypes.Int64Type,
						},
					},
				},
			},
			"radius_accounting_interim_interval": schema.Int64Attribute{
				MarkdownDescription: "The interval (in seconds) in which accounting information is updated and sent to the RADIUS accounting server",
				Optional:            true,
				CustomType:          jsontypes.Int64Type,
			},
			"radius_attribute_for_group_policies": schema.StringAttribute{
				MarkdownDescription: "Specify the RADIUS attribute used to look up group policies must be one of:\n 'Airespace-ACL-Name', 'Aruba-User-Role', 'Filter-Id' or 'Reply-Message'",
				Optional:            true,
				CustomType:          jsontypes.StringType,
				Validators: []validator.String{
					stringvalidator.OneOf("Airespace-ACL-Name", "Aruba-User-Role", "Filter-Id", "Reply-Message"),
				},
			},
			"ip_assignment_mode": schema.StringAttribute{
				MarkdownDescription: "The client IP assignment mode",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"use_vlan_tagging": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether VLAN tagging is used.",
				Optional:            true,
				CustomType:          jsontypes.BoolType,
			},
			"concentrator_network_id": schema.StringAttribute{
				MarkdownDescription: "The concentrator to use when the ipAssignmentMode is 'Layer 3 roaming with a concentrator' or 'VPN'",
				Optional:            true,
				CustomType:          jsontypes.StringType,
			},
			"secondary_concentrator_network_id": schema.StringAttribute{
				MarkdownDescription: "The secondary concentrator to use when the ipAssignmentMode is 'VPN'",
				Optional:            true,
				CustomType:          jsontypes.StringType,
			},
			"disassociate_clients_on_vpn_fail_over": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether clients should be disassociated during VPN failover.",
				Optional:            true,
				CustomType:          jsontypes.BoolType,
			},
			"vlan_id": schema.Int64Attribute{
				MarkdownDescription: "The VLAN ID used for VLAN tagging",
				Optional:            true,
				CustomType:          jsontypes.Int64Type,
			},
			"default_vlan_id": schema.Int64Attribute{
				MarkdownDescription: "The default VLAN ID used for 'all other APs'",
				Optional:            true,
				CustomType:          jsontypes.Int64Type,
			},
			"ap_tags_and_vlan_ids": schema.ListNestedAttribute{
				Optional:            true,
				MarkdownDescription: "A set of AP tags and corresponding VLAN IDs.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"tags": schema.ListAttribute{
							MarkdownDescription: "Array of AP tags.",
							Optional:            true,
							ElementType:         jsontypes.StringType,
						},
						"vlan_id": schema.Int64Attribute{
							MarkdownDescription: "VLAN ID associated with the AP tags.",
							Optional:            true,
							CustomType:          jsontypes.Int64Type,
						},
					},
				},
			},
			"walled_garden_enabled": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether a walled garden is enabled for the SSID.",
				Optional:            true,
				CustomType:          jsontypes.BoolType,
			},
			"walled_garden_ranges": schema.ListAttribute{
				MarkdownDescription: "List of Walled Garden ranges",
				Optional:            true,
				ElementType:         jsontypes.StringType,
			},
			"gre": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "GRE (Generic Routing Encapsulation) tunnel configuration.",
				Attributes: map[string]schema.Attribute{
					"concentrator": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "GRE tunnel concentrator configuration.",
						Attributes: map[string]schema.Attribute{
							"host": schema.StringAttribute{
								MarkdownDescription: "The GRE concentrator host.",
								Optional:            true,
								CustomType:          jsontypes.StringType,
							},
						},
					},
					"key": schema.Int64Attribute{
						MarkdownDescription: "The GRE key.",
						Optional:            true,
						CustomType:          jsontypes.Int64Type,
					},
				},
			},
			"radius_override": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether RADIUS attributes can override other settings.",
				Optional:            true,
				CustomType:          jsontypes.BoolType,
			},
			"radius_guest_vlan_enabled": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether the RADIUS guest VLAN is enabled.",
				Optional:            true,
				CustomType:          jsontypes.BoolType,
			},
			"radius_guest_vlan_id": schema.Int64Attribute{
				MarkdownDescription: "VLAN ID of the RADIUS Guest VLAN",
				Optional:            true,
				CustomType:          jsontypes.Int64Type,
			},
			"min_bit_rate": schema.Float64Attribute{
				MarkdownDescription: "The minimum bitrate in Mbps of this SSID in the default indoor RF profile",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.Float64Type,
			},
			"band_selection": schema.StringAttribute{
				MarkdownDescription: "The client-serving radio frequencies of this SSID in the default indoor RF profile",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"per_client_bandwidth_limit_up": schema.Int64Attribute{
				MarkdownDescription: "The upload bandwidth limit in Kbps",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.Int64Type,
			},
			"per_client_bandwidth_limit_down": schema.Int64Attribute{
				MarkdownDescription: "The download bandwidth limit in Kbps",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.Int64Type,
			},
			"per_ssid_bandwidth_limit_up": schema.Int64Attribute{
				MarkdownDescription: "The total upload bandwidth limit in Kbps",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.Int64Type,
			},
			"per_ssid_bandwidth_limit_down": schema.Int64Attribute{
				MarkdownDescription: "The total download bandwidth limit in Kbps",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.Int64Type,
			},
			"lan_isolation_enabled": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether LAN isolation is enabled.",
				Optional:            true,
				CustomType:          jsontypes.BoolType,
			},
			"visible": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether the SSID is visible.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.BoolType,
			},
			"available_on_all_aps": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether the SSID is available on all access points.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.BoolType,
			},
			"availability_tags": schema.ListAttribute{
				MarkdownDescription: "List of availability tags for the SSID.",
				Optional:            true,
				Computed:            true,
				ElementType:         jsontypes.StringType,
			},
			"mandatory_dhcp_enabled": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether mandatory DHCP is enabled.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.BoolType,
			},
			"adult_content_filtering_enabled": schema.BoolAttribute{
				MarkdownDescription: "Indicates whether adult content filtering is enabled.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.BoolType,
			},
			"dns_rewrite": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "DNS rewrite configuration.",
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "Indicates whether DNS rewrite is enabled.",
						Optional:            true,
						CustomType:          jsontypes.BoolType,
					},
					"dns_custom_nameservers": schema.ListAttribute{
						MarkdownDescription: "List of custom DNS nameservers.",
						Optional:            true,
						ElementType:         jsontypes.StringType,
					},
				},
			},
			"speed_burst": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Speed burst configuration.",
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "Indicates whether speed burst is enabled.",
						Optional:            true,
						CustomType:          jsontypes.BoolType,
					},
				},
			},
			"named_vlans": schema.SingleNestedAttribute{
				Optional:            true,
				MarkdownDescription: "Configuration for named VLANs.",
				Attributes: map[string]schema.Attribute{
					"tagging": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "Tagging configuration for named VLANs.",
						Attributes: map[string]schema.Attribute{
							"enabled": schema.BoolAttribute{
								MarkdownDescription: "Indicates whether VLAN tagging is enabled.",
								Optional:            true,
								CustomType:          jsontypes.BoolType,
							},
							"default_vlan_name": schema.StringAttribute{
								MarkdownDescription: "The default VLAN name.",
								Optional:            true,
								CustomType:          jsontypes.StringType,
							},
							"by_ap_tags": schema.ListNestedAttribute{
								Optional:            true,
								Computed:            true,
								MarkdownDescription: "Sets of AP tags and corresponding VLAN names.",
								NestedObject: schema.NestedAttributeObject{
									Attributes: map[string]schema.Attribute{
										"tags": schema.ListAttribute{
											MarkdownDescription: "Array of AP tags.",
											Optional:            true,
											ElementType:         jsontypes.StringType,
										},
										"vlan_name": schema.StringAttribute{
											MarkdownDescription: "VLAN name associated with the AP tags.",
											Optional:            true,
											CustomType:          jsontypes.StringType,
										},
									},
								},
							},
						},
					},
					"radius": schema.SingleNestedAttribute{
						Optional:            true,
						MarkdownDescription: "RADIUS configuration for named VLANs.",
						Attributes: map[string]schema.Attribute{
							"guest_vlan": schema.SingleNestedAttribute{
								Optional:            true,
								MarkdownDescription: "Guest VLAN configuration for RADIUS.",
								Attributes: map[string]schema.Attribute{
									"enabled": schema.BoolAttribute{
										MarkdownDescription: "Indicates whether the RADIUS guest VLAN is enabled.",
										Optional:            true,
										CustomType:          jsontypes.BoolType,
									},
									"name": schema.StringAttribute{
										MarkdownDescription: "Name of the RADIUS guest VLAN.",
										Optional:            true,
										CustomType:          jsontypes.StringType,
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

func (r *NetworksWirelessSsidsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *NetworksWirelessSsidsResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	payload, payloadDiag := NetworksWirelessSsidsPayload(ctx, data)
	if payloadDiag.HasError() {
		resp.Diagnostics.AddError("Resource Payload Error", fmt.Sprintf("\n%v", payloadDiag))
		return
	}

	number := strconv.Itoa(int(data.Number.ValueInt64()))
	_, httpResp, err := r.client.WirelessApi.UpdateNetworkWirelessSsid(context.Background(), data.NetworkID.ValueString(), number).UpdateNetworkWirelessSsidRequest(payload).Execute()

	// If there was an error during API call, add it to diagnostics.
	if err != nil {
		resp.Diagnostics.AddError(
			"Create HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
	}

	// Check for API success response code
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Create Unexpected HTTP Response Status Code",
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
	if err = json.NewDecoder(httpResp.Body).Decode(&data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

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
	number := strconv.Itoa(int(data.Number.ValueInt64()))
	_, httpResp, err := r.client.WirelessApi.GetNetworkWirelessSsid(context.Background(), data.NetworkID.ValueString(), number).Execute()

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

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

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

	payload, payloadDiag := NetworksWirelessSsidsPayload(ctx, data)
	if payloadDiag.HasError() {
		resp.Diagnostics.AddError("Resource Payload Error", fmt.Sprintf("\n%v", payloadDiag))
		return
	}
	number := strconv.Itoa(int(data.Number.ValueInt64()))
	_, httpResp, err := r.client.WirelessApi.UpdateNetworkWirelessSsid(context.Background(), data.NetworkID.ValueString(), number).UpdateNetworkWirelessSsidRequest(payload).Execute()

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

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

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

	payload := *openApiClient.NewUpdateNetworkWirelessSsidRequest()
	payload.SetEnabled(false)
	payload.SetName("")
	payload.SetAuthMode("open")
	payload.SetVlanId(1)

	number := strconv.Itoa(int(data.Number.ValueInt64()))
	_, httpResp, err := r.client.WirelessApi.UpdateNetworkWirelessSsid(context.Background(), data.NetworkID.ValueString(), number).UpdateNetworkWirelessSsidRequest(payload).Execute()

	// If there was an error during API call, add it to diagnostics.
	if err != nil {
		resp.Diagnostics.AddError(
			"Delete HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
	}

	// Check for API success response code
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Delete Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// If there were any errors up to this point, log the state data and return.
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("State Data", fmt.Sprintf("\n%v", data))
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
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("number"), idParts[1])...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func NetworksWirelessSsidsPayload(ctx context.Context, data *NetworksWirelessSsidsResourceModel) (openApiClient.UpdateNetworkWirelessSsidRequest, diag.Diagnostics) {
	var err diag.Diagnostics

	payload := *openApiClient.NewUpdateNetworkWirelessSsidRequest()

	// Name
	if !data.Name.IsNull() && !data.Name.IsUnknown() {
		payload.SetName(data.Name.ValueString())
	}

	// Enabled
	if !data.Enabled.IsNull() && !data.Enabled.IsUnknown() {
		payload.SetEnabled(data.Enabled.ValueBool())
	}

	// AuthMode
	if !data.AuthMode.IsNull() && !data.AuthMode.IsUnknown() {
		payload.SetAuthMode(data.AuthMode.ValueString())
	}

	// EnterpriseAdminAccess
	if !data.EnterpriseAdminAccess.IsNull() && !data.EnterpriseAdminAccess.IsUnknown() {
		payload.SetEnterpriseAdminAccess(data.EnterpriseAdminAccess.ValueString())
	}

	// EncryptionMode
	if !data.EncryptionMode.IsNull() && !data.EncryptionMode.IsUnknown() {
		payload.SetEncryptionMode(data.EncryptionMode.ValueString())
	}

	// Psk
	if !data.Psk.IsNull() && !data.Psk.IsUnknown() {
		payload.SetPsk(data.Psk.ValueString())
	}

	// WpaEncryptionMode
	if !data.WpaEncryptionMode.IsNull() && !data.WpaEncryptionMode.IsUnknown() {
		payload.SetWpaEncryptionMode(data.WpaEncryptionMode.ValueString())
	}

	// Dot11W
	if !data.Dot11W.IsNull() && !data.Dot11W.IsUnknown() {
		var dot11wData Dot11WConfig
		err = data.Dot11W.As(ctx, &dot11wData, basetypes.ObjectAsOptions{})
		if err != nil {
			return openApiClient.UpdateNetworkWirelessSsidRequest{}, diag.Diagnostics{
				diag.NewErrorDiagnostic("Failed to unmarshal dot11w settings",
					fmt.Sprintf("%v", err)),
			}
		}

		dot11w, diags := dot11wData.ToAPIModel(ctx)
		if diags.HasError() {
			return openApiClient.UpdateNetworkWirelessSsidRequest{}, diags
		}

		if dot11w != nil {
			payload.SetDot11w(*dot11w)
		}
	}

	// Dot11R
	if !data.Dot11R.IsNull() && !data.Dot11R.IsUnknown() {
		var dot11rData Dot11RConfig
		err = data.Dot11R.As(ctx, &dot11rData, basetypes.ObjectAsOptions{})
		if err != nil {
			return openApiClient.UpdateNetworkWirelessSsidRequest{}, diag.Diagnostics{
				diag.NewErrorDiagnostic("Failed to unmarshal dot11r",
					fmt.Sprintf("%v", err)),
			}
		}

		dot11r, diags := dot11rData.ToAPIModel(ctx)
		if diags.HasError() {
			return openApiClient.UpdateNetworkWirelessSsidRequest{}, diag.Diagnostics{
				diag.NewErrorDiagnostic("Failed to unmarshal dot11r",
					fmt.Sprintf("%v", err)),
			}
		}

		if dot11r != nil {
			payload.SetDot11r(*dot11r)
		}
	}

	// SplashPage
	if !data.SplashPage.IsNull() && !data.SplashPage.IsUnknown() {
		payload.SetSplashPage(data.SplashPage.ValueString())

		if data.SplashPage.ValueString() == "Sponsored guest" {
			// SplashGuestSponsorDomains
			if !data.SplashGuestSponsorDomains.IsNull() && !data.SplashGuestSponsorDomains.IsUnknown() {
				var domains []string
				for _, domain := range data.SplashGuestSponsorDomains.Elements() {
					domains = append(domains, domain.String())
					payload.SetSplashGuestSponsorDomains(domains)
				}
			}
		}
	}

	// Oauth
	if !data.Oauth.IsNull() && !data.Oauth.IsUnknown() {
		var oauthData OauthConfig
		oauthErr := data.Oauth.As(ctx, &oauthData, basetypes.ObjectAsOptions{})
		if oauthErr != nil {
			return openApiClient.UpdateNetworkWirelessSsidRequest{}, diag.Diagnostics{
				diag.NewErrorDiagnostic("Failed to unmarshal oauth",
					fmt.Sprintf("%v", oauthErr)),
			}
		}

		oauth, oauthDiags := oauthData.ToAPIModel(ctx)
		if oauthDiags.HasError() {
			return openApiClient.UpdateNetworkWirelessSsidRequest{}, oauthDiags
		}

		if oauth != nil {
			payload.SetOauth(*oauth)
		}

	}

	// LocalRadius
	if !data.LocalRadius.IsNull() && !data.LocalRadius.IsUnknown() {

		var localRadiusData LocalRadiusConfig
		localRadiusDataErr := data.LocalRadius.As(ctx, &localRadiusData, basetypes.ObjectAsOptions{})
		if localRadiusDataErr != nil {
			return openApiClient.UpdateNetworkWirelessSsidRequest{}, diag.Diagnostics{
				diag.NewErrorDiagnostic("Failed to unmarshal localRadius",
					fmt.Sprintf("%v", localRadiusDataErr)),
			}
		}

		localRadius, localRadiusDiags := localRadiusData.ToAPIModel(ctx)
		if localRadiusDiags.HasError() {
			return openApiClient.UpdateNetworkWirelessSsidRequest{}, localRadiusDiags
		}

		if localRadius != nil {
			payload.SetLocalRadius(*localRadius)
		}

		/*
			localRadius := *openApiClient.NewUpdateNetworkWirelessSsidRequestLocalRadius()

			var localRadiusData LocalRadiusConfig
			err = data.Ldap.As(ctx, localRadiusData, basetypes.ObjectAsOptions{})
			if err != nil {
				d.AddError("Failed to unmarshal localRadius settings", fmt.Sprintf("%v", err))
			}

			// CacheTimeout
			localRadius.SetCacheTimeout(int32(localRadiusData.CacheTimeout.ValueInt64()))

			// PasswordAuthentication
			passwordAuthentication := *openApiClient.NewUpdateNetworkWirelessSsidRequestLocalRadiusPasswordAuthentication()
			passwordAuthentication.SetEnabled(localRadiusData.PasswordAuthentication.Enabled.ValueBool())
			localRadius.SetPasswordAuthentication(passwordAuthentication)

			//  CertificateAuthentication
			certificateAuthentication := *openApiClient.NewUpdateNetworkWirelessSsidRequestLocalRadiusCertificateAuthentication()
			certificateAuthentication.SetEnabled(localRadiusData.CertificateAuthentication.Enabled.ValueBool())
			certificateAuthentication.SetUseLdap(localRadiusData.CertificateAuthentication.UseLdap.ValueBool())
			certificateAuthentication.SetUseOcsp(localRadiusData.CertificateAuthentication.UseOcsp.ValueBool())
			certificateAuthentication.SetOcspResponderUrl(localRadiusData.CertificateAuthentication.OcspResponderUrl.ValueString())

			contents := *openApiClient.NewUpdateNetworkWirelessSsidRequestLocalRadiusCertificateAuthenticationClientRootCaCertificate()
			contents.SetContents(localRadiusData.CertificateAuthentication.ClientRootCaCertificate.Contents.ValueString())
			certificateAuthentication.SetClientRootCaCertificate(contents)

			localRadius.SetCertificateAuthentication(certificateAuthentication)

			payload.SetLocalRadius(localRadius)
		*/
	}

	// Ldap
	if !data.Ldap.IsNull() && !data.Ldap.IsUnknown() {

		var ldapData LdapConfig
		ldapErr := data.Ldap.As(ctx, &ldapData, basetypes.ObjectAsOptions{})
		if ldapErr != nil {
			return openApiClient.UpdateNetworkWirelessSsidRequest{}, diag.Diagnostics{
				diag.NewErrorDiagnostic("Failed to unmarshal LDAP config",
					fmt.Sprintf("%v", ldapErr)),
			}
		}

		ldap, ldapDiags := ldapData.ToAPIModel(ctx)
		if ldapDiags.HasError() {
			return openApiClient.UpdateNetworkWirelessSsidRequest{}, ldapDiags
		}

		if ldap != nil {
			payload.SetLdap(*ldap)
		}

		/*
			ldapPayload := *openApiClient.NewUpdateNetworkWirelessSsidRequestLdap()

			var ldapData LdapConfig
			err = data.Ldap.As(ctx, ldapData, basetypes.ObjectAsOptions{})
			if err != nil {
				d.AddError("Failed to unmarshal ldap settings", fmt.Sprintf("%v", err))
			}

			// Servers
			var serversPayload []*openApiClient.UpdateNetworkWirelessSsidRequestLdapServersInner

			for _, server := range ldapData.Servers {
				var serverPayload *openApiClient.UpdateNetworkWirelessSsidRequestLdapServersInner

				serverPayload.SetHost(server.Host.ValueString())
				serverPayload.SetPort(int32(server.Port.ValueInt64()))

				serversPayload = append(serversPayload, serverPayload)
			}

			//  Credentials
			var credentialsPayload openApiClient.UpdateNetworkWirelessSsidRequestLdapCredentials
			credentialsPayload.SetPassword(ldapData.Credentials.Password.ValueString())
			credentialsPayload.SetDistinguishedName(ldapData.Credentials.DistinguishedName.ValueString())
			payload.Ldap.SetCredentials(credentialsPayload)

			// BaseDistinguishedName
			payload.Ldap.SetBaseDistinguishedName(ldapData.BaseDistinguishedName.ValueString())

			// ServerCaCertificate
			var serverCaCertificatePayload openApiClient.UpdateNetworkWirelessSsidRequestLdapServerCaCertificate
			serverCaCertificatePayload.SetContents(ldapData.ServerCaCertificate.Contents.ValueString())
			payload.Ldap.SetServerCaCertificate(serverCaCertificatePayload)

			payload.SetLdap(ldapPayload)

		*/
	}

	// ActiveDirectory
	if !data.ActiveDirectory.IsNull() && !data.ActiveDirectory.IsUnknown() {
		var adData ActiveDirectoryConfig
		adErr := data.ActiveDirectory.As(ctx, &adData, basetypes.ObjectAsOptions{})
		if adErr != nil {
			return openApiClient.UpdateNetworkWirelessSsidRequest{}, diag.Diagnostics{
				diag.NewErrorDiagnostic("Failed to unmarshal Active Directory config",
					fmt.Sprintf("%v", adErr)),
			}
		}

		ad, adDiags := adData.ToAPIModel(ctx)
		if adDiags.HasError() {
			return openApiClient.UpdateNetworkWirelessSsidRequest{}, adDiags
		}

		if ad != nil {
			payload.SetActiveDirectory(*ad)
		}

		/*
			adPayload := *openApiClient.NewUpdateNetworkWirelessSsidRequestActiveDirectory()

				adServers := make([]openApiClient.UpdateNetworkWirelessSsidRequestActiveDirectoryServersInner, 0)
				var adServersData ActiveDirectoryConfig
				err = data.ActiveDirectory.As(ctx, adServersData, basetypes.ObjectAsOptions{})
				if err != nil {
					d.AddError("Failed to unmarshal Active Directory settings", fmt.Sprintf("%v", err))
				}
				for _, server := range adServersData.Servers {
					adServer := *openApiClient.NewUpdateNetworkWirelessSsidRequestActiveDirectoryServersInner(server.Host.ValueString())
					adServer.SetPort(int32(server.Port.ValueInt64()))
					adServers = append(adServers, adServer)
				}
				adPayload.SetServers(adServers)

				adCredential := *openApiClient.NewUpdateNetworkWirelessSsidRequestActiveDirectoryCredentials()
				adCredential.SetLogonName(adServersData.Credentials.DistinguishedName.ValueString())
				adCredential.SetPassword(adServersData.Credentials.Password.ValueString())
				adPayload.SetCredentials(adCredential)

				payload.SetActiveDirectory(adPayload)
		*/
	}

	// RadiusServers
	if !data.RadiusServers.IsNull() && !data.RadiusServers.IsUnknown() {

		var radiusServersData []RadiusServerConfig
		radiusServersErr := data.RadiusServers.ElementsAs(ctx, &radiusServersData, false)
		if radiusServersErr != nil {
			return openApiClient.UpdateNetworkWirelessSsidRequest{}, diag.Diagnostics{
				diag.NewErrorDiagnostic("Failed to unmarshal RadiusServers",
					fmt.Sprintf("%v", radiusServersErr)),
			}
		}

		var radiusServersPayload []openApiClient.UpdateNetworkWirelessSsidRequestRadiusServersInner
		for _, serverConfig := range radiusServersData {
			radiusServer, radiusServerDiags := serverConfig.ToAPIModel(ctx)
			if radiusServerDiags.HasError() {
				return openApiClient.UpdateNetworkWirelessSsidRequest{}, radiusServerDiags
			}
			radiusServersPayload = append(radiusServersPayload, *radiusServer)
		}

		if len(radiusServersPayload) > 0 {
			payload.SetRadiusServers(radiusServersPayload)
		}
	}

	// RadiusProxyEnabled
	if !data.RadiusProxyEnabled.IsNull() && !data.RadiusProxyEnabled.IsUnknown() {
		payload.SetRadiusProxyEnabled(data.RadiusProxyEnabled.ValueBool())
	}

	// RadiusTestingEnabled
	if !data.RadiusTestingEnabled.IsNull() && !data.RadiusTestingEnabled.IsUnknown() {
		payload.SetRadiusTestingEnabled(data.RadiusTestingEnabled.ValueBool())
	}

	// RadiusCalledStationId
	if !data.RadiusCalledStationId.IsNull() && !data.RadiusCalledStationId.IsUnknown() {
		payload.SetRadiusCalledStationId(data.RadiusCalledStationId.ValueString())
	}

	// RadiusAuthenticationNasId
	if !data.RadiusAuthenticationNasId.IsNull() && !data.RadiusAuthenticationNasId.IsUnknown() {
		payload.SetRadiusAuthenticationNasId(data.RadiusAuthenticationNasId.ValueString())
	}

	// RadiusServerTimeout
	if !data.RadiusServerTimeout.IsNull() && !data.RadiusServerTimeout.IsUnknown() {
		payload.SetRadiusServerTimeout(int32(data.RadiusServerTimeout.ValueInt64()))
	}

	// RadiusServerAttemptsLimit
	if !data.RadiusServerAttemptsLimit.IsNull() && !data.RadiusServerAttemptsLimit.IsUnknown() {
		payload.SetRadiusServerAttemptsLimit(int32(data.RadiusServerAttemptsLimit.ValueInt64()))
	}

	// RadiusFallbackEnabled
	if !data.RadiusFallbackEnabled.IsNull() && !data.RadiusFallbackEnabled.IsUnknown() {
		payload.SetRadiusFallbackEnabled(data.RadiusFallbackEnabled.ValueBool())
	}

	// RadiusCoaEnabled
	if !data.RadiusCoaEnabled.IsNull() && !data.RadiusCoaEnabled.IsUnknown() {
		payload.SetRadiusCoaEnabled(data.RadiusCoaEnabled.ValueBool())
	}

	// RadiusFailOverPolicy
	if !data.RadiusFailOverPolicy.IsNull() && !data.RadiusFailOverPolicy.IsUnknown() {
		payload.SetRadiusFailoverPolicy(data.RadiusFailOverPolicy.ValueString())
	}

	// RadiusLoadBalancingPolicy
	if !data.RadiusLoadBalancingPolicy.IsNull() && !data.RadiusLoadBalancingPolicy.IsUnknown() {
		payload.SetRadiusLoadBalancingPolicy(data.RadiusLoadBalancingPolicy.ValueString())
	}

	// RadiusAccountingEnabled
	if !data.RadiusAccountingEnabled.IsNull() && !data.RadiusAccountingEnabled.IsUnknown() {
		payload.SetRadiusAccountingEnabled(data.RadiusAccountingEnabled.ValueBool())
	}

	// RadiusAccountingServers
	if !data.RadiusAccountingServers.IsNull() && !data.RadiusAccountingServers.IsUnknown() {
		var radiusAccServersData []RadiusServerConfig
		err = data.RadiusAccountingServers.ElementsAs(ctx, radiusAccServersData, false)
		if err != nil {

		}

		radiusAccServers := make([]openApiClient.UpdateNetworkWirelessSsidRequestRadiusAccountingServersInner, len(radiusAccServersData))

		for i, server := range radiusAccServersData {
			radiusAccServer := *openApiClient.NewUpdateNetworkWirelessSsidRequestRadiusAccountingServersInner(server.Host.ValueString())
			radiusAccServer.SetSecret(server.Secret.ValueString())
			radiusAccServer.SetCaCertificate(server.CaCertificate.ValueString())
			radiusAccServer.SetPort(int32(server.Port.ValueInt64()))
			radiusAccServer.SetRadsecEnabled(server.RadSecEnabled.ValueBool())
			radiusAccServers[i] = radiusAccServer
		}
		payload.SetRadiusAccountingServers(radiusAccServers)
	}

	// RadiusAccountingInterimInterval
	if !data.RadiusAccountingInterimInterval.IsNull() && !data.RadiusAccountingInterimInterval.IsUnknown() {
		payload.SetRadiusAccountingInterimInterval(int32(data.RadiusAccountingInterimInterval.ValueInt64()))
	}

	// RadiusAttributeForGroupPolicies
	if !data.RadiusAttributeForGroupPolicies.IsNull() && !data.RadiusAttributeForGroupPolicies.IsUnknown() {
		payload.SetRadiusAttributeForGroupPolicies(data.RadiusAttributeForGroupPolicies.ValueString())
	}

	// IpAssignmentMode
	if !data.IpAssignmentMode.IsNull() && !data.IpAssignmentMode.IsUnknown() {
		payload.SetIpAssignmentMode(data.IpAssignmentMode.ValueString())
	}

	// UseVlanTagging
	if !data.UseVlanTagging.IsNull() && !data.UseVlanTagging.IsUnknown() {
		payload.SetUseVlanTagging(data.UseVlanTagging.ValueBool())
	}

	// ConcentratorNetworkId
	if !data.ConcentratorNetworkId.IsNull() && !data.ConcentratorNetworkId.IsUnknown() {
		payload.SetConcentratorNetworkId(data.ConcentratorNetworkId.ValueString())
	}

	// SecondaryConcentratorNetworkId
	if !data.SecondaryConcentratorNetworkId.IsNull() && !data.SecondaryConcentratorNetworkId.IsUnknown() {
		payload.SetSecondaryConcentratorNetworkId(data.SecondaryConcentratorNetworkId.ValueString())
	}

	// DisassociateClientsOnVpnFailOver
	if !data.DisassociateClientsOnVpnFailOver.IsNull() && !data.DisassociateClientsOnVpnFailOver.IsUnknown() {
		payload.SetDisassociateClientsOnVpnFailover(data.DisassociateClientsOnVpnFailOver.ValueBool())
	}

	// VlanId
	if !data.VlanId.IsNull() && !data.VlanId.IsUnknown() {
		payload.SetVlanId(int32(data.VlanId.ValueInt64()))
	}

	// DefaultVlanId
	if !data.DefaultVlanId.IsNull() && !data.DefaultVlanId.IsUnknown() {
		payload.SetDefaultVlanId(int32(data.DefaultVlanId.ValueInt64()))
	}

	// ApTagsAndVlanIds
	if !data.ApTagsAndVlanIds.IsNull() && !data.ApTagsAndVlanIds.IsUnknown() {

		var tagAndVlanIds []TagAndVlanId
		tagAndVlanIdsDiags := data.ApTagsAndVlanIds.ElementsAs(ctx, &tagAndVlanIds, false)
		if tagAndVlanIdsDiags.HasError() {
			return openApiClient.UpdateNetworkWirelessSsidRequest{}, tagAndVlanIdsDiags
		}

		var apTagsAndVlanIdsPayload []openApiClient.UpdateNetworkWirelessSsidRequestApTagsAndVlanIdsInner
		for _, tagAndVlanId := range tagAndVlanIds {
			tagAndVlanIdModel, tagAndVlanIdDiags := tagAndVlanId.ToAPIModel(ctx)
			if tagAndVlanIdDiags.HasError() {
				return openApiClient.UpdateNetworkWirelessSsidRequest{}, tagAndVlanIdDiags
			}
			apTagsAndVlanIdsPayload = append(apTagsAndVlanIdsPayload, *tagAndVlanIdModel)
		}

		payload.SetApTagsAndVlanIds(apTagsAndVlanIdsPayload)

	}

	// WalledGardenEnabled
	if !data.WalledGardenEnabled.IsNull() && !data.WalledGardenEnabled.IsUnknown() {
		payload.SetWalledGardenEnabled(data.WalledGardenEnabled.ValueBool())
	}

	// WalledGardenRanges
	if !data.WalledGardenRanges.IsNull() && !data.WalledGardenRanges.IsUnknown() {
		var ranges []string
		for _, walledGarden := range data.WalledGardenRanges.Elements() {
			ranges = append(ranges, walledGarden.String())
		}
		payload.SetWalledGardenRanges(ranges)
	}

	// Gre
	var greData GreConfig
	err = data.Gre.As(ctx, &greData, basetypes.ObjectAsOptions{})
	if err != nil {
		// handle error
	}

	gre, diags := greData.ToAPIModel(ctx)
	if diags.HasError() {
		// Add diagnostics to your response
	}

	if gre != nil {
		payload.SetGre(*gre)
	}

	// RadiusOverride
	if !data.RadiusOverride.IsNull() && !data.RadiusOverride.IsUnknown() {
		payload.SetRadiusOverride(data.RadiusOverride.ValueBool())
	}

	// RadiusGuestVlanEnabled
	if !data.RadiusGuestVlanEnabled.IsNull() && !data.RadiusGuestVlanEnabled.IsUnknown() {
		payload.SetRadiusGuestVlanEnabled(data.RadiusGuestVlanEnabled.ValueBool())
	}

	// RadiusGuestVlanId
	if !data.RadiusGuestVlanId.IsNull() && !data.RadiusGuestVlanId.IsUnknown() {
		payload.SetRadiusGuestVlanId(int32(data.RadiusGuestVlanId.ValueInt64()))
	}

	// MinBitrate
	if !data.MinBitrate.IsNull() && !data.MinBitrate.IsUnknown() {
		payload.SetMinBitrate(float32(data.MinBitrate.ValueFloat64()))
	}

	// BandSelection
	if !data.BandSelection.IsNull() && !data.BandSelection.IsUnknown() {
		payload.SetBandSelection(data.BandSelection.ValueString())
	}

	// PerClientBandwidthLimitUp
	if !data.PerClientBandwidthLimitUp.IsNull() && !data.PerClientBandwidthLimitUp.IsUnknown() {
		payload.SetPerClientBandwidthLimitUp(int32(data.PerClientBandwidthLimitUp.ValueInt64()))
	}

	// PerClientBandwidthLimitDown
	if !data.PerClientBandwidthLimitDown.IsNull() && !data.PerClientBandwidthLimitDown.IsUnknown() {
		payload.SetPerClientBandwidthLimitDown(int32(data.PerClientBandwidthLimitDown.ValueInt64()))
	}

	// PerSsidBandwidthLimitUp
	if !data.PerSsidBandwidthLimitUp.IsNull() && !data.PerSsidBandwidthLimitUp.IsUnknown() {
		payload.SetPerSsidBandwidthLimitUp(int32(data.PerSsidBandwidthLimitUp.ValueInt64()))
	}

	// PerSsidBandwidthLimitDown
	if !data.PerSsidBandwidthLimitDown.IsNull() && !data.PerSsidBandwidthLimitDown.IsUnknown() {
		payload.SetPerSsidBandwidthLimitDown(int32(data.PerSsidBandwidthLimitDown.ValueInt64()))
	}

	// LanIsolationEnabled
	if !data.LanIsolationEnabled.IsNull() && !data.LanIsolationEnabled.IsUnknown() {
		payload.SetLanIsolationEnabled(data.LanIsolationEnabled.ValueBool())
	}

	// Visible
	if !data.Visible.IsNull() && !data.Visible.IsUnknown() {
		payload.SetVisible(data.Visible.ValueBool())
	}

	// AvailableOnAllAps
	if !data.AvailableOnAllAps.IsNull() && !data.AvailableOnAllAps.IsUnknown() {
		payload.SetAvailableOnAllAps(data.AvailableOnAllAps.ValueBool())
	}

	// AvailabilityTags
	if !data.AvailabilityTags.IsNull() && !data.AvailabilityTags.IsUnknown() {
		var tags []string

		for _, tag := range data.AvailabilityTags.Elements() {
			tags = append(tags, tag.String())
		}

		payload.SetAvailabilityTags(tags)
	}

	// MandatoryDhcpEnabled
	if !data.MandatoryDhcpEnabled.IsNull() && !data.MandatoryDhcpEnabled.IsUnknown() {
		payload.SetMandatoryDhcpEnabled(data.MandatoryDhcpEnabled.ValueBool())
	}

	// AdultContentFilteringEnabled
	if !data.AdultContentFilteringEnabled.IsNull() && !data.AdultContentFilteringEnabled.IsUnknown() {
		payload.SetAdultContentFilteringEnabled(data.AdultContentFilteringEnabled.ValueBool())
	}

	// DnsRewrite
	if !data.DnsRewrite.IsNull() && !data.DnsRewrite.IsUnknown() {
		var dnsRewriteData DnsRewriteConfig
		dnsRewriteDataErr := data.DnsRewrite.As(ctx, &dnsRewriteData, basetypes.ObjectAsOptions{})
		if dnsRewriteDataErr != nil {
			return openApiClient.UpdateNetworkWirelessSsidRequest{}, diag.Diagnostics{
				diag.NewErrorDiagnostic("Failed to unmarshal DNS rewrite config",
					fmt.Sprintf("%v", dnsRewriteDataErr)),
			}
		}

		dnsRewrite, dnsRewriteDataDiags := dnsRewriteData.ToAPIModel(ctx)
		if dnsRewriteDataDiags.HasError() {
			return openApiClient.UpdateNetworkWirelessSsidRequest{}, dnsRewriteDataDiags
		}

		if dnsRewrite != nil {
			payload.SetDnsRewrite(*dnsRewrite)
		}
	}

	// SpeedBurst
	if !data.SpeedBurst.IsNull() && !data.SpeedBurst.IsUnknown() {
		var speedBurstData SpeedBurstConfig
		speedBurstDataErr := data.SpeedBurst.As(ctx, &speedBurstData, basetypes.ObjectAsOptions{})
		if speedBurstDataErr != nil {
			return openApiClient.UpdateNetworkWirelessSsidRequest{}, diag.Diagnostics{
				diag.NewErrorDiagnostic("Failed to unmarshal Speed Burst config",
					fmt.Sprintf("%v", speedBurstDataErr)),
			}
		}

		speedBurst, speedBurstDataDiags := speedBurstData.ToAPIModel(ctx)
		if speedBurstDataDiags.HasError() {
			return openApiClient.UpdateNetworkWirelessSsidRequest{}, speedBurstDataDiags
		}

		if speedBurst != nil {
			payload.SetSpeedBurst(*speedBurst)
		}
	}

	return payload, nil
}

/*
	// NamedVlans
		namedVlansPayload := openApiClient.NewNamedVlansConfig()

		// Handling Tagging Configuration
		taggingPayload := openApiClient.NewNamedVlansTaggingConfig()
		taggingPayload.SetEnabled(data.NamedVlans.Tagging.Enabled.ValueBool())
		taggingPayload.SetDefaultVlanName(data.NamedVlans.Tagging.DefaultVlanName.ValueString())

		byApTags := make([]openApiClient.TagAndVlanName, len(data.NamedVlans.Tagging.ByApTags))
		for i, tagData := range data.NamedVlans.Tagging.ByApTags {
			tagAndVlanNamePayload := openApiClient.NewTagAndVlanName()
			tagsPayload := make([]string, len(tagData.Tags))
			for j, tag := range tagData.Tags {
				tagsPayload[j] = tag.Value.ValueString()
			}
			tagAndVlanNamePayload.SetTags(tagsPayload)
			tagAndVlanNamePayload.SetVlanName(tagData.VlanName.ValueString())
			byApTags[i] = *tagAndVlanNamePayload
		}
		taggingPayload.SetByApTags(byApTags)
		namedVlansPayload.SetTagging(*taggingPayload)

		// Handling Radius Configuration
		radiusPayload := openApiClient.NewNamedVlansRadiusConfig()
		guestVlanPayload := openApiClient.NewNamedVlansGuestVlanConfig()
		guestVlanPayload.SetEnabled(data.NamedVlans.Radius.GuestVlan.Enabled.ValueBool())
		guestVlanPayload.SetName(data.NamedVlans.Radius.GuestVlan.Name.ValueString())
		radiusPayload.SetGuestVlan(*guestVlanPayload)
		namedVlansPayload.SetRadius(*radiusPayload)

		payload.NamedVlans(*namedVlansPayload)
*/
