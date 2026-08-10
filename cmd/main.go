package main

import (
	"fmt"
	"admin-template/internal/fondation"
)

func main() {
	cfg := fondation.Load()

	fmt.Println("Application :", cfg.App.Name)
	fmt.Println("Environment :", cfg.App.Env)
	fmt.Println("Port        :", cfg.App.Port)

	fmt.Println()

	fmt.Println("Database Host :", cfg.DB.Host)
	fmt.Println("Database Port :", cfg.DB.Port)
	fmt.Println("Database User :", cfg.DB.User)
	fmt.Println("Database Name :", cfg.DB.Name)

	db := fondation.NewDB(cfg.DB)
	defer db.Close()

	fmt.Println("Database connected")
}
