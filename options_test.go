package srvgrpc

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"google.golang.org/grpc"
)

var _ = Describe("Options", func() {
	Describe("WithUnaryInterceptor", func() {
		It("should append the interceptor", func() {
			wantInterceptor := func(_ context.Context, _ interface{}, _ *grpc.UnaryServerInfo, _ grpc.UnaryHandler) (resp interface{}, err error) {
				return nil, nil
			}
			o := grpcOpts{}
			WithUnaryInterceptor(wantInterceptor)(&o)
			Expect(o.unaryInterceptors).To(HaveLen(1))
			Expect(fmt.Sprintf("%p", o.unaryInterceptors[0])).To(Equal(fmt.Sprintf("%p", wantInterceptor)))
		})
	})

	Describe("WithStreamInterceptor", func() {
		It("should append the interceptor", func() {
			wantInterceptor := func(_ interface{}, _ grpc.ServerStream, _ *grpc.StreamServerInfo, _ grpc.StreamHandler) error {
				return nil
			}
			o := grpcOpts{}
			WithStreamInterceptor(wantInterceptor)(&o)
			Expect(o.streamInterceptors).To(HaveLen(1))
			Expect(fmt.Sprintf("%p", o.streamInterceptors[0])).To(Equal(fmt.Sprintf("%p", wantInterceptor)))
		})
	})

	Describe("WithConnectionTimeout", func() {
		It("should set the connection timeout", func() {
			wantConnectionTimeout := time.Second * 123
			o := grpcOpts{}
			WithConnectionTimeout(wantConnectionTimeout)(&o)
			Expect(o.connectionTimeout).To(Equal(wantConnectionTimeout))
		})
	})
})
