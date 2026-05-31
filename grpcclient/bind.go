package grpcclient

import (
	"github.com/iVampireSP/beacon/contracts"
	"google.golang.org/grpc"
)

// BindStub registers a peer's generated gRPC stub as a singleton: the container
// dials the named peer through the pool (lazy, OpenTelemetry-instrumented) and
// hands callers the stub interface, so a service depends on the generated client
// interface — not on the pool or a concrete type — and stays testable. One line
// per peer, in the AppServiceProvider:
//
//	grpcclient.BindStub(app, "user", userv1.NewUserServiceClient)
//
// wrap is the generated New<Svc>Client constructor; its return type T is the key
// the container injects into whatever service asks for it.
func BindStub[T any](app contracts.Application, name string, wrap func(grpc.ClientConnInterface) T) {
	app.Singleton(func(pool *GRPCClientPool) (T, error) {
		conn, err := pool.Conn(name)
		if err != nil {
			var zero T
			return zero, err
		}
		return wrap(conn), nil
	})
}
