package servers

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"

	"github.com/deeptech-kmitl/Cicero-Backend/config"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type IServer interface {
	GetServer() *server
	Start()
}

type server struct {
	app *fiber.App
	cfg config.IConfig
	db  *sqlx.DB
}

func NewServer(cfg config.IConfig, db *sqlx.DB) IServer {
	return &server{
		db:  db,
		cfg: cfg,
		app: fiber.New(fiber.Config{
			AppName:      cfg.App().Name(),
			BodyLimit:    cfg.App().BodyLimit(),
			ReadTimeout:  cfg.App().ReadTimeout(),
			WriteTimeout: cfg.App().WriteTimeout(),
			JSONEncoder:  json.Marshal,
			JSONDecoder:  json.Unmarshal,
		}),
	}

}

func (s *server) Start() {
	// Middleware

	// Module
	api := s.app.Group("/api")

	modules := NewModule(api, s)
	modules.MonitorModule()

	//Graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		log.Println("server is shutting down...")
		_ = s.app.Shutdown()
	}()

	//Listen to host:port
	log.Printf("server is running at %v", s.cfg.App().Url())
	s.app.Listen(s.cfg.App().Url())

}

func (s *server) GetServer() *server {
	return s
}
