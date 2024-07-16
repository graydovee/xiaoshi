package server

import (
	"chatgpt/pkg/bot"
	"chatgpt/pkg/chatgpt"
	"chatgpt/pkg/config"
	log "github.com/sirupsen/logrus"
	easy "github.com/t-tomalak/logrus-easy-formatter"
	"gopkg.in/yaml.v2"
	"os"
)

type Server struct {
	config config.Config

	qqBot *bot.QQBot
	gpt   *chatgpt.ChatGPT
}

func NewServer(config string) (*Server, error) {
	cfgbin, err := os.ReadFile(config)
	if err != nil {
		return nil, err
	}
	s := &Server{}
	if err = yaml.Unmarshal(cfgbin, &s.config); err != nil {
		return nil, err
	}
	s.gpt = chatgpt.NewChatGPT(&s.config.ChatGpt)

	log.Info("qq qqBot init, id: ", s.config.QQBot.Id)
	s.qqBot, err = bot.NewQQBot(&s.config, s.gpt)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func init() {
	log.SetFormatter(&easy.Formatter{
		TimestampFormat: "2006-01-02 15:04:05",
		LogFormat:       "[zero][%time%][%lvl%]: %msg% \n",
	})
	log.SetLevel(log.DebugLevel)
}

func (s *Server) Start() error {
	if err := s.qqBot.Start(); err != nil {
		panic(err)
	}
	return nil
}
