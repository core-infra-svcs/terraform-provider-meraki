package device

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// ImportState function is used to import an existing resource into Terraform.
func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {

	// Pass through the ID directly from the ImportStateRequest to the ImportStateResponse
	resource.ImportStatePassthroughID(ctx, path.Root("serial"), req, resp)

}
