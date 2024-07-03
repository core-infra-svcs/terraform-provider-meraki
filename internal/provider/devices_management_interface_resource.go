package provider

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/utils"
	"net/http"
	"time"

	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &DevicesTestAccDevicesManagementInterfaceResourceResource{}
	_ resource.ResourceWithConfigure   = &DevicesTestAccDevicesManagementInterfaceResourceResource{}
	_ resource.ResourceWithImportState = &DevicesTestAccDevicesManagementInterfaceResourceResource{}
)

func NewDevicesTestAccDevicesManagementInterfaceResourceResource() resource.Resource {
	return &DevicesTestAccDevicesManagementInterfaceResourceResource{}
}

// DevicesTestAccDevicesManagementInterfaceResourceResource defines the resource implementation.
type DevicesTestAccDevicesManagementInterfaceResourceResource struct {
	client *openApiClient.APIClient
}

// DevicesTestAccDevicesManagementInterfaceResourceResourceModel describes the resource data model.
type DevicesTestAccDevicesManagementInterfaceResourceResourceModel struct {
	Id     types.String `tfsdk:"id"`
	Serial types.String `tfsdk:"serial" json:"serial"`
	Wan1   types.Object `tfsdk:"wan1" json:"wan1"`
	Wan2   types.Object `tfsdk:"wan2" json:"wan2"`
}

type DevicesTestAccDevicesManagementInterfaceResourceResourceModelWan struct {
	WanEnabled       types.String `tfsdk:"wan_enabled" json:"wanEnabled"`
	UsingStaticIp    types.Bool   `tfsdk:"using_static_ip" json:"usingStaticIp"`
	StaticIp         types.String `tfsdk:"static_ip" json:"staticIp"`
	StaticSubnetMask types.String `tfsdk:"static_subnet_mask" json:"staticSubnetMask"`
	StaticGatewayIp  types.String `tfsdk:"static_gateway_ip" json:"staticGatewayIp"`
	StaticDns        types.List   `tfsdk:"static_dns" json:"staticDns"`
	Vlan             types.Int64  `tfsdk:"vlan" json:"vlan,omitempty"`
}

type OutputDevicesTestAccDevicesManagementInterfaceModel struct {
	Wan1 types.Object `json:"wan1"`
	Wan2 types.Object `json:"wan2"`
}

type OutputDevicesTestAccDevicesManagementInterfaceModelWan struct {
	WanEnabled       string   `json:"wanEnabled,omitempty"`
	UsingStaticIp    bool     `json:"usingStaticIp,omitempty"`
	StaticIp         string   `json:"staticIp,omitempty"`
	StaticSubnetMask string   `json:"staticSubnetMask,omitempty"`
	StaticGatewayIp  string   `json:"staticGatewayIp,omitempty"`
	StaticDns        []string `json:"staticDns,omitempty"`
	Vlan             int64    `json:"vlan,omitempty"`
}

func WANData() map[string]attr.Type {
	return map[string]attr.Type{
		"wan_enabled":        types.StringType,
		"using_static_ip":    types.BoolType,
		"static_ip":          types.StringType,
		"static_subnet_mask": types.StringType,
		"static_gateway_ip":  types.StringType,
		"static_dns":         types.ListType{ElemType: types.StringType},
		"vlan":               types.Int64Type,
	}
}

