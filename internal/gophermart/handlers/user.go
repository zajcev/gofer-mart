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

func (h *Handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	respCode, token, err := addUser(ctx, user, h)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	defer cancel()
	if respCode == http.StatusOK {
		w.Header().Set("Authorization", token)
	}
	w.WriteHeader(respCode)
}

func (h *Handler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	token, respCode, err := auth(ctx, user, h)
	if err != nil {
		http.Error(w, err.Error(), respCode)
	}
	w.Header().Set("Authorization", token)
	w.WriteHeader(respCode)
}

func auth(ctx context.Context, u model.User, h *Handler) (string, int, error) {
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

func addUser(ctx context.Context, u model.User, h *Handler) (int, string, error) {
	if verifyLogin(ctx, u, h) {
		return http.StatusConflict, "", errors.New("login failed")
	}
	hash, err := hashPassword(u.Password)
	if err != nil {
		return http.StatusInternalServerError, "", err
	}
	h.db.AddUser(ctx, u.Login, hash)
	token := generateAuthToken()
	err = h.db.NewSession(ctx, u.Login, token)
	return http.StatusOK, token, err
}

func verifyPassword(ctx context.Context, u model.User, h *Handler) bool {
	dbPass := h.db.GetPassword(ctx, u.Login)
	hash, _ := hashPassword(u.Password)
	return dbPass == hash
}

func verifyLogin(ctx context.Context, u model.User, h *Handler) bool {
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
