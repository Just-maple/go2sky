package pool

import (
	"github.com/SkyAPM/go2sky"
)

func newTestReporter() (go2sky.Reporter, error) {
	return &testReporter{}, nil
}

type testReporter struct{}

func (lr *testReporter) Boot(service string, serviceInstance string) {}

func (lr *testReporter) Send(spans []go2sky.ReportedSpan) {}

func (lr *testReporter) Close() {}
