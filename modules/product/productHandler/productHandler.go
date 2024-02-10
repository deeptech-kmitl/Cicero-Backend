package productHandler

import (
	"fmt"
	"math"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/deeptech-kmitl/Cicero-Backend/config"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/entities"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/files"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/files/filesUsecase"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/product"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/product/productUsecase"
	"github.com/deeptech-kmitl/Cicero-Backend/pkg/utils"
	"github.com/gofiber/fiber/v2"
)

type productHandlerErrCode = string

const (
	FindOneProductErr       productHandlerErrCode = "product-001"
	AddProductErr           productHandlerErrCode = "product-002"
	DeleteProductErr        productHandlerErrCode = "product-003"
	findProductErr          productHandlerErrCode = "product-004"
	UpdateProductErr        productHandlerErrCode = "product-005"
	FindImageByProductIdErr productHandlerErrCode = "product-006"
	DeleteImageProductErr   productHandlerErrCode = "product-007"
	InsertImageProductErr   productHandlerErrCode = "product-008"
)

type IProductHandler interface {
	FindOneProduct(c *fiber.Ctx) error
	AddProduct(c *fiber.Ctx) error
	DeleteProduct(c *fiber.Ctx) error
	FindProduct(c *fiber.Ctx) error
	UpdateProduct(c *fiber.Ctx) error
	FindImageByProductId(c *fiber.Ctx) error
	DeleteImageProduct(c *fiber.Ctx) error
	InsertImageProduct(c *fiber.Ctx) error
}

type productHandler struct {
	cfg            config.IConfig
	productUsecase productUsecase.IProductUsecase
	fileUsecase    filesUsecase.IFilesUsecase
}

func ProductHandler(productUsecase productUsecase.IProductUsecase, fileUsecase filesUsecase.IFilesUsecase, cfg config.IConfig) IProductHandler {
	return &productHandler{
		productUsecase: productUsecase,
		cfg:            cfg,
		fileUsecase:    fileUsecase,
	}
}

func (h *productHandler) FindOneProduct(c *fiber.Ctx) error {
	prodId := strings.TrimSpace(c.Params("product_id"))
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

func (h *productHandler) AddProduct(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(AddProductErr),
			err.Error(),
		).Res()
	}

	productTitle, exists := form.Value["product_title"]
	if !exists || len(productTitle[0]) == 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(AddProductErr),
			"product_title is required",
		).Res()
	}
	productDesc, exists := form.Value["product_desc"]
	if !exists || len(productDesc[0]) == 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(AddProductErr),
			"product_desc is required",
		).Res()
	}
	productPrice, exists := form.Value["product_price"]
	if !exists || len(productPrice[0]) == 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(AddProductErr),
			"product_price is required",
		).Res()
	}
	productColor, exists := form.Value["product_color"]
	if !exists || len(productColor[0]) == 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(AddProductErr),
			"product_color is required",
		).Res()
	}
	productSize, exists := form.Value["product_size"]
	if !exists || len(productSize[0]) == 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(AddProductErr),
			"product_size is required",
		).Res()
	}
	productSex, exists := form.Value["product_sex"]
	if !exists || len(productSex[0]) == 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(AddProductErr),
			"product_sex is required",
		).Res()
	}
	productCategory, exists := form.Value["product_category"]
	if !exists || len(productCategory[0]) == 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(AddProductErr),
			"product_category is required",
		).Res()
	}
	images, exists := form.File["images"]
	if !exists || len(images) == 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(AddProductErr),
			"images is required",
		).Res()
	}

	productPriceFloat, err := strconv.ParseFloat(productPrice[0], 64)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(AddProductErr),
			"invalid product price",
		).Res()
	}

	if productPriceFloat < 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(AddProductErr),
			"product price must be greater than 0",
		).Res()
	}

	// files ext validation
	extMap := map[string]string{
		"png":  "png",
		"jpg":  "jpg",
		"jpeg": "jpeg",
	}

	req := make([]*files.FileReq, 0)

	for _, file := range images {
		// check file extension
		ext := strings.TrimPrefix(filepath.Ext(file.Filename), ".")
		if extMap[ext] != ext || extMap[ext] == "" {
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(AddProductErr),
				"invalid file extension",
			).Res()
		}
		// check file size
		if file.Size > int64(h.cfg.App().FileLimit()) {
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(AddProductErr),
				fmt.Sprintf("file size must less than %d MB", int(math.Ceil(float64(h.cfg.App().FileLimit())/math.Pow(1024, 2)))),
			).Res()
		}

		filename := utils.RandFileName(ext)
		req = append(req, &files.FileReq{
			File:        file,
			Destination: fmt.Sprintf("%s/%s", productTitle, filename),
			FileName:    filename,
			Extension:   ext,
		})
	}

	img, err := h.fileUsecase.UploadToGCP(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(AddProductErr),
			err.Error(),
		).Res()
	}

	prod := &product.AddProduct{
		ProductTitle:    productTitle[0],
		ProductDesc:     productDesc[0],
		ProductPrice:    productPriceFloat,
		ProductColor:    productColor[0],
		ProductSize:     productSize[0],
		ProductSex:      productSex[0],
		ProductCategory: productCategory[0],
		Images:          img,
	}

	result, err := h.productUsecase.AddProduct(prod)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(AddProductErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusCreated, result).Res()
}

