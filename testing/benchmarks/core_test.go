//go:build testing

package benchmarks

import (
	"context"
	"testing"

	"github.com/zoobz-io/capitan"
	"github.com/zoobz-io/sum"
)

// Service locator benchmarks

type benchService struct{}

// benchEventData is a struct type for benchmarking Event with sentinel.
type benchEventData struct {
	Message string
}

// benchEventInt is a struct type for benchmarking Event with int-like data.
type benchEventInt struct {
	Value int
}

func BenchmarkRegister(b *testing.B) {
	sum.Reset()
	b.Cleanup(sum.Reset)

	k := sum.Start()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		sum.Register[benchService](k, benchService{})
		sum.Unregister[benchService]()
	}
}

func BenchmarkUse(b *testing.B) {
	sum.Reset()
	b.Cleanup(sum.Reset)

	k := sum.Start()
	sum.Register[benchService](k, benchService{})

	ctx := context.Background()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, _ = sum.Use[benchService](ctx)
	}
}

func BenchmarkMustUse(b *testing.B) {
	sum.Reset()
	b.Cleanup(sum.Reset)

	k := sum.Start()
	sum.Register[benchService](k, benchService{})

	ctx := context.Background()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = sum.MustUse[benchService](ctx)
	}
}

func BenchmarkUseParallel(b *testing.B) {
	sum.Reset()
	b.Cleanup(sum.Reset)

	k := sum.Start()
	sum.Register[benchService](k, benchService{})

	ctx := context.Background()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = sum.Use[benchService](ctx)
		}
	})
}

// Event benchmarks

func BenchmarkEventEmit(b *testing.B) {
	signal := capitan.NewSignal("bench.emit", "Benchmark emit")
	event := sum.NewInfoEvent[benchEventData](signal)

	ctx := context.Background()
	data := benchEventData{Message: "test"}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		event.Emit(ctx, data)
	}
}

func BenchmarkEventEmitWithListener(b *testing.B) {
	signal := capitan.NewSignal("bench.emit.listen", "Benchmark emit with listener")
	event := sum.NewInfoEvent[benchEventData](signal)

	listener := event.Listen(func(ctx context.Context, data benchEventData) {
		// no-op
	})
	defer listener.Close()

	ctx := context.Background()
	data := benchEventData{Message: "test"}
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		event.Emit(ctx, data)
	}
}

func BenchmarkEventListen(b *testing.B) {
	signal := capitan.NewSignal("bench.listen", "Benchmark listen")
	event := sum.NewInfoEvent[benchEventData](signal)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		listener := event.Listen(func(ctx context.Context, data benchEventData) {})
		listener.Close()
	}
}

func BenchmarkEventEmitParallel(b *testing.B) {
	signal := capitan.NewSignal("bench.emit.parallel", "Benchmark emit parallel")
	event := sum.NewInfoEvent[benchEventInt](signal)

	listener := event.Listen(func(ctx context.Context, data benchEventInt) {
		// no-op
	})
	defer listener.Close()

	ctx := context.Background()
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			event.Emit(ctx, benchEventInt{Value: i})
			i++
		}
	})
}

// Baseline results (run date: [add date when running])
// BenchmarkRegister
// BenchmarkUse
// BenchmarkMustUse
// BenchmarkUseParallel
// BenchmarkEventEmit
// BenchmarkEventEmitWithListener
// BenchmarkEventListen
// BenchmarkEventEmitParallel
