package captcha

import (
	"context"
	"errors"
	"testing"
)

type mockValidateSignatureTask struct {
	executeFunc func(seed, signature string) error
}

func (m *mockValidateSignatureTask) Execute(seed, signature string) error {
	return m.executeFunc(seed, signature)
}

type mockCheckSeedTimestampTask struct {
	executeFunc func(seed string) error
}

func (m *mockCheckSeedTimestampTask) Execute(seed string) error {
	return m.executeFunc(seed)
}

type mockValidateUsedSeedTask struct {
	executeFunc func(ctx context.Context, seed string) error
}

func (m *mockValidateUsedSeedTask) Execute(ctx context.Context, seed string) error {
	return m.executeFunc(ctx, seed)
}

type mockVerifyPowTask struct {
	executeFunc func(seed, nonce string) error
}

func (m *mockVerifyPowTask) Execute(seed, nonce string) error {
	return m.executeFunc(seed, nonce)
}

type mockSaveUsedSeedTask struct {
	executeFunc func(ctx context.Context, seed string) error
}

func (m *mockSaveUsedSeedTask) Execute(ctx context.Context, seed string) error {
	return m.executeFunc(ctx, seed)
}

type mockGenerateCaptchaTask struct {
	executeFunc func() (string, string, string, error)
}

func (m *mockGenerateCaptchaTask) Execute() (string, string, string, error) {
	return m.executeFunc()
}

type mockSaveCaptchaTask struct {
	executeFunc func(ctx context.Context, id, value string) error
}

func (m *mockSaveCaptchaTask) Execute(ctx context.Context, id, value string) error {
	return m.executeFunc(ctx, id, value)
}

