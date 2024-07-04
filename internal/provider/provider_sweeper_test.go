package provider

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	openApiClient "github.com/meraki/dashboard-api-go/client"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

// Meraki Dashboard API Calls //

func getMerakiOrganization(ctx context.Context, client *openApiClient.APIClient, organizationId string) (*openApiClient.GetOrganizations200ResponseInner, error) {
	inlineResp, _, err := client.OrganizationsApi.GetOrganization(ctx, organizationId).Execute()
	if err != nil {
		return nil, fmt.Errorf("error getting organization from Meraki API: %s", err)
	}
	return inlineResp, nil
}

func getMerakiOrganizations(ctx context.Context, client *openApiClient.APIClient) ([]openApiClient.GetOrganizations200ResponseInner, error) {
	inlineResp, _, err := client.OrganizationsApi.GetOrganizations(ctx).Execute()
	if err != nil {
		return nil, fmt.Errorf("error getting organizations list from Meraki API: %s", err)
	}
	return inlineResp, nil
}

func deleteMerakiOrganization(ctx context.Context, client *openApiClient.APIClient, organization openApiClient.GetOrganizations200ResponseInner) error {
	tflog.Info(ctx, "Deleting organization", map[string]interface{}{
		"name": *organization.Name,
		"id":   *organization.Id,
	})

	httpRespOrg, err := client.OrganizationsApi.DeleteOrganization(ctx, *organization.Id).Execute()
	if err != nil {
		var responseBody string
		if httpRespOrg.Body != nil {
			bodyBytes, readErr := io.ReadAll(httpRespOrg.Body)
			if readErr == nil {
				responseBody = string(bodyBytes)
			} else {
				responseBody = fmt.Sprintf("Failed to read response body: %s", readErr)
			}
			err := httpRespOrg.Body.Close()
			if err != nil {
				return err
			}
		} else {
			responseBody = "No response body"
		}
		tflog.Error(ctx, "Error deleting organization", map[string]interface{}{
			"name":         *organization.Name,
			"id":           *organization.Id,
			"error":        err,
			"responseBody": responseBody,
		})
		return err
	}

	if httpRespOrg.StatusCode == 204 {
		tflog.Info(ctx, "Successfully deleted organization", map[string]interface{}{
			"name": *organization.Name,
			"id":   *organization.Id,
		})
	}
	return nil
}

func getMerakiNetworks(ctx context.Context, client *openApiClient.APIClient, organizationId string, perPage int32) ([]openApiClient.GetNetwork200Response, error) {
	inlineResp, _, err := client.NetworksApi.GetOrganizationNetworks(ctx, organizationId).PerPage(perPage).Execute()
	if err != nil {
		return nil, fmt.Errorf("error getting network list from organization: %s, error: %s", organizationId, err)
	}
	return inlineResp, nil
}

func deleteMerakiNetwork(ctx context.Context, client *openApiClient.APIClient, network openApiClient.GetNetwork200Response) error {
	retries := 3
	wait := 1
	var deletedFromMerakiPortal bool

	tflog.Info(ctx, "Deleting network", map[string]interface{}{
		"name": *network.Name,
		"id":   *network.Id,
	})

	for retries > 0 {
		httpResp, err := client.NetworksApi.DeleteNetwork(ctx, *network.Id).Execute()
		if err != nil {
			tflog.Error(ctx, "Error deleting network", map[string]interface{}{
				"networkID": network.Id,
				"error":     err,
			})
		}

		if httpResp.StatusCode == 204 {
			tflog.Info(ctx, "Successfully deleted network", map[string]interface{}{
				"name": *network.Name,
				"id":   *network.Id,
			})
			deletedFromMerakiPortal = true
			break
		} else {
			retries -= 1
			time.Sleep(time.Duration(wait) * time.Second)
			wait += 1
		}
	}

	if !deletedFromMerakiPortal {
		tflog.Error(ctx, "Failed to delete network", map[string]interface{}{
			"name": *network.Name,
			"id":   *network.Id,
		})
		return fmt.Errorf("failed to delete network: %s, id: %s", *network.Name, *network.Id)
	}

	return nil
}

func getMerakiAdmins(ctx context.Context, client *openApiClient.APIClient, organizationId string) ([]openApiClient.GetOrganizationAdmins200ResponseInner, error) {
	admins, httpResp, err := client.AdminsApi.GetOrganizationAdmins(ctx, organizationId).Execute()
	if err != nil {
		if httpResp != nil {
			var responseBody string
			if httpResp.Body != nil {
				bodyBytes, readErr := io.ReadAll(httpResp.Body)
				if readErr == nil {
					responseBody = string(bodyBytes)
				} else {
					responseBody = fmt.Sprintf("Failed to read response body: %s", readErr)
				}
				err := httpResp.Body.Close()
				if err != nil {
					return nil, err
				}
			} else {
				responseBody = "No response body"
			}
			return nil, fmt.Errorf("error getting admin list from organization: %s, HTTP Status: %d, error: %s, Response Body: %s", organizationId, httpResp.StatusCode, err, responseBody)
		}
		return nil, fmt.Errorf("error getting admin list from organization: %s, error: %s", organizationId, err)
	}
	return admins, nil
}

