package handlers

import (
	"context"
	"net/http"
	"soft-hsm/internal/auth/dto"
	"soft-hsm/internal/auth/services"
	"soft-hsm/internal/middleware"
)

type AuthHandler struct {
	authSerice *services.AuthService
	activationService *services.ActivationService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{authSerice: authService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterDTO

	if err := middleware.DecodeJSON(r, &req); err != nil {
		middleware.ErrorHandler(w, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	// Вызываем бизнес-логику
	resp, err := h.authSerice.Register(context.Background(), req)
	if err != nil {
		middleware.ErrorHandler(w, http.StatusInternalServerError, err, "Registration failed")
		return
	}

	middleware.JSONResponse(w, http.StatusCreated, resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginDTO

	if err := middleware.DecodeJSON(r, &req); err != nil {
		middleware.ErrorHandler(w, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	resp, err := h.authSerice.Login(context.Background(), req)
	if err != nil {
		middleware.ErrorHandler(w, http.StatusInternalServerError, err, "Login failed")
		return
	}

	middleware.JSONResponse(w, http.StatusOK, resp)
}

func (h *AuthHandler) Activate(w http.ResponseWriter, r *http.Request) {
	var req dto.ActivateDTO

	if err := middleware.DecodeJSON(r, &req); err != nil {
		middleware.ErrorHandler(w, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	resp, err := h.activationService.ActiveUser(context.Background(), req.ActivateToken)

	if err != nil {
		middleware.ErrorHandler(w, http.StatusInternalServerError, err, "Activation failed")
		return
	}

	middleware.JSONResponse(w, http.StatusOK, resp)
}
