package repository

import (
	"context"

	"admin-template/internal/core/user/domain"
)

// UserRepository: Bahasa Database
type UserRepository interface {
	Insert(ctx context.Context, user domain.User) (domain.User, error)
	GetByEmail(ctx context.Context, email string) (domain.User, error)
	GedById(ctx context.Context, userID int64) (domain.User, error)
}
