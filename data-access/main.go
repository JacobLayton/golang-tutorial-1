package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func main() {
	ctx := context.Background()

    // Build connection string
    dsn := fmt.Sprintf(
        "postgres://%s:%s@127.0.0.1:5432/recordings?sslmode=disable",
        os.Getenv("DBUSER"),
        os.Getenv("DBPASS"),
    )

	// Parse and configure connection
    cfg, err := pgxpool.ParseConfig(dsn)
    if err != nil {
        log.Fatal(err)
    }

	// Create connection pool
    db, err = pgxpool.NewWithConfig(ctx, cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

	// Verify connection
    if err = db.Ping(ctx); err != nil {
        log.Fatal(err)
    }

	fmt.Println("Connected!")
}