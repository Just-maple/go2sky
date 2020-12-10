package go2sky

import (
	"sync"
)

var rootSegmentSpanPool = &sync.Pool{
	New: func() interface{} {
		return new(rootSegmentSpan)
	},
}

var segmentSpanImplPool = &sync.Pool{
	New: func() interface{} {
		return new(segmentSpanImpl)
	},
}

type PoolSpan interface {
	PutPool()
}

var enablePool = true

func (rs *rootSegmentSpan) PutPool() {
	segmentSpanImplPool.Put(rs.segmentSpanImpl)
}

func (s *segmentSpanImpl) PutPool() {
	segmentSpanImplPool.Put(s)
}

func SetPoolEnable(status bool) {
	enablePool = status
}

func PutSpanPool(spans []ReportedSpan) {
	if !enablePool {
		return
	}
	for i := range spans {
		if poolSpan, ok := spans[i].(PoolSpan); ok {
			poolSpan.PutPool()
		}
	}
}
