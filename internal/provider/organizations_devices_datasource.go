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

// OrganizationsDevicesDataSource struct. If not, implement them.
var _ datasource.DataSource = &OrganizationsDevicesDataSource{}

// The NewOrganizationsDevicesDataSource function is a constructor for the data source. This function needs
// to be added to the list of Data Sources in provider.go: func (p *ScaffoldingProvider) DataSources.
// If it's not added, the provider won't be aware of this data source's existence.
func NewOrganizationsDevicesDataSource() datasource.DataSource {
	return &OrganizationsDevicesDataSource{}
}

// OrganizationsDevicesDataSource struct defines the structure for this data source.
// It includes an APIClient field for making requests to the Meraki API.
// If additional fields are required (e.g., for caching or for tracking internal state), add them here.
type OrganizationsDevicesDataSource struct {
	client *openApiClient.APIClient
}

// The OrganizationsDevicesDataSourceModel structure describes the data model.
// This struct is where you define all the attributes that are part of this data source's state.
type OrganizationsDevicesDataSourceModel struct {
	Id             jsontypes.String          `tfsdk:"id"`
	OrganizationId jsontypes.String          `tfsdk:"organization_id"`
	List           []OrganizationDeicesModel `tfsdk:"list"`
}

type OrganizationDeicesModel struct {
	Name        jsontypes.String                `tfsdk:"name"`
	Lat         jsontypes.Float64               `tfsdk:"lat"`
	Lng         jsontypes.Float64               `tfsdk:"lng"`
	Address     jsontypes.String                `tfsdk:"address"`
	Notes       jsontypes.String                `tfsdk:"notes"`
	Tags        jsontypes.Set[jsontypes.String] `tfsdk:"tags"`
	NetworkId   jsontypes.String                `tfsdk:"network_id"`
	Serial      jsontypes.String                `tfsdk:"serial"`
	Model       jsontypes.String                `tfsdk:"model"`
	Mac         jsontypes.String                `tfsdk:"mac"`
	LanIp       jsontypes.String                `tfsdk:"lan_ip"`
	Firmware    jsontypes.String                `tfsdk:"firmware"`
	ProductType jsontypes.String                `tfsdk:"product_type"`
}

// Metadata provides a way to define information about the data source.
// This method is called by the framework to retrieve metadata about the data source.
func (d *OrganizationsDevicesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizations_devices"
}

// Schema provides a way to define the structure of the data source data.
// It is called by the framework to get the schema of the data source.
func (d *OrganizationsDevicesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	// The Schema object defines the structure of the data source.
	resp.Schema = schema.Schema{

		// It should provide a clear and concise description of the data source.
		MarkdownDescription: "OrganizationsDevices",

		// The Attributes map describes the fields of the data source.
		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.StringType,
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "Organization ID",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 31),
				},
				CustomType: jsontypes.StringType,
			},
			"list": schema.SetNestedAttribute{
				Optional:    true,
				Computed:    true,
				Description: "List the devices",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "Name of the device",
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
						"lat": schema.Float64Attribute{
							MarkdownDescription: "Latitude of the device",
							Optional:            true,
							CustomType:          jsontypes.Float64Type,
						},
						"lng": schema.Float64Attribute{
							MarkdownDescription: "Longitude of the device",
							Optional:            true,
							CustomType:          jsontypes.Float64Type,
						},
						"address": schema.StringAttribute{
							MarkdownDescription: "Physical address of the device",
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
						"notes": schema.StringAttribute{
							MarkdownDescription: "Notes for the device, limited to 255 characters",
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
						"tags": schema.SetAttribute{
							CustomType:  jsontypes.SetType[jsontypes.String](),
							ElementType: jsontypes.StringType,
							Description: "List of tags assigned to the device",
							Computed:    true,
							Optional:    true,
						},
						"network_id": schema.StringAttribute{
							MarkdownDescription: "ID of the network the device belongs to",
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
						"serial": schema.StringAttribute{
							MarkdownDescription: "Serial number of the device",
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
						"model": schema.StringAttribute{
							MarkdownDescription: "Model of the device",
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
						"mac": schema.StringAttribute{
							MarkdownDescription: "MAC address of the device",
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
						"lan_ip": schema.StringAttribute{
							MarkdownDescription: "LAN IP address of the device",
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
						"firmware": schema.StringAttribute{
							MarkdownDescription: "Firmware version of the device",
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
						"product_type": schema.StringAttribute{
							MarkdownDescription: "Product type of the device",
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
					},
				},
			},
		},
	}
}

// Configure is a method of the data source interface that Terraform calls to provide the configured provider instance to the data source.
// It passes the DataSourceData that's been stored by the provider's ConfigureFunc.
func (d *OrganizationsDevicesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

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
func (d *OrganizationsDevicesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OrganizationsDevicesDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	// Check if there are any errors before proceeding.
	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := d.client.OrganizationsApi.GetOrganizationDevices(ctx, data.OrganizationId.ValueString()).Execute()
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
	if err = json.NewDecoder(httpResp.Body).Decode(&data.List); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read data source")
}
