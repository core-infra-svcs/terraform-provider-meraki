package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"

	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

var _ datasource.DataSource = &OrganizationsCellularGatewayUplinkStatusesDataSource{}

func NewOrganizationsCellularGatewayUplinkStatusesDataSource() datasource.DataSource {
	return &OrganizationsCellularGatewayUplinkStatusesDataSource{}
}

type OrganizationsCellularGatewayUplinkStatusesDataSource struct {
	client *openApiClient.APIClient
}

// The OrganizationsCellularGatewayUplinkStatusesDataSourceModel structure describes the data model.
// This struct is where you define all the attributes that are part of this data source's state.
type OrganizationsCellularGatewayUplinkStatusesDataSourceModel struct {
	Id             jsontypes.String                                                `tfsdk:"id"`
	OrganizationId jsontypes.String                                                `tfsdk:"organization_id"`
	PerPage        jsontypes.Int64                                                 `tfsdk:"per_page"`
	StartingAfter  jsontypes.String                                                `tfsdk:"starting_after"`
	EndingBefore   jsontypes.String                                                `tfsdk:"ending_before"`
	NetworkIds     []jsontypes.String                                              `tfsdk:"network_ids"`
	Serials        []jsontypes.String                                              `tfsdk:"serials"`
	Iccids         []jsontypes.String                                              `tfsdk:"iccids"`
	List           []OrganizationsCellularGatewayUplinkStatusesDataSourceModelList `tfsdk:"list"`
}

type OrganizationsCellularGatewayUplinkStatusesDataSourceModelList struct {
	NetworkId      jsontypes.String `tfsdk:"network_id" json:"networkId,omitempty"`
	Serial         jsontypes.String `tfsdk:"serial"`
	Model          jsontypes.String `tfsdk:"model"`
	LastReportedAt jsontypes.String `tfsdk:"last_reported_at"`
	Uplinks        []Uplink         `tfsdk:"uplinks"`
}

type Uplink struct {
	Interface      jsontypes.String `tfsdk:"interface"`
	Status         jsontypes.String `tfsdk:"status"`
	Ip             jsontypes.String `tfsdk:"ip"`
	Provider       jsontypes.String `tfsdk:"provider"`
	PublicIp       jsontypes.String `tfsdk:"public_ip"`
	Model          jsontypes.String `tfsdk:"model"`
	SignalStat     SignalStat       `tfsdk:"signal_stat"`
	ConnectionType jsontypes.String `tfsdk:"connection_type"`
	Apn            jsontypes.String `tfsdk:"apn"`
	Gateway        jsontypes.String `tfsdk:"gateway"`
	Dns1           jsontypes.String `tfsdk:"dns1"`
	Dns2           jsontypes.String `tfsdk:"dns2"`
	SignalType     jsontypes.String `tfsdk:"signal_type"`
	Iccid          jsontypes.String `tfsdk:"iccid"`
}

type SignalStat struct {
	Rsrp jsontypes.String `tfsdk:"rsrp"`
	Rsrq jsontypes.String `tfsdk:"rsrq"`
}

// Metadata provides a way to define information about the data source.
// This method is called by the framework to retrieve metadata about the data source.
func (d *OrganizationsCellularGatewayUplinkStatusesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_organizations_cellular_gateway_uplink_statuses"
}

