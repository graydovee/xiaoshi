package bot

import (
	"chatgpt/pkg/chatgpt"
	"chatgpt/pkg/config"
	"chatgpt/pkg/util"
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/driver"
	"github.com/wdvxdr1123/ZeroBot/extension/shell"
	"github.com/wdvxdr1123/ZeroBot/message"
	"strconv"
	"strings"
	"time"
)

type QQBot struct {
	ctx    context.Context
	cancel context.CancelFunc

	group   util.Map[int64, *chatgpt.ChatSession]
	private util.Map[int64, *chatgpt.ChatSession]
	id      string

	cfg *config.Config
	gpt *chatgpt.ChatGPT
}

func NewQQBot(cfg *config.Config, gpt *chatgpt.ChatGPT) (*QQBot, error) {
	b := &QQBot{
		cfg: cfg,
		gpt: gpt,
		id:  strconv.FormatInt(cfg.QQBot.Id, 10),
	}
	return b, nil
}

func (s *QQBot) Start() error {
	if s.ctx != nil {
		return nil
	}
	var drivers []zero.Driver
	if s.cfg.QQBot.Ws != nil {
		wsUrl := fmt.Sprintf("ws://%s:%d", strings.TrimPrefix(s.cfg.QQBot.Ws.Addr, "ws://"), s.cfg.QQBot.Ws.Port)
		drivers = append(drivers, driver.NewWebSocketClient(wsUrl, s.cfg.QQBot.Ws.Token))
	}

	s.registerHandler()
	zero.RunAndBlock(&zero.Config{
		NickName:      s.cfg.QQBot.Zero.NickNames,
		CommandPrefix: "/",
		SuperUsers:    s.cfg.QQBot.Zero.SuperUsers,
		Driver:        drivers,
	}, nil)
	return nil
}

func (s *QQBot) registerHandler() {
	zero.OnCommand("config").Handle(s.onCommand).SetBlock(true)
	zero.OnMessage(zero.OnlyToMe).Handle(s.onMessage).SecondPriority()
}

func (s *QQBot) getSession(ctx *zero.Ctx) *chatgpt.ChatSession {
	if ctx.Event.GroupID != 0 {
		session, _ := s.group.LoadOrStore(ctx.Event.GroupID, func() *chatgpt.ChatSession {
			c := chatgpt.NewChat(s.gpt, chatgpt.NewMemoryLimitHistory(s.cfg.ChatGpt.Session.ExpireSeconds, time.Second*time.Duration(s.cfg.ChatGpt.Session.ExpireSeconds)))
			c.SetPrompt(chatgpt.DefaultPrompt)
			return c
		})
		return session
	} else if ctx.Event.UserID != 0 {
		session, _ := s.private.LoadOrStore(ctx.Event.UserID, func() *chatgpt.ChatSession {
			c := chatgpt.NewChat(s.gpt, chatgpt.NewMemoryLimitHistory(-1, time.Second*time.Duration(s.cfg.ChatGpt.Session.ExpireSeconds)))
			c.SetPrompt(chatgpt.DefaultPrompt)
			return c
		})
		return session
	}
	return nil
}

func (s *QQBot) onMessage(ctx *zero.Ctx) {
	text := ctx.Event.Message.ExtractPlainText()
	if text == "" {
		return
	}
	chatSession := s.getSession(ctx)
	if chatSession == nil {
		log.Error("chat session is nil")
		return
	}

	// chat response
	response, err := chatSession.GetResponse(text)
	if err != nil {
		log.Error("gen response error: ", err)
		return
	}
	ctx.Send(response)
	return
}

func (s *QQBot) onCommand(ctx *zero.Ctx) {
	chatSession := s.getSession(ctx)
	if chatSession == nil {
		log.Error("chat session is nil")
		return
	}

	arguments := shell.Parse(ctx.State["args"].(string))
	out, err := RunCmd(BuildCommand(chatSession), arguments)
	if err != nil {
		log.Error("execute sub command error: ", err)
		return
	}
	ctx.Send(message.Text(out))
}
