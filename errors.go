package sterrors

import (
	"fmt"
	"runtime"
	"slices"
)

type errorWithFrames struct {
	error
	frames []runtime.Frame
}

func (e *errorWithFrames) Unwrap() error {
	return e.error
}

func (e *errorWithFrames) Frames() []runtime.Frame {
	return e.frames
}

func (e *errorWithFrames) Error() string {
	if e.error == nil {
		return "<nil>"
	}
	return e.error.Error()
}

func WithFrames(err error) error {
	return &errorWithFrames{err, newFrames(1)}
}

func WithFramesSkip(err error, addSkip uint16) error {
	return &errorWithFrames{err, newFrames(1 + addSkip)}
}

func newFrames(addSkip uint16) []runtime.Frame {
	var pcs [64]uintptr
	n := runtime.Callers(2+int(addSkip), pcs[:])

	fr := runtime.CallersFrames(pcs[:n])

	frames := make([]runtime.Frame, 0, n)
	for more := true; more; {
		var frame runtime.Frame
		frame, more = fr.Next()
		frames = append(frames, frame)
	}

	return frames
}

func HasFrames(err error) bool {
	if x, ok := err.(interface{ Frames() []runtime.Frame }); ok {
		if len(x.Frames()) != 0 {
			return true
		}
	}
	if x, ok := err.(interface{ Unwrap() error }); ok {
		if HasFrames(x.Unwrap()) {
			return true
		}
	}
	if x, ok := err.(interface{ Unwrap() []error }); ok {
		if slices.ContainsFunc(x.Unwrap(), HasFrames) {
			return true
		}
	}
	return false
}

func StackTraces(err error) []error {
	return collect([]error{}, err)
}

func collect(errs []error, err error) []error {
	if x, ok := err.(interface{ Frames() []runtime.Frame }); ok {
		if len(x.Frames()) != 0 {
			errs = append(errs, err)
		}
	}
	if x, ok := err.(interface{ Unwrap() error }); ok {
		errs = collect(errs, x.Unwrap())
	}
	if x, ok := err.(interface{ Unwrap() []error }); ok {
		for _, e := range x.Unwrap() {
			errs = collect(errs, e)
		}
	}
	return errs
}

func Format(errs []error) []map[string]any {
	stackTraces := make([]map[string]any, 0, len(errs))

	for _, err := range errs {
		x, ok := err.(interface{ Frames() []runtime.Frame })
		if !ok {
			continue
		}

		stackTrace := make([]string, len(x.Frames()))
		for i, v := range x.Frames() {
			stackTrace[i] = fmt.Sprintf("%s(%s:%d)", v.Function, v.File, v.Line)
		}

		stackTraces = append(stackTraces, map[string]any{
			"error":      err.Error(),
			"stackTrace": stackTrace,
		})
	}

	return stackTraces
}

type errorFrames struct {
	s      string
	frames []runtime.Frame
}

func (e *errorFrames) Error() string {
	return e.s
}

func (e *errorFrames) Frames() []runtime.Frame {
	return e.frames
}

func New(text string) error {
	return &errorFrames{text, newFrames(1)}
}

func NewFrames(msg string, addSkip uint16) error {
	return &errorFrames{msg, newFrames(1 + addSkip)}
}
