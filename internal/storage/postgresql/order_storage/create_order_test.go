package orderstorage

import (
	"L0/internal/model"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreateOrder(t *testing.T) {
	// Create a mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create mock database connection: %v", err)
	}
	defer db.Close()

	// Create a storage instance with the mock database
	storage := &Storage{
		db: db,
	}

	// Create a test order
	orderUID := uuid.New().String()
	order := &model.Order{
		OrderUID:    orderUID,
		TrackNumber: "WBILMTESTTRACK",
		Entry:       "WBIL",
		Delivery: model.Delivery{
			Name:   "Test Testov",
			Phone:  "+9720000000",
			Zip:    "2639809",
			City:   "Kiryat Mozkin",
			Adress: "Ploshad Mira 15",
			Region: "Kraiot",
			Email:  "test@gmail.com",
		},
		Payment: model.Payment{
			Transaction:  "b563feb7b2b84b6test",
			RequestID:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDT:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items: []model.Item{
			{
				ChrtID:      9934930,
				TrackNumber: "WBILMTESTTRACK",
				Price:       453,
				Rid:         "ab4219087a764ae0btest",
				Name:        "Mascaras",
				Sale:        30,
				Size:        "5",
				TotalPrice:  317,
				NmID:        2389212,
				Brand:       "VivienneSabo",
				Status:      202,
			},
		},
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test",
		DeliveryService:   "meest",
		Shardkey:          "9",
		SmID:              99,
		DateCreated:       time.Now().Format("2006-01-02T15:04:05Z"),
		OofShard:          "1",
	}

	tests := []struct {
		name      string
		expectErr bool
		mockFunc  func()
	}{
		{
			"Correct test",
			false,
			func() {
				// Set up expectations for the mock
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO orders").WithArgs(order.OrderUID, order.TrackNumber, order.Entry, order.Locale,
					order.InternalSignature, order.CustomerID, order.DeliveryService, order.Shardkey,
					order.SmID, order.DateCreated, order.OofShard,
				).WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("INSERT INTO deliveries").WithArgs(
					order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
					order.Delivery.City, order.Delivery.Adress, order.Delivery.Region, order.Delivery.Email,
				).WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("INSERT INTO payments").WithArgs(
					order.OrderUID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency,
					order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDT, order.Payment.Bank,
					order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee,
				).WillReturnResult(sqlmock.NewResult(1, 1))

				for i := range order.Items {
					mock.ExpectExec("INSERT INTO items").WithArgs(
						order.OrderUID, order.Items[i].ChrtID, order.Items[i].TrackNumber, order.Items[i].Price,
						order.Items[i].Rid, order.Items[i].Name, order.Items[i].Sale, order.Items[i].Size,
						order.Items[i].TotalPrice, order.Items[i].NmID, order.Items[i].Brand, order.Items[i].Status,
					).WillReturnResult(sqlmock.NewResult(1, 1))
				}

				mock.ExpectCommit()
			},
		},
		{
			"Orders err test",
			true,
			func() {
				// Set up expectations for the mock
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO orders").WithArgs(order.OrderUID, order.TrackNumber, order.Entry, order.Locale,
					order.InternalSignature, order.CustomerID, order.DeliveryService, order.Shardkey,
					order.SmID, order.DateCreated, order.OofShard,
				).WillReturnError(errors.New("some error"))
				mock.ExpectRollback()
			},
		},
		{
			"Deliveries err test",
			true,
			func() {
				// Set up expectations for the mock
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO orders").WithArgs(order.OrderUID, order.TrackNumber, order.Entry, order.Locale,
					order.InternalSignature, order.CustomerID, order.DeliveryService, order.Shardkey,
					order.SmID, order.DateCreated, order.OofShard,
				).WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("INSERT INTO deliveries").WithArgs(
					order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
					order.Delivery.City, order.Delivery.Adress, order.Delivery.Region, order.Delivery.Email,
				).WillReturnError(errors.New("some error"))

				mock.ExpectRollback()
			},
		},
		{
			"Payments err test",
			true,
			func() {
				// Set up expectations for the mock
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO orders").WithArgs(order.OrderUID, order.TrackNumber, order.Entry, order.Locale,
					order.InternalSignature, order.CustomerID, order.DeliveryService, order.Shardkey,
					order.SmID, order.DateCreated, order.OofShard,
				).WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("INSERT INTO deliveries").WithArgs(
					order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
					order.Delivery.City, order.Delivery.Adress, order.Delivery.Region, order.Delivery.Email,
				).WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("INSERT INTO payments").WithArgs(
					order.OrderUID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency,
					order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDT, order.Payment.Bank,
					order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee,
				).WillReturnError(errors.New("some error"))

				mock.ExpectRollback()
			},
		},
		{
			"Items err test",
			true,
			func() {
				// Set up expectations for the mock
				mock.ExpectBegin()
				mock.ExpectExec("INSERT INTO orders").WithArgs(order.OrderUID, order.TrackNumber, order.Entry, order.Locale,
					order.InternalSignature, order.CustomerID, order.DeliveryService, order.Shardkey,
					order.SmID, order.DateCreated, order.OofShard,
				).WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("INSERT INTO deliveries").WithArgs(
					order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip,
					order.Delivery.City, order.Delivery.Adress, order.Delivery.Region, order.Delivery.Email,
				).WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectExec("INSERT INTO payments").WithArgs(
					order.OrderUID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency,
					order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDT, order.Payment.Bank,
					order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee,
				).WillReturnResult(sqlmock.NewResult(1, 1))

				for i := range order.Items {
					mock.ExpectExec("INSERT INTO items").WithArgs(
						order.OrderUID, order.Items[i].ChrtID, order.Items[i].TrackNumber, order.Items[i].Price,
						order.Items[i].Rid, order.Items[i].Name, order.Items[i].Sale, order.Items[i].Size,
						order.Items[i].TotalPrice, order.Items[i].NmID, order.Items[i].Brand, order.Items[i].Status,
					).WillReturnError(errors.New("some error"))
				}

				mock.ExpectRollback()
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockFunc()

			orderStorageUID, err := storage.CreateOrder(context.Background(), order)

			if tt.expectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, orderStorageUID, orderUID)
			}
			// we make sure that all expectations were met
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}

		})
	}
}
