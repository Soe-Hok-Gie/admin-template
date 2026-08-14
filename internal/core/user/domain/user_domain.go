package domain

import "time"

// tabel database
type User struct {
	ID        int64
	Name      string
	Email     string
	Password  string
	Status    string
	RoleID    int64
	RoleName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
