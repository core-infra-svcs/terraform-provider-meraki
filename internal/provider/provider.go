package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
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
	EncryptionKey         types.String `tfsdk:"encryption_key"`
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
			"encryption_key": schema.StringAttribute{
				Optional:            true,
				Description:         "Encryption key for encrypting sensitive values.",
				MarkdownDescription: "Encryption key for encrypting sensitive values.",
			},
		},
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &CiscoMerakiProvider{
			version: version,
		}
	}
}
