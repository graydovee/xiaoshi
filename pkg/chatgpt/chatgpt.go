package chatgpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"os"
)

const (
	gpturl       = "https://api.openai.com/v1/chat/completions"
	model3d50301 = "gpt-3.5-turbo-0301"
	model3d5     = "gpt-3.5-turbo"
)

type ChatBot interface {
	GetResponse(msg []Message) (*Response, error)
}

type ChatGPT struct {
	apiKey string
	client *http.Client
}

func NewChatGPT(apiKey string) *ChatGPT {
	// 创建 Client 对象，并使用 Transport 对象
	return &ChatGPT{
		apiKey: apiKey,
		client: &http.Client{
			Transport: GetTransport(),
		},
	}
}

func GetTransport() *http.Transport {
	proxyStr := os.Getenv("https_proxy")
	if proxyStr == "" {
		log.Debug("https_proxy is not set")
		return nil
	}

	// 解析代理地址
	proxyURL, err := url.Parse(proxyStr)
	if err != nil {
		log.Debug("Error parsing proxy URL:", err)
		return nil
	}
	log.Debug("use proxy: ", proxyStr)

	// 创建 Transport 对象
	return &http.Transport{Proxy: http.ProxyURL(proxyURL)}
}

func (c *ChatGPT) GetResponse(msg []Message) (*Response, error) {
	log.Debugf("answer for: %v", msg)
	request := &RequestBody{
		Model:       model3d5,
		Messages:    msg,
		Temperature: 0.7,
		N:           1,
		Stop:        "\n",
		MaxTokens:   200,
	}

	requestBody, err := json.Marshal(request)
	log.Debug(string(requestBody))

	req, err := http.NewRequest("POST", gpturl, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	// 设置 HTTP 请求标头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	// 发送 HTTP 请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// 处理 HTTP 响应
	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	var response Response
	if err = json.Unmarshal(responseBytes, &response); err != nil {
		return nil, err
	}
	log.Debug("response: ", string(responseBytes))

	return &response, nil
}

type RepeatedBot struct {
}

func (r RepeatedBot) GetResponse(msg []Message) (*Response, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	log.Info(string(data))
	return &Response{
		Choices: []Choice{
			{
				Message: Message{
					Role:    RoleAssistant,
					Content: msg[len(msg)-1].Content,
				},
			},
		},
	}, nil
}
