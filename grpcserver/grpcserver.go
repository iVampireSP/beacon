// Package grpcserver is the inbound counterpart of grpcclient: an
// OpenTelemetry-instrumented gRPC server that implements transport.TransportServer, so
// the HTTP kernel runs it. The otelgrpc server handler extracts the W3C
// traceparent from incoming metadata, continuing the caller's trace; the address
// comes from config (grpc.host:grpc.port). An app registers its services onto
// Registrar() before serving.
package grpcserver

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/iVampireSP/beacon/config"
	"github.com/iVampireSP/beacon/logger"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

// GRPCServer wraps *grpc.Server with listen + graceful shutdown. It satisfies
// transport.TransportServer.
type GRPCServer struct {
	srv  *grpc.Server
	addr string
}

// NewGRPCServer creates an instrumented gRPC server listening on the address from config
// (grpc.host:grpc.port).
func NewGRPCServer(opts ...grpc.ServerOption) *GRPCServer {
	addr := fmt.Sprintf("%s:%d",
		config.String("grpc.host", "0.0.0.0"),
		config.Int("grpc.port", 9000),
	)
	return NewGRPCServerAt(addr, opts...)
}

// NewGRPCServerAt creates an instrumented gRPC server bound to an explicit address.
func NewGRPCServerAt(addr string, opts ...grpc.ServerOption) *GRPCServer {
	opts = append(opts, grpc.StatsHandler(otelgrpc.NewServerHandler()))
	return &GRPCServer{srv: grpc.NewServer(opts...), addr: addr}
}

// Registrar exposes the underlying server so services register themselves, e.g.
// helloworldv1.RegisterHelloWorldServiceServer(srv.Registrar(), impl).
func (s *GRPCServer) Registrar() grpc.ServiceRegistrar { return s.srv }

// Start listens and serves, blocking until Stop is called. Satisfies transport.Server.
func (s *GRPCServer) Start(_ context.Context) error {
	lis, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	logger.Info("gRPC server listening", "addr", s.addr)
	if err := s.srv.Serve(lis); err != nil && !errors.Is(err, grpc.ErrServerStopped) {
		return err
	}
	return nil
}

// Stop drains the server gracefully, force-stopping if ctx's deadline is hit.
// Satisfies transport.TransportServer.
func (s *GRPCServer) Stop(ctx context.Context) error {
	logger.Info("gRPC server shutting down")
	stopped := make(chan struct{})
	go func() {
		s.srv.GracefulStop()
		close(stopped)
	}()
	select {
	case <-stopped:
	case <-ctx.Done():
		s.srv.Stop()
	}
	return nil
}
