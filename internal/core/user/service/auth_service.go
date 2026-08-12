package service

import (
	"admin-template/internal/core/user/dto"
	"context"
)

type AuthService interface {
	Register(ctx context.Context, input dto.RegisterRequest) (dto.UserResponse, error)
	Login(ctx context.Context, input dto.LoginRequest) (dto.LoginResponse, error)
}
