package sqlstorage

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type Storage struct { // TODO
	DB *pgx.Conn

	// //! TODO: заглушка, нужно убрать
	// app.Storage
}

func New(db *pgx.Conn) *Storage {
	return &Storage{DB: db}
}

func (s *Storage) Connect(ctx context.Context) error {
	// TODO
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	// TODO
	return nil
}
