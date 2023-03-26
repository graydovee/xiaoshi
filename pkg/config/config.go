package config

type Config struct {
	ChatGpt ChatGPT `json:"chatGpt" yaml:"chatGpt"`
	QQBot   QQBot   `json:"qqBot" yaml:"qqBot"`
}

type ChatGPT struct {
	ApiKey string `json:"apiKey" yaml:"apiKey"`
}

type QQBot struct {
	Id string `json:"id" yaml:"id"`
}
