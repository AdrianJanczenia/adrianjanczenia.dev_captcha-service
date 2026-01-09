package task

import (
	"context"
	"errors"
	"testing"
	"time"

	appErrors "github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/logic/errors"
)

type mockSaveUsedSeedRedisClient struct {
	setFunc func(ctx context.Context, key string, value interface{}, expiration time.Duration) error
}

func (m *mockSaveUsedSeedRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return m.setFunc(ctx, key, value, expiration)
}

func TestSaveUsedSeedTask_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		m := &mockSaveUsedSeedRedisClient{
			setFunc: func(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
				return nil
			},
		}
		task := NewMarkSeedUsedTask(m, 5)
		if err := task.Execute(ctx, "seed"); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		m := &mockSaveUsedSeedRedisClient{
			setFunc: func(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
				return errors.New("fail")
			},
		}
		task := NewMarkSeedUsedTask(m, 5)
		if err := task.Execute(ctx, "seed"); err != appErrors.ErrInternalServerError {
			t.Errorf("expected ErrInternalServerError, got %v", err)
		}
	})
}
