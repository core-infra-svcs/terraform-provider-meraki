package provider

import (
	"context"
	"fmt"
	apiclient "github.com/core-infra-svcs/dashboard-api-go/client"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &OrganizationsAdaptivePolicyAclsDataSource{}

func NewOrganizationsAdaptivePolicyAclsDataSource() datasource.DataSource {
	return &OrganizationsAdaptivePolicyAclsDataSource{}
}

// OrganizationsAdaptivePolicyAclsDataSource defines the data source implementation.
type OrganizationsAdaptivePolicyAclsDataSource struct {
	client *apiclient.APIClient
}

// OrganizationsAdaptivePolicyAclsDataSourceModel describes the data source data model.
type OrganizationsAdaptivePolicyAclsDataSourceModel struct {
	Id    types.String                                    `tfsdk:"id"`
	OrgId types.String                                    `tfsdk:"organization_id"`
	List  []OrganizationAdaptivePolicyAclsDataSourceModel `tfsdk:"list"`
}

// OrganizationAdaptivePolicyAclsDataSourceModel describes the acl data source data model.
type OrganizationAdaptivePolicyAclsDataSourceModel struct {
	AclId       types.String                                         `tfsdk:"acl_id"`
	Name        types.String                                         `tfsdk:"name"`
	Description types.String                                         `tfsdk:"description"`
	IpVersion   types.String                                         `tfsdk:"ip_version"`
	Rules       []OrganizationAdaptivePolicyAclsDataSourceModelRules `tfsdk:"rules"`
	CreatedAt   types.String                                         `tfsdk:"created_at"`
	UpdatedAt   types.String                                         `tfsdk:"updated_at"`
}

type OrganizationAdaptivePolicyAclsDataSourceModelRules struct {
	Policy   types.String `tfsdk:"policy"`
	Protocol types.String `tfsdk:"protocol"`
	SrcPort  types.String `tfsdk:"src_port"`
	DstPort  types.String `tfsdk:"dst_port"`
}

func (d *OrganizationsAdaptivePolicyAclsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizations_adaptive_policy_acls"
}

func (d *OrganizationsAdaptivePolicyAclsDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{

		MarkdownDescription: "OrganizationsAdaptivePolicyAcls data source - get all list of  acls in an organization",
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
			"list": {
				MarkdownDescription: "List of organization acls",
				Optional:            true,
				Attributes: tfsdk.SetNestedAttributes(map[string]tfsdk.Attribute{
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
				})},
		},
	}, nil
}

func (d *OrganizationsAdaptivePolicyAclsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *OrganizationsAdaptivePolicyAclsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OrganizationsAdaptivePolicyAclsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	// Check for required parameters
	if len(data.OrgId.ValueString()) < 1 {
		resp.Diagnostics.AddError("Missing organizationId", fmt.Sprintf("Value: %s", data.OrgId.ValueString()))
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}

	inlineResp, httpResp, err := d.client.OrganizationsApi.GetOrganizationAdaptivePolicyAcls(context.Background(), data.OrgId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read datasource",
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
		resp.Diagnostics.AddError("State Data", fmt.Sprintf("\n%v", data))
		return
	}

	// Save data into Terraform state
	data.Id = types.StringValue("example-id")

	// adaptivePolicies attribute
	if adaptivePolicies := inlineResp; adaptivePolicies != nil {

		for _, inlineRespValue := range adaptivePolicies {
			var adaptivePolicy OrganizationAdaptivePolicyAclsDataSourceModel

			// id attribute
			if id := inlineRespValue["aclId"]; id != nil {
				adaptivePolicy.AclId = types.StringValue(id.(string))
			} else {
				adaptivePolicy.AclId = types.StringNull()
			}

			// description attribute
			if description := inlineRespValue["description"]; description != nil {
				adaptivePolicy.Description = types.StringValue(description.(string))
			} else {
				adaptivePolicy.Description = types.StringNull()
			}

			// ipVersion attribute
			if ipVersion := inlineRespValue["ipVersion"]; ipVersion != nil {
				adaptivePolicy.IpVersion = types.StringValue(ipVersion.(string))
			} else {
				adaptivePolicy.IpVersion = types.StringNull()
			}

			// rules attribute
			if rules := inlineRespValue["rules"]; rules != nil {
				for _, v := range rules.([]interface{}) {
					rule := v.(map[string]interface{})
					var ruleResult OrganizationAdaptivePolicyAclsDataSourceModelRules

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
					adaptivePolicy.Rules = append(adaptivePolicy.Rules, ruleResult)
				}

			}

			// updatedAt attribute
			if createdAt := inlineRespValue["createdAt"]; createdAt != nil {
				adaptivePolicy.CreatedAt = types.StringValue(createdAt.(string))
			} else {
				adaptivePolicy.CreatedAt = types.StringNull()
			}

			// updatedAt attribute
			if updatedAt := inlineRespValue["updatedAt"]; updatedAt != nil {
				adaptivePolicy.UpdatedAt = types.StringValue(updatedAt.(string))
			} else {
				adaptivePolicy.UpdatedAt = types.StringNull()
			}

			// append adaptivePolicy to list of adaptivePolicies
			data.List = append(data.List, adaptivePolicy)
		}

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read data source")
}
