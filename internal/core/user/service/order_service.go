package service

import (
	"admin-template/internal/core/user/dto"
	"context"
)

type OrderService interface {
	CreateOrder(ctx context.Context, req dto.OrderRequest) (*dto.OrderResponse, error)
	ProcessWebhook(notification dto.WebhookNotification) error
	CheckStatus(OrderID string) (string, error)
}
