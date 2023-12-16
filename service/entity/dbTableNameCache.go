package entity

import "sync"

type DbTableNameCache struct {
	cache map[string]bool
	lock  sync.RWMutex
}

func HasTableCache(value string) bool {
	dbTableNameCache.lock.RLock()
	defer dbTableNameCache.lock.RUnlock()

	_, ok := dbTableNameCache.cache[value]
	return ok
}

func SetHasTableCache(value string) {
	dbTableNameCache.lock.Lock()
	defer dbTableNameCache.lock.Unlock()

	if dbTableNameCache.cache == nil {
		dbTableNameCache.cache = make(map[string]bool)
	}

	dbTableNameCache.cache[value] = true
}

func ClearHasTableCache() {
	dbTableNameCache.lock.Lock()
	defer dbTableNameCache.lock.Unlock()

	dbTableNameCache.cache = make(map[string]bool)
}
