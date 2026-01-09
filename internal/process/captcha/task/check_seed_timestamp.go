package task

import (
	"strconv"
	"strings"
	"time"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/logic/errors"
)

type CheckSeedTimestampTask struct {
	ttlMinutes int
}

func NewCheckSeedTimestampTask(ttl int) *CheckSeedTimestampTask {
	return &CheckSeedTimestampTask{
		ttlMinutes: ttl,
	}
}

func (t *CheckSeedTimestampTask) Execute(seed string) error {
	parts := strings.Split(seed, ":")
	if len(parts) != 2 {
		return errors.ErrInvalidInput
	}

	ts, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return errors.ErrInvalidInput
	}

	issuedAt := time.Unix(ts, 0)
	if time.Since(issuedAt) > time.Duration(t.ttlMinutes)*time.Minute {
		return errors.ErrPowExpired
	}

	return nil
}
