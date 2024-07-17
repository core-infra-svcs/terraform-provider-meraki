package utils

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"io"
)

// Encrypt encrypts the given plaintext with the provided key.
func Encrypt(key, text string) (string, error) {
	block, err := aes.NewCipher([]byte(createHash(key)))
	if err != nil {
		return "", err
	}
	plaintext := []byte(text)
	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
	iv := ciphertext[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}
	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(ciphertext[aes.BlockSize:], plaintext)
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts the given ciphertext with the provided key.
func Decrypt(key, cryptoText string) (string, error) {
	block, err := aes.NewCipher([]byte(createHash(key)))
	if err != nil {
		return "", err
	}
	ciphertext, err := base64.URLEncoding.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}
	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)
	return string(ciphertext), nil
}

// createHash creates a hash of the given key.
func createHash(key string) string {
	hash := sha256.New()
	hash.Write([]byte(key))
	return base64.URLEncoding.EncodeToString(hash.Sum(nil))
}

type SensitivePlanModifier struct {
	encryptionKey string
}

func (m SensitivePlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// If password is not set in the plan, preserve the value from the state
	// Preserve the value from the state if the config value is null or unknown
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		resp.PlanValue = req.StateValue
		return
	}

	var err error
	var planValue string

	if m.encryptionKey != "" {
		// Encrypt the new password and set it in the plan
		planValue, err = Encrypt(m.encryptionKey, req.ConfigValue.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Encrypting Value",
				fmt.Sprintf("Could not Encrypt the password: %s", err),
			)
			return
		}
	} else {
		// If encryption key is not set, use the plain value directly
		planValue = req.ConfigValue.ValueString()
	}

	resp.PlanValue = types.StringValue(planValue)
}

func (m SensitivePlanModifier) Description(ctx context.Context) string {
	return "Custom plan modifier to hash or Encrypt sensitive strings."
}

func (m SensitivePlanModifier) MarkdownDescription(ctx context.Context) string {
	return "Custom plan modifier to hash or Encrypt sensitive strings."
}

func NewSensitivePlanModifier(encryptionKey string) planmodifier.String {
	return SensitivePlanModifier{encryptionKey: encryptionKey}
}
