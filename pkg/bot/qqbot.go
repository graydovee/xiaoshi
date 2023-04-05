package bot

import (
	"bytes"
	"chatgpt/pkg/chatgpt"
	"chatgpt/pkg/util"
	"context"
	"fmt"
	"github.com/dezhishen/onebot-sdk/pkg/api"
	"github.com/dezhishen/onebot-sdk/pkg/config"
	"github.com/dezhishen/onebot-sdk/pkg/event"
	"github.com/dezhishen/onebot-sdk/pkg/model"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strings"
	"time"
)

//var _ Bot = &QQBot{}

type QQBot struct {
	ctx    context.Context
	cancel context.CancelFunc

	cfgFile string
	cfg     *config.OnebotConfig

	ApiClient   *api.OnebotApiClient
	EventClient *event.OnebotEventClient

	group   util.Map[int64, *chatgpt.ChatSession]
	private util.Map[int64, *chatgpt.ChatSession]
	id      string

	gpt *chatgpt.ChatGPT
}

func NewQQBot(cfgFile string, id string, gpt *chatgpt.ChatGPT) (*QQBot, error) {
	b := &QQBot{
		cfgFile: cfgFile,
		id:      id,
		gpt:     gpt,
	}
	conf, err := config.LoadConfig(b.cfgFile)
	if err != nil {
		return nil, err
	}
	b.cfg = conf
	b.ApiClient, err = api.NewOnebotApiClient(b.cfg.Api)
	if err != nil {
		return nil, err
	}
	b.EventClient, err = event.NewOnebotEventCli(b.cfg.Event)
	if err != nil {
		return nil, err
	}

	b.EventClient.ListenMessageGroup(b.onGroupMsg)
	b.EventClient.ListenMessagePrivate(b.onPrivateMsg)
	return b, nil
}

func (s *QQBot) Start() error {
	if s.ctx != nil {
		return nil
	}
	s.ctx, s.cancel = context.WithCancel(context.Background())
	return s.startListen()
}

func (s *QQBot) startListen() error {
	s.EventClient.ListenRequestFriend(func(data model.EventRequestFriend) error {
		// friend add
		accept := strings.Contains(strings.ToLower(data.Comment), "chatgpt")
		return s.ApiClient.SetFriendAddRequest(data.Flag, accept, "")
	})

	if err := s.EventClient.StartListenWithCtx(s.ctx); err != nil {
		return err
	}
	return nil
}

func (s *QQBot) SendGroupMsg(groupId int64, message string) error {
	msg := &model.GroupMsg{
		GroupId: groupId,
		Message: []*model.MessageSegment{
			{
				Type: "text",
				Data: &model.MessageElementText{
					Text: message,
				},
			},
		},
	}
	result, err := s.ApiClient.SendGroupMsg(msg)
	if err != nil {
		return err
	}
	if result.Retcode != 200 && result.Retcode != 0 {
		return fmt.Errorf("send group message error, code: %d, msg: %s", result.Retcode, result.Msg)
	}
	return nil
}

func (s *QQBot) SendPrivateMsg(userId int64, message string) error {
	msg := &model.PrivateMsg{
		UserId: userId,
		Message: []*model.MessageSegment{
			{
				Type: "text",
				Data: &model.MessageElementText{
					Text: message,
				},
			},
		},
	}
	result, err := s.ApiClient.SendPrivateMsg(msg)
	if err != nil {
		return err
	}
	if result.Retcode != 200 && result.Retcode != 0 {
		return fmt.Errorf("send group message error, code: %d, msg: %s", result.Retcode, result.Msg)
	}
	return nil
}

func (s *QQBot) onGroupMsg(msg model.EventMessageGroup) error {
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
	if at.Qq != s.id {
		return nil
	}

	chat, _ := s.group.LoadOrStore(msg.GroupId, func() *chatgpt.ChatSession {
		c := chatgpt.NewChat(s.gpt, chatgpt.NewMemoryLimitHistory(16, time.Minute*10))
		c.SetPrompt(chatgpt.DefaultPrompt)
		return c
	})

	var response string
	var err error
	if cmd := strings.TrimSpace(text.Text); strings.HasPrefix(cmd, "/") {
		response = s.onCommand(cmd, chat)
	} else {
		response, err = chat.GetResponse(text.Text)
		if err != nil {
			return err
		}
	}

	log.Infof("send group response: %s for %d", response, msg.GroupId)
	return s.SendGroupMsg(msg.GroupId, response)
}

func (s *QQBot) onPrivateMsg(msg model.EventMessagePrivate) error {
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
		c := chatgpt.NewChat(s.gpt, chatgpt.NewMemoryLimitHistory(-1, time.Minute*10))
		c.SetPrompt(chatgpt.DefaultPrompt)
		return c
	})

	var response string
	var err error
	if cmd := strings.TrimSpace(text.Text); strings.HasPrefix(cmd, "/") {
		response = s.onCommand(cmd, chat)
	} else {
		response, err = chat.GetResponse(text.Text)
		if err != nil {
			return err
		}
	}

	log.Infof("send private response: %s for %d", response, msg.UserId)
	return s.SendPrivateMsg(msg.UserId, response)
}

const (
	CmdRole = "role"
)

func (s *QQBot) onCommand(cmdStr string, chat *chatgpt.ChatSession) string {
	out := bytes.NewBuffer(nil)

	root := cobra.Command{
		Use: "/",
	}
	root.SetArgs(splitArgs(cmdStr[1:]))
	root.SetOut(out)
	list := &cobra.Command{
		Use:   "role [role]",
		Short: "切换至预设角色",
		Long: func() string {
			roleList := bytes.NewBuffer(nil)
			p := util.NewPrinter(roleList)
			p.Println("切换至预设角色\n")
			p.Println("预设角色列表：")
			c := 1
			for role := range chatgpt.RoleMap {
				p.Println(c, ". ", role)
				c++
			}
			return roleList.String()
		}(),
		RunE: func(cmd *cobra.Command, args []string) error {
			p := util.NewPrinter(cmd.OutOrStdout())
			if len(args) == 0 {
				return fmt.Errorf("角色名为空")
			}
			roleDetail, ok := chatgpt.RoleMap[args[0]]
			if ok {
				chat.SetPrompt(roleDetail)
				p.Println("角色切换至：", args[0])
			} else {
				p.Printf("角色: %s 不存在\n", args[0])
			}
			return nil
		},
	}
	root.AddCommand(list)
	promt := &cobra.Command{
		Use:   "promt [detail]",
		Short: "手动设定人格",
		Long:  "手动设定人格",
		RunE: func(cmd *cobra.Command, args []string) error {
			p := util.NewPrinter(cmd.OutOrStdout())
			if len(args) == 0 {
				return fmt.Errorf("设定为空")
			}
			chat.SetPrompt(strings.Join(args, " "))
			p.Println("设定角色完成")
			return nil
		},
	}
	root.AddCommand(promt)
	if err := root.Execute(); err != nil {
		log.Error("execute sub command error: ", err)
	}
	return out.String()
}

func splitArgs(s string) []string {
	var args []string
	for _, argRaw := range strings.Split(strings.TrimSpace(s), " ") {
		arg := strings.TrimSpace(argRaw)
		if arg != "" {
			args = append(args, arg)
		}
	}
	return args
}
