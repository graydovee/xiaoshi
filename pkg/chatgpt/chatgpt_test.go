package chatgpt

import (
	"chatgpt/pkg/config"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
	"testing"
)

func newChat() *ChatGPT {
	c := config.Config{}
	file, err := os.ReadFile("/mnt/c/Projects/chatgpt/chatgpt/configs/secrets.yaml")
	if err != nil {
		panic(err)
	}
	if err = yaml.Unmarshal(file, &c); err != nil {
		panic(err)
	}
	return NewChatGPT(c.ChatGpt.ApiKey)
}

func TestChatGPT_GetResponse(t *testing.T) {
	log.SetLevel(log.DebugLevel)

	response, err := newChat().GetResponse([]Message{
		{
			Role:    RoleUser,
			Content: "Hello",
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(response)
	if len(response.Choices) == 0 {
		panic("no choices")
	}
}
