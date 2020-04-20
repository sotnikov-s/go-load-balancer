package iterator

import (
	"fmt"
	"sort"

	"github.com/sotnikov-s/go-load-balancer/proxy"
)

// NewLeastConnections accepts a number of proxies to be used in the LeastConnections iterator
// and returns the LeastConnections instance itself
func NewLeastConnections(proxies ...*proxy.Proxy) Iterator {
	bunch := make(commonProxiesBunch, 0, len(proxies))
	for _, p := range proxies {
		bunch = append(bunch, p)
	}
	return &LeastConnections{
		proxies: bunch,
	}
}

// LeastConnections is an iterator that returns as next the least loaded available proxy
type LeastConnections struct {
	proxies commonProxiesBunch
}

// Next returns the least loaded available proxy
func (r *LeastConnections) Next() (*proxy.Proxy, error) {
	if r.proxies.Len() != 0 {
		proxies := append(r.proxies[:0:0], r.proxies...)
		sort.SliceStable(proxies, func(i, j int) bool {
			return proxies[i].GetLoad() < proxies[j].GetLoad()
		})
		return getAvailableProxy(proxies, 0)
	}

	return nil, fmt.Errorf("no proxies set")
}
