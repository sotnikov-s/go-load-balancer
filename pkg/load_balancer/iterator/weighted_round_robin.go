package iterator

import (
	"sync"

	"github.com/sotnikov-s/go-load-balancer/pkg/proxy"
)

// NewWeightedRoundRobin accepts a map from a proxy to its weight and returns the iterator which
// will switch between proxies depending on their weight
func NewWeightedRoundRobin(proxies map[*proxy.Proxy]int32) Iterator {
	weightedProxies := make([]*proxyWithWeight, 0, len(proxies))
	for p, w := range proxies {
		weightedProxies = append(weightedProxies, &proxyWithWeight{
			Proxy: p,
			w:     w,
		})
	}

	return &WeightedRoundRobin{
		proxies: weightedProxies,
	}
}

// WeightedRoundRobin is like the round robin iterator but with possibility to set
// weights to proxies
type WeightedRoundRobin struct {
	proxies []*proxyWithWeight

	mu       sync.Mutex
	current  int32
	reqCount int32
}

// Next returns the next proxy or the current one depending on its usage
func (lb *WeightedRoundRobin) Next() *proxy.Proxy {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	// TODO: use proxy.IsAvailable()
	currentProxy := lb.proxies[lb.current]
	if lb.reqCount < currentProxy.w {
		lb.reqCount++
	} else {
		lb.current = (lb.current + 1) % int32(len(lb.proxies))
		lb.reqCount = 1
	}
	return lb.proxies[lb.current].Proxy
}

// proxyWithWeight is the wrapper over the proxy struct with the proxy instance weight
type proxyWithWeight struct {
	*proxy.Proxy
	w int32
}
