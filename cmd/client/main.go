package main

import (
	"context"
	"log"
	"os"

	"github.com/dhruvsoni1802/Tailscale-Home-Network/internal/client"
	"github.com/dhruvsoni1802/Tailscale-Home-Network/internal/config"
	"tailscale.com/tsnet"
)

func main() {
	cfg := config.DefaultClient()
	cfg.AuthKey = os.Getenv("TS_AUTHKEY")

	if cfg.AuthKey != "" {
		log.Println("auth key found, registering client node")
	} else {
		log.Println("no auth key, loading from saved state")
	}

	//Creating a tailscale node which acts as a client
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

	log.Printf("client node is up: %s (%v)", cfg.Hostname, status.TailscaleIPs)

	//Creating a new client instance and starting the client UI
	c := client.New(ts, cfg)
	if err := c.Start(); err != nil {
		log.Fatalf("client error: %v", err)
	}

}