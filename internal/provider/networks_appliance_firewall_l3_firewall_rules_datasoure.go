package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/core-infra-svcs/terraform-provider-meraki/tools"

	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &NetworksApplianceFirewallL3FirewallRulesDataSource{}

func NewNetworksApplianceFirewallL3FirewallRulesDataSource() datasource.DataSource {
	return &NetworksApplianceFirewallL3FirewallRulesDataSource{}
}

// NetworksApplianceFirewallL3FirewallRulesDataSource defines the data source implementation.
type NetworksApplianceFirewallL3FirewallRulesDataSource struct {
	client *openApiClient.APIClient
}

type NetworksApplianceFirewallL3FirewallRulesDataSourceModel struct {
	Id                jsontypes.String                                              `tfsdk:"id"`
	NetworkId         jsontypes.String                                              `tfsdk:"network_id" json:"network_id"`
	SyslogDefaultRule jsontypes.Bool                                                `tfsdk:"syslog_default_rule"`
	List              []NetworksApplianceFirewallL3FirewallRulesDataSourceModelList `tfsdk:"rules" json:"rules"`
}

// NetworksApplianceFirewallL3FirewallRulesDataSourceModelList describes the data source data model.
type NetworksApplianceFirewallL3FirewallRulesDataSourceModelList struct {
	Comment       jsontypes.String `tfsdk:"comment"`
	DestCidr      jsontypes.String `tfsdk:"dest_cidr"`
	DestPort      jsontypes.String `tfsdk:"dest_port"`
	Policy        jsontypes.String `tfsdk:"policy"`
	Protocol      jsontypes.String `tfsdk:"protocol"`
	SrcPort       jsontypes.String `tfsdk:"src_port"`
	SrcCidr       jsontypes.String `tfsdk:"src_cidr"`
	SysLogEnabled jsontypes.Bool   `tfsdk:"syslog_enabled"`
}

func (d *NetworksApplianceFirewallL3FirewallRulesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_appliance_firewall_l3_firewall_rules"
}

func (d *NetworksApplianceFirewallL3FirewallRulesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Get l3 firewall rules",

		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.StringType,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Network ID",
				Required:            true,
				CustomType:          jsontypes.StringType,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 31),
				},
			},
			"syslog_default_rule": schema.BoolAttribute{
				MarkdownDescription: "Log the special default rule (boolean value - enable only if you've configured a syslog server) (optional)",
				Optional:            true,
				CustomType:          jsontypes.BoolType,
			},
			"rules": schema.SetNestedAttribute{
				Optional: true,
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"comment": schema.StringAttribute{
							MarkdownDescription: "Description of the rule (optional)",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"dest_cidr": schema.StringAttribute{
							MarkdownDescription: "Comma-separated list of destination IP address(es) (in IP or CIDR notation), fully-qualified domain names (FQDN) or 'Any'",
							Required:            true,
							CustomType:          jsontypes.StringType,
						},
						"dest_port": schema.StringAttribute{
							MarkdownDescription: "Comma-separated list of destination port(s) (integer in the range 1-65535), or 'Any'",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"src_cidr": schema.StringAttribute{
							MarkdownDescription: "Comma-separated list of source IP address(es) (in IP or CIDR notation), or 'any' (note: FQDN not supported for source addresses)",
							Required:            true,
							CustomType:          jsontypes.StringType,
						},
						"src_port": schema.StringAttribute{
							MarkdownDescription: "Comma-separated list of source port(s) (integer in the range 1-65535), or 'Any'",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"policy": schema.StringAttribute{
							MarkdownDescription: "'allow' or 'deny' traffic specified by this rule",
							Required:            true,
							CustomType:          jsontypes.StringType,
						},
						"protocol": schema.StringAttribute{
							MarkdownDescription: "The type of protocol (must be 'tcp', 'udp', 'icmp', 'icmp6', 'Any', or 'any')",
							Required:            true,
							CustomType:          jsontypes.StringType,
							Validators: []validator.String{
								stringvalidator.OneOf([]string{"tcp", "udp", "icmp", "icmp6", "Any", "any"}...),
							},
						},
						"syslog_enabled": schema.BoolAttribute{
							MarkdownDescription: "Log this rule to syslog (true or false, boolean value) - only applicable if a syslog has been configured (optional)",
							Optional:            true,
							CustomType:          jsontypes.BoolType,
						},
					},
				},
			},
		},
	}
}

func (d *NetworksApplianceFirewallL3FirewallRulesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NetworksApplianceFirewallL3FirewallRulesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var data *NetworksApplianceFirewallL3FirewallRulesDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := d.client.ApplianceApi.GetNetworkApplianceFirewallL3FirewallRules(context.Background(), data.NetworkId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
	}

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
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	// Check if the default rule is nil
	if data.SyslogDefaultRule.IsNull() {
		data.SyslogDefaultRule = jsontypes.BoolValue(false)
	}

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}
