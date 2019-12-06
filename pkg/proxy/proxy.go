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
		origin:      addr,
		proxy:       httputil.NewSingleHostReverseProxy(addr),
		isAvailable: true,
	}
}

// Proxy is a simple http proxy entity
type Proxy struct {
	origin *url.URL
	proxy  *httputil.ReverseProxy

	// TODO: move all health related stuff to distinct struct
	healthMutex   *sync.RWMutex
	healthChecker func(addr *url.URL) bool
	healthCancel  chan struct{}
	isAvailable   bool
}

func (p *Proxy) WithHealthCheck(checkFunc func(addr *url.URL) bool, period time.Duration) *Proxy {
	if p.healthChecker != nil {
		close(p.healthCancel)
	}
	p.healthMutex = &sync.RWMutex{}
	p.healthChecker = checkFunc
	p.healthCancel = make(chan struct{})
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
		isAvailable := p.healthChecker(p.origin)
		p.healthMutex.Lock()
		defer p.healthMutex.Unlock()
		p.isAvailable = isAvailable
	}

	t := time.NewTicker(period)
	for {
		select {
		case <-t.C:
			checkHealth()
		case <-p.healthCancel:
			t.Stop()
			return
		default:
			runtime.Gosched()
		}
	}
}

// ServeHTTP proxies incoming requests
func (p Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.proxy.ServeHTTP(w, r)
}
