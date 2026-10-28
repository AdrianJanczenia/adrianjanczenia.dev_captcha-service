package redis

import (
	_ "embed"

	"github.com/go-redis/redis/v8"
)

//go:embed scripts/verify_captcha.lua
var verifyCaptchaSource string

var verifyCaptchaScript = redis.NewScript(verifyCaptchaSource)
