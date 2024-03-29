package orderUsecase

import (
	"fmt"

	"github.com/deeptech-kmitl/Cicero-Backend/config"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/order"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/order/orderRepository"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/users/usersRepositories"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/users/usersUsecases"
)

type IOrderUsecase interface {
	AddOrder(req *order.AddOrderReq) (string, error)
	GetOrderByUserId(userId string) []*order.GetOrderByUserId
	GetOneOrderById(orderId string) (*order.GetOneOrderById, error)
}

type orderUsecase struct {
	orderRepo   orderRepository.IOrderRepository
	userUsecase usersUsecases.IUserUsecase
}

func OrderUsecase(orderRepo orderRepository.IOrderRepository, userRepo usersRepositories.IUsersRepository, cfg config.IConfig) IOrderUsecase {
	return &orderUsecase{
		orderRepo:   orderRepo,
		userUsecase: usersUsecases.UserUsecase(userRepo, cfg),
	}
}

func (u *orderUsecase) AddOrder(req *order.AddOrderReq) (string, error) {

	productsOrder, err := u.userUsecase.GetCart(req.UserId)
	if err != nil {
		return "", err
	}

	if len(productsOrder) == 0 {
		return "", fmt.Errorf("cart is empty")

	}

	req.Status = "pending"

	orders := &order.OrderProducts{
		Products: productsOrder,
	}

	orderId, err := u.orderRepo.AddOrder(req, orders)
	if err != nil {
		return "", err
	}

	return orderId, nil
}

func (u *orderUsecase) GetOrderByUserId(userId string) []*order.GetOrderByUserId {
	return u.orderRepo.GetOrderByUserId(userId)
}

func (u *orderUsecase) GetOneOrderById(orderId string) (*order.GetOneOrderById, error) {
	return u.orderRepo.GetOneOrderById(orderId)
}
