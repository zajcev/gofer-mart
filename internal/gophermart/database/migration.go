package database

import (
	"context"
	"database/sql"
	_ "embed"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"path/filepath"
)

//go:embed scripts/000001_tables.up.sql
var file string

var db *pgxpool.Pool

func Init(ctx context.Context, DBUrl string) error {
	d, err := pgxpool.New(ctx, DBUrl)
	if err != nil {
		log.Printf("Error while connect to database: %v", err)
		return err
	}
	db = d
	migration(DBUrl)
	return nil
}

func migration(DBUrl string) {
	wd, _ := os.Getwd()
	filePath := filepath.Join(wd, "internal/gophermart/database/scripts/")
	d, _ := sql.Open("postgres", DBUrl)
	driver, _ := postgres.WithInstance(d, &postgres.Config{})
	m, err := migrate.NewWithDatabaseInstance(
		"file:///"+filePath, "postgres", driver)
	if err != nil {
		log.Fatal(err)
	}
	if err = m.Up(); err != nil {
		log.Printf("Migrating database: %v", err)
	} else {
		log.Print("Migration complete")
	}
}
