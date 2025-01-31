package base

import (
	"bytes"
	"fmt"
	"runtime"
	"strconv"
)

type Exception struct {
	error             // original error
	StatusCode int    // error http status code
	Code       string // error code
	Message    string // standard error message to user
	MessageCn  string // Chinese error message
	Public     bool   // public to user
	RawError   string // original error text
}

func NewException(statusCode int, textCode, msg string, msgCn string) Exception {
	return Exception{
		error:      nil,
		StatusCode: statusCode,
		Code:       textCode,
		Message:    msg,
		MessageCn:  msgCn,
		Public:     false,
		RawError:   "",
	}
}

func (e Exception) Error() string {
	return fmt.Sprintf("StatusCode: %d, Code: %s, Message: %s, RawError: %s", e.StatusCode, e.Code, e.Message, e.RawError)
}

func (e Exception) WithRawError(err error) Exception {
	e.error = err
	e.RawError = ""
	if e.error != nil {
		e.RawError = e.error.Error()
	}
	return e
}

func (e Exception) WithMessage(msg string) Exception {
	e.Message = msg
	return e
}

func PanicIf(expression bool, exception Exception) {
	if expression {
		panic(exception)
	}
}

func PanicIfErr(err error, exception Exception) {
	if err != nil {
		panic(exception)
	}
}

func Catch(f func()) (e *Exception) {
	defer func() {
		if panicData := recover(); panicData != nil {
			if except, ok := panicData.(Exception); ok {
				e = &except
			} else {
				temp := InternalError.WithRawError(fmt.Errorf("%v", panicData))
				e = &temp
			}
			return
		}
	}()
	f()
	return
}

func CatchRetry(times int, f func()) (e *Exception) {
	for i := 0; i < times; i++ {
		e = Catch(f)
		if e == nil {
			break
		}
	}
	return
}

func GetStack() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, 2048))
	pc := make([]uintptr, 16)
	n := runtime.Callers(3, pc)
	if n == 0 {
		return make([]byte, 0)
	}
	pc = pc[:n]
	frames := runtime.CallersFrames(pc)

	for {
		frame, more := frames.Next()
		buf.WriteString(frame.Function)
		buf.WriteByte('\n')
		buf.WriteByte('\t')
		buf.WriteString(frame.File)
		buf.WriteByte(':')
		buf.WriteString(strconv.Itoa(frame.Line))
		buf.WriteByte('\n')
		if !more {
			break
		}
	}
	return buf.Bytes()
}
