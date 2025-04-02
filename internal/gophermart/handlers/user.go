package handlers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/zajcev/gofer-mart/internal/gophermart/model"
	"net/http"
	"time"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	respCode := user.AddUser()
	w.WriteHeader(respCode)
	if respCode == http.StatusOK {
		w.Header().Set("Authorization", generateAuthToken())
	}
}
func LoginUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	if !user.VerifyLogin() || !user.VerifyPassword() {
		http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Authorization", generateAuthToken())
	w.WriteHeader(http.StatusOK)
}

func generateAuthToken() string {
	token := sha256.Sum256([]byte(time.Now().String() + "added491416943647d33486a67a0ec8f"))
	return hex.EncodeToString(token[:])
}
