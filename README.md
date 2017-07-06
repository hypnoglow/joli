# joli

[![GoDoc](https://godoc.org/github.com/hypnoglow/joli?status.svg)](https://godoc.org/github.com/hypnoglow/joli)

A simple job processor library.

It just provides an easy way to define an N-sized queue and 
a K-sized worker pool to concurrently run queued jobs.

## Purpose

The example usage can be to offload any `http.Handler` from jobs 
that are waiting for an available worker and blocking the code.

See `example/main.go` for details.

## License

[MIT](https://github.com/hypnoglow/joli/blob/master/LICENSE)