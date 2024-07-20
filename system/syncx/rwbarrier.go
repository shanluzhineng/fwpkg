package syncx

import (
	"sync"
)

func GuardRLock(lock *sync.RWMutex, fn func()) {
	lock.RLock()
	defer lock.RUnlock()
	fn()
}

func GuardRLockAs[T any](lock *sync.RWMutex, fn func() T) T {
	lock.RLock()
	defer lock.RUnlock()
	return fn()
}

func GuardLock(lock *sync.RWMutex, fn func()) {
	lock.Lock()
	defer lock.Unlock()
	fn()
}

func GuardLockAs[T any](lock *sync.RWMutex, fn func() T) T {
	lock.Lock()
	defer lock.Unlock()
	return fn()
}
