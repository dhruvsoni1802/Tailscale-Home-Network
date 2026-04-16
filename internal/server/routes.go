package server

import "net/http"

func (s *Server) registerRoutes(mux *http.ServeMux) {
	// health check needs no auth
	mux.HandleFunc("GET /health", s.handleHealth)

	// all file operations are wrapped with auth middleware
	mux.Handle("POST /upload", s.authMiddleware(http.HandlerFunc(s.handleUpload)))
	mux.Handle("GET /download/{filename}", s.authMiddleware(http.HandlerFunc(s.handleDownload)))
	mux.Handle("GET /files", s.authMiddleware(http.HandlerFunc(s.handleList)))
	mux.Handle("DELETE /files/{filename}", s.authMiddleware(http.HandlerFunc(s.handleDelete)))
}