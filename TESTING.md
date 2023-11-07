# TESTING.md

This document provides guidelines on how to set up an integration testing environment for the Meraki Terraform Provider. Integration tests are essential to ensure that the provider interacts correctly with the Meraki Dashboard API and hardware.

## Requirements

Before you begin, you will need:

- **Meraki Hardware**: Access to physical Meraki devices including MX (Security & SD-WAN), MS (Switches), and MR (Wireless Access Points).
- **Licenses**: Valid licenses for the Meraki hardware to operate.
- **Meraki Organization**: A dedicated Meraki organization with API access enabled is required for testing purposes.
- **API Access**: The Meraki Dashboard API must be enabled within your organization. You can enable this feature on the Meraki dashboard under Organization > Settings.

## Environment Variables

The Meraki Terraform Provider's integration tests rely on several environment variables. These variables must be set before running the tests:

- `MERAKI_DASHBOARD_API_KEY`: Your Meraki Dashboard API Key.
- `TF_ACC=1`: This enables Terraform acceptance tests.
- `TF_ACC_MERAKI_MG_SERIAL`: The serial number of your Meraki MG device.
- `TF_ACC_MERAKI_MR_SERIAL`: The serial number of your Meraki MR device.
- `TF_ACC_MERAKI_MS_SERIAL`: The serial number of your Meraki MS device.
- `TF_ACC_MERAKI_MX_SERIAL`: The serial number of your Meraki MX device.
- `TF_ACC_MERAKI_MX_LICENSE`: The license key for your Meraki MX device.
- `TF_ACC_MERAKI_ORDER_NUMBER`: The order number associated with your Meraki devices.
- `TF_ACC_MERAKI_ORGANIZATION_ID`: The ID of your dedicated Meraki organization.

## Setting Up Your Testing Environment

To set up your integration testing environment, follow these steps:

1. **Set Environment Variables**:

   Export the required environment variables in your shell session. Replace the placeholders with actual values from your Meraki setup.

   ```shell
   export TF_ACC_MERAKI_DASHBOARD_API_KEY='your_dashboard_api_key'
   export TF_ACC=1
   export TF_ACC_MERAKI_MG_SERIAL='your_meraki_mg_serial_number'
   export TF_ACC_MERAKI_MR_SERIAL='your_meraki_mr_serial_number'
   export TF_ACC_MERAKI_MS_SERIAL='your_meraki_ms_serial_number'
   export TF_ACC_MERAKI_MX_SERIAL='your_meraki_mx_serial_number'
   export TF_ACC_MERAKI_MX_LICENSE='your_meraki_mx_license_key'
   export TF_ACC_MERAKI_ORDER_NUMBER='your_meraki_order_number'
   export TF_ACC_MERAKI_ORGANIZATION_ID='your_meraki_organization_id'
   ```

2. **Write Acceptance Tests**:

   Create a test file within the `./meraki/provider` directory that will use these environment variables to configure the provider and run the tests.

3. **Running the Tests**:

   From the root of the repository, execute the acceptance tests using the `make` command:

   ```shell
   # test all resources
   make testacc
   
   # test a single resource
   make test NAME={TEST_RESOURCE_NAME} 
   
   # clean up test networks left in the test organization
   make sweep
   
   ```

   The `make` command will trigger the `go test` command along with any necessary flags and arguments to run the acceptance tests.

## Best Practices for Integration Testing

- **Security**: Treat your API keys and other sensitive data with care. Ensure they are not hard-coded in your tests or committed to version control.
- **Cost Management**: Remember that running acceptance tests on real devices could incur costs and changes to your environment. Always review the tests to understand the actions they perform.
- **Clean Up**: Write tests that clean up resources after they run to prevent unnecessary charges and clutter in your Meraki environment.
- **Monitoring**: Keep an eye on the Meraki dashboard while tests are running. This will help you understand the changes being made and troubleshoot any issues that arise.

## Troubleshooting

If you encounter any issues during testing, consider the following steps:

- Verify that all required environment variables are correctly set.
- Check your Meraki dashboard to ensure that the API access is functioning as expected.
- Review the Terraform provider logs for any error messages or indications of what may be going wrong.

With your testing environment properly set up, you can proceed to run integration tests on your Meraki Terraform Provider codebase confidently.