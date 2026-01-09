package task

import (
	"context"
	"testing"
	"time"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/logic/errors"
	taskCaptcha "github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/process/captcha/task"
)

type mockValidateCaptchaRedisClient struct {
	setFunc func(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	delFunc func(ctx context.Context, key string) error
}

func (m *mockValidateCaptchaRedisClient) Set(ctx context.Context, key string, v interface{}, e time.Duration) error {
	return m.setFunc(ctx, key, v, e)
}
func (m *mockValidateCaptchaRedisClient) Del(ctx context.Context, key string) error {
	return m.delFunc(ctx, key)
}

func TestValidateCaptchaTask_Execute(t *testing.T) {
	ctx := context.Background()

	t.Run("correct value", func(t *testing.T) {
		m := &mockValidateCaptchaRedisClient{setFunc: func(ctx context.Context, key string, v interface{}, e time.Duration) error { return nil }}
		task := NewValidateCaptchaTask(m, 3)
		state := &taskCaptcha.Captcha{Value: "123", Solved: false}
		err := task.Execute(ctx, "id", "123", state)
		if err != nil || !state.Solved {
			t.Errorf("expected solved, got err=%v", err)
		}
	})

	t.Run("wrong value, tries left", func(t *testing.T) {
		m := &mockValidateCaptchaRedisClient{setFunc: func(ctx context.Context, key string, v interface{}, e time.Duration) error { return nil }}
		task := NewValidateCaptchaTask(m, 3)
		state := &taskCaptcha.Captcha{Value: "123", TriesLeft: 2}
		err := task.Execute(ctx, "id", "wrong", state)
		if err != errors.ErrInvalidCaptchaValue || state.TriesLeft != 1 {
			t.Errorf("wrong behavior: err=%v, tries=%d", err, state.TriesLeft)
		}
	})

	t.Run("wrong value, no tries left", func(t *testing.T) {
		m := &mockValidateCaptchaRedisClient{delFunc: func(ctx context.Context, key string) error { return nil }}
		task := NewValidateCaptchaTask(m, 3)
		state := &taskCaptcha.Captcha{Value: "123", TriesLeft: 1}
		err := task.Execute(ctx, "id", "wrong", state)
		if err != errors.ErrNoTriesLeft {
			t.Errorf("expected ErrNoTriesLeft, got %v", err)
		}
	})
}