func deleteMerakiAdmin(ctx context.Context, client *openApiClient.APIClient, organizationId string, admin openApiClient.GetOrganizationAdmins200ResponseInner) error {
	tflog.Info(ctx, "Deleting admin", map[string]interface{}{"email": *admin.Email, "id": *admin.Id})

	httpResp, err := client.AdminsApi.DeleteOrganizationAdmin(ctx, organizationId, *admin.Id).Execute()
	if err != nil {
		if httpResp != nil {
			var responseBody string
			if httpResp.Body != nil {
				bodyBytes, readErr := io.ReadAll(httpResp.Body)
				if readErr == nil {
					responseBody = string(bodyBytes)
				} else {
					responseBody = fmt.Sprintf("Failed to read response body: %s", readErr)
				}
				err := httpResp.Body.Close()
				if err != nil {
					return err
				}
			} else {
				responseBody = "No response body"
			}
			tflog.Error(ctx, "Error deleting admin", map[string]interface{}{
				"email":        *admin.Email,
				"id":           *admin.Id,
				"organization": organizationId,
				"status":       httpResp.StatusCode,
				"error":        err,
				"responseBody": responseBody,
			})
		} else {
			tflog.Error(ctx, "Error deleting admin", map[string]interface{}{
				"email":        *admin.Email,
				"id":           *admin.Id,
				"organization": organizationId,
				"error":        err,
			})
		}
		return err
	}

	if httpResp.StatusCode == http.StatusNoContent {
		tflog.Info(ctx, "Successfully deleted admin", map[string]interface{}{"email": *admin.Email, "id": *admin.Id})
	} else {
		var responseBody string
		if httpResp.Body != nil {
			bodyBytes, readErr := io.ReadAll(httpResp.Body)
			if readErr == nil {
				responseBody = string(bodyBytes)
			} else {
				responseBody = fmt.Sprintf("Failed to read response body: %s", readErr)
			}
			err := httpResp.Body.Close()
			if err != nil {
				return err
			}
		} else {
			responseBody = "No response body"
		}
		tflog.Error(ctx, "Failed to delete admin", map[string]interface{}{
			"email":        *admin.Email,
			"id":           *admin.Id,
			"status":       httpResp.StatusCode,
			"responseBody": responseBody,
		})
		return fmt.Errorf("failed to delete admin: %s, id: %s", *admin.Email, *admin.Id)
	}

	return nil
}

// Terraform Sweepers //

func sweepMerakiNetworks(ctx context.Context, client *openApiClient.APIClient, organizationId string) error {
	tflog.Info(ctx, "Starting network sweeper for organization", map[string]interface{}{"organization": organizationId})

	perPage := int32(100000)
	networks, err := getMerakiNetworks(ctx, client, organizationId, perPage)
	if err != nil {
		return err
	}

	for _, network := range networks {
		if strings.HasPrefix(*network.Name, "test_acc") {
			if err := deleteMerakiNetwork(ctx, client, network); err != nil {
				tflog.Error(ctx, "Failed to delete network", map[string]interface{}{
					"name":  *network.Name,
					"id":    *network.Id,
					"error": err,
				})
			}
		}
	}

	tflog.Info(ctx, "Finished running network sweeper")
	return nil
}

func sweepMerakiAdmins(ctx context.Context, client *openApiClient.APIClient, organizationId string) error {
	tflog.Info(ctx, "Starting admin sweeper for organization", map[string]interface{}{"organization": organizationId})

	admins, err := getMerakiAdmins(ctx, client, organizationId)
	if err != nil {
		return err
	}

	for _, admin := range admins {
		if admin.Email == nil || admin.Id == nil {
			tflog.Warn(ctx, "Skipping admin with missing email or ID", map[string]interface{}{"admin": admin})
			continue
		}

		if strings.HasPrefix(*admin.Email, "test_acc") {
			if err := deleteMerakiAdmin(ctx, client, organizationId, admin); err != nil {
				tflog.Error(ctx, "Failed to delete admin", map[string]interface{}{
					"email": *admin.Email,
					"id":    *admin.Id,
					"error": err,
				})
			}
		}
	}

	tflog.Info(ctx, "Finished running admin sweeper")
	return nil
}

