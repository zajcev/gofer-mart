package handlers

import (
	"encoding/json"
	"net/http"
)

func (h *Handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("Authorization")
	if token == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	userID, err := getUserID(r.Context(), token, h)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	balance, err := h.db.GetUserBalance(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	resp, err := json.Marshal(&balance)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(resp)
	if err != nil {
		return
	}
}
