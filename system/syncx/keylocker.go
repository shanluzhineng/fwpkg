package syncx

import (
	"sync"
)

// KeyLock 用于根据 key 来加锁和解锁
type KeyLock struct {
	lockMap    map[string]*sync.Mutex
	globalLock sync.Mutex // 用于保护 lockMap
}

var keyLocker *KeyLock

// GetSingleKeyLock 获取 KeyLock 单例（不存在则创建）
func GetSingleKeyLock() *KeyLock {
	if keyLocker == nil {
		keyLocker = &KeyLock{
			lockMap: make(map[string]*sync.Mutex),
		}
	}
	return keyLocker
}

// Lock 根据 key 加锁
func (kl *KeyLock) Lock(key string) {
	kl.globalLock.Lock()
	defer kl.globalLock.Unlock()
	if _, ok := kl.lockMap[key]; !ok {
		kl.lockMap[key] = &sync.Mutex{}
	}
	kl.lockMap[key].Lock()
}

// Unlock 根据 key 解锁
func (kl *KeyLock) Unlock(key string) {
	kl.globalLock.Lock()
	defer kl.globalLock.Unlock()
	if lock, ok := kl.lockMap[key]; ok {
		lock.Unlock()
		// 可以选择在这里删除锁，或者在其他地方进行清理
		delete(kl.lockMap, key)
	}
}
