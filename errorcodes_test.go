package errors

import (
	"fmt"
	"testing"
)

func TestErrorCodes(t *testing.T) {
	// fmt.Print(MapErrorCode())
	err := doSomething()
	if err != nil {
		et := err.(*withStack)
		// x := et.Error()
		et.stack.StackTrace()
		fmt.Printf("%+v", err)
	}
}

func doSomething() error {
	return WrapfErrorCode(doSomethingElse(), 2000, "failed to do something")
}

func doSomethingElse() error {
	return New("failed to do something else")

}
