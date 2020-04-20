package iterator

import (
	"sync"

	"github.com/sotnikov-s/go-load-balancer/proxy"
)

// NewWeightedRoundRobin accepts a map from a proxy to its weight and returns the iterator which
// will switch between proxies depending on their weight
func NewWeightedRoundRobin(proxies map[*proxy.Proxy]int32) Iterator {
	bunch := make(weightedProxiesBunch, 0, len(proxies))
	for p, w := range proxies {
		bunch = append(bunch, &proxyWithWeight{
			Proxy:  p,
			weight: w,
		})
	}

	return &WeightedRoundRobin{
		proxies: bunch,
	}
}

// WeightedRoundRobin is like the round robin iterator but with possibility to set
// weights to proxies
type WeightedRoundRobin struct {
	proxies weightedProxiesBunch

	mu       sync.Mutex
	current  int32
	reqCount int32
}

// Next returns the next proxy or the current one depending on its usage
func (w *WeightedRoundRobin) Next() (*proxy.Proxy, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	currentProxy := w.proxies[w.current]
	if w.reqCount < currentProxy.weight {
		w.reqCount++
	} else {
		w.current = (w.current + 1) % int32(len(w.proxies))
		w.reqCount = 1
	}
	return getAvailableProxy(w.proxies, int(w.current))
}
