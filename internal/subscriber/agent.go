package subscriber

import "sync"

type Agent struct {
	mu     sync.Mutex
	subs   []chan string
	quit   chan struct{}
	closed bool
}

func NewAgent() *Agent {
	return &Agent{
		subs: make([]chan string, 0),
		quit: make(chan struct{}),
	}
}

func (b *Agent) Publish(msg string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return
	}

	for _, ch := range b.subs {
		ch <- msg
	}
}

func (b *Agent) Subscribe() <-chan string {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return nil
	}

	ch := make(chan string)
	b.subs = append(b.subs, ch)
	return ch
}

func (b *Agent) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return
	}

	b.closed = true
	close(b.quit)

	for _, ch := range b.subs {
		close(ch)
	}
}
