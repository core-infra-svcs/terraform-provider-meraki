package device

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
	"net/http"
	"strings"
	"time"
)

func mapPayload(plan *ResourceModel) (openApiClient.UpdateDeviceRequest, diag.Diagnostics) {
	var diags diag.Diagnostics
	payload := openApiClient.NewUpdateDeviceRequest()

	//    Name
	if !plan.Name.IsNull() && !plan.Name.IsUnknown() {
		payload.SetName(plan.Name.ValueString())
	}

	//    Tags
	if !plan.Tags.IsNull() && !plan.Tags.IsUnknown() {
		var tags []string
		for _, attribute := range plan.Tags.Elements() {
			tag := fmt.Sprint(strings.Trim(attribute.String(), "\""))
			tags = append(tags, tag)
		}
		payload.SetTags(tags)
	}

	//    Lat
	if !plan.Lat.IsNull() && !plan.Lat.IsUnknown() {
		payload.SetLat(float32(plan.Lat.ValueFloat64()))

	}

	//    Lng
	if !plan.Lng.IsNull() && !plan.Lng.IsUnknown() {
		payload.SetLng(float32(plan.Lng.ValueFloat64()))
	}

	//    Address
	if !plan.Address.IsNull() && !plan.Address.IsUnknown() {
		payload.SetAddress(plan.Address.ValueString())
	}

	//    Notes
	if !plan.Notes.IsNull() && !plan.Notes.IsUnknown() {
		payload.SetNotes(plan.Notes.ValueString())
	}

	//    MoveMapMarker
	if !plan.MoveMapMarker.IsNull() && !plan.MoveMapMarker.IsUnknown() {
		payload.SetMoveMapMarker(plan.MoveMapMarker.ValueBool())
	}

	//    SwitchProfileId
	if !plan.SwitchProfileId.IsNull() && !plan.SwitchProfileId.IsUnknown() {
		payload.SetSwitchProfileId(plan.SwitchProfileId.ValueString())
	}

	//    FloorPlanId
	if !plan.FloorPlanId.IsNull() && !plan.FloorPlanId.IsUnknown() {
		payload.SetFloorPlanId(plan.FloorPlanId.ValueString())
	}

	return *payload, diags

}

