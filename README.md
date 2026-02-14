# server-grpc

A gRPC server wrapper with built-in support for:
- Prometheus metrics via `go-grpc-middleware`
- OpenTelemetry tracing
- Graceful shutdown
- Custom interceptors
- Configurable connection timeouts

## Installation

```bash
go get github.com/jamillosantos/server-grpc
```

## Features

- **Automatic Instrumentation**: Built-in Prometheus metrics and OpenTelemetry tracing
- **Flexible Configuration**: Configure bind address, timeouts, and custom interceptors
- **Graceful Shutdown**: Properly handles server lifecycle with `Listen()`, `Close()`, and `IsReady()`
- **Interceptor Chaining**: Easily add custom unary and stream interceptors
- **Service Interface**: Implements a standard service interface with `Name()`, `Listen()`, `Close()`, and `IsReady()` methods

## Quick Start

```go
srv := srvgrpc.NewGRPCServer(
	"My gRPC Server",
	func(s *grpc.Server) error {
		// Register your gRPC services here
		// pb.RegisterYourServiceServer(s, &yourServiceImpl{})
		return nil
	},
)

ctx := context.Background()
if err := srv.Listen(ctx); err != nil {
	log.Fatal(err)
}
```

## Examples

Complete working examples are available in the [`examples/`](./examples) directory:

- **[Basic Usage](./examples/basic)** - Minimal setup with graceful shutdown
- **[Custom Configuration](./examples/custom-config)** - Custom bind address, timeouts, and interceptors
- **[Complete Service](./examples/complete)** - Full example with service implementation and readiness checks

To run an example:

```bash
go run examples/basic/main.go
go run examples/custom-config/main.go
go run examples/complete/main.go
```

## Configuration Options

| Option | Description | Default |
|--------|-------------|---------|
| `WithBindAddress(addr string)` | Sets the bind address for the server | `:50051` |
| `WithConnectionTimeout(duration time.Duration)` | Sets the connection timeout | `15s` |
| `WithUnaryInterceptor(interceptors ...grpc.UnaryServerInterceptor)` | Adds custom unary interceptors | None |
| `WithStreamInterceptor(interceptors ...grpc.StreamServerInterceptor)` | Adds custom stream interceptors | None |

## Built-in Observability

The server automatically includes:

- **Prometheus Metrics**: Request counts, latencies, and error rates via `go-grpc-middleware/providers/prometheus`
- **OpenTelemetry Tracing**: Distributed tracing support via `otelgrpc`

These are configured automatically and don't require additional setup.

## API

### `NewGRPCServer(name string, registerer Registerer, opts ...Option) *GRPCServer`
Creates a new gRPC server instance.

### `Listen(ctx context.Context) error`
Starts the gRPC server. This is non-blocking and returns immediately after starting.

### `Close(ctx context.Context) error`
Gracefully stops the gRPC server and waits for it to shut down.

### `IsReady(ctx context.Context) error`
Returns `nil` if the server is ready to accept requests, or `ErrNotReady` otherwise.

### `Name() string`
Returns the human-readable name of the server.

## License

[Add your license here]