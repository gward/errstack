package errstack

import (
	"fmt"
	"io"
	"runtime"
)

// ErrorStack is an error object with a stack trace and layers of context.
type ErrorStack struct {
	msg string
	cause error
	stack StackTrace
}

// New creates a new ErrorStack with a single message and a new stack trace,
// pointing to the caller of New().
func New(msg string) ErrorStack {
	return ErrorStack{
		msg: msg,
		stack: Callers(2),
	}
}

// Errorf creates a new ErrorStack with a message generated from format and args.
func Errorf(format string, args ...interface{}) ErrorStack {
	msg := fmt.Sprintf(format, args...)
	return ErrorStack{
		msg: msg,
		stack: Callers(2),
	}
}

func (err ErrorStack) Error() string {
	return err.msg
}

// Cause returns the immediate predecessor in the chain of causation
// to this error, or nil if there is no known cause.
func (err ErrorStack) Cause() error {
	return err.cause
}

func (err ErrorStack) StackTrace() StackTrace {
	return err.stack
}

func (err ErrorStack) Format(state fmt.State, verb rune) {
	switch verb {
	case 'v':
		if state.Flag('+') {
			err.WriteStack(state)
			return
		}
		fallthrough
	case 's':
		io.WriteString(state, err.msg)
	case 'q':
		fmt.Fprintf(state, "%q", err.msg)
	}
}

// WriteStack writes the chain of stack traces in err to writer.
func (err ErrorStack) WriteStack(writer io.Writer) {
	stopFunction := "runtime.main"
	for _, line := range err.FormatStack(&stopFunction) {
		writer.Write([]byte(line + "\n"))
	}
}

// FormatStack returns the chain of stack traces as a slice of ready-to-print
// lines of text.
func (err ErrorStack) FormatStack(stopFunction *string) []string {
	firstFunction := runtime.FuncForPC(err.stack[0]).Name()
	lines := err.stack.FormatStack(err.msg, stopFunction)

	// chain stack traces to the next underlying error
	if esCause, ok := err.cause.(ErrorStack); ok {
		lines = append(lines, "")
		lines = append(lines, esCause.FormatStack(&firstFunction)...)
	}

	return lines
}

// WrapChain returns a new ErrorStack that wraps an existing error,
// adding a stack trace to it. If cause is already an ErrorStack,
// then the stack traces will be chained when the error is formatted.
func WrapChain(cause error, msg string) ErrorStack {
	return wrap(cause, msg)
}

// WrapOptional returns cause unchanged if cause is already an ErrorStack.
// Otherwise, it wraps cause in a new ErrorStack that adds a stack trace.
func WrapOptional(cause error, msg string) ErrorStack {
	// cause already has a stack trace: preserve it and ignore msg
	if esCause, ok := cause.(ErrorStack); ok {
		return esCause
	}
	return wrap(cause, msg)
}

// WrapTruncate returns cause with its stack trace truncated to the
// current caller, if cause is already an ErrorStack. Otherwise, it
// wraps cause in a new ErrorStack and adds a stack trace.
func WrapTruncate(cause error, msg string) ErrorStack {
	if esCause, ok := cause.(ErrorStack); ok {
		esCause.stack = Callers(2)
		return esCause
	}
	return wrap(cause, msg)
}

func wrap(cause error, msg string) ErrorStack {
	return ErrorStack{
		msg: msg,
		cause: cause,
		stack: Callers(3),
	}
}

// Wrap is an alias for WrapChain. If you want different default
// error-wrapping behaviour in your application, just replace Wrap
// with a different function (e.g. WrapOptional or WrapTruncate)
// during startup.
var Wrap = WrapChain

type Causer interface {
	Cause() error
}

// Cause returns the underlying cause of the error, if possible.
// An error value has a cause if it implements Causer.
//
// If the error does not implement Cause, the original error will
// be returned. If the error is nil, nil will be returned without further
// investigation.
func Cause(err error) error {
	var previous error
	for err != nil {
		cause, ok := err.(Causer)
		if !ok {
			break
		}
		previous = err
		err = cause.Cause()
	}
	if err != nil {
		return err
	}
	return previous
}
