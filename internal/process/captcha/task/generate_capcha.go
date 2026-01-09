package task

import (
	"github.com/mojocn/base64Captcha"
)

type GenerateCaptchaTask struct {
	store base64Captcha.Store
}

func NewGenerateCaptchaTask() *GenerateCaptchaTask {
	return &GenerateCaptchaTask{
		store: base64Captcha.DefaultMemStore,
	}
}

func (t *GenerateCaptchaTask) Execute() (string, string, string, error) {
	driver := base64Captcha.NewDriverString(
		80,
		240,
		60,
		base64Captcha.OptionShowSineLine|base64Captcha.OptionShowSlimeLine,
		6,
		"1234567890ABCDEFGHJKLMNOPQRSTUVWXYZ",
		nil,
		nil,
		nil,
	)

	c := base64Captcha.NewCaptcha(driver, t.store)

	return c.Generate()
}
