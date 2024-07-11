package organizations

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// OrganizationsInventoryDevicesDataSource struct. If not, implement them.
var _ datasource.DataSource = &OrganizationsInventoryDevicesDataSource{}

// The NewOrganizationsInventoryDevicesDataSource function is a constructor for the data source. This function needs
func NewOrganizationsInventoryDevicesDataSource() datasource.DataSource {
	return &OrganizationsInventoryDevicesDataSource{}
}

// OrganizationsInventoryDevicesDataSource struct defines the structure for this data source.
type OrganizationsInventoryDevicesDataSource struct {
	client *openApiClient.APIClient
}

// The OrganizationsInventoryDevicesDataSourceModel structure describes the data model.
type OrganizationsInventoryDevicesDataSourceModel struct {
	Id             jsontypes.String                                              `tfsdk:"id"`
	OrganizationID jsontypes.String                                              `tfsdk:"organization_id"`
	List           []OrganizationsInventoryDevicesDataSourceModelInventoryDevice `tfsdk:"list"`
}

type OrganizationsInventoryDevicesDataSourceModelInventoryDevice struct {
	Mac               jsontypes.String                                                     `tfsdk:"mac" json:"mac"`
	Serial            jsontypes.String                                                     `tfsdk:"serial" json:"serial"`
	Name              jsontypes.String                                                     `tfsdk:"name" json:"name"`
	Model             jsontypes.String                                                     `tfsdk:"model" json:"model"`
	NetworkId         jsontypes.String                                                     `tfsdk:"network_id" json:"networkId"`
	OrderNumber       jsontypes.String                                                     `tfsdk:"order_number" json:"orderNumber"`
	ClaimedAt         jsontypes.String                                                     `tfsdk:"claimed_at" json:"claimedAt"`
	LicenseExpiryDate jsontypes.String                                                     `tfsdk:"license_expiry_date" json:"licenseExpiryDate"`
	Tags              []jsontypes.String                                                   `tfsdk:"tags" json:"tags"`
	ProductType       jsontypes.String                                                     `tfsdk:"product_type" json:"productType"`
	CountryCode       jsontypes.String                                                     `tfsdk:"country_code" json:"countryCode"`
	Details           []OrganizationsInventoryDevicesDataSourceModelInventoryDeviceDetails `tfsdk:"details" json:"details"`
}

type OrganizationsInventoryDevicesDataSourceModelInventoryDeviceDetails struct {
	Name  jsontypes.String `tfsdk:"name" json:"name"`
	Value jsontypes.String `tfsdk:"value" json:"value"`
}

// Metadata provides a way to define information about the data source.
// This method is called by the framework to retrieve metadata about the data source.
func (d *OrganizationsInventoryDevicesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_organizations_inventory_devices"
}

