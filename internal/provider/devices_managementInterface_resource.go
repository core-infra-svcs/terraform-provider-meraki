package provider

import (
	"encoding/json"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &DevicesManagementinterfaceResource{}
	_ resource.ResourceWithConfigure   = &DevicesManagementinterfaceResource{}
	_ resource.ResourceWithImportState = &DevicesManagementinterfaceResource{}
)

func NewDevicesManagementinterfaceResource() resource.Resource {
	return &DevicesManagementinterfaceResource{}
}

// DevicesManagementinterfaceResource defines the resource implementation.
type DevicesManagementinterfaceResource struct {
	client *openApiClient.APIClient
}

// DevicesManagementinterfaceResourceModel describes the resource data model.
type DevicesManagementinterfaceResourceModel struct {
	Id     jsontypes.String             `tfsdk:"id"`
	Serial jsontypes.String             `tfsdk:"serial"`
	Wan1   DeviceManagementInterfaceWan `tfsdk:"wan1"`
	Wan2   DeviceManagementInterfaceWan `tfsdk:"wan2"`
}

type DeviceManagementInterfaceWan struct {
	WanEnabled       jsontypes.String `tfsdk:"wan_enabled"`
	UsingStaticIp    jsontypes.Bool   `tfsdk:"using_static_ip"`
	StaticIp         jsontypes.String `tfsdk:"static_ip"`
	StaticSubnetMask jsontypes.String `tfsdk:"static_subnet_mask"`
	StaticGatewayIp  jsontypes.String `tfsdk:"static_gateway_ip"`
	StaticDns        types.List       `tfsdk:"static_dns"`
	Vlan             jsontypes.Int64  `tfsdk:"vlan"`
}

func (r *DevicesManagementinterfaceResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_devices_management_interface"
}

func (r *DevicesManagementinterfaceResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{

		MarkdownDescription: "DevicesManagementinterface",
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

func (r *DevicesManagementinterfaceResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DevicesManagementinterfaceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *DevicesManagementinterfaceResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	object14 := openApiClient.NewInlineObject14()

	updateDeviceManagementInterfaceRequest := object14

	vlan := int32(data.Wan1.Vlan.ValueInt64())
	staticDNS := []string{}
	for _, dns := range data.Wan1.StaticDns.Elements() {
		staticDNS = append(staticDNS, dns.String())
	}
	wan1 := openApiClient.DevicesSerialManagementInterfaceWan1{
		WanEnabled:       data.Wan1.WanEnabled.ValueStringPointer(),
		UsingStaticIp:    data.Wan1.UsingStaticIp.ValueBoolPointer(),
		StaticIp:         data.Wan1.StaticIp.ValueStringPointer(),
		StaticGatewayIp:  data.Wan1.StaticGatewayIp.ValueStringPointer(),
		StaticSubnetMask: data.Wan1.StaticSubnetMask.ValueStringPointer(),
		StaticDns:        staticDNS,
		Vlan:             &vlan,
	}

	vlan = int32(data.Wan2.Vlan.ValueInt64())
	staticDNS = []string{}
	for _, dns := range data.Wan2.StaticDns.Elements() {
		staticDNS = append(staticDNS, dns.String())
	}
	wan2 := openApiClient.DevicesSerialManagementInterfaceWan2{
		WanEnabled:       data.Wan2.WanEnabled.ValueStringPointer(),
		UsingStaticIp:    data.Wan2.UsingStaticIp.ValueBoolPointer(),
		StaticIp:         data.Wan2.StaticIp.ValueStringPointer(),
		StaticGatewayIp:  data.Wan2.StaticGatewayIp.ValueStringPointer(),
		StaticSubnetMask: data.Wan2.StaticSubnetMask.ValueStringPointer(),
		StaticDns:        staticDNS,
		Vlan:             &vlan,
	}

	updateDeviceManagementInterfaceRequest.Wan1 = &wan1
	updateDeviceManagementInterfaceRequest.Wan2 = &wan2

	_, httpResp, err := r.client.ManagementInterfaceApi.UpdateDeviceManagementInterface(context.Background(), data.Serial.ValueString()).UpdateDeviceManagementInterface(*updateDeviceManagementInterfaceRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// collect diagnostics
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
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
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%s", data))
		return
	}

	// Save data into Terraform state
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")
}

func (r *DevicesManagementinterfaceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *DevicesManagementinterfaceResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := r.client.DevicesApi.GetDeviceManagementInterface(context.Background(), data.Serial.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to read resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// collect diagnostics
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
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
		return
	}

	// Save data into Terraform state
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *DevicesManagementinterfaceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *DevicesManagementinterfaceResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	object14 := openApiClient.NewInlineObject14()

	updateDeviceManagementInterfaceRequest := object14

	vlan := int32(data.Wan1.Vlan.ValueInt64())
	staticDNS := []string{}
	for _, dns := range data.Wan1.StaticDns.Elements() {
		staticDNS = append(staticDNS, dns.String())
	}
	wan1 := openApiClient.DevicesSerialManagementInterfaceWan1{
		WanEnabled:       data.Wan1.WanEnabled.ValueStringPointer(),
		UsingStaticIp:    data.Wan1.UsingStaticIp.ValueBoolPointer(),
		StaticIp:         data.Wan1.StaticIp.ValueStringPointer(),
		StaticGatewayIp:  data.Wan1.StaticGatewayIp.ValueStringPointer(),
		StaticSubnetMask: data.Wan1.StaticSubnetMask.ValueStringPointer(),
		StaticDns:        staticDNS,
		Vlan:             &vlan,
	}

	vlan = int32(data.Wan2.Vlan.ValueInt64())
	staticDNS = []string{}
	for _, dns := range data.Wan2.StaticDns.Elements() {
		staticDNS = append(staticDNS, dns.String())
	}
	wan2 := openApiClient.DevicesSerialManagementInterfaceWan2{
		WanEnabled:       data.Wan2.WanEnabled.ValueStringPointer(),
		UsingStaticIp:    data.Wan2.UsingStaticIp.ValueBoolPointer(),
		StaticIp:         data.Wan2.StaticIp.ValueStringPointer(),
		StaticGatewayIp:  data.Wan2.StaticGatewayIp.ValueStringPointer(),
		StaticSubnetMask: data.Wan2.StaticSubnetMask.ValueStringPointer(),
		StaticDns:        staticDNS,
		Vlan:             &vlan,
	}

	updateDeviceManagementInterfaceRequest.Wan1 = &wan1
	updateDeviceManagementInterfaceRequest.Wan2 = &wan2

	_, httpResp, err := r.client.ManagementInterfaceApi.UpdateDeviceManagementInterface(context.Background(), data.Serial.ValueString()).UpdateDeviceManagementInterface(*updateDeviceManagementInterfaceRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// collect diagnostics
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
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
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%s", data))
		return
	}

	// Save data into Terraform state
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *DevicesManagementinterfaceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *DevicesManagementinterfaceResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	object14 := openApiClient.NewInlineObject14()

	updateDeviceManagementInterfaceRequest := object14

	vlan := int32(data.Wan1.Vlan.ValueInt64())
	staticDNS := []string{}
	for _, dns := range data.Wan1.StaticDns.Elements() {
		staticDNS = append(staticDNS, dns.String())
	}
	wan1 := openApiClient.DevicesSerialManagementInterfaceWan1{
		WanEnabled:       data.Wan1.WanEnabled.ValueStringPointer(),
		UsingStaticIp:    data.Wan1.UsingStaticIp.ValueBoolPointer(),
		StaticIp:         data.Wan1.StaticIp.ValueStringPointer(),
		StaticGatewayIp:  data.Wan1.StaticGatewayIp.ValueStringPointer(),
		StaticSubnetMask: data.Wan1.StaticSubnetMask.ValueStringPointer(),
		StaticDns:        staticDNS,
		Vlan:             &vlan,
	}

	vlan = int32(data.Wan2.Vlan.ValueInt64())
	staticDNS = []string{}
	for _, dns := range data.Wan2.StaticDns.Elements() {
		staticDNS = append(staticDNS, dns.String())
	}
	wan2 := openApiClient.DevicesSerialManagementInterfaceWan2{
		WanEnabled:       data.Wan2.WanEnabled.ValueStringPointer(),
		UsingStaticIp:    data.Wan2.UsingStaticIp.ValueBoolPointer(),
		StaticIp:         data.Wan2.StaticIp.ValueStringPointer(),
		StaticGatewayIp:  data.Wan2.StaticGatewayIp.ValueStringPointer(),
		StaticSubnetMask: data.Wan2.StaticSubnetMask.ValueStringPointer(),
		StaticDns:        staticDNS,
		Vlan:             &vlan,
	}

	updateDeviceManagementInterfaceRequest.Wan1 = &wan1
	updateDeviceManagementInterfaceRequest.Wan2 = &wan2

	_, httpResp, err := r.client.ManagementInterfaceApi.UpdateDeviceManagementInterface(context.Background(), data.Serial.ValueString()).UpdateDeviceManagementInterface(*updateDeviceManagementInterfaceRequest).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// collect diagnostics
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
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
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%s", data))
		return
	}

	// Save data into Terraform state
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	data.Id = jsontypes.StringValue("example-id")

	resp.State.RemoveResource(ctx)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")
}

func (r *DevicesManagementinterfaceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("serial"), req.ID)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