func (h *productHandler) DeleteProduct(c *fiber.Ctx) error {
	prodId := strings.TrimSpace(c.Params("product_id"))
	result, err := h.productUsecase.DeleteProduct(prodId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(DeleteProductErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, result).Res()
}

func (h *productHandler) FindProduct(c *fiber.Ctx) error {
	req := &product.ProductFilter{
		PaginationReq: &entities.PaginationReq{},
		SortReq:       &entities.SortReq{},
	}

	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findProductErr),
			err.Error(),
		).Res()
	}

	products := h.productUsecase.FindProduct(req)
	return entities.NewResponse(c).Success(fiber.StatusOK, products).Res()
}

func (h *productHandler) UpdateProduct(c *fiber.Ctx) error {
	req := new(product.UpdateProduct)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(UpdateProductErr),
			err.Error(),
		).Res()
	}

	result, err := h.productUsecase.UpdateProduct(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(UpdateProductErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, result).Res()
}

func (h *productHandler) FindImageByProductId(c *fiber.Ctx) error {
	prodId := strings.TrimSpace(c.Params("product_id"))
	result, err := h.productUsecase.FindImageByProductId(prodId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(FindImageByProductIdErr),
			err.Error(),
		).Res()
	}
	return entities.NewResponse(c).Success(fiber.StatusOK, result).Res()
}

func (h *productHandler) DeleteImageProduct(c *fiber.Ctx) error {
	imageId := strings.TrimSpace(c.Params("image_id"))

	result, err := h.productUsecase.DeleteImageProduct(imageId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(DeleteImageProductErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, result).Res()
}

func (h *productHandler) InsertImageProduct(c *fiber.Ctx) error {
	prodId := strings.TrimSpace(c.Params("product_id"))
	form, err := c.MultipartForm()
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(InsertImageProductErr),
			err.Error(),
		).Res()
	}

	images, exists := form.File["images"]
	if !exists || len(images) == 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(InsertImageProductErr),
			"images is required",
		).Res()
	}

	// files ext validation
	extMap := map[string]string{
		"png":  "png",
		"jpg":  "jpg",
		"jpeg": "jpeg",
	}

	req := make([]*files.FileReq, 0)

	for _, file := range images {
		// check file extension
		ext := strings.TrimPrefix(filepath.Ext(file.Filename), ".")
		if extMap[ext] != ext || extMap[ext] == "" {
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(InsertImageProductErr),
				"invalid file extension",
			).Res()
		}
		// check file size
		if file.Size > int64(h.cfg.App().FileLimit()) {
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(InsertImageProductErr),
				fmt.Sprintf("file size must less than %d MB", int(math.Ceil(float64(h.cfg.App().FileLimit())/math.Pow(1024, 2)))),
			).Res()
		}

		filename := utils.RandFileName(ext)
		req = append(req, &files.FileReq{
			File:        file,
			Destination: fmt.Sprintf("%s/%s", prodId, filename),
			FileName:    filename,
			Extension:   ext,
		})
	}

	img, err := h.fileUsecase.UploadToGCP(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(InsertImageProductErr),
			err.Error(),
		).Res()
	}

	result, err := h.productUsecase.InsertImageProduct(img, prodId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(InsertImageProductErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, result).Res()
}
