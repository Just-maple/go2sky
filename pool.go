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
	if enablePool {
		rootSegmentSpanPool.Put(rs)
	}
}

func (s *segmentSpanImpl) PutPool() {
	if enablePool {
		segmentSpanImplPool.Put(s)
	}
}

func SetPoolEnable(status bool) {
	enablePool = status
}
