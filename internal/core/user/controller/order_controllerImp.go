package controller

import (
	"admin-template/internal/core/user/dto"
	"admin-template/internal/core/user/middleware"
	"encoding/json"
	"net/http"
)

type OrderControllerImp struct {
}

func (controller *OrderControllerImp) Checkout(writer http.ResponseWriter, request *http.Request) {

	ctx := request.Context()
	//panggil userId
	userID, ok := request.Context().Value(middleware.UserIDKey)(int64)
	if !ok {
		writer.Header().Set("content-type", "aplication/json")
		writer.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(writer).Encode(dto.Response{
			Code:   http.StatusUnauthorized,
			Status: "Unauthorized",
			Data:   nil,
		})
		return
	}

}
