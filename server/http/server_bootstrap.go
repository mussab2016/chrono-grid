package http

import (
	"context"
	"io/fs"
	"unisched/server/internal/services"
)

var srv *Server

func StartHTTPServer(host string, port string, distFS fs.FS, svc *services.Services) error {

	// Create HTTPS server
	srv = NewServer(host, port, distFS, svc)
	// Start the server
	return srv.Start()
}

func ShutdownHTTPServer(ctx context.Context) error {
	return srv.Stop(ctx)
}
