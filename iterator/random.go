package iterator

import (
	"fmt"
	"math/rand"

	"github.com/sotnikov-s/go-load-balancer/proxy"
)

// NewRandom accepts a number of proxies to be used in the Random iterator
// and returns the Random iterator instance itself. The seed represents a
// func to initialize the random number generator state.
func NewRandom(seed func(), proxies ...*proxy.Proxy) Iterator {
	bunch := make(commonProxiesBunch, 0, len(proxies))
	for _, p := range proxies {
		bunch = append(bunch, p)
	}
	seed()
	return &Random{
		proxies: bunch,
	}
}

// Random chooses a random proxy of the bunch and uses it as the handler
type Random struct {
	proxies commonProxiesBunch
}

// Next returns a random proxy of the bunch
func (r *Random) Next() (*proxy.Proxy, error) {
	if r.proxies.Len() > 0 {
		n := rand.Intn(r.proxies.Len())
		return getAvailableProxy(r.proxies, n)
	}
	return nil, fmt.Errorf("no proxies set")
}
