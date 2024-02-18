package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &DevicesCellularGatewayLanResource{}
	_ resource.ResourceWithConfigure   = &DevicesCellularGatewayLanResource{}
	_ resource.ResourceWithImportState = &DevicesCellularGatewayLanResource{}
)

func NewDevicesCellularGatewayLanResource() resource.Resource {
	return &DevicesCellularGatewayLanResource{}
}

// DevicesCellularGatewayLanResource defines the resource implementation.
type DevicesCellularGatewayLanResource struct {
	client *openApiClient.APIClient
}

// DevicesCellularGatewayLanResourceModel describes the resource data model.
type DevicesCellularGatewayLanResourceModel struct {
	Id                 jsontypes.String `tfsdk:"id"`
	Serial             jsontypes.String `tfsdk:"serial"`
	DeviceName         jsontypes.String `tfsdk:"device_name" json:"deviceName"`
	DeviceLanIp        jsontypes.String `tfsdk:"device_lan_ip" json:"deviceLanIp"`
	DeviceSubnet       jsontypes.String `tfsdk:"device_subnet" json:"deviceSubnet"`
	FixedIpAssignments types.Set        `tfsdk:"fixed_ip_assignments" json:"fixedIpAssignments"`
	ReservedIpRanges   types.Set        `tfsdk:"reserved_ip_ranges" json:"reservedIpRanges"`
}

type DevicesCellularGatewayLanResourceModelFixedIpAssignments struct {
	Mac  jsontypes.String `tfsdk:"mac" json:"mac"`
	Name jsontypes.String `tfsdk:"name" json:"name"`
	Ip   jsontypes.String `tfsdk:"ip" json:"ip"`
}

type DevicesCellularGatewayLanResourceModelReservedIpRanges struct {
	Start   jsontypes.String `tfsdk:"start" json:"start"`
	End     jsontypes.String `tfsdk:"end" json:"end"`
	Comment jsontypes.String `tfsdk:"comment" json:"comment"`
}

func (r *DevicesCellularGatewayLanResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_devices_cellular_gateway_lan"
}

func (r *DevicesCellularGatewayLanResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{

		MarkdownDescription: "Manage Cellular Gateway Lan",
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
			"device_lan_ip": schema.StringAttribute{
				MarkdownDescription: "LAN IP of the device",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"device_subnet": schema.StringAttribute{
				MarkdownDescription: "Subnet of the device",
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

func (r *DevicesCellularGatewayLanResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *DevicesCellularGatewayLanResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *DevicesCellularGatewayLanResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	payload, payloadDiag := DevicesCellularGatewayLanResourcePayload(ctx, data)
	if payloadDiag.HasError() {
		resp.Diagnostics.AddError("Resource Payload Error", fmt.Sprintf("\n%v", payloadDiag))
		return
	}

	_, httpResp, err := r.client.CellularGatewayApi.UpdateDeviceCellularGatewayLan(ctx, data.Serial.ValueString()).UpdateDeviceCellularGatewayLanRequest(payload).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Create HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
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

	// If there were any errors up to this point, log the plan data and return.
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
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

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")
}

func (r *DevicesCellularGatewayLanResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *DevicesCellularGatewayLanResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := r.client.CellularGatewayApi.GetDeviceCellularGatewayLan(ctx, data.Serial.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Read HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
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

	// If there were any errors up to this point, log the plan data and return.
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
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

func (r *DevicesCellularGatewayLanResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *DevicesCellularGatewayLanResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	payload, payloadDiag := DevicesCellularGatewayLanResourcePayload(ctx, data)
	if payloadDiag.HasError() {
		resp.Diagnostics.AddError("Resource Payload Error", fmt.Sprintf("\n%v", payloadDiag))
		return
	}

	_, httpResp, err := r.client.CellularGatewayApi.UpdateDeviceCellularGatewayLan(ctx, data.Serial.ValueString()).UpdateDeviceCellularGatewayLanRequest(payload).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Update HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
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

	// If there were any errors up to this point, log the plan data and return.
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
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

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *DevicesCellularGatewayLanResource)3A	1 Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *DevicesCellularGatewayLanResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	payload := *openApiClient.NewUpdateDeviceCellularGatewayLanRequest()

	_, httpResp, err := r.client.CellularGatewayApi.UpdateDeviceCellularGatewayLan(ctx, data.Serial.ValueString()).UpdateDeviceCellularGatewayLanRequest(payload).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Delete HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
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

	// If there were any errors up to this point, log the plan data and return.
	if resp.Diagnostics.HasError() {
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%v", data))
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

	resp.State.RemoveResource(ctx)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")
}

func (r *DevicesCellularGatewayLanResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("serial"), req.ID)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

func DevicesCellularGatewayLanResourcePayload(ctx context.Context, data *DevicesCellularGatewayLanResourceModel) (openApiClient.UpdateDeviceCellularGatewayLanRequest, diag.Diagnostics) {

	payload := *openApiClient.NewUpdateDeviceCellularGatewayLanRequest()

	// Reserved IP Ranges
	if !data.ReservedIpRanges.IsUnknown() && !data.ReservedIpRanges.IsNull() {
		var reservedIpRanges []openApiClient.UpdateDeviceCellularGatewayLanRequestReservedIpRangesInner
		diags := data.ReservedIpRanges.ElementsAs(ctx, &reservedIpRanges, false)
		if diags.HasError() {
			return payload, diags
		}
		payload.SetReservedIpRanges(reservedIpRanges)
	}

	// Fixed IP Assignment
	if !data.FixedIpAssignments.IsUnknown() && !data.FixedIpAssignments.IsNull() {
		var fixedAssignments []openApiClient.UpdateDeviceCellularGatewayLanRequestFixedIpAssignmentsInner

		diags := data.FixedIpAssignments.ElementsAs(ctx, &fixedAssignments, false)
		if diags.HasError() {
			return payload, diags
		}
		payload.SetFixedIpAssignments(fixedAssignments)
	}

	return payload, nil
}
