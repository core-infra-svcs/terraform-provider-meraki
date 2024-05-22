package provider

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	openApiClient "github.com/meraki/dashboard-api-go/client"
	"net/http"
	"os"
	"strings"
	"testing"
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

func TestMain(m *testing.M) {
	exitCode := m.Run()

	if exitCode != 0 {
		organizationID := os.Getenv("TF_ACC_MERAKI_ORGANIZATION_ID")
		if organizationID == "" {
			fmt.Println("TF_ACC_MERAKI_ORGANIZATION_ID must be set for sweeper to run")
			os.Exit(exitCode)
		}

		fmt.Println("Running sweeper due to test failures...")
		err := sweepMerakiNetwork(organizationID)
		if err != nil {
			fmt.Printf("Error running network sweeper: %s\n", err)
		} else {
			fmt.Println("Network sweeper ran successfully")
		}

		err = sweepMerakiOrganization(organizationID)
		if err != nil {
			fmt.Printf("Error running organization sweeper: %s\n", err)
		} else {
			fmt.Println("Organization sweeper ran successfully")
		}
	}

	os.Exit(exitCode)
}

func sweepMerakiNetwork(organization string) error {
	fmt.Printf("Starting network sweeper for organization: %s\n", organization)

	client, err := SweeperHTTPClient()
	if err != nil {
		return fmt.Errorf("error getting http client: %s", err)
	}

	retries := 3
	wait := 1
	var deletedFromMerakiPortal bool

	perPage := int32(100000)
	inlineResp, _, err := client.NetworksApi.GetOrganizationNetworks(context.Background(), organization).PerPage(perPage).Execute()
	if err != nil {
		return fmt.Errorf("error getting network list from organization:%s \nerror: %s", organization, err)
	}

	for _, merakiNetwork := range inlineResp {
		if strings.HasPrefix(*merakiNetwork.Name, "test_acc") {
			fmt.Printf("Deleting network: %s, id: %s\n", *merakiNetwork.Name, *merakiNetwork.Id)

			for retries > 0 {
				httpResp, err2 := client.NetworksApi.DeleteNetwork(context.Background(), *merakiNetwork.Id).Execute()
				if err2 != nil {
					fmt.Printf("Error deleting network from organization:%s \nerror: %s\n", organization, err2)
				}

				if httpResp.StatusCode == 204 {
					fmt.Printf("Successfully deleted network: %s, id: %s\n", *merakiNetwork.Name, *merakiNetwork.Id)
					deletedFromMerakiPortal = true
					break
				} else {
					retries -= 1
					time.Sleep(time.Duration(wait) * time.Second)
					wait += 1
				}
			}

			if !deletedFromMerakiPortal {
				fmt.Printf("Failed to delete network: %s, id: %s\n", *merakiNetwork.Name, *merakiNetwork.Id)
			}
		}
	}
	fmt.Println("Finished running network sweeper")
	return nil
}

func sweepMerakiOrganization(organization string) error {
	fmt.Printf("Starting organization sweeper for organization: %s\n", organization)

	client, err := SweeperHTTPClient()
	if err != nil {
		return fmt.Errorf("error getting http client: %s", err)
	}

	inlineResp, _, err := client.OrganizationsApi.GetOrganizations(context.Background()).Execute()
	if err != nil {
		return fmt.Errorf("error getting organizations list from Meraki API: %s\n", err)
	}

	for _, merakiOrganization := range inlineResp {
		if strings.HasPrefix(*merakiOrganization.Name, "test") {
			fmt.Printf("Deleting organization: %s, id: %s\n", *merakiOrganization.Name, *merakiOrganization.Id)

			perPage := int32(100000)
			inlineRespNetwork, _, err1 := client.NetworksApi.GetOrganizationNetworks(context.Background(), *merakiOrganization.Id).PerPage(perPage).Execute()
			if err1 != nil {
				return fmt.Errorf("error getting network list from organization:%s error: %s\n", *merakiOrganization.Id, err1)
			}

			for _, merakiNetwork := range inlineRespNetwork {
				fmt.Printf("Deleting network: %s, id: %s\n", *merakiNetwork.Name, *merakiNetwork.Id)

				networkHttpResp, err2 := client.NetworksApi.DeleteNetwork(context.Background(), *merakiNetwork.Id).Execute()
				if err2 != nil {
					fmt.Printf("%v\n", networkHttpResp)
				}

				if networkHttpResp.StatusCode == 204 {
					fmt.Printf("Successfully deleted network: %s, id: %s\n", *merakiNetwork.Name, *merakiNetwork.Id)
				}
			}

			httpRespOrg, err3 := client.OrganizationsApi.DeleteOrganization(context.Background(), *merakiOrganization.Id).Execute()
			if err3 != nil {
				fmt.Printf("%v\n", httpRespOrg.Body)
			}
			if httpRespOrg.StatusCode == 204 {
				fmt.Printf("Successfully deleted organization: %s, id: %s\n", *merakiOrganization.Name, *merakiOrganization.Id)
			}
		}
	}

	fmt.Println("Finished running organization sweeper")
	return nil
}
