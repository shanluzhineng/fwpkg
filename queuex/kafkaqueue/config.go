package kafkaqueue

const (
	firstOffset = "first"
	lastOffset  = "last"
)

type KqConf struct {
	Brokers    []string `json:"brokers"`
	Group      string   `json:"group"`
	Topic      string   `json:"topic"`
	Offset     string   `json:"offset,"`
	Conns      int      `json:"conns,"`
	Consumers  int      `json:"consumers,"`
	Processors int      `json:"processors,"`
	// 10K
	MinBytes int `json:"minBytes,"`
	// 10M
	MaxBytes int `json:"maxBytes,"`
}
