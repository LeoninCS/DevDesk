package CodeShare

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	MaxEntries     = 1000
	MaxContentSize = 100000
)

type Code struct {
	Author      string `json:"author"`
	Language    string `json:"language"`
	Content     string `json:"content"`
	Hash        string `json:"hash"`
	DestroyTime int64  `json:"destroy_time"`
}

type CodeShare struct {
	mu    sync.RWMutex
	store map[string]*Code
	order []string

	rdb *redis.Client
	ctx context.Context
}

func New(redisAddr, redisPass string, redisDB int) *CodeShare {
	var rdb *redis.Client
	if redisAddr != "" {
		rdb = redis.NewClient(&redis.Options{
			Addr:     redisAddr,
			Password: redisPass,
			DB:       redisDB,
		})
	}

	cs := &CodeShare{
		store: make(map[string]*Code),
		order: make([]string, 0, MaxEntries),
		rdb:   rdb,
		ctx:   context.Background(),
	}

	go func() {
		for range time.Tick(3 * time.Hour) {
			cs.cleanExpired()
		}
	}()

	return cs
}

// 上传
func (cs *CodeShare) Upload(author, lang, content string, ttl int64) *Code {
	if len(content) > MaxContentSize {
		return nil
	}

	destroy := time.Now().Unix() + ttl
	data := fmt.Sprintf("%s|%s|%s|%d", author, lang, content, destroy)
	hashBytes := sha256.Sum256([]byte(data))
	hash := hex.EncodeToString(hashBytes[:])

	code := &Code{
		Author:      author,
		Language:    lang,
		Content:     content,
		Hash:        hash,
		DestroyTime: destroy,
	}

	cs.mu.Lock()
	defer cs.mu.Unlock()

	// 内存存储
	if len(cs.store) >= MaxEntries {
		oldest := cs.order[0]
		delete(cs.store, oldest)
		cs.order = cs.order[1:]
	}
	cs.store[hash] = code
	cs.order = append(cs.order, hash)

	// 同步写入 Redis
	if cs.rdb != nil {
		data, _ := json.Marshal(code)
		expire := time.Duration(ttl) * time.Second
		cs.rdb.Set(cs.ctx, "code:"+hash, data, expire)
	}

	return code
}

// 获取
func (cs *CodeShare) Get(hash string) (*Code, bool) {
	cs.mu.RLock()
	code, ok := cs.store[hash]
	cs.mu.RUnlock()

	if ok && time.Now().Unix() <= code.DestroyTime {
		return code, true
	}

	// 本地没有，从 Redis 拉取
	if cs.rdb != nil {
		val, err := cs.rdb.Get(cs.ctx, "code:"+hash).Result()
		if err == nil {
			var c Code
			if json.Unmarshal([]byte(val), &c) == nil {
				if time.Now().Unix() <= c.DestroyTime {
					// 放回内存缓存
					cs.mu.Lock()
					cs.store[hash] = &c
					cs.order = append(cs.order, hash)
					cs.mu.Unlock()
					return &c, true
				}
			}
		}
	}

	return nil, false
}

// 清理过期的
func (cs *CodeShare) cleanExpired() {
	now := time.Now().Unix()
	cs.mu.Lock()
	defer cs.mu.Unlock()

	newOrder := cs.order[:0]
	for _, h := range cs.order {
		c := cs.store[h]
		if c == nil || c.DestroyTime < now {
			delete(cs.store, h)
			continue
		}
		newOrder = append(newOrder, h)
	}
	cs.order = newOrder
}
