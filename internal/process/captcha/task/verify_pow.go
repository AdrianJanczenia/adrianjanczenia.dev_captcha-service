package task

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/logic/errors"
)

type VerifyPowTask struct {
	difficulty int
}

func NewVerifyPowTask(d int) *VerifyPowTask {
	return &VerifyPowTask{
		difficulty: d,
	}
}

func (t *VerifyPowTask) Execute(seed, nonce string) error {
	data := seed + nonce
	hash := sha256.Sum256([]byte(data))
	hashStr := hex.EncodeToString(hash[:])

	prefix := strings.Repeat("0", t.difficulty)
	if !strings.HasPrefix(hashStr, prefix) {
		return errors.ErrInsufficientWork
	}

	return nil
}
