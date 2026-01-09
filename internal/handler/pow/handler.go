package pow

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/logic/errors"
	process "github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/process/pow"
)

type PowProcess interface {
	Process(ctx context.Context) (*process.Response, error)
}

type Handler struct {
	process PowProcess
}

func NewHandler(p PowProcess) *Handler {
	return &Handler{
		process: p,
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		errors.WriteJSON(w, errors.ErrMethodNotAllowed)
		return
	}

	resp, err := h.process.Process(r.Context())
	if err != nil {
		errors.WriteJSON(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