// Schema provides a way to define the structure of the data source data.
// It is called by the framework to get the schema of the data source.
func (d *OrganizationsInventoryDevicesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	// The Schema object defines the structure of the data source.
	resp.Schema = schema.Schema{
		MarkdownDescription: "OrganizationsInventoryDevices",
		// The Attributes map describes the fields of the data source.
		Attributes: map[string]schema.Attribute{
			// Every data source must have an ID attribute. This is computed by the framework.
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.StringType,
			},
			"organization_id": schema.StringAttribute{
				Required:   true,
				CustomType: jsontypes.StringType,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 31),
				},
			},
			"list": schema.SetNestedAttribute{
				Computed:    true,
				Description: "",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"mac": schema.StringAttribute{
							Computed:            true,
							CustomType:          jsontypes.StringType,
							Optional:            true,
							MarkdownDescription: "MAC address of the device",
						},
						"serial": schema.StringAttribute{
							Computed:            true,
							CustomType:          jsontypes.StringType,
							Optional:            true,
							MarkdownDescription: "serial number of the devicee",
						},
						"name": schema.StringAttribute{
							Computed:            true,
							CustomType:          jsontypes.StringType,
							Optional:            true,
							MarkdownDescription: "Name of the device",
						},
						"model": schema.StringAttribute{
							Computed:            true,
							CustomType:          jsontypes.StringType,
							Optional:            true,
							MarkdownDescription: "Model type of the device",
						},
						"network_id": schema.StringAttribute{
							Computed:            true,
							CustomType:          jsontypes.StringType,
							Optional:            true,
							MarkdownDescription: "Network Id of the device",
						},
						"order_number": schema.StringAttribute{
							Computed:            true,
							CustomType:          jsontypes.StringType,
							Optional:            true,
							MarkdownDescription: "Order number of the device",
						},
						"claimed_at": schema.StringAttribute{
							Computed:            true,
							CustomType:          jsontypes.StringType,
							Optional:            true,
							MarkdownDescription: "Claimed time of the device",
						},
						"license_expiry_date": schema.StringAttribute{
							Computed:            true,
							CustomType:          jsontypes.StringType,
							Optional:            true,
							MarkdownDescription: "License expiry date of the device",
						},
						"tags": schema.SetAttribute{
							Computed:            true,
							CustomType:          jsontypes.SetType[jsontypes.String](),
							ElementType:         jsontypes.StringType,
							Optional:            true,
							MarkdownDescription: "Device tags",
						},
						"product_type": schema.StringAttribute{
							Computed:            true,
							CustomType:          jsontypes.StringType,
							Optional:            true,
							MarkdownDescription: "Product type of the device",
						},
						"country_code": schema.StringAttribute{
							Computed:            true,
							CustomType:          jsontypes.StringType,
							Optional:            true,
							MarkdownDescription: "",
						},
						"details": schema.SetNestedAttribute{
							Computed:    true,
							Description: "",
							NestedObject: schema.NestedAttributeObject{
								Attributes: map[string]schema.Attribute{
									"name": schema.StringAttribute{
										Computed:            true,
										CustomType:          jsontypes.StringType,
										MarkdownDescription: "",
									},
									"value": schema.StringAttribute{
										Computed:            true,
										CustomType:          jsontypes.StringType,
										MarkdownDescription: "",
									},
								},
							}},
					},
				},
			},
		},
	}
}

// Configure is a method of the data source interface that Terraform calls to provide the configured provider instance to the data source.
// It passes the DataSourceData that's been stored by the provider's ConfigureFunc.
func (d *OrganizationsInventoryDevicesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

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
func (d *OrganizationsInventoryDevicesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OrganizationsInventoryDevicesDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	// Check if there are any errors before proceeding.
	if resp.Diagnostics.HasError() {
		return
	}

	// Remember to handle any potential errors.
	inlineResp, httpResp, err := d.client.InventoryApi.GetOrganizationInventoryDevices(ctx, data.OrganizationID.ValueString()).PerPage(1000).Execute()
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
	}

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("State Data", fmt.Sprintf("\n%s", data))
		return
	}

	data.Id = jsontypes.StringValue("example-id")
	for _, inlineRespDevice := range inlineResp {

		var deviceData OrganizationsInventoryDevicesDataSourceModelInventoryDevice

		deviceData.Mac = jsontypes.StringValue(inlineRespDevice.GetMac())
		deviceData.Serial = jsontypes.StringValue(inlineRespDevice.GetSerial())
		deviceData.Name = jsontypes.StringValue(inlineRespDevice.GetName())
		deviceData.Model = jsontypes.StringValue(inlineRespDevice.GetModel())
		deviceData.NetworkId = jsontypes.StringValue(inlineRespDevice.GetNetworkId())
		deviceData.OrderNumber = jsontypes.StringValue(inlineRespDevice.GetOrderNumber())
		deviceData.ClaimedAt = jsontypes.StringValue(inlineRespDevice.GetClaimedAt().String())
		deviceData.LicenseExpiryDate = jsontypes.StringValue(inlineRespDevice.GetLicenseExpirationDate().String())

		var tags []jsontypes.String
		for _, tag := range inlineRespDevice.Tags {
			tags = append(tags, jsontypes.StringValue(tag))
		}
		deviceData.Tags = tags

		deviceData.ProductType = jsontypes.StringValue(inlineRespDevice.GetProductType())
		deviceData.CountryCode = jsontypes.StringValue(inlineRespDevice.GetCountryCode())

		data.List = append(data.List, deviceData)
	}

	// Now set the final state of the data source.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read data source")
}
