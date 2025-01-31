package logs

import (
	"github.com/dcalsky/gogo/common/logid"
	"github.com/rs/zerolog"
)

type LogidHook struct {
}

func (s LogidHook) Run(e *zerolog.Event, level zerolog.Level, msg string) {
	ctx := e.GetCtx()
	spanId := logid.GetLogId(ctx)
	e.Str("logid", spanId)
}
