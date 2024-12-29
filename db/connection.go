package db

import (
	"context"
	"database/sql"
	_ "embed"
	"log/slog"

	_ "github.com/mattn/go-sqlite3"
)

//go:embed schema/schema.sql
var ddl string

var Sqlc *sql.DB // DB is a package level variable to hold the connection pool

var CallQueries *Queries

func run() error {
	ctx := context.Background()
	var err error
	Sqlc, err = sql.Open("sqlite3", "file:app.db?cache=shared&mode=rwc")
	if err != nil {
		return err
	}

	// create tables
	if _, err := Sqlc.ExecContext(ctx, ddl); err != nil {
		slog.Warn(err.Error())
	}

	CallQueries = New(Sqlc) // Initialize the queries

	return nil
}

func InitDB() {
	if err := run(); err != nil {
		slog.Error(err.Error())
	}
}
