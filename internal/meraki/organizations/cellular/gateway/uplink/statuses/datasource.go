package statuses

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

var _ datasource.DataSource = &DataSource{}

func NewDataSource() datasource.DataSource {
	return &DataSource{}
}

type DataSource struct {
	client *openApiClient.APIClient
}

// Metadata provides a way to define information about the data source.
// This method is called by the framework to retrieve metadata about the data source.
func (d *DataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_organizations_cellular_gateway_uplink_statuses"
}

// Configure is a method of the data source interface that Terraform calls to provide the configured provider instance to the data source.
// It passes the DataSourceData that's been stored by the provider's ConfigureFunc.
func (d *DataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

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
func (d *DataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data dataSourceModel

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
