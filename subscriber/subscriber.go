package subscriber

import (
	"L0/internal/model"
	"L0/internal/usecase/interfaces"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/nats-io/stan.go"
)

func addMessageToDB(ctx context.Context, orderUsecase interfaces.OrderUsecase, order *model.Order) error {
	orderUID, err := orderUsecase.CreateOrder(ctx, order)
	if err != nil {
		return err
	}
	log.Printf("Order received with UUID: %s\n", orderUID)
	return nil
}

func addMessageToOrderCache(order *model.Order) {
	model.OrderCache[order.OrderUID] = *order
}

func dataProcessing(data []byte) (model.Order, error) {
	order := model.Order{}
	// Parsing of the received data
	err := json.Unmarshal(data, &order)
	if err != nil {
		return model.Order{}, fmt.Errorf("invalid data received: %w", err)
	}
	return order, nil
}

func StartSubscribe(ctx context.Context, channel string, orderUsecase interfaces.OrderUsecase) (stan.Conn, stan.Subscription) {
	clusterID := "test-cluster"
	clientID := "consumer-1"
	natsURL := "nats://127.0.0.1:4222"

	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL))
	if err != nil {
		log.Fatalf("Error connecting to the NATS Streaming server: %v", err)
	}
	sub, err := sc.Subscribe(channel, func(msg *stan.Msg) {
		order, err := dataProcessing(msg.Data)
		if err != nil {
			log.Printf("ERROR: %v", err)
			return
		}
		err = addMessageToDB(ctx, orderUsecase, &order)
		if err != nil {
			log.Printf("ERROR: %v", err)
			return
		}
		addMessageToOrderCache(&order)

	}, stan.DurableName("durable-subscriber"))
	if err != nil {
		log.Fatalf("Channel subscription error: %v", err)
	}

	return sc, sub
}
