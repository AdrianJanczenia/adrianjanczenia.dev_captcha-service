package pow

import (
	"context"
	"errors"
	"testing"
)

type mockCreateSignedSeedTask struct {
	executeFunc func() (string, string, error)
}

func (m *mockCreateSignedSeedTask) Execute() (string, string, error) {
	return m.executeFunc()
}

func TestProcess_Pow(t *testing.T) {
	tests := []struct {
		name     string
		mockFunc func() (string, string, error)
		wantSeed string
		wantSig  string
		wantErr  bool
	}{
		{
			name: "success",
			mockFunc: func() (string, string, error) {
				return "seed-123", "sig-123", nil
			},
			wantSeed: "seed-123",
			wantSig:  "sig-123",
			wantErr:  false,
		},
		{
			name: "error",
			mockFunc: func() (string, string, error) {
				return "", "", errors.New("fail")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewProcess(&mockCreateSignedSeedTask{executeFunc: tt.mockFunc})
			res, err := p.Process(context.Background())

			if (err != nil) != tt.wantErr {
				t.Errorf("wantErr = %v, got %v", tt.wantErr, err)
			}
			if !tt.wantErr {
				if res.Seed != tt.wantSeed || res.Signature != tt.wantSig {
					t.Errorf("unexpected response: %+v", res)
				}
			}
		})
	}
}
