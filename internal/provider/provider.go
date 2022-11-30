package provider

import (
	"context"
	"os"

	apiclient "github.com/core-infra-svcs/dashboard-api-go/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure ScaffoldingProvider satisfies various provider interfaces.
var _ provider.Provider = &ScaffoldingProvider{}
var _ provider.ProviderWithMetadata = &ScaffoldingProvider{}

// ScaffoldingProvider defines the provider implementation.
type ScaffoldingProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ScaffoldingProviderModel describes the provider data model.
type ScaffoldingProviderModel struct {
	ApiKey  types.String `tfsdk:"apikey"`
	BaseUrl types.String `tfsdk:"baseurl"`
}

func (p *ScaffoldingProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "meraki"
	resp.Version = p.version
}

func (p *ScaffoldingProvider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"apikey": {
				Type:      types.StringType,
				Optional:  true,
				Sensitive: true,
			},
			"baseurl": {
				Type:     types.StringType,
				Optional: true,
			},
		},
	}, nil
}

func (p *ScaffoldingProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data ScaffoldingProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// api key
	var apiKey string
	if data.ApiKey.IsNull() {
		apiKey = os.Getenv("MERAKI_DASHBOARD_API_KEY")
		if apiKey == "" {
			// Error vs warning - empty value must stop execution
			resp.Diagnostics.AddError(
				"Unable to set apiKey from env vars",
				"api key must be set",
			)
			return
		}
	}

	// base url
	var baseUrl string
	if data.BaseUrl.IsNull() {
		baseUrl = "https://api.meraki.com/api/v1"
	}

	// Example client configuration for data sources and resources

	// If the upstream provider SDK or HTTP client requires configuration, such
	// as authentication or logging, this is a great opportunity to do so.
	configuration := apiclient.NewConfiguration()
	configuration.AddDefaultHeader("X-Cisco-Meraki-API-Key", apiKey)
	configuration.Host = os.Getenv(baseUrl)

	client := apiclient.NewAPIClient(configuration)
	resp.DataSourceData = client
	resp.ResourceData = client

}

func (p *ScaffoldingProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewOrganizationResource,
		NewOrganizationsAdminResource,
	}
}

func (p *ScaffoldingProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewOrganizationsDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ScaffoldingProvider{
			version: version,
		}
	}
}