// mapApiResponseToState updates the resource state with the provided api data.
func mapApiResponseToState(ctx context.Context, state *ResourceModel, httpResp *http.Response) diag.Diagnostics {
	var diags diag.Diagnostics

	httpRespMap, httpRespMapErr := utils.ExtractResponseToMap(httpResp)
	if httpRespMapErr != nil {
		diags.AddError("Failed to Unmarshal HttpResp", httpRespMapErr.Error())
	}

	// "name": "My AP",
	if state.Name.IsNull() || state.Name.IsUnknown() {
		name, err := utils.ExtractStringAttr(httpRespMap, "name")
		if err != nil {
			diags.Append(err...)
		}
		state.Name = name
	}

	//  "lat": 37.4180951010362,
	if state.Lat.IsNull() || state.Lat.IsUnknown() {
		lat, err := utils.ExtractFloat64Attr(httpRespMap, "lat")
		if err != nil {
			diags = append(diags, err...)

		}
		state.Lat = lat
	}

	//  "lng": -122.098531723022,
	if state.Lng.IsNull() || state.Lng.IsUnknown() {
		lng, err := utils.ExtractFloat64Attr(httpRespMap, "lng")
		if err != nil {
			diags = append(diags, err...)
		}
		state.Lng = lng
	}

	//  "address": "1600 Pennsylvania Ave",
	if state.Address.IsNull() || state.Address.IsUnknown() {
		address, err := utils.ExtractStringAttr(httpRespMap, "address")
		if err != nil {
			diags = append(diags, err...)
		}
		state.Address = address
	}

	//  "notes": "My AP's note",
	if state.Notes.IsNull() || state.Notes.IsUnknown() {
		notes, err := utils.ExtractStringAttr(httpRespMap, "notes")
		if err != nil {
			diags = append(diags, err...)
		}
		state.Notes = notes
	}

	//  "tags": [
	//    " recently-added "
	//  ],
	if state.Tags.IsNull() || state.Tags.IsUnknown() {
		tags, err := utils.ExtractListStringAttr(httpRespMap, "tags")
		if err != nil {
			diags = append(diags, err...)
		}

		state.Tags = tags
	}

	//  "networkId": "N_24329156",
	if state.NetworkId.IsNull() || state.NetworkId.IsUnknown() {
		networkId, err := utils.ExtractStringAttr(httpRespMap, "networkId")
		if err != nil {
			diags = append(diags, err...)
		}
		state.NetworkId = networkId
	}

	//  "serial": "Q234-ABCD-5678",
	if state.Serial.IsNull() || state.Serial.IsUnknown() {
		diags.AddError("Missing Serial", "Missing Serial After State Update")
	}

	//  "model": "MR34",
	if state.Model.IsNull() || state.Model.IsUnknown() {
		model, err := utils.ExtractStringAttr(httpRespMap, "model")
		if err != nil {
			diags = append(diags, err...)
		}
		state.Model = model
	}

	//  "mac": "00:11:22:33:44:55",
	if state.Mac.IsNull() || state.Mac.IsUnknown() {
		mac, err := utils.ExtractStringAttr(httpRespMap, "mac")
		if err != nil {
			diags = append(diags, err...)
		}
		state.Mac = mac
	}

	//  "lanIp": "1.2.3.4",
	if state.LanIp.IsNull() || state.LanIp.IsUnknown() {
		lanIp, err := utils.ExtractStringAttr(httpRespMap, "lanIp")
		if err != nil {
			diags = append(diags, err...)
		}
		state.LanIp = lanIp
	}

	//  "firmware": "wireless-25-14",
	if state.Firmware.IsNull() || state.Firmware.IsUnknown() {
		firmware, err := utils.ExtractStringAttr(httpRespMap, "firmware")
		if err != nil {
			diags = append(diags, err...)
		}
		state.Firmware = firmware
	}

	//  "floorPlanId": "g_2176982374",
	if state.FloorPlanId.IsNull() || state.FloorPlanId.IsUnknown() {
		floorPlanId, err := utils.ExtractStringAttr(httpRespMap, "floorPlanId")
		if err != nil {
			diags = append(diags, err...)
		}
		state.FloorPlanId = floorPlanId
	}

	//  "details": [
	//    {
	//      "name": "Catalyst serial",
	//      "value": "123ABC"
	//    }
	//  ],
	if state.Details.IsNull() || state.Details.IsUnknown() {

		detailAttr := map[string]attr.Type{
			"name":  types.StringType,
			"value": types.StringType,
		}

		detailsAttrs := types.ObjectType{AttrTypes: detailAttr}

		_, ok := httpRespMap["details"].([]map[string]interface{})
		if ok {

			detailsList, err := utils.ExtractListAttr(httpRespMap, "details", detailsAttrs)
			if err.HasError() {
				tflog.Error(context.Background(), fmt.Sprintf("%s", err.Errors()))
			}

			state.Details = detailsList

		} else {
			detailsArrayObjNull := types.ListNull(detailsAttrs)
			state.Details = detailsArrayObjNull
		}

	}

	//  "beaconIdParams": {
	//    "uuid": "00000000-0000-0000-0000-000000000000",
	//    "major": 5,
	//    "minor": 3
	if state.BeaconIdParams.IsNull() || state.BeaconIdParams.IsUnknown() {
		beaconIdParamsAttrs := map[string]attr.Type{
			"uuid":  types.StringType,
			"major": types.Int64Type,
			"minor": types.Int64Type,
		}

		beaconIdParamsResp, ok := httpRespMap["beaconIdParams"].(map[string]interface{})
		if ok {
			var beaconIdParams BeaconIdParamsModel

			// uuid
			uuid, err := utils.ExtractStringAttr(beaconIdParamsResp, "uuid")
			if err.HasError() {
				diags.AddError("uuid Attr", fmt.Sprintf("%s", err.Errors()))
			}

			beaconIdParams.Uuid = uuid

			// major
			major, err := utils.ExtractInt32Attr(beaconIdParamsResp, "major")
			if err.HasError() {
				diags.AddError("major Attr", fmt.Sprintf("%s", err.Errors()))
			}

			beaconIdParams.Major = major

			// minor
			minor, err := utils.ExtractInt32Attr(beaconIdParamsResp, "minor")
			if err.HasError() {
				diags.AddError("minor Attr", fmt.Sprintf("%s", err.Errors()))
			}

			beaconIdParams.Minor = minor

			beaconIdParamsObj, err := types.ObjectValueFrom(ctx, beaconIdParamsAttrs, beaconIdParams)
			if err.HasError() {
				diags.AddError("beaconIdParamsObj Attr", fmt.Sprintf("%s", err.Errors()))
			}

			state.BeaconIdParams = beaconIdParamsObj
		} else {
			beaconIdParamsObjNull := types.ObjectNull(beaconIdParamsAttrs)
			state.BeaconIdParams = beaconIdParamsObjNull
		}

	}

	// url
	if state.Url.IsNull() || state.Url.IsUnknown() {
		url, err := utils.ExtractStringAttr(httpRespMap, "url")
		if err.HasError() {
			diags.AddError("url Attr", fmt.Sprintf("%s", err.Errors()))
		}
		state.Url = url
	}

	// SwitchProfileId
	if state.SwitchProfileId.IsNull() || state.SwitchProfileId.IsUnknown() {
		switchProfileId, err := utils.ExtractStringAttr(httpRespMap, "switchProfileId")
		if err.HasError() {
			diags.AddError("switchProfileId Attr", fmt.Sprintf("%s", err.Errors()))
		}
		state.SwitchProfileId = switchProfileId
	}

	// MoveMapMarker
	if state.MoveMapMarker.IsNull() || state.MoveMapMarker.IsUnknown() {
		moveMapMarker, err := utils.ExtractBoolAttr(httpRespMap, "moveMapMarker")
		if err.HasError() {
			diags.AddError("moveMapMarker Attr", fmt.Sprintf("%s", err.Errors()))
		}
		state.MoveMapMarker = moveMapMarker
	}

	// Set ID for the new resource.
	if state.Id.IsNull() || state.Id.IsUnknown() {
		state.Id = state.Serial
	}

	return diags
}

