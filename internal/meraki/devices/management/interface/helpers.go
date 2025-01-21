package _interface

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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

	// Deserialize the Object into the WANModel
	wanObject.As(ctx, &wan, basetypes.ObjectAsOptions{})

	// Perform manual validation of required attributes
	attributes := wanObject.Attributes()

	// Ensure "wan_enabled" exists and is not null
	if wanEnabled, ok := attributes["wan_enabled"]; !ok || wanEnabled.IsNull() {
		diags.AddError("Validation Error", "WAN configuration must include 'wan_enabled'")
		return nil, diags
	}

	payload := &client.UpdateDeviceManagementInterfaceRequestWan1{
		WanEnabled:       wan.WanEnabled.ValueStringPointer(),
		UsingStaticIp:    wan.UsingStaticIp.ValueBoolPointer(),
		StaticIp:         wan.StaticIp.ValueStringPointer(),
		StaticSubnetMask: wan.StaticSubnetMask.ValueStringPointer(),
		StaticGatewayIp:  wan.StaticGatewayIp.ValueStringPointer(),
		StaticDns:        utils.FlattenList(wan.StaticDns),
	}

	// Explicit VLAN mapping with type validation and null handling
	if !wan.Vlan.IsNull() {
		vlan := int32(wan.Vlan.ValueInt64())
		if vlan == 0 {
			// Set the VLAN to null if it's set to 0
			payload.Vlan = nil
		} else {
			payload.Vlan = &vlan
		}
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

	// Explicit VLAN mapping with type validation and null handling
	if !wan.Vlan.IsNull() {
		vlan := int32(wan.Vlan.ValueInt64())
		if vlan == 0 {
			// Set the VLAN to null if it's set to 0
			payload.Vlan = nil
		} else {
			payload.Vlan = &vlan
		}
	}

	return payload, diags
}

/* API Call Abstractions   */

func CallCreateAPI(ctx context.Context, client *client.APIClient, data ResourceModel) (ResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Generate API payload
	payload := GenerateCreatePayload(ctx, data)

	// Call the CREATE API
	apiResponse, httpResp, err := client.ManagementInterfaceApi.UpdateDeviceManagementInterface(ctx, data.Serial.ValueString()).
		UpdateDeviceManagementInterfaceRequest(*payload).Execute()
	if err := utils.HandleAPIError(ctx, httpResp, err, &diags); err != nil {
		return ResourceModel{}, diags
	}

	// Marshal API response into state
	state, apiDiags := MarshalStateFromAPI(ctx, apiResponse, data)
	diags.Append(apiDiags...)

	// Set the ID to the serial value
	state.Id = types.StringValue(data.Serial.ValueString())

	// Retain the serial in the state
	state.Serial = data.Serial

	return state, diags
}

func CallReadAPIDataSource(ctx context.Context, client *client.APIClient, data DataSourceModel) (ResourceModel, diag.Diagnostics) {
	state, diags := CallReadAPI(ctx, client, ResourceModel(data))

	return state, diags
}

func CallReadAPI(ctx context.Context, client *client.APIClient, data ResourceModel) (ResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Call the READ API
	apiResponse, httpResp, err := client.ManagementInterfaceApi.GetDeviceManagementInterface(ctx, data.Serial.ValueString()).Execute()
	if err := utils.HandleAPIError(ctx, httpResp, err, &diags); err != nil {
		return ResourceModel{}, diags
	}

	// Marshal API response into state
	state, apiDiags := MarshalStateFromAPI(ctx, apiResponse, data)
	diags.Append(apiDiags...)

	// Retain the serial in the state
	state.Serial = data.Serial

	// Set the ID to the serial value
	state.Id = types.StringValue(data.Serial.ValueString())

	return state, diags
}

func CallUpdateAPI(ctx context.Context, client *client.APIClient, data ResourceModel) (ResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Generate API payload
	payload := GenerateUpdatePayload(ctx, data)

	// Call the UPDATE API
	apiResponse, httpResp, err := client.ManagementInterfaceApi.UpdateDeviceManagementInterface(ctx, data.Serial.ValueString()).
		UpdateDeviceManagementInterfaceRequest(*payload).Execute()
	if err := utils.HandleAPIError(ctx, httpResp, err, &diags); err != nil {
		return ResourceModel{}, diags
	}

	// Marshal API response into state
	state, apiDiags := MarshalStateFromAPI(ctx, apiResponse, data)
	diags.Append(apiDiags...)

	// Retain the serial in the state
	state.Serial = data.Serial

	// Set the ID to the serial value
	state.Id = types.StringValue(data.Serial.ValueString())

	return state, diags
}

