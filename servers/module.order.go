package servers

import (
	"github.com/deeptech-kmitl/Cicero-Backend/modules/order/orderHandler"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/order/orderRepository"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/order/orderUsecase"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/users/usersRepositories"
)

type IOrderModule interface {
	Init()
	Repository() orderRepository.IOrderRepository
	Usecase() orderUsecase.IOrderUsecase
	Handler() orderHandler.IOrderHandler
}

type orderModule struct {
	*moduleFactory
	repository orderRepository.IOrderRepository
	usecase    orderUsecase.IOrderUsecase
	handler    orderHandler.IOrderHandler
}

func (m *moduleFactory) OrderModule() IOrderModule {
	userRepository := usersRepositories.UsersRepository(m.s.db)
	orderRepository := orderRepository.OrderRepository(m.s.db)
	orderUsecase := orderUsecase.OrderUsecase(orderRepository, userRepository, m.s.cfg)
	orderHandler := orderHandler.OrderHandler(orderUsecase)
	return &orderModule{
		moduleFactory: m,
		repository:    orderRepository,
		usecase:       orderUsecase,
		handler:       orderHandler,
	}
}

func (m *orderModule) Init() {
	router := m.r.Group("/order")

	router.Post("/", m.mid.JwtAuth(), m.mid.Authorize(1), m.handler.AddOrder)
	router.Get("/:user_id", m.mid.JwtAuth(), m.mid.ParamsCheck(), m.mid.Authorize(1), m.handler.GetOrderByUserId)

}

func (p *orderModule) Repository() orderRepository.IOrderRepository { return p.repository }
func (p *orderModule) Usecase() orderUsecase.IOrderUsecase          { return p.usecase }
func (p *orderModule) Handler() orderHandler.IOrderHandler          { return p.handler }
