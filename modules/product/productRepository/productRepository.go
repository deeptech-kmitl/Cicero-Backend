package productRepository

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/deeptech-kmitl/Cicero-Backend/modules/entities"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/files"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/product"
	"github.com/deeptech-kmitl/Cicero-Backend/modules/product/productPattern"
	"github.com/jmoiron/sqlx"
)

type IProductRepository interface {
	FindOneProduct(prodId string) (*product.Product, error)
	InsertProduct(req *product.AddProduct) (*product.Product, error)
	DeleteProduct(productId string) error
	FindProduct(req *product.ProductFilter) ([]*product.Product, int)
	UpdateProduct(req *product.UpdateProduct) (*product.Product, error)
	FindImageByProductId(productId string) ([]*entities.ImageRes, error)
	DeleteImageProduct(imageId string) error
	InsertImageProduct(images []*files.FileRes, prodId string) error
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

func (r *productRepository) UpdateProduct(req *product.UpdateProduct) (*product.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	queryWhereStack := make([]string, 0)
	values := make([]any, 0)
	lastIndex := 1

	query := `
	UPDATE "Product" SET`

	if req.ProductCategory != "" {
		values = append(values, req.ProductCategory)

		queryWhereStack = append(queryWhereStack, fmt.Sprintf(`
		"product_category" = $%d?`, lastIndex))

		lastIndex++
	}

	if req.ProductColor != "" {
		values = append(values, req.ProductColor)

		queryWhereStack = append(queryWhereStack, fmt.Sprintf(`
		"product_color" = $%d?`, lastIndex))

		lastIndex++
	}

	if req.ProductDesc != "" {
		values = append(values, req.ProductDesc)

		queryWhereStack = append(queryWhereStack, fmt.Sprintf(`
		"product_desc" = $%d?`, lastIndex))

		lastIndex++
	}

	if req.ProductSex != "" {
		values = append(values, req.ProductSex)

		queryWhereStack = append(queryWhereStack, fmt.Sprintf(`
		"product_sex" = $%d?`, lastIndex))

		lastIndex++
	}
	if req.ProductSize != "" {
		values = append(values, req.ProductSize)

		queryWhereStack = append(queryWhereStack, fmt.Sprintf(`
		"product_size" = $%d?`, lastIndex))

		lastIndex++
	}
	if req.ProductTitle != "" {
		values = append(values, req.ProductTitle)

		queryWhereStack = append(queryWhereStack, fmt.Sprintf(`
		"product_title" = $%d?`, lastIndex))

		lastIndex++
	}
	if req.ProductPrice > 0 {
		values = append(values, req.ProductPrice)

		queryWhereStack = append(queryWhereStack, fmt.Sprintf(`
		"product_price" = $%d?`, lastIndex))

		lastIndex++
	}

	values = append(values, req.Id)

	queryClose := fmt.Sprintf(`
	WHERE "id" = $%d;`, lastIndex)

	for i := range queryWhereStack {
		if i != len(queryWhereStack)-1 {
			query += strings.Replace(queryWhereStack[i], "?", ",", 1)
		} else {
			query += strings.Replace(queryWhereStack[i], "?", "", 1)
		}
	}
	query += queryClose

	if _, err := r.db.ExecContext(ctx, query, values...); err != nil {
		return nil, fmt.Errorf("update product failed: %v", err)
	}

	product, err := r.FindOneProduct(req.Id)
	if err != nil {
		return nil, err
	}
	return product, nil

}

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

func (r *productRepository) FindImageByProductId(productId string) ([]*entities.ImageRes, error) {
	query := `
	SELECT
		"id",
		"filename",
		"url"
	FROM "Image"
	WHERE "product_id" = $1;`

	images := make([]*entities.ImageRes, 0)
	if err := r.db.Select(&images, query, productId); err != nil {
		return nil, fmt.Errorf("find images failed: %v", err)
	}
	return images, nil
}

// delete image by image id
func (r *productRepository) DeleteImageProduct(imageId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15) // Timeout of 15 seconds
	defer cancel()
	query := `DELETE FROM "Image" WHERE "id" = $1;`

	if _, err := r.db.ExecContext(ctx, query, imageId); err != nil {
		return fmt.Errorf("delete image product failed: %v", err)
	}

	return nil
}

func (r *productRepository) InsertImageProduct(images []*files.FileRes, prodId string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	query := `
	INSERT INTO "Image" (
		"filename",
		"url",
		"product_id"
	)
	VALUES`

	valueStack := make([]any, 0)
	var index int
	for i := range images {
		valueStack = append(valueStack,
			images[i].FileName,
			images[i].Url,
			prodId,
		)

		if i != len(images)-1 {
			query += fmt.Sprintf(`
			($%d, $%d, $%d),`, index+1, index+2, index+3)
		} else {
			query += fmt.Sprintf(`
			($%d, $%d, $%d);`, index+1, index+2, index+3)
		}
		index += 3
	}

	if _, err := r.db.ExecContext(
		ctx,
		query,
		valueStack...,
	); err != nil {
		return fmt.Errorf("insert images product failed: %v", err)
	}
	return nil
}
