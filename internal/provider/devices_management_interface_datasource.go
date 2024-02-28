package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &DevicesManagementInterfaceDatasource{}

func NewDevicesManagementInterfaceDatasource() datasource.DataSource {
	return &DevicesManagementInterfaceDatasource{}
}

// DevicesManagementInterfaceDatasource defines the resource implementation.
type DevicesManagementInterfaceDatasource struct {
	client *openApiClient.APIClient
}

// DevicesManagementInterfaceDatasourceModel describes the resource data model.
type DevicesManagementInterfaceDatasourceModel struct {
	Id     jsontypes.String `tfsdk:"id"`
	Serial jsontypes.String `tfsdk:"serial"`
	Wan1   types.Object     `tfsdk:"wan1"`
	Wan2   types.Object     `tfsdk:"wan2"`
}

type DevicesManagementInterfaceDatasourceModelModelWan struct {
	WanEnabled       jsontypes.String `tfsdk:"wan_enabled"`
	UsingStaticIp    jsontypes.Bool   `tfsdk:"using_static_ip"`
	StaticIp         jsontypes.String `tfsdk:"static_ip"`
	StaticSubnetMask jsontypes.String `tfsdk:"static_subnet_mask"`
	StaticGatewayIp  jsontypes.String `tfsdk:"static_gateway_ip"`
	StaticDns        types.List       `tfsdk:"static_dns" json:"staticDns"`
	Vlan             jsontypes.Int64  `tfsdk:"vlan" json:"vlan,omitempty"`
}

func (r *DevicesManagementInterfaceDatasource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_devices_management_interface"
}

func (r *DevicesManagementInterfaceDatasource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{

		MarkdownDescription: "Manage the management interface settings for a device",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.StringType,
			},
			"serial": schema.StringAttribute{
				MarkdownDescription: "Serial number",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"wan1": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"wan_enabled": schema.StringAttribute{
						MarkdownDescription: "Enable or disable the interface (only for MX devices). Valid values are 'enabled', 'disabled', and 'not configured'.",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"using_static_ip": schema.BoolAttribute{
						MarkdownDescription: "Configure the interface to have static IP settings or use DHCP.",
						Optional:            true,
						CustomType:          jsontypes.BoolType,
					},
					"vlan": schema.Int64Attribute{
						MarkdownDescription: "The VLAN that management traffic should be tagged with. Applies whether usingStaticIp is true or false.",
						Optional:            true,
						CustomType:          jsontypes.Int64Type,
					},
					"static_ip": schema.StringAttribute{
						MarkdownDescription: "The IP the device should use on the WAN.",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"static_subnet_mask": schema.StringAttribute{
						MarkdownDescription: "The subnet mask for the WAN.",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"static_gateway_ip": schema.StringAttribute{
						MarkdownDescription: "The IP of the gateway on the WAN.",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"static_dns": schema.ListAttribute{
						MarkdownDescription: "Up to two DNS IPs.",
						Optional:            true,
						ElementType:         jsontypes.StringType,
					},
				},
			},
			"wan2": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"wan_enabled": schema.StringAttribute{
						MarkdownDescription: "Enable or disable the interface (only for MX devices). Valid values are 'enabled', 'disabled', and 'not configured'.",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"using_static_ip": schema.BoolAttribute{
						MarkdownDescription: "Configure the interface to have static IP settings or use DHCP.",
						Optional:            true,
						CustomType:          jsontypes.BoolType,
					},
					"vlan": schema.Int64Attribute{
						MarkdownDescription: "The VLAN that management traffic should be tagged with. Applies whether usingStaticIp is true or false.",
						Optional:            true,
						CustomType:          jsontypes.Int64Type,
					},
					"static_ip": schema.StringAttribute{
						MarkdownDescription: "The IP the device should use on the WAN.",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"static_subnet_mask": schema.StringAttribute{
						MarkdownDescription: "The subnet mask for the WAN.",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"static_gateway_ip": schema.StringAttribute{
						MarkdownDescription: "The IP of the gateway on the WAN.",
						Optional:            true,
						CustomType:          jsontypes.StringType,
					},
					"static_dns": schema.ListAttribute{
						MarkdownDescription: "Up to two DNS IPs.",
						Optional:            true,
						ElementType:         jsontypes.StringType,
					},
				},
			},
		},
	}
}

func (r *DevicesManagementInterfaceDatasource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

func (r *DevicesManagementInterfaceDatasource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data *DevicesManagementInterfaceDatasourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := r.client.ConfigureApi.GetDeviceManagementInterface(context.Background(), data.Serial.ValueString()).Execute()

	// Check for API success response code
	if httpResp.StatusCode == 404 {
		resp.Diagnostics.AddWarning(
			"No Management interface information found in API",
			tools.HttpDiagnostics(httpResp),
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
				tools.HttpDiagnostics(httpResp),
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
		if err = json.NewDecoder(httpResp.Body).Decode(&data); err != nil {
			resp.Diagnostics.AddError(
				"JSON decoding error",
				fmt.Sprintf("%v\n", err.Error()),
			)
			return
		}
	}

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}
