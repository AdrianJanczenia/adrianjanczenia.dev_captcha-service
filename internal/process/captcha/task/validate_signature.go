package task

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/logic/errors"
)

type ValidateSignatureTask struct {
	hmacSecret string
}

func NewValidateSignatureTask(hmacSecret string) *ValidateSignatureTask {
	return &ValidateSignatureTask{
		hmacSecret: hmacSecret,
	}
}

func (t *ValidateSignatureTask) Execute(seed, signature string) error {
	h := hmac.New(sha256.New, []byte(t.hmacSecret))
	h.Write([]byte(seed))
	expected := hex.EncodeToString(h.Sum(nil))

	if expected != signature {
		return errors.ErrInvalidSignature
	}

	return nil
}
