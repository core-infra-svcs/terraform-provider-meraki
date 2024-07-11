package provider

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/meraki/administered"
	devices2 "github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/meraki/devices"
	networks2 "github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/meraki/networks"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/meraki/networks/appliance"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/meraki/networks/switch"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/meraki/networks/wireless"
	organizations2 "github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/meraki/organizations"
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
					"Example: `https://api.meraki.com`" +
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
					"Example: `/api/v1`",
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

// BearerAuthTransport Custom transport to add bearer token in the Authorization header
type BearerAuthTransport struct {
	Transport *http.Transport
	Token     string
}

func (t *BearerAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
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
	authenticatedTransport := &BearerAuthTransport{
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
		organizations2.NewOrganizationResource,
		organizations2.NewOrganizationSamlResource,
		organizations2.NewOrganizationsSamlIdpResource,
		organizations2.NewOrganizationsAdminResource,
		organizations2.NewOrganizationsAdaptivePolicyAclResource,
		networks2.NewNetworkResource,
		organizations2.NewOrganizationsSamlRolesResource,
		_switch.NewNetworksSwitchSettingsResource,
		networks2.NewNetworksSnmpResource,
		organizations2.NewOrganizationsSnmpResource,
		networks2.NewNetworksSettingsResource,
		appliance.NewNetworksApplianceFirewallL3FirewallRulesResource,
		appliance.NewNetworksApplianceFirewallL7FirewallRulesResource,
		organizations2.NewOrganizationsApplianceVpnVpnFirewallRulesResource,
		networks2.NewNetworksTrafficAnalysisResource,
		networks2.NewNetworksNetflowResource,
		networks2.NewNetworksSyslogServersResource,
		appliance.NewNetworksApplianceVlansSettingsResource,
		appliance.NewNetworksApplianceSettingsResource,
		appliance.NewNetworksApplianceFirewallSettingsResource,
		_switch.NewNetworksSwitchQosRuleResource,
		_switch.NewNetworksSwitchDscpToCosMappingsResource,
		_switch.NewNetworksSwitchMtuResource,
		networks2.NewNetworksGroupPolicyResource,
		organizations2.NewOrganizationsLicenseResource,
		wireless.NewNetworksWirelessSsidsFirewallL3FirewallRulesResource,
		wireless.NewNetworksWirelessSsidsFirewallL7FirewallRulesResource,
		devices2.NewDevicesResource,
		organizations2.NewOrganizationsClaimResource,
		networks2.NewNetworksDevicesClaimResource,
		appliance.NewNetworkApplianceStaticRoutesResource,
		networks2.NewNetworksCellularGatewaySubnetPoolResource,
		networks2.NewNetworksCellularGatewayUplinkResource,
		wireless.NewNetworksWirelessSsidsSplashSettingsResource,
		devices2.NewDevicesCellularSimsResource,
		devices2.NewDevicesTestAccDevicesManagementInterfaceResourceResource,
		appliance.NewNetworksApplianceVpnSiteToSiteVpnResource,
		_switch.NewDevicesSwitchPortsCycleResource,
		appliance.NewNetworksApplianceTrafficShapingUplinkBandWidthResource,
		appliance.NewNetworksApplianceVLANResource,
		_switch.NewDevicesSwitchPortResource,
		appliance.NewNetworksAppliancePortsResource,
		wireless.NewNetworksWirelessSsidsResource,
		networks2.NewNetworksStormControlResource,
	}
}

func (p *CiscoMerakiProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		organizations2.NewOrganizationsDataSource,
		organizations2.NewOrganizationsNetworksDataSource,
		administered.NewAdministeredIdentitiesMeDataSource,
		devices2.NewNetworkDevicesDataSource,
		organizations2.NewOrganizationsAdminsDataSource,
		organizations2.NewOrganizationsSamlIdpsDataSource,
		organizations2.NewOrganizationsInventoryDevicesDataSource,
		organizations2.NewOrganizationsAdaptivePolicyAclsDataSource,
		organizations2.NewOrganizationsSamlRolesDataSource,
		networks2.NewNetworkGroupPoliciesDataSource,
		appliance.NewNetworksAppliancePortsDataSource,
		appliance.NewNetworksApplianceVLANsDatasource,
		organizations2.NewOrganizationsCellularGatewayUplinkStatusesDataSource,
		organizations2.NewOrganizationsLicensesDataSource,
		appliance.NewNetworksApplianceVlansSettingsDatasource,
		_switch.NewDevicesSwitchPortsStatusesDataSource,
		appliance.NewDevicesApplianceDhcpSubnetsDataSource,
		wireless.NewNetworksWirelessSsidsDataSource,
		_switch.NewNetworksSwitchQosRulesDataSource,
		appliance.NewNetworksApplianceVpnSiteToSiteVpnDatasource,
		_switch.NewNetworksSwitchMtuDataSource,
		devices2.NewDevicesManagementInterfaceDatasource,
		appliance.NewNetworksApplianceFirewallL3FirewallRulesDataSource,
		networks2.NewNetworksSwitchStormControlDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &CiscoMerakiProvider{
			version: version,
		}
	}
}
