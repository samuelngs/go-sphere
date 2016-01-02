package sphere

import (
	"encoding/json"
	"hash/fnv"
	"sync"
)

// ChannelShardCount nums of shard
var ChannelShardCount = 32

// ChannelMap is a "thread" safe map of type string:Anything.
// To avoid lock bottlenecks this map is dived to several (ChannelShardCount) map shards.
type ChannelMap []*ChannelMapShared

// ChannelMapShared is a "thread" safe string to anything map.
type ChannelMapShared struct {
	items        map[string]*Channel
	sync.RWMutex // Read Write mutex, guards access to internal map.
}

// ChannelTuple used by the Iter & IterBuffered functions to wrap two variables together over a channel,
type ChannelTuple struct {
	Key string
	Val *Channel
}

// NewChannelMap Creates a new concurrent map.
func NewChannelMap() ChannelMap {
	m := make(ChannelMap, ChannelShardCount)
	for i := 0; i < ChannelShardCount; i++ {
		m[i] = &ChannelMapShared{items: make(map[string]*Channel)}
	}
	return m
}

// GetShard returns shard under given key
func (m ChannelMap) GetShard(key string) *ChannelMapShared {
	hasher := fnv.New32()
	hasher.Write([]byte(key))
	return m[uint(hasher.Sum32())%uint(ChannelShardCount)]
}

// Set sets the given value under the specified key.
func (m *ChannelMap) Set(key string, value *Channel) {
	// Get map shard.
	shard := m.GetShard(key)
	shard.Lock()
	defer shard.Unlock()
	shard.items[key] = value
}

// SetIfAbsent sets the given value under the specified key if no value was associated with it.
func (m *ChannelMap) SetIfAbsent(key string, value *Channel) bool {
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
func (m ChannelMap) Get(key string) (*Channel, bool) {
	// Get shard
	shard := m.GetShard(key)
	shard.RLock()
	defer shard.RUnlock()

	// Get item from shard.
	val, ok := shard.items[key]
	return val, ok
}

// Count returns the number of elements within the map.
func (m ChannelMap) Count() int {
	count := 0
	for i := 0; i < ChannelShardCount; i++ {
		shard := m[i]
		shard.RLock()
		count += len(shard.items)
		shard.RUnlock()
	}
	return count
}

// Has looks up an item under specified key
func (m *ChannelMap) Has(key string) bool {
	// Get shard
	shard := m.GetShard(key)
	shard.RLock()
	defer shard.RUnlock()

	// See if element is within shard.
	_, ok := shard.items[key]
	return ok
}

// Remove removes an element from the map.
func (m *ChannelMap) Remove(key string) {
	// Try to get shard.
	shard := m.GetShard(key)
	shard.Lock()
	defer shard.Unlock()
	delete(shard.items, key)
}

// IsEmpty checks if map is empty.
func (m *ChannelMap) IsEmpty() bool {
	return m.Count() == 0
}

// Iter returns an iterator which could be used in a for range loop.
func (m ChannelMap) Iter() <-chan ChannelTuple {
	ch := make(chan ChannelTuple)
	go func() {
		// Foreach shard.
		for _, shard := range m {
			// Foreach key, value pair.
			shard.RLock()
			for key, val := range shard.items {
				ch <- ChannelTuple{key, val}
			}
			shard.RUnlock()
		}
		close(ch)
	}()
	return ch
}

// IterBuffered returns a buffered iterator which could be used in a for range loop.
func (m ChannelMap) IterBuffered() <-chan ChannelTuple {
	ch := make(chan ChannelTuple, m.Count())
	go func() {
		// Foreach shard.
		for _, shard := range m {
			// Foreach key, value pair.
			shard.RLock()
			for key, val := range shard.items {
				ch <- ChannelTuple{key, val}
			}
			shard.RUnlock()
		}
		close(ch)
	}()
	return ch
}

// MarshalJSON reviles ChannelMap "private" variables to json marshal.
func (m ChannelMap) MarshalJSON() ([]byte, error) {
	// Create a temporary map, which will hold all item spread across shards.
	tmp := make(map[string]*Channel)

	// Insert items to temporary map.
	for item := range m.Iter() {
		tmp[item.Key] = item.Val
	}
	return json.Marshal(tmp)
}
