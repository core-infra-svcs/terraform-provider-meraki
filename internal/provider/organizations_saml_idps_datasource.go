package provider

import (
	"context"
	"fmt"
	apiclient "github.com/core-infra-svcs/dashboard-api-go/client"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
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
	client *apiclient.APIClient
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

func (d *OrganizationsSamlIdpsDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "OrganizationsSamlIdps data source - ",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description:         "Example identifier needed for terraform",
				MarkdownDescription: "Example identifier needed for terraform",
				Type:                types.StringType,
				Required:            false,
				Optional:            false,
				Computed:            true,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"organization_id": {
				Description:         "Organization ID",
				MarkdownDescription: "Organization ID",
				Type:                types.StringType,
				Required:            true,
				Optional:            false,
				Computed:            false,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"list": {
				Description:         "List of Saml IDPs",
				MarkdownDescription: "List of Saml IDPs",
				Optional:            true,
				Attributes: tfsdk.SetNestedAttributes(map[string]tfsdk.Attribute{
					"consumer_url": {
						Description:         "URL that is consuming SAML Identity Provider (IdP)",
						MarkdownDescription: "URL that is consuming SAML Identity Provider (IdP)",
						Type:                types.StringType,
						Required:            false,
						Optional:            true,
						Computed:            true,
						Sensitive:           false,
						Attributes:          nil,
						DeprecationMessage:  "",
						Validators:          nil,
						PlanModifiers:       nil,
					},
					"idp_id": {
						Description:         "ID associated with the SAML Identity Provider (IdP)",
						MarkdownDescription: "ID associated with the SAML Identity Provider (IdP)",
						Type:                types.StringType,
						Required:            false,
						Optional:            true,
						Computed:            true,
						Sensitive:           false,
						Attributes:          nil,
						DeprecationMessage:  "",
						Validators:          nil,
						PlanModifiers:       nil,
					},
					"slo_logout_url": {
						Description:         "Dashboard will redirect users to this URL when they sign out.",
						MarkdownDescription: "Dashboard will redirect users to this URL when they sign out.",
						Type:                types.StringType,
						Required:            false,
						Optional:            true,
						Computed:            true,
						Sensitive:           false,
						Attributes:          nil,
						DeprecationMessage:  "",
						Validators:          nil,
						PlanModifiers:       nil,
					},
					"x_509_cert_sha1_fingerprint": {
						Description:         "Fingerprint (SHA1) of the SAML certificate provided by your Identity Provider (IdP). This will be used for encryption / validation.",
						MarkdownDescription: "Fingerprint (SHA1) of the SAML certificate provided by your Identity Provider (IdP). This will be used for encryption / validation.",
						Type:                types.StringType,
						Required:            false,
						Optional:            true,
						Computed:            true,
						Sensitive:           false,
						Attributes:          nil,
						DeprecationMessage:  "",
						Validators:          nil,
						PlanModifiers:       nil,
					},
				}),
			},
		},
	}, nil
}

func (d *OrganizationsSamlIdpsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*apiclient.APIClient)

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
