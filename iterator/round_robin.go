package iterator

import (
	"sync/atomic"

	"github.com/sotnikov-s/go-load-balancer/proxy"
)

// NewRoundRobin accepts a number of proxies to be used in the RoundRobin iterator
// and returns the RoundRobin instance itself
func NewRoundRobin(proxies ...*proxy.Proxy) Iterator {
	bunch := make(commonProxiesBunch, 0, len(proxies))
	for _, p := range proxies {
		bunch = append(bunch, p)
	}
	return &RoundRobin{
		proxies: bunch,
		current: -1,
	}
}

// RoundRobin is the most straightforward iterator which redirects requests to its
// proxies consequentially and cyclically. It's also usually called "next in loop"
type RoundRobin struct {
	proxies commonProxiesBunch
	current int32
}

// Next returns the next in the loop proxy
func (r *RoundRobin) Next() (*proxy.Proxy, error) {
	next := atomic.AddInt32(&r.current, 1) % int32(len(r.proxies))
	atomic.StoreInt32(&r.current, next)
	return getAvailableProxy(r.proxies, int(next))
}
