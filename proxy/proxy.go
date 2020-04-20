package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync/atomic"
	"time"

	"github.com/sotnikov-s/go-load-balancer/proxy/health"
)

// NewProxy is the Proxy constructor
func NewProxy(addr *url.URL) *Proxy {
	return &Proxy{
		proxy:  httputil.NewSingleHostReverseProxy(addr),
		health: health.NewProxyHealth(addr),
	}
}

// Proxy is a simple http proxy entity
type Proxy struct {
	health *health.ProxyHealth
	proxy  *httputil.ReverseProxy
	load   int32
}

// ServeHTTP proxies incoming requests
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt32(&p.load, 1)
	defer atomic.AddInt32(&p.load, -1)
	p.proxy.ServeHTTP(w, r)
}

// IsAvailable returns whether the proxy origin was successfully connected at the last check time.
func (p *Proxy) IsAvailable() bool {
	return p.health.IsAvailable()
}

// SetHealthCheck sets the passed check func as the algorithm of checking the origin availability
func (p *Proxy) SetHealthCheck(check func(addr *url.URL) bool, period time.Duration) {
	p.health.SetHealthCheck(check, period)
}

// GetLoad returns the number of requests being served by the proxy at the moment
func (p *Proxy) GetLoad() int32 {
	return atomic.LoadInt32(&p.load)
}
