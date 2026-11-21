package http

import (
	"context"
	"fmt"
	"io/fs"
	"net/http"
	"time"
	"unisched/server/internal/services"
)

// Server holds the HTTP server
type Server struct {
	srv *http.Server
}

// NewServer creates a new HTTPS server
func NewServer(host string, port string, distFS fs.FS, svc *services.Services) *Server {

	router := NewRouter(distFS, svc)

	srv := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", host, port),
		Handler:      router,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 600 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{srv: srv}
}

// Start runs the server
func (s *Server) Start() error {
	return s.srv.ListenAndServe() //return s.srv.ListenAndServeTLS(s.certFile, s.keyFile) HTTPS for alone server. And HTTP for nginx
}

// Stop gracefully shuts down the server
func (s *Server) Stop(ctx context.Context) error {
	if s.srv != nil {
		return s.srv.Shutdown(ctx)
	}
	return nil
}
