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
var _ datasource.DataSource = &NetworksSwitchStormControlDataSource{}

func NewNetworksSwitchStormControlDataSource() datasource.DataSource {
	return &NetworksSwitchStormControlDataSource{}
}

// NetworksSwitchStormControlDataSource defines the resource implementation.
type NetworksSwitchStormControlDataSource struct {
	client *openApiClient.APIClient
}

// NetworksSwitchStormControlDataSourceModel describes the resource data model.
type NetworksSwitchStormControlDataSourceModel struct {
	Id                      jsontypes.String `tfsdk:"id" json:"id"`
	NetworkId               jsontypes.String `tfsdk:"network_id" json:"network_id"`
	BroadcastThreshold      jsontypes.Int64  `tfsdk:"broadcast_threshold" json:"broadcastThreshold"`
	MulticastThreshold      jsontypes.Int64  `tfsdk:"multicast_threshold" json:"multicastThreshold"`
	UnknownUnicastThreshold jsontypes.Int64  `tfsdk:"unknown_unicast_threshold" json:"unknownUnicastThreshold"`
}

func (r *NetworksSwitchStormControlDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_storm_control"
}

func (r *NetworksSwitchStormControlDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Networks Switch Storm Control DataSource resource for updating Storm Control",
		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{
				MarkdownDescription: "Example identifier",
				Computed:            true,
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
			"broadcast_threshold": schema.Int64Attribute{
				MarkdownDescription: "Broadcast Threshold",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.Int64Type,
			},
			"multicast_threshold": schema.Int64Attribute{
				MarkdownDescription: "Multicast Threshold",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.Int64Type,
			},
			"unknown_unicast_threshold": schema.Int64Attribute{
				Description: "Unknown Unicast Threshold",
				Optional:    true,
				Computed:    true,
				CustomType:  jsontypes.Int64Type,
			},
		},
	}
}

func (r *NetworksSwitchStormControlDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *NetworksSwitchStormControlDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *NetworksSwitchStormControlDataSourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := r.client.ConfigureApi.GetNetworkSwitchStormControl(ctx, data.NetworkId.ValueString()).Execute()
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

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	} else {
		resp.Diagnostics.Append()
	}

	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON Decode issue",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
		return
	}

	data.Id = data.NetworkId

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}
