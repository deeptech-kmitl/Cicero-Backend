package productUsecase

import (
	"github.com/deeptech-kmitl/Cicero-Backend/config"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/product"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/product/productRepository"
)

type IProductUsecase interface {
	FindOneProduct(prodId string) (*product.Product, error)
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

func (u *productUsecase) FindOneProduct(prodId string) (*product.Product, error) {
	result, err := u.productsRepository.FindOneProduct(prodId)
	if err != nil {
		return nil, err
	}
	return result, nil
}
