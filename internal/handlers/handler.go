// internal/handlers/handler.go
package handlers

import (
	"kafka-blog-backend/internal/config"
	"kafka-blog-backend/internal/database"
	"kafka-blog-backend/internal/services"
)

// Handler is the main handler struct that contains all sub-handlers
type Handlers struct {
	Auth *AuthHandler
	// Post    *PostHandler
	// Comment *CommentHandler
	// Tag     *TagHandler
}

func New(db database.Service, cfg *config.Config, authService *services.AuthService) *Handlers {
	return &Handlers{
		Auth: NewAuthHandler(authService),
		// Post:    NewPostHandler(db),
		// Comment: NewCommentHandler(db),
		// Tag:     NewTagHandler(db),
	}
}
