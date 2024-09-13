package queuex

type (
	// A Consumer interface represents a consumer that can consume string messages.
	Consumer interface {
		//消费指定key的队列
		Consume(string) error
		OnEvent(event interface{})
	}

	// ConsumerFactory defines the factory to generate consumers.
	ConsumerFactory func() (Consumer, error)
)
