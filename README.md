# errstack: Go errors with chained stack traces

errstack is a Go package for creating errors with stack traces.
It is heavily inspired by github.com/pkg/errors,
as well as chained exceptions as seen in Java and Python 3.

There are two basic operations with errstack:
creating a new error, and wrapping an existing one to add a stack trace.

## Creating a new error

To create a new error, simply use the `New()` constructor:

```
import "github.com/gward/errstack"

err := errstack.New("something is broken")
```

This error object behaves very similarly to an error from the standard `errors` package,
but it remembers the call stack from where it was created.
This allows it to generate a stack trace on demand.

errstack also provides an `Errorf()` function,
very similar to `fmt.Errorf()` from the standard library:

```
err := errstack.Errorf("%s is broken", "something")
```

Both of these constructors return an instance of `ErrorStack`.

## Wrapping an existing error

If you are dealing with a library whose errors do not include stack traces,
you may want to add a stack trace to those errors.
Use the `Wrap()` function for this:

```
result, err := otherlib.DoSomething(...)
if err != nil {
	err := errstack.Wrap(err, "failed to do something")
	return err
}
```

Like the constructors, `Wrap()` returns an `ErrorStack` object.

By default, repeated application of `Wrap()` results in a chain of stack traces:

```
outermost error:
testapp.layer1()
	testapp/main.go:23
testapp.main()
	testapp/main.go:11
runtime.main()
	/usr/lib/go-1.17/src/runtime/proc.go:255

middle error:
testapp.layer3()
	testapp/main.go:47
testapp.layer2()
	testapp/main.go:43
testapp.layer1()
	testapp/main.go:21

innermost error:
testapp.layer4()
	testapp/main.go:53
testapp.layer3()
	testapp/main.go:45
```

This tells you that the error originated in `layer4()`, at line 53
(_or_ line 53 received an error without a stack trace, and wrapped it).
That error was passed up from layer4() to layer3(), which wrapped it at line 47.
Then it propagated up the stack to layer1(), which wrapped it again at line 23.

(You may be familiar with chained stack traces from exceptions in
Java or Python 3.)

## Accessing a stack trace

There are several ways to extract the stack trace from an `ErrorStack`.
You can format it with `%+v`:

```
if err != nil {
	fmt.Printf(os.Stderr, "got an error: %+v\n", err)
}
```

(The formatted string is a multiline string with no trailing newline.)

This is handy because it works with any `error` object --
there is no need for a type assertion or type switch.

If you already know the error object is an `ErrorStack`,
you can use `WriteStack()` or `FormatStack()`:

```
// In real life, you would do this in a less panic-prone way.
esErr := err.(errstack.ErrorStack)

// Write a stack trace to a Writer (with trailing newline).
esErr.WriteStack(os.Stderr)

// Get a slice of lines that can be written.
lines := esErr.FormatStack()
for _, line := range lines {
	fmt.Fprintln(os.Stderr, line)
}
```

Finally, you can use either `StackTrace()` or `StackChain()`
to get an object representing a single stack trace or
a chain of stack traces.

`StackTrace()` returns a `StackTrace` object, which represents a single stack trace.
For a multiply wrapped error that chains several stack traces,
this will be the _outermost_ stack trace, which is usually not what you want.
Use `StackChain()` to get the full chain of stack traces as a `StackChain` object.

See the API reference for more details on working with
`StackTrace` and `StackChain` objects.
