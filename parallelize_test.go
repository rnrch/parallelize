package parallelize

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUntil(t *testing.T) {
	tests := []struct {
		pieces      int
		parallelism int
	}{
		{
			pieces:      1000,
			parallelism: 0,
		},
		{
			pieces:      1000,
			parallelism: 20,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("test case %d", i), func(t *testing.T) {
			seen := make([]int32, tt.pieces)
			ctx := context.Background()
			Until(ctx, tt.pieces, func(p int) {
				atomic.AddInt32(&seen[p], 1)
			}, WithParallelism(tt.parallelism))

			wantSeen := make([]int32, tt.pieces)
			for i := 0; i < tt.pieces; i++ {
				wantSeen[i] = 1
			}
			assert.Equal(t, wantSeen, seen)
		})
	}
}

func BenchmarkUntil(b *testing.B) {
	tests := []struct {
		pieces      int
		parallelism int
	}{
		{
			pieces:      1000000,
			parallelism: 0,
		},
		{
			pieces:      1000000,
			parallelism: 100,
		},
	}
	for i, tt := range tests {
		b.Run(fmt.Sprintf("test case %d", i), func(b *testing.B) {
			ctx := context.Background()
			seen := make([]bool, tt.pieces)
			b.ResetTimer()
			for c := 0; c < b.N; c++ {
				Until(ctx, tt.pieces, func(p int) {
					seen[p] = prime(p)
				}, WithParallelism(tt.parallelism))
			}
			b.StopTimer()
			want := []bool{false, false, true, true, false, true, false, true, false, false, false, true}
			assert.Equal(b, want, seen[:len(want)])
		})
	}
}

func prime(p int) bool {
	if p <= 1 {
		return false
	}
	for i := 2; i*i <= p; i++ {
		if p%i == 0 {
			return false
		}
	}
	return true
}
