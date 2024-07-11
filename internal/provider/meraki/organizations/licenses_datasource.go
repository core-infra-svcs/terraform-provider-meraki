package organizations

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

var _ datasource.DataSource = &OrganizationsLicensesDataSource{}

func NewOrganizationsLicensesDataSource() datasource.DataSource {
	return &OrganizationsLicensesDataSource{}
}

// OrganizationsLicensesDataSource struct defines the structure for this data source.
// It includes an APIClient field for making requests to the Meraki API.
type OrganizationsLicensesDataSource struct {
	client *openApiClient.APIClient
}

// The OrganizationsLicensesDataSourceModel structure describes the data model.
// This struct is where you define all the attributes that are part of this data source's state.
type OrganizationsLicensesDataSourceModel struct {
	Id             jsontypes.String                           `tfsdk:"id"`
	OrganizationId jsontypes.String                           `tfsdk:"organization_id"`
	PerPage        jsontypes.Int64                            `tfsdk:"per_page"`
	StartingAfter  jsontypes.String                           `tfsdk:"starting_after"`
	EndingBefore   jsontypes.String                           `tfsdk:"ending_before"`
	DeviceSerial   jsontypes.String                           `tfsdk:"device_serial"`
	NetworkId      jsontypes.String                           `tfsdk:"network_id"`
	State          jsontypes.String                           `tfsdk:"state"`
	List           []OrganizationsLicensesDataSourceModelList `tfsdk:"list"`
}

type OrganizationsLicensesDataSourceModelList struct {
	Id                        jsontypes.String                                                                      `tfsdk:"id"`
	LicenseType               jsontypes.String                                                                      `tfsdk:"license_type"`
	LicenseKey                jsontypes.String                                                                      `tfsdk:"license_key"`
	OrderNumber               jsontypes.String                                                                      `tfsdk:"order_number"`
	DeviceSerial              jsontypes.String                                                                      `tfsdk:"device_serial"`
	NetworkId                 jsontypes.String                                                                      `tfsdk:"network_id"`
	State                     jsontypes.String                                                                      `tfsdk:"state"`
	SeatCount                 jsontypes.Int64                                                                       `tfsdk:"seat_count"`
	TotalDurationInDays       jsontypes.Int64                                                                       `tfsdk:"total_duration_in_days"`
	DurationInDays            jsontypes.Int64                                                                       `tfsdk:"duration_in_days"`
	PermanentlyQueuedLicenses []openApiClient.GetOrganizationLicenses200ResponseInnerPermanentlyQueuedLicensesInner `tfsdk:"permanently_queued_licenses"`
	ClaimDate                 jsontypes.String                                                                      `tfsdk:"claim_date"`
	ActivationDate            jsontypes.String                                                                      `tfsdk:"activation_date"`
	ExpirationDate            jsontypes.String                                                                      `tfsdk:"expiration_date"`
	HeadLicenseId             jsontypes.String                                                                      `tfsdk:"head_license_id"`
}

// Metadata provides a way to define information about the data source.
// This method is called by the framework to retrieve metadata about the data source.
func (d *OrganizationsLicensesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_organizations_licenses"
}

