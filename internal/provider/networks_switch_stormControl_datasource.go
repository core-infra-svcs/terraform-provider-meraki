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
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &NetworksSwitchStormcontrolDataSource{}

func NewNetworksSwitchStormcontrolDataSource() datasource.DataSource {
	return &NetworksSwitchStormcontrolDataSource{}
}

// NetworksSwitchStormcontrolDataSource defines the data source implementation.
type NetworksSwitchStormcontrolDataSource struct {
	client *openApiClient.APIClient
}

// NetworksSwitchStormcontrolDataSourceModel describes the data source data model.
type NetworksSwitchStormcontrolDataSourceModel struct {
	Id        jsontypes.String `tfsdk:"id"`
	NetworkId jsontypes.String `tfsdk:"network_id"`

	BroadcastThreshold      jsontypes.Int64 `tfsdk:"broadcast_threshold" json:"broadcastThreshold"`
	MulticastThreshold      jsontypes.Int64 `tfsdk:"multicast_threshold" json:"multicastThreshold"`
	UnknownUnicastThreshold jsontypes.Int64 `tfsdk:"unknown_unicast_threshold" json:"unknownUnicastThreshold"`
}

func (d *NetworksSwitchStormcontrolDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_switch_storm_control"
}

func (d *NetworksSwitchStormcontrolDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "NetworksSwitchStormcontrol",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.StringType,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Networkd ID",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 31),
				},
				CustomType: jsontypes.StringType,
			},
			"broadcast_threshold": schema.Int64Attribute{
				MarkdownDescription: "Broadcast Threshold",
				Description:         "Percentage (1 to 99) of total available port bandwidth for broadcast traffic type. Default value 100 percent rate is to clear the configuration.",
				Optional:            true,
				CustomType:          jsontypes.Int64Type,
			},
			"multicast_threshold": schema.Int64Attribute{
				MarkdownDescription: "Multicast Threshold",
				Description:         "Percentage (1 to 99) of total available port bandwidth for multicast traffic type. Default value 100 percent rate is to clear the configuration.",
				Optional:            true,
				CustomType:          jsontypes.Int64Type,
			},
			"unknown_unicast_threshold": schema.Int64Attribute{
				MarkdownDescription: "Unknown Unicast Threshold",
				Description:         "Percentage (1 to 99) of total available port bandwidth for unknown unicast (dlf-destination lookup failure) traffic type. Default value 100 percent rate is to clear the configuration.",
				Optional:            true,
				CustomType:          jsontypes.Int64Type,
			},
		},
	}
}

func (d *NetworksSwitchStormcontrolDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (d *NetworksSwitchStormcontrolDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data NetworksSwitchStormcontrolDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := d.client.SwitchApi.GetNetworkSwitchStormControl(context.Background(), data.NetworkId.ValueString()).Execute()
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
		resp.Diagnostics.AddError("State Data", fmt.Sprintf("\n%s", data))
		return
	}

	// Save data into Terraform state
	data.Id = jsontypes.StringValue("example-id")
	if err = json.NewDecoder(httpResp.Body).Decode(&data); err != nil {
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
