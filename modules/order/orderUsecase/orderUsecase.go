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
	AddOrder(req *order.AddOrderReq) error
	GetOrderByUserId(userId string) []*order.GetOrderByUserId
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

func (u *orderUsecase) AddOrder(req *order.AddOrderReq) error {

	productsOrder, err := u.userUsecase.GetCart(req.UserId)
	if err != nil {
		return err
	}

	if len(productsOrder) == 0 {
		return fmt.Errorf("cart is empty")

	}

	req.Status = "pending"

	orders := &order.OrderProducts{
		Products: productsOrder,
	}

	if err := u.orderRepo.AddOrder(req, orders); err != nil {
		return err
	}

	return nil
}

func (u *orderUsecase) GetOrderByUserId(userId string) []*order.GetOrderByUserId {
	return u.orderRepo.GetOrderByUserId(userId)
}
