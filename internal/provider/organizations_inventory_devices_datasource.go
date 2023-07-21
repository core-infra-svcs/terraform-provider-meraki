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

// If necessary, modify the existing functions to better organize the code and improve readability.

// 1. Meraki OpenAPI Specification: This provides the complete API documentation in the OpenAPI format. It can be used to understand the API endpoints, parameters, and responses.
// URL: https://raw.githubusercontent.com/meraki/openapi/master/openapi/spec2.json
//
// 2. Meraki Postman Documentation: This is a collection of Postman requests that demonstrate how to interact with the Meraki API. It can be useful for testing and understanding how to structure API requests.
// URL: https://documenter.getpostman.com/view/897512/SzYXYfmJ#intro
//
// 3. Dashboard-api-go Documentation: This is the Go client library for the Meraki API. It provides Go types and functions that make it easier to call the API from Go code.
// URL: https://github.com/meraki/dashboard-api-go/tree/main/client/docs

// The below var(s) ensures the provider defined types fully satisfy the required interface(s) for a data source.
// OrganizationsInventoryDevicesDataSource struct. If not, implement them.
var _ datasource.DataSource = &OrganizationsInventoryDevicesDataSource{}

// The NewOrganizationsInventoryDevicesDataSource function is a constructor for the data source. This function needs
// to be added to the list of Data Sources in provider.go: func (p *ScaffoldingProvider) DataSources.
// TODO: Ensure this function is added to the list of data sources in provider.go.
// If it's not added, the provider won't be aware of this data source's existence.
func NewOrganizationsInventoryDevicesDataSource() datasource.DataSource {
	return &OrganizationsInventoryDevicesDataSource{}
}

type InventoryDevice struct {
	Mac               jsontypes.String                `tfsdk:"mac" json:"mac"`
	Serial            jsontypes.String                `tfsdk:"serial" json:"serial"`
	Name              jsontypes.String                `tfsdk:"name" json:"name"`
	Model             jsontypes.String                `tfsdk:"model" json:"model"`
	NetworkId         jsontypes.String                `tfsdk:"network_id" json:"networkId"`
	OrderNumber       jsontypes.String                `tfsdk:"order_number" json:"orderNumber"`
	ClaimedAt         jsontypes.String                `tfsdk:"claimed_at" json:"claimedAt"`
	Tags              jsontypes.Set[jsontypes.String] `tfsdk:"tags" json:"tags"`
	ProductType       jsontypes.String                `tfsdk:"product_type" json:"productType"`
	LicenseExpiryDate jsontypes.String                `tfsdk:"license_expiry_date" json:"licenseExpiryDate"`
}

// OrganizationsInventoryDevicesDataSource struct defines the structure for this data source.
// It includes an APIClient field for making requests to the Meraki API.
// TODO: Ensure this structure includes all necessary fields to represent the state of this data source.
// If additional fields are required (e.g., for caching or for tracking internal state), add them here.
type OrganizationsInventoryDevicesDataSource struct {
	client *openApiClient.APIClient
}