// Schema provides a way to define the structure of the data source data.
// It is called by the framework to get the schema of the data source.
func (d *OrganizationsLicensesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	// The Schema object defines the structure of the data source.
	resp.Schema = schema.Schema{

		// It should provide a clear and concise description of the data source.
		MarkdownDescription: "List Organization Licenses",

		// The Attributes map describes the fields of the data source.
		Attributes: map[string]schema.Attribute{

			// Every data source must have an ID attribute. This is computed by the framework.
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
			"device_serial": schema.StringAttribute{
				MarkdownDescription: "Filter the licenses to those assigned to a particular device. Returned in the same order that they are queued to the device",
				CustomType:          jsontypes.StringType,
				Optional:            true,
				Computed:            true,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Filter the licenses to those assigned in a particular network",
				CustomType:          jsontypes.StringType,
				Optional:            true,
				Computed:            true,
			},
			"state": schema.StringAttribute{
				MarkdownDescription: "Filter the licenses to those in a particular state. Can be one of 'active', 'expired', 'expiring', 'unused', 'unusedActive' or 'recentlyQueued'",
				CustomType:          jsontypes.StringType,
				Optional:            true,
				Computed:            true,
			},
			"list": schema.SetNestedAttribute{
				MarkdownDescription: "List of organization acls",
				Optional:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"license_id": schema.StringAttribute{
							MarkdownDescription: "License ID",
							CustomType:          jsontypes.StringType,
							Required:            true,
						},
						"device_serial": schema.StringAttribute{
							MarkdownDescription: "The serial number of the device to assign this license to. Set this to null to unassign the license. If a different license is already active on the device, this parameter will control queueing/dequeuing this license.",
							CustomType:          jsontypes.StringType,
							Required:            true,
						},
						"license_type": schema.StringAttribute{
							MarkdownDescription: "License Type.",
							CustomType:          jsontypes.StringType,
							Optional:            true,
							Computed:            true,
						},
						"license_key": schema.StringAttribute{
							MarkdownDescription: "License Key.",
							CustomType:          jsontypes.StringType,
							Optional:            true,
							Computed:            true,
						},
						"order_number": schema.StringAttribute{
							MarkdownDescription: "Order SsidNumber.",
							CustomType:          jsontypes.StringType,
							Optional:            true,
							Computed:            true,
						},
						"network_id": schema.StringAttribute{
							MarkdownDescription: "ID of the network the license is assigned to.",
							CustomType:          jsontypes.StringType,
							Optional:            true,
							Computed:            true,
						},
						"state": schema.StringAttribute{
							MarkdownDescription: "The state of the license. All queued licenses have a status of `recentlyQueued`.",
							CustomType:          jsontypes.StringType,
							Optional:            true,
							Computed:            true,
						},
						"claim_date": schema.StringAttribute{
							MarkdownDescription: "The date the license was claimed into the organization.",
							CustomType:          jsontypes.StringType,
							Optional:            true,
							Computed:            true,
						},
						"activation_date": schema.StringAttribute{
							MarkdownDescription: "The date the license started burning.",
							CustomType:          jsontypes.StringType,
							Optional:            true,
							Computed:            true,
						},
						"expiration_date": schema.StringAttribute{
							MarkdownDescription: "The date the license will expire.",
							CustomType:          jsontypes.StringType,
							Optional:            true,
							Computed:            true,
						},
						"head_license_id": schema.StringAttribute{
							MarkdownDescription: "The id of the head license this license is queued behind. If there is no head license, it returns nil.",
							CustomType:          jsontypes.StringType,
							Optional:            true,
							Computed:            true,
						},
						"seat_count": schema.Int64Attribute{
							MarkdownDescription: "The number of seats of the license. Only applicable to SM licenses.",
							CustomType:          jsontypes.Int64Type,
							Optional:            true,
							Computed:            true,
						},
						"total_duration_in_days": schema.Int64Attribute{
							MarkdownDescription: "The duration of the license plus all permanently queued licenses associated with it.",
							CustomType:          jsontypes.Int64Type,
							Optional:            true,
							Computed:            true,
						},
						"duration_in_days": schema.Int64Attribute{
							MarkdownDescription: "The duration of the individual license.",
							CustomType:          jsontypes.Int64Type,
							Optional:            true,
							Computed:            true,
						},
						"permanently_queued_licenses": schema.SingleNestedAttribute{
							MarkdownDescription: "DEPRECATED List of permanently queued licenses attached to the license. Instead, use /organizations/{organizationId}/licenses?deviceSerial= to retrieved queued licenses for a given device.",
							Optional:            true,
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									MarkdownDescription: "Permanently queued license ID.",
									Optional:            true,
									CustomType:          jsontypes.StringType,
								},
								"license_type": schema.StringAttribute{
									MarkdownDescription: "License type.",
									Optional:            true,
									CustomType:          jsontypes.StringType,
								},
								"license_key": schema.StringAttribute{
									MarkdownDescription: "License key.",
									Optional:            true,
									CustomType:          jsontypes.StringType,
								},
								"order_number": schema.StringAttribute{
									MarkdownDescription: "Order number.",
									Optional:            true,
									CustomType:          jsontypes.StringType,
								},
								"duration_in_days": schema.Int64Attribute{
									MarkdownDescription: "The duration of the individual license.",
									Optional:            true,
									CustomType:          jsontypes.Int64Type,
								},
							},
						},
					},
				},
			},
		},
	}
}

// Configure is a method of the data source interface that Terraform calls to provide the configured provider instance to the data source.
// It passes the DataSourceData that's been stored by the provider's ConfigureFunc.
func (d *OrganizationsLicensesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

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
func (d *OrganizationsLicensesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data OrganizationsLicensesDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	// Check if there are any errors before proceeding.
	if resp.Diagnostics.HasError() {
		return
	}

	request := d.client.LicensesApi.GetOrganizationLicenses(context.Background(), data.OrganizationId.ValueString())

	if !data.PerPage.IsUnknown() {
		request.PerPage(int32(data.PerPage.ValueInt64()))
	}
	if !data.StartingAfter.IsUnknown() {
		request.StartingAfter(data.StartingAfter.ValueString())
	}
	if !data.EndingBefore.IsUnknown() {
		request.EndingBefore(data.EndingBefore.ValueString())
	}
	if !data.DeviceSerial.IsUnknown() {
		request.DeviceSerial(data.DeviceSerial.ValueString())
	}
	if !data.NetworkId.IsUnknown() {
		request.NetworkId(data.NetworkId.ValueString())
	}
	if !data.State.IsUnknown() {
		request.State(data.State.ValueString())
	}

	inlineResp, httpResp, err := request.Execute()

	// If there was an error during API call, add it to diagnostics.
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			utils.HttpDiagnostics(httpResp),
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

	for _, license := range inlineResp {
		var licenseData OrganizationsLicensesDataSourceModelList
		licenseData.LicenseType = jsontypes.StringValue(license.GetLicenseType())
		licenseData.LicenseKey = jsontypes.StringValue(license.GetLicenseKey())
		licenseData.OrderNumber = jsontypes.StringValue(license.GetOrderNumber())
		licenseData.DeviceSerial = jsontypes.StringValue(license.GetDeviceSerial())
		licenseData.NetworkId = jsontypes.StringValue(license.GetNetworkId())
		licenseData.State = jsontypes.StringValue(license.GetState())
		licenseData.ClaimDate = jsontypes.StringValue(license.GetClaimDate())
		licenseData.ActivationDate = jsontypes.StringValue(license.GetActivationDate())
		licenseData.ExpirationDate = jsontypes.StringValue(license.GetExpirationDate())
		licenseData.HeadLicenseId = jsontypes.StringValue(license.GetHeadLicenseId())
		licenseData.SeatCount = jsontypes.Int64Value(int64(license.GetSeatCount()))
		licenseData.TotalDurationInDays = jsontypes.Int64Value(int64(license.GetTotalDurationInDays()))
		licenseData.DurationInDays = jsontypes.Int64Value(int64(license.GetDurationInDays()))
		licenseData.PermanentlyQueuedLicenses = license.GetPermanentlyQueuedLicenses()
		data.List = append(data.List, licenseData)
	}

	// Set ID for the data source.
	data.Id = jsontypes.StringValue("example-id")

	// Now set the final state of the data source.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read data source")
}
