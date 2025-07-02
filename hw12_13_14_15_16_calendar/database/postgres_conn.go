package database

import (
	"context"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib" // revive:disable:blank-imports
	"github.com/jmoiron/sqlx"
)

const (
	MaxLifetime  = time.Minute
	MaxIdleConns = 5
	MaxOpenConns = 20
)

func MustConnectDatabase(ctx context.Context, dsn string) *sqlx.DB {
	db, err := sqlx.ConnectContext(ctx, "pgx", dsn)
	if err != nil {
		panic("Unable to connect to database: " + err.Error())
	}
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(MaxLifetime)
	db.SetMaxIdleConns(MaxIdleConns)
	return db
}