func DevicesManagementInterfaceStateWan(rawResp map[string]interface{}, wanKey string) (types.Object, diag.Diagnostics) {
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
		wanEnabled, err := utils.ExtractStringAttr(d, "wanEnabled")
		if err != nil {
			diags.Append(err...)
		}
		wan.WanEnabled = wanEnabled

		// using_static_ip
		usingStaticIp, err := utils.ExtractBoolAttr(d, "usingStaticIp")
		if err != nil {
			diags.Append(err...)
		}
		wan.UsingStaticIp = usingStaticIp

		// static_ip
		staticIp, err := utils.ExtractStringAttr(d, "staticIp")
		if err != nil {
			diags.Append(err...)
		}
		wan.StaticIp = staticIp

		// static_subnet_mask
		staticSubnetMask, err := utils.ExtractStringAttr(d, "staticSubnetMask")
		if err != nil {
			diags.Append(err...)
		}
		wan.StaticSubnetMask = staticSubnetMask

		// static_gateway_ip
		staticGatewayIp, err := utils.ExtractStringAttr(d, "staticGatewayIp")
		if err != nil {
			diags.Append(err...)
		}
		wan.StaticGatewayIp = staticGatewayIp

		// static_dns
		staticDns, err := utils.ExtractListStringAttr(d, "staticDns")
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

func updateDevicesManagementInterfaceResourceState(ctx context.Context, state *DevicesTestAccDevicesManagementInterfaceResourceResourceModel, data map[string]interface{}, httpResp *http.Response) diag.Diagnostics {
	var diags diag.Diagnostics

	rawResp, err := tools.ExtractResponseToMap(httpResp)
	if err != nil {
		diags.AddError("Failed to Unmarshal HttpResp", err.Error())
		return diags
	}

	// ID
	if state.Id.IsNull() || state.Id.IsUnknown() {
		state.Id, diags = utils.ExtractStringAttr(rawResp, "id")
		if diags.HasError() {
			diags.AddError("ID Attribute", "")
			return diags
		}
	}

	// Serial
	if state.Serial.IsNull() || state.Serial.IsUnknown() {
		state.Serial, diags = utils.ExtractStringAttr(rawResp, "serial")
		if diags.HasError() {
			diags.AddError("Serial Attribute", "")
			return diags
		}
	}

	// Wan1
	state.Wan1, diags = DevicesManagementInterfaceStateWan(rawResp, "wan1")
	if diags.HasError() {
		diags.AddError("Wan1 Attribute", "")
		return diags
	}

	// Wan2
	state.Wan2, diags = DevicesManagementInterfaceStateWan(rawResp, "wan2")
	if diags.HasError() {
		diags.AddError("Wan2 Attribute", "")
		return diags
	}

	// Log the updated state
	tflog.Debug(ctx, "Updated state", map[string]interface{}{
		"state": state,
	})

	return diags
}

func (r *DevicesTestAccDevicesManagementInterfaceResourceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_devices_management_interface"
}

func (r *DevicesTestAccDevicesManagementInterfaceResourceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{

		MarkdownDescription: "Manage the management interface settings for a device",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: types.StringType,
			},
			"serial": schema.StringAttribute{
				MarkdownDescription: "Serial number",
				Required:            true,
				CustomType:          types.StringType,
			},
			"wan1": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"wan_enabled": schema.StringAttribute{
						MarkdownDescription: "Enable or disable the interface (only for MX devices). Valid values are 'enabled', 'disabled', and 'not configured'.",
						Optional:            true,
						Computed:            true,
						CustomType:          types.StringType,
					},
					"using_static_ip": schema.BoolAttribute{
						MarkdownDescription: "Configure the interface to have static IP settings or use DHCP.",
						Optional:            true,
						Computed:            true,
						CustomType:          types.BoolType,
					},
					"vlan": schema.Int64Attribute{
						MarkdownDescription: "The VLAN that management traffic should be tagged with. Applies whether usingStaticIp is true or false.",
						Optional:            true,
						Computed:            true,
						CustomType:          types.Int64Type,
					},
					"static_ip": schema.StringAttribute{
						MarkdownDescription: "The IP the device should use on the WAN.",
						Optional:            true,
						Computed:            true,
						CustomType:          types.StringType,
					},
					"static_subnet_mask": schema.StringAttribute{
						MarkdownDescription: "The subnet mask for the WAN.",
						Optional:            true,
						Computed:            true,
						CustomType:          types.StringType,
					},
					"static_gateway_ip": schema.StringAttribute{
						MarkdownDescription: "The IP of the gateway on the WAN.",
						Optional:            true,
						Computed:            true,
						CustomType:          types.StringType,
					},
					"static_dns": schema.ListAttribute{
						MarkdownDescription: "Up to two DNS IPs.",
						Optional:            true,
						Computed:            true,
						ElementType:         types.StringType,
					},
				},
			},
			"wan2": schema.SingleNestedAttribute{
				Optional: true,
				Computed: true,
				Attributes: map[string]schema.Attribute{
					"wan_enabled": schema.StringAttribute{
						MarkdownDescription: "Enable or disable the interface (only for MX devices). Valid values are 'enabled', 'disabled', and 'not configured'.",
						Optional:            true,
						Computed:            true,
						CustomType:          types.StringType,
					},
					"using_static_ip": schema.BoolAttribute{
						MarkdownDescription: "Configure the interface to have static IP settings or use DHCP.",
						Optional:            true,
						Computed:            true,
						CustomType:          types.BoolType,
					},
					"vlan": schema.Int64Attribute{
						MarkdownDescription: "The VLAN that management traffic should be tagged with. Applies whether usingStaticIp is true or false.",
						Optional:            true,
						Computed:            true,
						CustomType:          types.Int64Type,
					},
					"static_ip": schema.StringAttribute{
						MarkdownDescription: "The IP the device should use on the WAN.",
						Optional:            true,
						Computed:            true,
						CustomType:          types.StringType,
					},
					"static_subnet_mask": schema.StringAttribute{
						MarkdownDescription: "The subnet mask for the WAN.",
						Optional:            true,
						Computed:            true,
						CustomType:          types.StringType,
					},
					"static_gateway_ip": schema.StringAttribute{
						MarkdownDescription: "The IP of the gateway on the WAN.",
						Optional:            true,
						Computed:            true,
						CustomType:          types.StringType,
					},
					"static_dns": schema.ListAttribute{
						MarkdownDescription: "Up to two DNS IPs.",
						Optional:            true,
						Computed:            true,
						ElementType:         types.StringType,
					},
				},
			},
		},
	}
}

