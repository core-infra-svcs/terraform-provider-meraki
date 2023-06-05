package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strings"

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
	openApiClient "github.com/meraki/dashboard-api-go/client"
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
	OrgId       jsontypes.String                     `tfsdk:"organization_id" json:"organizationId"`
	AclId       jsontypes.String                     `tfsdk:"acl_id" json:"aclId"`
	Name        jsontypes.String                     `tfsdk:"name"`
	Description jsontypes.String                     `tfsdk:"description"`
	IpVersion   jsontypes.String                     `tfsdk:"ip_version" json:"ipVersion"`
	Rules       []OrganizationsAdaptivePolicyAclRule `tfsdk:"rules"`
	CreatedAt   jsontypes.String                     `tfsdk:"created_at" json:"createdAt"`
	UpdatedAt   jsontypes.String                     `tfsdk:"updated_at" json:"updatedAt"`
}

// OrganizationsAdaptivePolicyAclRule  describes the rules data model
type OrganizationsAdaptivePolicyAclRule struct {
	Policy   jsontypes.String `tfsdk:"policy"`
	Protocol jsontypes.String `tfsdk:"protocol"`
	SrcPort  jsontypes.String `tfsdk:"src_port" json:"srcPort"`
	DstPort  jsontypes.String `tfsdk:"dst_port" json:"dstPort"`
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
				CustomType:          jsontypes.StringType,
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
				CustomType:          jsontypes.StringType,
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
				CustomType:          jsontypes.StringType,
			},
			"description": schema.StringAttribute{
				MarkdownDescription: "Description of the adaptive policy ACL",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"ip_version": schema.StringAttribute{
				MarkdownDescription: "IP version of adaptive policy ACL. One of: 'any', 'ipv4' or 'ipv6",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
				Validators: []validator.String{
					stringvalidator.OneOf([]string{"any", "ipv4", "ipv6"}...),
					stringvalidator.LengthAtLeast(3),
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
							CustomType:          jsontypes.StringType,
						},
						"protocol": schema.StringAttribute{
							MarkdownDescription: "",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"src_port": schema.StringAttribute{
							MarkdownDescription: "",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"dst_port": schema.StringAttribute{
							MarkdownDescription: "",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
					},
				},
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
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
	var rules []openApiClient.OrganizationsOrganizationIdAdaptivePolicyAclsRules1
	for _, attribute := range data.Rules {

		var rule openApiClient.OrganizationsOrganizationIdAdaptivePolicyAclsRules1
		rule.Protocol = attribute.Protocol.ValueString()
		rule.Policy = attribute.Policy.ValueString()

		srcPort := attribute.SrcPort.ValueString()
		rule.SrcPort = &srcPort

		dstPort := attribute.DstPort.ValueString()
		rule.DstPort = &dstPort

		rules = append(rules, rule)
	}

	// payload
	createOrganizationsAdaptivePolicyAcl := *openApiClient.NewInlineObject171(data.Name.ValueString(), rules, data.IpVersion.ValueString())
	createOrganizationsAdaptivePolicyAcl.SetDescription(data.Description.ValueString())

	_, httpResp, err := r.client.OrganizationsApi.CreateOrganizationAdaptivePolicyAcl(context.Background(), data.OrgId.ValueString()).CreateOrganizationAdaptivePolicyAcl(createOrganizationsAdaptivePolicyAcl).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	// collect diagnostics
	tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)

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
	data.Id = types.StringValue("example-id")

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

	_, httpResp, err := r.client.OrganizationsApi.GetOrganizationAdaptivePolicyAcl(context.Background(), data.OrgId.ValueString(), data.AclId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	// collect diagnostics
	tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)

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
	data.Id = types.StringValue("example-id")

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
	var rules []openApiClient.OrganizationsOrganizationIdAdaptivePolicyAclsRules1
	for _, attribute := range data.Rules {

		var rule openApiClient.OrganizationsOrganizationIdAdaptivePolicyAclsRules1
		rule.Protocol = attribute.Protocol.ValueString()
		rule.Policy = attribute.Policy.ValueString()

		srcPort := attribute.SrcPort.ValueString()
		rule.SrcPort = &srcPort

		dstPort := attribute.DstPort.ValueString()
		rule.DstPort = &dstPort

		rules = append(rules, rule)
	}

	// payload
	createOrganizationsAdaptivePolicyAcl := *openApiClient.NewInlineObject172()
	createOrganizationsAdaptivePolicyAcl.SetName(data.Name.ValueString())
	createOrganizationsAdaptivePolicyAcl.SetDescription(data.Description.ValueString())
	createOrganizationsAdaptivePolicyAcl.SetRules(rules)
	createOrganizationsAdaptivePolicyAcl.SetIpVersion(data.IpVersion.ValueString())

	_, httpResp, err := r.client.OrganizationsApi.UpdateOrganizationAdaptivePolicyAcl(context.Background(), data.OrgId.ValueString(), data.AclId.ValueString()).UpdateOrganizationAdaptivePolicyAcl(createOrganizationsAdaptivePolicyAcl).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	// collect diagnostics
	tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)

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
	data.Id = types.StringValue("example-id")

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
		return
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
		return
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%s", data))
		return
	}
	resp.State.RemoveResource(ctx)
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
