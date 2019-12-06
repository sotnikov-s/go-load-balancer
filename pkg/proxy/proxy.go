package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/sotnikov-s/go-load-balancer/pkg/proxy/health"
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
}

// ServeHTTP proxies incoming requests
func (p Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
