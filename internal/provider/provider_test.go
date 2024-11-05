package provider

import (
	"fmt"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// Check if the environment variable is set to run sweepers independently
	if os.Getenv("RUN_SWEEPERS") == "1" {
		runSweepers()
		return
	}

	// Run tests
	exitCode := m.Run()

	// If tests failed, run sweepers to clean up
	if exitCode != 0 {
		fmt.Println("Tests failed, running sweepers for cleanup")
		runSweepers()
	}

	os.Exit(exitCode)
}
