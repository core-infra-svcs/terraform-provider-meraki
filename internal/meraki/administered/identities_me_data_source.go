package administered

import (
	"context"
	"fmt"
	jsontypes2 "github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &AdministeredIdentitiesMeDataSource{}

func NewAdministeredIdentitiesMeDataSource() datasource.DataSource {
	return &AdministeredIdentitiesMeDataSource{}
}

// AdministeredIdentitiesMeDataSource defines the data source implementation.
type AdministeredIdentitiesMeDataSource struct {
	client *openApiClient.APIClient
}

// AdministeredIdentitiesMeDataSourceModel describes the data source data model.
type AdministeredIdentitiesMeDataSourceModel struct {
	Id                          types.String      `tfsdk:"id"`
	AuthenticationApiKeyCreated jsontypes2.Bool   `tfsdk:"authentication_api_key_created"`
	AuthenticationMode          jsontypes2.String `tfsdk:"authentication_mode"`
	AuthenticationSaml          jsontypes2.Bool   `tfsdk:"authentication_saml_enabled"`
	AuthenticationTwofactor     jsontypes2.Bool   `tfsdk:"authentication_two_factor_enabled"`
	Email                       jsontypes2.String `tfsdk:"email"`
	LastUsedDashboardAt         jsontypes2.String `tfsdk:"last_used_dashboard_at"`
	Name                        jsontypes2.String `tfsdk:"name"`
}

func (d *AdministeredIdentitiesMeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_administered_identities_me"
}

func (d *AdministeredIdentitiesMeDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Returns the identity of the current user",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Username",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes2.StringType,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "User email",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes2.StringType,
			},
			"last_used_dashboard_at": schema.StringAttribute{
				MarkdownDescription: "Last seen active on Dashboard UI",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes2.StringType,
			},
			"authentication_mode": schema.StringAttribute{
				MarkdownDescription: "Authentication mode",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes2.StringType,
			},
			"authentication_api_key_created": schema.BoolAttribute{
				MarkdownDescription: "If API key is created for this user",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes2.BoolType,
			},
			"authentication_two_factor_enabled": schema.BoolAttribute{
				MarkdownDescription: "If twoFactor authentication is enabled for this user",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes2.BoolType,
			},
			"authentication_saml_enabled": schema.BoolAttribute{
				MarkdownDescription: "If SAML authentication is enabled for this user",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes2.BoolType,
			},
		},
	}
}

func (d *AdministeredIdentitiesMeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*openApiClient.APIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *AdministeredIdentitiesMeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data AdministeredIdentitiesMeDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	inlineResp, httpResp, err := d.client.AdministeredApi.GetAdministeredIdentitiesMe(context.Background()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			utils.HttpDiagnostics(httpResp),
		)
		return
	}

	// Check for API success inlineResp code
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
		return
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	}

	data.Id = types.StringValue("example-id")
	data.Name = jsontypes2.StringValue(inlineResp.GetName())
	data.Email = jsontypes2.StringValue(inlineResp.GetEmail())
	data.LastUsedDashboardAt = jsontypes2.StringValue(inlineResp.GetLastUsedDashboardAt().Format(time.RFC3339))
	data.AuthenticationMode = jsontypes2.StringValue(inlineResp.Authentication.GetMode())
	data.AuthenticationApiKeyCreated = jsontypes2.BoolValue(inlineResp.Authentication.Api.Key.GetCreated())
	data.AuthenticationTwofactor = jsontypes2.BoolValue(inlineResp.Authentication.TwoFactor.GetEnabled())
	data.AuthenticationSaml = jsontypes2.BoolValue(inlineResp.Authentication.Saml.GetEnabled())

	// Write logs using the tflog package
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
