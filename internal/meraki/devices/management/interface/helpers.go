package _interface

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"net/http"

	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	client "github.com/meraki/dashboard-api-go/client"
)

/* Payload Helper Methods */

// GenerateCreatePayload generates the payload for the CREATE operation.
func GenerateCreatePayload(ctx context.Context, data ResourceModel) *client.UpdateDeviceManagementInterfaceRequest {
	tflog.Debug(ctx, "Generating Create payload")
	return generateBasePayload(ctx, data)
}

// GenerateUpdatePayload generates the payload for the UPDATE operation.
func GenerateUpdatePayload(ctx context.Context, data ResourceModel) *client.UpdateDeviceManagementInterfaceRequest {
	tflog.Debug(ctx, "Generating Update payload")
	return generateBasePayload(ctx, data)
}

// GenerateDeletePayload generates a blank/default payload for the DELETE operation.
func GenerateDeletePayload(ctx context.Context, data ResourceModel) *client.UpdateDeviceManagementInterfaceRequest {
	tflog.Debug(ctx, "Generating Delete payload")
	return &client.UpdateDeviceManagementInterfaceRequest{
		Wan1: &client.UpdateDeviceManagementInterfaceRequestWan1{},
		Wan2: &client.UpdateDeviceManagementInterfaceRequestWan2{},
	}
}

// generateBasePayload maps the state data into the API request payload.
func generateBasePayload(ctx context.Context, data ResourceModel) *client.UpdateDeviceManagementInterfaceRequest {
	payload := &client.UpdateDeviceManagementInterfaceRequest{}

	// Map WAN1 configuration
	if !data.Wan1.IsNull() {
		wan1, diags := generateWan1Payload(ctx, data.Wan1)
		if diags.HasError() {
			tflog.Error(ctx, "Failed to map WAN1 configuration")
		}
		payload.Wan1 = wan1
	}

	// Map WAN2 configuration
	if !data.Wan2.IsNull() {
		wan2, diags := generateWan2Payload(ctx, data.Wan2)
		if diags.HasError() {
			tflog.Error(ctx, "Failed to map WAN2 configuration")
		}
		payload.Wan2 = wan2
	}

	return payload
}

/* Attribute Mapping Helpers */

// mapWAN maps a WAN configuration from the Terraform state to the API payload.
func generateWan1Payload(ctx context.Context, wanObject types.Object) (*client.UpdateDeviceManagementInterfaceRequestWan1, diag.Diagnostics) {
	var wan WANModel
	var diags diag.Diagnostics

	wanObject.As(ctx, &wan, basetypes.ObjectAsOptions{})

	payload := &client.UpdateDeviceManagementInterfaceRequestWan1{
		WanEnabled:       wan.WanEnabled.ValueStringPointer(),
		UsingStaticIp:    wan.UsingStaticIp.ValueBoolPointer(),
		StaticIp:         wan.StaticIp.ValueStringPointer(),
		StaticSubnetMask: wan.StaticSubnetMask.ValueStringPointer(),
		StaticGatewayIp:  wan.StaticGatewayIp.ValueStringPointer(),
		StaticDns:        utils.FlattenList(wan.StaticDns),
	}

	// Explicit Vlan mapping with type validation
	if !wan.Vlan.IsNull() {
		vlan := int32(wan.Vlan.ValueInt64())
		payload.Vlan = &vlan
	}

	return payload, diags
}

// mapWAN maps a WAN configuration from the Terraform state to the API payload.
func generateWan2Payload(ctx context.Context, wanObject types.Object) (*client.UpdateDeviceManagementInterfaceRequestWan2, diag.Diagnostics) {
	var wan WANModel
	var diags diag.Diagnostics

	wanObject.As(ctx, &wan, basetypes.ObjectAsOptions{})

	payload := &client.UpdateDeviceManagementInterfaceRequestWan2{
		WanEnabled:       wan.WanEnabled.ValueStringPointer(),
		UsingStaticIp:    wan.UsingStaticIp.ValueBoolPointer(),
		StaticIp:         wan.StaticIp.ValueStringPointer(),
		StaticSubnetMask: wan.StaticSubnetMask.ValueStringPointer(),
		StaticGatewayIp:  wan.StaticGatewayIp.ValueStringPointer(),
		StaticDns:        utils.FlattenList(wan.StaticDns),
	}

	// Explicit Vlan mapping with type validation
	if !wan.Vlan.IsNull() {
		vlan := int32(wan.Vlan.ValueInt64())
		payload.Vlan = &vlan
	}

	return payload, diags
}

/* API Call Abstractions   */

func CallCreateAPI(ctx context.Context, client *client.APIClient, payload *client.UpdateDeviceManagementInterfaceRequest, serial string) (map[string]interface{}, *http.Response, error) {
	return client.ManagementInterfaceApi.UpdateDeviceManagementInterface(ctx, serial).UpdateDeviceManagementInterfaceRequest(*payload).Execute()
}

func CallReadAPI(ctx context.Context, client *client.APIClient, serial string) (map[string]interface{}, *http.Response, error) {
	return client.ManagementInterfaceApi.GetDeviceManagementInterface(ctx, serial).Execute()
}

func CallUpdateAPI(ctx context.Context, client *client.APIClient, payload *client.UpdateDeviceManagementInterfaceRequest, serial string) (map[string]interface{}, *http.Response, error) {
	return client.ManagementInterfaceApi.UpdateDeviceManagementInterface(ctx, serial).UpdateDeviceManagementInterfaceRequest(*payload).Execute()
}

