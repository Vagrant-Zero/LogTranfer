package kafka

import (
	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
	"gocode/LogTranfer/common"
	"gocode/LogTranfer/es"
	"gocode/LogTranfer/model"
	"time"
)

var (
	kafkaConsumer sarama.Consumer
)

func InitKafkaConsumerConnection(config model.KafkaConfig) (err error) {
	kafkaConsumer, err = sarama.NewConsumer([]string{config.Address}, nil)
	if err != nil {
		logrus.Errorf("Failed to connect Kafka consumer, err = %v", err.Error())
		return err
	}
	logrus.Infof("Init Kafka consumer connection success!")
	// 后台协程消费消息
	go consumeMessage(config.Topic)
	return nil
}

func consumeMessage(topic string) {
	logrus.Infof("Kafka consumer begin to poll message.")
	// 获取topic下的所有分区列表
	partitionList, err := kafkaConsumer.Partitions(topic)
	if err != nil {
		logrus.Errorf("Failed to get Kafka partitions, err = %v", err.Error())
		return
	}
	// 轮询ES是否初始化完成
	for {
		if common.GetSignal() {
			break
		}
		time.Sleep(time.Second)
	}
	logrus.Infof("ES is ready, Kafka consumer begin to consume message.")
	for partition := range partitionList {
		// 遍历所有分区，消费消息
		consumeMessageSinglePartition(partition, topic)
	}
}

func consumeMessageSinglePartition(partition int, topic string) {
	logrus.Infof("The partition is %d, topic is %s", partition, topic)
	// 对每个分区创建一个消费者
	partitionConsumer, err := kafkaConsumer.ConsumePartition(topic, int32(partition), sarama.OffsetOldest)
	if err != nil {
		logrus.Errorf("Failed to create kafka partition consumer, err = %v", err.Error())
		return
	}
	// 异步从每个分区消费数据
	go func(pc sarama.PartitionConsumer) {
		for msg := range pc.Messages() {
			logrus.Infof("Get a kafka msg: {Partition: %d, Offset: %d, key: %s, value: %s}",
				msg.Partition, msg.Offset, msg.Key, msg.Value)
			// 发送消息到ES
			es.PutLogDataIntoChan(msg)
		}
		defer pc.AsyncClose()
	}(partitionConsumer)
}
