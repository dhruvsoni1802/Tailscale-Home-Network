package client

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/dhruvsoni1802/Tailscale-Home-Network/internal/config"
	"tailscale.com/tsnet"
)

//Client has two HTTP handlers: 1 for React UI and another for API requests to storage node using tsnet
type Client struct {
	tsnet      *tsnet.Server
	httpClient *http.Client
	cfg        config.ClientConfig
}

// New creates a new client instance
func New(ts *tsnet.Server, cfg config.ClientConfig) *Client {
	return &Client{
		tsnet: ts,
		httpClient: ts.HTTPClient(),
		cfg:        cfg,
	}
}

// Start starts the client
func (c *Client) Start() error {
	if err := os.MkdirAll(c.cfg.StateDir, 0755); err != nil {
		return fmt.Errorf("failed to create state dir: %w", err)
	}

	mux := http.NewServeMux()
	c.registerRoutes(mux)

	// bind to localhost only — not accessible outside this device
	ln, err := net.Listen("tcp", c.cfg.LocalPort)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", c.cfg.LocalPort, err)
	}

	log.Printf("client UI available at http://localhost%s", c.cfg.LocalPort)

	return http.Serve(ln, mux)
}
