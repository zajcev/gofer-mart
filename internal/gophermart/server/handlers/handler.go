package handlers

import "github.com/zajcev/gofer-mart/internal/gophermart/storage"

type Handler struct {
	db *storage.DBService
}

func NewHandler(db *storage.DBService) *Handler {
	return &Handler{db: db}
}
