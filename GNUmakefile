default: testacc

# Run acceptance tests
.PHONY: testacc
testacc:
	TF_ACC=1 go test ./... -v $(TESTARGS) -timeout 120m

# Test a single resource. Usage: make test NAME=TestAccOrganizationsNetworkResource
.PHONY: test
test:
ifndef NAME
	$(error Resource/DataSource test NAME is not set. Usage: make test NAME=TestAccOrganizationsNetworkResource)
endif
	TF_ACC=1 go test ./... -v -run $(NAME) -timeout 120m

# Clean up resources with a Terraform sweeper
.PHONY: sweep
sweep:
	TF_ACC=1 go test ./... -v -sweep=$(TF_ACC_MERAKI_ORGANIZATION_ID) -sweep-run='meraki_network'
