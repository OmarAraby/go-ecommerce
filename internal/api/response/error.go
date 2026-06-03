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
