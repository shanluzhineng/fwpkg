package queue

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-redis/redis/v8"
)

type Queue struct {
	options *queueOptions
}

type queueOptions struct {
	client *redis.Client
	//queue name
	queueName string
	//key
	queueKey string
}

func defaultQueueOptions() *queueOptions {
	return &queueOptions{
		queueKey: "mq::queues",
	}
}

type QueueOption func(q *Queue)

func WithRedisClient(client *redis.Client) QueueOption {
	return func(q *Queue) {
		q.options.client = client
	}
}

func WithRedisConn(host, port, password string, redisDb int, connOpts ...func(*redis.Options)) QueueOption {
	return func(q *Queue) {
		options := &redis.Options{
			Addr:     fmt.Sprintf("%s:%s", host, port),
			Password: password,
			DB:       redisDb,
		}
		redisClient := redis.NewClient(options)
		for _, eachOpt := range connOpts {
			eachOpt(options)
		}
		q.options.client = redisClient
	}
}

func WithQueueName(queueName string) QueueOption {
	return func(q *Queue) {
		q.options.queueName = queueName
	}
}

// new a queue
func NewQueue(opts ...QueueOption) (*Queue, error) {
	queue := &Queue{
		options: defaultQueueOptions(),
	}
	for _, eachOpt := range opts {
		eachOpt(queue)
	}
	err := setupQueue(queue)
	if err != nil {
		return nil, err
	}
	return queue, nil
}

func setupQueue(q *Queue) error {
	if q.options.client == nil {
		return errors.New("必须指定redis连接参数")
	}
	if len(q.options.queueName) <= 0 {
		return errors.New("必须指定queueName参数值")
	}
	context := context.TODO()
	err := q.options.client.SAdd(context, q.options.queueKey, q.options.queueName).Err()
	if err != nil {
		return err
	}
	return nil
}
