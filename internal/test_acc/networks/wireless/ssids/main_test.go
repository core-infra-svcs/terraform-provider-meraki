package ssids

import (
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/test_acc"
	"testing"
)

func TestMain(m *testing.M) {
	// Call the shared EntryPoint function to ensure setup/cleanup
	test_acc.EntryPoint(m)
}
