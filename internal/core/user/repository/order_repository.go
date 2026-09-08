package repository

import (
	"admin-template/internal/core/user/domain"
	"context"
)

type OrderRepository interface {
	SaveOrder(ctx context.Context, order domain.Order) error
	UpdateStatus(status domain.StatusOrder) error
	GetStatus(OrderID string) (string, error)
}
