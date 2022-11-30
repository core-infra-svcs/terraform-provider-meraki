package provider

import (
	"context"
	"encoding/json"
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
var _ resource.Resource = &OrganizationsAdaptivepolicyAclResource{}
var _ resource.ResourceWithImportState = &OrganizationsAdaptivepolicyAclResource{}

func NewOrganizationsAdaptivepolicyAclResource() resource.Resource {
	return &OrganizationsAdaptivepolicyAclResource{}
}

// OrganizationsAdaptivepolicyAclResource defines the resource implementation.
type OrganizationsAdaptivepolicyAclResource struct {
	client *apiclient.APIClient
}

// OrganizationsAdaptivepolicyAclResourceModel describes the resource data model.
type OrganizationsAdaptivepolicyAclResourceModel struct {
	Id          types.String `tfsdk:"id"`
	AclId       types.String `tfsdk:"aclid"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	IpVersion   types.String `tfsdk:"ipversion"`
	Rules       []RulesData  `tfsdk:"rules"`
}

// AclInfo  describes the acl data model
type AclInfo struct {
	Name        string
	AclId       string
	Description string
	IpVersion   string
	Rules       []RulesData
}

// RulesData  describes the rules data model
type RulesData struct {
	Policy   string  `tfsdk:"policy"`
	Protocol string  `tfsdk:"protocol"`
	SrcPort  *string `tfsdk:"srcport"`
	DstPort  *string `tfsdk:"dstport"`
}

func (r *OrganizationsAdaptivepolicyAclResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizations_adaptivePolicy_acl"
}

func (r *OrganizationsAdaptivepolicyAclResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		// This description is used by the documentation generator and the language server.
		MarkdownDescription: "OrganizationsAdaptivepolicyAcl resource  Manage the acls for an organization",
		Attributes: map[string]tfsdk.Attribute{

			"id": {
				MarkdownDescription: "Organization ID",
				Type:                types.StringType,
				Required:            true,
				Optional:            false,
				Computed:            false,
			},
			"aclid": {
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
				Required:            true,
				Optional:            false,
				Computed:            false,
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
				Required:            true,
				Optional:            false,
				Computed:            false,
				Sensitive:           false,
				Attributes:          nil,
				DeprecationMessage:  "",
				Validators:          nil,
				PlanModifiers:       nil,
			},
			"ipversion": {
				Description:         "IP version of adaptive policy ACL. One of: any, ipv4 or ipv6",
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
			"rules": {
				Description: "An ordered array of the adaptive policy ACL rules.",
				Required:    true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
					"policy": {
						Description:         "'allow' or 'deny' traffic specified by this rule.",
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
					"protocol": {
						Description:         "The type of protocol (must be 'tcp', 'udp', 'icmp' or 'any').",
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
					"srcport": {
						Description:         "Source port. Must be in the format of single port: '1', port list: '1,2' or port range: '1-10', and in the range of 1-65535, or 'any'. Default is 'any'.",
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
					"dstport": {
						Description:         "Destination port. Must be in the format of single port: '1', port list: '1,2' or port range: '1-10', and in the range of 1-65535, or 'any'. Default is 'any'.",
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

func (r *OrganizationsAdaptivepolicyAclResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *OrganizationsAdaptivepolicyAclResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *OrganizationsAdaptivepolicyAclResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var v []apiclient.OrganizationsOrganizationIdAdaptivePolicyAclsRules
	if data.Rules != nil && len(data.Rules) == 0 {
		resp.Diagnostics.AddError("Rules should not be empty.", fmt.Sprintf("rules: %v", data.Rules))
		return
	} else {
		for _, ruledata := range data.Rules {
			var r apiclient.OrganizationsOrganizationIdAdaptivePolicyAclsRules
			r.DstPort = ruledata.DstPort
			r.Policy = ruledata.Policy
			r.Protocol = ruledata.Protocol
			r.SrcPort = ruledata.SrcPort
			v = append(v, r)

		}
	}

	createOrganizationsAdaptivepolicyAcl := *apiclient.NewInlineObject169(data.Name.ValueString(), v, data.IpVersion.ValueString())
	createOrganizationsAdaptivepolicyAcl.SetDescription(data.Description.ValueString())
	inlineResp, httpResp, err := r.client.OrganizationsApi.CreateOrganizationAdaptivePolicyAcl(context.Background(), data.Id.ValueString()).CreateOrganizationAdaptivePolicyAcl(createOrganizationsAdaptivepolicyAcl).Execute()
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

	// Convert map to struct
	result, err := ConvertToSingleAclData(inlineResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Convert map to struct",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append()

	data.IpVersion = types.StringValue(result.IpVersion)
	data.Name = types.StringValue(result.Name)
	data.AclId = types.StringValue(result.AclId)
	data.Description = types.StringValue(result.Description)
	if data.Rules != nil {
		data.Rules = result.Rules
	}

	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OrganizationsAdaptivepolicyAclResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *OrganizationsAdaptivepolicyAclResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.AclId.IsUnknown() || data.AclId.IsNull() {

		inlineGetAclResp, httpResp, err := r.client.OrganizationsApi.GetOrganizationAdaptivePolicyAcls(context.Background(), data.Id.ValueString()).Execute()
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

		// Convert map to list of acl data
		acldata, err := ConvertToSingleAclDataList(inlineGetAclResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to Convert map to list of acl data",
				fmt.Sprintf("%v\n", err.Error()),
			)
		}

		// Check for errors after diagnostics collected
		if resp.Diagnostics.HasError() {
			return
		}

		resp.Diagnostics.Append()

		aclsearch := false

		for _, acl := range acldata {
			if data.Name.ValueString() == acl.Name {
				data.AclId = types.StringValue(acl.AclId)
				data.IpVersion = types.StringValue(acl.IpVersion)
				data.Name = types.StringValue(acl.Name)
				data.AclId = types.StringValue(acl.AclId)
				data.Description = types.StringValue(acl.Description)
				if acl.Rules != nil {
					data.Rules = acl.Rules
				}
				aclsearch = true
			}
		}

		if !aclsearch {
			resp.Diagnostics.AddError("No Acl details found. if you want to update acl name please add aclid field in resource declaration", fmt.Sprintf("acl name: %s", data.Name))
			return
		}
	} else {

		inlineResp, httpResp, err := r.client.OrganizationsApi.GetOrganizationAdaptivePolicyAcl(context.Background(), data.Id.ValueString(), data.AclId.ValueString()).Execute()
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

		// Convert map to struct
		result, err := ConvertToSingleAclData(inlineResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to Convert map to struct",
				fmt.Sprintf("%v\n", err.Error()),
			)
		}

		// Check for errors after diagnostics collected
		if resp.Diagnostics.HasError() {
			return
		}

		resp.Diagnostics.Append()

		data.IpVersion = types.StringValue(result.IpVersion)
		data.Name = types.StringValue(result.Name)
		data.AclId = types.StringValue(result.AclId)
		data.Description = types.StringValue(result.Description)
		if result.Rules != nil {
			data.Rules = result.Rules
		}

	}
	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OrganizationsAdaptivepolicyAclResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *OrganizationsAdaptivepolicyAclResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.AclId.IsUnknown() || data.AclId.IsNull() {

		inlinegetResp, httpResp, err := r.client.OrganizationsApi.GetOrganizationAdaptivePolicyAcls(context.Background(), data.Id.ValueString()).Execute()
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

		// Convert map to list of acl data
		acldata, err := ConvertToSingleAclDataList(inlinegetResp)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to Convert map to list of acl data",
				fmt.Sprintf("%v\n", err.Error()),
			)
		}

		// Check for errors after diagnostics collected
		if resp.Diagnostics.HasError() {
			return
		}
		resp.Diagnostics.Append()

		aclsearch := false

		for _, acl := range acldata {
			if data.Name.ValueString() == acl.Name {
				data.AclId = types.StringValue(acl.AclId)
				aclsearch = true
			}
		}

		if !aclsearch {
			resp.Diagnostics.AddError("No Acl details found. if you want to update acl name please add aclid field in resource declaration", fmt.Sprintf("acl name: %s", data.Name))
			return
		}
	}

	var v []apiclient.OrganizationsOrganizationIdAdaptivePolicyAclsRules
	if data.Rules != nil {
		if len(data.Rules) != 0 {
			for _, ruledata := range data.Rules {
				var r apiclient.OrganizationsOrganizationIdAdaptivePolicyAclsRules
				r.DstPort = ruledata.DstPort
				r.Policy = ruledata.Policy
				r.Protocol = ruledata.Protocol
				r.SrcPort = ruledata.SrcPort
				v = append(v, r)

			}
		} else {
			resp.Diagnostics.AddError("rules should not be empty", fmt.Sprintf("rules: %v", data.Rules))
			return
		}

	}
	updateOrganizationAdaptivePolicyAcl := *apiclient.NewInlineObject170()
	updateOrganizationAdaptivePolicyAcl.SetRules(v)
	updateOrganizationAdaptivePolicyAcl.SetName(data.Name.ValueString())
	updateOrganizationAdaptivePolicyAcl.SetIpVersion(data.IpVersion.ValueString())
	updateOrganizationAdaptivePolicyAcl.SetDescription(data.Description.ValueString())
	inlineUpdatedResp, httpResp, err := r.client.OrganizationsApi.UpdateOrganizationAdaptivePolicyAcl(context.Background(), data.Id.ValueString(), data.AclId.ValueString()).UpdateOrganizationAdaptivePolicyAcl(updateOrganizationAdaptivePolicyAcl).Execute()
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

	aclUpdatedData, err := ConvertToSingleAclData(inlineUpdatedResp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Convert map to struct",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	}
	resp.Diagnostics.Append()

	data.IpVersion = types.StringValue(aclUpdatedData.IpVersion)
	data.Name = types.StringValue(aclUpdatedData.Name)
	data.AclId = types.StringValue(aclUpdatedData.AclId)
	data.Description = types.StringValue(aclUpdatedData.Description)
	if data.Rules != nil {
		data.Rules = aclUpdatedData.Rules
	}

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *OrganizationsAdaptivepolicyAclResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *OrganizationsAdaptivepolicyAclResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.client.OrganizationsApi.DeleteOrganizationAdaptivePolicyAcl(context.Background(), data.Id.ValueString(), data.AclId.ValueString()).Execute()
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
	resp.Diagnostics.Append()

	resp.State.RemoveResource(ctx)

}

func (r *OrganizationsAdaptivepolicyAclResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: id,name. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), idParts[1])...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Convert to ACL Data
func ConvertToSingleAclData(inlineResp map[string]interface{}) (AclInfo, error) {

	var aclData AclInfo
	// Convert map to json string
	jsongetStr, err := json.Marshal(inlineResp)
	if err != nil {
		return aclData, err

	}
	// Convert json string to struct
	if err := json.Unmarshal(jsongetStr, &aclData); err != nil {
		return aclData, err
	}

	return aclData, err

}

// Convert to ACL Data List
func ConvertToSingleAclDataList(inlineResp []map[string]interface{}) ([]AclInfo, error) {

	var aclDataList []AclInfo
	// Convert map to json string
	jsongetStr, err := json.Marshal(inlineResp)
	if err != nil {
		return aclDataList, err

	}
	// Convert json string to struct
	if err := json.Unmarshal(jsongetStr, &aclDataList); err != nil {
		return aclDataList, err
	}

	return aclDataList, err

}
