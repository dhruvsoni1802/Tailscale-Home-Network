package config

import (
	"os"
	"path/filepath"
)

type ServerConfig struct {
	Hostname   string
	StateDir   string
	AuthKey    string
	StorageDir string
}

type ClientConfig struct {
	Hostname    string
	StateDir    string
	AuthKey     string
	StorageNode string
	LocalPort   string
}


func DefaultServer() ServerConfig {
	return ServerConfig{
		Hostname:   "storage-node",
		StateDir:   "./tsnet-state",
		StorageDir: "./storage-data",
	}
}

func DefaultClient() ClientConfig {
	home, _ := os.UserHomeDir()


	return ClientConfig{
		Hostname:    "tailstore-client",
		StateDir:    filepath.Join(home, ".tailstore", "tsnet-state"),
		StorageNode: "http://storage-node:8080",
		LocalPort:   ":4000",
	}
}