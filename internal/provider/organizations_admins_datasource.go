package provider

import (
	"context"
	"encoding/json"
	"fmt"

	apiclient "github.com/core-infra-svcs/dashboard-api-go/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
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
	Id   types.String                        `tfsdk:"id"`
	List []OrganizationAdminsDataSourceModel `tfsdk:"list"`
}

// OrganizationAdminsDataSourceModel describes the data source data model.
type OrganizationAdminsDataSourceModel struct {
	Email                types.String `tfsdk:"email"`
	Name                 types.String `tfsdk:"name"`
	Id                   types.String `tfsdk:"id"`
	AuthenticationMethod types.String `tfsdk:"authentication_method"`
	OrgAccess            types.String `tfsdk:"orgaccess"`
	AccountStatus        types.String `tfsdk:"account_status"`
	TwoFactorAuthEnabled types.Bool   `tfsdk:"two_factor_auth_enabled"`
	HasApiKey            types.Bool   `tfsdk:"has_api_key"`
	LastActive           types.String `tfsdk:"last_active"`
	Networks             []Network    `tfsdk:"networks"`
	Tags                 []Tag        `tfsdk:"tags"`
}

type Network struct {
	Id     types.String `tfsdk:"id"`
	Access types.String `tfsdk:"access"`
}

type Tag struct {
	Tag    types.String `tfsdk:"tag"`
	Access types.String `tfsdk:"access"`
}

type AdminData struct {
	Name                 string
	Email                string
	Id                   string
	OrgAccess            string
	AuthenticationMethod string
	AccountStatus        string
	TwoFactorAuthEnabled bool
	HasApiKey            bool
	LastActive           string
	Networks             []Network
	Tags                 []Tag
}

func (d *OrganizationsAdminsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizations_admins"
}

func (d *OrganizationsAdminsDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "OrganizationsAdmins data source - get all list of  admins in an organization",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				MarkdownDescription: "org id",
				Description:         "Organization Id",
				Type:                types.StringType,
				Required:            true,
			},
			"list": {
				MarkdownDescription: "List of organization admins",
				Optional:            true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"orgaccess": {
						Description:         "Organization Access",
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
					"authentication_method": {
						Description:         "Authentication method",
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
					"email": {
						Description:         "Email of the dashboard administrator",
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
						Description:         "name of the dashboard administrator",
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
						Description:         "id of the organization",
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
					"two_factor_auth_enabled": {
						Description:         "Two Factor Auth Enabled or Not",
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
					"account_status": {
						Description:         "Account Status",
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
					"has_api_key": {
						Description:         "Api key exists or not",
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
					"last_active": {
						Description:         "Last Time Active",
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
					"tags": {
						Description: "list of tags that the dashboard administrator has privileges on.",
						Computed:    true,
						Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
							"tag": {
								Description:         "tag",
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
							"access": {
								Description:         "access",
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
						})},
					"networks": {
						Description: "list of networks that the dashboard administrator has privileges on.",
						Computed:    true,
						Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
							"id": {
								Description:         "network id ",
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
							"access": {
								Description:         "network access",
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
						})},
				}),
			},
		},
	}, nil
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

	if resp.Diagnostics.HasError() {
		return
	}
	response, r, err := d.client.AdminsApi.GetOrganizationAdmins(context.Background(), data.Id.Value).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error when calling read: %v\n", r),
			"Could not complete read request: "+err.Error(),
		)
		return
	}

	var admins []AdminData

	// Convert map to json string
	jsonStr, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct
	if err := json.Unmarshal(jsonStr, &admins); err != nil {
		fmt.Println(err)
	}

	for _, adminData := range admins {

		var result OrganizationAdminsDataSourceModel
		result.Name = types.String{Value: adminData.Name}
		result.Email = types.String{Value: adminData.Email}
		result.Id = types.String{Value: adminData.Id}
		result.OrgAccess = types.String{Value: adminData.OrgAccess}
		result.AuthenticationMethod = types.String{Value: adminData.AuthenticationMethod}
		result.LastActive = types.String{Value: adminData.LastActive}
		result.AccountStatus = types.String{Value: adminData.AccountStatus}
		result.TwoFactorAuthEnabled = types.Bool{Value: adminData.TwoFactorAuthEnabled}
		result.HasApiKey = types.Bool{Value: adminData.HasApiKey}
		for _, network := range adminData.Networks {
			var networkData Network
			networkData.Id = types.String{Value: network.Id.Value}
			networkData.Access = types.String{Value: network.Access.Value}
			result.Networks = append(result.Networks, networkData)
		}
		for _, tag := range adminData.Tags {
			var tagData Tag
			tagData.Tag = types.String{Value: tag.Tag.Value}
			tagData.Access = types.String{Value: tag.Access.Value}
			result.Tags = append(result.Tags, tagData)
		}
		data.List = append(data.List, result)
	}

	data.Id = types.String{Value: "example-id"}

	// Write logs using the tflog package
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
