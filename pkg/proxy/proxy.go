package proxy

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"runtime"
	"sync"
	"time"
)

// NewProxy is the Proxy constructor
func NewProxy(addr *url.URL) *Proxy {
	return &Proxy{
		origin:       addr,
		ReverseProxy: httputil.NewSingleHostReverseProxy(addr),
		isAvailable:  true,
	}
}

// Proxy is a simple http proxy entity
type Proxy struct {
	origin *url.URL
	*httputil.ReverseProxy

	// TODO: move all health related stuff to distinct struct
	healthMutex   *sync.RWMutex
	healthChecker func(addr *url.URL) bool
	isAvailable   bool
}

func (p *Proxy) WithHealthCheck(checkFunc func(addr *url.URL) bool, period time.Duration) *Proxy {
	// TODO: fix mutex race
	p.healthMutex = &sync.RWMutex{}
	p.healthChecker = checkFunc
	go p.runHealthCheck(period)

	return p
}

func (p *Proxy) IsAvailable() bool {
	p.healthMutex.RLock()
	defer p.healthMutex.RUnlock()
	return p.isAvailable
}

func (p *Proxy) runHealthCheck(period time.Duration) {
	checkHealth := func() {
		p.healthMutex.Lock()
		defer p.healthMutex.Unlock()
		p.isAvailable = p.healthChecker(p.origin)
	}

	// TODO: use cancel channel to stop the goroutine
	t := time.NewTicker(period)
	for {
		select {
		case <-t.C:
			checkHealth()
		default:
			runtime.Gosched()
		}
	}
}

// ServeHTTP proxies incoming requests
func (p Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.ReverseProxy.ServeHTTP(w, r)
}
