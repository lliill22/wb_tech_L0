package storage

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
	"wb_tech_L0/internal/config"

	"github.com/jackc/pgx/v5"
)

const (
	orderFields    = `order_uid, track_number, entry, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, date_created, oof_shard`
	deliveryFields = `order_uid, name, phone, zip, city, address, region, email`
	paymentFields  = `order_uid, transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee`
	itemsFields    = `order_uid, chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status`
)

type Database struct {
	DB *pgx.Conn
}

func NewDB(cfg config.Postgres) (*Database, error) {
	var conn *pgx.Conn
	var err error
	retries := 10
	// postgres://postgres:postgres@postgres:5432/postgres
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	for i := range retries {
		conn, err = pgx.Connect(context.Background(), connStr)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to db %d/%d: %v", i+1, retries, err)
		time.Sleep(2 * time.Second)
	}

	db := &Database{DB: conn}
	if err != nil {
		db.DB.Close(context.Background())
		return nil, err
	}

	return db, nil
}

func SaveOrder(ctx context.Context, db *sql.DB, order Order) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("start transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			log.Println("transaction rollback due to error:", err)
			_ = tx.Rollback()
		}
	}()

	ordersQuerry := `INSERT INTO orders (` + orderFields + `) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
	// Вставка в orders
	_, err = tx.ExecContext(ctx, ordersQuerry, order.OrderUID, order.TrackNumber, order.Entry, order.Locale, order.InternalSignature,
		order.CustomerID, order.DeliveryService, order.Shardkey, order.SmID, order.DateCreated, order.OofShard)
	if err != nil {
		return fmt.Errorf("insert orders: %w", err)
	}

	deliveryQuerry := `INSERT INTO delivery (` + deliveryFields + `) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`
	// Вставка в delivery
	_, err = tx.ExecContext(ctx, deliveryQuerry,
		order.OrderUID, order.Delivery.Name, order.Delivery.Phone, order.Delivery.Zip, order.Delivery.City,
		order.Delivery.Address, order.Delivery.Region, order.Delivery.Email)
	if err != nil {
		return fmt.Errorf("insert delivery: %w", err)
	}

	paymentQuerry := `INSERT INTO payment (` + paymentFields + `) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`

	// Вставка в payment
	_, err = tx.ExecContext(ctx, paymentQuerry, order.OrderUID, order.Payment.Transaction, order.Payment.RequestID, order.Payment.Currency, order.Payment.Provider, order.Payment.Amount, order.Payment.PaymentDt, order.Payment.Bank, order.Payment.DeliveryCost, order.Payment.GoodsTotal, order.Payment.CustomFee)
	if err != nil {
		return fmt.Errorf("insert payment: %w", err)
	}

	itemsQuerry := `INSERT INTO items (` + itemsFields + `) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12)`
	// Вставка в items
	for _, item := range order.Items {
		_, err = tx.ExecContext(ctx, itemsQuerry, order.OrderUID, item.ChrtID, item.TrackNumber, item.Price, item.Rid, item.Name,
			item.Sale, item.Size, item.TotalPrice, item.NmID, item.Brand, item.Status)
		if err != nil {
			return fmt.Errorf("insert item: %w", err)
		}
	}

	if err = tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	log.Println("Order saved successfully:", order.OrderUID)
	return nil
}
