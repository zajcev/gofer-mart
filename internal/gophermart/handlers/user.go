package handlers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/zajcev/gofer-mart/internal/gophermart/database"
	"github.com/zajcev/gofer-mart/internal/gophermart/model"
	"log"
	"net/http"
	"time"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	respCode, token := addUser(ctx, user)
	defer cancel()
	if respCode == http.StatusOK {
		w.Header().Set("Authorization", token)
	}
	w.WriteHeader(respCode)
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()

	token, respCode := auth(ctx, user)
	if respCode == http.StatusOK {
		w.Header().Set("Authorization", token)
		w.WriteHeader(respCode)
	} else {
		w.WriteHeader(respCode)
	}
}

func auth(ctx context.Context, u model.User) (string, int) {
	if verifyLogin(ctx, u) && verifyPassword(ctx, u) {
		token := generateAuthToken()
		database.NewSession(ctx, u.Login, token)
		return token, http.StatusOK
	} else {
		return "", http.StatusUnauthorized
	}
}

func addUser(ctx context.Context, u model.User) (int, string) {
	if verifyLogin(ctx, u) {
		return http.StatusConflict, ""
	}
	hash, err := hashPassword(u.Password)
	if err != nil {
		log.Printf("error hashing password: %v", err)
		return http.StatusInternalServerError, ""
	}
	database.AddUser(ctx, u.Login, hash)
	token := generateAuthToken()
	database.NewSession(ctx, u.Login, token)
	return http.StatusOK, token
}

func verifyPassword(ctx context.Context, u model.User) bool {
	dbPass := database.GetPassword(ctx, u.Login)
	hash, _ := hashPassword(u.Password)
	if dbPass == hash {
		return true
	}
	return false
}

func verifyLogin(ctx context.Context, u model.User) bool {
	login := database.GetLogin(ctx, u.Login)
	if login == "" {
		return false
	}
	return true
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
