package orderRepository

import (
	"context"
	"fmt"
	"time"

	"github.com/deeptech-kmitl/Cicero-Backend/modules/order"
	"github.com/jmoiron/sqlx"
)

type IOrderRepository interface {
	AddOrder(req *order.AddOrderReq, products *order.OrderProducts) error
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

	if _, err := r.db.ExecContext(ctx, query, req.UserId, req.Total, req.Status, products, req.Address, req.PaymentDetail); err != nil {
		return fmt.Errorf("add order: %v", err)
	}

	return nil
}
