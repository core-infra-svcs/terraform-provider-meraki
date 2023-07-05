package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
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
var _ resource.Resource = &NetworksSwitchQosRulesResource{}
var _ resource.ResourceWithImportState = &NetworksSwitchQosRulesResource{}

func NewNetworksSwitchQosRulesResource() resource.Resource {
	return &NetworksSwitchQosRulesResource{}
}

// NetworksSwitchQosRulesResource defines the resource implementation.
type NetworksSwitchQosRulesResource struct {
	client *openApiClient.APIClient
}

// NetworksSwitchQosRulesResourceModel describes the resource data model.
type NetworksSwitchQosRulesResourceModel struct {
	Id           jsontypes.String  `tfsdk:"id"`
	NetworkId    jsontypes.String  `tfsdk:"network_id" json:"network_id"`
	QosRulesId   jsontypes.String  `tfsdk:"qos_rules_id" json:"id"`
	Vlan         jsontypes.Int64   `tfsdk:"vlan" json:"vlan"`
	Dscp         jsontypes.Int64   `tfsdk:"dscp" json:"dscp"`
	DstPort      jsontypes.Float64 `tfsdk:"dst_port" json:"dstPort"`
	SrcPort      jsontypes.Float64 `tfsdk:"src_port" json:"srcPort"`
	DstPortRange jsontypes.String  `tfsdk:"dst_port_range" json:"dstPortRange"`
	Protocol     jsontypes.String  `tfsdk:"protocol" json:"protocol"`
	SrcPortRange jsontypes.String  `tfsdk:"src_port_range" json:"srcPortRange"`
}

func (r *NetworksSwitchQosRulesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_switch_qos_rules"
}

func (r *NetworksSwitchQosRulesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "NetworksSwitchQosRules resource for updating network switch qos rules.",
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
			"qos_rules_id": schema.StringAttribute{
				MarkdownDescription: "Qos Rules Id",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vlan": schema.Int64Attribute{
				MarkdownDescription: "The VLAN of the incoming packet. A null value will match any VLAN.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.Int64Type,
			},
			"dscp": schema.Int64Attribute{
				MarkdownDescription: "DSCP tag. Set this to -1 to trust incoming DSCP. Default value is 0.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.Int64Type,
			},
			"dst_port": schema.Float64Attribute{
				MarkdownDescription: "The destination port of the incoming packet. Applicable only if protocol is TCP or UDP.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.Float64Type,
			},
			"src_port": schema.Float64Attribute{
				MarkdownDescription: "The source port of the incoming packet. Applicable only if protocol is TCP or UDP.",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.Float64Type,
			},
			"dst_port_range": schema.StringAttribute{
				MarkdownDescription: "The destination port range of the incoming packet. Applicable only if protocol is set to TCP or UDP. Example: 70-80",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"protocol": schema.StringAttribute{
				MarkdownDescription: "The protocol of the incoming packet. Can be one of ANY, TCP or UDP. Default value is ANY",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"src_port_range": schema.StringAttribute{
				MarkdownDescription: "The source port range of the incoming packet. Applicable only if protocol is set to TCP or UDP. Example: 70-80",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
		},
	}
}

func (r *NetworksSwitchQosRulesResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NetworksSwitchQosRulesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *NetworksSwitchQosRulesResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	createNetworksSwitchQosRules := *openApiClient.NewCreateNetworkSwitchQosRuleRequest(int32(data.Vlan.ValueInt64()))

	if !data.Dscp.IsUnknown() {
		createNetworksSwitchQosRules.SetDscp(int32(data.Dscp.ValueInt64()))
	}
	if !data.DstPort.IsUnknown() {
		createNetworksSwitchQosRules.SetDstPort(int32(data.DstPort.ValueFloat64()))
	}
	if !data.Protocol.IsUnknown() {
		createNetworksSwitchQosRules.SetProtocol(data.Protocol.ValueString())
	}
	if !data.DstPortRange.IsUnknown() {
		createNetworksSwitchQosRules.SetDstPortRange(data.DstPortRange.String())
	}
	if !data.SrcPortRange.IsUnknown() {
		createNetworksSwitchQosRules.SetSrcPortRange(data.SrcPortRange.ValueString())
	}
	if !data.SrcPort.IsUnknown() {
		createNetworksSwitchQosRules.SetSrcPort(int32(data.SrcPort.ValueFloat64()))
	}

	inlineResp, httpResp, err := r.client.QosRulesApi.CreateNetworkSwitchQosRule(context.Background(), data.NetworkId.ValueString()).CreateNetworkSwitchQosRuleRequest(createNetworksSwitchQosRules).Execute()
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
	if httpResp.StatusCode != 201 {
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

	// save inlineResp data into Terraform state.
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON Decode issue",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
		return
	}

	if dstPortRange := inlineResp["dstPortRange"]; dstPortRange != nil {
		data.DstPortRange = jsontypes.StringValue(dstPortRange.(string))
	} else {
		data.DstPortRange = jsontypes.StringNull()
	}
	if srcPortRange := inlineResp["srcPortRange"]; srcPortRange != nil {
		data.SrcPortRange = jsontypes.StringValue(srcPortRange.(string))
	} else {
		data.SrcPortRange = jsontypes.StringNull()
	}
	if srcPort := inlineResp["srcPort"]; srcPort != nil {
		data.SrcPort = jsontypes.Float64Value(srcPort.(float64))
	} else {
		data.SrcPort = jsontypes.Float64Null()
	}
	if dstPort := inlineResp["dstPort"]; dstPort != nil {
		data.DstPort = jsontypes.Float64Value(dstPort.(float64))
	} else {
		data.DstPort = jsontypes.Float64Null()
	}
	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")
}

func (r *NetworksSwitchQosRulesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetworksSwitchQosRulesResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	inlineResp, httpResp, err := r.client.QosRulesApi.GetNetworkSwitchQosRule(ctx, data.NetworkId.ValueString(), data.QosRulesId.ValueString()).Execute()
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

	// save inlineResp data into Terraform state.
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON Decode issue",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
		return
	}

	data.DstPortRange = jsontypes.StringValue(inlineResp.GetDstPortRange())

	data.SrcPortRange = jsontypes.StringValue(inlineResp.GetSrcPortRange())

	data.SrcPort = jsontypes.Float64Value(float64(inlineResp.GetSrcPort()))

	data.DstPort = jsontypes.Float64Value(float64(inlineResp.GetDstPort()))

	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *NetworksSwitchQosRulesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var data *NetworksSwitchQosRulesResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	updateNetworksSwitchQosRules := *openApiClient.NewUpdateNetworkSwitchQosRuleRequest()

	if !data.Vlan.IsUnknown() {
		updateNetworksSwitchQosRules.SetVlan(int32(data.Vlan.ValueInt64()))
	}
	if !data.Dscp.IsUnknown() {
		updateNetworksSwitchQosRules.SetDscp(int32(data.Dscp.ValueInt64()))
	}
	if !data.DstPort.IsUnknown() {
		updateNetworksSwitchQosRules.SetDstPort(int32(data.DstPort.ValueFloat64()))
	}
	if !data.Protocol.IsUnknown() {
		updateNetworksSwitchQosRules.SetProtocol(data.Protocol.ValueString())
	}
	if !data.DstPortRange.IsUnknown() {
		updateNetworksSwitchQosRules.SetDstPortRange(data.DstPortRange.String())
	}
	if !data.SrcPortRange.IsUnknown() {
		updateNetworksSwitchQosRules.SetSrcPortRange(data.SrcPortRange.ValueString())
	}
	if !data.SrcPort.IsUnknown() {
		updateNetworksSwitchQosRules.SetSrcPort(int32(data.SrcPort.ValueFloat64()))
	}

	inlineResp, httpResp, err := r.client.QosRulesApi.UpdateNetworkSwitchQosRule(context.Background(), data.NetworkId.ValueString(), data.QosRulesId.ValueString()).UpdateNetworkSwitchQosRuleRequest(updateNetworksSwitchQosRules).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to update resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// collect diagnostics
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
	}

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

	// save inlineResp data into Terraform state.
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON Decode issue",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
		return
	}

	if dstPortRange := inlineResp["dstPortRange"]; dstPortRange != nil {
		data.DstPortRange = jsontypes.StringValue(dstPortRange.(string))
	} else {
		data.DstPortRange = jsontypes.StringNull()
	}
	if srcPortRange := inlineResp["srcPortRange"]; srcPortRange != nil {
		data.SrcPortRange = jsontypes.StringValue(srcPortRange.(string))
	} else {
		data.SrcPortRange = jsontypes.StringNull()
	}
	if srcPort := inlineResp["srcPort"]; srcPort != nil {
		data.SrcPort = jsontypes.Float64Value(srcPort.(float64))
	} else {
		data.SrcPort = jsontypes.Float64Null()
	}
	if dstPort := inlineResp["dstPort"]; dstPort != nil {
		data.DstPort = jsontypes.Float64Value(dstPort.(float64))
	} else {
		data.DstPort = jsontypes.Float64Null()
	}
	data.Id = jsontypes.StringValue("example-id")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *NetworksSwitchQosRulesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data *NetworksSwitchQosRulesResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	httpResp, err := r.client.QosRulesApi.DeleteNetworkSwitchQosRule(ctx, data.NetworkId.ValueString(), data.QosRulesId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Delete resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}

	// collect diagnostics
	if httpResp != nil {
		tools.CollectHttpDiagnostics(ctx, &resp.Diagnostics, httpResp)
	}

	if httpResp.StatusCode != 204 {
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

func (r *NetworksSwitchQosRulesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: network_id, quos_rule_id. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("qos_rule_id"), idParts[1])...)
	if resp.Diagnostics.HasError() {
		return
	}
}
