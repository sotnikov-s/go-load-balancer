package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
)

// NewProxy is the Proxy constructor
func NewProxy(addr *url.URL) *Proxy {
	return &Proxy{
		ReverseProxy: httputil.NewSingleHostReverseProxy(addr),
	}
}

// Proxy is a simple http proxy entity
type Proxy struct {
	*httputil.ReverseProxy
}

// ServeHTTP proxies incoming requests
func (p Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.ReverseProxy.ServeHTTP(w, r)
}
