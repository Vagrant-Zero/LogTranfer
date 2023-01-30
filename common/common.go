package common

import (
	"github.com/go-ini/ini"
	"github.com/sirupsen/logrus"
	"gocode/LogTranfer/model"
)

var (
	signalChan chan bool = make(chan bool, 1)
)

func ReadConfig(path string) (*model.Config, error, bool) {
	var cfg = new(model.Config)
	err := ini.MapTo(cfg, path)
	if err != nil {
		return nil, err, true
	}
	logrus.Infof("read config success!")
	return cfg, nil, false
}

// SendSignal 发送信号
func SendSignal() {
	signalChan <- true
}

// GetSignal 接收信号
func GetSignal() bool {
	return <-signalChan
}
