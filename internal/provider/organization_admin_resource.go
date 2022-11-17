package provider

import (
	"context"
	"encoding/json"
	"fmt"

	apiclient "github.com/core-infra-svcs/dashboard-api-go/client"

	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &OrganizationsAdminResource{}
var _ resource.ResourceWithImportState = &OrganizationsAdminResource{}

func NewOrganizationsAdminResource() resource.Resource {
	return &OrganizationsAdminResource{}
}

// OrganizationsAdminResource defines the resource implementation.
type OrganizationsAdminResource struct {
	client *apiclient.APIClient
}

// OrganizationsAdminResourceModel describes the resource data model.
type OrganizationsAdminResourceModel struct {
	Id                   types.String           `tfsdk:"id"`
	AdminId              types.String           `tfsdk:"adminid"`
	Name                 types.String           `tfsdk:"name"`
	Email                types.String           `tfsdk:"email"`
	OrgAccess            types.String           `tfsdk:"orgaccess"`
	AuthenticationMethod types.String           `tfsdk:"authentication_method"`
	AccountStatus        types.String           `tfsdk:"account_status"`
	TwoFactorAuthEnabled types.Bool             `tfsdk:"two_factor_auth_enabled"`
	HasApiKey            types.Bool             `tfsdk:"has_api_key"`
	LastActive           types.String           `tfsdk:"last_active"`
	Tags                 []AdminResourceTag     `tfsdk:"tags"`
	Networks             []AdminResourceNetwork `tfsdk:"networks"`
}

// AdminResourceTag  describes the tag data model
type AdminResourceTag struct {
	Tag    string `tfsdk:"tag"`
	Access string `tfsdk:"access"`
}

// AdminResourceNetwork  describes the network data model
type AdminResourceNetwork struct {
	Id     string `tfsdk:"id"`
	Access string `tfsdk:"access"`
}

// AdminResourceInfo  describes the resource data model
type AdminResourceInfo struct {
	Name                 string
	Email                string
	Id                   string
	AdminId              string
	OrgAccess            string
	AuthenticationMethod string
	Tags                 []AdminResourceTag
	Networks             []AdminResourceNetwork
	AccountStatus        string
	TwoFactorAuthEnabled bool
	HasApiKey            bool
	LastActive           string
}

func (r *OrganizationsAdminResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizations_admin"
}

func (r *OrganizationsAdminResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "Organization Admin resource - Manage the admins for an organization",
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Description:         "meraki organization Id",
				MarkdownDescription: "",
				Type:                types.StringType,
				Required:            true,
			},
			"adminid": {
				Description:         "id of dashboard administrator",
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
				Required:            true,
				Optional:            false,
				Computed:            false,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"orgaccess": {
				Description:         "Organization Access",
				MarkdownDescription: "",
				Type:                types.StringType,
				Required:            true,
				Optional:            false,
				Computed:            false,
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
				Required:            true,
				Optional:            false,
				Computed:            false,
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
				Optional:    true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"tag": {
						Description:         "administrator tag",
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
						Description:         "administrator tag access",
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
				Optional:    true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"id": {
						Description:         "administrator network id ",
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
						Description:         "administrator network access",
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
		},
	}, nil
}

func (r *OrganizationsAdminResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*apiclient.APIClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	r.client = client
}

func (r *OrganizationsAdminResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *OrganizationsAdminResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Creating and Validating Payload for Creating Database Adminstrator
	createOrganizationAdmin := *apiclient.NewInlineObject176(data.Email.ValueString(), data.Name.ValueString(), data.OrgAccess.ValueString())

	if data.Tags != nil {
		if len(data.Tags) != 0 {
			var t []apiclient.OrganizationsOrganizationIdAdminsTags
			for _, tag := range data.Tags {
				var tagData apiclient.OrganizationsOrganizationIdAdminsTags
				tagData.Tag = tag.Tag
				tagData.Access = tag.Access
				t = append(t, tagData)
			}
			createOrganizationAdmin.SetTags(t)
		} else {
			resp.Diagnostics.AddError("tags should not be empty. Add atleast one tag and access fields or else remove tags field ", fmt.Sprintf("tags: %s", data.Tags))
			return
		}
	}

	if data.Networks != nil {
		if len(data.Networks) != 0 {
			var n []apiclient.OrganizationsOrganizationIdAdminsNetworks
			for _, network := range data.Networks {
				var networkData apiclient.OrganizationsOrganizationIdAdminsNetworks
				networkData.Id = network.Id
				networkData.Access = network.Access
				n = append(n, networkData)
			}
			createOrganizationAdmin.SetNetworks(n)
		} else {
			resp.Diagnostics.AddError("networks should not be empty. Add atleast one id and access fields or else remove networks field ", fmt.Sprintf("networks: %s", data.Networks))
			return
		}

	}

	if !data.AuthenticationMethod.IsUnknown() {
		createOrganizationAdmin.SetAuthenticationMethod(data.AuthenticationMethod.ValueString())
	}

	inlineResp, httpResp, err := r.client.AdminsApi.CreateOrganizationAdmin(context.Background(), data.Id.ValueString()).CreateOrganizationAdmin(createOrganizationAdmin).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// Check for API success response code
	if httpResp.StatusCode != 201 {
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

	resp.Diagnostics.Append()

	var admindata AdminResourceInfo

	// Convert map to json string
	jsonStr, err := json.Marshal(inlineResp)
	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct
	if err := json.Unmarshal(jsonStr, &admindata); err != nil {
		fmt.Println(err)
	}

	data.Name = types.StringValue(admindata.Name)
	data.Email = types.StringValue(admindata.Email)
	data.AdminId = types.StringValue(admindata.Id)
	data.OrgAccess = types.StringValue(admindata.OrgAccess)
	data.AuthenticationMethod = types.StringValue(admindata.AuthenticationMethod)
	data.AccountStatus = types.StringValue(admindata.AccountStatus)
	data.TwoFactorAuthEnabled = types.BoolValue(admindata.TwoFactorAuthEnabled)
	data.HasApiKey = types.BoolValue(admindata.HasApiKey)
	data.LastActive = types.StringValue(admindata.LastActive)
	if data.Tags != nil {
		data.Tags = admindata.Tags
	}
	if data.Networks != nil {
		data.Networks = admindata.Networks
	}

	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OrganizationsAdminResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *OrganizationsAdminResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	inlineResp, httpResp, err := r.client.AdminsApi.GetOrganizationAdmins(context.Background(), data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read resource",
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

	var adminresource []AdminResourceInfo

	// Convert map to json string
	jsonStr, err := json.Marshal(inlineResp)
	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct
	if err := json.Unmarshal(jsonStr, &adminresource); err != nil {
		fmt.Println(err)
	}

	for _, admindata := range adminresource {

		if admindata.Email == data.Email.ValueString() {

			data.Name = types.StringValue(admindata.Name)
			data.Email = types.StringValue(admindata.Email)
			data.AdminId = types.StringValue(admindata.Id)
			data.OrgAccess = types.StringValue(admindata.OrgAccess)
			data.AuthenticationMethod = types.StringValue(admindata.AuthenticationMethod)
			data.AccountStatus = types.StringValue(admindata.AccountStatus)
			data.TwoFactorAuthEnabled = types.BoolValue(admindata.TwoFactorAuthEnabled)
			data.HasApiKey = types.BoolValue(admindata.HasApiKey)
			data.LastActive = types.StringValue(admindata.LastActive)
			if data.Tags != nil {
				data.Tags = admindata.Tags
			}
			if data.Networks != nil {
				data.Networks = admindata.Networks
			}

		}

	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OrganizationsAdminResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *OrganizationsAdminResourceModel
	var state *OrganizationsAdminResourceModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}
	// Creating and Validating Payload for Updating Database Adminstrator
	updateOrganizationAdmin := *apiclient.NewInlineObject177()
	if data.Tags != nil {
		if len(data.Tags) != 0 {
			var t []apiclient.OrganizationsOrganizationIdAdminsTags
			for _, tag := range data.Tags {
				var tagData apiclient.OrganizationsOrganizationIdAdminsTags
				tagData.Tag = tag.Tag
				tagData.Access = tag.Access
				t = append(t, tagData)

			}
			updateOrganizationAdmin.SetTags(t)
		} else {
			resp.Diagnostics.AddError("tags should not be empty. Add atleast one tag and access fields or else remove tags field ", fmt.Sprintf("tags: %s", data.Tags))
			return
		}
	}
	if data.Networks != nil {
		if len(data.Networks) != 0 {
			var n []apiclient.OrganizationsOrganizationIdAdminsNetworks
			for _, network := range data.Networks {
				var networkData apiclient.OrganizationsOrganizationIdAdminsNetworks
				networkData.Id = network.Id
				networkData.Access = network.Access
				n = append(n, networkData)
			}
			updateOrganizationAdmin.SetNetworks(n)
		} else {
			resp.Diagnostics.AddError("networks should not be empty. Add atleast one id and access fields or else remove networks field ", fmt.Sprintf("networks: %s", data.Networks))
			return
		}

	}
	updateOrganizationAdmin.SetName(data.Name.ValueString())
	updateOrganizationAdmin.SetOrgAccess(data.OrgAccess.ValueString())
	inlineResp, httpResp, err := r.client.AdminsApi.UpdateOrganizationAdmin(context.Background(), data.Id.ValueString(), state.AdminId.ValueString()).UpdateOrganizationAdmin(updateOrganizationAdmin).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update resource",
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
		return
	}

	var admindata AdminResourceInfo

	// Convert map to json string
	jsonStr, err := json.Marshal(inlineResp)
	if err != nil {
		fmt.Println(err)
	}
	// Convert json string to struct
	if err := json.Unmarshal(jsonStr, &admindata); err != nil {
		fmt.Println(err)
	}

	data.Name = types.StringValue(admindata.Name)
	data.Email = types.StringValue(admindata.Email)
	data.AdminId = types.StringValue(admindata.Id)
	data.OrgAccess = types.StringValue(admindata.OrgAccess)
	data.AuthenticationMethod = types.StringValue(admindata.AuthenticationMethod)
	data.AccountStatus = types.StringValue(admindata.AccountStatus)
	data.TwoFactorAuthEnabled = types.BoolValue(admindata.TwoFactorAuthEnabled)
	data.HasApiKey = types.BoolValue(admindata.HasApiKey)
	data.LastActive = types.StringValue(admindata.LastActive)
	if data.Tags != nil {
		data.Tags = admindata.Tags
	}
	if data.Networks != nil {
		data.Networks = admindata.Networks
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OrganizationsAdminResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *OrganizationsAdminResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.client.AdminsApi.DeleteOrganizationAdmin(context.Background(), data.Id.ValueString(), data.AdminId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to delete resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// Check for API success response code
	if httpResp.StatusCode != 204 {
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

	resp.State.RemoveResource(ctx)

}

func (r *OrganizationsAdminResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
