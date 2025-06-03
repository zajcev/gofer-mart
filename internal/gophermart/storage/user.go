package storage

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/spf13/cast"
	"github.com/zajcev/gofer-mart/internal/gophermart/storage/scripts"
	"log"
)

func (s *DBService) AddUser(ctx context.Context, login string, pass string) {
	_, err := s.db.Exec(ctx, scripts.AddUser, login, pass)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			log.Printf("PGError: %v", pgErr)
		} else {
			log.Printf("Error while addUser: Error=%v, login=%v", err, login)
		}
	}
}

func (s *DBService) GetLogin(ctx context.Context, login string) string {
	row, err := s.db.Query(ctx, scripts.GetLogin, login)
	if err != nil {
		return ""
	}
	defer row.Close()

	if !row.Next() {
		return ""
	}

	if err = row.Err(); err != nil {
		return ""
	}

	var value interface{}
	err = row.Scan(&value)
	if err != nil {
		return ""
	}

	return cast.ToString(value)
}

func (s *DBService) GetPassword(ctx context.Context, login string) string {
	row, err := s.db.Query(ctx, scripts.GetPassword, login)
	if err != nil {
		return ""
	}
	defer row.Close()

	if !row.Next() {
		return ""
	}

	if err = row.Err(); err != nil {
		return ""
	}

	var value interface{}
	err = row.Scan(&value)
	if err != nil {
		return ""
	}

	return cast.ToString(value)
}

func (s *DBService) NewSession(ctx context.Context, login string, token string) error {
	var id int
	row := s.db.QueryRow(ctx, scripts.GetUserIDByLogin, login)
	err := row.Scan(&id)
	if err != nil {
		return err
	}
	if id != 0 {
		_, err := s.db.Exec(ctx, scripts.AddSession, id, token)
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

func (s *DBService) GetUserID(ctx context.Context, login string) (int, error) {
	var id interface{}
	rowID, _ := s.db.Query(ctx, scripts.GetUserIDByLogin, login)
	if rowID.Err() != nil {
		return 0, rowID.Err()
	}
	defer rowID.Close()
	if !rowID.Next() {
		return 0, errors.New("No tokens returned for login: " + login)
	}
	err := rowID.Scan(&id)
	if err != nil {
		return 0, err
	}
	return cast.ToInt(id), nil
}

func (s *DBService) GetUserIDByToken(ctx context.Context, token string) (int, error) {
	var id interface{}
	rowID, _ := s.db.Query(ctx, scripts.GetUserIDByToken, token)
	if rowID.Err() != nil {
		return 0, rowID.Err()
	}
	defer rowID.Close()
	if !rowID.Next() {
		return 0, errors.New("No session found for token: " + token)
	}
	err := rowID.Scan(&id)
	if err != nil {
		return 0, err
	}
	return cast.ToInt(id), nil
}
