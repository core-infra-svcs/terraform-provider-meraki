package _interface

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

// Ensure provider defined types fully satisfy framework interfaces
var (
	_ resource.Resource                = &Resource{}
	_ resource.ResourceWithConfigure   = &Resource{}
	_ resource.ResourceWithImportState = &Resource{}
)

func NewResource() resource.Resource {
	return &Resource{}
}

// Resource defines the resource implementation.
type Resource struct {
	client *openApiClient.APIClient
}

func (r *Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_devices_management_interface"
}

func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = GetResourceSchema
}

func (r *Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

// Create implements the CREATE operation.
func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API payload
	payload := GenerateCreatePayload(ctx, data)

	// Call the CREATE API (PUT)
	apiResponse, httpResp, err := CallCreateAPI(ctx, r.client, payload, data.Serial.ValueString())
	if err := utils.HandleAPIError(ctx, httpResp, err, &resp.Diagnostics); err != nil {
		return
	}

	// Marshal API response into state
	state, diags := MarshalStateFromAPI(ctx, apiResponse)
	resp.Diagnostics.Append(diags...)

	resp.State.Set(ctx, &state)
}

// Read implements the READ operation.
func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Call the READ API (GET)
	apiResponse, httpResp, err := CallReadAPI(ctx, r.client, data.Serial.ValueString())
	if err := utils.HandleAPIError(ctx, httpResp, err, &resp.Diagnostics); err != nil {
		resp.State.RemoveResource(ctx)
		return
	}

	// Marshal API response into state
	state, diags := MarshalStateFromAPI(ctx, apiResponse)
	resp.Diagnostics.Append(diags...)

	resp.State.Set(ctx, &state)
}

// Update implements the UPDATE operation.
func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data ResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Generate API payload
	payload := GenerateUpdatePayload(ctx, data)

	// Call the UPDATE API (PUT)
	apiResponse, httpResp, err := CallUpdateAPI(ctx, r.client, payload, data.Serial.ValueString())
	if err := utils.HandleAPIError(ctx, httpResp, err, &resp.Diagnostics); err != nil {
		return
	}

	// Marshal API response into state
	state, diags := MarshalStateFromAPI(ctx, apiResponse)
	resp.Diagnostics.Append(diags...)

	resp.State.Set(ctx, &state)
}

// Delete implements the DELETE operation.
func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data ResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Generate blank/default payload for DELETE
	payload := GenerateDeletePayload(ctx, data)

	// Call the DELETE API (PUT with blank payload)
	httpResp, err := CallDeleteAPI(ctx, r.client, payload, data.Serial.ValueString())
	if err := utils.HandleAPIError(ctx, httpResp, err, &resp.Diagnostics); err != nil {
		return
	}

	// Remove resource from state
	resp.State.RemoveResource(ctx)
}

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("serial"), req.ID)...)

	if resp.Diagnostics.HasError() {
		return
	}
}
