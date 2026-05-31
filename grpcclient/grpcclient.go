// Package grpcclient dials peer services over gRPC. There is NO service
// registry: a peer's address comes from config "services.<name>". In production
// that is a Kubernetes/mesh DNS name (Istio does discovery, load balancing and
// mTLS); locally it is whatever config points at (e.g. localhost:9001), so the
// service runs as a plain binary with no cluster. Connections are OpenTelemetry-
// instrumented, so the W3C traceparent propagates across the call and the trace
// continues into the callee.
package grpcclient

import (
	"errors"
	"fmt"
	"sync"

	"github.com/iVampireSP/beacon/config"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GRPCClientPool is a pool of lazily-dialed, instrumented gRPC client connections to
// peer services, keyed by logical service name.
type GRPCClientPool struct {
	mu    sync.Mutex
	conns map[string]*grpc.ClientConn
}

// NewGRPCClientPool creates an empty client pool.
func NewGRPCClientPool() *GRPCClientPool {
	return &GRPCClientPool{conns: make(map[string]*grpc.ClientConn)}
}

// Conn returns a cached client connection to the named service, dialing it on
// first use using the address from config "services.<name>". A generated stub is
// created over it, e.g. userv1.NewUserServiceClient(clients.Conn("user")).
// Transport credentials are plaintext: in production the mesh terminates mTLS;
// locally there is no TLS.
func (c *GRPCClientPool) Conn(name string) (*grpc.ClientConn, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if conn, ok := c.conns[name]; ok {
		return conn, nil
	}

	addr := config.String("services." + name)
	if addr == "" {
		return nil, fmt.Errorf("grpcclient: no address for service %q — set config services.%s", name, name)
	}

	conn, err := grpc.NewClient(addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
	)
	if err != nil {
		return nil, fmt.Errorf("grpcclient: dial %q (%s): %w", name, addr, err)
	}
	c.conns[name] = conn
	return conn, nil
}

// Close closes every pooled connection.
func (c *GRPCClientPool) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var errs []error
	for name, conn := range c.conns {
		if err := conn.Close(); err != nil {
			errs = append(errs, fmt.Errorf("%s: %w", name, err))
		}
		delete(c.conns, name)
	}
	return errors.Join(errs...)
}
