package ports

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
	"io"
	"net/http"
	"time"
)

var _ datasource.DataSource = &PortsDataSource{}

func NewDevicesSwitchPortsStatusesDataSource() datasource.DataSource {
	return &PortsDataSource{}
}

type PortsDataSource struct {
	client *openApiClient.APIClient
}

func (d *PortsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {

	resp.TypeName = req.ProviderTypeName + "_devices_switch_ports"
}

// Schema provides a way to define the structure of the data source data.
// It is called by the framework to get the schema of the data source.
func (d *PortsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {

	// The Schema object defines the structure of the data source.
	resp.Schema = schema.Schema{

		MarkdownDescription: "Ports the switch ports for a switch",

		// The Attributes map describes the fields of the data source.
		Attributes: map[string]schema.Attribute{

			// Every data source must have an ID attribute. This is computed by the framework.
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.StringType,
			},
			"serial": schema.StringAttribute{
				MarkdownDescription: "A list of serial numbers. The returned devices will be filtered to only include these serials.",
				CustomType:          jsontypes.StringType,
				Required:            true,
			},
			"ports": schema.SetNestedAttribute{
				MarkdownDescription: "Ports of switch ports",
				Optional:            true,
				Computed:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: portsDataSourceSchema,
				},
			},
		},
	}
}

// Configure is a method of the data source interface that Terraform calls to provide the configured provider instance to the data source.
// It passes the DataSourceData that's been stored by the provider's ConfigureFunc.
func (d *PortsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {

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
func (d *PortsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data PortsDataSourceModel
	var diags diag.Diagnostics
	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	// Check if there are any errors before proceeding.
	if resp.Diagnostics.HasError() {
		return
	}

	maxRetries := d.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(d.client.GetConfig().Retry4xxErrorWaitTime)

	// usage of CustomHttpRequestRetry with a slice of strongly typed structs
	apiCallSlice := func() ([]openApiClient.GetDeviceSwitchPorts200ResponseInner, *http.Response, error, diag.Diagnostics) {
		inline, httpResp, err := d.client.SwitchApi.GetDeviceSwitchPorts(ctx, data.Serial.ValueString()).Execute()
		return inline, httpResp, err, diags
	}

	resultSlice, httpRespSlice, errSlice, tfDiags := utils.CustomHttpRequestRetryStronglyTyped(ctx, maxRetries, retryDelay, apiCallSlice)
	if errSlice != nil {

		if tfDiags.HasError() {
			fmt.Printf(": %s\n", tfDiags.Errors())
		}

		fmt.Printf("Error creating switch port: %s\n", errSlice)
		if httpRespSlice != nil {
			var responseBody string
			if httpRespSlice.Body != nil {
				bodyBytes, readErr := io.ReadAll(httpRespSlice.Body)
				if readErr == nil {
					responseBody = string(bodyBytes)
				}
			}
			fmt.Printf("Failed to create resource. HTTP Status Code: %d, Response Body: %s\n", httpRespSlice.StatusCode, responseBody)
		}
		return
	}

	// Type assert apiResp to the expected []openApiClient.GetDeviceSwitchPorts200ResponseInner type
	inlineResp, ok := any(resultSlice).([]openApiClient.GetDeviceSwitchPorts200ResponseInner)
	if !ok {
		fmt.Println("Failed to assert API response type to []openApiClient.GetDeviceSwitchPorts200ResponseInner. Please ensure the API response structure matches the expected type.")
		return
	}

	// If it's not what you expect, add an error to diagnostics.
	if httpRespSlice.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpRespSlice.StatusCode),
		)
	}

	// If there were any errors up to this point, log the state data and return.
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("State Data", fmt.Sprintf("\n%v", data))
		return
	}

	for _, i := range inlineResp {

		var portData *PortResourceModel

		// Use typedApiResp with the correct type for further processing
		port, portDiags := PortResourceState(ctx, &i, portData)
		if portDiags.HasError() {
			resp.Diagnostics.AddError("Resource Response Error", fmt.Sprintf("\n%v", portDiags))
			return
		}

		data.Ports = append(data.Ports, *port)
	}

	// Now set the final state of the data source.
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read data source")
}
