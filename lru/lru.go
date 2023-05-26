// Package lru implement a LRU datastructure
// LRU means Least Frequently Used
package lru

import "container/list"

type entry struct {
	key   string
	value Value
}

type Cache struct {
	// Max alistowded memory
	maxBytes int64

	// Already used memory
	nbytes int64

	// Data type of list is entry, then we can get key in O(1) when updating
	list *list.List

	// Map fron key to listnode
	cache map[string]*list.Element

	// Callback when one data is evicted
	OnEvicted func(key string, value Value)
}

// New is the Constructor of Cache
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		list:      list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// Get returns the value of key if exists, and update the frequency of key
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.list.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// Add adds a key to the cache if doesn`t exist, else modify key`s value
func (c *Cache) Add(key string, value Value) {
	if ele, ok := c.cache[key]; ok {
		c.list.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		ele := c.list.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// Len the number of cache entries
func (c *Cache) Len() int {
	return c.list.Len()
}

// RemoveOldest removes the oldest item
func (c *Cache) RemoveOldest() {
	ele := c.list.Back()
	if ele != nil {
		c.list.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

type Value interface {
	Len() int
}
