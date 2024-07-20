package redisx

import "github.com/go-redis/redis/v8"

type IRedisStringService interface {
	//获取key值
	StringGet(key string, opts ...RedisValueOption) IRedisValue
	//获取多个key值
	StringMGet(opts []RedisValueOption, keys ...string) (map[string]IRedisValue, error)
	//设置值
	StringSet(key string, value interface{}, opts ...RedisValueOption) error

	//自增id
	KeyIncr(key string, opts ...RedisValueOption) (int64, error)
	KeyIncrBy(key string, value int64, opts ...RedisValueOption) (int64, error)
	GetIncr(key string, opts ...RedisValueOption) (int64, error)
	//自减id
	KeyDecr(key string, opts ...RedisValueOption) (int64, error)
	KeyDecrBy(key string, decrement int64, opts ...RedisValueOption) (int64, error)
	GetDecr(key string, opts ...RedisValueOption) (int64, error)
}

var _ IRedisStringService = (*RedisStringService)(nil)

type RedisStringService struct {
	*RedisKeyService
}

func NewRedisStringService(options *RedisOptions) IRedisStringService {
	s := &RedisStringService{
		RedisKeyService: NewRedisKeyService(options),
	}
	return s
}

func (s *RedisStringService) StringGet(key string, opts ...RedisValueOption) IRedisValue {
	options := s.options.createRedisValueOptions()
	options.applyOption(opts...)

	b, err := s.options.client.Get(options.ctx, options.appendKeyPrefix(key)).Bytes()
	if err != nil {
		if err == redis.Nil {
			return newNilRedisValue()
		}
		return newErrRedisValue(err)
	}
	return newRedisValue(b, s.options.Unmarshal)
}

// 获取多个key值
func (s *RedisStringService) StringMGet(opts []RedisValueOption, keys ...string) (map[string]IRedisValue, error) {
	options := s.options.createRedisValueOptions()
	options.applyOption(opts...)

	result := make(map[string]IRedisValue)
	normalizedKeys := make([]string, 0)
	for _, eachKey := range keys {
		currentNormalizedKey := options.appendKeyPrefix(eachKey)
		normalizedKeys = append(normalizedKeys, currentNormalizedKey)
		result[eachKey] = newNilRedisValue()
	}
	b := s.options.client.MGet(options.ctx, normalizedKeys...)
	if err := b.Err(); err != nil {
		if err == redis.Nil {
			return result, nil
		}
		return nil, err
	}
	for i, eachKey := range keys {
		val := b.Val()[i]
		if val == nil {
			result[eachKey] = newNilRedisValue()
			continue
		}
		currentRedisValue, ok := val.(string)
		if !ok {
			result[eachKey] = newNilRedisValue()
			continue
		}
		result[eachKey] = newRedisValue([]byte(currentRedisValue), s.options.Unmarshal)
	}
	return result, nil
}

func (s *RedisStringService) StringSet(key string, value interface{}, opts ...RedisValueOption) error {
	options := s.options.createRedisValueOptions()
	options.applyOption(opts...)

	data, err := options.marshal(value)
	if err != nil {
		return err
	}

	//有效期
	ttl := Redis_NoExpiration_TTL
	if options.ttl != nil {
		ttl = *options.ttl
	}
	return s.options.client.Set(options.ctx, options.appendKeyPrefix(key), data, ttl).Err()
}

func (s *RedisStringService) KeyIncr(key string, opts ...RedisValueOption) (int64, error) {
	options := s.options.createRedisValueOptions()
	options.applyOption(opts...)

	result := s.options.client.Incr(options.ctx, options.appendKeyPrefix(key))
	err := result.Err()
	if err != nil {
		return -1, err
	}
	return result.Val(), nil
}

func (s *RedisStringService) KeyIncrBy(key string, value int64, opts ...RedisValueOption) (int64, error) {
	options := s.options.createRedisValueOptions()
	options.applyOption(opts...)

	result := s.options.client.IncrBy(options.ctx, options.appendKeyPrefix(key), value)
	err := result.Err()
	if err != nil {
		return -1, err
	}
	return result.Val(), nil
}

func (s *RedisStringService) GetIncr(key string, opts ...RedisValueOption) (int64, error) {
	return s._keyInt64Value(key, opts...)
}

func (s *RedisStringService) KeyDecr(key string, opts ...RedisValueOption) (int64, error) {
	options := s.options.createRedisValueOptions()
	options.applyOption(opts...)

	result := s.options.client.Decr(options.ctx, options.appendKeyPrefix(key))
	err := result.Err()
	if err != nil {
		return -1, err
	}
	return result.Val(), nil
}

func (s *RedisStringService) KeyDecrBy(key string, decrement int64, opts ...RedisValueOption) (int64, error) {
	options := s.options.createRedisValueOptions()
	options.applyOption(opts...)

	result := s.options.client.DecrBy(options.ctx, options.appendKeyPrefix(key), decrement)
	err := result.Err()
	if err != nil {
		return -1, err
	}
	return result.Val(), nil
}

func (s *RedisStringService) GetDecr(key string, opts ...RedisValueOption) (int64, error) {
	return s._keyInt64Value(key, opts...)
}

func (s *RedisStringService) _keyInt64Value(key string, opts ...RedisValueOption) (int64, error) {
	options := s.options.createRedisValueOptions()
	options.applyOption(opts...)

	value := s.StringGet(key, opts...)
	return value.ValToInt64()
}
