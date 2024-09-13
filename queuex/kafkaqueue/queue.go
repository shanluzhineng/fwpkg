package kafkaqueue

import (
	"context"
	"fmt"
	"io"
	"time"

	kafka "github.com/segmentio/kafka-go"
	"github.com/shanluzhineng/fwpkg/queuex"
	"github.com/shanluzhineng/fwpkg/system/log"
	"github.com/shanluzhineng/threadingx/service"
	"github.com/shanluzhineng/threadingx/threading"
)

const (
	defaultCommitInterval = time.Second
	defaultMaxWait        = time.Second
	defaultQueueCapacity  = 1000
)

type (
	ConsumeHandle func(key, value string) error

	ConsumeHandler interface {
		Consume(key, value string) error
	}

	queueOptions struct {
		commitInterval time.Duration
		queueCapacity  int
		maxWait        time.Duration
	}

	QueueOption func(*queueOptions)

	kafkaQueue struct {
		c                KqConf
		consumer         *kafka.Reader
		handler          ConsumeHandler
		channel          chan kafka.Message
		producerRoutines *threading.RoutineGroup
		consumerRoutines *threading.RoutineGroup
	}

	kafkaQueues struct {
		queues []queuex.MessageQueue
		group  *service.ServiceGroup
	}
)

func MustNewQueue(c KqConf, handler ConsumeHandler, opts ...QueueOption) queuex.MessageQueue {
	q, _ := NewQueue(c, handler, opts...)
	return q
}

func NewQueue(c KqConf, handler ConsumeHandler, opts ...QueueOption) (queuex.MessageQueue, error) {
	var options queueOptions
	for _, opt := range opts {
		opt(&options)
	}
	ensureQueueOptions(c, &options)

	if c.Conns < 1 {
		c.Conns = 1
	}
	q := kafkaQueues{
		group: service.NewServiceGroup(),
	}
	for i := 0; i < c.Conns; i++ {
		q.queues = append(q.queues, newKafkaQueue(c, handler, options))
	}

	return q, nil
}

func newKafkaQueue(c KqConf, handler ConsumeHandler, options queueOptions) queuex.MessageQueue {
	var offset int64
	if c.Offset == firstOffset {
		offset = kafka.FirstOffset
	} else {
		offset = kafka.LastOffset
	}
	consumer := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        c.Brokers,
		GroupID:        c.Group,
		Topic:          c.Topic,
		StartOffset:    offset,
		MinBytes:       c.MinBytes,
		MaxBytes:       c.MaxBytes,
		MaxWait:        options.maxWait,
		CommitInterval: options.commitInterval,
		QueueCapacity:  options.queueCapacity,
	})

	return &kafkaQueue{
		c:                c,
		consumer:         consumer,
		handler:          handler,
		channel:          make(chan kafka.Message),
		producerRoutines: threading.NewRoutineGroup(),
		consumerRoutines: threading.NewRoutineGroup(),
	}
}

// 启动服务
func (q *kafkaQueue) Start() {
	q.startConsumers()
	q.startProducers()

	q.producerRoutines.Wait()
	close(q.channel)
	q.consumerRoutines.Wait()
}

func (q *kafkaQueue) Stop() {
	q.consumer.Close()
}

func (q *kafkaQueue) consumeOne(key, val string) error {
	err := q.handler.Consume(key, val)
	return err
}

func (q *kafkaQueue) startConsumers() {
	for i := 0; i < q.c.Processors; i++ {
		q.consumerRoutines.Run(func() {
			for msg := range q.channel {
				if err := q.consumeOne(string(msg.Key), string(msg.Value)); err != nil {
					log.Logger.Error(fmt.Sprintf("Error on consuming: %s, error: %v", string(msg.Value), err))
				}
				q.consumer.CommitMessages(context.Background(), msg)
			}
		})
	}
}

func (q *kafkaQueue) startProducers() {
	for i := 0; i < q.c.Consumers; i++ {
		q.producerRoutines.Run(func() {
			for {
				msg, err := q.consumer.FetchMessage(context.Background())
				//io.EOF means consumer closed
				//io.ErrClosedPipe means committing messages on the consumer,
				//kafka will refire the message on uncommitted messages,ignore
				if err == io.EOF || err == io.ErrClosedPipe {
					return
				}
				if err != nil {
					log.Logger.Error(fmt.Sprintf("Error on reading message, %q", err.Error()))
					continue
				}
				q.channel <- msg
			}
		})
	}
}

func (q kafkaQueues) Start() {
	for _, each := range q.queues {
		q.group.Add(each)
	}
	q.group.Start()
}

func (q kafkaQueues) Stop() {
	q.group.Stop()
}

// 设置commitInterval参数
func WithCommitInterval(interval time.Duration) QueueOption {
	return func(options *queueOptions) {
		options.commitInterval = interval
	}
}

func WithQueueCapacity(queueCapacity int) QueueOption {
	return func(options *queueOptions) {
		options.queueCapacity = queueCapacity
	}
}

func WithHandle(handle ConsumeHandle) ConsumeHandler {
	return innerConsumeHandler{
		handle: handle,
	}
}

func WithMaxWait(wait time.Duration) QueueOption {
	return func(options *queueOptions) {
		options.maxWait = wait
	}
}

type innerConsumeHandler struct {
	handle ConsumeHandle
}

func (ch innerConsumeHandler) Consume(k, v string) error {
	return ch.handle(k, v)
}

func ensureQueueOptions(c KqConf, options *queueOptions) {
	if options.commitInterval == 0 {
		options.commitInterval = defaultCommitInterval
	}
	if options.queueCapacity == 0 {
		options.queueCapacity = defaultQueueCapacity
	}
	if options.maxWait == 0 {
		options.maxWait = defaultMaxWait
	}
}
