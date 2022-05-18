package srvgrpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestNewGRPCServer(t *testing.T) {
	wantName := "name"
	wantRegisterer := func(server *grpc.Server) error { return nil }
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

		srv := NewGRPCServer("grpc server", func(s *grpc.Server) error {
			return nil
		}, WithBindAddress("localhost:9091"))

		require.NotNil(t, srv)

		err := srv.IsReady(ctx)
		assert.ErrorIs(t, err, ErrNotReady)

		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()
			defer srv.Close(ctx)

			require.Eventually(t, func() bool {
				ctx, cancelFunc := context.WithTimeout(ctx, time.Second)
				defer cancelFunc()

				conn, err := grpc.DialContext(ctx, "localhost:9091", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithReturnConnectionError())
				if err != nil {
					return false
				}
				defer func() {
					_ = conn.Close()
				}()

				return true
			}, time.Second, time.Millisecond*100, "server was could not be connected")

		}()

		err = srv.Listen(ctx)
		require.NoError(t, err)

		err = srv.IsReady(ctx)
		assert.NoError(t, err)

		wg.Wait()
	})

	t.Run("should fail when the port is busy", func(t *testing.T) {
		ctx := context.Background()

		lis, err := net.Listen("tcp", "localhost:9091")
		require.NoError(t, err)
		defer lis.Close()

		srv := NewGRPCServer("grpc server", func(s *grpc.Server) error {
			return nil
		}, WithBindAddress("localhost:9091"))

		defer srv.Close(ctx)

		err = srv.Listen(ctx)
		assert.Error(t, err)
	})

	t.Run("should fail when the registerer fails", func(t *testing.T) {
		ctx := context.Background()

		wantErr := errors.New("random error")

		srv := NewGRPCServer("grpc server", func(s *grpc.Server) error {
			return wantErr
		}, WithBindAddress("localhost:9091"))

		defer srv.Close(ctx)

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
