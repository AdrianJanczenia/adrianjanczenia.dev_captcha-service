package task

import (
	"context"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/logic/errors"
)

type CaptchaVerifier interface {
	VerifyCaptcha(ctx context.Context, key, answer string) (string, error)
}

type VerifyCaptchaTask struct {
	client CaptchaVerifier
}

func NewVerifyCaptchaTask(client CaptchaVerifier) *VerifyCaptchaTask {
	return &VerifyCaptchaTask{client: client}
}

func (t *VerifyCaptchaTask) Execute(ctx context.Context, id, answer string) error {
	result, err := t.client.VerifyCaptcha(ctx, "captcha:"+id, answer)
	if err != nil {
		return errors.ErrInternalServerError
	}
	switch result {
	case "solved":
		return nil
	case "not_found":
		return errors.ErrCaptchaNotFound
	case "invalid_value":
		return errors.ErrInvalidCaptchaValue
	case "no_tries":
		return errors.ErrNoTriesLeft
	default:
		return errors.ErrInternalServerError
	}
}
