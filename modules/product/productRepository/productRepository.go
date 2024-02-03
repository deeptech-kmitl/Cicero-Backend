package productRepository

import "github.com/jmoiron/sqlx"

type IProductRepository interface {
}

type productRepository struct {
	db *sqlx.DB
}

func ProductRepository(db *sqlx.DB) IProductRepository {
	return &productRepository{
		db: db,
	}
}

func (r *productRepository) FindOneProduct(prodId string) {

}
