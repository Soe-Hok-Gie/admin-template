package service

import (
	"admin-template/internal/core/user/domain"
	"admin-template/internal/core/user/dto"
	"admin-template/internal/core/user/repository"
	"admin-template/internal/core/user/utils"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthServiceImp struct {
	UserRepository repository.UserRepository
}

var (
	// error global yang dipakai setiap kali register gagal (name & password are required)
	ErrInvalidInput      = errors.New("name & password are required")
	ErrUsernameExists    = errors.New("duplikat name")
	ErrInvalidCredential = errors.New("error credencial")
	ErrInternal          = errors.New("internal server error")
)

func (service *AuthServiceImp) Register(ctx context.Context, input dto.RegisterRequest) (dto.UserResponse, error) {
	if input.Name == "" || input.Password == "" {
		return dto.UserResponse{}, ErrInvalidInput
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.UserResponse{}, fmt.Errorf("failed hash :%w", err)
	}
	userDomain := domain.User{
		Name:      input.Name,
		Password:  string(hashedPassword),
		RoleID:    1,
		CreatedAt: time.Now(),
	}

	saveUser, err := service.UserRepository.Insert(ctx, userDomain)
	if err != nil {
		if repository.IsDuplicateKeyError(err) {
			return dto.UserResponse{}, ErrUsernameExists
		}
	}
	userResponse := dto.UserResponse{
		Name:      saveUser.Name,
		CreatedAt: saveUser.CreatedAt,
	}
	return userResponse, nil
}

func (service *AuthServiceImp) Login(ctx context.Context, input dto.LoginRequest) (dto.LoginResponse, error) {

	user, err := service.UserRepository.GetByEmail(ctx, input.Name)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("username not found!")
			return dto.LoginResponse{}, ErrInvalidCredential
		}
		return dto.LoginResponse{}, ErrInternal
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)) != nil {
		fmt.Println("error credential")
		return dto.LoginResponse{}, ErrInvalidCredential
	}
	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return dto.LoginResponse{
			AccessToken: token,
		}, nil
	}
	return dto.LoginResponse{}, nil
}
