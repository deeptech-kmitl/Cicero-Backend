package orderHandler

import (
	"github.com/deeptech-kmitl/Cicero-Backend/modules/entities"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/order"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/order/orderUsecase"
	"github.com/gofiber/fiber/v2"
)

type orderHandlerErrCode = string

const (
	addOrderErr orderHandlerErrCode = "order-001"
)

type IOrderHandler interface {
	AddOrder(c *fiber.Ctx) error
}

type orderHandler struct {
	orderUsecase orderUsecase.IOrderUsecase
}

func OrderHandler(orderUsecase orderUsecase.IOrderUsecase) IOrderHandler {
	return &orderHandler{orderUsecase}
}

func (h *orderHandler) AddOrder(c *fiber.Ctx) error {

	//bodyparser
	req := new(order.AddOrderReq)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(addOrderErr),
			err.Error(),
		).Res()
	}

	//add order
	if err := h.orderUsecase.AddOrder(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(addOrderErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, "add order success").Res()
}
