# Local Development Environment Setup

Setting up a local development environment is a crucial step for contributing to the Meraki Terraform Provider project. Below you'll find detailed instructions on how to prepare your development environment.

## Requirements

- [Go](https://golang.org/doc/install) version 1.20 or higher
- [Terraform](https://www.terraform.io/downloads.html) version 1.0 or higher

## Go Environment Setup

To contribute to the provider's codebase, you need to have Go installed on your system.

1. **Install Go**:

   Use Homebrew to install Go on macOS:

   ```shell
   brew install go
   ```

2. **Configure Go Environment**:

   Set up your Go environment by configuring the `GOBIN` path:

   ```shell
   echo $GOBIN
   export GOBIN=$(go env GOPATH)/bin
   ```

   > **Apple M1/ARM-based processor note**: Add an extra environment variable to mitigate potential issues with asynchronous preemption:

   ```shell
   export GODEBUG=asyncpreemptoff=1
   ```

## Terraform Development Environment Setup

After setting up Go, you need to install Terraform to work on the provider.

1. **Install Terraform**:

   ```shell
   brew install terraform
   ```

   > **Apple M1/ARM-based processor note**: If necessary, upgrade Terraform specifically for ARM architecture:

   ```shell
   arch -arm64 brew upgrade terraform
   ```

2. **Clone the Provider Repository**:

   Clone the official repository to start working on the provider:

   ```shell
   git clone https://github.com/core-infra-svcs/terraform-provider-meraki
   cd terraform-provider-meraki
   ```

3. **Terraform Development Overrides**:

   Set up Terraform development overrides by creating a `.terraformrc` file:

   ```shell
   touch $HOME/.terraformrc
   vim $HOME/.terraformrc
   ```

   Insert the following content:

   ```hcl
   provider_installation {
     dev_overrides {
       "core-infra-svcs/meraki" = "$(go env GOPATH)/bin"
     }
     direct {}
   }
   ```

4. **Download Required Binaries**:

   Fetch the dependencies:

   ```shell
   go mod download
   ```

## Compiling the Provider

Compile the provider using the Go install command:

```shell
go install
```

This will place the provider binary in your `$GOPATH/bin` directory.

To update the documentation, run:

```shell
go generate
```

For acceptance tests (which may incur costs due to resource creation), run:

```shell
make testacc
```

## Adding Dependencies

This provider is managed with Go modules.

To add a new dependency, for example, `github.com/author/dependency`:

```shell
go get github.com/author/dependency
go mod tidy
```

Commit the changes to `go.mod` and `go.sum` after modification.

## Meraki Development Environment

To work on the Meraki provider, set up your environment:

1. **Create Terraform Configuration Files**:

   For instance, `meraki_vars.auto.tfvars` and `terraform.tfvars` to handle environment variables.

   ```hcl
   # meraki_vars.auto.tfvars
   variable "MERAKI_DASHBOARD_API_KEY" {
     type      = string
     sensitive = true
   }
   
   variable "host" {
     type    = string
     default = "api.meraki.com"
   }
   ```

   ```hcl
   # terraform.tfvars
   MERAKI_DASHBOARD_API_KEY = var.MERAKI_DASHBOARD_API_KEY
   ```

2. **Export Environment Variable**:

   Set the Meraki dashboard API key in your shell:

   ```shell
   export MERAKI_DASHBOARD_API_KEY="your_api_key_here"
   ```

3. **Run Terraform**:

   With the environment variable set, you can now run Terraform commands such as `terraform init`, `terraform plan`, and `terraform apply`.

Ensure that you manage your API key securely and do not expose it in your code or version control.

With your environment configured, you're ready to contribute to the Meraki Terraform Provider. Happy coding!