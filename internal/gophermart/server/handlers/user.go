package handlers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/zajcev/gofer-mart/internal/gophermart/model"
	"net/http"
	"time"
)

type AuthProvider interface {
	GetUserIDByToken(ctx context.Context, token string) (int, error)
}

type AuthStorage struct {
	DB AuthProvider
}

func NewAuthStorage(db AuthProvider) AuthStorage {
	return AuthStorage{DB: db}
}

type UserStorage interface {
	AddUser(ctx context.Context, login string, pass string)
	GetLogin(ctx context.Context, login string) string
	GetPassword(ctx context.Context, login string) string
	NewSession(ctx context.Context, login string, token string) error
	GetUserID(ctx context.Context, login string) (int, error)
	GetUserIDByToken(ctx context.Context, token string) (int, error)
}
type UserHandler struct {
	db UserStorage
}

func NewUserHandler(db UserStorage) *UserHandler {
	return &UserHandler{db: db}
}

func (uh *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	respCode, token, err := addUser(ctx, user, uh)
	if err != nil {
		http.Error(w, err.Error(), respCode)
		return
	}

	if respCode == http.StatusOK {
		w.Header().Set("Authorization", token)
	}
	w.WriteHeader(respCode)
}

func (uh *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	token, respCode, err := auth(ctx, user, uh)
	if err != nil {
		http.Error(w, err.Error(), respCode)
	}
	w.Header().Set("Authorization", token)
	w.WriteHeader(respCode)
}

func auth(ctx context.Context, u model.User, h *UserHandler) (string, int, error) {
	if verifyLogin(ctx, u, h) && verifyPassword(ctx, u, h) {
		token := generateAuthToken()
		err := h.db.NewSession(ctx, u.Login, token)
		if err != nil {
			return "", http.StatusInternalServerError, err
		}
		return token, http.StatusOK, nil
	} else {
		return "", http.StatusUnauthorized, errors.New("login failed")
	}
}

func addUser(ctx context.Context, u model.User, h *UserHandler) (int, string, error) {
	if verifyLogin(ctx, u, h) {
		return http.StatusConflict, "", errors.New("user already exists")
	}
	hash, err := hashPassword(u.Password)
	if err != nil {
		return http.StatusInternalServerError, "", err
	}
	h.db.AddUser(ctx, u.Login, hash)
	token := generateAuthToken()
	err = h.db.NewSession(ctx, u.Login, token)
	if err != nil {
		return http.StatusInternalServerError, "", err
	}
	return http.StatusOK, token, nil
}

func verifyPassword(ctx context.Context, u model.User, h *UserHandler) bool {
	dbPass := h.db.GetPassword(ctx, u.Login)
	hash, _ := hashPassword(u.Password)
	return dbPass == hash
}

func verifyLogin(ctx context.Context, u model.User, h *UserHandler) bool {
	login := h.db.GetLogin(ctx, u.Login)
	return login != ""
}

func generateAuthToken() string {
	token := sha256.Sum256([]byte(time.Now().String() + "added491416943647d33486a67a0ec8f"))
	return hex.EncodeToString(token[:])
}

func hashPassword(password string) (string, error) {
	h := sha256.New()
	h.Write([]byte(password))
	return hex.EncodeToString(h.Sum(nil)), nil
}
