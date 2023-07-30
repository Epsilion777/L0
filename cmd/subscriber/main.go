package main

import (
	"L0/internal/delivery/rest"
	"L0/internal/storage/postgresql"
	orderstorage "L0/internal/storage/postgresql/order_storage"
	"L0/internal/usecase/order"
	"L0/subscriber"
	"context"

	"github.com/gin-gonic/gin"
)

func main() {
	channelNats := "foo"

	ctx := context.Background()

	postgresDB := postgresql.InitPostgres(ctx)
	defer postgresDB.Close()

	orderStorage := orderstorage.NewStorage(postgresDB)
	orderUsecase := order.NewUsecase(orderStorage)
	orderUsecase.RecoverCache(ctx)

	sc, sub := subscriber.StartSubscribe(ctx, channelNats, orderUsecase)
	defer sc.Close()
	defer sub.Unsubscribe()
	defer sub.Close()

	router := gin.Default()
	rest.NewHandler(orderUsecase, router)
	router.Run(":8081")
}
