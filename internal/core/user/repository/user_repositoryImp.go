package repository

import (
	"context"
	"database/sql"

	"admin-template/internal/core/user/domain"

	"github.com/go-sql-driver/mysql"
)

type userRepositoryImp struct {
	DB *sql.DB
}

func NewUserRepository(DB *sql.DB) UserRepository {
	return &userRepositoryImp{DB: DB}
}

func (repository *userRepositoryImp) Insert(ctx context.Context, user domain.User) (domain.User, error) {
	script := "INSERT INTO users (name,email,password,status,role_id,created_at,updated_at) Values (?,?,?,?,?,?,?)"
	result, err := repository.DB.ExecContext(ctx, script, user.Name, user.Email, user.Password, user.Status, user.RoleID, user.CreatedAt, user.UpdatedAt)
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
	script := "SELECT id, name, email,password, status, role_id, created_at FROM users WHERE email = ?"
	row := repository.DB.QueryRowContext(ctx, script, email)

	var u domain.User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Status, &u.RoleID, &u.CreatedAt)
	if err != nil {
		return u, err
	}
	return u, nil
}

func (repository *userRepositoryImp) GedById(ctx context.Context, userID int64) (domain.User, error) {
	script := "SELECT id, name, email, password, status, role, created_at, updated_at FROM users WHERE id = ?"
	row := repository.DB.QueryRowContext(ctx, script, userID)

	var u domain.User
	err := row.Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Status, &u.RoleID, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, sql.ErrNoRows
		}
	}
	return u, nil

}
