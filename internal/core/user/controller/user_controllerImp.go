package controller

import (
	"admin-template/internal/core/user/dto"
	"admin-template/internal/core/user/service"
	"encoding/json"
	"errors"
	"net/http"
)

type UserControllerImp struct {
	authservice service.AuthService
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
