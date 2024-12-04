package _interface

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"net/http"
)

//// mapAPIResponseToState maps the API response to the Terraform data source state.
//func mapAPIResponseToState(ctx context.Context, rawResponse map[string]interface{}) (resourceModel, diag.Diagnostics) {
//	var diags diag.Diagnostics
//	var state resourceModel
//
//	tflog.Trace(ctx, "[management_interface] Mapping API response to state")
//
//	// Map top-level fields
//	id, idDiags := utils.ExtractStringAttr(rawResponse, "id")
//	serial, serialDiags := utils.ExtractStringAttr(rawResponse, "serial")
//	diags = append(diags, idDiags...)
//	diags = append(diags, serialDiags...)
//
//	state.Id = id
//	state.Serial = serial
//
//	// Map nested DDNSHostnames
//	ddnsHostnames, ddnsDiags := mapDDNSHostnames(ctx, rawResponse["ddnsHostnames"])
//	diags = append(diags, ddnsDiags...)
//	state.DDNSHostnames = ddnsHostnames
//
//	// Map WAN1 and WAN2
//	wan1, wan1Diags := mapWAN(ctx, rawResponse, "wan1")
//	wan2, wan2Diags := mapWAN(ctx, rawResponse, "wan2")
//	diags = append(diags, wan1Diags...)
//	diags = append(diags, wan2Diags...)
//
//	state.Wan1 = wan1
//	state.Wan2 = wan2
//
//	tflog.Debug(ctx, "[management_interface] Mapped state", map[string]interface{}{
//		"state": state,
//	})
//
//	return state, diags
//}

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

func resourceWanState(rawResp map[string]interface{}, wanKey string, wanEnabledPlan string) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var wan wanModel

	wanAttrs := map[string]attr.Type{
		"wan_enabled":        types.StringType,
		"using_static_ip":    types.BoolType,
		"static_ip":          types.StringType,
		"static_subnet_mask": types.StringType,
		"static_gateway_ip":  types.StringType,
		"static_dns":         types.ListType{ElemType: types.StringType},
		"vlan":               types.Int64Type,
	}

	if d, ok := rawResp[wanKey].(map[string]interface{}); ok {
		// wan_enabled
		wanEnabled, err := utils.ExtractStringAttr(d, "wanEnabled")

		if wanEnabledPlan == "not configured" && wanEnabled.IsNull() {
			wanEnabled = types.StringValue("not configured")
		}

		if wanEnabledPlan == "" && wanEnabled.IsNull() {
			wanEnabled = types.StringValue("")
		}

		if err != nil {
			diags.Append(err...)
		}
		wan.WanEnabled = wanEnabled

		// using_static_ip
		usingStaticIp, err := utils.ExtractBoolAttr(d, "usingStaticIp")
		if err != nil {
			diags.Append(err...)
		}
		wan.UsingStaticIp = usingStaticIp

		// static_ip
		staticIp, err := utils.ExtractStringAttr(d, "staticIp")
		if err != nil {
			diags.Append(err...)
		}
		wan.StaticIp = staticIp

		// static_subnet_mask
		staticSubnetMask, err := utils.ExtractStringAttr(d, "staticSubnetMask")
		if err != nil {
			diags.Append(err...)
		}
		wan.StaticSubnetMask = staticSubnetMask

		// static_gateway_ip
		staticGatewayIp, err := utils.ExtractStringAttr(d, "staticGatewayIp")
		if err != nil {
			diags.Append(err...)
		}
		wan.StaticGatewayIp = staticGatewayIp

		// static_dns
		staticDns, err := utils.ExtractListStringAttr(d, "staticDns")
		if err != nil {
			diags.Append(err...)
		}
		wan.StaticDns = staticDns

		// vlan
		if vlanValue, exists := d["vlan"]; exists && vlanValue != nil {
			switch v := vlanValue.(type) {
			case float64:
				wan.Vlan = types.Int64Value(int64(v))
			case int64:
				wan.Vlan = types.Int64Value(v)
			case int:
				wan.Vlan = types.Int64Value(int64(v))
			default:
				wan.Vlan = types.Int64Null()
				diags.AddError("Type Error", fmt.Sprintf("Unsupported type for vlan attribute: %T", v))
			}
		} else {
			wan.Vlan = types.Int64Null()
		}

		// Log the extracted vlan value
		tflog.Debug(context.Background(), "Extracted vlan", map[string]interface{}{
			"vlan": wan.Vlan.ValueInt64(),
		})

	} else {
		WanNull := types.ObjectNull(wanAttrs)
		return WanNull, diags
	}

	wanObj, err := types.ObjectValueFrom(context.Background(), wanAttrs, wan)
	if err != nil {
		diags.Append(err...)
	}

	return wanObj, diags
}

