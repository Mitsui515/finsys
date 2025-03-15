package middleware

import (
	"context"
	"time"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

func Logger() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		start := time.Now()
		path := string(c.Request.URI().Path())
		method := string(c.Request.Method())

		hlog.CtxInfof(ctx, "Started %s %s", method, path)
		c.Next(ctx)

		latency := time.Since(start)
		statusCode := c.Response.StatusCode()
		hlog.CtxInfof(ctx, "Completed %s %s %d %v", method, path, statusCode, latency)
	}
}
