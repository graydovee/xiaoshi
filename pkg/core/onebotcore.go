package core

import (
	"context"
	"github.com/CuteReimu/onebot"
	"github.com/graydovee/xiaoshi/pkg/chat"
	"github.com/graydovee/xiaoshi/pkg/config"
	"github.com/graydovee/xiaoshi/pkg/mcp"
	"github.com/graydovee/xiaoshi/pkg/util"
	"golang.org/x/time/rate"
	"log/slog"
	"strconv"
	"strings"
	"time"
)

type OneBotCore struct {
	ctx    context.Context
	cancel context.CancelFunc

	group   util.Map[int64, *chat.Session]
	private util.Map[int64, *chat.Session]
	id      string

	cfg    *config.Config
	chatAi *mcp.ChatBot
}

func NewOneBotCore(cfg *config.Config) (*OneBotCore, error) {
	slog.Info("qq qqBot init", "id", cfg.OneBot.Id)
	bot, err := mcp.NewChatBot(&cfg.MCP)
	if err != nil {
		return nil, err
	}
	core := &OneBotCore{
		cfg:    cfg,
		chatAi: bot,
		id:     strconv.FormatInt(cfg.OneBot.Id, 10),
	}
	return core, nil
}

func (s *OneBotCore) Start() error {
	if err := s.chatAi.Initialize(); err != nil {
		return err
	}

	b, err := onebot.Connect(s.cfg.OneBot.Ws.Addr, s.cfg.OneBot.Ws.Port, onebot.WsChannelAll, s.cfg.OneBot.Ws.Token, s.cfg.OneBot.Id, false)
	if err != nil {
		panic(err)
	}

	// 设置限流策略为：令牌桶容量为10，每秒放入一个令牌，超过的消息直接丢弃
	limitCfg := s.cfg.OneBot.Limit
	if limitCfg != nil && limitCfg.Enabled {
		b.SetLimiter("drop", rate.NewLimiter(rate.Limit(limitCfg.Frequency), limitCfg.Bucket))
	}
	b.ListenGroupMessage(func(message *onebot.GroupMessage) bool {
		msg := message.Message
		if !s.IsToMe(msg) {
			return true
		}

		text := ExtractPlainText(message.Message)
		session := s.getSession(message.GroupId, IdTypeGroup)
		response, err := s.generateResponse(text, session)
		if err != nil {
			slog.Error("gen response error: ", err)
			return true
		}

		var ret onebot.MessageChain
		ret = append(ret, &onebot.Text{Text: response})
		_, err = b.SendGroupMessage(message.GroupId, ret)
		if err != nil {
			slog.Error("发送失败", "error", err)
		}
		return true
	})
	b.ListenPrivateMessage(func(message *onebot.PrivateMessage) bool {
		text := ExtractPlainText(message.Message)
		session := s.getSession(message.UserId, IdTypeGroup)
		response, err := s.generateResponse(text, session)
		if err != nil {
			slog.Error("gen response error: ", err)
			return true
		}

		var ret onebot.MessageChain
		ret = append(ret, &onebot.Text{Text: response})
		_, err = b.SendPrivateMessage(message.UserId, ret)
		if err != nil {
			slog.Error("发送失败", "error", err)
		}
		return true
	})
	return nil
}

const (
	IdTypeGroup = "group"
	IdTypeUser  = "user"
)

func (s *OneBotCore) getSession(id int64, idType string) *chat.Session {
	switch idType {
	case IdTypeGroup:
		session, _ := s.group.LoadOrStore(id, func() *chat.Session {
			history := mcp.NewMemoryLimitHistory(s.cfg.Memory.MessageLimit, time.Second*time.Duration(s.cfg.Memory.ExpireSeconds))
			if s.cfg.MCP.SystemPrompt != "" {
				history.SetSystemPrompt(s.cfg.MCP.SystemPrompt)
			}
			c := chat.NewChat(s.chatAi, history)
			c.SetPrompt(s.cfg.MCP.SystemPrompt)
			return c
		})
		return session
	case IdTypeUser:
		session, _ := s.private.LoadOrStore(id, func() *chat.Session {
			history := mcp.NewMemoryLimitHistory(-1, time.Second*time.Duration(s.cfg.Memory.ExpireSeconds))
			if s.cfg.MCP.SystemPrompt != "" {
				history.SetSystemPrompt(s.cfg.MCP.SystemPrompt)
			}
			c := chat.NewChat(s.chatAi, mcp.NewMemoryLimitHistory(-1, time.Second*time.Duration(s.cfg.Memory.ExpireSeconds)))
			c.SetPrompt(s.cfg.MCP.SystemPrompt)
			return c
		})
		return session
	}
	return nil
}

func (s *OneBotCore) IsToMe(m onebot.MessageChain) bool {
	for _, val := range m {
		if val.GetMessageType() == "at" {
			at := val.(*onebot.At)
			if at.QQ == strconv.FormatInt(s.cfg.OneBot.Id, 10) {
				return true
			}
		}
	}
	return false
}

func (s *OneBotCore) generateResponse(text string, session *chat.Session) (string, error) {
	if strings.HasPrefix(text, "/") {
		text = strings.TrimPrefix(text, "/")

		arguments := ParseShell(text)
		return RunCmd(s.BuildCommand(session), arguments)
	} else {
		return session.GetResponse(text)
	}
}
