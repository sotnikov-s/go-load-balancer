package loadbalancer

import (
	"net/http"
	"sync"

	"github.com/sotnikov-s/go-load-balancer/pkg/proxy"
)

// NewWeightedRoundRobinLoadBalancer is the WeightedRoundRobinLoadBalancer constructor
func NewWeightedRoundRobinLoadBalancer(proxies map[*proxy.Proxy]int32) LoadBalancer {
	weightedProxies := make([]*proxyWithWeight, 0, len(proxies))
	for p, w := range proxies {
		weightedProxies = append(weightedProxies, &proxyWithWeight{
			Proxy: p,
			w:     w,
		})
	}

	return &WeightedRoundRobinLoadBalancer{
		proxies: weightedProxies,
	}
}

// WeightedRoundRobinLoadBalancer is like the simple round robin load balancer but with possibility to set
// weights to proxies
type WeightedRoundRobinLoadBalancer struct {
	proxies []*proxyWithWeight

	mu       sync.Mutex
	current  int32
	reqCount int32
}

// ServeHTTP uses the next proxy to serve the http request
func (lb *WeightedRoundRobinLoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lb.proxies[lb.next()].ServeHTTP(w, r)
}

// next sets the load balancer to the next proxy
func (lb *WeightedRoundRobinLoadBalancer) next() int32 {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	// TODO: use proxy.IsAvailable()
	currentProxy := lb.proxies[lb.current]
	if lb.reqCount < currentProxy.w {
		lb.reqCount++
		return lb.current
	}

	lb.current = (lb.current + 1) % int32(len(lb.proxies))
	lb.reqCount = 1
	return lb.current
}

// proxyWithWeight is the wrapper over the proxy struct to make it possible to set its weight
type proxyWithWeight struct {
	*proxy.Proxy
	w int32
}
