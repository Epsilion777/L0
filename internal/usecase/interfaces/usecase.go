package interfaces

import (
	"L0/internal/model"
	"context"
)

type OrderUsecase interface {
	CreateOrder(ctx context.Context, order *model.Order) (string, error)
	RecoverCache(ctx context.Context) error
}
