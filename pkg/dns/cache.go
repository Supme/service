package dns

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

type MXCache struct {
	expirationDuration time.Duration
	caches             map[string]domainStruct
	mu                 sync.RWMutex
}

type domainStruct struct {
	expiration time.Time
	value      []*net.MX
	err        error
}

var (
	ErrorHostDoesNotHaveMX       = fmt.Errorf("host does not have MX")
	ErrorMXHostInReservedIPRange = fmt.Errorf("MX host in reserved IP range")
	ErrorOther                   = fmt.Errorf("other error")
)

func NewMXCache(ttlSecond int) *MXCache {
	cache := MXCache{
		expirationDuration: time.Duration(ttlSecond) * time.Second,
		caches:             map[string]domainStruct{},
	}
	go func() {
		t := time.NewTicker(time.Second)
		defer t.Stop()
		for range t.C {
			cache.clearExpired()
		}
	}()
	return &cache
}

func (c *MXCache) clearExpired() {
	var keyForRemove []string
	c.mu.RLock()
	for k := range c.caches {
		if c.caches[k].expiration.Second() >= time.Now().Second() {
			keyForRemove = append(keyForRemove, k)
		}
	}
	c.mu.RUnlock()
	if len(keyForRemove) > 0 {
		c.mu.Lock()
		for k := range keyForRemove {
			delete(c.caches, keyForRemove[k])
		}
		c.mu.Unlock()
	}
}

func (c *MXCache) Get(domain string) ([]*net.MX, error) {
	domain = strings.ToLower(strings.TrimSpace(domain))
	c.mu.RLock()
	mxCached, ok := c.caches[domain]
	c.mu.RUnlock()
	if !ok {
		mx, err := c.getFromDNS(domain)
		mxCached.value = mx
		mxCached.err = err
		mxCached.expiration = time.Now().Add(c.expirationDuration)
		c.mu.Lock()
		c.caches[domain] = mxCached
		c.mu.Unlock()
	}
	return mxCached.value, mxCached.err
}

func (c *MXCache) getFromDNS(domain string) ([]*net.MX, error) {
	var (
		mx  []*net.MX
		err error
	)

	for tries := 0; tries < 10; tries++ {
		mx, err = net.LookupMX(domain)
		if err == nil {
			break
		}
		if strings.HasSuffix(err.Error(), "no such host") {
			return mx, ErrorHostDoesNotHaveMX
		}
		time.Sleep(500 * time.Millisecond)
	}
	// Other error?
	if err != nil {
		return mx, ErrorOther
	}

	if len(mx) == 0 {
		return mx, ErrorHostDoesNotHaveMX
	}

	bad := 0
	for i := range mx {
		a, err := net.LookupIP(mx[i].Host)
		if err != nil {
			bad++
			continue
		}
		for n := range a {
			if a[n].IsLoopback() || a[n].IsMulticast() || a[n].Equal(net.IPv4bcast) || a[n].IsUnspecified() {
				bad++
			}
		}
	}
	if len(mx)-bad < 1 {
		return mx, ErrorMXHostInReservedIPRange
	}

	return mx, nil

}
