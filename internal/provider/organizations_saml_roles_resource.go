package provider

import (
	"context"
	"fmt"
	"strings"

	openApiClient "github.com/core-infra-svcs/dashboard-api-go/client"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &OrganizationsSamlrolesResource{}
var _ resource.ResourceWithImportState = &OrganizationsSamlrolesResource{}

func NewOrganizationsSamlrolesResource() resource.Resource {
	return &OrganizationsSamlrolesResource{}
}

// OrganizationsSamlrolesResource defines the resource implementation.
type OrganizationsSamlrolesResource struct {
	client *openApiClient.APIClient
}

// OrganizationsSamlrolesResourceModel describes the resource data model.
type OrganizationsSamlrolesResourceModel struct {
	Id        types.String                                `tfsdk:"id"`
	OrgId     types.String                                `tfsdk:"organization_id"`
	RoleId    types.String                                `tfsdk:"role_id"`
	Role      types.String                                `tfsdk:"role"`
	OrgAccess types.String                                `tfsdk:"org_access"`
	Tags      []OrganizationsSamlRoleResourceModelTag     `tfsdk:"tags"`
	Networks  []OrganizationsSamlRoleResourceModelNetwork `tfsdk:"networks"`
}

type OrganizationsSamlRoleResourceModelNetwork struct {
	Id     types.String `tfsdk:"id"`
	Access types.String `tfsdk:"access"`
}

type OrganizationsSamlRoleResourceModelTag struct {
	Tag    types.String `tfsdk:"tag"`
	Access types.String `tfsdk:"access"`
}

func (r *OrganizationsSamlrolesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizations_saml_role"
}

func (r *OrganizationsSamlrolesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage the saml roles in this organization",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "Organization ID",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 31),
				},
			},
			"role_id": schema.StringAttribute{
				MarkdownDescription: "Saml Role ID",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"role": schema.StringAttribute{
				MarkdownDescription: "The role of the SAML administrator",
				Required:            true,
			},
			"org_access": schema.StringAttribute{
				MarkdownDescription: "The privilege of the SAML administrator on the organization. Can be one of 'none', 'read-only', 'full' or 'enterprise'",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"full", "read-only", "enterprise", "none"}...),
					stringvalidator.LengthAtLeast(4),
				},
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
						},
						"access": schema.StringAttribute{
							MarkdownDescription: "",
							Optional:            true,
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
						},
						"access": schema.StringAttribute{
							MarkdownDescription: "",
							Optional:            true,
						},
					},
				},
			},
		},
	}
}

func (r *OrganizationsSamlrolesResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OrganizationsSamlrolesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *OrganizationsSamlrolesResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if len(data.OrgId.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing OrganizationId", fmt.Sprintf("Value: %s", data.OrgId.ValueString()))
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	createOrganizationSamlRole := *openApiClient.NewInlineObject216(data.Role.ValueString(), data.OrgAccess.ValueString())

	// Tags

	if len(data.Tags) > 0 {
		var tags []openApiClient.OrganizationsOrganizationIdSamlRolesTags
		for _, attribute := range data.Tags {
			var tag openApiClient.OrganizationsOrganizationIdSamlRolesTags
			tag.Tag = attribute.Tag.ValueString()
			tag.Access = attribute.Access.ValueString()
			tags = append(tags, tag)
		}
		createOrganizationSamlRole.SetTags(tags)
	}

	// Networks
	if len(data.Networks) > 0 {
		var networks []openApiClient.OrganizationsOrganizationIdSamlRolesNetworks
		for _, attribute := range data.Networks {
			var network openApiClient.OrganizationsOrganizationIdSamlRolesNetworks
			network.Id = attribute.Id.ValueString()
			network.Access = attribute.Access.ValueString()
			networks = append(networks, network)
		}
		createOrganizationSamlRole.SetNetworks(networks)
	}

	inlineResp, httpResp, err := r.client.OrganizationsApi.CreateOrganizationSamlRole(context.Background(), data.OrgId.ValueString()).CreateOrganizationSamlRole(createOrganizationSamlRole).Execute()

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
			fmt.Sprintf("%v%v", httpResp.StatusCode, inlineResp),
		)
	}

	// collect diagnostics
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%s", data))
		return
	}

	extractHttpResponseOrganizationSamlRoleResource(ctx, inlineResp, data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")
}

