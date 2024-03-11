package provider

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// Ensure CiscoMerakiProvider satisfies various provider interfaces.
var _ provider.Provider = &CiscoMerakiProvider{}

// CiscoMerakiProvider defines the provider implementation.
type CiscoMerakiProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// CiscoMerakiProviderModel describes the provider data model.
type CiscoMerakiProviderModel struct {
	LoggingEnabled        types.Bool   `tfsdk:"logging_enabled"`
	ApiKey                types.String `tfsdk:"api_key"`
	BaseUrl               types.String `tfsdk:"base_url"`
	BasePath              types.String `tfsdk:"base_path"`
	CertificatePath       types.String `tfsdk:"certificate_path"`
	Proxy                 types.String `tfsdk:"proxy"`
	SingleRequestTimeout  types.Int64  `tfsdk:"single_request_timeout"`
	MaximumRetries        types.Int64  `tfsdk:"maximum_retries"`
	Nginx429RetryWaitTime types.Int64  `tfsdk:"nginx_429_retry_wait_time"`
	WaitOnRateLimit       types.Bool   `tfsdk:"wait_on_rate_limit"`
}

func (p *CiscoMerakiProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "meraki"
	resp.Version = p.version
}

func (p *CiscoMerakiProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Terraform Provider Meraki is a declarative infrastructure management tool for the Cisco Meraki Dashboard API.",
		Attributes: map[string]schema.Attribute{
			"logging_enabled": schema.BoolAttribute{
				Description: "Display http client debug messages in console",
				Optional:    true,
			},
			"api_key": schema.StringAttribute{
				Description: "Meraki Dashboard API Key",
				Optional:    true,
				Sensitive:   true,
			},
			"base_url": schema.StringAttribute{
				Description: "Endpoint for Meraki Dashboard API",
				MarkdownDescription: "The API version must be specified in the URL:" +
					"Example: `https://api.meraki.com" +
					"For organizations hosted in the China dashboard, use: `https://api.meraki.cn/v1`",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^(https:\/\/)(?:[a-zA-Z0-9]{1,62}(?:[-\.][a-zA-Z0-9]{1,62})+)(:\d+)?$`),
						"The API version must be specified in the URL. Example: https://api.meraki.com",
					)},
			},
			"base_path": schema.StringAttribute{
				Description: "API version prefix to be appended after the base URL",
				MarkdownDescription: "The API version to be specified in the URL:" +
					"Example: `/api/v1",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^\/api\/v1$`),
						"The API version to be specified in the URL. Example: /api/v1",
					)},
			},
			"certificate_path": schema.StringAttribute{
				Description: "Path for TLS/SSL certificate verification if behind local proxy",
				Optional:    true,
				Sensitive:   true,
			},
			"proxy": schema.StringAttribute{
				Description: "Proxy server and port, if needed, for HTTPS",
				Optional:    true,
			},
			"single_request_timeout": schema.Int64Attribute{
				Description: "Maximum number of seconds for each API call",
				Optional:    true,
			},
			"maximum_retries": schema.Int64Attribute{
				Description: "Retry up to this many times when encountering 429s or other server-side errors",
				Optional:    true,
			},
			"nginx_429_retry_wait_time": schema.Int64Attribute{
				Description: "Nginx 429 retry wait time",
				Optional:    true,
			},
			"wait_on_rate_limit": schema.BoolAttribute{
				Description: "Retry if 429 rate limit error encountered",
				Optional:    true,
			},
		},
	}
}

// Custom transport to add bearer token in the Authorization header
type bearerAuthTransport struct {
	Transport *http.Transport
	Token     string
}

func (t *bearerAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	// Add the bearer token to the Authorization header
	req.Header.Set("Authorization", "Bearer "+t.Token)

	// Use the underlying transport to perform the actual request
	return t.Transport.RoundTrip(req)
}

