package cellular

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	openApiClient "github.com/meraki/dashboard-api-go/client"
)

var (
	_ resource.Resource                = &Resource{} // Terraform resource interface
	_ resource.ResourceWithConfigure   = &Resource{} // Interface for resources with configuration methods
	_ resource.ResourceWithImportState = &Resource{} // Interface for resources with import state functionality
)

func NewResource() resource.Resource {
	return &Resource{}
}

type Resource struct {
	client *openApiClient.APIClient // APIClient instance for making API requests
}

func (r *Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_devices_cellular_sims"
}

// Schema defines the schema for the resource.
func (r *Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = GetResourceSchema
}

func (r *Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	client, ok := req.ProviderData.(*openApiClient.APIClient)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T", req.ProviderData),
		)
		return
	}
	r.client = client
}

func (r *Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan resourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload, diags := mapModelToApiPayload(&plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := r.client.CellularApi.UpdateDeviceCellularSims(ctx, plan.Serial.ValueString())
	apiResp, httpResp, err := apiReq.UpdateDeviceCellularSimsRequest(*payload).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Failed to create cellular SIMs resource", fmt.Sprintf("Error: %s, HTTP Response: %v", err.Error(), httpResp))
		return
	}

	resp.Diagnostics.Append(mapApiResponseToModel(apiResp, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state resourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := r.client.CellularApi.GetDeviceCellularSims(ctx, state.Serial.ValueString())
	apiResp, httpResp, err := apiReq.Execute()
	if err != nil {
		resp.Diagnostics.AddError("Failed to read cellular SIMs resource", fmt.Sprintf("Error: %s, HTTP Response: %v", err.Error(), httpResp))
		return
	}

	resp.Diagnostics.Append(mapApiResponseToModel(apiResp, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state resourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	payload, diags := mapModelToApiPayload(&plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	apiReq := r.client.CellularApi.UpdateDeviceCellularSims(ctx, state.Serial.ValueString())
	apiResp, httpResp, err := apiReq.UpdateDeviceCellularSimsRequest(*payload).Execute()
	if err != nil {
		resp.Diagnostics.AddError("Failed to update cellular SIMs resource", fmt.Sprintf("Error: %s, HTTP Response: %v", err.Error(), httpResp))
		return
	}

	resp.Diagnostics.Append(mapApiResponseToModel(apiResp, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state resourceModel

	// Read the current state into the resource model
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create a payload with default/blank values
	resetPayload := openApiClient.NewUpdateDeviceCellularSimsRequest()

	// Set default values for Sims
	resetPayload.SetSims([]openApiClient.UpdateDeviceCellularSimsRequestSimsInner{})

	// Set default values for SimFailOver
	resetPayload.SetSimFailover(openApiClient.UpdateDeviceCellularSimsRequestSimFailover{
		Enabled: openApiClient.PtrBool(false),
		Timeout: openApiClient.PtrInt32(0),
	})

	// Call the Update API to reset the resource
	apiReq := r.client.CellularApi.UpdateDeviceCellularSims(ctx, state.Serial.ValueString())
	_, httpResp, err := apiReq.UpdateDeviceCellularSimsRequest(*resetPayload).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to delete cellular SIMs resource",
			fmt.Sprintf("Error: %s, HTTP Response: %v", err.Error(), httpResp),
		)
		return
	}

	// Remove the resource from state
	resp.State.RemoveResource(ctx)
}
