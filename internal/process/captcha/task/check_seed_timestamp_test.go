package task

import (
	"fmt"
	"testing"
	"time"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/logic/errors"
)

func TestCheckTimestampTask_Execute(t *testing.T) {
	ttl := 5
	task := NewCheckSeedTimestampTask(ttl)

	t.Run("fresh", func(t *testing.T) {
		seed := fmt.Sprintf("id:%d", time.Now().Unix())
		err := task.Execute(seed)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("expired", func(t *testing.T) {
		old := time.Now().Add(-10 * time.Minute).Unix()
		seed := fmt.Sprintf("id:%d", old)
		err := task.Execute(seed)
		if err != errors.ErrPowExpired {
			t.Errorf("expected ErrPowExpired, got %v", err)
		}
	})

	t.Run("invalid format", func(t *testing.T) {
		err := task.Execute("invalid-seed")
		if err != errors.ErrInvalidInput {
			t.Errorf("expected ErrInvalidInput, got %v", err)
		}
	})
}
