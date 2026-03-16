// Package srvgrpc provides a gRPC server wrapper with support for interceptors, metrics, and tracing.
package srvgrpc

import (
	"context"
	"errors"
	"net"
	"sync"
	"sync/atomic"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

// ErrNotReady is returned when the server is not ready to accept requests.
var ErrNotReady = errors.New("service is not ready")

type grpcOpts struct {
	bindAddress        string
	unaryInterceptors  []grpc.UnaryServerInterceptor
	streamInterceptors []grpc.StreamServerInterceptor
	connectionTimeout  time.Duration
}

func defaultOpts() grpcOpts {
	return grpcOpts{
		bindAddress:        ":50051",
		unaryInterceptors:  make([]grpc.UnaryServerInterceptor, 0),
		streamInterceptors: make([]grpc.StreamServerInterceptor, 0),
		connectionTimeout:  time.Second * 15,
	}
}

// Registerer is a function that registers gRPC services with the server.
type Registerer func(*grpc.Server) error

// GRPCServer is a wrapper around a gRPC server that implements the application.Service interface.
type GRPCServer struct {
	name        string
	server      *grpc.Server
	registerer  Registerer
	opts        []Option
	bindAddress string
	ready       atomic.Value
	serverWg    sync.WaitGroup
}

// NewGRPCServer creates a new gRPC server with the given name, service registerer, and options.
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
func (g *GRPCServer) Listen(ctx context.Context) error {
	o := defaultOpts()

	for _, option := range g.opts {
		option(&o)
	}

	lis, err := net.Listen("tcp", o.bindAddress)
	if err != nil {
		return err
	}
	g.bindAddress = lis.Addr().String()

	srvMetrics := prometheus.NewServerMetrics()

	unaryInterceptors := []grpc.UnaryServerInterceptor{
		srvMetrics.UnaryServerInterceptor(),
	}
	unaryInterceptors = append(unaryInterceptors, o.unaryInterceptors...)

	streamInterceptors := []grpc.StreamServerInterceptor{
		srvMetrics.StreamServerInterceptor(),
	}
	streamInterceptors = append(streamInterceptors, o.streamInterceptors...)

	g.server = grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
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

	g.serverWg.Add(1)
	go func() {
		defer func() {
			g.serverWg.Done()
			g.ready.Store(false)
		}()
		_ = g.server.Serve(lis)
	}()

	go func() {
		dialer := &net.Dialer{}
		for {
			conn, err := dialer.DialContext(ctx, "tcp", g.bindAddress)
			if err == nil {
				_ = conn.Close()
				g.ready.Store(true)
				return
			}
			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}()

	return nil
}

// Close will stop this server.
//
// If the server has not started, or is already stopped, this should do nothing and just return nil.
func (g *GRPCServer) Close(_ context.Context) error {
	if g.server != nil {
		g.server.GracefulStop()
	}
	g.serverWg.Wait()
	return nil
}

// IsReady will return true if the service is ready to accept requests. This is compliant with the
// github.com/jamillosantos/application library.
func (g *GRPCServer) IsReady(_ context.Context) error {
	if v := g.ready.Load(); v == nil || v == false {
		return ErrNotReady
	}
	return nil
}
