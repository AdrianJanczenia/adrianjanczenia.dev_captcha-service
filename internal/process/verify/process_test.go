package verify_captcha

import (
	"context"
	"errors"
	"testing"
)

type fakeVerifyTask struct{ err error }

func (f fakeVerifyTask) Execute(_ context.Context, id, value string) error {
	if id != "id" || value != "answer" {
		return errors.New("unexpected arguments")
	}

	return f.err
}

func TestProcess(t *testing.T) {
	for _, want := range []error{nil, errors.New("invalid answer")} {
		p := NewProcess(fakeVerifyTask{err: want})

		response, err := p.Process(context.Background(), Request{CaptchaID: "id", CaptchaValue: "answer"})
		if err != want {
			t.Fatalf("Process() error = %v, want %v", err, want)
		}
		if want == nil && (response == nil || response.CaptchaID != "id") {
			t.Fatalf("Process() response = %#v", response)
		}
		if want != nil && response != nil {
			t.Fatalf("Process() response after error = %#v", response)
		}
	}
}
