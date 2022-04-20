package srvgrpc

import (
	"time"

	"google.golang.org/grpc"
)

type Option func(o *grpcOpts)

// WithBindAddress sets the bind address used to open the GRPC server.
func WithBindAddress(bindAddress string) Option {
	return func(o *grpcOpts) {
		o.bindAddress = bindAddress
	}
}

// WithUnaryInterceptor set unary interceptors that will be used.
func WithUnaryInterceptor(interceptors ...grpc.UnaryServerInterceptor) Option {
	return func(o *grpcOpts) {
		o.unaryInterceptors = interceptors
	}
}

// WithStreamInterceptor set stream interceptors that will be used.
func WithStreamInterceptor(interceptors ...grpc.StreamServerInterceptor) Option {
	return func(o *grpcOpts) {
		o.streamInterceptors = interceptors
	}
}

// WithConnectionTimeout set the connection timeout using the grpc.ConnectionTimeout option. The default value is
// 15 seconds.
func WithConnectionTimeout(connectionTimeout time.Duration) Option {
	return func(o *grpcOpts) {
		o.connectionTimeout = connectionTimeout
	}
}
