package store

import (
	"sort"
	"sync"
	"time"
)

// Store defines requirements for an store implementation
type Store interface {
	Set(key, value string) error
	Get(key string) (string, bool)
	Delete(key string) error
	Cap() int
	GetLastModifiedKeys() []string
}

// MemoryStore uses a sync.Map to implement the Store interface.
type MemoryStore struct {
	s sync.Map

	lastMod sync.Map
	lastAcc sync.Map

	// TODO: This could be a circular buff
	cachedAcc sort.IntSlice

	cap int
	mc  sync.Mutex
}

// NewMemoryStore creates a new instance of memoryStore
// with the internal map initialized.
func NewMemoryStore(cap int) *MemoryStore {
	return &MemoryStore{
		cap: cap,
	}
}

// Cap returns available capacity
func (m *MemoryStore) Cap() int {
	return m.cap
}

// Set receives key and value strings and saves the key/value in the internal map.
func (m *MemoryStore) Set(key, value string) error {
	if m.cap < len(value) {
		m.clean(len(value))
	}

	m.s.Store(key, value)

	// Increase num
	m.trackModifield(key)

	// Reduce capacity
	m.mc.Lock()
	m.cap -= len(value)
	m.mc.Unlock()

	return nil
}

// Get receives a key string and return the value and a boolean.
func (m *MemoryStore) Get(key string) (string, bool) {
	v, ok := m.s.Load(key)
	if !ok {
		return "", false
	}

	// Push last accessed item
	m.trackAccess(key)

	return v.(string), ok
}

// Delete receives a key string and deletes its value from the internal map.
func (m *MemoryStore) Delete(key string) error {
	v, _ := m.Get(key)
	m.s.Delete(key)

	m.cleanAccess(key)

	m.mc.Lock()
	m.cap += len(v)
	m.mc.Unlock()
	return nil
}

func (m *MemoryStore) GetLastModifiedKeys() []string {
	var (
		i sort.IntSlice
		r []string
	)
	m.lastAcc.Range(func(key interface{}, value interface{}) bool {
		i = append(i, key.(int))
		return true
	})
	sort.Sort(sort.Reverse(i))
	for _, t := range i {
		v, _ := m.lastMod.Load(t)
		r = append(r, v.(string))
	}
	return r
}

func (m *MemoryStore) trackModifield(key string) {
	m.lastMod.Store(time.Now().Unix(), key)
}

func (m *MemoryStore) trackAccess(key string) {
	t := int(time.Now().Unix())
	m.lastAcc.Store(t, key)
	m.cachedAcc = append(m.cachedAcc, t)
}

func (m *MemoryStore) cleanAccess(key string) {
	m.lastAcc.Range(func(t interface{}, value interface{}) bool {
		if value.(string) == key {
			m.lastAcc.Delete(key)
			return false
		}
		return true
	})

	m.lastMod.Range(func(t interface{}, value interface{}) bool {
		if value.(string) == key {
			m.lastMod.Delete(key)
			return false
		}
		return true
	})
}

func (m *MemoryStore) clean(size int) {
	if !sort.IsSorted(m.cachedAcc) {
		m.lastAcc.Range(func(key interface{}, value interface{}) bool {
			m.cachedAcc = append(m.cachedAcc, key.(int))
			return true
		})
		sort.Sort(m.cachedAcc)
	}

	for m.cap < size && len(m.cachedAcc) > 0 {
		k, _ := m.lastAcc.Load(m.cachedAcc[0])
		m.cachedAcc = append(m.cachedAcc[:0], m.cachedAcc[1:]...)
		m.Delete(k.(string))
	}
}
