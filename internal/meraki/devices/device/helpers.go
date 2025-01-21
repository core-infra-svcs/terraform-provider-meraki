package device

import (
	"context"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
	"sort"
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

	payload := openApiClient.UpdateDeviceRequest{}
	payload.SetName("")
	payload.SetNotes("")
	payload.SetAddress("")
	//payload.SetFloorPlanId("")
	//payload.SetMoveMapMarker(false)
	payload.SetTags(nil)
	payload.SetLat(0)
	payload.SetLng(0)

	return &payload
}

// Generate base payload for API operations
func generateBasePayload(ctx context.Context, data ResourceModel) *openApiClient.UpdateDeviceRequest {
	payload := openApiClient.NewUpdateDeviceRequest()

	// Set Name
	if !data.Name.IsUnknown() && !data.Name.IsNull() {
		payload.SetName(data.Name.ValueString())
	}

	// Set Tags
	if !data.Tags.IsUnknown() && !data.Tags.IsNull() {
		var tags []string
		for _, element := range data.Tags.Elements() {
			if str, ok := element.(types.String); ok && !str.IsUnknown() && !str.IsNull() {
				tags = append(tags, str.ValueString())
			}
		}
		sort.Strings(tags)
		payload.SetTags(tags)
	}

	// Set Latitude - OMIT if null
	if !data.Lat.IsUnknown() && !data.Lat.IsNull() {
		payload.SetLat(float32(data.Lat.ValueFloat64()))
	}

	// Set Longitude - OMIT if null
	if !data.Lng.IsUnknown() && !data.Lng.IsNull() {
		payload.SetLng(float32(data.Lng.ValueFloat64()))
	}

	// Set Address - OMIT if null
	if !data.Address.IsUnknown() && !data.Address.IsNull() {
		payload.SetAddress(data.Address.ValueString())
	}

	// Set Notes - OMIT if null
	if !data.Notes.IsUnknown() && !data.Notes.IsNull() {
		payload.SetNotes(data.Notes.ValueString())
	}

	// Set MoveMapMarker (default false if null)
	if !data.MoveMapMarker.IsUnknown() && !data.MoveMapMarker.IsNull() {
		payload.SetMoveMapMarker(data.MoveMapMarker.ValueBool())
	}

	if data.MoveMapMarker.IsNull() {
		payload.SetMoveMapMarker(false)
	}

	// Set SwitchProfileId - OMIT if null
	if !data.SwitchProfileId.IsUnknown() && !data.SwitchProfileId.IsNull() {
		payload.SetSwitchProfileId(data.SwitchProfileId.ValueString())
	}

	// Set FloorPlanId - OMIT if null
	if !data.FloorPlanId.IsUnknown() && !data.FloorPlanId.IsNull() {
		payload.SetFloorPlanId(data.FloorPlanId.ValueString())
	}

	return payload
}

// CallCreateAPI Call the CREATE API
func CallCreateAPI(ctx context.Context, client *openApiClient.APIClient, data ResourceModel) (ResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	payload := GenerateCreatePayload(ctx, data)

	utils.LogPayload(ctx, payload)

	// Call the API
	i, httpResp, err := client.DevicesApi.UpdateDevice(ctx, data.Serial.ValueString()).
		UpdateDeviceRequest(*payload).Execute()

	utils.LogResponseBody(ctx, httpResp)

	if err := utils.HandleAPIError(ctx, httpResp, err, &diags); err != nil {
		tflog.Error(ctx, "CREATE operation failed", map[string]interface{}{
			"error":         err,
			"response_code": httpResp.StatusCode,
		})
		return ResourceModel{}, diags
	}

	// Map API response to Terraform state
	state, diags := MarshalStateFromAPI(ctx, i)
	return state, diags
}

