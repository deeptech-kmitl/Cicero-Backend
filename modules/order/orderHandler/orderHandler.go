package orderHandler

import (
	"strings"

	"github.com/deeptech-kmitl/Cicero-Backend/modules/entities"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/order"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/order/orderUsecase"
	"github.com/gofiber/fiber/v2"
)

type orderHandlerErrCode = string

const (
	addOrderErr         orderHandlerErrCode = "order-001"
	getOrderByUserIdErr orderHandlerErrCode = "order-002"
	getOneOrderByIdErr  orderHandlerErrCode = "order-003"
)

type IOrderHandler interface {
	AddOrder(c *fiber.Ctx) error
	GetOrderByUserId(c *fiber.Ctx) error
	GetOneOrderById(c *fiber.Ctx) error
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
	orderId, err := h.orderUsecase.AddOrder(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(addOrderErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, orderId).Res()
}

func (h *orderHandler) GetOrderByUserId(c *fiber.Ctx) error {
	userId := strings.TrimSpace(c.Params("user_id"))

	//get order by user id
	order := h.orderUsecase.GetOrderByUserId(userId)

	return entities.NewResponse(c).Success(fiber.StatusOK, order).Res()
}

func (h *orderHandler) GetOneOrderById(c *fiber.Ctx) error {
	orderId := strings.TrimSpace(c.Params("order_id"))

	//get order by order id
	order, err := h.orderUsecase.GetOneOrderById(orderId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(getOneOrderByIdErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, order).Res()
}
