package redisx

import "github.com/go-redis/redis/v8"

type IRedisHashService interface {
	HashGet(key string, field string, opts ...RedisValueOption) IRedisValue
	HashGetAll(key string, opts ...RedisValueOption) (RedisValueMap, error)

	HashSet(key string, values map[string]interface{}, opts ...RedisValueOption) error
}

type RedisHashService struct {
	*RedisKeyService
}

func NewRedisHashService(options *RedisOptions) IRedisHashService {
	s := &RedisHashService{
		RedisKeyService: NewRedisKeyService(options),
	}
	return s
}

// get field from hash
func (s *RedisHashService) HashGet(key string, field string, opts ...RedisValueOption) IRedisValue {
	options := s.options.createRedisValueOptions()
	options.applyOption(opts...)

	b, err := s.options.client.HGet(options.ctx, options.appendKeyPrefix(key), field).Bytes()
	if err != nil {
		if err == redis.Nil {
			return newNilRedisValue()
		}
		return newErrRedisValue(err)
	}
	return newRedisValue(b, s.options.Unmarshal)
}

// get all value from hash
func (s *RedisHashService) HashGetAll(key string, opts ...RedisValueOption) (RedisValueMap, error) {
	options := s.options.createRedisValueOptions()
	options.applyOption(opts...)

	result := RedisValueMap{}
	b := s.options.client.HGetAll(options.ctx, options.appendKeyPrefix(key))
	if err := b.Err(); err != nil {
		if err == redis.Nil {
			return result, nil
		}
		return nil, err
	}

	valueList := b.Val()
	if len(valueList) <= 0 {
		return result, nil
	}

	for eachKey, eachValue := range valueList {
		result[eachKey] = newRedisValue([]byte(eachValue), s.options.Unmarshal)
	}
	return result, nil
}

func (s *RedisHashService) HashSet(key string, values map[string]interface{}, opts ...RedisValueOption) error {
	options := s.options.createRedisValueOptions()
	options.applyOption()

	data := make(map[string]interface{})
	for eachKey, eachValue := range values {
		currentSValue, err := s.options.Marshal(eachValue)
		if err != nil {
			return err
		}
		data[eachKey] = string(currentSValue)
	}
	return s.options.client.HSet(options.ctx, options.appendKeyPrefix(key), data).Err()
}
