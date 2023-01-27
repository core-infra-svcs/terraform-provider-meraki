package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	openApiClient "github.com/core-infra-svcs/dashboard-api-go/client"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
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
	Id                   jsontypes.String                         `tfsdk:"id"`
	OrgId                jsontypes.String                         `tfsdk:"organization_id" json:"organizationId"`
	AdminId              jsontypes.String                         `tfsdk:"admin_id" json:"id"`
	Name                 jsontypes.String                         `tfsdk:"name"`
	Email                jsontypes.String                         `tfsdk:"email"`
	OrgAccess            jsontypes.String                         `tfsdk:"org_access" json:"orgAccess"`
	AccountStatus        jsontypes.String                         `tfsdk:"account_status" json:"accountStatus"`
	TwoFactorAuthEnabled jsontypes.Bool                           `tfsdk:"two_factor_auth_enabled" json:"twoFactorAuthEnabled"`
	HasApiKey            jsontypes.Bool                           `tfsdk:"has_api_key" json:"hasApiKey"`
	LastActive           jsontypes.String                         `tfsdk:"last_active" json:"lastActive"`
	Tags                 []OrganizationsAdminResourceModelTag     `tfsdk:"tags" json:"tags"`
	Networks             []OrganizationsAdminResourceModelNetwork `tfsdk:"networks" json:"networks"`
	AuthenticationMethod jsontypes.String                         `tfsdk:"authentication_method" json:"authenticationMethod"`
}

type OrganizationsAdminResourceModelTag struct {
	Tag    jsontypes.String `tfsdk:"tag" json:"tag"`
	Access jsontypes.String `tfsdk:"access" json:"access"`
}

type OrganizationsAdminResourceModelNetwork struct {
	Id     jsontypes.String `tfsdk:"id" json:"id"`
	Access jsontypes.String `tfsdk:"access" json:"access"`
}

func (r *OrganizationsAdminResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizations_admin"
}

func (r *OrganizationsAdminResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage the dashboard administrators in this organization",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				CustomType: jsontypes.StringType,
				Computed:   true,
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "Organization ID",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
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
				CustomType:          jsontypes.StringType,
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
				CustomType:          jsontypes.StringType,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "The email of the dashboard administrator. This attribute can not be updated.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"org_access": schema.StringAttribute{
				MarkdownDescription: "The privilege of the dashboard administrator on the organization. Can be one of 'full', 'read-only', 'enterprise' or 'none'",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"full", "read-only", "enterprise", "none"}...),
					stringvalidator.LengthAtLeast(4),
				},
			},
			"account_status": schema.StringAttribute{
				MarkdownDescription: "",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"two_factor_auth_enabled": schema.BoolAttribute{
				MarkdownDescription: "",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.BoolType,
			},
			"has_api_key": schema.BoolAttribute{
				MarkdownDescription: "",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.BoolType,
			},
			"last_active": schema.StringAttribute{
				MarkdownDescription: "",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
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
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"access": schema.StringAttribute{
							MarkdownDescription: "",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
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
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"access": schema.StringAttribute{
							MarkdownDescription: "",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
					},
				},
			},
			"authentication_method": schema.StringAttribute{
				MarkdownDescription: "The method of authentication the user will use to sign in to the Meraki dashboard. Can be one of 'Email' or 'Cisco SecureX Sign-On'. The default is Email authentication",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
				Validators: []validator.String{

					stringvalidator.OneOf([]string{"Email", "Cisco SecureX Sign-On"}...),
					stringvalidator.LengthAtLeast(5),
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
	if len(data.Tags) > 0 {
		var tags []openApiClient.OrganizationsOrganizationIdAdminsTags
		for _, attribute := range data.Tags {
			var tag openApiClient.OrganizationsOrganizationIdAdminsTags
			tag.Tag = attribute.Tag.ValueString()
			tag.Access = attribute.Access.ValueString()
			tags = append(tags, tag)
		}
		createOrganizationAdmin.SetTags(tags)
	}

	// Networks
	if len(data.Networks) > 0 {
		var networks []openApiClient.OrganizationsOrganizationIdAdminsNetworks
		for _, attribute := range data.Networks {
			var network openApiClient.OrganizationsOrganizationIdAdminsNetworks
			network.Id = attribute.Id.ValueString()
			network.Access = attribute.Access.ValueString()
			networks = append(networks, network)
		}
		createOrganizationAdmin.SetNetworks(networks)
	}

	if !data.AuthenticationMethod.IsNull() {
		createOrganizationAdmin.SetAuthenticationMethod(data.AuthenticationMethod.ValueString())
	}

	_, httpResp, err := r.client.AdminsApi.CreateOrganizationAdmin(context.Background(), data.OrgId.ValueString()).CreateOrganizationAdmin(createOrganizationAdmin).Execute()
	// collect diagnostics
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	// Check for API success response code
	if httpResp.StatusCode != 201 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
		return
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%s", data))
		return
	}

	// Save data into Terraform state
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *OrganizationsAdminResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *OrganizationsAdminResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := r.client.AdminsApi.GetOrganizationAdmins(context.Background(), data.OrgId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
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
		return
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	}

	var admins []OrganizationsAdminResourceModel
	if err = json.NewDecoder(httpResp.Body).Decode(&admins); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	// There is no single GET ADMIN endpoint, so we must GET a list of all admins and search by adminId.
	for _, admin := range admins {
		// Match id found in tf state
		if admin.AdminId.ValueString() == data.AdminId.ValueString() {
			admin.Id = jsontypes.StringValue("example-id")
			admin.OrgId = data.OrgId // this value does not exist in the api response

			// Save updated data into Terraform state
			resp.Diagnostics.Append(resp.State.Set(ctx, &admin)...)
		}
	}
}

func (r *OrganizationsAdminResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *OrganizationsAdminResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Creating and Validating Payload for Creating Administrator
	updateOrganizationAdmin := *openApiClient.NewInlineObject177()
	updateOrganizationAdmin.SetName(data.Name.ValueString())
	updateOrganizationAdmin.SetOrgAccess(data.OrgAccess.ValueString())

	// Tags
	if len(data.Tags) > 0 {
		var tags []openApiClient.OrganizationsOrganizationIdAdminsTags
		for _, attribute := range data.Tags {
			var tag openApiClient.OrganizationsOrganizationIdAdminsTags
			tag.Tag = attribute.Tag.ValueString()
			tag.Access = attribute.Access.ValueString()
			tags = append(tags, tag)
		}
		updateOrganizationAdmin.SetTags(tags)
	}

	// Networks
	if len(data.Networks) > 0 {
		var networks []openApiClient.OrganizationsOrganizationIdAdminsNetworks
		for _, attribute := range data.Networks {
			var network openApiClient.OrganizationsOrganizationIdAdminsNetworks
			network.Id = attribute.Id.ValueString()
			network.Access = attribute.Access.ValueString()
			networks = append(networks, network)
		}
		updateOrganizationAdmin.SetNetworks(networks)
	}

	_, httpResp, err := r.client.AdminsApi.UpdateOrganizationAdmin(context.Background(), data.OrgId.ValueString(), data.AdminId.ValueString()).UpdateOrganizationAdmin(updateOrganizationAdmin).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
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
		return
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%s", data))
		return
	}

	// Save data into Terraform state
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	data.Id = jsontypes.StringValue("example-id")

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

	httpResp, err := r.client.AdminsApi.DeleteOrganizationAdmin(context.Background(), data.OrgId.ValueString(), data.AdminId.ValueString()).Execute()
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to delete resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
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
	}

	resp.State.RemoveResource(ctx)

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
