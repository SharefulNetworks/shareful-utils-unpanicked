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
	unpanicked.RunSafe(func() { panic("boom") }, func(rec any, stack []byte) {
		log.Printf("recovered: %v\n%s", rec, stack)
	})
}
```

Additional wrappers:

- `RunSafeCtx(ctx, fn, hook)`: skips execution if `ctx.Done()` is triggered before running `fn`, otherwise behaves like `RunSafe`.
- `RunSafeWithWG(wg, fn, hook)`: wraps `fn`, recovers panics, and always calls `wg.Done()` on exit.

If you pass `nil` for `hook`, Unpanicked logs the recovered panic via `log.Printf` along with the stack trace.

## Development
- Run tests: `go test ./...`
