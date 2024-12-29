package db

import (
	"context"
	"database/sql"
	_ "embed"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var ddl string

func run() error {
	ctx := context.Background()
	db, err := sql.Open("sqlite3", "file:app.db?cache=shared&mode=rwc")
	if err != nil {
		return err
	}

	// create tables
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return err
	}

	queries := New(db)

	// List all users
	users, err := queries.ListUsers(ctx)
	if err != nil {
		return err
	}
	log.Println("Users:", users)

	return nil
}

func InitDB() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
