package local_cache

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

const (
	NoExpiration time.Duration = -1

	DefaultExpiration time.Duration = 0
)

type kv struct {
	key   string
	value interface{}
}

type Item struct {
	Object     interface{}
	Expiration int64
}

type LocalCache struct {
	*localCache
}

type localCache struct {
	defaultExpiration time.Duration
	data              map[string]Item
	mu                sync.RWMutex
	onEvicted         func(string, interface{})
	janitor           *janitor
}

type janitor struct {
	Interval time.Duration
	stop     chan bool
}

func (i *Item) Expired() bool {
	if i.Expiration == 0 {
		return false
	}
	return i.Expiration < time.Now().UnixNano()
}

func (lc *localCache) Set(k string, v interface{}, d time.Duration) {
	var e int64
	if d == DefaultExpiration {
		d = lc.defaultExpiration
	}
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}
	lc.mu.Lock()
	lc.data[k] = Item{
		Object:     v,
		Expiration: e,
	}
	lc.mu.Unlock()
}

func (lc *localCache) set(k string, v interface{}, d time.Duration) {
	var e int64
	if d == DefaultExpiration {
		d = lc.defaultExpiration
	}
	if d > 0 {
		e = time.Now().Add(d).UnixNano()
	}
	lc.data[k] = Item{
		Object:     v,
		Expiration: e,
	}
}

func (lc *localCache) Get(k string) (interface{}, bool) {
	lc.mu.Lock()
	item, found := lc.data[k]
	if !found {
		lc.mu.Unlock()
		return nil, false
	}
	if item.Expiration > 0 && item.Expiration < time.Now().UnixNano() {
		lc.mu.Unlock()
		return nil, false
	}
	lc.mu.Unlock()
	return item.Object, true
}

func (lc *localCache) get(k string) (interface{}, bool) {
	item, found := lc.data[k]
	if !found {
		return nil, false
	}
	if item.Expiration > 0 && item.Expiration < time.Now().UnixNano() {
		return nil, false
	}
	return item.Object, true
}

func (lc *localCache) Add(k string, v interface{}, d time.Duration) error {
	lc.mu.Lock()
	_, found := lc.get(k)
	if found {
		return fmt.Errorf("Item:%s has already exit", k)
	}
	lc.set(k, v, d)
	lc.mu.Unlock()
	return nil
}

func (lc *localCache) Replace(k string, v interface{}, d time.Duration) error {
	lc.mu.Lock()
	_, found := lc.get(k)
	if !found {
		lc.mu.Unlock()
		return fmt.Errorf("Item:%s dosen't exit", k)
	}
	lc.set(k, v, d)
	lc.mu.Unlock()
	return nil
}

func (lc *localCache) Increament(k string, n int64) error {
	lc.mu.Lock()
	v, found := lc.data[k]
	if !found || v.Expired() {
		lc.mu.Unlock()
		return fmt.Errorf("Item not found or expired")
	}
	switch v.Object.(type) {
	case int:
		v.Object = v.Object.(int) + int(n)
	default:
		lc.mu.Unlock()
		return fmt.Errorf("not support value tyepe")
	}
	lc.data[k] = v
	lc.mu.Unlock()
	return nil
}

func (lc *localCache) Delete(k string) {
	lc.mu.Lock()
	v, evicted := lc.delete(k)
	lc.mu.Unlock()
	if evicted {
		lc.onEvicted(k, v)
	}
}

func (lc *localCache) delete(k string) (interface{}, bool) {
	if lc.onEvicted != nil {
		if v, found := lc.data[k]; found {
			delete(lc.data, k)
			return v.Object, true
		}
	}
	delete(lc.data, k)
	return nil, false
}

func (lc *localCache) DeleteExpired() {
	var evictedItems []kv
	timeNow := time.Now().UnixNano()
	lc.mu.Lock()
	for k, v := range lc.data {
		if v.Expiration > 0 && v.Expiration < timeNow {
			v, evicted := lc.delete(k)
			if evicted {
				evictedItems = append(evictedItems, kv{k, v})
			}
		}
	}
	lc.mu.Unlock()
	for _, v := range evictedItems {
		lc.onEvicted(v.key, v.value)
	}
}

func (lc *localCache) GetData() map[string]Item {
	lc.mu.Lock()
	defer lc.mu.Unlock()
	return lc.data
}

func (lc *localCache) GetCount() int {
	lc.mu.Lock()
	n := len(lc.data)
	lc.mu.Unlock()
	return n
}

func (lc *localCache) OnEvicted(f func(string, interface{})) {
	lc.mu.Lock()
	lc.onEvicted = f
	lc.mu.Unlock()
}

func (lc *localCache) Flush() {
	lc.mu.Lock()
	lc.data = map[string]Item{}
	lc.mu.Unlock()
}

func (j *janitor) Run(lc *localCache) {
	j.stop = make(chan bool)
	ticker := time.NewTicker(j.Interval)
	for {
		select {
		case <-ticker.C:
			lc.DeleteExpired()
		case <-j.stop:
			ticker.Stop()
			return

		}
	}
}

func stopJanitor(lc *LocalCache) {
	lc.janitor.stop <- true
}

func runJanitor(lc *localCache, d time.Duration) {
	j := &janitor{
		Interval: d,
	}
	lc.janitor = j
	go j.Run(lc)
}

func newLocalCache(d time.Duration, m map[string]Item) *localCache {
	if d == 0 {
		d = -1
	}
	lc := &localCache{
		defaultExpiration: d,
		data:              m,
	}
	return lc
}

func newCacheWithJanitor(defaultExpiration time.Duration, cleanupInterval time.Duration, data map[string]Item) *LocalCache {
	lc := newLocalCache(defaultExpiration, data)
	C := &LocalCache{lc}
	if cleanupInterval > 0 {
		runJanitor(lc, cleanupInterval)
		runtime.SetFinalizer(C, stopJanitor)
	}
	return C
}

func New(defaultExpiration, cleanupInterval time.Duration) *LocalCache {
	data := make(map[string]Item)
	return newCacheWithJanitor(defaultExpiration, cleanupInterval, data)
}
