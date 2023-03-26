package chatgpt

import (
	"sync"
	"time"
)

type History interface {
	AddHistory(msg Message)
	GetHistory() []Message
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

func NewMemoryLimitHistory(limit int, timeout time.Duration) *MemoryLimitHistory {
	return &MemoryLimitHistory{limit: limit, timeout: timeout}
}

func (m *MemoryLimitHistory) AddHistory(msg Message) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.history = append(m.history, historyMessage{
		msg:     msg,
		created: time.Now(),
	})
}

func (m *MemoryLimitHistory) GetHistory() []Message {
	m.mu.Lock()
	defer m.mu.Unlock()
	var msg []Message
	deadLine := time.Now().Add(-m.timeout)
	for _, hmsg := range m.history {
		if hmsg.created.Before(deadLine) {
			continue
		}
		msg = append(msg, hmsg.msg)
	}
	if m.limit >= 0 && len(msg) > m.limit {
		msg = msg[len(msg)-m.limit:]
	}
	return msg
}
