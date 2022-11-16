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
	CloudRegion    types.String `tfsdk:"cloud_region"`
	OrgId          types.String `tfsdk:"organization_id"`
	LicensingModel types.String `tfsdk:"licensing_model"`
	Name           types.String `tfsdk:"name"`
	Url            types.String `tfsdk:"url"`
}

func (d *OrganizationsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizations"
}

func (d *OrganizationsDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Organizations data source",

		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "Example identifier",
				Type:                types.StringType,
				Computed:            true,
			},
			"list": {
				MarkdownDescription: "List of organizations",
				Optional:            true,
				Attributes: tfsdk.SetNestedAttributes(map[string]tfsdk.Attribute{
					"api_enabled": {
						Description:         "Enable API access",
						MarkdownDescription: "Enable API access",
						Type:                types.BoolType,
						Required:            false,
						Optional:            true,
						Computed:            true,
						Sensitive:           false,
						Attributes:          nil,
						DeprecationMessage:  "",
						Validators:          nil,
						PlanModifiers:       nil,
					},
					"cloud_region": {
						Description:         "Region info",
						MarkdownDescription: "Region info",
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
					"organization_id": {
						Description:         "Organization ID",
						MarkdownDescription: "Organization ID",
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
					"licensing_model": {
						Description:         "Organization licensing model. Can be 'co-term', 'per-device', or 'subscription'.",
						MarkdownDescription: "Organization licensing model. Can be 'co-term', 'per-device', or 'subscription'.",
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
					"name": {
						Description:         "Organization name",
						MarkdownDescription: "Organization name",
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
					"url": {
						Description:         "Organization URL",
						MarkdownDescription: "Organization URL",
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
