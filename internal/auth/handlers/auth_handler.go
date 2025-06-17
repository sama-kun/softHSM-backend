package handlers

import (
	"context"
	"fmt"
	"net/http"
	"soft-hsm/internal/auth/dto"
	"soft-hsm/internal/auth/services"
	"soft-hsm/internal/common/validators"
	"soft-hsm/internal/middleware"
)

type AuthHandler struct {
	authSerice        *services.AuthService
	activationService *services.ActivationService
}

func NewAuthHandler(authService *services.AuthService, activationService *services.ActivationService) *AuthHandler {
	return &AuthHandler{authSerice: authService, activationService: activationService}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterDTO

	if err := middleware.DecodeJSON(r, &req); err != nil {
		middleware.ErrorHandler(w, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	resp, err := h.authSerice.Register(context.Background(), req)
	if err != nil {
		middleware.ErrorHandler(w, http.StatusInternalServerError, err, "Registration failed")
		return
	}

	middleware.JSONResponse(w, http.StatusCreated, resp)
}

// func (h *AuthHandler) SetMasterPassword(w http.ResponseWriter, r *http.Request) {
// 	var req dto.SetMasterPassword
// 	user, err := middleware.GetUserFromContext(r)

// 	if err != nil {
// 		// http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 		middleware.ErrorHandler(w, http.StatusUnauthorized, err, "Unauthorized")
// 		return
// 	}

// 	if err := middleware.DecodeJSON(r, &req); err != nil {
// 		middleware.ErrorHandler(w, http.StatusBadRequest, err, "Invalid request body")
// 		return
// 	}

// 	if err := validators.ValidateStruct(req); err != nil {
// 		middleware.ErrorHandler(w, http.StatusBadRequest, err, "invalid input")
// 		return
// 	}

// 	// resp, err := h.authSerice.SetMasterPassword(context.Background(), int64(user.Id), req.MasterPassword)
// 	if err != nil {
// 		middleware.ErrorHandler(w, http.StatusInternalServerError, err, "Set Master Password failed")
// 		return
// 	}

// 	middleware.JSONResponse(w, http.StatusOK, resp)
// }

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginDTO

	if err := middleware.DecodeJSON(r, &req); err != nil {
		fmt.Println(err)
		middleware.ErrorHandler(w, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	resp, err := h.authSerice.Login(context.Background(), req)
	if err != nil {
		fmt.Println(err)
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

func (h *AuthHandler) CheckMasterPassword(w http.ResponseWriter, r *http.Request) {
	var req dto.SetMasterPassword
	// user, err := middleware.GetUserFromContext(r)

	// if err != nil {
	// 	// http.Error(w, "Unauthorized", http.StatusUnauthorized)
	// 	fmt.Println(err)
	// 	middleware.ErrorHandler(w, http.StatusUnauthorized, err, "Unauthorized")
	// 	return
	// }

	if err := middleware.DecodeJSON(r, &req); err != nil {
		fmt.Println(err)
		middleware.ErrorHandler(w, http.StatusBadRequest, err, "Invalid request body")
		return
	}

	if err := validators.ValidateStruct(req); err != nil {
		fmt.Println(err)
		middleware.ErrorHandler(w, http.StatusBadRequest, err, "invalid input")
		return
	}

	resp, err := h.authSerice.CheckMasterPassword(context.Background(), req.SessionToken, req.Otp)

	if err != nil {
		fmt.Println(err)
		middleware.ErrorHandler(w, http.StatusBadRequest, err, "cannot generate token")
		return
	}

	middleware.JSONResponse(w, http.StatusOK, resp)
}

func (h *AuthHandler) ResetPassword(w http.ResponseWriter, r *http.Request) {
	var req dto.ResetPasswordDTO
	if err := middleware.DecodeJSON(r, &req); err != nil {
		fmt.Println(err)
		middleware.ErrorHandler(w, http.StatusBadRequest, err, "Invalid request body")
		return
	}
	fmt.Println("REQ:", req)
	user, err := middleware.GetUserFromContext(r)

	if err != nil {
		middleware.ErrorHandler(w, http.StatusUnauthorized, err, "Unauthorized")
		return
	}

	if err := validators.ValidateStruct(req); err != nil {
		middleware.ErrorHandler(w, http.StatusBadRequest, err, "invalid input")
		return
	}

	resp, err := h.authSerice.ResetPassword(context.Background(), int64(user.Id), req)

	fmt.Println("Error", err)

	middleware.JSONResponse(w, http.StatusOK, resp)
}
