package task

import (
	"context"
	"fmt"
	"time"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/logic/errors"
)

type SaveUsedSeedRedisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
}

type SaveUsedSeedTask struct {
	client SaveUsedSeedRedisClient
	ttl    int
}

func NewMarkSeedUsedTask(c SaveUsedSeedRedisClient, ttl int) *SaveUsedSeedTask {
	return &SaveUsedSeedTask{
		client: c,
		ttl:    ttl,
	}
}

func (t *SaveUsedSeedTask) Execute(ctx context.Context, seed string) error {
	key := fmt.Sprintf("pow:%s", seed)

	err := t.client.Set(ctx, key, "1", time.Duration(t.ttl)*time.Minute)
	if err != nil {
		return errors.ErrInternalServerError
	}

	return nil
}
