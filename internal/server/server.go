package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dhruvsoni1802/Tailscale-Home-Network/internal/storage"
	"tailscale.com/client/local"
	"tailscale.com/tsnet"
)

// Server is the main server struct
type Server struct {
	tsnet       *tsnet.Server
	http        *http.Server
	storage     *storage.Manager
	localClient *local.Client
}

// New creates a new server instance
func New(ts *tsnet.Server, store *storage.Manager) (*Server, error) {

	// Getting the local client
	lc, err := ts.LocalClient()
	if err != nil {
		return nil, fmt.Errorf("failed to get local client: %w", err)
	}
	return &Server{
		tsnet:       ts,
		storage:     store,
		localClient: lc,
	}, nil
}

// Start starts the server
func (s *Server) Start() error {
	// listen only on the tailnet interface
	ln, err := s.tsnet.Listen("tcp", ":8080")
	if err != nil {
		return fmt.Errorf("failed to listen on tailnet: %w", err)
	}

	mux := http.NewServeMux()
	s.registerRoutes(mux)

	s.http = &http.Server{Handler: mux}

	log.Println("HTTP server listening on :8080")

	// blocks until the server stops
	return s.http.Serve(ln)
}