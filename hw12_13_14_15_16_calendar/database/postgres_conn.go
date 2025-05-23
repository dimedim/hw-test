package database

import (
	"context"
	"time"

	"github.com/dimedim/hw-test/hw12_13_14_15_calendar/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func MustConnectDatabase(ctx context.Context, cfg *config.Config) *sqlx.DB {

	// db, err := sql.Open("pgx", cfg.GetPostgresDSN())
	// if err != nil {
	// 	panic("Unable to connect to database: " + err.Error())
	// }
	// pgx.Open

	db, err := sqlx.Connect("pgx", cfg.GetPostgresDSN())
	if err != nil {
		panic("Unable to connect to database: " + err.Error())
	}
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(time.Minute)
	// pool, err := pgxpool.New(ctx, cfg.GetPostgresDSN())
	// if err != nil {
	// 	panic("Unable to connect to database: " + err.Error())
	// }

	return db
}
