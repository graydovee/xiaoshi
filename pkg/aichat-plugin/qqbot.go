package aichat_plugin

import (
	"context"
	"fmt"
	"git.graydove.cn/graydove/xiaoshi.git/pkg/chatgpt"
	"git.graydove.cn/graydove/xiaoshi.git/pkg/config"
	"git.graydove.cn/graydove/xiaoshi.git/pkg/util"
	log "github.com/sirupsen/logrus"
	zero "github.com/wdvxdr1123/ZeroBot"
	"github.com/wdvxdr1123/ZeroBot/extension/shell"
	"github.com/wdvxdr1123/ZeroBot/message"
	"strconv"
	"time"
)

type AIBot struct {
	ctx    context.Context
	cancel context.CancelFunc

	group   util.Map[int64, *chatgpt.ChatSession]
	private util.Map[int64, *chatgpt.ChatSession]
	id      string

	cfg *config.Config
	gpt *chatgpt.ChatGPT

	prompt *chatgpt.Prompt
}

func InitAIBot(cfg *config.Config) {
	gpt := chatgpt.NewChatGPT(&cfg.ChatGpt)

	log.Info("qq qqBot init, id: ", cfg.QQBot.Id)

	aiBot := &AIBot{
		cfg: cfg,
		gpt: gpt,
		id:  strconv.FormatInt(cfg.QQBot.Id, 10),
	}
	Register(aiBot)
	return
}

func (s *AIBot) getSession(ctx *zero.Ctx) *chatgpt.ChatSession {
	if ctx.Event.GroupID != 0 {
		session, _ := s.group.LoadOrStore(ctx.Event.GroupID, func() *chatgpt.ChatSession {
			c := chatgpt.NewChat(s.gpt, chatgpt.NewMemoryLimitHistory(s.cfg.ChatGpt.Session.ExpireSeconds, time.Second*time.Duration(s.cfg.ChatGpt.Session.ExpireSeconds)))
			c.SetPrompt(s.prompt.GetPrompt()...)
			return c
		})
		return session
	} else if ctx.Event.UserID != 0 {
		session, _ := s.private.LoadOrStore(ctx.Event.UserID, func() *chatgpt.ChatSession {
			c := chatgpt.NewChat(s.gpt, chatgpt.NewMemoryLimitHistory(-1, time.Second*time.Duration(s.cfg.ChatGpt.Session.ExpireSeconds)))
			c.SetPrompt(s.prompt.GetPrompt()...)
			return c
		})
		return session
	}
	return nil
}

func (s *AIBot) OnMessage(ctx *zero.Ctx) {
	text := ctx.Event.Message.ExtractPlainText()
	if text == "" {
		return
	}
	chatSession := s.getSession(ctx)
	if chatSession == nil {
		log.Error("chat session is nil")
		return
	}

	name := ctx.Event.Sender.Name()

	// chat response
	response, err := chatSession.GetResponse(fmt.Sprintf("%s: %s", name, text))
	if err != nil {
		log.Error("gen response error: ", err)
		return
	}
	ctx.Send(response)
	return
}

func (s *AIBot) OnCommand(ctx *zero.Ctx) {
	chatSession := s.getSession(ctx)
	if chatSession == nil {
		log.Error("chat session is nil")
		return
	}

	arguments := shell.Parse(ctx.State["args"].(string))
	out, err := RunCmd(s.BuildCommand(chatSession), arguments)
	if err != nil {
		log.Error("execute sub command error: ", err)
		return
	}
	ctx.Send(message.Text(out))
}
