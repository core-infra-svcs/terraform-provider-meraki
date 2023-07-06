package provider

import (
	"context"
	"os"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
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

const (
	BaseUrl = "api.meraki.com"

	// TODO - Expose HTTP Client Configuration
	//	SingleRequestTimeout = ""
	//	CertificatePath = ""
	//	RequestsProxy = ""
	//	WaitOnRateLimit = ""
	//	Nginx429RetryWaitTime = ""
	//	ActionBatchRetryWaitTime = ""
	//	Retry4xxError = ""
	//	Retry4xxErrorWaitTime = ""
	//	MaximumRetries = ""
	//	Simulate = ""
	//	BeGeoId = ""
	//	Caller = ""
	//	UseIteratorForGetPage = ""

)

// CiscoMerakiProviderModel describes the provider data model.
type CiscoMerakiProviderModel struct {
	ApiKey  types.String `tfsdk:"api_key"`
	BaseUrl types.String `tfsdk:"base_url"`

	// TODO - Expose HTTP Client Configuration
	// SingleRequestTimeout types.String `tfsdk:"single_request_timeout"`
	//	CertificatePath           types.String `tfsdk:"certificate_path"`
	//	RequestsProxy             types.String `tfsdk:"requests_proxy"`
	//	WaitOnRateLimit           types.String `tfsdk:"wait_on_rate_limit"`
	//	Nginx429RetryWaitTime     types.String `tfsdk:"nginx_429_retry_wait_time"`
	//	ActionBatchRetryWaitTime  types.String `tfsdk:"action_batch_retry_wait_time"`
	//	Retry4xxError             types.String `tfsdk:"retry_4xx_error"`
	//	Retry4xxErrorWaitTime     types.String `tfsdk:"retry_4xx_error_wait_time"`
	//	MaximumRetries            types.String `tfsdk:"maximum_retries"`
	//	Simulate 		  types.String `tfsdk:" simulate"`
	//	BeGeoId  		  types.String `tfsdk:"be_geo_id"`
	//	Caller                	  types.String `tfsdk:"caller"`
	//	UseIteratorForGetPage 	  types.String `tfsdk:"use_iterator_for_get_page"`
}

func (p *CiscoMerakiProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "meraki"
	resp.Version = p.version
}

func (p *CiscoMerakiProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"api_key": schema.StringAttribute{
				MarkdownDescription: "Meraki Dashboard API Key",
				Optional:            true,
			},
			"base_url": schema.StringAttribute{
				Description: "Endpoint for Meraki Dashboard API",
				MarkdownDescription: "The API version must be specified in the URL:" +
					"Example: `api.meraki.com" +
					"For organizations hosted in the China dashboard, use: `api.meraki.cn`" +
					"Defaults to `" + BaseUrl + "`.",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`(?:[a-zA-Z0-9]{1,62}(?:[-\.][a-zA-Z0-9]{1,62})+)(:\d+)?$`),
						"The API version must be specified in the URL. Example: api.meraki.com",
					),
				},
			},
		},
	}
}

func (p *CiscoMerakiProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data CiscoMerakiProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Check environment variables
	apiKey := os.Getenv("MERAKI_DASHBOARD_API_KEY")

	// Check configuration data, which should take precedence over
	// environment variable data, if found.
	if data.ApiKey.ValueString() != "" {
		apiKey = data.ApiKey.ValueString()
	}

	// Return err if API key is missing
	if apiKey == "" {
		resp.Diagnostics.AddError(
			"Missing API Key",
			"api key must be set either through the provider config or environmental variables.",
		)
		return
	}

	// base url
	var baseUrl string
	if data.BaseUrl.IsNull() {
		baseUrl = BaseUrl
	} else {
		baseUrl = data.BaseUrl.ValueString()
	}

	// If the upstream provider SDK or HTTP client requires configuration, such
	// as authentication or logging, this is a great opportunity to do so.
	// TODO - Add API version, retry limit and other such values to HTTP client and expose here
	configuration := openApiClient.NewConfiguration()
	configuration.AddDefaultHeader("X-Cisco-Meraki-API-Key", apiKey)
	configuration.Host = baseUrl
	configuration.UserAgent = configuration.UserAgent + "terraform" + p.version

	// enable debug for provider development
	if p.version == "dev" {
		configuration.Debug = true
	}

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
		NewNetworksSwitchQosRulesResource,
		NewNetworksSwitchDscpToCosMappingsResource,
		NewNetworksSwitchMtuResource,
		NewNetworksGroupPolicyResource,
		NewOrganizationsLicenseResource,
		NewNetworksWirelessSsidsFirewallL3FirewallRulesResource,
		NewNetworksWirelessSsidsFirewallL7FirewallRulesResource,
		NewDevicesResource,
		NewOrganizationsClaimResource,
		NewNetworksDevicesClaimResource,
    NewNetworksCellularGatewaySubnetPoolResource,
    NewNetworksCellularGatewayUplinkResource,
    NewNetworksWirelessSsidsSplashSettingsResource,
    NewDevicesCellularSimsResource,
	}
}

func (p *CiscoMerakiProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewOrganizationsDataSource,
		NewOrganizationsNetworksDataSource,
		NewAdministeredIdentitiesMeDataSource,
		NewOrganizationsAdminsDataSource,
		NewOrganizationsSamlIdpsDataSource,
		NewOrganizationsAdaptivePolicyAclsDataSource,
		NewOrganizationsSamlRolesDataSource,
		NewNetworkGroupPoliciesDataSource,
		NewOrganizationsCellularGatewayUplinkStatusesDataSource,
		NewOrganizationsLicensesDataSource,
		NewDevicesSwitchPortsStatusesDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &CiscoMerakiProvider{
			version: version,
		}
	}
}
