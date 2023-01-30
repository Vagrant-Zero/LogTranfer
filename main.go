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
