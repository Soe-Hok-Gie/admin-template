package controller

import "net/http"

type UserController interface {
	Register(writer http.ResponseWriter, request *http.Request)
	Login(writer http.ResponseWriter, request *http.Request)
	Profile(writer http.ResponseWriter, request *http.Request)
}
