package glang

import (
	"bytes"
	"errors"
	"fmt"
	"runtime"
	"strconv"
)

func panicDataToError(data any) error {
	switch v := data.(type) {
	case error:
		return v
	case string:
		return errors.New(v)
	case fmt.Stringer:
		return errors.New(v.String())
	default:
		return fmt.Errorf("%v", v)
	}
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
