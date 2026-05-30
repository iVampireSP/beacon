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
	"github.com/iVampireSP/beacon/registry"
	"golang.org/x/sync/errgroup"
)

// App is the service runtime. It runs a set of servers under one service
// identity, registers the instance with the registry (if configured), and
// blocks until an OS signal — or a server error — triggers a graceful
// shutdown. It mirrors Kratos's kratos.App lifecycle.
type App struct {
	opts     options
	ctx      context.Context
	cancel   context.CancelFunc
	mu       sync.Mutex
	instance *registry.ServiceInstance
}

// New creates an App. Unless overridden by options, it fills in a random
// instance ID and the default termination signals (SIGTERM, SIGQUIT, SIGINT).
func New(opts ...Option) *App {
	o := options{
		ctx:              context.Background(),
		sigs:             []os.Signal{syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT},
		registrarTimeout: 10 * time.Second,
		stopTimeout:      10 * time.Second,
	}
	if id, err := uuid.NewUUID(); err == nil {
		o.id = id.String()
	}
	for _, opt := range opts {
		opt(&o)
	}
	ctx, cancel := context.WithCancel(o.ctx)
	return &App{opts: o, ctx: ctx, cancel: cancel}
}

// ID returns the instance ID.
func (a *App) ID() string { return a.opts.id }

// Name returns the service name.
func (a *App) Name() string { return a.opts.name }

// Version returns the service version.
func (a *App) Version() string { return a.opts.version }

// Run starts every server, registers the instance, runs the start hooks, and
// then blocks until a shutdown signal or a server error. It returns nil on a
// clean shutdown.
func (a *App) Run() error {
	instance, err := a.buildInstance()
	if err != nil {
		return err
	}
	a.mu.Lock()
	a.instance = instance
	a.mu.Unlock()

	if err = runHooks(a.ctx, a.opts.beforeStart); err != nil {
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

	if a.opts.registrar != nil {
		rctx, cancel := context.WithTimeout(a.ctx, a.opts.registrarTimeout)
		defer cancel()
		if err = a.opts.registrar.Register(rctx, instance); err != nil {
			return err
		}
	}

	if err = runHooks(a.ctx, a.opts.afterStart); err != nil {
		return err
	}
	logger.Info("service started", "id", instance.ID, "name", instance.Name, "version", instance.Version)

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

	if err = eg.Wait(); err != nil && !errors.Is(err, context.Canceled) {
		return err
	}
	return runHooks(a.ctx, a.opts.afterStop)
}

// Stop gracefully shuts the App down: it runs the stop hooks, deregisters the
// instance, and cancels the App context, which signals all servers to stop.
func (a *App) Stop() error {
	logger.Info("service stopping", "id", a.opts.id, "name", a.opts.name)

	if err := runHooks(a.ctx, a.opts.beforeStop); err != nil {
		return err
	}

	a.mu.Lock()
	instance := a.instance
	a.mu.Unlock()
	if a.opts.registrar != nil && instance != nil {
		ctx, cancel := context.WithTimeout(a.ctx, a.opts.registrarTimeout)
		defer cancel()
		if err := a.opts.registrar.Deregister(ctx, instance); err != nil {
			return err
		}
	}

	if a.cancel != nil {
		a.cancel()
	}
	return nil
}

// buildInstance assembles the ServiceInstance from the configured identity and
// the endpoints — explicit ones if set, otherwise those reported by servers
// implementing Endpointer.
func (a *App) buildInstance() (*registry.ServiceInstance, error) {
	endpoints := make([]string, 0, len(a.opts.endpoints))
	for _, e := range a.opts.endpoints {
		endpoints = append(endpoints, e.String())
	}
	if len(endpoints) == 0 {
		for _, srv := range a.opts.servers {
			e, ok := srv.(Endpointer)
			if !ok {
				continue
			}
			u, err := e.Endpoint()
			if err != nil {
				return nil, err
			}
			endpoints = append(endpoints, u.String())
		}
	}
	return &registry.ServiceInstance{
		ID:        a.opts.id,
		Name:      a.opts.name,
		Version:   a.opts.version,
		Metadata:  a.opts.metadata,
		Endpoints: endpoints,
	}, nil
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
