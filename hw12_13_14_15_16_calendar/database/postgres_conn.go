package database

import (
	"context"
	"time"

	"github.com/dimedim/hw-test/hw12_13_14_15_16_calendar/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib" // revive:disable:blank-imports
	"github.com/jmoiron/sqlx"
)

const (
	MaxLifetime  = time.Minute
	MaxIdleConns = 5
	MaxOpenConns = 20
)

func MustConnectDatabase(ctx context.Context, cfg *config.Config) *sqlx.DB {
	db, err := sqlx.ConnectContext(ctx, "pgx", cfg.GetPostgresDSN())
	if err != nil {
		panic("Unable to connect to database: " + err.Error())
	}
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(MaxLifetime)
	db.SetMaxIdleConns(MaxIdleConns)
	return db
}
