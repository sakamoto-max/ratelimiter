package interceptors

import (
	"context"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

var reqIdKey ctxKey = "X-Request-ID"
var timeStartKey ctxKey = "time_start"

type RequestId struct{}

func (r *RequestId) UnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {

	reqId, ok := ctx.Value(reqIdKey).(string)
	if !ok || reqId == "" {
		reqId = uuid.New().String()
		ctx = context.WithValue(ctx, reqIdKey, reqId)
	}

	timeStart := time.Now()

	ctx = context.WithValue(ctx, timeStartKey, &timeStart)

	return handler(ctx, req)
}

func GetReqId(ctx context.Context) string {
	val, ok := ctx.Value(reqIdKey).(string)
	if !ok {
		return ""
	}
	return val
}

func GetTimeStart(ctx context.Context) *time.Time {
	val, ok := ctx.Value(timeStartKey).(*time.Time)
	if !ok {
		return nil
	}

	return val
}
