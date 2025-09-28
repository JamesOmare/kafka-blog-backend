package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"blog-api/internal/config"
	"blog-api/internal/database"
)

type Handler struct {
	db     database.Service
	config *config.Config
}

func New(db database.Service, config *config.Config) *Handler {
	return &Handler{
		db:     db,
		config: config,
	}
}

func (h *Handler) jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func (h *Handler) errorResponse(w http.ResponseWriter, status int, message string) {
	h.jsonResponse(w, status, map[string]string{
		"error": message,
	})
}

func (h *Handler) getUserIDFromContext(ctx context.Context) int32 {
	if userID, ok := ctx.Value("userID").(int32); ok {
		return userID
	}
	return 0
}

func (h *Handler) generateSlug(title string) string {
	// Simple slug generation - you might want to use a proper slug library
	slug := strings.ToLower(title)
	slug = strings.ReplaceAll(slug, " ", "-")
	return slug
}