func TestProcess_Captcha(t *testing.T) {
	tests := []struct {
		name                 string
		validateSigFunc      func(string, string) error
		checkTimestampFunc   func(string) error
		validateUsedSeedFunc func(context.Context, string) error
		verifyPowFunc        func(string, string) error
		saveUsedSeedFunc     func(context.Context, string) error
		generateCaptchaFunc  func() (string, string, string, error)
		saveCaptchaFunc      func(context.Context, string, string) error
		wantErr              error
		wantId               string
		wantImg              string
	}{
		{
			name:                 "successful process",
			validateSigFunc:      func(s, sig string) error { return nil },
			checkTimestampFunc:   func(s string) error { return nil },
			validateUsedSeedFunc: func(ctx context.Context, s string) error { return nil },
			verifyPowFunc:        func(s, n string) error { return nil },
			saveUsedSeedFunc:     func(ctx context.Context, s string) error { return nil },
			generateCaptchaFunc:  func() (string, string, string, error) { return "id-1", "img-1", "ans-1", nil },
			saveCaptchaFunc:      func(ctx context.Context, id, val string) error { return nil },
			wantErr:              nil,
			wantId:               "id-1",
			wantImg:              "img-1",
		},
		{
			name:                 "signature validation error",
			validateSigFunc:      func(s, sig string) error { return errors.New("sig error") },
			checkTimestampFunc:   func(s string) error { return nil },
			validateUsedSeedFunc: func(ctx context.Context, s string) error { return nil },
			verifyPowFunc:        func(s, n string) error { return nil },
			saveUsedSeedFunc:     func(ctx context.Context, s string) error { return nil },
			generateCaptchaFunc:  func() (string, string, string, error) { return "", "", "", nil },
			saveCaptchaFunc:      func(ctx context.Context, id, val string) error { return nil },
			wantErr:              errors.New("sig error"),
		},
		{
			name:                 "timestamp check error",
			validateSigFunc:      func(s, sig string) error { return nil },
			checkTimestampFunc:   func(s string) error { return errors.New("expired") },
			validateUsedSeedFunc: func(ctx context.Context, s string) error { return nil },
			verifyPowFunc:        func(s, n string) error { return nil },
			saveUsedSeedFunc:     func(ctx context.Context, s string) error { return nil },
			generateCaptchaFunc:  func() (string, string, string, error) { return "", "", "", nil },
			saveCaptchaFunc:      func(ctx context.Context, id, val string) error { return nil },
			wantErr:              errors.New("expired"),
		},
		{
			name:                 "seed already used error",
			validateSigFunc:      func(s, sig string) error { return nil },
			checkTimestampFunc:   func(s string) error { return nil },
			validateUsedSeedFunc: func(ctx context.Context, s string) error { return errors.New("double spend") },
			verifyPowFunc:        func(s, n string) error { return nil },
			saveUsedSeedFunc:     func(ctx context.Context, s string) error { return nil },
			generateCaptchaFunc:  func() (string, string, string, error) { return "", "", "", nil },
			saveCaptchaFunc:      func(ctx context.Context, id, val string) error { return nil },
			wantErr:              errors.New("double spend"),
		},
		{
			name:                 "pow verification error",
			validateSigFunc:      func(s, sig string) error { return nil },
			checkTimestampFunc:   func(s string) error { return nil },
			validateUsedSeedFunc: func(ctx context.Context, s string) error { return nil },
			verifyPowFunc:        func(s, n string) error { return errors.New("invalid work") },
			saveUsedSeedFunc:     func(ctx context.Context, s string) error { return nil },
			generateCaptchaFunc:  func() (string, string, string, error) { return "", "", "", nil },
			saveCaptchaFunc:      func(ctx context.Context, id, val string) error { return nil },
			wantErr:              errors.New("invalid work"),
		},
		{
			name:                 "save used seed error",
			validateSigFunc:      func(s, sig string) error { return nil },
			checkTimestampFunc:   func(s string) error { return nil },
			validateUsedSeedFunc: func(ctx context.Context, s string) error { return nil },
			verifyPowFunc:        func(s, n string) error { return nil },
			saveUsedSeedFunc:     func(ctx context.Context, s string) error { return errors.New("db error") },
			generateCaptchaFunc:  func() (string, string, string, error) { return "", "", "", nil },
			saveCaptchaFunc:      func(ctx context.Context, id, val string) error { return nil },
			wantErr:              errors.New("db error"),
		},
		{
			name:                 "captcha generation error",
			validateSigFunc:      func(s, sig string) error { return nil },
			checkTimestampFunc:   func(s string) error { return nil },
			validateUsedSeedFunc: func(ctx context.Context, s string) error { return nil },
			verifyPowFunc:        func(s, n string) error { return nil },
			saveUsedSeedFunc:     func(ctx context.Context, s string) error { return nil },
			generateCaptchaFunc:  func() (string, string, string, error) { return "", "", "", errors.New("gen fail") },
			saveCaptchaFunc:      func(ctx context.Context, id, val string) error { return nil },
			wantErr:              errors.New("gen fail"),
		},
		{
			name:                 "captcha save error",
			validateSigFunc:      func(s, sig string) error { return nil },
			checkTimestampFunc:   func(s string) error { return nil },
			validateUsedSeedFunc: func(ctx context.Context, s string) error { return nil },
			verifyPowFunc:        func(s, n string) error { return nil },
			saveUsedSeedFunc:     func(ctx context.Context, s string) error { return nil },
			generateCaptchaFunc:  func() (string, string, string, error) { return "id-1", "img-1", "ans-1", nil },
			saveCaptchaFunc:      func(ctx context.Context, id, val string) error { return errors.New("save fail") },
			wantErr:              errors.New("save fail"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewProcess(
				&mockValidateSignatureTask{executeFunc: tt.validateSigFunc},
				&mockCheckSeedTimestampTask{executeFunc: tt.checkTimestampFunc},
				&mockValidateUsedSeedTask{executeFunc: tt.validateUsedSeedFunc},
				&mockVerifyPowTask{executeFunc: tt.verifyPowFunc},
				&mockSaveUsedSeedTask{executeFunc: tt.saveUsedSeedFunc},
				&mockGenerateCaptchaTask{executeFunc: tt.generateCaptchaFunc},
				&mockSaveCaptchaTask{executeFunc: tt.saveCaptchaFunc},
			)

			resp, err := p.Process(context.Background(), Request{Seed: "seed", Signature: "sig", Nonce: "nonce"})

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
			if resp.CaptchaImg != tt.wantImg {
				t.Errorf("Process() CaptchaImg = %v, want %v", resp.CaptchaImg, tt.wantImg)
			}
		})
	}
}
