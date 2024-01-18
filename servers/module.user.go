package servers

import (
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
	repository := usersRepositories.UsersRepositoryHandler(m.s.db)
	usecase := usersUsecases.UserUsecaseHandler(repository, m.s.cfg)
	handler := usersHandlers.UsersHandler(m.s.cfg, usecase)
	return &userModule{
		moduleFactory: m,
		repository:    repository,
		usecase:       usecase,
		handler:       handler,
	}
}

func (m *userModule) Init() {
	router := m.r.Group("/users")

	router.Post("/signup", m.handler.SignUpCustomer)
	router.Post("/signup-admin", m.mid.JwtAuth(), m.mid.Authorize(2), m.handler.SignUpAdmin)
	router.Post("/signin", m.handler.SignIn)
	router.Post("/signout", m.handler.SignOut)
	router.Get("/:user_id", m.mid.JwtAuth(), m.mid.ParamsCheck(), m.handler.GetUserProfile)
}

func (p *userModule) Repository() usersRepositories.IUsersRepository { return p.repository }
func (p *userModule) Usecase() usersUsecases.IUserUsecase            { return p.usecase }
func (p *userModule) Handler() usersHandlers.IUsersHandler           { return p.handler }
