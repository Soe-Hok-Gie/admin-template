package repository

import (
	"admin-template/internal/core/user/domain"
	"context"
	"database/sql"
)

type orderRepositoryImp struct {
	DB *sql.DB
}

func (repository *orderRepositoryImp) SaveOrder(ctx context.Context, order domain.Order) error {
	script := "INSERT INTO orders (order_id,amount, status) VALUES (?,?,?)"
	_, err := repository.DB.ExecContext(ctx, script, order.OrderID, order.Amount, order.Status)
	if err != nil {
		return err
	}
	return nil
}
