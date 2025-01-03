package provider

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	openApiClient "github.com/meraki/dashboard-api-go/client"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func SweeperHTTPClient() (*openApiClient.APIClient, error) {
	configuration := openApiClient.NewConfiguration()
	configuration.UserAgent = configuration.UserAgent + " terraform/dev-sweeper"
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
	}
	authenticatedTransport := &BearerAuthTransport{
		Transport: transport,
	}
	authenticatedTransport.Token = os.Getenv("MERAKI_DASHBOARD_API_KEY")
	retryClient := retryablehttp.NewClient()
	retryClient.HTTPClient.Transport = authenticatedTransport
	configuration.HTTPClient = retryClient.HTTPClient
	client := openApiClient.NewAPIClient(configuration)
	return client, nil
}

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
	fmt.Println("Deleting organization", map[string]interface{}{
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
		fmt.Println("Error deleting organization", map[string]interface{}{
			"name":         *organization.Name,
			"id":           *organization.Id,
			"error":        err,
			"responseBody": responseBody,
		})
		return err
	}

	if httpRespOrg.StatusCode == 204 {
		fmt.Println("Successfully deleted organization", map[string]interface{}{
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

	fmt.Println("Deleting network", map[string]interface{}{
		"name": *network.Name,
		"id":   *network.Id,
	})

	for retries > 0 {
		httpResp, err := client.NetworksApi.DeleteNetwork(ctx, *network.Id).Execute()
		if err != nil {
			fmt.Println("Error deleting network", map[string]interface{}{
				"networkID": network.Id,
				"error":     err,
			})
		}

		if httpResp.StatusCode == 204 {
			fmt.Println("Successfully deleted network", map[string]interface{}{
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
		fmt.Println("Failed to delete network", map[string]interface{}{
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
	fmt.Println("Deleting admin", map[string]interface{}{"email": *admin.Email, "id": *admin.Id})

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
			fmt.Println("Error deleting admin", map[string]interface{}{
				"email":        *admin.Email,
				"id":           *admin.Id,
				"organization": organizationId,
				"status":       httpResp.StatusCode,
				"error":        err,
				"responseBody": responseBody,
			})
		} else {
			fmt.Println("Error deleting admin", map[string]interface{}{
				"email":        *admin.Email,
				"id":           *admin.Id,
				"organization": organizationId,
				"error":        err,
			})
		}
		return err
	}

	if httpResp.StatusCode == http.StatusNoContent {
		fmt.Println("Successfully deleted admin", map[string]interface{}{"email": *admin.Email, "id": *admin.Id})
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
		fmt.Println("Failed to delete admin", map[string]interface{}{
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

func SweepMerakiNetworks(ctx context.Context, client *openApiClient.APIClient, organizationId string) error {
	fmt.Println("Starting network sweeper for organization", map[string]interface{}{"organization": organizationId})

	perPage := int32(100000)
	networks, err := getMerakiNetworks(ctx, client, organizationId, perPage)
	if err != nil {
		return err
	}

	for _, network := range networks {
		if strings.HasPrefix(*network.Name, "test_acc") {
			if err := deleteMerakiNetwork(ctx, client, network); err != nil {
				fmt.Println("Failed to delete network", map[string]interface{}{
					"name":  *network.Name,
					"id":    *network.Id,
					"error": err,
				})
			}
		}
	}

	fmt.Println("Finished running network sweeper")
	return nil
}

func SweepMerakiAdmins(ctx context.Context, client *openApiClient.APIClient, organizationId string) error {
	fmt.Println("Starting admin sweeper for organization", map[string]interface{}{"organization": organizationId})

	admins, err := getMerakiAdmins(ctx, client, organizationId)
	if err != nil {
		return err
	}

	for _, admin := range admins {
		if admin.Email == nil || admin.Id == nil {
			fmt.Println("Skipping admin with missing email or ID", map[string]interface{}{"admin": admin})
			continue
		}

		if strings.HasPrefix(*admin.Email, "test_acc") {
			if err := deleteMerakiAdmin(ctx, client, organizationId, admin); err != nil {
				fmt.Println("Failed to delete admin", map[string]interface{}{
					"email": *admin.Email,
					"id":    *admin.Id,
					"error": err,
				})
			}
		}
	}

	fmt.Println("Finished running admin sweeper")
	return nil
}

func SweepMerakiOrganizations(ctx context.Context, client *openApiClient.APIClient) error {
	fmt.Println("Starting organizations sweeper")

	organizations, err := getMerakiOrganizations(ctx, client)
	if err != nil {
		return err
	}

	for _, organization := range organizations {
		if strings.HasPrefix(*organization.Name, "test_acc") {

			// First, sweep networks and admins within the organization
			if err := SweepMerakiNetworks(ctx, client, *organization.Id); err != nil {
				fmt.Println("Failed to sweep networks", map[string]interface{}{
					"organization": *organization.Name,
					"id":           *organization.Id,
					"error":        err,
				})
				continue
			}
			if err := SweepMerakiAdmins(ctx, client, *organization.Id); err != nil {
				fmt.Println("Failed to sweep admins", map[string]interface{}{
					"organization": *organization.Name,
					"id":           *organization.Id,
					"error":        err,
				})
				continue
			}
			// Finally, delete the organization
			if err := deleteMerakiOrganization(ctx, client, organization); err != nil {
				fmt.Println("Failed to delete organization", map[string]interface{}{
					"name":  *organization.Name,
					"id":    *organization.Id,
					"error": err,
				})
			}
		}
	}

	fmt.Println("Finished running organizations sweeper")
	return nil
}

func SweepMerakiOrganization(ctx context.Context, client *openApiClient.APIClient, organizationId string) error {
	fmt.Println("Starting organization sweeper")

	organization, err := getMerakiOrganization(ctx, client, organizationId)
	if err != nil {
		return err
	}

	if err := SweepMerakiNetworks(ctx, client, *organization.Id); err != nil {
		fmt.Println("Failed to sweep networks", map[string]interface{}{
			"organization": *organization.Name,
			"id":           *organization.Id,
			"error":        err,
		})
	}
	if err := SweepMerakiAdmins(ctx, client, *organization.Id); err != nil {
		fmt.Println("Failed to sweep admins", map[string]interface{}{
			"organization": *organization.Name,
			"id":           *organization.Id,
			"error":        err,
		})
	}

	fmt.Println("Finished running organization sweeper")
	return nil
}

// Sweeper Definitions //

func init() {
	ctx := context.Background()

	resource.AddTestSweepers("meraki_networks", &resource.Sweeper{
		Name: "meraki_networks",
		F: func(organizationId string) error {
			fmt.Println("Running meraki_networks sweeper")
			client, err := SweeperHTTPClient()
			if err != nil {
				fmt.Println("Error creating HTTP client", map[string]interface{}{"error": err})
				return err
			}
			return SweepMerakiNetworks(ctx, client, organizationId)
		},
	})

	resource.AddTestSweepers("meraki_admins", &resource.Sweeper{
		Name: "meraki_admins",
		F: func(organizationId string) error {
			fmt.Println("Running meraki_admins sweeper")
			client, err := SweeperHTTPClient()
			if err != nil {
				fmt.Println("Error creating HTTP client", map[string]interface{}{"error": err})
				return err
			}
			return SweepMerakiAdmins(ctx, client, organizationId)
		},
	})

	resource.AddTestSweepers("meraki_organization", &resource.Sweeper{
		Name: "meraki_organization",
		F: func(organizationId string) error {
			fmt.Println("Running meraki_organization sweeper")
			client, err := SweeperHTTPClient()
			if err != nil {
				fmt.Println("Error creating HTTP client", map[string]interface{}{"error": err})
				return err
			}
			return SweepMerakiOrganization(ctx, client, organizationId)
		},
	})

	resource.AddTestSweepers("meraki_organizations", &resource.Sweeper{
		Name: "meraki_organizations",
		F: func(organizationId string) error {
			fmt.Println("Running meraki_organizations sweeper")
			client, err := SweeperHTTPClient()
			if err != nil {
				fmt.Println("Error creating HTTP client", map[string]interface{}{"error": err})
				return err
			}
			return SweepMerakiOrganizations(ctx, client)
		},
	})
}

// runSweepers function can be invoked independently or after tests
func runSweepers() {
	ctx := context.Background()
	fmt.Println("Running Sweepers...")

	organizationId := os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")
	if organizationId == "" {
		fmt.Println("TF_ACC_MERAKI_ORGANIZATION_ID must be set for sweeper to run")
		os.Exit(1)
	}

	client, clientErr := SweeperHTTPClient()
	if clientErr != nil {
		fmt.Println("Error getting HTTP client", map[string]interface{}{
			"error": clientErr,
		})
		os.Exit(1)
	}

	// Sweep a Specific Static Organization
	err := SweepMerakiOrganization(ctx, client, organizationId)
	if err != nil {
		fmt.Println("Error running organization sweeper", map[string]interface{}{
			"error": err,
		})
	} else {
		fmt.Println("Organization sweeper ran successfully")
	}

	// Targeted "test_acc" Organizations Sweeper
	err = SweepMerakiOrganizations(ctx, client)
	if err != nil {
		fmt.Println("Error running organizations sweeper", map[string]interface{}{
			"error": err,
		})
	} else {
		fmt.Println("Organizations sweeper ran successfully")
	}
}
