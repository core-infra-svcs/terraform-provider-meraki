package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"net/http"
	"strings"
	"time"

	"github.com/core-infra-svcs/terraform-provider-meraki/internal/provider/jsontypes"
	"github.com/core-infra-svcs/terraform-provider-meraki/tools"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

var (
	_ resource.Resource                = &NetworksDevicesClaimResource{}
	_ resource.ResourceWithImportState = &NetworksDevicesClaimResource{}
	_ resource.ResourceWithConfigure   = &NetworksDevicesClaimResource{}
)

func NewNetworksDevicesClaimResource() resource.Resource {
	return &NetworksDevicesClaimResource{}
}

type NetworksDevicesClaimResource struct {
	client *openApiClient.APIClient
}

type NetworksDevicesClaimResourceModel struct {
	Id        jsontypes.String   `tfsdk:"id"`
	NetworkId jsontypes.String   `tfsdk:"network_id"`
	Serials   []jsontypes.String `tfsdk:"serials"`
}

func (r *NetworksDevicesClaimResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_devices_claim"
}

func (r *NetworksDevicesClaimResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Claim devices into a network",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:   true,
				CustomType: jsontypes.StringType,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Network Id",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 31),
				},
				CustomType: jsontypes.StringType,
			},
			"serials": schema.SetAttribute{
				MarkdownDescription: "The serials of the devices that should be claimed",
				ElementType:         jsontypes.StringType,
				CustomType:          jsontypes.SetType[jsontypes.String](),
				Computed:            true,
				Optional:            true,
			},
		},
	}
}

func (r *NetworksDevicesClaimResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*openApiClient.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *openApiClient.APIClient, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *NetworksDevicesClaimResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data NetworksDevicesClaimResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var serials []string
	for _, serial := range data.Serials {
		serials = append(serials, serial.ValueString())
	}

	claimNetworkDevices := *openApiClient.NewClaimNetworkDevicesRequest(serials)

	maxRetries := r.client.GetConfig().MaximumRetries
	retryDelay := r.client.GetConfig().Retry4xxErrorWaitTime

	httpResp, err := r.client.NetworksApi.ClaimNetworkDevices(ctx, data.NetworkId.ValueString()).ClaimNetworkDevicesRequest(claimNetworkDevices).Execute()
	retries := 0
	remaining := maxRetries - retries
	for retries < maxRetries && httpResp != nil && httpResp.StatusCode == http.StatusBadRequest {
		tflog.Warn(ctx, "Retrying Create API call", map[string]interface{}{
			"maxRetries":        maxRetries,
			"retryDelay":        retryDelay,
			"remainingAttempts": remaining,
			"httpStatusCode":    httpResp.StatusCode,
		})
		time.Sleep(time.Duration(retryDelay) * time.Second)
		httpResp, err = r.client.NetworksApi.ClaimNetworkDevices(ctx, data.NetworkId.ValueString()).ClaimNetworkDevicesRequest(claimNetworkDevices).Execute()
		retries++
	}

	if err != nil {
		resp.Diagnostics.AddError("Error claiming devices", err.Error())
		if httpResp != nil {
			resp.Diagnostics.AddError("HTTP Response", tools.HttpDiagnostics(httpResp))
		}
		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError("Unexpected HTTP Response Status Code", fmt.Sprintf("%v", httpResp.StatusCode))
		return
	}

	data.Id = jsontypes.StringValue(data.NetworkId.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "Created resource")
}

func (r *NetworksDevicesClaimResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data NetworksDevicesClaimResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	inlineResp, httpResp, err := r.client.DevicesApi.GetNetworkDevices(ctx, data.NetworkId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error reading devices", err.Error())
		if httpResp != nil {
			resp.Diagnostics.AddError("HTTP Response", tools.HttpDiagnostics(httpResp))
		}
		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError("Unexpected HTTP Response Status Code", fmt.Sprintf("%v", httpResp.StatusCode))
		return
	}

	var respSerials []jsontypes.String
	for _, device := range inlineResp {
		if serial, sOk := device["serial"].(string); sOk {
			respSerials = append(respSerials, jsontypes.StringValue(serial))
		}
	}

	data.Serials = respSerials
	data.Id = jsontypes.StringValue(data.NetworkId.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "Read resource")
}

func (r *NetworksDevicesClaimResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state NetworksDevicesClaimResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var planSerials, stateSerials, serialsToAdd, serialsToRemove []string
	for _, serial := range data.Serials {
		planSerials = append(planSerials, serial.ValueString())
	}
	for _, serial := range state.Serials {
		stateSerials = append(stateSerials, serial.ValueString())
	}

	serialsToAdd = difference(planSerials, stateSerials)
	serialsToRemove = difference(stateSerials, planSerials)

	claimNetworkDevices := *openApiClient.NewClaimNetworkDevicesRequest(serialsToAdd)
	httpResp, err := r.client.NetworksApi.ClaimNetworkDevices(ctx, data.NetworkId.ValueString()).ClaimNetworkDevicesRequest(claimNetworkDevices).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error claiming devices", err.Error())
		if httpResp != nil {
			resp.Diagnostics.AddError("HTTP Response", tools.HttpDiagnostics(httpResp))
		}
		return
	}

	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError("Unexpected HTTP Response Status Code", fmt.Sprintf("%v", httpResp.StatusCode))
		return
	}

	for _, serial := range serialsToRemove {
		removeNetworkDevices := *openApiClient.NewRemoveNetworkDevicesRequest(strings.Trim(serial, "\""))
		httpResp, err := r.client.NetworksApi.RemoveNetworkDevices(ctx, data.NetworkId.ValueString()).RemoveNetworkDevicesRequest(removeNetworkDevices).Execute()
		if err != nil {
			resp.Diagnostics.AddError("Error removing devices", err.Error())
			if httpResp != nil {
				resp.Diagnostics.AddError("HTTP Response", tools.HttpDiagnostics(httpResp))
			}
			return
		}

		if httpResp.StatusCode != 204 {
			resp.Diagnostics.AddError("Unexpected HTTP Response Status Code", fmt.Sprintf("%v", httpResp.StatusCode))
			return
		}
	}

	data.Id = jsontypes.StringValue(data.NetworkId.ValueString())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "Updated resource")
}

func (r *NetworksDevicesClaimResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data NetworksDevicesClaimResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	for _, serial := range data.Serials {
		removeNetworkDevices := *openApiClient.NewRemoveNetworkDevicesRequest(strings.Trim(serial.ValueString(), "\""))
		httpResp, err := r.client.NetworksApi.RemoveNetworkDevices(ctx, data.NetworkId.ValueString()).RemoveNetworkDevicesRequest(removeNetworkDevices).Execute()
		if err != nil {
			resp.Diagnostics.AddError("Error removing devices", err.Error())
			if httpResp != nil {
				resp.Diagnostics.AddError("HTTP Response", tools.HttpDiagnostics(httpResp))
			}
			return
		}

		if httpResp.StatusCode != 204 {
			resp.Diagnostics.AddError("Unexpected HTTP Response Status Code", fmt.Sprintf("%v", httpResp.StatusCode))
			return
		}
	}

	resp.State.RemoveResource(ctx)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "Deleted resource")
}

func (r *NetworksDevicesClaimResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("network_id"), req, resp)
	if resp.Diagnostics.HasError() {
		return
	}
}

func difference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}
