package servers

import (
	"github.com/deeptech-kmitl/Cicero-Backend/modules/monitor/monitorHandlers"
	"github.com/gofiber/fiber/v2"
)

type IModuleFactory interface {
	MonitorModule()
}

type moduleFactory struct {
	r fiber.Router
	s *server
}

func NewModule(r fiber.Router, s *server) IModuleFactory {
	return &moduleFactory{
		r: r,
		s: s,
	}
}

func (m *moduleFactory) MonitorModule() {
	monitorHandler := monitorHandlers.MonitorHandler(m.s.cfg)
	m.r.Get("/", monitorHandler.HealthCheck)
}
