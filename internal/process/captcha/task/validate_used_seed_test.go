package task

import (
	"context"
	"errors"
	"testing"

	appErrors "github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/logic/errors"
)

type mockValidateUsedSeedRedisClient struct {
	existsFunc func(ctx context.Context, key string) (bool, error)
}

func (m *mockValidateUsedSeedRedisClient) Exists(ctx context.Context, key string) (bool, error) {
	return m.existsFunc(ctx, key)
}

func TestValidateUsedSeedTask_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("new seed", func(t *testing.T) {
		m := &mockValidateUsedSeedRedisClient{
			existsFunc: func(ctx context.Context, key string) (bool, error) {
				return false, nil
			},
		}
		task := NewValidateUsedSeedTask(m)
		if err := task.Execute(ctx, "seed"); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("already used", func(t *testing.T) {
		m := &mockValidateUsedSeedRedisClient{
			existsFunc: func(ctx context.Context, key string) (bool, error) {
				return true, nil
			},
		}
		task := NewValidateUsedSeedTask(m)
		if err := task.Execute(ctx, "seed"); err != appErrors.ErrSeedAlreadyUsed {
			t.Errorf("expected ErrSeedAlreadyUsed, got %v", err)
		}
	})

	t.Run("store error", func(t *testing.T) {
		m := &mockValidateUsedSeedRedisClient{
			existsFunc: func(ctx context.Context, key string) (bool, error) {
				return false, errors.New("redis fail")
			},
		}
		task := NewValidateUsedSeedTask(m)
		if err := task.Execute(ctx, "seed"); err != appErrors.ErrInternalServerError {
			t.Errorf("expected ErrInternalServerError, got %v", err)
		}
	})
}
