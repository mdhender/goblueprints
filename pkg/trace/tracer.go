/*
 * Go Programming Blueprints - 2nd ed - by Mat Ryer
 *
 * Copyright (c) 2021 Michael D Henderson
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package trace

import (
	"fmt"
	"io"
)

// Tracer describes an object capable of tracing events through code.
type Tracer interface {
	Trace(...interface{})
}

// New creates a new Tracer that will write the output to
// the specified io.Writer.
func New(w io.Writer) Tracer {
	return &tracer{out: w}
}

// Off creates a Tracer that will ignore calls to Trace.
func Off() Tracer {
	return &nilTracer{}
}

// tracer is a Tracer that writes to an io.Writer.
type tracer struct {
	out io.Writer
}

// Trace writes the arguments to this Tracer's io.Writer.
func (t *tracer) Trace(a ...interface{}) {
	_, _ = fmt.Fprint(t.out, a...)
	_, _ = fmt.Fprintln(t.out)
}

// nilTracer is a Tracer that ignores all calls to Trace.
type nilTracer struct{}

// Trace for a nilTracer does nothing.
func (t *nilTracer) Trace(a ...interface{}) {}
