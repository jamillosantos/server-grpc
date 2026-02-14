// Package main demonstrates a complete service implementation with the helloworld example.
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"

	srvgrpc "github.com/jamillosantos/server-grpc"
)

// GreeterService implements the helloworld.GreeterServer interface
type GreeterService struct {
	pb.UnimplementedGreeterServer
}

func (s *GreeterService) SayHello(_ context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Received request from: %s", req.GetName())
	return &pb.HelloReply{
		Message: "Hello, " + req.GetName() + "!",
	}, nil
}

func NewServer() *srvgrpc.GRPCServer {
	greeter := &GreeterService{}

	return srvgrpc.NewGRPCServer(
		"Greeter Service",
		func(s *grpc.Server) error {
			// Register the service with the gRPC server
			pb.RegisterGreeterServer(s, greeter)
			log.Println("Greeter service registered")
			return nil
		},
		srvgrpc.WithBindAddress(":9090"),
	)
}

func main() {
	srv := NewServer()

	ctx := context.Background()
	if err := srv.Listen(ctx); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

	log.Printf("Server '%s' is listening on :9090", srv.Name())

	// Check readiness
	if err := srv.IsReady(ctx); err != nil {
		log.Printf("Server not ready: %v", err)
	} else {
		log.Println("Server is ready to accept requests")
	}

	// Wait for shutdown signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down gracefully...")
	if err := srv.Close(ctx); err != nil {
		log.Printf("Error closing server: %v", err)
	}
	log.Println("Server stopped")
}