func CallDeleteAPI(ctx context.Context, client *client.APIClient, data ResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	// Generate DELETE payload
	payload := GenerateDeletePayload(ctx, data)

	// Call the DELETE API
	_, httpResp, err := client.ManagementInterfaceApi.UpdateDeviceManagementInterface(ctx, data.Serial.ValueString()).
		UpdateDeviceManagementInterfaceRequest(*payload).Execute()
	if err := utils.HandleAPIError(ctx, httpResp, err, &diags); err != nil {
		return diags
	}

	return diags
}

/* State Transformation    */

func MarshalStateFromAPI(ctx context.Context, apiResponse map[string]interface{}, data ResourceModel) (ResourceModel, diag.Diagnostics) {
	var state ResourceModel
	var diags diag.Diagnostics

	tflog.Debug(ctx, "Mapping API response to Terraform state", map[string]interface{}{
		"api_response": apiResponse,
	})

	// Map Serial/Id
	state.Id = types.StringValue(data.Serial.ValueString())
	state.Serial = types.StringValue(data.Serial.ValueString())

	// Map DDNS Hostnames
	if ddnsRaw, ok := apiResponse["ddnsHostnames"].(map[string]interface{}); ok {
		ddns, ddnsDiags := mapDDNS(ctx, ddnsRaw)
		diags.Append(ddnsDiags...)
		state.DDNSHostnames = ddns
	} else {
		state.DDNSHostnames = types.ObjectNull(DdnsHostnamesType)
	}

	// Map WAN1
	if wan1Raw, ok := apiResponse["wan1"].(map[string]interface{}); ok {
		wan1, wan1Diags := mapWAN(ctx, wan1Raw)
		diags.Append(wan1Diags...)
		state.Wan1 = wan1
	} else {
		state.Wan1 = types.ObjectNull(WANType)
	}

	// Map WAN2
	if wan2Raw, ok := apiResponse["wan2"].(map[string]interface{}); ok {
		wan2, wan2Diags := mapWAN(ctx, wan2Raw)
		diags.Append(wan2Diags...)
		state.Wan2 = wan2
	} else {
		state.Wan2 = types.ObjectNull(WANType)
	}

	return state, diags
}

func mapWAN(ctx context.Context, raw map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics
	wanValues := map[string]attr.Value{
		"wan_enabled":        utils.SafeStringAttr(raw, "wanEnabled"),
		"using_static_ip":    utils.SafeBoolAttr(raw, "usingStaticIp"),
		"static_ip":          utils.SafeStringAttr(raw, "staticIp"), // Defaults to null if missing
		"static_subnet_mask": utils.SafeStringAttr(raw, "staticSubnetMask"),
		"static_gateway_ip":  utils.SafeStringAttr(raw, "staticGatewayIp"),
		"static_dns":         utils.SafeListStringAttr(raw, "staticDns"), // Defaults to an empty list
		"vlan":               utils.SafeInt64Attr(raw, "vlan"),
	}

	// Create the ObjectValue using WANType
	wanObject, err := types.ObjectValue(WANType, wanValues)
	if err != nil {
		diags.AddError("State Mapping Error", fmt.Sprintf("Failed to map WAN: %s", err))
	}
	return wanObject, diags
}

// mapWAN safely maps the WAN data into a types.Object
func mapDDNS(ctx context.Context, raw map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Convert raw attributes to attr.Values
	ddnsValues := map[string]attr.Value{
		"active_ddns_hostname": types.StringValue(raw["activeDdnsHostname"].(string)),
		"ddns_hostname_wan1":   types.StringValue(raw["ddnsHostnameWan1"].(string)),
		"ddns_hostname_wan2":   types.StringValue(raw["ddnsHostnameWan2"].(string)),
	}

	// Create the ObjectValue using DdnsHostnamesType
	ddnsObject, err := types.ObjectValue(DdnsHostnamesType, ddnsValues)
	if err != nil {
		diags.AddError("State Mapping Error", fmt.Sprintf("Failed to map DDNS hostnames: %s", err))
	}

	return ddnsObject, diags
}
