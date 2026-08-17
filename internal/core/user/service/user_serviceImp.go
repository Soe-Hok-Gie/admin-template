package service

import (
	"admin-template/internal/core/user/domain"
	"admin-template/internal/core/user/dto"
	"admin-template/internal/core/user/repository"
	"context"
)

type userServiceImp struct {
	UserRepository repository.UserRepository
}

func NewUserService(userRepository repository.UserRepository) UserService {
	return &userServiceImp{UserRepository: userRepository}
}

func (service *userServiceImp) GetProfile(ctx context.Context, userID int64) (domain.User, dto.UserResponse, error) {
	//ambil data asli dari domain
	userDomain, err := service.UserRepository.GedById(ctx, userID)
	if err != nil {
		return domain.User{}, dto.UserResponse{}, err
	}
	//konversi RoleID(int64) pada domain menjadi Role (string) pada DTO
	var roleName string
	switch userDomain.RoleID {
	case 1:
		roleName = "admin"
	case 2:
		roleName = "user"
	default:
		roleName = "guest"

	}

	//mapping data domain ke dto (field password tidak ikut)
	userDto := dto.UserResponse{
		ID:    userDomain.ID,
		Name:  userDomain.Name,
		Email: userDomain.Email,
		Role:  roleName,
	}
	return userDomain, userDto, err
}
