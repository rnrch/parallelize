package parallelize

import (
	"context"
	"math"
	"sync"
)

type options struct {
	parallelism int
}

type Options func(*options)

type DoFunc func(piece int)

func WithParallelism(parallelism int) func(*options) {
	return func(o *options) {
		if parallelism < 1 {
			parallelism = 1
		}
		o.parallelism = parallelism
	}
}

// chunkSizeFor returns a chunk size for the given number of items to use for
// parallel work. The size aims to produce good CPU utilization.
// returns max(1, min(sqrt(n), n/Parallelism))
func (o *options) chunkSizeFor(n int) int {
	s := int(math.Sqrt(float64(n)))

	if r := n/o.parallelism + 1; s > r {
		s = r
	} else if s < 1 {
		s = 1
	}
	return s
}

// Until is a wrapper around workqueue.ParallelizeUntil to use in scheduling algorithms.
func Until(ctx context.Context, pieces int, do DoFunc, opts ...Options) {
	o := options{parallelism: 16}
	for _, opt := range opts {
		opt(&o)
	}
	until(ctx, o.parallelism, pieces, o.chunkSizeFor(pieces), do)
}

func until(ctx context.Context, workers int, pieces int, chunkSize int, do DoFunc) {
	if pieces == 0 {
		return
	}
	if chunkSize < 1 {
		chunkSize = 1
	}
	chunks := (pieces + chunkSize - 1) / chunkSize
	toProcess := make(chan int, chunks)
	for i := 0; i < chunks; i++ {
		toProcess <- i
	}
	close(toProcess)
	var stop <-chan struct{}
	if ctx != nil {
		stop = ctx.Done()
	}
	if chunks < workers {
		workers = chunks
	}
	wg := sync.WaitGroup{}
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			for chunk := range toProcess {
				start := chunk * chunkSize
				end := start + chunkSize
				if end > pieces {
					end = pieces
				}
				for p := start; p < end; p++ {
					select {
					case <-stop:
						return
					default:
						do(p)
					}
				}
			}
		}()
	}
	wg.Wait()
}
