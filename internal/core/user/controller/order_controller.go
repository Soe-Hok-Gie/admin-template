package controller

import "net/http"

type OrderController interface {
	Checkout(writer http.ResponseWriter, request *http.Request)
	Webhook(writer http.ResponseWriter, request *http.Request)
	Checkstatus(writer http.ResponseWriter, request *http.Request)
}
