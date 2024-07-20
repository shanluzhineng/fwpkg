package redisx

import (
	"context"
	"time"

	redis "github.com/go-redis/redis/v8"
)

const (
	//redis过期时间(30分钟)
	Redis_30M_TTL time.Duration = time.Duration(30) * time.Minute
	//redis过期时间(24小时)
	Redis_24Hour_TTL time.Duration = time.Duration(24) * time.Hour
	//1个月的有效期
	Redis_Expiration_Month_TTL time.Duration = time.Duration(30) * time.Hour * 24
	//无过期时间
	Redis_NoExpiration_TTL time.Duration = time.Duration(-1)
)

type (
	UnmarshalFunc func([]byte, interface{}) error
	MarshalFunc   func(interface{}) ([]byte, error)
)

// 模块参数配置
type RedisOptions struct {
	client    *redis.Client
	KeyPrefix string
	//默认的有效期时间,如果不设，则表示没有有效期
	DefaultTTL *time.Duration

	Unmarshal UnmarshalFunc
	Marshal   MarshalFunc
}

// 创建默认的配置项
func NewRedisOptions(client *redis.Client) *RedisOptions {
	noTTL := Redis_NoExpiration_TTL
	return &RedisOptions{
		client:     client,
		DefaultTTL: &noTTL,
		Unmarshal:  _unmarshal,
		Marshal:    _marshal,
	}
}

func (o *RedisOptions) createRedisValueOptions() *RedisValueOptions {
	valueOptions := newRedisValueOptions()
	valueOptions.ttl = o.DefaultTTL
	valueOptions.marshal = o.Marshal
	valueOptions.unmarshal = o.Unmarshal
	valueOptions.ctx = context.TODO()
	valueOptions.keyPrefix = o.KeyPrefix
	return valueOptions
}
