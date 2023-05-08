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
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &NetworksSwitchStormcontrolResource{}
	_ resource.ResourceWithConfigure   = &NetworksSwitchStormcontrolResource{}
	_ resource.ResourceWithImportState = &NetworksSwitchStormcontrolResource{}
)

// TODO - This function needs to be added to the list of Data Sources in provider.go: func (p *ScaffoldingProvider) Resources
// TODO - Otherwise the provider has no idea this resource exists.
func NewNetworksSwitchStormcontrolResource() resource.Resource {
	return &NetworksSwitchStormcontrolResource{}
}

// NetworksSwitchStormcontrolResource defines the resource implementation.
type NetworksSwitchStormcontrolResource struct {
	client *openApiClient.APIClient
}

// NetworksSwitchStormcontrolResourceModel describes the resource data model.
type NetworksSwitchStormcontrolResourceModel struct {
	Id        jsontypes.String `tfsdk:"id"`
	NetworkId jsontypes.String `tfsdk:"network_id"`

	BroadcastThreshold      jsontypes.Int64 `tfsdk:"broadcast_threshold" json:"broadcastThreshold"`
	MulticastThreshold      jsontypes.Int64 `tfsdk:"multicast_threshold" json:"multicastThreshold"`
	UnknownUnicastThreshold jsontypes.Int64 `tfsdk:"unknown_unicast_threshold" json:"unknownUnicastThreshold"`
}

func (r *NetworksSwitchStormcontrolResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_switch_storm_control"
}

func (r *NetworksSwitchStormcontrolResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "NetworksSwitchStormcontrol",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.StringType,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Networkd ID",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
				Validators: []validator.String{
					stringvalidator.LengthBetween(8, 31),
				},
				CustomType: jsontypes.StringType,
			},
			"broadcast_threshold": schema.Int64Attribute{
				MarkdownDescription: "Broadcast Threshold",
				Description:         "Percentage (1 to 99) of total available port bandwidth for broadcast traffic type. Default value 100 percent rate is to clear the configuration.",
				Optional:            true,
				CustomType:          jsontypes.Int64Type,
			},
			"multicast_threshold": schema.Int64Attribute{
				MarkdownDescription: "Multicast Threshold",
				Description:         "Percentage (1 to 99) of total available port bandwidth for multicast traffic type. Default value 100 percent rate is to clear the configuration.",
				Optional:            true,
				CustomType:          jsontypes.Int64Type,
			},
			"unknown_unicast_threshold": schema.Int64Attribute{
				MarkdownDescription: "Unknown Unicast Threshold",
				Description:         "Percentage (1 to 99) of total available port bandwidth for unknown unicast (dlf-destination lookup failure) traffic type. Default value 100 percent rate is to clear the configuration.",
				Optional:            true,
				CustomType:          jsontypes.Int64Type,
			},
		},
	}
}

func (r *NetworksSwitchStormcontrolResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *NetworksSwitchStormcontrolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data *NetworksSwitchStormcontrolResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	object137 := openApiClient.NewInlineObject137()
	object137.SetMulticastThreshold(int32(data.MulticastThreshold.ValueInt64()))
	object137.SetBroadcastThreshold(int32(data.BroadcastThreshold.ValueInt64()))
	object137.SetUnknownUnicastThreshold(int32(data.UnknownUnicastThreshold.ValueInt64()))

	_, httpResp, err := r.client.SwitchApi.UpdateNetworkSwitchStormControl(context.Background(), data.NetworkId.ValueString()).UpdateNetworkSwitchStormControl(*object137).Execute()
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

func (r *NetworksSwitchStormcontrolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data *NetworksSwitchStormcontrolResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	_, httpResp, err := r.client.SwitchApi.GetNetworkSwitchStormControl(context.Background(), data.NetworkId.ValueString()).Execute()
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

func (r *NetworksSwitchStormcontrolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data *NetworksSwitchStormcontrolResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	object137 := openApiClient.NewInlineObject137()
	object137.SetMulticastThreshold(int32(data.MulticastThreshold.ValueInt64()))
	object137.SetBroadcastThreshold(int32(data.BroadcastThreshold.ValueInt64()))
	object137.SetUnknownUnicastThreshold(int32(data.UnknownUnicastThreshold.ValueInt64()))

	_, httpResp, err := r.client.SwitchApi.UpdateNetworkSwitchStormControl(context.Background(), data.NetworkId.ValueString()).UpdateNetworkSwitchStormControl(*object137).Execute()
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

func (r *NetworksSwitchStormcontrolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data *NetworksSwitchStormcontrolResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	object137 := openApiClient.NewInlineObject137()
	object137.SetMulticastThreshold(int32(data.MulticastThreshold.ValueInt64()))
	object137.SetBroadcastThreshold(int32(data.BroadcastThreshold.ValueInt64()))
	object137.SetUnknownUnicastThreshold(int32(data.UnknownUnicastThreshold.ValueInt64()))

	_, httpResp, err := r.client.SwitchApi.UpdateNetworkSwitchStormControl(context.Background(), data.NetworkId.ValueString()).UpdateNetworkSwitchStormControl(*object137).Execute()
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
	tflog.Trace(ctx, "resource removed")
}

func (r *NetworksSwitchStormcontrolResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), req.ID)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
