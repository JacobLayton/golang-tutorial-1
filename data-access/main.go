package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5"
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

	alb, err := albumById(2)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Album found: %v\n", alb)

	albID, err := addAlbum(Album{
		Title: "The White Album",
		Artist: "the Beatles",
		Price: 10.99,
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("ID of added album: %v\n", albID)
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

// albumByID queries for the album with the specified ID.
func albumById(id int64) (Album, error) {
	// an album to hold data from the returned row
	var alb Album

	row := db.QueryRow(context.Background(), "SELECT * FROM album WHERE id = $1", id)
	if err := row.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
		if err == pgx.ErrNoRows {
			return alb, fmt.Errorf("albumById %d: no such album", id)
		}
		return alb, fmt.Errorf("albumsById %d: %v", id, err)
	}
	return alb, nil
}

// addAlbum adds the specified album to the database,
// returning the album ID of the new entry
func addAlbum(alb Album) (int64, error) {
    var id int64

    err := db.QueryRow(context.Background(),
        "INSERT INTO album (title, artist, price) VALUES ($1, $2, $3) RETURNING id",
        alb.Title,
        alb.Artist,
        alb.Price,
    ).Scan(&id)

    if err != nil {
        return 0, fmt.Errorf("addAlbum: %w", err)
    }

    return id, nil
}