// The OrganizationsInventoryDevicesDataSourceModel structure describes the data model.
// This struct is where you define all the attributes that are part of this data source's state.
type OrganizationsInventoryDevicesDataSourceModel struct {
	Id             jsontypes.String                `tfsdk:"id"`
	OrganizationID jsontypes.String                `tfsdk:"organization_id"`
	List           []InventoryDevice               `tfsdk:"list"`
	PerPage        jsontypes.Int64                 `tfsdk:"per_page"`
	StartingAfter  jsontypes.String                `tfsdk:"starting_after"`
	EndingBefore   jsontypes.String                `tfsdk:"ending_before"`
	UsedState      jsontypes.String                `tfsdk:"used_state"`
	Search         jsontypes.String                `tfsdk:"search"`
	TagsFilterType jsontypes.String                `tfsdk:"tags_filter_type"`
	Macs           jsontypes.Set[jsontypes.String] `tfsdk:"macs"`
	NetworkIDs     jsontypes.Set[jsontypes.String] `tfsdk:"network_ids"`
	Serials        jsontypes.Set[jsontypes.String] `tfsdk:"serials"`
	Models         jsontypes.Set[jsontypes.String] `tfsdk:"models"`
	OrderNumbers   jsontypes.Set[jsontypes.String] `tfsdk:"order_numbers"`
	Tags           jsontypes.Set[jsontypes.String] `tfsdk:"tags"`
	ProductTypes   jsontypes.Set[jsontypes.String] `tfsdk:"product_types"`
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
				Computed:   true,
				Optional:   true,
				CustomType: jsontypes.StringType,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 31),
				},
			},
			"per_page": schema.StringAttribute{
				Computed:   true,
				Optional:   true,
				CustomType: jsontypes.StringType,
			},
			"starting_after": schema.StringAttribute{
				Computed:   true,
				Optional:   true,
				CustomType: jsontypes.StringType,
			},
			"ending_before": schema.StringAttribute{
				Computed:   true,
				Optional:   true,
				CustomType: jsontypes.StringType,
			},
			"used_state": schema.StringAttribute{
				Computed:   true,
				Optional:   true,
				CustomType: jsontypes.StringType,
			},
			"search": schema.StringAttribute{
				Computed:   true,
				Optional:   true,
				CustomType: jsontypes.StringType,
			},
			"tags_filter_type": schema.StringAttribute{
				Computed:   true,
				Optional:   true,
				CustomType: jsontypes.StringType,
			},
			"macs": schema.SetAttribute{
				Computed:            true,
				CustomType:          jsontypes.SetType[jsontypes.String](),
				ElementType:         jsontypes.StringType,
				Optional:            true,
				MarkdownDescription: "Device tags",
			},
			"network_ids": schema.SetAttribute{
				Computed:            true,
				CustomType:          jsontypes.SetType[jsontypes.String](),
				ElementType:         jsontypes.StringType,
				Optional:            true,
				MarkdownDescription: "Device tags",
			},
			"serials": schema.SetAttribute{
				Computed:            true,
				CustomType:          jsontypes.SetType[jsontypes.String](),
				ElementType:         jsontypes.StringType,
				Optional:            true,
				MarkdownDescription: "Device tags",
			},
			"models": schema.SetAttribute{
				Computed:            true,
				CustomType:          jsontypes.SetType[jsontypes.String](),
				ElementType:         jsontypes.StringType,
				Optional:            true,
				MarkdownDescription: "Device tags",
			},
			"order_numbers": schema.SetAttribute{
				Computed:            true,
				CustomType:          jsontypes.SetType[jsontypes.String](),
				ElementType:         jsontypes.StringType,
				Optional:            true,
				MarkdownDescription: "Device tags",
			},
			"tags": schema.SetAttribute{
				Computed:            true,
				CustomType:          jsontypes.SetType[jsontypes.String](),
				ElementType:         jsontypes.StringType,
				Optional:            true,
				MarkdownDescription: "Device tags",
			},
			"product_types": schema.SetAttribute{
				Computed:            true,
				CustomType:          jsontypes.SetType[jsontypes.String](),
				ElementType:         jsontypes.StringType,
				Optional:            true,
				MarkdownDescription: "Device tags",
			},
			"list": schema.SetNestedAttribute{
				Optional:    true,
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
						"license_expiry_date": schema.StringAttribute{
							Computed:            true,
							CustomType:          jsontypes.StringType,
							Optional:            true,
							MarkdownDescription: "License expiry date of the device",
						},
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
	devices := d.client.InventoryApi.GetOrganizationInventoryDevices(ctx, data.OrganizationID.ValueString())
	if !data.PerPage.IsUnknown() {
		devices.PerPage(int32(data.PerPage.ValueInt64()))
	}
	if !data.StartingAfter.IsUnknown() {
		devices.StartingAfter(data.StartingAfter.ValueString())
	}
	if !data.EndingBefore.IsUnknown() {
		devices.EndingBefore(data.EndingBefore.ValueString())
	}
	if !data.UsedState.IsUnknown() {
		devices.UsedState(data.UsedState.ValueString())
	}
	if !data.Search.IsUnknown() {
		devices.Search(data.Search.ValueString())
	}
	if !data.TagsFilterType.IsUnknown() {
		devices.TagsFilterType(data.TagsFilterType.ValueString())
	}
	if !data.Macs.IsUnknown() {
		var macs []string
		for _, value := range data.Macs.Elements() {
			macs = append(macs, value.String())
		}
		devices.Macs(macs)
	}
	if !data.NetworkIDs.IsUnknown() {
		var networkIds []string
		for _, value := range data.NetworkIDs.Elements() {
			networkIds = append(networkIds, value.String())
		}
		devices.NetworkIds(networkIds)
	}
	if !data.Serials.IsUnknown() {
		var serials []string
		for _, value := range data.Serials.Elements() {
			serials = append(serials, value.String())
		}
		devices.Serials(serials)
	}
	if !data.Models.IsUnknown() {
		var models []string
		for _, value := range data.Models.Elements() {
			models = append(models, value.String())
		}
		devices.Models(models)
	}
	if !data.OrderNumbers.IsUnknown() {
		var orderNumbers []string
		for _, value := range data.OrderNumbers.Elements() {
			orderNumbers = append(orderNumbers, value.String())
		}
		devices.OrderNumbers(orderNumbers)
	}
	if !data.Tags.IsUnknown() {
		var tags []string
		for _, value := range data.Tags.Elements() {
			tags = append(tags, value.String())
		}
		devices.Tags(tags)
	}
	if !data.ProductTypes.IsUnknown() {
		var productTypes []string
		for _, value := range data.ProductTypes.Elements() {
			productTypes = append(productTypes, value.String())
		}
		devices.ProductTypes(productTypes)
	}

	_, httpResp, err := devices.Execute()
	// If there was an error during API call, add it to diagnostics.
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read data source",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// Collect any HTTP diagnostics that might be useful for debugging.
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
	}

	// If it's not what you expect, add an error to diagnostics.
	// TODO: Check the HTTP response status code matches the API endpoint.
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

	// Decode the HTTP response body into your data model.
	// If there's an error, add it to diagnostics.
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
