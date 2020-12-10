package pool

import (
	"context"
	"testing"

	"github.com/SkyAPM/go2sky"
)

func do(t *go2sky.Tracer) {
	sp, ctx, _ := t.CreateLocalSpan(context.Background())
	var h string
	sp2, _ := t.CreateExitSpan(ctx, "test", "peer", func(header string) error {
		h = header
		return nil
	})
	sp3, ctx, _ := t.CreateEntrySpan(ctx, "test", func() (string, error) {
		return h, nil
	})
	sp3.End()
	for i := 0; i < 5; i++ {
		var tmpSp go2sky.Span
		tmpSp, ctx, _ = t.CreateLocalSpan(ctx)
		defer tmpSp.End()
	}
	sp2.End()
	sp.End()
}

func newTracer() *go2sky.Tracer {
	rp, _ := newTestReporter()
	t, _ := go2sky.NewTracer("test", go2sky.WithReporter(rp))
	return t
}

func run(b *testing.B, p bool) {
	t := newTracer()
	b.ReportAllocs()
	b.ResetTimer()
	if p {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				do(t)
			}
		})
	} else {
		for i := 0; i < b.N; i++ {
			do(t)
		}
	}
}

func BenchmarkDisablePoolP(b *testing.B) {
	go2sky.SetPoolEnable(false)
	run(b, true)
}

//5s BenchmarkDisablePoolP-12    	 1488094	      4190 ns/op	    6118 B/op	      70 allocs/op

func BenchmarkPoolP(b *testing.B) {
	go2sky.SetPoolEnable(true)
	run(b, true)
}

//5s BenchmarkPoolP-12    	 1727493	      3634 ns/op	    4134 B/op	      62 allocs/op
