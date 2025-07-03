package migrations

import (
	"context"
	"database/sql"

	goose "github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upGoFileTest, downGoFileTest)
}

func upGoFileTest(ctx context.Context, tx *sql.Tx) error {
	_ = ctx
	_ = tx
	// slog.Info("ЗАГЛУШКА", slog.String("Миграции из го", "Проблема в том, что дропать тоже надо через го"))
	// _, err := tx.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS utest (
	// id INT PRIMARY KEY
	// );`)
	// if err != nil {
	// 	return err
	// }
	return nil
}

func downGoFileTest(ctx context.Context, tx *sql.Tx) error {
	_ = ctx
	_ = tx
	// slog.Info("ЗАГЛУШКА", slog.String("Дроп из го", "Ну и дроп соответственно тоже через го"))
	// _, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS utest;`)
	// if err != nil {
	// 	return err
	// }
	return nil
}
