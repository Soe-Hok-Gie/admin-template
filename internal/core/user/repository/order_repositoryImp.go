package repository

import (
	"admin-template/internal/core/user/domain"
	"context"
	"database/sql"
)

type orderRepositoryImp struct {
	DB *sql.DB
}

func NewOrderRepository(DB *sql.DB) OrderRepository {
	return &orderRepositoryImp{
		DB: DB,
	}
}

func (repository *orderRepositoryImp) SaveOrder(ctx context.Context, order domain.Order) error {
	script := "INSERT INTO orders (order_id,user_id, amount, status) VALUES (?,?,?,?)"
	_, err := repository.DB.ExecContext(ctx, script, order.OrderID, order.UserID, order.Amount, order.Status)
	if err != nil {
		return err
	}
	return nil
}

func (repository *orderRepositoryImp) UpdateStatus(status domain.StatusOrder) error {
	script := "UPDATE orders SET status=? WHERE id=?"
	_, err := repository.DB.Exec(script, status.OrderID, status.Status)
	if err != nil {
		return err
	}
	return nil

}
