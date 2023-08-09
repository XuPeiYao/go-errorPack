package errorPack

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

var StackTraceLength = 32

type IStringable interface {
	String() string
}

type ErrorPack[T any] struct {
	InnerError error
	Data       *T
	stackTrace *string
}

func NewErrorPack[T any](innerErr error, data T) *ErrorPack[T] {
	return &ErrorPack[T]{
		InnerError: innerErr,
		Data:       &data,
	}
}

func (ep *ErrorPack[T]) MemberwiseClone() *ErrorPack[T] {
	return memberwiseClone(ep).(*ErrorPack[T])
}

func (ep *ErrorPack[T]) StackTrace() string {
	if ep.stackTrace == nil || len(*ep.stackTrace) == 0 {
		return ""
	}
	return *ep.stackTrace
}

func (ep *ErrorPack[T]) Throw() {
	// ref. https://github.com/uber-go/zap/blob/416e66ad83ebde35df0e09f02e65ac149e193b0e/stacktrace.go#L46-L110
	epi := ep.MemberwiseClone()

	var pcs []uintptr = make([]uintptr, StackTraceLength)
	pcs_count := runtime.Callers(2, pcs)

	frames := runtime.CallersFrames(pcs[:pcs_count])

	var sb strings.Builder
	i := 0
	for frame, more := frames.Next(); more; frame, more = frames.Next() {
		if i > 0 {
			sb.WriteString("\n")
		}

		sb.WriteString(frame.Function)
		sb.WriteString("\n")
		sb.WriteString("\t")
		sb.WriteString(frame.File)
		sb.WriteString(":")
		sb.WriteString(fmt.Sprintf("%d", frame.Line))

		i++
	}

	stackTrace := sb.String()
	epi.stackTrace = &stackTrace

	panic(epi)
}

func (ep *ErrorPack[T]) Error() string {
	dataStr := ""

	dataType := reflect.TypeOf(ep.Data)

	if dataType.Implements(reflect.TypeOf((*IStringable)(nil)).Elem()) {
		dataStr = "\r\n" + any(ep.Data).(IStringable).String()
	}

	stackTrace := ep.StackTrace()
	if len(stackTrace) > 0 {
		stackTrace = "\r\n" + stackTrace
	}
	if ep.InnerError == nil {
		return fmt.Sprintf("%s[%s]%s%s", "ErrorPack", getType(ep.Data), dataStr, stackTrace)
	} else {
		return fmt.Sprintf("%s[%s]\r\n%s%s%s", "ErrorPack", getType(ep.Data), dataStr, ep.InnerError.Error(), stackTrace)
	}
}
