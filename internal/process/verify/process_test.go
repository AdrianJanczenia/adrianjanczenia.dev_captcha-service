package verify_captcha

import (
	"context"
	"errors"
	"testing"

	captcha "github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/process/captcha/task"
)

type mockReadCaptchaTask struct {
	executeFunc func(ctx context.Context, id string) (*captcha.Captcha, error)
}

func (m *mockReadCaptchaTask) Execute(ctx context.Context, id string) (*captcha.Captcha, error) {
	return m.executeFunc(ctx, id)
}

type mockValidateCaptchaTask struct {
	executeFunc func(ctx context.Context, id, val string, c *captcha.Captcha) error
}

func (m *mockValidateCaptchaTask) Execute(ctx context.Context, id, val string, c *captcha.Captcha) error {
	return m.executeFunc(ctx, id, val, c)
}

func TestProcess_Verify(t *testing.T) {
	tests := []struct {
		name         string
		readFunc     func(context.Context, string) (*captcha.Captcha, error)
		validateFunc func(context.Context, string, string, *captcha.Captcha) error
		wantErr      error
		wantId       string
	}{
		{
			name: "successful verification",
			readFunc: func(ctx context.Context, id string) (*captcha.Captcha, error) {
				return &captcha.Captcha{}, nil
			},
			validateFunc: func(ctx context.Context, id, val string, c *captcha.Captcha) error {
				return nil
			},
			wantErr: nil,
			wantId:  "test-id",
		},
		{
			name: "captcha not found error",
			readFunc: func(ctx context.Context, id string) (*captcha.Captcha, error) {
				return nil, errors.New("not found")
			},
			validateFunc: func(ctx context.Context, id, val string, c *captcha.Captcha) error {
				return nil
			},
			wantErr: errors.New("not found"),
			wantId:  "",
		},
		{
			name: "validation logic error",
			readFunc: func(ctx context.Context, id string) (*captcha.Captcha, error) {
				return &captcha.Captcha{}, nil
			},
			validateFunc: func(ctx context.Context, id, val string, c *captcha.Captcha) error {
				return errors.New("invalid value")
			},
			wantErr: errors.New("invalid value"),
			wantId:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewProcess(
				&mockReadCaptchaTask{executeFunc: tt.readFunc},
				&mockValidateCaptchaTask{executeFunc: tt.validateFunc},
			)

			resp, err := p.Process(context.Background(), Request{CaptchaId: "test-id", CaptchaValue: "test-val"})

			if tt.wantErr != nil {
				if err == nil || err.Error() != tt.wantErr.Error() {
					t.Errorf("Process() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			if err != nil {
				t.Fatalf("Process() unexpected error: %v", err)
			}

			if resp.CaptchaId != tt.wantId {
				t.Errorf("Process() CaptchaId = %v, want %v", resp.CaptchaId, tt.wantId)
			}
		})
	}
}
