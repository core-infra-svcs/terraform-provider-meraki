package provider

import (
	"context"
	"fmt"

	apiclient "github.com/core-infra-svcs/dashboard-api-go/client"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &AdministeredIdentitiesMeDataSource{}

func NewAdministeredIdentitiesMeDataSource() datasource.DataSource {
	return &AdministeredIdentitiesMeDataSource{}
}

// AdministeredIdentitiesMeDataSource defines the data source implementation.
type AdministeredIdentitiesMeDataSource struct {
	client *apiclient.APIClient
}

// AdministeredIdentitiesMeDataSourceModel describes the data source data model.
type AdministeredIdentitiesMeDataSourceModel struct {
	Id                          types.String `tfsdk:"id"`
	AuthenticationApiKeyCreated types.Bool   `tfsdk:"authentication_api_key_created"`
	AuthenticationMode          types.String `tfsdk:"authentication_mode"`
	AuthenticationSaml          types.Bool   `tfsdk:"authentication_saml"`
	AuthenticationTwofactor     types.Bool   `tfsdk:"authentication_two_factor"`
	Email                       types.String `tfsdk:"email"`
	LastUsedDashboardAt         types.String `tfsdk:"last_used_dashboard_at"`
	Name                        types.String `tfsdk:"name"`
}

func (d *AdministeredIdentitiesMeDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_administered_identities_me"
}

func (d *AdministeredIdentitiesMeDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		MarkdownDescription: "AdministeredIdentitiesMe data source - Returns the identity of the current user",
		Attributes: map[string]tfsdk.Attribute{

			"id": {
				MarkdownDescription: "Example identifier",
				Type:                types.StringType,
				Computed:            true,
			},

			"authentication_api_key_created": {
				Description:         "API authentication Key",
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
			"authentication_mode": {
				Description:         "Authentication mode",
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
			"authentication_saml": {
				Description:         "SAML authentication",
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
			"authentication_two_factor": {
				Description:         "TwoFactor authentication",
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
			"email": {
				Description:         "User email",
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
			"last_used_dashboard_at": {
				Description:         "Last seen active on Dashboard UI",
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
				Description:         "Username",
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
		},
	}, nil
}

func (d *AdministeredIdentitiesMeDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *AdministeredIdentitiesMeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data AdministeredIdentitiesMeDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	response, r, err := d.client.AdministeredApi.GetAdministeredIdentitiesMe(context.Background()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error when calling read: %v\n", r),
			"Could not complete read request: "+err.Error(),
		)
		return
	}

	data.Id = types.String{Value: "example-id"}
	data.Name = types.String{Value: response.GetName()}
	data.Email = types.String{Value: response.GetEmail()}
	data.LastUsedDashboardAt = types.String{Value: response.GetLastUsedDashboardAt().String()}
	data.AuthenticationMode = types.String{Value: response.Authentication.GetMode()}
	data.AuthenticationApiKeyCreated = types.Bool{Value: response.Authentication.Api.Key.GetCreated()}
	data.AuthenticationSaml = types.Bool{Value: response.Authentication.Saml.GetEnabled()}
	data.AuthenticationTwofactor = types.Bool{Value: response.Authentication.TwoFactor.GetEnabled()}

	// Write logs using the tflog package
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
