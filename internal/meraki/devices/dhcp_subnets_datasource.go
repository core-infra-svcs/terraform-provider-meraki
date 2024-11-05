package devices

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// The below var(s) ensures the provider defined types fully satisfy the required interface(s) for a data source.
// DevicesApplianceDhcpSubnetsDataSource struct. If not, implement them.
var _ datasource.DataSource = &DevicesApplianceDhcpSubnetsDataSource{}

// The NewDevicesApplianceDhcpSubnetsDataSource function is a constructor for the data source. This function needs
// to be added to the list of Data Sources in provider.go: func (p *ScaffoldingProvider) DataSources.
// If it's not added, the provider won't be aware of this data source's existence.
func NewDevicesApplianceDhcpSubnetsDataSource() datasource.DataSource {
	return &DevicesApplianceDhcpSubnetsDataSource{}
}

// DevicesApplianceDhcpSubnetsDataSource struct defines the structure for this data source.
// It includes an APIClient field for making requests to the Meraki API.
// If additional fields are required (e.g., for caching or for tracking internal state), add them here.
type DevicesApplianceDhcpSubnetsDataSource struct {
	client *openApiClient.APIClient
}

// The DevicesApplianceDhcpSubnetsDataSourceModel structure describes the data model.
// This struct is where you define all the attributes that are part of this data source's state.
type DevicesApplianceDhcpSubnetsDataSourceModel struct {

	// The Id field is mandatory for all data sources. It's used for data source identification and is required
	// for the acceptance tests to run.
	Id     jsontypes.String                                 `tfsdk:"id"`
	Serial jsontypes.String                                 `tfsdk:"serial"`
	List   []DevicesApplianceDhcpSubnetsDataSourceModelList `tfsdk:"list"`

	// Each of the remaining fields represents an attribute of this data source. They should match the attributes
	// defined in the tfsdk.Schema for this data source.
}

type DevicesApplianceDhcpSubnetsDataSourceModelList struct {
	Subnet    jsontypes.String `tfsdk:"subnet" json:"subnet"`
	VlanId    jsontypes.Int64  `tfsdk:"vlan_id" json:"vlanId"`
	UsedCount jsontypes.Int64  `tfsdk:"used_count" json:"usedCount"`
	FreeCount jsontypes.Int64  `tfsdk:"free_count" json:"freeCount"`
}

// Metadata provides a way to define information about the data source.
// This method is called by the framework to retrieve metadata about the data source.
func (d *DevicesApplianceDhcpSubnetsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_devices_appliance_dhcp_subnets"
}

// Schema provides a way to define the structure of the data source data.
// It is called by the framework to get the schema of the data source.
func (d *DevicesApplianceDhcpSubnetsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	// The Schema object defines the structure of the data source.
	resp.Schema = schema.Schema{
		// It should provide a clear and concise description of the data source.
		MarkdownDescription: "Return the DHCP subnet information for an appliance",

		// The Attributes map describes the fields of the data source.
		Attributes: map[string]schema.Attribute{

			// Every data source must have an ID attribute. This is computed by the framework.
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.StringType,
			},
			"serial": schema.StringAttribute{
				Required:   true,
				CustomType: jsontypes.StringType,
			},
			"list": schema.ListNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "List of DHCP subnets",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"subnet": schema.StringAttribute{
							MarkdownDescription: "Network ID",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.StringType,
						},
						"vlan_id": schema.Int64Attribute{
							MarkdownDescription: "VLAN ID",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.Int64Type,
						},
						"used_count": schema.Int64Attribute{
							MarkdownDescription: "used count",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.Int64Type,
						},
						"free_count": schema.Int64Attribute{
							MarkdownDescription: "free count",
							Optional:            true,
							Computed:            true,
							CustomType:          jsontypes.Int64Type,
						},
					},
				},
			},
		},
	}
}

// Configure is a method of the data source interface that Terraform calls to provide the configured provider instance to the data source.
// It passes the DataSourceData that's been stored by the provider's ConfigureFunc.
func (d *DevicesApplianceDhcpSubnetsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

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
func (d *DevicesApplianceDhcpSubnetsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data DevicesApplianceDhcpSubnetsDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	// Check if there are any errors before proceeding.
	if resp.Diagnostics.HasError() {
		return
	}

	// Remember to handle any potential errors.
	_, httpResp, err := d.client.SubnetsApi.GetDeviceApplianceDhcpSubnets(ctx, data.Serial.ValueString()).Execute()

	// If it's not what you expect, add an error to diagnostics.
	if httpResp.StatusCode == 404 {
		resp.Diagnostics.AddWarning(
			"Device likely does not support DHCP subnets",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	} else if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)

		// HTTP 400 counts as an error so moving this here
		// If there was an error during API call, add it to diagnostics.
		if err != nil {
			resp.Diagnostics.AddError(
				"HTTP Client Failure",
				utils.HttpDiagnostics(httpResp),
			)
			return
		}

		// If there were any errors up to this point, log the state data and return.
		if resp.Diagnostics.HasError() {
			resp.Diagnostics.AddError("State Data", fmt.Sprintf("\n%v", data))
			return
		}

		// Decode the HTTP response body into your data model.
		// If there's an error, add it to diagnostics.
		if err = json.NewDecoder(httpResp.Body).Decode(&data.List); err != nil {
			resp.Diagnostics.AddError(
				"JSON decoding error",
				fmt.Sprintf("%v\n", err.Error()),
			)
			return
		}

	}

	// Set ID for the data source.
	data.Id = jsontypes.StringValue("example-id")

	// Now set the final state of the data source.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read data source")
}
