package bandwidth

import (
	"context"
	"encoding/json"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/jsontypes"
	"io"
)

func extractHttpResponseUplinkBandwidthResource(ctx context.Context, httpRespBody io.ReadCloser, apiResponse *resourceModelApiResponse, data *resourceModel) (*resourceModel, error) {

	if err := json.NewDecoder(httpRespBody).Decode(apiResponse); err != nil {
		return data, err
	}

	data.BandwidthLimitCellularLimitDown = apiResponse.UplinkBandwidthLimits.Cellular.LimitDown
	data.BandwidthLimitCellularLimitUp = apiResponse.UplinkBandwidthLimits.Cellular.LimitUp
	data.BandwidthLimitWan2LimitDown = apiResponse.UplinkBandwidthLimits.Wan2.LimitDown
	data.BandwidthLimitWan2LimitUp = apiResponse.UplinkBandwidthLimits.Wan2.LimitUp
	data.BandwidthLimitWan1LimitDown = apiResponse.UplinkBandwidthLimits.Wan1.LimitDown
	data.BandwidthLimitWan1LimitUp = apiResponse.UplinkBandwidthLimits.Wan1.LimitUp

	if data.BandwidthLimitWan1LimitDown.IsUnknown() {
		data.BandwidthLimitWan1LimitDown = jsontypes.Int64Null()
	}
	if data.BandwidthLimitWan1LimitUp.IsUnknown() {
		data.BandwidthLimitWan1LimitUp = jsontypes.Int64Null()
	}
	if data.BandwidthLimitWan2LimitDown.Int64Value.IsUnknown() {
		data.BandwidthLimitWan2LimitDown = jsontypes.Int64Null()
	}
	if data.BandwidthLimitWan2LimitUp.IsUnknown() {
		data.BandwidthLimitWan2LimitUp = jsontypes.Int64Null()
	}
	if data.BandwidthLimitCellularLimitDown.IsUnknown() {
		data.BandwidthLimitCellularLimitDown = jsontypes.Int64Null()
	}
	if data.BandwidthLimitCellularLimitUp.IsUnknown() {
		data.BandwidthLimitCellularLimitUp = jsontypes.Int64Null()
	}

	return data, nil
}
