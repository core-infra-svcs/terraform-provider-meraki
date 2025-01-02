package provider

import (
	"context"
	"fmt"
	"github.com/core-infra-svcs/terraform-provider-meraki/internal/utils"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

/*
func TestProviderEncryptionFeature(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { provider.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: TestAccProtoV6ProviderFactories,
		CheckDestroy:             testAccCheckProviderDestroy,
		Steps: []resource.TestStep{

			// Provider with encryption key
			{
				Config: testProviderConfigWithEncryptionKey("my_secret_encryption_key"),
				Check: resource.ComposeTestCheckFunc(
					testCheckEncryptionDecryptionWithKey("my_secret_encryption_key"),
				),
			},

			// Provider without encryption key
			{
				Config: testProviderConfigWithoutEncryptionKey(),
				Check: resource.ComposeTestCheckFunc(
					testCheckEncryptionDecryptionWithoutKey(),
				),
			},
		},
	})
}

*/

func testAccCheckProviderDestroy(s *terraform.State) error {
	return nil
}

func testProviderConfigWithEncryptionKey(encryptionKey string) string {
	return fmt.Sprintf(`
provider "meraki" {
  encryption_key = "%s"
}
`, encryptionKey)
}

func testProviderConfigWithoutEncryptionKey() string {
	return `
provider "meraki" {
  // No encryption key
}
`
}

// Define a custom type for context keys
type contextKey string

const encryptionKeyContextKey contextKey = "encryption_key"

func testCheckEncryptionDecryptionWithKey(encryptionKey string) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		// Set a value in the context using the custom key type
		ctx := context.WithValue(context.Background(), encryptionKeyContextKey, encryptionKey)

		// Retrieve the value from the context
		if v, ok := ctx.Value(encryptionKeyContextKey).(string); ok {
			fmt.Println("Encryption Key from context:", v)
		} else {
			fmt.Println("Encryption Key not found in context")
		}

		// Test encryption
		encrypted, err := utils.Encrypt(encryptionKey, "supersecret")
		if err != nil {
			return fmt.Errorf("error encrypting: %s, context: %v", err, ctx)
		}

		// Test decryption
		decrypted, err := utils.Decrypt(encryptionKey, encrypted)
		if err != nil {
			return fmt.Errorf("error decrypting: %s, context: %v", err, ctx)
		}

		if decrypted != "supersecret" {
			return fmt.Errorf("decrypted value does not match original, got: %s", decrypted)
		}

		return nil
	}
}

func testCheckEncryptionDecryptionWithoutKey() resource.TestCheckFunc {
	return func(s *terraform.State) error {
		// Test encryption without key (should return original value or handle accordingly)
		encrypted, err := utils.Encrypt("", "supersecret")
		if err != nil {
			return fmt.Errorf("unexpected error when encrypting without key: %s", err)
		}

		// Since there's no encryption key, we assume the value should be unchanged
		if encrypted != "supersecret" {
			return fmt.Errorf("expected encrypted value to be unchanged when no key is provided, got: %s", encrypted)
		}

		// Test decryption without key (should return original value or handle accordingly)
		decrypted, err := utils.Decrypt("", "supersecret")
		if err != nil {
			return fmt.Errorf("unexpected error when decrypting without key: %s", err)
		}

		// Since there's no encryption key, we assume the value should be unchanged
		if decrypted != "supersecret" {
			return fmt.Errorf("expected decrypted value to be unchanged when no key is provided, got: %s", decrypted)
		}

		return nil
	}
}
