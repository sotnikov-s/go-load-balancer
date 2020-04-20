package iterator

import (
	"fmt"

	"github.com/sotnikov-s/go-load-balancer/proxy"
)

type proxiesBunch interface {
	Len() int
	Get(idx int) *proxy.Proxy
}

// getAvailableProxy walks through the proxies and returns the first available one starting from
// the one at the marker index. If no available proxy was found, it returns an error
func getAvailableProxy(proxies proxiesBunch, marker int) (*proxy.Proxy, error) {
	for i := 0; i < proxies.Len(); i++ {
		tryProxy := (marker + i) % proxies.Len()
		p := proxies.Get(tryProxy)
		if p.IsAvailable() {
			return p, nil
		}
	}
	return nil, fmt.Errorf("all proxies are unavailable")
}

// commonProxiesBunch is the simplest proxiesBunch implementation
type commonProxiesBunch []*proxy.Proxy

func (b commonProxiesBunch) Len() int                 { return len(b) }
func (b commonProxiesBunch) Get(idx int) *proxy.Proxy { return b[idx] }

// weightedProxiesBunch is a bunch of the weighted proxies
type weightedProxiesBunch []*proxyWithWeight

func (b weightedProxiesBunch) Len() int                 { return len(b) }
func (b weightedProxiesBunch) Get(idx int) *proxy.Proxy { return b[idx].Proxy }

// proxyWithWeight is the wrapper over the proxy struct with the proxy instance weight
type proxyWithWeight struct {
	*proxy.Proxy
	weight int32
}
