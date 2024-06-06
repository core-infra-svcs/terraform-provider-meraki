# Meraki Terraform Provider Development Workflow

This repository contains the development workflow for the open-source Meraki Terraform Provider. The workflow is designed to ensure code quality, facilitate peer review, and enforce the approval process for all contributors.

In order to fully leverage this automated workflow it is essential to use the [Conventional Commit](https://www.conventionalcommits.org/en/v1.0.0/) specification for all commit messages.

## Workflow Overview

Our development workflow consists of the following GitHub Actions workflows:

- `resource-development-branch.yml`: This workflow runs linting checks on every branch change to ensure code quality, and tests for a single resource, and is triggered on push events for all branches except the main branch.
- `release-please-action.yml`: This workflow creates a release pull request and git tags for the Terraform provider.
- `terraform-provider-release.yml`: This workflow publishes assets for release when a tag is created.

## Workflow Details

### `integration-test-branch.yml`

This workflow is responsible for running linting checks on the codebase.
This workflow is triggered on push events for all branches except the main branch. It runs tests for a single resource. The workflow consists of the following steps:

1. Checks out the repository.
2. Sets up the Go environment.
3. Runs the `golangci-lint` action to perform linting checks on the codebase.
4. Builds the project by running `go build` and ensures the project builds successfully.
5. Generates code by running `go generate` and checks for any unexpected differences in directories after code generation.
6. Runs acceptance tests for your single resource or data source using Terraform CLI. It sets the `TF_ACC` environment variable to enable acceptance tests and provides the `MERAKI_DASHBOARD_API_KEY` secret as an environment variable.
7. Checks the test results and fails the workflow if the tests fail.

### `release-please-action.yml`

This workflow is responsible for creating a release pull request and git tags for the Terraform provider. It is triggered on push events to the main branch. The workflow includes the following steps:

1. Checks out the repository.
2. Runs the `release-please-action` to create a release pull request and git tags. It uses the `release-type` and `package-name` parameters to specify the type of release and the package name.
3. Tags major and minor versions of the release and pushes the tags to the repository as defined in [release-please-action](https://github.com/marketplace/actions/release-please-action#release-types-supported).

### `terraform-provider-release.yml`

This workflow publishes assets for release when a tag is created with the pattern "v*". It includes the following steps:

1. Checks out the repository.
2. Sets up the Go environment.
3. Runs acceptance tests for all resources and data sources using go test. It sets the `TF_ACC` environment variable to enable acceptance tests and provides the `MERAKI_DASHBOARD_API_KEY` secret as an environment variable.
4. Checks the test results and fails the workflow if the tests fail.
5. Imports the GPG key using the `crazy-max/ghaction-import-gpg` action.
6. Runs GoReleaser using the `goreleaser/goreleaser-action` to build and package the release assets.
7. Uses the `GPG_PRIVATE_KEY` and `PASSPHRASE` secrets for GPG key handling.
8. Uses the `GITHUB_TOKEN` secret for authentication and authorization.

## Branch Protections

- Branch protections are set up in the repository to enforce the peer review and approval process.

- Pull requests must meet the defined criteria, including an approving review and passing all required CI tests and checks.

- This ensures that external contributors' changes go through a rigorous review and testing process before being merged into the main branch.


## Peer Review and Approval Process

- Each pull request submitted by a contributor goes through a peer review process.

- Peer reviewers examine the code changes, provide feedback, and approve the pull request if the changes meet the project's standards.

- Contributors can create a pull request from their development branch to the main branch only after all local test cases have passed.

- The pull request triggers the `integration-tests-branch-pr.yml` workflow, running tests for all resources included in the pull request.

- The pull request cannot be merged until it passes all the tests and receives the required approvals.

- Once the pull request is merged into main the `release-please-action.yml` workflow will create a new version major/minor and a release PR.



- The final approval comes after the release PR is merged into main with the latest git tag version which starts the `terraform-provider-release.yml` workflow of generating the binaries which will be published to Terraform registry.

