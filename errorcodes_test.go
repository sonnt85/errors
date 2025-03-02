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
	fmt.Printf("%#v\n", ErrorCodesMap())
	fmt.Print(Json(GetStandardErrorCode(Errors.AccessDeniedError)))
}

func TestCodeStr(t *testing.T) {
	err := WithErrorCode(2000, doSomethingElse(), "failed to do something")
	tests := []struct {
		err     error
		minsize int
		prefix  string
		want    string
	}{
		{err, 5, "ERR", "ERR02000"},
		{err, 6, "", "002000"},
		{err, 3, "", "2000"},
		{nil, 5, "ERR", ""},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("err: %v, minsize: %d, prefix: %s", tt.err, tt.minsize, tt.prefix), func(t *testing.T) {
			if got := CodeStr(tt.err, tt.minsize, tt.prefix); got != tt.want {
				t.Errorf("CodeStr() = %v, want %v", got, tt.want)
			}
		})
	}
}
