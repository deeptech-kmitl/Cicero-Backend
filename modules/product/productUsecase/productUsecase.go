package productUsecase

import (
	"github.com/deeptech-kmitl/Cicero-Backend/config"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/product/productRepository"
)

type IProductUsecase interface {
}

type productUsecase struct {
	cfg                config.IConfig
	productsRepository productRepository.IProductRepository
}

func ProductUsecase(productsRepository productRepository.IProductRepository, cfg config.IConfig) IProductUsecase {
	return &productUsecase{
		productsRepository: productsRepository,
		cfg:                cfg,
	}
}
