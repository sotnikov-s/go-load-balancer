package health

import (
	"net"
	"net/url"
	"sync"
	"time"
)

// NewProxyHealth is the ProxyHealth constructor
func NewProxyHealth(origin *url.URL) *ProxyHealth {
	h := &ProxyHealth{
		origin: origin,
		check:  defaultHealthCheck,
		period: defaultHealthCheckPeriod,
	}
	go h.run()

	return h
}

// ProxyHealth is looking after the proxy origin availability using either a set by
// ProxyHealth.SetHealthCheck check function or the defaultHealthCheck func.
type ProxyHealth struct {
	origin *url.URL

	mu          sync.Mutex
	check       func(addr *url.URL) bool
	period      time.Duration
	cancel      chan struct{}
	isAvailable bool
}

// IsAvailable returns whether the proxy origin was successfully connected at the last check time.
func (h *ProxyHealth) IsAvailable() bool {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.isAvailable
}

// SetHealthCheck sets the passed check func as the algorithm of checking the origin availability and
// calls for it with interval defined with the passed period variable.
func (h *ProxyHealth) SetHealthCheck(check func(addr *url.URL) bool, period time.Duration) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.stop()
	h.check = check
	h.period = period
	h.run()
}

// run runs the check func in a new goroutine.
func (h *ProxyHealth) run() {
	checkHealth := func() {
		h.mu.Lock()
		defer h.mu.Unlock()
		isAvailable := h.check(h.origin)
		h.isAvailable = isAvailable
	}

	h.cancel = make(chan struct{})
	go func() {
		checkHealth()
		t := time.NewTicker(h.period)
		for {
			select {
			case <-t.C:
				checkHealth()
			case <-h.cancel:
				t.Stop()
				return
			}
		}
	}()
}

// stop stops the currently rinning check func.
func (h *ProxyHealth) stop() {
	if h.cancel != nil {
		h.cancel <- struct{}{}
		close(h.cancel)
		h.cancel = nil
	}
}

// defaultHealthCheck is the default most simple check function
var defaultHealthCheck = func(addr *url.URL) bool {
	conn, err := net.DialTimeout("tcp", addr.Host, defaultHealthCheckTimeout)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

var (
	defaultHealthCheckTimeout = 10 * time.Second
	defaultHealthCheckPeriod  = 10 * time.Second
)
