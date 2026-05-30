package transport_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/iVampireSP/beacon/transport"
)

// fakeServer blocks in Start until Stop is called, recording both.
type fakeServer struct {
	started chan struct{}
	done    chan struct{}
	stopped atomic.Bool
}

func newFakeServer() *fakeServer {
	return &fakeServer{started: make(chan struct{}), done: make(chan struct{})}
}

func (s *fakeServer) Start(context.Context) error {
	close(s.started)
	<-s.done
	return nil
}

func (s *fakeServer) Stop(context.Context) error {
	s.stopped.Store(true)
	close(s.done)
	return nil
}

func TestAppRunStop(t *testing.T) {
	srv := newFakeServer()
	var beforeStart, afterStart, beforeStop, afterStop atomic.Bool

	app := transport.New(
		transport.ID("test-id"),
		transport.Name("test-svc"),
		transport.Version("v1"),
		transport.Servers(srv),
		transport.BeforeStart(func(context.Context) error { beforeStart.Store(true); return nil }),
		transport.AfterStart(func(context.Context) error { afterStart.Store(true); return nil }),
		transport.BeforeStop(func(context.Context) error { beforeStop.Store(true); return nil }),
		transport.AfterStop(func(context.Context) error { afterStop.Store(true); return nil }),
	)

	done := make(chan error, 1)
	go func() { done <- app.Run() }()

	select {
	case <-srv.started:
	case <-time.After(2 * time.Second):
		t.Fatal("server did not start")
	}

	if err := app.Stop(); err != nil {
		t.Fatalf("Stop: %v", err)
	}

	select {
	case err := <-done:
		if err != nil {
			t.Fatalf("Run: %v", err)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Run did not return after Stop")
	}

	if !srv.stopped.Load() {
		t.Error("server was not stopped")
	}
	for name, hook := range map[string]*atomic.Bool{
		"beforeStart": &beforeStart, "afterStart": &afterStart,
		"beforeStop": &beforeStop, "afterStop": &afterStop,
	} {
		if !hook.Load() {
			t.Errorf("%s hook did not run", name)
		}
	}
}
