package cellular

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Call the API to fetch the resource using the import ID
	apiReq := r.client.CellularApi.GetDeviceCellularSims(ctx, req.ID)
	apiResp, httpResp, err := apiReq.Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to import cellular SIMs resource",
			fmt.Sprintf("Error: %s, HTTP Response: %v", err.Error(), httpResp),
		)
		return
	}

	// Map the API response to the Terraform model
	var state resourceModel
	resp.Diagnostics.Append(mapApiResponseToModel(apiResp, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save the imported state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
