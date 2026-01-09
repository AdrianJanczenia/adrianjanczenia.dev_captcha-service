package task

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/logic/errors"
	captcha "github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/process/captcha/task"
)

type ReadCaptchaRedisClient interface {
	Get(ctx context.Context, key string) (string, error)
}

type ReadCaptchaTask struct {
	client ReadCaptchaRedisClient
}

func NewFetchCaptchaTask(c ReadCaptchaRedisClient) *ReadCaptchaTask {
	return &ReadCaptchaTask{
		client: c,
	}
}

func (t *ReadCaptchaTask) Execute(ctx context.Context, id string) (*captcha.Captcha, error) {
	key := fmt.Sprintf("captcha:%s", id)

	data, err := t.client.Get(ctx, key)
	if err != nil {
		return nil, errors.ErrCaptchaNotFound
	}

	var c captcha.Captcha
	if err := json.Unmarshal([]byte(data), &c); err != nil {
		return nil, errors.ErrInternalServerError
	}

	return &c, nil
}
