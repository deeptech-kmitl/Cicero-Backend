package middlewareHandler

import (
	"fmt"
	"strings"

	"github.com/deeptech-kmitl/Cicero-Backend/config"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/entities"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/middlewares/middlewareUsecase"
	"github.com/deeptech-kmitl/Cicero-Backend/pkg/auth"
	"github.com/deeptech-kmitl/Cicero-Backend/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type middlewareHandlersErrCode string

const (
	routerCheckErr middlewareHandlersErrCode = "middleware-001"
	jwtAuthErr     middlewareHandlersErrCode = "middleware-002"
	paramsCheckErr middlewareHandlersErrCode = "middleware-003"
	authorizeErr   middlewareHandlersErrCode = "middleware-004"
)

type IMiddlewaresHandler interface {
	Cors() fiber.Handler
	RouterCheck() fiber.Handler
	Logger() fiber.Handler
	JwtAuth() fiber.Handler
	ParamsCheck() fiber.Handler
	Authorize(expectRoleId ...int) fiber.Handler
}

type middlewaresHandler struct {
	cfg               config.IConfig
	middlewareUsecase middlewareUsecase.IMiddlewaresUsecase
}

func MiddlewaresHandler(cfg config.IConfig, usecase middlewareUsecase.IMiddlewaresUsecase) IMiddlewaresHandler {
	return &middlewaresHandler{
		cfg:               cfg,
		middlewareUsecase: usecase,
	}
}

// กำหนด cors ให้กับ api
func (h *middlewaresHandler) Cors() fiber.Handler {
	return cors.New(cors.Config{
		Next:             cors.ConfigDefault.Next,
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowHeaders:     "",
		AllowCredentials: false,
		ExposeHeaders:    "",
		MaxAge:           0,
	})
}

// ตรวจสอบว่ามี router นี้หรือไม่
func (h *middlewaresHandler) RouterCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return entities.NewResponse(c).Error(
			fiber.ErrNotFound.Code,
			string(routerCheckErr),
			"router not found",
		).Res()
	}
}

func (h *middlewaresHandler) Logger() fiber.Handler {
	return logger.New(logger.Config{
		Format:     "${time} [${ip}] ${status} - ${method} ${path}\n",
		TimeFormat: "02/01/2006",
		TimeZone:   "Bangkok/Asia",
	})
}

// แกะ token และตรวจสอบว่า Login อยู่หรือไม่
func (h *middlewaresHandler) JwtAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := strings.TrimPrefix(c.Get("Authorization"), "Bearer ")
		result, err := auth.ParseToken(h.cfg.Jwt(), token)
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(jwtAuthErr),
				err.Error(),
			).Res()
		}

		claims := result.Claims
		// check token in db
		check := h.middlewareUsecase.FindAccessToken(claims.Id, token)
		if !check {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(jwtAuthErr),
				"You Are Not Logged In",
			).Res()
		}

		// set UserId and roleId to locals
		c.Locals("userId", claims.Id)
		c.Locals("userRoleId", claims.RoleId)
		return c.Next()
	}
}

// ป้องกันการเข้าถึงข้อมูลของคนอื่น ต้องมาคู่กับ JwtAuth
func (h *middlewaresHandler) ParamsCheck() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Locals("userId")
		fmt.Println(userId, c.Params("user_id"))
		if c.Locals("userRoleId").(int) == 2 {
			return c.Next()
		}
		if c.Params("user_id") != userId {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(paramsCheckErr),
				"no permission to access",
			).Res()
		}
		return c.Next()
	}

}

// ตรวจสอบว่า user มีสิทธิ์เข้าถึง api นี้หรือไม่ ต้องมาคู่กับ JwtAuth
func (h *middlewaresHandler) Authorize(expectRoleId ...int) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRoleId, ok := c.Locals("userRoleId").(int)
		if !ok {
			return entities.NewResponse(c).Error(
				fiber.ErrUnauthorized.Code,
				string(authorizeErr),
				"user role id is not int",
			).Res()
		}

		// หาจำนวน role ทั้งหมดในระบบ
		roles, err := h.middlewareUsecase.FindRole()
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrInternalServerError.Code,
				string(authorizeErr),
				err.Error(),
			).Res()
		}

		sum := 0
		for _, v := range expectRoleId {
			sum += v
		}

		expectedValueBinary := utils.BinaryConverter(sum, len(roles))
		userValueBinary := utils.BinaryConverter(userRoleId, len(roles))
		// value in db of higher role must be greater than lower role
		//  1 = 01 => customer = 1 (in db) ทำได้คนเดียว
		//  2 = 10 => admin = 2 (in db) ทำได้คนเดียว
		//  1,2 = 3 = 11 => ทำได้ทั้งสองคน

		// examnple
		// userRoleId =          0 1 => customer
		// expectedValueBinary = 1 0 => admin
		// 0 != 1 && 1 != 0 => false

		for i := range userValueBinary {
			if userValueBinary[i] == 1 && expectedValueBinary[i] == 1 {
				return c.Next()
			}
		}

		return entities.NewResponse(c).Error(
			fiber.ErrUnauthorized.Code,
			string(authorizeErr),
			"no permission to access",
		).Res()
	}
}
