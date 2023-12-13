package cache

import (
	"errors"
	model "nats-service/app/model"
	"sync"
)

// var cache Cache

// type Item struct {
// 	Value      interface{}
// 	Created    time.Time
// 	Expiration int64
// }

type Cache struct {
	mu sync.RWMutex
	sync.RWMutex
	items map[int]model.Order
}

func New() *Cache {
	// if len(cache.items) != 0 {
	// 	return &cache
	// }
	items := make(map[int]model.Order)
	cache := Cache{
		items: items,
	}

	// if cleanupInterval > 0 {
	// 	cache.StartGC() // данный метод рассматривается ниже
	// }

	return &cache
}

func (c *Cache) Copy() *Cache {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Создаем новый кэш
	copiedCache := &Cache{
		items: make(map[int]model.Order, len(c.items)),
	}

	// Копируем элементы из оригинального кэша
	for key, value := range c.items {
		copiedCache.items[key] = value
	}

	return copiedCache
}

func (c *Cache) Set(key int, value model.Order) {

	c.Lock()

	defer c.Unlock()

	c.items[key] = value
}

func (c *Cache) Len() int {
	return len(c.items)
}

func (c *Cache) Get(key int) (*model.Order, bool) {
	c.RLock()
	defer c.RUnlock()

	item, found := c.items[key]

	if !found {
		return nil, false
	}

	return &item, true
}

func (c *Cache) Delete(key int) error {

	c.Lock()

	defer c.Unlock()

	if _, found := c.items[key]; !found {
		return errors.New("Key not found")
	}

	delete(c.items, key)

	return nil
}

// func (c *Cache) StartGC() {
// 	go c.GC()
// }

// func (c *Cache) GC() {

// 	for {
// 		// ожидаем время установленное в cleanupInterval
// 		<-time.After(c.cleanupInterval)

// 		if c.items == nil {
// 			return
// 		}

// 		// Ищем элементы с истекшим временем жизни и удаляем из хранилища
// 		if keys := c.expiredKeys(); len(keys) != 0 {
// 			c.clearItems(keys)

// 		}
// 	}
// }

// func (c *Cache) expiredKeys() (keys []int) {

// 	c.RLock()

// 	defer c.RUnlock()

// 	for k, i := range c.items {
// 		if time.Now().UnixNano() > i.Expiration && i.Expiration > 0 {
// 			keys = append(keys, k)
// 		}
// 	}

// 	return
// }

func (c *Cache) clearItems(keys []int) {

	c.Lock()

	defer c.Unlock()

	for _, k := range keys {
		delete(c.items, k)
	}
}