func CallDeleteAPI(ctx context.Context, client *client.APIClient, payload *client.UpdateDeviceManagementInterfaceRequest, serial string) (*http.Response, error) {
	_, resp, err := client.ManagementInterfaceApi.UpdateDeviceManagementInterface(ctx, serial).UpdateDeviceManagementInterfaceRequest(*payload).Execute()
	return resp, err
}

/* State Transformation    */

func MarshalStateFromAPI(ctx context.Context, apiResponse map[string]interface{}) (ResourceModel, diag.Diagnostics) {
	var state ResourceModel
	var diags diag.Diagnostics

	tflog.Debug(ctx, "Mapping API response to Terraform state", map[string]interface{}{
		"api_response": apiResponse,
	})

	// Map Serial
	if serial, ok := apiResponse["serial"].(string); ok {
		state.Serial = types.StringValue(serial)
	} else {
		state.Serial = types.StringNull()
	}

	// Map DDNS Hostnames
	if ddnsRaw, ok := apiResponse["ddnsHostnames"].(map[string]interface{}); ok {
		ddnsValues := map[string]attr.Value{
			"active_ddns_hostname": utils.SafeStringAttr(ddnsRaw, "activeDdnsHostname"),
			"ddns_hostname_wan1":   utils.SafeStringAttr(ddnsRaw, "ddnsHostnameWan1"),
			"ddns_hostname_wan2":   utils.SafeStringAttr(ddnsRaw, "ddnsHostnameWan2"),
		}

		ddnsObject, err := types.ObjectValue(DdnsHostnamesType, ddnsValues)
		if err != nil {
			diags.AddError("State Mapping Error", fmt.Sprintf("Failed to map DDNS hostnames: %s", err))
		}
		state.DDNSHostnames = ddnsObject
	} else {
		state.DDNSHostnames = types.ObjectNull(DdnsHostnamesType)
	}

	// Map WAN1
	if wan1Raw, ok := apiResponse["wan1"].(map[string]interface{}); ok {
		wan1Values := map[string]attr.Value{
			"wan_enabled":        utils.SafeStringAttr(wan1Raw, "wanEnabled"),
			"using_static_ip":    utils.SafeBoolAttr(wan1Raw, "usingStaticIp"),
			"static_ip":          utils.SafeStringAttr(wan1Raw, "staticIp"),
			"static_subnet_mask": utils.SafeStringAttr(wan1Raw, "staticSubnetMask"),
			"static_gateway_ip":  utils.SafeStringAttr(wan1Raw, "staticGatewayIp"),
			"static_dns":         utils.SafeListStringAttr(wan1Raw, "staticDns"),
			"vlan":               utils.SafeInt64Attr(wan1Raw, "vlan"),
		}

		wan1Object, err := types.ObjectValue(WANType, wan1Values)
		if err != nil {
			diags.AddError("State Mapping Error", fmt.Sprintf("Failed to map WAN1: %s", err))
		}
		state.Wan1 = wan1Object
	} else {
		state.Wan1 = types.ObjectNull(WANType)
	}

	// Map WAN2
	if wan2Raw, ok := apiResponse["wan2"].(map[string]interface{}); ok {
		wan2Values := map[string]attr.Value{
			"wan_enabled":        utils.SafeStringAttr(wan2Raw, "wanEnabled"),
			"using_static_ip":    utils.SafeBoolAttr(wan2Raw, "usingStaticIp"),
			"static_ip":          utils.SafeStringAttr(wan2Raw, "staticIp"),
			"static_subnet_mask": utils.SafeStringAttr(wan2Raw, "staticSubnetMask"),
			"static_gateway_ip":  utils.SafeStringAttr(wan2Raw, "staticGatewayIp"),
			"static_dns":         utils.SafeListStringAttr(wan2Raw, "staticDns"),
			"vlan":               utils.SafeInt64Attr(wan2Raw, "vlan"),
		}

		wan2Object, err := types.ObjectValue(WANType, wan2Values)
		if err != nil {
			diags.AddError("State Mapping Error", fmt.Sprintf("Failed to map WAN2: %s", err))
		}
		state.Wan2 = wan2Object
	} else {
		state.Wan2 = types.ObjectNull(WANType)
	}

	return state, diags
}

// mapWAN safely maps the WAN data into a types.Object
func mapWAN(ctx context.Context, raw map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	wanValues := map[string]attr.Value{
		"wan_enabled":        utils.SafeStringAttr(raw, "wanEnabled"),
		"using_static_ip":    utils.SafeBoolAttr(raw, "usingStaticIp"),
		"static_ip":          utils.SafeStringAttr(raw, "staticIp"),
		"static_subnet_mask": utils.SafeStringAttr(raw, "staticSubnetMask"),
		"static_gateway_ip":  utils.SafeStringAttr(raw, "staticGatewayIp"),
		"static_dns":         utils.SafeListStringAttr(raw, "staticDns"),
		"vlan":               utils.SafeInt64Attr(raw, "vlan"),
	}

	wanObject, err := types.ObjectValue(WANType, wanValues)
	if err != nil {
		diags.AddError("State Mapping Error", "Failed to map WAN configuration")
	}
	return wanObject, diags
}
