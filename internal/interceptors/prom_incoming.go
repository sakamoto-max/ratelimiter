package interceptors

import (
	"context"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	rateLimiterPb "github.com/sakamoto-max/rate_limiter_proto/shared/rate_limiter"
	"google.golang.org/grpc"
)

type PromIncoming struct {
	RlOpts            prometheus.Counter
	OwnerOpts         *prometheus.CounterVec
	OwnerResourceOpts *prometheus.CounterVec
	IpOpts            *prometheus.CounterVec
	IpResourceOpts    *prometheus.CounterVec
}

func NewPromIncoming(reg *prometheus.Registry) *PromIncoming {
	m := &PromIncoming{
		RlOpts: promauto.With(reg).NewCounter(prometheus.CounterOpts{
			Name: "rate_limiter_requests_total",
			Help: "Total requests to rate_limiter_server",
		}),
		OwnerOpts: promauto.With(reg).NewCounterVec(prometheus.CounterOpts{
			Name: "rate_limiter_requests_total_by_owner",
			Help: "Total requests to rate_limiter_server by owner",
		}, []string{"owner_name"}),
		OwnerResourceOpts: promauto.With(reg).NewCounterVec(prometheus.CounterOpts{
			Name: "rate_limiter_requests_total_by_owner_resource",
			Help: "Total requests to rate_limiter_server by owner with resource",
		}, []string{"owner_name", "resource_name"}),
		IpOpts: promauto.With(reg).NewCounterVec(prometheus.CounterOpts{
			Name: "rate_limiter_requests_total_by_owner_ip",
			Help: "Total requests to rate_limiter_server from owner by an ip",
		}, []string{"owner_name", "user_ip"}),
		IpResourceOpts: promauto.With(reg).NewCounterVec(prometheus.CounterOpts{
			Name: "rate_limiter_requests_total_by_owner_ip_resource",
			Help: "Total requests to rate_limiter_server from owner by an ip in a route",
		}, []string{"owner_name", "user_ip", "resource_name"}),
	}

	return m
}

func (m *PromIncoming) UnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {

	reqNew := req.(*rateLimiterPb.CheckRequest)
	ownerName := GetOwnerName(ctx)
	m.RlOpts.Inc()
	m.OwnerOpts.WithLabelValues(ownerName).Inc()
	m.OwnerResourceOpts.WithLabelValues(ownerName, reqNew.ResourceName).Inc()
	m.IpOpts.WithLabelValues(ownerName, reqNew.ClientIp).Inc()
	m.IpResourceOpts.WithLabelValues(ownerName, reqNew.ClientIp, reqNew.ResourceName).Inc()

	return handler(ctx, req)
}
