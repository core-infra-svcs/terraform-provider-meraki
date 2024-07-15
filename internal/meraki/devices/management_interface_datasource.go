package devices

import (
	"context"
	"fmt"
	utils2 "github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"

	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
	"net/http"
	"time"
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
	Id     types.String `tfsdk:"id"`
	Serial types.String `tfsdk:"serial" json:"serial"`
	Wan1   types.Object `tfsdk:"wan1" json:"wan1"`
	Wan2   types.Object `tfsdk:"wan2" json:"wan2"`
}

type DevicesManagementInterfaceDatasourceModelModelWan struct {
	WanEnabled       types.String `tfsdk:"wan_enabled" json:"wanEnabled"`
	UsingStaticIp    types.Bool   `tfsdk:"using_static_ip" json:"usingStaticIp"`
	StaticIp         types.String `tfsdk:"static_ip" json:"staticIp"`
	StaticSubnetMask types.String `tfsdk:"static_subnet_mask" json:"staticSubnetMask"`
	StaticGatewayIp  types.String `tfsdk:"static_gateway_ip" json:"staticGatewayIp"`
	StaticDns        types.List   `tfsdk:"static_dns" json:"staticDns"`
	Vlan             types.Int64  `tfsdk:"vlan" json:"vlan,omitempty"`
}

func DevicesManagementInterfaceDatasourceStateWan(rawResp map[string]interface{}, wanKey string) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var wan DevicesTestAccDevicesManagementInterfaceResourceResourceModelWan

	wanAttrs := map[string]attr.Type{
		"wan_enabled":        types.StringType,
		"using_static_ip":    types.BoolType,
		"static_ip":          types.StringType,
		"static_subnet_mask": types.StringType,
		"static_gateway_ip":  types.StringType,
		"static_dns":         types.ListType{ElemType: types.StringType},
		"vlan":               types.Int64Type,
	}

	if d, ok := rawResp[wanKey].(map[string]interface{}); ok {
		// wan_enabled
		wanEnabled, err := utils2.ExtractStringAttr(d, "wanEnabled")
		if err != nil {
			diags.Append(err...)
		}
		wan.WanEnabled = wanEnabled

		// using_static_ip
		usingStaticIp, err := utils2.ExtractBoolAttr(d, "usingStaticIp")
		if err != nil {
			diags.Append(err...)
		}
		wan.UsingStaticIp = usingStaticIp

		// static_ip
		staticIp, err := utils2.ExtractStringAttr(d, "staticIp")
		if err != nil {
			diags.Append(err...)
		}
		wan.StaticIp = staticIp

		// static_subnet_mask
		staticSubnetMask, err := utils2.ExtractStringAttr(d, "staticSubnetMask")
		if err != nil {
			diags.Append(err...)
		}
		wan.StaticSubnetMask = staticSubnetMask

		// static_gateway_ip
		staticGatewayIp, err := utils2.ExtractStringAttr(d, "staticGatewayIp")
		if err != nil {
			diags.Append(err...)
		}
		wan.StaticGatewayIp = staticGatewayIp

		// static_dns
		staticDns, err := utils2.ExtractListStringAttr(d, "staticDns")
		if err != nil {
			diags.Append(err...)
		}
		wan.StaticDns = staticDns

		// vlan
		if vlanValue, exists := d["vlan"]; exists && vlanValue != nil {
			switch v := vlanValue.(type) {
			case float64:
				wan.Vlan = types.Int64Value(int64(v))
			case int64:
				wan.Vlan = types.Int64Value(v)
			case int:
				wan.Vlan = types.Int64Value(int64(v))
			default:
				wan.Vlan = types.Int64Null()
				diags.AddError("Type Error", fmt.Sprintf("Unsupported type for vlan attribute: %T", v))
			}
		} else {
			wan.Vlan = types.Int64Null()
		}

		// Log the extracted vlan value
		tflog.Debug(context.Background(), "Extracted vlan", map[string]interface{}{
			"vlan": wan.Vlan.ValueInt64(),
		})

	} else {
		WanNull := types.ObjectNull(wanAttrs)
		return WanNull, diags
	}

	wanObj, err := types.ObjectValueFrom(context.Background(), wanAttrs, wan)
	if err != nil {
		diags.Append(err...)
	}

	return wanObj, diags
}

