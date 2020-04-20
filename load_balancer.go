package loadbalancer

import (
	"net/http"

	"github.com/sotnikov-s/go-load-balancer/iterator"
)

// NewLoadBalancer returns the LoadBalancer instance with the specified iterator
func NewLoadBalancer(iterator iterator.Iterator) *LoadBalancer {
	return &LoadBalancer{iter: iterator}
}

// LoadBalancer handles requests distributing them between one or more proxies
type LoadBalancer struct {
	iter iterator.Iterator
}

// ServeHTTP handles the request by the next proxy gotten from the iterator
func (l *LoadBalancer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p, err := l.iter.Next()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("The server didn't respond"))
		return
	}
	p.ServeHTTP(w, r)
}
