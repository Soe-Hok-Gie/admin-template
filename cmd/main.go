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
	"os"
	"strings"

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

	// //payment order
	// orderRepository := repository.NewOrderRepository(db)
	// orderService := service.NewOrderService(orderRepository)
	// orderController := controller.NewOrderController(orderService)

	// payment order
	orderRepository := repository.NewOrderRepository(db)
	serverKey := strings.TrimSpace(os.Getenv("KEY_MIDTRANS"))
	if serverKey == "" {
		log.Fatal("KEY_MIDTRANS belum diset")
	}
	orderService := service.NewOrderService(
		orderRepository,
		serverKey,
	)
	orderController := controller.NewOrderController(orderService)

	r := mux.NewRouter()
	//rute auth
	r.HandleFunc("/auth/register", userController.Register).Methods("POST")
	r.HandleFunc("/auth/login", userController.Login).Methods("POST")

	//webhook
	r.HandleFunc("/payment/notification", orderController.Webhook).Methods("POST")

	//rute user, harus proteksi dengan middleware
	jwtMiddleware := middleware.JWTMiddleware()
	r.Handle("/auth/login/profile", jwtMiddleware(http.HandlerFunc(userController.Profile))).Methods("GET")
	//checkout
	r.Handle("/orders/checkout", jwtMiddleware(http.HandlerFunc(orderController.Checkout))).Methods("POST")

	fmt.Println("Database connected")
	// Menggabungkan tanda ":" dengan port aplikasi (8080) dari .env
	serverAddress := ":" + cfg.App.Port
	fmt.Println("Server Go berjalan di port", serverAddress)
	log.Fatal(http.ListenAndServe(serverAddress, r))

}
