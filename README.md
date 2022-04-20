# server-grpc

`server-grpc` implements a service that runs a GRPC server.

## Example

```go
package grpc

import (
	srvgrpc "github.com/jamillosantos/server-grpc"
	userpbv1beta1 "github.com/<whatever>/protorepo/build/go/rpc/users/v1beta1"
	"google.golang.org/grpc"
)

func New() *srvgrpc.GRPCServer {
	return srvgrpc.NewGRPCServer(func(s *grpc.Server) error {
		userpbv1beta1.RegisterUserServiceServer(s, NewUserServiceServer())
		return nil
	}, srvgrpc.GRPCServerInitOpts{}).
		WithName("GRPC Server").
		WithConfig(srvgrpc.GRPCServerConfig{
			BindAddress: ":9090",
		})
}
```