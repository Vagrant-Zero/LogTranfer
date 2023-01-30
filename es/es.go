package es

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
	"gocode/LogTranfer/common"
	"gocode/LogTranfer/model"
)

type ESClient struct {
	client      *elastic.Client
	logDataChan chan *sarama.ConsumerMessage
}

var (
	esClient ESClient
)

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

// 发送到ES的消息
type message struct {
	Topic string `json:"topic"`
	Data  string `json:"data"`
}

func sendMessage(config model.EsConfig) {
	for {
		select {
		case msg := <-esClient.logDataChan:
			logrus.Infof("Get msg: {%s} from Kafka", string(msg.Value))
			m := message{
				Topic: config.Index,
				Data:  string(msg.Value),
			}
			resp, err := esClient.client.Index().Index(config.Index).BodyJson(m).Do(context.Background())
			if err != nil {
				logrus.Errorf("Insert msg failed, err = %v", err.Error())
				return
			}
			logrus.Infof("Insert a msg{id : %s, type : %s} into %s success", resp.Id, resp.Type, resp.Index)
		}
	}
}

func PutLogDataIntoChan(msg *sarama.ConsumerMessage) {
	esClient.logDataChan <- msg
}
