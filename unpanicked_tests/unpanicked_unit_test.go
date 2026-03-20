package unpanickedtests_test

import (
	"context"
	"sync"
	"testing"
	"time"

	unpanicked "github.com/SharefulNetworks/shareful-utils-unpanicked/unpanicked"
)

func TestRunSafeCallsHookOnPanic(t *testing.T) {
	t.Helper()

	var (
		called    bool
		recovered any
		stack     []byte
	)

	hook := func(r any, s []byte) {
		called = true
		recovered = r
		stack = s
	}

	unpanicked.RunSafe(func() { panic("boom") }, hook)

	if !called {
		t.Fatal("expected hook to be called")
	}
	if recovered != "boom" {
		t.Fatalf("unexpected recover value: %v", recovered)
	}
	if len(stack) == 0 {
		t.Fatal("expected stack to be populated")
	}
}

func TestRunSafeCtxSkipsWhenContextDone(t *testing.T) {
	t.Helper()

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	ran := false
	hookCalled := false

	unpanicked.RunSafeCtx(ctx, func() { ran = true }, func(_ any, _ []byte) { hookCalled = true })

	if ran {
		t.Fatal("expected function not to run when context is done")
	}
	if hookCalled {
		t.Fatal("expected hook not to run when function is skipped")
	}
}

func TestRunSafeWithWGSignalsAndRecovers(t *testing.T) {
	t.Helper()

	var (
		wg         sync.WaitGroup
		hookCalled bool
		recovered  any
	)

	wg.Add(1)

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	go unpanicked.RunSafeWithWG(&wg, func() { panic("boom") }, func(r any, _ []byte) {
		hookCalled = true
		recovered = r
	})

	select {
	case <-done:
	case <-time.After(200 * time.Millisecond):
		t.Fatal("waitgroup did not complete in time")
	}

	if !hookCalled {
		t.Fatal("expected hook to be called")
	}
	if recovered != "boom" {
		t.Fatalf("unexpected recover value: %v", recovered)
	}
}
