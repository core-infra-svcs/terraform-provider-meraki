package provider

import (
	"context"
	"encoding/json"
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
var _ datasource.DataSource = &OrganizationsAdminsDataSource{}

func NewOrganizationsAdminsDataSource() datasource.DataSource {
	return &OrganizationsAdminsDataSource{}
}

// OrganizationsAdminsDataSource defines the data source implementation.
type OrganizationsAdminsDataSource struct {
	client *apiclient.APIClient
}

// OrganizationsAdminsDataSourceModel describes the data source data model.
type OrganizationsAdminsDataSourceModel struct {
	Id    types.String                        `tfsdk:"id"`
	OrgId types.String                        `tfsdk:"organization_id"`
	List  []OrganizationAdminsDataSourceModel `tfsdk:"list"`
}

// OrganizationAdminsDataSourceModel describes the data source data model.
type OrganizationAdminsDataSourceModel struct {
	Id                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	Email                types.String `tfsdk:"email"`
	OrgAccess            types.String `tfsdk:"org_access"`
	AccountStatus        types.String `tfsdk:"account_status"`
	TwoFactorAuthEnabled types.Bool   `tfsdk:"two_factor_auth_enabled"`
	HasApiKey            types.Bool   `tfsdk:"has_api_key"`
	LastActive           types.String `tfsdk:"last_active"`
	Tags                 []Tag        `tfsdk:"tags"`
	Networks             []Network    `tfsdk:"networks"`
	AuthenticationMethod types.String `tfsdk:"authentication_method"`
}

type Network struct {
	Id     types.String `tfsdk:"id"`
	Access types.String `tfsdk:"access"`
}

type Tag struct {
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
			},
			"list": schema.ListNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"admin_id": schema.StringAttribute{
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
							Validators: []validator.String{
								stringvalidator.ExactlyOneOf(
									path.MatchRoot("full"),
									path.MatchRoot("read-only"),
									path.MatchRoot("enterprise"),
									path.MatchRoot("none"),
								),
							},
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
						"tags": schema.ListNestedAttribute{
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
						"networks": schema.ListNestedAttribute{
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
							Validators: []validator.String{
								stringvalidator.ExactlyOneOf(
									path.MatchRoot("Email"),
									path.MatchRoot("Cisco SecureX Sign-On"),
								),
							},
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

func (d *OrganizationsAdminsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OrganizationsAdminsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	// Check for required parameters
	if len(data.OrgId.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing organizationId", fmt.Sprintf("Value: %s", data.OrgId.ValueString()))
		return
	}

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
	tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)

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

			// id attribute
			if id := inlineRespValue["id"]; id != nil {
				admin.Id = types.StringValue(id.(string))
			} else {
				admin.Id = types.StringNull()
			}

			// name attribute
			if name := inlineRespValue["name"]; name != nil {
				admin.Name = types.StringValue(name.(string))
			} else {
				admin.Name = types.StringNull()
			}

			// email attribute
			if email := inlineRespValue["email"]; email != nil {
				admin.Email = types.StringValue(email.(string))
			} else {
				admin.Email = types.StringNull()
			}

			// orgAccess attribute
			if orgAccess := inlineRespValue["orgAccess"]; orgAccess != nil {
				admin.OrgAccess = types.StringValue(orgAccess.(string))
			} else {
				admin.OrgAccess = types.StringNull()
			}

			// accountStatus attribute
			if accountStatus := inlineRespValue["accountStatus"]; accountStatus != nil {
				admin.AccountStatus = types.StringValue(accountStatus.(string))
			} else {
				admin.AccountStatus = types.StringNull()
			}

			// twoFactorAuthEnabled attribute
			if twoFactorAuthEnabled := inlineRespValue["twoFactorAuthEnabled"]; twoFactorAuthEnabled != nil {
				admin.TwoFactorAuthEnabled = types.BoolValue(twoFactorAuthEnabled.(bool))
			} else {
				admin.TwoFactorAuthEnabled = types.BoolNull()
			}

			// hasApiKey attribute
			if hasApiKey := inlineRespValue["hasApiKey"]; hasApiKey != nil {
				admin.HasApiKey = types.BoolValue(hasApiKey.(bool))
			} else {
				admin.HasApiKey = types.BoolNull()
			}

			// lastActive attribute
			if lastActive := inlineRespValue["lastActive"]; lastActive != nil {
				admin.LastActive = types.StringValue(lastActive.(string))
			} else {
				admin.LastActive = types.StringNull()
			}

			// tags attribute
			if tags := inlineRespValue["tags"]; tags != nil {
				for _, tv := range tags.([]interface{}) {
					var tag Tag
					_ = json.Unmarshal([]byte(tv.(string)), &tag)
					admin.Tags = append(admin.Tags, tag)
				}
			} else {
				admin.Tags = nil
			}

			// networks attribute
			if networks := inlineRespValue["networks"]; networks != nil {
				for _, tv := range networks.([]interface{}) {
					var network Network
					_ = json.Unmarshal([]byte(tv.(string)), &network)
					admin.Networks = append(admin.Networks, network)
				}
			} else {
				admin.Networks = nil
			}

			// authenticationMethod attribute
			if authenticationMethod := inlineRespValue["authenticationMethod"]; authenticationMethod != nil {
				admin.AuthenticationMethod = types.StringValue(authenticationMethod.(string))
			} else {
				admin.AuthenticationMethod = types.StringNull()
			}

			// append admin to list of admins
			data.List = append(data.List, admin)
		}

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read data source")
}
