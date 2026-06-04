package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/OmarAraby/go-ecommerce/internal/api/response"
	"github.com/OmarAraby/go-ecommerce/internal/api/validate"
	userapp "github.com/OmarAraby/go-ecommerce/internal/application/services/user"
	"github.com/OmarAraby/go-ecommerce/internal/domain"
)

type Handler struct {
	svc userapp.Service
}

func NewHandler(svc userapp.Service) *Handler {
	return &Handler{svc: svc}
}

// POST /auth/register
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name     string `json:"name"     validate:"required,min=2,max=100"`
		Email    string `json:"email"    validate:"required,email"`
		Password string `json:"password" validate:"required,min=8"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	if errs := validate.Check(req); errs != nil {
		response.ValidationError(w, errs)
		return
	}
	user, err := h.svc.Register(r.Context(), req.Name, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrConflict) {
			response.Conflict(w, "email already registered")
			return
		}
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusCreated, user)
}

// POST /auth/login
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email    string `json:"email"    validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	if errs := validate.Check(req); errs != nil {
		response.ValidationError(w, errs)
		return
	}
	result, err := h.svc.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidCredentials) {
			response.Unauthorized(w, "invalid email or password")
			return
		}
		response.InternalError(w)
		return
	}
	response.JSON(w, http.StatusOK, result)
}

// POST /auth/refresh
func (h *Handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	if errs := validate.Check(req); errs != nil {
		response.ValidationError(w, errs)
		return
	}
	pair, err := h.svc.Refresh(r.Context(), req.RefreshToken)
	if err != nil {
		response.Unauthorized(w, "invalid or expired refresh token")
		return
	}
	response.JSON(w, http.StatusOK, pair)
}

// POST /auth/logout
func (h *Handler) Logout(w http.ResponseWriter, r *http.Request) {
	var req struct {
		RefreshToken string `json:"refresh_token" validate:"required"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	_ = h.svc.Logout(r.Context(), req.RefreshToken)
	w.WriteHeader(http.StatusNoContent)
}

// POST /auth/forgot-password
func (h *Handler) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Email string `json:"email" validate:"required,email"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	if errs := validate.Check(req); errs != nil {
		response.ValidationError(w, errs)
		return
	}
	token, err := h.svc.ForgotPassword(r.Context(), req.Email)
	if err != nil {
		response.InternalError(w)
		return
	}
	// In production: send token via email, return generic message only.
	response.JSON(w, http.StatusOK, map[string]string{
		"message":     "if this email is registered you will receive a reset link",
		"reset_token": token, // remove in production
	})
}

// POST /auth/reset-password
func (h *Handler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Token       string `json:"token"        validate:"required"`
		NewPassword string `json:"new_password" validate:"required,min=8"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid request body")
		return
	}
	if errs := validate.Check(req); errs != nil {
		response.ValidationError(w, errs)
		return
	}
	if err := h.svc.ResetPassword(r.Context(), req.Token, req.NewPassword); err != nil {
		response.BadRequest(w, "invalid or expired reset token")
		return
	}
	response.JSON(w, http.StatusOK, map[string]string{"message": "password reset successfully"})
}
