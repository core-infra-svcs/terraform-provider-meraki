package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
	var data *DevicesTestAccDevicesManagementInterfaceResourceResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	payload := openApiClient.NewUpdateDeviceManagementInterfaceRequest()

	if !data.Wan1.IsNull() {

		var wan1Plan DevicesTestAccDevicesManagementInterfaceResourceResourceModelWan
		data.Wan1.As(ctx, &wan1Plan, basetypes.ObjectAsOptions{})
		var staticDNS []string
		wan1Plan.StaticDns.ElementsAs(ctx, &staticDNS, false)
		wan1 := openApiClient.UpdateDeviceManagementInterfaceRequestWan1{
			WanEnabled:       wan1Plan.WanEnabled.ValueStringPointer(),
			UsingStaticIp:    wan1Plan.UsingStaticIp.ValueBoolPointer(),
			StaticIp:         wan1Plan.StaticIp.ValueStringPointer(),
			StaticGatewayIp:  wan1Plan.StaticGatewayIp.ValueStringPointer(),
			StaticSubnetMask: wan1Plan.StaticSubnetMask.ValueStringPointer(),
			StaticDns:        staticDNS,
		}
		if !wan1Plan.Vlan.IsNull() {
			var vlan = int32(wan1Plan.Vlan.ValueInt64())
			wan1.Vlan = &vlan
		}

		payload.Wan1 = &wan1

	}
	if !data.Wan2.IsNull() {
		var wan2Plan DevicesTestAccDevicesManagementInterfaceResourceResourceModelWan
		data.Wan2.As(ctx, &wan2Plan, basetypes.ObjectAsOptions{})
		var staticDNS []string
		wan2Plan.StaticDns.ElementsAs(ctx, &staticDNS, false)
		wan2 := openApiClient.UpdateDeviceManagementInterfaceRequestWan2{
			WanEnabled:       wan2Plan.WanEnabled.ValueStringPointer(),
			UsingStaticIp:    wan2Plan.UsingStaticIp.ValueBoolPointer(),
			StaticIp:         wan2Plan.StaticIp.ValueStringPointer(),
			StaticGatewayIp:  wan2Plan.StaticGatewayIp.ValueStringPointer(),
			StaticSubnetMask: wan2Plan.StaticSubnetMask.ValueStringPointer(),
			StaticDns:        staticDNS,
		}
		if !wan2Plan.Vlan.IsNull() {
			var vlan = int32(wan2Plan.Vlan.ValueInt64())
			wan2.Vlan = &vlan
		}
		payload.Wan2 = &wan2
	}

	inlineResp, httpResp, err := r.client.ManagementInterfaceApi.UpdateDeviceManagementInterface(context.Background(), data.Serial.ValueString()).UpdateDeviceManagementInterfaceRequest(*payload).Execute()

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

	var wan1Output OutputDevicesTestAccDevicesManagementInterfaceModelWan
	jsonData, _ := json.Marshal(inlineResp["wan1"].(map[string]interface{}))
	json.Unmarshal(jsonData, &wan1Output)

	var wan2Output OutputDevicesTestAccDevicesManagementInterfaceModelWan
	jsonData, _ = json.Marshal(inlineResp["wan2"].(map[string]interface{}))
	json.Unmarshal(jsonData, &wan2Output)

	staticDnsList1, diags := types.ListValueFrom(ctx, types.StringType, wan1Output.StaticDns)
	if diags.HasError() {
		return
	}

	if staticDnsList1.IsUnknown() {
		staticDnsList1 = basetypes.NewListNull(types.StringType)
	}

	staticDnsList2, diags := types.ListValueFrom(ctx, types.StringType, wan2Output.StaticDns)
	if diags.HasError() {
		return
	}

	if staticDnsList2.IsUnknown() {
		staticDnsList2 = basetypes.NewListNull(types.StringType)
	}

	wan1Map, _ := basetypes.NewObjectValue(WANData(), map[string]attr.Value{
		"wan_enabled":        basetypes.NewStringValue(wan1Output.WanEnabled),
		"using_static_ip":    basetypes.NewBoolValue(wan1Output.UsingStaticIp),
		"static_ip":          basetypes.NewStringValue(wan1Output.StaticIp),
		"static_subnet_mask": basetypes.NewStringValue(wan1Output.StaticSubnetMask),
		"static_gateway_ip":  basetypes.NewStringValue(wan1Output.StaticGatewayIp),
		"static_dns":         staticDnsList1,
		"vlan":               basetypes.NewInt64Value(wan1Output.Vlan),
	})

	objectVal1, diags := types.ObjectValueFrom(ctx, WANData(), wan1Map)
	if diags.HasError() {
		return
	}

	wan2Map, _ := basetypes.NewObjectValue(WANData(), map[string]attr.Value{
		"wan_enabled":        basetypes.NewStringValue(wan2Output.WanEnabled),
		"using_static_ip":    basetypes.NewBoolValue(wan2Output.UsingStaticIp),
		"static_ip":          basetypes.NewStringValue(wan2Output.StaticIp),
		"static_subnet_mask": basetypes.NewStringValue(wan2Output.StaticSubnetMask),
		"static_gateway_ip":  basetypes.NewStringValue(wan2Output.StaticGatewayIp),
		"static_dns":         staticDnsList2,
		"vlan":               basetypes.NewInt64Value(wan2Output.Vlan),
	})

	objectVal2, diags := types.ObjectValueFrom(ctx, WANData(), wan2Map)
	if diags.HasError() {
		return
	}

	data.Wan1 = objectVal1
	data.Wan2 = objectVal2

	if data.Wan1.IsUnknown() {
		data.Wan1 = types.ObjectNull(WANData())
	}
	if data.Wan2.IsUnknown() {
		data.Wan2 = types.ObjectNull(WANData())
	}

	data.Id = types.StringValue("example-id")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")
}

