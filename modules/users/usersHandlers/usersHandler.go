package usersHandlers

import (
	"fmt"
	"math"
	"path/filepath"
	"strings"

	"github.com/deeptech-kmitl/Cicero-Backend/config"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/entities"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/files"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/files/filesUsecase"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/users"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/users/usersUsecases"
	"github.com/deeptech-kmitl/Cicero-Backend/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type userHandlerErrCode = string

const (
	signUpCustomerErr    userHandlerErrCode = "users-001"
	signInErr            userHandlerErrCode = "users-002"
	signOutErr           userHandlerErrCode = "users-003"
	signUpAdminErr       userHandlerErrCode = "users-004"
	getUserProfileErr    userHandlerErrCode = "users-005"
	updateUserProfileErr userHandlerErrCode = "users-006"
	WishlistErr          userHandlerErrCode = "users-007"
	GetWishlistErr       userHandlerErrCode = "users-008"
	AddCartErr           userHandlerErrCode = "users-009"
	RemoveCartErr        userHandlerErrCode = "users-010"
	GetCartErr           userHandlerErrCode = "users-011"
	DecreaseQtyCartErr   userHandlerErrCode = "users-012"
	IncreaseQtyCartErr   userHandlerErrCode = "users-013"
	UpdateSizeCartErr    userHandlerErrCode = "users-014"
)

type IUsersHandler interface {
	SignUpCustomer(c *fiber.Ctx) error
	SignUpAdmin(c *fiber.Ctx) error
	SignIn(c *fiber.Ctx) error
	SignOut(c *fiber.Ctx) error
	GetUserProfile(c *fiber.Ctx) error
	UpdateUserProfile(c *fiber.Ctx) error
	Wishlist(c *fiber.Ctx) error
	GetWishlist(c *fiber.Ctx) error
	AddCart(c *fiber.Ctx) error
	RemoveCart(c *fiber.Ctx) error
	GetCart(c *fiber.Ctx) error
	DecreaseQtyCart(c *fiber.Ctx) error
	IncreaseQtyCart(c *fiber.Ctx) error
	UpdateSizeCart(c *fiber.Ctx) error
}

type usersHandler struct {
	cfg         config.IConfig
	userUsecase usersUsecases.IUserUsecase
	fileUsecase filesUsecase.IFilesUsecase
}

func UsersHandler(cfg config.IConfig, UserUsecase usersUsecases.IUserUsecase, fileUsecase filesUsecase.IFilesUsecase) IUsersHandler {
	return &usersHandler{
		cfg:         cfg,
		userUsecase: UserUsecase,
		fileUsecase: fileUsecase,
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
		switch err.Error() {
		case "phone number has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(signUpCustomerErr),
				err.Error(),
			).Res()
		case "email has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(signUpCustomerErr),
				err.Error(),
			).Res()

		default:
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(signUpCustomerErr),
				err.Error(),
			).Res()
		}
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
		switch err.Error() {
		case "phone number has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(signUpAdminErr),
				err.Error(),
			).Res()
		case "email has been used":
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(signUpAdminErr),
				err.Error(),
			).Res()

		default:
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(signUpCustomerErr),
				err.Error(),
			).Res()
		}
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

