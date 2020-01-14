package iterator

import (
	"fmt"
	"sync"

	"github.com/sotnikov-s/go-load-balancer/pkg/proxy"
)

// NewWeightedRoundRobin accepts a map from a proxy to its weight and returns the iterator which
// will switch between proxies depending on their weight
func NewWeightedRoundRobin(proxies map[*proxy.Proxy]int32) Iterator {
	weightedProxies := make([]*proxyWithWeight, 0, len(proxies))
	for p, w := range proxies {
		weightedProxies = append(weightedProxies, &proxyWithWeight{
			Proxy:  p,
			weight: w,
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
	return w.getAvailableProxy(int(w.current))
}

// getAvailableProxy walks through the proxies and returns the first available one starting from
// the one at the marker index. If no available proxy was found, it returns an error
func (w *WeightedRoundRobin) getAvailableProxy(marker int) (*proxy.Proxy, error) {
	for i := 0; i < len(w.proxies); i++ {
		tryProxy := (marker + i) % len(w.proxies)
		p := w.proxies[tryProxy]
		if p.IsAvailable() {
			return p.Proxy, nil
		}
	}
	return nil, fmt.Errorf("all proxies are unavailable")
}

// proxyWithWeight is the wrapper over the proxy struct with the proxy instance weight
type proxyWithWeight struct {
	*proxy.Proxy
	weight int32
}
