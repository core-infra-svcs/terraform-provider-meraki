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
var _ datasource.DataSource = &OrganizationsNetworksDataSource{}

func NewOrganizationsNetworksDataSource() datasource.DataSource {
	return &OrganizationsNetworksDataSource{}
}

// OrganizationsNetworksDataSource defines the data source implementation.
type OrganizationsNetworksDataSource struct {
	client *apiclient.APIClient
}

// OrganizationsNetworksDataSourceModel describes the data source data model.
type OrganizationsNetworksDataSourceModel struct {
	Id                      types.String                               `tfsdk:"id"`
	OrgId                   types.String                               `tfsdk:"organization_id"`             // path var
	ConfigTemplateId        types.String                               `tfsdk:"config_template_id"`          // query params
	IsBoundToConfigTemplate types.Bool                                 `tfsdk:"is_bound_to_config_template"` // query params
	Tags                    []OrganizationsNetworksDataSourceModelTag  `tfsdk:"tags"`                        // query params
	TagsFilterType          types.String                               `tfsdk:"tagsFilterType"`              // query params
	List                    []OrganizationsNetworksDataSourceModelList `tfsdk:"list"`
}

type OrganizationsNetworksDataSourceModelList struct {
	Id                      types.String   `tfsdk:"id"`
	OrganizationId          types.String   `tfsdk:"organizationId"`
	Name                    types.String   `tfsdk:"name"`
	ProductTypes            []types.String `tfsdk:"product_types"`
	TimeZone                types.String   `tfsdk:"time_zone"`
	Tags                    []types.String `tfsdk:"tags"`
	EnrollmentString        types.String   `tfsdk:"enrollment_string"`
	Url                     types.String   `tfsdk:"url"`
	Notes                   types.String   `tfsdk:"notes"`
	isBoundToConfigTemplate types.Bool     `tfsdk:"is_bound_to_config_template"`
}

type OrganizationsNetworksDataSourceModelTag struct {
	Tag types.String `tfsdk:"tag"`
}

func (d *OrganizationsNetworksDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizations_networks"
}

func (d *OrganizationsNetworksDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "OrganizationsNetworks resource - ",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description:         "Example identifier",
				MarkdownDescription: "Example identifier",
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
				Required:            false,
				Optional:            true,
				Computed:            true,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"config_template_id": {
				Description:         "config_template_id",
				MarkdownDescription: "config_template_id",
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
			"is_bound_to_config_template": {
				Description:         "is_bound_to_config_template",
				MarkdownDescription: "is_bound_to_config_template",
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
			"tags": {
				Description:         "list of tags",
				MarkdownDescription: "list of tags",
				Computed:            true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"tag": {
						Description:         "tag",
						MarkdownDescription: "tag",
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
				})},
			"tagsFilterType": {
				Description:         "tagsFilterType",
				MarkdownDescription: "tagsFilterType",
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
			"list": {
				MarkdownDescription: "OrganizationsNetworksResourceModelList of networks",
				Optional:            true,
				Attributes: tfsdk.SetNestedAttributes(map[string]tfsdk.Attribute{
					"id": {
						Description:         "network_id",
						MarkdownDescription: "network_id",
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
						Description:         "network name",
						MarkdownDescription: "network name",
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
					"product_types": {
						Description:         "product_types",
						MarkdownDescription: "product_types",
						Type:                types.ListType{ElemType: types.StringType},
						Required:            false,
						Optional:            true,
						Computed:            true,
						Sensitive:           false,
						Attributes:          nil,
						DeprecationMessage:  "",
						Validators:          nil,
						PlanModifiers:       nil,
					},
					"time_zone": {
						Description:         "time_zone",
						MarkdownDescription: "time_zone",
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
					"tags": {
						Description:         "tags",
						MarkdownDescription: "tags",
						Type:                types.ListType{ElemType: types.StringType},
						Required:            false,
						Optional:            true,
						Computed:            true,
						Sensitive:           false,
						Attributes:          nil,
						DeprecationMessage:  "",
						Validators:          nil,
						PlanModifiers:       nil,
					},
					"enrollment_string": {
						Description:         "enrollment_string",
						MarkdownDescription: "enrollment_string",
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
						Description:         "URL",
						MarkdownDescription: "URL",
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
					"notes": {
						Description:         "notes",
						MarkdownDescription: "notes",
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
					"is_bound_to_config_template": {
						Description:         "is_bound_to_config_template",
						MarkdownDescription: "is_bound_to_config_template",
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
				}),
			},
		},
	}, nil
}

func (d *OrganizationsNetworksDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *OrganizationsNetworksDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OrganizationsNetworksDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	// Check for required parameters
	if len(data.OrgId.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing OrganizationId", fmt.Sprintf("Value: %s", data.OrgId.ValueString()))
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Tags
	var tags []string
	if len(data.Tags) < 0 {
		for _, tag := range data.Tags {
			tags = append(tags, tag.Tag.ValueString())
		}
	}

	perPage := int32(100000) // int32 | The number of entries per page returned. Acceptable range is 3 - 100000. Default is 1000. (optional)
	startingAfter := ""      // string | A token used by the server to indicate the start of the page. Often this is a timestamp or an ID but it is not limited to those. This parameter should not be defined by client applications. The link for the first, last, prev, or next page in the HTTP Link header should define it. (optional)
	endingBefore := ""       // string | A token used by the server to indicate the end of the page. Often this is a timestamp or an ID but it is not limited to those. This parameter should not be defined by client applications. The link for the first, last, prev, or next page in the HTTP Link header should define it. (optional)

	// TODO -- PAGINATION
	inlineResp, httpResp, err := d.client.OrganizationsApi.GetOrganizationNetworks(context.Background(),
		data.OrgId.ValueString()).ConfigTemplateId(data.ConfigTemplateId.ValueString()).IsBoundToConfigTemplate(
		data.IsBoundToConfigTemplate.ValueBool()).Tags(tags).TagsFilterType(
		data.TagsFilterType.ValueString()).PerPage(perPage).StartingAfter(
		startingAfter).EndingBefore(endingBefore).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// Check for API success response code
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
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%s", data))
		return
	}

	// save inlineResp data into Terraform state.
	data.Id = types.StringValue("example-id")

	for _, network := range inlineResp {
		var result OrganizationsNetworksDataSourceModelList

		result.Id = types.StringValue(network.GetId())
		result.OrganizationId = types.StringValue(network.GetOrganizationId())
		result.Name = types.StringValue(network.GetName())

		// ProductTypes
		var productTypes []types.String
		for _, productType := range network.ProductTypes {
			productTypes = append(productTypes, types.StringValue(productType))
		}

		result.ProductTypes = productTypes

		result.TimeZone = types.StringValue(network.GetTimeZone())

		// Tags
		var tags []types.String
		for _, tag := range network.ProductTypes {
			tags = append(tags, types.StringValue(tag))
		}

		result.Tags = tags

		result.EnrollmentString = types.StringValue(network.GetEnrollmentString())
		result.Url = types.StringValue(network.GetUrl())
		result.Notes = types.StringValue(network.GetNotes())
		result.isBoundToConfigTemplate = types.BoolValue(network.GetIsBoundToConfigTemplate())

		data.List = append(data.List, result)
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read data source")
}
