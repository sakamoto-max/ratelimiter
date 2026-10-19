package interceptors

import (
	"context"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sakamoto-max/ratelimiter/internal/opentel"
	"go.opentelemetry.io/otel/attribute"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var reqIdKey ctxKey = "X-Request-ID"
var timeStartKey ctxKey = "time_start"

type RequestId struct {
	totalTimeTaken prometheus.Histogram
}

func NewRequestId(reg *prometheus.Registry) *RequestId {
	return &RequestId{
		totalTimeTaken: promauto.With(reg).NewHistogram(prometheus.HistogramOpts{
			Name: "request_duration_in_seconds",
			Help: "Total time spent by a request",
			Buckets: []float64{
				0.005,
				0.01,
				0.015,
				0.025,
				0.05,
				0.075,
				0.1,
				0.25,
				0.5, // 500 ms
				1,   // 1 s
				2,   // 2 s
			},
		},
		),
	}
}

// all the requests first hits this interceptor
// it starts the timer for calculating the total time spent by a request
// it gets the request id from the metadata
// if the request id is not present in the metadata, it generates a new request id
func (r *RequestId) UnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {

	ctx, span := opentel.CreateNewSpan(ctx, "req_id_interceptor")

	defer span.End()

	timeStart := time.Now()

	var reqId string

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, ErrTokenIsMissing
	}

	reqIds := md.Get("req_id")
	if len(reqIds) == 0 {

		reqId = uuid.New().String()
		ctx = context.WithValue(ctx, reqIdKey, reqId)
		span.SetAttributes(
			attribute.String("requestid.givenby", "ratelimiter"),
			attribute.String("request.id", reqId),
		)
		
		} else {
			
			reqId = reqIds[0]
			span.SetAttributes(
			attribute.String("requestid.givenby", "gateway"),
			attribute.String("request.id", reqId),
		)
	}

	ctx = context.WithValue(ctx, timeStartKey, &timeStart)

	defer func() {
		totalTime := time.Since(timeStart)

		r.totalTimeTaken.Observe(float64(totalTime.Seconds()))

		if totalTime.Seconds() > 1 {
			// log.Printf("request took %v seconds", totalTime.Seconds())
			log.Printf("request with id : %v took %v seconds", reqId, totalTime.Seconds())
		}
	}()

	resp, err := handler(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
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
