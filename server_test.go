package srvgrpc

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var _ = Describe("GRPCServer", func() {
	Describe("NewGRPCServer", func() {
		It("should set name, registerer and opts", func() {
			wantName := "name"
			wantRegisterer := func(_ *grpc.Server) error { return nil }
			wantOpt1 := WithBindAddress("bind:9090")
			wantOpt2 := WithConnectionTimeout(time.Second)
			s := NewGRPCServer(wantName, wantRegisterer, wantOpt1, wantOpt2)
			Expect(s.name).To(Equal(wantName))
			Expect(fmt.Sprintf("%p", s.registerer)).To(Equal(fmt.Sprintf("%p", wantRegisterer)))
			Expect(s.opts).To(HaveLen(2))
			Expect(fmt.Sprintf("%p", s.opts[0])).To(Equal(fmt.Sprintf("%p", wantOpt1)))
			Expect(fmt.Sprintf("%p", s.opts[1])).To(Equal(fmt.Sprintf("%p", wantOpt2)))
		})
	})

	Describe("Listen", func() {
		It("should start the server", func() {
			ctx := context.TODO()

			srv := NewGRPCServer("grpc server", func(_ *grpc.Server) error {
				return nil
			}, WithBindAddress("localhost:9091"))

			Expect(srv).NotTo(BeNil())

			err := srv.IsReady(ctx)
			Expect(err).To(MatchError(ErrNotReady))

			err = srv.Listen(ctx)
			Expect(err).NotTo(HaveOccurred())

			DeferCleanup(func() {
				_ = srv.Close(ctx)
			})

			Eventually(func() error {
				return srv.IsReady(ctx)
			}, 2*time.Second, 100*time.Millisecond).Should(Succeed(), "server never became ready")

			conn, err := grpc.NewClient("localhost:9091", grpc.WithTransportCredentials(insecure.NewCredentials()))
			Expect(err).NotTo(HaveOccurred())
			DeferCleanup(func() {
				_ = conn.Close()
			})
		})

		It("should fail when the port is busy", func() {
			ctx := context.Background()

			lis, err := net.Listen("tcp", "localhost:9092")
			Expect(err).NotTo(HaveOccurred())
			DeferCleanup(func() {
				_ = lis.Close()
			})

			srv := NewGRPCServer("grpc server", func(_ *grpc.Server) error {
				return nil
			}, WithBindAddress("localhost:9092"))

			DeferCleanup(func() {
				_ = srv.Close(ctx)
			})

			err = srv.Listen(ctx)
			Expect(err).To(HaveOccurred())
		})

		It("should fail when the registerer fails", func() {
			ctx := context.Background()

			wantErr := errors.New("random error")

			srv := NewGRPCServer("grpc server", func(_ *grpc.Server) error {
				return wantErr
			}, WithBindAddress("localhost:9093"))

			DeferCleanup(func() {
				_ = srv.Close(ctx)
			})

			err := srv.Listen(ctx)
			Expect(err).To(MatchError(wantErr))
		})
	})

	Describe("Name", func() {
		It("should return the server name", func() {
			wantName := "name"
			s := &GRPCServer{name: wantName}
			Expect(s.Name()).To(Equal(wantName))
		})
	})
})