func updateResourceState(ctx context.Context, state *resourceModel, data map[string]interface{}, httpResp *http.Response, wan1EnabledPlan string, wan2EnabledPlan string) diag.Diagnostics {
	var diags diag.Diagnostics

	rawResp, err := utils.ExtractResponseToMap(httpResp)
	if err != nil {
		diags.AddError("Failed to Unmarshal HttpResp", err.Error())
		return diags
	}

	// ID
	if state.Id.IsNull() || state.Id.IsUnknown() {
		state.Id, diags = utils.ExtractStringAttr(rawResp, "id")
		if diags.HasError() {
			diags.AddError("ID Attribute", "")
			return diags
		}
	}

	// Serial
	if state.Serial.IsNull() || state.Serial.IsUnknown() {
		state.Serial, diags = utils.ExtractStringAttr(rawResp, "serial")
		if diags.HasError() {
			diags.AddError("Serial Attribute", "")
			return diags
		}
	}

	// Map nested DDNSHostnames
	ddnsHostnames, ddnsDiags := mapDDNSHostnames(ctx, data["ddnsHostnames"])
	diags = append(diags, ddnsDiags...)
	state.DDNSHostnames = ddnsHostnames

	// Wan1
	state.Wan1, diags = resourceWanState(rawResp, "wan1", wan1EnabledPlan)
	if diags.HasError() {
		diags.AddError("Wan1 Attribute", "")
		return diags
	}

	// Wan2
	state.Wan2, diags = resourceWanState(rawResp, "wan2", wan2EnabledPlan)
	if diags.HasError() {
		diags.AddError("Wan2 Attribute", "")
		return diags
	}

	// Log the updated state
	tflog.Debug(ctx, "Updated state", map[string]interface{}{
		"state": state,
	})

	return diags
}

func datasourceWanState(rawResp map[string]interface{}, wanKey string) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	var wan wanModel

	wanAttrs := map[string]attr.Type{
		"wan_enabled":        types.StringType,
		"using_static_ip":    types.BoolType,
		"static_ip":          types.StringType,
		"static_subnet_mask": types.StringType,
		"static_gateway_ip":  types.StringType,
		"static_dns":         types.ListType{ElemType: types.StringType},
		"vlan":               types.Int64Type,
	}

	if d, ok := rawResp[wanKey].(map[string]interface{}); ok {
		// wan_enabled
		wanEnabled, err := utils.ExtractStringAttr(d, "wanEnabled")
		if err != nil {
			diags.Append(err...)
		}
		wan.WanEnabled = wanEnabled

		// using_static_ip
		usingStaticIp, err := utils.ExtractBoolAttr(d, "usingStaticIp")
		if err != nil {
			diags.Append(err...)
		}
		wan.UsingStaticIp = usingStaticIp

		// static_ip
		staticIp, err := utils.ExtractStringAttr(d, "staticIp")
		if err != nil {
			diags.Append(err...)
		}
		wan.StaticIp = staticIp

		// static_subnet_mask
		staticSubnetMask, err := utils.ExtractStringAttr(d, "staticSubnetMask")
		if err != nil {
			diags.Append(err...)
		}
		wan.StaticSubnetMask = staticSubnetMask

		// static_gateway_ip
		staticGatewayIp, err := utils.ExtractStringAttr(d, "staticGatewayIp")
		if err != nil {
			diags.Append(err...)
		}
		wan.StaticGatewayIp = staticGatewayIp

		// static_dns
		staticDns, err := utils.ExtractListStringAttr(d, "staticDns")
		if err != nil {
			diags.Append(err...)
		}
		wan.StaticDns = staticDns

		// vlan
		if vlanValue, exists := d["vlan"]; exists && vlanValue != nil {
			switch v := vlanValue.(type) {
			case float64:
				wan.Vlan = types.Int64Value(int64(v))
			case int64:
				wan.Vlan = types.Int64Value(v)
			case int:
				wan.Vlan = types.Int64Value(int64(v))
			default:
				wan.Vlan = types.Int64Null()
				diags.AddError("Type Error", fmt.Sprintf("Unsupported type for vlan attribute: %T", v))
			}
		} else {
			wan.Vlan = types.Int64Null()
		}

		// Log the extracted vlan value
		tflog.Debug(context.Background(), "Extracted vlan", map[string]interface{}{
			"vlan": wan.Vlan.ValueInt64(),
		})

	} else {
		WanNull := types.ObjectNull(wanAttrs)
		return WanNull, diags
	}

	wanObj, err := types.ObjectValueFrom(context.Background(), wanAttrs, wan)
	if err != nil {
		diags.Append(err...)
	}

	return wanObj, diags
}

func updateDatasourceState(ctx context.Context, state *resourceModel, data map[string]interface{}, httpResp *http.Response) diag.Diagnostics {
	var diags diag.Diagnostics

	rawResp, err := utils.ExtractResponseToMap(httpResp)
	if err != nil {
		diags.AddError("Failed to Unmarshal HttpResp", err.Error())
		return diags
	}

	// ID
	if state.Id.IsNull() || state.Id.IsUnknown() {
		state.Id, diags = utils.ExtractStringAttr(rawResp, "id")
		if diags.HasError() {
			diags.AddError("ID Attribute", "")
			return diags
		}
	}

	// Serial
	if state.Serial.IsNull() || state.Serial.IsUnknown() {
		state.Serial, diags = utils.ExtractStringAttr(rawResp, "serial")
		if diags.HasError() {
			diags.AddError("Serial Attribute", "")
			return diags
		}
	}

	// Wan1
	state.Wan1, diags = datasourceWanState(rawResp, "wan1")
	if diags.HasError() {
		diags.AddError("Wan1 Attribute", "")
		return diags
	}

	// Wan2
	state.Wan2, diags = datasourceWanState(rawResp, "wan2")
	if diags.HasError() {
		diags.AddError("Wan2 Attribute", "")
		return diags
	}

	// Log the state after updating
	tflog.Debug(ctx, "State after update", map[string]interface{}{
		"state": state,
	})

	return diags
}
