package productRepository

import (
	"encoding/json"
	"fmt"

	"github.com/deeptech-kmitl/Cicero-Backend/modules/entities"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/product"
	"github.com/jmoiron/sqlx"
)

type IProductRepository interface {
	FindOneProduct(prodId string) (*product.Product, error)
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
