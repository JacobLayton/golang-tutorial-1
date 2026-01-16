package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

type Album struct {
	ID int64
	Title  string
	Artist string
	Price  float32
}

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

	albums, err := albumsByArtist("John Coltrane")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Albums found: %v\n", albums)
}

// albumsByArtist queries for albums that have the specified artist name.
func albumsByArtist(name string) ([]Album, error) {
	// an album slice to hold data from recurned rows
	var albums []Album

	// rows, err := db.Query("Select * FROM album WHERE artist = ?", name)
	rows, err := db.Query(context.Background(), "SELECT * FROM album WHERE artist = $1", name)
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	defer rows.Close()
	// Loop through rows, using Scan to assign column data to struct fields.
	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
		}
		albums = append(albums, alb)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}
	return albums, nil
}