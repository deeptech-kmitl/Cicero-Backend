package servers

import (
	"github.com/deeptech-kmitl/Cicero-Backend/modules/middlewares/middlewareHandler"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/middlewares/middlewareRepository"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/middlewares/middlewareUsecase"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/monitor/monitorHandlers"
	"github.com/gofiber/fiber/v2"
)

type IModuleFactory interface {
	MonitorModule()
	UserModule() IUserModule
}

type moduleFactory struct {
	r   fiber.Router
	s   *server
	mid middlewareHandler.IMiddlewaresHandler
}

func NewModule(r fiber.Router, s *server, mid middlewareHandler.IMiddlewaresHandler) IModuleFactory {
	return &moduleFactory{
		r:   r,
		s:   s,
		mid: mid,
	}
}

func (m *moduleFactory) MonitorModule() {
	monitorHandler := monitorHandlers.MonitorHandler(m.s.cfg)
	m.r.Get("/", monitorHandler.HealthCheck)
}

func InitMiddlewares(s *server) middlewareHandler.IMiddlewaresHandler {
	repository := middlewareRepository.MiddlewaresRepository(s.db)
	usecase := middlewareUsecase.MiddlewaresUsecase(repository)
	return middlewareHandler.MiddlewaresHandler(s.cfg, usecase)
}
