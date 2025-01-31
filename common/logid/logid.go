package logid

import (
	"context"
)

const (
	LogIdKey = "GOGO_LOG_ID"
)

func ContextWithLogId(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, LogIdKey, id)
}

func GetLogId(ctx context.Context) string {
	return getStringFromContext(ctx, LogIdKey)
}

func getStringFromContext(ctx context.Context, key string) string {
	if ctx == nil {
		return ""
	}

	v := ctx.Value(key)
	if v == nil {
		return ""
	}
	str, _ := v.(string)
	return str
}
