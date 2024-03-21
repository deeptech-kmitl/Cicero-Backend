package product

import (
	"github.com/deeptech-kmitl/Cicero-Backend/modules/entities"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/files"
)

type GetAllProduct struct {
	Id              string               `db:"id" json:"id" form:"id"`
	ProductTitle    string               `db:"product_title" json:"product_title" form:"product_title"`
	ProductPrice    float64              `db:"product_price" json:"product_price" form:"product_price"`
	ProductColor    string               `db:"product_color" json:"product_color" form:"product_color"`
	ProductSize     []*string            `db:"product_size" json:"product_size" form:"product_size"`
	ProductSex      string               `db:"product_sex" json:"product_sex" form:"product_sex"`
	ProductDesc     string               `db:"product_desc" json:"product_desc" form:"product_desc"`
	ProductCategory string               `db:"product_category" json:"product_category" form:"product_category"`
	ProductStock    int                  `db:"product_stock" json:"product_stock" form:"product_stock"`
	Images          []*entities.ImageRes `json:"images" form:"images"`
}

type Product struct {
	Id              string               `db:"id" json:"id" form:"id"`
	ProductTitle    string               `db:"product_title" json:"product_title" form:"product_title"`
	ProductPrice    float64              `db:"product_price" json:"product_price" form:"product_price"`
	ProductColor    string               `db:"product_color" json:"product_color" form:"product_color"`
	ProductSize     string               `db:"product_size" json:"product_size" form:"product_size"`
	ProductSex      string               `db:"product_sex" json:"product_sex" form:"product_sex"`
	ProductDesc     string               `db:"product_desc" json:"product_desc" form:"product_desc"`
	ProductCategory string               `db:"product_category" json:"product_category" form:"product_category"`
	ProductStock    int                  `db:"product_stock" json:"product_stock" form:"product_stock"`
	Images          []*entities.ImageRes `json:"images" form:"images"`
}

type AddProduct struct {
	Id              string           `db:"id" json:"id" form:"id"`
	ProductTitle    string           `db:"product_title" json:"product_title" form:"product_title"`
	ProductPrice    float64          `db:"product_price" json:"product_price" form:"product_price"`
	ProductColor    string           `db:"product_color" json:"product_color" form:"product_color"`
	ProductSize     string           `db:"product_size" json:"product_size" form:"product_size"`
	ProductSex      string           `db:"product_sex" json:"product_sex" form:"product_sex"`
	ProductDesc     string           `db:"product_desc" json:"product_desc" form:"product_desc"`
	ProductCategory string           `db:"product_category" json:"product_category" form:"product_category"`
	ProductStock    int              `db:"product_stock" json:"product_stock" form:"product_stock"`
	Images          []*files.FileRes `json:"images" form:"images"`
}

type ProductFilter struct {
	Id     string `json:"id" query:"id"`
	Search string `json:"search" query:"search"` // search by title and description
	*entities.PaginationReq
	*entities.SortReq
}

type UpdateProduct struct {
	Id              string           `db:"id" json:"id" form:"id"`
	ProductTitle    string           `db:"product_title" json:"product_title" form:"product_title"`
	ProductPrice    float64          `db:"product_price" json:"product_price" form:"product_price"`
	ProductColor    string           `db:"product_color" json:"product_color" form:"product_color"`
	ProductSize     string           `db:"product_size" json:"product_size" form:"product_size"`
	ProductSex      string           `db:"product_sex" json:"product_sex" form:"product_sex"`
	ProductDesc     string           `db:"product_desc" json:"product_desc" form:"product_desc"`
	ProductCategory string           `db:"product_category" json:"product_category" form:"product_category"`
	ProductStock    int              `db:"product_stock" json:"product_stock" form:"product_stock"`
	Images          []*files.FileRes `json:"images" form:"images"`
}
