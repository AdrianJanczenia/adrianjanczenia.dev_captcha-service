package task

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type CreateSignedSeedTask struct {
	hmacSecret string
}

func NewCreateSignedSeedTask(hmacSecret string) *CreateSignedSeedTask {
	return &CreateSignedSeedTask{
		hmacSecret: hmacSecret,
	}
}

func (t *CreateSignedSeedTask) Execute() (string, string, error) {
	timestamp := time.Now().Unix()
	id := uuid.New().String()
	seed := fmt.Sprintf("%s:%d", id, timestamp)

	h := hmac.New(sha256.New, []byte(t.hmacSecret))
	h.Write([]byte(seed))

	signature := hex.EncodeToString(h.Sum(nil))

	return seed, signature, nil
}
