package database

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/spf13/cast"
	"log"
)

func AddUser(ctx context.Context, login string, pass string) {
	_, err := db.Exec(ctx, addUser, login, pass)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Printf("PGError: %v", pgErr)
		} else {
			log.Printf("Error while addUser: Error=%v, login=%v", err, login)
		}
	}
}

func GetLogin(ctx context.Context, login string) string {
	row, err := db.Query(ctx, getLogin, login)
	if err != nil {
		log.Printf("Error while executing query: %v", err)
		return ""
	}
	defer row.Close()

	if !row.Next() {
		log.Printf("No rows returned for login: %v", login)
		return ""
	}

	if err = row.Err(); err != nil {
		log.Printf("Error after row.Next(): %v", err)
		return ""
	}

	var value interface{}
	err = row.Scan(&value)
	if err != nil {
		log.Printf("Error while scan login value: %v", err)
		return ""
	}

	return cast.ToString(value)
}

func GetPassword(ctx context.Context, login string) string {
	row, err := db.Query(ctx, getPassword, login)
	if err != nil {
		log.Printf("Error while executing query: %v", err)
		return ""
	}
	defer row.Close()

	if !row.Next() {
		log.Printf("No passwords returned for login: %v", login)
		return ""
	}

	if err = row.Err(); err != nil {
		log.Printf("Error after row.Next(): %v", err)
		return ""
	}

	var value interface{}
	err = row.Scan(&value)
	if err != nil {
		log.Printf("Error while scan password value: %v", err)
		return ""
	}

	return cast.ToString(value)
}

func NewSession(ctx context.Context, login string, token string) error {
	id, err := getUserId(ctx, login)
	if err != nil {
		return err
	}
	if id != 0 {
		_, err := db.Exec(ctx, addSession, id, token)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) {
				log.Printf("PGError: %v", pgErr)
			} else {
				log.Printf("Error while addSession: Error=%v, login=%v", err, login)
			}
		}
	}
	return nil
}

func getUserId(ctx context.Context, login string) (int, error) {
	var id interface{}
	rowId, _ := db.Query(ctx, getUserIdByLogin, login)
	if rowId.Err() != nil {
		log.Printf("Error while execute query: %v", rowId.Err())
		return 0, rowId.Err()
	}
	defer rowId.Close()
	if !rowId.Next() {
		log.Printf("No tokens returned for login: %v", login)
		return 0, errors.New("No tokens returned for login: " + login)
	}
	err := rowId.Scan(&id)
	if err != nil {
		log.Printf("Error while scan token value: %v", rowId.Err())
		return 0, err
	}
	return cast.ToInt(id), nil
}

func GetUserIdByToken(ctx context.Context, token string) (int, error) {
	var id interface{}
	rowId, _ := db.Query(ctx, getUserIdByToken, token)
	if rowId.Err() != nil {
		log.Printf("Error while execute query: %v", rowId.Err())
		return 0, rowId.Err()
	}
	defer rowId.Close()
	if !rowId.Next() {
		log.Println("No session found for token " + token)
		return 0, errors.New("No session found for token: " + token)
	}
	err := rowId.Scan(&id)
	if err != nil {
		log.Printf("Error while scan id value: %v", rowId.Err())
		return 0, err
	}
	return cast.ToInt(id), nil
}
