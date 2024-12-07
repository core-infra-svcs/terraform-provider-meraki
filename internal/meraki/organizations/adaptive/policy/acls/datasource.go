package acls

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &DataSource{}

func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

type DataSource struct {
	client *openApiClient.APIClient
}

type dataSourceModel struct {
	Id   jsontypes.String      `tfsdk:"id" json:"organization_id"`
	List []dataSourceModelList `tfsdk:"list"`
}

type dataSourceModelList struct {
	AclId       jsontypes.String       `tfsdk:"acl_id" json:"AclId"`
	Name        jsontypes.String       `tfsdk:"name"`
	Description jsontypes.String       `tfsdk:"description"`
	IpVersion   jsontypes.String       `tfsdk:"ip_version" json:"IpVersion"`
	Rules       []dataSourceModelRules `tfsdk:"rules"`
	CreatedAt   jsontypes.String       `tfsdk:"created_at" json:"createdAt"`
	UpdatedAt   jsontypes.String       `tfsdk:"updated_at" json:"updatedAt"`
}

type dataSourceModelRules struct {
	Policy   jsontypes.String `tfsdk:"policy"`
	Protocol jsontypes.String `tfsdk:"protocol"`
	SrcPort  jsontypes.String `tfsdk:"src_port" json:"srcPort"`
	DstPort  jsontypes.String `tfsdk:"dst_port" json:"dstPort"`
}

func (d *DataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizations_adaptive_policy_acls"
}

func (d *DataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Ports adaptive policy ACLs in a organization",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.StringType,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 31),
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

func (d *DataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data dataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := d.client.OrganizationsApi.GetOrganizationAdaptivePolicyAcls(context.Background(), data.Id.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			utils.HttpDiagnostics(httpResp),
		)
		return
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
		resp.Diagnostics.AddError("State Data", fmt.Sprintf("\n%v", data))
		return
	}

	// Save data into Terraform state
	if err = json.NewDecoder(httpResp.Body).Decode(&data.List); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read data source")
}
