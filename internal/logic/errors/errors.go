package errors

import (
	"encoding/json"
	"errors"
	"net/http"
)

type AppError struct {
	HTTPStatus int
	Slug       string
	Err        error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Slug
}

func (e *AppError) Unwrap() error {
	return e.Err
}

var (
	ErrInternalServerError = &AppError{HTTPStatus: http.StatusInternalServerError, Slug: "error_captcha_server"}
	ErrInvalidSignature    = &AppError{HTTPStatus: http.StatusForbidden, Slug: "error_pow_signature"}
	ErrSeedAlreadyUsed     = &AppError{HTTPStatus: http.StatusConflict, Slug: "error_pow_double_spend"}
	ErrPowExpired          = &AppError{HTTPStatus: http.StatusGone, Slug: "error_pow_expired"}
	ErrInsufficientWork    = &AppError{HTTPStatus: http.StatusBadRequest, Slug: "error_pow_work"}
	ErrCaptchaNotFound     = &AppError{HTTPStatus: http.StatusNotFound, Slug: "error_captcha_not_found"}
	ErrInvalidCaptchaValue = &AppError{HTTPStatus: http.StatusBadRequest, Slug: "error_captcha_invalid"}
	ErrNoTriesLeft         = &AppError{HTTPStatus: http.StatusGone, Slug: "error_captcha_expired"}
	ErrInvalidInput        = &AppError{HTTPStatus: http.StatusBadRequest, Slug: "error_message"}
	ErrMethodNotAllowed    = &AppError{HTTPStatus: http.StatusMethodNotAllowed, Slug: "error_message"}
)

func WriteJSON(w http.ResponseWriter, err error) {
	var appErr *AppError
	if !errors.As(err, &appErr) {
		appErr = ErrInternalServerError
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(appErr.HTTPStatus)
	json.NewEncoder(w).Encode(map[string]string{
		"error": appErr.Slug,
	})
}
