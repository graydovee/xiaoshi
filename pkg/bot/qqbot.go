package bot

import (
	"context"
	"fmt"
	"github.com/dezhishen/onebot-sdk/pkg/api"
	"github.com/dezhishen/onebot-sdk/pkg/config"
	"github.com/dezhishen/onebot-sdk/pkg/event"
	"github.com/dezhishen/onebot-sdk/pkg/model"
	"strings"
)

//var _ Bot = &QQBot{}

type QQBot struct {
	ctx    context.Context
	cancel context.CancelFunc

	cfgFile string
	cfg     *config.OnebotConfig

	ApiClient   *api.OnebotApiClient
	EventClient *event.OnebotEventClient

	recvChan chan model.EventMessageGroup
}

func NewQQBot(cfgFile string) (*QQBot, error) {
	b := &QQBot{
		cfgFile:  cfgFile,
		recvChan: make(chan model.EventMessageGroup),
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
