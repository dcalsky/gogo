package logs

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"os"
)

var defaultLogger zerolog.Logger

func init() {
	logLevel := os.Getenv("log-level")
	if logLevel == "" {
		logLevel = "info"
	}
	lv, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		lv = zerolog.InfoLevel
	}
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	stdWriter := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "20060102T15:04:05.000Z07"}
	defaultLogger = zerolog.New(stdWriter).Level(lv).Hook(LogidHook{}).With().Timestamp().Logger()
}

func Fatalf(ctx context.Context, template string, args ...any) {
	defaultLogger.Fatal().Ctx(ctx).Msg(fmt.Sprintf(template, args...))
}

func Errorf(ctx context.Context, template string, args ...any) {
	defaultLogger.Error().Ctx(ctx).Msg(fmt.Sprintf(template, args...))
}

func Warnf(ctx context.Context, template string, args ...any) {
	defaultLogger.Warn().Ctx(ctx).Msg(fmt.Sprintf(template, args...))
}

func Infof(ctx context.Context, template string, args ...any) {
	defaultLogger.Info().Ctx(ctx).Msg(fmt.Sprintf(template, args...))
}

func Debugf(ctx context.Context, template string, args ...any) {
	defaultLogger.Debug().Ctx(ctx).Msg(fmt.Sprintf(template, args...))
}
