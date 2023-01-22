package provider

import (
	"context"
	"fmt"
	openApiClient "github.com/core-infra-svcs/dashboard-api-go/client"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
	client *openApiClient.APIClient
}

// OrganizationsNetworksDataSourceModel describes the data source data model.
type OrganizationsNetworksDataSourceModel struct {
	Id                      types.String `tfsdk:"id"`
	OrgId                   types.String `tfsdk:"organization_id"`
	ConfigTemplateId        types.String `tfsdk:"config_template_id"`
	IsBoundToConfigTemplate types.Bool   `tfsdk:"is_bound_to_config_template"`
	Tags                    types.Set    `tfsdk:"tags"`
	TagsFilterType          types.String `tfsdk:"tags_filter_type"`

	List []OrganizationsNetworksDataSourceModelList `tfsdk:"list"`
}

type OrganizationsNetworksDataSourceModelList struct {
	Id                      types.String `tfsdk:"id" json:"id"`
	OrganizationId          types.String `tfsdk:"organization_id" json:"organizationId"`
	Name                    types.String `tfsdk:"name" json:"name"`
	ProductTypes            types.Set    `tfsdk:"product_types" json:"productTypes"`
	TimeZone                types.String `tfsdk:"time_zone" json:"timeZone"`
	Tags                    types.Set    `tfsdk:"tags" json:"tags"`
	EnrollmentString        types.String `tfsdk:"enrollment_string"  json:"enrollmentString"`
	Url                     types.String `tfsdk:"url" json:"url"`
	Notes                   types.String `tfsdk:"notes" json:"notes"`
	isBoundToConfigTemplate types.Bool   `tfsdk:"is_bound_to_config_template" json:"isBoundToConfigTemplate"`
}

func (d *OrganizationsNetworksDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizations_networks"
}

