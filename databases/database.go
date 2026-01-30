package databases

import (
	"database/sql"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

func Init(connectionString string) (*sql.DB, error) {
	// Open database via stdlib always prepared statement "stmtcache_xxx" already exists
	// db, err := sql.Open("pgx", connectionString)
	// if err != nil {
	// 	return nil, err
	// }
	cfg, err := pgx.ParseConfig(connectionString)
	if err != nil {
		return nil, err
	}

	cfg.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	db := stdlib.OpenDB(*cfg)

	// Set connection pool settings (optional tapi recommended)
	db.SetMaxOpenConns(5)
	db.SetMaxIdleConns(3)
	db.SetConnMaxLifetime(2 * time.Minute)

	// Test connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	log.Println("Database connected successfully")
	return db, nil
}
