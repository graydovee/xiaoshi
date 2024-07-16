package chatgpt

import (
	"fmt"
	"git.graydove.cn/graydove/xiaoshi.git/pkg/config"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
	"testing"
)

func newChat() *ChatGPT {
	c := config.Config{}
	file, err := os.ReadFile("C:\\Projects\\chatgpt\\chatgpt\\configs\\secrets.yaml")
	if err != nil {
		panic(err)
	}
	if err = yaml.Unmarshal(file, &c); err != nil {
		panic(err)
	}
	return NewChatGPT(&c.ChatGpt)
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

func TestChatGPT_GenImage(t *testing.T) {
	log.SetLevel(log.DebugLevel)
	chat := newChat()
	image, err := chat.GenImage("contains cat and dog")
	if err != nil {
		panic(err)
	}
	fmt.Println(image)
}

func TestChatGPT_download(t *testing.T) {
	chat := newChat()
	url := "https://oaidalleapiprodscus.blob.core.windows.net/private/org-3Jq2EQOaizQRK3XpFpFqI6U8/user-eWqfNbhfUimQuyi4l7aO9Smb/img-n9MjUCgoq5rIHnNQUAk147lp.png?st=2023-04-05T06%3A52%3A52Z&se=2023-04-05T08%3A52%3A52Z&sp=r&sv=2021-08-06&sr=b&rscd=inline&rsct=image/png&skoid=6aaadede-4fb3-4698-a8f6-684d7786b067&sktid=a48cca56-e6da-484e-a814-9c849652bcb3&skt=2023-04-05T06%3A42%3A32Z&ske=2023-04-06T06%3A42%3A32Z&sks=b&skv=2021-08-06&sig=cgGzTfKsRJZm6I8mhlVgqWSUjv7JSb80/kFKiHW683c%3D"
	if err := chat.download(url, "./test.png"); err != nil {
		panic(err)
	}
	if err := os.Remove("./test.png"); err != nil {
		panic(err)
	}
}
