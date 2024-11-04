package db

import (
	"container/list"
	"sync"
)

type Cache[K comparable, V any] interface {
	Get(key K) (V, bool)
	Put(key K, value V)
	Evict()
	Len() int
}

type Lru[K comparable, V any] struct {
	mu       sync.Mutex
	capacity int
	cache    map[K]*list.Element
	list     *list.List
}

type entry[K comparable, V any] struct {
	key   K
	value V
}

func NewCache[K comparable, V any](capacity int) Cache[K, V] {
	return &Lru[K, V]{
		capacity: capacity,
		cache:    make(map[K]*list.Element),
		list:     list.New(),
	}
}

func (c *Lru[K, V]) Get(key K) (V, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.cache[key]; ok {
		c.list.MoveToFront(elem)
		return elem.Value.(*entry[K, V]).value, true
	}

	var zero V
	return zero, false
}

func (c *Lru[K, V]) Put(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.cache[key]; ok {
		c.list.MoveToFront(elem)
		elem.Value.(*entry[K, V]).value = value
		return
	}

	if c.list.Len() == c.capacity {
		backElem := c.list.Back()
		c.list.Remove(backElem)
		delete(c.cache, backElem.Value.(*entry[K, V]).key)
	}

	newEntry := &entry[K, V]{key, value}
	elem := c.list.PushFront(newEntry)
	c.cache[key] = elem
}

func (lru *Lru[K, V]) Evict() {
	element := lru.list.Back()
	if element != nil {
		pair := element.Value.(*entry[K, V])
		delete(lru.cache, pair.key)
		lru.list.Remove(element)
	}
}

func (lru *Lru[K, V]) Len() int {
	return lru.list.Len()
}
