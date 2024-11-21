package cellular

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Log the import operation
	tflog.Debug(ctx, "Starting import for cellular SIMs resource", map[string]interface{}{
		"import_id": req.ID,
	})

	// Call the API to fetch the resource using the import ID (serial)
	apiReq := r.client.CellularApi.GetDeviceCellularSims(ctx, req.ID)
	apiResp, httpResp, err := apiReq.Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to import cellular SIMs resource",
			fmt.Sprintf("Error: %s, HTTP Response: %v, Import ID: %s", err.Error(), httpResp, req.ID),
		)
		return
	}

	// Map the API response to the Terraform model
	var state resourceModel
	mapDiags := mapApiResponseToModel(apiResp, &state)
	resp.Diagnostics.Append(mapDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Save the imported state
	saveDiags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(saveDiags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Log success
	tflog.Debug(ctx, "Successfully imported cellular SIMs resource", map[string]interface{}{
		"import_id": req.ID,
	})
}
