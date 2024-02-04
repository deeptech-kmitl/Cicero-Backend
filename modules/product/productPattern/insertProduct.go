package productPattern

import (
	"context"
	"fmt"
	"time"

	"github.com/deeptech-kmitl/Cicero-Backend/modules/product"
	"github.com/jmoiron/sqlx"
)

type IInsertProductBuilder interface {
	initTransaction() error
	insertProduct() error
	insertImages() error
	commit() error
	getProductId() string
}

type insertProductBuilder struct {
	db  *sqlx.DB
	tx  *sqlx.Tx
	req *product.AddProduct
}

func InsertProductBuilder(db *sqlx.DB, req *product.AddProduct) IInsertProductBuilder {
	return &insertProductBuilder{
		db:  db,
		req: req,
	}
}

func (b *insertProductBuilder) initTransaction() error {
	tx, err := b.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	b.tx = tx

	return nil
}

func (b *insertProductBuilder) getProductId() string {
	return b.req.Id
}

func (b *insertProductBuilder) insertProduct() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	query := `
	INSERT INTO "Product" 
	(
		"product_title",
		"product_desc",
		"product_price",
		"product_color",
		"product_size",
		"product_sex",
		"product_category"
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING "id";`

	if err := b.tx.QueryRowxContext(
		ctx,
		query,
		b.req.ProductTitle,
		b.req.ProductDesc,
		b.req.ProductPrice,
		b.req.ProductColor,
		b.req.ProductSize,
		b.req.ProductSex,
		b.req.ProductCategory,
	).Scan(&b.req.Id); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert product failed: %v", err)
	}
	return nil
}

func (b *insertProductBuilder) insertImages() error {
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
	for i := range b.req.Images {
		valueStack = append(valueStack,
			b.req.Images[i].FileName,
			b.req.Images[i].Url,
			b.req.Id,
		)

		if i != len(b.req.Images)-1 {
			query += fmt.Sprintf(`
			($%d, $%d, $%d),`, index+1, index+2, index+3)
		} else {
			query += fmt.Sprintf(`
			($%d, $%d, $%d);`, index+1, index+2, index+3)
		}
		index += 3
	}

	if _, err := b.tx.ExecContext(
		ctx,
		query,
		valueStack...,
	); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert images failed: %v", err)
	}
	return nil
}

func (b *insertProductBuilder) commit() error {
	if err := b.tx.Commit(); err != nil {
		return fmt.Errorf("commit failed: %v", err)
	}
	return nil
}

type insertProductEngineer struct {
	builder IInsertProductBuilder
}

func InsertProductEngineer(builder IInsertProductBuilder) *insertProductEngineer {
	return &insertProductEngineer{
		builder: builder,
	}
}

func (en *insertProductEngineer) InsertProduct() (string, error) {
	if err := en.builder.initTransaction(); err != nil {
		return "", err
	}

	if err := en.builder.insertProduct(); err != nil {
		return "", err
	}

	if err := en.builder.insertImages(); err != nil {
		return "", err
	}

	if err := en.builder.commit(); err != nil {
		return "", err
	}
	return en.builder.getProductId(), nil
}
