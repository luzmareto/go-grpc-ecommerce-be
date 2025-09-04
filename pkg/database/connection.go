package database

import (
	"context"

	_ "github.com/lib/pq"

	"database/sql"
)

func ConnectDB(ctx context.Context, connstr string) *sql.DB {
	db, err := sql.Open("postgres", connstr)
	if err != nil {
		panic(err)
	}

	err = db.PingContext(ctx)
	if err != nil {
		panic(err)
	}

	return db
}
