package service

import (
	"errors"
	"sync"
)

var (
	ErrDocNotFound = errors.New("document not found")
)

type Markdown struct {
	mu   sync.RWMutex
	Docs map[string]*Document
}

type Document struct {
	// 文档的唯一 hash，用来分享
	Hash string

	mu      sync.RWMutex
	Content string
	Clients map[chan string]struct{}
}

func NewMarkdown() *Markdown {
	return &Markdown{
		Docs: make(map[string]*Document),
	}
}

func (m *Markdown) NewDocument() (*Document, error) {
	var hash string
	for {
		hash = GetHash(10)
		m.mu.RLock()
		_, exists := m.Docs[hash]
		m.mu.RUnlock()
		if !exists {
			break
		}
	}

	doc := &Document{
		Hash:    hash,
		Content: "",
		Clients: make(map[chan string]struct{}),
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.Docs[hash] = doc
	return doc, nil
}

func (m *Markdown) GetDocument(hash string) (*Document, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	doc, ok := m.Docs[hash]
	return doc, ok
}

func (m *Markdown) UpdateDocument(hash, content string) (*Document, error) {
	doc, ok := m.GetDocument(hash)
	if !ok {
		return nil, ErrDocNotFound
	}
	doc.SetContent(content)
	return doc, nil
}

func (d *Document) GetContent() string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.Content
}

func (d *Document) SetContent(content string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	d.Content = content

	for ch := range d.Clients {
		select {
		case ch <- content:
		default:
		}
	}

}

func (d *Document) AddClient() chan string {
	ch := make(chan string, 10)

	d.mu.Lock()
	d.Clients[ch] = struct{}{}
	// 注册时先推一次当前内容
	initial := d.Content
	d.mu.Unlock()

	// 将当前内容先发出去
	go func() {
		if initial != "" {
			ch <- initial
		}
	}()

	return ch
}

// RemoveClient 注销客户端，并关闭通道
func (d *Document) RemoveClient(ch chan string) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if _, ok := d.Clients[ch]; ok {
		delete(d.Clients, ch)
		close(ch)
	}
}
