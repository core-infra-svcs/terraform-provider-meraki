package provider

import (
	"context"
	"encoding/json"
	"fmt"
	openApiClient "github.com/core-infra-svcs/dashboard-api-go/client"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &OrganizationsAdaptivePolicyAclsDataSource{}

func NewOrganizationsAdaptivePolicyAclsDataSource() datasource.DataSource {
	return &OrganizationsAdaptivePolicyAclsDataSource{}
}

// OrganizationsAdaptivePolicyAclsDataSource defines the data source implementation.
type OrganizationsAdaptivePolicyAclsDataSource struct {
	client *openApiClient.APIClient
}

// OrganizationsAdaptivePolicyAclsDataSourceModel describes the data source data model.
type OrganizationsAdaptivePolicyAclsDataSourceModel struct {
	Id    jsontypes.String                                `tfsdk:"id"`
	OrgId jsontypes.String                                `tfsdk:"organization_id"`
	List  []OrganizationAdaptivePolicyAclsDataSourceModel `tfsdk:"list"`
}

// OrganizationAdaptivePolicyAclsDataSourceModel describes the acl data source data model.
type OrganizationAdaptivePolicyAclsDataSourceModel struct {
	AclId       jsontypes.String                                     `tfsdk:"acl_id" json:"AclId"`
	Name        jsontypes.String                                     `tfsdk:"name"`
	Description jsontypes.String                                     `tfsdk:"description"`
	IpVersion   jsontypes.String                                     `tfsdk:"ip_version" json:"IpVersion"`
	Rules       []OrganizationAdaptivePolicyAclsDataSourceModelRules `tfsdk:"rules"`
	CreatedAt   jsontypes.String                                     `tfsdk:"created_at" json:"createdAt"`
	UpdatedAt   jsontypes.String                                     `tfsdk:"updated_at" json:"updatedAt"`
}

type OrganizationAdaptivePolicyAclsDataSourceModelRules struct {
	Policy   jsontypes.String `tfsdk:"policy"`
	Protocol jsontypes.String `tfsdk:"protocol"`
	SrcPort  jsontypes.String `tfsdk:"src_port" json:"srcPort"`
	DstPort  jsontypes.String `tfsdk:"dst_port" json:"dstPort"`
}

func (d *OrganizationsAdaptivePolicyAclsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizations_adaptive_policy_acls"
}

func (d *OrganizationsAdaptivePolicyAclsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List adaptive policy ACLs in a organization",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.StringType,
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "Organization ID",
				Optional:            true,
				CustomType:          jsontypes.StringType,
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 31),
				},
			},
			"list": schema.ListNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"acl_id": schema.StringAttribute{
							MarkdownDescription: "ACL ID",
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the adaptive policy ACL",
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
						"description": schema.StringAttribute{
							MarkdownDescription: "Description of the adaptive policy ACL",
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
						"ip_version": schema.StringAttribute{
							MarkdownDescription: "IP version of adaptive policy ACL. One of: 'any', 'ipv4' or 'ipv6",
							Optional:            true,
							CustomType:          jsontypes.StringType,
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
										CustomType:          jsontypes.StringType,
									},
									"protocol": schema.StringAttribute{
										MarkdownDescription: "",
										Optional:            true,
										CustomType:          jsontypes.StringType,
									},
									"src_port": schema.StringAttribute{
										MarkdownDescription: "",
										Optional:            true,
										CustomType:          jsontypes.StringType,
									},
									"dst_port": schema.StringAttribute{
										MarkdownDescription: "",
										Optional:            true,
										CustomType:          jsontypes.StringType,
									},
								},
							},
						},
						"created_at": schema.StringAttribute{
							MarkdownDescription: "",
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
						"updated_at": schema.StringAttribute{
							MarkdownDescription: "",
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
					},
				},
			},
		},
	}
}

func (d *OrganizationsAdaptivePolicyAclsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*openApiClient.APIClient)

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
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("State Data", fmt.Sprintf("\n%v", data))
		return
	}

	// Save data into Terraform state
	data.Id = jsontypes.StringValue("example-id")

	// adaptivePolicies attribute
	if adaptivePolicies := inlineResp; adaptivePolicies != nil {

		for _, inlineRespValue := range adaptivePolicies {
			var adaptivePolicy OrganizationAdaptivePolicyAclsDataSourceModel

			// TODO - Workaround until json.RawMessage is implemented in HTTP client
			b, err := json.Marshal(inlineRespValue)
			if err != nil {
				resp.Diagnostics.AddError(
					"Failed to marshal API response",
					fmt.Sprintf("%v", err),
				)
			}
			if err := json.Unmarshal(b, &adaptivePolicy); err != nil {
				resp.Diagnostics.AddError(
					"Failed to unmarshal API response",
					fmt.Sprintf("Unmarshal error%v", err),
				)
			}

			// append adaptivePolicy to list of adaptivePolicies
			data.List = append(data.List, adaptivePolicy)
		}

	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read data source")
}
