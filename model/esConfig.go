package model

type EsConfig struct {
	Address         string `ini:"address"`
	Index           string `ini:"index"`
	LogDataChanSize int    `ini:"logDataChan_size"`
}
