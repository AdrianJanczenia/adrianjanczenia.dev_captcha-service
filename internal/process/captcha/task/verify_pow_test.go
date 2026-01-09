package task

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"strings"
	"testing"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/logic/errors"
)

func TestVerifyPowTask_Execute(t *testing.T) {
	difficulty := 4
	task := NewVerifyPowTask(difficulty)
	seed := "test-seed"

	t.Run("valid work", func(t *testing.T) {
		var foundNonce string
		prefix := strings.Repeat("0", difficulty)

		for i := 0; i < 1000000; i++ {
			n := strconv.Itoa(i)
			hash := sha256.Sum256([]byte(seed + n))
			if strings.HasPrefix(hex.EncodeToString(hash[:]), prefix) {
				foundNonce = n
				break
			}
		}

		if foundNonce == "" {
			t.Fatal("could not find valid nonce for test")
		}

		err := task.Execute(seed, foundNonce)
		if err != nil {
			t.Errorf("unexpected error for nonce %s: %v", foundNonce, err)
		}
	})

	t.Run("invalid work", func(t *testing.T) {
		err := task.Execute(seed, "invalid-nonce-12345")
		if err != errors.ErrInsufficientWork {
			t.Errorf("expected ErrInsufficientWork, got %v", err)
		}
	})
}