func (r *DevicesTestAccDevicesManagementInterfaceResourceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *DevicesTestAccDevicesManagementInterfaceResourceResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	inlineResp, httpResp, err := r.client.DevicesApi.GetDeviceManagementInterface(context.Background(), data.Serial.ValueString()).Execute()

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

	var wan1Output OutputDevicesTestAccDevicesManagementInterfaceModelWan
	jsonData, _ := json.Marshal(inlineResp["wan1"].(map[string]interface{}))
	json.Unmarshal(jsonData, &wan1Output)

	var wan2Output OutputDevicesTestAccDevicesManagementInterfaceModelWan
	jsonData, _ = json.Marshal(inlineResp["wan2"].(map[string]interface{}))
	json.Unmarshal(jsonData, &wan2Output)

	staticDnsList1, diags := types.ListValueFrom(ctx, types.StringType, wan1Output.StaticDns)
	if diags.HasError() {
		return
	}

	if staticDnsList1.IsUnknown() {
		staticDnsList1 = basetypes.NewListNull(types.StringType)
	}

	staticDnsList2, diags := types.ListValueFrom(ctx, types.StringType, wan2Output.StaticDns)
	if diags.HasError() {
		return
	}

	if staticDnsList2.IsUnknown() {
		staticDnsList2 = basetypes.NewListNull(types.StringType)
	}

	wan1Map, _ := basetypes.NewObjectValue(WANData(), map[string]attr.Value{
		"wan_enabled":        basetypes.NewStringValue(wan1Output.WanEnabled),
		"using_static_ip":    basetypes.NewBoolValue(wan1Output.UsingStaticIp),
		"static_ip":          basetypes.NewStringValue(wan1Output.StaticIp),
		"static_subnet_mask": basetypes.NewStringValue(wan1Output.StaticSubnetMask),
		"static_gateway_ip":  basetypes.NewStringValue(wan1Output.StaticGatewayIp),
		"static_dns":         staticDnsList1,
		"vlan":               basetypes.NewInt64Value(wan1Output.Vlan),
	})

	objectVal1, diags := types.ObjectValueFrom(ctx, WANData(), wan1Map)
	if diags.HasError() {
		return
	}

	wan2Map, _ := basetypes.NewObjectValue(WANData(), map[string]attr.Value{
		"wan_enabled":        basetypes.NewStringValue(wan2Output.WanEnabled),
		"using_static_ip":    basetypes.NewBoolValue(wan2Output.UsingStaticIp),
		"static_ip":          basetypes.NewStringValue(wan2Output.StaticIp),
		"static_subnet_mask": basetypes.NewStringValue(wan2Output.StaticSubnetMask),
		"static_gateway_ip":  basetypes.NewStringValue(wan2Output.StaticGatewayIp),
		"static_dns":         staticDnsList2,
		"vlan":               basetypes.NewInt64Value(wan2Output.Vlan),
	})

	objectVal2, diags := types.ObjectValueFrom(ctx, WANData(), wan2Map)
	if diags.HasError() {
		return
	}

	data.Wan1 = objectVal1
	data.Wan2 = objectVal2

	if data.Wan1.IsUnknown() {
		data.Wan1 = types.ObjectNull(WANData())
	}
	if data.Wan2.IsUnknown() {
		data.Wan2 = types.ObjectNull(WANData())
	}

	data.Id = types.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *DevicesTestAccDevicesManagementInterfaceResourceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *DevicesTestAccDevicesManagementInterfaceResourceResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	payload := openApiClient.NewUpdateDeviceManagementInterfaceRequest()

	if !data.Wan1.IsNull() {

		var wan1Plan DevicesTestAccDevicesManagementInterfaceResourceResourceModelWan
		data.Wan1.As(ctx, &wan1Plan, basetypes.ObjectAsOptions{})
		var staticDNS []string
		wan1Plan.StaticDns.ElementsAs(ctx, &staticDNS, false)
		wan1 := openApiClient.UpdateDeviceManagementInterfaceRequestWan1{
			WanEnabled:       wan1Plan.WanEnabled.ValueStringPointer(),
			UsingStaticIp:    wan1Plan.UsingStaticIp.ValueBoolPointer(),
			StaticIp:         wan1Plan.StaticIp.ValueStringPointer(),
			StaticGatewayIp:  wan1Plan.StaticGatewayIp.ValueStringPointer(),
			StaticSubnetMask: wan1Plan.StaticSubnetMask.ValueStringPointer(),
			StaticDns:        staticDNS,
		}
		if !wan1Plan.Vlan.IsNull() {
			var vlan = int32(wan1Plan.Vlan.ValueInt64())
			wan1.Vlan = &vlan
		}

		payload.Wan1 = &wan1

	}
	if !data.Wan2.IsNull() {
		var wan2Plan DevicesTestAccDevicesManagementInterfaceResourceResourceModelWan
		data.Wan2.As(ctx, &wan2Plan, basetypes.ObjectAsOptions{})
		var staticDNS2 []string
		wan2Plan.StaticDns.ElementsAs(ctx, &staticDNS2, false)
		wan2 := openApiClient.UpdateDeviceManagementInterfaceRequestWan2{
			WanEnabled:       wan2Plan.WanEnabled.ValueStringPointer(),
			UsingStaticIp:    wan2Plan.UsingStaticIp.ValueBoolPointer(),
			StaticIp:         wan2Plan.StaticIp.ValueStringPointer(),
			StaticGatewayIp:  wan2Plan.StaticGatewayIp.ValueStringPointer(),
			StaticSubnetMask: wan2Plan.StaticSubnetMask.ValueStringPointer(),
			StaticDns:        staticDNS2,
		}
		if !wan2Plan.Vlan.IsNull() {
			var vlan = int32(wan2Plan.Vlan.ValueInt64())
			wan2.Vlan = &vlan
		}
		payload.Wan2 = &wan2
	}

	inlineResp, httpResp, err := r.client.ManagementInterfaceApi.UpdateDeviceManagementInterface(context.Background(), data.Serial.ValueString()).UpdateDeviceManagementInterfaceRequest(*payload).Execute()
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

	var wan1Output OutputDevicesTestAccDevicesManagementInterfaceModelWan
	jsonData, _ := json.Marshal(inlineResp["wan1"].(map[string]interface{}))
	json.Unmarshal(jsonData, &wan1Output)

	var wan2Output OutputDevicesTestAccDevicesManagementInterfaceModelWan
	jsonData, _ = json.Marshal(inlineResp["wan2"].(map[string]interface{}))
	json.Unmarshal(jsonData, &wan2Output)

	staticDnsList1, diags := types.ListValueFrom(ctx, types.StringType, wan1Output.StaticDns)
	if diags.HasError() {
		return
	}

	if staticDnsList1.IsUnknown() {
		staticDnsList1 = basetypes.NewListNull(types.StringType)
	}

	staticDnsList2, diags := types.ListValueFrom(ctx, types.StringType, wan2Output.StaticDns)
	if diags.HasError() {
		return
	}

	if staticDnsList2.IsUnknown() {
		staticDnsList2 = basetypes.NewListNull(types.StringType)
	}

	wan1Map, _ := basetypes.NewObjectValue(WANData(), map[string]attr.Value{
		"wan_enabled":        basetypes.NewStringValue(wan1Output.WanEnabled),
		"using_static_ip":    basetypes.NewBoolValue(wan1Output.UsingStaticIp),
		"static_ip":          basetypes.NewStringValue(wan1Output.StaticIp),
		"static_subnet_mask": basetypes.NewStringValue(wan1Output.StaticSubnetMask),
		"static_gateway_ip":  basetypes.NewStringValue(wan1Output.StaticGatewayIp),
		"static_dns":         staticDnsList1,
		"vlan":               basetypes.NewInt64Value(wan1Output.Vlan),
	})

	objectVal1, diags := types.ObjectValueFrom(ctx, WANData(), wan1Map)
	if diags.HasError() {
		return
	}

	wan2Map, _ := basetypes.NewObjectValue(WANData(), map[string]attr.Value{
		"wan_enabled":        basetypes.NewStringValue(wan2Output.WanEnabled),
		"using_static_ip":    basetypes.NewBoolValue(wan2Output.UsingStaticIp),
		"static_ip":          basetypes.NewStringValue(wan2Output.StaticIp),
		"static_subnet_mask": basetypes.NewStringValue(wan2Output.StaticSubnetMask),
		"static_gateway_ip":  basetypes.NewStringValue(wan2Output.StaticGatewayIp),
		"static_dns":         staticDnsList2,
		"vlan":               basetypes.NewInt64Value(wan2Output.Vlan),
	})

	objectVal2, diags := types.ObjectValueFrom(ctx, WANData(), wan2Map)
	if diags.HasError() {
		return
	}

	data.Wan1 = objectVal1
	data.Wan2 = objectVal2

	if data.Wan1.IsUnknown() {
		data.Wan1 = types.ObjectNull(WANData())
	}
	if data.Wan2.IsUnknown() {
		data.Wan2 = types.ObjectNull(WANData())
	}

	data.Id = types.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
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

	data.Id = types.StringValue("example-id")

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
