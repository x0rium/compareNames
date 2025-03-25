package matcher

import (
	"sort"
	"strings"
	"time"
)

// NewCache создает новый экземпляр кэша
func NewCache(maxSize int, ttl time.Duration) *Cache {
	return &Cache{
		items:   make(map[string]CacheItem),
		keys:    make([]string, 0, maxSize),
		maxSize: maxSize,
		TTL:     ttl,
	}
}

// Get получает элемент из кэша
func (c *Cache) Get(key string) (MatchResult, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, ok := c.items[key]
	if !ok {
		return MatchResult{}, false
	}

	// Проверяем, не устарел ли элемент
	if time.Since(item.CreateTime) > c.TTL {
		return MatchResult{}, false
	}

	return item.Result, true
}

// Put добавляет элемент в кэш
func (c *Cache) Put(key string, result MatchResult) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Если кэш достиг максимального размера, удаляем самый старый элемент (LRU)
	if len(c.items) >= c.maxSize {
		if len(c.keys) > 0 {
			oldestKey := c.keys[0]
			delete(c.items, oldestKey)
			c.keys = c.keys[1:]
		}
	}

	// Добавляем новый элемент
	c.items[key] = CacheItem{
		Result:     result,
		CreateTime: time.Now(),
	}
	c.keys = append(c.keys, key)
}

// Update обновляет LRU кэш (перемещает ключ в конец списка)
func (c *Cache) Update(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	// Ищем ключ в списке
	for i, k := range c.keys {
		if k == key {
			// Удаляем ключ из текущей позиции
			c.keys = append(c.keys[:i], c.keys[i+1:]...)
			// Добавляем ключ в конец списка
			c.keys = append(c.keys, key)
			break
		}
	}
}

// Генерирует ключ для кэша результатов сравнения
func (m *NameMatcher) cacheKey(name1, name2 string, attrs Attributes) string {
	// Сортируем имена для обеспечения одинакового ключа независимо от порядка аргументов
	names := []string{strings.ToLower(name1), strings.ToLower(name2)}
	sort.Strings(names)

	// Формируем базовый ключ из имен
	key := names[0] + "||" + names[1] + "||"

	// Если есть дополнительные атрибуты, добавляем их к ключу
	if attrs != nil {
		// Сортируем атрибуты по имени для стабильности ключа
		attrNames := make([]string, 0, len(attrs))
		for name := range attrs {
			attrNames = append(attrNames, name)
		}
		sort.Strings(attrNames)

		// Добавляем атрибуты к ключу
		for _, name := range attrNames {
			if attrs[name].Match {
				key += name + ":true||"
			} else {
				key += name + ":false||"
			}
		}
	}

	return key
}
