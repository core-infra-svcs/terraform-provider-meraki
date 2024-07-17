package test_acc

import "testing"

func TestMain(m *testing.M) {
	// Call the shared EntryPoint function to ensure setup/cleanup
	EntryPoint(m)
}
