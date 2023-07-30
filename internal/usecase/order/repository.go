package order

import "L0/internal/usecase/interfaces"

type Usecase struct {
	orderStorage interfaces.OrderStorage
}

func NewUsecase(orderStorage interfaces.OrderStorage) *Usecase {
	return &Usecase{
		orderStorage: orderStorage,
	}
}
