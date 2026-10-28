package task

import (
	"context"
	"errors"
	"testing"
	"time"

	appErrors "github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/logic/errors"
)

type mockSaveUsedSeedRedisClient struct {
	setNXFunc func(ctx context.Context, key string, value any, expiration time.Duration) (bool, error)
}

func (m *mockSaveUsedSeedRedisClient) SetNX(ctx context.Context, key string, value any, expiration time.Duration) (bool, error) {
	return m.setNXFunc(ctx, key, value, expiration)
}

func TestSaveUsedSeedTask_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		m := &mockSaveUsedSeedRedisClient{
			setNXFunc: func(ctx context.Context, key string, value any, expiration time.Duration) (bool, error) {
				if key != "pow:seed" || expiration != 5*time.Minute {
					t.Fatalf("unexpected key or expiration: %q, %s", key, expiration)
				}
				return true, nil
			},
		}
		task := NewSaveUsedSeedTask(m, 5)
		if err := task.Execute(ctx, "seed"); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("already used", func(t *testing.T) {
		m := &mockSaveUsedSeedRedisClient{
			setNXFunc: func(context.Context, string, any, time.Duration) (bool, error) {
				return false, nil
			},
		}
		if err := NewSaveUsedSeedTask(m, 5).Execute(ctx, "seed"); !errors.Is(err, appErrors.ErrSeedAlreadyUsed) {
			t.Fatalf("expected ErrSeedAlreadyUsed, got %v", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		m := &mockSaveUsedSeedRedisClient{
			setNXFunc: func(ctx context.Context, key string, value any, expiration time.Duration) (bool, error) {
				return false, errors.New("fail")
			},
		}
		task := NewSaveUsedSeedTask(m, 5)
		if err := task.Execute(ctx, "seed"); !errors.Is(err, appErrors.ErrInternalServerError) {
			t.Errorf("expected ErrInternalServerError, got %v", err)
		}
	})
}
