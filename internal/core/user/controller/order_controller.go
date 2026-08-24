package controller

import "net/http"

type OrderController interface {
	Checkout(writer http.ResponseWriter, request *http.Request)
}
