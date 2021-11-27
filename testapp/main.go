// A small standalone command-line tool for manually testing errstack.

package main

import (
	"fmt"
	"os"

	"github.com/gward/errstack"
)

func main() {
	errstack.Wrap = errstack.WrapChain

	err := layer1()
	fmt.Fprintf(os.Stderr, "Fprintf(%%v):\n%v\n", err)
	fmt.Fprintf(os.Stderr, "Fprintf(%%s):\n%s\n", err)
	fmt.Fprintf(os.Stderr, "Fprintf(%%q):\n%q\n", err)
	fmt.Fprintf(os.Stderr, "Fprintf(%%+v):\n%+v\n", err)

	esErr := err.(errstack.ErrorStack)
	cause := esErr.Cause()
	fmt.Fprintf(os.Stderr, "immediate cause: %T %v\n", cause, cause)

	cause = errstack.Cause(err)
	fmt.Fprintf(os.Stderr, "root cause: %T %v\n", cause, cause)

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
