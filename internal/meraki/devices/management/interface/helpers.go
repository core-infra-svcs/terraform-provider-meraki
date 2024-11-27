package _interface

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// mapAPIResponseToState maps the API response to the Terraform data source state.
func mapAPIResponseToState(ctx context.Context, rawResponse map[string]interface{}) (resourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var state resourceModel

	tflog.Trace(ctx, "[management_interface] Mapping API response to state")

	// Map top-level fields
	id, idDiags := utils.ExtractStringAttr(rawResponse, "id")
	serial, serialDiags := utils.ExtractStringAttr(rawResponse, "serial")
	diags = append(diags, idDiags...)
	diags = append(diags, serialDiags...)

	state.Id = id
	state.Serial = serial

	// Map nested DDNSHostnames
	ddnsHostnames, ddnsDiags := mapDDNSHostnames(ctx, rawResponse["ddnsHostnames"])
	diags = append(diags, ddnsDiags...)
	state.DDNSHostnames = ddnsHostnames

	// Map WAN1 and WAN2
	wan1, wan1Diags := mapWAN(ctx, rawResponse, "wan1")
	wan2, wan2Diags := mapWAN(ctx, rawResponse, "wan2")
	diags = append(diags, wan1Diags...)
	diags = append(diags, wan2Diags...)

	state.Wan1 = wan1
	state.Wan2 = wan2

	tflog.Debug(ctx, "[management_interface] Mapped state", map[string]interface{}{
		"state": state,
	})

	return state, diags
}

// mapDDNSHostnames maps the DDNS hostnames from the API response.
func mapDDNSHostnames(ctx context.Context, raw interface{}) (types.Object, diag.Diagnostics) {
	tflog.Trace(ctx, "[management_interface] Mapping DDNS hostnames")
	if raw == nil {
		return types.ObjectNull(ddnsHostnamesType), nil
	}

	ddnsRaw, ok := raw.(map[string]interface{})
	if !ok {
		return types.ObjectNull(ddnsHostnamesType), diag.Diagnostics{
			diag.NewErrorDiagnostic("Mapping Error", "Expected a map for DDNS hostnames"),
		}
	}

	activeDDNS, activeDDNSDiags := utils.ExtractStringAttr(ddnsRaw, "activeDdnsHostname")
	wan1DDNS, wan1DDNSDiags := utils.ExtractStringAttr(ddnsRaw, "ddnsHostnameWan1")
	wan2DDNS, wan2DDNSDiags := utils.ExtractStringAttr(ddnsRaw, "ddnsHostnameWan2")

	diags := append(activeDDNSDiags, wan1DDNSDiags...)
	diags = append(diags, wan2DDNSDiags...)

	ddnsAttrValues := map[string]attr.Value{
		"active_ddns_hostname": activeDDNS,
		"ddns_hostname_wan1":   wan1DDNS,
		"ddns_hostname_wan2":   wan2DDNS,
	}

	ddnsObject, err := types.ObjectValue(ddnsHostnamesType, ddnsAttrValues)
	if err.HasError() {
		diags = append(diags, err...)
	}
	return ddnsObject, diags
}

// mapWAN maps a single WAN configuration (WAN1 or WAN2) from the API response.
func mapWAN(ctx context.Context, rawResponse map[string]interface{}, wanKey string) (types.Object, diag.Diagnostics) {
	tflog.Trace(ctx, fmt.Sprintf("[management_interface] Mapping %s", wanKey))
	if rawResponse[wanKey] == nil {
		return types.ObjectNull(wanType), nil
	}

	wanRaw, ok := rawResponse[wanKey].(map[string]interface{})
	if !ok {
		return types.ObjectNull(wanType), diag.Diagnostics{
			diag.NewErrorDiagnostic("Mapping Error", fmt.Sprintf("Expected a map for %s", wanKey)),
		}
	}

	wanEnabled, wanEnabledDiags := utils.ExtractStringAttr(wanRaw, "wanEnabled")
	usingStaticIP, usingStaticIPDiags := utils.ExtractBoolAttr(wanRaw, "usingStaticIp")
	staticIP, staticIPDiags := utils.ExtractStringAttr(wanRaw, "staticIp")
	staticSubnet, staticSubnetDiags := utils.ExtractStringAttr(wanRaw, "staticSubnetMask")
	gatewayIP, gatewayIPDiags := utils.ExtractStringAttr(wanRaw, "staticGatewayIp")
	dns, dnsDiags := utils.ExtractListStringAttr(wanRaw, "staticDns")
	vlan, vlanDiags := utils.ExtractInt64Attr(wanRaw, "vlan")

	diags := append(wanEnabledDiags, usingStaticIPDiags...)
	diags = append(diags, staticIPDiags...)
	diags = append(diags, staticSubnetDiags...)
	diags = append(diags, gatewayIPDiags...)
	diags = append(diags, dnsDiags...)
	diags = append(diags, vlanDiags...)

	wanAttrValues := map[string]attr.Value{
		"wan_enabled":        wanEnabled,
		"using_static_ip":    usingStaticIP,
		"static_ip":          staticIP,
		"static_subnet_mask": staticSubnet,
		"static_gateway_ip":  gatewayIP,
		"static_dns":         dns,
		"vlan":               vlan,
	}

	wanObject, err := types.ObjectValue(wanType, wanAttrValues)
	if err.HasError() {
		diags = append(diags, err...)
	}
	return wanObject, diags
}
