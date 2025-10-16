package main
import (
	"fmt"
	"sync"
	"time"
)

// ### 任务5：内存缓存系统
// **目标**：掌握map、mutex、接口设计
// **描述**：实现一个线程安全的内存缓存系统，支持TTL和LRU淘汰策略

// **流程提示**：
// 1. 设计缓存接口（Get, Set, Delete, Clear）
// 2. 实现基于map的存储结构
// 3. 添加读写锁保证线程安全
// 4. 实现TTL过期机制
// 5. 实现LRU淘汰算法
// 6. 添加统计功能（命中率、大小等）

type CacheConfig struct {
	Size     int
	TTL      time.Duration
	LRU      bool
	LFU      bool
	SizeCount int
}

type CacheItem struct {
	Value     interface{}
	ExpireAt  time.Time
	HitCount  int
}
type Cache struct {
	config CacheConfig
	mu     sync.RWMutex
	items  map[string]*CacheItem
}
func (c *Cache) Keys() []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	keys := make([]string, 0, len(c.items))
	for key := range c.items {
		keys = append(keys, key)
	}
	return keys
}

func (c *Cache) Size() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.config.SizeCount
}

func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]*CacheItem)
	c.config.SizeCount = 0
}

func InitCache(config CacheConfig) *Cache {
	return &Cache{
		config: config,
		items:  make(map[string]*CacheItem),
	}
}
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	item, found := c.items[key]
	if !found {
		return nil, false
	}
	if item.ExpireAt.Before(time.Now()) {
		go delete(c.items, key)
		return nil, false
	}
	item.HitCount++
	return item.Value, true
}
func (c *Cache ) Set(Key string,value interface{},ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.config.SizeCount + 1 > c.config.Size {
		var miniKey string
			var minimumHitCount int
			var olderTime time.Time
		if c.config.LFU {
			for key, item := range c.items {
				if item.HitCount <= minimumHitCount  {
					miniKey = key
					minimumHitCount = item.HitCount
					if minimumHitCount == item.HitCount {
						if item.ExpireAt.Before(olderTime) {
							miniKey = key
							olderTime = item.ExpireAt
						}
					}
				}
			}
		
	}else {
		for key, item := range c.items {
			if item.ExpireAt.Before(olderTime) {
					miniKey = key
					olderTime = item.ExpireAt
			}
		}
				
	}
		delete(c.items, miniKey)
}
	
	if _, found := c.items[Key]; found {
		c.config.SizeCount--
	}
	c.items[Key] = &CacheItem{
		Value:     value,
		ExpireAt:  time.Now().Add(ttl),
		HitCount:  0,
	}
	c.config.SizeCount++
 
}

func main8() {
	fmt.Println("Go 语言学习任务")
	fmt.Println("==================")

	newCache := InitCache(CacheConfig{
		Size:     10,
		TTL:      time.Minute,
		LRU:      true,
		LFU:      false,
	})
	newCache.Set("key1", "value1", time.Minute)
	newCache.Set("key2", "value2", time.Minute)
	newCache.Set("key3", "value3", time.Minute)
	newCache.Set("key4", "value4", time.Minute)
	newCache.Set("key5", "value5", time.Minute)
	newCache.Set("key6", "value6", time.Minute)
	newCache.Set("key7", "value7", time.Minute)
	newCache.Set("key8", "value8", time.Minute)
	newCache.Set("key9", "value9", time.Minute)
	newCache.Set("key10", "value10", time.Minute)
	fmt.Println(newCache.Get("key1"))
	fmt.Println(newCache.items)
	time.Sleep(time.Minute)
	fmt.Println(newCache.items)
	newCache.Set("key11", "value11", time.Minute)
	fmt.Println(newCache.items)
	
}