package handlers

import "github.com/zajcev/gofer-mart/internal/gophermart/database"

type Handler struct {
	db *database.DBService
}

func NewHandler(db *database.DBService) *Handler {
	return &Handler{db: db}
}
