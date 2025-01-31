package logs

import (
	"context"
	"github.com/dcalsky/gogo/common/idgen"
	"github.com/dcalsky/gogo/common/logid"
	"testing"
)

func TestLogging(t *testing.T) {
	ctx := context.Background()
	ctx = logid.ContextWithLogId(ctx, idgen.New[string]())
	Infof(ctx, "[logging test] name: %s, age: %d", "gogo", 18)
}
