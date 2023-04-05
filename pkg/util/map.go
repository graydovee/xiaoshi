package util

import "sync"

type Map[k comparable, V any] struct {
	m sync.Map
}

func (m *Map[K, V]) Load(key K) (V, bool) {
	v, ok := m.m.Load(key)
	if !ok {
		vp := new(V)
		return *vp, false
	}
	return v.(V), true
}

func (m *Map[K, V]) Store(key K, value V) {
	m.m.Store(key, value)
}

func (m *Map[K, V]) Delete(key K) {
	m.m.Delete(key)
}

func (m *Map[K, V]) Range(f func(key K, value V) bool) {
	m.m.Range(func(key, value interface{}) bool {
		return f(key.(K), value.(V))
	})
}

func (m *Map[K, V]) LoadOrStore(key K, newV func() V) (actual V, loaded bool) {
	v, loaded := m.m.Load(key)
	if loaded {
		return v.(V), loaded
	}
	v, loaded = m.m.LoadOrStore(key, newV())
	return v.(V), loaded
}
