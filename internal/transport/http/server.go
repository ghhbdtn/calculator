package http

import (
	"context"
	"embed"
	"log"
	"net/http"
	"time"

	"calculator/internal/app/service"
	"github.com/gorilla/mux"
)

//go:embed *
var FS embed.FS

type Server struct {
	router     *mux.Router
	calculator *service.Calculator
	httpServer *http.Server
}

func NewServer(calculator *service.Calculator) *Server {
	s := &Server{
		router:     mux.NewRouter(),
		calculator: calculator,
	}

	handler := NewHandler(calculator)
	handler.RegisterRoutes(s.router)
	s.registerSwaggerRoutes()

	return s
}

func (s *Server) registerSwaggerRoutes() {
	// Swagger UI
	s.router.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/",
		http.FileServer(http.FS(FS))))

	// Swagger JSON
	s.router.HandleFunc("/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		file, _ := FS.ReadFile("swagger.json")
		w.Header().Set("Content-Type", "application/json")
		w.Write(file)
	})
}

func (s *Server) Start(addr string) error {
	s.httpServer = &http.Server{
		Addr:         addr,
		Handler:      s.router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("HTTP server starting on %s", addr)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	if s.httpServer != nil {
		return s.httpServer.Shutdown(ctx)
	}
	return nil
}
