package safemap

import (
	"sync"
)

func New(inputMap map[interface{}]interface{}) *SafeMap {
	return &SafeMap{
		lock: new(sync.RWMutex),
		data: inputMap,
	}
}

type SafeMap struct {
	lock *sync.RWMutex
	data map[interface{}]interface{}
}

//Get whole map
func (sm *SafeMap) GetMap() map[interface{}]interface{} {
	sm.lock.RLock()
	defer sm.lock.RUnlock()

	return sm.data
}

//Get from maps return the k's value
func (sm *SafeMap) Get(k interface{}) interface{} {
	sm.lock.RLock()
	defer sm.lock.RUnlock()

	if val, ok := sm.data[k]; ok {
		return val
	}
	return nil
}

// Maps the given key and value. Returns bool
func (sm *SafeMap) Set(k interface{}, v interface{}) bool {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	sm.data[k] = v

	return true
}

// Returns true if k is exist in the map.
func (sm *SafeMap) Check(k interface{}) bool {
	sm.lock.RLock()
	defer sm.lock.RUnlock()

	if _, ok := sm.data[k]; !ok {
		return false
	}

	return true
}

func (sm *SafeMap) Delete(k interface{}) {
	sm.lock.Lock()
	defer sm.lock.Unlock()

	delete(sm.data, k)
}
