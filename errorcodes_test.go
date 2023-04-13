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

func TestUpdate(t *testing.T) {
	type ErrorCodes struct {
		ErrorCode2000 ErrorCode `default:"2000"`
		ErrorCode2001 ErrorCode `default:"2001"`
		ErrorCode2002 ErrorCode `default:"2002"`
		ErrorCode2003 ErrorCode `default:"2003"`
	}
	errs := new(ErrorCodes)
	ErrorCodesUpdate(errs)
	fmt.Printf("%+v", ErrorCodesMap())
}
