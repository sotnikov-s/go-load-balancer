package health

import (
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTrueIsAvailable(t *testing.T) {
	ts := httptest.NewServer(nil)
	defer ts.Close()
	origin, _ := url.Parse(ts.URL)
	h := NewProxyHealth(origin)

	h.SetHealthCheck(func(_ *url.URL) bool {
		return true
	}, time.Second)

	h.Stop()
	assert.Equalf(t, h.isAvailable, h.IsAvailable(), "unequal isAvailable value and IsAvailable() result")
	assert.Equal(t, true, h.IsAvailable())
}

func TestFalseIsAvailable(t *testing.T) {
	ts := httptest.NewServer(nil)
	defer ts.Close()
	origin, _ := url.Parse(ts.URL)
	h := NewProxyHealth(origin)

	h.SetHealthCheck(func(_ *url.URL) bool {
		return false
	}, time.Second)

	h.Stop()
	assert.Equalf(t, h.isAvailable, h.IsAvailable(), "unequal isAvailable value and IsAvailable() result")
	assert.Equal(t, false, h.IsAvailable())
}

func TestStop(t *testing.T) {
	ts := httptest.NewServer(nil)
	defer ts.Close()
	origin, _ := url.Parse(ts.URL)
	h := NewProxyHealth(origin)

	h.SetHealthCheck(func(_ *url.URL) bool {
		return true
	}, time.Second)
	h.Stop()

	assert.Nilf(t, h.cancel, "after the Stop() call the h.cancel chan should be nil")
}

func TestMultipleSetHealthCheck(t *testing.T) {
	ts := httptest.NewServer(nil)
	defer ts.Close()
	origin, _ := url.Parse(ts.URL)
	h := NewProxyHealth(origin)

	h.SetHealthCheck(func(_ *url.URL) bool {
		return true
	}, time.Second)
	assert.Equal(t, true, h.IsAvailable())

	h.SetHealthCheck(func(_ *url.URL) bool {
		return false
	}, time.Second)
	assert.Equal(t, false, h.IsAvailable())

	h.SetHealthCheck(func(_ *url.URL) bool {
		return true
	}, time.Second)
	assert.Equal(t, true, h.IsAvailable())

	h.Stop()
}

func TestHealthCheckContinuity(t *testing.T) {
	ts := httptest.NewServer(nil)
	defer ts.Close()
	origin, _ := url.Parse(ts.URL)
	h := NewProxyHealth(origin)

	ch := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			got := <-ch
			assert.Equal(t, i, got)
		}
		h.Stop()
	}()

	var res int
	h.SetHealthCheck(func(_ *url.URL) bool {
		ch <- res
		res++
		return true
	}, 10*time.Millisecond)
}
