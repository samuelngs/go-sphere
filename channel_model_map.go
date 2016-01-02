package sphere

import (
	"encoding/json"
	"hash/fnv"
	"sync"
)

// ChannelModelShardCount nums of shard
var ChannelModelShardCount = 32

// ChannelModelMap is a "thread" safe map of type string:Anything.
// To avoid lock bottlenecks this map is dived to several (ChannelModelShardCount) map shards.
type ChannelModelMap []*ChannelModelMapShared

// ChannelModelMapShared is a "thread" safe string to anything map.
type ChannelModelMapShared struct {
	items        map[string]IChannels
	sync.RWMutex // Read Write mutex, guards access to internal map.
}

// ChannelModelTuple used by the Iter & IterBuffered functions to wrap two variables together over a channel,
type ChannelModelTuple struct {
	Key string
	Val IChannels
}

// NewChannelModelMap Creates a new concurrent map.
func NewChannelModelMap() ChannelModelMap {
	m := make(ChannelModelMap, ChannelModelShardCount)
	for i := 0; i < ChannelModelShardCount; i++ {
		m[i] = &ChannelModelMapShared{items: make(map[string]IChannels)}
	}
	return m
}

// GetShard returns shard under given key
func (m ChannelModelMap) GetShard(key string) *ChannelModelMapShared {
	hasher := fnv.New32()
	hasher.Write([]byte(key))
	return m[uint(hasher.Sum32())%uint(ChannelModelShardCount)]
}

// Set sets the given value under the specified key.
func (m *ChannelModelMap) Set(key string, value IChannels) {
	// Get map shard.
	shard := m.GetShard(key)
	shard.Lock()
	defer shard.Unlock()
	shard.items[key] = value
}

// SetIfAbsent sets the given value under the specified key if no value was associated with it.
func (m *ChannelModelMap) SetIfAbsent(key string, value IChannels) bool {
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
func (m ChannelModelMap) Get(key string) (IChannels, bool) {
	// Get shard
	shard := m.GetShard(key)
	shard.RLock()
	defer shard.RUnlock()

	// Get item from shard.
	val, ok := shard.items[key]
	return val, ok
}

// Count returns the number of elements within the map.
func (m ChannelModelMap) Count() int {
	count := 0
	for i := 0; i < ChannelModelShardCount; i++ {
		shard := m[i]
		shard.RLock()
		count += len(shard.items)
		shard.RUnlock()
	}
	return count
}

// Has looks up an item under specified key
func (m *ChannelModelMap) Has(key string) bool {
	// Get shard
	shard := m.GetShard(key)
	shard.RLock()
	defer shard.RUnlock()

	// See if element is within shard.
	_, ok := shard.items[key]
	return ok
}

// Remove removes an element from the map.
func (m *ChannelModelMap) Remove(key string) {
	// Try to get shard.
	shard := m.GetShard(key)
	shard.Lock()
	defer shard.Unlock()
	delete(shard.items, key)
}

// IsEmpty checks if map is empty.
func (m *ChannelModelMap) IsEmpty() bool {
	return m.Count() == 0
}

// Iter returns an iterator which could be used in a for range loop.
func (m ChannelModelMap) Iter() <-chan ChannelModelTuple {
	ch := make(chan ChannelModelTuple)
	go func() {
		// Foreach shard.
		for _, shard := range m {
			// Foreach key, value pair.
			shard.RLock()
			for key, val := range shard.items {
				ch <- ChannelModelTuple{key, val}
			}
			shard.RUnlock()
		}
		close(ch)
	}()
	return ch
}

// IterBuffered returns a buffered iterator which could be used in a for range loop.
func (m ChannelModelMap) IterBuffered() <-chan ChannelModelTuple {
	ch := make(chan ChannelModelTuple, m.Count())
	go func() {
		// Foreach shard.
		for _, shard := range m {
			// Foreach key, value pair.
			shard.RLock()
			for key, val := range shard.items {
				ch <- ChannelModelTuple{key, val}
			}
			shard.RUnlock()
		}
		close(ch)
	}()
	return ch
}

// MarshalJSON reviles ChannelModelMap "private" variables to json marshal.
func (m ChannelModelMap) MarshalJSON() ([]byte, error) {
	// Create a temporary map, which will hold all item spread across shards.
	tmp := make(map[string]IChannels)

	// Insert items to temporary map.
	for item := range m.Iter() {
		tmp[item.Key] = item.Val
	}
	return json.Marshal(tmp)
}
