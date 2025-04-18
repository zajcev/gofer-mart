package database

import (
	"context"
	"database/sql"
	_ "embed"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"log"
	"os"
	"path/filepath"
)

//go:embed scripts/000001_tables.up.sql
var file string

type DB interface {
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
}

type DBService struct {
	db DB
}

func NewDBService(db DB) *DBService {
	return &DBService{db: db}
}

func Migration(DBUrl string) {
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
