package dto

import "time"

type RegisterRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
type LoginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
type UserResponse struct {
	Name      string    `json:"Name"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginResponse struct {
	AccessToken string `json: "access_token"`
}
