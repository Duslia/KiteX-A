package zlog

import (
	"context"
	"math/rand"

	"github.com/cloudwego/kitex/pkg/endpoint"
	"github.com/cloudwego/kitex/pkg/klog"
)

func LogMiddleware(mwCtx context.Context) endpoint.Middleware { // middleware builder

	logger := mwCtx.Value(endpoint.CtxLoggerKey).(klog.FormatLogger) // get the logger

	return func(next endpoint.Endpoint) endpoint.Endpoint { // middleware
		return func(ctx context.Context, request, response interface{}) error {
			// 采样 1%
			if rand.Intn(100) < 1 {
				logger.Debugf("Request is %v", request)
			}
			err := next(ctx, request, response)

			return err
		}
	}
}
