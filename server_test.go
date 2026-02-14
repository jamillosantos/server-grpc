package srvgrpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestNewGRPCServer(t *testing.T) {
	wantName := "name"
	wantRegisterer := func(_ *grpc.Server) error { return nil }
	wantOpt1 := WithBindAddress("bind:9090")
	wantOpt2 := WithConnectionTimeout(time.Second)
	s := NewGRPCServer(wantName, wantRegisterer, wantOpt1, wantOpt2)
	assert.Equal(t, wantName, s.name)
	assert.Equal(t, fmt.Sprintf("%p", wantRegisterer), fmt.Sprintf("%p", s.registerer))
	require.Len(t, s.opts, 2)
	assert.Equal(t, fmt.Sprintf("%p", wantOpt1), fmt.Sprintf("%p", s.opts[0]))
	assert.Equal(t, fmt.Sprintf("%p", wantOpt2), fmt.Sprintf("%p", s.opts[1]))
}

func TestGRPCServer_Listen(t *testing.T) {
	t.Run("should start the server", func(t *testing.T) {
		ctx := context.TODO()

		srv := NewGRPCServer("grpc server", func(_ *grpc.Server) error {
			return nil
		}, WithBindAddress("localhost:9091"))

		require.NotNil(t, srv)

		err := srv.IsReady(ctx)
		assert.ErrorIs(t, err, ErrNotReady)

		err = srv.Listen(ctx)
		require.NoError(t, err)

		defer func() {
			_ = srv.Close(ctx)
		}()

		// Wait for server to be ready
		require.Eventually(t, func() bool {
			err := srv.IsReady(ctx)
			return err == nil
		}, time.Second*2, time.Millisecond*100, "server never became ready")

		// Verify we can connect
		conn, err := grpc.NewClient("localhost:9091", grpc.WithTransportCredentials(insecure.NewCredentials()))
		require.NoError(t, err)
		defer func() {
			_ = conn.Close()
		}()
	})

	t.Run("should fail when the port is busy", func(t *testing.T) {
		ctx := context.Background()

		lis, err := net.Listen("tcp", "localhost:9092")
		require.NoError(t, err)
		defer func() {
			_ = lis.Close()
		}()

		srv := NewGRPCServer("grpc server", func(_ *grpc.Server) error {
			return nil
		}, WithBindAddress("localhost:9092"))

		defer func() {
			_ = srv.Close(ctx)
		}()

		err = srv.Listen(ctx)
		assert.Error(t, err)
	})

	t.Run("should fail when the registerer fails", func(t *testing.T) {
		ctx := context.Background()

		wantErr := errors.New("random error")

		srv := NewGRPCServer("grpc server", func(_ *grpc.Server) error {
			return wantErr
		}, WithBindAddress("localhost:9093"))

		defer func() {
			_ = srv.Close(ctx)
		}()

		err := srv.Listen(ctx)
		assert.ErrorIs(t, err, wantErr)
	})
}

func TestGRPCServer_Name(t *testing.T) {
	wantName := "name"
	s := &GRPCServer{
		name: wantName,
	}
	assert.Equal(t, wantName, s.Name())
}
