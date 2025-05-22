package database

import (
	"context"

	"github.com/dimedim/hw-test/hw12_13_14_15_calendar/internal/config"
	"github.com/jackc/pgx/v5"
)

func InitDatabase(ctx context.Context, cfg *config.Config) *pgx.Conn {
	conn, err := pgx.Connect(ctx, cfg.GetPostgresDSN())
	if err != nil {
		panic("Unable to connect to database: " + err.Error())
	}
	defer conn.Close(ctx)

	return conn
}
