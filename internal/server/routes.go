package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	custommiddleware "kafka-blog-backend/internal/middleware"
)

func (s *Server) setupRoutes() http.Handler {
	r := chi.NewRouter()
	// r.Use(middleware.Logger)

	// r.Use(cors.Handler(cors.Options{
	// 	AllowedOrigins:   []string{"https://*", "http://*"},
	// 	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
	// 	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
	// 	AllowCredentials: true,
	// 	MaxAge:           300,
	// }))
	r.Route("/api/v1", func(r chi.Router) {
		// Public routes
		r.Group(func(r chi.Router) {
			// Auth routes
			r.Post("/auth/register", s.handlers.Register)
			r.Post("/auth/login", s.handlers.Login)

			// Public blog routes
			r.Get("/posts", s.handlers.GetPublishedPosts)
			r.Get("/posts/{slug}", s.handlers.GetPostBySlug)
			r.Get("/posts/{id}/comments", s.handlers.GetPostComments)
			r.Get("/tags", s.handlers.GetTags)
		})

		// Protected routes
		r.Group(func(r chi.Router) {
			r.Use(custommiddleware.JWTAuth(s.config.JWTSecret))

			// User profile routes
			r.Get("/profile", s.handlers.GetProfile)
			r.Put("/profile", s.handlers.UpdateProfile)

			// Comment routes (authenticated users)
			r.Post("/posts/{id}/comments", s.handlers.CreateComment)
			r.Put("/comments/{id}", s.handlers.UpdateComment)
			r.Delete("/comments/{id}", s.handlers.DeleteComment)
		})

		// Author/Admin routes
		r.Group(func(r chi.Router) {
			r.Use(custommiddleware.JWTAuth(s.config.JWTSecret))
			r.Use(custommiddleware.RequireRole("author"))

			// Post management
			r.Get("/admin/posts", s.handlers.GetUserPosts)
			r.Post("/posts", s.handlers.CreatePost)
			r.Put("/posts/{id}", s.handlers.UpdatePost)
			r.Delete("/posts/{id}", s.handlers.DeletePost)

			// Tag management
			r.Post("/tags", s.handlers.CreateTag)
			r.Put("/tags/{id}", s.handlers.UpdateTag)
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
