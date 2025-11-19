package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Nikitaannusewicz/carwash-crm/internal/config"
	"github.com/Nikitaannusewicz/carwash-crm/internal/database"
	"github.com/Nikitaannusewicz/carwash-crm/internal/modules/identity"
)

type Server struct {
	cfg    *config.Config
	db     *database.DB
	router *http.ServeMux
}

func NewServer(cfg *config.Config, db *database.DB) *Server {
	mux := http.NewServeMux()

	// 1. Create Repo (Pass the raw sql.DB)
	userRepo := identity.NewPostgresRepository(db.DB)
	// 2. Create Service (Inject Repo)
	userService := identity.NewService(userRepo)
	// 3. Create Handler (Inject service)
	userHandler := identity.NewHandler(userService, cfg.JWTSecret)
	// 4. Register Routes
	userHandler.RegisterRoutes(mux)

	s := &Server{
		cfg:    cfg,
		db:     db,
		router: mux,
	}

	s.registerRoutes()

	return s
}

func (s *Server) registerRoutes() {
	// Std health check
	s.router.HandleFunc("GET /health", s.handleHealthCheck)

	s.router.Handle("GET /api/v1/me", s.AuthMiddleware(http.HandlerFunc(s.handleMe)))
}

func (s *Server) handleMe(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey)
	role := r.Context().Value(roleKey)

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
