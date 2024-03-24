package orderRepository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/deeptech-kmitl/Cicero-Backend/modules/order"
	"github.com/jmoiron/sqlx"
)

type IOrderRepository interface {
	AddOrder(req *order.AddOrderReq, products *order.OrderProducts) error
	GetOrderByUserId(userId string) (*order.GetOrderByUserId, error)
}

type orderRepository struct {
	db *sqlx.DB
}

func OrderRepository(db *sqlx.DB) IOrderRepository {
	return &orderRepository{db}
}

func (r *orderRepository) AddOrder(req *order.AddOrderReq, products *order.OrderProducts) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction add order failed: %v", err)
	}

	query := `
	INSERT INTO "Order" (
		"user_id",
		"total",
		"status",
		"products",
		"address",
		"payment_detail"
	)
	VALUES ($1, $2, $3, $4, $5, $6);
	`

	if _, err := tx.ExecContext(ctx, query, req.UserId, req.Total, req.Status, products, req.Address, req.PaymentDetail); err != nil {
		tx.Rollback()
		return fmt.Errorf("add order: %v", err)
	}

	//delete all products in cart by user_id
	query = `
	DELETE FROM "Cart"
	WHERE "user_id" = $1;
	`

	if _, err := tx.ExecContext(ctx, query, req.UserId); err != nil {
		tx.Rollback()
		return fmt.Errorf("delete all products in cart by user_id: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit add order failed: %v", err)
	}
	return nil
}

func (r *orderRepository) GetOrderByUserId(userId string) (*order.GetOrderByUserId, error) {
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	query := `
	SELECT
    	to_jsonb("t")
	FROM 
	(
    SELECT
        "o"."id",
        "o"."user_id",
        "o"."total",
        "o"."status",
        "o"."products"
    FROM "Order" "o"
    WHERE "o"."user_id" = $1
	) AS "t";`

	bytes := make([]byte, 0)
	order := new(order.GetOrderByUserId)

	if err := r.db.Get(&bytes, query, userId); err != nil {
		return nil, fmt.Errorf("cannot get order by user id: %w", err)
	}

	if err := json.Unmarshal(bytes, &order); err != nil {
		return nil, fmt.Errorf("unmarshal order failed: %v", err)
	}

	return order, nil
}
