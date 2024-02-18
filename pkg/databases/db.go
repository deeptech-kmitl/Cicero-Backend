package databases

import (
	"log"

	"github.com/deeptech-kmitl/Cicero-Backend/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func DbConnect(cfg config.IDbConfig) *sqlx.DB {
	// Connect
	db, err := sqlx.Connect("pgx", cfg.Url())
	if err != nil {
		log.Fatalf("connect to db failed: %v\n", err)
	}
	// Set time zone
	_, err = db.Exec("SET TIME ZONE 'Asia/Bangkok'")
	if err != nil {
		log.Fatalf("set time zone failed: %v\n", err)
	}

	db.DB.SetMaxOpenConns(cfg.MaxOpenConns())
	return db
}
