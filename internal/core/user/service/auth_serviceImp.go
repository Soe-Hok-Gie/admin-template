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
	"strconv"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type AuthServiceImp struct {
	UserRepository repository.UserRepository
}

func NewAuthService(userRepository repository.UserRepository) AuthService {
	return &AuthServiceImp{UserRepository: userRepository}
}

var (
	// error global yang dipakai setiap kali register gagal (name & password are required)
	ErrInvalidInput      = errors.New("name & password are required")
	ErrUsernameExists    = errors.New("duplikat name")
	ErrInvalidCredential = errors.New("error credencial")
	ErrInternal          = errors.New("internal server error")
)

func (service *AuthServiceImp) Register(ctx context.Context, input dto.RegisterRequest) (dto.AuthResponse, error) {
	if input.Name == "" || input.Password == "" {
		return dto.AuthResponse{}, ErrInvalidInput
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.AuthResponse{}, fmt.Errorf("failed hash :%w", err)
	}
	userDomain := domain.User{
		Name:      input.Name,
		Email:     input.Email,
		Password:  string(hashedPassword),
		Status:    "active",
		RoleID:    input.RoleID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	saveUser, err := service.UserRepository.Insert(ctx, userDomain)
	if err != nil {
		if repository.IsDuplicateKeyError(err) {
			return dto.AuthResponse{}, ErrUsernameExists
		}
		// JANGAN DIABAIKAN! Jika database menolak (misal karena role_id tidak ada),
		// logic harus berhenti di sini dan langsung mengembalikan error ke Postman
		return dto.AuthResponse{}, err
	}
	userResponse := dto.AuthResponse{
		ID:        saveUser.ID,
		Name:      saveUser.Name,
		Email:     saveUser.Email,
		RoleID:    saveUser.RoleID,
		CreatedAt: saveUser.CreatedAt,
	}
	return userResponse, nil
}

func (service *AuthServiceImp) Login(ctx context.Context, input dto.LoginRequest) (dto.LoginResponse, error) {

	user, err := service.UserRepository.GetByEmail(ctx, input.Email)

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
		return dto.LoginResponse{}, err
	}
	return dto.LoginResponse{
		AccessToken: token,
		User: dto.UserResponse{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      strconv.Itoa(int(user.RoleID)),
			CreatedAt: user.CreatedAt,
		},
	}, nil
}
