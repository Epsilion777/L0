package orderstorage

import (
	"L0/internal/model"
	"context"
	"fmt"
)

const getOrdersQuery = `SELECT * FROM orders`
const getDeliveiesQuery = `SELECT * FROM deliveries`
const getPaymentsQuery = `SELECT * FROM payments`
const getItemsQuery = `SELECT * FROM items`

func (p *Storage) RecoverCache(ctx context.Context) error {
	tx, err := p.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()

	// Uploading orders to the cache
	rows, err := tx.Query(getOrdersQuery)
	if err != nil {
		return fmt.Errorf("request execution error (getOrdersQuery): %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var order model.Order
		err := rows.Scan(&order.OrderUID, &order.TrackNumber, &order.Entry, &order.Locale, &order.InternalSignature,
			&order.CustomerID, &order.DeliveryService, &order.Shardkey, &order.SmID, &order.DateCreated,
			&order.OofShard)
		if err != nil {
			return fmt.Errorf("error when scanning orders lines: %w", err)
		}
		if _, ok := model.OrderCache[order.OrderUID]; !ok {
			model.OrderCache[order.OrderUID] = order
		}
	}

	// Uploading deliveries to the cache
	rows, err = tx.Query(getDeliveiesQuery)
	if err != nil {
		return fmt.Errorf("request execution error (getDeliveiesQuery): %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var delivery model.Delivery
		var orderUID string
		err := rows.Scan(&orderUID, &delivery.Name, &delivery.Phone, &delivery.Zip, &delivery.City, &delivery.Adress,
			&delivery.Region, &delivery.Email)
		if err != nil {
			return fmt.Errorf("error when scanning orders lines: %w", err)
		}
		if val, ok := model.OrderCache[orderUID]; ok {
			val.Delivery = delivery
			model.OrderCache[orderUID] = val
		}
	}

	// Uploading payments to the cache
	rows, err = tx.Query(getPaymentsQuery)
	if err != nil {
		return fmt.Errorf("request execution error (getPaymentsQuery): %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var payment model.Payment
		var orderUID string
		err := rows.Scan(&orderUID, &payment.Transaction, &payment.RequestID, &payment.Currency, &payment.Provider,
			&payment.Amount, &payment.PaymentDT, &payment.Bank, &payment.DeliveryCost, &payment.GoodsTotal,
			&payment.CustomFee)
		if err != nil {
			return fmt.Errorf("error when scanning orders lines: %w", err)
		}
		if val, ok := model.OrderCache[orderUID]; ok {
			val.Payment = payment
			model.OrderCache[orderUID] = val
		}
	}

	// Uploading items to the cache
	rows, err = tx.Query(getItemsQuery)
	if err != nil {
		return fmt.Errorf("request execution error (getItemsQuery): %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var item model.Item
		var orderUID string
		err := rows.Scan(&orderUID, &item.ChrtID, &item.TrackNumber, &item.Price, &item.Rid, &item.Name, &item.Sale,
			&item.Size, &item.TotalPrice, &item.NmID, &item.Brand, &item.Status)
		if err != nil {
			return fmt.Errorf("error when scanning orders lines: %w", err)
		}
		if val, ok := model.OrderCache[orderUID]; ok {
			val.Items = append(val.Items, item)
			model.OrderCache[orderUID] = val
		}
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("transaction confirmation error: %w", err)
	}

	fmt.Println("Data has been successfully uploaded from the database to the cache")

	return nil
}
