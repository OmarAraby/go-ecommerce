package users

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/OmarAraby/go-ecommerce/internal/api/response"
	"github.com/OmarAraby/go-ecommerce/internal/api/validate"
	userapp "github.com/OmarAraby/go-ecommerce/internal/application/services/user"
	"github.com/OmarAraby/go-ecommerce/internal/domain"
	"github.com/OmarAraby/go-ecommerce/internal/infrastructure/auth"
)

type Handler struct {
	svc userapp.Service
}

func NewHandler(svc userapp.Service) *Handler {
	return &Handler{svc: svc}
}

// GET /users/me
func (h *Handler) GetProfile(w http.ResponseWriter, r *http.Request) {
	claims := claimsFromCtx(r)
	user, err := h.svc.GetProfile(r.Context(), claims.UserID)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, user)
}

// PUT /users/me
func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name" validate:"required,min=2,max=100"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	if errs := validate.Check(req); errs != nil {
		response.ValidationError(w, errs)
		return
	}
	claims := claimsFromCtx(r)
	user, err := h.svc.UpdateProfile(r.Context(), claims.UserID, req.Name)
	if err != nil {
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, user)
}

// PUT /users/me/email
func (h *Handler) ChangeEmail(w http.ResponseWriter, r *http.Request) {
	var req struct {
		NewEmail        string `json:"new_email"        validate:"required,email"`
		CurrentPassword string `json:"current_password" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	if errs := validate.Check(req); errs != nil {
		response.ValidationError(w, errs)
		return
	}
	claims := claimsFromCtx(r)
	user, err := h.svc.ChangeEmail(r.Context(), claims.UserID, req.NewEmail, req.CurrentPassword)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			response.Unauthorized(w, "incorrect password")
			return
		}
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, user)
}

// PUT /users/me/password
func (h *Handler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		CurrentPassword string `json:"current_password" validate:"required"`
		NewPassword     string `json:"new_password"     validate:"required,min=8"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	if errs := validate.Check(req); errs != nil {
		response.ValidationError(w, errs)
		return
	}
	claims := claimsFromCtx(r)
	if err := h.svc.ChangePassword(r.Context(), claims.UserID, req.CurrentPassword, req.NewPassword); err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			response.Unauthorized(w, "incorrect current password")
			return
		}
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "password changed successfully"})
}

func claimsFromCtx(r *http.Request) *auth.Claims {
	return r.Context().Value(auth.ContextKey).(*auth.Claims)
}
