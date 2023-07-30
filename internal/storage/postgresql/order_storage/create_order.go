package orderstorage

import (
	"L0/internal/model"
	"context"
	"fmt"
)

const createOrder = `
	INSERT INTO 
	orders (order_uid, track_number, entry, locale, internal_signature, customer_id,
	delivery_service, shardkey, sm_id, date_created, oof_shard) 
	VALUES 
	($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
`
const createDelivery = `
	INSERT INTO
	deliveries
	(order_uid, name, phone, zip, city, address, region, email)
	VALUES 
	($1, $2, $3, $4, $5, $6, $7, $8)
`
const createPayment = `
	INSERT INTO 
	payments
	(order_uid, transaction, request_id, currency, provider, amount,
	payment_dt, bank, delivery_cost, goods_total, custom_fee)
	VALUES
	($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
`
const createItem = `
	INSERT INTO
	items 
	(order_uid, chrt_id, track_number, price, rid, name,
	sale, size, total_price, nm_id, brand, status)
	VALUES
	($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
`

func (p *Storage) CreateOrder(ctx context.Context, order *model.Order) (string, error) {
	tx, err := p.db.Begin()
	if err != nil {
		return "", fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Adding data to the orders table
	_, err = tx.Exec(createOrder,
		order.OrderUID, order.TrackNumber, order.Entry, order.Locale,
		order.InternalSignature, order.CustomerID, order.DeliveryService, order.Shardkey,
		order.SmID, order.DateCreated, order.OofShard)
	if err != nil {
		return "", fmt.Errorf("error when adding an entry to orders: %w", err)
	}

	// Adding data to the deliveries table
	_, err = tx.Exec(createDelivery,
		order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
		order.Delivery.City, order.Delivery.Adress, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		return "", fmt.Errorf("error when adding an entry to deliveries: %w", err)
	}

	// Adding data to the payments table
	_, err = tx.Exec(createPayment,
		order.OrderUID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency,
		order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDT, order.Payment.Bank,
		order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		return "", fmt.Errorf("error when adding an entry to payments: %w", err)
	}

	// Adding data to the items table
	for i := range order.Items {
		_, err = tx.Exec(createItem,
			order.OrderUID, order.Items[i].ChrtID, order.Items[i].TrackNumber, order.Items[i].Price,
			order.Items[i].Rid, order.Items[i].Name, order.Items[i].Sale, order.Items[i].Size,
			order.Items[i].TotalPrice, order.Items[i].NmID, order.Items[i].Brand, order.Items[i].Status)
		if err != nil {
			return "", fmt.Errorf("error when adding an entry to items: %w", err)
		}
	}

	err = tx.Commit()
	if err != nil {
		return "", fmt.Errorf("transaction confirmation error: %w", err)
	}

	return order.OrderUID, nil
}
