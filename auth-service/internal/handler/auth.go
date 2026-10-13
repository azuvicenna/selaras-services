package handler

import (
	"context"

	authv1 "github.com/azuvicenna/selaras-services/auth-service/gen/auth/v1"
	"github.com/azuvicenna/selaras-services/auth-service/internal/domain"
)

type AuthHandler struct {
	authv1.UnimplementedAuthServiceServer
	uc domain.AuthUsecase
}

func NewAuthHandler(uc domain.AuthUsecase) *AuthHandler {
	return &AuthHandler{uc: uc}
}

// -----------------------------------------------------------------------------
// Captcha
// -----------------------------------------------------------------------------

func (h *AuthHandler) GenerateCaptcha(ctx context.Context, _ *authv1.GenerateCaptchaRequest) (*authv1.GenerateCaptchaResponse, error) {
	challenge, err := h.uc.GenerateCaptcha(ctx)
	if err != nil {
		return nil, toStatusError(err)
	}
	return &authv1.GenerateCaptchaResponse{
		CaptchaId:    challenge.ID,
		CaptchaImage: challenge.Image,
	}, nil
}

// -----------------------------------------------------------------------------
// Authentication
// -----------------------------------------------------------------------------

func (h *AuthHandler) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	user, tokens, err := h.uc.Login(ctx, req.GetUsername(), req.GetPassword(), req.GetCaptchaId(), req.GetCaptchaCode())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &authv1.LoginResponse{
		User:   toProtoUser(user),
		Tokens: toProtoTokenPair(tokens),
	}, nil
}

func (h *AuthHandler) RefreshToken(ctx context.Context, req *authv1.RefreshTokenRequest) (*authv1.RefreshTokenResponse, error) {
	tokens, err := h.uc.RefreshToken(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &authv1.RefreshTokenResponse{
		Tokens: toProtoTokenPair(tokens),
	}, nil
}

func (h *AuthHandler) ValidateToken(ctx context.Context, req *authv1.ValidateTokenRequest) (*authv1.ValidateTokenResponse, error) {
	valid, user, err := h.uc.ValidateToken(ctx, req.GetToken())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &authv1.ValidateTokenResponse{
		IsValid: valid,
		User:    toProtoUser(user),
	}, nil
}

func (h *AuthHandler) Logout(ctx context.Context, req *authv1.LogoutRequest) (*authv1.LogoutResponse, error) {
	if err := h.uc.Logout(ctx, req.GetRefreshToken()); err != nil {
		return nil, toStatusError(err)
	}
	return &authv1.LogoutResponse{
		Success: true,
	}, nil
}

// -----------------------------------------------------------------------------
// User Management
// -----------------------------------------------------------------------------

func (h *AuthHandler) CreateUser(ctx context.Context, req *authv1.CreateUserRequest) (*authv1.CreateUserResponse, error) {
	userEntity := domain.User{
		Username:       req.GetUsername(),
		Email:          req.GetEmail(),
		FullName:       req.GetFullName(),
		Role:           domain.Role(req.GetRole()),
		PractitionerID: req.GetPractitionerId(),
	}

	created, err := h.uc.CreateUser(ctx, userEntity, req.GetPassword())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &authv1.CreateUserResponse{
		User: toProtoUser(created),
	}, nil
}

func (h *AuthHandler) GetUserByID(ctx context.Context, req *authv1.GetUserByIDRequest) (*authv1.GetUserByIDResponse, error) {
	user, err := h.uc.GetUserByID(ctx, req.GetId())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &authv1.GetUserByIDResponse{
		User: toProtoUser(user),
	}, nil
}

func (h *AuthHandler) GetUserByUsername(ctx context.Context, req *authv1.GetUserByUsernameRequest) (*authv1.GetUserByUsernameResponse, error) {
	user, err := h.uc.GetUserByUsername(ctx, req.GetUsername())
	if err != nil {
		return nil, toStatusError(err)
	}
	return &authv1.GetUserByUsernameResponse{
		User: toProtoUser(user),
	}, nil
}

func (h *AuthHandler) ListUsers(ctx context.Context, req *authv1.ListUsersRequest) (*authv1.ListUsersResponse, error) {
	filter := domain.UserFilter{
		Query:    req.GetQuery(),
		Role:     domain.Role(req.GetRole()),
		IsActive: req.IsActive,
		Limit:    int(req.GetLimit()),
		Offset:   int(req.GetOffset()),
	}

	users, total, err := h.uc.ListUsers(ctx, filter)
	if err != nil {
		return nil, toStatusError(err)
	}

	return &authv1.ListUsersResponse{
		Users: toProtoList(users, toProtoUser),
		Total: total,
	}, nil
}

func (h *AuthHandler) UpdateUser(ctx context.Context, req *authv1.UpdateUserRequest) (*authv1.UpdateUserResponse, error) {
	userEntity := domain.User{
		ID:             req.GetId(),
		FullName:       req.GetFullName(),
		Email:          req.GetEmail(),
		Role:           domain.Role(req.GetRole()),
		PractitionerID: req.GetPractitionerId(),
		IsActive:       req.GetIsActive(),
	}

	updated, err := h.uc.UpdateUser(ctx, userEntity)
	if err != nil {
		return nil, toStatusError(err)
	}

	return &authv1.UpdateUserResponse{
		User: toProtoUser(updated),
	}, nil
}

func (h *AuthHandler) ChangePassword(ctx context.Context, req *authv1.ChangePasswordRequest) (*authv1.ChangePasswordResponse, error) {
	if err := h.uc.ChangePassword(ctx, req.GetUserId(), req.GetOldPassword(), req.GetNewPassword()); err != nil {
		return nil, toStatusError(err)
	}
	return &authv1.ChangePasswordResponse{
		Success: true,
	}, nil
}

func (h *AuthHandler) ResetPassword(ctx context.Context, req *authv1.ResetPasswordRequest) (*authv1.ResetPasswordResponse, error) {
	if err := h.uc.ResetPassword(ctx, req.GetUserId(), req.GetNewPassword()); err != nil {
		return nil, toStatusError(err)
	}
	return &authv1.ResetPasswordResponse{
		Success: true,
	}, nil
}

// -----------------------------------------------------------------------------
// Helpers
// -----------------------------------------------------------------------------

func toProtoUser(u *domain.User) *authv1.User {
	if u == nil {
		return nil
	}
	return &authv1.User{
		Id:             u.ID,
		Username:       u.Username,
		Email:          u.Email,
		FullName:       u.FullName,
		Role:           authv1.Role(u.Role),
		PractitionerId: u.PractitionerID,
		IsActive:       u.IsActive,
		LastLoginAt:    toTimestampPtr(u.LastLoginAt),
		CreatedAt:      toTimestamp(u.CreatedAt),
		UpdatedAt:      toTimestamp(u.UpdatedAt),
	}
}

func toProtoTokenPair(tp *domain.TokenPair) *authv1.TokenPair {
	if tp == nil {
		return nil
	}
	return &authv1.TokenPair{
		AccessToken:  tp.AccessToken,
		RefreshToken: tp.RefreshToken,
		ExpiresIn:    tp.ExpiresIn,
		TokenType:    tp.TokenType,
	}
}
