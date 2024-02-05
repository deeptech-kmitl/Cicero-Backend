package productUsecase

import (
	"math"

	"github.com/deeptech-kmitl/Cicero-Backend/config"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/entities"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/product"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/product/productRepository"
)

type IProductUsecase interface {
	FindOneProduct(prodId string) (*product.Product, error)
	AddProduct(req *product.AddProduct) (*product.Product, error)
	DeleteProduct(prodId string) (string, error)
	FindProduct(req *product.ProductFilter) *entities.PaginateRes
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

func (u *productUsecase) FindProduct(req *product.ProductFilter) *entities.PaginateRes {
	if req.Page < 1 {
		req.Page = 1
	}
	if req.Limit < 3 {
		req.Limit = 3
	}

	if req.OrderBy == "" {
		req.OrderBy = "title"
	}
	if req.Sort == "" {
		req.Sort = "ASC"
	}

	products, count := u.productsRepository.FindProduct(req)
	return &entities.PaginateRes{
		Data:      products,
		TotalItem: count,
		Page:      req.Page,
		Limit:     req.Limit,
		TotalPage: int(math.Ceil(float64(count) / float64(req.Limit))),
	}

}

func (u *productUsecase) AddProduct(req *product.AddProduct) (*product.Product, error) {
	product, err := u.productsRepository.InsertProduct(req)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (u *productUsecase) DeleteProduct(prodId string) (string, error) {
	if err := u.productsRepository.DeleteProduct(prodId); err != nil {
		return "", err
	}

	return "Product deleted", nil
}
