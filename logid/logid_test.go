package logid

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogId(t *testing.T) {
	id := "example"
	ctx := context.Background()
	ctx = ContextWithLogId(ctx, id)
	ctxLogId := ctx.Value(LogIdKey)
	require.NotNil(t, ctxLogId)
	require.Equal(t, id, ctxLogId.(string))

	logIdStr := GetLogId(ctx)
	require.Equal(t, id, logIdStr)
}