func (r *DevicesTestAccDevicesManagementInterfaceResourceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DevicesTestAccDevicesManagementInterfaceResourceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data DevicesTestAccDevicesManagementInterfaceResourceResourceModel

	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	payload := openApiClient.NewUpdateDeviceManagementInterfaceRequest()

	if !data.Wan1.IsNull() && !data.Wan1.IsUnknown() {
		var wan1Plan DevicesTestAccDevicesManagementInterfaceResourceResourceModelWan
		data.Wan1.As(ctx, &wan1Plan, basetypes.ObjectAsOptions{})
		var staticDNS1 []string
		wan1Plan.StaticDns.ElementsAs(ctx, &staticDNS1, false)

		wan1 := openApiClient.UpdateDeviceManagementInterfaceRequestWan1{}
		wan1.SetWanEnabled(wan1Plan.WanEnabled.ValueString())
		wan1.SetStaticDns(staticDNS1)
		wan1.SetStaticGatewayIp(wan1Plan.StaticGatewayIp.ValueString())
		wan1.SetStaticSubnetMask(wan1Plan.StaticSubnetMask.ValueString())
		wan1.SetStaticIp(wan1Plan.StaticIp.ValueString())
		wan1.SetUsingStaticIp(wan1Plan.UsingStaticIp.ValueBool())

		if !wan1Plan.Vlan.IsNull() {
			vlan := int32(wan1Plan.Vlan.ValueInt64())
			wan1.Vlan = &vlan
		}

		tflog.Debug(ctx, "Wan1 payload before API call", map[string]interface{}{
			"wan1": wan1,
		})

		payload.Wan1 = &wan1
	}

	if !data.Wan2.IsNull() && !data.Wan2.IsUnknown() {
		var wan2Plan DevicesTestAccDevicesManagementInterfaceResourceResourceModelWan
		data.Wan2.As(ctx, &wan2Plan, basetypes.ObjectAsOptions{})
		var staticDNS2 []string
		wan2Plan.StaticDns.ElementsAs(ctx, &staticDNS2, false)

		wan2 := openApiClient.UpdateDeviceManagementInterfaceRequestWan2{}
		wan2.SetWanEnabled(wan2Plan.WanEnabled.ValueString())
		wan2.SetStaticDns(staticDNS2)
		wan2.SetStaticGatewayIp(wan2Plan.StaticGatewayIp.ValueString())
		wan2.SetStaticSubnetMask(wan2Plan.StaticSubnetMask.ValueString())
		wan2.SetStaticIp(wan2Plan.StaticIp.ValueString())
		wan2.SetUsingStaticIp(wan2Plan.UsingStaticIp.ValueBool())

		if !wan2Plan.Vlan.IsNull() {
			vlan := int32(wan2Plan.Vlan.ValueInt64())
			wan2.Vlan = &vlan
		}

		tflog.Debug(ctx, "Wan2 payload", map[string]interface{}{
			"wan2": wan2,
		})

		payload.Wan2 = &wan2
	}

	if data.Serial.IsNull() || data.Serial.IsUnknown() {
		resp.Diagnostics.AddError(
			"Serial Number Not Found",
			"The serial number must be provided to create the device management interface.",
		)
		return
	}

	serial := data.Serial.ValueString()
	if serial == "" {
		resp.Diagnostics.AddError(
			"Serial Number Empty",
			"The serial number provided is empty. Ensure the serial number is set correctly.",
		)
		return
	}

	maxRetries := r.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(r.client.GetConfig().Retry4xxErrorWaitTime)

	inlineResp, httpResp, err := tools.CustomHttpRequestRetry[map[string]interface{}](ctx, maxRetries, retryDelay, func() (map[string]interface{}, *http.Response, error) {
		inline, respHttp, err := r.client.ManagementInterfaceApi.UpdateDeviceManagementInterface(context.Background(), data.Serial.ValueString()).UpdateDeviceManagementInterfaceRequest(*payload).Execute()
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
			tools.HttpDiagnostics(httpResp),
		)
	} else if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
		if err != nil {
			resp.Diagnostics.AddError(
				"HTTP Client Failure",
				tools.HttpDiagnostics(httpResp),
			)
			return
		}
		if resp.Diagnostics.HasError() {
			resp.Diagnostics.AddError("State Data", fmt.Sprintf("\n%v", data))
			return
		}
	}

	diags = updateDevicesManagementInterfaceResourceState(ctx, &data, inlineResp, httpResp)
	data.Id = types.StringValue(data.Serial.ValueString())
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	tflog.Trace(ctx, "create resource")
}

