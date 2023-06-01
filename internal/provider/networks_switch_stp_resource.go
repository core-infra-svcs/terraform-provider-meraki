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
	_ resource.Resource                = &NetworksSwitchStpResource{}
	_ resource.ResourceWithConfigure   = &NetworksSwitchStpResource{}
	_ resource.ResourceWithImportState = &NetworksSwitchStpResource{}
)

func NewNetworksSwitchStpResource() resource.Resource {
	return &NetworksSwitchStpResource{}
}

// NetworksSwitchStpResource defines the resource implementation.
type NetworksSwitchStpResource struct {
	client *openApiClient.APIClient
}

// NetworksSwitchStpResourceModel describes the resource data model.
type NetworksSwitchStpResourceModel struct {
	Id                jsontypes.String    `tfsdk:"id"`
	NetworkId         jsontypes.String    `tfsdk:"network_id"`
	RstpEnabled       jsontypes.Bool      `tfsdk:"rstp_enabled" json:"rstpEnabled"`
	StpBridgePriority []STPBridgePriority `tfsdk:"stp_bridge_priority" json:"stpBridgePriority"`
}

type STPBridgePriority struct {
	Switches       jsontypes.Set[jsontypes.String] `tfsdk:"switches" json:"switches"`
	StpPriority    jsontypes.Int64                 `tfsdk:"stp_priority" json:"stpPriority"`
	Stacks         jsontypes.Set[jsontypes.String] `tfsdk:"stacks" json:"stacks"`
	SwitchProfiles jsontypes.Set[jsontypes.String] `tfsdk:"switch_profiles" json:"switchProfiles"`
}

func (r *NetworksSwitchStpResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_switch_stp"
}

func (r *NetworksSwitchStpResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "NetworksSwitchStp",
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
			},
			"rstp_enabled": schema.BoolAttribute{
				MarkdownDescription: "The spanning tree protocol status in network",
				Optional:            false,
				Computed:            false,
				Required:            true,
				CustomType:          jsontypes.BoolType,
			},
			"stp_bridge_priority": schema.ListNestedAttribute{
				MarkdownDescription: "STP bridge priority for switches/stacks or switch profiles. An empty array will clear the STP bridge priority settings.",
				Optional:            false,
				Computed:            false,
				Required:            true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"switches": schema.SetAttribute{
							CustomType:          jsontypes.SetType[jsontypes.String](),
							ElementType:         jsontypes.StringType,
							MarkdownDescription: "List of switch serial numbers",
							Optional:            true,
							Computed:            true,
						},
						"stp_priority": schema.Int64Attribute{
							MarkdownDescription: "STP priority for switch, stacks, or switch profiles",
							Optional:            false,
							Computed:            false,
							Required:            true,
							CustomType:          jsontypes.Int64Type,
						},
						"stacks": schema.SetAttribute{
							CustomType:          jsontypes.SetType[jsontypes.String](),
							ElementType:         jsontypes.StringType,
							MarkdownDescription: "List of stack IDs",
							Optional:            true,
							Computed:            true,
						},
						"switch_profiles": schema.SetAttribute{
							CustomType:          jsontypes.SetType[jsontypes.String](),
							ElementType:         jsontypes.StringType,
							MarkdownDescription: "List of switch profile IDs",
							Optional:            true,
							Computed:            true,
						},
					},
				},
			},
		},
	}
}

