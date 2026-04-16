package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

type Manager struct {
	baseDir string
}

// FileInfo holds metadata about a stored file
type FileInfo struct {
	Name       string `json:"name"`
	Size       int64  `json:"size"`
	ModifiedAt string `json:"modified_at"`
}

//Helper function to return the storage path for a specific user and creating the storage directory if it doesn't exist
func (m *Manager) userDir(user string) (string, error) {
	dir := filepath.Join(m.baseDir, user)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create user dir: %w", err)
	}

	return dir, nil
}


func NewManager(baseDir string) (*Manager, error) {
	
	//Creating the storage directory if it doesn't exist
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage dir: %w", err)
	}

	return &Manager{baseDir: baseDir}, nil
}

func (m *Manager) Save(user,filename string, src io.Reader) error {
	// Stripping any directory components from filename to prevent path traversal attacks
	safe := filepath.Base(filename)

	if safe == "." || safe == "/" {
		return fmt.Errorf("invalid filename")
	}

	dir, err := m.userDir(user)
	if err != nil {
		return err
	}

	dst, err := os.Create(filepath.Join(dir, safe))
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	
	defer dst.Close()

	// Streaming from request body to disk hence this never loads full file into memory
	if _, err := io.Copy(dst, src); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

func (m *Manager) GetPath(user,filename string) (string, error) {
	safe := filepath.Base(filename)
	if safe == "." || safe == "/" {
		return "", fmt.Errorf("invalid filename")
	}

	dir, err := m.userDir(user)
	if err != nil {
		return "", err
	}

	fullPath := filepath.Join(dir, safe)

	// Checking if the file actually exists before returning the path
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file not found: %s", safe)
	}

	return fullPath, nil
}

func (m *Manager) List(user string) ([]FileInfo, error) {

	dir, err := m.userDir(user)
	if err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read storage dir: %w", err)
	}

	files := []FileInfo{}

	for _, entry := range entries {
		// Skipping subfolders as we only track files
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		files = append(files, FileInfo{
			Name:       info.Name(),
			Size:       info.Size(),
			ModifiedAt: info.ModTime().Format(time.RFC3339),
		})
	}

	return files, nil
}

func (m *Manager) Delete(user,filename string) error {
	safe := filepath.Base(filename)
	if safe == "." || safe == "/" {
		return fmt.Errorf("invalid filename")
	}

	dir, err := m.userDir(user)
	if err != nil {
		return err
	}

	fullPath := filepath.Join(dir, safe)

	if err := os.Remove(fullPath); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("file not found: %s", safe)
		}
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}