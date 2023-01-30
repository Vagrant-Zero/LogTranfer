package main

import (
	"context"
	"github.com/Shopify/sarama"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
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
	sendMessage(config)
	return nil
}

func sendMessage(config model.EsConfig) {
	for {
		select {
		case msg := <-esClient.logDataChan:
			logrus.Infof("Get msg: {%s} from Kafka", string(msg.Value))
			resp, err := esClient.client.Index().Index(config.Index).BodyJson(msg.Value).Do(context.Background())
			if err != nil {
				logrus.Errorf("Insert a msg{id : %s, type : %s} into %s failed, err = %v", resp.Id, resp.Type, resp.Index, err.Error())
				return
			}
			logrus.Infof("Insert a msg{id : %s, type : %s} into %s success", resp.Id, resp.Type, resp.Index)
		}
	}
}
