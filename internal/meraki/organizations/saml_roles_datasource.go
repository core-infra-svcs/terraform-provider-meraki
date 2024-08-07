package organizations

import (
	"context"
	"encoding/json"
	"fmt"
	jsontypes2 "github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &OrganizationsSamlRolesDataSource{}

func NewOrganizationsSamlRolesDataSource() datasource.DataSource {
	return &OrganizationsSamlRolesDataSource{}
}

// OrganizationsSamlRolesDataSource defines the data source implementation.
type OrganizationsSamlRolesDataSource struct {
	client *openApiClient.APIClient
}

// OrganizationsSamlRolesDataSourceModel describes the data source data model.
type OrganizationsSamlRolesDataSourceModel struct {
	Id    jsontypes2.String                     `tfsdk:"id"`
	OrgId jsontypes2.String                     `tfsdk:"organization_id"`
	List  []OrganizationSamlRoleDataSourceModel `tfsdk:"list"`
}

// OrganizationSamlRoleDataSourceModel describes the data source data model.
type OrganizationSamlRoleDataSourceModel struct {
	Id        jsontypes2.String                             `tfsdk:"id"`
	Role      jsontypes2.String                             `tfsdk:"role" json:"role"`
	OrgAccess jsontypes2.String                             `tfsdk:"org_access" json:"orgAccess"`
	Tags      []OrganizationsSamlRoleDataSourceModelTag     `tfsdk:"tags" json:"tags"`
	Networks  []OrganizationsSamlRoleDataSourceModelNetwork `tfsdk:"networks" json:"networks"`
}

type OrganizationsSamlRoleDataSourceModelNetwork struct {
	Id     jsontypes2.String `tfsdk:"id" json:"id"`
	Access jsontypes2.String `tfsdk:"access" json:"access"`
}

type OrganizationsSamlRoleDataSourceModelTag struct {
	Tag    jsontypes2.String `tfsdk:"tag" json:"tag"`
	Access jsontypes2.String `tfsdk:"access" json:"access"`
}

func (d *OrganizationsSamlRolesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizations_saml_roles"
}

func (d *OrganizationsSamlRolesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List the saml roles in this organization",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes2.StringType,
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "Organization ID",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 31),
				},
				CustomType: jsontypes2.StringType,
			},
			"list": schema.SetNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Role ID",
							Optional:            true,
							CustomType:          jsontypes2.StringType,
						},
						"role": schema.StringAttribute{
							MarkdownDescription: "The role of the SAML administrator",
							Optional:            true,
							CustomType:          jsontypes2.StringType,
						},
						"org_access": schema.StringAttribute{
							MarkdownDescription: "The privilege of the SAML administrator on the organization. Can be one of 'none', 'read-only', 'full' or 'enterprise'",
							Optional:            true,
							CustomType:          jsontypes2.StringType,
						},
						"tags": schema.SetNestedAttribute{
							Description: "The list of tags that the SAML administrator has privleges on.",
							Optional:    true,
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"tag": schema.StringAttribute{
										MarkdownDescription: "",
										Optional:            true,
										CustomType:          jsontypes2.StringType,
									},
									"access": schema.StringAttribute{
										MarkdownDescription: "",
										Optional:            true,
										CustomType:          jsontypes2.StringType,
									},
								},
							},
						},
						"networks": schema.SetNestedAttribute{
							Description: "The list of networks that the SAML administrator has privileges on.",
							Optional:    true,
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										MarkdownDescription: "",
										Optional:            true,
										CustomType:          jsontypes2.StringType,
									},
									"access": schema.StringAttribute{
										MarkdownDescription: "",
										Optional:            true,
										CustomType:          jsontypes2.StringType,
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (d *OrganizationsSamlRolesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *OrganizationsSamlRolesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OrganizationsSamlRolesDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := d.client.OrganizationsApi.GetOrganizationSamlRoles(context.Background(), data.OrgId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			utils.HttpDiagnostics(httpResp),
		)
	}

	// Check for API success inlineResp code
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("State Data", fmt.Sprintf("\n%s", data))
		return
	}

	// Save data into Terraform state
	if err = json.NewDecoder(httpResp.Body).Decode(&data.List); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	data.Id = jsontypes2.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read data source")
}
