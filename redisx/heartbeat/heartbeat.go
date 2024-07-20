package redisx

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type Heartbeat struct {
	stop chan chan struct{}

	redisClient *redis.Client
	options     *Options
}

type HeartbeatError struct {
	RedisError error
	Count      int
}

type Options struct {
	Interval     time.Duration
	HeartbeatKey string
}

func NewHeartbeat(redisClient *redis.Client, opts ...Options) *Heartbeat {
	options := &Options{
		Interval:     time.Second,
		HeartbeatKey: "mq::connection::heartbeat",
	}
	if len(opts) > 0 {
		options = &opts[0]
	}
	b := &Heartbeat{
		redisClient: redisClient,
		stop:        make(chan chan struct{}),
		options:     options,
	}
	return b
}

// 启动心跳，本函数会阻塞，调用会应改开启协程来调用
func (b *Heartbeat) Start(errCallback func(err HeartbeatError)) {
	errCount := 0

	ticker := time.NewTicker(b.options.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			//continue
		case c := <-b.stop:
			close(c)
			return
		}

		err := b.hitHeartbeart()
		if err == nil {
			errCount = 0
			continue
		}

		errCount++
		if errCallback != nil {
			invokeFunc := func() {
				defer func() {
					if funcErr := recover(); funcErr != nil {
						fmt.Println("invoke heartbeart error callback occur panic", funcErr)
					}
				}()
				errCallback(HeartbeatError{
					RedisError: err,
					Count:      errCount,
				})
			}
			invokeFunc()
		}
	}
}

func (b *Heartbeat) Stop() error {
	if b.stop == nil {
		return nil
	}

	heartbeatStopped := make(chan struct{})
	b.stop <- heartbeatStopped
	<-heartbeatStopped
	b.stop = nil

	_, err := b.redisClient.Del(context.TODO(), b.options.HeartbeatKey).Result()
	if err != nil {
		return err
	}
	return nil
}

func (b *Heartbeat) hitHeartbeart() error {
	context := context.TODO()
	return b.redisClient.Set(context, b.options.HeartbeatKey, "ok", b.options.Interval).Err()
}
