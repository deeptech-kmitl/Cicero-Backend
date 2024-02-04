package productHandler

import (
	"github.com/deeptech-kmitl/Cicero-Backend/config"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/entities"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/product/productUsecase"
	"github.com/gofiber/fiber/v2"
)

type productHandlerErrCode = string

const (
	FindOneProductErr productHandlerErrCode = "product-001"
)

type IProductHandler interface {
	FindOneProduct(c *fiber.Ctx) error
}

type productHandler struct {
	cfg            config.IConfig
	productUsecase productUsecase.IProductUsecase
}

func ProductHandler(productUsecase productUsecase.IProductUsecase, cfg config.IConfig) IProductHandler {
	return &productHandler{
		productUsecase: productUsecase,
		cfg:            cfg,
	}
}

func (h *productHandler) FindOneProduct(c *fiber.Ctx) error {
	prodId := c.Params("product_id")
	result, err := h.productUsecase.FindOneProduct(prodId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(FindOneProductErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, result).Res()

}
