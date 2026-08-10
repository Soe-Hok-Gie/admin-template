package repository

import (
	"context"

	"admin-template/internal/core/user/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
}
