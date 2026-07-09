package interceptors

import "github.com/prometheus/client_golang/prometheus"

type Interceptors struct {
	Auth            *Auth
	PromIncoming    *PromIncoming
	PromOutgoing    *PromOutgoing
	CanonicalLogger *CanonicalLogger
	RequestId       *RequestId
}

func NewInterceptors(reg *prometheus.Registry) *Interceptors {
	return &Interceptors{
		Auth:            &Auth{},
		PromIncoming:    NewPromIncoming(reg),
		PromOutgoing:    NewPromOutgoing(reg),
		CanonicalLogger: &CanonicalLogger{},
		RequestId:       &RequestId{},
	}
}
