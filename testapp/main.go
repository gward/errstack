// A small standalone command-line tool for manually testing errstack.

package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/gward/errstack"
)

func main() {
	err := layer1()
	fmt.Fprintf(os.Stderr, "%+v\n", err)

	esErr := err.(errstack.ErrorStack)
	fmt.Printf("raw stack trace:\n")
	for _, pc := range(esErr.StackTrace()) {
		fmt.Printf("  %x\n", pc)
	}

	fmt.Printf("CallersFrames():\n")
	frames := runtime.CallersFrames([]uintptr(esErr.StackTrace()))
	for {
		frame, more := frames.Next()
		fmt.Printf("  %s, %s:%d\n", frame.Function, frame.File, frame.Line)
		if !more {
			break
		}
	}

	fmt.Printf("WriteStack():\n")
	esErr.WriteStack(os.Stdout)
}

func layer1() error {
	err := layer2()
	if err != nil {
		return errstack.Wrap(err, "outermost error")
	}
	return nil
}

func layer2() error {
	return layer3()
}

func layer3() error {
	err := layer4()
	if err != nil {
		return errstack.Wrap(err, "middle error")
	}
	return nil
}

func layer4() error {
	return errstack.New("innermost error")
}
