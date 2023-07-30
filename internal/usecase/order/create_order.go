package order

import (
	"L0/internal/model"
	"context"
)

func (uc *Usecase) CreateOrder(ctx context.Context, order *model.Order) (string, error) {
	return uc.orderStorage.CreateOrder(ctx, order)
}
