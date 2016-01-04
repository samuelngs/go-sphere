package sphere

import (
	"encoding/json"
	"hash/fnv"
	"sync"
)

// eventModelShardCount nums of shard
var eventModelShardCount = 32

// eventmodelmap is a "thread" safe map of type string:Anything.
// To avoid lock bottlenecks this map is dived to several (eventModelShardCount) map shards.
type eventmodelmap []*eventmodelmapshared

// eventmodelmapshared is a "thread" safe string to anything map.
type eventmodelmapshared struct {
	items        map[string]IEvents
	sync.RWMutex // Read Write mutex, guards access to internal map.
}

// neweventModelMap Creates a new concurrent map.
func newEventModelMap() eventmodelmap {
	m := make(eventmodelmap, eventModelShardCount)
	for i := 0; i < eventModelShardCount; i++ {
		m[i] = &eventmodelmapshared{items: make(map[string]IEvents)}
	}
	return m
}

// eventmodeltuple used by the Iter & IterBuffered functions to wrap two variables together over a event,
type eventmodeltuple struct {
	Key string
	Val IEvents
}

// GetShard returns shard under given key
func (m eventmodelmap) GetShard(key string) *eventmodelmapshared {
	hasher := fnv.New32()
	hasher.Write([]byte(key))
	return m[uint(hasher.Sum32())%uint(eventModelShardCount)]
}

// Set sets the given value under the specified key.
func (m *eventmodelmap) Set(key string, value IEvents) {
	// Get map shard.
	shard := m.GetShard(key)
	shard.Lock()
	defer shard.Unlock()
	shard.items[key] = value
}

// SetIfAbsent sets the given value under the specified key if no value was associated with it.
func (m *eventmodelmap) SetIfAbsent(key string, value IEvents) bool {
	// Get map shard.
	shard := m.GetShard(key)
	shard.Lock()
	defer shard.Unlock()
	_, ok := shard.items[key]
	if !ok {
		shard.items[key] = value
	}
	return !ok
}

// Get retrieves an element from map under given key.
func (m eventmodelmap) Get(key string) (IEvents, bool) {
	// Get shard
	shard := m.GetShard(key)
	shard.RLock()
	defer shard.RUnlock()

	// Get item from shard.
	val, ok := shard.items[key]
	return val, ok
}

// Count returns the number of elements within the map.
func (m eventmodelmap) Count() int {
	count := 0
	for i := 0; i < eventModelShardCount; i++ {
		shard := m[i]
		shard.RLock()
		count += len(shard.items)
		shard.RUnlock()
	}
	return count
}

// Has looks up an item under specified key
func (m *eventmodelmap) Has(key string) bool {
	// Get shard
	shard := m.GetShard(key)
	shard.RLock()
	defer shard.RUnlock()

	// See if element is within shard.
	_, ok := shard.items[key]
	return ok
}

// Remove removes an element from the map.
func (m *eventmodelmap) Remove(key string) {
	// Try to get shard.
	shard := m.GetShard(key)
	shard.Lock()
	defer shard.Unlock()
	delete(shard.items, key)
}

// IsEmpty checks if map is empty.
func (m *eventmodelmap) IsEmpty() bool {
	return m.Count() == 0
}

// Iter returns an iterator which could be used in a for range loop.
func (m eventmodelmap) Iter() <-chan eventmodeltuple {
	ch := make(chan eventmodeltuple)
	go func() {
		// Foreach shard.
		for _, shard := range m {
			// Foreach key, value pair.
			shard.RLock()
			for key, val := range shard.items {
				ch <- eventmodeltuple{key, val}
			}
			shard.RUnlock()
		}
		close(ch)
	}()
	return ch
}

// IterBuffered returns a buffered iterator which could be used in a for range loop.
func (m eventmodelmap) IterBuffered() <-chan eventmodeltuple {
	ch := make(chan eventmodeltuple, m.Count())
	go func() {
		// Foreach shard.
		for _, shard := range m {
			// Foreach key, value pair.
			shard.RLock()
			for key, val := range shard.items {
				ch <- eventmodeltuple{key, val}
			}
			shard.RUnlock()
		}
		close(ch)
	}()
	return ch
}

// MarshalJSON reviles eventmodelmap "private" variables to json marshal.
func (m eventmodelmap) MarshalJSON() ([]byte, error) {
	// Create a temporary map, which will hold all item spread across shards.
	tmp := make(map[string]IEvents)

	// Insert items to temporary map.
	for item := range m.Iter() {
		tmp[item.Key] = item.Val
	}
	return json.Marshal(tmp)
}
