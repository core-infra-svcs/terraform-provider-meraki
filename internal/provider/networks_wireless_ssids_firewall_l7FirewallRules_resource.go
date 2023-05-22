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
	_ resource.Resource                = &NetworksWirelessSsidsFirewallL7firewallrulesResource{}
	_ resource.ResourceWithConfigure   = &NetworksWirelessSsidsFirewallL7firewallrulesResource{}
	_ resource.ResourceWithImportState = &NetworksWirelessSsidsFirewallL7firewallrulesResource{}
)

func NewNetworksWirelessSsidsFirewallL7firewallrulesResource() resource.Resource {
	return &NetworksWirelessSsidsFirewallL7firewallrulesResource{}
}

// NetworksWirelessSsidsFirewallL7firewallrulesResource defines the resource implementation.
type NetworksWirelessSsidsFirewallL7firewallrulesResource struct {
	client *openApiClient.APIClient
}

// NetworksWirelessSsidsFirewallL7firewallrulesResourceModel describes the resource data model.
type NetworksWirelessSsidsFirewallL7firewallrulesResourceModel struct {
	Id        jsontypes.String `tfsdk:"id"`
	NetworkId jsontypes.String `tfsdk:"network_id"`
	Number    jsontypes.String `tfsdk:"number"`
	Rules     []Rule           `tfsdk:"rules"`
}

type Rule struct {
	Policy jsontypes.String `tfsdk:"policy"`
	Type   jsontypes.String `tfsdk:"type"`
	Value  jsontypes.String `tfsdk:"value"`
}

func (r *NetworksWirelessSsidsFirewallL7firewallrulesResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_wireless_ssids_firewall_l7_firewall_rules"
}

func (r *NetworksWirelessSsidsFirewallL7firewallrulesResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{

		MarkdownDescription: "NetworksWirelessSsidsFirewallL7firewallrules",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.StringType,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Network ID",
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
				MarkdownDescription: "Number",
				Optional:            true,
				Computed:            true,
				CustomType:          jsontypes.StringType,
			},
			"rules": schema.ListNestedAttribute{
				MarkdownDescription: "Rules",
				Optional:            false,
				Computed:            false,
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"policy": schema.StringAttribute{
							MarkdownDescription: "policy",
							Optional:            false,
							Computed:            false,
							Required:            true,
							CustomType:          jsontypes.StringType,
						},
						"type": schema.StringAttribute{
							MarkdownDescription: "type",
							Optional:            false,
							Computed:            false,
							Required:            true,
							CustomType:          jsontypes.StringType,
						},
						"value": schema.StringAttribute{
							MarkdownDescription: "value",
							Optional:            false,
							Computed:            false,
							Required:            true,
							CustomType:          jsontypes.StringType,
						},
					},
				},
			},
		},
	}
}

func (r *NetworksWirelessSsidsFirewallL7firewallrulesResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NetworksWirelessSsidsFirewallL7firewallrulesResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *NetworksWirelessSsidsFirewallL7firewallrulesResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	object158 := openApiClient.NewInlineObject158()
	v := []openApiClient.NetworksNetworkIdWirelessSsidsNumberFirewallL7FirewallRulesRules{}
	for _, inputRule := range data.Rules {
		rule := openApiClient.NewNetworksNetworkIdWirelessSsidsNumberFirewallL7FirewallRulesRules()
		rule.SetType(inputRule.Type.String())
		rule.SetPolicy(inputRule.Policy.String())
		rule.SetValue(inputRule.Value.String())

		v = append(v, *rule)
	}
	object158.SetRules(v)

	_, httpResp, err := r.client.WirelessApi.UpdateNetworkWirelessSsidFirewallL7FirewallRules(ctx, data.NetworkId.ValueString(), data.Number.ValueString()).UpdateNetworkWirelessSsidFirewallL7FirewallRules(*object158).Execute()

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

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "create resource")
}

func (r *NetworksWirelessSsidsFirewallL7firewallrulesResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetworksWirelessSsidsFirewallL7firewallrulesResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := r.client.WirelessApi.GetNetworkWirelessSsidFirewallL7FirewallRules(context.Background(), data.NetworkId.ValueString(), data.Number.ValueString()).Execute()
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
	if err = json.NewDecoder(httpResp.Body).Decode(data); err != nil {
		resp.Diagnostics.AddError(
			"JSON decoding error",
			fmt.Sprintf("%v\n", err.Error()),
		)
		return
	}
	data.Id = jsontypes.StringValue("example-id")
	// Save updated data into Terraform state

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "read resource")
}

func (r *NetworksWirelessSsidsFirewallL7firewallrulesResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *NetworksWirelessSsidsFirewallL7firewallrulesResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	object158 := openApiClient.NewInlineObject158()
	v := []openApiClient.NetworksNetworkIdWirelessSsidsNumberFirewallL7FirewallRulesRules{}
	for _, inputRule := range data.Rules {
		rule := openApiClient.NewNetworksNetworkIdWirelessSsidsNumberFirewallL7FirewallRulesRules()
		rule.SetType(inputRule.Type.String())
		rule.SetPolicy(inputRule.Policy.String())
		rule.SetValue(inputRule.Value.String())

		v = append(v, *rule)
	}
	object158.SetRules(v)

	_, httpResp, err := r.client.WirelessApi.UpdateNetworkWirelessSsidFirewallL7FirewallRules(ctx, data.NetworkId.ValueString(), data.Number.ValueString()).UpdateNetworkWirelessSsidFirewallL7FirewallRules(*object158).Execute()

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

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "updated resource")
}

func (r *NetworksWirelessSsidsFirewallL7firewallrulesResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *NetworksWirelessSsidsFirewallL7firewallrulesResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	object158 := openApiClient.NewInlineObject158()
	v := []openApiClient.NetworksNetworkIdWirelessSsidsNumberFirewallL7FirewallRulesRules{}
	for _, inputRule := range data.Rules {
		rule := openApiClient.NewNetworksNetworkIdWirelessSsidsNumberFirewallL7FirewallRulesRules()
		rule.SetType(inputRule.Type.String())
		rule.SetPolicy(inputRule.Policy.String())
		rule.SetValue(inputRule.Value.String())

		v = append(v, *rule)
	}
	object158.SetRules(v)

	_, httpResp, err := r.client.WirelessApi.UpdateNetworkWirelessSsidFirewallL7FirewallRules(ctx, data.NetworkId.ValueString(), data.Number.ValueString()).UpdateNetworkWirelessSsidFirewallL7FirewallRules(*object158).Execute()

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

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)

	// Write logs using the tflog package
	tflog.Trace(ctx, "removed resource")
}

func (r *NetworksWirelessSsidsFirewallL7firewallrulesResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: network_id, number. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("number"), idParts[1])...)
	if resp.Diagnostics.HasError() {
		return
	}
}
