package provider

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
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
	Id                   jsontypes.String `tfsdk:"id"`
	Serial               jsontypes.String `tfsdk:"serial"`
	Wan1WanEnabled       jsontypes.String `tfsdk:"wan1_wan_enabled"`
	Wan1UsingStaticIp    jsontypes.Bool   `tfsdk:"wan1_using_static_ip"`
	Wan1StaticIp         jsontypes.String `tfsdk:"wan1_static_ip"`
	Wan1StaticSubnetMask jsontypes.String `tfsdk:"wan1_static_subnet_mask"`
	Wan1StaticGatewayIp  jsontypes.String `tfsdk:"wan1_static_gateway_ip"`
	Wan1StaticDns        types.List       `tfsdk:"wan1_static_dns"`
	Wan1Vlan             jsontypes.Int64  `tfsdk:"wan1_vlan"`
	Wan2WanEnabled       jsontypes.String `tfsdk:"wan2_wan_enabled"`
	Wan2UsingStaticIp    jsontypes.Bool   `tfsdk:"wan2_using_static_ip"`
	Wan2Vlan             jsontypes.Int64  `tfsdk:"wan2_vlan"`
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
			"wan1_wan_enabled": schema.StringAttribute{
				MarkdownDescription: "Enable or disable the interface (only for MX devices). Valid values are 'enabled', 'disabled', and 'not configured'.",
				Optional:            true,
				CustomType:          jsontypes.StringType,
			},
			"wan1_using_static_ip": schema.BoolAttribute{
				MarkdownDescription: "Configure the interface to have static IP settings or use DHCP.",
				Optional:            true,
				CustomType:          jsontypes.BoolType,
			},
			"wan1_vlan": schema.Int64Attribute{
				MarkdownDescription: "The VLAN that management traffic should be tagged with. Applies whether usingStaticIp is true or false.",
				Optional:            true,
				CustomType:          jsontypes.Int64Type,
			},
			"wan1_static_ip": schema.StringAttribute{
				MarkdownDescription: "The IP the device should use on the WAN.",
				Optional:            true,
				CustomType:          jsontypes.StringType,
			},
			"wan1_static_subnet_mask": schema.StringAttribute{
				MarkdownDescription: "The subnet mask for the WAN.",
				Optional:            true,
				CustomType:          jsontypes.StringType,
			},
			"wan1_static_gateway_ip": schema.StringAttribute{
				MarkdownDescription: "The IP of the gateway on the WAN.",
				Optional:            true,
				CustomType:          jsontypes.StringType,
			},
			"wan1_static_dns": schema.ListNestedAttribute{
				MarkdownDescription: "Up to two DNS IPs.",
				Optional:            true,
			},
			"wan2_wan_enabled": schema.StringAttribute{
				MarkdownDescription: "Enable or disable the interface (only for MX devices). Valid values are 'enabled', 'disabled', and 'not configured'.",
				Optional:            true,
				CustomType:          jsontypes.StringType,
			},
			"wan2_using_static_ip": schema.BoolAttribute{
				MarkdownDescription: "Configure the interface to have static IP settings or use DHCP.",
				Optional:            true,
				CustomType:          jsontypes.BoolType,
			},
			"wan2_vlan": schema.Int64Attribute{
				MarkdownDescription: "The VLAN that management traffic should be tagged with. Applies whether usingStaticIp is true or false.",
				Optional:            true,
				CustomType:          jsontypes.Int64Type,
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

	wan1 := openApiClient.NewDevicesSerialManagementInterfaceWan1()
	wan2 := openApiClient.NewDevicesSerialManagementInterfaceWan2()
	wan1.SetStaticGatewayIp(data.Wan1StaticGatewayIp.ValueString())
	wan1.SetWanEnabled(data.Wan1WanEnabled.ValueString())
	wan1.SetStaticSubnetMask(data.Wan1StaticSubnetMask.ValueString())
	wan1.SetStaticIp(data.Wan1StaticIp.ValueString())
	wan1.SetUsingStaticIp(data.Wan1UsingStaticIp.ValueBool())
	wan1.SetVlan(int32(data.Wan1Vlan.ValueInt64()))

	staticDNS := []string{}
	for _, dns := range data.Wan1StaticDns.Elements() {
		staticDNS = append(staticDNS, dns.String())
	}
	wan1.SetStaticDns(staticDNS)

	//wan2.SetStaticGatewayIp(data.Wan2StaticGatewayIp.ValueString())
	wan2.SetWanEnabled(data.Wan2WanEnabled.ValueString())
	//wan2.SetStaticSubnetMask(data.Wan2StaticSubnetMask.ValueString())
	//wan2.SetStaticIp(data.Wan2StaticIp.ValueString())
	wan2.SetUsingStaticIp(data.Wan2UsingStaticIp.ValueBool())
	wan2.SetVlan(int32(data.Wan2Vlan.ValueInt64()))

	//wan2StaticDNS := []string{}
	//for _, dns := range data.Wan2StaticDns.Elements() {
	//	wan2StaticDNS = append(wan2StaticDNS, dns.String())
	//}
	//wan2.SetStaticDns(wan2StaticDNS)

	deviceNetworkInterface := openApiClient.NewInlineObject14()
	deviceNetworkInterface.SetWan1(*wan1)
	deviceNetworkInterface.SetWan2(*wan2)

	_, httpResp, err := r.client.DevicesApi.UpdateDeviceManagementInterface(ctx, data.Serial.ValueString()).UpdateDeviceManagementInterface(*deviceNetworkInterface).Execute()
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

	wan1 := openApiClient.NewDevicesSerialManagementInterfaceWan1()
	wan2 := openApiClient.NewDevicesSerialManagementInterfaceWan2()
	wan1.SetStaticGatewayIp(data.Wan1StaticGatewayIp.ValueString())
	wan1.SetWanEnabled(data.Wan1WanEnabled.ValueString())
	wan1.SetStaticSubnetMask(data.Wan1StaticSubnetMask.ValueString())
	wan1.SetStaticIp(data.Wan1StaticIp.ValueString())
	wan1.SetUsingStaticIp(data.Wan1UsingStaticIp.ValueBool())
	wan1.SetVlan(int32(data.Wan1Vlan.ValueInt64()))

	staticDNS := []string{}
	for _, dns := range data.Wan1StaticDns.Elements() {
		staticDNS = append(staticDNS, dns.String())
	}
	wan1.SetStaticDns(staticDNS)

	//wan2.SetStaticGatewayIp(data.Wan2StaticGatewayIp.ValueString())
	wan2.SetWanEnabled(data.Wan2WanEnabled.ValueString())
	//wan2.SetStaticSubnetMask(data.Wan2StaticSubnetMask.ValueString())
	//wan2.SetStaticIp(data.Wan2StaticIp.ValueString())
	wan2.SetUsingStaticIp(data.Wan2UsingStaticIp.ValueBool())
	wan2.SetVlan(int32(data.Wan2Vlan.ValueInt64()))

	//wan2StaticDNS := []string{}
	//for _, dns := range data.Wan2StaticDns.Elements() {
	//	wan2StaticDNS = append(wan2StaticDNS, dns.String())
	//}
	//wan2.SetStaticDns(wan2StaticDNS)

	deviceNetworkInterface := openApiClient.NewInlineObject14()
	deviceNetworkInterface.SetWan1(*wan1)
	deviceNetworkInterface.SetWan2(*wan2)

	_, httpResp, err := r.client.DevicesApi.UpdateDeviceManagementInterface(ctx, data.Serial.ValueString()).UpdateDeviceManagementInterface(*deviceNetworkInterface).Execute()
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

	wan1 := openApiClient.NewDevicesSerialManagementInterfaceWan1()
	wan2 := openApiClient.NewDevicesSerialManagementInterfaceWan2()
	wan1.SetStaticGatewayIp(data.Wan1StaticGatewayIp.ValueString())
	wan1.SetWanEnabled(data.Wan1WanEnabled.ValueString())
	wan1.SetStaticSubnetMask(data.Wan1StaticSubnetMask.ValueString())
	wan1.SetStaticIp(data.Wan1StaticIp.ValueString())
	wan1.SetUsingStaticIp(data.Wan1UsingStaticIp.ValueBool())
	wan1.SetVlan(int32(data.Wan1Vlan.ValueInt64()))

	staticDNS := []string{}
	for _, dns := range data.Wan1StaticDns.Elements() {
		staticDNS = append(staticDNS, dns.String())
	}
	wan1.SetStaticDns(staticDNS)

	//wan2.SetStaticGatewayIp(data.Wan2StaticGatewayIp.ValueString())
	wan2.SetWanEnabled(data.Wan2WanEnabled.ValueString())
	//wan2.SetStaticSubnetMask(data.Wan2StaticSubnetMask.ValueString())
	//wan2.SetStaticIp(data.Wan2StaticIp.ValueString())
	wan2.SetUsingStaticIp(data.Wan2UsingStaticIp.ValueBool())
	wan2.SetVlan(int32(data.Wan2Vlan.ValueInt64()))

	//wan2StaticDNS := []string{}
	//for _, dns := range data.Wan2StaticDns.Elements() {
	//	wan2StaticDNS = append(wan2StaticDNS, dns.String())
	//}
	//wan2.SetStaticDns(wan2StaticDNS)

	deviceNetworkInterface := openApiClient.NewInlineObject14()
	deviceNetworkInterface.SetWan1(*wan1)
	deviceNetworkInterface.SetWan2(*wan2)

	_, httpResp, err := r.client.DevicesApi.UpdateDeviceManagementInterface(ctx, data.Serial.ValueString()).UpdateDeviceManagementInterface(*deviceNetworkInterface).Execute()
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

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

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
