package chatgpt

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"strings"
)

type Session interface {
	GetResponse(question string) (string, error)
	GetImageResponse(prompt string) ([]string, error)
	History() History
}

type ChatSession struct {
	chatBot ChatBot

	prompt  []Message
	history History
}

func (c *ChatSession) GetImageResponse(prompt string) ([]string, error) {
	image, err := c.chatBot.GenImage(prompt)
	if err != nil {
		return nil, err
	}

	var res []string
	for _, d := range image.Data {
		if d.File != "" {
			res = append(res, d.File)
		}
	}
	if len(res) == 0 {
		return nil, fmt.Errorf("no image url")
	}
	return res, nil
}

func NewChat(chatBot ChatBot, history History) *ChatSession {
	return &ChatSession{
		chatBot: chatBot,
		history: history,
	}
}

func (c *ChatSession) SetPrompt(prompts ...string) {
	var promptMessage []Message
	for _, prompt := range prompts {
		promptMessage = append(promptMessage, Message{
			Role:    RoleSystem,
			Content: prompt,
		})
	}
	c.prompt = promptMessage
}

func (c *ChatSession) ChatBot() ChatBot {
	return c.chatBot
}

func (c *ChatSession) GetResponse(question string) (string, error) {
	history := c.history.GetHistory()
	msg := Message{
		Role:    RoleUser,
		Content: question,
	}
	c.history.AddHistory(msg)

	response, err := c.chatBot.GetResponse(append(c.prompt, append(history, msg)...))
	if err != nil {
		return "", err
	}

	if len(response.Choices) == 0 {
		log.Info(response)
		return "", fmt.Errorf("no choices")
	}

	var replys []string
	for _, choice := range response.Choices {
		replys = append(replys, choice.Message.Content)
	}
	reply := strings.Join(replys, "\n")
	c.history.AddHistory(Message{
		Role:    RoleAssistant,
		Content: reply,
	})

	return reply, nil
}

func (c *ChatSession) History() History {
	return c.history
}
