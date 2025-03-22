package handlers

import (
	"context"
	"net/http"
	"soft-hsm/internal/middleware"
	"soft-hsm/internal/user/services"
)

type UserHandler struct {
	userService *services.UserService
}

func NewUserHandler(userService *services.UserService ) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) Me(w http.ResponseWriter, r *http.Request){
	user, err := middleware.GetUserFromContext(r);

	if err != nil {
		// http.Error(w, "Unauthorized", http.StatusUnauthorized)
		middleware.ErrorHandler(w, http.StatusUnauthorized, err, "Unauthorized")
		return
	}

	resp, err := h.userService.Me(context.Background(), int64(user.Id))
	if err != nil {
		middleware.ErrorHandler(w, http.StatusInternalServerError, err, "Registration failed")
		return
	}

	middleware.JSONResponse(w, http.StatusAccepted, resp)
}