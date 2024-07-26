package chatgpt

import (
	"sync"
	"time"
)

type History interface {
	AddHistory(msg Message)
	GetHistory() []Message
	SetLimit(limit int)
	Clear()
}

type historyMessage struct {
	msg     Message
	created time.Time
}

type MemoryLimitHistory struct {
	mu sync.Mutex

	history []historyMessage

	limit   int
	timeout time.Duration
}

func (m *MemoryLimitHistory) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.history = nil
}

func NewMemoryLimitHistory(limit int, timeout time.Duration) *MemoryLimitHistory {
	return &MemoryLimitHistory{limit: limit, timeout: timeout}
}

func (m *MemoryLimitHistory) AddHistory(msg Message) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.cleanExpired()

	m.history = append(m.history, historyMessage{
		msg:     msg,
		created: time.Now(),
	})
}

func (m *MemoryLimitHistory) GetHistory() []Message {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.cleanExpired()

	var validHistory []Message
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
