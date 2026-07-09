package interceptors

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	rateLimiterPb "github.com/sakamoto-max/rate_limiter_proto/shared/rate_limiter"
	"google.golang.org/grpc"
)

type PromOutgoing struct {
	RlOpts            prometheus.Counter
	RlOptsWithStatus  *prometheus.CounterVec
	OwnerOpts         *prometheus.CounterVec
	OwnerResourceOpts *prometheus.CounterVec
	IpOpts            *prometheus.CounterVec
	IpResourceOpts    *prometheus.CounterVec
}

func NewPromOutgoing(reg *prometheus.Registry) *PromOutgoing {
	return &PromOutgoing{
		RlOpts: promauto.With(reg).NewCounter(prometheus.CounterOpts{
			Name: "rate_limiter_server_outgoing_requests_total",
			Help: "Total outgoing requests from rate_limiter_server",
		}),
		RlOptsWithStatus: promauto.With(reg).NewCounterVec(prometheus.CounterOpts{
			Name: "rate_limiter_server_outgoing_requests_total_with_status",
			Help: "Total outgoing requests from rate_limiter_server with status allowed or blocked",
		}, []string{"status"}),
		OwnerOpts: promauto.With(reg).NewCounterVec(prometheus.CounterOpts{
			Name: "rate_limiter_server_outgoing_requests_with_status_for_owner",
			Help: "Total outgoing requests from rate_limiter_server with status allowed or blocked",
		}, []string{"owner_name", "status"}),
		OwnerResourceOpts: promauto.With(reg).NewCounterVec(prometheus.CounterOpts{
			Name: "rate_limiter_server_outgoing_requests_with_status_for_owner_resource",
			Help: "Total outgoing requests from rate_limiter_server with status allowed or blocked",
		}, []string{"owner_name", "resource_name", "status"}),
		IpOpts: promauto.With(reg).NewCounterVec(prometheus.CounterOpts{
			Name: "rate_limiter_server_outgoing_requests_with_status_for_owner_ip",
			Help: "Total outgoing requests from rate_limiter_server with status allowed or blocked",
		}, []string{"owner_name", "user_ip", "status"}),
		IpResourceOpts: promauto.With(reg).NewCounterVec(prometheus.CounterOpts{
			Name: "rate_limiter_server_outgoing_requests_with_status_for_owner_ip_resource",
			Help: "Total outgoing requests from rate_limiter_server with status allowed or blocked",
		}, []string{"owner_name", "user_ip", "resource_name", "status"}),
	}
}

func (m *PromOutgoing) UnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {

	reqNew := req.(*rateLimiterPb.CheckRequest)

	resp, err := handler(ctx, req)
	if err != nil {
		return nil, err
	}

	r := resp.(*rateLimiterPb.CheckResponse)

	var status string

	if r.Allowed {
		status = "allowed"
	} else {
		status = "blocked"
	}

	m.RlOpts.Inc()
	m.RlOptsWithStatus.With(prometheus.Labels{"status": status}).Inc()
	m.OwnerOpts.With(prometheus.Labels{"owner_name": GetOwnerName(ctx), "status": status}).Inc()
	m.OwnerResourceOpts.With(prometheus.Labels{"owner_name": GetOwnerName(ctx), "resource_name": reqNew.ResourceName, "status": status}).Inc()
	m.IpOpts.With(prometheus.Labels{"owner_name": GetOwnerName(ctx), "user_ip": reqNew.ClientIp, "status": status}).Inc()
	m.IpResourceOpts.With(prometheus.Labels{"owner_name": GetOwnerName(ctx), "user_ip": reqNew.ClientIp, "resource_name": reqNew.ResourceName, "status": status}).Inc()

	return resp, nil
}
