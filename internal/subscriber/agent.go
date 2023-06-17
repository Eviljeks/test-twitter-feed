package subscriber

import "sync"

type Agent struct {
	mu     sync.Mutex
	subs   map[chan string]struct{}
	quit   chan struct{}
	closed bool
}

func NewAgent() *Agent {
	return &Agent{
		subs: make(map[chan string]struct{}, 0),
		quit: make(chan struct{}),
	}
}

func (b *Agent) Publish(msg string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return
	}

	for ch := range b.subs {
		ch <- msg
	}
}

func (b *Agent) Subscribe() chan string {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return nil
	}

	ch := make(chan string)
	b.subs[ch] = struct{}{}
	return ch
}

func (b *Agent) Unsubscribe(ch chan string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return
	}

	close(ch)
	delete(b.subs, ch)
}

func (b *Agent) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()

	if b.closed {
		return
	}

	b.closed = true
	close(b.quit)

	for ch := range b.subs {
		close(ch)
	}
}
