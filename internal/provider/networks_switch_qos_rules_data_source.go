package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
	"net/http"
	"time"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &NetworksSwitchQosRulesDataSource{}

func NewNetworksSwitchQosRulesDataSource() datasource.DataSource {
	return &NetworksSwitchQosRulesDataSource{}
}

// NetworksSwitchQosRulesDataSource defines the resource implementation.
type NetworksSwitchQosRulesDataSource struct {
	client *openApiClient.APIClient
}

// NetworksSwitchQosRulesDataSourceModel describes the resource data model.
type NetworksSwitchQosRulesDataSourceModel struct {
	Id        jsontypes.String                             `tfsdk:"id" json:"-"`
	NetworkId jsontypes.String                             `tfsdk:"network_id" json:"network_id"`
	List      []NetworksSwitchQosRulesDataSourceModelRules `tfsdk:"list"`
}

// NetworksSwitchQosRulesDataSourceModelRules describes the resource data model.
type NetworksSwitchQosRulesDataSourceModelRules struct {
	QosRulesId   jsontypes.String  `tfsdk:"qos_rule_id" json:"id"`
	Vlan         jsontypes.Int64   `tfsdk:"vlan" json:"vlan"`
	Dscp         jsontypes.Int64   `tfsdk:"dscp" json:"dscp"`
	DstPort      jsontypes.Float64 `tfsdk:"dst_port" json:"dstPort"`
	SrcPort      jsontypes.Float64 `tfsdk:"src_port" json:"srcPort"`
	DstPortRange jsontypes.String  `tfsdk:"dst_port_range" json:"dstPortRange"`
	Protocol     jsontypes.String  `tfsdk:"protocol" json:"protocol"`
	SrcPortRange jsontypes.String  `tfsdk:"src_port_range" json:"srcPortRange"`
}

func (r *NetworksSwitchQosRulesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_switch_qos_rules"
}

func (r *NetworksSwitchQosRulesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "NetworksSwitchQosRule resource for updating network switch qos rule.",
		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{
				MarkdownDescription: "Qos Rules data source Id",
				Computed:            true,
				Optional:            true,
				CustomType:          jsontypes.StringType,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Network Id",
				Required:            true,
				CustomType:          jsontypes.StringType,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 31),
				},
			},
			"list": schema.ListNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "list of switch qos rules",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"qos_rule_id": schema.StringAttribute{
							MarkdownDescription: "Qos Rules Id",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"vlan": schema.Int64Attribute{
							MarkdownDescription: "The VLAN of the incoming packet. A null value will match any VLAN.",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.Int64Type,
						},
						"dscp": schema.Int64Attribute{
							MarkdownDescription: "DSCP tag. Set this to -1 to trust incoming DSCP. Default value is 0.",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.Int64Type,
						},
						"dst_port": schema.Float64Attribute{
							MarkdownDescription: "The destination port of the incoming packet. Applicable only if protocol is TCP or UDP.",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.Float64Type,
						},
						"src_port": schema.Float64Attribute{
							MarkdownDescription: "The source port of the incoming packet. Applicable only if protocol is TCP or UDP.",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.Float64Type,
						},
						"dst_port_range": schema.StringAttribute{
							MarkdownDescription: "The destination port range of the incoming packet. Applicable only if protocol is set to TCP or UDP. Example: 70-80",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"protocol": schema.StringAttribute{
							MarkdownDescription: "The protocol of the incoming packet. Can be one of ANY, TCP or UDP. Default value is ANY",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"src_port_range": schema.StringAttribute{
							MarkdownDescription: "The source port range of the incoming packet. Applicable only if protocol is set to TCP or UDP. Example: 70-80",
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

func (r *NetworksSwitchQosRulesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *NetworksSwitchQosRulesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *NetworksSwitchQosRulesDataSourceModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	maxRetries := r.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(r.client.GetConfig().Retry4xxErrorWaitTime)

	inlineResp, httpResp, err := tools.CustomHttpRequestRetry[[]map[string]interface{}](ctx, maxRetries, retryDelay, func() ([]map[string]interface{}, *http.Response, error) {
		inline, respHttp, err := r.client.QosRulesApi.GetNetworkSwitchQosRules(ctx, data.NetworkId.ValueString()).Execute()
		return inline, respHttp, err
	})

	if err != nil {
		tflog.Error(ctx, "HTTP Call Failed", map[string]interface{}{
			"error": err.Error(),
		})
		resp.Diagnostics.AddError(
			"HTTP Call Failed",
			fmt.Sprintf("Details: %s", err.Error()),
		)
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
	}

	// Check for API success response code
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	var rules []NetworksSwitchQosRulesDataSourceModelRules

	// Iterate through the inline response
	for _, item := range inlineResp {
		var rule NetworksSwitchQosRulesDataSourceModelRules

		// Convert the map item to a JSON string
		itemJSON, itemJSONErr := json.Marshal(item)
		if itemJSONErr != nil {
			resp.Diagnostics.AddError(
				"Failed to Convert the map item to a JSON string",
				fmt.Sprintf("%v", itemJSONErr),
			)
			return
		}

		// Unmarshal JSON string into NetworksSwitchQosRulesDataSourceModelRules
		err = json.Unmarshal(itemJSON, &rule)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to Unmarshal JSON string into rule data model",
				fmt.Sprintf("%v", err),
			)
			return
		}

		// Append the unmarshal rule to the list
		rules = append(rules, rule)
	}

	data.List = rules

	data.Id = data.NetworkId

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}
