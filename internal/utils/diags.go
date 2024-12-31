package utils

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"io"
	"net/http"
	"testing"
)

// LogCurrentState LogCurrentStat Custom test check function to log the current state
func LogCurrentState(t *testing.T) resource.TestCheckFunc {
	/*
		// Example test step
		{
			Config: testAccDevicesClaimResourceConfigCreateNetwork(os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")),
			Check: resource.ComposeAggregateTestCheckFunc(
				tools.LogCurrentState(t), // Add this line to log the state
				resource.TestCheckResourceAttr("meraki_network.test", "name", "test_acc_network_claim_device"),
			),
		},
	*/

	return func(s *terraform.State) error {
		for _, rs := range s.RootModule().Resources {
			rsType := fmt.Sprintf("[DEBUG] Resource: %s", rs.Type)
			fmt.Println(rsType)
			for key, value := range rs.Primary.Attributes {
				rsAttr := fmt.Sprintf("[DEBUG] %s: %s", key, value)
				fmt.Println(rsAttr)
			}
		}
		return nil
	}
}

// LogPayload logs the request payload before making an API call.
func LogPayload(ctx context.Context, payload interface{}) {
	payloadJSON, err := json.Marshal(payload)
	if err != nil {
		tflog.Warn(ctx, "Failed to serialize payload to JSON", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}
	tflog.Debug(ctx, "Request Payload", map[string]interface{}{
		"payload": string(payloadJSON),
	})
}

// LogResponseBody logs the raw response body from the API.
func LogResponseBody(ctx context.Context, httpResp *http.Response) {
	if httpResp == nil || httpResp.Body == nil {
		tflog.Warn(ctx, "Response body is nil or empty")
		return
	}

	defer httpResp.Body.Close()
	bodyBytes, err := io.ReadAll(httpResp.Body)
	if err != nil {
		tflog.Warn(ctx, "Failed to read response body", map[string]interface{}{
			"error": err.Error(),
		})
		return
	}

	// Log the raw body for debugging
	tflog.Debug(ctx, "API Response Body", map[string]interface{}{
		"body": string(bodyBytes),
	})

	// Reset the body for further use by re-wrapping it
	httpResp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
}
