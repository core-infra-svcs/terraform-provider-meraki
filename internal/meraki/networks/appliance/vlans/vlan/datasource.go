package vlan

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
	"io"
	"strconv"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &Datasource{}

func NewNDatasource() datasource.DataSource {
	return &Datasource{}
}

// Datasource defines the resource implementation.
type Datasource struct {
	client *openApiClient.APIClient
}

func (r *Datasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_appliance_vlans"
}

func (r *Datasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *Datasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *datasourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	inlineResp, httpResp, err := r.client.ApplianceApi.GetNetworkApplianceVlans(ctx, data.NetworkId.ValueString()).Execute()
	if err != nil && httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"HTTP Client Read Failure",
			utils.HttpDiagnostics(httpResp),
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

	// Assuming httpResp is your *http.Response object
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		// Handle error: unable to read the response body
		resp.Diagnostics.AddError("Read Error", fmt.Sprintf("Unable to read HTTP response body: %v", err))
		return
	}

	// Define a struct to specifically capture the ID from the JSON data
	type HttpRespID struct {
		ID   int64  `json:"id"`
		Name string `json:"name"`
	}

	var jsonResponse []HttpRespID
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		// Handle error: JSON parsing error
		resp.Diagnostics.AddError("JSON Parsing Error", fmt.Sprintf("Error parsing JSON data for ID field: %v", err))
	}

	for _, inRespData := range inlineResp {

		vlanData := NetworksApplianceVLANModel{}
		vlanData.NetworkId = types.StringValue(data.NetworkId.ValueString())

		// Workaround for Id bug in client.GetNetworkApplianceVlans200ResponseInner
		for _, jsonInRespData := range jsonResponse {
			if jsonInRespData.Name == inRespData.GetName() {

				/*
					// Convert string to int64
							vlanId, err := strconv.ParseInt(idStr, 10, 64)
							if err != nil {
								resp.AddError("VlanId Conversion Error", fmt.Sprintf("Error converting VlanId '%s' to int64: %v", idStr, err))

				*/
				vlanData.VlanId = types.Int64Value(jsonInRespData.ID)
				data.Id = types.StringValue(fmt.Sprintf("%s,%v", data.NetworkId.ValueString(), strconv.FormatInt(jsonInRespData.ID, 10)))
			}
		}

		payloadRespDiags := DatasourceReadHttpResponse(ctx, &vlanData, &inRespData)
		if payloadRespDiags != nil {
			resp.Diagnostics.Append(payloadRespDiags...)
		}

		data.List = append(data.List, vlanData)

	}

	data.Id = types.StringValue(data.NetworkId.ValueString())

	// Check for errors after diagnostics collected
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}
