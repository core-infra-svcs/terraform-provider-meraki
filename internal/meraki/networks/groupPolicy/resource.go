package groupPolicy

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/meraki/dashboard-api-go/client"
	"io"
	"net/http"
	"strings"
	"time"
)

// Resource defines the resource implementation.
type Resource struct {
	client *client.APIClient
}

func NewResource() resource.Resource {
	return &Resource{}
}

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &Resource{}
var _ resource.ResourceWithConfigure = &Resource{}
var _ resource.ResourceWithImportState = &Resource{}

// Schema defines the schema for the resource.
func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: ResourceSchema,
	}
}

func (r *Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_networks_group_policy"
}

// Configure configures the resource with the API client.
func (r *Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*client.APIClient)

}

// Create handles the creation of the group policy.
func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan resourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload, groupPolicyDiags := fromPayload(&plan)
	if groupPolicyDiags.HasError() {

		resp.Diagnostics.AddError(
			"Error creating group policy payload",
			fmt.Sprintf("Unexpected error: %s", groupPolicyDiags),
		)
		return
	}

	maxRetries := r.client.GetConfig().MaximumRetries
	retryDelay := time.Duration(r.client.GetConfig().Retry4xxErrorWaitTime)

	// API call function to be passed to retryOn4xx
	apiCall := func() (map[string]interface{}, *http.Response, error) {
		time.Sleep(retryDelay * time.Millisecond)

		inline, httpResp, err := r.client.NetworksApi.CreateNetworkGroupPolicy(ctx, plan.NetworkId.ValueString()).CreateNetworkGroupPolicyRequest(payload).Execute()
		return inline, httpResp, err
	}
	// Use retryOn4xx for the API call as the meraki API backend returns HTTP 400 messages as a result of collision issues with rapid creation of postgres GroupPolicyIds.
	inlineResp, httpResp, err := utils.CustomHttpRequestRetry(ctx, maxRetries, retryDelay, apiCall)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating group policy",
			fmt.Sprintf("Could not create group policy, unexpected error: %s", err),
		)

		if httpResp != nil {
			var responseBody string
			if httpResp != nil && httpResp.Body != nil {
				bodyBytes, readErr := io.ReadAll(httpResp.Body)
				if readErr == nil {
					responseBody = string(bodyBytes)
				}
			}
			tflog.Error(ctx, "Failed to create resource", map[string]interface{}{
				"error":          err.Error(),
				"httpStatusCode": httpResp.StatusCode,
				"responseBody":   responseBody,
			})
			resp.Diagnostics.AddError(
				"Error creating group policy",
				fmt.Sprintf("HTTP Response: %v\nResponse Body: %s", httpResp, responseBody),
			)
		}
		return
	}

	diags = ToTerraformState(ctx, &plan, inlineResp)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Read handles reading the group policy.
func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	inlineResp, httpResp, err := r.client.NetworksApi.GetNetworkGroupPolicy(ctx, state.NetworkId.ValueString(), state.GroupPolicyId.ValueString()).Execute()
	if err != nil {
		var responseBody string
		if httpResp != nil && httpResp.Body != nil {
			bodyBytes, readErr := io.ReadAll(httpResp.Body)
			if readErr == nil {
				responseBody = string(bodyBytes)
			}
		}

		resp.Diagnostics.AddError(
			"Error reading group policy",
			fmt.Sprintf("Could not read group policy, unexpected error: %s\nHTTP Response: %v\nResponse Body: %s", err, httpResp, responseBody),
		)
		return
	}

	// Update the state with the new state
	diags = ToTerraformState(ctx, &state, inlineResp)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update handles updating the group policy.
func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan resourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state resourceModel

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	groupPolicy, groupPolicyDiags := fromPayload(&plan)
	if groupPolicyDiags.HasError() {

		resp.Diagnostics.AddError(
			"Error creating group policy payload",
			fmt.Sprintf("Unexpected error: %s", groupPolicyDiags),
		)
		return
	}

	groupPolicyUpdate := client.UpdateNetworkGroupPolicyRequest{
		Name:                      &groupPolicy.Name,
		Scheduling:                groupPolicy.Scheduling,
		Bandwidth:                 groupPolicy.Bandwidth,
		FirewallAndTrafficShaping: groupPolicy.FirewallAndTrafficShaping,
		ContentFiltering:          groupPolicy.ContentFiltering,
		SplashAuthSettings:        groupPolicy.SplashAuthSettings,
		VlanTagging:               groupPolicy.VlanTagging,
		BonjourForwarding:         groupPolicy.BonjourForwarding,
	}

	inlineResp, httpResp, err := r.client.NetworksApi.UpdateNetworkGroupPolicy(ctx, plan.NetworkId.ValueString(), state.GroupPolicyId.ValueString()).UpdateNetworkGroupPolicyRequest(groupPolicyUpdate).Execute()
	if err != nil {
		var responseBody string
		if httpResp != nil && httpResp.Body != nil {
			bodyBytes, readErr := io.ReadAll(httpResp.Body)
			if readErr == nil {
				responseBody = string(bodyBytes)
			}
		}

		resp.Diagnostics.AddError(
			"Error updating group policy",
			fmt.Sprintf("Could not update group policy, unexpected error: %s\nHTTP Response: %v\nResponse Body: %s", err, httpResp, responseBody),
		)
		return
	}

	// Update the state with the new plan
	diags = ToTerraformState(ctx, &plan, inlineResp)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Delete handles deleting the group policy.
func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state resourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check for empty GroupPolicyId
	if state.GroupPolicyId.IsNull() || state.GroupPolicyId.IsUnknown() {
		resp.Diagnostics.AddError("DELETE, Received empty GroupPolicy.", fmt.Sprintf("Name: %s, Id: %s", state.Name.ValueString(), state.GroupPolicyId.ValueString()))
		return
	}

	httpResp, err := r.client.NetworksApi.DeleteNetworkGroupPolicy(ctx, state.NetworkId.ValueString(), state.GroupPolicyId.ValueString()).Execute()
	if err != nil {
		var responseBody string
		if httpResp != nil && httpResp.Body != nil {
			bodyBytes, readErr := io.ReadAll(httpResp.Body)
			if readErr == nil {
				responseBody = string(bodyBytes)
			}
		}
		resp.Diagnostics.AddError(
			"Error deleting group policy",
			fmt.Sprintf("Could not delete group policy, unexpected error: %s\nHTTP Response: %v\nResponse Body: %s", err, httpResp, responseBody),
		)
		return
	}

	resp.State.RemoveResource(ctx)

}

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: network_id, group_policy_id. Got: %q", req.ID),
		)
		return
	}
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("group_policy_id"), idParts[1])...)
	if resp.Diagnostics.HasError() {
		return
	}
}
