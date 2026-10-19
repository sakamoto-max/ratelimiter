package interceptors

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	rateLimiterPb "github.com/sakamoto-max/rate_limiter_proto/shared/rate_limiter"
	"github.com/sakamoto-max/ratelimiter/internal/opentel"
	"google.golang.org/grpc"
)

type logData struct {
	ReqId           string `json:"req_Id"`
	Resource        string `json:"resource"`
	Allowed         bool   `json:"allowed"`
	RemainingTokens int64  `json:"remaining_tokens"`
	RetryAfter      string `json:"retry_after"`
	Latency         string `json:"latency"`
}

type errorLogData struct {
	ReqId    string `json:"req_Id"`
	Resource string `json:"resource"`
	Error    string `json:"error"`
}

type CanonicalLogger struct {
	timeTaken prometheus.Histogram
}

func NewCanonicalLogger(reg *prometheus.Registry) *CanonicalLogger {
	return &CanonicalLogger{
		timeTaken: promauto.With(reg).NewHistogram(prometheus.HistogramOpts{
			Name: "time_taken_for_canonical_logger_in_seconds",
			Help: "the total time taken for canonical logger",
			Buckets: []float64{
				0.0005, // 0.5 ms
				0.001,  // 1 ms
				0.0025, // 2.5 ms
				0.005,  // 5 ms
				0.01,   // 10 ms
				0.025,  // 25 ms
				0.05,   // 50 ms
				0.1,    // 100 ms
				0.5,    // 500 ms
				1,      // 1 s
				2,      // 2 s
			},
		},
		),
	}
}

func (l *CanonicalLogger) UnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {

	ctx, span := opentel.CreateNewSpan(ctx, "logger_interceptor")
	defer span.End()

	loggerTimeStart := time.Now()

	reqNew := req.(*rateLimiterPb.CheckRequest)

	resp, err := handler(ctx, req)
	if err != nil {
		logJson, _ := json.Marshal(errorLogData{
			ReqId:    GetReqId(ctx),
			Resource: reqNew.ResourceName,
			Error:    err.Error(),
		})

		log.Println(string(logJson))
		return nil, err
	}

	r := resp.(*rateLimiterPb.CheckResponse)

	timeStart := GetTimeStart(ctx)

	newLogInMicroseconds := logData{
		ReqId:           GetReqId(ctx),
		Resource:        reqNew.ResourceName,
		Allowed:         r.Allowed,
		RemainingTokens: r.TokensLeft,
		RetryAfter:      fmt.Sprintf("%v s", r.TryAfter),
		Latency:         fmt.Sprintf("%v micro seconds", time.Since(*timeStart).Microseconds()),
	}

	logJson, _ := json.Marshal(newLogInMicroseconds)

	log.Println(string(logJson))

	l.timeTaken.Observe(float64(time.Since(loggerTimeStart).Seconds()))

	return resp, nil
}
