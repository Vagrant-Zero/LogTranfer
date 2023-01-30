package model

type KafkaConfig struct {
	Address string `ini:"address"`
	Topic   string `ini:"topic"`
}
