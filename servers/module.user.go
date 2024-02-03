package servers

import (
	"github.com/deeptech-kmitl/Cicero-Backend/modules/files/filesUsecase"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/users/usersHandlers"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/users/usersRepositories"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/users/usersUsecases"
)

type IUserModule interface {
	Init()
	Repository() usersRepositories.IUsersRepository
	Usecase() usersUsecases.IUserUsecase
	Handler() usersHandlers.IUsersHandler
}

type userModule struct {
	*moduleFactory
	repository usersRepositories.IUsersRepository
	usecase    usersUsecases.IUserUsecase
	handler    usersHandlers.IUsersHandler
}

func (m *moduleFactory) UserModule() IUserModule {
	fileUsecase := filesUsecase.FilesUsecase(m.s.cfg)
	userRepository := usersRepositories.UsersRepository(m.s.db)
	userUsecase := usersUsecases.UserUsecase(userRepository, m.s.cfg)
	userHandler := usersHandlers.UsersHandler(m.s.cfg, userUsecase, fileUsecase)
	return &userModule{
		moduleFactory: m,
		repository:    userRepository,
		usecase:       userUsecase,
		handler:       userHandler,
	}
}

func (m *userModule) Init() {
	router := m.r.Group("/users")

	router.Post("/signup", m.handler.SignUpCustomer)
	router.Post("/signup-admin", m.mid.JwtAuth(), m.mid.Authorize(2), m.handler.SignUpAdmin)
	router.Post("/signin", m.handler.SignIn)
	router.Post("/signout", m.mid.JwtAuth(), m.handler.SignOut)
	router.Get("/:user_id", m.mid.JwtAuth(), m.mid.ParamsCheck(), m.handler.GetUserProfile)
	router.Put("/:user_id", m.mid.JwtAuth(), m.mid.ParamsCheck(), m.handler.UpdateUserProfile)
	router.Post("/:user_id/wishlist/:product_id", m.mid.JwtAuth(), m.mid.ParamsCheck(), m.handler.Wishlist)
	router.Get("/wishlist/:user_id", m.mid.JwtAuth(), m.mid.ParamsCheck(), m.handler.GetWishlist)
	router.Post("/cart/:user_id", m.mid.JwtAuth(), m.mid.ParamsCheck(), m.handler.AddCart)
	router.Delete("/cart/:user_id/:product_id", m.mid.JwtAuth(), m.mid.ParamsCheck(), m.handler.RemoveCart)
	router.Get("/cart/:user_id", m.mid.JwtAuth(), m.mid.ParamsCheck(), m.handler.GetCart)
}

func (p *userModule) Repository() usersRepositories.IUsersRepository { return p.repository }
func (p *userModule) Usecase() usersUsecases.IUserUsecase            { return p.usecase }
func (p *userModule) Handler() usersHandlers.IUsersHandler           { return p.handler }
