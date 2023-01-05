package provider

import (
	"context"
	"fmt"
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
	"strings"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &OrganizationsAdaptivePolicyAclResource{}
	_ resource.ResourceWithConfigure   = &OrganizationsAdaptivePolicyAclResource{}
	_ resource.ResourceWithImportState = &OrganizationsAdaptivePolicyAclResource{}
)

func NewOrganizationsAdaptivePolicyAclResource() resource.Resource {
	return &OrganizationsAdaptivePolicyAclResource{}
}

// OrganizationsAdaptivePolicyAclResource defines the resource implementation.
type OrganizationsAdaptivePolicyAclResource struct {
	client *openApiClient.APIClient
}

// OrganizationsAdaptivePolicyAclResourceModel describes the resource data model.
type OrganizationsAdaptivePolicyAclResourceModel struct {
	Id          types.String                         `tfsdk:"id"`
	OrgId       types.String                         `tfsdk:"organization_id"`
	AclId       types.String                         `tfsdk:"acl_id"`
	Name        types.String                         `tfsdk:"name"`
	Description types.String                         `tfsdk:"description"`
	IpVersion   types.String                         `tfsdk:"ip_version"`
	Rules       []OrganizationsAdaptivePolicyAclRule `tfsdk:"rules"`
	CreatedAt   types.String                         `tfsdk:"created_at"`
	UpdatedAt   types.String                         `tfsdk:"updated_at"`
}

// OrganizationsAdaptivePolicyAclRule  describes the rules data model
type OrganizationsAdaptivePolicyAclRule struct {
	Policy   types.String `tfsdk:"policy"`
	Protocol types.String `tfsdk:"protocol"`
	SrcPort  types.String `tfsdk:"src_port"`
	DstPort  types.String `tfsdk:"dst_port"`
}

func (r *OrganizationsAdaptivePolicyAclResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizations_adaptive_policy_acl"
}

func (r *OrganizationsAdaptivePolicyAclResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manage adaptive policy ACLs in a organization",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
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
			"acl_id": schema.StringAttribute{
				MarkdownDescription: "ACL ID",
				Optional:            true,
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(4, 31),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the adaptive policy ACL",
				Optional:            true,
				Computed:            true,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the adaptive policy ACL",
				Optional:            true,
				Computed:            true,
			},
			"ip_version": schema.StringAttribute{
				MarkdownDescription: "IP version of adaptive policy ACL. One of: 'any', 'ipv4' or 'ipv6",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("any"),
						path.MatchRoot("ipv4"),
						path.MatchRoot("ipv6"),
					),
				},
			},
			"rules": schema.ListNestedAttribute{
				Description: "An ordered array of the adaptive policy ACL rules. An empty array will clear the rules.",
				Optional:    true,
				Computed:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"policy": schema.StringAttribute{
							MarkdownDescription: "",
							Optional:            true,
							Computed:            true,
						},
						"protocol": schema.StringAttribute{
							MarkdownDescription: "",
							Optional:            true,
							Computed:            true,
						},
						"src_port": schema.StringAttribute{
							MarkdownDescription: "",
							Optional:            true,
							Computed:            true,
						},
						"dst_port": schema.StringAttribute{
							MarkdownDescription: "",
							Optional:            true,
							Computed:            true,
						},
					},
				},
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "",
				Optional:            true,
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "",
				Optional:            true,
				Computed:            true,
			},
		},
	}
}

func (r *OrganizationsAdaptivePolicyAclResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OrganizationsAdaptivePolicyAclResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *OrganizationsAdaptivePolicyAclResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// rules
	var rules []openApiClient.OrganizationsOrganizationIdAdaptivePolicyAclsRules
	for _, attribute := range data.Rules {

		var rule openApiClient.OrganizationsOrganizationIdAdaptivePolicyAclsRules
		rule.Protocol = attribute.Protocol.ValueString()
		rule.Policy = attribute.Policy.ValueString()

		srcPort := attribute.SrcPort.ValueString()
		rule.SrcPort = &srcPort

		dstPort := attribute.DstPort.ValueString()
		rule.DstPort = &dstPort

		rules = append(rules, rule)
	}

	// payload
	createOrganizationsAdaptivePolicyAcl := *openApiClient.NewInlineObject169(data.Name.ValueString(), rules, data.IpVersion.ValueString())
	createOrganizationsAdaptivePolicyAcl.SetDescription(data.Description.ValueString())

	inlineResp, httpResp, err := r.client.OrganizationsApi.CreateOrganizationAdaptivePolicyAcl(context.Background(), data.OrgId.ValueString()).CreateOrganizationAdaptivePolicyAcl(createOrganizationsAdaptivePolicyAcl).Execute()
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
	extractHttpResponseOrganizationAdaptivePolicyAclResource(ctx, inlineResp, data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")
}

func (r *OrganizationsAdaptivePolicyAclResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *OrganizationsAdaptivePolicyAclResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	inlineResp, httpResp, err := r.client.OrganizationsApi.GetOrganizationAdaptivePolicyAcl(context.Background(), data.OrgId.ValueString(), data.AclId.ValueString()).Execute()
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
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%s", data))
		return
	}

	// Save data into Terraform state
	extractHttpResponseOrganizationAdaptivePolicyAclResource(ctx, inlineResp, data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *OrganizationsAdaptivePolicyAclResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *OrganizationsAdaptivePolicyAclResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// rules
	var rules []openApiClient.OrganizationsOrganizationIdAdaptivePolicyAclsRules
	for _, attribute := range data.Rules {

		var rule openApiClient.OrganizationsOrganizationIdAdaptivePolicyAclsRules
		rule.Protocol = attribute.Protocol.ValueString()
		rule.Policy = attribute.Policy.ValueString()

		srcPort := attribute.SrcPort.ValueString()
		rule.SrcPort = &srcPort

		dstPort := attribute.DstPort.ValueString()
		rule.DstPort = &dstPort

		rules = append(rules, rule)
	}

	// payload
	createOrganizationsAdaptivePolicyAcl := *openApiClient.NewInlineObject170()
	createOrganizationsAdaptivePolicyAcl.SetName(data.Name.ValueString())
	createOrganizationsAdaptivePolicyAcl.SetDescription(data.Description.ValueString())
	createOrganizationsAdaptivePolicyAcl.SetRules(rules)
	createOrganizationsAdaptivePolicyAcl.SetIpVersion(data.IpVersion.ValueString())

	inlineResp, httpResp, err := r.client.OrganizationsApi.UpdateOrganizationAdaptivePolicyAcl(context.Background(), data.OrgId.ValueString(), data.AclId.ValueString()).UpdateOrganizationAdaptivePolicyAcl(createOrganizationsAdaptivePolicyAcl).Execute()
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
	extractHttpResponseOrganizationAdaptivePolicyAclResource(ctx, inlineResp, data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *OrganizationsAdaptivePolicyAclResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *OrganizationsAdaptivePolicyAclResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.client.OrganizationsApi.DeleteOrganizationAdaptivePolicyAcl(context.Background(), data.OrgId.ValueString(), data.AclId.ValueString()).Execute()
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
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%s", data))
		return
	}
	resp.State.RemoveResource(ctx)
}

func extractHttpResponseOrganizationAdaptivePolicyAclResource(ctx context.Context, inlineRespValue map[string]interface{}, data *OrganizationsAdaptivePolicyAclResourceModel) *OrganizationsAdaptivePolicyAclResourceModel {

	// save into the Terraform state
	data.Id = types.StringValue("example-id")

	// id attribute
	if id := inlineRespValue["aclId"]; id != nil {
		data.AclId = types.StringValue(id.(string))
	} else {
		data.AclId = types.StringNull()
	}

	// description attribute
	if description := inlineRespValue["description"]; description != nil {
		data.Description = types.StringValue(description.(string))
	} else {
		data.Description = types.StringNull()
	}

	// ipVersion attribute
	if ipVersion := inlineRespValue["ipVersion"]; ipVersion != nil {
		data.IpVersion = types.StringValue(ipVersion.(string))
	} else {
		data.IpVersion = types.StringNull()
	}

	// rules attribute
	if rules := inlineRespValue["rules"]; rules != nil {
		data.Rules = nil // prevents duplicate rule entries
		for _, v := range rules.([]interface{}) {
			rule := v.(map[string]interface{})
			var ruleResult OrganizationsAdaptivePolicyAclRule

			// policy attribute
			if policy := rule["policy"]; policy != nil {
				ruleResult.Policy = types.StringValue(policy.(string))
			} else {
				ruleResult.Policy = types.StringNull()
			}

			// protocol attribute
			if protocol := rule["protocol"]; protocol != nil {
				ruleResult.Protocol = types.StringValue(protocol.(string))
			} else {
				ruleResult.Protocol = types.StringNull()
			}

			// srcPort attribute
			if srcPort := rule["srcPort"]; srcPort != nil {
				ruleResult.SrcPort = types.StringValue(srcPort.(string))
			} else {
				ruleResult.SrcPort = types.StringNull()
			}

			// dstPort attribute
			if dstPort := rule["dstPort"]; dstPort != nil {
				ruleResult.DstPort = types.StringValue(dstPort.(string))
			} else {
				ruleResult.DstPort = types.StringNull()
			}
			data.Rules = append(data.Rules, ruleResult)
		}

	}

	// updatedAt attribute
	if createdAt := inlineRespValue["createdAt"]; createdAt != nil {
		data.CreatedAt = types.StringValue(createdAt.(string))
	} else {
		data.CreatedAt = types.StringNull()
	}

	// updatedAt attribute
	if updatedAt := inlineRespValue["updatedAt"]; updatedAt != nil {
		data.UpdatedAt = types.StringValue(updatedAt.(string))
	} else {
		data.UpdatedAt = types.StringNull()
	}

	return data
}

func (r *OrganizationsAdaptivePolicyAclResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: organization_id, acl_id. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("organization_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("acl_id"), idParts[1])...)
	if resp.Diagnostics.HasError() {
		return
	}

}
