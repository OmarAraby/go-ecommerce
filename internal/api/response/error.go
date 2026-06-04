package response

import "net/http"

// Error codes — machine-readable, used by clients to handle errors programmatically.
const (
	CodeNotFound      = "NOT_FOUND"
	CodeBadRequest    = "BAD_REQUEST"
	CodeInternalError = "INTERNAL_ERROR"
	CodeConflict      = "CONFLICT"
	CodeUnauthorized  = "UNAUTHORIZED"
	CodeForbidden     = "FORBIDDEN"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func Error(w http.ResponseWriter, status int, code, message string) {
	JSON(w, status, ErrorResponse{Code: code, Message: message})
}

func NotFound(w http.ResponseWriter, message string) {
	Error(w, http.StatusNotFound, CodeNotFound, message)
}

func BadRequest(w http.ResponseWriter, message string) {
	Error(w, http.StatusBadRequest, CodeBadRequest, message)
}

func InternalError(w http.ResponseWriter) {
	Error(w, http.StatusInternalServerError, CodeInternalError, "internal server error")
}

func Unauthorized(w http.ResponseWriter, message string) {
	Error(w, http.StatusUnauthorized, CodeUnauthorized, message)
}

func Forbidden(w http.ResponseWriter, message string) {
	Error(w, http.StatusForbidden, CodeForbidden, message)
}

func Conflict(w http.ResponseWriter, message string) {
	Error(w, http.StatusConflict, CodeConflict, message)
}
