package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{
		"status": "ok",
		"time":   time.Now().Format(time.RFC3339),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleUpload(w http.ResponseWriter, r *http.Request) {
	// 32MB max in memory, rest streamed to disk. This is to prevent the server from running out of memory.
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "failed to parse form", http.StatusBadRequest)
		return
	}

	// Getting the file from the form data
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file field missing", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Saving the file to the storage which in our case if the local disk
	if err := s.storage.Save(header.Filename, file); err != nil {
		log.Printf("upload error: %v", err)
		http.Error(w, "failed to save file", http.StatusInternalServerError)
		return
	}

	log.Printf("uploaded: %s (%d bytes)", header.Filename, header.Size)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":   "uploaded",
		"filename": header.Filename,
	})
}

func (s *Server) handleDownload(w http.ResponseWriter, r *http.Request) {
	filename := r.PathValue("filename")
	if filename == "" {
		http.Error(w, "filename is required", http.StatusBadRequest)
		return
	}

	// Getting the path of the file from the storage
	filePath, err := s.storage.GetPath(filename)
	
	if err != nil {
		// File not found is a client error, not a server error
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	log.Printf("downloading: %s", filename)

	// Handling Content-Type, Content-Length, range requests, and streaming
	http.ServeFile(w, r, filePath)
}


