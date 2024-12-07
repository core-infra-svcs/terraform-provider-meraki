package claim

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	openApiClient "github.com/meraki/dashboard-api-go/client"
	"io"
	"net/http"
	"strings"
)

func handleError(ctx context.Context, err error, httpResp *http.Response, message string, resp *resource.CreateResponse) {
	resp.Diagnostics.AddError(
		message,
		fmt.Sprintf("Could not perform operation, unexpected error: %s", err),
	)
	if httpResp != nil {
		var responseBody string
		if httpResp.Body != nil {
			bodyBytes, readErr := io.ReadAll(httpResp.Body)
			if readErr == nil {
				responseBody = string(bodyBytes)
			}
		}
		tflog.Error(ctx, "API Request failed", map[string]interface{}{
			"error":          err.Error(),
			"httpStatusCode": httpResp.StatusCode,
			"responseBody":   responseBody,
		})
		resp.Diagnostics.AddError(
			"API Request failed",
			fmt.Sprintf("HTTP Response: %v\nResponse Body: %s", httpResp, responseBody),
		)
	}
}

// Extracts strings from a given attribute set
func extractSerials(serials types.Set) []string {
	var results []string
	for _, serial := range serials.Elements() {
		results = append(results, strings.Trim(serial.String(), "\""))
	}

	return results
}

// difference returns elements in 'a' that are not in 'b'.
func difference(a, b []string) []string {
	bMap := make(map[string]bool)
	for _, item := range b {
		bMap[item] = true
	}

	var diff []string
	for _, item := range a {
		if _, found := bMap[item]; !found {
			diff = append(diff, item)
		}
	}
	return diff
}

// Handles claiming or removing devices based on the provided action
func manageDeviceClaims(ctx context.Context, client *openApiClient.APIClient, networkID string, serials []string, isAdd bool, resp *resource.UpdateResponse) error {
	var httpResp *http.Response
	var err error

	if isAdd {
		claimRequest := *openApiClient.NewClaimNetworkDevicesRequest(serials)
		httpResp, err = client.NetworksApi.ClaimNetworkDevices(ctx, networkID).ClaimNetworkDevicesRequest(claimRequest).Execute()
	} else {
		for _, serial := range serials {
			removeRequest := *openApiClient.NewRemoveNetworkDevicesRequest(serial)
			httpResp, err = client.NetworksApi.RemoveNetworkDevices(ctx, networkID).RemoveNetworkDevicesRequest(removeRequest).Execute()
		}
	}

	if err != nil {
		resp.Diagnostics.AddError("Error managing devices", err.Error())
		if httpResp != nil {
			resp.Diagnostics.AddError("HTTP Response", utils.HttpDiagnostics(httpResp))
		}
		return fmt.Errorf("failed to manage devices: %w", err)
	}

	expectedStatusCode := 200
	if !isAdd {
		expectedStatusCode = 204
	}

	if httpResp.StatusCode != expectedStatusCode {
		resp.Diagnostics.AddError("Unexpected HTTP Response Status Code", fmt.Sprintf("Expected status code %d but got %v", expectedStatusCode, httpResp.StatusCode))
		return fmt.Errorf("unexpected status code: %d", httpResp.StatusCode)
	}
	return nil
}

// Handles the API call to remove devices from a network
func removeDevices(ctx context.Context, client *openApiClient.APIClient, networkID string, serials []string, resp *resource.DeleteResponse) error {

	for _, serial := range serials {
		removeNetworkDevices := *openApiClient.NewRemoveNetworkDevicesRequest(serial)
		httpResp, err := client.NetworksApi.RemoveNetworkDevices(ctx, networkID).RemoveNetworkDevicesRequest(removeNetworkDevices).Execute()
		if err != nil {
			resp.Diagnostics.AddError("Error removing devices", err.Error())
			return fmt.Errorf("failed to remove devices: %w", err)
		}

		if httpResp.StatusCode != 204 {
			resp.Diagnostics.AddError("Unexpected HTTP Response Status Code", fmt.Sprintf("Expected 204 but received %d", httpResp.StatusCode))
			return fmt.Errorf("unexpected status code: %d", httpResp.StatusCode)
		}

	}
	return nil
}

func mergeSerials(planSerials []string, serialsToAdd []string) []string {
	// Create a map to keep track of the existing serials in planSerials
	serialMap := make(map[string]bool)
	for _, serial := range planSerials {
		serialMap[serial] = true
	}

	// Add only the unique serials from serialsToAdd
	for _, serial := range serialsToAdd {
		if !serialMap[serial] {
			planSerials = append(planSerials, serial)
			serialMap[serial] = true
		}
	}

	return planSerials
}
