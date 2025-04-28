package mcp

import (
	"github.com/openai/openai-go"
	"sync"
	"time"
)

type History interface {
	GetSystemPrompt() *openai.ChatCompletionMessageParamUnion
	SetSystemPrompt(msg string)
	AddHistory(msg openai.ChatCompletionMessageParamUnion)
	GetHistory() []openai.ChatCompletionMessageParamUnion
	SetLimit(limit int)
	Clear()
}

var _ History = &MemoryLimitHistory{}

type historyMessage struct {
	msg     openai.ChatCompletionMessageParamUnion
	created time.Time
}

type MemoryLimitHistory struct {
	mu sync.Mutex

	systemMsg *openai.ChatCompletionMessageParamUnion

	history []historyMessage

	limit   int
	timeout time.Duration
}

func (m *MemoryLimitHistory) GetSystemPrompt() *openai.ChatCompletionMessageParamUnion {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.systemMsg
}

func (m *MemoryLimitHistory) SetSystemPrompt(msg string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	sysMsg := openai.AssistantMessage(msg)
	m.systemMsg = &sysMsg
}

func (m *MemoryLimitHistory) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.history = nil
}

func NewMemoryLimitHistory(limit int, timeout time.Duration) *MemoryLimitHistory {
	return &MemoryLimitHistory{limit: limit, timeout: timeout}
}

func (m *MemoryLimitHistory) AddHistory(msg openai.ChatCompletionMessageParamUnion) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.cleanExpired()

	m.history = append(m.history, historyMessage{
		msg:     msg,
		created: time.Now(),
	})
}

func (m *MemoryLimitHistory) GetHistory() []openai.ChatCompletionMessageParamUnion {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.cleanExpired()

	var validHistory []openai.ChatCompletionMessageParamUnion
	if m.systemMsg != nil {
		validHistory = append(validHistory, *m.systemMsg)
	}
	for _, hmsg := range m.history {
		validHistory = append(validHistory, hmsg.msg)
	}
	return validHistory
}

func (m *MemoryLimitHistory) SetLimit(limit int) {
	m.limit = limit
}

func (m *MemoryLimitHistory) cleanExpired() {
	if m.timeout == 0 {
		return
	}
	deadLine := time.Now().Add(-m.timeout)

	for i, hmsg := range m.history {
		if hmsg.created.After(deadLine) {
			m.history = m.history[i:]
			break
		}
	}

	if m.limit >= 0 && len(m.history) > m.limit {
		m.history = m.history[len(m.history)-m.limit:]
	}
}
