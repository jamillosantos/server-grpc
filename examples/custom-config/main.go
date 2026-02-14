// Package main demonstrates custom configuration with interceptors.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	srvgrpc "github.com/jamillosantos/server-grpc"
	"google.golang.org/grpc"
)

func main() {
	srv := srvgrpc.NewGRPCServer(
		"Custom gRPC Server",
		func(_ *grpc.Server) error {
			// Register your services
			return nil
		},
		// Configure bind address
		srvgrpc.WithBindAddress(":9090"),

		// Configure connection timeout
		srvgrpc.WithConnectionTimeout(30*time.Second),

		// Add custom unary interceptors
		srvgrpc.WithUnaryInterceptor(myUnaryInterceptor),

		// Add custom stream interceptors
		srvgrpc.WithStreamInterceptor(myStreamInterceptor),
	)

	ctx := context.Background()
	if err := srv.Listen(ctx); err != nil {
		log.Fatal(err)
	}

	log.Printf("Server '%s' listening on :9090", srv.Name())

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down gracefully...")
	if err := srv.Close(ctx); err != nil {
		log.Printf("Error closing server: %v", err)
	}
}

func myUnaryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("Unary call to: %s", info.FullMethod)
	return handler(ctx, req)
}

func myStreamInterceptor(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	log.Printf("Stream call to: %s", info.FullMethod)
	return handler(srv, ss)
}
