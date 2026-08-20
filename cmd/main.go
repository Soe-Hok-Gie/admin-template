package main

import (
	"admin-template/internal/core/user/controller"
	"admin-template/internal/core/user/middleware"
	"admin-template/internal/core/user/repository"
	"admin-template/internal/core/user/service"
	"admin-template/internal/fondation"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	cfg := fondation.Load()
	db := fondation.NewDB(cfg.DB)
	defer db.Close()

	userRepository := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepository)
	userService := service.NewUserService(userRepository)
	userController := controller.NewUserController(authService, userService)

	r := mux.NewRouter()
	//rute auth
	r.HandleFunc("/auth/register", userController.Register).Methods("POST")
	r.HandleFunc("/auth/login", userController.Login).Methods("POST")

	//rute user, harus proteksi dengan middleware
	jwtMiddleware := middleware.JWTMiddleware()
	r.Handle("/auth/login/profile", jwtMiddleware(http.HandlerFunc(userController.Profile))).Methods("GET")

	fmt.Println("Database connected")
	// Menggabungkan tanda ":" dengan port aplikasi (8080) dari .env
	serverAddress := ":" + cfg.App.Port
	fmt.Println("Server Go berjalan di port", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, r))

}
