package task

import (
	"github.com/mojocn/base64Captcha"
)

type GenerateCaptchaTask struct {
}

func NewGenerateCaptchaTask() *GenerateCaptchaTask {
	return &GenerateCaptchaTask{}
}

// The library requires a Store even though this service persists the returned answer in Redis and never reads it through base64Captcha
type discardStore struct{}

func (discardStore) Set(string, string) error         { return nil }
func (discardStore) Get(string, bool) string          { return "" }
func (discardStore) Verify(string, string, bool) bool { return false }

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

	c := base64Captcha.NewCaptcha(driver, discardStore{})

	return c.Generate()
}
