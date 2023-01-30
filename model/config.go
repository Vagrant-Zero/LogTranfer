package model

type Config struct {
	KafkaConfig `ini:"kafka"`
	EsConfig    `ini:"es"`
}
