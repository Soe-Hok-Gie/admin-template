package controller

import (
	"admin-template/internal/core/user/dto"
	"admin-template/internal/core/user/middleware"
	"admin-template/internal/core/user/service"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type UserControllerImp struct {
	authservice service.AuthService
	userservice service.UserService
}

func NewUserController(
	authservice service.AuthService,
	userservice service.UserService,
) UserController {
	return &UserControllerImp{
		authservice: authservice,
		userservice: userservice,
	}
}

func (controller *UserControllerImp) Register(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	var req dto.RegisterRequest
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		writer.Header().Set("content-type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(dto.Response{
			Code:   http.StatusBadRequest,
			Status: "Bad Request",
			Data:   "Input Salah",
		})
		return
	}
	user, err := controller.authservice.Register(ctx, req)
	if err != nil {
		status := http.StatusInternalServerError
		msg := "Internal server error"

		// Mapping error ke pesan & status code
		switch {
		case errors.Is(err, service.ErrInvalidInput):
			status = http.StatusBadRequest
			msg = "Username & password are required"
		case errors.Is(err, service.ErrUsernameExists):
			status = http.StatusConflict
			msg = "username already exists"
		}
		writer.WriteHeader(status)
		json.NewEncoder(writer).Encode(dto.Response{
			Code:   status,
			Status: http.StatusText(status),
			Data:   msg,
		})
		return
	}
	writer.Header().Set("content-type", "application/jso")
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(dto.Response{
		Code:   http.StatusCreated,
		Status: "Create",
		Data:   user,
	})
}

func (controller *UserControllerImp) Login(writer http.ResponseWriter, request *http.Request) {

	ctx := request.Context()

	var input dto.LoginRequest
	if err := json.NewDecoder(request.Body).Decode(&input); err != nil {
		writer.Header().Set("content-type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(dto.Response{
			Code:   http.StatusBadRequest,
			Status: "Bad Request",
			Data:   "input salah",
		})
		return
	}
	userResponse, err := controller.authservice.Login(ctx, input)
	if err != nil {
		writer.Header().Set("Content-Type", "application/json")
		//errorcredensial
		if err == service.ErrInvalidCredential {
			writer.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(writer).Encode(dto.Response{
				Code:   http.StatusUnauthorized,
				Status: "unauthorize",
				Data:   "Username or password failed",
			})
			return
		}
		writer.Header().Set("content-type", "application/json")
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(dto.Response{
			Code:   http.StatusInternalServerError,
			Status: "500",
			Data:   "internal server error",
		})
		return
	}

	http.SetCookie(writer, &http.Cookie{
		Name:     "token",
		Value:    userResponse.AccessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		MaxAge:   900,
	})
	writer.Header().Set("content-type", "application/json")
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(dto.Response{
		Code:   http.StatusOK,
		Status: "ok",
		Data:   userResponse.User,
	})

}

func (controller *UserControllerImp) Profile(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()
	//panggil middleware
	userID, ok := request.Context().Value(middleware.UserIDKey).(int64)
	if !ok {
		writer.Header().Set("content-type", "application/json")
		writer.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(writer).Encode(dto.Response{
			Code:   http.StatusUnauthorized,
			Status: "unauthorized",
			Data:   nil,
		})
		return
	}
	// Gunakan blank identifier (_) untuk MENOLAK/MEMBUANG domain.User yang berisi password,jadi kita pakai data dari dto
	_, userDto, err := controller.userservice.GetProfile(ctx, userID)
	if err != nil {
		fmt.Println("Profil error:", err)
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		return
	}
	writer.Header().Set("content-type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(dto.Response{
		Code:   http.StatusCreated,
		Status: "Created",
		Data:   userDto,
	})
}
