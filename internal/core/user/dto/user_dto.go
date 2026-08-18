package dto

import "time"

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	RoleID   int64  `json:"role_id"`
}
type LoginRequest struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
type UserResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
}
type AuthResponse struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	RoleID    int64     `json:"role_id"` // Mengembalikan angka langsung dari DB
	CreatedAt time.Time `json:"created_at"`
}

type LoginResponse struct {
	User        UserResponse `json:"user"`
	AccessToken string       `json: "access_token"`
}

// standart response data
type Response struct {
	Code   int         `json:"code"`
	Status string      `json:"status"`
	Data   interface{} `json:"data"`
}
