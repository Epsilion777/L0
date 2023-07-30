package order

import (
	"context"
)

func (uc *Usecase) RecoverCache(ctx context.Context) error {
	return uc.orderStorage.RecoverCache(ctx)
}
