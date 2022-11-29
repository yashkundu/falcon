package utils

import (
	"container/list"
	"sync"
	"time"
)

// ephemeral-cache -- stays temporarily for some time

type timeKey struct {
	inTime time.Time
	key    interface{}
}

type ECache struct {
	keyQueue *list.List
	repeats  map[interface{}]int
	datas    map[interface{}]interface{}
	mu       *sync.Mutex
	ttl      time.Duration
}

func NewECache(ttl time.Duration) *ECache {
	ec := new(ECache)
	ec.datas = make(map[interface{}]interface{})
	ec.repeats = make(map[interface{}]int)
	ec.keyQueue = list.New()
	ec.mu = new(sync.Mutex)
	return ec
}

func (ec *ECache) Set(key, data interface{}) {
	ec.mu.Lock()
	defer ec.mu.Unlock()
	ec.refresh()

	_, exist := ec.datas[key]
	ec.datas[key] = data
	tkey := timeKey{
		inTime: time.Now(),
		key:    key,
	}
	ec.keyQueue.PushBack(tkey)

	if exist {
		val, ok := ec.repeats[key]
		if ok {
			ec.repeats[key] = val + 1
		} else {
			ec.repeats[key] = 1
		}
	}

}

func (ec *ECache) Get(key interface{}) (interface{}, bool) {
	ec.mu.Lock()
	defer ec.mu.Unlock()
	ec.refresh()
	d, ok := ec.datas[key]
	return d, ok
}

func (ec *ECache) Len() int {
	ec.mu.Lock()
	defer ec.mu.Unlock()
	ec.refresh()
	return len(ec.datas)
}

func (ec *ECache) Update() {
	ec.mu.Lock()
	defer ec.mu.Unlock()
	ec.refresh()
}

// removes the expired items
func (ec *ECache) refresh() {
	if ec.keyQueue.Len() == 0 {
		return
	}
	for {
		item := ec.keyQueue.Front()
		if item == nil {
			return
		}
		tKey := item.Value.(timeKey)
		if time.Now().Sub(tKey.inTime) > ec.ttl {
			ec.keyQueue.Remove(item)
			if val, ok := ec.repeats[tKey.key]; !ok {
				delete(ec.datas, tKey.key)
			} else {
				if val <= 1 {
					delete(ec.repeats, tKey.key)
				} else {
					ec.repeats[tKey.key] = val - 1
				}
			}
			continue
		} else {
			break
		}
	}
}
