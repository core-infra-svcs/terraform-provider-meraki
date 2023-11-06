# Local Development Environment Setup

Follow the instructions below to set up your local development environment for the Meraki Terraform Provider.
## Requirements

- [Go](https://golang.org/doc/install) >= 1.18
- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
### Go Environment

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine.

1. Install Go by running the following command:

   ```shell
   brew install go
   ```

2. Set the `GOBIN` environment variable by executing the following commands:

   ```shell
   echo $GOBIN
   export GOBIN=$GOPATH/bin
   ```

   > Note: If you are using an Apple M1/ARM-based processor, add the following environment variable as well:

   ```shell
   export GODEBUG=asyncpreemptoff=1
   ```

### Terraform Development Environment

1. Install Terraform by running the following command:

   ```shell
   brew install terraform
   ```

   > Note: If you are using an Apple M1/ARM-based processor, upgrade Terraform using the following command:

   ```shell
   arch -arm64 brew upgrade terraform
   ```

2. Clone the repository: [https://github.com/core-infra-svcs/terraform-provider-meraki](https://github.com/core-infra-svcs/terraform-provider-meraki)

3. Set up the Terraform development overrides by creating and editing the `.terraformrc` file:

   ```shell
   touch $HOME/.terraformrc
   vim $HOME/.terraformrc
   ```

   Add the following configuration to the file:

   ```text
   provider_installation {
       dev_overrides {
           "$REPOSITORY/$PROVIDER" = "/Users/$USER/go/bin"
       }
   }
   ```

3. Download the required binaries by executing the following commands:

   ```shell
   go mod download
   ```

## Compiling the Provider

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up to date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

### Meraki Development Environment

To set up your Meraki development environment for Terraform, follow these steps:

1. Create the following files to simplify running Terraform:

   **`meraki_vars.auto.tf`**
   ```hcl
   variable "MERAKI_DASHBOARD_API_KEY" {
       type      = string
       sensitive = true
   }

   variable "host" {
       type     = string
       default  = "api.meraki.com"
   }
   ```

   **`terraform.tfvars`**
   ```hcl
   MERAKI_DASHBOARD_API_KEY = var.MERAKI_DASHBOARD_API_KEY
   ```

   In the `meraki_vars.auto.tf` file, we define the necessary variables for working with the Meraki provider. The `MERAKI_DASHBOARD_API_KEY` variable is marked as sensitive to protect the sensitive password.

   In the `terraform.tfvars` file, we reference the `MERAKI_DASHBOARD_API_KEY` variable using the `var.MERAKI_DASHBOARD_API_KEY` syntax. This allows us to pull the value from the environment variable.

2. Set the environment variable:

   Before running Terraform commands, make sure to set the `MERAKI_DASHBOARD_API_KEY` environment variable with your sensitive password. You can do this in your shell or by using a script. For example, in Unix-based systems, you can run the following command:

   ```shell
   export MERAKI_DASHBOARD_API_KEY="01234567890123456789"
   ```

   Replace `"01234567890123456789"` with your actual Meraki dashboard API key.

   This command sets the `MERAKI_DASHBOARD_API_KEY` environment variable to the desired sensitive password.

3. Run Terraform commands:

   With the environment variable set and the necessary files in place, you can now run Terraform commands as usual. Terraform will automatically retrieve the sensitive password from the environment variable and use it during execution.

   For example, you can run `terraform init` to initialize the working directory and then proceed with other Terraform commands like `terraform plan` and `terraform apply`.

   Remember to keep the environment variable containing the sensitive password secure and avoid printing or exposing it inadvertently.

That's it! You have successfully set up your local development environment for the Meraki Terraform Provider. Now you can start contributing to the project. Happy coding!