package reporter

import (
	"github.com/SkyAPM/go2sky"
)

type wrap struct {
	go2sky.Reporter
}

func (w wrap) Send(spans []go2sky.ReportedSpan) {
	w.Reporter.Send(spans)
	go2sky.PutSpanPool(spans)
}

func WrapPoolReporter(reporter go2sky.Reporter) go2sky.Reporter {
	return wrap{Reporter: reporter}
}
