package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/glebarez/go-sqlite"
	"gitlab.com/robertpelloni/HyperNexus/internal/cloud"
)

func main() {
	// Parse command line flags
	port := flag.Int("port", 7778, "Port to listen on")
	dbPath := flag.String("db", "/data/hypernexus.db", "Path to database")
	backupDir := flag.String("backup-dir", "/var/backups/hypernexus", "Path to backup directory")
	flag.Parse()

	// Initialize database
	db, err := sql.Open("sqlite", *dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	// Create cloud server
	server := cloud.NewCloudServer(db, *backupDir)

	// Initialize server
	if err := server.Initialize(); err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	// Create router
	mux := http.NewServeMux()

	// Register routes
	server.RegisterRoutes(mux)

	// Create HTTP server
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", *port),
		Handler: mux,
	}

	// Start server
	go func() {
		log.Printf("[Cloud] Starting HyperNexus Cloud on port %d", *port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("[Cloud] Shutting down server...")
	httpServer.Close()
	log.Println("[Cloud] Server stopped")
}
