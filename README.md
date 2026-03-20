# shareful-utils-unpanicked
Unpanicked is a tiny Go helper that runs arbitrary functions from a safe context so panics are recovered instead of crashing your process.

## Installation
```
go get github.com/SharefulNetworks/shareful-utils-unpanicked/unpanicked
```

## Usage
```go
package main

import (
	"log"
	"github.com/SharefulNetworks/shareful-utils-unpanicked/unpanicked"
)

func main() {
	// Protect the main goroutine by wrapping top-level logic.
	unpanicked.RunSafe(MyTopLevelAppFunc, func(rec any, stack []byte) {
		log.Printf("recovered: %v\n%s", rec, stack)
	})
}

func MyTopLevelAppFunc() {
	// Application code that might panic.
	panic("boom")
}
```

Additional wrappers:

- `RunSafeCtx(ctx, fn, hook)`: skips execution if `ctx.Done()` is triggered before running `fn`, otherwise behaves like `RunSafe`.
- `RunSafeWithWG(wg, fn, hook)`: wraps `fn`, recovers panics, and always calls `wg.Done()` on exit.

If you pass `nil` for `hook`, Unpanicked logs the recovered panic via `log.Printf` along with the stack trace.

### Wrapping goroutines
```go
package main

import (
	"log"
	"sync"

	"github.com/SharefulNetworks/shareful-utils-unpanicked/unpanicked"
)

func main() {
	var wg sync.WaitGroup

	wg.Add(1)
	go unpanicked.RunSafeWithWG(&wg, func() {
		// Worker logic that might panic.
		panic("worker boom")
	}, func(rec any, stack []byte) {
		log.Printf("worker recovered: %v\n%s", rec, stack)
	})

	wg.Wait()
}
```

### Production hardening
- Wrap entry-point logic in `main` to guard the main goroutine from fatal panics.
- Wrap goroutine boundaries (e.g., `go unpanicked.RunSafe(func(){ ... }, hook)`) so panics in spawned workers are contained and reported.

## Development
- Run tests: `go test ./...`

## Authors
Giles Thompson    
giles@shareful.net


## License
MIT License.