// CallReadAPI Call the READ API
func CallReadAPI(ctx context.Context, client *openApiClient.APIClient, data ResourceModel) (ResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics

	i, httpResp, err := client.DevicesApi.GetDevice(ctx, data.Serial.ValueString()).Execute()

	utils.LogResponseBody(ctx, httpResp)

	if err := utils.HandleAPIError(ctx, httpResp, err, &diags); err != nil {
		tflog.Error(ctx, "READ operation failed", map[string]interface{}{
			"error":         err,
			"response_code": httpResp.StatusCode,
		})
		return ResourceModel{}, diags
	}

	state, diags := MarshalStateFromAPI(ctx, i)
	return state, diags
}

// CallUpdateAPI Call the UPDATE API
func CallUpdateAPI(ctx context.Context, client *openApiClient.APIClient, data ResourceModel) (ResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	payload := GenerateUpdatePayload(ctx, data)

	utils.LogPayload(ctx, payload)

	// Call the API
	i, httpResp, err := client.DevicesApi.UpdateDevice(ctx, data.Serial.ValueString()).
		UpdateDeviceRequest(*payload).Execute()

	utils.LogResponseBody(ctx, httpResp)

	if err := utils.HandleAPIError(ctx, httpResp, err, &diags); err != nil {
		tflog.Error(ctx, "UPDATE operation failed", map[string]interface{}{
			"error":         err,
			"response_code": httpResp.StatusCode,
		})
		return ResourceModel{}, diags
	}

	state, diags := MarshalStateFromAPI(ctx, i)
	return state, diags
}

// CallDeleteAPI Call the DELETE API
func CallDeleteAPI(ctx context.Context, client *openApiClient.APIClient, data ResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics
	payload := GenerateDeletePayload(ctx, data)

	utils.LogPayload(ctx, payload)

	// Call the API
	_, httpResp, err := client.DevicesApi.UpdateDevice(ctx, data.Serial.ValueString()).
		UpdateDeviceRequest(*payload).Execute()

	utils.LogResponseBody(ctx, httpResp)

	if err := utils.HandleAPIError(ctx, httpResp, err, &diags); err != nil {
		tflog.Error(ctx, "DELETE operation failed", map[string]interface{}{
			"error":         err,
			"response_code": httpResp.StatusCode,
		})
		return diags
	}

	return diags
}

// MarshalStateFromAPI marshals the API response (map) to the Terraform state model.
func MarshalStateFromAPI(ctx context.Context, apiData map[string]interface{}) (ResourceModel, diag.Diagnostics) {
	var diags diag.Diagnostics
	var data ResourceModel

	tflog.Trace(ctx, "[device] Starting to marshal API response to Terraform state")

	// Use utils.ExtractStringAttr for string mappings
	data.Serial, _ = utils.ExtractStringAttr(apiData, "serial")
	data.Id = data.Serial

	data.Name, _ = utils.ExtractStringAttr(apiData, "name")
	data.Address, _ = utils.ExtractStringAttr(apiData, "address")
	data.Notes, _ = utils.ExtractStringAttr(apiData, "notes")
	data.LanIp, _ = utils.ExtractStringAttr(apiData, "lanIp")
	data.Model, _ = utils.ExtractStringAttr(apiData, "model")

	value, exists := apiData["moveMapMarker"]
	if exists && value != nil {
		data.MoveMapMarker = types.BoolValue(value.(bool))
	} else {
		data.MoveMapMarker = types.BoolValue(false)
	}

	// Handle coordinates
	data.Lat, _ = utils.ExtractFloat64Attr(apiData, "lat")
	data.Lng, _ = utils.ExtractFloat64Attr(apiData, "lng")

	// Extract lists
	// Extract lists
	data.Tags, _ = utils.ExtractListStringAttr(apiData, "tags")

	// Convert []attr.Value to []string
	tags := make([]string, len(data.Tags.Elements()))
	for i, v := range data.Tags.Elements() {
		tags[i] = v.(basetypes.StringValue).ValueString()
	}

	// Sort the converted []string
	sort.Strings(tags)

	// Handle complex nested details separately
	data.Details, _ = MarshalDetailsFromAPI(ctx, apiData)

	data.BeaconIdParams, diags = MarshalBeaconIdParamsFromAPI(ctx, apiData)

	tflog.Trace(ctx, "[device] Successfully marshaled API response to Terraform state")
	return data, diags
}

