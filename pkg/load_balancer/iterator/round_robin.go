package iterator

import (
	"fmt"
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
func (r *RoundRobin) Next() (*proxy.Proxy, error) {
	next := atomic.AddInt32(&r.current, 1) % int32(len(r.proxies))
	atomic.StoreInt32(&r.current, next)
	return r.getAvailableProxy(int(next))
}

// getAvailableProxy walks through the proxies and returns the first available one starting from
// the one at the marker index. If no available proxy was found, it returns an error
func (r *RoundRobin) getAvailableProxy(marker int) (*proxy.Proxy, error) {
	for i := 0; i < len(r.proxies); i++ {
		tryProxy := (marker + i) % len(r.proxies)
		p := r.proxies[tryProxy]
		if p.IsAvailable() {
			return p, nil
		}
	}
	return nil, fmt.Errorf("all proxies are unavailable")
}