// Schema provides a way to define the structure of the data source data.
// It is called by the framework to get the schema of the data source.
func (d *OrganizationsCellularGatewayUplinkStatusesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	// The Schema object defines the structure of the data source.
	resp.Schema = schema.Schema{

		MarkdownDescription: "OrganizationsCellularGatewayUplinkStatuses",

		// The Attributes map describes the fields of the data source.
		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.StringType,
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "Organization Id",
				CustomType:          jsontypes.StringType,
				Required:            true,
			},
			"per_page": schema.Int64Attribute{
				MarkdownDescription: "The number of entries per page returned. Acceptable range is 3 - 1000. Default is 1000.",
				CustomType:          jsontypes.Int64Type,
				Optional:            true,
				Computed:            true,
			},
			"starting_after": schema.StringAttribute{
				MarkdownDescription: "A token used by the server to indicate the start of the page. Often this is a timestamp or an ID but it is not limited to those. This parameter should not be defined by client applications. The link for the first, last, prev, or next page in the HTTP Link header should define it.",
				CustomType:          jsontypes.StringType,
				Optional:            true,
				Computed:            true,
			},
			"ending_before": schema.StringAttribute{
				MarkdownDescription: "A token used by the server to indicate the end of the page. Often this is a timestamp or an ID but it is not limited to those. This parameter should not be defined by client applications. The link for the first, last, prev, or next page in the HTTP Link header should define it.",
				CustomType:          jsontypes.StringType,
				Optional:            true,
				Computed:            true,
			},
			"serials": schema.SetAttribute{
				MarkdownDescription: "A list of serial numbers. The returned devices will be filtered to only include these serials.",
				ElementType:         jsontypes.StringType,
				CustomType:          jsontypes.SetType[jsontypes.String](),
				Computed:            true,
				Optional:            true,
			},
			"network_ids": schema.SetAttribute{
				MarkdownDescription: "A list of network IDs. The returned devices will be filtered to only include these networks.",
				ElementType:         jsontypes.StringType,
				CustomType:          jsontypes.SetType[jsontypes.String](),
				Computed:            true,
				Optional:            true,
			},
			"iccids": schema.SetAttribute{
				MarkdownDescription: "A list of ICCIDs. The returned devices will be filtered to only include these ICCIDs.",
				ElementType:         jsontypes.StringType,
				CustomType:          jsontypes.SetType[jsontypes.String](),
				Computed:            true,
				Optional:            true,
			},
			"list": schema.SetNestedAttribute{
				MarkdownDescription: "List of organization acls",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"network_id": schema.StringAttribute{
							MarkdownDescription: "ID of the network.",
							CustomType:          jsontypes.StringType,
							Optional:            true,
							Computed:            true,
						},
						"serial": schema.StringAttribute{
							MarkdownDescription: "Serial number of the device.",
							CustomType:          jsontypes.StringType,
							Optional:            true,
							Computed:            true,
						},
						"model": schema.StringAttribute{
							MarkdownDescription: "Device model.",
							CustomType:          jsontypes.StringType,
							Optional:            true,
							Computed:            true,
						},
						"last_reported_at": schema.StringAttribute{
							MarkdownDescription: "Last reported time for the device.",
							CustomType:          jsontypes.StringType,
							Optional:            true,
							Computed:            true,
						},
						"uplinks": schema.ListNestedAttribute{
							MarkdownDescription: "Uplinks info.",
							Optional:            true,
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"interface": schema.StringAttribute{
										MarkdownDescription: "Uplink interface.",
										CustomType:          jsontypes.StringType,
										Optional:            true,
										Computed:            true,
									},
									"status": schema.StringAttribute{
										MarkdownDescription: "Uplink status.",
										CustomType:          jsontypes.StringType,
										Optional:            true,
										Computed:            true,
									},
									"ip": schema.StringAttribute{
										MarkdownDescription: "Uplink ip.",
										CustomType:          jsontypes.StringType,
										Optional:            true,
										Computed:            true,
									},
									"provider": schema.StringAttribute{
										MarkdownDescription: "Network Provider.",
										CustomType:          jsontypes.StringType,
										Optional:            true,
										Computed:            true,
									},
									"public_ip": schema.StringAttribute{
										MarkdownDescription: "Public IP.",
										CustomType:          jsontypes.StringType,
										Optional:            true,
										Computed:            true,
									},
									"model": schema.StringAttribute{
										MarkdownDescription: "Uplink model.",
										CustomType:          jsontypes.StringType,
										Optional:            true,
										Computed:            true,
									},
									"signal_stat": schema.SingleNestedAttribute{
										Optional: true,
										Attributes: map[string]schema.Attribute{
											"rsrp": schema.StringAttribute{
												MarkdownDescription: "Reference Signal Received Power.",
												CustomType:          jsontypes.StringType,
												Optional:            true,
												Computed:            true,
											},
											"rsrq": schema.StringAttribute{
												MarkdownDescription: "Reference Signal Received Quality.",
												CustomType:          jsontypes.StringType,
												Optional:            true,
												Computed:            true,
											},
										},
									},
									"connection_type": schema.StringAttribute{
										MarkdownDescription: "Connection Type.",
										CustomType:          jsontypes.StringType,
										Optional:            true,
										Computed:            true,
									},
									"apn": schema.StringAttribute{
										MarkdownDescription: "Access Point Name.",
										CustomType:          jsontypes.StringType,
										Optional:            true,
										Computed:            true,
									},
									"gateway": schema.StringAttribute{
										MarkdownDescription: "Gateway IP.",
										CustomType:          jsontypes.StringType,
										Optional:            true,
										Computed:            true,
									},
									"dns1": schema.StringAttribute{
										MarkdownDescription: "Primary DNS IP.",
										CustomType:          jsontypes.StringType,
										Optional:            true,
										Computed:            true,
									},
									"dns2": schema.StringAttribute{
										MarkdownDescription: "Secondary DNS IP.",
										CustomType:          jsontypes.StringType,
										Optional:            true,
										Computed:            true,
									},
									"signal_type": schema.StringAttribute{
										MarkdownDescription: "Signal Type.",
										CustomType:          jsontypes.StringType,
										Optional:            true,
										Computed:            true,
									},
									"iccid": schema.StringAttribute{
										MarkdownDescription: "Integrated Circuit Card Identification SsidNumber.",
										CustomType:          jsontypes.StringType,
										Optional:            true,
										Computed:            true,
									},
								},
							},
						},
					}}},
		},
	}
}

