package redisx

import (
	"context"
	"fmt"
	"strconv"
	"sync/atomic"

	"github.com/go-redis/redis/v8"
)

const (
	RedisRandLen    = 16
	tolerance       = 500 // milliseconds
	millisPerSecond = 1000
	lockCommand     = `if redis.call("GET", KEYS[1]) == ARGV[1] then
    redis.call("SET", KEYS[1], ARGV[1], "PX", ARGV[2])
    return "OK"
else
    return redis.call("SET", KEYS[1], ARGV[1], "NX", "PX", ARGV[2])
end`
	delCommand = `local val = redis.call("GET", KEYS[1])
if val then
	if val == ARGV[1] then
		return redis.call("DEL", KEYS[1])
	else
		return 0
	end
else
	return -1
end`
)

// A RedisLock is a redis lock.
type RedisLock struct {
	Store   *redis.Client
	Seconds uint32
	Key     string
	Value   string
}

// NewRedisLock returns a RedisLock.
func NewRedisLock(store *redis.Client, key string, value string, expire uint32) *RedisLock {
	return &RedisLock{
		Store:   store,
		Key:     key,
		Value:   value,
		Seconds: expire,
	}
}

// RedisAcquire Lua script方式加锁
func (rl *RedisLock) RedisAcquire() (bool, error) {
	seconds := atomic.LoadUint32(&rl.Seconds)
	resp, err := rl.Store.Eval(context.Background(), lockCommand, []string{rl.Key}, []string{
		rl.Value, strconv.Itoa(int(seconds)*millisPerSecond + tolerance),
	}).Result()
	if err == redis.Nil {
		return false, err
	} else if err != nil {
		return false, err
	} else if resp == nil {
		return false, nil
	}

	reply, ok := resp.(string)
	if ok && reply == "OK" {
		return true, nil
	}

	return false, nil
}

// RedisRelease releases the lock.
// @ return
// @   int64: Released number
// @	 0 - exist key but given a wrong token
// @	 -1 - key is not exist
// @	 -2 - other errors
func (rl *RedisLock) RedisRelease() (int64, error) {
	v, err := rl.Store.Eval(context.Background(), delCommand, []string{rl.Key}, []string{rl.Value}).Int64()
	if err != nil {
		return -2, err
	}
	return v, nil
}

// SetExpire sets the expiration.
func (rl *RedisLock) SetExpire(seconds int) {
	atomic.StoreUint32(&rl.Seconds, uint32(seconds))
}

func (rl *RedisLock) ToString() string {
	return fmt.Sprintf("key=%s, v=%s, seconds=%d", rl.Key, rl.Value, rl.Seconds)
}
