package server

import "net/http"

func (s *Server) registerRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /health", s.handleHealth)
	mux.HandleFunc("POST /upload", s.handleUpload)
	mux.HandleFunc("GET /download/{filename}", s.handleDownload)
	
}