func (r *NetworksSwitchStpResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NetworksSwitchStpResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *NetworksSwitchStpResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	object138 := openApiClient.NewInlineObject138()
	object138.SetRstpEnabled(data.RstpEnabled.ValueBool())
	var stpBridgePriority []openApiClient.NetworksNetworkIdSwitchStpStpBridgePriority
	for _, d := range data.StpBridgePriority {
		priority := openApiClient.NewNetworksNetworkIdSwitchStpStpBridgePriority(int32(d.StpPriority.ValueInt64()))
		stacks := []string{}
		for _, stack := range d.Stacks.Elements() {
			stacks = append(stacks, stack.String())
		}
		switches := []string{}
		for _, switchs := range d.Switches.Elements() {
			switches = append(switches, switchs.String())
		}
		priority.SetSwitches(switches)
		priority.SetStacks(stacks)
		switcheProfiles := []string{}
		for _, switchProfile := range d.SwitchProfiles.Elements() {
			switcheProfiles = append(switcheProfiles, switchProfile.String())
		}
		priority.SetSwitchProfiles(switcheProfiles)
		stpBridgePriority = append(stpBridgePriority, *priority)
	}
	object138.SetStpBridgePriority(stpBridgePriority)
	_, httpResp, err := r.client.SwitchApi.UpdateNetworkSwitchStp(ctx, data.NetworkId.ValueString()).UpdateNetworkSwitchStp(*object138).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}
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

	// save into the Terraform state.
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

func (r *NetworksSwitchStpResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetworksSwitchStpResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := r.client.SwitchApi.GetNetworkSwitchStp(context.Background(), data.NetworkId.ValueString()).Execute()
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

	// save into the Terraform state.
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

func (r *NetworksSwitchStpResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *NetworksSwitchStpResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	object138 := openApiClient.NewInlineObject138()
	object138.SetRstpEnabled(data.RstpEnabled.ValueBool())
	var stpBridgePriority []openApiClient.NetworksNetworkIdSwitchStpStpBridgePriority
	for _, d := range data.StpBridgePriority {
		priority := openApiClient.NewNetworksNetworkIdSwitchStpStpBridgePriority(int32(d.StpPriority.ValueInt64()))
		stacks := []string{}
		for _, stack := range d.Stacks.Elements() {
			stacks = append(stacks, stack.String())
		}
		switches := []string{}
		for _, switchs := range d.Switches.Elements() {
			switches = append(switches, switchs.String())
		}
		priority.SetSwitches(switches)
		priority.SetStacks(stacks)
		switchProfiles := []string{}
		for _, switchProfile := range d.SwitchProfiles.Elements() {
			switchProfiles = append(switchProfiles, switchProfile.String())
		}
		priority.SetSwitchProfiles(switchProfiles)
		stpBridgePriority = append(stpBridgePriority, *priority)
	}
	object138.SetStpBridgePriority(stpBridgePriority)
	_, httpResp, err := r.client.SwitchApi.UpdateNetworkSwitchStp(ctx, data.NetworkId.String()).UpdateNetworkSwitchStp(*object138).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to Update resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}
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

	// save into the Terraform state.
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

func (r *NetworksSwitchStpResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *NetworksSwitchStpResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	object138 := openApiClient.NewInlineObject138()
	object138.SetRstpEnabled(data.RstpEnabled.ValueBool())
	var stpBridgePriority []openApiClient.NetworksNetworkIdSwitchStpStpBridgePriority
	for _, d := range data.StpBridgePriority {
		priority := openApiClient.NewNetworksNetworkIdSwitchStpStpBridgePriority(int32(d.StpPriority.ValueInt64()))
		stacks := []string{}
		for _, stack := range d.Stacks.Elements() {
			stacks = append(stacks, stack.String())
		}
		switches := []string{}
		for _, switchs := range d.Switches.Elements() {
			switches = append(switches, switchs.String())
		}
		priority.SetSwitches(switches)
		priority.SetStacks(stacks)
		switcheProfiles := []string{}
		for _, switchProfile := range d.SwitchProfiles.Elements() {
			switcheProfiles = append(switcheProfiles, switchProfile.String())
		}
		priority.SetSwitchProfiles(switcheProfiles)
		stpBridgePriority = append(stpBridgePriority, *priority)
	}
	object138.SetStpBridgePriority(stpBridgePriority)
	_, httpResp, err := r.client.SwitchApi.UpdateNetworkSwitchStp(ctx, data.NetworkId.String()).UpdateNetworkSwitchStp(*object138).Execute()

	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to create resource",
			fmt.Sprintf("%v\n", err.Error()),
		)
	}
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

	// save into the Terraform state.
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

func (r *NetworksSwitchStpResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), req.ID)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
