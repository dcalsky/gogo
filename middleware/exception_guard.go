package middleware

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/dcalsky/gogo/base"
	"github.com/dcalsky/gogo/common/logs"
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

func HertzExceptionGuard() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		defer func() {
			if panicData := recover(); panicData != nil {
				panicErr := panicDataToError(panicData)
				logs.Errorf(ctx, "[ExceptionGuard] panic data: %s. stack: \n%s", panicErr.Error(), base.GetStack())
				args := base.GetArgs(ctx, c)
				e, ok := panicData.(base.Exception)
				if !ok {
					e = base.InternalError.WithRawError(panicErr)
				}
				base.RespondError(args, c, e)
				return
			}
		}()
		c.Next(ctx)
	}
}
