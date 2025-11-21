package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Nikitaannusewicz/carwash-crm/internal/config"
	"github.com/Nikitaannusewicz/carwash-crm/internal/database"
	"github.com/Nikitaannusewicz/carwash-crm/internal/middleware"
	"github.com/Nikitaannusewicz/carwash-crm/internal/modules/identity"
	"github.com/Nikitaannusewicz/carwash-crm/internal/modules/operations"
)

type Server struct {
	cfg    *config.Config
	db     *database.DB
	router *http.ServeMux
}

func NewServer(cfg *config.Config, db *database.DB) *Server {
	mux := http.NewServeMux()

	// 1. Identity Module
	userRepo := identity.NewPostgresRepository(db.DB)
	userService := identity.NewService(userRepo)
	userHandler := identity.NewHandler(userService, cfg.JWTSecret)
	userHandler.RegisterRoutes(mux)

	// 2. Operations Module
	opsRepo := operations.NewPostgresRepository(db.DB)
	opsService := operations.NewService(opsRepo)
	opsHandler := operations.NewHandler(opsService)

	s := &Server{
		cfg:    cfg,
		db:     db,
		router: mux,
	}

	mux.Handle("POST /api/v1/locations", s.AuthMiddleware(http.HandlerFunc(opsHandler.HandleCreateLocation)))
	mux.Handle("POST /api/v1/locations/{id}/bays", s.AuthMiddleware(http.HandlerFunc(opsHandler.HandleCreateBay)))
	mux.Handle("POST /api/v1/service", s.AuthMiddleware(http.HandlerFunc(opsHandler.HandleCreateService)))

	s.registerRoutes()

	return s
}

func (s *Server) registerRoutes() {
	// Std health check
	s.router.HandleFunc("GET /health", s.handleHealthCheck)

	s.router.Handle("GET /api/v1/me", s.AuthMiddleware(http.HandlerFunc(s.handleMe)))
}

func (s *Server) handleMe(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(middleware.UserIDKey)
	role := r.Context().Value(middleware.RoleKey)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, `{"user_id": %v, "role": "%v", "message": "You are authorized!"}`, userID, role)
}

func (s *Server) handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok", "message": "Carwash CRM API is running"}`))
}

func (s *Server) Run() error {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", s.cfg.Port),
		Handler:      s.router,
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	fmt.Printf("Starting server on port %d... \n", s.cfg.Port)
	return server.ListenAndServe()
}
