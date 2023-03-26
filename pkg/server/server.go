package server

import (
	"chatgpt/pkg/bot"
	"chatgpt/pkg/chatgpt"
	"chatgpt/pkg/config"
	"github.com/dezhishen/onebot-sdk/pkg/model"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"time"
)

type Server struct {
	secrets config.Config

	bot *bot.QQBot
	gpt *chatgpt.ChatGPT

	group   Map[int64, *chatgpt.ChatSession]
	private Map[int64, *chatgpt.ChatSession]
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
	s.gpt = chatgpt.NewChatGPT(s.secrets.ChatGpt.ApiKey)

	log.Info("qq bot init, id: ", s.secrets.QQBot.Id)
	s.bot, err = bot.NewQQBot(config)
	if err != nil {
		return nil, err
	}
	s.initListener()
	return s, nil
}

func (s *Server) initListener() {
	s.bot.EventClient.ListenMessageGroup(s.onGroupMsg)
	s.bot.EventClient.ListenMessagePrivate(s.onPrivateMsg)
}

func (s *Server) onGroupMsg(msg model.EventMessageGroup) error {
	log.Info("receive group msg: ", msg)
	if len(msg.Message) != 2 {
		return nil
	}
	msgs := msg.Message
	if msgs[0].Type != model.CQTypeAt && msgs[1].Type != model.CQTypeText {
		return nil
	}
	at, ok := msgs[0].Data.(*model.MessageElementAt)
	if !ok {
		log.Infof("at type error: %T", msgs[0].Data)
		return nil
	}
	text, ok := msgs[1].Data.(*model.MessageElementText)
	if !ok {
		log.Infof("message type error: %T", msgs[0].Data)
		return nil
	}
	if at.Qq != s.secrets.QQBot.Id {
		return nil
	}

	chat, _ := s.group.LoadOrStore(msg.GroupId, func() *chatgpt.ChatSession {
		return chatgpt.NewChat(s.gpt, chatgpt.NewMemoryLimitHistory(10, time.Minute*10))
	})

	response, err := chat.GetResponse(text.Text)
	if err != nil {
		return err
	}
	log.Infof("send group response: %s for %d", response, msg.GroupId)
	return s.bot.SendGroupMsg(msg.GroupId, response)
}

func (s *Server) onPrivateMsg(msg model.EventMessagePrivate) error {
	log.Info("receive private msg: ", msg)
	if len(msg.Message) < 1 {
		return nil
	}
	m := msg.Message[0]
	if m.Type != model.CQTypeText {
		return nil
	}
	text, ok := m.Data.(*model.MessageElementText)
	if !ok {
		log.Infof("message type error: %T", m.Data)
		return nil
	}

	chat, _ := s.private.LoadOrStore(msg.UserId, func() *chatgpt.ChatSession {
		return chatgpt.NewChat(s.gpt, chatgpt.NewMemoryLimitHistory(-1, time.Minute*10))
	})

	response, err := chat.GetResponse(text.Text)
	if err != nil {
		return err
	}
	log.Infof("send private response: %s for %d", response, msg.UserId)
	return s.bot.SendPrivateMsg(msg.UserId, response)
}

func (s *Server) Start() error {
	if err := s.bot.Start(); err != nil {
		panic(err)
	}
	return nil
}
