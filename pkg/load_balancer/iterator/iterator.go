package iterator

import "github.com/sotnikov-s/go-load-balancer/pkg/proxy"

// Iterator is the iterator pattern implementation created to iterate over proxies
type Iterator interface {
	Next() *proxy.Proxy
}
