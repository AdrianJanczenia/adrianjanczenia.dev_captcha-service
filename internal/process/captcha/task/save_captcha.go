package task

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/logic/errors"
)

type Captcha struct {
	Value     string `json:"value"`
	TriesLeft int    `json:"triesLeft"`
	Solved    bool   `json:"solved"`
}

type SaveCaptchaRedisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
}

type SaveCaptchaTask struct {
	client     SaveCaptchaRedisClient
	ttlMinutes int
	maxTries   int
}

func NewSaveCaptchaTask(c SaveCaptchaRedisClient, ttl, max int) *SaveCaptchaTask {
	return &SaveCaptchaTask{
		client:     c,
		ttlMinutes: ttl,
		maxTries:   max,
	}
}

func (t *SaveCaptchaTask) Execute(ctx context.Context, id, value string) error {
	captcha := Captcha{
		Value:     value,
		TriesLeft: t.maxTries,
		Solved:    false,
	}

	data, err := json.Marshal(captcha)
	if err != nil {
		return errors.ErrInternalServerError
	}

	key := fmt.Sprintf("captcha:%s", id)

	if err := t.client.Set(ctx, key, string(data), time.Duration(t.ttlMinutes)*time.Minute); err != nil {
		return errors.ErrInternalServerError
	}

	return nil
}
