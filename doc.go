// Package errstack provides an error type with stack traces.
//
// Use New() or Errorf() to create a new error with a stack trace.
// Use WrapChain(), WrapOptional(), or WrapTruncate() to add a stack
// trace to an existing error object.
//
// See https://github.com/gward/errstack for more information.
package errstack
