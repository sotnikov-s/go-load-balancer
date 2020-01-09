package loadbalancer

import (
	"net/http"
)

// LoadBalancer handles requests distributing them between one or more proxies
type LoadBalancer interface {
	http.Handler
}
