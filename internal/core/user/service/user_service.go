package service

import (
	"admin-template/internal/core/user/domain"
	"admin-template/internal/core/user/dto"
	"context"
)

type UserService interface {
	GetProfile(ctx context.Context, userID int64) (domain.User, dto.UserResponse, error)
}
