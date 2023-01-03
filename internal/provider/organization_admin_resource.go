package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"strings"

	openApiClient "github.com/core-infra-svcs/dashboard-api-go/client"

	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &OrganizationsAdminResource{}
	_ resource.ResourceWithConfigure   = &OrganizationsAdminResource{}
	_ resource.ResourceWithImportState = &OrganizationsAdminResource{}
)

func NewOrganizationsAdminResource() resource.Resource {
	return &OrganizationsAdminResource{}
}

// OrganizationsAdminResource defines the resource implementation.
type OrganizationsAdminResource struct {
	client *openApiClient.APIClient
}

// OrganizationsAdminResourceModel describes the resource data model.
type OrganizationsAdminResourceModel struct {
	Id                   types.String                             `tfsdk:"id"`
	OrgId                types.String                             `tfsdk:"organization_id"`
	AdminId              types.String                             `tfsdk:"admin_id"`
	Name                 types.String                             `tfsdk:"name"`
	Email                types.String                             `tfsdk:"email"`
	OrgAccess            types.String                             `tfsdk:"org_access"`
	AccountStatus        types.String                             `tfsdk:"account_status"`
	TwoFactorAuthEnabled types.Bool                               `tfsdk:"two_factor_auth_enabled"`
	HasApiKey            types.Bool                               `tfsdk:"has_api_key"`
	LastActive           types.String                             `tfsdk:"last_active"`
	Tags                 []OrganizationsAdminResourceModelTag     `tfsdk:"tags"`
	Networks             []OrganizationsAdminResourceModelNetwork `tfsdk:"networks"`
	AuthenticationMethod types.String                             `tfsdk:"authentication_method"`
}

type OrganizationsAdminResourceModelNetwork struct {
	Id     types.String `tfsdk:"id"`
	Access types.String `tfsdk:"access"`
}

type OrganizationsAdminResourceModelTag struct {
	Tag    types.String `tfsdk:"tag"`
	Access types.String `tfsdk:"access"`
}

func (r *OrganizationsAdminResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizations_admin"
}

func (r *OrganizationsAdminResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage the dashboard administrators in this organization",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Example identifier",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "Organization ID",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 31),
				},
			},
			"admin_id": schema.StringAttribute{
				MarkdownDescription: "Admin ID",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 31),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the dashboard administrator",
				Optional:            true,
				Computed:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "The email of the dashboard administrator. This attribute can not be updated.",
				Optional:            true,
				Computed:            true,
			},
			"org_access": schema.StringAttribute{
				MarkdownDescription: "The privilege of the dashboard administrator on the organization. Can be one of 'full', 'read-only', 'enterprise' or 'none'",
				Optional:            true,
				Computed:            true,
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
				Computed:            true,
			},
			"two_factor_auth_enabled": schema.BoolAttribute{
				MarkdownDescription: "",
				Optional:            true,
				Computed:            true,
			},
			"has_api_key": schema.BoolAttribute{
				MarkdownDescription: "",
				Optional:            true,
				Computed:            true,
			},
			"last_active": schema.StringAttribute{
				MarkdownDescription: "",
				Optional:            true,
				Computed:            true,
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
	}
}

func (r *OrganizationsAdminResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*openApiClient.APIClient)

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

	// Creating and Validating Payload for Creating Administrator
	createOrganizationAdmin := *openApiClient.NewInlineObject176(
		data.Email.ValueString(),
		data.Name.ValueString(),
		data.OrgAccess.ValueString())

	// Tags
	if len(data.Tags) < 0 {
		var tags []openApiClient.OrganizationsOrganizationIdAdminsTags
		for _, attribute := range data.Tags {
			var tag openApiClient.OrganizationsOrganizationIdAdminsTags
			tag.Tag = attribute.Tag.ValueString()
			tag.Access = attribute.Tag.ValueString()
			tags = append(tags, tag)
		}
		createOrganizationAdmin.SetTags(tags)
	}

	// Networks
	if len(data.Networks) < 0 {
		var networks []openApiClient.OrganizationsOrganizationIdAdminsNetworks
		for _, attribute := range data.Networks {
			var network openApiClient.OrganizationsOrganizationIdAdminsNetworks
			network.Id = attribute.Id.ValueString()
			network.Access = attribute.Access.ValueString()
			networks = append(networks, network)
		}
		createOrganizationAdmin.SetNetworks(networks)
	}

	if data.AuthenticationMethod.IsNull() != true {
		createOrganizationAdmin.SetAuthenticationMethod(data.AuthenticationMethod.ValueString())
	}

	inlineResp, httpResp, err := r.client.AdminsApi.CreateOrganizationAdmin(context.Background(), data.OrgId.ValueString()).CreateOrganizationAdmin(createOrganizationAdmin).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// collect diagnostics
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
	}

	// Check for API success response code
	if httpResp.StatusCode != 201 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%s", data))
		return
	}

	// Save data into Terraform state
	extractHttpResponseOrganizationAdminResource(ctx, inlineResp, data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")
}

