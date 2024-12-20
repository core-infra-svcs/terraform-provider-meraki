package device

import (
	"context"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	openApiClient "github.com/meraki/dashboard-api-go/client"
	"net/http"
)

// GenerateCreatePayload Generate payload for CREATE operation
func GenerateCreatePayload(ctx context.Context, data ResourceModel) *openApiClient.UpdateDeviceRequest {
	return generateBasePayload(ctx, data)
}

// GenerateUpdatePayload Generate payload for UPDATE operation
func GenerateUpdatePayload(ctx context.Context, data ResourceModel) *openApiClient.UpdateDeviceRequest {
	return generateBasePayload(ctx, data)
}

// GenerateDeletePayload Generate blank payload for DELETE operation
func GenerateDeletePayload(ctx context.Context, data ResourceModel) *openApiClient.UpdateDeviceRequest {
	return &openApiClient.UpdateDeviceRequest{}
}

// Generate base payload for API operations
func generateBasePayload(ctx context.Context, data ResourceModel) *openApiClient.UpdateDeviceRequest {
	payload := openApiClient.NewUpdateDeviceRequest()

	if !data.Name.IsNull() {
		payload.SetName(data.Name.ValueString())
	}

	if !data.Tags.IsNull() && !data.Tags.IsUnknown() {
		var tags []string
		for _, element := range data.Tags.Elements() {
			if str, ok := element.(types.String); ok && !str.IsNull() && !str.IsUnknown() {
				tags = append(tags, str.ValueString())
			}
		}
		payload.SetTags(tags)
	}
	if !data.Lat.IsNull() {
		payload.SetLat(float32(data.Lat.ValueFloat64()))
	}

	if !data.Lng.IsNull() {
		payload.SetLng(float32(data.Lng.ValueFloat64()))
	}

	if !data.Address.IsNull() {
		payload.SetAddress(data.Address.ValueString())
	}

	if !data.Notes.IsNull() {
		payload.SetNotes(data.Notes.ValueString())
	}

	if !data.MoveMapMarker.IsNull() {
		payload.SetMoveMapMarker(data.MoveMapMarker.ValueBool())
	}

	if !data.SwitchProfileId.IsNull() {
		payload.SetSwitchProfileId(data.SwitchProfileId.ValueString())
	}

	if !data.FloorPlanId.IsNull() {
		payload.SetFloorPlanId(data.FloorPlanId.ValueString())
	}

	return payload
}

// CallCreateAPI Call the CREATE API
func CallCreateAPI(ctx context.Context, client *openApiClient.APIClient, data ResourceModel) (ResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	payload := GenerateCreatePayload(ctx, data)

	_, httpResp, err := client.DevicesApi.UpdateDevice(ctx, data.Serial.ValueString()).UpdateDeviceRequest(*payload).Execute()
	if err := utils.HandleAPIError(ctx, httpResp, err, &diags); err != nil {
		return ResourceModel{}, diags
	}

	state, diags := MarshalStateFromAPI(ctx, httpResp, data)
	return state, diags
}

// CallReadAPI Call the READ API
func CallReadAPI(ctx context.Context, client *openApiClient.APIClient, data ResourceModel) (ResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// API call without using the response directly
	_, httpResp, err := client.DevicesApi.GetDevice(ctx, data.Serial.ValueString()).Execute()
	if err := utils.HandleAPIError(ctx, httpResp, err, &diags); err != nil {
		return ResourceModel{}, diags
	}

	// Use httpResp to marshal the state
	state, diags := MarshalStateFromAPI(ctx, httpResp, data)
	return state, diags
}

// CallUpdateAPI Call the UPDATE API
func CallUpdateAPI(ctx context.Context, client *openApiClient.APIClient, data ResourceModel) (ResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	payload := GenerateUpdatePayload(ctx, data)

	_, httpResp, err := client.DevicesApi.UpdateDevice(ctx, data.Serial.ValueString()).UpdateDeviceRequest(*payload).Execute()
	if err := utils.HandleAPIError(ctx, httpResp, err, &diags); err != nil {
		return ResourceModel{}, diags
	}

	state, diags := MarshalStateFromAPI(ctx, httpResp, data)
	return state, diags
}

// CallDeleteAPI Call the DELETE API
func CallDeleteAPI(ctx context.Context, client *openApiClient.APIClient, data ResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	payload := GenerateDeletePayload(ctx, data)

	_, httpResp, err := client.DevicesApi.UpdateDevice(ctx, data.Serial.ValueString()).UpdateDeviceRequest(*payload).Execute()
	if err := utils.HandleAPIError(ctx, httpResp, err, &diags); err != nil {
		return diags
	}

	return diags
}

// MarshalStateFromAPI Map API response to Terraform state
func MarshalStateFromAPI(ctx context.Context, httpResp *http.Response, data ResourceModel) (ResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Extract API response into a map
	apiData, err := utils.ExtractResponseToMap(httpResp)
	if err != nil {
		diags.AddError("Failed to parse API response", err.Error())
		return ResourceModel{}, diags
	}

	// Map attributes from API response to Terraform state
	mapAttribute := func(apiKey string, stateField *types.String, optional bool) {
		if value, ok := apiData[apiKey]; ok && value != nil {
			if strValue, ok := value.(string); ok {
				*stateField = types.StringValue(strValue)
			}
		} else if !optional {
			*stateField = types.StringNull()
		}
	}

	mapFloat64Attribute := func(apiKey string, stateField *types.Float64) {
		if value, ok := apiData[apiKey]; ok && value != nil {
			if floatValue, ok := value.(float64); ok {
				*stateField = types.Float64Value(floatValue)
			}
		} else {
			*stateField = types.Float64Null()
		}
	}

	mapBoolAttribute := func(apiKey string, stateField *types.Bool) {
		if value, ok := apiData[apiKey]; ok && value != nil {
			if boolValue, ok := value.(bool); ok {
				*stateField = types.BoolValue(boolValue)
			}
		} else {
			*stateField = types.BoolNull()
		}
	}

	// Map individual attributes
	mapAttribute("name", &data.Name, true)
	mapFloat64Attribute("lat", &data.Lat)
	mapFloat64Attribute("lng", &data.Lng)
	mapAttribute("address", &data.Address, true)
	mapAttribute("notes", &data.Notes, true)
	mapAttribute("networkId", &data.NetworkId, true)
	mapAttribute("serial", &data.Serial, false)
	mapAttribute("model", &data.Model, true)
	mapAttribute("mac", &data.Mac, true)
	mapAttribute("lanIp", &data.LanIp, true)
	mapAttribute("firmware", &data.Firmware, true)
	mapAttribute("floorPlanId", &data.FloorPlanId, true)
	mapAttribute("url", &data.Url, true)
	mapAttribute("switchProfileId", &data.SwitchProfileId, true)
	mapBoolAttribute("moveMapMarker", &data.MoveMapMarker)

	// Set the resource ID to the serial value
	data.Id = data.Serial

	return data, diags
}
