package loadbalancer

import (
	"net/http"

	"github.com/sotnikov-s/go-load-balancer/pkg/load_balancer/iterator"
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
	l.iter.Next().ServeHTTP(w, r)
}
