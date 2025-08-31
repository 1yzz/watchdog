package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"watchdog/api"
	"watchdog/config"
	"watchdog/server"
)

func main() {
	cfg, db, err := config.LoadWithEntClient()
	if err != nil {
		log.Fatalf("Failed to load config and connect to database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	port := fmt.Sprintf(":%d", cfg.Server.Port)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()

	watchdogServer := server.NewWatchdogServer(db)
	api.RegisterWatchdogServiceServer(s, watchdogServer)

	reflection.Register(s)

	// Check if running in service mode (non-interactive)
	isService := os.Getenv("WATCHDOG_SERVICE_MODE") == "1" || !isTerminal()

	go func() {
		fmt.Printf("gRPC server listening at %v\n", lis.Addr())
		fmt.Println("Database configuration:")
		fmt.Printf("  Host: %s:%d\n", cfg.Database.Host, cfg.Database.Port)
		fmt.Printf("  Database: %s\n", cfg.Database.Database)
		fmt.Printf("  Username: %s\n", cfg.Database.Username)

		if err := s.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	if isService {
		// Service mode: run indefinitely
		fmt.Println("Running in service mode - server will continue running")
		select {} // Block forever
	} else {
		// Interactive mode: wait for signals
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)

		<-c
		fmt.Println("\nShutting down gRPC server...")
		s.GracefulStop()
		fmt.Println("Server stopped")
	}
}

// isTerminal checks if the process is running in a terminal
func isTerminal() bool {
	fileInfo, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return (fileInfo.Mode() & os.ModeCharDevice) != 0
}
