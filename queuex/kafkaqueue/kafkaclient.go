package kafkaqueue

import (
	"fmt"
	//"net"
	//"strconv"

	kq "github.com/segmentio/kafka-go"
	"github.com/shanluzhineng/fwpkg/system/log"
)

type IKafkaClient interface {
	EnsureTopicCreated(topic string) error
	Close() error
}

type KafkaClient struct {
	controllerConn *kq.Conn
	Brokers        string
	Topic          string
}

func NewKafkaClient(brokers string) (IKafkaClient, error) {
	// to connect to the kafka leader via an existing non-leader connection rather than using DialLeader
	conn, err := kq.Dial("tcp", brokers)
	if err != nil {
		err = fmt.Errorf("连接到kafka时出现异常,详细异常信息:%s", err.Error())
		log.Logger.Error(err.Error())
		return nil, err
	}
	defer conn.Close()
	// 连接重复了？
	//controller, err := conn.Controller()
	//if err != nil {
	//	err = fmt.Errorf("连接到kafka时出现异常,详细异常信息:%s", err.Error())
	//	log.Logger.Error(err.Error())
	//	return nil, err
	//}
	//var controllerConn *kq.Conn
	//controllerConn, err = kq.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	//if err != nil {
	//	err = fmt.Errorf("连接到kafka时出现异常,详细异常信息:%s", err.Error())
	//	log.Logger.Error(err.Error())
	//	return nil, err
	//}
	kafkaClient := &KafkaClient{
		Brokers:        brokers,
		controllerConn: conn,
	}
	return kafkaClient, nil
}

// 确保topic创建完成
func (c *KafkaClient) EnsureTopicCreated(topic string) error {
	//By default kafka has the auto.create.topics.enable='true' (KAFKA_AUTO_CREATE_TOPICS_ENABLE='true' in the wurstmeister/kafka kafka docker image)
	topicConfigs := []kq.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     1,
			ReplicationFactor: 1,
		},
	}
	err := c.controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		err = fmt.Errorf("在创建topic时出现异常,topic:%s,详细异常信息:%s", topic, err.Error())
		log.Logger.Error(err.Error())
		return nil
	}
	return nil
}

func (c *KafkaClient) Close() error {
	if c.controllerConn == nil {
		return nil
	}
	return c.controllerConn.Close()
}