func (d *OrganizationsNetworksDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "Organization ID",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 31),
				},
			},
			"config_template_id": schema.StringAttribute{
				MarkdownDescription: "config_template_id",
				Optional:            true,
				Computed:            true,
			},
			"is_bound_to_config_template": schema.BoolAttribute{
				MarkdownDescription: "",
				Optional:            true,
				Computed:            true,
			},
			"tags": schema.SetAttribute{
				Description: "Network tags",
				ElementType: types.StringType,
				Computed:    true,
				Optional:    true,
			},
			"tags_filter_type": schema.StringAttribute{
				MarkdownDescription: "",
				Optional:            true,
				Computed:            true,
			},
			"list": schema.ListNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "",
							Optional:            true,
							Computed:            true,
						},
						"organization_id": schema.StringAttribute{
							MarkdownDescription: "Organization ID",
							Optional:            true,
							Computed:            true,
							Validators: []validator.String{
								stringvalidator.LengthBetween(8, 31),
							},
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Network name",
							Optional:            true,
							Computed:            true,
						},
						"product_types": schema.SetAttribute{
							ElementType: types.StringType,
							Optional:    true,
							Computed:    true,
							Validators: []validator.Set{
								setvalidator.ValueStringsAre(
									stringvalidator.OneOf([]string{"appliance", "switch", "wireless", "systemsManager", "camera", "cellularGateway", "sensor"}...),
									stringvalidator.LengthAtLeast(5),
								),
							},
						},
						"timezone": schema.StringAttribute{
							MarkdownDescription: "Timezone of the network",
							Optional:            true,
							Computed:            true,
						},
						"tags": schema.SetAttribute{
							Description: "Network tags",
							ElementType: types.StringType,
							Computed:    true,
							Optional:    true,
						},
						"enrollment_string": schema.StringAttribute{
							MarkdownDescription: "A unique identifier which can be used for device enrollment or easy access through the Meraki SM Registration page or the Self Service Portal. Once enabled, a network enrollment strings can be changed but they cannot be deleted.",
							Optional:            true,
							Computed:            true,
							Validators: []validator.String{
								stringvalidator.NoneOf([]string{";", ":", "@", "=", "&", "$", "!", "‘", "“", ",", "?", ".", "(", ")", "{", "}", "[", "]", "\\", "*", "+", "/", "#", "<", ">", "|", "^", "%"}...),
								stringvalidator.LengthBetween(4, 50),
							},
						},
						"url": schema.StringAttribute{
							MarkdownDescription: "URL to the network Dashboard UI",
							Optional:            true,
							Computed:            true,
						},
						"notes": schema.StringAttribute{
							MarkdownDescription: "Notes for the network",
							Optional:            true,
							Computed:            true,
						},
						"is_bound_to_config_template": schema.BoolAttribute{
							MarkdownDescription: "If the network is bound to a config template",
							Optional:            true,
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *OrganizationsNetworksDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *OrganizationsNetworksDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OrganizationsNetworksDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	/*

		// Tags
		tags := []string{""} // []string | An optional parameter to filter networks by tags. The filtering is case-sensitive. If tags are included, 'tagsFilterType' should also be included (see below). (optional)
		if !data.Tags.IsUnknown() {
			tags = append(tags, data.Tags.String())
		}

		// tagsFilterType
		tagsFilterType := "" // string | An optional parameter of value 'withAnyTags' or 'withAllTags' to indicate whether to return networks which contain ANY or ALL of the included tags. If no type is included, 'withAnyTags' will be selected. (optional)
		if !data.TagsFilterType.IsUnknown() {
			tagsFilterType = data.TagsFilterType.ValueString()
		}

		// configTemplateId
		configTemplateId := "" // string | An optional parameter that is the ID of a config template. Will return all networks bound to that template. (optional)
		if !data.ConfigTemplateId.IsUnknown() {
			configTemplateId = data.ConfigTemplateId.ValueString()
		}

		//
		isBoundToConfigTemplate := false // bool | An optional parameter to filter config template bound networks. If configTemplateId is set, this cannot be false. (optional)
		if data.IsBoundToConfigTemplate.ValueBool() {
			isBoundToConfigTemplate = data.IsBoundToConfigTemplate.ValueBool()
		}

	*/

	perPage := int32(100000) // int32 | The number of entries per page returned. Acceptable range is 3 - 100000. Default is 1000. (optional)
	inlineResp, httpResp, err := d.client.OrganizationsApi.GetOrganizationNetworks(context.Background(),
		data.OrgId.ValueString()).PerPage(perPage).Execute()
	// .ConfigTemplateId(configTemplateId).IsBoundToConfigTemplate(isBoundToConfigTemplate).Tags(tags).TagsFilterType(tagsFilterType)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read datasource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	fmt.Printf("%v#", inlineResp)

	// Check for API success inlineResp code
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// collect diagnostics
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	} else {
		resp.Diagnostics.Append()
	}

	// save inlineResp data into Terraform state.
	data.Id = types.StringValue("example-id")
	/*
		for _, network := range inlineResp {
			var result OrganizationsNetworksDataSourceModelList

			result.Id = types.StringValue(network.GetId())
			result.OrganizationId = types.StringValue(network.GetOrganizationId())
			result.Name = types.StringValue(network.GetName())
			result.TimeZone = types.StringValue(network.GetTimeZone())


				// product types attribute
					if network.ProductTypes != nil {
						var pt []attr.Value
						for _, productTypeResp := range network.ProductTypes {
							pt = append(pt, types.StringValue(productTypeResp))
						}
						result.ProductTypes, _ = types.SetValue(types.StringType, pt)
					}

					// tags attribute
					if network.Tags != nil {
						var t []attr.Value
						for _, TagsTypeResp := range network.Tags {
							t = append(t, types.StringValue(TagsTypeResp))
						}
						result.Tags, _ = types.SetValue(types.StringType, t)
					}


			result.EnrollmentString = types.StringValue(network.GetEnrollmentString())
			result.Url = types.StringValue(network.GetUrl())
			result.Notes = types.StringValue(network.GetNotes())
			result.isBoundToConfigTemplate = types.BoolValue(network.GetIsBoundToConfigTemplate())

			data.List = append(data.List, result)
		}

	*/

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read data source")
}
