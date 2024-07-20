package redisx

import (
	"time"

	redis "github.com/go-redis/redis/v8"
)

// 对redis中的key进行操作的接口
type IRedisKeyService interface {
	GetRedisClient() *redis.Client

	//搜索符合pattern的key列表
	Keys(keyPattern string, opts ...RedisValueOption) ([]string, error)
	ExistKey(key string, opts ...RedisValueOption) (bool, error)
	DeleteKey(key string, opts ...RedisValueOption) error
	//删除指定前缀的所有key
	DeleteKeys(pattern string, opts ...RedisValueOption) error
	//设置key的过期时间
	KeyExpire(key string, duration time.Duration, opts ...RedisValueOption) error
	//获取key的过期时间
	KeyTimeToLive(key string, opts ...RedisValueOption) *time.Duration
}

var _ IRedisKeyService = (*RedisKeyService)(nil)

// 默认的IRedisKeyService实现
type RedisKeyService struct {
	options *RedisOptions
}

func NewRedisKeyService(options *RedisOptions) *RedisKeyService {
	return &RedisKeyService{
		options: options,
	}
}

func (s *RedisKeyService) GetRedisClient() *redis.Client {
	return s.options.client
}

// 搜索符合pattern的key列表
func (s *RedisKeyService) Keys(keyPattern string, opts ...RedisValueOption) ([]string, error) {
	options := s.options.createRedisValueOptions()
	options.applyOption(opts...)

	v := s.options.client.Keys(options.ctx, keyPattern)
	if v.Err() != nil {
		return make([]string, 0), v.Err()
	}
	return v.Val(), nil
}

func (s *RedisKeyService) ExistKey(key string, opts ...RedisValueOption) (bool, error) {
	options := s.options.createRedisValueOptions()
	options.applyOption(opts...)

	v := s.options.client.Exists(options.ctx, s.appendKeyPrefix(key, options))
	if v.Err() != nil {
		return false, v.Err()
	}
	return v.Val() > 0, nil
}

func (s *RedisKeyService) DeleteKey(key string, opts ...RedisValueOption) error {
	options := s.options.createRedisValueOptions()
	options.applyOption(opts...)

	return s.options.client.Del(options.ctx, s.appendKeyPrefix(key, options)).Err()
}

// 删除指定前缀的所有key
func (s *RedisKeyService) DeleteKeys(pattern string, opts ...RedisValueOption) error {
	if len(pattern) <= 0 {
		return nil
	}
	delWithPrefix := `
	local keys = redis.call('keys', ARGV[1]) 
                for i=1,#keys,5000 do 
                redis.call('del', unpack(keys, i, math.min(i+4999, #keys)))
                end
	`
	options := s.options.createRedisValueOptions()
	options.applyOption(opts...)

	cmd := s.options.client.Eval(options.ctx, delWithPrefix, []string{s.appendKeyPrefix(pattern, options)})
	return cmd.Err()
}

// 设置key的过期时间
func (s *RedisKeyService) KeyExpire(key string, duration time.Duration, opts ...RedisValueOption) error {
	options := s.options.createRedisValueOptions()
	options.applyOption(opts...)

	return s.options.client.Expire(options.ctx, s.appendKeyPrefix(key, options), duration).Err()
}

// 获取key的过期时间
func (s *RedisKeyService) KeyTimeToLive(key string, opts ...RedisValueOption) *time.Duration {
	options := s.options.createRedisValueOptions()
	options.applyOption(opts...)

	cmd := s.options.client.TTL(options.ctx, s.appendKeyPrefix(key, options))
	ttl, err := cmd.Result()
	if err != nil {
		return nil
	}
	return &ttl
}

func (s *RedisKeyService) appendKeyPrefix(key string, options *RedisValueOptions) string {
	if !options.withoutPrefixKey {
		return options.appendKeyPrefix(key)
	}
	return key
}
