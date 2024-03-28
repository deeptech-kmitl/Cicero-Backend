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
	AddOrder(req *order.AddOrderReq, products *order.OrderProducts) (string, error)
	GetOrderByUserId(userId string) []*order.GetOrderByUserId
	GetOneOrderById(orderId string) (*order.GetOneOrderById, error)
}

type orderRepository struct {
	db *sqlx.DB
}

func OrderRepository(db *sqlx.DB) IOrderRepository {
	return &orderRepository{db}
}

func (r *orderRepository) AddOrder(req *order.AddOrderReq, products *order.OrderProducts) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("begin transaction add order failed: %v", err)
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
	VALUES ($1, $2, $3, $4, $5, $6)
	RETURNING "id";
	`

	var orderId string

	if err := tx.QueryRowxContext(ctx, query, req.UserId, req.Total, req.Status, products, req.Address, req.PaymentDetail).Scan(&orderId); err != nil {
		tx.Rollback()
		return "", fmt.Errorf("add order: %v", err)
	}

	//delete all products in cart by user_id
	query = `
	DELETE FROM "Cart"
	WHERE "user_id" = $1;
	`

	if _, err := tx.ExecContext(ctx, query, req.UserId); err != nil {
		tx.Rollback()
		return "", fmt.Errorf("delete all products in cart by user_id: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return "", fmt.Errorf("commit add order failed: %v", err)
	}
	return orderId, nil
}

func (r *orderRepository) GetOrderByUserId(userId string) []*order.GetOrderByUserId {
	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	query := `
	SELECT
		COALESCE(array_to_json(array_agg("t")), '[]'::json)
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
	orderData := make([]*order.GetOrderByUserId, 0)

	if err := r.db.Get(&bytes, query, userId); err != nil {
		return make([]*order.GetOrderByUserId, 0)
	}

	if err := json.Unmarshal(bytes, &orderData); err != nil {
		return make([]*order.GetOrderByUserId, 0)
	}

	return orderData
}

func (r *orderRepository) GetOneOrderById(orderId string) (*order.GetOneOrderById, error) {
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
			"o"."products",
			"o"."address",
			"o"."payment_detail",
			"o"."created_at"
		FROM "Order" "o"
		WHERE "o"."id" = $1
	) AS "t";`

	orderBytes := make([]byte, 0)
	orderData := new(order.GetOneOrderById)
	if err := r.db.Get(&orderBytes, query, orderId); err != nil {
		switch err.Error() {
		case "sql: no rows in result set":
			return nil, fmt.Errorf("order not found: %v", err)
		default:
			return nil, fmt.Errorf("get one order failed: %v", err)
		}
	}
	if err := json.Unmarshal(orderBytes, &orderData); err != nil {
		return nil, fmt.Errorf("unmarshal order failed: %v", err)
	}

	return orderData, nil
}
