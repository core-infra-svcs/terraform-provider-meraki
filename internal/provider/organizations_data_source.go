package provider

import (
	"context"
	"fmt"
	apiclient "github.com/core-infra-svcs/dashboard-api-go/client"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &OrganizationsDataSource{}

func NewOrganizationsDataSource() datasource.DataSource {
	return &OrganizationsDataSource{}
}

// OrganizationsDataSource defines the data source implementation.
type OrganizationsDataSource struct {
	client *apiclient.APIClient
}

// OrganizationsDataSourceModel describes the data source data model.
type OrganizationsDataSourceModel struct {
	Id   types.String                  `tfsdk:"id"`
	List []OrganizationDataSourceModel `tfsdk:"list"`
}

// OrganizationDataSourceModel describes the data source data model.
type OrganizationDataSourceModel struct {
	ApiEnabled     types.Bool   `tfsdk:"api_enabled"`
	CloudRegion    types.String `tfsdk:"cloud_region_name"`
	OrgId          types.String `tfsdk:"organization_id"`
	LicensingModel types.String `tfsdk:"licensing_model"`
	Name           types.String `tfsdk:"name"`
	Url            types.String `tfsdk:"url"`
}

func (d *OrganizationsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizations"
}

func (d *OrganizationsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "List the organizations that the user has privileges on",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
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
						},
						"cloud_region_name": schema.StringAttribute{
							MarkdownDescription: "Name of region",
							Optional:            true,
						},
						"organization_id": schema.StringAttribute{
							MarkdownDescription: "Organization ID",
							Optional:            true,
						},
						"licensing_model": schema.StringAttribute{
							MarkdownDescription: "Organization licensing model. Can be 'co-term', 'per-device', or 'subscription'.",
							Optional:            true,
							Validators: []validator.String{
								stringvalidator.ExactlyOneOf(
									path.MatchRoot("co-term"),
									path.MatchRoot("per-device"),
									path.MatchRoot("subscription"),
								),
							},
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "\"Organization name",
							Optional:            true,
						},
						"url": schema.StringAttribute{
							MarkdownDescription: "\"Organization URL",
							Optional:            true,
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

	for _, organization := range inlineResp {
		var result OrganizationDataSourceModel

		result.OrgId = types.StringValue(organization.GetId())
		result.Name = types.StringValue(organization.GetName())
		result.Url = types.StringValue(organization.GetUrl())
		result.ApiEnabled = types.BoolValue(*organization.GetApi().Enabled)
		result.LicensingModel = types.StringValue(*organization.GetLicensing().Model)
		result.CloudRegion = types.StringValue(organization.Cloud.Region.GetName())

		data.List = append(data.List, result)
	}

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read datasource")
}
