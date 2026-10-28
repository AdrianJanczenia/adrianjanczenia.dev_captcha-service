package verify_captcha

import (
	"context"
)

type VerifyCaptchaTask interface {
	Execute(ctx context.Context, id, val string) error
}

type Request struct {
	CaptchaID    string `json:"captchaId"`
	CaptchaValue string `json:"captchaValue"`
}

type Response struct {
	CaptchaID string `json:"captchaId"`
}

type Process struct {
	verifyCaptchaTask VerifyCaptchaTask
}

func NewProcess(task VerifyCaptchaTask) *Process {
	return &Process{verifyCaptchaTask: task}
}

func (p *Process) Process(ctx context.Context, req Request) (*Response, error) {
	if err := p.verifyCaptchaTask.Execute(ctx, req.CaptchaID, req.CaptchaValue); err != nil {
		return nil, err
	}

	return &Response{
		CaptchaID: req.CaptchaID,
	}, nil
}
