package pow

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	processPow "github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/process/pow"
)

type mockPowProcess struct {
	processFunc func(ctx context.Context) (*processPow.Response, error)
}

func (m *mockPowProcess) Process(ctx context.Context) (*processPow.Response, error) {
	return m.processFunc(ctx)
}

func TestHandler_Pow(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		mockFunc   func(context.Context) (*processPow.Response, error)
		wantStatus int
		wantSlug   string
	}{
		{
			name:   "success",
			method: http.MethodGet,
			mockFunc: func(ctx context.Context) (*processPow.Response, error) {
				return &processPow.Response{Seed: "s", Signature: "sig"}, nil
			},
			wantStatus: http.StatusOK,
		},
		{
			name:   "method not allowed",
			method: http.MethodPost,
			mockFunc: func(ctx context.Context) (*processPow.Response, error) {
				return nil, nil
			},
			wantStatus: http.StatusMethodNotAllowed,
			wantSlug:   "error_message",
		},
		{
			name:   "process error",
			method: http.MethodGet,
			mockFunc: func(ctx context.Context) (*processPow.Response, error) {
				return nil, errors.New("internal fail")
			},
			wantStatus: http.StatusInternalServerError,
			wantSlug:   "error_captcha_server",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(&mockPowProcess{processFunc: tt.mockFunc})
			req := httptest.NewRequest(tt.method, "/pow", nil)
			rr := httptest.NewRecorder()

			h.Handle(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("status = %v, want %v", rr.Code, tt.wantStatus)
			}
			if tt.wantSlug != "" {
				var resp map[string]string
				json.NewDecoder(rr.Body).Decode(&resp)
				if resp["error"] != tt.wantSlug {
					t.Errorf("slug = %v, want %v", resp["error"], tt.wantSlug)
				}
			}
		})
	}
}
