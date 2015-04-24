package timedmap

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type item struct {
	value   interface{}
	ttl     int
	basettl int
}

type TimedMap struct {
	m      map[string]*item
	lock   sync.Mutex
	ticker *time.Ticker
}

func NewTimedMap() *TimedMap {
	m := &TimedMap{}
	m = startTTL(m)

	return m
}

func startTTL(m *TimedMap) *TimedMap {
	m.m = make(map[string]*item)

	m.ticker = time.NewTicker(time.Second)
	go func() {
		for {
			select {
			case <-m.ticker.C:
				m.lock.Lock()
				for k, v := range m.m {
					// if ttl < 0 then permanent key
					if v.ttl < 0 {
						continue
					} else {
						v.ttl = v.ttl - 1
					}

					if v.ttl == 0 {
						delete(m.m, k)
					}
				}
				m.lock.Unlock()
			}
		}

	}()

	return m
}

func (m *TimedMap) Len() int {
	return len(m.m)
}

func (m *TimedMap) Put(key string, value interface{}) (interface{}, error) {
	// use default ttl
	// default ttl is -1
	// -1 means permanent

	m.lock.Lock()
	defer m.lock.Unlock()

	ret := &item{}

	// check key already exist
	i, ok := m.m[key]
	if ok {
		// key already exist
		ret = i
	} else {
		// make new key
		if m.m == nil {
			m = startTTL(m)
		}
	}

	i = &item{value: value, ttl: -1, basettl: -1}
	m.m[key] = i

	return ret.value, nil
}

func (m *TimedMap) PutRaw(key string, value interface{}, ttl int) (interface{}, error) {
	// use given ttl

	if ttl < 0 {
		return nil, errors.New("ttl should bigger than 0")
	}

	m.lock.Lock()
	defer m.lock.Unlock()

	ret := &item{}

	// check key already exist
	i, ok := m.m[key]
	if ok {
		// key already exist
		ret = i
	} else {
		// make new key
		if m.m == nil {
			m = startTTL(m)
		}
	}

	new := &item{value: value, ttl: ttl, basettl: ttl}
	m.m[key] = new

	return ret.value, nil
}

func (m *TimedMap) Get(key string) (interface{}, error) {
	if i, ok := m.m[key]; ok {
		// Yeah, found the key !
		// update ttl
		m.lock.Lock()
		defer m.lock.Unlock()

		i.ttl = i.basettl
		return i.value, nil
	} else {
		// ooops can not find the key
		return nil, errors.New("This map contains no mapping for the key")
	}
}

func (m *TimedMap) GetTTL(key string) (int, error) {
	if i, ok := m.m[key]; ok {
		// Yeah, found the key !
		// return ttl only without update
		return i.ttl, nil
	} else {
		// ooops can not find the key
		return 0, errors.New("This map Contains no mappging for the key")
	}
}

func (m *TimedMap) Touch(key string) error {
	if i, ok := m.m[key]; ok {
		// Yeah, found the key !
		m.lock.Lock()
		defer m.lock.Unlock()
		i.ttl = i.basettl
		return nil
	} else {
		// ooops can not find the key
		return errors.New("This map Contains no mappging for the key")
	}
}

func (m *TimedMap) ContainsKey(key string) bool {
	if _, ok := m.m[key]; ok {
		// Yeah, key exist!
		return true
	} else {
		// ooops can not find the key
		return false
	}
}

func (m *TimedMap) Remove(key string) (interface{}, error) {
	if i, ok := m.m[key]; ok {
		// Yeah, key exist!
		// delete the key
		ret := i.value
		m.lock.Lock()
		defer m.lock.Unlock()

		delete(m.m, key)
		return ret, nil
	} else {
		// ooops can not find the key
		return nil, nil
	}
}

func (m *TimedMap) Clear() {
	m.lock.Lock()
	for k, _ := range m.m {
		delete(m.m, k)
	}
	m.lock.Unlock()
}

func (m *TimedMap) PrintMap() {
	for k, v := range m.m {
		fmt.Printf("%s : %v - TTL: %d\n", k, v.value, v.ttl)
	}
}
