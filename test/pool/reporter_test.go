package pool

import (
	"github.com/SkyAPM/go2sky"
	"github.com/SkyAPM/go2sky/reporter/grpc/common"
	agentv3 "github.com/SkyAPM/go2sky/reporter/grpc/language-agent"
)

func newTestReporter() (go2sky.Reporter, error) {
	return &testReporter{sendCh: make(chan *agentv3.SegmentObject, 10000)}, nil
}

type testReporter struct {
	sendCh chan *agentv3.SegmentObject
}

func (r *testReporter) Boot(service string, serviceInstance string) {
	go func() {
		for {
			s := <-r.sendCh
			_ = s
		}
	}()
}

func (r *testReporter) Send(spans []go2sky.ReportedSpan) {
	spanSize := len(spans)
	if spanSize < 1 {
		return
	}
	rootSpan := spans[spanSize-1]
	rootCtx := rootSpan.Context()
	segmentObject := &agentv3.SegmentObject{
		TraceId:        rootCtx.TraceID,
		TraceSegmentId: rootCtx.SegmentID,
		Spans:          make([]*agentv3.SpanObject, spanSize),
	}
	for i, s := range spans {
		var tags []*common.KeyStringValuePair
		var logs []*agentv3.Log

		if go2sky.PoolState() {
			copy(tags, s.Tags())
			copy(logs, s.Logs())
		} else {
			tags = s.Tags()
			logs = s.Logs()
		}

		spanCtx := s.Context()
		segmentObject.Spans[i] = &agentv3.SpanObject{
			SpanId:        spanCtx.SpanID,
			ParentSpanId:  spanCtx.ParentSpanID,
			StartTime:     s.StartTime(),
			EndTime:       s.EndTime(),
			OperationName: s.OperationName(),
			Peer:          s.Peer(),
			SpanType:      s.SpanType(),
			SpanLayer:     s.SpanLayer(),
			ComponentId:   s.ComponentID(),
			IsError:       s.IsError(),
			Tags:          tags,
			Logs:          logs,
		}
		srr := make([]*agentv3.SegmentReference, 0)
		if i == (spanSize-1) && spanCtx.ParentSpanID > -1 {
			srr = append(srr, &agentv3.SegmentReference{
				RefType:              agentv3.RefType_CrossThread,
				TraceId:              spanCtx.TraceID,
				ParentTraceSegmentId: spanCtx.ParentSegmentID,
				ParentSpanId:         spanCtx.ParentSpanID,
			})
		}
		if len(s.Refs()) > 0 {
			for _, tc := range s.Refs() {
				srr = append(srr, &agentv3.SegmentReference{
					RefType:                  agentv3.RefType_CrossProcess,
					TraceId:                  spanCtx.TraceID,
					ParentTraceSegmentId:     tc.ParentSegmentID,
					ParentSpanId:             tc.ParentSpanID,
					ParentService:            tc.ParentService,
					ParentServiceInstance:    tc.ParentServiceInstance,
					ParentEndpoint:           tc.ParentEndpoint,
					NetworkAddressUsedAtPeer: tc.AddressUsedAtClient,
				})
			}
		}
		segmentObject.Spans[i].Refs = srr
	}
	select {
	case r.sendCh <- segmentObject:
	default:
	}
}

func (r *testReporter) Close() {}
