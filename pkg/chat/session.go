package chat

import (
	"context"
	"github.com/graydovee/xiaoshi/pkg/mcp"
)

type Session struct {
	chatBot *mcp.ChatBot
	history mcp.History
}

func NewChat(chatBot *mcp.ChatBot, history mcp.History) *Session {
	return &Session{
		chatBot: chatBot,
		history: history,
	}
}

func (c *Session) SetPrompt(prompts ...string) {
	c.history.GetHistory()
}

func (c *Session) ChatBot() *mcp.ChatBot {
	return c.chatBot
}

func (c *Session) GetResponse(question string) (string, error) {
	response, err := c.chatBot.Completion(context.Background(), question, c.history)
	if err != nil {
		return "", err
	}

	reply := mcp.ReadCompletionResponse(response)
	return reply, nil
}

func (c *Session) History() mcp.History {
	return c.history
}
