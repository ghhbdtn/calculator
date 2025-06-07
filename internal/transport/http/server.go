package http

import (
	"context"
	"log"
	"net/http"
	"time"

	"calculator/internal/app/service"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

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
	swaggerHandler := httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // URL к документации
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)

	// Маршрут для Swagger UI
	s.router.PathPrefix("/swagger/").Handler(swaggerHandler)
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
