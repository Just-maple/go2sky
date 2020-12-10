package go2sky

import (
	"sync"
)

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
	rs.segmentSpanImpl.PutPool()
}

func (s *segmentSpanImpl) PutPool() {
	s.defaultSpan.Refs = s.defaultSpan.Refs[:0]
	s.defaultSpan.Logs = s.defaultSpan.Logs[:0]
	s.defaultSpan.Tags = s.defaultSpan.Tags[:0]
	segmentSpanImplPool.Put(s)
}

func PoolState() bool {
	return enablePool
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