// MarshalDetailsFromAPI handles extracting and marshaling the 'details' attribute from the API response.
func MarshalDetailsFromAPI(ctx context.Context, apiData map[string]interface{}) (types.List, diag.Diagnostics) {
	var diags diag.Diagnostics
	var detailsElements []attr.Value

	if detailsList, exists := apiData["details"].([]interface{}); exists {
		for _, item := range detailsList {
			detail := item.(map[string]interface{})

			objValue, objDiags := types.ObjectValue(
				map[string]attr.Type{
					"name":  types.StringType,
					"value": types.StringType,
				},
				map[string]attr.Value{
					"name":  types.StringValue(detail["name"].(string)),
					"value": types.StringValue(detail["value"].(string)),
				},
			)

			// Capture diagnostics if there's an error
			diags.Append(objDiags...)

			if objDiags.HasError() {
				tflog.Error(ctx, "Failed to create ObjectValue for details attribute", map[string]interface{}{
					"diagnostics": objDiags,
				})
				continue
			}

			// Append only the ObjectValue to the list
			detailsElements = append(detailsElements, objValue)
		}
	}

	// Create ListValue from ObjectValues
	listValue, listDiags := types.ListValue(
		types.ObjectType{
			AttrTypes: map[string]attr.Type{
				"name":  types.StringType,
				"value": types.StringType,
			},
		}, detailsElements)

	// Append list-level diagnostics
	diags.Append(listDiags...)

	return listValue, diags
}

func MarshalBeaconIdParamsFromAPI(ctx context.Context, apiData map[string]interface{}) (types.Object, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Define attribute types for beacon_id_params
	beaconIdParamsAttrTypes := map[string]attr.Type{
		"beacon_id": types.StringType,
		"major":     types.Int64Type,
		"minor":     types.Int64Type,
		"proximity": types.StringType,
		"uuid":      types.StringType,
	}

	// Check if beacon_id_params exists in the API response
	beaconParams, ok := apiData["beacon_id_params"].(map[string]interface{})
	if !ok || len(beaconParams) == 0 {
		// Return an empty/null object if no data is found
		return types.ObjectNull(beaconIdParamsAttrTypes), diags
	}

	// Extract fields from the beaconParams map
	beaconId := types.StringNull()
	major := types.Int64Null()
	minor := types.Int64Null()
	proximity := types.StringNull()
	uuid := types.StringNull()

	if v, exists := beaconParams["beacon_id"]; exists && v != nil {
		beaconId = types.StringValue(v.(string))
	}
	if v, exists := beaconParams["major"]; exists && v != nil {
		major = types.Int64Value(int64(v.(float64)))
	}
	if v, exists := beaconParams["minor"]; exists && v != nil {
		minor = types.Int64Value(int64(v.(float64)))
	}
	if v, exists := beaconParams["proximity"]; exists && v != nil {
		proximity = types.StringValue(v.(string))
	}
	if v, exists := beaconParams["uuid"]; exists && v != nil {
		uuid = types.StringValue(v.(string))
	}

	// Construct the ObjectValue for beacon_id_params
	beaconIdObject, objDiags := types.ObjectValue(
		beaconIdParamsAttrTypes,
		map[string]attr.Value{
			"beacon_id": beaconId,
			"major":     major,
			"minor":     minor,
			"proximity": proximity,
			"uuid":      uuid,
		},
	)
	diags.Append(objDiags...)

	if objDiags.HasError() {
		tflog.Error(ctx, "Failed to marshal beacon_id_params", map[string]interface{}{
			"diagnostics": objDiags,
		})
		return types.ObjectNull(beaconIdParamsAttrTypes), diags
	}

	return beaconIdObject, diags
}