func (p *CiscoMerakiProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data CiscoMerakiProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Get http retryClient variables and default values
	configuration := openApiClient.NewConfiguration()

	// Debug
	if p.version == "dev" {

		// always enable debug for provider development
		configuration.Debug = true

	} else if data.LoggingEnabled.ValueBool() {

		// check if user enabled debug in the provider
		configuration.Debug = data.LoggingEnabled.ValueBool()
	}

	// MERAKI BASE URL
	if !data.BaseUrl.IsNull() {
		baseUrl, err := url.Parse(data.BaseUrl.ValueString())
		if err == nil {
			configuration.Servers = openApiClient.ServerConfigurations{
				{
					URL:         baseUrl.String() + "/{basePath}",
					Description: "No description provided",
					Variables: map[string]openApiClient.ServerVariable{
						"basePath": {
							Description:  "Meraki API Go Client",
							DefaultValue: data.BasePath.ValueString(),
						},
					},
				},
			}
		}
	}

	// UserAgent
	configuration.UserAgent = configuration.UserAgent + " terraform/" + p.version

	// Set certificate path
	if !data.CertificatePath.IsNull() {
		configuration.CertificatePath = data.CertificatePath.ValueString()
	}

	// Proxy
	if !data.Proxy.IsNull() {
		configuration.RequestsProxy = data.Proxy.ValueString()
	}

	// SingleRequestTimeout
	if !data.SingleRequestTimeout.IsNull() {
		configuration.SingleRequestTimeout = int(data.SingleRequestTimeout.ValueInt64())
	}

	// MaximumRetries
	if !data.MaximumRetries.IsNull() {
		configuration.MaximumRetries = int(data.MaximumRetries.ValueInt64())
	}

	// Nginx429RetryWaitTime
	if !data.Nginx429RetryWaitTime.IsNull() {
		configuration.Nginx429RetryWaitTime = int(data.Nginx429RetryWaitTime.ValueInt64())
	}

	// New custom retryable retryClient
	retryClient := retryablehttp.NewClient()

	// Custom Transport
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
	}

	// add certificate to retryClient if certificate path isn't empty
	if configuration.CertificatePath != "" {

		// Load the certificate file
		certFile := configuration.CertificatePath
		cert, err := os.ReadFile(certFile)
		if err != nil {
			e := fmt.Sprintf("%v", err.Error())
			tflog.Error(ctx, e)
		}

		// Create a certificate pool and add the certificate
		certPool := x509.NewCertPool()
		certPool.AppendCertsFromPEM(cert)

		// Create a custom Cert pool with the certificate and add TLS configuration to transport
		transport.TLSClientConfig.RootCAs = certPool

	}

	if configuration.RequestsProxy != "" {
		proxyUrl, err := url.Parse(configuration.RequestsProxy)
		if err == nil {
			transport.Proxy = http.ProxyURL(proxyUrl)
		}
	}

	// Set single request timeout in transport
	retryClient.HTTPClient.Timeout = time.Duration(configuration.SingleRequestTimeout) * time.Second

	retryClient.RetryMax = configuration.MaximumRetries

	retryClient.RetryWaitMax = time.Duration(configuration.Nginx429RetryWaitTime) * time.Second

	configuration.UserAgent = configuration.UserAgent + "terraform" + p.version

	// Set Bearer Token in transport
	authenticatedTransport := &bearerAuthTransport{
		Transport: transport,
	}

	// MERAKI DASHBOARD API KEY
	if !data.ApiKey.IsNull() {
		authenticatedTransport.Token = data.ApiKey.ValueString()
	} else {
		authenticatedTransport.Token = os.Getenv("MERAKI_DASHBOARD_API_KEY")
	}

	retryClient.HTTPClient.Transport = authenticatedTransport
	configuration.HTTPClient = retryClient.HTTPClient

	client := openApiClient.NewAPIClient(configuration)

	resp.DataSourceData = client
	resp.ResourceData = client

}

func (p *CiscoMerakiProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewOrganizationResource,
		NewOrganizationSamlResource,
		NewOrganizationsSamlIdpResource,
		NewOrganizationsAdminResource,
		NewOrganizationsAdaptivePolicyAclResource,
		NewNetworkResource,
		NewOrganizationsSamlRolesResource,
		NewNetworksSwitchSettingsResource,
		NewOrganizationsSnmpResource,
		NewNetworksSettingsResource,
		NewNetworksApplianceFirewallL3FirewallRulesResource,
		NewNetworksApplianceFirewallL7FirewallRulesResource,
		NewOrganizationsApplianceVpnVpnFirewallRulesResource,
		NewNetworksTrafficAnalysisResource,
		NewNetworksNetflowResource,
		NewNetworksSyslogServersResource,
		NewNetworksApplianceVlansSettingsResource,
		NewNetworksApplianceSettingsResource,
		NewNetworksApplianceFirewallSettingsResource,
		NewNetworksSwitchQosRuleResource,
		NewNetworksSwitchDscpToCosMappingsResource,
		NewNetworksSwitchMtuResource,
		NewNetworksGroupPolicyResource,
		NewOrganizationsLicenseResource,
		NewNetworksWirelessSsidsFirewallL3FirewallRulesResource,
		NewNetworksWirelessSsidsFirewallL7FirewallRulesResource,
		NewDevicesResource,
		NewOrganizationsClaimResource,
		NewNetworksDevicesClaimResource,
		NewNetworkApplianceStaticRoutesResource,
		NewNetworksCellularGatewaySubnetPoolResource,
		NewNetworksCellularGatewayUplinkResource,
		NewNetworksWirelessSsidsSplashSettingsResource,
		NewDevicesCellularSimsResource,
		NewDevicesTestAccDevicesManagementInterfaceResourceResource,
		NewNetworksApplianceVpnSiteToSiteVpnResource,
		NewDevicesSwitchPortsCycleResource,
		NewNetworksApplianceTrafficShapingUplinkBandWidthResource,
		NewNetworksApplianceVLANsResource,
		NewDevicesSwitchPortResource,
		NewNetworksAppliancePortsResource,
		NewNetworksWirelessSsidsResource,
		NewNetworksStormControlResource,
	}
}

func (p *CiscoMerakiProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewOrganizationsDataSource,
		NewOrganizationsNetworksDataSource,
		NewAdministeredIdentitiesMeDataSource,
		NewOrganizationsAdminsDataSource,
		NewOrganizationsSamlIdpsDataSource,
		NewOrganizationsInventoryDevicesDataSource,
		NewOrganizationsAdaptivePolicyAclsDataSource,
		NewOrganizationsSamlRolesDataSource,
		NewNetworkGroupPoliciesDataSource,
		NewNetworksAppliancePortsDataSource,
		NewOrganizationsCellularGatewayUplinkStatusesDataSource,
		NewOrganizationsLicensesDataSource,
		NewDevicesSwitchPortsStatusesDataSource,
		NewDevicesApplianceDhcpSubnetsDataSource,
		NewNetworksWirelessSsidsDataSource,
		NewNetworksSwitchQosRulesDataSource,
		NewNetworksApplianceVpnSiteToSiteVpnDatasource,
		NewNetworksSwitchMtuDataSource,
		NewDevicesManagementInterfaceDatasource,
		NewNetworksApplianceFirewallL3FirewallRulesDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &CiscoMerakiProvider{
			version: version,
		}
	}
}
