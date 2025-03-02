package errors

import (
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/sonnt85/godefault"
	"github.com/sonnt85/strcase"
)

type ErrorCode int

type ErrorCodes struct {
	Success                  ErrorCode `default:"0"`
	NotFound                 ErrorCode `default:"1"`
	AlreadyExist             ErrorCode `default:"2"`
	PermissionDenied         ErrorCode `default:"3"`
	InvalidArgument          ErrorCode `default:"4"`
	InternalServerError      ErrorCode `default:"5"`
	ServiceUnavailable       ErrorCode `default:"6"`
	Unauthorized             ErrorCode `default:"7"`
	AuthFailure              ErrorCode `default:"8"`
	TooManyRequests          ErrorCode `default:"9"`
	UnsuccessfulInstallation ErrorCode `default:"10"`
	UnspecifiedParam         ErrorCode `default:"11"`
	InvalidValue             ErrorCode `default:"12"`
	CanNotStartProg          ErrorCode `default:"13"`
	AlreadyRunning           ErrorCode `default:"14"`
	ConnectFailure           ErrorCode `default:"15"`
	ExitSuccessfully         ErrorCode `default:"16"`

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
	UnspecFail                ErrorCode `defaut:"1027"`
	NoAddrsAvail              ErrorCode `defaut:"1028"`
	NoBinding                 ErrorCode `defaut:"1029"`
	NotOnLink                 ErrorCode `defaut:"1030"`
	UseMulticast              ErrorCode `defaut:"1031"`
	NoPrefixAvail             ErrorCode `defaut:"1032"`
	// RFC 5007
	UnknownQueryType ErrorCode `defaut:"1033"`
	MalformedQuery   ErrorCode `defaut:"1034"`
	NotConfigured    ErrorCode `defaut:"1035"`
	NotAllowed       ErrorCode `defaut:"1036"`
	// RFC 5460
	QueryTerminated ErrorCode `defaut:"1037"`
	// RFC 7653
	DataMissing          ErrorCode `defaut:"1038"`
	CatchUpComplete      ErrorCode `defaut:"1039"`
	NotSupported         ErrorCode `defaut:"1040"`
	TLSConnectionRefused ErrorCode `defaut:"1041"`
	// RFC 8156
	AddressInUse               ErrorCode `defaut:"1042"`
	ConfigurationConflict      ErrorCode `defaut:"1043"`
	MissingBindingInformation  ErrorCode `defaut:"1044"`
	OutdatedBindingInformation ErrorCode `defaut:"1045"`
	ServerShuttingDown         ErrorCode `defaut:"1046"`
	DNSUpdateNotSupported      ErrorCode `defaut:"1047"`
	ExcessiveTimeSkew          ErrorCode `defaut:"1048"`
	UnknownError               ErrorCode `default:"9999"`
}

var messages map[ErrorCode]string

// Storage all Error codes
var Errors *ErrorCodes
var UserErrors interface{}

//	func camelToNormal(s string) string {
//		var buf bytes.Buffer
//		for i, r := range s {
//			if unicode.IsUpper(r) && i > 0 {
//				buf.WriteRune(' ')
//			}
//			buf.WriteRune(unicode.ToLower(r))
//		}
//		return buf.String()
//	}
func init() {
	Init()
}

func Init() {
	Errors = new(ErrorCodes)
	godefault.SetDefaults(Errors)
	errs := reflect.TypeOf(*Errors)
	messages = map[ErrorCode]string{}
	for i := 0; i < errs.NumField(); i++ {
		field := errs.Field(i)
		fieldName := field.Name
		fieldValue := reflect.ValueOf(*Errors).FieldByName(fieldName).Interface()
		errMsg := field.Tag.Get("errmsg")
		if errMsg == "" {
			errMsg = strcase.ToDelimited(fieldName, ' ')
		}
		messages[fieldValue.(ErrorCode)] = errMsg
	}
}

// Add Errorcode from errorCodeStruct [may be struct pointer or not] to Global Errorcode [Messages]
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

// Return all registered ErrorCode
func ErrorCodesMap() map[ErrorCode]string {
	return messages
}

// WithErrorCode annotates err with a new message.
// If errCause is nil, WithErrorCode returns nil.
func WithErrorCode(code int, errCause error, message string) error {
	if errCause == nil {
		return nil
	}
	return &withErrorCode{
		cause: errCause,
		msg:   message,
		code:  code,
	}
}

func GetStandardErrorCode(code ErrorCode) error {
	return &withErrorCode{
		msg:  messages[code],
		code: int(code),
	}
}

// WithStandardErrorCode annotates err with a new message.
// If errCause is nil, WithErrorCode returns nil.
func WithStandardErrorCode(code ErrorCode, errCause error) error {
	if errCause == nil {
		return nil
	}
	return &withErrorCode{
		cause: errCause,
		msg:   messages[code],
		code:  int(code),
	}
}

// WithStandardErrorCode annotates err(convert from cause) with a new message.
// If causeString is "", WithErrorCode returns nil.
func WithStandardErrorCodeCauseString(code ErrorCode, causeString string) error {
	if causeString == "" {
		return nil
	}
	return &withErrorCode{
		cause: fmt.Errorf(causeString),
		msg:   messages[code],
		code:  int(code),
	}
}

// If causeString is "", WithErrorCode returns nil.
func WithStandardSucces(causeString string) error {
	if causeString == "" {
		return nil
	}
	return &withErrorCode{
		cause: fmt.Errorf(causeString),
		msg:   messages[Errors.Success],
		code:  int(Errors.Success),
	}
}

// WithStandardErrorCode annotates err(convert from format cause) with a new message.
// If err is nil, WithErrorCode returns nil.
func WithStandardErrorfCodeCause(code ErrorCode, format string, args ...string) error {
	if format == "" {
		return nil
	}
	return &withErrorCode{
		cause: fmt.Errorf(format, args),
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

// WithErrorCodef annotates err (convert from format args) with the format specifier.
// If format == “” WithErrorCodef returns nil.
func WithErrorCodefCause(code int, format string, args ...interface{}) error {
	if format == "" {
		return nil
	}
	return &withErrorCode{
		cause: fmt.Errorf(format, args...),
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

func isJSON(s string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(s), &js) == nil
}

func (w *withErrorCode) Json(msg_is_json ...bool) string {
	cause := ""
	if w.cause != nil {
		cause = strings.ReplaceAll(w.cause.Error(), `"`, `\"`)
	}

	var msgStr string
	if len(msg_is_json) != 0 && msg_is_json[0] && isJSON(w.msg) {
		msgStr = w.msg
	} else {
		msgStr = fmt.Sprintf(`"%s"`, strings.ReplaceAll(w.msg, `"`, `\"`))
	}

	return fmt.Sprintf(`{"code": %d, "msg" : %s, "cause" : "%s"}`, w.code, msgStr, cause)
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

func Json(err error, msg_is_json ...bool) string {
	type code interface {
		Json(msg_is_json ...bool) string
	}
	if err != nil {
		errcode, ok := err.(code)
		if ok {
			return errcode.Json(msg_is_json...)
		}
	}
	return ""
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
