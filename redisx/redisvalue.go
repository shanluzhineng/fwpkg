package redisx

import (
	"context"
	"time"
)

// redis值
type IRedisValue interface {
	//获取值是否存在
	Exist() bool
	Bytes() []byte
	Err() error

	//转换为目标值
	ToValue(val interface{}) error

	ValToString() string
	ValToInt() (int, error)
	ValToInt32() (int32, error)
	ValToInt64() (int64, error)
	ValToBool() (bool, error)
	ValToTime() (time.Time, error)
}

type RedisValueOption func(*RedisValueOptions)

type RedisValueOptions struct {
	//ttl,如果不指定则不超时
	ttl *time.Duration
	//是否不处理prefix
	withoutPrefixKey bool
	ctx              context.Context
	keyPrefix        string
	//在get时如果为空，是否自动load
	loadIfEmpty bool

	unmarshal UnmarshalFunc
	marshal   MarshalFunc
}

// 确保redis key包含了指定的前缀
func (o *RedisValueOptions) appendKeyPrefix(key string) string {
	if len(o.keyPrefix) <= 0 {
		return key
	}
	return ensureStartWith(key, o.keyPrefix)
}

func newRedisValueOptions() *RedisValueOptions {
	return &RedisValueOptions{
		withoutPrefixKey: false,
		loadIfEmpty:      true,
	}
}

func (o *RedisValueOptions) applyOption(opts ...RedisValueOption) {
	for _, eachOpt := range opts {
		eachOpt(o)
	}
}

// 值带上ttl
func WithTTL(ttl time.Duration) RedisValueOption {
	return func(rvo *RedisValueOptions) {
		rvo.ttl = &ttl
	}
}

func WithExpiredTime(expiredTime time.Time) RedisValueOption {
	return func(rvo *RedisValueOptions) {
		d := time.Until(expiredTime)
		if d < 0 {
			return
		}
		rvo.ttl = &d
	}
}

// key不带上prefix
func WithoutPrefixKey() RedisValueOption {
	return func(rvo *RedisValueOptions) {
		rvo.withoutPrefixKey = true
	}
}

// 不自动load
func SuppressedLoadIfEmpty() RedisValueOption {
	return func(rvo *RedisValueOptions) {
		rvo.loadIfEmpty = false
	}
}

// 指定key前缀
func WithKeyPrefix(keyPrefix string) RedisValueOption {
	return func(rvo *RedisValueOptions) {
		rvo.keyPrefix = keyPrefix
	}
}

func WithContext(ctx context.Context) RedisValueOption {
	return func(rvo *RedisValueOptions) {
		rvo.ctx = ctx
	}
}

func WithUnmarshalFunc(unmarshal UnmarshalFunc) RedisValueOption {
	return func(rvo *RedisValueOptions) {
		rvo.unmarshal = unmarshal
	}
}

func WithMarshalFunc(marshal MarshalFunc) RedisValueOption {
	return func(rvo *RedisValueOptions) {
		rvo.marshal = marshal
	}
}
