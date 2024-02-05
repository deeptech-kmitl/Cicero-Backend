package productRepository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/deeptech-kmitl/Cicero-Backend/modules/entities"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/product"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/product/productPattern"
	"github.com/jmoiron/sqlx"
)

type IProductRepository interface {
	FindOneProduct(prodId string) (*product.Product, error)
	InsertProduct(req *product.AddProduct) (*product.Product, error)
	DeleteProduct(productId string) error
	FindProduct(req *product.ProductFilter) ([]*product.Product, int)
}

type productRepository struct {
	db *sqlx.DB
}

func ProductRepository(db *sqlx.DB) IProductRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) FindOneProduct(prodId string) (*product.Product, error) {
	query := `
	SELECT
		to_jsonb("t")
	FROM (
		SELECT
			"p"."id",
			"p"."product_title",
			"p"."product_desc",
			"p"."product_price",
			"p"."product_color",
			"p"."product_size",
			"p"."product_sex",
			"p"."product_category",
			(
				SELECT
					COALESCE(array_to_json(array_agg("it")), '[]'::json)
				FROM (
					SELECT
						"i"."id",
						"i"."filename",
						"i"."url"
					FROM "Image" "i"
					WHERE "i"."product_id" = "p"."id"
				) AS "it"
			) AS "images"
		FROM "Product" "p"
		WHERE "p"."id" = $1
		LIMIT 1
	) AS "t";
	`

	//COALESCE(array_to_json(array_agg("it")), '[]'::json)
	// คือ ถ้าไม่มีข้อมูล(null) ให้ return '[]'::json แทน

	// array_agg คือ การรวมข้อมูลใน array ที่มีค่าเหมือนกันเป็น 1 row
	// array_to_json คือ การแปลง array เป็น json
	// to_jsonb คือ การแปลง NON-JSON เป็น jsonb

	productBytes := make([]byte, 0)
	product := &product.Product{
		Images: make([]*entities.ImageRes, 0), //เวลาสร้าง struct ใหม่ แล้วข้างในมี array ให้ make array ไว้เลยเพื่อป้องกัน null pointer
	}
	if err := r.db.Get(&productBytes, query, prodId); err != nil {
		return nil, fmt.Errorf("get product failed: %v", err)
	}
	if err := json.Unmarshal(productBytes, &product); err != nil {
		return nil, fmt.Errorf("unmarshal product failed: %v", err)
	}

	return product, nil

}

func (r *productRepository) InsertProduct(req *product.AddProduct) (*product.Product, error) {
	builder := productPattern.InsertProductBuilder(r.db, req)
	productId, err := productPattern.InsertProductEngineer(builder).InsertProduct()
	if err != nil {
		return nil, fmt.Errorf("insert product failed: %v", err)
	}

	product, err := r.FindOneProduct(productId)
	if err != nil {
		return nil, fmt.Errorf("find product failed: %v", err)
	}

	return product, nil

}

// func (r *productRepository) UpdateProduct(req *products.Products) (*products.Products, error) {
// 	builder := productsPatterns.UpdateProductBuilder(r.db, req, r.fileUsecase, r.cfg)
// 	engineer := productsPatterns.UpdateProductEngineer(builder)

// 	if err := engineer.UpdateProduct(); err != nil {
// 		return nil, err
// 	}

// 	product, err := r.FindOneProduct(req.Id)
// 	if err != nil {
// 		return nil,  err
// 	}
// 	return product, nil

// }

func (r *productRepository) DeleteProduct(productId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15) // Timeout of 15 seconds
	defer cancel()
	query := `DELETE FROM "Product" WHERE "id" = $1;`

	if _, err := r.db.ExecContext(ctx, query, productId); err != nil {
		return fmt.Errorf("delete product failed: %v", err)
	}

	return nil
}

func (r *productRepository) FindProduct(req *product.ProductFilter) ([]*product.Product, int) {
	builder := productPattern.FindProductBuilder(r.db, req)
	engineer := productPattern.FindProductEngineer(builder)

	result := engineer.FindProduct().Result()
	count := engineer.CountProduct().Count()

	engineer.FindProduct().PrintQuery()

	return result, count
}
