package websocket

import "sync"

type Map[K comparable, V any] struct {
	m  map[K]V
	mx *sync.RWMutex
}

func NewMap[K comparable, V any]() Map[K, V] {
	return Map[K, V]{
		m:  make(map[K]V),
		mx: new(sync.RWMutex),
	}
}

func (m Map[K, V]) Get(key K) (V, bool) {
	m.mx.RLock()
	val, ok := m.m[key]
	m.mx.RUnlock()
	return val, ok
}

func (m Map[K, V]) Delete(key K) {
	m.mx.Lock()
	delete(m.m, key)
	m.mx.Unlock()
}

func (m Map[K, V]) Set(key K, value V) {
	m.mx.Lock()
	m.m[key] = value
	m.mx.Unlock()
}

func (m Map[K, V]) Range(handler func(key K, value V) (error, bool)) error {
	if handler == nil {
		return nil
	}
	m.mx.Lock()
	defer m.mx.Unlock()

	for k, v := range m.m {
		err, br := handler(k, v)
		if err != nil {
			return err
		}
		if br {
			return nil
		}
	}
	return nil
}
