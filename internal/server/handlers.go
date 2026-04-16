package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
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
	user := getUser(r)
	device := getDevice(r)

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
	if err := s.storage.Save(user, header.Filename, file); err != nil {
		log.Printf("upload error user: %s, device: %s, filename: %s, error: %v", user, device, header.Filename, err)
		http.Error(w, "failed to save file", http.StatusInternalServerError)
		return
	}

	log.Printf("uploaded by user: %s, device: %s, filename: %s (%d bytes)", user, device, header.Filename, header.Size)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":   "uploaded",
		"filename": header.Filename,
	})
}

func (s *Server) handleDownload(w http.ResponseWriter, r *http.Request) {
	user := getUser(r)
	device := getDevice(r)

	filename := r.PathValue("filename")
	if filename == "" {
		http.Error(w, "filename is required", http.StatusBadRequest)
		return
	}

	// Getting the path of the file from the storage
	filePath, err := s.storage.GetPath(user, filename)
	
	if err != nil {
		// File not found is a client error, not a server error
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	log.Printf("downloading by user: %s, device: %s, filename: %s", user, device, filename)

	// Handling Content-Type, Content-Length, range requests, and streaming
	http.ServeFile(w, r, filePath)
}

func (s *Server) handleList(w http.ResponseWriter, r *http.Request) {
	user := getUser(r)
	device := getDevice(r)
	
	files, err := s.storage.List(user)
	if err != nil {
		log.Printf("list error user: %s, device: %s, error: %v", user, device, err)
		http.Error(w, "failed to list files", http.StatusInternalServerError)
		return
	}

	resp := map[string]any{
		"files": files,
		"total": len(files),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
	user := getUser(r)
	device := getDevice(r)
	
	filename := r.PathValue("filename")
	if filename == "" {
		http.Error(w, "filename is required", http.StatusBadRequest)
		return
	}

	if err := s.storage.Delete(user, filename); err != nil {
		// Distinguishing between not found and a real server error
		if strings.Contains(err.Error(), "file not found") {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		log.Printf("delete error user: %s, device: %s, filename: %s, error: %v", user, device, filename, err)
		http.Error(w, "failed to delete file", http.StatusInternalServerError)
		return
	}

	log.Printf("deleted by user: %s, device: %s, filename: %s", user, device, filename)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status":   "deleted",
		"filename": filename,
	})
}