// RetryAPICall encapsulates the retry logic for making API calls.
func RetryAPICall(ctx context.Context, maxRetries int, retryDelay time.Duration, apiCall func() (map[string]interface{}, *http.Response, error)) (map[string]interface{}, *http.Response, error) {
	for i := 0; i <= maxRetries; i++ {
		resp, httpResp, err := apiCall()
		if err == nil {
			return resp, httpResp, nil
		}

		if httpResp != nil && httpResp.StatusCode >= 400 && httpResp.StatusCode < 500 {
			tflog.Warn(ctx, "4xx error encountered; retrying", map[string]interface{}{
				"attempt":        i + 1,
				"httpStatusCode": httpResp.StatusCode,
				"error":          err.Error(),
			})
			time.Sleep(retryDelay)
		} else {
			return nil, httpResp, err
		}
	}
	return nil, nil, fmt.Errorf("exceeded maximum retry attempts")
}

// HandleAPICall manages the full lifecycle of an API call, including retries and error logging.
func HandleAPICall(ctx context.Context, client *openApiClient.APIClient, apiCall func() (map[string]interface{}, *http.Response, error)) (map[string]interface{}, *http.Response, error) {
	maxRetries := client.GetConfig().MaximumRetries
	retryDelay := time.Duration(client.GetConfig().Retry4xxErrorWaitTime)

	resp, httpResp, err := RetryAPICall(ctx, maxRetries, retryDelay, apiCall)

	return resp, httpResp, err
}
