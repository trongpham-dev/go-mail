package memcache

import (
	"go-mail/common"
	"sync"
	"time"
)

type Caching interface {
	Write(k string, value []uint32)
	Read(k string) []uint32
}

type caching struct {
	store  map[string][]uint32
	locker *sync.RWMutex
}

func NewCaching() *caching {
	return &caching{
		store:  make(map[string][]uint32),
		locker: new(sync.RWMutex),
	}
}

func (c *caching) Write(k string, value []uint32) {
	c.locker.Lock()
	defer c.locker.Unlock()
	c.store[k] = value
}

func (c *caching) Read(k string) []uint32 {
	c.locker.RLock()
	defer c.locker.RUnlock()
	return c.store[k]
}

func (c *caching) WriteTTL(k string, value []uint32, exp int) {
	c.locker.Lock()
	defer c.locker.Unlock()
	c.store[k] = value

	go func() {
		defer common.AppRecover()
		<-time.NewTimer(time.Second * time.Duration(exp)).C
		c.Write(k, nil)
	}()

}
