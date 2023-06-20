package provider

import (
	"context"
	"encoding/json"
	"fmt"
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
	"strings"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &NetworksWirelessSsidsVpnResource{}
	_ resource.ResourceWithConfigure   = &NetworksWirelessSsidsVpnResource{}
	_ resource.ResourceWithImportState = &NetworksWirelessSsidsVpnResource{}
)

func NewNetworksWirelessSsidsVpnResource() resource.Resource {
	return &NetworksWirelessSsidsVpnResource{}
}

// NetworksWirelessSsidsVpnResource defines the resource implementation.
type NetworksWirelessSsidsVpnResource struct {
	client *openApiClient.APIClient
}

// NetworksWirelessSsidsVpnResourceModel describes the resource data model.
type NetworksWirelessSsidsVpnResourceModel struct {
	Id        jsontypes.String `tfsdk:"id"`
	NetworkId jsontypes.String `tfsdk:"network_id"`
	Number    jsontypes.String `tfsdk:"number"`

	Concentrator Concentrator `tfsdk:"concentrator" json:"concentrator"`
	Failover     Failover     `tfsdk:"failover" json:"failover"`
	SplitTunnel  SplitTunnel  `tfsdk:"split_tunnel" json:"splitTunnel"`
}

type Concentrator struct {
	NetworkId jsontypes.String `tfsdk:"network_id" json:"networkId"`
	VlanId    jsontypes.Int64  `tfsdk:"vlan_id" json:"vlanId"`
	Name      jsontypes.String `tfsdk:"name" json:"name"`
}

type SplitTunnel struct {
	Enabled jsontypes.Bool `tfsdk:"enabled" json:"enabled"`
	Rules   []TunnelRule   `tfsdk:"rules" json:"rules"`
}

type TunnelRule struct {
	Protocol jsontypes.String `tfsdk:"protocol" json:"protocol,omitempty"`
	DestCidr jsontypes.String `tfsdk:"dest_cidr" json:"destCidr"`
	DestPort jsontypes.String `tfsdk:"dest_port" json:"destPort"`
	Policy   jsontypes.String `tfsdk:"policy" json:"policy"`
	Comment  jsontypes.String `tfsdk:"comment" json:"comment"`
}

type Failover struct {
	RequestIp         jsontypes.String `tfsdk:"request_ip" json:"requestIp"`
	HeartbeatInterval jsontypes.Int64  `tfsdk:"heartbeat_interval" json:"heartbeatInterval"`
	IdleTimeout       jsontypes.Int64  `tfsdk:"idle_timeout" json:"idleTimeout"`
}

func (r *NetworksWirelessSsidsVpnResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_wireless_ssids_vpn"
}

func (r *NetworksWirelessSsidsVpnResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{

		MarkdownDescription: "NetworksWirelessSsidsVpn",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.StringType,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Organization ID",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 31),
				},
			},
			"number": schema.StringAttribute{
				MarkdownDescription: "Organization ID",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"concentrator": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"network_id": schema.StringAttribute{
						MarkdownDescription: "The first IP in the reserved range",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.StringType,
					},
					"vlan_id": schema.Int64Attribute{
						MarkdownDescription: "The last IP in the reserved range",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.Int64Type,
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "A text comment for the reserved range",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.StringType,
					},
				},
				Required: true,
			},
			"failover": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"request_ip": schema.StringAttribute{
						MarkdownDescription: "The first IP in the reserved range",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.StringType,
					},
					"idle_timeout": schema.Int64Attribute{
						MarkdownDescription: "The last IP in the reserved range",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.Int64Type,
					},
					"heartbeat_interval": schema.Int64Attribute{
						MarkdownDescription: "The last IP in the reserved range",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.Int64Type,
					},
				},
				Required: true,
			},
			"split_tunnel": schema.SingleNestedAttribute{
				Attributes: map[string]schema.Attribute{
					"enabled": schema.BoolAttribute{
						MarkdownDescription: "The first IP in the reserved range",
						Optional:            true,
						Computed:            true,
						CustomType:          jsontypes.BoolType,
					},
					"rules": schema.SetNestedAttribute{
						Optional:    true,
						Computed:    false,
						Description: "",
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"protocol": schema.StringAttribute{
									MarkdownDescription: "Organization ID",
									Optional:            true,
									Computed:            true,
									CustomType:          jsontypes.StringType,
								},
								"dest_cidr": schema.StringAttribute{
									MarkdownDescription: "Organization ID",
									Optional:            true,
									Computed:            true,
									CustomType:          jsontypes.StringType,
								},
								"dest_port": schema.StringAttribute{
									MarkdownDescription: "Organization ID",
									Optional:            true,
									Computed:            true,
									CustomType:          jsontypes.StringType,
								},
								"policy": schema.StringAttribute{
									MarkdownDescription: "Organization ID",
									Optional:            true,
									Computed:            true,
									CustomType:          jsontypes.StringType,
								},
								"comment": schema.StringAttribute{
									MarkdownDescription: "Organization ID",
									Optional:            true,
									Computed:            true,
									CustomType:          jsontypes.StringType,
								},
							},
						},
					},
				},
				Required: true,
			},
		},
	}
}

func (r *NetworksWirelessSsidsVpnResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NetworksWirelessSsidsVpnResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *NetworksWirelessSsidsVpnResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	object166 := openApiClient.NewInlineObject166()
	concentrator := *openApiClient.NewNetworksNetworkIdWirelessSsidsNumberVpnConcentrator()
	concentrator.SetVlanId(int32(data.Concentrator.VlanId.ValueInt64()))
	concentrator.SetNetworkId(data.Concentrator.NetworkId.ValueString())
	object166.SetConcentrator(concentrator)
	failover := *openApiClient.NewNetworksNetworkIdWirelessSsidsNumberVpnFailover()
	failover.SetIdleTimeout(int32(data.Failover.IdleTimeout.ValueInt64()))
	failover.SetHeartbeatInterval(int32(data.Failover.HeartbeatInterval.ValueInt64()))
	failover.SetRequestIp(data.Failover.RequestIp.ValueString())
	object166.SetFailover(failover)
	splitTunnel := *openApiClient.NewNetworksNetworkIdWirelessSsidsNumberVpnSplitTunnel()
	splitTunnel.SetEnabled(data.SplitTunnel.Enabled.ValueBool())
	splitTunnels := []openApiClient.NetworksNetworkIdWirelessSsidsNumberVpnSplitTunnelRules{}
	for _, rule := range data.SplitTunnel.Rules {
		splitRule := *openApiClient.NewNetworksNetworkIdWirelessSsidsNumberVpnSplitTunnelRules(rule.DestCidr.ValueString(), rule.Policy.ValueString())
		splitRule.SetComment(rule.Comment.ValueString())
		splitRule.SetDestPort(rule.DestPort.ValueString())
		splitRule.SetProtocol(rule.Protocol.ValueString())
		splitTunnels = append(splitTunnels, splitRule)
	}
	splitTunnel.SetRules(splitTunnels)
	object166.SetSplitTunnel(splitTunnel)

	_, httpResp, err := r.client.SsidsApi.UpdateNetworkWirelessSsidVpn(ctx, data.NetworkId.ValueString(), data.Number.ValueString()).UpdateNetworkWirelessSsidVpn(*object166).Execute()

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

func (r *NetworksWirelessSsidsVpnResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetworksWirelessSsidsVpnResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := r.client.SsidsApi.GetNetworkWirelessSsidVpn(ctx, data.NetworkId.ValueString(), data.Number.ValueString()).Execute()
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

func (r *NetworksWirelessSsidsVpnResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *NetworksWirelessSsidsVpnResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	object166 := openApiClient.NewInlineObject166()
	concentrator := *openApiClient.NewNetworksNetworkIdWirelessSsidsNumberVpnConcentrator()
	concentrator.SetVlanId(int32(data.Concentrator.VlanId.ValueInt64()))
	concentrator.SetNetworkId(data.Concentrator.NetworkId.ValueString())
	object166.SetConcentrator(concentrator)
	failover := *openApiClient.NewNetworksNetworkIdWirelessSsidsNumberVpnFailover()
	failover.SetIdleTimeout(int32(data.Failover.IdleTimeout.ValueInt64()))
	failover.SetHeartbeatInterval(int32(data.Failover.HeartbeatInterval.ValueInt64()))
	failover.SetRequestIp(data.Failover.RequestIp.ValueString())
	object166.SetFailover(failover)
	splitTunnel := *openApiClient.NewNetworksNetworkIdWirelessSsidsNumberVpnSplitTunnel()
	splitTunnel.SetEnabled(data.SplitTunnel.Enabled.ValueBool())
	splitTunnels := []openApiClient.NetworksNetworkIdWirelessSsidsNumberVpnSplitTunnelRules{}
	for _, rule := range data.SplitTunnel.Rules {
		splitRule := *openApiClient.NewNetworksNetworkIdWirelessSsidsNumberVpnSplitTunnelRules(rule.DestCidr.ValueString(), rule.Policy.ValueString())
		splitRule.SetComment(rule.Comment.ValueString())
		splitRule.SetDestPort(rule.DestPort.ValueString())
		splitRule.SetProtocol(rule.Protocol.ValueString())
		splitTunnels = append(splitTunnels, splitRule)
	}
	splitTunnel.SetRules(splitTunnels)
	object166.SetSplitTunnel(splitTunnel)

	_, httpResp, err := r.client.SsidsApi.UpdateNetworkWirelessSsidVpn(ctx, data.NetworkId.ValueString(), data.Number.ValueString()).UpdateNetworkWirelessSsidVpn(*object166).Execute()

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

func (r *NetworksWirelessSsidsVpnResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *NetworksWirelessSsidsVpnResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	object166 := openApiClient.NewInlineObject166()
	concentrator := *openApiClient.NewNetworksNetworkIdWirelessSsidsNumberVpnConcentrator()
	concentrator.SetVlanId(int32(data.Concentrator.VlanId.ValueInt64()))
	concentrator.SetNetworkId(data.Concentrator.NetworkId.ValueString())
	object166.SetConcentrator(concentrator)
	failover := *openApiClient.NewNetworksNetworkIdWirelessSsidsNumberVpnFailover()
	failover.SetIdleTimeout(int32(data.Failover.IdleTimeout.ValueInt64()))
	failover.SetHeartbeatInterval(int32(data.Failover.HeartbeatInterval.ValueInt64()))
	failover.SetRequestIp(data.Failover.RequestIp.ValueString())
	object166.SetFailover(failover)
	splitTunnel := *openApiClient.NewNetworksNetworkIdWirelessSsidsNumberVpnSplitTunnel()
	splitTunnel.SetEnabled(data.SplitTunnel.Enabled.ValueBool())
	splitTunnels := []openApiClient.NetworksNetworkIdWirelessSsidsNumberVpnSplitTunnelRules{}
	for _, rule := range data.SplitTunnel.Rules {
		splitRule := *openApiClient.NewNetworksNetworkIdWirelessSsidsNumberVpnSplitTunnelRules(rule.DestCidr.ValueString(), rule.Policy.ValueString())
		splitRule.SetComment(rule.Comment.ValueString())
		splitRule.SetDestPort(rule.DestPort.ValueString())
		splitRule.SetProtocol(rule.Protocol.ValueString())
		splitTunnels = append(splitTunnels, splitRule)
	}
	splitTunnel.SetRules(splitTunnels)
	object166.SetSplitTunnel(splitTunnel)

	_, httpResp, err := r.client.SsidsApi.UpdateNetworkWirelessSsidVpn(ctx, data.NetworkId.ValueString(), data.Number.ValueString()).UpdateNetworkWirelessSsidVpn(*object166).Execute()

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
	tflog.Trace(ctx, "deleted resource")
}

func (r *NetworksWirelessSsidsVpnResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: organization_id, admin_id. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("number"), idParts[1])...)
	if resp.Diagnostics.HasError() {
		return
	}
}
