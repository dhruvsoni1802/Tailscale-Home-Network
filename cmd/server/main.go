package main

import (
	"context"
	"log"
	"os"

	"github.com/dhruvsoni1802/Tailscale-Home-Network/internal/config"
	"github.com/dhruvsoni1802/Tailscale-Home-Network/internal/server"
	"github.com/dhruvsoni1802/Tailscale-Home-Network/internal/storage"
	"tailscale.com/tsnet"
)


func main() {
	cfg := config.Default()
	cfg.AuthKey = os.Getenv("TS_AUTHKEY")

	if cfg.AuthKey != "" {
		log.Println("auth key found, registering node")
	} else {
		log.Println("no auth key, loading from saved state")
	}

	//Creating a tailscale node which acts as a storage server
	ts := &tsnet.Server{
		Hostname: cfg.Hostname,
		Dir:      cfg.StateDir,
		AuthKey:  cfg.AuthKey,
	}
	
	defer ts.Close()

	//Until the node has fully joined the tailnet network, we need to block
	status, err := ts.Up(context.Background())
	if err != nil {
		log.Fatalf("failed to join tailnet: %v", err)
	}

	log.Printf("node is up: %s (%v)", cfg.Hostname, status.TailscaleIPs)

	//Creating a storage manager
	store, err := storage.NewManager(cfg.StorageDir)
	if err != nil {
		log.Fatalf("failed to init storage: %v", err)
	}

	//Starting an HTTP server on top of tailnet node
	srv := server.New(ts, store)
	if err := srv.Start(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}