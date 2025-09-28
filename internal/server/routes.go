package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"

	custommiddleware "kafka-blog-backend/internal/middleware"
)

func (s *Server) setupRoutes() http.Handler {
	r := chi.NewRouter()

	// Apply global middleware
	s.setupMiddleware()

	// Get JWT auth instance from auth service
	tokenAuth := s.authService.GetTokenAuth()

	// For debugging/example purposes, generate and print a sample jwt token
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{
		"user_id":  123,
		"username": "testuser",
		"email":    "test@example.com",
		"role":     "user",
	})
	fmt.Printf("DEBUG: a sample jwt is %s\n\n", tokenString)

	r.Route("/api/v1", func(r chi.Router) {
		// Public routes - No JWT required
		r.Group(func(r chi.Router) {
			// Auth routes
			r.Post("/auth/register", s.handlers.Auth.Register)
			r.Post("/auth/login", s.handlers.Auth.Login)

			// Public blog routes
			// r.Get("/posts", s.handlers.GetPublishedPosts)
			// r.Get("/posts/{slug}", s.handlers.GetPostBySlug)
			// r.Get("/posts/{id}/comments", s.handlers.GetPostComments)
			// r.Get("/tags", s.handlers.GetTags)
		})

		// Protected routes - JWT required
		r.Group(func(r chi.Router) {
			// Step 1: Verifier middleware - extracts, decodes, verifies JWT token
			// Sets jwtauth.TokenCtxKey and jwtauth.ErrorCtxKey in context
			// Searches for token in: 1) Authorization: Bearer header, 2) jwt cookie
			r.Use(jwtauth.Verifier(tokenAuth))

			// Step 2: Authenticator middleware - responds with 401 for invalid tokens
			// Allows valid tokens to pass through to handlers
			r.Use(jwtauth.Authenticator(tokenAuth))

			// User profile routes
			// r.Get("/profile", s.handlers.GetProfile)
			// r.Put("/profile", s.handlers.UpdateProfile)

			// Comment routes (authenticated users)
			// r.Post("/posts/{id}/comments", s.handlers.CreateComment)
			// r.Put("/comments/{id}", s.handlers.UpdateComment)
			// r.Delete("/comments/{id}", s.handlers.DeleteComment)
		})

		// Author/Admin routes - JWT + Role required
		r.Group(func(r chi.Router) {
			// Step 1: Verify JWT token
			r.Use(jwtauth.Verifier(tokenAuth))

			// Step 2: Authenticate JWT token (401 for invalid)
			r.Use(jwtauth.Authenticator(tokenAuth))

			// Step 3: Check user role (403 for insufficient permissions)
			r.Use(custommiddleware.RequireRole("author"))

			// Post management
			// r.Get("/admin/posts", s.handlers.GetUserPosts)
			// r.Post("/posts", s.handlers.CreatePost)
			// r.Put("/posts/{id}", s.handlers.UpdatePost)
			// r.Delete("/posts/{id}", s.handlers.DeletePost)

			// Tag management
			// r.Post("/tags", s.handlers.CreateTag)
			// r.Put("/tags/{id}", s.handlers.UpdateTag)
		})
	})

	r.Get("/", s.HelloWorldHandler)
	r.Get("/health", s.healthHandler)

	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}
