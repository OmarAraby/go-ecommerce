package response

import (
	"encoding/json"
	"net/http"
)

type validationErrorResponse struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Errors  map[string]string `json:"errors"`
}

// ValidationError returns 422 with a per-field error map.
func ValidationError(w http.ResponseWriter, errors map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusUnprocessableEntity)
	json.NewEncoder(w).Encode(validationErrorResponse{
		Code:    "VALIDATION_ERROR",
		Message: "validation failed",
		Errors:  errors,
	})
}
