# golang load balancer

A load balancer is a device that, behaving as a reverse proxy, uses a number of remote servers to improve computing performance by distributing workloads across the servers.

* [Algorithms](#algorithms)
* [Schema](#schema)
* [Example](#example)

---

## Algorithms

The load balancer supports a number of load balancing algorithms:
- round robin
- weighted round robin
- random _(to be implemented)_
- least connections _(to be implemented)_

## Schema

![schema](https://i.ibb.co/tPJT5WN/Screenshot-2020-01-10-at-14-12-55.png)

## Example

Here comes a simple round robin load balancer instantiation and usage:

```golang
func main() {
	shortRespUrl, err := url.Parse("http://127.0.0.1:8081")
	if err != nil {
		log.Fatal(err)
	}
	p1 := proxy.NewProxy(shortRespUrl)

	longRespUrl, err := url.Parse("http://127.0.0.1:8082")
	if err != nil {
		log.Fatal(err)
	}
	p2 := proxy.NewProxy(longRespUrl)

	lb := loadbalancer.NewLoadBalancer(iterator.NewRoundRobin(p1, p2))
	log.Printf("load balancer started at port :8080")
	go func() {
		log.Fatal(http.ListenAndServe(":8080", lb))
	}()

	for i := 0; i < 5; i++ {
		func() {
			r, _ := http.Get("http://127.0.0.1:8080")
			b, _ := ioutil.ReadAll(r.Body)
			log.Printf("got %d resp: %s", i, string(b))
		}()
	}
}
```

Output:
```text
2019/12/04 11:59:23 load balancer started at port :8080
2019/12/04 11:59:24 got 0 resp: --- short resp ---
2019/12/04 11:59:24 got 1 resp: --------- long resp ---------
2019/12/04 11:59:24 got 2 resp: --- short resp ---
2019/12/04 11:59:24 got 3 resp: --------- long resp ---------
2019/12/04 11:59:24 got 4 resp: --- short resp ---
```