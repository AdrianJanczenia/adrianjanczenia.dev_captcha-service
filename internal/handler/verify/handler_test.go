package verify

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	appErrors "github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/logic/errors"
	processVerify "github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/process/verify"
)

type mockVerifyProcess struct {
	processFunc func(ctx context.Context, req processVerify.Request) (*processVerify.Response, error)
}

func (m *mockVerifyProcess) Process(ctx context.Context, req processVerify.Request) (*processVerify.Response, error) {
	return m.processFunc(ctx, req)
}

func TestHandler_Verify(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		body       interface{}
		mockFunc   func(context.Context, processVerify.Request) (*processVerify.Response, error)
		wantStatus int
		wantSlug   string
	}{
		{
			name:   "valid request",
			method: http.MethodPost,
			body:   processVerify.Request{CaptchaId: "id-123", CaptchaValue: "val-123"},
			mockFunc: func(ctx context.Context, req processVerify.Request) (*processVerify.Response, error) {
				return &processVerify.Response{CaptchaId: "id-123"}, nil
			},
			wantStatus: http.StatusOK,
		},
		{
			name:       "wrong method",
			method:     http.MethodGet,
			body:       nil,
			wantStatus: http.StatusMethodNotAllowed,
			wantSlug:   appErrors.ErrMethodNotAllowed.Slug,
		},
		{
			name:       "invalid json",
			method:     http.MethodPost,
			body:       "{invalid-json}",
			wantStatus: http.StatusBadRequest,
			wantSlug:   appErrors.ErrInvalidInput.Slug,
		},
		{
			name:   "process returns app error",
			method: http.MethodPost,
			body:   processVerify.Request{CaptchaId: "id-123", CaptchaValue: "wrong"},
			mockFunc: func(ctx context.Context, req processVerify.Request) (*processVerify.Response, error) {
				return nil, appErrors.ErrInvalidCaptchaValue
			},
			wantStatus: http.StatusBadRequest,
			wantSlug:   appErrors.ErrInvalidCaptchaValue.Slug,
		},
		{
			name:   "process returns unknown error",
			method: http.MethodPost,
			body:   processVerify.Request{CaptchaId: "id-123", CaptchaValue: "val"},
			mockFunc: func(ctx context.Context, req processVerify.Request) (*processVerify.Response, error) {
				return nil, errors.New("db error")
			},
			wantStatus: http.StatusInternalServerError,
			wantSlug:   appErrors.ErrInternalServerError.Slug,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(&mockVerifyProcess{processFunc: tt.mockFunc})

			var body []byte
			if s, ok := tt.body.(string); ok {
				body = []byte(s)
			} else if tt.body != nil {
				body, _ = json.Marshal(tt.body)
			}

			req := httptest.NewRequest(tt.method, "/verify", bytes.NewBuffer(body))
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
