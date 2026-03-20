package unpanicked


import (
    "context"
    "log"
    "runtime/debug"
)

// PanicHook defines the signature of the callback function that handles panics.
type PanicHook func(rec any, stack []byte)

// RunSafe runs fn, recovering panics and invoking hook (or log.Printf if hook is nil).
func RunSafe(fn func(), hook PanicHook) {
    defer func() {
        if r := recover(); r != nil {
            if hook != nil {
                hook(r, debug.Stack())
                return
            }
            log.Printf("panic recovered: %v\n%s", r, debug.Stack())
        }
    }()
    fn()
}

// RunSafeCtx checks ctx before running fn; recovers panics the same way.
func RunSafeCtx(ctx context.Context, fn func(), hook PanicHook) {
    if ctx != nil {
        select {
        case <-ctx.Done():
            return
        default:
        }
    }
    RunSafe(fn, hook)
}

// RunSafeWithWG wraps fn with recover and signals wg.Done() on exit.
func RunSafeWithWG(wg interface{ Done() }, fn func(), hook PanicHook) {
    defer func() {
        if wg != nil {
            wg.Done()
        }
    }()
    RunSafe(fn, hook)
}