func (r *DevicesTestAccDevicesManagementInterfaceResourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data DevicesTestAccDevicesManagementInterfaceResourceResourceModel

	// Read Terraform prior state data into the model

	diags := req.State.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//inlineResp, httpResp, err := r.client.DevicesApi.GetDeviceManagementInterface(context.Background(), data.Serial.ValueString()).Execute()

	maxRetries := r.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(r.client.GetConfig().Retry4xxErrorWaitTime)

	inlineResp, httpResp, err := tools.CustomHttpRequestRetry[map[string]interface{}](ctx, maxRetries, retryDelay, func() (map[string]interface{}, *http.Response, error) {
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

	}

	diags = updateDevicesManagementInterfaceResourceState(ctx, &data, inlineResp, httpResp)

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

func (r *DevicesTestAccDevicesManagementInterfaceResourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data DevicesTestAccDevicesManagementInterfaceResourceResourceModel

	// Read Terraform plan data into the model
	diags := req.Plan.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	payload := openApiClient.NewUpdateDeviceManagementInterfaceRequest()

	if !data.Wan1.IsNull() && !data.Wan1.IsUnknown() {
		var wan1Plan DevicesTestAccDevicesManagementInterfaceResourceResourceModelWan
		data.Wan1.As(ctx, &wan1Plan, basetypes.ObjectAsOptions{})
		var staticDNS1 []string
		wan1Plan.StaticDns.ElementsAs(ctx, &staticDNS1, false)

		wan1 := openApiClient.UpdateDeviceManagementInterfaceRequestWan1{}
		wan1.SetWanEnabled(wan1Plan.WanEnabled.ValueString())
		wan1.SetStaticDns(staticDNS1)
		wan1.SetStaticGatewayIp(wan1Plan.StaticGatewayIp.ValueString())
		wan1.SetStaticSubnetMask(wan1Plan.StaticSubnetMask.ValueString())
		wan1.SetStaticIp(wan1Plan.StaticIp.ValueString())
		wan1.SetUsingStaticIp(wan1Plan.UsingStaticIp.ValueBool())

		if !wan1Plan.Vlan.IsNull() {
			vlan := int32(wan1Plan.Vlan.ValueInt64())
			wan1.Vlan = &vlan
		}

		tflog.Debug(ctx, "Wan1 payload before API call", map[string]interface{}{
			"wan1": wan1,
		})

		payload.Wan1 = &wan1
	}

	if !data.Wan2.IsNull() && !data.Wan2.IsUnknown() {
		var wan2Plan DevicesTestAccDevicesManagementInterfaceResourceResourceModelWan
		data.Wan2.As(ctx, &wan2Plan, basetypes.ObjectAsOptions{})
		var staticDNS2 []string
		wan2Plan.StaticDns.ElementsAs(ctx, &staticDNS2, false)

		wan2 := openApiClient.UpdateDeviceManagementInterfaceRequestWan2{}
		wan2.SetWanEnabled(wan2Plan.WanEnabled.ValueString())
		wan2.SetStaticDns(staticDNS2)
		wan2.SetStaticGatewayIp(wan2Plan.StaticGatewayIp.ValueString())
		wan2.SetStaticSubnetMask(wan2Plan.StaticSubnetMask.ValueString())
		wan2.SetStaticIp(wan2Plan.StaticIp.ValueString())
		wan2.SetUsingStaticIp(wan2Plan.UsingStaticIp.ValueBool())

		if !wan2Plan.Vlan.IsNull() {
			vlan := int32(wan2Plan.Vlan.ValueInt64())
			wan2.Vlan = &vlan
		}

		tflog.Debug(ctx, "Wan2 payload before API call", map[string]interface{}{
			"wan2": wan2,
		})

		payload.Wan2 = &wan2
	}

	maxRetries := r.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(r.client.GetConfig().Retry4xxErrorWaitTime)

	inlineResp, httpResp, err := tools.CustomHttpRequestRetry[map[string]interface{}](ctx, maxRetries, retryDelay, func() (map[string]interface{}, *http.Response, error) {
		inline, respHttp, err := r.client.ManagementInterfaceApi.UpdateDeviceManagementInterface(context.Background(), data.Serial.ValueString()).UpdateDeviceManagementInterfaceRequest(*payload).Execute()
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

	}

	diags = updateDevicesManagementInterfaceResourceState(ctx, &data, inlineResp, httpResp)

	data.Id = types.StringValue(data.Serial.ValueString())

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Write logs using the tflog package
	tflog.Debug(ctx, "Updated state after API call", map[string]interface{}{
		"data": data,
	})

	diags = resp.State.Set(ctx, &data)
	resp.Diagnostics.Append(diags...)
	tflog.Trace(ctx, "updated resource")
}

func (r *DevicesTestAccDevicesManagementInterfaceResourceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *DevicesTestAccDevicesManagementInterfaceResourceResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	payload := openApiClient.NewUpdateDeviceManagementInterfaceRequest()
	wan1 := openApiClient.UpdateDeviceManagementInterfaceRequestWan1{}
	wan2 := openApiClient.UpdateDeviceManagementInterfaceRequestWan2{}

	payload.Wan1 = &wan1
	payload.Wan2 = &wan2

	_, httpResp, err := r.client.ManagementInterfaceApi.UpdateDeviceManagementInterface(context.Background(), data.Serial.ValueString()).UpdateDeviceManagementInterfaceRequest(*payload).Execute()

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
	}

	resp.State.RemoveResource(ctx)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")
}

func (r *DevicesTestAccDevicesManagementInterfaceResourceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("serial"), req.ID)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
