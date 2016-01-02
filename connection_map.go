package sphere

import (
	"encoding/json"
	"hash/fnv"
	"sync"
)

// ConnectionShardCount nums of shard
var ConnectionShardCount = 32

// ConnectionMap is a "thread" safe map of type string:Anything.
// To avoid lock bottlenecks this map is dived to several (ConnectionShardCount) map shards.
type ConnectionMap []*ConnectionMapShared

// ConnectionMapShared is a "thread" safe string to anything map.
type ConnectionMapShared struct {
	items        map[string]*Connection
	sync.RWMutex // Read Write mutex, guards access to internal map.
}

// ConnectionTuple used by the Iter & IterBuffered functions to wrap two variables together over a channel,
type ConnectionTuple struct {
	Key string
	Val *Connection
}

// NewConnectionMap Creates a new concurrent map.
func NewConnectionMap() ConnectionMap {
	m := make(ConnectionMap, ConnectionShardCount)
	for i := 0; i < ConnectionShardCount; i++ {
		m[i] = &ConnectionMapShared{items: make(map[string]*Connection)}
	}
	return m
}

// GetShard returns shard under given key
func (m ConnectionMap) GetShard(key string) *ConnectionMapShared {
	hasher := fnv.New32()
	hasher.Write([]byte(key))
	return m[uint(hasher.Sum32())%uint(ConnectionShardCount)]
}

// Set sets the given value under the specified key.
func (m *ConnectionMap) Set(key string, value *Connection) {
	// Get map shard.
	shard := m.GetShard(key)
	shard.Lock()
	defer shard.Unlock()
	shard.items[key] = value
}

// SetIfAbsent sets the given value under the specified key if no value was associated with it.
func (m *ConnectionMap) SetIfAbsent(key string, value *Connection) bool {
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
func (m ConnectionMap) Get(key string) (*Connection, bool) {
	// Get shard
	shard := m.GetShard(key)
	shard.RLock()
	defer shard.RUnlock()

	// Get item from shard.
	val, ok := shard.items[key]
	return val, ok
}

// Count returns the number of elements within the map.
func (m ConnectionMap) Count() int {
	count := 0
	for i := 0; i < ConnectionShardCount; i++ {
		shard := m[i]
		shard.RLock()
		count += len(shard.items)
		shard.RUnlock()
	}
	return count
}

// Has looks up an item under specified key
func (m *ConnectionMap) Has(key string) bool {
	// Get shard
	shard := m.GetShard(key)
	shard.RLock()
	defer shard.RUnlock()

	// See if element is within shard.
	_, ok := shard.items[key]
	return ok
}

// Remove removes an element from the map.
func (m *ConnectionMap) Remove(key string) {
	// Try to get shard.
	shard := m.GetShard(key)
	shard.Lock()
	defer shard.Unlock()
	delete(shard.items, key)
}

// IsEmpty checks if map is empty.
func (m *ConnectionMap) IsEmpty() bool {
	return m.Count() == 0
}

// Iter returns an iterator which could be used in a for range loop.
func (m ConnectionMap) Iter() <-chan ConnectionTuple {
	ch := make(chan ConnectionTuple)
	go func() {
		// Foreach shard.
		for _, shard := range m {
			// Foreach key, value pair.
			shard.RLock()
			for key, val := range shard.items {
				ch <- ConnectionTuple{key, val}
			}
			shard.RUnlock()
		}
		close(ch)
	}()
	return ch
}

// IterBuffered returns a buffered iterator which could be used in a for range loop.
func (m ConnectionMap) IterBuffered() <-chan ConnectionTuple {
	ch := make(chan ConnectionTuple, m.Count())
	go func() {
		// Foreach shard.
		for _, shard := range m {
			// Foreach key, value pair.
			shard.RLock()
			for key, val := range shard.items {
				ch <- ConnectionTuple{key, val}
			}
			shard.RUnlock()
		}
		close(ch)
	}()
	return ch
}

// MarshalJSON reviles ConnectionMap "private" variables to json marshal.
func (m ConnectionMap) MarshalJSON() ([]byte, error) {
	// Create a temporary map, which will hold all item spread across shards.
	tmp := make(map[string]*Connection)

	// Insert items to temporary map.
	for item := range m.Iter() {
		tmp[item.Key] = item.Val
	}
	return json.Marshal(tmp)
}
