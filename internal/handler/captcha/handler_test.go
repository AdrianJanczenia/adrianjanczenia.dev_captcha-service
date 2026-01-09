package captcha

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	appErrors "github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/logic/errors"
	process "github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/process/captcha"
)

type mockCaptchaProcess struct {
	processFunc func(ctx context.Context, req process.Request) (*process.Response, error)
}

func (m *mockCaptchaProcess) Process(ctx context.Context, req process.Request) (*process.Response, error) {
	return m.processFunc(ctx, req)
}

func TestHandler_Captcha(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		body       interface{}
		mockFunc   func(context.Context, process.Request) (*process.Response, error)
		wantStatus int
		wantSlug   string
	}{
		{
			name:   "success",
			method: http.MethodPost,
			body:   process.Request{Seed: "s", Signature: "sig", Nonce: "n"},
			mockFunc: func(ctx context.Context, req process.Request) (*process.Response, error) {
				return &process.Response{CaptchaId: "id-123", CaptchaImg: "img-data"}, nil
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "method not allowed",
			method:     http.MethodGet,
			body:       nil,
			wantStatus: http.StatusMethodNotAllowed,
			wantSlug:   "error_message",
		},
		{
			name:       "invalid json",
			method:     http.MethodPost,
			body:       "invalid-json",
			wantStatus: http.StatusBadRequest,
			wantSlug:   "error_message",
		},
		{
			name:   "process returns app error",
			method: http.MethodPost,
			body:   process.Request{Seed: "s", Signature: "sig", Nonce: "n"},
			mockFunc: func(ctx context.Context, req process.Request) (*process.Response, error) {
				return nil, appErrors.ErrInvalidSignature
			},
			wantStatus: http.StatusForbidden,
			wantSlug:   appErrors.ErrInvalidSignature.Slug,
		},
		{
			name:   "process returns unknown error",
			method: http.MethodPost,
			body:   process.Request{Seed: "s", Signature: "sig", Nonce: "n"},
			mockFunc: func(ctx context.Context, req process.Request) (*process.Response, error) {
				return nil, errors.New("unexpected fail")
			},
			wantStatus: http.StatusInternalServerError,
			wantSlug:   appErrors.ErrInternalServerError.Slug,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(&mockCaptchaProcess{processFunc: tt.mockFunc})

			var body []byte
			if s, ok := tt.body.(string); ok {
				body = []byte(s)
			} else if tt.body != nil {
				body, _ = json.Marshal(tt.body)
			}

			req := httptest.NewRequest(tt.method, "/captcha", bytes.NewBuffer(body))
			rr := httptest.NewRecorder()

			h.Handle(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("Handle() status = %v, wantStatus %v", rr.Code, tt.wantStatus)
			}

			if tt.wantSlug != "" {
				var resp map[string]string
				json.NewDecoder(rr.Body).Decode(&resp)
				if resp["error"] != tt.wantSlug {
					t.Errorf("Handle() error slug = %v, wantSlug %v", resp["error"], tt.wantSlug)
				}
			}
		})
	}
}
