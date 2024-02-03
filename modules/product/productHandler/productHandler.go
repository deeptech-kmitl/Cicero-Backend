package productHandler

import (
	"github.com/deeptech-kmitl/Cicero-Backend/config"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/product/productUsecase"
)

type IProductHandler interface {
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
