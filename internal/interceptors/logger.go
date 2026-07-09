package interceptors

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	rateLimiterPb "github.com/sakamoto-max/rate_limiter_proto/shared/rate_limiter"
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
	
}

func (l *CanonicalLogger) UnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {

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

	return resp, nil
}

