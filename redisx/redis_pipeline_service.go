package redisx

type IRedisPipelineService interface {
}

// 默认的IRedisPipelineService实现
type RedisPipelineService struct {
	options *RedisOptions
}

func NewRedisPipelineService(options *RedisOptions) *RedisPipelineService {
	return &RedisPipelineService{
		options: options,
	}
}

func (s *RedisPipelineService) PipelineSetData(list map[string]interface{}, opts ...RedisValueOption) error {
	options := s.options.createRedisValueOptions()
	options.applyOption(opts...)

	//有效期
	ttl := Redis_NoExpiration_TTL
	if options.ttl != nil {
		ttl = *options.ttl
	}
	pipe := s.options.client.Pipeline()
	for key, v := range list {
		pipe.Set(options.ctx, options.appendKeyPrefix(key), v, ttl)
	}

	_, err := pipe.Exec(options.ctx)
	return err
}
