package claim

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
	"io"
	"net/http"
	"strings"
	"time"
)

var (
	_ resource.Resource                = &Resource{}
	_ resource.ResourceWithImportState = &Resource{}
	_ resource.ResourceWithConfigure   = &Resource{}
)

func NewResource() resource.Resource {
	return &Resource{}
}

type Resource struct {
	client *openApiClient.APIClient
}

type resourceModel struct {
	Id        types.String `tfsdk:"id"`
	NetworkId types.String `tfsdk:"network_id"`
	Serials   types.Set    `tfsdk:"serials"`
}

func (r *Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_devices_claim"
}

func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Claim devices into a network",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"network_id": schema.StringAttribute{
				MarkdownDescription: "Network Id",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthBetween(1, 31),
				},
			},
			"serials": schema.SetAttribute{
				MarkdownDescription: "The serials of the devices that should be claimed",
				ElementType:         types.StringType,
				Computed:            true,
				Optional:            true,
			},
		},
	}
}

func (r *Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data resourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var serialsToClaim []string
	var serialsUnclaimed []string
	var serialsAlreadyClaimed []string
	var claimedSerials []string

	for _, serial := range data.Serials.Elements() {
		serialsToClaim = append(serialsToClaim, strings.Trim(serial.String(), "\""))
	}

	// Get current devices in the network to check if they are already claimed
	currentDevicesResp, _, err := r.client.DevicesApi.GetNetworkDevices(ctx, data.NetworkId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Error fetching current network devices", fmt.Sprintf("Could not fetch current network devices: %s", err))
		return
	}

	// Update list of verified claimed devices
	for _, device := range currentDevicesResp {
		serial, ok := device["serial"].(string)
		if ok {
			serialsAlreadyClaimed = append(serialsAlreadyClaimed, serial)
		}
	}

	// Convert claimed devices list into a map for quick lookup
	existsInSerialsAlreadyClaimed := make(map[string]bool)
	for _, item := range serialsAlreadyClaimed {
		existsInSerialsAlreadyClaimed[item] = true
	}

	// Determine unclaimed serials
	for _, serial := range serialsToClaim {
		if !existsInSerialsAlreadyClaimed[serial] {
			serialsUnclaimed = append(serialsUnclaimed, serial)
		}
	}

	if len(serialsUnclaimed) > 0 {
		claimNetworkDevices := *openApiClient.NewClaimNetworkDevicesRequest(serialsUnclaimed)

		maxRetries := r.client.GetConfig().MaximumRetries
		retryDelay := time.Duration(r.client.GetConfig().Retry4xxErrorWaitTime)

		apiCall := func() (interface{}, *http.Response, error) {
			time.Sleep(retryDelay * time.Millisecond)
			httpResp, err := r.client.NetworksApi.ClaimNetworkDevices(ctx, data.NetworkId.ValueString()).ClaimNetworkDevicesRequest(claimNetworkDevices).Execute()
			return httpResp.Body, httpResp, err
		}

		claimDevicesResp, httpResp, err := utils.CustomHttpRequestRetry(ctx, maxRetries, retryDelay, apiCall)
		if err != nil {
			handleError(ctx, err, httpResp, "Error claiming devices", resp)
			return
		}

		if httpResp.StatusCode != 200 {
			resp.Diagnostics.AddError("Unexpected HTTP Response Status Code", fmt.Sprintf("Received status code: %v", httpResp.StatusCode))
			return
		}

		// Decode the response from claiming devices
		var claimResponse map[string]interface{}
		if err := json.NewDecoder(claimDevicesResp.(io.Reader)).Decode(&claimResponse); err != nil {
			resp.Diagnostics.AddError("Error decoding response", fmt.Sprintf("Error decoding response: %v", err))
			return
		}

		if serials, ok := claimResponse["serials"].([]interface{}); ok {
			for _, serial := range serials {
				if s, ok := serial.(string); ok {
					claimedSerials = append(claimedSerials, s)
				}
			}
		}
	} else {
		tflog.Info(ctx, "All devices are already claimed in the network", map[string]interface{}{
			"network_id": data.NetworkId.ValueString(),
		})
	}

	// Convert the claimed serials list into types.StringList
	claimedSerialsList, claimedSerialsListErr := types.SetValueFrom(ctx, types.StringType, claimedSerials)
	if claimedSerialsListErr.HasError() {
		resp.Diagnostics.Append(claimedSerialsListErr...)
		return
	}
	data.Serials = claimedSerialsList

	data.Id = data.NetworkId

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "Created resource", map[string]interface{}{
		"network_id": data.NetworkId.ValueString(),
	})
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data resourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Call the API to get network devices
	devicesResp, httpResp, err := r.client.DevicesApi.GetNetworkDevices(ctx, data.NetworkId.ValueString()).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Failed to read list of serials", err.Error())
		return
	}

	// Verify HTTP status code
	if httpResp.StatusCode != 200 {
		resp.Diagnostics.AddError("Unexpected HTTP response status code", fmt.Sprintf("Received status code: %d", httpResp.StatusCode))
		return
	}

	// Process the response and collect serial numbers
	var respSerials []string
	for _, device := range devicesResp {
		if serial, ok := device["serial"].(string); ok {
			respSerials = append(respSerials, serial)
		}
	}

	// Convert serial numbers to a type that can be stored in the Terraform state
	serialsList, serialsListErr := types.SetValueFrom(ctx, types.StringType, respSerials)
	if serialsListErr != nil {
		resp.Diagnostics.AddError("Failed to convert serial numbers for Terraform state", fmt.Sprintf("%s", serialsListErr))
		return
	}
	data.Serials = serialsList

	// Set the read data ID and store the state
	data.Id = types.StringValue(data.NetworkId.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "Successfully read network device resource", map[string]interface{}{
		"network_id": data.NetworkId.ValueString(),
	})
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data, state resourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract serials from the plan and state
	planSerials := extractSerials(data.Serials)
	stateSerials := extractSerials(state.Serials)

	// Determine which serials to add or remove
	serialsToAdd := difference(planSerials, stateSerials)
	serialsToRemove := difference(stateSerials, planSerials)

	// Claim new devices
	if len(serialsToAdd) > 0 {
		if err := manageDeviceClaims(ctx, r.client, data.NetworkId.ValueString(), serialsToAdd, true, resp); err != nil {
			return
		}
	}

	// Remove devices
	if len(serialsToRemove) > 0 {
		if err := manageDeviceClaims(ctx, r.client, data.NetworkId.ValueString(), serialsToRemove, false, resp); err != nil {
			return
		}
	}

	// Ensure the state is updated correctly with the new serials
	claimedSerials := mergeSerials(planSerials, serialsToAdd)

	updatedSerialsList, serialsListErr := types.SetValueFrom(ctx, types.StringType, claimedSerials)
	if serialsListErr != nil {
		resp.Diagnostics.AddError("Failed to convert serial numbers for Terraform state", fmt.Sprintf("%s", serialsListErr))
		return
	}
	data.Serials = updatedSerialsList

	// Update the ID and state
	data.Id = types.StringValue(data.NetworkId.ValueString())
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "Updated network devices resource", map[string]interface{}{
		"network_id": data.NetworkId.ValueString(),
		"added":      serialsToAdd,
		"removed":    serialsToRemove,
	})
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data resourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Extract serials from the state
	serials := extractSerials(data.Serials)

	if len(serials) > 0 {
		// Attempt to remove all devices
		if err := removeDevices(ctx, r.client, data.NetworkId.ValueString(), serials, resp); err != nil {
			return
		}
	}

	// Confirm removal of the resource from state
	resp.State.RemoveResource(ctx)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Trace(ctx, "Successfully deleted network devices resource", map[string]interface{}{
		"network_id": data.NetworkId.ValueString(),
	})
}
