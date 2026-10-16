package handler

import (
	"encoding/json"
	"net/http"

	authv1 "github.com/azuvicenna/selaras-services/auth-service/gen/auth/v1"
	"github.com/azuvicenna/selaras-services/bff-mobile/internal/client"
	"github.com/azuvicenna/selaras-services/bff-mobile/internal/middleware"
	"github.com/azuvicenna/selaras-services/bff-mobile/internal/response"
)

type AuthHandler struct {
	client *client.ClientManager
}

func NewAuthHandler(cm *client.ClientManager) *AuthHandler {
	return &AuthHandler{client: cm}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	var req authv1.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid JSON payload", err.Error())
		return
	}

	res, err := h.client.Auth.Login(r.Context(), &req)
	if err != nil {
		response.WriteGRPCError(w, err)
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Login successful", res)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		response.WriteError(w, http.StatusMethodNotAllowed, "Method not allowed", nil)
		return
	}

	var req authv1.RefreshTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.WriteError(w, http.StatusBadRequest, "Invalid JSON payload", err.Error())
		return
	}

	res, err := h.client.Auth.RefreshToken(r.Context(), &req)
	if err != nil {
		response.WriteGRPCError(w, err)
		return
	}

	response.WriteSuccess(w, http.StatusOK, "Token refreshed successfully", res)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r.Context())
	if userID == "" {
		response.WriteError(w, http.StatusUnauthorized, "User context not found", nil)
		return
	}

	res, err := h.client.Auth.GetUserByID(r.Context(), &authv1.GetUserByIDRequest{Id: userID})
	if err != nil {
		response.WriteGRPCError(w, err)
		return
	}

	response.WriteSuccess(w, http.StatusOK, "User profile retrieved", res.User)
}
