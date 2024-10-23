package tracetest

import "testing"

func TestGoatCallTracer(t *testing.T) {
	testCallTracer("callTracer", "goat_call_tracer", t)
}
