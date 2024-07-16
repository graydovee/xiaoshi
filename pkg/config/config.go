package config

type Config struct {
	ChatGpt ChatGPT `json:"chatGpt" yaml:"chatGpt"`
	QQBot   QQBot   `json:"qqBot" yaml:"qqBot"`
}

type ChatGPT struct {
	ApiKey   string      `json:"apiKey" yaml:"apiKey"`
	ImageDir string      `json:"imageDir" yaml:"imageDir"`
	ApiUrl   string      `json:"apiUrl" yaml:"apiUrl"`
	Model    string      `json:"model" yaml:"model"`
	Session  ChatSession `json:"session" yaml:"session"`
}

type ChatSession struct {
	MessageLimit  int `json:"messageLimit" yaml:"messageLimit"`
	ExpireSeconds int `json:"expireSeconds" yaml:"expireSeconds"`
}

type QQBot struct {
	Id     int64      `json:"id" yaml:"id"`
	Ws     *WebSocket `json:"ws,omitempty" yaml:"ws,omitempty"`
	Zero   ZeroConfig `json:"zero" yaml:"zero"`
	WebGui WebGui     `json:"webGui" yaml:"webGui"`
}

type ZeroConfig struct {
	SuperUsers    []int64  `json:"superUsers" yaml:"superUsers"`
	NickNames     []string `json:"nickNames" yaml:"nickNames"`
	CommandPrefix string   `json:"commandPrefix" yaml:"commandPrefix"`
}

type WebSocket struct {
	Addr  string `json:"addr" yaml:"addr"`
	Port  int    `json:"port" yaml:"port"`
	Token string `json:"token" yaml:"token"`
}

type WebGui struct {
	Host string `json:"host" yaml:"host"`
	Port int    `json:"port" yaml:"port"`
}
