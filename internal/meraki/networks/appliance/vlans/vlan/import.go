package vlan

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
	"strings"
)

func (r *Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
	idParts := strings.Split(req.ID, ",")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			fmt.Sprintf("Expected import identifier with format: network_id, vlan_id. Got: %q", req.ID),
		)
		return
	}

	// ensure vlanId is formatted properly
	str := idParts[1]

	// Convert the string to int64
	i, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		resp.Diagnostics.AddError(
			"Failed to convert vlanId to integer",
			fmt.Sprintf("Expected import identifier with format: network_id, vlan_id. Got: %q", req.ID),
		)
	}

	// Convert the int64 to types.Int64Value
	vlanId := types.Int64Value(i)

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("network_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vlan_id"), vlanId)...)
	if resp.Diagnostics.HasError() {
		return
	}
}
