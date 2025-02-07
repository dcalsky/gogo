package middleware

import (
	"context"
	"github.com/dcalsky/gogo/ghttp"
	"github.com/dcalsky/gogo/idgen"
	"github.com/dcalsky/gogo/logid"
	"github.com/dcalsky/gogo/logs"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
)

func HertzTrafficLogger(maxRequestLen, maxResponseLen uint) app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		ctx = logid.ContextWithLogId(ctx, idgen.New[string]())
		reqStr := ghttp.DumpHertzRequest(&c.Request)
		startAt := time.Now()
		if len(reqStr) > int(maxRequestLen) {
			reqStr = reqStr[:maxRequestLen]
		}
		logs.Infof(ctx, "[TrafficLogger.Request] %s, body:\n%s", startAt.Format(time.RFC3339), reqStr)
		c.Next(ctx)
		resp := c.GetResponse()
		respStr := resp.Body()
		if len(respStr) > int(maxResponseLen) {
			respStr = respStr[:maxResponseLen]
		}
		logs.Infof(ctx, "[TrafficLogger.Response] status code: %d, cost: %d ms, body:\n%s", resp.StatusCode(), time.Since(startAt).Milliseconds(), respStr)
	}
}
