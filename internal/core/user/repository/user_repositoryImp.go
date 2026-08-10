package repository

import (
	"context"

	"admin-template/internal/core/user/domain"
)

type userRepositoryImp struct {
}

func (UserRepository *userRepositoryImp) Create(ctx context.Context, user *domain.User) error {

	return nil

}
