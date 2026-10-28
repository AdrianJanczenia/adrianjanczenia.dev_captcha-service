package task

import (
	"context"
	"errors"
	"testing"

	appErrors "github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/logic/errors"
)

type fakeCaptchaVerifier struct {
	result string
	err    error
}

func (f *fakeCaptchaVerifier) VerifyCaptcha(_ context.Context, key, answer string) (string, error) {
	if key != "captcha:id" || answer != "answer" {
		return "", errors.New("unexpected arguments")
	}
	return f.result, f.err
}

func TestVerifyCaptchaTask(t *testing.T) {
	for _, tc := range []struct {
		name   string
		result string
		err    error
		want   error
	}{
		{name: "solved", result: "solved"},
		{name: "missing", result: "not_found", want: appErrors.ErrCaptchaNotFound},
		{name: "wrong answer", result: "invalid_value", want: appErrors.ErrInvalidCaptchaValue},
		{name: "last attempt", result: "no_tries", want: appErrors.ErrNoTriesLeft},
		{name: "invalid state", result: "invalid_state", want: appErrors.ErrInternalServerError},
		{name: "redis unavailable", err: errors.New("unavailable"), want: appErrors.ErrInternalServerError},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := NewVerifyCaptchaTask(&fakeCaptchaVerifier{result: tc.result, err: tc.err}).Execute(context.Background(), "id", "answer")
			if got != tc.want {
				t.Fatalf("Execute() = %v, want %v", got, tc.want)
			}
		})
	}
}
