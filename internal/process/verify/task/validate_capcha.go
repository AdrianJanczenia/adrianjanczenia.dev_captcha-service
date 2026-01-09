package task

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/logic/errors"
	captcha "github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/process/captcha/task"
)

type ValidateCaptchaRedisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Del(ctx context.Context, key string) error
}

type ValidateCaptchaTask struct {
	client     ValidateCaptchaRedisClient
	ttlMinutes int
}

func NewValidateCaptchaTask(c ValidateCaptchaRedisClient, ttl int) *ValidateCaptchaTask {
	return &ValidateCaptchaTask{
		client:     c,
		ttlMinutes: ttl,
	}
}

func (t *ValidateCaptchaTask) Execute(ctx context.Context, id, value string, captcha *captcha.Captcha) error {
	key := fmt.Sprintf("capcha:%s", id)

	if captcha.Value != value {
		captcha.TriesLeft--
		if captcha.TriesLeft <= 0 {
			t.client.Del(ctx, key)
			return errors.ErrNoTriesLeft
		}

		data, _ := json.Marshal(captcha)
		err := t.client.Set(ctx, key, string(data), time.Duration(t.ttlMinutes)*time.Minute)
		if err != nil {
			return errors.ErrInternalServerError
		}

		return errors.ErrInvalidCaptchaValue
	}
	captcha.Solved = true

	data, err := json.Marshal(captcha)
	if err != nil {
		return errors.ErrInternalServerError
	}

	return t.client.Set(ctx, key, string(data), time.Duration(t.ttlMinutes)*time.Minute)
}