func updateDevicesManagementInterfaceDatasourceState(ctx context.Context, state *DevicesManagementInterfaceDatasourceModel, data map[string]interface{}, httpResp *http.Response) diag.Diagnostics {
	var diags diag.Diagnostics

	rawResp, err := utils2.ExtractResponseToMap(httpResp)
	if err != nil {
		diags.AddError("Failed to Unmarshal HttpResp", err.Error())
		return diags
	}

	// ID
	if state.Id.IsNull() || state.Id.IsUnknown() {
		state.Id, diags = utils2.ExtractStringAttr(rawResp, "id")
		if diags.HasError() {
			diags.AddError("ID Attribute", "")
			return diags
		}
	}

	// Serial
	if state.Serial.IsNull() || state.Serial.IsUnknown() {
		state.Serial, diags = utils2.ExtractStringAttr(rawResp, "serial")
		if diags.HasError() {
			diags.AddError("Serial Attribute", "")
			return diags
		}
	}

	// Wan1
	state.Wan1, diags = DevicesManagementInterfaceDatasourceStateWan(rawResp, "wan1")
	if diags.HasError() {
		diags.AddError("Wan1 Attribute", "")
		return diags
	}

	// Wan2
	state.Wan2, diags = DevicesManagementInterfaceDatasourceStateWan(rawResp, "wan2")
	if diags.HasError() {
		diags.AddError("Wan2 Attribute", "")
		return diags
	}

	// Log the state after updating
	tflog.Debug(ctx, "State after update", map[string]interface{}{
		"state": state,
	})

	return diags
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
				CustomType: types.StringType,
			},
			"serial": schema.StringAttribute{
				MarkdownDescription: "Serial number",
				Optional:            true,
				Computed:            true,
				CustomType:          types.StringType,
			},
			"wan1": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"wan_enabled": schema.StringAttribute{
						MarkdownDescription: "Enable or disable the interface (only for MX devices). Valid values are 'enabled', 'disabled', and 'not configured'.",
						Optional:            true,
						CustomType:          types.StringType,
					},
					"using_static_ip": schema.BoolAttribute{
						MarkdownDescription: "Configure the interface to have static IP settings or use DHCP.",
						Optional:            true,
						CustomType:          types.BoolType,
					},
					"vlan": schema.Int64Attribute{
						MarkdownDescription: "The VLAN that management traffic should be tagged with. Applies whether usingStaticIp is true or false.",
						Optional:            true,
						CustomType:          types.Int64Type,
					},
					"static_ip": schema.StringAttribute{
						MarkdownDescription: "The IP the device should use on the WAN.",
						Optional:            true,
						CustomType:          types.StringType,
					},
					"static_subnet_mask": schema.StringAttribute{
						MarkdownDescription: "The subnet mask for the WAN.",
						Optional:            true,
						CustomType:          types.StringType,
					},
					"static_gateway_ip": schema.StringAttribute{
						MarkdownDescription: "The IP of the gateway on the WAN.",
						Optional:            true,
						CustomType:          types.StringType,
					},
					"static_dns": schema.ListAttribute{
						MarkdownDescription: "Up to two DNS IPs.",
						Optional:            true,
						ElementType:         types.StringType,
					},
				},
			},
			"wan2": schema.SingleNestedAttribute{
				Optional: true,
				Attributes: map[string]schema.Attribute{
					"wan_enabled": schema.StringAttribute{
						MarkdownDescription: "Enable or disable the interface (only for MX devices). Valid values are 'enabled', 'disabled', and 'not configured'.",
						Optional:            true,
						CustomType:          types.StringType,
					},
					"using_static_ip": schema.BoolAttribute{
						MarkdownDescription: "Configure the interface to have static IP settings or use DHCP.",
						Optional:            true,
						CustomType:          types.BoolType,
					},
					"vlan": schema.Int64Attribute{
						MarkdownDescription: "The VLAN that management traffic should be tagged with. Applies whether usingStaticIp is true or false.",
						Optional:            true,
						CustomType:          types.Int64Type,
					},
					"static_ip": schema.StringAttribute{
						MarkdownDescription: "The IP the device should use on the WAN.",
						Optional:            true,
						CustomType:          types.StringType,
					},
					"static_subnet_mask": schema.StringAttribute{
						MarkdownDescription: "The subnet mask for the WAN.",
						Optional:            true,
						CustomType:          types.StringType,
					},
					"static_gateway_ip": schema.StringAttribute{
						MarkdownDescription: "The IP of the gateway on the WAN.",
						Optional:            true,
						CustomType:          types.StringType,
					},
					"static_dns": schema.ListAttribute{
						MarkdownDescription: "Up to two DNS IPs.",
						Optional:            true,
						ElementType:         types.StringType,
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
	var data DevicesManagementInterfaceDatasourceModel

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	//inlineResp, httpResp, err := r.client.DevicesApi.GetDeviceManagementInterface(context.Background(), data.Serial.ValueString()).Execute()

	maxRetries := r.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(r.client.GetConfig().Retry4xxErrorWaitTime)

	inlineResp, httpResp, err := utils2.CustomHttpRequestRetry[map[string]interface{}](ctx, maxRetries, retryDelay, func() (map[string]interface{}, *http.Response, error) {
		inline, respHttp, err := r.client.DevicesApi.GetDeviceManagementInterface(context.Background(), data.Serial.ValueString()).Execute()
		return inline, respHttp, err
	})
	if err != nil {
		tflog.Error(ctx, "HTTP Call Failed", map[string]interface{}{
			"error": err.Error(),
		})
		resp.Diagnostics.AddError(
			"HTTP Call Failed",
			fmt.Sprintf("Details: %s", err.Error()),
		)
	}

	// Check for API success response code
	if httpResp.StatusCode == 404 {
		resp.Diagnostics.AddWarning(
			"No Management interface information found in API",
			utils2.HttpDiagnostics(httpResp),
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
				utils2.HttpDiagnostics(httpResp),
			)
			return
		}

		// If there were any errors up to this point, log the state data and return.
		if resp.Diagnostics.HasError() {
			resp.Diagnostics.AddError("State Data", fmt.Sprintf("\n%v", data))
			return
		}

	}

	diags = updateDevicesManagementInterfaceDatasourceState(ctx, &data, inlineResp, httpResp)

	data.Id = types.StringValue(data.Serial.ValueString())

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}
