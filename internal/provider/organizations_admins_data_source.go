package provider

import (
	"context"
	"encoding/json"
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
var _ datasource.DataSource = &OrganizationsAdminsDataSource{}

func NewOrganizationsAdminsDataSource() datasource.DataSource {
	return &OrganizationsAdminsDataSource{}
}

// OrganizationsAdminsDataSource defines the data source implementation.
type OrganizationsAdminsDataSource struct {
	client *openApiClient.APIClient
}

// OrganizationsAdminsDataSourceModel describes the data source data model.
type OrganizationsAdminsDataSourceModel struct {
	Id    types.String                        `tfsdk:"id"`
	OrgId types.String                        `tfsdk:"organization_id"`
	List  []OrganizationAdminsDataSourceModel `tfsdk:"list"`
}

// OrganizationAdminsDataSourceModel describes the data source data model.
type OrganizationAdminsDataSourceModel struct {
	Id                   types.String                               `tfsdk:"id"`
	Name                 types.String                               `tfsdk:"name"`
	Email                types.String                               `tfsdk:"email"`
	OrgAccess            types.String                               `tfsdk:"org_access"`
	AccountStatus        types.String                               `tfsdk:"account_status"`
	TwoFactorAuthEnabled types.Bool                                 `tfsdk:"two_factor_auth_enabled"`
	HasApiKey            types.Bool                                 `tfsdk:"has_api_key"`
	LastActive           types.String                               `tfsdk:"last_active"`
	Tags                 []OrganizationAdminsDataSourceModelTag     `tfsdk:"tags"`
	Networks             []OrganizationAdminsDataSourceModelNetwork `tfsdk:"networks"`
	AuthenticationMethod types.String                               `tfsdk:"authentication_method"`
}

type OrganizationAdminsDataSourceModelNetwork struct {
	Id     types.String `tfsdk:"id"`
	Access types.String `tfsdk:"access"`
}

type OrganizationAdminsDataSourceModelTag struct {
	Tag    types.String `tfsdk:"tag"`
	Access types.String `tfsdk:"access"`
}

func (d *OrganizationsAdminsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizations_admins"
}

func (d *OrganizationsAdminsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List the dashboard administrators in this organization",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "Organization ID",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 31),
				},
			},
			"list": schema.SetNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							MarkdownDescription: "Admin ID",
							Optional:            true,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "The name of the dashboard administrator",
							Optional:            true,
						},
						"email": schema.StringAttribute{
							MarkdownDescription: "The email of the dashboard administrator. This attribute can not be updated.",
							Optional:            true,
						},
						"org_access": schema.StringAttribute{
							MarkdownDescription: "The privilege of the dashboard administrator on the organization. Can be one of 'full', 'read-only', 'enterprise' or 'none'",
							Optional:            true,
						},
						"account_status": schema.StringAttribute{
							MarkdownDescription: "",
							Optional:            true,
						},
						"two_factor_auth_enabled": schema.BoolAttribute{
							MarkdownDescription: "",
							Optional:            true,
						},
						"has_api_key": schema.BoolAttribute{
							MarkdownDescription: "",
							Optional:            true,
						},
						"last_active": schema.StringAttribute{
							MarkdownDescription: "",
							Optional:            true,
						},
						"tags": schema.SetNestedAttribute{
							Description: "The list of tags that the dashboard administrator has privileges on",
							Optional:    true,
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"tag": schema.StringAttribute{
										MarkdownDescription: "",
										Optional:            true,
									},
									"access": schema.StringAttribute{
										MarkdownDescription: "",
										Optional:            true,
									},
								},
							},
						},
						"networks": schema.SetNestedAttribute{
							Description: "The list of networks that the dashboard administrator has privileges on",
							Optional:    true,
							Computed:    true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"id": schema.StringAttribute{
										MarkdownDescription: "",
										Optional:            true,
									},
									"access": schema.StringAttribute{
										MarkdownDescription: "",
										Optional:            true,
									},
								},
							},
						},
						"authentication_method": schema.StringAttribute{
							MarkdownDescription: "The method of authentication the user will use to sign in to the Meraki dashboard. Can be one of 'Email' or 'Cisco SecureX Sign-On'. The default is Email authentication",
							Optional:            true,
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (d *OrganizationsAdminsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *OrganizationsAdminsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OrganizationsAdminsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	inlineResp, httpResp, err := d.client.AdminsApi.GetOrganizationAdmins(context.Background(), data.OrgId.ValueString()).Execute()
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
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("State Data", fmt.Sprintf("\n%s", data))
		return
	}

	// Save data into Terraform state
	data.Id = types.StringValue("example-id")

	// admins attribute
	if admins := inlineResp; admins != nil {

		for _, inlineRespValue := range admins {
			var admin OrganizationAdminsDataSourceModel

			admin.Id = tools.MapStringValue(inlineRespValue, "id", &resp.Diagnostics)
			admin.Name = tools.MapStringValue(inlineRespValue, "name", &resp.Diagnostics)
			admin.Email = tools.MapStringValue(inlineRespValue, "email", &resp.Diagnostics)
			admin.OrgAccess = tools.MapStringValue(inlineRespValue, "orgAccess", &resp.Diagnostics)
			admin.AccountStatus = tools.MapStringValue(inlineRespValue, "accountStatus", &resp.Diagnostics)
			admin.TwoFactorAuthEnabled = tools.MapBoolValue(inlineRespValue, "twoFactorAuthEnabled", &resp.Diagnostics)
			admin.HasApiKey = tools.MapBoolValue(inlineRespValue, "hasApiKey", &resp.Diagnostics)
			admin.LastActive = tools.MapStringValue(inlineRespValue, "lastActive", &resp.Diagnostics)
			admin.AuthenticationMethod = tools.MapStringValue(inlineRespValue, "authenticationMethod", &resp.Diagnostics)

			// TODO - use tools.Map funcs for nested tags & networks data
			// tags attribute
			if tags := inlineRespValue["tags"]; tags != nil {
				for _, tv := range tags.([]interface{}) {
					var tag OrganizationAdminsDataSourceModelTag
					_ = json.Unmarshal([]byte(tv.(string)), &tag)
					admin.Tags = append(admin.Tags, tag)
				}
			} else {
				admin.Tags = nil
			}

			// networks attribute
			if networks := inlineRespValue["networks"]; networks != nil {
				for _, tv := range networks.([]interface{}) {
					var network OrganizationAdminsDataSourceModelNetwork
					_ = json.Unmarshal([]byte(tv.(string)), &network)
					admin.Networks = append(admin.Networks, network)
				}
			} else {
				admin.Networks = nil
			}

			// append admin to list of admins
			data.List = append(data.List, admin)
		}
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read data source")
}