func (h *usersHandler) UpdateUserProfile(c *fiber.Ctx) error {
	avatarFile := make([]*files.FileReq, 0)
	userId := strings.Trim(c.Params("user_id"), " ")

	form, err := c.MultipartForm()
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(updateUserProfileErr),
			err.Error(),
		).Res()
	}

	email := ""
	if values, exists := form.Value["email"]; exists && len(values) > 0 {
		email = values[0]
	}

	lastName := ""
	if values, exists := form.Value["lname"]; exists && len(values) > 0 {
		lastName = values[0]
	}

	firstName := ""
	if values, exists := form.Value["fname"]; exists && len(values) > 0 {
		firstName = values[0]
	}

	phone := ""
	if values, exists := form.Value["phone"]; exists && len(values) > 0 {
		phone = values[0]
	}

	// avatar := make([]*multipart.FileHeader, 0)
	avatarUrl := ""
	if avatar, exists := form.File["avatar"]; exists {
		// avatar = file

		if len(avatar) > 1 {
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(updateUserProfileErr),
				"avatar must be one file",
			).Res()
		}

		// // files ext validation
		extMap := map[string]string{
			"png":  "png",
			"jpg":  "jpg",
			"jpeg": "jpeg",
		}

		// check file extension
		ext := strings.TrimPrefix(filepath.Ext(avatar[0].Filename), ".")
		if extMap[ext] != ext || extMap[ext] == "" {
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(updateUserProfileErr),
				"invalid file extension",
			).Res()
		}
		// 	// check file size
		if avatar[0].Size > int64(h.cfg.App().FileLimit()) {
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(updateUserProfileErr),
				fmt.Sprintf("file size must less than %d MB", int(math.Ceil(float64(h.cfg.App().FileLimit())/math.Pow(1024, 2)))),
			).Res()
		}

		filename := utils.RandFileName(ext)
		avatarFile = append(avatarFile, &files.FileReq{
			File:        avatar[0],
			Destination: fmt.Sprintf("%s/%s", userId, filename),
			FileName:    filename,
			Extension:   ext,
		})

		result, err := h.fileUsecase.UploadToGCP(avatarFile)
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(updateUserProfileErr),
				err.Error(),
			).Res()
		}

		avatarUrl = result[0].Url
	}

	req := &users.UserUpdate{
		Id:        userId,
		Avatar:    avatarUrl,
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Phone:     phone,
	}

	res, err := h.userUsecase.UpdateUserProfile(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(updateUserProfileErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, res).Res()
}

func (h *usersHandler) Wishlist(c *fiber.Ctx) error {
	userId := strings.Trim(c.Params("user_id"), " ")
	prodId := strings.Trim(c.Params("product_id"), " ")

	result, err := h.userUsecase.Wishlist(userId, prodId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(WishlistErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, result).Res()

}

func (h *usersHandler) GetWishlist(c *fiber.Ctx) error {
	userId := strings.Trim(c.Params("user_id"), " ")

	result, err := h.userUsecase.GetWishlist(userId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(GetWishlistErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, result).Res()

}

func (h *usersHandler) AddCart(c *fiber.Ctx) error {
	req := new(users.AddCartReq)

	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(AddCartErr),
			err.Error(),
		).Res()
	}

	result, err := h.userUsecase.AddCart(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(AddCartErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, result).Res()
}

func (h *usersHandler) RemoveCart(c *fiber.Ctx) error {
	userId := strings.Trim(c.Params("user_id"), " ")
	prodId := strings.Trim(c.Params("product_id"), " ")

	result, err := h.userUsecase.RemoveCart(userId, prodId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(RemoveCartErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, result).Res()
}

func (h *usersHandler) GetCart(c *fiber.Ctx) error {
	userId := strings.Trim(c.Params("user_id"), " ")

	result, err := h.userUsecase.GetCart(userId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(GetCartErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, result).Res()
}

func (h *usersHandler) DecreaseQtyCart(c *fiber.Ctx) error {
	userId := strings.Trim(c.Params("user_id"), " ")
	prodId := strings.Trim(c.Params("product_id"), " ")

	qty, err := h.userUsecase.DecreaseQtyCart(userId, prodId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(DecreaseQtyCartErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, qty).Res()
}

func (h *usersHandler) IncreaseQtyCart(c *fiber.Ctx) error {
	userId := strings.Trim(c.Params("user_id"), " ")
	prodId := strings.Trim(c.Params("product_id"), " ")

	qty, err := h.userUsecase.IncreaseQtyCart(userId, prodId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(IncreaseQtyCartErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, qty).Res()
}

func (h *usersHandler) UpdateSizeCart(c *fiber.Ctx) error {
	req := new(users.UpdateSizeReq)

	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(UpdateSizeCartErr),
			err.Error(),
		).Res()
	}

	size, err := h.userUsecase.UpdateSizeCart(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(UpdateSizeCartErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, size).Res()
}
