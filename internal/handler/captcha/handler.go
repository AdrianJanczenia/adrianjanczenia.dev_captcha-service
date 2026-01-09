package captcha

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/logic/errors"
	process "github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/process/captcha"
)

type CaptchaProcess interface {
	Process(ctx context.Context, req process.Request) (*process.Response, error)
}

type Handler struct {
	process CaptchaProcess
}

func NewHandler(p CaptchaProcess) *Handler {
	return &Handler{
		process: p,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		errors.WriteJSON(w, errors.ErrMethodNotAllowed)
		return
	}

	var req process.Request
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		errors.WriteJSON(w, errors.ErrInvalidInput)
		return
	}

	resp, err := h.process.Process(r.Context(), req)
	if err != nil {
		errors.WriteJSON(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
