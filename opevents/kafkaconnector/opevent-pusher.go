package kafkaconnector

import (
	"fmt"
	"time"

	"github.com/shanluzhineng/fwpkg/opevents/pkg"
	"github.com/shanluzhineng/fwpkg/system/log"

	"github.com/shanluzhineng/configurationx"
	"github.com/shanluzhineng/fwpkg/queuex/kafkaqueue"
)

const (
	Topic_OpEventLog string = "opeventlogs"
)

type opEventLogKafkaPusher struct {
	kafkaPusher *kafkaqueue.Pusher
}

func newOpEventLogKafkaPusher() *opEventLogKafkaPusher {
	kafkaOptions := configurationx.GetInstance().Kafka.GetDefaultOptions()
	producerOptions := kafkaOptions.GetProducer(Topic_OpEventLog)
	if producerOptions == nil {
		err := fmt.Errorf("producers中没有配置topic为 %s 的消息生产者", Topic_OpEventLog)
		log.Logger.Error(err.Error())
		panic(err)
	}
	//确保topic已经被创建完成
	ensureTopicCreated(kafkaOptions.Brokers[0], Topic_OpEventLog)
	pusher := kafkaqueue.NewPusher(kafkaOptions.Brokers,
		Topic_OpEventLog,
		kafkaqueue.WithFlushInterval(time.Millisecond*time.Duration(producerOptions.FlushInterval)))

	return &opEventLogKafkaPusher{
		kafkaPusher: pusher,
	}
}

// push一组opevent到kafka中
func (p *opEventLogKafkaPusher) PushOpEvents(item *pkg.OpEventLog) error {
	if p.kafkaPusher == nil {
		return fmt.Errorf("必须先初始化Pusher")
	}
	datas := string(item.Bytes())
	err := p.kafkaPusher.Push(datas)
	if err != nil {
		log.Logger.Error(fmt.Sprintf("在将数据push到kafka中时出现异常,数据:%s,详细异常信息:%s", datas, err.Error()))
	}
	return err
}

func ensureTopicCreated(brokers string, topic string) {
	client, err := kafkaqueue.NewKafkaClient(brokers)
	if err != nil {
		panic(err)
	}
	defer client.Close()
	err = client.EnsureTopicCreated(topic)
	if err != nil {
		panic(err)
	}
}
