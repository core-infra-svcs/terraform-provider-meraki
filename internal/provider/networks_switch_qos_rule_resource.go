package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"net/http"
	"strings"
	"time"

	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
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
var _ resource.Resource = &NetworksSwitchQosRuleResource{}
var _ resource.ResourceWithImportState = &NetworksSwitchQosRuleResource{}

func NewNetworksSwitchQosRuleResource() resource.Resource {
	return &NetworksSwitchQosRuleResource{}
}

// NetworksSwitchQosRuleResource defines the resource implementation.
type NetworksSwitchQosRuleResource struct {
	client *openApiClient.APIClient
}

// NetworksSwitchQosRuleResourceModel describes the resource data model.
type NetworksSwitchQosRuleResourceModel struct {
	Id           jsontypes.String  `tfsdk:"id" json:"-"`
	NetworkId    jsontypes.String  `tfsdk:"network_id" json:"network_id"`
	QosRulesId   jsontypes.String  `tfsdk:"qos_rule_id" json:"id"`
	Vlan         jsontypes.Int64   `tfsdk:"vlan" json:"vlan"`
	Dscp         jsontypes.Int64   `tfsdk:"dscp" json:"dscp"`
	DstPort      jsontypes.Float64 `tfsdk:"dst_port" json:"dstPort"`
	SrcPort      jsontypes.Float64 `tfsdk:"src_port" json:"srcPort"`
	DstPortRange jsontypes.String  `tfsdk:"dst_port_range" json:"dstPortRange"`
	Protocol     jsontypes.String  `tfsdk:"protocol" json:"protocol"`
	SrcPortRange jsontypes.String  `tfsdk:"src_port_range" json:"srcPortRange"`
}

func (r *NetworksSwitchQosRuleResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_switch_qos_rule"
}

func (r *NetworksSwitchQosRuleResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "NetworksSwitchQosRule resource for updating network switch qos rule.",
		Attributes: map[string]schema.Attribute{

			"id": schema.StringAttribute{
				MarkdownDescription: "Qos Rule Id",
				Computed:            true,
				Optional:            true,
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
			"qos_rule_id": schema.StringAttribute{
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

func (r *NetworksSwitchQosRuleResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NetworksSwitchQosRuleResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *NetworksSwitchQosRuleResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload := *openApiClient.NewCreateNetworkSwitchQosRuleRequest(int32(data.Vlan.ValueInt64()))

	// Set payload fields
	if !data.Dscp.IsUnknown() {
		payload.SetDscp(int32(data.Dscp.ValueInt64()))
	}
	if !data.DstPort.IsUnknown() {
		payload.SetDstPort(int32(data.DstPort.ValueFloat64()))
	}
	if !data.Protocol.IsUnknown() {
		payload.SetProtocol(data.Protocol.ValueString())
	}
	if !data.DstPortRange.IsUnknown() {
		payload.SetDstPortRange(data.DstPortRange.String())
	}
	if !data.SrcPortRange.IsUnknown() {
		payload.SetSrcPortRange(data.SrcPortRange.ValueString())
	}
	if !data.SrcPort.IsUnknown() {
		payload.SetSrcPort(int32(data.SrcPort.ValueFloat64()))
	}

	// Log the payload before sending
	tflog.Debug(ctx, "Creating QoS Rule with payload", map[string]interface{}{
		"payload": payload,
	})

	// Retry logic for 500 Internal Server Error
	const maxRetries500 = 3
	const retryDelay500 = 3 * time.Second

	// retry variables for all other errors
	maxRetries := r.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(r.client.GetConfig().Retry4xxErrorWaitTime)

	apiCall := func() (map[string]interface{}, *http.Response, error) {
		return r.client.QosRulesApi.CreateNetworkSwitchQosRule(context.Background(), data.NetworkId.ValueString()).CreateNetworkSwitchQosRuleRequest(payload).Execute()
	}

	var inlineResp map[string]interface{}
	var httpResp *http.Response
	var err error

	for attempt := 0; attempt < maxRetries500; attempt++ {
		inlineResp, httpResp, err = tools.CustomHttpRequestRetry[map[string]interface{}](ctx, maxRetries, retryDelay, apiCall)

		if err == nil && httpResp.StatusCode == 201 {
			// Success, break the retry loop
			break
		}

		if httpResp != nil && httpResp.StatusCode == 500 {
			tflog.Warn(ctx, "Received 500 Internal Server Error, retrying...", map[string]interface{}{
				"attempt": attempt + 1,
			})
			time.Sleep(retryDelay500)
		} else {
			// Break on other errors
			break
		}
	}

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
	if httpResp.StatusCode != 201 {
		resp.Diagnostics.AddError(
			"Unexpected HTTP Response Status Code",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
	}

	// Decode the response into the model
	if err = json.NewDecoder(httpResp.Body).Decode(&data); err != nil {
		resp.Diagnostics.AddError(
			"JSON Decode issue",
			fmt.Sprintf("%v", httpResp.StatusCode),
		)
		return
	}

	// Populate the data model with the response
	if iD := inlineResp["id"]; iD != nil {
		data.QosRulesId = jsontypes.StringValue(iD.(string))
	} else {
		data.QosRulesId = jsontypes.StringNull()
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

	data.Id = jsontypes.StringValue(data.NetworkId.ValueString() + "," + data.QosRulesId.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")
}

func (r *NetworksSwitchQosRuleResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetworksSwitchQosRuleResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.QosRulesId.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing QoS Rule Id",
			fmt.Sprintf("%v", data.QosRulesId.ValueString()),
		)
		return
	}

	maxRetries := r.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(r.client.GetConfig().Retry4xxErrorWaitTime)

	inlineResp, httpResp, err := tools.CustomHttpRequestRetry[*openApiClient.GetNetworkSwitchQosRule200Response](ctx, maxRetries, retryDelay, func() (*openApiClient.GetNetworkSwitchQosRule200Response, *http.Response, error) {
		inline, respHttp, err := r.client.QosRulesApi.GetNetworkSwitchQosRule(ctx, data.NetworkId.ValueString(), data.QosRulesId.ValueString()).Execute()
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

	id := inlineResp.GetId()
	if id == "" {
		resp.Diagnostics.AddError(
			"Missing QoS Rule Id from read payload",
			fmt.Sprintf("%v", inlineResp.GetId()),
		)
		return
	}

	data.QosRulesId = jsontypes.StringValue(id)

	if data.QosRulesId.IsUnknown() || data.QosRulesId.IsNull() {
		resp.Diagnostics.AddError(
			"Missing QoS Rule Id from read payload",
			fmt.Sprintf("%v", inlineResp.GetId()),
		)
		return
	}

	data.Id = jsontypes.StringValue(data.NetworkId.ValueString() + "," + data.QosRulesId.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *NetworksSwitchQosRuleResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	var stateData *NetworksSwitchQosRuleResourceModel

	// Read Terraform state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &stateData)...)

	var data *NetworksSwitchQosRuleResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	payload := *openApiClient.NewUpdateNetworkSwitchQosRuleRequest()

	if !data.Vlan.IsUnknown() {
		payload.SetVlan(int32(data.Vlan.ValueInt64()))
	}
	if !data.Dscp.IsUnknown() {
		payload.SetDscp(int32(data.Dscp.ValueInt64()))
	}
	if !data.DstPort.IsUnknown() {
		payload.SetDstPort(int32(data.DstPort.ValueFloat64()))
	}
	if !data.Protocol.IsUnknown() {
		payload.SetProtocol(data.Protocol.ValueString())
	}
	if !data.DstPortRange.IsUnknown() {
		payload.SetDstPortRange(data.DstPortRange.String())
	}
	if !data.SrcPortRange.IsUnknown() {
		payload.SetSrcPortRange(data.SrcPortRange.ValueString())
	}
	if !data.SrcPort.IsUnknown() {
		payload.SetSrcPort(int32(data.SrcPort.ValueFloat64()))
	}

	maxRetries := r.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(r.client.GetConfig().Retry4xxErrorWaitTime)

	inlineResp, httpResp, err := tools.CustomHttpRequestRetry[map[string]interface{}](ctx, maxRetries, retryDelay, func() (map[string]interface{}, *http.Response, error) {
		inline, respHttp, err := r.client.QosRulesApi.UpdateNetworkSwitchQosRule(ctx, data.NetworkId.ValueString(), stateData.QosRulesId.ValueString()).UpdateNetworkSwitchQosRuleRequest(payload).Execute()
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

	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
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

	data.Id = jsontypes.StringValue(data.NetworkId.ValueString() + "," + data.QosRulesId.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *NetworksSwitchQosRuleResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	var data *NetworksSwitchQosRuleResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	//httpResp, err := r.client.QosRulesApi.DeleteNetworkSwitchQosRule(ctx, data.NetworkId.ValueString(), data.QosRulesId.ValueString()).Execute()

	maxRetries := r.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(r.client.GetConfig().Retry4xxErrorWaitTime)

	_, httpResp, err := tools.CustomHttpRequestRetry[map[string]interface{}](ctx, maxRetries, retryDelay, func() (map[string]interface{}, *http.Response, error) {
		respHttp, err := r.client.QosRulesApi.DeleteNetworkSwitchQosRule(ctx, data.NetworkId.ValueString(), data.QosRulesId.ValueString()).Execute()
		return nil, respHttp, err
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

	if err != nil {
		resp.Diagnostics.AddError(
			"HTTP Client Failure",
			tools.HttpDiagnostics(httpResp),
		)
		return
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

func (r *NetworksSwitchQosRuleResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: network_id, qos_rule_id. Got: %q", req.ID),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("qos_rule_id"), idParts[1])...)
	if resp.Diagnostics.HasError() {
		return
	}
}
