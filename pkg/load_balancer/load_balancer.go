package loadbalancer

import (
	"net/http"
	"sync/atomic"

	"github.com/sotnikov-s/go-load-balancer/pkg/proxy"
)

// LoadBalancer handles requests distributing them between one or more proxies
type LoadBalancer interface {
	http.Handler
}

// NewRoundRobinLoadBalancer is the RoundRobinLoadBalancer constructor
func NewRoundRobinLoadBalancer(proxies ...*proxy.Proxy) LoadBalancer {
	return &RoundRobinLoadBalancer{
		proxies: proxies,
		current: -1,
	}
}

// RoundRobinLoadBalancer is the most straightforward load balancer which redirects requests
// to its proxies consequentially and cyclically. It's also usually called "next in loop"
type RoundRobinLoadBalancer struct {
	proxies []*proxy.Proxy
	current int32
}

// ServeHTTP uses the next proxy to serve the http request
func (lb *RoundRobinLoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lb.proxies[lb.next()].ServeHTTP(w, r)
}

// next sets the load balancer to the next proxy
func (lb *RoundRobinLoadBalancer) next() int32 {
	// TODO: use proxy.IsAvailable()
	next := atomic.AddInt32(&lb.current, 1) % int32(len(lb.proxies))
	atomic.StoreInt32(&lb.current, next)
	return next
}
