package transport

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	"github.com/iVampireSP/beacon/logger"
	"golang.org/x/sync/errgroup"
)

// TransportApp is the service runtime: it runs a set of servers under one service
// identity and blocks until an OS signal — or a server error — triggers a
// graceful shutdown, with before/after start/stop hooks. It mirrors the run
// half of Kratos's kratos.TransportApp lifecycle.
//
// Service registration/discovery is deliberately NOT here: in production the
// service mesh (Istio) / Kubernetes handles discovery, LB and mTLS; locally a
// peer's address comes from config. So there is no registry to run.
type TransportApp struct {
	opts   options
	ctx    context.Context
	cancel context.CancelFunc
}

// NewTransportApp creates an App. Unless overridden by options, it fills in a random
// instance ID and the default termination signals (SIGTERM, SIGQUIT, SIGINT).
func NewTransportApp(opts ...Option) *TransportApp {
	o := options{
		ctx:         context.Background(),
		sigs:        []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		stopTimeout: 10 * time.Second,
	}
	if id, err := uuid.NewUUID(); err == nil {
		o.id = id.String()
	}
	for _, opt := range opts {
		opt(&o)
	}
	ctx, cancel := context.WithCancel(o.ctx)
	return &TransportApp{opts: o, ctx: ctx, cancel: cancel}
}

// ID returns the instance ID.
func (a *TransportApp) ID() string { return a.opts.id }

// Name returns the service name.
func (a *TransportApp) Name() string { return a.opts.name }

// Version returns the service version.
func (a *TransportApp) Version() string { return a.opts.version }

// Run starts every server, runs the start hooks, and blocks until a shutdown
// signal or a server error. It returns nil on a clean shutdown.
func (a *TransportApp) Run() error {
	if err := runHooks(a.ctx, a.opts.beforeStart); err != nil {
		return err
	}

	eg, ctx := errgroup.WithContext(a.ctx)
	var wg sync.WaitGroup
	for _, srv := range a.opts.servers {
		// Stop the server once the App context is cancelled.
		eg.Go(func() error {
			<-ctx.Done()
			stopCtx, cancel := context.WithTimeout(context.Background(), a.opts.stopTimeout)
			defer cancel()
			return srv.Stop(stopCtx)
		})
		wg.Add(1)
		eg.Go(func() error {
			wg.Done()
			return srv.Start(a.ctx)
		})
	}
	wg.Wait()

	if err := runHooks(a.ctx, a.opts.afterStart); err != nil {
		return err
	}
	logger.Info("service started", "id", a.opts.id, "name", a.opts.name, "version", a.opts.version)

	c := make(chan os.Signal, 1)
	signal.Notify(c, a.opts.sigs...)
	eg.Go(func() error {
		select {
		case <-ctx.Done():
			return nil
		case <-c:
			return a.Stop()
		}
	})

	if err := eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return runHooks(a.ctx, a.opts.afterStop)
}

// Stop gracefully shuts the App down: it runs the stop hooks, then cancels the
// App context, which signals all servers to stop.
func (a *TransportApp) Stop() error {
	logger.Info("service stopping", "id", a.opts.id, "name", a.opts.name)

	if err := runHooks(a.ctx, a.opts.beforeStop); err != nil {
		return err
	}
	if a.cancel != nil {
		a.cancel()
	}
	return nil
}

// runHooks runs lifecycle hooks in order, stopping at the first error.
func runHooks(ctx context.Context, hooks []func(context.Context) error) error {
	for _, fn := range hooks {
		if err := fn(ctx); err != nil {
			return err
		}
	}
	return nil
}
