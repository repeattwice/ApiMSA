package servs

import (
	"context"
	"os"

	"github.com/jackc/pgx/v5"
)

func GetBDConnection(ctx context.Context) (*pgx.Conn, error) { // вроде готово, можно проверить и улучшить
	port := os.Getenv("DB_PORT")
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	connection := "postgres://" + user + ":" + password + "@" + host + ":" + port + "/postgres?sslmode=disable"
	return pgx.Connect(ctx, connection)
}

func CreateTables(ctx context.Context, conn *pgx.Conn) { // готово
	sqlForUser := `
	CREATE TABLE IF NOT EXISTS users(
		user_id SERIAL PRIMARY KEY,
		user_name VARCHAR NOT NULL,
		last_name VARCHAR NOT NULL,
		email VARCHAR
	);
	`
	conn.Exec(ctx, sqlForUser)

	sqlForItems := `
	CREATE TABLE IF NOT EXISTS items(
		item_name VARCHAR NOT PRIMARY KEY,
		item_price INTEGER NOT NULL
	);
	`
	conn.Exec(ctx, sqlForItems)

	sqlForCort := `
	CREATE TABLE IF NOT EXISTS cort(
		user_id INTEGER REFERENCES useers(user_id) ON DELETE CASCADE,
		item_name_in_cort VARCHAR REFERENCES items(item_name),
	);
	`
	conn.Exec(ctx, sqlForCort)
}
