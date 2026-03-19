# shareful-utils-unpanicked
Unpanicked is a tiny, lightweight Go library that exports a series of wrapper function variants that all take a third-party function and callback, as input, and immediately runs the said function from a safe context,any panics that would have ordinarily terminated the entire process are instead caught and reported via the specified callback.
