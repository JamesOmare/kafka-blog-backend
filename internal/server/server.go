package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	_ "github.com/joho/godotenv/autoload"

	"kafka-blog-backend/internal/config"
	"kafka-blog-backend/internal/database"
	"kafka-blog-backend/internal/handlers"
	custommiddleware "kafka-blog-backend/internal/middleware"
)

type Server struct {
	port     int
	db       database.Service
	router   chi.Router
	handlers *handlers.Handler
	config   *config.Config
}

func NewServer(cfg *config.Config, db database.Service) *http.Server {
	port, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		port = 8080 // fallback/default port
	}

	h := handlers.New(db, cfg)
	s := &Server{
		port:     port,
		db:       db,
		router:   chi.NewRouter(),
		handlers: h,
		config:   cfg,
	}

	s.setupMiddleware()
	s.setupRoutes()

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.port),
		Handler:      s.setupRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

func (s *Server) setupMiddleware() {
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.RequestID)
	s.router.Use(custommiddleware.CORS)
	s.router.Use(middleware.Timeout(60 * time.Second))
}
