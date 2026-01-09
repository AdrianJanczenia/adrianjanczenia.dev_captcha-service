package task

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/logic/errors"
)

func TestValidateSignatureTask_Execute(t *testing.T) {
	secret := "secret"
	task := NewValidateSignatureTask(secret)
	seed := "uuid:12345678"

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(seed))
	validSig := hex.EncodeToString(h.Sum(nil))

	t.Run("valid", func(t *testing.T) {
		err := task.Execute(seed, validSig)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("invalid", func(t *testing.T) {
		err := task.Execute(seed, "wrong-sig")
		if err != errors.ErrInvalidSignature {
			t.Errorf("expected ErrInvalidSignature, got %v", err)
		}
	})
}
