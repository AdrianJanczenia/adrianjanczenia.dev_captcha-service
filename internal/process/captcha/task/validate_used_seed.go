package task

import (
	"context"
	"fmt"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/logic/errors"
)

type ValidateUsedSeedRedisClient interface {
	Exists(ctx context.Context, key string) (bool, error)
}

type ValidateUsedSeed struct {
	client ValidateUsedSeedRedisClient
}

func NewValidateUsedSeedTask(c ValidateUsedSeedRedisClient) *ValidateUsedSeed {
	return &ValidateUsedSeed{
		client: c,
	}
}

func (t *ValidateUsedSeed) Execute(ctx context.Context, seed string) error {
	key := fmt.Sprintf("pow:%s", seed)

	exists, err := t.client.Exists(ctx, key)
	if err != nil {
		return errors.ErrInternalServerError
	}
	if exists {
		return errors.ErrSeedAlreadyUsed
	}

	return nil
}
