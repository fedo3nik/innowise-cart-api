package db

import (
	"runtime"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/spf13/viper"
)

func InitDb() (*sqlx.DB, error) {
	path := viper.Get("DB_URL").(string)

	pool, err := sqlx.Open("postgres", path)

	pool.SetMaxIdleConns(1)
	pool.SetConnMaxLifetime(2 * time.Minute)
	pool.SetMaxOpenConns(runtime.NumCPU())

	if err != nil {
		return nil, err
	}

	return pool, nil
}
