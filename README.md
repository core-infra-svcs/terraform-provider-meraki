# Terraform Provider Meraki

The Meraki Terraform Provider is a declarative tool that enables teams and individuals to automate their workflows and manage Cisco Meraki network infrastructure using Terraform. With this provider, you can define and manage Meraki organizations, networks, devices, and other resources as code, providing simplicity, scalability, and repeatability in your automation strategy.

Additional provider information available on the [Terraform Registry](https://registry.terraform.io/providers/core-infra-svcs/meraki/latest).

## Features

- Provision and manage Meraki organizations, networks, devices, and more through infrastructure as code.
- Configure and control various aspects of your Meraki infrastructure, including network settings, security policies, and device configurations.
- Leverage the power of Terraform to plan, apply, and manage changes to your Meraki environment in a controlled and auditable manner.
- Enable collaboration and version control for your Meraki configurations, allowing teams to work together efficiently and track changes over time.


## Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0

## Getting Started

To start using the Meraki Terraform Provider, follow these steps:

1. Install Terraform: Make sure you have Terraform installed on your local machine. You can download and install Terraform from the [official website](https://www.terraform.io/downloads.html).

2. Configure the Meraki Terraform Provider: Set up the provider block in your Terraform configuration file, specifying your Meraki API key and base URL.

3. Define Meraki resources: Create resource blocks in your configuration file to describe the desired state of your Meraki infrastructure. You can create organizations, networks, devices, and more.

4. Initialize and apply changes: Run `terraform init` to initialize your Terraform configuration. Then, use `terraform plan` to preview the changes that will be applied, and `terraform apply` to apply the changes and provision the resources.

## Usage Documentation

To use the Meraki Terraform Provider in your Terraform configuration, you need to configure the required provider and define resources specific to the Meraki platform. Follow the steps below to get started:

### Configuration

1. Add the provider block to your Terraform configuration file (e.g., `main.tf`):

   ```hcl
   terraform {
     required_providers {
       meraki = {
         source = "core-infra-svcs/meraki"
       }
     }
   }

   provider "meraki" {
     api_key  = var.MERAKI_DASHBOARD_API_KEY
     base_url = var.MERAKI_DASHBOARD_API_URL
   }
   ```

   Replace `var.MERAKI_DASHBOARD_API_KEY` and `var.MERAKI_DASHBOARD_API_URL` with your own API key and base URL values.

2. Define Meraki resources in your configuration file. For example, you can create a new Meraki organization and network:

   ```hcl
   // Create new Meraki organization
   resource "meraki_organization" "demo" {
     name            = "example"
     api_enabled     = true
     licensing_model = "co-term"
   }

   // Create a new Network
   resource "meraki_network" "demo" {
     depends_on      = [meraki_organization.demo]
     organization_id = resource.meraki_organization.demo.organization_id
     product_types   = ["appliance", "switch", "wireless"]
     tags            = ["cisco", "meraki", "terraform"]
     name            = "Main Office"
     timezone        = "America/Los_Angeles"
     notes           = "example demo network"
   }
   ```

   Customize the configuration according to your requirements, such as providing the appropriate values for the organization name, product types, tags, etc.

### Commands

1. Initialize the Terraform configuration by running the following command in your project directory:

   ```shell
   terraform init
   ```

2. Plan and preview the changes that will be applied to your Meraki environment:

   ```shell
   terraform plan
   ```

3. Apply the changes to create the Meraki organization and network:

   ```shell
   terraform apply
   ```

   Review the changes and confirm by typing `yes` when prompted.

4. Monitor the Terraform output for any errors or warnings. Once the process completes successfully, the Meraki organization and network will be created according to your configuration.

Now you have successfully used the Meraki Terraform Provider to provision resources in your Meraki environment. You can further customize your configuration to manage additional Meraki resources or update existing ones by modifying the Terraform configuration file.

Remember to always review and verify the changes before applying them to your production environment.

For even more detailed information and usage examples, please refer to the following documentation resources:

- [Meraki Provider Documentation](./docs): Explore the available resources and data sources provided by the Meraki Terraform Provider.


## Contributing

Contributions are welcome! If you are interested in contributing to the Meraki Terraform Provider, please refer to the [Contributing Guidelines](./CONTRIBUTING.md) for detailed instructions on how to get started.

See the [Getting Started Document](.github/workflow-docs/getting-started.md) for detailed instructions.

## License

The Meraki Terraform Provider is open-source and licensed under the [Mozilla Public License Version 2.0](./LICENSE). Feel free to use, modify, and distribute the provider according to the terms of the license.

## Support

If you encounter any issues, have questions, or need assistance, please [create an issue](https://github.com/core-infra-svcs/terraform-provider-meraki/issues) on the GitHub repository. Our community and maintainers will be happy to help you.

## Acknowledgements

We would like to express our gratitude to the contributors who have made this project possible. Your contributions and feedback are highly appreciated and valuable to the Meraki Terraform Provider community.

## Disclaimer

This project is not officially supported by Cisco or Meraki. It is maintained and supported by a community of enthusiastic engineers and developers.
