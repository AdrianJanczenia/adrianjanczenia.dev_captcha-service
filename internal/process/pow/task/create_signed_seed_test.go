package task

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"
)

func TestCreateSignedSeedTask_Execute(t *testing.T) {
	secret := "test-secret"
	task := NewCreateSignedSeedTask(secret)

	seed, signature, err := task.Execute()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	parts := strings.Split(seed, ":")
	if len(parts) != 2 {
		t.Errorf("invalid seed format: %s", seed)
	}

	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(seed))
	expected := hex.EncodeToString(h.Sum(nil))

	if signature != expected {
		t.Errorf("invalid signature; got %s, want %s", signature, expected)
	}
}
