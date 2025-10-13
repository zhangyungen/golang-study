package tssync

import (
	"sync"
)

/**
 * 线程安全锁 压测发现是伪命题
 */
type miniMutex struct {
	lockMap map[string]*sync.Mutex
	rw      sync.RWMutex
}

var MiniMutex = &miniMutex{
	lockMap: make(map[string]*sync.Mutex),
}

func getLock(key string) *sync.Mutex {
	value, ok := MiniMutex.lockMap[key]
	if ok {
		return value
	}

	// 如果不存在，则使用写锁创建
	MiniMutex.rw.Lock()
	defer MiniMutex.rw.Unlock()
	// 双检查，避免在获取写锁的过程中已经被其他goroutine创建
	value, ok = MiniMutex.lockMap[key]
	if ok {
		return value
	}
	mutex := new(sync.Mutex)
	MiniMutex.lockMap[key] = mutex
	return mutex
}

func TryLock(key string) bool {
	return getLock(key).TryLock()
}

func Lock(key string) {
	getLock(key).Lock()
}

func UnLock(key string) {
	getLock(key).Unlock()
}

func RLock(key string) {
	getLock(key).Lock()
}

func RUnlock(key string) {
	getLock(key).Unlock()
}