func (r *OrganizationsAdminResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *OrganizationsAdminResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	inlineResp, httpResp, err := r.client.AdminsApi.GetOrganizationAdmins(context.Background(), data.OrgId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// collect diagnostics
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
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
		return
	} else {
		resp.Diagnostics.Append()
	}

	// get admin
	for _, admin := range inlineResp {

		// Match id found in tf state
		if adminId := admin["id"]; adminId == data.AdminId.ValueString() {

			// Save data into Terraform state
			extractHttpResponseOrganizationAdminResource(ctx, admin, data)
		}
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OrganizationsAdminResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *OrganizationsAdminResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if len(data.AdminId.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing AdminId", fmt.Sprintf("AdminId: %s", data.AdminId.ValueString()))
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Creating and Validating Payload for Creating Administrator
	updateOrganizationAdmin := *openApiClient.NewInlineObject177()
	updateOrganizationAdmin.SetName(data.Name.ValueString())
	updateOrganizationAdmin.SetOrgAccess(data.OrgAccess.ValueString())

	// Tags
	if len(data.Tags) < 0 {
		var tags []openApiClient.OrganizationsOrganizationIdAdminsTags
		for _, attribute := range data.Tags {
			var tag openApiClient.OrganizationsOrganizationIdAdminsTags
			tag.Tag = attribute.Tag.ValueString()
			tag.Access = attribute.Tag.ValueString()
			tags = append(tags, tag)
		}
		updateOrganizationAdmin.SetTags(tags)
	}

	// Networks
	if len(data.Networks) < 0 {
		var networks []openApiClient.OrganizationsOrganizationIdAdminsNetworks
		for _, attribute := range data.Networks {
			var network openApiClient.OrganizationsOrganizationIdAdminsNetworks
			network.Id = attribute.Id.ValueString()
			network.Access = attribute.Access.ValueString()
			networks = append(networks, network)
		}
		updateOrganizationAdmin.SetNetworks(networks)
	}

	inlineResp, httpResp, err := r.client.AdminsApi.UpdateOrganizationAdmin(context.Background(), data.OrgId.ValueString(), data.AdminId.ValueString()).UpdateOrganizationAdmin(updateOrganizationAdmin).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// collect diagnostics
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
	}

	// Check for API success response code
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%s", data))
		return
	}

	// Save data into Terraform state
	extractHttpResponseOrganizationAdminResource(ctx, inlineResp, data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OrganizationsAdminResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *OrganizationsAdminResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.client.AdminsApi.DeleteOrganizationAdmin(context.Background(), data.OrgId.ValueString(), data.AdminId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to delete resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// collect diagnostics
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
	}

	// Check for API success response code
	if httpResp.StatusCode != 204 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	} else {
		resp.Diagnostics.Append()
	}

	resp.State.RemoveResource(ctx)

}

func extractHttpResponseOrganizationAdminResource(ctx context.Context, inlineRespValue map[string]interface{}, data *OrganizationsAdminResourceModel) *OrganizationsAdminResourceModel {

	// save into the Terraform state
	data.Id = types.StringValue("example-id")

	// id attribute
	if id := inlineRespValue["id"]; id != nil {
		data.AdminId = types.StringValue(id.(string))
	} else {
		data.AdminId = types.StringNull()
	}

	// name attribute
	if name := inlineRespValue["name"]; name != nil {
		data.Name = types.StringValue(name.(string))
	} else {
		data.Name = types.StringNull()
	}

	// email attribute
	if email := inlineRespValue["email"]; email != nil {
		data.Email = types.StringValue(email.(string))
	} else {
		data.Email = types.StringNull()
	}

	// orgAccess attribute
	if orgAccess := inlineRespValue["orgAccess"]; orgAccess != nil {
		data.OrgAccess = types.StringValue(orgAccess.(string))
	} else {
		data.OrgAccess = types.StringNull()
	}

	// accountStatus attribute
	if accountStatus := inlineRespValue["accountStatus"]; accountStatus != nil {
		data.AccountStatus = types.StringValue(accountStatus.(string))
	} else {
		data.AccountStatus = types.StringNull()
	}

	// twoFactorAuthEnabled attribute
	if twoFactorAuthEnabled := inlineRespValue["twoFactorAuthEnabled"]; twoFactorAuthEnabled != nil {
		data.TwoFactorAuthEnabled = types.BoolValue(twoFactorAuthEnabled.(bool))
	} else {
		data.TwoFactorAuthEnabled = types.BoolNull()
	}

	// hasApiKey attribute
	if hasApiKey := inlineRespValue["hasApiKey"]; hasApiKey != nil {
		data.HasApiKey = types.BoolValue(hasApiKey.(bool))
	} else {
		data.HasApiKey = types.BoolNull()
	}

	// lastActive attribute
	if lastActive := inlineRespValue["lastActive"]; lastActive != nil {
		data.LastActive = types.StringValue(lastActive.(string))
	} else {
		data.LastActive = types.StringNull()
	}

	// tags attribute
	if tags := inlineRespValue["tags"]; tags != nil {
		for _, tv := range tags.([]interface{}) {
			var tag OrganizationsAdminResourceModelTag
			_ = json.Unmarshal([]byte(tv.(string)), &tag)
			data.Tags = append(data.Tags, tag)
		}
	} else {
		data.Tags = nil
	}

	// networks attribute
	if networks := inlineRespValue["networks"]; networks != nil {
		for _, tv := range networks.([]interface{}) {
			var network OrganizationsAdminResourceModelNetwork
			_ = json.Unmarshal([]byte(tv.(string)), &network)
			data.Networks = append(data.Networks, network)
		}
	} else {
		data.Networks = nil
	}

	// authenticationMethod attribute
	if authenticationMethod := inlineRespValue["authenticationMethod"]; authenticationMethod != nil {
		data.AuthenticationMethod = types.StringValue(authenticationMethod.(string))
	} else {
		data.AuthenticationMethod = types.StringNull()
	}

	return data
}

func (r *OrganizationsAdminResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: organization_id, admin_id. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("admin_id"), idParts[1])...)
	if resp.Diagnostics.HasError() {
		return
	}
}
