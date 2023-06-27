package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &DevicesCellulargatewayLanResource{}
	_ resource.ResourceWithConfigure   = &DevicesCellulargatewayLanResource{}
	_ resource.ResourceWithImportState = &DevicesCellulargatewayLanResource{}
)

func NewDevicesCellularGatewayLanResource() resource.Resource {
	return &DevicesCellulargatewayLanResource{}
}

// DevicesCellulargatewayLanResource defines the resource implementation.
type DevicesCellulargatewayLanResource struct {
	client *openApiClient.APIClient
}

// DevicesCellulargatewayLanResourceModel describes the resource data model.
type DevicesCellulargatewayLanResourceModel struct {
	Id     jsontypes.String `tfsdk:"id"`
	Serial jsontypes.String `tfsdk:"serial"`

	DeviceName         jsontypes.String `tfsdk:"device_name" json:"deviceName"`
	DeviceLanIp        jsontypes.String `tfsdk:"device_lan_ip" json:"deviceLanIp"`
	DeviceSubnet       jsontypes.String `tfsdk:"device_subnet" json:"deviceSubnet"`
	FixedIpAssignments []struct {
		Mac  jsontypes.String `tfsdk:"mac" json:"mac"`
		Name jsontypes.String `tfsdk:"name" json:"name"`
		Ip   jsontypes.String `tfsdk:"ip" json:"ip"`
	} `tfsdk:"fixed_ip_assignments" json:"fixedIpAssignments"`
	ReservedIpRanges []struct {
		Start   jsontypes.String `tfsdk:"start" json:"start"`
		End     jsontypes.String `tfsdk:"end" json:"end"`
		Comment jsontypes.String `tfsdk:"comment" json:"comment"`
	} `tfsdk:"reserved_ip_ranges" json:"reservedIpRanges"`
}

func (r *DevicesCellulargatewayLanResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_devices_cellular_gateway_lan"
}

func (r *DevicesCellulargatewayLanResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{

		MarkdownDescription: "DevicesCellulargatewayLan",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.StringType,
			},
			"serial": schema.StringAttribute{
				MarkdownDescription: "Serial",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"device_name": schema.StringAttribute{
				MarkdownDescription: "The name of a device",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"fixed_ip_assignments": schema.SetNestedAttribute{
				Description: "list of all fixed IP assignments for a single MG'.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"name": schema.StringAttribute{
							MarkdownDescription: "A descriptive name of the assignment",
							Computed:            true,
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
						"mac": schema.StringAttribute{
							MarkdownDescription: "The MAC address of the server or device that hosts the internal resource that you wish to receive the specified IP address.",
							Computed:            true,
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
						"ip": schema.StringAttribute{
							MarkdownDescription: "The IP address you want to assign to a specific server or device\n\n.",
							Computed:            true,
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
					},
				},
			},
			"reserved_ip_ranges": schema.SetNestedAttribute{
				Description: "list of all reserved IP ranges for a single MG.",
				Optional:    true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"start": schema.StringAttribute{
							MarkdownDescription: "Starting IP included in the reserved range of IPs",
							Computed:            true,
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
						"end": schema.StringAttribute{
							MarkdownDescription: "Ending IP included in the reserved range of IPs.",
							Computed:            true,
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
						"comment": schema.StringAttribute{
							MarkdownDescription: "Comment explaining the reserved IP range.",
							Computed:            true,
							Optional:            true,
							CustomType:          jsontypes.StringType,
						},
					},
				},
			},
		},
	}
}

func (r *DevicesCellulargatewayLanResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DevicesCellulargatewayLanResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *DevicesCellulargatewayLanResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	object10 := openApiClient.NewInlineObject10()

	v := []openApiClient.DevicesSerialCellularGatewayLanReservedIpRanges{}

	for _, reservedIP := range data.ReservedIpRanges {
		ranges := openApiClient.NewDevicesSerialCellularGatewayLanReservedIpRanges(reservedIP.Start.ValueString(), reservedIP.End.ValueString(), reservedIP.Comment.ValueString())
		v = append(v, *ranges)
	}

	object10.SetReservedIpRanges(v)

	fixedAssignments := []openApiClient.DevicesSerialCellularGatewayLanFixedIpAssignments{}

	for _, fixedAssignment := range data.FixedIpAssignments {
		assignment := openApiClient.NewDevicesSerialCellularGatewayLanFixedIpAssignments(fixedAssignment.Ip.ValueString(), fixedAssignment.Mac.ValueString())
		fixedAssignments = append(fixedAssignments, *assignment)
	}
	object10.SetFixedIpAssignments(fixedAssignments)

	_, httpResp, err := r.client.CellularGatewayApi.UpdateDeviceCellularGatewayLan(ctx, data.Serial.ValueString()).UpdateDeviceCellularGatewayLan(*object10).Execute()

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

func (r *DevicesCellulargatewayLanResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *DevicesCellulargatewayLanResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := r.client.CellularGatewayApi.GetDeviceCellularGatewayLan(ctx, data.Serial.ValueString()).Execute()

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

func (r *DevicesCellulargatewayLanResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *DevicesCellulargatewayLanResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	object10 := openApiClient.NewInlineObject10()

	v := []openApiClient.DevicesSerialCellularGatewayLanReservedIpRanges{}

	for _, reservedIP := range data.ReservedIpRanges {
		ranges := openApiClient.NewDevicesSerialCellularGatewayLanReservedIpRanges(reservedIP.Start.ValueString(), reservedIP.End.ValueString(), reservedIP.Comment.ValueString())
		v = append(v, *ranges)
	}

	object10.SetReservedIpRanges(v)

	fixedAssignments := []openApiClient.DevicesSerialCellularGatewayLanFixedIpAssignments{}

	for _, fixedAssignment := range data.FixedIpAssignments {
		assignment := openApiClient.NewDevicesSerialCellularGatewayLanFixedIpAssignments(fixedAssignment.Ip.ValueString(), fixedAssignment.Mac.ValueString())
		fixedAssignments = append(fixedAssignments, *assignment)
	}
	object10.SetFixedIpAssignments(fixedAssignments)

	_, httpResp, err := r.client.CellularGatewayApi.UpdateDeviceCellularGatewayLan(ctx, data.Serial.ValueString()).UpdateDeviceCellularGatewayLan(*object10).Execute()

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

func (r *DevicesCellulargatewayLanResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *DevicesCellulargatewayLanResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	object10 := openApiClient.NewInlineObject10()

	_, httpResp, err := r.client.CellularGatewayApi.UpdateDeviceCellularGatewayLan(ctx, data.Serial.ValueString()).UpdateDeviceCellularGatewayLan(*object10).Execute()

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
	tflog.Trace(ctx, "removed resource")
}

func (r *DevicesCellulargatewayLanResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("serial"), req.ID)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
