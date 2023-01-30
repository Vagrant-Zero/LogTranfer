# LogTranfer
对应于LogAgent（收集日志），该项目用于从kafka中读取数据，发送到ES，便于后续的数据分析，如kibana。该项目参考七米老师的思路，非常感谢七米老师！
## 1. 流程分析

1. 读取配置文件
2. 初始化kafka
3. 初始化ES
4. 后台协程消费kafka消息
5. 后台协程向ES发送消息

```go
package main

import (
	"github.com/sirupsen/logrus"
	"gocode/LogTranfer/common"
	"gocode/LogTranfer/es"
	"gocode/LogTranfer/kafka"
)

func main() {
	// 1. 加载配置文件
	cfg, err, done := common.ReadConfig("./config/logtransfer.ini")
	if done {
		logrus.Errorf("Read config falied, err = %v", err.Error())
		return
	}
	// 2. 初始化kafka连接
	err = kafka.InitKafkaConsumerConnection(cfg.KafkaConfig)
	if err != nil {
		logrus.Errorf("Init kafka consumer failed, err = %v", err.Error())
		return
	}
	// 3. 初始化es连接
	err = es.InitEsConnection(cfg.EsConfig)
	if err != nil {
		logrus.Errorf("Init es connection failed, err = %v", err.Error())
		return
	}
	// 阻塞
	select {}
}
```

## 2. 解决ES和kafka之间的通信问题

+ 问题：kafka和ES的启动有先后顺序，但是两者交互的channel需要等待ES初始化完成后才能用于通信

+ 解决：在common中定义了信号收发函数，相当于起到锁的作用，在kafka消费消息的逻辑中轮询ES是否初始化完成

```go
package kafka
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
```

```go
package es
func InitEsConnection(config model.EsConfig) (err error) {
	client, err := elastic.NewClient(elastic.SetURL(config.Address))
	if err != nil {
		return err
	}
	logDataChan := make(chan *sarama.ConsumerMessage, config.LogDataChanSize)
	esClient = ESClient{
		client:      client,
		logDataChan: logDataChan,
	}
	logrus.Infof("connect to es success!")
	common.SendSignal() // 通知kafka，ES已经准备完成
	// 启用后台协程来接收kafka的数据并发送到ES
	go sendMessage(config)
	return nil
}
```
