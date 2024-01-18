package usersHandlers

import (
	"strings"

	"github.com/deeptech-kmitl/Cicero-Backend/config"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/entities"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/users"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/users/usersUsecases"
	"github.com/gofiber/fiber/v2"
)

type userHandlerErrCode = string

const (
	signUpCustomerErr     userHandlerErrCode = "users-001"
	signInErr             userHandlerErrCode = "users-002"
	refreshPassportErr    userHandlerErrCode = "users-003"
	signOutErr            userHandlerErrCode = "users-004"
	signUpAdminErr        userHandlerErrCode = "users-005"
	generateAdminTokenErr userHandlerErrCode = "users-006"
	getUserProfileErr     userHandlerErrCode = "users-007"
)

type IUsersHandler interface {
	SignUpCustomer(c *fiber.Ctx) error
	SignUpAdmin(c *fiber.Ctx) error
	SignIn(c *fiber.Ctx) error
	SignOut(c *fiber.Ctx) error
	GetUserProfile(c *fiber.Ctx) error
}

type usersHandler struct {
	cfg         config.IConfig
	userUsecase usersUsecases.IUserUsecase
}

func UsersHandler(cfg config.IConfig, UserUsecase usersUsecases.IUserUsecase) IUsersHandler {
	return &usersHandler{
		cfg:         cfg,
		userUsecase: UserUsecase,
	}
}

func (h *usersHandler) SignUpCustomer(c *fiber.Ctx) error {
	// Request Body parser
	req := new(users.UserRegisterReq)

	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpCustomerErr),
			err.Error(),
		).Res()
	}
	// Email validation
	if !req.IsEmail() {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpCustomerErr),
			"email is invalid",
		).Res()
	}

	// Insert user
	result, err := h.userUsecase.InsertCustomer(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpCustomerErr),
			err.Error(),
		).Res()
		// switch err.Error() {
		// case "username has been used":
		// 	return entities.NewResponse(c).Error(
		// 		fiber.ErrBadRequest.Code,
		// 		string(signUpCustomerErr),
		// 		err.Error(),
		// 	).Res()
		// case "email has been used":
		// 	return entities.NewResponse(c).Error(
		// 		fiber.ErrBadRequest.Code,
		// 		string(signUpCustomerErr),
		// 		err.Error(),
		// 	).Res()

		// default:
		// 	return entities.NewResponse(c).Error(
		// 		fiber.ErrInternalServerError.Code,
		// 		string(signUpCustomerErr),
		// 		err.Error(),
		// 	).Res()
		// }
	}

	return entities.NewResponse(c).Success(fiber.StatusCreated, result).Res()
}

func (h *usersHandler) SignUpAdmin(c *fiber.Ctx) error {
	// Request Body parser
	req := new(users.UserRegisterReq)

	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpAdminErr),
			err.Error(),
		).Res()
	}
	// Email validation
	if !req.IsEmail() {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpAdminErr),
			"email is invalid",
		).Res()
	}

	// Insert user
	result, err := h.userUsecase.InsertAdmin(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpAdminErr),
			err.Error(),
		).Res()
		// switch err.Error() {
		// case "username has been used":
		// 	return entities.NewResponse(c).Error(
		// 		fiber.ErrBadRequest.Code,
		// 		string(signUpAdminErr),
		// 		err.Error(),
		// 	).Res()
		// case "email has been used":
		// 	return entities.NewResponse(c).Error(
		// 		fiber.ErrBadRequest.Code,
		// 		string(signUpAdminErr),
		// 		err.Error(),
		// 	).Res()

		// default:
		// 	return entities.NewResponse(c).Error(
		// 		fiber.ErrInternalServerError.Code,
		// 		string(signUpCustomerErr),
		// 		err.Error(),
		// 	).Res()
		// }
	}

	return entities.NewResponse(c).Success(fiber.StatusCreated, result).Res()
}

func (h *usersHandler) SignIn(c *fiber.Ctx) error {
	req := new(users.UserCredential)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signInErr),
			err.Error(),
		).Res()
	}

	result, err := h.userUsecase.GetPassport(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signInErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, result).Res()
}

func (h *usersHandler) SignOut(c *fiber.Ctx) error {
	req := new(users.UserRemoveCredential)

	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signOutErr),
			err.Error(),
		).Res()
	}

	if err := h.userUsecase.DeleteOauth(req.OauthId); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signOutErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, nil).Res()
}

func (h *usersHandler) GetUserProfile(c *fiber.Ctx) error {
	userId := strings.Trim(c.Params("user_id"), " ")

	result, err := h.userUsecase.GetUserProfile(userId)
	if err != nil {
		switch err.Error() {
		case "get user failed: sql: no rows in result set":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(getUserProfileErr),
				err.Error(),
			).Res()
		default:
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(getUserProfileErr),
				err.Error(),
			).Res()

		}

	}

	return entities.NewResponse(c).Success(fiber.StatusOK, result).Res()
}
