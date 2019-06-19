package trace

import (
	"fmt"
	"io"
)

// Tracer is ...
type Tracer interface {
	Trace(...interface{})
}

type tracer struct {
	out io.Writer
}

// New is ...
func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

func (t *tracer) Trace(a ...interface{}) {
	t.out.Write([]byte(fmt.Sprint(a...)))
	t.out.Write([]byte("\n"))
}

type nilTracer struct{}

func (t *nilTracer) Trace(a ...interface{}) {}

// Off is ..
func Off() Tracer {
	return &nilTracer{}
}