func sweepMerakiOrganizations(ctx context.Context, client *openApiClient.APIClient) error {
	tflog.Info(ctx, "Starting organizations sweeper")

	organizations, err := getMerakiOrganizations(ctx, client)
	if err != nil {
		return err
	}

	for _, organization := range organizations {
		if strings.HasPrefix(*organization.Name, "test_acc") {

			// First, sweep networks and admins within the organization
			if err := sweepMerakiNetworks(ctx, client, *organization.Id); err != nil {
				tflog.Error(ctx, "Failed to sweep networks", map[string]interface{}{
					"organization": *organization.Name,
					"id":           *organization.Id,
					"error":        err,
				})
				continue
			}
			if err := sweepMerakiAdmins(ctx, client, *organization.Id); err != nil {
				tflog.Error(ctx, "Failed to sweep admins", map[string]interface{}{
					"organization": *organization.Name,
					"id":           *organization.Id,
					"error":        err,
				})
				continue
			}
			// Finally, delete the organization
			if err := deleteMerakiOrganization(ctx, client, organization); err != nil {
				tflog.Error(ctx, "Failed to delete organization", map[string]interface{}{
					"name":  *organization.Name,
					"id":    *organization.Id,
					"error": err,
				})
			}
		}
	}

	tflog.Info(ctx, "Finished running organizations sweeper")
	return nil
}

func sweepMerakiOrganization(ctx context.Context, client *openApiClient.APIClient, organizationId string) error {
	tflog.Info(ctx, "Starting organization sweeper")

	organization, err := getMerakiOrganization(ctx, client, organizationId)
	if err != nil {
		return err
	}

	if err := sweepMerakiNetworks(ctx, client, *organization.Id); err != nil {
		tflog.Error(ctx, "Failed to sweep networks", map[string]interface{}{
			"organization": *organization.Name,
			"id":           *organization.Id,
			"error":        err,
		})
	}
	if err := sweepMerakiAdmins(ctx, client, *organization.Id); err != nil {
		tflog.Error(ctx, "Failed to sweep admins", map[string]interface{}{
			"organization": *organization.Name,
			"id":           *organization.Id,
			"error":        err,
		})
	}

	tflog.Info(ctx, "Finished running organization sweeper")
	return nil
}

func SweeperHTTPClient() (*openApiClient.APIClient, error) {
	configuration := openApiClient.NewConfiguration()
	configuration.UserAgent = configuration.UserAgent + " terraform/dev-sweeper"
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
	}
	authenticatedTransport := &bearerAuthTransport{
		Transport: transport,
	}
	authenticatedTransport.Token = os.Getenv("MERAKI_DASHBOARD_API_KEY")
	retryClient := retryablehttp.NewClient()
	retryClient.HTTPClient.Transport = authenticatedTransport
	configuration.HTTPClient = retryClient.HTTPClient
	client := openApiClient.NewAPIClient(configuration)
	return client, nil
}

func init() {

	resource.AddTestSweepers("meraki_networks", &resource.Sweeper{
		Name: "meraki_networks",
		F: func(organizationId string) error {
			ctx := context.Background()
			client, err := SweeperHTTPClient()
			if err != nil {
				return err
			}
			return sweepMerakiNetworks(ctx, client, organizationId)
		},
	})

	resource.AddTestSweepers("meraki_admins", &resource.Sweeper{
		Name: "meraki_admins",
		F: func(organizationId string) error {
			ctx := context.Background()
			client, err := SweeperHTTPClient()
			if err != nil {
				return err
			}
			return sweepMerakiAdmins(ctx, client, organizationId)
		},
	})

	resource.AddTestSweepers("meraki_organization", &resource.Sweeper{
		Name: "meraki_organization",
		F: func(organizationId string) error {
			ctx := context.Background()
			client, err := SweeperHTTPClient()
			if err != nil {
				return err
			}
			return sweepMerakiOrganization(ctx, client, organizationId)
		},
	})

	resource.AddTestSweepers("meraki_organizations", &resource.Sweeper{
		Name: "meraki_organizations",
		F: func(organizationId string) error {
			ctx := context.Background()
			client, err := SweeperHTTPClient()
			if err != nil {
				return err
			}
			return sweepMerakiOrganizations(ctx, client)
		},
	})
}

func TestMain(m *testing.M) {
	exitCode := m.Run()

	ctx := context.Background()

	if exitCode != 0 {
		organizationId := os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")
		if organizationId == "" {
			tflog.Error(ctx, "TF_ACC_MERAKI_ORGANIZATION_ID must be set for sweeper to run")
			os.Exit(exitCode)
		}

		client, clientErr := SweeperHTTPClient()
		if clientErr != nil {
			tflog.Error(ctx, "Error getting HTTP client", map[string]interface{}{
				"error": clientErr,
			})
		}

		// Any Specified Organization Sweeper
		tflog.Info(ctx, "Running terraform sweepers due to test failures...")
		err := sweepMerakiOrganization(ctx, client, organizationId)
		if err != nil {
			tflog.Error(ctx, "Error running organization sweeper", map[string]interface{}{
				"error": err,
			})
		} else {
			tflog.Info(ctx, "Organization sweeper ran successfully")
		}

		// Targeted "test_acc" Organizations Sweeper
		err = sweepMerakiOrganizations(ctx, client)
		if err != nil {
			tflog.Error(ctx, "Error running organizations sweeper", map[string]interface{}{
				"error": err,
			})
		} else {
			tflog.Info(ctx, "Organizations sweeper ran successfully")
		}
	}

	os.Exit(exitCode)
}
