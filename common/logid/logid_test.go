package logid_test

import (
	"context"
	"github.com/dcalsky/gogo/common/logid"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLogId(t *testing.T) {
	logid := "example"
	ctx := context.Background()
	ctx = logid.ContextWithLogId(ctx, logid)
	ctxLogId := ctx.Value(logid.LogIdKey)
	require.NotNil(t, ctxLogId)
	require.Equal(t, logid, ctxLogId.(string))

	logIdStr := logid.GetLogId(ctx)
	require.Equal(t, logid, logIdStr)
}
