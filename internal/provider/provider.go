package provider

import (
	"context"
	apiclient "github.com/core-infra-svcs/dashboard-api-golang/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"os"
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
	Endpoint types.String `tfsdk:"endpoint"`
	Host     types.String `tfsdk:"host"`
	Path     types.String `tfsdk:"path"`
	ApiKey   types.String `tfsdk:"apikey"`
}

func (p *ScaffoldingProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "meraki"
	resp.Version = p.version
}

func (p *ScaffoldingProvider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"endpoint": {
				MarkdownDescription: "Example provider attribute",
				Optional:            true,
				Type:                types.StringType,
			},
			"path": {
				Type:     types.StringType,
				Optional: true,
				Computed: true,
			},
			"host": {
				Type:     types.StringType,
				Optional: true,
				Computed: true,
			},
			"apikey": {
				Type:      types.StringType,
				Required:  true,
				Sensitive: true,
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

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }

	// User must provide url to the provider
	var host string
	if data.Host.Null {
		host = "api.meraki.com"
	} else {
		host = data.Host.Value
	}

	if host == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find host",
			"host cannot be an empty string",
		)
		return
	}

	// must provide api path to the provider
	var path string
	if data.Path.Null {
		path = "/api/v1"
	} else {
		path = data.Path.Value
	}

	if path == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find path",
			"path cannot be an empty string",
		)
		return
	}

	// must provide apikey to the provider
	var apikey string
	if data.ApiKey.Unknown {
		// Cannot connect to client with an unknown value
		resp.Diagnostics.AddError(
			"Unable to create client",
			"Cannot use unknown value as apikey",
		)
		return
	}

	if data.ApiKey.Null {
		apikey = os.Getenv("MERAKI_DASHBOARD_API_KEY")
	} else {
		apikey = data.ApiKey.Value
	}

	if apikey == "" {
		// Error vs warning - empty value must stop execution
		resp.Diagnostics.AddError(
			"Unable to find apikey",
			"apikey cannot be an empty string",
		)
		return
	}

	// Example client configuration for data sources and resources

	// If the upstream provider SDK or HTTP client requires configuration, such
	// as authentication or logging, this is a great opportunity to do so.
	configuration := apiclient.NewConfiguration()
	configuration.AddDefaultHeader("X-Cisco-Meraki-API-Key", apikey)

	client := apiclient.NewAPIClient(configuration)
	resp.DataSourceData = client
	resp.ResourceData = client

}

func (p *ScaffoldingProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewExampleResource,
		NewOrganizationResource,
	}
}

func (p *ScaffoldingProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewExampleDataSource,
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
