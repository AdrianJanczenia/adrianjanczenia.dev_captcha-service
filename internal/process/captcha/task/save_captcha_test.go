package task

import (
	"context"
	"errors"
	"testing"
	"time"

	appErrors "github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/logic/errors"
)

type mockSaveCaptchaRedisClient struct {
	setFunc func(ctx context.Context, key string, value interface{}, expiration time.Duration) error
}

func (m *mockSaveCaptchaRedisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return m.setFunc(ctx, key, value, expiration)
}

func TestSaveCaptchaTask_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		m := &mockSaveCaptchaRedisClient{
			setFunc: func(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
				return nil
			},
		}
		task := NewSaveCaptchaTask(m, 3, 3)
		if err := task.Execute(ctx, "id", "answer"); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		m := &mockSaveCaptchaRedisClient{
			setFunc: func(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
				return errors.New("fail")
			},
		}
		task := NewSaveCaptchaTask(m, 3, 3)
		if err := task.Execute(ctx, "id", "answer"); err != appErrors.ErrInternalServerError {
			t.Errorf("expected ErrInternalServerError, got %v", err)
		}
	})
}
