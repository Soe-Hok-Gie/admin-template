package controller

import (
	"admin-template/internal/core/user/dto"
	"admin-template/internal/core/user/middleware"
	"admin-template/internal/core/user/service"
	"encoding/json"
	"net/http"
)

type OrderControllerImp struct {
	orderService service.OrderService
}

func (controller *OrderControllerImp) Checkout(writer http.ResponseWriter, request *http.Request) {
	ctx := request.Context()

	var req dto.OrderRequest
	if err := json.NewDecoder(request.Body).Decode(&req); err != nil {
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(dto.Response{
			Code:   http.StatusBadRequest,
			Status: "Bad Request",
			Data:   "invalid request",
		})
		return
	}

	//panggil userId
	userID, ok := request.Context().Value(middleware.UserIDKey).(int64)
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

	//Masukkan userID yang sudah tervalidasi ke dalam struct req
	req.UserID = userID

	response, err := controller.orderService.CreateOrder(ctx, req)
	writer.Header().Set("content-type", "application/json")
	if err != nil {
		// ERROR 1: validasi bisnis (stok habis atau bad request)
		if err.Error() == "out of stock" || err.Error() == "product not found" {
			writer.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(writer).Encode(dto.Response{
				Code:   http.StatusBadRequest,
				Status: "Bad Request",
				Data:   err.Error(),
			})
			return
		}
		//Error 2 : Server Error
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(dto.Response{
			Code:   http.StatusInternalServerError,
			Status: "500",
			Data:   "server error: " + err.Error(),
		})
		return
	}

	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(dto.Response{
		Code:   http.StatusCreated,
		Status: "create",
		Data:   response,
	})
}

func (controller *OrderControllerImp) Webhook(writer http.ResponseWriter, request *http.Request) {

	var notification dto.WebhookNotification
	if err := json.NewDecoder(request.Body).Decode(&notification); err != nil {
		writer.Header().Set("content-type", "application/json")
		writer.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(writer).Encode(dto.Response{
			Code:   http.StatusBadRequest,
			Status: "bad request",
			Data:   "400",
		})
		return
	}

	err := controller.orderService.ProcessWebhook(notification)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(writer).Encode(dto.Response{
			Code:   http.StatusInternalServerError,
			Status: "server error",
			Data:   "500",
		})
		return
	}
	writer.WriteHeader(http.StatusOK)
	json.NewEncoder(writer).Encode(dto.Response{
		Code:   http.StatusOK,
		Status: "OK",
		Data:   "succes",
	})
	return

}
