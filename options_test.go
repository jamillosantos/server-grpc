package srvgrpc

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestWithUnaryInterceptor(t *testing.T) {
	wantInterceptor := func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		return nil, nil
	}
	o := grpcOpts{}
	WithUnaryInterceptor(wantInterceptor)(&o)
	require.Len(t, o.unaryInterceptors, 1)
	assert.Equal(t, fmt.Sprintf("%p", wantInterceptor), fmt.Sprintf("%p", o.unaryInterceptors[0]))
}

func TestWithStreamInterceptor(t *testing.T) {
	wantInterceptor := func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		return nil
	}
	o := grpcOpts{}
	WithStreamInterceptor(wantInterceptor)(&o)
	require.Len(t, o.streamInterceptors, 1)
	assert.Equal(t, fmt.Sprintf("%p", wantInterceptor), fmt.Sprintf("%p", o.streamInterceptors[0]))
}

func TestWithConnectionTimeout(t *testing.T) {
	wantConnectionTimeout := time.Second * 123
	o := grpcOpts{}
	WithConnectionTimeout(wantConnectionTimeout)(&o)
	assert.Equal(t, wantConnectionTimeout, o.connectionTimeout)
}
