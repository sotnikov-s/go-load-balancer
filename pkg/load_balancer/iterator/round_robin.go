package iterator

import (
	"sync/atomic"

	"github.com/sotnikov-s/go-load-balancer/pkg/proxy"
)

// NewRoundRobin accepts a number of proxies to be used in the RoundRobin iterator
// and returns the RoundRobin instance itself
func NewRoundRobin(proxies ...*proxy.Proxy) Iterator {
	return &RoundRobin{
		proxies: proxies,
		current: -1,
	}
}

// RoundRobin is the most straightforward iterator which redirects requests to its
// proxies consequentially and cyclically. It's also usually called "next in loop"
type RoundRobin struct {
	proxies []*proxy.Proxy
	current int32
}

// Next returns the next in the loop proxy
func (lb *RoundRobin) Next() *proxy.Proxy {
	// TODO: use proxy.IsAvailable()
	next := atomic.AddInt32(&lb.current, 1) % int32(len(lb.proxies))
	atomic.StoreInt32(&lb.current, next)
	return lb.proxies[next]
}
