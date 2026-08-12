package controller

import "net/http"

type UserController interface {
	Register(writer http.ResponseWriter, request *http.Request)
}
