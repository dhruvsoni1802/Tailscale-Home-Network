package client

import (
	"net/http"

	"github.com/dhruvsoni1802/Tailscale-Home-Network/ui"
)

func (c *Client) registerRoutes(mux *http.ServeMux) {
	mux.Handle("GET /api/files", corsMiddleware(http.HandlerFunc(c.handleList)))
	mux.Handle("POST /api/upload", corsMiddleware(http.HandlerFunc(c.handleUpload)))
	mux.Handle("GET /api/download/{filename}", corsMiddleware(http.HandlerFunc(c.handleDownload)))
	mux.Handle("DELETE /api/files/{filename}", corsMiddleware(http.HandlerFunc(c.handleDelete)))
	mux.Handle("GET /api/health", corsMiddleware(http.HandlerFunc(c.handleHealth)))

	// serve the embedded React UI for everything that isn't an API route
	mux.Handle("/", ui.Handler())
}