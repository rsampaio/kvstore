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
	GetSortedLastModifiedList() TimeTrackingAsc
}

// At define time access or modify time tracking format
type At struct {
	Key string
	At  int64
}

// TimeTrackingDesc is a list of time that can be sorted descending
type TimeTrackingDesc []At

// TimeTrackingDesc is a list of time that can be sorted ascending
type TimeTrackingAsc []At

// Implement sort interface
func (a TimeTrackingDesc) Len() int           { return len(a) }
func (a TimeTrackingDesc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a TimeTrackingDesc) Less(i, j int) bool { return a[i].At < a[j].At }

// Implement sort interface
func (a TimeTrackingAsc) Len() int           { return len(a) }
func (a TimeTrackingAsc) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a TimeTrackingAsc) Less(i, j int) bool { return a[i].At > a[j].At }

// MemoryStore uses a sync.Map to implement the Store interface.
type MemoryStore struct {
	s sync.Map

	cap int
	mc  sync.Mutex

	// Last modified items to stream in order
	lastMod TimeTrackingAsc
	mm      sync.Mutex

	// Last accessed items to clean least accessed item
	lastAcc TimeTrackingDesc
	ma      sync.Mutex
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
	m.initAccess(key)

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

	m.mc.Lock()
	m.cap += len(v)
	m.mc.Unlock()

	m.ma.Lock()
	m.removeFromAccess(key)
	m.ma.Unlock()

	m.mm.Lock()
	m.removeFromModified(key)
	m.mm.Unlock()
	return nil
}

func (m *MemoryStore) GetSortedLastModifiedList() TimeTrackingAsc {
	sort.Sort(m.lastMod)
	return m.lastMod
}

func (m *MemoryStore) clean(size int) {
	sort.Sort(m.lastAcc)
	for m.cap < size {
		m.Delete(m.lastAcc[0].Key)
	}
}

func (m *MemoryStore) removeFromModified(key string) {
	for i, v := range m.lastMod {
		if v.Key == key {
			m.lastMod = append(m.lastMod[:i], m.lastMod[i+1:]...)
			break
		}
	}
}

func (m *MemoryStore) removeFromAccess(key string) {
	for i, v := range m.lastAcc {
		if v.Key == key {
			m.lastAcc = append(m.lastAcc[:i], m.lastAcc[i+1:]...)
			break
		}
	}
}

func (m *MemoryStore) initAccess(key string) {
	m.ma.Lock()
	defer m.ma.Unlock()
	m.removeFromAccess(key)
	m.lastAcc = append(m.lastAcc, At{Key: key, At: 0})
}

func (m *MemoryStore) trackAccess(key string) {
	m.ma.Lock()
	defer m.ma.Unlock()
	m.removeFromAccess(key)
	m.lastAcc = append(m.lastAcc, At{Key: key, At: time.Now().Unix()})
}

func (m *MemoryStore) trackModifield(key string) {
	m.mm.Lock()
	defer m.mm.Unlock()
	m.removeFromModified(key)
	m.lastMod = append(m.lastMod, At{Key: key, At: time.Now().Unix()})
}
