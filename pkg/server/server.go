package server

import (
	"chatgpt/pkg/bot"
	"chatgpt/pkg/chatgpt"
	"chatgpt/pkg/config"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

type Server struct {
	secrets config.Config

	qqBot *bot.QQBot
	gpt   *chatgpt.ChatGPT
}

func NewServer(config string) (*Server, error) {
	secretsBin, err := os.ReadFile(filepath.Join("./configs", "secrets.yaml"))
	if err != nil {
		return nil, err
	}
	s := &Server{}
	if err = yaml.Unmarshal(secretsBin, &s.secrets); err != nil {
		return nil, err
	}
	s.gpt = chatgpt.NewChatGPT(&s.secrets.ChatGpt)

	log.Info("qq qqBot init, id: ", s.secrets.QQBot.Id)
	s.qqBot, err = bot.NewQQBot(config, s.secrets.QQBot.Id, s.gpt)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Server) Start() error {
	if err := s.qqBot.Start(); err != nil {
		panic(err)
	}
	return nil
}
