package chatgpt

import (
	"bytes"
	"chatgpt/pkg/config"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultApiUrl = "https://api.openai.com"
	uriChat       = "/v1/chat/completions"
	uriImageGen   = "/v1/images/generations"
)

type ChatBot interface {
	SetModel(mode string)
	GetResponse(msg []Message) (*ChatResponse, error)
	GenImage(prompt string) (*GenImageResponse, error)
}

var _ ChatBot = &ChatGPT{}

type ChatGPT struct {
	apiKey   string
	imageDir string
	apiUrl   string
	model    string
	client   *http.Client
}

func (c *ChatGPT) SetModel(mode string) {
	c.model = mode
}

func (c *ChatGPT) buildUrl(uri string) string {
	return fmt.Sprintf("%s/%s", strings.TrimSuffix(c.apiUrl, "/"), strings.TrimPrefix(uri, "/"))
}

func NewChatGPT(gpt *config.ChatGPT) *ChatGPT {
	var u string
	if gpt.ApiUrl == "" {
		u = defaultApiUrl
	} else {
		u = gpt.ApiUrl
	}
	// 创建 Client 对象，并使用 Transport 对象
	return &ChatGPT{
		apiKey:   gpt.ApiKey,
		imageDir: gpt.ImageDir,
		model:    gpt.Model,
		apiUrl:   u,
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

func (c *ChatGPT) GetResponse(msg []Message) (*ChatResponse, error) {
	log.Debugf("answer for: %v", msg)
	request := &ChatRequestBody{
		Model:     c.model,
		Messages:  msg,
		N:         1,
		MaxTokens: 256,
	}

	var response ChatResponse
	if err := c.post(c.buildUrl(uriChat), request, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (c *ChatGPT) GenImage(prompt string) (*GenImageResponse, error) {
	log.Debugf("gen image for: %v", prompt)
	request := &GenImageRequestBody{
		Prompt: prompt,
		N:      1,
		Size:   "512x512",
	}

	var response GenImageResponse
	if err := c.post(c.buildUrl(uriImageGen), request, &response); err != nil {
		return nil, err
	}

	for i, data := range response.Data {
		id := uuid.New().String()
		fileName := filepath.Join(c.imageDir, id+".png")
		if err := c.download(data.Url, fileName); err != nil {
			log.Error("save image error: ", err)
			continue
		}
		response.Data[i].File = fileName
	}
	return &response, nil
}

func (c *ChatGPT) download(url, fileName string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (c *ChatGPT) post(url string, body any, response any) error {
	requestBody, err := json.Marshal(body)
	log.Debug(string(requestBody))

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return err
	}

	// 设置 HTTP 请求标头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	// 发送 HTTP 请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// 处理 HTTP 响应
	responseBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	log.Debug("responseBytes: ", string(responseBytes))

	if response != nil {
		return json.Unmarshal(responseBytes, response)
	}
	return nil
}

type RepeatedBot struct {
}

func (r RepeatedBot) SetModel(mode string) {
}

func (r RepeatedBot) GenImage(prompt string) (*GenImageResponse, error) {
	return nil, fmt.Errorf("unsupported")
}

var _ ChatBot = &RepeatedBot{}

func (r RepeatedBot) GetResponse(msg []Message) (*ChatResponse, error) {
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, err
	}
	log.Info(string(data))
	return &ChatResponse{
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
