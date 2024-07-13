package utils

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/crypto/pbkdf2"
)

func createHash(key string) []byte {
	return pbkdf2.Key([]byte(key), []byte("salt"), 4096, 32, sha256.New)
}

func Encrypt(key, text string) (string, error) {
	if key == "" {
		// If no key is provided, return the original text
		return text, nil
	}

	block, err := aes.NewCipher(createHash(key))
	if err != nil {
		return "", err
	}

	plaintext := []byte(text)
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return hex.EncodeToString(ciphertext), nil
}

func Decrypt(key, cryptoText string) (string, error) {
	if key == "" {
		// If no key is provided, return the original text
		return cryptoText, nil
	}

	block, err := aes.NewCipher(createHash(key))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	data, err := hex.DecodeString(cryptoText)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
func hashString(s string) string {
	h := sha256.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

type SensitivePlanModifier struct {
	encryptionKey string
}

func (m SensitivePlanModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// If password is not set in the plan, preserve the value from the state
	if req.ConfigValue.IsNull() {
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
		// Hash the new password and set it in the plan if no encryption key is provided
		planValue = hashString(req.ConfigValue.ValueString())
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