// Configure is a method of the data source interface that Terraform calls to provide the configured provider instance to the data source.
// It passes the DataSourceData that's been stored by the provider's ConfigureFunc.
func (d *OrganizationsCellularGatewayUplinkStatusesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

	// The provider must be properly configured before it can be used.
	if req.ProviderData == nil {
		return
	}

	// Here we expect the provider data to be of type *openApiClient.APIClient.
	client, ok := req.ProviderData.(*openApiClient.APIClient)

	// This is a fatal error and the provider cannot proceed without it.
	// If you see this error, it means there is an issue with the provider setup.
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	// This allows the data source to use the configured provider for any API calls it needs to make.
	d.client = client
}

// Read method is responsible for reading an existing data source's state.
func (d *OrganizationsCellularGatewayUplinkStatusesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OrganizationsCellularGatewayUplinkStatusesDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	// Check if there are any errors before proceeding.
	if resp.Diagnostics.HasError() {
		return
	}

	request := d.client.CellularGatewayApi.GetOrganizationCellularGatewayUplinkStatuses(context.Background(), data.OrganizationId.ValueString())

	if !data.PerPage.IsUnknown() {
		request.PerPage(int32(data.PerPage.ValueInt64()))
	}
	if !data.StartingAfter.IsUnknown() {
		request.StartingAfter(data.StartingAfter.ValueString())
	}
	if !data.EndingBefore.IsUnknown() {
		request.EndingBefore(data.EndingBefore.ValueString())
	}

	if len(data.Serials) > 0 {
		var serials []string
		for _, serial := range data.Serials {
			serials = append(serials, serial.String())
		}
		request.Serials(serials)
	}

	if len(data.NetworkIds) > 0 {
		var networkIds []string
		for _, networkId := range data.NetworkIds {
			networkIds = append(networkIds, networkId.String())
		}
		request.NetworkIds(networkIds)
	}

	if len(data.Iccids) > 0 {
		var iccids []string
		for _, iccid := range data.Iccids {
			iccids = append(iccids, iccid.String())
		}
		request.Iccids(iccids)
	}

	_, httpResp, err := request.Execute()

	// If there was an error during API call, add it to diagnostics.
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
	}

	// If it's not what you expect, add an error to diagnostics.
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// If there were any errors up to this point, log the state data and return.
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("State Data", fmt.Sprintf("\n%v", data))
		return
	}

	// save into the Terraform state.
	if err = json.NewDecoder(httpResp.Body).Decode(&data.List); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	// Set ID for the data source.
	data.Id = jsontypes.StringValue("example-id")

	// Now set the final state of the data source.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read data source")
}
