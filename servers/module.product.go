package servers

import (
	"github.com/deeptech-kmitl/Cicero-Backend/modules/files/filesUsecase"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/product/productHandler"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/product/productRepository"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/product/productUsecase"
)

type IProductModule interface {
	Init()
	Repository() productRepository.IProductRepository
	Usecase() productUsecase.IProductUsecase
	Handler() productHandler.IProductHandler
}

type productModule struct {
	*moduleFactory
	repo    productRepository.IProductRepository
	usecase productUsecase.IProductUsecase
	handler productHandler.IProductHandler
}

func (m *moduleFactory) ProductModule() IProductModule {
	fileUsecase := filesUsecase.FilesUsecase(m.s.cfg)
	repo := productRepository.ProductRepository(m.s.db)
	usecase := productUsecase.ProductUsecase(repo, m.s.cfg)
	handler := productHandler.ProductHandler(usecase, fileUsecase, m.s.cfg)

	return &productModule{
		moduleFactory: m,
		usecase:       usecase,
		handler:       handler,
	}
}

func (m *productModule) Init() {
	router := m.r.Group("/product")

	router.Get("/search", m.handler.FindProduct)
	router.Get("/:product_id", m.handler.FindOneProduct)
	router.Post("/", m.mid.JwtAuth(), m.mid.Authorize(2), m.handler.AddProduct)
	router.Delete("/:product_id", m.mid.JwtAuth(), m.mid.Authorize(2), m.handler.DeleteProduct)
	router.Put("/", m.mid.JwtAuth(), m.mid.Authorize(2), m.handler.UpdateProduct)
	router.Get("/image/:product_id", m.handler.FindImageByProductId)

}

func (f *productModule) Repository() productRepository.IProductRepository { return f.repo }
func (f *productModule) Usecase() productUsecase.IProductUsecase          { return f.usecase }
func (f *productModule) Handler() productHandler.IProductHandler          { return f.handler }
