package provider

import (
	"context"
	"fmt"
	"strings"

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
var _ resource.Resource = &OrganizationsAdaptivePolicyAclResource{}
var _ resource.ResourceWithImportState = &OrganizationsAdaptivePolicyAclResource{}

func NewOrganizationsAdaptivePolicyAclResource() resource.Resource {
	return &OrganizationsAdaptivePolicyAclResource{}
}

// OrganizationsAdaptivePolicyAclResource defines the resource implementation.
type OrganizationsAdaptivePolicyAclResource struct {
	client *apiclient.APIClient
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

func (r *OrganizationsAdaptivePolicyAclResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "OrganizationsAdaptivePolicyAcl resource  Manage the acls for an organization",
		Attributes: map[string]tfsdk.Attribute{

			"id": {
				Description:         "Example identifier",
				MarkdownDescription: "Example identifier",
				Type:                types.StringType,
				Required:            false,
				Optional:            false,
				Computed:            true,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"organization_id": {
				Description:         "Organization Id",
				MarkdownDescription: "The Id of the organization",
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
			"acl_id": {
				Description:         "Acl ID",
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
				Description:         "Name of the adaptive policy ACL",
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
			"description": {
				Description:         "Description of the adaptive policy ACL",
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
			"ip_version": {
				Description:         "IP version of adaptive policy ACL. One of: any, ipv4 or ipv6",
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
			"rules": {
				Description: "An ordered array of the adaptive policy ACL rules.",
				Optional:    true,
				Computed:    true,
				// Type: types.ListType{ElemType: types.SetType{ElemType: types.StringType}},
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"policy": {
						Description:         "'allow' or 'deny' traffic specified by this rule.",
						MarkdownDescription: "'allow' or 'deny' traffic specified by this rule.",
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
					"protocol": {
						Description:         "The type of protocol (must be 'tcp', 'udp', 'icmp' or 'any').",
						MarkdownDescription: "The type of protocol (must be 'tcp', 'udp', 'icmp' or 'any').",
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
					"src_port": {
						Description:         "Source port. Must be in the format of single port: '1', port list: '1,2' or port range: '1-10', and in the range of 1-65535, or 'any'. Default is 'any'.",
						MarkdownDescription: "Source port. Must be in the format of single port: '1', port list: '1,2' or port range: '1-10', and in the range of 1-65535, or 'any'. Default is 'any'.",
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
					"dst_port": {
						Description:         "Destination port. Must be in the format of single port: '1', port list: '1,2' or port range: '1-10', and in the range of 1-65535, or 'any'. Default is 'any'.",
						MarkdownDescription: "Destination port. Must be in the format of single port: '1', port list: '1,2' or port range: '1-10', and in the range of 1-65535, or 'any'. Default is 'any'.",
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
			"created_at": {
				Description:         "rule created timestamp",
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
			"updated_at": {
				Description:         "last updated timestamp",
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

func (r *OrganizationsAdaptivePolicyAclResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OrganizationsAdaptivePolicyAclResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *OrganizationsAdaptivePolicyAclResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// Check for required parameters
	if len(data.OrgId.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing organizationId on create", fmt.Sprintf("Value: %s", data.OrgId.ValueString()))
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// rules
	var rules []apiclient.OrganizationsOrganizationIdAdaptivePolicyAclsRules
	for _, attribute := range data.Rules {

		var rule apiclient.OrganizationsOrganizationIdAdaptivePolicyAclsRules
		rule.Protocol = attribute.Protocol.ValueString()
		rule.Policy = attribute.Policy.ValueString()

		srcPort := attribute.SrcPort.ValueString()
		rule.SrcPort = &srcPort

		dstPort := attribute.DstPort.ValueString()
		rule.DstPort = &dstPort

		rules = append(rules, rule)
	}

	// payload
	createOrganizationsAdaptivePolicyAcl := *apiclient.NewInlineObject169(data.Name.ValueString(), rules, data.IpVersion.ValueString())
	createOrganizationsAdaptivePolicyAcl.SetDescription(data.Description.ValueString())

	inlineResp, httpResp, err := r.client.OrganizationsApi.CreateOrganizationAdaptivePolicyAcl(context.Background(), data.OrgId.ValueString()).CreateOrganizationAdaptivePolicyAcl(createOrganizationsAdaptivePolicyAcl).Execute()
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

	// Check for required parameters
	if len(data.OrgId.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing organizationId on read", fmt.Sprintf("Value: %s", data.OrgId.ValueString()))
	}

	if len(data.AclId.ValueString()) == 0 {
		resp.Diagnostics.AddError("Missing acl Id on read", fmt.Sprintf("Value: %v", data.AclId.ValueString()))
	}

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
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%s", data))
		return
	}

	// Save data into Terraform state
	extractHttpResponseOrganizationAdaptivePolicyAclResource(ctx, inlineResp, data)
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

}

func (r *OrganizationsAdaptivePolicyAclResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *OrganizationsAdaptivePolicyAclResourceModel
	var stateData *OrganizationsAdaptivePolicyAclResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)

	// Check for required parameters
	if len(data.OrgId.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing organizationId on Update", fmt.Sprintf("Value: %s", data.OrgId.ValueString()))
	}

	// Check state data if missing from plan
	if len(data.AclId.ValueString()) < 1 {
		data.AclId = stateData.AclId
	}

	if len(data.AclId.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing AclId on update", fmt.Sprintf("AclId: %s", data.AclId.ValueString()))
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// rules
	var rules []apiclient.OrganizationsOrganizationIdAdaptivePolicyAclsRules
	for _, attribute := range data.Rules {

		var rule apiclient.OrganizationsOrganizationIdAdaptivePolicyAclsRules
		rule.Protocol = attribute.Protocol.ValueString()
		rule.Policy = attribute.Policy.ValueString()

		srcPort := attribute.SrcPort.ValueString()
		rule.SrcPort = &srcPort

		dstPort := attribute.DstPort.ValueString()
		rule.DstPort = &dstPort

		rules = append(rules, rule)
	}

	// payload
	createOrganizationsAdaptivePolicyAcl := *apiclient.NewInlineObject170()
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

	// Check for required parameters
	if len(data.OrgId.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing organizationId on Delete", fmt.Sprintf("Value: %s", data.OrgId.ValueString()))
	}

	if len(data.AclId.ValueString()) == 0 {
		resp.Diagnostics.AddError("Missing acl Id on delete", fmt.Sprintf("Value: %v", data.AclId.ValueString()))
	}

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
