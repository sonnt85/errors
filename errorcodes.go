package errors

import (
	"fmt"
	"io"
	"reflect"

	"github.com/sonnt85/godefault"
	"github.com/sonnt85/strcase"
)

type ErrorCode int

type ErrorCodes struct {
	NotFound            ErrorCode `default:"0"`
	AlreadyExist        ErrorCode `default:"1"`
	PermissionDenied    ErrorCode `default:"2"`
	InvalidArgument     ErrorCode `default:"3"`
	InternalServerError ErrorCode `default:"4"`
	ServiceUnavailable  ErrorCode `default:"5"`
	Unauthorized        ErrorCode `default:"6"`
	AuthFailure         ErrorCode `default:"7"`
	TooManyRequests     ErrorCode `default:"8"`

	DatabaseConnectionError   ErrorCode `default:"1000"`
	EmailSendingError         ErrorCode `default:"1001"`
	InvalidInputFormatError   ErrorCode `default:"1002"`
	ConfigNotFoundError       ErrorCode `default:"1003"`
	InvalidConfigError        ErrorCode `default:"1004"`
	ResourceNotFoundError     ErrorCode `default:"1005"`
	ResourceUnavailableError  ErrorCode `default:"1006"`
	AccessDeniedError         ErrorCode `default:"1007"`
	NetworkError              ErrorCode `default:"1008"`
	TimeoutError              ErrorCode `default:"1009"`
	ExternalDependencyError   ErrorCode `default:"1010"`
	InvalidCredentialsError   ErrorCode `default:"1011"`
	ExpiredTokenError         ErrorCode `default:"1012"`
	UnsupportedMediaTypeError ErrorCode `default:"1013"`
	BadRequestError           ErrorCode `default:"1014"`
	ForbiddenError            ErrorCode `default:"1015"`
	MethodNotAllowedError     ErrorCode `default:"1016"`
	NotAcceptableError        ErrorCode `default:"1017"`
	ConflictError             ErrorCode `default:"1018"`
	UnprocessableEntityError  ErrorCode `default:"1019"`
	NotImplemented            ErrorCode `default:"1020"`
	ServiceDiscoveryError     ErrorCode `default:"1021"`
	DNSLookupError            ErrorCode `default:"1022"`
	SSLCertificateError       ErrorCode `default:"1023"`
	ConnectionRefused         ErrorCode `default:"1024"`
	OptionsError              ErrorCode `default:"1025"`
	UnknownError              ErrorCode `default:"9999"`
}

var messages map[ErrorCode]string

var Errors *ErrorCodes
var UserErrors interface{}

// func camelToNormal(s string) string {
// 	var buf bytes.Buffer
// 	for i, r := range s {
// 		if unicode.IsUpper(r) && i > 0 {
// 			buf.WriteRune(' ')
// 		}
// 		buf.WriteRune(unicode.ToLower(r))
// 	}
// 	return buf.String()
// }

func init() {
	Errors = new(ErrorCodes)
	godefault.SetDefaults(Errors)
	errs := reflect.TypeOf(*Errors)
	messages = map[ErrorCode]string{}
	for i := 0; i < errs.NumField(); i++ {
		field := errs.Field(i)
		fieldName := field.Name
		fieldValue := reflect.ValueOf(*Errors).FieldByName(fieldName).Interface()
		messages[fieldValue.(ErrorCode)] = strcase.ToDelimited(fieldName, ' ')
	}
}

func ErrorCodesUpdate(errorCodeStruct interface{}) {
	godefault.SetDefaults(errorCodeStruct)
	val := reflect.ValueOf(errorCodeStruct)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	if val.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		fieldValue := val.Field(i).Interface()
		if ec, ok := fieldValue.(ErrorCode); ok {
			messages[ec] = strcase.ToDelimited(field.Name, ' ')
		}
	}
}

func ErrorCodesMap() map[ErrorCode]string {
	return messages
}

// WithErrorCode annotates err with a new message.
// If err is nil, WithErrorCode returns nil.
func WithErrorCode(code int, err error, message string) error {
	if err == nil {
		return nil
	}
	return &withErrorCode{
		cause: err,
		msg:   message,
		code:  code,
	}
}

// WithStandardErrorCode annotates err with a new message.
// If err is nil, WithErrorCode returns nil.
func WithStandardErrorCode(code ErrorCode, err error) error {
	if err == nil {
		return nil
	}
	return &withErrorCode{
		cause: err,
		msg:   messages[code],
		code:  int(code),
	}
}

// WithErrorCodef annotates err with the format specifier.
// If err is nil, WithErrorCodef returns nil.
func WithErrorCodef(code int, err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &withErrorCode{
		cause: err,
		msg:   fmt.Sprintf(format, args...),
		code:  code,
	}
}

// WithErrorCodef annotates err with the format specifier.
// If err is nil, WithErrorCodef returns nil.
// func WithStandardErrorCodef(code ErrorCode, err error, format string, args ...interface{}) error {
// 	if err == nil {
// 		return nil
// 	}
// 	return &withErrorCode{
// 		cause: err,
// 		msg:   fmt.Sprintf(format, args...),
// 		code:  int(code),
// 	}
// }

type withErrorCode struct {
	code  int
	cause error
	msg   string
}

func (w *withErrorCode) Error() string {
	return fmt.Sprintf("[%d]%s\n%s", w.code, w.msg, w.cause.Error())
}

func (w *withErrorCode) Cause() error { return w.cause }

// Unwrap provides compatibility for Go 1.13 error chains.
func (w *withErrorCode) Unwrap() error { return w.cause }

func (w *withErrorCode) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			fmt.Fprintf(s, "%+v\n", w.Cause())
			io.WriteString(s, w.msg)
			return
		}
		fallthrough
	case 's', 'q':
		io.WriteString(s, w.Error())
	}
}

func (w *withErrorCode) Code() int { return w.code }

func Code(err error) int {
	type code interface {
		Code() int
	}
	if err != nil {
		errcode, ok := err.(code)
		if ok {
			return errcode.Code()
		}
	}
	return -1
}

// Wrap returns an error annotating err with a stack trace
// at the point Wrap is called, and the supplied message.
// If err is nil, Wrap returns nil.
func WrapWithErrorCode(err error, code int, message string) error {
	if err == nil {
		return nil
	}
	err = &withErrorCode{
		cause: err,
		msg:   message,
		code:  code,
	}
	return &withStack{
		err,
		callers(),
	}
}

// Wrapf returns an error annotating err with a stack trace
// at the point Wrapf is called, and the format specifier.
// If err is nil, Wrapf returns nil.
func WrapfErrorCode(err error, code int, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	err = &withErrorCode{
		cause: err,
		msg:   fmt.Sprintf(format, args...),
		code:  code,
	}
	return &withStack{
		err,
		callers(),
	}
}

func StackTraceErr(errc error) StackTrace {
	type stackTracer interface {
		StackTrace() StackTrace
	}

	err, ok := Cause(errc).(stackTracer)
	if !ok {
		return nil
	}
	return err.StackTrace()
	// st := err.StackTrace()
	// return fmt.Sprintf("%+v", st[0:2]) // top two frames
}