func (r *OrganizationsSamlrolesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *OrganizationsSamlrolesResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if len(data.OrgId.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing OrganizationId", fmt.Sprintf("Value: %s", data.OrgId.ValueString()))
		return
	}

	if len(data.RoleId.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing RoleId", fmt.Sprintf("Value: %s", data.RoleId.ValueString()))
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	inlineResp, httpResp, err := r.client.OrganizationsApi.GetOrganizationSamlRole(ctx, data.OrgId.ValueString(), data.RoleId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

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
	}

	extractHttpResponseOrganizationSamlRoleResource(ctx, inlineResp, data)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *OrganizationsSamlrolesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data *OrganizationsSamlrolesResourceModel
	var stateData *OrganizationsSamlrolesResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)

	if len(data.OrgId.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing OrganizationId", fmt.Sprintf("Value: %s", data.OrgId.ValueString()))
		return
	}

	// Check state for required attribute
	if len(data.RoleId.ValueString()) < 1 {
		data.RoleId = stateData.RoleId
	}

	if len(data.RoleId.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing AdminId", fmt.Sprintf("RoleId: %s", data.RoleId.ValueString()))
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	updateOrganizationSamlRole := *openApiClient.NewInlineObject217()
	updateOrganizationSamlRole.SetRole(data.Role.ValueString())
	updateOrganizationSamlRole.SetOrgAccess(data.OrgAccess.ValueString())

	// Tags
	if len(data.Tags) > 0 {
		var tags []openApiClient.OrganizationsOrganizationIdSamlRolesTags
		for _, attribute := range data.Tags {
			var tag openApiClient.OrganizationsOrganizationIdSamlRolesTags
			tag.Tag = attribute.Tag.ValueString()
			tag.Access = attribute.Access.ValueString()
			tags = append(tags, tag)
		}
		updateOrganizationSamlRole.SetTags(tags)
	}

	// Networks
	if len(data.Networks) > 0 {
		var networks []openApiClient.OrganizationsOrganizationIdSamlRolesNetworks
		for _, attribute := range data.Networks {
			var network openApiClient.OrganizationsOrganizationIdSamlRolesNetworks
			network.Id = attribute.Id.ValueString()
			network.Access = attribute.Access.ValueString()
			networks = append(networks, network)
		}
		updateOrganizationSamlRole.SetNetworks(networks)
	}

	inlineResp, httpResp, err := r.client.OrganizationsApi.UpdateOrganizationSamlRole(context.Background(), data.OrgId.ValueString(), data.RoleId.ValueString()).UpdateOrganizationSamlRole(updateOrganizationSamlRole).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

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
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%s", data))
		return
	}

	data.Id = types.StringValue("example-id")
	data.OrgId = types.StringValue(inlineResp.GetId())
	data.RoleId = types.StringValue(inlineResp.GetId())
	data.Role = types.StringValue(inlineResp.GetRole())
	data.OrgAccess = types.StringValue(inlineResp.GetOrgAccess())

	// tags attribute
	if tags := inlineResp.Tags; tags != nil {
		for _, tv := range tags {
			var tag OrganizationsSamlRoleResourceModelTag
			tag.Tag = types.StringValue(*tv.Tag)
			tag.Access = types.StringValue(*tv.Access)
			data.Tags = append(data.Tags, tag)
		}
	} else {
		data.Tags = nil
	}

	// networks attribute
	if networks := inlineResp.Networks; networks != nil {
		for _, nw := range networks {
			var network OrganizationsSamlRoleResourceModelNetwork
			network.Id = types.StringValue(*nw.Id)
			network.Access = types.StringValue(*nw.Access)
			data.Networks = append(data.Networks, network)
		}
	} else {
		data.Networks = nil
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *OrganizationsSamlrolesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *OrganizationsSamlrolesResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if len(data.OrgId.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing OrganizationId", fmt.Sprintf("Value: %s", data.OrgId.ValueString()))
		return
	}

	if len(data.RoleId.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing RoleId", fmt.Sprintf("Value: %s", data.RoleId.ValueString()))
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.client.OrganizationsApi.DeleteOrganizationSamlRole(context.Background(), data.OrgId.ValueString(), data.RoleId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to delete resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	if httpResp.StatusCode != 204 {
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
	}

	resp.State.RemoveResource(ctx)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")

}

func (r *OrganizationsSamlrolesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: organization_id, role_id. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("role_id"), idParts[1])...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func extractHttpResponseOrganizationSamlRoleResource(ctx context.Context, inlineRespValue map[string]interface{}, data *OrganizationsSamlrolesResourceModel) *OrganizationsSamlrolesResourceModel {
	// save into the Terraform state
	data.Id = types.StringValue("example-id")

	// role id attribute
	if role_id := inlineRespValue["id"]; role_id != nil {
		data.RoleId = types.StringValue(role_id.(string))
	} else {
		data.RoleId = types.StringNull()
	}

	// role name attribute
	if role := inlineRespValue["role"]; role != nil {
		data.Role = types.StringValue(role.(string))
	} else {
		data.Role = types.StringNull()
	}

	// orgAccess attribute
	if orgAccess := inlineRespValue["orgAccess"]; orgAccess != nil {
		data.OrgAccess = types.StringValue(orgAccess.(string))
	} else {
		data.OrgAccess = types.StringNull()
	}

	// tags attribute
	if tags := inlineRespValue["tags"]; tags != nil {
		data.Tags = nil
		for _, tv := range tags.([]interface{}) {
			var ruleResult OrganizationsSamlRoleResourceModelTag
			rule := tv.(map[string]interface{})

			// policy attribute
			if policy := rule["tag"]; policy != nil {
				ruleResult.Tag = types.StringValue(policy.(string))
			} else {
				ruleResult.Tag = types.StringNull()
			}
			if policy := rule["access"]; policy != nil {
				ruleResult.Access = types.StringValue(policy.(string))
			} else {
				ruleResult.Access = types.StringNull()
			}
			data.Tags = append(data.Tags, ruleResult)
		}
	} else {
		data.Tags = nil
	}

	// networks attribute
	if networks := inlineRespValue["networks"]; networks != nil {
		data.Networks = nil
		for _, nv := range networks.([]interface{}) {
			var network OrganizationsSamlRoleResourceModelNetwork
			rulen := nv.(map[string]interface{})
			// policy attribute
			if policyn := rulen["id"]; policyn != nil {
				network.Id = types.StringValue(policyn.(string))
			} else {
				network.Id = types.StringNull()
			}
			if policyp := rulen["access"]; policyp != nil {
				network.Access = types.StringValue(policyp.(string))
			} else {
				network.Access = types.StringNull()
			}
			data.Networks = append(data.Networks, network)
		}
	} else {
		data.Networks = nil

	}

	return data
}
