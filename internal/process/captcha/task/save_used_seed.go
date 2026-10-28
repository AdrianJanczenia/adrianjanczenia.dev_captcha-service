package task

import (
	"context"
	"fmt"
	"time"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/logic/errors"
)

type SaveUsedSeedRedisClient interface {
	SetNX(ctx context.Context, key string, value any, expiration time.Duration) (bool, error)
}

type SaveUsedSeedTask struct {
	client SaveUsedSeedRedisClient
	ttl    int
}

func NewSaveUsedSeedTask(c SaveUsedSeedRedisClient, ttl int) *SaveUsedSeedTask {
	return &SaveUsedSeedTask{
		client: c,
		ttl:    ttl,
	}
}

func (t *SaveUsedSeedTask) Execute(ctx context.Context, seed string) error {
	key := fmt.Sprintf("pow:%s", seed)

	created, err := t.client.SetNX(ctx, key, "1", time.Duration(t.ttl)*time.Minute)
	if err != nil {
		return errors.ErrInternalServerError
	}
	if !created {
		return errors.ErrSeedAlreadyUsed
	}

	return nil
}
