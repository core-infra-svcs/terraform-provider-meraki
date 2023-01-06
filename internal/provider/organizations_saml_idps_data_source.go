package provider

import (
	"context"
	"fmt"
	openApiClient "github.com/core-infra-svcs/dashboard-api-go/client"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &OrganizationsSamlIdpsDataSource{}

func NewOrganizationsSamlIdpsDataSource() datasource.DataSource {
	return &OrganizationsSamlIdpsDataSource{}
}

// OrganizationsSamlIdpsDataSource defines the data source implementation.
type OrganizationsSamlIdpsDataSource struct {
	client *openApiClient.APIClient
}

type OrganizationsSamlIdpsDataSourceModel struct {
	Id             types.String                          `tfsdk:"id"`
	OrganizationId types.String                          `tfsdk:"organization_id"`
	List           []OrganizationsSamlIdpDataSourceModel `tfsdk:"list"`
}

// OrganizationsSamlIdpDataSourceModel describes the data source data model.
type OrganizationsSamlIdpDataSourceModel struct {
	ConsumerUrl             types.String `tfsdk:"consumer_url"`
	IdpId                   types.String `tfsdk:"idp_id"`
	SloLogOutUrl            types.String `tfsdk:"slo_logout_url"`
	X509CertSha1FingerPrint types.String `tfsdk:"x_509_cert_sha1_fingerprint"`
}

func (d *OrganizationsSamlIdpsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizations_saml_idps"
}

func (d *OrganizationsSamlIdpsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List the SAML IdPs in your organization.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "Organization ID",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 31),
				},
			},
			"list": schema.ListNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"consumer_url": schema.StringAttribute{
							MarkdownDescription: "URL that is consuming SAML Identity Provider (IdP)",
							Optional:            true,
						},
						"idp_id": schema.StringAttribute{
							MarkdownDescription: "ID associated with the SAML Identity Provider (IdP)",
							Optional:            true,
						},
						"slo_logout_url": schema.StringAttribute{
							MarkdownDescription: "Dashboard will redirect users to this URL when they sign out.",
							Optional:            true,
						},
						"x_509_cert_sha1_fingerprint": schema.StringAttribute{
							MarkdownDescription: "Fingerprint (SHA1) of the SAML certificate provided by your Identity Provider (IdP). This will be used for encryption / validation.",
							Optional:            true,
						},
					},
				},
			},
		},
	}
}

func (d *OrganizationsSamlIdpsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *OrganizationsSamlIdpsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OrganizationsSamlIdpsDataSourceModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Initialize provider client and make API call
	inlineResp, httpResp, err := d.client.OrganizationsApi.GetOrganizationSamlIdps(context.Background(), data.OrganizationId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read datasource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// Check for API success inlineResp code
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// collect diagnostics
	tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	} else {
		resp.Diagnostics.Append()
	}

	// save inlineResp data into Terraform state.
	data.Id = types.StringValue("example-id")

	if resp.Diagnostics.HasError() {
		return
	}

	for _, samlIdp := range inlineResp {
		var result OrganizationsSamlIdpDataSourceModel

		result.IdpId = types.StringValue(samlIdp.GetIdpId())
		result.ConsumerUrl = types.StringValue(samlIdp.GetConsumerUrl())
		result.SloLogOutUrl = types.StringValue(samlIdp.GetSloLogoutUrl())
		result.X509CertSha1FingerPrint = types.StringValue(samlIdp.GetX509certSha1Fingerprint())
		data.List = append(data.List, result)
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read datasource")
}
