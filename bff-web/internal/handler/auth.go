package handler

import (
	"net/http"
	"strconv"

	authv1 "github.com/azuvicenna/selaras-services/auth-service/gen/auth/v1"
	"github.com/azuvicenna/selaras-services/bff-web/internal/client"
	"github.com/azuvicenna/selaras-services/bff-web/internal/middleware"
)

type AuthHandler struct {
	cm *client.ClientManager
}

func NewAuthHandler(cm *client.ClientManager) *AuthHandler {
	return &AuthHandler{cm: cm}
}

func (h *AuthHandler) GenerateCaptcha(w http.ResponseWriter, r *http.Request) {
	resp, err := h.cm.Auth.GenerateCaptcha(r.Context(), &authv1.GenerateCaptchaRequest{})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req authv1.LoginRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_BODY", err.Error())
		return
	}

	resp, err := h.cm.Auth.Login(r.Context(), &req)
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp)
}

func (h *AuthHandler) RefreshToken(w http.ResponseWriter, r *http.Request) {
	var req authv1.RefreshTokenRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_BODY", err.Error())
		return
	}

	resp, err := h.cm.Auth.RefreshToken(r.Context(), &req)
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req authv1.LogoutRequest
	_ = ReadJSON(r, &req)

	resp, err := h.cm.Auth.Logout(r.Context(), &req)
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp)
}

func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	claims, ok := middleware.GetUserClaims(r.Context())
	if !ok {
		WriteError(w, http.StatusUnauthorized, "UNAUTHORIZED", "User claims not found in context")
		return
	}

	resp, err := h.cm.Auth.GetUserByID(r.Context(), &authv1.GetUserByIDRequest{Id: claims.UserID})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.User)
}

func (h *AuthHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	resp, err := h.cm.Auth.ListUsers(r.Context(), &authv1.ListUsersRequest{
		Query:  r.URL.Query().Get("search"),
		Limit:  int32(pageSize),
		Offset: int32(offset),
	})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}

	WritePaginated(w, resp.Users, page, pageSize, resp.Total)
}

func (h *AuthHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req authv1.CreateUserRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_BODY", err.Error())
		return
	}

	resp, err := h.cm.Auth.CreateUser(r.Context(), &req)
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteCreated(w, resp.User)
}

func (h *AuthHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "id is required")
		return
	}

	resp, err := h.cm.Auth.GetUserByID(r.Context(), &authv1.GetUserByIDRequest{Id: id})
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.User)
}

func (h *AuthHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		WriteError(w, http.StatusBadRequest, "INVALID_PARAM", "id is required")
		return
	}

	var req authv1.UpdateUserRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_BODY", err.Error())
		return
	}
	req.Id = id

	resp, err := h.cm.Auth.UpdateUser(r.Context(), &req)
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp.User)
}

func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var req authv1.ChangePasswordRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteError(w, http.StatusBadRequest, "INVALID_BODY", err.Error())
		return
	}

	claims, ok := middleware.GetUserClaims(r.Context())
	if ok && req.UserId == "" {
		req.UserId = claims.UserID
	}

	resp, err := h.cm.Auth.ChangePassword(r.Context(), &req)
	if err != nil {
		WriteGRPCError(w, err)
		return
	}
	WriteSuccess(w, resp)
}
