package sphere

import (
	"encoding/json"
	"hash/fnv"
	"sync"
)

// connectionShardCount nums of shard
var connectionShardCount = 32

// connectionmap is a "thread" safe map of type string:Anything.
// To avoid lock bottlenecks this map is dived to several (connectionShardCount) map shards.
type connectionmap []*connectionmapshared

// connectionmapshared is a "thread" safe string to anything map.
type connectionmapshared struct {
	items        map[string]*Connection
	sync.RWMutex // Read Write mutex, guards access to internal map.
}

// connectiontuple used by the Iter & IterBuffered functions to wrap two variables together over a channel,
type connectiontuple struct {
	Key string
	Val *Connection
}

// newConnectionMap Creates a new concurrent map.
func newConnectionMap() connectionmap {
	m := make(connectionmap, connectionShardCount)
	for i := 0; i < connectionShardCount; i++ {
		m[i] = &connectionmapshared{items: make(map[string]*Connection)}
	}
	return m
}

// GetShard returns shard under given key
func (m connectionmap) GetShard(key string) *connectionmapshared {
	hasher := fnv.New32()
	hasher.Write([]byte(key))
	return m[uint(hasher.Sum32())%uint(connectionShardCount)]
}

// Set sets the given value under the specified key.
func (m *connectionmap) Set(key string, value *Connection) {
	// Get map shard.
	shard := m.GetShard(key)
	shard.Lock()
	defer shard.Unlock()
	shard.items[key] = value
}

// SetIfAbsent sets the given value under the specified key if no value was associated with it.
func (m *connectionmap) SetIfAbsent(key string, value *Connection) bool {
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
func (m connectionmap) Get(key string) (*Connection, bool) {
	// Get shard
	shard := m.GetShard(key)
	shard.RLock()
	defer shard.RUnlock()

	// Get item from shard.
	val, ok := shard.items[key]
	return val, ok
}

// Count returns the number of elements within the map.
func (m connectionmap) Count() int {
	count := 0
	for i := 0; i < connectionShardCount; i++ {
		shard := m[i]
		shard.RLock()
		count += len(shard.items)
		shard.RUnlock()
	}
	return count
}

// Has looks up an item under specified key
func (m *connectionmap) Has(key string) bool {
	// Get shard
	shard := m.GetShard(key)
	shard.RLock()
	defer shard.RUnlock()

	// See if element is within shard.
	_, ok := shard.items[key]
	return ok
}

// Remove removes an element from the map.
func (m *connectionmap) Remove(key string) {
	// Try to get shard.
	shard := m.GetShard(key)
	shard.Lock()
	defer shard.Unlock()
	delete(shard.items, key)
}

// IsEmpty checks if map is empty.
func (m *connectionmap) IsEmpty() bool {
	return m.Count() == 0
}

// Iter returns an iterator which could be used in a for range loop.
func (m connectionmap) Iter() <-chan connectiontuple {
	ch := make(chan connectiontuple)
	go func() {
		// Foreach shard.
		for _, shard := range m {
			// Foreach key, value pair.
			shard.RLock()
			for key, val := range shard.items {
				ch <- connectiontuple{key, val}
			}
			shard.RUnlock()
		}
		close(ch)
	}()
	return ch
}

// IterBuffered returns a buffered iterator which could be used in a for range loop.
func (m connectionmap) IterBuffered() <-chan connectiontuple {
	ch := make(chan connectiontuple, m.Count())
	go func() {
		// Foreach shard.
		for _, shard := range m {
			// Foreach key, value pair.
			shard.RLock()
			for key, val := range shard.items {
				ch <- connectiontuple{key, val}
			}
			shard.RUnlock()
		}
		close(ch)
	}()
	return ch
}

// MarshalJSON reviles connectionmap "private" variables to json marshal.
func (m connectionmap) MarshalJSON() ([]byte, error) {
	// Create a temporary map, which will hold all item spread across shards.
	tmp := make(map[string]*Connection)

	// Insert items to temporary map.
	for item := range m.Iter() {
		tmp[item.Key] = item.Val
	}
	return json.Marshal(tmp)
}
