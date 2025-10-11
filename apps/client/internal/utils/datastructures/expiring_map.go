package datastructures

import (
	"sync"
	"time"
)

type ExpiringMap[K comparable, V any] struct {
	vals        map[K]V
	expirations map[K]time.Time
	mutex       sync.Mutex
}

func (m *ExpiringMap[K, V]) SetExpiring(key K, val V, expiresIn time.Duration) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if expiresIn != 0 {
		m.expirations[key] = time.Now().Add(expiresIn)
	} else {
		m.expirations[key] = time.Time{}
	}

	m.vals[key] = val
}

func (m *ExpiringMap[K, V]) Set(key K, val V) {
	m.SetExpiring(key, val, 0)
}

func (m *ExpiringMap[K, V]) Get(key K) *V {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	return m.get(key)
}

func (m *ExpiringMap[K, V]) Delete(key K) *V {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	val := m.get(key)
	if val == nil {
		return nil
	}

	m.delete(key)

	return val
}

func (m *ExpiringMap[K, V]) delete(key K) {
	delete(m.vals, key)
	delete(m.expirations, key)
}

func (m *ExpiringMap[K, V]) get(key K) *V {
	expireTime, ok := m.expirations[key]
	if !ok {
		return nil
	}

	if !expireTime.IsZero() && expireTime.Before(time.Now()) {
		m.delete(key)

		return nil
	}

	val := m.vals[key]

	return &val
}
