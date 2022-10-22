package provider

import (
	"context"
	"fmt"
	apiclient "github.com/core-infra-svcs/dashboard-api-golang/client"
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
	Api_enabled     types.Bool   `tfsdk:"api_enabled"`
	Cloud_region    types.String `tfsdk:"cloud_region"`
	Id              types.String `tfsdk:"id"`
	Licensing_model types.String `tfsdk:"licensing_model"`
	Name            types.String `tfsdk:"name"`
	Url             types.String `tfsdk:"url"`
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
						MarkdownDescription: "",
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
						MarkdownDescription: "",
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
					"id": {
						Description:         "Organization ID",
						MarkdownDescription: "",
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
						MarkdownDescription: "",
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
						MarkdownDescription: "",
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
						MarkdownDescription: "",
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

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// If applicable, this is a great opportunity to initialize any necessary
	// provider client data and make a call using it.
	response, r, err := d.client.OrganizationsApi.GetOrganizations(context.Background()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error when calling read: %v\n", r),
			"Could not complete read request: "+err.Error(),
		)
		return
	}

	// Map response body to resource schema attribute
	for _, organization := range response {

		var result OrganizationDataSourceModel

		result.Id = types.String{Value: organization.GetId()}
		result.Name = types.String{Value: organization.GetName()}
		result.Url = types.String{Value: organization.GetUrl()}
		result.Api_enabled = types.Bool{Value: *organization.GetApi().Enabled}
		result.Licensing_model = types.String{Value: *organization.GetLicensing().Model}
		result.Cloud_region = types.String{Value: organization.Cloud.Region.GetName()}

		data.List = append(data.List, result)
	}

	// For the purposes of this example code, hardcoding a response value to
	// save into the Terraform state.
	data.Id = types.String{Value: "example-id"}

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
