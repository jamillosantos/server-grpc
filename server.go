package srvgrpc

import (
	"context"
	"net"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

type grpcOpts struct {
	bindAddress        string
	unaryInterceptors  []grpc.UnaryServerInterceptor
	streamInterceptors []grpc.StreamServerInterceptor
	connectionTimeout  time.Duration
}

func defaultOpts() grpcOpts {
	return grpcOpts{
		bindAddress:        ":9090",
		unaryInterceptors:  make([]grpc.UnaryServerInterceptor, 0),
		streamInterceptors: make([]grpc.StreamServerInterceptor, 0),
		connectionTimeout:  time.Second * 15,
	}
}

type Registerer func(*grpc.Server) error

type GRPCServer struct {
	name       string
	server     *grpc.Server
	registerer Registerer
	opts       []Option
}

func NewGRPCServer(name string, registerer Registerer, opts ...Option) *GRPCServer {
	return &GRPCServer{
		name:       name,
		registerer: registerer,
		opts:       opts,
	}
}

// Name will return a human identifiable name for this service. Ex: Postgresql Connection.
func (g *GRPCServer) Name() string {
	return g.name
}

// Listen will start the server and will block until the service is closed.
//
// If the services is already listining, this should return an error ErrAlreadyListening.
func (g *GRPCServer) Listen(_ context.Context) error {
	o := defaultOpts()

	for _, option := range g.opts {
		option(&o)
	}

	lis, err := net.Listen("tcp", o.bindAddress)
	if err != nil {
		return err
	}

	unaryInterceptors := []grpc.UnaryServerInterceptor{
		otelgrpc.UnaryServerInterceptor(),
	}
	unaryInterceptors = append(unaryInterceptors, o.unaryInterceptors...)

	streamInterceptors := []grpc.StreamServerInterceptor{
		otelgrpc.StreamServerInterceptor(),
	}
	streamInterceptors = append(streamInterceptors, o.streamInterceptors...)

	g.server = grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(unaryInterceptors...),
		),
		grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(streamInterceptors...),
		),
		grpc.ConnectionTimeout(o.connectionTimeout),
	)
	err = g.registerer(g.server)
	if err != nil {
		return err
	}

	return g.server.Serve(lis)
}

// Close will stop this server.
//
// If the server has not started, or is already stopped, this should do nothing and just return nil.
func (g *GRPCServer) Close(_ context.Context) error {
	if g.server != nil {
		g.server.GracefulStop()
	}
	return nil
}
