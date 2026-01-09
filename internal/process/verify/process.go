package verify_captcha

import (
	"context"

	captcha "github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/process/captcha/task"
)

type ReadCaptchaTask interface {
	Execute(ctx context.Context, id string) (*captcha.Captcha, error)
}

type ValidateCaptchaTask interface {
	Execute(ctx context.Context, id, val string, captcha *captcha.Captcha) error
}

type Request struct {
	CaptchaId    string `json:"captchaId"`
	CaptchaValue string `json:"captchaValue"`
}

type Response struct {
	CaptchaId string `json:"captchaId"`
}

type Process struct {
	readCaptchaTask     ReadCaptchaTask
	validateCaptchaTask ValidateCaptchaTask
}

func NewProcess(readCaptchaTask ReadCaptchaTask, validateCaptchaTask ValidateCaptchaTask) *Process {
	return &Process{
		readCaptchaTask:     readCaptchaTask,
		validateCaptchaTask: validateCaptchaTask,
	}
}

func (p *Process) Process(ctx context.Context, req Request) (*Response, error) {
	c, err := p.readCaptchaTask.Execute(ctx, req.CaptchaId)
	if err != nil {
		return nil, err
	}

	if err := p.validateCaptchaTask.Execute(ctx, req.CaptchaId, req.CaptchaValue, c); err != nil {
		return nil, err
	}

	return &Response{
		CaptchaId: req.CaptchaId,
	}, nil
}
