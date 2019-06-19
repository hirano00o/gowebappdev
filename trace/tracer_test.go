package trace

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Error("want not nil, have nil")
	} else {
		msg := "hello, test world"
		tracer.Trace(msg)
		if buf.String() != msg+"\n" {
			t.Errorf("want %s\\n, have %s", msg, buf.String())
		}
	}
}

func TestOff(t *testing.T) {
	var silentTracer Tracer = Off()
	silentTracer.Trace("test")
}
