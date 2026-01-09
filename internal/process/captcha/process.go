package captcha

import (
	"context"
)

type ValidateSignatureTask interface {
	Execute(seed, signature string) error
}

type CheckSeedTimestampTask interface {
	Execute(seed string) error
}

type ValidateUsedSeedTask interface {
	Execute(ctx context.Context, seed string) error
}

type VerifyPowTask interface {
	Execute(seed, nonce string) error
}

type SaveUsedSeedTask interface {
	Execute(ctx context.Context, seed string) error
}

type GenerateCaptchaTask interface {
	Execute() (string, string, string, error)
}

type SaveCaptchaTask interface {
	Execute(ctx context.Context, id, value string) error
}

type Request struct {
	Seed      string `json:"seed"`
	Signature string `json:"signature"`
	Nonce     string `json:"nonce"`
}

type Response struct {
	CaptchaId  string `json:"captchaId"`
	CaptchaImg string `json:"captchaImg"`
}

type Process struct {
	validateSignatureTask  ValidateSignatureTask
	checkSeedTimestampTask CheckSeedTimestampTask
	validateUsedSeedTask   ValidateUsedSeedTask
	verifyPowTask          VerifyPowTask
	saveUsedSeedTask       SaveUsedSeedTask
	generateCaptchaTask    GenerateCaptchaTask
	saveCaptchaTask        SaveCaptchaTask
}

func NewProcess(
	validateSignatureTask ValidateSignatureTask,
	checkSeedTimestampTask CheckSeedTimestampTask,
	validateUsedSeedTask ValidateUsedSeedTask,
	verifyPowTask VerifyPowTask,
	saveUsedSeedTask SaveUsedSeedTask,
	generateCaptchaTask GenerateCaptchaTask,
	saveCaptchaTask SaveCaptchaTask,
) *Process {
	return &Process{
		validateSignatureTask:  validateSignatureTask,
		checkSeedTimestampTask: checkSeedTimestampTask,
		validateUsedSeedTask:   validateUsedSeedTask,
		verifyPowTask:          verifyPowTask,
		saveUsedSeedTask:       saveUsedSeedTask,
		generateCaptchaTask:    generateCaptchaTask,
		saveCaptchaTask:        saveCaptchaTask,
	}
}

func (p *Process) Process(ctx context.Context, req Request) (*Response, error) {
	if err := p.validateSignatureTask.Execute(req.Seed, req.Signature); err != nil {
		return nil, err
	}

	if err := p.checkSeedTimestampTask.Execute(req.Seed); err != nil {
		return nil, err
	}

	if err := p.validateUsedSeedTask.Execute(ctx, req.Seed); err != nil {
		return nil, err
	}

	if err := p.verifyPowTask.Execute(req.Seed, req.Nonce); err != nil {
		return nil, err
	}

	if err := p.saveUsedSeedTask.Execute(ctx, req.Seed); err != nil {
		return nil, err
	}

	id, b64s, answer, err := p.generateCaptchaTask.Execute()
	if err != nil {
		return nil, err
	}

	if err := p.saveCaptchaTask.Execute(ctx, id, answer); err != nil {
		return nil, err
	}

	return &Response{
		CaptchaId:  id,
		CaptchaImg: b64s,
	}, nil
}
