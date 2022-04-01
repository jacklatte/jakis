package jakis

import "sync"

type Dict interface {
	Get(key string) string
	Set(key string, val string)
}

type SimpleMap struct {
	dict map[string]string
	mu   sync.Mutex
}

func (m *SimpleMap) Get(key string) string {
	m.mu.Lock()
	val := m.dict[key]
	m.mu.Unlock()
	return val
}

func (m *SimpleMap) Set(key string, val string) {
	m.mu.Lock()
	m.dict[key] = val
	m.mu.Unlock()
}

func NewSimpleMap() *SimpleMap {
	return &SimpleMap{}
}

type SyncMap struct {
	dict sync.Map
}

func (m *SyncMap) Get(key string) string {
	output, _ := m.dict.Load(key)
	val, _ := output.(string)
	return val
}

func (m *SyncMap) Set(key string, val string) {
	m.dict.Store(key, val)
}

func NewSyncMap() *SyncMap {
	return &SyncMap{}
}
