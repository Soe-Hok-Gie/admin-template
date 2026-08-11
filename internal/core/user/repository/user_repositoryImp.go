package repository

import (
	"context"
	"database/sql"
	"log"

	"admin-template/internal/core/user/domain"

	"github.com/go-sql-driver/mysql"
)

type userRepositoryImp struct {
	DB *sql.DB
}

func (repository *userRepositoryImp) Insert(ctx context.Context, user domain.User) (domain.User, error) {
	script := "INSERT INTO users (name,email,password,role,status,created_at) Values (?,?,?,?,?,?)"
	result, err := repository.DB.ExecContext(ctx, script, user.Name, user.Email, user.Password, user.RoleID, user.Status, user.CreatedAt)
	if err != nil {
		if IsDuplicateKeyError(err) {
			return user, err
		}
		return user, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return user, err
	}
	user.ID = id
	return user, nil

}

// Helper untuk cek duplicate key MySQL
func IsDuplicateKeyError(err error) bool {
	if me, ok := err.(*mysql.MySQLError); ok {
		return me.Number == 1062
	}
	return false
}

func (repository *userRepositoryImp) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	script := "SELECT id, name, email, role FROM users WHERE name = ?"
	row := repository.DB.QueryRowContext(ctx, script, email)

	var u domain.User
	err := row.Scan(&u.ID, &u.Name, &u, email, &u.RoleID)
	if err != nil {
		log.Println("user not found:", err)
	}
	return u, nil
}
