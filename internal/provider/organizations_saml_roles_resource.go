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
	Id        jsontypes.String                            `tfsdk:"id"`
	OrgId     jsontypes.String                            `tfsdk:"organization_id" json:"organization_id"`
	RoleId    jsontypes.String                            `tfsdk:"role_id" json:"id"`
	Role      jsontypes.String                            `tfsdk:"role" json:"role"`
	OrgAccess jsontypes.String                            `tfsdk:"org_access" json:"org_access"`
	Tags      []OrganizationsSamlRoleResourceModelTag     `tfsdk:"tags" json:"tags"`
	Networks  []OrganizationsSamlRoleResourceModelNetwork `tfsdk:"networks" json:"networks"`
}

type OrganizationsSamlRoleResourceModelTag struct {
	Tag    jsontypes.String `tfsdk:"tag" json:"tag"`
	Access jsontypes.String `tfsdk:"access" json:"access"`
}

type OrganizationsSamlRoleResourceModelNetwork struct {
	Id     jsontypes.String `tfsdk:"id" json:"id"`
	Access jsontypes.String `tfsdk:"access" json:"access"`
}

func (r *OrganizationsSamlrolesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizations_saml_role"
}

func (r *OrganizationsSamlrolesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage the saml roles in this organization",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.StringType,
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "Organization ID",
				Required:            true,
				CustomType:          jsontypes.StringType,
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
				CustomType:          jsontypes.StringType,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 31),
				},
			},
			"role": schema.StringAttribute{
				MarkdownDescription: "The role of the SAML administrator",
				Required:            true,
				CustomType:          jsontypes.StringType,
			},
			"org_access": schema.StringAttribute{
				MarkdownDescription: "The privilege of the SAML administrator on the organization. Can be one of 'none', 'read-only', 'full' or 'enterprise'",
				Required:            true,
				CustomType:          jsontypes.StringType,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"full", "read-only", "enterprise", "none"}...),
					stringvalidator.LengthAtLeast(4),
				},
			},
			"tags": schema.SetNestedAttribute{
				Description: "The list of tags that the SAML administrator has privleges on.",
				Optional:    true,
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
				Description: "The list of networks that the SAML administrator has privileges on.",
				Optional:    true,
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

	_, httpResp, err := r.client.OrganizationsApi.CreateOrganizationSamlRole(context.Background(), data.OrgId.ValueString()).CreateOrganizationSamlRole(createOrganizationSamlRole).Execute()

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
		return
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

	// Save data into Terraform state
	if err = json.NewDecoder(httpResp.Body).Decode(&data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")
}

func (r *OrganizationsSamlrolesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *OrganizationsSamlrolesResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := r.client.OrganizationsApi.GetOrganizationSamlRole(ctx, data.OrgId.ValueString(), data.RoleId.ValueString()).Execute()
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

	// Save data into Terraform state
	if err = json.NewDecoder(httpResp.Body).Decode(&data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *OrganizationsSamlrolesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data *OrganizationsSamlrolesResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

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

	_, httpResp, err := r.client.OrganizationsApi.UpdateOrganizationSamlRole(context.Background(), data.OrgId.ValueString(), data.RoleId.ValueString()).UpdateOrganizationSamlRole(updateOrganizationSamlRole).Execute()
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
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *OrganizationsSamlrolesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *OrganizationsSamlrolesResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

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

	// collect diagnostics
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
	}

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

/*
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
*/
