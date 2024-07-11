package networks

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &NetworksNetflowResource{}
var _ resource.ResourceWithImportState = &NetworksNetflowResource{}

func NewNetworksNetflowResource() resource.Resource {
	return &NetworksNetflowResource{}
}

// NetworksNetflowResource defines the resource implementation.
type NetworksNetflowResource struct {
	client *openApiClient.APIClient
}

// NetworksNetflowResourceModel describes the resource data model.
type NetworksNetflowResourceModel struct {
	Id               jsontypes.String `tfsdk:"id"`
	NetworkId        jsontypes.String `tfsdk:"network_id" json:"network_id"`
	ReportingEnabled jsontypes.Bool   `tfsdk:"reporting_enabled" json:"reportingEnabled"`
	CollectorIp      jsontypes.String `tfsdk:"collector_ip" json:"collectorIp"`
	CollectorPort    jsontypes.Int64  `tfsdk:"collector_port" json:"collectorPort"`
	EtaEnabled       jsontypes.Bool   `tfsdk:"eta_enabled" json:"etaEnabled"`
	EtaDstPort       jsontypes.Int64  `tfsdk:"eta_dst_port" json:"etaDstPort"`
}

func (r *NetworksNetflowResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_netflow"
}

func (r *NetworksNetflowResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "NetworksNetflow resource for updating networks netflow.",
		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{
				MarkdownDescription: "Example identifier",
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Network Id",
				Required:            true,
				CustomType:          jsontypes.StringType,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 31),
				},
			},
			"reporting_enabled": schema.BoolAttribute{
				MarkdownDescription: "Boolean indicating whether NetFlow traffic reporting is enabled (true) or disabled (false).",
				Required:            true,
				CustomType:          jsontypes.BoolType,
			},
			"collector_ip": schema.StringAttribute{
				MarkdownDescription: "The IPv4 address of the NetFlow collector.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"collector_port": schema.Int64Attribute{
				MarkdownDescription: "The port that the NetFlow collector will be listening on.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.Int64Type,
			},
			"eta_enabled": schema.BoolAttribute{
				MarkdownDescription: "Boolean indicating whether Encrypted Traffic Analytics is enabled (true) or disabled (false).",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.BoolType,
			},
			"eta_dst_port": schema.Int64Attribute{
				MarkdownDescription: "The port that the Encrypted Traffic Analytics collector will be listening on.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.Int64Type,
			},
		},
	}
}

func (r *NetworksNetflowResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NetworksNetflowResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *NetworksNetflowResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworkNetflow := *openApiClient.NewUpdateNetworkNetflowRequest()
	if !data.CollectorIp.IsUnknown() {
		updateNetworkNetflow.SetCollectorIp(data.CollectorIp.ValueString())
	}
	if !data.ReportingEnabled.IsUnknown() {
		updateNetworkNetflow.SetReportingEnabled(data.ReportingEnabled.ValueBool())
	}
	if !data.CollectorPort.IsUnknown() {
		updateNetworkNetflow.SetCollectorPort(int32(data.CollectorPort.ValueInt64()))
	}
	if !data.EtaEnabled.IsUnknown() {
		updateNetworkNetflow.SetEtaEnabled(data.EtaEnabled.ValueBool())
	}
	if !data.EtaDstPort.IsUnknown() {
		updateNetworkNetflow.SetEtaDstPort(int32(data.EtaDstPort.ValueInt64()))
	}

	_, httpResp, err := r.client.NetworksApi.UpdateNetworkNetflow(ctx, data.NetworkId.ValueString()).UpdateNetworkNetflowRequest(updateNetworkNetflow).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
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

	if data.CollectorIp.IsUnknown() {
		data.CollectorIp = jsontypes.StringNull()
	}
	if data.CollectorPort.IsUnknown() {
		data.CollectorPort = jsontypes.Int64Null()
	}

	if data.EtaDstPort.IsUnknown() {
		data.EtaDstPort = jsontypes.Int64Null()
	}
	if data.EtaEnabled.IsUnknown() {
		data.EtaEnabled = jsontypes.BoolNull()
	}

	data.Id = jsontypes.StringValue(data.NetworkId.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")
}

func (r *NetworksNetflowResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetworksNetflowResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := r.client.NetworksApi.GetNetworkNetflow(ctx, data.NetworkId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
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

	// Save data into Terraform state
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}

	data.Id = jsontypes.StringValue(data.NetworkId.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *NetworksNetflowResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data *NetworksNetflowResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworkNetflow := *openApiClient.NewUpdateNetworkNetflowRequest()
	if !data.CollectorIp.IsUnknown() {
		updateNetworkNetflow.SetCollectorIp(data.CollectorIp.ValueString())
	}
	if !data.ReportingEnabled.IsUnknown() {
		updateNetworkNetflow.SetReportingEnabled(data.ReportingEnabled.ValueBool())
	}
	if !data.CollectorPort.IsUnknown() {
		updateNetworkNetflow.SetCollectorPort(int32(data.CollectorPort.ValueInt64()))
	}
	if !data.EtaEnabled.IsUnknown() {
		updateNetworkNetflow.SetEtaEnabled(data.EtaEnabled.ValueBool())
	}
	if !data.EtaDstPort.IsUnknown() {
		updateNetworkNetflow.SetEtaDstPort(int32(data.EtaDstPort.ValueInt64()))
	}

	_, httpResp, err := r.client.NetworksApi.UpdateNetworkNetflow(ctx, data.NetworkId.ValueString()).UpdateNetworkNetflowRequest(updateNetworkNetflow).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
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

	if data.CollectorIp.IsUnknown() {
		data.CollectorIp = jsontypes.StringNull()
	}
	if data.CollectorPort.IsUnknown() {
		data.CollectorPort = jsontypes.Int64Null()
	}

	if data.EtaDstPort.IsUnknown() {
		data.EtaDstPort = jsontypes.Int64Null()
	}
	if data.EtaEnabled.IsUnknown() {
		data.EtaEnabled = jsontypes.BoolNull()
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *NetworksNetflowResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data *NetworksNetflowResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworkNetflow := *openApiClient.NewUpdateNetworkNetflowRequest()
	updateNetworkNetflow.SetReportingEnabled(false)
	updateNetworkNetflow.CollectorPort = nil
	updateNetworkNetflow.CollectorIp = nil
	updateNetworkNetflow.SetEtaEnabled(false)
	updateNetworkNetflow.EtaDstPort = nil

	_, httpResp, err := r.client.NetworksApi.UpdateNetworkNetflow(ctx, data.NetworkId.ValueString()).UpdateNetworkNetflowRequest(updateNetworkNetflow).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
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
		resp.Diagnostics.AddError("Plan Data", fmt.Sprintf("\n%s", data))
		return
	}

	resp.State.RemoveResource(ctx)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")

}

func (r *NetworksNetflowResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), req.ID)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
