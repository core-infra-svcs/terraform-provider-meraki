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
var _ datasource.DataSource = &OrganizationsAdaptivepolicyAclsDataSource{}

func NewOrganizationsAdaptivepolicyAclsDataSource() datasource.DataSource {
	return &OrganizationsAdaptivepolicyAclsDataSource{}
}

// OrganizationsAdaptivepolicyAclsDataSource defines the data source implementation.
type OrganizationsAdaptivepolicyAclsDataSource struct {
	client *apiclient.APIClient
}

// OrganizationsAdaptivepolicyAclsDataSourceModel describes the data source data model.
type OrganizationsAdaptivepolicyAclsDataSourceModel struct {
	Id   types.String                                    `tfsdk:"id"`
	List []OrganizationAdaptivepolicyAclsDataSourceModel `tfsdk:"list"`
}

// OrganizationAdaptivepolicyAclsDataSourceModel describes the acl data source data model.
type OrganizationAdaptivepolicyAclsDataSourceModel struct {
	AclId       types.String `tfsdk:"aclid"`
	Name        types.String `tfsdk:"name"`
	Description types.String `tfsdk:"description"`
	IpVersion   types.String `tfsdk:"ipversion"`
	Rules       []RulesData  `tfsdk:"rules"`
}

func (d *OrganizationsAdaptivepolicyAclsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizations_adaptivePolicy_acls"
}

func (d *OrganizationsAdaptivepolicyAclsDataSource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{

		MarkdownDescription: "OrganizationsAdaptivepolicyAcls data source - get all list of  acls in an organization",
		Attributes: map[string]tfsdk.Attribute{

			"id": {
				MarkdownDescription: "org id",
				Description:         "Organization Id",
				Type:                types.StringType,
				Required:            true,
			},
			"list": {
				MarkdownDescription: "List of organization acls",
				Optional:            true,
				Attributes: tfsdk.ListNestedAttributes(map[string]tfsdk.Attribute{
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
					"ipversion": {
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
				})},
		},
	}, nil
}

func (d *OrganizationsAdaptivepolicyAclsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *OrganizationsAdaptivepolicyAclsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OrganizationsAdaptivepolicyAclsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	inlineGetAclResp, httpResp, err := d.client.OrganizationsApi.GetOrganizationAdaptivePolicyAcls(context.Background(), data.Id.ValueString()).Execute()
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
	fmt.Println(acldata)

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append()

	for _, acl := range acldata {
		var result OrganizationAdaptivepolicyAclsDataSourceModel
		result.AclId = types.StringValue(acl.AclId)
		result.Name = types.StringValue(acl.Name)
		result.Description = types.StringValue(acl.Description)
		result.IpVersion = types.StringValue(acl.IpVersion)
		if acl.Rules != nil {
			result.Rules = acl.Rules
		}

		data.List = append(data.List, result)
	}

	// Write logs using the tflog package
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
