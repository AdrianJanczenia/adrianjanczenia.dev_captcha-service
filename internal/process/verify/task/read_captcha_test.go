package task

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	appErrors "github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/logic/errors"
	taskCaptcha "github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/process/captcha/task"
)

type mockReadCaptchaRedisClient struct {
	getFunc func(ctx context.Context, key string) (string, error)
}

func (m *mockReadCaptchaRedisClient) Get(ctx context.Context, key string) (string, error) {
	return m.getFunc(ctx, key)
}

func TestReadCaptchaTask_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		state := taskCaptcha.Captcha{Value: "123", TriesLeft: 3}
		data, _ := json.Marshal(state)
		m := &mockReadCaptchaRedisClient{
			getFunc: func(ctx context.Context, key string) (string, error) {
				return string(data), nil
			},
		}
		task := NewFetchCaptchaTask(m)
		res, err := task.Execute(ctx, "id")
		if err != nil || res.Value != "123" {
			t.Errorf("unexpected: err=%v, res=%+v", err, res)
		}
	})

	t.Run("not found", func(t *testing.T) {
		m := &mockReadCaptchaRedisClient{
			getFunc: func(ctx context.Context, key string) (string, error) {
				return "", errors.New("not found")
			},
		}
		task := NewFetchCaptchaTask(m)
		_, err := task.Execute(ctx, "id")
		if err != appErrors.ErrCaptchaNotFound {
			t.Errorf("expected ErrCaptchaNotFound, got %v", err)
		}
	})
}
