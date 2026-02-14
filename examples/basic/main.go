// Package main demonstrates basic usage of the server-grpc wrapper.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	srvgrpc "github.com/jamillosantos/server-grpc"
	"google.golang.org/grpc"
)

func main() {
	// Create a new gRPC server
	srv := srvgrpc.NewGRPCServer(
		"My gRPC Server",
		func(_ *grpc.Server) error {
			// Register your gRPC services here
			// pb.RegisterYourServiceServer(s, &yourServiceImpl{})
			return nil
		},
	)

	// Start the server (non-blocking)
	ctx := context.Background()
	if err := srv.Listen(ctx); err != nil {
		log.Fatal(err)
	}

	log.Printf("Server '%s' started successfully", srv.Name())

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down gracefully...")
	if err := srv.Close(ctx); err != nil {
		log.Printf("Error closing server: %v", err)
	}
}
