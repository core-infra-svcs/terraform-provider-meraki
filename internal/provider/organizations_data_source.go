package provider

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &OrganizationsDataSource{}

func NewOrganizationsDataSource() datasource.DataSource {
	return &OrganizationsDataSource{}
}

// OrganizationsDataSource defines the data source implementation.
type OrganizationsDataSource struct {
	client *openApiClient.APIClient
}

// OrganizationsDataSourceModel describes the data source data model.
type OrganizationsDataSourceModel struct {
	Id   jsontypes.String                   `tfsdk:"id"`
	List []OrganizationsDataSourceModelList `tfsdk:"list"`
}

// OrganizationsDataSourceModelList describes the data source data model.
type OrganizationsDataSourceModelList struct {
	ApiEnabled     jsontypes.Bool   `tfsdk:"api_enabled"`
	CloudRegion    jsontypes.String `tfsdk:"cloud_region_name"`
	OrgId          jsontypes.String `tfsdk:"organization_id"`
	LicensingModel jsontypes.String `tfsdk:"licensing_model"`
	Name           jsontypes.String `tfsdk:"name"`
	Url            jsontypes.String `tfsdk:"url"`
}

func (d *OrganizationsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizations"
}

func (d *OrganizationsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List the organizations that the user has privileges on",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.StringType,
			},
			"list": schema.ListNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"api_enabled": schema.BoolAttribute{
							MarkdownDescription: "Enable API access",
							Optional:            true,
							CustomType:          jsontypes.BoolType,
						},
						"cloud_region_name": schema.StringAttribute{
							MarkdownDescription: "Name of region",
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
						"organization_id": schema.StringAttribute{
							MarkdownDescription: "Organization ID",
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
						"licensing_model": schema.StringAttribute{
							MarkdownDescription: "Organization licensing model. Can be 'co-term', 'per-device', or 'subscription'.",
							Optional:            true,
							CustomType:          jsontypes.StringType,
							Validators: []validator.String{
								stringvalidator.OneOf([]string{"co-term", "per-device", "subscription"}...),
								stringvalidator.LengthAtLeast(7),
							},
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Organization name",
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
						"url": schema.StringAttribute{
							MarkdownDescription: "Organization URL",
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
					},
				},
			},
		},
	}
}

func (d *OrganizationsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *OrganizationsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OrganizationsDataSourceModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Initialize provider client and make API call
	inlineResp, httpResp, err := d.client.OrganizationsApi.GetOrganizations(context.Background()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
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

	// save inlineResp data into Terraform state.
	data.Id = jsontypes.StringValue("example-id")

	for _, organization := range inlineResp {
		var result OrganizationsDataSourceModelList

		result.OrgId = jsontypes.StringValue(organization.GetId())
		result.Name = jsontypes.StringValue(organization.GetName())
		result.Url = jsontypes.StringValue(organization.GetUrl())
		result.ApiEnabled = jsontypes.BoolValue(*organization.GetApi().Enabled)
		result.LicensingModel = jsontypes.StringValue(*organization.GetLicensing().Model)
		result.CloudRegion = jsontypes.StringValue(organization.Cloud.Region.GetName())

		data.List = append(data.List, result)
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read datasource")
}
