package client

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

//Helper function to forward a request to the storage server and pipes the response back
func (c *Client) forward(w http.ResponseWriter, r *http.Request, method, path string, body io.Reader, contentType string) {
	url := fmt.Sprintf("%s%s", c.cfg.StorageNode, path)

	req, err := http.NewRequestWithContext(r.Context(), method, url, body)
	if err != nil {
		http.Error(w, "failed to build request", http.StatusInternalServerError)
		return
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	// httpClient routes this through `the tsnet tunnel to storage node
	resp, err := c.httpClient.Do(req)
	if err != nil {
		log.Printf("forward error: %v", err)
		http.Error(w, "failed to reach storage server", http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// pipe response headers and body straight back to the client UI
	w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}

func (c *Client) handleHealth(w http.ResponseWriter, r *http.Request) {
	c.forward(w, r, http.MethodGet, "/health", nil, "")
}

func (c *Client) handleList(w http.ResponseWriter, r *http.Request) {
	c.forward(w, r, http.MethodGet, "/files", nil, "")
}

func (c *Client) handleDownload(w http.ResponseWriter, r *http.Request) {
	filename := r.PathValue("filename")
	c.forward(w, r, http.MethodGet, "/download/"+filename, nil, "")
}

func (c *Client) handleDelete(w http.ResponseWriter, r *http.Request) {
	filename := r.PathValue("filename")
	c.forward(w, r, http.MethodDelete, "/files/"+filename, nil, "")
}

func (c *Client) handleUpload(w http.ResponseWriter, r *http.Request) {
	// pipe the multipart body straight to the storage node. No need to parse it, we simply forward the raw bytes.
	c.forward(w, r, http.MethodPost, "/upload", r.Body, r.Header.Get("Content-Type"))
}